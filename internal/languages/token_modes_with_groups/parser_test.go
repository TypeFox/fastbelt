// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package token_modes_with_groups

import (
	"testing"

	"github.com/stretchr/testify/require"
	"typefox.dev/fastbelt/test"
)

func TestXX(t *testing.T) {
	fixture := test.New(t, CreateServices())
	doc := fixture.Parse("xx")
	require.Len(t, doc.Document.Tokens, 2)
}
