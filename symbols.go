// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

import (
	"iter"
	"reflect"

	"typefox.dev/fastbelt/util/extiter"
)

type AstNodeDescription struct {
	URI         URI
	Node        AstNode
	Name        string
	NameSegment *TextSegment
	FullSegment *TextSegment
}

func NewAstNodeDescription(node AstNode, name string, nameSegment, fullSegment *TextSegment) *AstNodeDescription {
	doc := node.Document()
	return &AstNodeDescription{
		URI:         doc.URI,
		Node:        node,
		Name:        name,
		NameSegment: nameSegment,
		FullSegment: fullSegment,
	}
}

var EmptyAstNodeDescriptions = extiter.Empty[*AstNodeDescription]()

// [SymbolContainers] is a service that is able to generate new [SymbolContainer] items for the current language.
// It is used for the [Document.LocalSymbols], [Document.ExportedSymbols], and [Document.ImportedSymbols] fields.
type SymbolContainers interface {
	New() SymbolContainer
}

// Shorthand for a sequence of symbol descriptions.
type SymbolSeq = iter.Seq[*AstNodeDescription]

// A [SymbolContainer] is an efficient data structure for storing and querying symbol descriptions.
// References usually need to query symbols for specific AST types.
// Language specific implementations optimize for this by indexing descriptions by the type of their AST node.
type SymbolContainer interface {
	// Attempts to put the given description into the container.
	// Returns true if the description was added.
	// The container is allowed to reject descriptions that it does not want to hold, for example because they are of the wrong type.
	Put(desc *AstNodeDescription) bool
	// Returns an iterator over all descriptions in the container.
	All() SymbolSeq
	// Returns an iterator over all descriptions in the container that have a node of the given type.
	Type(targetType reflect.Type) SymbolSeq
}

var EmptySymbolContainer SymbolContainer = &emptySymbolContainer{}

type emptySymbolContainer struct{}

func (c *emptySymbolContainer) Put(desc *AstNodeDescription) bool {
	// This container is immutable
	panic("Cannot put description into empty symbol container")
}

func (c *emptySymbolContainer) All() SymbolSeq {
	return EmptyAstNodeDescriptions
}

func (c *emptySymbolContainer) Type(targetType reflect.Type) SymbolSeq {
	return EmptyAstNodeDescriptions
}

type mergedSymbolContainer struct {
	containers iter.Seq[SymbolContainer]
}

func (c *mergedSymbolContainer) Put(desc *AstNodeDescription) bool {
	// This container is immutable
	panic("Cannot put description into merged symbol container")
}

func (c *mergedSymbolContainer) All() SymbolSeq {
	return extiter.FlatMap(c.containers, func(container SymbolContainer) SymbolSeq {
		return container.All()
	})
}

func (c *mergedSymbolContainer) Type(targetType reflect.Type) SymbolSeq {
	return extiter.FlatMap(c.containers, func(container SymbolContainer) SymbolSeq {
		return container.Type(targetType)
	})
}

// Merges multiple symbol containers into one. The resulting container is immutable and reflects the combined contents of all input containers.
func MergeSymbolContainers(containers iter.Seq[SymbolContainer]) SymbolContainer {
	if extiter.IsEmpty(containers) {
		return EmptySymbolContainer
	}
	return &mergedSymbolContainer{
		containers: containers,
	}
}

// [LocalSymbols] are used for lexical scoping within a document.
type LocalSymbols interface {
	// Returns the [SymbolContainer] for the given AST node, which contains all symbols that are locally visible in that node.
	For(node AstNode) SymbolContainer
}
