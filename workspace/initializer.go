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

// Initializer initializes the documents of a workspace.
type Initializer interface {
	// Initialize walks the given workspace folders, reads files matching the
	// configured extensions, and creates documents for them.
	Initialize(ctx context.Context, folders []lsp.WorkspaceFolder) error
}

// DefaultInitializer is the default implementation of [Initializer].
type DefaultInitializer struct {
	sc *service.Container
}

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
	// TODO parse the document and collect exported symbols
	// We don't need that right now because all documents are rebuilt on every change
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
