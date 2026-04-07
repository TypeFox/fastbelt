// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package e2e_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// TestScaffoldModule_endToEnd runs the fastbelt CLI against a temp directory, then npm build
// steps and a short-lived go run of the generated LSP. Skipped under -short or when node/npm
// are unavailable. Requires network for go get inside scaffold.
func TestScaffoldModule_endToEnd(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e scaffold test in short mode")
	}
	skipIfMissingTool(t, "go")
	skipIfMissingTool(t, "node")
	skipIfMissingTool(t, "npm")

	fastbeltBin := filepath.Join(t.TempDir(), fastbeltExe())
	workDir := filepath.Join(t.TempDir(), "work")
	require.NoError(t, os.MkdirAll(workDir, 0755))

	modulePath := "example.com/fastbelte2e/e2e" + strconv.FormatInt(time.Now().UnixNano(), 10)
	// Go module paths use slash syntax; path.Base, not filepath.Base (Windows).
	moduleRoot := filepath.Join(workDir, path.Base(modulePath))
	repoRoot := repoRoot(t)

	execCmd := func(t *testing.T, cmd *exec.Cmd) {
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			t.Fatalf("%s: %v\n%s", cmd.String(), err, out.String())
		}
	}

	t.Run("build fastbelt", func(t *testing.T) {
		cmd := exec.Command("go", "build", "-o", fastbeltBin, "./cmd/fastbelt")
		cmd.Dir = repoRoot
		execCmd(t, cmd)
	})

	t.Run("scaffold module", func(t *testing.T) {
		scaffold := exec.Command(fastbeltBin, "scaffold", "-module", modulePath, "-language", "E2E Lang")
		scaffold.Dir = workDir
		// Local fastbelt builds often embed a non-proxy pseudo-version (+dirty); pin a resolvable version for go get.
		scaffold.Env = append(os.Environ(), "FASTBELT_SCAFFOLD_FASTBELT_GO_VERSION=latest")
		execCmd(t, scaffold)
		require.DirExists(t, moduleRoot)
	})

	for _, step := range []struct {
		name string
		args []string
	}{
		{"npm install", []string{"install"}},
		{"npm run build", []string{"run", "build"}},
		{"npm run bundle", []string{"run", "bundle"}},
	} {
		step := step
		t.Run(step.name, func(t *testing.T) {
			cmd := exec.Command("npm", step.args...)
			cmd.Dir = moduleRoot
			execCmd(t, cmd)
		})
	}

	t.Run("go run lsp", func(t *testing.T) {
		lspDir := singleCmdSubdir(t, moduleRoot)

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		// go run requires a "./" prefix so the path is treated as a package directory, not an import path.
		lspPkg := "./" + filepath.ToSlash(filepath.Join("cmd", lspDir))
		cmd := exec.CommandContext(ctx, "go", "run", lspPkg)
		cmd.Dir = moduleRoot
		stdinR, stdinW := io.Pipe()
		cmd.Stdin = stdinR
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		cmd.Stdout = io.Discard

		require.NoError(t, cmd.Start())
		waitErr := make(chan error, 1)
		go func() { waitErr <- cmd.Wait() }()

		select {
		case <-time.After(5 * time.Second):
			cancel()
		case err := <-waitErr:
			_ = stdinW.Close()
			if err != nil {
				t.Fatalf("LSP exited early: %v stderr=%q", err, stderr.String())
			}
			t.Fatal("LSP exited cleanly before timeout (expected it to block reading stdio)")
		}
		_ = stdinW.Close()
		err := <-waitErr
		require.Error(t, err, "expected LSP process to stop after context cancel; stderr=%q", stderr.String())
		// CommandContext may surface context.Canceled, or Wait may report signal: killed.
		ok := errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded)
		if !ok {
			var exitErr *exec.ExitError
			ok = errors.As(err, &exitErr)
		}
		require.True(t, ok, "unexpected wait result: %v stderr=%q", err, stderr.String())
	})
}

func repoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok)
	// internal/e2e_test/*.go -> repository root
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func fastbeltExe() string {
	if runtime.GOOS == "windows" {
		return "fastbelt.exe"
	}
	return "fastbelt"
}

func skipIfMissingTool(t *testing.T, name string) {
	t.Helper()
	if _, err := exec.LookPath(name); err != nil {
		t.Skipf("%s not in PATH: %v", name, err)
	}
}

func singleCmdSubdir(t *testing.T, moduleRoot string) string {
	t.Helper()
	require.DirExists(t, moduleRoot)
	cmdDir := filepath.Join(moduleRoot, "cmd")
	entries, err := os.ReadDir(cmdDir)
	require.NoError(t, err)
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	require.Len(t, names, 1, "expected exactly one cmd/* subdirectory, got %v", names)
	return names[0]
}
