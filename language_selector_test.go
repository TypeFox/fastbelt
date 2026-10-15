// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDocumentSelectorWithPatterns(t *testing.T) {
	selector := NewDocumentSelectorWithPatterns("greeting", "**/*.hello", "*.hi")

	assert.True(t, selector.Matcher(ParseURI("file:///ws/a.hello")))
	// A pattern without a '/' matches the file name in any directory.
	assert.True(t, selector.Matcher(ParseURI("file:///ws/sub/a.hi")))
	assert.False(t, selector.Matcher(ParseURI("file:///ws/a.bye")))
}
