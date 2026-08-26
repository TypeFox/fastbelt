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
	if !sc.Has[DocumentManager]() {
		sc.Put(NewDefaultDocumentManager(sc))
	}
	if !sc.Has[FileSystem]() {
		sc.Put(NewDiskFileSystem())
	}
	if !sc.Has[Initializer]() {
		sc.Put(NewDefaultInitializer(sc))
	}
	if !sc.Has[IncludeFilter]() {
		sc.Put(NewDefaultIncludeFilter(sc))
	}
	if !sc.Has[Lock]() {
		sc.Put(NewDefaultLock(sc))
	}
	if !sc.Has[DocumentUpdater]() {
		sc.Put(NewDefaultDocumentUpdater(sc))
	}
	if !sc.Has[DocumentChangeImpact]() {
		sc.Put(NewDefaultDocumentChangeImpact(sc))
	}
	if !sc.Has[Builder]() {
		sc.Put(NewDefaultBuilder(sc))
	}
	if !sc.Has[DocumentValidator]() {
		sc.Put(NewDefaultDocumentValidator(sc))
	}
}
