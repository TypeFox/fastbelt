// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"github.com/gobwas/glob"
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/service"
)

// DocumentMatcher assigns documents to a language, either by LSP language ID
// or by matching the document URI.
type DocumentMatcher struct {
	// LanguageID is the LSP language ID of the language.
	LanguageID string
	// Matcher returns true if the document with the given URI belongs to the language.
	Matcher func(uri core.URI) bool
}

// NewDocumentMatcherWithPatterns returns a DocumentMatcher for the given language ID
// that matches document URI paths against a list of glob patterns.
func NewDocumentMatcherWithPatterns(languageID string, patterns ...string) DocumentMatcher {
	globs := make([]*glob.Pattern, len(patterns))
	for i, pattern := range patterns {
		globs[i] = glob.MustCompile(pattern)
	}
	return DocumentMatcher{LanguageID: languageID, Matcher: func(uri core.URI) bool {
		for _, glob := range globs {
			if glob.Match(uri.Path()) {
				return true
			}
		}
		return false
	}}
}

// DefaultLanguageSelector is the default implementation of the LanguageSelector.
// It holds one DocumentMatcher per language, in the order of the languages of the build.
// Note: The list lives here because the service container has no multi-binding.
//
// Documents managed by the DocumentManager are matched by their language ID first
// and fall back to the URI matchers. URIs of unmanaged documents are matched by
// the URI matchers only.
type DefaultLanguageSelector struct {
	sc        *service.Container
	selectors []DocumentMatcher
}

// NewDefaultLanguageSelector returns a new DefaultLanguageSelector
// that matches documents against the given matchers in order.
func NewDefaultLanguageSelector(sc *service.Container, selectors ...DocumentMatcher) *DefaultLanguageSelector {
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
