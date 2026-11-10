// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package linking

import "typefox.dev/fastbelt/util/service"

// SetupDefaultServices sets up the default services for the linking package.
// If any service is already set, it's not overwritten.
func SetupDefaultServices(sc *service.Container) {
	if !sc.Has[SymbolExporter]() {
		sc.Put(NewDefaultSymbolExporter(sc))
	}
	if !sc.Has[SymbolImporter]() {
		sc.Put(NewDefaultSymbolImporter(sc))
	}
	if !sc.Has[Linker]() {
		sc.Put(NewDefaultLinker(sc))
	}
	if !sc.Has[LocalSymbolsProvider]() {
		sc.Put(NewDefaultLocalSymbolsProvider(sc))
	}
	if !sc.Has[ReferenceDescriptionsProvider]() {
		sc.Put(NewDefaultReferenceDescriptionsProvider(sc))
	}
	if !sc.Has[ReferenceDescriber]() {
		sc.Put(NewDefaultReferenceDescriber(sc))
	}
}
