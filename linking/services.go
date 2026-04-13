// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package linking

import "typefox.dev/fastbelt/generated"

type LinkingSrv struct {
	ExportedSymbolsProvider       ExportedSymbolsProvider
	ExportedSymbolDescriber       ExportedSymbolDescriber
	ImportedSymbolsProvider       ImportedSymbolsProvider
	LocalSymbolsProvider          LocalSymbolsProvider
	LocalSymbolDescriber          LocalSymbolDescriber
	ReferenceDescriptionsProvider ReferenceDescriptionsProvider
	ReferenceDescriber            ReferenceDescriber
	Namer                         Namer
	Linker                        Linker
}

type LinkingSrvCont interface {
	generated.GeneratedSrvCont
	Linking() *LinkingSrv
}

type LinkingSrvContBlock struct {
	linking LinkingSrv
}

func (b *LinkingSrvContBlock) Linking() *LinkingSrv {
	return &b.linking
}

func CreateDefaultServices(srv LinkingSrvCont) {
	linking := srv.Linking()
	if linking.ExportedSymbolsProvider == nil {
		linking.ExportedSymbolsProvider = NewDefaultExportedSymbolsProvider(srv)
	}
	if linking.ExportedSymbolDescriber == nil {
		linking.ExportedSymbolDescriber = NewDefaultExportedSymbolDescriber(srv)
	}
	if linking.ImportedSymbolsProvider == nil {
		linking.ImportedSymbolsProvider = NewDefaultImportedSymbolsProvider(srv)
	}
	if linking.LocalSymbolsProvider == nil {
		linking.LocalSymbolsProvider = NewDefaultLocalSymbolsProvider(srv)
	}
	if linking.LocalSymbolDescriber == nil {
		linking.LocalSymbolDescriber = NewDefaultLocalSymbolDescriber(srv)
	}
	if linking.ReferenceDescriptionsProvider == nil {
		linking.ReferenceDescriptionsProvider = NewDefaultReferenceDescriptionsProvider(srv)
	}
	if linking.ReferenceDescriber == nil {
		linking.ReferenceDescriber = NewDefaultReferenceDescriber()
	}
	if linking.Namer == nil {
		linking.Namer = NewDefaultNamer()
	}
	if linking.Linker == nil {
		linking.Linker = NewDefaultLinker()
	}
}
