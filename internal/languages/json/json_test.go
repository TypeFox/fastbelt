// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package json

import (
	"encoding/json/v2"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestUnmarshalZeroFieldInterface checks the adherence to json.UnmarshalerFrom's contract:
// json v2 requires an UnmarshalerFrom implementation to read exactly one value;
// with nothing consumed, json.Unmarshal fails with "must read or write exactly one
// value" for any zero-field marker interface.
// ([Empty] here, or e.g. arithmetics.Expression / arithmetics.Statement)
func TestUnmarshalZeroFieldInterface(t *testing.T) {
	t.Run("json.Unmarshal", func(t *testing.T) {
		err := json.Unmarshal([]byte(`{"$type":"Empty"}`), NewEmpty())
		require.NoError(t, err)
	})

	t.Run("UnmarshalValue", func(t *testing.T) {
		value, err := UnmarshalValue[Empty]([]byte(`{"$type":"Empty"}`))
		require.NoError(t, err)
		require.NotNil(t, value)
		require.IsType(t, NewEmpty(), value)
	})
}
