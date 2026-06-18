// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"context"
	"io/fs"
	"os"
	"strings"

	core "typefox.dev/fastbelt"
)

// DirEntry represents a file system entry, which can be a file, directory, or/and symbolic link.
type DirEntry struct {
	// IsFile is true if the entry is a regular file.
	IsFile bool
	// IsDir is true if the entry is a directory.
	IsDir bool
	// IsSymlink is true if the entry is a symbolic link.
	IsSymlink bool
	// URI is the URI of the entry.
	URI core.URI
}

// FileSystem defines the interface for reading data from an arbitrary file system.
type FileSystem interface {
	// Exists checks if the given URI exists in the file system.
	Exists(ctx context.Context, uri core.URI) bool
	// Stat returns the file system entry for the given URI, or an error if it does not exist.
	Stat(ctx context.Context, uri core.URI) (DirEntry, error)
	// ReadFile reads the content of the file at the given URI and returns it as a byte slice.
	ReadFile(ctx context.Context, uri core.URI) ([]byte, error)
	// ReadDir reads the contents of the directory at the given URI and returns a slice of DirEntry for its entries.
	ReadDir(ctx context.Context, uri core.URI) ([]DirEntry, error)
}

// SkipDir indicates that the current directory should be skipped.
var SkipDir = fs.SkipDir

// SkipAll indicates that the entire file system walk should be aborted.
var SkipAll = fs.SkipAll

func WalkFileSystem(ctx context.Context, fs FileSystem, uri core.URI, fn func(DirEntry) error) error {
	entry, err := fs.Stat(ctx, uri)
	if err != nil {
		return err
	} else {
		err = walkDir(ctx, fs, entry, fn)
	}
	if err == SkipAll || err == SkipDir {
		return nil
	}
	return err
}

func walkDir(ctx context.Context, fs FileSystem, entry DirEntry, fn func(DirEntry) error) error {
	if ctx.Err() != nil {
		// Abort the walk if the context is already cancelled.
		return ctx.Err()
	}
	if err := fn(entry); err != nil || !entry.IsDir {
		if err == SkipDir && entry.IsDir {
			return nil // Skip this dir's contents, not an error
		}
		return err
	}
	entries, err := fs.ReadDir(ctx, entry.URI)
	if err != nil {
		return err
	}
	for _, child := range entries {
		if err := walkDir(ctx, fs, child, fn); err != nil {
			if err == SkipDir {
				break // SkipDir from a child stops the remaining siblings
			}
			return err
		}
	}
	return nil
}

// DiskFileSystem is a FileSystem implementation that reads from the local disk using the os package.
type DiskFileSystem struct{}

// NewDiskFileSystem creates a new instance of DiskFileSystem.
func NewDiskFileSystem() FileSystem {
	return &DiskFileSystem{}
}

// Exists checks if the given URI exists in the local file system.
func (s *DiskFileSystem) Exists(ctx context.Context, uri core.URI) bool {
	_, err := s.Stat(ctx, uri)
	return err == nil
}

// Stat returns the file system entry for the given URI, or an error if it does not exist.
func (s *DiskFileSystem) Stat(ctx context.Context, uri core.URI) (DirEntry, error) {
	if uri.Scheme() != core.FileScheme {
		return DirEntry{}, os.ErrNotExist
	}
	path := uri.FilePath()
	info, err := os.Lstat(path)
	if err != nil {
		return DirEntry{}, err
	}
	return DirEntry{
		IsFile:    info.Mode().IsRegular(),
		IsDir:     info.IsDir(),
		IsSymlink: info.Mode()&os.ModeSymlink != 0,
		URI:       uri,
	}, nil
}

// ReadFile reads the content of the file at the given URI and returns it as a byte slice.
func (s *DiskFileSystem) ReadFile(ctx context.Context, uri core.URI) ([]byte, error) {
	if uri.Scheme() != core.FileScheme {
		return nil, os.ErrNotExist
	}
	return os.ReadFile(uri.FilePath())
}

// ReadDir reads the contents of the directory at the given URI and returns a slice of DirEntry for its entries.
func (s *DiskFileSystem) ReadDir(ctx context.Context, uri core.URI) ([]DirEntry, error) {
	if uri.Scheme() != core.FileScheme {
		return nil, os.ErrNotExist
	}
	path := uri.FilePath()
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	result := make([]DirEntry, len(entries))
	for i, entry := range entries {
		result[i] = DirEntry{
			IsFile:    entry.Type().IsRegular(),
			IsDir:     entry.IsDir(),
			IsSymlink: entry.Type()&os.ModeSymlink != 0,
			URI:       uri.JoinPath(entry.Name()),
		}
	}
	return result, nil
}

// virtualDirEntry is a single node in the file system trie. A node is either a
// directory (children != nil) or a file (children == nil, content holds the
// bytes). The empty string key in the parent's children map is never used;
// children are keyed by path segment.
type virtualDirEntry struct {
	name     string
	content  []byte                      // non-nil => file
	children map[string]*virtualDirEntry // non-nil => directory
}

func (e *virtualDirEntry) isDir() bool { return e.children != nil }

// VirtualFileSystem is an in-memory FileSystem. Each URI scheme has its own
// directory trie rooted at an implicit root directory.
//
// Note: not safe for concurrent use; wrap with a workspace lock if shared.
type VirtualFileSystem struct {
	schemeRoots map[string]*virtualDirEntry
}

func NewVirtualFileSystem() *VirtualFileSystem {
	return &VirtualFileSystem{
		schemeRoots: make(map[string]*virtualDirEntry),
	}
}

// segments splits a URI path into non-empty path segments.
func segments(uri core.URI) []string {
	return strings.FieldsFunc(uri.Path(), func(r rune) bool { return r == '/' })
}

// find walks the trie for the given URI and returns the node, or nil if any
// segment is missing or traverses through a file.
func (s *VirtualFileSystem) find(uri core.URI) *virtualDirEntry {
	node, ok := s.schemeRoots[uri.Scheme()]
	if !ok {
		return nil
	}
	for _, seg := range segments(uri) {
		if !node.isDir() {
			return nil
		}
		next, ok := node.children[seg]
		if !ok {
			return nil
		}
		node = next
	}
	return node
}

// rootFor returns the root directory for a scheme, creating it on first use.
func (s *VirtualFileSystem) rootFor(scheme string) *virtualDirEntry {
	root, ok := s.schemeRoots[scheme]
	if !ok {
		root = &virtualDirEntry{children: map[string]*virtualDirEntry{}}
		s.schemeRoots[scheme] = root
	}
	return root
}

// mkdirAll walks segs from the scheme root, creating missing directories. It
// returns an error if a segment collides with an existing file.
func (s *VirtualFileSystem) mkdirAll(scheme string, segs []string) (*virtualDirEntry, error) {
	node := s.rootFor(scheme)
	for _, seg := range segs {
		child, ok := node.children[seg]
		if !ok {
			child = &virtualDirEntry{name: seg, children: map[string]*virtualDirEntry{}}
			node.children[seg] = child
		} else if !child.isDir() {
			return nil, os.ErrExist
		}
		node = child
	}
	return node, nil
}

func entryFor(node *virtualDirEntry, uri core.URI) DirEntry {
	return DirEntry{
		IsFile: !node.isDir(),
		IsDir:  node.isDir(),
		URI:    uri,
	}
}

func (s *VirtualFileSystem) Exists(ctx context.Context, uri core.URI) bool {
	return s.find(uri) != nil
}

func (s *VirtualFileSystem) Stat(ctx context.Context, uri core.URI) (DirEntry, error) {
	node := s.find(uri)
	if node == nil {
		return DirEntry{}, os.ErrNotExist
	}
	return entryFor(node, uri), nil
}

func (s *VirtualFileSystem) ReadFile(ctx context.Context, uri core.URI) ([]byte, error) {
	node := s.find(uri)
	if node == nil || node.isDir() {
		return nil, os.ErrNotExist
	}
	return node.content, nil
}

func (s *VirtualFileSystem) ReadDir(ctx context.Context, uri core.URI) ([]DirEntry, error) {
	node := s.find(uri)
	if node == nil || !node.isDir() {
		return nil, os.ErrNotExist
	}
	result := make([]DirEntry, 0, len(node.children))
	for name, child := range node.children {
		result = append(result, entryFor(child, uri.JoinPath(name)))
	}
	return result, nil
}

func (s *VirtualFileSystem) Mkdir(uri core.URI) error {
	segs := segments(uri)
	if len(segs) == 0 {
		return os.ErrExist // root always exists
	}
	parent, err := s.mkdirAll(uri.Scheme(), segs[:len(segs)-1])
	if err != nil {
		return err
	}
	last := segs[len(segs)-1]
	if _, ok := parent.children[last]; ok {
		return os.ErrExist
	}
	parent.children[last] = &virtualDirEntry{name: last, children: map[string]*virtualDirEntry{}}
	return nil
}

func (s *VirtualFileSystem) WriteFile(uri core.URI, content []byte) error {
	segs := segments(uri)
	if len(segs) == 0 {
		return os.ErrInvalid // cannot write the root
	}
	parent, err := s.mkdirAll(uri.Scheme(), segs[:len(segs)-1])
	if err != nil {
		return err
	}
	last := segs[len(segs)-1]
	if existing, ok := parent.children[last]; ok && existing.isDir() {
		return os.ErrInvalid // target is a directory
	}
	if content == nil {
		content = []byte{}
	}
	parent.children[last] = &virtualDirEntry{name: last, content: content}
	return nil
}

func (s *VirtualFileSystem) Delete(uri core.URI) error {
	segs := segments(uri)
	if len(segs) == 0 {
		if _, ok := s.schemeRoots[uri.Scheme()]; !ok {
			return os.ErrNotExist
		}
		delete(s.schemeRoots, uri.Scheme())
		return nil
	}
	dir, ok := s.schemeRoots[uri.Scheme()]
	if !ok {
		return os.ErrNotExist
	}
	for _, seg := range segs[:len(segs)-1] {
		child, ok := dir.children[seg]
		if !ok || !child.isDir() {
			return os.ErrNotExist
		}
		dir = child
	}
	last := segs[len(segs)-1]
	if _, ok := dir.children[last]; !ok {
		return os.ErrNotExist
	}
	delete(dir.children, last)
	return nil
}
