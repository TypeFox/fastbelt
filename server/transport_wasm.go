// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

//go:build js && wasm

package server

import (
	"context"
	"io"
	"syscall/js"

	"golang.org/x/exp/jsonrpc2"
	"typefox.dev/fastbelt/util/service"
)

func SetupWasmServices(sc *service.Container) {
	if !service.Has[jsonrpc2.Binder](sc) {
		service.Put[jsonrpc2.Binder](sc, NewWasmBinder(sc))
	}
	if !service.Has[jsonrpc2.Dialer](sc) {
		service.Put[jsonrpc2.Dialer](sc, NewWasmDialer())
	}
}

// WasmReceive is the name of the global JS function that the host calls to deliver incoming messages.
const WasmReceive = "wasmReceive"

// WasmSend is the name of the global JS function that the server calls to send outgoing messages.
const WasmSend = "wasmSend"

// WasmReady is the name of the global JS function that the server calls to signal it is ready to receive messages.
const WasmReady = "wasmReady"

// WasmDialer implements [jsonrpc2.Dialer] for communication with a JavaScript
// host (typically a Web Worker) instead of stdio. Messages are exchanged as
// bare JSON strings over the JS bridge:
//
//   - Outgoing messages (server -> client) are passed to the JS function
//     globalThis.wasmSend(message).
//   - Incoming messages (client -> server) are delivered by the JS host calling
//     the exported globalThis.wasmReceive(message) function.
//
// The host is expected to translate between JSON-RPC message objects (used by
// vscode-jsonrpc's BrowserMessageReader/Writer over postMessage) and the JSON
// strings exchanged here.
type WasmDialer struct{}

// NewWasmDialer creates a dialer that bridges JSON-RPC over the JavaScript host.
func NewWasmDialer() *WasmDialer {
	return &WasmDialer{}
}

func (d *WasmDialer) Dial(ctx context.Context) (io.ReadWriteCloser, error) {
	conn := &wasmReadWriteCloser{
		incoming: make(chan []byte, 1024),
		closed:   make(chan struct{}),
	}

	// Expose the receive entry point to the JS host. The host calls this with
	// each incoming JSON-RPC message (as a string) when it arrives.
	receive := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) == 0 {
			return nil
		}
		select {
		case conn.incoming <- []byte(args[0].String()):
		case <-conn.closed:
		}
		return nil
	})
	js.Global().Set(WasmReceive, receive)

	// Signal the host that the server is ready to receive messages. The host
	// flushes any messages it buffered before this point.
	if ready := js.Global().Get(WasmReady); ready.Type() == js.TypeFunction {
		ready.Invoke()
	}

	return conn, nil
}

// wasmReadWriteCloser adapts the JS bridge to an [io.ReadWriteCloser] consumed
// by jsonrpc2. It is paired with [jsonrpc2.RawFramer] so each Write carries one
// complete JSON message and incoming messages are decoded as a JSON stream.
type wasmReadWriteCloser struct {
	incoming chan []byte
	buf      []byte
	closed   chan struct{}
}

func (c *wasmReadWriteCloser) Read(p []byte) (int, error) {
	if len(c.buf) == 0 {
		select {
		case b := <-c.incoming:
			c.buf = b
		case <-c.closed:
			return 0, io.EOF
		}
	}
	n := copy(p, c.buf)
	c.buf = c.buf[n:]
	return n, nil
}

func (c *wasmReadWriteCloser) Write(p []byte) (int, error) {
	send := js.Global().Get(WasmSend)
	if send.Type() != js.TypeFunction {
		return 0, io.ErrClosedPipe
	}
	send.Invoke(string(p))
	return len(p), nil
}

func (c *wasmReadWriteCloser) Close() error {
	select {
	case <-c.closed:
	default:
		close(c.closed)
	}
	return nil
}

// WasmBinder behaves like [DefaultBinder] but uses [jsonrpc2.RawFramer] so that
// messages are exchanged as bare JSON values rather than with Content-Length
// headers, which is the natural format for the postMessage-based JS bridge.
type WasmBinder struct {
	DefaultBinder
}

// NewWasmBinder creates a binder for the WASM/browser transport.
func NewWasmBinder(sc *service.Container) *WasmBinder {
	return &WasmBinder{DefaultBinder{sc: sc}}
}

func (b *WasmBinder) Bind(ctx context.Context, conn *jsonrpc2.Connection) (jsonrpc2.ConnectionOptions, error) {
	opts, err := b.DefaultBinder.Bind(ctx, conn)
	if err != nil {
		return opts, err
	}
	opts.Framer = jsonrpc2.RawFramer()
	return opts, nil
}
