// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package linking

import (
	"context"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/collections"
)

type ReferenceDescriptionsProvider interface {
	Provide(ctx context.Context, document *core.Document)
}

type DefaultReferenceDescriptionsProvider struct {
	srv LinkingSrvCont
}

func NewDefaultReferenceDescriptionsProvider(srv LinkingSrvCont) ReferenceDescriptionsProvider {
	return &DefaultReferenceDescriptionsProvider{
		srv: srv,
	}
}

func (p *DefaultReferenceDescriptionsProvider) Provide(ctx context.Context, document *core.Document) {
	descriptions := collections.NewMultiMap[core.AstNode, *core.ReferenceDescription]()
	describer := p.srv.Linking().ReferenceDescriber
	for _, ref := range document.References {
		node := ref.RefNode(ctx)
		if node != nil {
			description := describer.Describe(ctx, ref)
			if description != nil {
				descriptions.Put(node, description)
			}
		}
	}
	document.ReferenceDescriptions = core.NewReferenceDescriptionsFromMap(descriptions)
}

type ReferenceDescriber interface {
	Describe(ctx context.Context, ref core.UntypedReference) *core.ReferenceDescription
}

type DefaultReferenceDescriber struct{}

func NewDefaultReferenceDescriber() ReferenceDescriber {
	return &DefaultReferenceDescriber{}
}

func (d *DefaultReferenceDescriber) Describe(ctx context.Context, ref core.UntypedReference) *core.ReferenceDescription {
	source := ref.Owner()
	target := ref.RefNode(ctx)
	segment := ref.Segment()
	if source == nil || target == nil || segment == nil {
		return nil
	}
	return core.NewReferenceDescription(source, target, segment)
}
