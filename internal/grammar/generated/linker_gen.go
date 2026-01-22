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

type FastbeltReferenceLinker interface {
	LinkInterfaceExtends(ctx context.Context, reference *core.Reference[Interface]) (*core.AstNodeDescription, *core.ReferenceError)
	LinkReferenceTypeType(ctx context.Context, reference *core.Reference[Interface]) (*core.AstNodeDescription, *core.ReferenceError)
	LinkParserRuleReturnType(ctx context.Context, reference *core.Reference[Interface]) (*core.AstNodeDescription, *core.ReferenceError)
	LinkCrossRefType(ctx context.Context, reference *core.Reference[Interface]) (*core.AstNodeDescription, *core.ReferenceError)
	LinkRuleCallRule(ctx context.Context, reference *core.Reference[AbstractRule]) (*core.AstNodeDescription, *core.ReferenceError)
	LinkActionType(ctx context.Context, reference *core.Reference[Interface]) (*core.AstNodeDescription, *core.ReferenceError)
	LinkActionProperty(ctx context.Context, reference *core.Reference[Field]) (*core.AstNodeDescription, *core.ReferenceError)
}

type DefaultFastbeltReferenceLinker struct {
	srv FastbeltLinkingSrvCont
}

func NewDefaultFastbeltReferenceLinker(srv FastbeltLinkingSrvCont) *DefaultFastbeltReferenceLinker {
	return &DefaultFastbeltReferenceLinker{srv: srv}
}

func (l *DefaultFastbeltReferenceLinker) LinkInterfaceExtends(ctx context.Context, reference *core.Reference[Interface]) (*core.AstNodeDescription, *core.ReferenceError) {
	scope := l.srv.FastbeltLinking().ScopeProvider.ScopeInterfaceExtends(ctx, reference)
	return core.DefaultLink(scope, reference.Text)
}

func (l *DefaultFastbeltReferenceLinker) LinkReferenceTypeType(ctx context.Context, reference *core.Reference[Interface]) (*core.AstNodeDescription, *core.ReferenceError) {
	scope := l.srv.FastbeltLinking().ScopeProvider.ScopeReferenceTypeType(ctx, reference)
	return core.DefaultLink(scope, reference.Text)
}

func (l *DefaultFastbeltReferenceLinker) LinkParserRuleReturnType(ctx context.Context, reference *core.Reference[Interface]) (*core.AstNodeDescription, *core.ReferenceError) {
	scope := l.srv.FastbeltLinking().ScopeProvider.ScopeParserRuleReturnType(ctx, reference)
	return core.DefaultLink(scope, reference.Text)
}

func (l *DefaultFastbeltReferenceLinker) LinkCrossRefType(ctx context.Context, reference *core.Reference[Interface]) (*core.AstNodeDescription, *core.ReferenceError) {
	scope := l.srv.FastbeltLinking().ScopeProvider.ScopeCrossRefType(ctx, reference)
	return core.DefaultLink(scope, reference.Text)
}

func (l *DefaultFastbeltReferenceLinker) LinkRuleCallRule(ctx context.Context, reference *core.Reference[AbstractRule]) (*core.AstNodeDescription, *core.ReferenceError) {
	scope := l.srv.FastbeltLinking().ScopeProvider.ScopeRuleCallRule(ctx, reference)
	return core.DefaultLink(scope, reference.Text)
}

func (l *DefaultFastbeltReferenceLinker) LinkActionType(ctx context.Context, reference *core.Reference[Interface]) (*core.AstNodeDescription, *core.ReferenceError) {
	scope := l.srv.FastbeltLinking().ScopeProvider.ScopeActionType(ctx, reference)
	return core.DefaultLink(scope, reference.Text)
}

func (l *DefaultFastbeltReferenceLinker) LinkActionProperty(ctx context.Context, reference *core.Reference[Field]) (*core.AstNodeDescription, *core.ReferenceError) {
	scope := l.srv.FastbeltLinking().ScopeProvider.ScopeActionProperty(ctx, reference)
	return core.DefaultLink(scope, reference.Text)
}

type FastbeltReferencesConstructor interface {
	InterfaceExtends(owner core.AstNode, token *core.Token) *core.Reference[Interface]
	ReferenceTypeType(owner core.AstNode, token *core.Token) *core.Reference[Interface]
	ParserRuleReturnType(owner core.AstNode, token *core.Token) *core.Reference[Interface]
	CrossRefType(owner core.AstNode, token *core.Token) *core.Reference[Interface]
	RuleCallRule(owner core.AstNode, token *core.Token) *core.Reference[AbstractRule]
	ActionType(owner core.AstNode, token *core.Token) *core.Reference[Interface]
	ActionProperty(owner core.AstNode, token *core.Token) *core.Reference[Field]
}

type DefaultFastbeltReferencesConstructor struct {
	srv FastbeltLinkingSrvCont
}

func NewDefaultFastbeltReferencesConstructor(srv FastbeltLinkingSrvCont) *DefaultFastbeltReferencesConstructor {
	return &DefaultFastbeltReferencesConstructor{srv: srv}
}

func (g *DefaultFastbeltReferencesConstructor) InterfaceExtends(owner core.AstNode, token *core.Token) *core.Reference[Interface] {
	fn := g.srv.FastbeltLinking().ReferenceLinker.LinkInterfaceExtends
	return core.NewReference(owner, token, fn)
}

func (g *DefaultFastbeltReferencesConstructor) ReferenceTypeType(owner core.AstNode, token *core.Token) *core.Reference[Interface] {
	fn := g.srv.FastbeltLinking().ReferenceLinker.LinkReferenceTypeType
	return core.NewReference(owner, token, fn)
}

func (g *DefaultFastbeltReferencesConstructor) ParserRuleReturnType(owner core.AstNode, token *core.Token) *core.Reference[Interface] {
	fn := g.srv.FastbeltLinking().ReferenceLinker.LinkParserRuleReturnType
	return core.NewReference(owner, token, fn)
}

func (g *DefaultFastbeltReferencesConstructor) CrossRefType(owner core.AstNode, token *core.Token) *core.Reference[Interface] {
	fn := g.srv.FastbeltLinking().ReferenceLinker.LinkCrossRefType
	return core.NewReference(owner, token, fn)
}

func (g *DefaultFastbeltReferencesConstructor) RuleCallRule(owner core.AstNode, token *core.Token) *core.Reference[AbstractRule] {
	fn := g.srv.FastbeltLinking().ReferenceLinker.LinkRuleCallRule
	return core.NewReference(owner, token, fn)
}

func (g *DefaultFastbeltReferencesConstructor) ActionType(owner core.AstNode, token *core.Token) *core.Reference[Interface] {
	fn := g.srv.FastbeltLinking().ReferenceLinker.LinkActionType
	return core.NewReference(owner, token, fn)
}

func (g *DefaultFastbeltReferencesConstructor) ActionProperty(owner core.AstNode, token *core.Token) *core.Reference[Field] {
	fn := g.srv.FastbeltLinking().ReferenceLinker.LinkActionProperty
	return core.NewReference(owner, token, fn)
}
