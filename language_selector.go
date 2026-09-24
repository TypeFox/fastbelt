// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

import (
	"sync"

	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/util/glob"
	"typefox.dev/fastbelt/util/service"
)

// DocumentMatcher is a function that returns true if the given document URI matches a
// specific pattern. It is used by [DocumentSelector] to match documents by URI path.
type DocumentMatcher func(uri URI) bool

// DocumentSelector matches a document by LSP language id and/or a glob over the
// document URI path.
type DocumentSelector struct {
	LanguageID string
	Matcher    DocumentMatcher
}

// NewDocumentSelector returns a [DocumentSelector] that matches documents by
// language id and a custom matcher function.
func NewDocumentSelector(languageID string, matcher DocumentMatcher) DocumentSelector {
	return DocumentSelector{LanguageID: languageID, Matcher: matcher}
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
// It resolves the document's language id from the [textdoc.Store] first and
// matches it against every selector's language id; only when no selector claims
// that id does it fall back to the URI-path globs (in selector order). The store
// is looked up lazily from the container so the selector can be registered
// before the container is sealed.
type DefaultLanguageSelector struct {
	sc        *service.Container
	selectors []DocumentSelector

	storeOnce sync.Once
	store     textdoc.Store // nil when no store is registered
}

// NewDefaultLanguageSelector returns a new [DefaultLanguageSelector]
// that matches documents against the given selectors in order.
func NewDefaultLanguageSelector(sc *service.Container, selectors ...DocumentSelector) *DefaultLanguageSelector {
	return &DefaultLanguageSelector{sc: sc, selectors: selectors}
}

// textdocStore returns the registered [textdoc.Store], resolving it from the
// container on first use. Select runs on the build hot path (once per document
// for the lexer, the parser and the completion parser), so the container lookup
// is not repeated.
func (s *DefaultLanguageSelector) textdocStore() textdoc.Store {
	s.storeOnce.Do(func() {
		if s.sc == nil {
			return
		}
		if store, err := service.Get[textdoc.Store](s.sc); err == nil {
			s.store = store
		}
	})
	return s.store
}

func (s *DefaultLanguageSelector) Select(uri URI) (int, string) {
	// If a document handle already exists for the given URI, its (client
	// supplied) language id takes precedence over any path glob.
	if store := s.textdocStore(); store != nil {
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
