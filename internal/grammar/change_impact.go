// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/collections"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
)

// changeImpactImpl marks every grammar file of a folder as affected when one
// of them changes. All .fb files of a folder form one grammar, so folder-wide
// validations (unique names, token coverage, ...) of the siblings must run
// again even when they hold no reference into the changed file.
type changeImpactImpl struct {
	workspace.DocumentChangeImpact
}

func newChangeImpactImpl(sc *service.Container) workspace.DocumentChangeImpact {
	return &changeImpactImpl{DocumentChangeImpact: workspace.NewDefaultDocumentChangeImpact(sc)}
}

func (s *changeImpactImpl) Affected(doc *core.Document, changedURIs collections.Set[string]) bool {
	for changed := range changedURIs {
		if sameFolder(doc.URI, core.ParseURI(changed)) {
			return true
		}
	}
	return s.DocumentChangeImpact.Affected(doc, changedURIs)
}
