package generated

import (
	"typefox.dev/fastbelt/linking"
	"typefox.dev/fastbelt/workspace"
)

type FastbeltLinkingSrv struct {
	ScopeProvider       FastbeltScopeProvider
	Linker              FastbeltLinker
	ReferencesGenerator FastbeltReferenceGenerator
}

type FastbeltLinkingSrvCont interface {
	linking.LinkingSrvCont
	FastbeltLinking() *FastbeltLinkingSrv
}

type FastbeltLinkingSrvContBlock struct {
	fastbeltLinking FastbeltLinkingSrv
}

func (b *FastbeltLinkingSrvContBlock) FastbeltLinking() *FastbeltLinkingSrv {
	return &b.fastbeltLinking
}

type FastbeltGeneratedSrvCont interface {
	workspace.GeneratedSrvCont
	FastbeltLinkingSrvCont
}

func CreateDefaultServices(srv FastbeltGeneratedSrvCont) {
	linking := srv.FastbeltLinking()
	if linking.ScopeProvider == nil {
		linking.ScopeProvider = NewDefaultFastbeltScopeProvider(srv)
	}
	if linking.Linker == nil {
		linking.Linker = NewDefaultFastbeltLinker(srv)
	}
	if linking.ReferencesGenerator == nil {
		linking.ReferencesGenerator = NewDefaultFastbeltReferenceGenerator(srv)
	}
	generated := srv.Generated()
	if generated.Lexer == nil {
		generated.Lexer = NewLexer()
	}
	if generated.Parser == nil {
		generated.Parser = NewFastbeltParser(srv)
	}
}
