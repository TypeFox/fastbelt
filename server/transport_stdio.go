// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"
	"io"
	"os"

	"golang.org/x/exp/jsonrpc2"
	"typefox.dev/fastbelt/util/service"
)

func SetupStdioServices(sc *service.Container) {
	if !service.Has[jsonrpc2.Binder](sc) {
		service.Put[jsonrpc2.Binder](sc, NewDefaultBinder(sc))
	}
	if !service.Has[jsonrpc2.Dialer](sc) {
		service.Put[jsonrpc2.Dialer](sc, NewStdioDialer())
	}
}

// StdioDialer implements jsonrpc2.Dialer for stdio communication
type StdioDialer struct{}

func NewStdioDialer() *StdioDialer {
	return &StdioDialer{}
}

func (d *StdioDialer) Dial(ctx context.Context) (io.ReadWriteCloser, error) {
	return &stdioReadWriteCloser{
		Reader: os.Stdin,
		Writer: os.Stdout,
	}, nil
}

// stdioReadWriteCloser combines stdin/stdout into a ReadWriteCloser
type stdioReadWriteCloser struct {
	io.Reader
	io.Writer
}

func (rw *stdioReadWriteCloser) Close() error {
	// stdin/stdout don't need explicit closing
	return nil
}
