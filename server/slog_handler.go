// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"
	"log/slog"

	"typefox.dev/fastbelt/util/service"
	"typefox.dev/lsp"
)

// SlogHandler is a [slog.Handler] that sends log messages to the
// LSP client via the connection. Note that this handler does not
// take attributes or groups into account when writing the message
// to the language client.
type SlogHandler struct {
	sc *service.Container
}

// NewSlogHandler creates a new slog handler that writes to the LSP connection.
func NewSlogHandler(sc *service.Container) slog.Handler {
	return &SlogHandler{sc: sc}
}

// Initializes the default slog handler to send log messages to the LSP client.
func (h *SlogHandler) OnServerInitialize(_ *lsp.ParamInitialize) {
	slog.SetDefault(slog.New(h))
}

func (h *SlogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	// Enable all log levels. You can customize this to filter based on level or context.
	return true
}

func (h *SlogHandler) Handle(ctx context.Context, record slog.Record) error {
	var lspLevel lsp.MessageType
	switch {
	case record.Level < slog.LevelInfo:
		lspLevel = lsp.Log
	case record.Level < slog.LevelWarn:
		lspLevel = lsp.Info
	case record.Level < slog.LevelError:
		lspLevel = lsp.Warning
	default:
		lspLevel = lsp.Error
	}
	return h.write(ctx, lspLevel, record.Message)
}

// No-op implementation. [SlogHandler] does not support attributes.
func (h *SlogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

// No-op implementation. [SlogHandler] does not support groups.
func (h *SlogHandler) WithGroup(name string) slog.Handler {
	return h
}

func (h *SlogHandler) write(ctx context.Context, level lsp.MessageType, msg string) error {
	conn, err := service.Get[*Connection](h.sc)
	if err != nil {
		return err
	}
	if conn.Value == nil {
		return nil // Connection not established; skip logging
	}
	lspConn := lsp.ClientDispatcher(conn.Value)
	return lspConn.LogMessage(ctx, &lsp.LogMessageParams{
		Type:    level,
		Message: msg,
	})
}
