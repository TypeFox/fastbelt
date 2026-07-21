// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/linking"
	"typefox.dev/fastbelt/util/service"
)

// FoundName is the result of a [NameFinder.Find] call. Either field may be
// nil if the search could not determine the corresponding unit. Source and
// Target are the same unit when the search started on a symbol's declaration
// name rather than on a reference to it.
type FoundName struct {
	// The unit that contains the token that was used to start the search.
	Source core.StringUnit
	// The unit that contains the name of the symbol that was found.
	Target core.StringUnit
}

// NameFinder is responsible for finding the source and target [core.StringUnit] for a given token.
// It is used by various LSP services to find the name of a referenced/given symbol.
// Adopters should customize this service if they want to change how names are found in LSP services.
// Downstream LSP services will automatically use the new implementation.
type NameFinder interface {
	// Find returns the source and target [core.StringUnit] for the given tokens.
	// This method accepts two tokens, as in some cases (i.e. if the cursor is between two tokens),
	// there may be two relevant tokens that could be used to find a name.
	// Use [core.TokenSlice.SearchOffset2] to get the tokens at a given offset.
	Find(ctx context.Context, first, second *core.Token) FoundName
}

// DefaultNameFinder is the default implementation of [NameFinder].
// If a given token belongs to a cross-reference, it resolves the reference
// and returns the name unit of the target symbol; if the token lies within
// the name of a declaration, it returns that name unit as both source and
// target.
type DefaultNameFinder struct {
	sc *service.Container
}

// NewDefaultNameFinder creates a new default name finder.
func NewDefaultNameFinder(sc *service.Container) NameFinder {
	return &DefaultNameFinder{sc: sc}
}

func (nf *DefaultNameFinder) Find(ctx context.Context, first, second *core.Token) FoundName {
	firstResult := nf.forToken(ctx, first)
	// Early guard: if the first token already resolves to a name, return it immediately
	if firstResult.Target != nil {
		return firstResult
	}
	// Try the second token now
	secondResult := nf.forToken(ctx, second)
	if secondResult.Target != nil {
		return secondResult
	}
	// Neither token could be resolved to a name.
	// Try to return the one that has a source unit.
	if firstResult.Source != nil {
		return firstResult
	} else if secondResult.Source != nil {
		return secondResult
	}
	// Found nothing, return empty result
	return FoundName{}
}

func (nf *DefaultNameFinder) forToken(ctx context.Context, token *core.Token) FoundName {
	if token == nil {
		return FoundName{}
	}
	ref := core.ReferenceOfToken(token)
	if ref != nil {
		// The token is a reference, try to resolve it and return the target name unit
		unit := ref.Unit()
		ref.Resolve(ctx) // Ensure the reference is resolved before accessing the target
		refDescription := ref.Description()
		if refDescription == nil {
			// Reference could not be resolved, but return the source unit
			return FoundName{Source: unit}
		}
		return FoundName{Source: unit, Target: refDescription.Name}
	} else {
		// Not a reference, try to find the name segment that contains the given token
		node := token.Owner()
		if node == nil {
			return FoundName{}
		}
		nameUnit := linking.Name(node)
		if nameUnit == nil {
			return FoundName{}
		}
		segment := nameUnit.Segment()
		if segment == nil || token.TextSegment.Indices.Start < segment.Indices.Start || token.TextSegment.Indices.End > segment.Indices.End {
			return FoundName{} // The token is not within the name segment, i.e. not a name
		}
		// Source and target are the same
		return FoundName{Source: nameUnit, Target: nameUnit}
	}
}
