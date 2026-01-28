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
