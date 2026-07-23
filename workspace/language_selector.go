// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/service"
)

// DefaultLanguageSelector selects a language by matching the document against an
// ordered list of [core.DocumentSelector] values, returning the first match's index.
// It takes the language id of the document that the [DocumentManager] holds for
// the URI first and matches it against every selector's language id; only when
// there is no such document or no selector claims its id does it fall back to
// the URI-path globs (in selector order).
//
// The managed document is the one that is being or was last built, so the lexer,
// the parser and request-time consumers such as completion resolve the language
// the document was built with, even while the text-document store already holds
// a handle with another language id.
type DefaultLanguageSelector struct {
	sc        *service.Container
	selectors []core.DocumentSelector
}

// NewDefaultLanguageSelector returns a new [DefaultLanguageSelector]
// that matches documents against the given selectors in order.
func NewDefaultLanguageSelector(sc *service.Container, selectors ...core.DocumentSelector) *DefaultLanguageSelector {
	return &DefaultLanguageSelector{sc: sc, selectors: selectors}
}

func (s *DefaultLanguageSelector) Select(uri core.URI) (int, string) {
	if documents, err := service.Get[DocumentManager](s.sc); err == nil {
		if doc := documents.Get(uri); doc != nil {
			if languageID := doc.TextDoc.LanguageID(); languageID != "" {
				for i, sel := range s.selectors {
					if sel.LanguageID == languageID {
						return i, sel.LanguageID
					}
				}
			}
		}
	}
	for i, sel := range s.selectors {
		if sel.Matcher != nil && sel.Matcher(uri) {
			return i, sel.LanguageID
		}
	}
	return -1, ""
}
