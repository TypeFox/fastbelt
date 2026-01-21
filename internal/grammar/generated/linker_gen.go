package generated

import (
	"context"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/linking"
)

type FastbeltScopeProvider interface {
	ScopeInterfaceExtends(ctx context.Context, reference *core.Reference[Interface]) core.Scope
	ScopeReferenceTypeType(ctx context.Context, reference *core.Reference[Interface]) core.Scope
	ScopeParserRuleReturnType(ctx context.Context, reference *core.Reference[Interface]) core.Scope
	ScopeCrossRefType(ctx context.Context, reference *core.Reference[Interface]) core.Scope
	ScopeRuleCallRule(ctx context.Context, reference *core.Reference[AbstractRule]) core.Scope
	ScopeActionType(ctx context.Context, reference *core.Reference[Interface]) core.Scope
	ScopeActionProperty(ctx context.Context, reference *core.Reference[Field]) core.Scope
}

type DefaultFastbeltScopeProvider struct {
	srv FastbeltLinkingSrvCont
}

func NewDefaultFastbeltScopeProvider(srv FastbeltLinkingSrvCont) *DefaultFastbeltScopeProvider {
	return &DefaultFastbeltScopeProvider{srv: srv}
}

func (s *DefaultFastbeltScopeProvider) ScopeInterfaceExtends(ctx context.Context, reference *core.Reference[Interface]) core.Scope {
	return linking.LocalScopeOfType[Interface](reference.Owner, s.srv.Linking().SymbolTable.LocalSymbols)
}

func (s *DefaultFastbeltScopeProvider) ScopeReferenceTypeType(ctx context.Context, reference *core.Reference[Interface]) core.Scope {
	return linking.LocalScopeOfType[Interface](reference.Owner, s.srv.Linking().SymbolTable.LocalSymbols)
}

func (s *DefaultFastbeltScopeProvider) ScopeParserRuleReturnType(ctx context.Context, reference *core.Reference[Interface]) core.Scope {
	return linking.LocalScopeOfType[Interface](reference.Owner, s.srv.Linking().SymbolTable.LocalSymbols)
}

func (s *DefaultFastbeltScopeProvider) ScopeCrossRefType(ctx context.Context, reference *core.Reference[Interface]) core.Scope {
	return linking.LocalScopeOfType[Interface](reference.Owner, s.srv.Linking().SymbolTable.LocalSymbols)
}

func (s *DefaultFastbeltScopeProvider) ScopeRuleCallRule(ctx context.Context, reference *core.Reference[AbstractRule]) core.Scope {
	return linking.LocalScopeOfType[AbstractRule](reference.Owner, s.srv.Linking().SymbolTable.LocalSymbols)
}

func (s *DefaultFastbeltScopeProvider) ScopeActionType(ctx context.Context, reference *core.Reference[Interface]) core.Scope {
	return linking.LocalScopeOfType[Interface](reference.Owner, s.srv.Linking().SymbolTable.LocalSymbols)
}

func (s *DefaultFastbeltScopeProvider) ScopeActionProperty(ctx context.Context, reference *core.Reference[Field]) core.Scope {
	return linking.LocalScopeOfType[Field](reference.Owner, s.srv.Linking().SymbolTable.LocalSymbols)
}

type FastbeltLinker interface {
	LinkInterfaceExtends(ctx context.Context, reference *core.Reference[Interface]) (*core.AstNodeDescription, *core.ReferenceError)
	LinkReferenceTypeType(ctx context.Context, reference *core.Reference[Interface]) (*core.AstNodeDescription, *core.ReferenceError)
	LinkParserRuleReturnType(ctx context.Context, reference *core.Reference[Interface]) (*core.AstNodeDescription, *core.ReferenceError)
	LinkCrossRefType(ctx context.Context, reference *core.Reference[Interface]) (*core.AstNodeDescription, *core.ReferenceError)
	LinkRuleCallRule(ctx context.Context, reference *core.Reference[AbstractRule]) (*core.AstNodeDescription, *core.ReferenceError)
	LinkActionType(ctx context.Context, reference *core.Reference[Interface]) (*core.AstNodeDescription, *core.ReferenceError)
	LinkActionProperty(ctx context.Context, reference *core.Reference[Field]) (*core.AstNodeDescription, *core.ReferenceError)
}

type DefaultFastbeltLinker struct {
	srv FastbeltLinkingSrvCont
}

func NewDefaultFastbeltLinker(srv FastbeltLinkingSrvCont) *DefaultFastbeltLinker {
	return &DefaultFastbeltLinker{srv: srv}
}

func (l *DefaultFastbeltLinker) LinkInterfaceExtends(ctx context.Context, reference *core.Reference[Interface]) (*core.AstNodeDescription, *core.ReferenceError) {
	scope := l.srv.FastbeltLinking().ScopeProvider.ScopeInterfaceExtends(ctx, reference)
	return core.DefaultLink(scope, reference.Text)
}

func (l *DefaultFastbeltLinker) LinkReferenceTypeType(ctx context.Context, reference *core.Reference[Interface]) (*core.AstNodeDescription, *core.ReferenceError) {
	scope := l.srv.FastbeltLinking().ScopeProvider.ScopeReferenceTypeType(ctx, reference)
	return core.DefaultLink(scope, reference.Text)
}

func (l *DefaultFastbeltLinker) LinkParserRuleReturnType(ctx context.Context, reference *core.Reference[Interface]) (*core.AstNodeDescription, *core.ReferenceError) {
	scope := l.srv.FastbeltLinking().ScopeProvider.ScopeParserRuleReturnType(ctx, reference)
	return core.DefaultLink(scope, reference.Text)
}

func (l *DefaultFastbeltLinker) LinkCrossRefType(ctx context.Context, reference *core.Reference[Interface]) (*core.AstNodeDescription, *core.ReferenceError) {
	scope := l.srv.FastbeltLinking().ScopeProvider.ScopeCrossRefType(ctx, reference)
	return core.DefaultLink(scope, reference.Text)
}

func (l *DefaultFastbeltLinker) LinkRuleCallRule(ctx context.Context, reference *core.Reference[AbstractRule]) (*core.AstNodeDescription, *core.ReferenceError) {
	scope := l.srv.FastbeltLinking().ScopeProvider.ScopeRuleCallRule(ctx, reference)
	return core.DefaultLink(scope, reference.Text)
}

func (l *DefaultFastbeltLinker) LinkActionType(ctx context.Context, reference *core.Reference[Interface]) (*core.AstNodeDescription, *core.ReferenceError) {
	scope := l.srv.FastbeltLinking().ScopeProvider.ScopeActionType(ctx, reference)
	return core.DefaultLink(scope, reference.Text)
}

func (l *DefaultFastbeltLinker) LinkActionProperty(ctx context.Context, reference *core.Reference[Field]) (*core.AstNodeDescription, *core.ReferenceError) {
	scope := l.srv.FastbeltLinking().ScopeProvider.ScopeActionProperty(ctx, reference)
	return core.DefaultLink(scope, reference.Text)
}

type FastbeltReferenceGenerator interface {
	InterfaceExtends(owner core.AstNode, token *core.Token) *core.Reference[Interface]
	ReferenceTypeType(owner core.AstNode, token *core.Token) *core.Reference[Interface]
	ParserRuleReturnType(owner core.AstNode, token *core.Token) *core.Reference[Interface]
	CrossRefType(owner core.AstNode, token *core.Token) *core.Reference[Interface]
	RuleCallRule(owner core.AstNode, token *core.Token) *core.Reference[AbstractRule]
	ActionType(owner core.AstNode, token *core.Token) *core.Reference[Interface]
	ActionProperty(owner core.AstNode, token *core.Token) *core.Reference[Field]
}

type DefaultFastbeltReferenceGenerator struct {
	srv FastbeltLinkingSrvCont
}

func NewDefaultFastbeltReferenceGenerator(srv FastbeltLinkingSrvCont) *DefaultFastbeltReferenceGenerator {
	return &DefaultFastbeltReferenceGenerator{srv: srv}
}

func (g *DefaultFastbeltReferenceGenerator) InterfaceExtends(owner core.AstNode, token *core.Token) *core.Reference[Interface] {
	fn := g.srv.FastbeltLinking().Linker.LinkInterfaceExtends
	return core.NewReference(owner, token, fn)
}

func (g *DefaultFastbeltReferenceGenerator) ReferenceTypeType(owner core.AstNode, token *core.Token) *core.Reference[Interface] {
	fn := g.srv.FastbeltLinking().Linker.LinkReferenceTypeType
	return core.NewReference(owner, token, fn)
}

func (g *DefaultFastbeltReferenceGenerator) ParserRuleReturnType(owner core.AstNode, token *core.Token) *core.Reference[Interface] {
	fn := g.srv.FastbeltLinking().Linker.LinkParserRuleReturnType
	return core.NewReference(owner, token, fn)
}

func (g *DefaultFastbeltReferenceGenerator) CrossRefType(owner core.AstNode, token *core.Token) *core.Reference[Interface] {
	fn := g.srv.FastbeltLinking().Linker.LinkCrossRefType
	return core.NewReference(owner, token, fn)
}

func (g *DefaultFastbeltReferenceGenerator) RuleCallRule(owner core.AstNode, token *core.Token) *core.Reference[AbstractRule] {
	fn := g.srv.FastbeltLinking().Linker.LinkRuleCallRule
	return core.NewReference(owner, token, fn)
}

func (g *DefaultFastbeltReferenceGenerator) ActionType(owner core.AstNode, token *core.Token) *core.Reference[Interface] {
	fn := g.srv.FastbeltLinking().Linker.LinkActionType
	return core.NewReference(owner, token, fn)
}

func (g *DefaultFastbeltReferenceGenerator) ActionProperty(owner core.AstNode, token *core.Token) *core.Reference[Field] {
	fn := g.srv.FastbeltLinking().Linker.LinkActionProperty
	return core.NewReference(owner, token, fn)
}
