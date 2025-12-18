package linking

import (
	core "typefox.dev/fastbelt"
)

type NameProvider interface {
	GetName(node core.AstNode) (string, *core.Token)
}

type DefaultNameProvider struct{}

func NewDefaultNameProvider() NameProvider {
	return &DefaultNameProvider{}
}

func (n *DefaultNameProvider) GetName(node core.AstNode) (string, *core.Token) {
	if namedNode, ok := node.(core.NamedNode); ok {
		return namedNode.Name(), namedNode.NameToken()
	} else {
		return "", nil
	}
}
