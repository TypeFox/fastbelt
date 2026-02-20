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
	"typefox.dev/lsp"
)

// Initializer traverses workspace folders, finds language files, and registers
// them in the DocumentManager.
type Initializer interface {
	// Initialize walks the given workspace folders, reads files matching the
	// configured extensions, and creates documents for them.
	Initialize(ctx context.Context, folders []lsp.WorkspaceFolder) error
}

// DefaultInitializer is the default implementation of Initializer.
type DefaultInitializer struct {
	srv WorkspaceSrvCont
}

// NewDefaultInitializer creates a new default initializer.
func NewDefaultInitializer(srv WorkspaceSrvCont) *DefaultInitializer {
	return &DefaultInitializer{srv: srv}
}

// Initialize walks each workspace folder, reads files whose extension matches
// FileExtensions, and registers the resulting documents in the DocumentManager.
func (i *DefaultInitializer) Initialize(ctx context.Context, folders []lsp.WorkspaceFolder) error {
	ws := i.srv.Workspace()
	if ws.LanguageID == "" {
		log.Print("workspace LanguageID is not set")
	}
	if len(ws.FileExtensions) == 0 {
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
			if !i.matchesExtension(ext) {
				return nil
			}
			i.loadFile(path)
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (i *DefaultInitializer) loadFile(path string) {
	content, err := os.ReadFile(path)
	if err != nil {
		log.Printf("failed to read file %s: %v", path, err)
		return
	}
	uri := core.FileURI(path)
	textDoc, err := textdoc.NewFile(uri.DocumentURI(), i.srv.Workspace().LanguageID, 0, string(content))
	if err != nil {
		log.Printf("failed to create text document for %s: %v", path, err)
		return
	}
	doc := core.NewDocument(textDoc)
	i.srv.Workspace().DocumentManager.Set(doc)
	// TODO parse the document and collect exported symbols
	// We don't need that right now because all documents are rebuilt on every change
}

func (i *DefaultInitializer) matchesExtension(ext string) bool {
	for _, e := range i.srv.Workspace().FileExtensions {
		if e == ext {
			return true
		}
	}
	return false
}
