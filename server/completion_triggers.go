// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

// CompletionTriggers declares the characters that should automatically open
// the completion popup (e.g. "." in JavaScript, ":" in YAML). The default
// implementation returns nil - completion is then only opened by the
// client's manual shortcut (Ctrl+Space in VS Code).
//
// The DefaultLanguageServer.Initialize handler merges these into the LSP
// CompletionProvider capability so the client knows when to call.
type CompletionTriggers interface {
	TriggerCharacters() []string
}

// DefaultCompletionTriggers returns nil - no auto-open characters.
type DefaultCompletionTriggers struct{}

// NewDefaultCompletionTriggers returns the no-op trigger set.
func NewDefaultCompletionTriggers() CompletionTriggers {
	return &DefaultCompletionTriggers{}
}

// TriggerCharacters returns nil.
func (*DefaultCompletionTriggers) TriggerCharacters() []string {
	return nil
}
