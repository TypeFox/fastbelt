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
// Once all documents are collected, the Initialize function should start
// a workspace build using those documents.
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
	builder, err := service.Get[Builder](s.sc)
	if err != nil {
		log.Print("workspace builder is not set")
		return nil
	}
	docs := slices.Collect(service.MustGet[DocumentManager](s.sc).All())
	builder.Build(ctx, docs)
	return nil
}

func (s *DefaultInitializer) loadFile(ctx context.Context, uri core.URI, fs FileSystem, languageID LanguageID) {
	docs := service.MustGet[DocumentManager](s.sc)
	if docs.Has(uri) {
		// Document already loaded, maybe through overlapping workspaces
		return
	}
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
	docs.Set(doc)
}
