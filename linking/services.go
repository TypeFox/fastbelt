// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package linking

import "typefox.dev/fastbelt/util/service"

// SetupDefaultServices sets up the default services for the linking package.
// If any service is already set, it's not overwritten.
func SetupDefaultServices(sc *service.Container) {
	if !service.Has[SymbolExporter](sc) {
		service.Put(sc, NewDefaultSymbolExporter(sc))
	}
	if !service.Has[SymbolImporter](sc) {
		service.Put(sc, NewDefaultSymbolImporter(sc))
	}
	if !service.Has[Linker](sc) {
		service.Put(sc, NewDefaultLinker(sc))
	}
	if !service.Has[LocalSymbolsProvider](sc) {
		service.Put(sc, NewDefaultLocalSymbolsProvider(sc))
	}
	if !service.Has[ReferenceDescriptionsProvider](sc) {
		service.Put(sc, NewDefaultReferenceDescriptionsProvider(sc))
	}
	if !service.Has[ReferenceDescriber](sc) {
		service.Put(sc, NewDefaultReferenceDescriber(sc))
	}
}
