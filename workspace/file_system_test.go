// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"context"
	"errors"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	core "typefox.dev/fastbelt"
)

// walkFS builds a small tree and returns the sorted set of visited URIs,
// applying fn (may be nil) to control the walk.
func walkVisited(t *testing.T, fs *VirtualFileSystem, root core.URI, fn func(DirEntry) error) []string {
	t.Helper()
	var visited []string
	err := WalkFileSystem(context.Background(), fs, root, func(e DirEntry) error {
		visited = append(visited, e.URI.String())
		if fn != nil {
			return fn(e)
		}
		return nil
	})
	require.NoError(t, err)
	sort.Strings(visited)
	return visited
}

func TestWalkFileSystemVisitsAll(t *testing.T) {
	fs := NewVirtualFileSystem()
	require.NoError(t, fs.WriteFile(core.FileURI("/a/b/c.txt"), []byte("c")))
	require.NoError(t, fs.WriteFile(core.FileURI("/a/d.txt"), []byte("d")))

	visited := walkVisited(t, fs, core.FileURI("/a"), nil)
	assert.Equal(t, []string{
		"file:///a",
		"file:///a/b",
		"file:///a/b/c.txt",
		"file:///a/d.txt",
	}, visited)
}

func TestWalkFileSystemVisitsAllInAlphabeticalOrder(t *testing.T) {
	fs := NewVirtualFileSystem()
	require.NoError(t, fs.WriteFile(core.FileURI("/workspace/a.txt"), []byte("a")))
	require.NoError(t, fs.WriteFile(core.FileURI("/workspace/d.txt"), []byte("d")))
	require.NoError(t, fs.WriteFile(core.FileURI("/workspace/c.txt"), []byte("c")))
	require.NoError(t, fs.WriteFile(core.FileURI("/workspace/b.txt"), []byte("b")))

	visited := walkVisited(t, fs, core.FileURI("/workspace"), nil)
	assert.Equal(t, []string{
		"file:///workspace",
		"file:///workspace/a.txt",
		"file:///workspace/b.txt",
		"file:///workspace/c.txt",
		"file:///workspace/d.txt",
	}, visited)
}

func TestWalkFileSystemSingleFile(t *testing.T) {
	fs := NewVirtualFileSystem()
	require.NoError(t, fs.WriteFile(core.FileURI("/a/b/c.txt"), []byte("c")))

	visited := walkVisited(t, fs, core.FileURI("/a/b/c.txt"), nil)
	assert.Equal(t, []string{"file:///a/b/c.txt"}, visited)
}

func TestWalkFileSystemSkipDir(t *testing.T) {
	fs := NewVirtualFileSystem()
	require.NoError(t, fs.WriteFile(core.FileURI("/a/b/c.txt"), []byte("c")))
	require.NoError(t, fs.WriteFile(core.FileURI("/a/d.txt"), []byte("d")))

	// Skip directory b: its contents must not be visited, but siblings are.
	visited := walkVisited(t, fs, core.FileURI("/a"), func(e DirEntry) error {
		if e.URI.String() == "file:///a/b" {
			return SkipDir
		}
		return nil
	})
	assert.Equal(t, []string{
		"file:///a",
		"file:///a/b",
		"file:///a/d.txt",
	}, visited)
}

func TestWalkFileSystemSkipFile(t *testing.T) {
	fs := NewVirtualFileSystem()
	require.NoError(t, fs.WriteFile(core.FileURI("/a/b/c.txt"), []byte("c")))
	require.NoError(t, fs.WriteFile(core.FileURI("/a/d.txt"), []byte("d")))
	require.NoError(t, fs.WriteFile(core.FileURI("/a/e.txt"), []byte("e")))

	// Skip file d: will prevent further siblings from being visited
	visited := walkVisited(t, fs, core.FileURI("/a"), func(e DirEntry) error {
		if e.URI.String() == "file:///a/d.txt" {
			return SkipDir
		}
		return nil
	})
	assert.Equal(t, []string{
		"file:///a",
		"file:///a/b",
		"file:///a/b/c.txt",
		"file:///a/d.txt",
	}, visited)
}

func TestWalkFileSystemSkipAll(t *testing.T) {
	fs := NewVirtualFileSystem()
	require.NoError(t, fs.WriteFile(core.FileURI("/a/b/c.txt"), []byte("c")))
	require.NoError(t, fs.WriteFile(core.FileURI("/a/d.txt"), []byte("d")))

	count := 0
	err := WalkFileSystem(context.Background(), fs, core.FileURI("/a"), func(e DirEntry) error {
		count++
		return SkipAll
	})
	require.NoError(t, err)
	assert.Equal(t, 1, count) // aborted after the root
}

func TestWalkFileSystemErrorPropagates(t *testing.T) {
	fs := NewVirtualFileSystem()
	require.NoError(t, fs.WriteFile(core.FileURI("/a/b/c.txt"), []byte("c")))

	sentinel := errors.New("boom")
	err := WalkFileSystem(context.Background(), fs, core.FileURI("/a"), func(e DirEntry) error {
		return sentinel
	})
	assert.ErrorIs(t, err, sentinel)
}

func TestWalkFileSystemMissingRoot(t *testing.T) {
	fs := NewVirtualFileSystem()
	err := WalkFileSystem(context.Background(), fs, core.FileURI("/nope"), func(e DirEntry) error {
		return nil
	})
	assert.Error(t, err)
}

func TestVirtualFileSystemWriteCreatesParents(t *testing.T) {
	ctx := context.Background()
	fs := NewVirtualFileSystem()
	file := core.FileURI("/a/b/c.txt")

	require.NoError(t, fs.WriteFile(file, []byte("hello")))

	content, err := fs.ReadFile(ctx, file)
	require.NoError(t, err)
	assert.Equal(t, []byte("hello"), content)

	// Intermediate parents are real, listable directories.
	dir := core.FileURI("/a/b")
	stat, err := fs.Stat(ctx, dir)
	require.NoError(t, err)
	assert.True(t, stat.IsDir)
}

func TestVirtualFileSystemReadDir(t *testing.T) {
	ctx := context.Background()
	fs := NewVirtualFileSystem()
	require.NoError(t, fs.WriteFile(core.FileURI("/a/b/c.txt"), []byte("hello")))

	entries, err := fs.ReadDir(ctx, core.FileURI("/a/b"))
	require.NoError(t, err)
	require.Len(t, entries, 1)
	assert.True(t, entries[0].IsFile)
}

func TestVirtualFileSystemOverwrite(t *testing.T) {
	ctx := context.Background()
	fs := NewVirtualFileSystem()
	file := core.FileURI("/a/b/c.txt")

	require.NoError(t, fs.WriteFile(file, []byte("hello")))
	require.NoError(t, fs.WriteFile(file, []byte("world")))

	content, err := fs.ReadFile(ctx, file)
	require.NoError(t, err)
	assert.Equal(t, []byte("world"), content)
}

func TestVirtualFileSystemMkdirExistingFails(t *testing.T) {
	fs := NewVirtualFileSystem()
	require.NoError(t, fs.WriteFile(core.FileURI("/a/b/c.txt"), []byte("hello")))

	assert.Error(t, fs.Mkdir(core.FileURI("/a/b")))
}

func TestVirtualFileSystemFileDirTypeMismatch(t *testing.T) {
	ctx := context.Background()
	fs := NewVirtualFileSystem()
	require.NoError(t, fs.WriteFile(core.FileURI("/a/b/c.txt"), []byte("hello")))
	dir := core.FileURI("/a/b")

	// Cannot overwrite a directory with a file.
	assert.Error(t, fs.WriteFile(dir, []byte("x")))
	// Cannot read a directory as a file.
	_, err := fs.ReadFile(ctx, dir)
	assert.Error(t, err)
}

func TestVirtualFileSystemDeleteRecursive(t *testing.T) {
	ctx := context.Background()
	fs := NewVirtualFileSystem()
	file := core.FileURI("/a/b/c.txt")
	require.NoError(t, fs.WriteFile(file, []byte("hello")))

	require.NoError(t, fs.Delete(core.FileURI("/a")))
	assert.False(t, fs.Exists(ctx, file))
	assert.False(t, fs.Exists(ctx, core.FileURI("/a/b")))
}

func TestVirtualFileSystemDeleteMissingFails(t *testing.T) {
	fs := NewVirtualFileSystem()
	assert.Error(t, fs.Delete(core.FileURI("/a")))
}
