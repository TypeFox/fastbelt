package linking

type LinkingSrv struct {
	SymbolTable LocalSymbolTableProvider
	Namer       Namer
	Linker      Linker
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
	if linking.SymbolTable == nil {
		linking.SymbolTable = NewDefaultLocalSymbolTableProvider(srv)
	}
	if linking.Namer == nil {
		linking.Namer = NewDefaultNamer()
	}
	if linking.Linker == nil {
		linking.Linker = NewDefaultLinker()
	}
}
