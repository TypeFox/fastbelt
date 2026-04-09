// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunner_ParseArgs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		args    []string
		want    string
		wantErr bool
	}{
		{name: "one path", args: []string{"/tmp/x.statemachine"}, want: "/tmp/x.statemachine"},
		{name: "zero args", args: nil, wantErr: true},
		{name: "empty args", args: []string{}, wantErr: true},
		{name: "too many", args: []string{"a", "b"}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := &Runner{ProgName: "statemachine"}
			err := r.ParseArgs(tt.args)
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), "usage:")
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, r.InputPath)
		})
	}
}

func TestRunner_LoadSource(t *testing.T) {
	t.Parallel()
	dir := t.TempDir()
	p := filepath.Join(dir, "m.statemachine")
	require.NoError(t, os.WriteFile(p, []byte("statemachine X\n"), 0644))
	r := &Runner{InputPath: p}
	require.NoError(t, r.LoadSource())
	assert.Equal(t, "statemachine X\n", string(r.Content))
	abs, absErr := filepath.Abs(p)
	require.NoError(t, absErr)
	assert.Equal(t, abs, r.SourcePath)
}

func TestRunner_valid_examples(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		sourcePath    string
		eventScript   string
		wantStdout    []string
		stderrNoLexer bool
	}{
		{
			name:        "traffic_light",
			sourcePath:  "../../traffic_light.statemachine",
			eventScript: "switchCapacity\nnext\n",
			wantStdout: []string{
				`State machine: "TrafficLight"`,
				`Start state: "PowerOff"`,
				`"switchCapacity" -> "RedLight"`,
				`"next" -> "GreenLight"`,
			},
			stderrNoLexer: true,
		},
		{
			name:        "elevator",
			sourcePath:  "../../elevator.statemachine",
			eventScript: "requestCar\natFloor\n",
			wantStdout: []string{
				`State machine: "Elevator"`,
				`Commands: bell`,
				`  - "Waiting"  actions: bell`,
				`Start state: "Waiting"`,
				`"requestCar" -> "Moving"`,
				`"atFloor" -> "Waiting"`,
				`  emit command "bell"`,
			},
			stderrNoLexer: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var out, errOut bytes.Buffer
			r := &Runner{
				InputPath:  tt.sourcePath,
				Stdout:     &out,
				Stderr:     &errOut,
				EventInput: strings.NewReader(tt.eventScript),
			}
			require.NoError(t, r.LoadSource())
			require.NoError(t, r.ParseAndValidate())
			require.NoError(t, r.Run())
			s := out.String()
			for _, frag := range tt.wantStdout {
				assert.Contains(t, s, frag, "stdout should contain %q", frag)
			}
			if tt.stderrNoLexer {
				assert.NotContains(t, errOut.String(), "lexer ")
				assert.NotContains(t, errOut.String(), "parser ")
			}
		})
	}
}
