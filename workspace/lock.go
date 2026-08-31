// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"context"
	"sync"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/service"
)

// Lock controls read/write access to workspace [core.Document] data.
//
// Builds run under the exclusive [Lock.Write]. LSP request handlers use
// [Lock.Read] for whole-workspace access, or [Lock.ReadAt] to run against a
// known set of documents as soon as those documents reach the required build
// states — potentially while a build is still writing other documents.
type Lock interface {
	// Write cancels any pending or in-progress write, then acquires an
	// exclusive lock and calls do with a fresh context. do is always called,
	// even if the context was cancelled by a newer write, so that document
	// mutations are never silently dropped. Cancelled write actions should
	// skip expensive work by checking ctx.Err().
	Write(ctx context.Context, do func(ctx context.Context))
	// Read acquires a shared lock on the whole workspace, calls do, then
	// releases the lock. It blocks while a write is in progress or pending.
	// Read returns ctx.Err() if ctx is cancelled while waiting.
	Read(ctx context.Context, do func(ctx context.Context)) error
	// ReadAt acquires a shared lock scoped to the given documents. It is
	// admitted as soon as every URI resolves to a document whose state covers
	// states, even while a build is writing. A URI with no document yet counts
	// as not ready and waits for the build that creates it. ReadAt returns
	// ctx.Err() if ctx is cancelled while waiting.
	//
	// If uris is empty, ReadAt is scoped to the whole workspace: it is
	// admitted once every document in the workspace has reached states, as
	// reported to [Lock.StateChanged] by the build.
	ReadAt(ctx context.Context, states core.DocumentState, uris []core.URI, do func(ctx context.Context)) error
	// StateChanged signals that document build states advanced, waking
	// pending ReadAt calls so they re-check their documents. current is the
	// build state that every document in the workspace is guaranteed to have
	// reached (the workspace floor), or 0 if no workspace-wide guarantee can
	// be made; the floor accumulates monotonically until the next write
	// acquires the lock and admits workspace-wide ReadAt calls.
	//
	// [Builder.Build] calls it once when the build phase begins, after every
	// per-document state advance, and with the reached floor after each phase
	// barrier. Other Lock users normally never call it.
	StateChanged(current core.DocumentState)
}

// DefaultLock is the default implementation of [Lock].
//
// Writes have priority: a pending write blocks new readers (including ReadAt),
// and starting a write cancels any write still in progress so the freshest
// edit wins. The next write acquires the lock only after all readers — shared
// and state-scoped — have drained, so readers never observe document resets.
type DefaultLock struct {
	sc   *service.Container
	docs DocumentManager // lazily resolved from sc on first ReadAt

	mu             sync.Mutex
	cond           *sync.Cond
	writeHeld      bool               // exclusive write phase is active
	writeWaiters   int                // number of goroutines waiting to acquire the write lock
	readers        int                // number of active shared lock holders (Read and ReadAt)
	stateAdvanced  bool               // current write entered the build phase; document states only advance
	workspaceState core.DocumentState // floor reached by every document; reset when a write acquires the lock
	readyCh        chan struct{}      // closed when !writeHeld && writeWaiters==0; replaced each cycle
	stateCh        chan struct{}      // closed and replaced whenever ReadAt admission may have changed
	cancelWrite    context.CancelFunc // cancels the current pending or in-progress write
}

// NewDefaultLock returns a [Lock] with write-priority scheduling and
// per-document state admission for ReadAt.
func NewDefaultLock(sc *service.Container) Lock {
	l := &DefaultLock{sc: sc}
	l.cond = sync.NewCond(&l.mu)
	l.readyCh = make(chan struct{})
	close(l.readyCh) // initially readable
	l.stateCh = make(chan struct{})
	return l
}

func (l *DefaultLock) Write(ctx context.Context, do func(ctx context.Context)) {
	ctx, cancel := context.WithCancel(ctx)

	l.mu.Lock()
	// Cancel any previous write that is still pending or in progress.
	if l.cancelWrite != nil {
		l.cancelWrite()
	}
	l.cancelWrite = cancel
	if !l.writeHeld && l.writeWaiters == 0 {
		// Transition from readable -> blocked: replace readyCh so incoming reads wait.
		l.readyCh = make(chan struct{})
	}
	l.writeWaiters++
	// Wait for the current write holder and all active readers to finish.
	// We always proceed even if our context was cancelled by a newer Write call:
	// do must run so that document mutations are never silently dropped.
	for l.readers > 0 || l.writeHeld {
		l.cond.Wait()
	}
	l.writeWaiters--
	l.writeHeld = true
	// Back in the mutation phase: ReadAt must not trust document states, and
	// the workspace floor no longer holds (documents may be reset or created).
	l.stateAdvanced = false
	l.workspaceState = 0
	l.mu.Unlock()

	defer func() {
		l.mu.Lock()
		l.writeHeld = false
		if l.writeWaiters == 0 {
			// No writer is waiting - unblock pending reads.
			close(l.readyCh)
		}
		l.wakeReadAtLocked()
		l.cond.Broadcast() // wake any writer waiting in cond.Wait
		l.mu.Unlock()
	}()

	do(ctx)
}

func (l *DefaultLock) Read(ctx context.Context, do func(ctx context.Context)) error {
	for {
		l.mu.Lock()
		if !l.writeHeld && l.writeWaiters == 0 {
			l.readers++
			l.mu.Unlock()
			break
		}
		ch := l.readyCh
		l.mu.Unlock()

		select {
		case <-ch: // state changed; re-check under the lock
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	defer l.releaseReader()

	do(ctx)
	return nil
}

func (l *DefaultLock) ReadAt(ctx context.Context, states core.DocumentState, uris []core.URI, do func(ctx context.Context)) error {
	for {
		l.mu.Lock()
		// Currently, no write is pending
		// A pending write has priority over readers
		if l.writeWaiters == 0 &&
			// The current write is ready to admit readers
			(!l.writeHeld || l.stateAdvanced) &&
			// Every requested document is at the requested states
			l.docsReadyLocked(states, uris) {
			l.readers++
			l.mu.Unlock()
			break
		}
		ch := l.stateCh
		l.mu.Unlock()

		select {
		case <-ch: // states advanced or a write cycle ended; re-check
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	defer l.releaseReader()

	do(ctx)
	return nil
}

func (l *DefaultLock) StateChanged(current core.DocumentState) {
	l.mu.Lock()
	l.stateAdvanced = true
	l.workspaceState = l.workspaceState.With(current)
	l.wakeReadAtLocked()
	l.mu.Unlock()
}

// wakeReadAtLocked wakes all ReadAt calls waiting for admission so they
// re-check their documents. Callers must hold l.mu.
func (l *DefaultLock) wakeReadAtLocked() {
	close(l.stateCh)
	l.stateCh = make(chan struct{})
}

// docsReadyLocked reports whether every URI resolves to a document whose state
// covers states. An empty uris list means the whole workspace and is checked
// against the floor reported through StateChanged, without scanning documents.
// Callers must hold l.mu.
func (l *DefaultLock) docsReadyLocked(states core.DocumentState, uris []core.URI) bool {
	if len(uris) == 0 {
		return l.workspaceState.Has(states)
	}
	if l.docs == nil {
		l.docs = service.MustGet[DocumentManager](l.sc)
	}
	for _, uri := range uris {
		doc := l.docs.Get(uri)
		if doc == nil || !doc.State().Has(states) {
			return false
		}
	}
	return true
}

func (l *DefaultLock) releaseReader() {
	l.mu.Lock()
	l.readers--
	if l.readers == 0 {
		l.cond.Broadcast() // wake any writer waiting for readers to drain
	}
	l.mu.Unlock()
}
