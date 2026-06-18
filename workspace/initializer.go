// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"context"
	"io/fs"
	"log"
	"os"
	"path/filepath"
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
	if !service.Has[LanguageID](s.sc) {
		log.Print("workspace LanguageID is not set")
	}
	if !service.Has[FileExtensions](s.sc) {
		log.Print("workspace FileExtensions is not set")
		return nil
	}
	for _, folder := range folders {
		root := core.ParseURI(folder.URI).Path()
		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if err := ctx.Err(); err != nil {
				return err
			}
			name := d.Name()
			if d.IsDir() {
				if strings.HasPrefix(name, ".") {
					return filepath.SkipDir
				}
				return nil
			}
			ext := filepath.Ext(name)
			if !s.matchesExtension(ext) {
				return nil
			}
			s.loadFile(path)
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *DefaultInitializer) loadFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Printf("failed to read file %s: %v", path, err)
		return
	}
	uri := core.FileURI(path)
	languageID := service.MustGet[LanguageID](s.sc)
	textDoc, err := textdoc.NewFile(uri.DocumentURI(), string(languageID), 0, string(content))
	if err != nil {
		log.Printf("failed to create text document for %s: %v", path, err)
		return
	}
	doc := core.NewDocument(textDoc)
	service.MustGet[DocumentManager](s.sc).Set(doc)
}

func (s *DefaultInitializer) matchesExtension(ext string) bool {
	extensions := service.MustGet[FileExtensions](s.sc)
	for _, e := range extensions {
		if e == ext {
			return true
		}
	}
	return false
}
