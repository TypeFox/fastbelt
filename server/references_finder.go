// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"
	"iter"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/linking"
	"typefox.dev/fastbelt/util/extiter"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
)

type FindReferencesOptions struct {
	// Whether to include the declaration of the symbol as a dedicated [core.ReferenceDescription] in the results.
	// If true, will contain the declaration as a reference with the same source and target node, and the range of the name as the first element.
	IncludeDeclaration bool
}

type ReferencesFinder interface {
	Find(ctx context.Context, target core.AstNode, options FindReferencesOptions) iter.Seq[*core.ReferenceDescription]
}

type DefaultReferencesFinder struct {
	sc *service.Container
}

func NewDefaultReferencesFinder(sc *service.Container) ReferencesFinder {
	return &DefaultReferencesFinder{sc: sc}
}

func (rf *DefaultReferencesFinder) Find(ctx context.Context, target core.AstNode, options FindReferencesOptions) iter.Seq[*core.ReferenceDescription] {
	sequences := []iter.Seq[*core.ReferenceDescription]{}
	if options.IncludeDeclaration {
		nameUnit := linking.Name(target)
		if nameUnit != nil {
			selfDescription := core.NewReferenceDescription(target, target, nameUnit.Segment())
			sequences = append(sequences, extiter.Of(selfDescription))
		}
	}
	documentManager := service.MustGet[workspace.DocumentManager](rf.sc)
	// Iterate through all documents and collect references to the symbol
	for doc := range documentManager.All() {
		refDescriptions := doc.ReferenceDescriptions.ForTarget(target)
		sequences = append(sequences, refDescriptions)
	}
	return extiter.Concat(sequences...)
}
