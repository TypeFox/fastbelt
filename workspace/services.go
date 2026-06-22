// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"typefox.dev/fastbelt/util/service"
)

// LanguageID is the LSP language identifier for documents in this workspace.
// Adopters must register a value with [service.Put] before sealing the
// container; [DefaultInitializer] and the LSP document sync layer read it
// when creating [core.Document] values.
type LanguageID string

// FileExtensions lists filename suffixes that belong to this language, each
// with a leading dot (for example []string{".statemachine"}). Adopters must
// register a value with [service.Put]; [DefaultInitializer] uses it to
// decide which files to load when a workspace folder is opened.
type FileExtensions []string

// SetupDefaultServices registers default workspace services in sc.
// It is idempotent: types that are already registered are left unchanged.
// Call it from a language's SetupServices after registering [LanguageID] and
// [FileExtensions], and before [service.Container.Seal].
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
	if !service.Has[DocumentChangeImpact](sc) {
		service.Put(sc, NewDefaultDocumentChangeImpact(sc))
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
