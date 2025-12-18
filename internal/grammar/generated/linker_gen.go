package generated

import (
	"context"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/linking"
)

type FastbeltScopeProvider interface {
	ScopeParserRuleReturnType(ctx context.Context, reference *core.Reference[Interface]) core.Scope
}

type DefaultFastbeltScopeProvider struct {
	srv FastbeltLinkingSrvCont
}

func NewDefaultFastbeltScopeProvider(srv FastbeltLinkingSrvCont) FastbeltScopeProvider {
	return &DefaultFastbeltScopeProvider{srv: srv}
}

func (s *DefaultFastbeltScopeProvider) ScopeParserRuleReturnType(ctx context.Context, reference *core.Reference[Interface]) core.Scope {
	return linking.LocalScopeOfType[Interface](reference.Owner, s.srv.Linking().SymbolTable.LocalSymbols)
}

type FastbeltLinker interface {
	LinkParserRuleReturnType(ctx context.Context, reference *core.Reference[Interface]) (*core.AstNodeDescription, *core.ReferenceError)
}

type DefaultFastbeltLinker struct {
	srv FastbeltLinkingSrvCont
}

func NewDefaultFastbeltLinker(srv FastbeltLinkingSrvCont) FastbeltLinker {
	return &DefaultFastbeltLinker{srv: srv}
}

func (l *DefaultFastbeltLinker) LinkParserRuleReturnType(ctx context.Context, reference *core.Reference[Interface]) (*core.AstNodeDescription, *core.ReferenceError) {
	scope := l.srv.FastbeltLinking().ScopeProvider.ScopeParserRuleReturnType(ctx, reference)
	return core.DefaultLink(scope, reference.Text)
}

type FastbeltReferenceGenerator interface {
	ParserRuleReturnType(owner core.AstNode, token *core.Token) *core.Reference[Interface]
}

type DefaultFastbeltReferenceGenerator struct {
	srv FastbeltLinkingSrvCont
}

func NewDefaultFastbeltReferenceGenerator(srv FastbeltLinkingSrvCont) FastbeltReferenceGenerator {
	return &DefaultFastbeltReferenceGenerator{srv: srv}
}

func (g *DefaultFastbeltReferenceGenerator) ParserRuleReturnType(owner core.AstNode, token *core.Token) *core.Reference[Interface] {
	fn := g.srv.FastbeltLinking().Linker.LinkParserRuleReturnType
	return core.NewReference(owner, token, fn)
}
