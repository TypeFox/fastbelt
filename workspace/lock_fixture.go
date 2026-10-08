// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

// This file contains various fixtures that help in testing the workspace lock-

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/service"
)

const shortWait = 20 * time.Millisecond
const longWait = 2 * time.Second

// signal is a one-shot notification between the test goroutine and lock
// operations. fire is idempotent.
type signal struct {
	ch   chan struct{}
	once sync.Once
}

func newSignal() *signal { return &signal{ch: make(chan struct{})} }

func (s *signal) fire() { s.once.Do(func() { close(s.ch) }) }

// await blocks until the signal fires, failing the test after longWait.
func (s *signal) await(t *testing.T, msg string) {
	t.Helper()
	select {
	case <-s.ch:
	case <-time.After(longWait):
		t.Fatal(msg)
	}
}

// assertPending verifies the signal does not fire within shortWait.
func (s *signal) assertPending(t *testing.T, msg string) {
	t.Helper()
	select {
	case <-s.ch:
		t.Fatal(msg)
	case <-time.After(shortWait):
	}
}

// lockHarness bundles a lock and its document manager with helpers to run
// lock operations in goroutines under explicit admission and release control.
type lockHarness struct {
	t    *testing.T
	lock Lock
	dm   DocumentManager
}

func newLockHarness(t *testing.T) *lockHarness {
	sc := service.NewContainer()
	dm := NewDefaultDocumentManager(sc)
	service.Put(sc, dm)
	sc.Seal()
	return &lockHarness{t: t, lock: NewDefaultLock(sc), dm: dm}
}

// doc creates a document with the given state and registers it.
func (h *lockHarness) doc(uri string, state core.DocumentState) *core.Document {
	h.t.Helper()
	doc, err := core.NewDocumentFromString(uri, "test", "")
	if err != nil {
		h.t.Fatal(err)
	}
	doc.SetState(state)
	h.dm.Set(doc)
	return doc
}

// lockOp is a lock operation (Write, Read, or ReadAt) running in its own
// goroutine. Its callback fires entered once the lock admits it, then holds
// the lock until release fires; done fires when the lock call has returned.
type lockOp struct {
	t       *testing.T
	entered *signal
	release *signal
	done    *signal
	ctx     context.Context // ctx passed to the callback; valid once entered fired
	err     error           // result of Read/ReadAt; valid once done fired
}

func (h *lockHarness) newOp() *lockOp {
	return &lockOp{t: h.t, entered: newSignal(), release: newSignal(), done: newSignal()}
}

// enter is the callback run under the lock: it records the callback context,
// reports admission, and holds the lock until the test releases it.
func (op *lockOp) enter(ctx context.Context) {
	op.ctx = ctx
	op.entered.fire()
	<-op.release.ch
}

// startWrite runs lock.Write in a goroutine and returns its handle.
func (h *lockHarness) startWrite(ctx context.Context) *lockOp {
	op := h.newOp()
	go func() {
		h.lock.Write(ctx, op.enter)
		op.done.fire()
	}()
	return op
}

// startRead runs lock.Read in a goroutine and returns its handle.
func (h *lockHarness) startRead(ctx context.Context) *lockOp {
	op := h.newOp()
	go func() {
		op.err = h.lock.Read(ctx, op.enter)
		op.done.fire()
	}()
	return op
}

// startReadAt runs lock.ReadAt in a goroutine and returns its handle.
func (h *lockHarness) startReadAt(ctx context.Context, states core.DocumentState, uris []core.URI) *lockOp {
	op := h.newOp()
	go func() {
		op.err = h.lock.ReadAt(ctx, states, uris, op.enter)
		op.done.fire()
	}()
	return op
}

// awaitEntered waits until the lock admits the operation.
func (op *lockOp) awaitEntered(msg string) {
	op.t.Helper()
	op.entered.await(op.t, msg)
}

// assertBlocked verifies the lock does not admit the operation within shortWait.
func (op *lockOp) assertBlocked(msg string) {
	op.t.Helper()
	op.entered.assertPending(op.t, msg)
}

// assertNotEntered verifies (without waiting) that the callback never ran.
func (op *lockOp) assertNotEntered(msg string) {
	op.t.Helper()
	select {
	case <-op.entered.ch:
		op.t.Fatal(msg)
	default:
	}
}

// awaitDone waits for the lock call to return; its result is left in op.err.
func (op *lockOp) awaitDone(msg string) {
	op.t.Helper()
	op.done.await(op.t, msg)
}

// finish releases the callback, waits for the lock call to return, and
// asserts that it succeeded. Operations expected to fail use awaitDone and
// inspect op.err instead.
func (op *lockOp) finish() {
	op.t.Helper()
	op.release.fire()
	op.awaitDone("lock operation did not finish after release")
	assert.NoError(op.t, op.err)
}
