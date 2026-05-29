// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

import (
	"iter"
	"reflect"

	"typefox.dev/fastbelt/util/extiter"
)

// SymbolDescription describes a named AST declaration that references can resolve to.
type SymbolDescription struct {
	// URI is the document URI where the symbol is declared.
	URI URI
	// Node is the AST node that declares the symbol.
	Node AstNode
	// Name is the source unit that provides the symbol's textual name.
	Name StringUnit
}

// NewSymbolDescription returns a [SymbolDescription] for node and name.
//
// The description URI is derived from node's document.
func NewSymbolDescription(node AstNode, name StringUnit) *SymbolDescription {
	doc := node.Document()
	return &SymbolDescription{
		URI:  doc.URI,
		Node: node,
		Name: name,
	}
}

// NewTokenSymbolDescription returns a [SymbolDescription] for a [NamedTokenNode].
//
// It uses [NamedTokenNode.NameToken] as the symbol name.
func NewTokenSymbolDescription(node NamedTokenNode) *SymbolDescription {
	return NewSymbolDescription(node, node.NameToken())
}

// NewCompositeNodeSymbolDescription returns a [SymbolDescription] for a [NamedCompositeNode].
//
// It uses [NamedCompositeNode.NameNode] as the symbol name.
func NewCompositeNodeSymbolDescription(node NamedCompositeNode) *SymbolDescription {
	return NewSymbolDescription(node, node.NameNode())
}

// EmptySymbolDescriptions is an empty [SymbolSeq] sentinel.
//
// It can be reused by implementations that have no symbols to return.
var EmptySymbolDescriptions = extiter.Empty[*SymbolDescription]()

// SymbolContainers is a service that is able to generate new [SymbolContainer] items
// for the current language.
// It is used for the [Document.LocalSymbols], [Document.ExportedSymbols], and
// [Document.ImportedSymbols] fields.
type SymbolContainers interface {
	// New returns a container instance for storing symbol descriptions.
	New() SymbolContainer
}

// SymbolSeq is a sequence of symbol descriptions.
type SymbolSeq = iter.Seq[*SymbolDescription]

// A SymbolContainer is an efficient data structure for storing and querying symbol descriptions.
// References usually need to query symbols for specific AST types.
// Language specific implementations optimize for this by indexing descriptions by the type of their AST node.
type SymbolContainer interface {
	// Put attempts to put the given description into the container.
	// Returns true if the description was added.
	// The container is allowed to reject descriptions that it does not want to hold, for example
	// because they are of the wrong type.
	Put(desc *SymbolDescription) bool
	// All returns an iterator over all descriptions in the container.
	All() SymbolSeq
	// ForType returns an iterator over all descriptions in the container whose node is of the given type.
	ForType(targetType reflect.Type) SymbolSeq
}

// EmptySymbolContainer is an immutable [SymbolContainer] with no symbols.
var EmptySymbolContainer SymbolContainer = &emptySymbolContainer{}

type emptySymbolContainer struct{}

func (c *emptySymbolContainer) Put(desc *SymbolDescription) bool {
	// This container is immutable, don't accept any descriptions.
	return false
}

func (c *emptySymbolContainer) All() SymbolSeq {
	return EmptySymbolDescriptions
}

func (c *emptySymbolContainer) ForType(targetType reflect.Type) SymbolSeq {
	return EmptySymbolDescriptions
}

// MergeSymbolContainers merges multiple symbol containers into one. The resulting container
// is immutable and reflects the combined contents of all input containers.
func MergeSymbolContainers(containers iter.Seq[SymbolContainer]) SymbolContainer {
	if extiter.IsEmpty(containers) {
		return EmptySymbolContainer
	}
	return &mergedSymbolContainer{
		containers: containers,
	}
}

type mergedSymbolContainer struct {
	containers iter.Seq[SymbolContainer]
}

func (c *mergedSymbolContainer) Put(desc *SymbolDescription) bool {
	// This container is immutable, don't accept any descriptions.
	return false
}

func (c *mergedSymbolContainer) All() SymbolSeq {
	return extiter.FlatMap(c.containers, func(container SymbolContainer) SymbolSeq {
		return container.All()
	})
}

func (c *mergedSymbolContainer) ForType(targetType reflect.Type) SymbolSeq {
	return extiter.FlatMap(c.containers, func(container SymbolContainer) SymbolSeq {
		return container.ForType(targetType)
	})
}

// LocalSymbols are used for lexical scoping within a document.
type LocalSymbols interface {
	// For returns the [SymbolContainer] for the given AST node, which contains all symbols
	// that are locally visible in that node.
	For(node AstNode) SymbolContainer
}
