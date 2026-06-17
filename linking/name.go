// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package linking

import (
	core "typefox.dev/fastbelt"
)

// Denominator can be implemented by AST node Impl structs to provide custom naming logic.
type Denominator interface {
	// Denominate determines the name of the receiver node.
	Denominate() core.StringUnit
}

// Name checks the 'Name' attribute of the given node and returns its value as a StringUnit,
// i.e. a Token or a CompositeNode, both of which can produce a string.
//
// Language-specific implementations can be provided by implementing the [Denominator] interface.
func Name(node core.AstNode) core.StringUnit {
	// Use the language-specific implementation associated with the node type if available
	if custom, ok := node.(Denominator); ok {
		return custom.Denominate()
	}

	if namedNode, ok := node.(core.NamedTokenNode); ok {
		// Unwrap the pointer to prevent nil issues
		if t := namedNode.NameToken(); t != nil {
			return t
		}
	} else if namedStringNode, ok := node.(core.NamedCompositeNode); ok {
		// Unwrap the pointer to prevent nil issues
		if cn := namedStringNode.NameNode(); cn != nil {
			return cn
		}
	}
	return nil
}
