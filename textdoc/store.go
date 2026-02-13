// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package textdoc

import (
	"sync"

	"typefox.dev/lsp"
)

// Store is responsible for storing cached file contents and overlays.
// It provides thread-safe access to document files and overlays.
type Store interface {
	// Get retrieves a handle by URI, checking overlays first, then files.
	// Returns nil if neither an overlay nor a file is stored.
	Get(uri lsp.DocumentURI) Handle
	// GetOverlay retrieves an overlay by URI. Returns nil if the overlay is not stored.
	GetOverlay(uri lsp.DocumentURI) *Overlay
	// GetFile retrieves a file by URI. Returns nil if the file is not stored.
	GetFile(uri lsp.DocumentURI) *File
	// AddOverlay adds or updates an overlay in the store.
	AddOverlay(overlay *Overlay)
	// AddFile adds or updates a file in the store.
	AddFile(file *File)
	// RemoveOverlay removes an overlay from the store by URI.
	RemoveOverlay(uri lsp.DocumentURI)
	// RemoveFile removes a file from the store by URI.
	RemoveFile(uri lsp.DocumentURI)
	// AllOverlays returns all stored overlays.
	AllOverlays() []*Overlay
	// AllFiles returns all stored files.
	AllFiles() []*File
	// KeysOverlays returns the URIs of all stored overlays.
	KeysOverlays() []lsp.DocumentURI
	// KeysFiles returns the URIs of all stored files.
	KeysFiles() []lsp.DocumentURI
}

// DefaultStore is an in-memory implementation of Store that manages files and overlays.
type DefaultStore struct {
	mu       sync.RWMutex
	overlays map[lsp.DocumentURI]*Overlay
	files    map[lsp.DocumentURI]*File
}

// NewDefaultStore creates a new store for files and overlays.
func NewDefaultStore() Store {
	return &DefaultStore{
		overlays: make(map[lsp.DocumentURI]*Overlay),
		files:    make(map[lsp.DocumentURI]*File),
	}
}

// Get retrieves a handle by URI, checking overlays first, then files.
// Returns nil if neither an overlay nor a file is stored.
func (s *DefaultStore) Get(uri lsp.DocumentURI) Handle {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if overlay, ok := s.overlays[uri]; ok {
		return overlay
	}
	if file, ok := s.files[uri]; ok {
		return file
	}
	return nil
}

// GetOverlay retrieves an overlay by URI. Returns nil if the overlay is not stored.
func (s *DefaultStore) GetOverlay(uri lsp.DocumentURI) *Overlay {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.overlays[uri]
}

// GetFile retrieves a file by URI. Returns nil if the file is not stored.
func (s *DefaultStore) GetFile(uri lsp.DocumentURI) *File {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.files[uri]
}

// AddOverlay adds or updates an overlay in the store.
func (s *DefaultStore) AddOverlay(overlay *Overlay) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.overlays[overlay.URI()] = overlay
}

// AddFile adds or updates a file in the store.
func (s *DefaultStore) AddFile(file *File) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.files[file.URI()] = file
}

// RemoveOverlay removes an overlay from the store by URI.
func (s *DefaultStore) RemoveOverlay(uri lsp.DocumentURI) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.overlays, uri)
}

// RemoveFile removes a file from the store by URI.
func (s *DefaultStore) RemoveFile(uri lsp.DocumentURI) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.files, uri)
}

// AllOverlays returns all stored overlays.
func (s *DefaultStore) AllOverlays() []*Overlay {
	s.mu.RLock()
	defer s.mu.RUnlock()

	docs := make([]*Overlay, 0, len(s.overlays))
	for _, doc := range s.overlays {
		docs = append(docs, doc)
	}
	return docs
}

// AllFiles returns all stored files.
func (s *DefaultStore) AllFiles() []*File {
	s.mu.RLock()
	defer s.mu.RUnlock()

	files := make([]*File, 0, len(s.files))
	for _, file := range s.files {
		files = append(files, file)
	}
	return files
}

// KeysOverlays returns the URIs of all stored overlays.
func (s *DefaultStore) KeysOverlays() []lsp.DocumentURI {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]lsp.DocumentURI, 0, len(s.overlays))
	for uri := range s.overlays {
		keys = append(keys, uri)
	}
	return keys
}

// KeysFiles returns the URIs of all stored files.
func (s *DefaultStore) KeysFiles() []lsp.DocumentURI {
	s.mu.RLock()
	defer s.mu.RUnlock()

	keys := make([]lsp.DocumentURI, 0, len(s.files))
	for uri := range s.files {
		keys = append(keys, uri)
	}
	return keys
}
