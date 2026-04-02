// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package linking

import (
	core "typefox.dev/fastbelt"
)

type Namer interface {
	Name(node core.AstNode) core.StringUnit
}

type DefaultNamer struct{}

func NewDefaultNamer() Namer {
	return &DefaultNamer{}
}

func (n *DefaultNamer) Name(node core.AstNode) core.StringUnit {
	if namedNode, ok := node.(core.NamedTokenNode); ok {
		return namedNode.NameToken()
	} else if namedStringNode, ok := node.(core.NamedStringNode); ok {
		return namedStringNode.NameNode()
	} else {
		return nil
	}
}
