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

// NewDefaultSignatureHelpTriggers returns the no-op trigger set.
func NewDefaultSignatureHelpTriggers() SignatureHelpTriggers {
	return NewSignatureHelpTriggers(nil, nil)
}

// signatureHelpTriggers is a fixed set of trigger and retrigger characters.
type signatureHelpTriggers struct {
	trigger, retrigger []string
}

// NewSignatureHelpTriggers returns a SignatureHelpTriggers that reports the
// given trigger and retrigger characters, e.g.
// NewSignatureHelpTriggers([]string{"("}, []string{","}).
func NewSignatureHelpTriggers(trigger, retrigger []string) SignatureHelpTriggers {
	return signatureHelpTriggers{trigger: trigger, retrigger: retrigger}
}

func (t signatureHelpTriggers) TriggerCharacters() []string {
	return t.trigger
}

func (t signatureHelpTriggers) RetriggerCharacters() []string {
	return t.retrigger
}
