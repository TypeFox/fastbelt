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
type DefaultInitializer struct {
	sc *service.Container
}

// NewDefaultInitializer returns an [Initializer] that scans workspace folders
// for files matching [FileExtensions].
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
	extensions, err := service.Get[FileExtensions](s.sc)
	if err != nil {
		log.Print("workspace FileExtensions is not set")
		return nil
	}
	for _, folder := range folders {
		root := core.ParseURI(folder.URI)
		err := WalkFileSystem(ctx, fs, root, func(entry DirEntry) error {
			name := filepath.Base(entry.URI.Path())
			if entry.IsDir {
				if strings.HasPrefix(name, ".") {
					return SkipDir
				}
				return nil
			}
			ext := filepath.Ext(name)
			if !slices.Contains(extensions, ext) {
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
