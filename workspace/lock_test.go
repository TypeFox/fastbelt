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
)

const shortWait = 20 * time.Millisecond
const longWait = 2 * time.Second

// TestReadRunsDoAndReturnsNil verifies the basic happy path: Read calls do and returns nil.
func TestReadRunsDoAndReturnsNil(t *testing.T) {
	lock := NewDefaultLock()
	called := false
	err := lock.Read(context.Background(), func(ctx context.Context) { called = true })
	assert.NoError(t, err)
	assert.True(t, called)
}

// TestConcurrentReads verifies that multiple Read calls can hold the lock simultaneously.
func TestConcurrentReads(t *testing.T) {
	lock := NewDefaultLock()

	const n = 10
	inside := make(chan struct{}, n)
	release := make(chan struct{})

	var wg sync.WaitGroup
	for range n {
		wg.Go(func() {
			err := lock.Read(context.Background(), func(ctx context.Context) {
				inside <- struct{}{}
				// Wait for the test to signal release before exiting
				// so all readers are inside simultaneously.
				<-release
			})
			assert.NoError(t, err)
		})
	}

	// All n readers must be inside simultaneously.
	for range n {
		select {
		case <-inside:
		case <-time.After(longWait):
			t.Fatal("timed out waiting for concurrent readers")
		}
	}
	close(release)
	wg.Wait()
}

// TestReadContextCancelledBeforeWrite verifies that Read returns ctx.Err() when
// the context is already cancelled before the read lock is acquired.
func TestReadContextCancelledBeforeAcquire(t *testing.T) {
	lock := NewDefaultLock()

	// Hold write phase so the reader must wait.
	inWrite := make(chan struct{})
	releaseWrite := make(chan struct{})
	go func() {
		lock.Write(context.Background(), func(ctx context.Context, downgrade func()) {
			close(inWrite)
			<-releaseWrite
		})
	}()
	<-inWrite

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // already cancelled

	called := false
	err := lock.Read(ctx, func(ctx context.Context) { called = true })

	assert.ErrorIs(t, err, context.Canceled)
	assert.False(t, called)

	close(releaseWrite)
}

// TestReadContextCancelledWhileWaiting verifies that a blocked Read returns
// ctx.Err() when its context is cancelled while waiting for a write to finish.
func TestReadContextCancelledWhileWaiting(t *testing.T) {
	lock := NewDefaultLock()

	inWrite := make(chan struct{})
	releaseWrite := make(chan struct{})
	go func() {
		lock.Write(context.Background(), func(ctx context.Context, downgrade func()) {
			close(inWrite)
			<-releaseWrite
		})
	}()
	<-inWrite

	ctx, cancel := context.WithCancel(context.Background())
	readErr := make(chan error, 1)
	called := false
	go func() {
		readErr <- lock.Read(ctx, func(ctx context.Context) { called = true })
	}()

	time.Sleep(shortWait)
	cancel()

	select {
	case err := <-readErr:
		assert.ErrorIs(t, err, context.Canceled)
		assert.False(t, called)
	case <-time.After(longWait):
		t.Fatal("Read did not return after context cancellation")
	}

	close(releaseWrite)
}

// TestWriteBlocksReadsUntilDowngrade verifies that reads are blocked during
// the write phase and unblocked once downgrade is called.
func TestWriteBlocksReadsUntilDowngrade(t *testing.T) {
	lock := NewDefaultLock()

	inWritePhase := make(chan struct{})
	doDowngrade := make(chan struct{})
	writeDoRunning := make(chan struct{})

	go func() {
		lock.Write(context.Background(), func(ctx context.Context, downgrade func()) {
			close(inWritePhase)
			<-doDowngrade
			downgrade()
			// Hold the read lock so we can observe reads running concurrently.
			<-writeDoRunning
		})
	}()
	<-inWritePhase

	// A read started during the write phase must block.
	readReached := make(chan struct{})
	go func() {
		err := lock.Read(context.Background(), func(ctx context.Context) { close(readReached) })
		assert.NoError(t, err)
	}()

	select {
	case <-readReached:
		t.Fatal("read proceeded before downgrade")
	case <-time.After(shortWait):
		// expected: blocked
	}

	// Signal downgrade - the read should now unblock.
	close(doDowngrade)

	select {
	case <-readReached:
		// good: read proceeded after downgrade
	case <-time.After(longWait):
		t.Fatal("read did not proceed after downgrade")
	}

	close(writeDoRunning)
}

// TestDowngradeAllowsReadsWhileDoStillRuns verifies that do continues running
// (e.g. validation / phase 3) after downgrade while readers proceed concurrently.
func TestDowngradeAllowsReadsWhileDoStillRuns(t *testing.T) {
	lock := NewDefaultLock()

	readProceedDone := make(chan struct{})
	doFinish := make(chan struct{})

	go func() {
		lock.Write(context.Background(), func(ctx context.Context, downgrade func()) {
			downgrade()
			// do is still running (phase 3) - wait for the concurrent reader.
			<-readProceedDone
			close(doFinish)
		})
	}()

	// The read should proceed immediately (downgrade was called synchronously above).
	err := lock.Read(context.Background(), func(ctx context.Context) {
		close(readProceedDone)
	})
	assert.NoError(t, err)

	select {
	case <-doFinish:
	case <-time.After(longWait):
		t.Fatal("write do did not finish after concurrent read completed")
	}
}

// TestDowngradeAtomicity is the critical correctness test: it verifies that no new
// write can acquire the lock between the call to downgrade and the end of do.
// Without atomic downgrade this test would be racy.
func TestDowngradeAtomicity(t *testing.T) {
	lock := NewDefaultLock()

	inReadPhase := make(chan struct{})
	releaseReadPhase := make(chan struct{})

	// First write: downgrades and then lingers in the read phase.
	go func() {
		lock.Write(context.Background(), func(ctx context.Context, downgrade func()) {
			downgrade() // transition to read phase
			close(inReadPhase)
			<-releaseReadPhase // hold the read lock
		})
	}()
	<-inReadPhase

	// Second write: must not start its do until the first write's do has returned.
	var mu sync.Mutex
	var events []string
	record := func(s string) {
		mu.Lock()
		events = append(events, s)
		mu.Unlock()
	}

	write2Started := make(chan struct{})
	go func() {
		lock.Write(context.Background(), func(ctx context.Context, downgrade func()) {
			record("write2")
			close(write2Started)
		})
	}()

	// Give the second write time to queue up while the read phase is still held.
	time.Sleep(shortWait)
	record("release")
	close(releaseReadPhase)

	select {
	case <-write2Started:
	case <-time.After(longWait):
		t.Fatal("second write never started")
	}

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, []string{"release", "write2"}, events,
		"write2 must not start before the read phase is released")
}

// TestWriteWaitsForActiveReaders verifies that Write only acquires the lock
// after all in-progress Read calls have completed.
func TestWriteWaitsForActiveReaders(t *testing.T) {
	lock := NewDefaultLock()

	readerInside := make(chan struct{})
	releaseReader := make(chan struct{})
	go func() {
		err := lock.Read(context.Background(), func(ctx context.Context) {
			close(readerInside)
			<-releaseReader
		})
		assert.NoError(t, err)
	}()
	<-readerInside

	var mu sync.Mutex
	var events []string
	record := func(s string) {
		mu.Lock()
		events = append(events, s)
		mu.Unlock()
	}

	writeStarted := make(chan struct{})
	go func() {
		lock.Write(context.Background(), func(ctx context.Context, downgrade func()) {
			record("write")
			close(writeStarted)
		})
	}()

	time.Sleep(shortWait)
	record("release")
	close(releaseReader)

	select {
	case <-writeStarted:
	case <-time.After(longWait):
		t.Fatal("write never started after reader finished")
	}

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, []string{"release", "write"}, events,
		"write must not start before the reader releases")
}

// TestNewWriteCancelsPendingWrite verifies that when a second Write call arrives
// while the first is still waiting to acquire the lock, the first still runs do
// (so document mutations are never lost) but receives a cancelled context so it
// can skip expensive work such as building.
func TestNewWriteCancelsPendingWrite(t *testing.T) {
	lock := NewDefaultLock()

	// Hold read lock so both writers must wait.
	readerAcquired := make(chan struct{})
	releaseReader := make(chan struct{})
	go func() {
		err := lock.Read(context.Background(), func(ctx context.Context) {
			close(readerAcquired)
			<-releaseReader
		})
		assert.NoError(t, err)
	}()
	<-readerAcquired

	g1Done := make(chan struct{})
	var g1Ctx context.Context
	go func() {
		lock.Write(context.Background(), func(ctx context.Context, downgrade func()) {
			g1Ctx = ctx
		})
		close(g1Done)
	}()
	time.Sleep(shortWait) // let G1 queue up in cond.Wait

	g2Done := make(chan struct{})
	var g2Ctx context.Context
	go func() {
		lock.Write(context.Background(), func(ctx context.Context, downgrade func()) {
			g2Ctx = ctx
		})
		close(g2Done)
	}()

	close(releaseReader)

	select {
	case <-g1Done:
	case <-time.After(longWait):
		t.Fatal("G1 did not complete")
	}
	select {
	case <-g2Done:
	case <-time.After(longWait):
		t.Fatal("G2 did not complete")
	}

	// G1 must have run (mutations must not be lost) but with a cancelled context.
	assert.ErrorIs(t, g1Ctx.Err(), context.Canceled, "G1 should have received a cancelled context")
	// G2 is the newest writer and must receive a live context.
	assert.NoError(t, g2Ctx.Err(), "G2 should have received a fresh context")
}

// TestNewWriteCancelsActiveWrite verifies that a second Write call cancels the
// first while its do callback is actively running, and then runs itself.
func TestNewWriteCancelsActiveWrite(t *testing.T) {
	lock := NewDefaultLock()

	firstRunning := make(chan struct{})

	go func() {
		lock.Write(context.Background(), func(ctx context.Context, downgrade func()) {
			close(firstRunning)
			<-ctx.Done() // block until cancelled
		})
	}()
	<-firstRunning

	// Second Write cancels the first and runs after it exits.
	secondCalled := false
	lock.Write(context.Background(), func(ctx context.Context, downgrade func()) {
		secondCalled = true
	})

	assert.True(t, secondCalled)
}

// TestDowngradeIsIdempotent verifies that calling downgrade multiple times
// does not panic, deadlock, or corrupt the lock state.
func TestDowngradeIsIdempotent(t *testing.T) {
	lock := NewDefaultLock()

	lock.Write(context.Background(), func(ctx context.Context, downgrade func()) {
		downgrade()
		downgrade()
		downgrade()
	})

	// Lock should be fully released; a subsequent read must succeed.
	done := make(chan struct{})
	go func() {
		err := lock.Read(context.Background(), func(ctx context.Context) { close(done) })
		assert.NoError(t, err)
	}()
	select {
	case <-done:
	case <-time.After(longWait):
		t.Fatal("lock was not released after idempotent downgrade calls")
	}
}

// TestWriteHasPriorityOverQueuedRead verifies the ordering:
// Write 1 enters -> Read queues up -> Write 2 arrives -> Write 2 runs before the Read.
func TestWriteHasPriorityOverQueuedRead(t *testing.T) {
	lock := NewDefaultLock()

	write1InDo := make(chan struct{})
	doDowngrade := make(chan struct{})
	write1Done := make(chan struct{})

	// Write 1: signal when in do, wait for downgrade cue, then linger in validation.
	go func() {
		lock.Write(context.Background(), func(ctx context.Context, downgrade func()) {
			close(write1InDo)
			<-doDowngrade
			downgrade()
			<-write1Done
		})
	}()
	<-write1InDo // Write 1 is now holding the write lock.

	// Read: queues up while Write 1 holds the lock.
	readDone := make(chan struct{})
	go func() {
		err := lock.Read(context.Background(), func(ctx context.Context) { close(readDone) })
		assert.NoError(t, err)
	}()
	time.Sleep(shortWait) // let the Read block on readyCh.

	// Write 2: arrives after the Read is already queued.
	write2Done := make(chan struct{})
	var mu sync.Mutex
	var events []string
	record := func(s string) {
		mu.Lock()
		events = append(events, s)
		mu.Unlock()
	}
	go func() {
		lock.Write(context.Background(), func(ctx context.Context, downgrade func()) {
			record("write2")
		})
		close(write2Done)
	}()
	time.Sleep(shortWait) // let Write 2 register as a waiter.

	// Release Write 1 through downgrade and then finish validation.
	close(doDowngrade)
	time.Sleep(shortWait) // let Write 1 downgrade and Write 2 attempt to acquire.
	close(write1Done)

	select {
	case <-write2Done:
	case <-time.After(longWait):
		t.Fatal("Write 2 did not complete")
	}
	select {
	case <-readDone:
	case <-time.After(longWait):
		t.Fatal("Read did not complete")
	}

	record("read")

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, []string{"write2", "read"}, events,
		"Write 2 must run before the queued Read")
}

// TestSafetyNetEnsuresDowngradeIsCalled verifies that the lock is fully released
// even when do never calls downgrade (e.g. it returns early on error).
func TestSafetyNetEnsuresDowngradeIsCalled(t *testing.T) {
	lock := NewDefaultLock()

	lock.Write(context.Background(), func(ctx context.Context, downgrade func()) {
		// Intentionally never call downgrade.
	})

	// The lock must be fully released; a subsequent read must succeed.
	done := make(chan struct{})
	go func() {
		err := lock.Read(context.Background(), func(ctx context.Context) { close(done) })
		assert.NoError(t, err)
	}()
	select {
	case <-done:
	case <-time.After(longWait):
		t.Fatal("lock was not released after do returned without calling downgrade")
	}
}
