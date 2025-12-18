package linking

type LinkingSrv struct {
	SymbolTable  LocalSymbolTableProvider
	NameProvider NameProvider
	Linker       Linker
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
		linking.SymbolTable = NewDefaultSymbolTable(srv)
	}
	if linking.NameProvider == nil {
		linking.NameProvider = NewDefaultNameProvider()
	}
	if linking.Linker == nil {
		linking.Linker = NewDefaultLinker()
	}
}
