// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package linking

import (
	"context"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/service"
)

// LocalSymbolsProvider is a service that computes local symbol information for documents.
type LocalSymbolsProvider interface {
	// LocalSymbols traverses the document's AST and computes the local symbol table.
	// The result is stored in the document's LocalSymbols field.
	LocalSymbols(ctx context.Context, document *core.Document) core.LocalSymbols
}

// DefaultLocalSymbolsProvider is the default implementation of LocalSymbolsProvider.
type DefaultLocalSymbolsProvider struct {
	sc *service.Container
}

func NewDefaultLocalSymbolsProvider(sc *service.Container) LocalSymbolsProvider {
	return &DefaultLocalSymbolsProvider{sc: sc}
}

func (s *DefaultLocalSymbolsProvider) LocalSymbols(ctx context.Context, document *core.Document) core.LocalSymbols {
	root := document.Root
	symbolContainers := service.MustGet[core.SymbolContainers](s.sc)
	localSymbols := NewDefaultLocalSymbols()
	containers := localSymbols.Symbols

	for node := range core.AllChildren(root) {
		desc, containerNode := DescribeLocal(node)
		if desc != nil && containerNode != nil {
			if symbols, exists := containers[containerNode]; exists {
				// Fast path for existing container: just put the new description into it
				symbols.Put(desc)
			} else {
				// No container yet for this node, create a new one
				newContainer := symbolContainers.New()
				if newContainer.Put(desc) {
					// Only add the new container if the description was successfully added
					containers[containerNode] = newContainer
				}
			}
		}
	}
	document.LocalSymbols = localSymbols
	return localSymbols
}

// DefaultLocalSymbols is the default implementation of [core.LocalSymbols].
type DefaultLocalSymbols struct {
	Symbols map[core.AstNode]core.SymbolContainer
}

// NewDefaultLocalSymbols creates a new [DefaultLocalSymbols] instance.
func NewDefaultLocalSymbols() *DefaultLocalSymbols {
	return &DefaultLocalSymbols{
		Symbols: make(map[core.AstNode]core.SymbolContainer),
	}
}

// NewDefaultLocalSymbolsFromMap creates a new [DefaultLocalSymbols] instance from a map of AST nodes to symbol containers.
func NewDefaultLocalSymbolsFromMap(m map[core.AstNode]core.SymbolContainer) *DefaultLocalSymbols {
	return &DefaultLocalSymbols{
		Symbols: m,
	}
}

// For returns the [SymbolContainer] for the given AST node, which contains all symbols that are locally visible in that node.
func (ls *DefaultLocalSymbols) For(node core.AstNode) core.SymbolContainer {
	if container, exists := ls.Symbols[node]; exists {
		return container
	} else {
		return core.EmptySymbolContainer
	}
}

// LocalSymbolDescriber can be implemented by AST node Impl structs to provide custom local symbol description logic.
type LocalSymbolDescriber interface {
	// DescribeLocal determines the name and other metadata of the receiver node as a locally visible symbol.
	// It also returns a container in which the symbol shall be visible. All direct or indirect children of the container
	// will include the symbol in their local scopes.
	//
	// If the node should not be visible locally, it returns nil.
	DescribeLocal() (*core.SymbolDescription, core.AstNode)
}

// DescribeLocal checks whether the given node is a symbol (i.e. has a name) and returns a description for it.
// A symbol is mapped to its direct container, making it visible in all its siblings.
//
// Language-specific implementations can be provided by implementing the [LocalSymbolDescriber] interface.
func DescribeLocal(node core.AstNode) (*core.SymbolDescription, core.AstNode) {
	// Use the language-specific implementation associated with the node type if available
	if custom, ok := node.(LocalSymbolDescriber); ok {
		return custom.DescribeLocal()
	}

	nameUnit := Name(node)
	if nameUnit != nil {
		container := node.Container()
		name := nameUnit.String()
		segment := nameUnit.Segment()
		fullSegment := node.Segment()
		desc := core.NewSymbolDescription(node, name, segment, fullSegment)
		return desc, container
	}
	return nil, nil
}
