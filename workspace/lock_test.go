// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	core "typefox.dev/fastbelt"
)

// TestReadRunsDoAndReturnsNil verifies the basic happy path: Read calls do and returns nil.
func TestReadRunsDoAndReturnsNil(t *testing.T) {
	h := newLockHarness(t)
	called := false
	err := h.lock.Read(context.Background(), func(ctx context.Context) { called = true })
	assert.NoError(t, err)
	assert.True(t, called)
}

// TestConcurrentReads verifies that multiple Read calls can hold the lock simultaneously.
func TestConcurrentReads(t *testing.T) {
	h := newLockHarness(t)

	const n = 10
	readers := make([]*lockOp, n)
	for i := range readers {
		readers[i] = h.startRead(context.Background())
	}
	// All n readers must be inside simultaneously: each is admitted while the
	// others still hold their read lock.
	for _, r := range readers {
		r.awaitEntered("timed out waiting for concurrent readers")
	}
	for _, r := range readers {
		r.finish()
	}
}

// TestReadContextCancelledBeforeAcquire verifies that Read returns ctx.Err() when
// the context is already cancelled before the read lock is acquired.
func TestReadContextCancelledBeforeAcquire(t *testing.T) {
	h := newLockHarness(t)

	// Hold the write lock so the reader must wait.
	w := h.startWrite(context.Background())
	w.awaitEntered("write was not admitted")

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // already cancelled

	called := false
	err := h.lock.Read(ctx, func(ctx context.Context) { called = true })

	assert.ErrorIs(t, err, context.Canceled)
	assert.False(t, called)

	w.finish()
}

// TestReadContextCancelledWhileWaiting verifies that a blocked Read returns
// ctx.Err() when its context is cancelled while waiting for a write to finish.
func TestReadContextCancelledWhileWaiting(t *testing.T) {
	h := newLockHarness(t)

	w := h.startWrite(context.Background())
	w.awaitEntered("write was not admitted")

	ctx, cancel := context.WithCancel(context.Background())
	r := h.startRead(ctx)
	r.assertBlocked("read proceeded while write was active")

	cancel()
	r.awaitDone("Read did not return after context cancellation")
	assert.ErrorIs(t, r.err, context.Canceled)
	r.assertNotEntered("read callback ran despite cancellation")

	w.finish()
}

// TestWriteBlocksReadsUntilDone verifies that reads are blocked while a write
// is active and unblocked once the write's do returns.
func TestWriteBlocksReadsUntilDone(t *testing.T) {
	h := newLockHarness(t)

	w := h.startWrite(context.Background())
	w.awaitEntered("write was not admitted")

	r := h.startRead(context.Background())
	r.assertBlocked("read proceeded while write was active")

	w.finish()
	r.awaitEntered("read did not proceed after write finished")
	r.finish()
}

// TestWriteWaitsForActiveReaders verifies that Write only acquires the lock
// after all in-progress Read calls have completed.
func TestWriteWaitsForActiveReaders(t *testing.T) {
	h := newLockHarness(t)

	r := h.startRead(context.Background())
	r.awaitEntered("read was not admitted")

	w := h.startWrite(context.Background())
	w.assertBlocked("write started before the reader released")

	r.finish()
	w.awaitEntered("write never started after reader finished")
	w.finish()
}

// TestNewWriteCancelsPendingWrite verifies that when a second Write call arrives
// while the first is still waiting to acquire the lock, the first still runs do
// (so document mutations are never lost) but receives a cancelled context so it
// can skip expensive work such as building.
func TestNewWriteCancelsPendingWrite(t *testing.T) {
	h := newLockHarness(t)

	// Hold a read lock so both writers must wait.
	r := h.startRead(context.Background())
	r.awaitEntered("read was not admitted")

	w1 := h.startWrite(context.Background())
	time.Sleep(shortWait) // let W1 register with the lock before W2 arrives
	w2 := h.startWrite(context.Background())

	// Admission order between two queued writers is unspecified: pre-release
	// both callbacks and assert only on the contexts they received.
	w1.release.fire()
	w2.release.fire()
	r.finish()
	w1.awaitDone("W1 did not complete")
	w2.awaitDone("W2 did not complete")

	// W1 must have run (mutations must not be lost) but with a cancelled context.
	assert.ErrorIs(t, w1.ctx.Err(), context.Canceled, "W1 should have received a cancelled context")
	// W2 is the newest writer and must receive a live context.
	assert.NoError(t, w2.ctx.Err(), "W2 should have received a fresh context")
}

// TestNewWriteCancelsActiveWrite verifies that a second Write call cancels the
// first while its do callback is actively running, and then runs itself.
func TestNewWriteCancelsActiveWrite(t *testing.T) {
	h := newLockHarness(t)

	w1 := h.startWrite(context.Background())
	w1.awaitEntered("W1 was not admitted")

	// Starting W2 cancels W1's context while W1 is still running.
	w2 := h.startWrite(context.Background())
	select {
	case <-w1.ctx.Done():
	case <-time.After(longWait):
		t.Fatal("W1's context was not cancelled by the new write")
	}

	w2.assertBlocked("W2 started before W1 finished")
	w1.finish()
	w2.awaitEntered("W2 was never admitted after W1 finished")
	w2.finish()
}

// TestWriteHasPriorityOverQueuedRead verifies the ordering:
// Write 1 enters -> Read queues up -> Write 2 arrives -> Write 2 runs before the Read.
func TestWriteHasPriorityOverQueuedRead(t *testing.T) {
	h := newLockHarness(t)

	w1 := h.startWrite(context.Background())
	w1.awaitEntered("W1 was not admitted")

	// Read: queues up while Write 1 holds the lock.
	r := h.startRead(context.Background())
	r.assertBlocked("read proceeded while W1 was active")

	// Write 2: arrives after the Read is already queued.
	w2 := h.startWrite(context.Background())
	time.Sleep(shortWait) // let W2 register as a waiter

	w1.finish()
	w2.awaitEntered("W2 was never admitted")
	r.assertNotEntered("the queued read overtook the pending write")
	w2.finish()
	r.awaitEntered("read did not proceed after W2 finished")
	r.finish()
}

// TestReadAtImmediateWhenDocumentReady verifies that ReadAt runs immediately
// when the workspace is idle and the document already has the requested states.
func TestReadAtImmediateWhenDocumentReady(t *testing.T) {
	h := newLockHarness(t)
	doc := h.doc("file:///a.test", core.DocStateParsed)

	called := false
	err := h.lock.ReadAt(context.Background(), core.DocStateParsed, []core.URI{doc.URI}, func(ctx context.Context) {
		called = true
	})
	assert.NoError(t, err)
	assert.True(t, called)
}

// TestReadAtAdmittedDuringBuildPhase is the core behavior: a ReadAt whose
// document has reached the requested state is admitted while a write is still
// running, whereas a plain Read stays blocked.
func TestReadAtAdmittedDuringBuildPhase(t *testing.T) {
	h := newLockHarness(t)
	doc := h.doc("file:///a.test", 0)

	w := h.startWrite(context.Background())
	w.awaitEntered("write was not admitted")
	// Build phase: parse completes.
	doc.SetState(core.DocStateParsed)
	h.lock.StateChanged(0)

	// ReadAt(Parsed) must be admitted while the write is still active.
	ra := h.startReadAt(context.Background(), core.DocStateParsed, []core.URI{doc.URI})
	ra.awaitEntered("ReadAt was not admitted during the build phase")
	ra.finish()

	// A plain Read must still be blocked.
	r := h.startRead(context.Background())
	r.assertBlocked("plain Read proceeded while write was active")

	w.finish()
	r.awaitEntered("read did not proceed after write finished")
	r.finish()
}

// TestReadAtBlockedDuringMutationPhase verifies that ReadAt does not trust
// document states while a write is in its mutation phase, even if the bits are
// set: the write may be about to reset the document.
func TestReadAtBlockedDuringMutationPhase(t *testing.T) {
	h := newLockHarness(t)
	doc := h.doc("file:///a.test", core.DocStateParsed)

	w := h.startWrite(context.Background())
	w.awaitEntered("write was not admitted")

	ra := h.startReadAt(context.Background(), core.DocStateParsed, []core.URI{doc.URI})
	ra.assertBlocked("ReadAt was admitted during the mutation phase")

	// Entering the build phase admits the ReadAt.
	h.lock.StateChanged(0)
	ra.awaitEntered("ReadAt was not admitted after the build phase began")
	ra.finish()

	w.finish()
}

// TestReadAtWaitsForRequestedState verifies that ReadAt blocks until the
// document reaches the requested states, then proceeds mid-write.
func TestReadAtWaitsForRequestedState(t *testing.T) {
	h := newLockHarness(t)
	doc := h.doc("file:///a.test", core.DocStateParsed)

	// ReadAt requires Linked, which the document does not have yet.
	ra := h.startReadAt(context.Background(), core.DocStateParsed|core.DocStateLinked, []core.URI{doc.URI})
	ra.assertBlocked("ReadAt proceeded before the document was linked")

	w := h.startWrite(context.Background())
	w.awaitEntered("write was not admitted")
	doc.SetState(doc.State().With(core.DocStateLinked))
	h.lock.StateChanged(0)

	ra.awaitEntered("ReadAt was not admitted after the document reached the state")
	ra.finish()
	w.finish()
}

// TestReadAtWaitsForUnknownDocument verifies that a URI with no document counts
// as not ready and that ReadAt is admitted once a build creates the document.
func TestReadAtWaitsForUnknownDocument(t *testing.T) {
	h := newLockHarness(t)
	uri := core.ParseURI("file:///new.test")

	ra := h.startReadAt(context.Background(), core.DocStateParsed, []core.URI{uri})
	ra.assertBlocked("ReadAt proceeded for an unknown document")

	// A build cycle creates and parses the document.
	w := h.startWrite(context.Background())
	w.awaitEntered("write was not admitted")
	h.doc("file:///new.test", core.DocStateParsed)
	h.lock.StateChanged(0)
	w.finish()

	ra.awaitEntered("ReadAt was not admitted after the document was created and parsed")
	ra.finish()
}

// TestReadAtContextCancelledWhileWaiting verifies that a blocked ReadAt returns
// ctx.Err() when cancelled.
func TestReadAtContextCancelledWhileWaiting(t *testing.T) {
	h := newLockHarness(t)
	uri := core.ParseURI("file:///missing.test")

	ctx, cancel := context.WithCancel(context.Background())
	ra := h.startReadAt(ctx, core.DocStateParsed, []core.URI{uri})
	ra.assertBlocked("ReadAt proceeded for a missing document")

	cancel()
	ra.awaitDone("ReadAt did not return after context cancellation")
	assert.ErrorIs(t, ra.err, context.Canceled)
	ra.assertNotEntered("ReadAt callback ran despite cancellation")
}

// TestReadAtBlockedByPendingWrite verifies write priority for ReadAt: once a
// new write is pending, ReadAt is not admitted even if its documents are ready.
func TestReadAtBlockedByPendingWrite(t *testing.T) {
	h := newLockHarness(t)
	doc := h.doc("file:///a.test", 0)

	w1 := h.startWrite(context.Background())
	w1.awaitEntered("W1 was not admitted")
	doc.SetState(core.DocStateParsed)
	h.lock.StateChanged(0)

	// Write 2 queues up behind Write 1.
	w2 := h.startWrite(context.Background())
	time.Sleep(shortWait) // let W2 register as a waiter

	// ReadAt must be blocked despite the document being ready.
	ra := h.startReadAt(context.Background(), core.DocStateParsed, []core.URI{doc.URI})
	ra.assertBlocked("ReadAt was admitted while a write was pending")

	w1.finish()
	w2.awaitEntered("W2 was never admitted")
	w2.finish()

	ra.awaitEntered("ReadAt was not admitted after the pending write finished")
	ra.finish()
}

// TestWriteWaitsForReadAtReaders verifies that a new write acquires the lock
// only after ReadAt readers admitted during the previous write have drained.
func TestWriteWaitsForReadAtReaders(t *testing.T) {
	h := newLockHarness(t)
	doc := h.doc("file:///a.test", 0)

	w1 := h.startWrite(context.Background())
	w1.awaitEntered("W1 was not admitted")
	doc.SetState(core.DocStateParsed)
	h.lock.StateChanged(0)

	// ReadAt is admitted mid-write and lingers.
	ra := h.startReadAt(context.Background(), core.DocStateParsed, []core.URI{doc.URI})
	ra.awaitEntered("ReadAt was not admitted during the build phase")

	w1.finish() // Write 1 finishes; the ReadAt reader is still active.

	w2 := h.startWrite(context.Background())
	w2.assertBlocked("W2 started before the ReadAt reader released")

	ra.finish()
	w2.awaitEntered("W2 was never admitted after the ReadAt reader released")
	w2.finish()
}

// TestReadAtStress hammers ReadAt against a document that is continuously
// reset and rebuilt. Run with -race: it verifies both the admission logic
// (data behind a requested state bit is always present) and the memory
// publication through the atomic document state.
func TestReadAtStress(t *testing.T) {
	h := newLockHarness(t)

	stable := h.doc("file:///stable.test", core.DocStateParsed)
	stable.ParserErrors = []*core.ParserError{}

	target := h.doc("file:///target.test", 0)

	stop := make(chan struct{})
	var wg sync.WaitGroup

	// Writer: continuously reset and rebuild the target document.
	wg.Go(func() {
		defer close(stop)
		for range 300 {
			h.lock.Write(context.Background(), func(ctx context.Context) {
				// Mutation phase: reset.
				target.SetState(0)
				target.ParserErrors = nil
				// Build phase: rebuild and publish.
				target.ParserErrors = []*core.ParserError{}
				target.SetState(core.DocStateParsed)
				h.lock.StateChanged(0)
			})
		}
	})

	// Readers: ReadAt must only ever observe fully published data.
	for _, doc := range []*core.Document{target, stable} {
		for range 2 {
			wg.Go(func() {
				for {
					select {
					case <-stop:
						return
					default:
					}
					ctx, cancel := context.WithTimeout(context.Background(), longWait)
					_ = h.lock.ReadAt(ctx, core.DocStateParsed, []core.URI{doc.URI}, func(ctx context.Context) {
						if doc.ParserErrors == nil {
							t.Error("ReadAt observed unpublished data: ParserErrors is nil despite Parsed state")
						}
					})
					cancel()
				}
			})
		}
	}

	wg.Wait()
}

// TestReadAtWorkspaceWide verifies the whole-workspace form (no URIs): it
// waits for the floor reported through StateChanged and is admitted mid-write
// once the floor covers the requested states.
func TestReadAtWorkspaceWide(t *testing.T) {
	h := newLockHarness(t)

	// The floor starts at 0, so a workspace-wide ReadAt must block.
	ra := h.startReadAt(context.Background(), core.DocStateParsed, nil)
	ra.assertBlocked("workspace-wide ReadAt proceeded before any floor was reported")

	w := h.startWrite(context.Background())
	w.awaitEntered("write was not admitted")
	h.lock.StateChanged(core.DocStateParsed | core.DocStateExportedSymbols)

	// Admitted mid-write: the floor covers Parsed while the write still runs.
	ra.awaitEntered("workspace-wide ReadAt was not admitted after the floor was reported")
	ra.finish()

	// A higher state than the floor must still block.
	linked := h.startReadAt(context.Background(), core.DocStateLinked, nil)
	linked.assertBlocked("workspace-wide ReadAt proceeded beyond the reported floor")

	h.lock.StateChanged(core.DocStateLinked)
	linked.awaitEntered("workspace-wide ReadAt was not admitted after the floor was raised")
	linked.finish()

	w.finish()
}

// TestReadAtWorkspaceWideFloorResetsOnNewWrite verifies that the floor is
// cleared when a new write acquires the lock, so workspace-wide ReadAt calls
// wait for the new build cycle to re-establish it.
func TestReadAtWorkspaceWideFloorResetsOnNewWrite(t *testing.T) {
	h := newLockHarness(t)

	// First write establishes a floor; after it ends the floor persists.
	w1 := h.startWrite(context.Background())
	w1.awaitEntered("W1 was not admitted")
	h.lock.StateChanged(core.DocStateParsed)
	w1.finish()
	err := h.lock.ReadAt(context.Background(), core.DocStateParsed, nil, func(ctx context.Context) {})
	assert.NoError(t, err, "floor must persist after the write ends")

	// Second write: the floor resets on acquisition (mutation phase).
	w2 := h.startWrite(context.Background())
	w2.awaitEntered("W2 was not admitted")

	ra := h.startReadAt(context.Background(), core.DocStateParsed, nil)
	ra.assertBlocked("workspace-wide ReadAt trusted a stale floor during a new write")

	// The write ends without re-establishing the floor: still blocked.
	w2.finish()
	ra.assertBlocked("workspace-wide ReadAt proceeded without a re-established floor")

	// A new build cycle re-establishes the floor.
	w3 := h.startWrite(context.Background())
	w3.awaitEntered("W3 was not admitted")
	h.lock.StateChanged(core.DocStateParsed)
	ra.awaitEntered("workspace-wide ReadAt was not admitted after the floor was re-established")
	ra.finish()
	w3.finish()
}
