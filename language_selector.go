// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

import (
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/util/glob"
	"typefox.dev/fastbelt/util/service"
)

type DocumentMatcher func(uri URI) bool

// DocumentFilter matches a document by LSP language id and/or a glob over the
// document URI path. A filter matches when EITHER a non-empty LanguageID equals
// the document's language id OR a non-empty Pattern matches the path (see
// [DocumentSelector.Matches]) — so one filter per language (its id plus its
// extension glob) is enough.
type DocumentSelector struct {
	LanguageID string
	Matcher    DocumentMatcher
}

func NewDocumentSelector(languageID string, matcher DocumentMatcher) DocumentSelector {
	return DocumentSelector{LanguageID: languageID, Matcher: matcher}
}

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

// Matches reports whether the selector matches a document with the given
// language id and URI path. Language id is tried first (open documents carry a
// client-supplied id), then the path glob (covers files loaded from disk).
func (d DocumentSelector) Matches(languageID string, uri URI) bool {
	if languageID != "" && d.LanguageID == languageID {
		return true
	}
	if d.Matcher != nil && d.Matcher(uri) {
		return true
	}
	return false
}

// LanguageSelector resolves a document URI to the index of the owning language
// (into the configured languages), or -1 if none match.
type LanguageSelector interface {
	Select(uri URI) (int, string)
}

// DefaultLanguageSelector selects a language by matching the document against an
// ordered list of [DocumentSelector] values, returning the first match's index.
// It resolves the document's language id from the [textdoc.Store] first, then
// falls back to the URI-path glob. The store is looked up lazily from the
// container so the selector can be registered before the container is sealed.
type DefaultLanguageSelector struct {
	sc        *service.Container
	selectors []DocumentSelector
}

// NewDefaultLanguageSelector returns a [LanguageSelector] that resolves the
// [textdoc.Store] from sc and matches documents against the given per-language
// selectors.
func NewDefaultLanguageSelector(sc *service.Container, selectors ...DocumentSelector) *DefaultLanguageSelector {
	return &DefaultLanguageSelector{sc: sc, selectors: selectors}
}

func (s *DefaultLanguageSelector) Select(uri URI) (int, string) {
	languageID := ""
	// If a document handle already exists for the given URI, use its language id.
	if store, err := service.Get[textdoc.Store](s.sc); err == nil {
		if handle := store.Get(uri.DocumentURI()); handle != nil {
			languageID = handle.LanguageID()
		}
	}
	for i, sel := range s.selectors {
		if sel.Matches(languageID, uri) {
			return i, sel.LanguageID
		}
	}
	return -1, ""
}
