// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

// SignatureHelpTriggers declares the characters that should open or refresh
// signature help (e.g. "(" to open, "," to re-trigger while typing arguments).
type SignatureHelpTriggers interface {
	TriggerCharacters() []string
	RetriggerCharacters() []string
}

// DefaultSignatureHelpTriggers returns nil for both trigger character sets.
type DefaultSignatureHelpTriggers struct{}

// NewDefaultSignatureHelpTriggers returns the no-op trigger set.
func NewDefaultSignatureHelpTriggers() SignatureHelpTriggers {
	return &DefaultSignatureHelpTriggers{}
}

func (*DefaultSignatureHelpTriggers) TriggerCharacters() []string {
	return nil
}

func (*DefaultSignatureHelpTriggers) RetriggerCharacters() []string {
	return nil
}
