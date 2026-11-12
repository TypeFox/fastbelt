// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package parallel

import "sync"

// Waiter is a one-shot latch: Await blocks until Close is called.
// The zero value is ready to use. A Waiter must not be copied after first use.
type Waiter struct {
	init  sync.Once
	close sync.Once
	ch    chan struct{}
}

func (w *Waiter) channel() chan struct{} {
	w.init.Do(func() { w.ch = make(chan struct{}) })
	return w.ch
}

// Close releases all current and future Await calls. It is safe to call more than once.
func (w *Waiter) Close() {
	w.close.Do(func() { close(w.channel()) })
}

// Await blocks until Close has been called.
func (w *Waiter) Await() {
	<-w.channel()
}

// Done returns a channel that is closed by Close, for use in select statements.
func (w *Waiter) Done() <-chan struct{} {
	return w.channel()
}
