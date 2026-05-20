// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	core "typefox.dev/fastbelt"
)

// NameProvider is a utility service for retrieving the name of an AstNode
// and the TextSegment containing the name for precise selection ranges.
type NameProvider interface {
	// GetName returns the name of the given AST node.
	// Returns empty string if the node has no name.
	GetName(node core.AstNode) string

	// GetNameSegment returns the text segment for the name portion of the node.
	// This is used for precise selection ranges in document symbols.
	// Returns nil if the name segment cannot be determined.
	GetNameSegment(node core.AstNode) *core.TextSegment
}

// DefaultNameProvider provides standard naming behavior using fastbelt's
// NamedNode and NamedTokenNode interfaces.
type DefaultNameProvider struct{}

func (p *DefaultNameProvider) GetName(node core.AstNode) string {
	if named, ok := node.(core.NamedNode); ok {
		return named.Name()
	}
	return ""
}

func (p *DefaultNameProvider) GetNameSegment(node core.AstNode) *core.TextSegment {
	if namedToken, ok := node.(core.NamedTokenNode); ok {
		token := namedToken.NameToken()
		if token != nil {
			return token.Segment()
		}
	}
	return nil
}
