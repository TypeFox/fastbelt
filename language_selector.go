// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

import (
	"path"
	"strings"

	"typefox.dev/fastbelt/util/glob"
)

// DocumentSelector matches a document by LSP language id and/or a predicate
// over the document URI (typically a glob over its path).
type DocumentSelector struct {
	LanguageID string
	Matcher    func(uri URI) bool
}

// NewDocumentSelectorWithPatterns returns a [DocumentSelector] that matches documents by
// language id and a list of glob patterns over the document URI path. The patterns
// are matched in order, and the first match wins. A pattern without a '/' (such as
// "*.hello") is matched against the file name only; any other pattern is matched
// against the absolute path, so it typically starts with "**/".
func NewDocumentSelectorWithPatterns(languageID string, patterns ...string) DocumentSelector {
	return DocumentSelector{LanguageID: languageID, Matcher: func(uri URI) bool {
		for _, pattern := range patterns {
			name := uri.Path()
			if !strings.Contains(pattern, "/") {
				name = path.Base(name)
			}
			if glob.Match(pattern, name) {
				return true
			}
		}
		return false
	}}
}

// LanguageSelector resolves a document URI to the index of the owning language
// (into the configured languages) and that language's id, or -1 if none match.
type LanguageSelector interface {
	Select(uri URI) (int, string)
}
