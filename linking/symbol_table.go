package linking

import (
	core "typefox.dev/fastbelt"
)

type LocalSymbols = []*core.AstNodeDescription

func SymbolsOfType[T core.AstNode](s LocalSymbols) LocalSymbols {
	result := LocalSymbols{}
	for _, desc := range s {
		if _, ok := desc.Node.(T); ok {
			result = append(result, desc)
		}
	}
	return result
}

func LocalScopeOfType[T core.AstNode](node core.AstNode, fn func(core.AstNode) LocalSymbols) core.Scope {
	symbols := fn(node)
	filtered := SymbolsOfType[T](symbols)
	elements := map[string][]*core.AstNodeDescription{}
	for _, desc := range filtered {
		name := desc.Name
		if _, ok := elements[name]; !ok {
			elements[name] = []*core.AstNodeDescription{}
		}
		elements[name] = append(elements[name], desc)
	}
	var outer core.Scope = nil
	if container := node.Container(); container != nil {
		outer = LocalScopeOfType[T](container, fn)
	}
	return core.NewMapScope(elements, outer)
}

type SymbolTable = map[core.AstNode]LocalSymbols

type LocalSymbolTableProvider interface {
	Compute(key string, root core.AstNode)
	Reset(key string)
	LocalSymbols(node core.AstNode) LocalSymbols
}

type DefaultSymbolTable struct {
	srv       LinkingSrvCont
	uriToNode map[string][]core.AstNode
	symbols   SymbolTable
}

func NewDefaultSymbolTable(srv LinkingSrvCont) LocalSymbolTableProvider {
	return &DefaultSymbolTable{
		srv:       srv,
		uriToNode: map[string][]core.AstNode{},
		symbols:   SymbolTable{},
	}
}

func (s *DefaultSymbolTable) Compute(key string, root core.AstNode) {
	s.Reset(key)
	nodes := []core.AstNode{}
	core.Traverse(root, func(node core.AstNode) {
		nodes = append(nodes, node)
		container := node.Container()
		if container != nil {
			name, nameToken := s.srv.Linking().NameProvider.GetName(node)
			if name != "" {
				var segment *core.TextSegment
				if nameToken != nil {
					segment = &nameToken.Segment
				}
				desc := core.NewAstNodeDescription(node, name, segment, node.Segment())
				symbols := s.symbols[container]
				if symbols == nil {
					symbols = LocalSymbols{}
				}
				s.symbols[container] = append(symbols, desc)
			}
		}
	})
	s.uriToNode[key] = nodes
}

func (s *DefaultSymbolTable) Reset(key string) {
	if nodes, ok := s.uriToNode[key]; ok {
		for _, node := range nodes {
			delete(s.symbols, node)
		}
		delete(s.uriToNode, key)
	}
}

func (s *DefaultSymbolTable) LocalSymbols(node core.AstNode) LocalSymbols {
	if descs, ok := s.symbols[node]; ok {
		return descs
	} else {
		return LocalSymbols{}
	}
}
