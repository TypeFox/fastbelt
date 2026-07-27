// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"typefox.dev/fastbelt/util/service"
)

// SetupDefaultServices registers default workspace services in sc.
// It is idempotent: types that are already registered are left unchanged.
// Call it from a language's SetupServices, before [service.Container.Seal].
func SetupDefaultServices(sc *service.Container) {
	if !service.Has[DocumentManager](sc) {
		service.Put(sc, NewDefaultDocumentManager(sc))
	}
	if !service.Has[FileSystem](sc) {
		service.Put(sc, NewDiskFileSystem())
	}
	if !service.Has[Initializer](sc) {
		service.Put(sc, NewDefaultInitializer(sc))
	}
	if !service.Has[IncludeFilter](sc) {
		service.Put(sc, NewDefaultIncludeFilter(sc))
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
	if !service.Has[DocumentValidator](sc) {
		service.Put(sc, NewDefaultDocumentValidator(sc))
	}
}
