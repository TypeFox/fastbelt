// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"typefox.dev/fastbelt/util/service"
)

// LanguageID is the identifier of the language managed by this workspace.
// It must be set by adopters and corresponds to the language ID used in the LSP protocol.
type LanguageID string

// FileExtensions contains the file extensions to include, with leading dot
// (e.g. []string{".statemachine"}). It must be set by adopters.
type FileExtensions []string

// SetupDefaultServices sets up the default services for the workspace package.
// If any service is already set, it's not overwritten.
func SetupDefaultServices(sc *service.Container) {
	if !service.Has[DocumentManager](sc) {
		service.Put(sc, NewDefaultDocumentManager(sc))
	}
	if !service.Has[Initializer](sc) {
		service.Put(sc, NewDefaultInitializer(sc))
	}
	if !service.Has[Lock](sc) {
		service.Put(sc, NewDefaultLock())
	}
	if !service.Has[DocumentUpdater](sc) {
		service.Put(sc, NewDefaultDocumentUpdater(sc))
	}
	if !service.Has[Builder](sc) {
		service.Put(sc, NewDefaultBuilder(sc))
	}
	if !service.Has[DocumentParser](sc) {
		service.Put(sc, NewDefaultDocumentParser(sc))
	}
	if !service.Has[DocumentValidator](sc) {
		service.Put(sc, NewDefaultDocumentValidator(sc))
	}
}
