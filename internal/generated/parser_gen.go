package generated

import (
	"github.com/TypeFox/langium-to-go/lexer"
	"github.com/TypeFox/langium-to-go/parser"
)

type Parser struct {
	state *parser.ParserState
}

func (p *Parser) Parse(tokens []*lexer.Token) Grammar {
	p.state = parser.NewParserState(tokens)
	return p.ParseGrammar()
}

func NewParser() *Parser {
	return &Parser{}
}

const (
	ActionDot_0 = iota + 1
	ActionLeftBrace_0
	ActionPropertyID_0
	ActionRightBrace_0
	ActionTypeID_0
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

var ActionLookahead10 = parser.LLkLookahead{
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_Dot_Idx},
	},
}

var AlternativesLookahead6 = parser.LLkLookahead{
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

var CrossRefLookahead9 = parser.LLkLookahead{
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_Colon_Idx},
	},
}

var ElementLookahead8 = parser.LLkLookahead{
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

var GroupLookahead7 = parser.LLkLookahead{
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
	{
		token := p.state.Consume(Keyword_grammar_Idx)
		node.WithToken(token, Grammargrammar_0)
	}
	{
		token := p.state.Consume(Token_ID_Idx)
		node.WithToken(token, GrammarNameID_0)
		if token != nil {
			node.WithName(token)
		}
	}
	{
		token := p.state.Consume(Keyword_Semicolon_Idx)
		node.WithToken(token, GrammarSemicolon_0)
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
	return node
}

func (p *Parser) ParseInterface() Interface {
	node := NewInterface()
	{
		token := p.state.Consume(Keyword_interface_Idx)
		node.WithToken(token, Interfaceinterface_0)
	}
	{
		token := p.state.Consume(Token_ID_Idx)
		node.WithToken(token, InterfaceNameID_0)
		if token != nil {
			node.WithName(token)
		}
	}
	if p.state.Lookahead(InterfaceLookahead1) == 0 {
		{
			token := p.state.Consume(Keyword_extends_Idx)
			node.WithToken(token, Interfaceextends_0)
		}
		{
			token := p.state.Consume(Token_ID_Idx)
			node.WithToken(token, InterfaceExtendsID_0)
			if token != nil {
				node.WithExtendsItem(token)
			}
		}
		for p.state.Lookahead(InterfaceLookahead2) == 0 {
			{
				token := p.state.Consume(Keyword_Comma_Idx)
				node.WithToken(token, InterfaceComma_0)
			}
			{
				token := p.state.Consume(Token_ID_Idx)
				node.WithToken(token, InterfaceExtendsID_1)
				if token != nil {
					node.WithExtendsItem(token)
				}
			}
		}
	}
	{
		token := p.state.Consume(Keyword_LeftBrace_Idx)
		node.WithToken(token, InterfaceLeftBrace_0)
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
		node.WithToken(token, InterfaceRightBrace_0)
	}
	return node
}

func (p *Parser) ParseField() Field {
	node := NewField()
	{
		token := p.state.Consume(Token_ID_Idx)
		node.WithToken(token, FieldNameID_0)
		if token != nil {
			node.WithName(token)
		}
	}
	if p.state.Lookahead(FieldLookahead4) == 0 {
		{
			token := p.state.Consume(Keyword_LeftBracket_Idx)
			node.WithToken(token, FieldArrayLeftBracket_0)
			if token != nil {
				node.WithArray(token)
			}
		}
		{
			token := p.state.Consume(Keyword_RightBracket_Idx)
			node.WithToken(token, FieldRightBracket_0)
		}
	}
	{
		token := p.state.Consume(Token_ID_Idx)
		node.WithToken(token, FieldTypeID_0)
		if token != nil {
			node.WithType(token)
		}
	}
	return node
}

func (p *Parser) ParseParserRule() ParserRule {
	node := NewParserRule()
	{
		token := p.state.Consume(Token_ID_Idx)
		node.WithToken(token, ParserRuleNameID_0)
		if token != nil {
			node.WithName(token)
		}
	}
	{
		token := p.state.Consume(Keyword_returns_Idx)
		node.WithToken(token, ParserRulereturns_0)
	}
	{
		token := p.state.Consume(Token_ID_Idx)
		node.WithToken(token, ParserRuleReturnTypeID_0)
		if token != nil {
			node.WithReturnType(token)
		}
	}
	{
		token := p.state.Consume(Keyword_Colon_Idx)
		node.WithToken(token, ParserRuleColon_0)
	}
	{
		result := p.ParseAlternatives()
		if result != nil {
			node.WithBody(result)
		}
	}
	{
		token := p.state.Consume(Keyword_Semicolon_Idx)
		node.WithToken(token, ParserRuleSemicolon_0)
	}
	return node
}

func (p *Parser) ParseToken() Token {
	node := NewToken()
	{
		if p.state.Lookahead(TokenLookahead5) == 0 {
			token := p.state.Consume(Keyword_hidden_Idx)
			node.WithToken(token, TokenTypehidden_0)
			if token != nil {
				node.WithType(token)
			}
		}
	}
	{
		token := p.state.Consume(Keyword_token_Idx)
		node.WithToken(token, Tokentoken_0)
	}
	{
		token := p.state.Consume(Token_ID_Idx)
		node.WithToken(token, TokenNameID_0)
		if token != nil {
			node.WithName(token)
		}
	}
	{
		token := p.state.Consume(Keyword_Colon_Idx)
		node.WithToken(token, TokenColon_0)
	}
	{
		token := p.state.Consume(Token_RegexLiteral_Idx)
		node.WithToken(token, TokenRegexpRegexLiteral_0)
		if token != nil {
			node.WithRegexp(token)
		}
	}
	{
		token := p.state.Consume(Keyword_Semicolon_Idx)
		node.WithToken(token, TokenSemicolon_0)
	}
	return node
}

func (p *Parser) ParseAlternatives() Alternatives {
	node := NewAlternatives()
	{
		result := p.ParseGroup()
		if result != nil {
			node.WithAltsItem(result)
		}
	}
	for p.state.Lookahead(AlternativesLookahead6) == 0 {
		{
			token := p.state.Consume(Keyword_Pipe_Idx)
			node.WithToken(token, AlternativesPipe_0)
		}
		{
			result := p.ParseGroup()
			if result != nil {
				node.WithAltsItem(result)
			}
		}
	}
	return node
}

func (p *Parser) ParseGroup() Group {
	node := NewGroup()
	{
		for ok := true; ok; ok = p.state.Lookahead(GroupLookahead7) == 0 {
			result := p.ParseElement()
			if result != nil {
				node.WithElementsItem(result)
			}
		}
	}
	return node
}

func (p *Parser) ParseElement() Element {
	node := NewElement()
	switch p.state.Lookahead(ElementLookaheadOr1) {
	case 0:
		{
			result := p.ParseKeyword()
			result.WithTokens(node.Tokens())
			node = result
		}
	case 1:
		{
			result := p.ParseAssignment()
			result.WithTokens(node.Tokens())
			node = result
		}
	case 2:
		{
			result := p.ParseRuleCall()
			result.WithTokens(node.Tokens())
			node = result
		}
	case 3:
		{
			result := p.ParseAction()
			result.WithTokens(node.Tokens())
			node = result
		}
	case 4:
		{
			token := p.state.Consume(Keyword_LeftParen_Idx)
			node.WithToken(token, ElementLeftParen_0)
		}
		{
			result := p.ParseAlternatives()
			result.WithTokens(node.Tokens())
			node = result
		}
		{
			token := p.state.Consume(Keyword_RightParen_Idx)
			node.WithToken(token, ElementRightParen_0)
		}
	}
	{
		if p.state.Lookahead(ElementLookahead8) == 0 {
			token := p.state.Consume(Token_Cardinality_Idx)
			node.WithToken(token, ElementCardinalityCardinality_0)
			if token != nil {
				node.WithCardinality(token)
			}
		}
	}
	return node
}

func (p *Parser) ParseKeyword() Keyword {
	node := NewKeyword()
	{
		token := p.state.Consume(Token_String_Idx)
		node.WithToken(token, KeywordValueString_0)
		if token != nil {
			node.WithValue(token)
		}
	}
	return node
}

func (p *Parser) ParseAssignment() Assignment {
	node := NewAssignment()
	{
		token := p.state.Consume(Token_ID_Idx)
		node.WithToken(token, AssignmentPropertyID_0)
		if token != nil {
			node.WithProperty(token)
		}
	}
	{
		token := p.state.Consume(Token_AssignmentOperator_Idx)
		node.WithToken(token, AssignmentOperatorAssignmentOperator_0)
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
	return node
}

func (p *Parser) ParseAssignable() Assignable {
	node := NewAssignable()
	switch p.state.Lookahead(AssignableLookaheadOr2) {
	case 0:
		{
			result := p.ParseKeyword()
			result.WithTokens(node.Tokens())
			node = result
		}
	case 1:
		{
			result := p.ParseRuleCall()
			result.WithTokens(node.Tokens())
			node = result
		}
	case 2:
		{
			result := p.ParseCrossRef()
			result.WithTokens(node.Tokens())
			node = result
		}
	}
	return node
}

func (p *Parser) ParseCrossRef() CrossRef {
	node := NewCrossRef()
	{
		token := p.state.Consume(Keyword_LeftBracket_Idx)
		node.WithToken(token, CrossRefLeftBracket_0)
	}
	{
		token := p.state.Consume(Token_ID_Idx)
		node.WithToken(token, CrossRefTypeID_0)
		if token != nil {
			node.WithType(token)
		}
	}
	if p.state.Lookahead(CrossRefLookahead9) == 0 {
		{
			token := p.state.Consume(Keyword_Colon_Idx)
			node.WithToken(token, CrossRefColon_0)
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
		node.WithToken(token, CrossRefRightBracket_0)
	}
	return node
}

func (p *Parser) ParseRuleCall() RuleCall {
	node := NewRuleCall()
	{
		token := p.state.Consume(Token_ID_Idx)
		node.WithToken(token, RuleCallRuleID_0)
		if token != nil {
			node.WithRule(token)
		}
	}
	return node
}

func (p *Parser) ParseAction() Action {
	node := NewAction()
	{
		token := p.state.Consume(Keyword_LeftBrace_Idx)
		node.WithToken(token, ActionLeftBrace_0)
	}
	{
		token := p.state.Consume(Token_ID_Idx)
		node.WithToken(token, ActionTypeID_0)
		if token != nil {
			node.WithType(token)
		}
	}
	if p.state.Lookahead(ActionLookahead10) == 0 {
		{
			token := p.state.Consume(Keyword_Dot_Idx)
			node.WithToken(token, ActionDot_0)
		}
		{
			token := p.state.Consume(Token_ID_Idx)
			node.WithToken(token, ActionPropertyID_0)
			if token != nil {
				node.WithProperty(token)
			}
		}
	}
	{
		token := p.state.Consume(Keyword_RightBrace_Idx)
		node.WithToken(token, ActionRightBrace_0)
	}
	return node
}
