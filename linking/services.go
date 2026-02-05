// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package linking

type LinkingSrv struct {
	LocalSymbolTableProvider     LocalSymbolTableProvider
	LocalSymbolTableItemProvider LocalSymbolTableItemProvider
	Namer                        Namer
	Linker                       Linker
}

type LinkingSrvCont interface {
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
	if linking.LocalSymbolTableProvider == nil {
		linking.LocalSymbolTableProvider = NewDefaultLocalSymbolTableProvider(srv)
	}
	if linking.LocalSymbolTableItemProvider == nil {
		linking.LocalSymbolTableItemProvider = NewDefaultLocalSymbolTableItemProvider(srv)
	}
	if linking.Namer == nil {
		linking.Namer = NewDefaultNamer()
	}
	if linking.Linker == nil {
		linking.Linker = NewDefaultLinker()
	}
}
