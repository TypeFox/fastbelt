package linking

import (
	"iter"
	"slices"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/extiter"
)

type LocalSymbols = iter.Seq[*core.AstNodeDescription]

func SymbolsOfType[T core.AstNode](s LocalSymbols) LocalSymbols {
	return extiter.Filter(s, func(desc *core.AstNodeDescription) bool {
		_, ok := desc.Node.(T)
		return ok
	})
}

func LocalScopeOfType[T core.AstNode](node core.AstNode, fn func(core.AstNode) LocalSymbols) core.Scope {
	symbols := fn(node)
	filtered := SymbolsOfType[T](symbols)
	var outer core.Scope = nil
	if container := node.Container(); container != nil {
		outer = LocalScopeOfType[T](container, fn)
	}
	if extiter.IsEmpty(filtered) {
		// Shortcut to generate fewer scopes
		if outer != nil {
			return outer
		} else {
			return core.EmptyScope
		}
	}
	return core.NewMapScopeFromSeq(filtered, outer)
}

type LocalSymbolTableProvider interface {
	// TODO: Replace "key" with Document structure once we have it
	Compute(key string, root core.AstNode)
	// TODO: Might not be required once we have a Document structure
	Reset(key string)
	LocalSymbols(node core.AstNode) LocalSymbols
}

// TODO: Refactor this once we have a Document structure
type DefaultLocalSymbolTableProvider struct {
	srv       LinkingSrvCont
	uriToNode map[string][]core.AstNode
	symbols   map[core.AstNode][]*core.AstNodeDescription
}

func NewDefaultLocalSymbolTableProvider(srv LinkingSrvCont) LocalSymbolTableProvider {
	return &DefaultLocalSymbolTableProvider{
		srv:       srv,
		uriToNode: map[string][]core.AstNode{},
		symbols:   map[core.AstNode][]*core.AstNodeDescription{},
	}
}

func (s *DefaultLocalSymbolTableProvider) Compute(key string, root core.AstNode) {
	s.Reset(key)
	nodes := []core.AstNode{}
	core.TraverseContent(root, func(node core.AstNode) {
		nodes = append(nodes, node)
		// Container is never nil for all but the root node
		container := node.Container()
		name, nameToken := s.srv.Linking().Namer.Name(node)
		if name != "" {
			var segment *core.TextSegment
			if nameToken != nil {
				segment = &nameToken.Segment
			}
			desc := core.NewAstNodeDescription(node, name, segment, node.Segment())
			symbols := s.symbols[container]
			s.symbols[container] = append(symbols, desc)
		}
	})
	s.uriToNode[key] = nodes
}

func (s *DefaultLocalSymbolTableProvider) Reset(key string) {
	if nodes, ok := s.uriToNode[key]; ok {
		for _, node := range nodes {
			delete(s.symbols, node)
		}
		delete(s.uriToNode, key)
	}
}

func (s *DefaultLocalSymbolTableProvider) LocalSymbols(node core.AstNode) LocalSymbols {
	if descs, ok := s.symbols[node]; ok {
		return slices.Values(descs)
	} else {
		return extiter.Empty[*core.AstNodeDescription]()
	}
}
