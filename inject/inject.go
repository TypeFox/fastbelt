// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package inject

// Services is a container for all services used for dependency injection.
// This is a fixed struct with known keys and known service types.
// Services are added as fields with their concrete types.
// To avoid import cycles, services are defined using interface{} and cast in usage sites.
//
// Library consumers can extend Services using struct embedding:
//
//	type MyAppServices struct {
//	    *inject.Services
//	    Database   *DatabaseService
//	    Cache      *CacheService
//	}
//
// This allows consumers to add their own services while maintaining access to base services.
type Services struct {
	// LSP services - using interface{} to avoid import cycles
	// Actual types are defined in the lsp package:
	// LanguageServerHandlers: *lsp.LanguageServerHandlers
	// LanguageServer: lsp.LanguageServer
	// ConnectionBinder: lsp.ConnectionBinder
	// ConnectionDialer: lsp.ConnectionDialer
	LanguageServerHandlers interface{}
	LanguageServer         interface{}
	ConnectionBinder       interface{}
	ConnectionDialer       interface{}
}

// NewServices creates a new Services container.
func NewServices() *Services {
	return &Services{}
}
