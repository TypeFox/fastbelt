// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"path"

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
	// changedURIs holds unencoded URI strings (see workspace.DocumentUpdater),
	// so the folder is compared on the string: Affected runs once per document
	// in the workspace and parsing every changed URI each time adds up.
	folder := path.Dir(doc.URI.StringUnencoded())
	for changed := range changedURIs {
		if path.Dir(changed) == folder {
			return true
		}
	}
	return s.DocumentChangeImpact.Affected(doc, changedURIs)
}
