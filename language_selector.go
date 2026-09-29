// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

import (
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/util/glob"
	"typefox.dev/fastbelt/util/service"
)

// DocumentSelector matches a document by LSP language id and/or a predicate
// over the document URI (typically a glob over its path).
type DocumentSelector struct {
	LanguageID string
	Matcher    func(uri URI) bool
}

// NewDocumentSelectorWithPatterns returns a [DocumentSelector] that matches documents by
// language id and a list of glob patterns over the document URI path. The patterns
// are matched in order, and the first match wins.
func NewDocumentSelectorWithPatterns(languageID string, patterns ...string) DocumentSelector {
	return DocumentSelector{LanguageID: languageID, Matcher: func(uri URI) bool {
		for _, pattern := range patterns {
			if glob.Match(pattern, uri.Path()) {
				return true
			}
		}
		return false
	}}
}

// LanguageSelector resolves a document URI to the index of the owning language
// (into the configured languages), or -1 if none match.
type LanguageSelector interface {
	Select(uri URI) (int, string)
}

// DefaultLanguageSelector selects a language by matching the document against an
// ordered list of [DocumentSelector] values, returning the first match's index.
// It resolves the document's language id from the [textdoc.Store] first and
// matches it against every selector's language id; only when no selector claims
// that id does it fall back to the URI-path globs (in selector order).
type DefaultLanguageSelector struct {
	sc        *service.Container
	selectors []DocumentSelector
}

// NewDefaultLanguageSelector returns a new [DefaultLanguageSelector]
// that matches documents against the given selectors in order.
func NewDefaultLanguageSelector(sc *service.Container, selectors ...DocumentSelector) *DefaultLanguageSelector {
	return &DefaultLanguageSelector{sc: sc, selectors: selectors}
}

func (s *DefaultLanguageSelector) Select(uri URI) (int, string) {
	// If a document handle already exists for the given URI, its (client
	// supplied) language id takes precedence over any path glob.
	if store, err := service.Get[textdoc.Store](s.sc); err == nil {
		if handle := store.Get(uri.DocumentURI()); handle != nil {
			if languageID := handle.LanguageID(); languageID != "" {
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
