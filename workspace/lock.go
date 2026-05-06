// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"context"
	"sync"
)

// Lock controls read/write access to the workspace.
type Lock interface {
	// Write cancels any pending or in-progress write, then acquires an exclusive
	// lock and calls do with a fresh context and a downgrade function. Calling
	// downgrade atomically transitions from the exclusive write lock to a shared
	// read lock, allowing reads to proceed while the caller's read phase continues.
	// downgrade is idempotent; a safety net ensures it is always called.
	Write(ctx context.Context, do func(ctx context.Context, downgrade func()))
	// Read acquires a shared lock, calls do, then releases the lock.
	// It blocks while a write is in progress or pending.
	// Returns ctx.Err() if the context is cancelled while waiting.
	Read(ctx context.Context, do func(ctx context.Context)) error
}

// DefaultLock is the default implementation of WorkspaceLock.
//
// The key property is atomic write-to-read downgrade: when downgrade is called,
// the caller atomically transitions from holding the exclusive write lock to holding
// a shared read lock, with no window in which a new writer could sneak in. This
// guarantees that phase 3 (validation) completes under a read lock before any
// subsequent write phase can begin.
//
// Cancellation only cancels the context passed to do - every write that enters
// Write always calls do, even if its context was cancelled by a newer write.
// This ensures that document mutations inside do (e.g. applying text changes)
// are never silently dropped; the cancelled do simply skips expensive work by
// checking ctx.Err() before each phase.
type DefaultLock struct {
	mu           sync.Mutex
	cond         *sync.Cond
	writeHeld    bool               // exclusive write phase is active
	writeWaiters int                // number of goroutines waiting to acquire the write lock
	readers      int                // number of active shared read lock holders
	readyCh      chan struct{}      // closed when !writeHeld && writeWaiters==0; replaced each cycle
	cancelWrite  context.CancelFunc // cancels the current pending or in-progress write
}

// NewDefaultLock creates a new default workspace lock.
func NewDefaultLock() Lock {
	l := &DefaultLock{}
	l.cond = sync.NewCond(&l.mu)
	l.readyCh = make(chan struct{})
	close(l.readyCh) // initially readable
	return l
}

func (l *DefaultLock) Write(ctx context.Context, do func(ctx context.Context, downgrade func())) {
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
	l.mu.Unlock()

	var once sync.Once
	downgrade := func() {
		once.Do(func() {
			l.mu.Lock()
			defer l.mu.Unlock()
			l.writeHeld = false
			// Downgrade: atomically acquire a read lock before releasing the write
			// lock. This leaves no window in which a new writer could start before
			// the caller's read phase (phase 3 / validation) has completed.
			l.readers++
			if l.writeWaiters == 0 {
				// No writer is waiting - unblock pending reads.
				close(l.readyCh)
			}
			l.cond.Broadcast() // wake any writer waiting in cond.Wait
		})
	}

	defer func() {
		downgrade() // safety net: ensures the write lock is always released
		l.mu.Lock()
		// Release the read lock that downgrade acquired.
		l.readers--
		if l.readers == 0 {
			l.cond.Broadcast() // wake any writer waiting for readers to drain
		}
		l.mu.Unlock()
	}()

	do(ctx, downgrade)
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

	defer func() {
		l.mu.Lock()
		l.readers--
		if l.readers == 0 {
			l.cond.Broadcast()
		}
		l.mu.Unlock()
	}()

	do(ctx)
	return nil
}
