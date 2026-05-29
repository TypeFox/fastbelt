// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

import (
	"iter"
	"slices"

	"typefox.dev/fastbelt/util/collections"
	"typefox.dev/fastbelt/util/extiter"
)

// DefaultLink resolves text in scope and returns the matching symbol description.
//
// DefaultLink is the standard implementation used by generated linker code.
// It returns the first match from [Scope.ElementByName], or a default
// [ReferenceError] if no symbol is found.
//
// DefaultLink panics if scope is nil.
// Use [EmptyScope] to represent "no visible symbols".
func DefaultLink(scope Scope, text string) (*SymbolDescription, *ReferenceError) {
	if scope == nil {
		panic("Scope cannot be nil. Return an empty scope (fastbelt.EmptyScope) instead.")
	}
	description := scope.ElementByName(text)
	if description == nil {
		return nil, defaultRefError(text)
	}
	return description, nil
}

func defaultRefError(text string) *ReferenceError {
	return NewReferenceError("Could not resolve reference to '" + text + "'.")
}

// Scope provides name-based lookup for symbol descriptions visible at a reference site.
//
// Implementations usually represent one lexical scope and optionally chain to an
// outer scope.
type Scope interface {
	// ElementByName returns one symbol for name, or nil when no symbol matches.
	//
	// When multiple symbols with the same name are visible, the returned element
	// is implementation-defined.
	ElementByName(name string) *SymbolDescription
	// ElementsByName returns all visible symbols named name.
	//
	// Implementations can return multiple symbols for overload-like scenarios or
	// duplicate declarations.
	ElementsByName(name string) iter.Seq[*SymbolDescription]
	// AllElements returns all symbols visible in this scope chain.
	AllElements() iter.Seq[*SymbolDescription]
}

type emptyScope struct{}

func (s *emptyScope) ElementByName(name string) *SymbolDescription {
	return nil
}

func (s *emptyScope) ElementsByName(name string) iter.Seq[*SymbolDescription] {
	return EmptySymbolDescriptions
}

func (s *emptyScope) AllElements() iter.Seq[*SymbolDescription] {
	return EmptySymbolDescriptions
}

// EmptyScope is a scope with no symbols.
//
// It is used as a non-nil sentinel for "no visible symbols", for example when
// a scope provider cannot compute applicable candidates.
var EmptyScope Scope = &emptyScope{}

// SeqScope is a scope backed by an iterator over symbol descriptions.
//
// SeqScope performs linear lookup over its local elements and optionally
// delegates misses to an outer scope.
type SeqScope struct {
	elements iter.Seq[*SymbolDescription]
	outer    Scope
}

// NewSeqScope creates a [SeqScope] from local elements and an optional outer scope.
//
// Local elements are searched before the outer scope.
func NewSeqScope(elements iter.Seq[*SymbolDescription], outer Scope) *SeqScope {
	return &SeqScope{
		elements: elements,
		outer:    outer,
	}
}

// ElementByName returns the first local symbol named name, then checks outer.
func (s *SeqScope) ElementByName(name string) *SymbolDescription {
	for desc := range s.elements {
		if desc.Name.String() == name {
			return desc
		}
	}
	if s.outer != nil {
		return s.outer.ElementByName(name)
	}
	return nil
}

// ElementsByName returns all local symbols named name followed by outer matches.
func (s *SeqScope) ElementsByName(name string) iter.Seq[*SymbolDescription] {
	matching := extiter.Filter(s.elements, func(desc *SymbolDescription) bool {
		return desc.Name.String() == name
	})
	if s.outer != nil {
		return extiter.Concat(matching, s.outer.ElementsByName(name))
	}
	return matching
}

// AllElements returns local elements followed by all outer elements.
func (s *SeqScope) AllElements() iter.Seq[*SymbolDescription] {
	if s.outer != nil {
		return extiter.Concat(s.elements, s.outer.AllElements())
	}
	return s.elements
}

// MapScope is a scope backed by a name-indexed multimap.
//
// MapScope is optimized for repeated name lookups by grouping local symbols by
// their string name and optionally chaining to an outer scope.
type MapScope struct {
	elements collections.MultiMap[string, *SymbolDescription]
	outer    Scope
}

// NewMapScope creates a [MapScope] from a name-indexed symbol multimap and an outer scope.
func NewMapScope(elements collections.MultiMap[string, *SymbolDescription], outer Scope) *MapScope {
	return &MapScope{
		elements: elements,
		outer:    outer,
	}
}

// NewMapScopeFromSlice builds a [MapScope] from a slice of symbols.
//
// Symbols are indexed by their [SymbolDescription.Name] string.
func NewMapScopeFromSlice(elements []*SymbolDescription, outer Scope) *MapScope {
	return NewMapScopeFromSeq(slices.Values(elements), outer)
}

// NewMapScopeFromSeq builds a [MapScope] from an iterator of symbols.
//
// Symbols are consumed once and indexed by their [SymbolDescription.Name] string.
func NewMapScopeFromSeq(elements iter.Seq[*SymbolDescription], outer Scope) *MapScope {
	elemMap := collections.NewMultiMap[string, *SymbolDescription]()
	for desc := range elements {
		elemMap.Put(desc.Name.String(), desc)
	}
	return NewMapScope(elemMap, outer)
}

// ElementByName returns one local symbol named name, then checks outer.
//
// If multiple local symbols have the same name, the first inserted symbol is returned.
func (s *MapScope) ElementByName(name string) *SymbolDescription {
	if elems, exists := s.elements.TryGet(name); exists && len(elems) > 0 {
		return elems[0]
	} else if s.outer != nil {
		return s.outer.ElementByName(name)
	}
	return nil
}

// ElementsByName returns all local symbols named name followed by outer matches.
func (s *MapScope) ElementsByName(name string) iter.Seq[*SymbolDescription] {
	elems := s.elements.Get(name)
	if len(elems) == 0 {
		if s.outer != nil {
			// Delegate directly to outer scope
			return s.outer.ElementsByName(name)
		} else {
			// No elements found and no outer scope
			return EmptySymbolDescriptions
		}
	} else {
		seq := slices.Values(elems)
		if s.outer == nil {
			// No outer scope, return only the local elements
			return seq
		}
		// Concatenate local elements with outer scope elements
		return extiter.Concat(seq, s.outer.ElementsByName(name))
	}
}

// AllElements returns all local elements and, when present, all outer elements.
//
// Local iteration order follows the underlying multimap and is not guaranteed.
func (s *MapScope) AllElements() iter.Seq[*SymbolDescription] {
	if s.elements.Size() == 0 {
		if s.outer != nil {
			// Delegate directly to outer scope
			return s.outer.AllElements()
		} else {
			return EmptySymbolDescriptions
		}
	} else {
		seq := s.elements.Values()
		if s.outer == nil {
			// No outer scope, return only the local elements
			return seq
		}
		// Concatenate local elements with outer scope elements
		return extiter.Concat(seq, s.outer.AllElements())
	}
}
