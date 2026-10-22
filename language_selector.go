// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

// LanguageSelector maps documents to the languages of a build.
// The default implementation lives in the workspace package.
type LanguageSelector interface {
	// Select returns the index of the language that owns the document with the
	// given URI and its language ID. Returns -1 and "" if no language matches.
	Select(uri URI) (int, string)
}
