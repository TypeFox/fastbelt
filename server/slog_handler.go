// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"
	"log/slog"
	"strings"

	"typefox.dev/lsp"
)

type slogHandler struct {
	srv ServerSrvCont
}

func NewSlogHandler(srv ServerSrvCont) slog.Handler {
	return &slogHandler{srv: srv}
}

func (h *slogHandler) Enabled(ctx context.Context, level slog.Level) bool {
	// Enable all log levels. You can customize this to filter based on level or context.
	return true
}

func (h *slogHandler) Handle(ctx context.Context, record slog.Record) error {
	sb := strings.Builder{}
	sb.WriteString(record.Time.Format("15:04:05.000"))
	sb.WriteString(" [")
	sb.WriteString(record.Level.String())
	sb.WriteString("] ")
	sb.WriteString(record.Message)
	return h.write(ctx, lsp.Info, sb.String())
}

func (h *slogHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// Ignore for now - we can implement attribute handling if needed. Just return the same handler.
	return h
}

func (h *slogHandler) WithGroup(name string) slog.Handler {
	// Ignore for now - we can implement group handling if needed. Just return the same handler.
	return h
}

func (h *slogHandler) write(ctx context.Context, level lsp.MessageType, msg string) error {
	conn := h.srv.Server().Connection
	if conn == nil {
		return nil // Connection not established; skip logging
	}
	lspConn := lsp.ClientDispatcher(conn)
	return lspConn.LogMessage(ctx, &lsp.LogMessageParams{
		Type:    level,
		Message: msg,
	})
}
