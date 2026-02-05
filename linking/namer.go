// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package linking

import (
	core "typefox.dev/fastbelt"
)

type Namer interface {
	Name(node core.AstNode) (string, *core.Token)
}

type DefaultNamer struct{}

func NewDefaultNamer() Namer {
	return &DefaultNamer{}
}

func (n *DefaultNamer) Name(node core.AstNode) (string, *core.Token) {
	if namedNode, ok := node.(core.NamedNode); ok {
		return namedNode.Name(), namedNode.NameToken()
	} else {
		return "", nil
	}
}
