// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"context"
	"log"
	"path/filepath"
	"strings"
	"sync"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/lsp"
)

// Initializer loads workspace files into [DocumentManager] when an LSP
// workspace folder is opened.
type Initializer interface {
	// Initialize walks folders on disk, reads files whose extension matches
	// [FileExtensions], and registers them with [DocumentManager]. Should be called
	// under [Lock.Write] to prevent concurrent write access to the documents.
	Initialize(ctx context.Context, folders []lsp.WorkspaceFolder) error
	// Ready returns a channel that is closed once all workspace files have been
	// loaded into the [DocumentManager]. The initial build may still be running
	// at that point: readers synchronize with it through [Lock], which the
	// initializer holds for loading the workspace. Code that may run before
	// the initialized notification arrives should wait on Ready first.
	Ready() <-chan struct{}
}

// IncludeFilter determines whether a file system entry should be included
// in the workspace initialization process. It is used by [DefaultInitializer]
// to filter files and directories.
type IncludeFilter interface {
	// Include returns true if the given file system entry should be included.
	// For files, the second return value is the language ID of the document.
	Include(entry DirEntry) (bool, string)
}

// DefaultIncludeFilter is the default implementation of [IncludeFilter].
// It includes files that the LanguageSelector maps to a language and skips
// hidden directories (names starting with ".").
type DefaultIncludeFilter struct {
	sc       *service.Container
	warnOnce sync.Once
}

// Include returns true if the given file system entry should be included.
// See [DefaultIncludeFilter] for more info.
func (s *DefaultIncludeFilter) Include(entry DirEntry) (bool, string) {
	if entry.IsDir {
		name := filepath.Base(entry.URI.Path())
		// Don't include hidden directories
		return !strings.HasPrefix(name, "."), ""
	}
	// Note: The file isn't managed yet, so the selector only matches by URI.
	if selector, err := s.sc.Get[core.LanguageSelector](); err == nil {
		index, languageId := selector.Select(entry.URI)
		if index >= 0 {
			return true, languageId
		}
	} else {
		s.warnOnce.Do(func() { log.Print("workspace LanguageSelector is not set, no files are loaded") })
	}
	return false, ""
}

// NewDefaultIncludeFilter returns a new instance of [DefaultIncludeFilter].
func NewDefaultIncludeFilter(sc *service.Container) IncludeFilter {
	return &DefaultIncludeFilter{sc: sc}
}

// NoopInitializer is an [Initializer] that does nothing.
// It can be used in environments where there is no file system to scan, such as the browser.
type NoopInitializer struct{}

// NewNoopInitializer creates a new instance of [NoopInitializer].
func NewNoopInitializer() Initializer {
	return NoopInitializer{}
}

// Initialize does nothing and returns nil.
func (NoopInitializer) Initialize(ctx context.Context, folders []lsp.WorkspaceFolder) error {
	return nil
}

var closedCh = func() chan struct{} { ch := make(chan struct{}); close(ch); return ch }()

// Ready returns an already closed channel: there is nothing to wait for.
func (NoopInitializer) Ready() <-chan struct{} {
	return closedCh
}

// DefaultInitializer is the default implementation of [Initializer].
// It scans the given workspace folders for files and matches them
// against the configured [IncludeFilter].
type DefaultInitializer struct {
	sc        *service.Container
	ready     chan struct{}
	readyOnce sync.Once
}

// NewDefaultInitializer returns an [Initializer] that scans workspace folders
// for files accepted by the registered [IncludeFilter].
func NewDefaultInitializer(sc *service.Container) Initializer {
	return &DefaultInitializer{sc: sc, ready: make(chan struct{})}
}

// Ready returns a channel that is closed once all workspace files are loaded.
func (s *DefaultInitializer) Ready() <-chan struct{} {
	return s.ready
}

// Marks the initializer as done, closing the Ready channel.
func (s *DefaultInitializer) readyUp() {
	s.readyOnce.Do(func() { close(s.ready) })
}

func (s *DefaultInitializer) Initialize(ctx context.Context, folders []lsp.WorkspaceFolder) error {
	// The initializer is always done once it exits this method
	defer s.readyUp()
	fs, err := s.sc.Get[FileSystem]()
	if err != nil {
		return nil
	}
	filter, err := s.sc.Get[IncludeFilter]()
	if err != nil {
		log.Print("workspace IncludeFilter is not set")
		return nil
	}
	for _, folder := range folders {
		root := core.ParseURI(folder.URI)
		err := WalkFileSystem(ctx, fs, root, func(entry DirEntry) error {
			included, languageID := filter.Include(entry)
			if entry.IsDir {
				if !included {
					// Skip this directory and its contents
					return SkipDir
				} else {
					return nil
				}
			} else if !included {
				// Just skip this file
				return nil
			}
			s.loadFile(ctx, entry.URI, fs, languageID)
			return nil
		})
		if err != nil {
			// Print the issue, do not return here, simply walk the rest of the folders
			log.Printf("Initializer encountered an error while walking %s: %v", root.StringUnencoded(), err)
		}
	}
	return nil
}

func (s *DefaultInitializer) loadFile(ctx context.Context, uri core.URI, fs FileSystem, languageID string) {
	docs := s.sc.MustGet[DocumentManager]()
	if docs.Has(uri) {
		// Document already loaded, maybe through overlapping workspaces
		return
	}
	content, err := fs.ReadFile(ctx, uri)
	if err != nil {
		log.Printf("failed to read file %s: %v", uri.String(), err)
		return
	}
	textDoc := textdoc.NewFile(uri.DocumentURI(), languageID, 0, string(content))
	doc := core.NewDocument(textDoc)
	docs.Set(doc)
}
