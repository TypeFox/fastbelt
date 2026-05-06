// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package linking

import (
	"context"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/collections"
	"typefox.dev/fastbelt/util/service"
)

// ReferenceDescriptionsProvider is a service that computes the reference descriptions for a document.
type ReferenceDescriptionsProvider interface {
	// ReferenceDescriptions computes the reference descriptions for a document.
	// The result is stored in the document's ReferenceDescriptions field.
	ReferenceDescriptions(ctx context.Context, document *core.Document) core.ReferenceDescriptions
}

// DefaultReferenceDescriptionsProvider is the default implementation of [ReferenceDescriptionsProvider].
type DefaultReferenceDescriptionsProvider struct {
	sc *service.Container
}

func NewDefaultReferenceDescriptionsProvider(sc *service.Container) ReferenceDescriptionsProvider {
	return &DefaultReferenceDescriptionsProvider{sc: sc}
}

func (s *DefaultReferenceDescriptionsProvider) ReferenceDescriptions(ctx context.Context, document *core.Document) core.ReferenceDescriptions {
	// Unlike AST node descriptions, reference descriptions can't be associated with specific node types,
	// so the [ReferenceDescriber] interface is used as a service.
	describer := service.MustGet[ReferenceDescriber](s.sc)
	descriptions := collections.NewMultiMap[core.AstNode, *core.ReferenceDescription]()
	for _, ref := range document.References {
		node := ref.RefNode(ctx)
		if node != nil {
			description := describer.DescribeReference(ctx, ref)
			if description != nil {
				descriptions.Put(node, description)
			}
		}
	}
	refDescriptions := core.NewReferenceDescriptionsFromMap(descriptions)
	document.ReferenceDescriptions = refDescriptions
	return refDescriptions
}

// ReferenceDescriber is a service that describes references.
type ReferenceDescriber interface {
	// DescribeReference describes metadata about a reference, like the source and target nodes.
	DescribeReference(ctx context.Context, ref core.UntypedReference) *core.ReferenceDescription
}

// DefaultReferenceDescriber is the default implementation of [ReferenceDescriber].
type DefaultReferenceDescriber struct {
	sc *service.Container
}

func NewDefaultReferenceDescriber(sc *service.Container) ReferenceDescriber {
	return &DefaultReferenceDescriber{sc: sc}
}

func (s *DefaultReferenceDescriber) DescribeReference(ctx context.Context, ref core.UntypedReference) *core.ReferenceDescription {
	source := ref.Owner()
	target := ref.RefNode(ctx)
	segment := ref.Segment()
	if source == nil || target == nil || segment == nil {
		return nil
	}
	return core.NewReferenceDescription(source, target, segment)
}
