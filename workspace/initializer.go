// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"context"
	"log"
	"path/filepath"
	"slices"
	"strings"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/lsp"
)

// Initializer loads workspace files into [DocumentManager] when an LSP
// workspace folder is opened.
//
// A build is not triggered automatically after initialization. Created
// documents are left in the initial document state until [Builder.Build] runs.
//
// In a language server context, the first workspace build is triggered when a
// file is opened. The textDocument/didOpen notification reaches the server's
// document sync layer, which calls [DocumentUpdater.Update]. The updater
// gathers every document that has not yet completed all build phases,
// including files loaded during initialization.
//
// Adopters who wish to pre-build the workspace before a file is opened should
// set a broader activation event for their IDE extension and invoke
// [DocumentUpdater.Update] after the server receives the "initialized"
// notification. In other usage contexts, such as a CLI, call the
// [Initializer] and then [Builder] directly.
type Initializer interface {
	// Initialize walks folders on disk, reads files whose extension matches
	// [FileExtensions], and registers them with [DocumentManager]. Hidden
	// directories (names starting with ".") are skipped.
	Initialize(ctx context.Context, folders []lsp.WorkspaceFolder) error
}

// IncludeFilter determines whether a file system entry should be included
// in the workspace initialization process. It is used by [DefaultInitializer]
// to filter files and directories.
type IncludeFilter interface {
	// Include returns true if the given file system entry should be included.
	Include(entry DirEntry) bool
}

// DefaultIncludeFilter is the default implementation of [IncludeFilter].
// It includes files whose extension matches [FileExtensions] and skips hidden
// directories (names starting with ".").
type DefaultIncludeFilter struct {
	sc *service.Container
}

// Include returns true if the given file system entry should be included.
// See [DefaultIncludeFilter] for more info.
func (s *DefaultIncludeFilter) Include(entry DirEntry) bool {
	if entry.IsDir {
		name := filepath.Base(entry.URI.Path())
		// Don't include hidden directories
		return !strings.HasPrefix(name, ".")
	}
	extensions, err := service.Get[FileExtensions](s.sc)
	ext := filepath.Ext(entry.URI.Path())
	if err != nil {
		// No file extensions are configured
		// Only include files without an extension
		return ext == ""
	}
	// Include files whose extension matches the configured file extensions
	return slices.Contains(extensions, ext)
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

// DefaultInitializer is the default implementation of [Initializer].
// It scans the given workspace folders for files and matches them
// against the configured [IncludeFilter].
type DefaultInitializer struct {
	sc *service.Container
}

// NewDefaultInitializer returns an [Initializer] that scans workspace folders
// for files matching [FileExtensions]. Will skip hidden directories (names starting with ".").
func NewDefaultInitializer(sc *service.Container) Initializer {
	return &DefaultInitializer{sc: sc}
}

func (s *DefaultInitializer) Initialize(ctx context.Context, folders []lsp.WorkspaceFolder) error {
	fs, err := service.Get[FileSystem](s.sc)
	if err != nil {
		return nil
	}
	languageID, err := service.Get[LanguageID](s.sc)
	if err != nil {
		log.Print("workspace LanguageID is not set")
		return nil
	}
	filter, err := service.Get[IncludeFilter](s.sc)
	if err != nil {
		log.Print("workspace IncludeFilter is not set")
		return nil
	}
	for _, folder := range folders {
		root := core.ParseURI(folder.URI)
		err := WalkFileSystem(ctx, fs, root, func(entry DirEntry) error {
			included := filter.Include(entry)
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
			return err
		}
	}
	return nil
}

func (s *DefaultInitializer) loadFile(ctx context.Context, uri core.URI, fs FileSystem, languageID LanguageID) {
	content, err := fs.ReadFile(ctx, uri)
	if err != nil {
		log.Printf("failed to read file %s: %v", uri.String(), err)
		return
	}
	textDoc, err := textdoc.NewFile(uri.DocumentURI(), string(languageID), 0, string(content))
	if err != nil {
		log.Printf("failed to create text document for %s: %v", uri.String(), err)
		return
	}
	doc := core.NewDocument(textDoc)
	service.MustGet[DocumentManager](s.sc).Set(doc)
}
