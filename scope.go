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

func DefaultLink(scope Scope, text string) (*AstNodeDescription, *ReferenceError) {
	if scope == nil {
		return nil, defaultRefError(text)
	}
	description := scope.ElementByName(text)
	if description == nil {
		return nil, defaultRefError(text)
	} else {
		return description, nil
	}
}

func defaultRefError(text string) *ReferenceError {
	return NewReferenceError("Could not resolve reference to '" + text + "'.")
}

type Scope interface {
	ElementByName(name string) *AstNodeDescription
	ElementsByName(name string) iter.Seq[*AstNodeDescription]
	AllElements() iter.Seq[*AstNodeDescription]
}

type emptyScope struct{}

func (s *emptyScope) ElementByName(name string) *AstNodeDescription {
	return nil
}

func (s *emptyScope) ElementsByName(name string) iter.Seq[*AstNodeDescription] {
	return EmptyAstNodeDescriptions
}

func (s *emptyScope) AllElements() iter.Seq[*AstNodeDescription] {
	return EmptyAstNodeDescriptions
}

var EmptyScope Scope = &emptyScope{}

type SeqScope struct {
	elements iter.Seq[*AstNodeDescription]
	outer    Scope
}

func NewSeqScope(elements iter.Seq[*AstNodeDescription], outer Scope) *SeqScope {
	return &SeqScope{
		elements: elements,
		outer:    outer,
	}
}

func (s *SeqScope) ElementByName(name string) *AstNodeDescription {
	for desc := range s.elements {
		if desc.Name == name {
			return desc
		}
	}
	if s.outer != nil {
		return s.outer.ElementByName(name)
	}
	return nil
}

func (s *SeqScope) ElementsByName(name string) iter.Seq[*AstNodeDescription] {
	matching := extiter.Filter(s.elements, func(desc *AstNodeDescription) bool {
		return desc.Name == name
	})
	if s.outer != nil {
		return extiter.Concat(matching, s.outer.ElementsByName(name))
	} else {
		return matching
	}
}

func (s *SeqScope) AllElements() iter.Seq[*AstNodeDescription] {
	if s.outer != nil {
		return extiter.Concat(s.elements, s.outer.AllElements())
	} else {
		return s.elements
	}
}

type MapScope struct {
	elements collections.MultiMap[string, *AstNodeDescription]
	outer    Scope
}

func NewMapScope(elements collections.MultiMap[string, *AstNodeDescription], outer Scope) *MapScope {
	return &MapScope{
		elements: elements,
		outer:    outer,
	}
}

func NewMapScopeFromSlice(elements []*AstNodeDescription, outer Scope) *MapScope {
	return NewMapScopeFromSeq(slices.Values(elements), outer)
}

func NewMapScopeFromSeq(elements iter.Seq[*AstNodeDescription], outer Scope) *MapScope {
	elemMap := collections.NewMultiMap[string, *AstNodeDescription]()
	for desc := range elements {
		elemMap.Put(desc.Name, desc)
	}
	return NewMapScope(elemMap, outer)
}

func (s *MapScope) ElementByName(name string) *AstNodeDescription {
	if elems, exists := s.elements.TryGet(name); exists && len(elems) > 0 {
		return elems[0]
	} else if s.outer != nil {
		return s.outer.ElementByName(name)
	}
	return nil
}

func (s *MapScope) ElementsByName(name string) iter.Seq[*AstNodeDescription] {
	elems := s.elements.Get(name)
	if len(elems) == 0 {
		if s.outer != nil {
			// Delegate directly to outer scope
			return s.outer.ElementsByName(name)
		} else {
			// No elements found and no outer scope
			return EmptyAstNodeDescriptions
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

func (s *MapScope) AllElements() iter.Seq[*AstNodeDescription] {
	if s.elements.Size() == 0 {
		if s.outer != nil {
			// Delegate directly to outer scope
			return s.outer.AllElements()
		} else {
			return EmptyAstNodeDescriptions
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
