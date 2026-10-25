// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package e2e_test

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
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
		scaffold := exec.Command(fastbeltBin, "scaffold", "--module", modulePath, "--language", "E2E Lang", "--vscode")
		scaffold.Dir = workDir
		// Use the local repo as a go mod replace so the scaffolded module builds against in-tree code.
		scaffold.Env = append(os.Environ(), "FASTBELT_SCAFFOLD_REPO="+repoRoot)
		execCmd(t, scaffold)
		require.DirExists(t, moduleRoot)
	})

	for _, step := range []struct {
		name string
		args []string
	}{
		// vsce is only needed for "npm run package", which this test never runs; it is by far the heaviest dependency.
		{"drop vsce", []string{"pkg", "delete", "devDependencies.@vscode/vsce"}},
		{"npm install", []string{"install", "--no-audit", "--no-fund", "--prefer-offline"}},
		{"npm run build", []string{"run", "build"}},
		{"npm run bundle", []string{"run", "bundle"}},
	} {
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
		cmd.WaitDelay = 10 * time.Second
		stdin, err := cmd.StdinPipe()
		require.NoError(t, err)
		stdout, err := cmd.StdoutPipe()
		require.NoError(t, err)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		require.NoError(t, cmd.Start())

		// Send an initialize request and wait for the first framed response instead of sleeping:
		// proves the server started and answers. EOF on stdout means it exited early.
		body := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"capabilities":{}}}`
		_, err = fmt.Fprintf(stdin, "Content-Length: %d\r\n\r\n%s", len(body), body)
		require.NoError(t, err)
		responded := make(chan error, 1)
		go func() {
			_, err := bufio.NewReader(stdout).ReadString('\n')
			responded <- err
		}()
		var readErr error
		select {
		case readErr = <-responded:
		case <-ctx.Done():
			readErr = ctx.Err()
		}
		// Stop the server and wait for exec's stderr copier before reading the buffer (data race otherwise).
		_ = stdin.Close()
		cancel()
		_ = cmd.Wait()
		require.NoError(t, readErr, "LSP produced no initialize response; stderr=%q", stderr.String())
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
