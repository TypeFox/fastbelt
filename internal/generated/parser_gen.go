package generated

import (
	"github.com/TypeFox/langium-to-go/core"
	"github.com/TypeFox/langium-to-go/parser"
)

type Parser struct {
	state *parser.ParserState
}

func (p *Parser) Parse(tokens []*core.Token) Grammar {
	p.state = parser.NewParserState(tokens)
	result := p.ParseGrammar()
	core.AssignContainers(result)
	return result
}

func NewParser() *Parser {
	return &Parser{}
}

const (
	ActionDot_0 = iota + 1
	ActionLeftBrace_0
	ActionOperatorAssignmentOperator_0
	ActionPropertyID_0
	ActionRightBrace_0
	ActionTypeID_0
	Actioncurrent_0
	AlternativesPipe_0
	AssignmentOperatorAssignmentOperator_0
	AssignmentPropertyID_0
	CrossRefColon_0
	CrossRefLeftBracket_0
	CrossRefRightBracket_0
	CrossRefTypeID_0
	ElementCardinalityCardinality_0
	ElementLeftParen_0
	ElementRightParen_0
	FieldArrayLeftBracket_0
	FieldNameID_0
	FieldRightBracket_0
	FieldTypeID_0
	GrammarNameID_0
	GrammarSemicolon_0
	Grammargrammar_0
	InterfaceComma_0
	InterfaceExtendsID_0
	InterfaceExtendsID_1
	InterfaceLeftBrace_0
	InterfaceNameID_0
	InterfaceRightBrace_0
	Interfaceextends_0
	Interfaceinterface_0
	KeywordValueString_0
	ParserRuleColon_0
	ParserRuleNameID_0
	ParserRuleReturnTypeID_0
	ParserRuleSemicolon_0
	ParserRulereturns_0
	RuleCallRuleID_0
	TokenColon_0
	TokenNameID_0
	TokenRegexpRegexLiteral_0
	TokenSemicolon_0
	TokenTypehidden_0
	Tokentoken_0
)

var ActionLookahead12 = parser.LLkLookahead{
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_Dot_Idx},
	},
}

var AlternativesLookahead6 = parser.LLkLookahead{
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_Pipe_Idx},
	},
}

var AlternativesLookahead7 = parser.LLkLookahead{
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_Pipe_Idx},
	},
}

var AssignableLookaheadOr2 = parser.LLkLookahead{
	parser.LookaheadOption{
		parser.LookaheadPath{Token_String_Idx},
	},
	parser.LookaheadOption{
		parser.LookaheadPath{Token_ID_Idx},
	},
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_LeftBracket_Idx},
	},
}

var CrossRefLookahead11 = parser.LLkLookahead{
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_Colon_Idx},
	},
}

var ElementLookahead10 = parser.LLkLookahead{
	parser.LookaheadOption{
		parser.LookaheadPath{Token_Cardinality_Idx},
	},
}

var ElementLookaheadOr1 = parser.LLkLookahead{
	parser.LookaheadOption{
		parser.LookaheadPath{Token_String_Idx},
	},
	parser.LookaheadOption{
		parser.LookaheadPath{Token_ID_Idx, Token_AssignmentOperator_Idx},
	},
	parser.LookaheadOption{
		parser.LookaheadPath{Token_ID_Idx},
	},
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_LeftBrace_Idx, Token_ID_Idx},
	},
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_LeftParen_Idx, Token_String_Idx},
		parser.LookaheadPath{Keyword_LeftParen_Idx, Token_ID_Idx},
		parser.LookaheadPath{Keyword_LeftParen_Idx, Token_ID_Idx},
		parser.LookaheadPath{Keyword_LeftParen_Idx, Keyword_LeftBrace_Idx},
		parser.LookaheadPath{Keyword_LeftParen_Idx, Keyword_LeftParen_Idx},
	},
}

var FieldLookahead4 = parser.LLkLookahead{
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_LeftBracket_Idx},
	},
}

var GrammarLookahead0 = parser.LLkLookahead{
	parser.LookaheadOption{
		parser.LookaheadPath{Token_ID_Idx},
		parser.LookaheadPath{Keyword_token_Idx},
		parser.LookaheadPath{Keyword_hidden_Idx},
		parser.LookaheadPath{Keyword_interface_Idx},
	},
}

var GrammarLookaheadOr0 = parser.LLkLookahead{
	parser.LookaheadOption{
		parser.LookaheadPath{Token_ID_Idx},
	},
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_token_Idx},
		parser.LookaheadPath{Keyword_hidden_Idx},
	},
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_interface_Idx},
	},
}

var GroupLookahead8 = parser.LLkLookahead{
	parser.LookaheadOption{
		parser.LookaheadPath{Token_String_Idx},
		parser.LookaheadPath{Token_ID_Idx},
		parser.LookaheadPath{Token_ID_Idx},
		parser.LookaheadPath{Keyword_LeftBrace_Idx},
		parser.LookaheadPath{Keyword_LeftParen_Idx},
	},
}

var GroupLookahead9 = parser.LLkLookahead{
	parser.LookaheadOption{
		parser.LookaheadPath{Token_String_Idx},
		parser.LookaheadPath{Token_ID_Idx},
		parser.LookaheadPath{Token_ID_Idx},
		parser.LookaheadPath{Keyword_LeftBrace_Idx},
		parser.LookaheadPath{Keyword_LeftParen_Idx},
	},
}

var InterfaceLookahead1 = parser.LLkLookahead{
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_extends_Idx},
	},
}

var InterfaceLookahead2 = parser.LLkLookahead{
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_Comma_Idx},
	},
}

var InterfaceLookahead3 = parser.LLkLookahead{
	parser.LookaheadOption{
		parser.LookaheadPath{Token_ID_Idx},
	},
}

var TokenLookahead5 = parser.LLkLookahead{
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_hidden_Idx},
	},
}

func (p *Parser) ParseGrammar() Grammar {
	node := NewGrammar()
	node.WithSegmentStartToken(p.state.LA(1))
	{
		token := p.state.Consume(Keyword_grammar_Idx)
		core.AssignToken(node, token, Grammargrammar_0)
	}
	{
		token := p.state.Consume(Token_ID_Idx)
		core.AssignToken(node, token, GrammarNameID_0)
		if token != nil {
			node.WithName(token)
		}
	}
	{
		token := p.state.Consume(Keyword_Semicolon_Idx)
		core.AssignToken(node, token, GrammarSemicolon_0)
	}
	for p.state.Lookahead(GrammarLookahead0) == 0 {
		switch p.state.Lookahead(GrammarLookaheadOr0) {
		case 0:
			{
				result := p.ParseParserRule()
				if result != nil {
					node.WithRulesItem(result)
				}
			}
		case 1:
			{
				result := p.ParseToken()
				if result != nil {
					node.WithTerminalsItem(result)
				}
			}
		case 2:
			{
				result := p.ParseInterface()
				if result != nil {
					node.WithInterfacesItem(result)
				}
			}
		}
	}
	node.WithSegmentEndToken(p.state.LA(0))
	return node
}

func (p *Parser) ParseInterface() Interface {
	node := NewInterface()
	node.WithSegmentStartToken(p.state.LA(1))
	{
		token := p.state.Consume(Keyword_interface_Idx)
		core.AssignToken(node, token, Interfaceinterface_0)
	}
	{
		token := p.state.Consume(Token_ID_Idx)
		core.AssignToken(node, token, InterfaceNameID_0)
		if token != nil {
			node.WithName(token)
		}
	}
	if p.state.Lookahead(InterfaceLookahead1) == 0 {
		{
			token := p.state.Consume(Keyword_extends_Idx)
			core.AssignToken(node, token, Interfaceextends_0)
		}
		{
			token := p.state.Consume(Token_ID_Idx)
			core.AssignToken(node, token, InterfaceExtendsID_0)
			if token != nil {
				node.WithExtendsItem(token)
			}
		}
		for p.state.Lookahead(InterfaceLookahead2) == 0 {
			{
				token := p.state.Consume(Keyword_Comma_Idx)
				core.AssignToken(node, token, InterfaceComma_0)
			}
			{
				token := p.state.Consume(Token_ID_Idx)
				core.AssignToken(node, token, InterfaceExtendsID_1)
				if token != nil {
					node.WithExtendsItem(token)
				}
			}
		}
	}
	{
		token := p.state.Consume(Keyword_LeftBrace_Idx)
		core.AssignToken(node, token, InterfaceLeftBrace_0)
	}
	{
		for p.state.Lookahead(InterfaceLookahead3) == 0 {
			result := p.ParseField()
			if result != nil {
				node.WithFieldsItem(result)
			}
		}
	}
	{
		token := p.state.Consume(Keyword_RightBrace_Idx)
		core.AssignToken(node, token, InterfaceRightBrace_0)
	}
	node.WithSegmentEndToken(p.state.LA(0))
	return node
}

func (p *Parser) ParseField() Field {
	node := NewField()
	node.WithSegmentStartToken(p.state.LA(1))
	{
		token := p.state.Consume(Token_ID_Idx)
		core.AssignToken(node, token, FieldNameID_0)
		if token != nil {
			node.WithName(token)
		}
	}
	if p.state.Lookahead(FieldLookahead4) == 0 {
		{
			token := p.state.Consume(Keyword_LeftBracket_Idx)
			core.AssignToken(node, token, FieldArrayLeftBracket_0)
			if token != nil {
				node.WithArray(token)
			}
		}
		{
			token := p.state.Consume(Keyword_RightBracket_Idx)
			core.AssignToken(node, token, FieldRightBracket_0)
		}
	}
	{
		token := p.state.Consume(Token_ID_Idx)
		core.AssignToken(node, token, FieldTypeID_0)
		if token != nil {
			node.WithType(token)
		}
	}
	node.WithSegmentEndToken(p.state.LA(0))
	return node
}

func (p *Parser) ParseParserRule() ParserRule {
	node := NewParserRule()
	node.WithSegmentStartToken(p.state.LA(1))
	{
		token := p.state.Consume(Token_ID_Idx)
		core.AssignToken(node, token, ParserRuleNameID_0)
		if token != nil {
			node.WithName(token)
		}
	}
	{
		token := p.state.Consume(Keyword_returns_Idx)
		core.AssignToken(node, token, ParserRulereturns_0)
	}
	{
		token := p.state.Consume(Token_ID_Idx)
		core.AssignToken(node, token, ParserRuleReturnTypeID_0)
		if token != nil {
			node.WithReturnType(token)
		}
	}
	{
		token := p.state.Consume(Keyword_Colon_Idx)
		core.AssignToken(node, token, ParserRuleColon_0)
	}
	{
		result := p.ParseAlternatives()
		if result != nil {
			node.WithBody(result)
		}
	}
	{
		token := p.state.Consume(Keyword_Semicolon_Idx)
		core.AssignToken(node, token, ParserRuleSemicolon_0)
	}
	node.WithSegmentEndToken(p.state.LA(0))
	return node
}

func (p *Parser) ParseToken() Token {
	node := NewToken()
	node.WithSegmentStartToken(p.state.LA(1))
	{
		if p.state.Lookahead(TokenLookahead5) == 0 {
			token := p.state.Consume(Keyword_hidden_Idx)
			core.AssignToken(node, token, TokenTypehidden_0)
			if token != nil {
				node.WithType(token)
			}
		}
	}
	{
		token := p.state.Consume(Keyword_token_Idx)
		core.AssignToken(node, token, Tokentoken_0)
	}
	{
		token := p.state.Consume(Token_ID_Idx)
		core.AssignToken(node, token, TokenNameID_0)
		if token != nil {
			node.WithName(token)
		}
	}
	{
		token := p.state.Consume(Keyword_Colon_Idx)
		core.AssignToken(node, token, TokenColon_0)
	}
	{
		token := p.state.Consume(Token_RegexLiteral_Idx)
		core.AssignToken(node, token, TokenRegexpRegexLiteral_0)
		if token != nil {
			node.WithRegexp(token)
		}
	}
	{
		token := p.state.Consume(Keyword_Semicolon_Idx)
		core.AssignToken(node, token, TokenSemicolon_0)
	}
	node.WithSegmentEndToken(p.state.LA(0))
	return node
}

func (p *Parser) ParseAlternatives() Element {
	node := NewElement()
	node.WithSegmentStartToken(p.state.LA(1))
	{
		result := p.ParseGroup()
		core.AssignTokens(result, node.Tokens())
		node = result
	}
	if p.state.Lookahead(AlternativesLookahead6) == 0 {
		{
			result := NewAlternatives()
			result.WithAltsItem(node)
			node = result
		}
		for ok := true; ok; ok = p.state.Lookahead(AlternativesLookahead7) == 0 {
			{
				token := p.state.Consume(Keyword_Pipe_Idx)
				core.AssignToken(node, token, AlternativesPipe_0)
			}
			{
				result := p.ParseGroup()
				if result != nil {
					node.(Alternatives).WithAltsItem(result)
				}
			}
		}
	}
	node.WithSegmentEndToken(p.state.LA(0))
	return node
}

func (p *Parser) ParseGroup() Element {
	node := NewElement()
	node.WithSegmentStartToken(p.state.LA(1))
	{
		result := p.ParseElement()
		core.AssignTokens(result, node.Tokens())
		node = result
	}
	if p.state.Lookahead(GroupLookahead8) == 0 {
		{
			result := NewGroup()
			result.WithElementsItem(node)
			node = result
		}
		{
			for ok := true; ok; ok = p.state.Lookahead(GroupLookahead9) == 0 {
				result := p.ParseElement()
				if result != nil {
					node.(Group).WithElementsItem(result)
				}
			}
		}
	}
	node.WithSegmentEndToken(p.state.LA(0))
	return node
}

func (p *Parser) ParseElement() Element {
	node := NewElement()
	node.WithSegmentStartToken(p.state.LA(1))
	switch p.state.Lookahead(ElementLookaheadOr1) {
	case 0:
		{
			result := p.ParseKeyword()
			core.AssignTokens(result, node.Tokens())
			node = result
		}
	case 1:
		{
			result := p.ParseAssignment()
			core.AssignTokens(result, node.Tokens())
			node = result
		}
	case 2:
		{
			result := p.ParseRuleCall()
			core.AssignTokens(result, node.Tokens())
			node = result
		}
	case 3:
		{
			result := p.ParseAction()
			core.AssignTokens(result, node.Tokens())
			node = result
		}
	case 4:
		{
			token := p.state.Consume(Keyword_LeftParen_Idx)
			core.AssignToken(node, token, ElementLeftParen_0)
		}
		{
			result := p.ParseAlternatives()
			core.AssignTokens(result, node.Tokens())
			node = result
		}
		{
			token := p.state.Consume(Keyword_RightParen_Idx)
			core.AssignToken(node, token, ElementRightParen_0)
		}
	}
	{
		if p.state.Lookahead(ElementLookahead10) == 0 {
			token := p.state.Consume(Token_Cardinality_Idx)
			core.AssignToken(node, token, ElementCardinalityCardinality_0)
			if token != nil {
				node.WithCardinality(token)
			}
		}
	}
	node.WithSegmentEndToken(p.state.LA(0))
	return node
}

func (p *Parser) ParseKeyword() Keyword {
	node := NewKeyword()
	node.WithSegmentStartToken(p.state.LA(1))
	{
		token := p.state.Consume(Token_String_Idx)
		core.AssignToken(node, token, KeywordValueString_0)
		if token != nil {
			node.WithValue(token)
		}
	}
	node.WithSegmentEndToken(p.state.LA(0))
	return node
}

func (p *Parser) ParseAssignment() Assignment {
	node := NewAssignment()
	node.WithSegmentStartToken(p.state.LA(1))
	{
		token := p.state.Consume(Token_ID_Idx)
		core.AssignToken(node, token, AssignmentPropertyID_0)
		if token != nil {
			node.WithProperty(token)
		}
	}
	{
		token := p.state.Consume(Token_AssignmentOperator_Idx)
		core.AssignToken(node, token, AssignmentOperatorAssignmentOperator_0)
		if token != nil {
			node.WithOperator(token)
		}
	}
	{
		result := p.ParseAssignable()
		if result != nil {
			node.WithValue(result)
		}
	}
	node.WithSegmentEndToken(p.state.LA(0))
	return node
}

func (p *Parser) ParseAssignable() Assignable {
	node := NewAssignable()
	node.WithSegmentStartToken(p.state.LA(1))
	switch p.state.Lookahead(AssignableLookaheadOr2) {
	case 0:
		{
			result := p.ParseKeyword()
			core.AssignTokens(result, node.Tokens())
			node = result
		}
	case 1:
		{
			result := p.ParseRuleCall()
			core.AssignTokens(result, node.Tokens())
			node = result
		}
	case 2:
		{
			result := p.ParseCrossRef()
			core.AssignTokens(result, node.Tokens())
			node = result
		}
	}
	node.WithSegmentEndToken(p.state.LA(0))
	return node
}

func (p *Parser) ParseCrossRef() CrossRef {
	node := NewCrossRef()
	node.WithSegmentStartToken(p.state.LA(1))
	{
		token := p.state.Consume(Keyword_LeftBracket_Idx)
		core.AssignToken(node, token, CrossRefLeftBracket_0)
	}
	{
		token := p.state.Consume(Token_ID_Idx)
		core.AssignToken(node, token, CrossRefTypeID_0)
		if token != nil {
			node.WithType(token)
		}
	}
	if p.state.Lookahead(CrossRefLookahead11) == 0 {
		{
			token := p.state.Consume(Keyword_Colon_Idx)
			core.AssignToken(node, token, CrossRefColon_0)
		}
		{
			result := p.ParseRuleCall()
			if result != nil {
				node.WithRule(result)
			}
		}
	}
	{
		token := p.state.Consume(Keyword_RightBracket_Idx)
		core.AssignToken(node, token, CrossRefRightBracket_0)
	}
	node.WithSegmentEndToken(p.state.LA(0))
	return node
}

func (p *Parser) ParseRuleCall() RuleCall {
	node := NewRuleCall()
	node.WithSegmentStartToken(p.state.LA(1))
	{
		token := p.state.Consume(Token_ID_Idx)
		core.AssignToken(node, token, RuleCallRuleID_0)
		if token != nil {
			node.WithRule(token)
		}
	}
	node.WithSegmentEndToken(p.state.LA(0))
	return node
}

func (p *Parser) ParseAction() Action {
	node := NewAction()
	node.WithSegmentStartToken(p.state.LA(1))
	{
		token := p.state.Consume(Keyword_LeftBrace_Idx)
		core.AssignToken(node, token, ActionLeftBrace_0)
	}
	{
		token := p.state.Consume(Token_ID_Idx)
		core.AssignToken(node, token, ActionTypeID_0)
		if token != nil {
			node.WithType(token)
		}
	}
	if p.state.Lookahead(ActionLookahead12) == 0 {
		{
			token := p.state.Consume(Keyword_Dot_Idx)
			core.AssignToken(node, token, ActionDot_0)
		}
		{
			token := p.state.Consume(Token_ID_Idx)
			core.AssignToken(node, token, ActionPropertyID_0)
			if token != nil {
				node.WithProperty(token)
			}
		}
		{
			token := p.state.Consume(Token_AssignmentOperator_Idx)
			core.AssignToken(node, token, ActionOperatorAssignmentOperator_0)
			if token != nil {
				node.WithOperator(token)
			}
		}
		{
			token := p.state.Consume(Keyword_current_Idx)
			core.AssignToken(node, token, Actioncurrent_0)
		}
	}
	{
		token := p.state.Consume(Keyword_RightBrace_Idx)
		core.AssignToken(node, token, ActionRightBrace_0)
	}
	node.WithSegmentEndToken(p.state.LA(0))
	return node
}
