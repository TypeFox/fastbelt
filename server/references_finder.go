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
	// If set to a non-nil value, the references finder will only return references that are located in the document with the given URI.
	//
	// If IncludeDeclaration is true, the declaration will only be included if it is located in the document with the given URI.
	TargetURI core.URI
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
		selfReference := findSelfReference(target, options)
		if selfReference != nil {
			sequences = append(sequences, extiter.Of(selfReference))
		}
	}
	documentManager := service.MustGet[workspace.DocumentManager](rf.sc)
	if options.TargetURI == nil {
		// Iterate through all documents and collect references to the symbol
		for doc := range documentManager.All() {
			refDescriptions := doc.ReferenceDescriptions.ForTarget(target)
			sequences = append(sequences, refDescriptions)
		}
	} else {
		// Only look for references in the document with the given URI
		doc := documentManager.Get(options.TargetURI)
		if doc != nil {
			refDescriptions := doc.ReferenceDescriptions.ForTarget(target)
			sequences = append(sequences, refDescriptions)
		}
	}

	return extiter.Concat(sequences...)
}

func findSelfReference(target core.AstNode, options FindReferencesOptions) *core.ReferenceDescription {
	if options.TargetURI != nil {
		document := target.Document()
		if document == nil || !document.URI.Equal(options.TargetURI) {
			return nil
		}
	}
	nameUnit := linking.Name(target)
	if nameUnit != nil {
		selfDescription := core.NewReferenceDescription(target, target, nameUnit.Segment())
		return selfDescription
	}
	return nil
}
