package generated

import (
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/parser"
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
	ActionOperatorEquals_0
	ActionOperatorPlusEquals_0
	ActionPropertyID_0
	ActionRightBrace_0
	ActionTypeID_0
	Actioncurrent_0
	AlternativesPipe_0
	AssignableAlternativesPipe_0
	AssignableLeftParen_0
	AssignableRightParen_0
	AssignmentOperatorEquals_0
	AssignmentOperatorPlusEquals_0
	AssignmentOperatorQuestionEquals_0
	AssignmentPropertyID_0
	CrossRefColon_0
	CrossRefLeftBracket_0
	CrossRefRightBracket_0
	CrossRefTypeID_0
	ElementCardinalityAsterisk_0
	ElementCardinalityPlus_0
	ElementCardinalityQuestion_0
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

var ActionLookahead14 = parser.LLkLookahead{
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_Dot_Idx},
	},
}

var ActionOperatorLookaheadOr6 = parser.LLkLookahead{
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_PlusEquals_Idx},
	},
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_Equals_Idx},
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

var AssignableAlternativesLookahead11 = parser.LLkLookahead{
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_Pipe_Idx},
	},
}

var AssignableAlternativesLookahead12 = parser.LLkLookahead{
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_Pipe_Idx},
	},
}

var AssignableLookaheadOr4 = parser.LLkLookahead{
	parser.LookaheadOption{
		parser.LookaheadPath{Token_String_Idx},
	},
	parser.LookaheadOption{
		parser.LookaheadPath{Token_ID_Idx},
	},
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_LeftBracket_Idx},
	},
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_LeftParen_Idx},
	},
}

var AssignableWithoutAltsLookaheadOr5 = parser.LLkLookahead{
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

var AssignmentOperatorLookaheadOr3 = parser.LLkLookahead{
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_PlusEquals_Idx},
	},
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_Equals_Idx},
	},
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_QuestionEquals_Idx},
	},
}

var CrossRefLookahead13 = parser.LLkLookahead{
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_Colon_Idx},
	},
}

var ElementCardinalityLookaheadOr2 = parser.LLkLookahead{
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_Asterisk_Idx},
	},
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_Plus_Idx},
	},
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_Question_Idx},
	},
}

var ElementLookahead10 = parser.LLkLookahead{
	parser.LookaheadOption{
		parser.LookaheadPath{Keyword_Asterisk_Idx},
		parser.LookaheadPath{Keyword_Plus_Idx},
		parser.LookaheadPath{Keyword_Question_Idx},
	},
}

var ElementLookaheadOr1 = parser.LLkLookahead{
	parser.LookaheadOption{
		parser.LookaheadPath{Token_String_Idx},
	},
	parser.LookaheadOption{
		parser.LookaheadPath{Token_ID_Idx, Keyword_PlusEquals_Idx},
		parser.LookaheadPath{Token_ID_Idx, Keyword_Equals_Idx},
		parser.LookaheadPath{Token_ID_Idx, Keyword_QuestionEquals_Idx},
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
	current := NewGrammar()
	current.WithSegmentStartToken(p.state.LA(1))
	{
		token := p.state.Consume(Keyword_grammar_Idx)
		core.AssignToken(current, token, Grammargrammar_0)
	}
	{
		token := p.state.Consume(Token_ID_Idx)
		core.AssignToken(current, token, GrammarNameID_0)
		if token != nil {
			current.WithName(token)
		}
	}
	{
		token := p.state.Consume(Keyword_Semicolon_Idx)
		core.AssignToken(current, token, GrammarSemicolon_0)
	}
	for p.state.Lookahead(GrammarLookahead0) == 0 {
		switch p.state.Lookahead(GrammarLookaheadOr0) {
		case 0:
			{
				result := p.ParseParserRule()
				if result != nil {
					current.WithRulesItem(result)
				}
			}
		case 1:
			{
				result := p.ParseToken()
				if result != nil {
					current.WithTerminalsItem(result)
				}
			}
		case 2:
			{
				result := p.ParseInterface()
				if result != nil {
					current.WithInterfacesItem(result)
				}
			}
		}
	}
	current.WithSegmentEndToken(p.state.LA(0))
	return current
}

func (p *Parser) ParseInterface() Interface {
	current := NewInterface()
	current.WithSegmentStartToken(p.state.LA(1))
	{
		token := p.state.Consume(Keyword_interface_Idx)
		core.AssignToken(current, token, Interfaceinterface_0)
	}
	{
		token := p.state.Consume(Token_ID_Idx)
		core.AssignToken(current, token, InterfaceNameID_0)
		if token != nil {
			current.WithName(token)
		}
	}
	if p.state.Lookahead(InterfaceLookahead1) == 0 {
		{
			token := p.state.Consume(Keyword_extends_Idx)
			core.AssignToken(current, token, Interfaceextends_0)
		}
		{
			token := p.state.Consume(Token_ID_Idx)
			core.AssignToken(current, token, InterfaceExtendsID_0)
			if token != nil {
				current.WithExtendsItem(token)
			}
		}
		for p.state.Lookahead(InterfaceLookahead2) == 0 {
			{
				token := p.state.Consume(Keyword_Comma_Idx)
				core.AssignToken(current, token, InterfaceComma_0)
			}
			{
				token := p.state.Consume(Token_ID_Idx)
				core.AssignToken(current, token, InterfaceExtendsID_1)
				if token != nil {
					current.WithExtendsItem(token)
				}
			}
		}
	}
	{
		token := p.state.Consume(Keyword_LeftBrace_Idx)
		core.AssignToken(current, token, InterfaceLeftBrace_0)
	}
	{
		for p.state.Lookahead(InterfaceLookahead3) == 0 {
			result := p.ParseField()
			if result != nil {
				current.WithFieldsItem(result)
			}
		}
	}
	{
		token := p.state.Consume(Keyword_RightBrace_Idx)
		core.AssignToken(current, token, InterfaceRightBrace_0)
	}
	current.WithSegmentEndToken(p.state.LA(0))
	return current
}

func (p *Parser) ParseField() Field {
	current := NewField()
	current.WithSegmentStartToken(p.state.LA(1))
	{
		token := p.state.Consume(Token_ID_Idx)
		core.AssignToken(current, token, FieldNameID_0)
		if token != nil {
			current.WithName(token)
		}
	}
	if p.state.Lookahead(FieldLookahead4) == 0 {
		{
			token := p.state.Consume(Keyword_LeftBracket_Idx)
			core.AssignToken(current, token, FieldArrayLeftBracket_0)
			if token != nil {
				current.WithArray(token)
			}
		}
		{
			token := p.state.Consume(Keyword_RightBracket_Idx)
			core.AssignToken(current, token, FieldRightBracket_0)
		}
	}
	{
		token := p.state.Consume(Token_ID_Idx)
		core.AssignToken(current, token, FieldTypeID_0)
		if token != nil {
			current.WithType(token)
		}
	}
	current.WithSegmentEndToken(p.state.LA(0))
	return current
}

func (p *Parser) ParseParserRule() ParserRule {
	current := NewParserRule()
	current.WithSegmentStartToken(p.state.LA(1))
	{
		token := p.state.Consume(Token_ID_Idx)
		core.AssignToken(current, token, ParserRuleNameID_0)
		if token != nil {
			current.WithName(token)
		}
	}
	{
		token := p.state.Consume(Keyword_returns_Idx)
		core.AssignToken(current, token, ParserRulereturns_0)
	}
	{
		token := p.state.Consume(Token_ID_Idx)
		core.AssignToken(current, token, ParserRuleReturnTypeID_0)
		if token != nil {
			current.WithReturnType(token)
		}
	}
	{
		token := p.state.Consume(Keyword_Colon_Idx)
		core.AssignToken(current, token, ParserRuleColon_0)
	}
	{
		result := p.ParseAlternatives()
		if result != nil {
			current.WithBody(result)
		}
	}
	{
		token := p.state.Consume(Keyword_Semicolon_Idx)
		core.AssignToken(current, token, ParserRuleSemicolon_0)
	}
	current.WithSegmentEndToken(p.state.LA(0))
	return current
}

func (p *Parser) ParseToken() Token {
	current := NewToken()
	current.WithSegmentStartToken(p.state.LA(1))
	{
		if p.state.Lookahead(TokenLookahead5) == 0 {
			token := p.state.Consume(Keyword_hidden_Idx)
			core.AssignToken(current, token, TokenTypehidden_0)
			if token != nil {
				current.WithType(token)
			}
		}
	}
	{
		token := p.state.Consume(Keyword_token_Idx)
		core.AssignToken(current, token, Tokentoken_0)
	}
	{
		token := p.state.Consume(Token_ID_Idx)
		core.AssignToken(current, token, TokenNameID_0)
		if token != nil {
			current.WithName(token)
		}
	}
	{
		token := p.state.Consume(Keyword_Colon_Idx)
		core.AssignToken(current, token, TokenColon_0)
	}
	{
		token := p.state.Consume(Token_RegexLiteral_Idx)
		core.AssignToken(current, token, TokenRegexpRegexLiteral_0)
		if token != nil {
			current.WithRegexp(token)
		}
	}
	{
		token := p.state.Consume(Keyword_Semicolon_Idx)
		core.AssignToken(current, token, TokenSemicolon_0)
	}
	current.WithSegmentEndToken(p.state.LA(0))
	return current
}

func (p *Parser) ParseAlternatives() Element {
	current := NewElement()
	current.WithSegmentStartToken(p.state.LA(1))
	{
		result := p.ParseGroup()
		core.MergeTokens(result, current.Tokens())
		current = result
	}
	if p.state.Lookahead(AlternativesLookahead6) == 0 {
		{
			result := NewAlternatives()
			result.WithSegment(current.Segment())
			result.WithAltsItem(current)
			current.WithSegmentEndToken(p.state.LA(0))
			current = result
		}
		current := current.(Alternatives)
		for ok := true; ok; ok = p.state.Lookahead(AlternativesLookahead7) == 0 {
			{
				token := p.state.Consume(Keyword_Pipe_Idx)
				core.AssignToken(current, token, AlternativesPipe_0)
			}
			{
				result := p.ParseGroup()
				if result != nil {
					current.WithAltsItem(result)
				}
			}
		}
	}
	current.WithSegmentEndToken(p.state.LA(0))
	return current
}

func (p *Parser) ParseGroup() Element {
	current := NewElement()
	current.WithSegmentStartToken(p.state.LA(1))
	{
		result := p.ParseElement()
		core.MergeTokens(result, current.Tokens())
		current = result
	}
	if p.state.Lookahead(GroupLookahead8) == 0 {
		{
			result := NewGroup()
			result.WithSegment(current.Segment())
			result.WithElementsItem(current)
			current.WithSegmentEndToken(p.state.LA(0))
			current = result
		}
		current := current.(Group)
		{
			for ok := true; ok; ok = p.state.Lookahead(GroupLookahead9) == 0 {
				result := p.ParseElement()
				if result != nil {
					current.WithElementsItem(result)
				}
			}
		}
	}
	current.WithSegmentEndToken(p.state.LA(0))
	return current
}

func (p *Parser) ParseElement() Element {
	current := NewElement()
	current.WithSegmentStartToken(p.state.LA(1))
	switch p.state.Lookahead(ElementLookaheadOr1) {
	case 0:
		{
			result := p.ParseKeyword()
			core.MergeTokens(result, current.Tokens())
			current = result
		}
	case 1:
		{
			result := p.ParseAssignment()
			core.MergeTokens(result, current.Tokens())
			current = result
		}
	case 2:
		{
			result := p.ParseRuleCall()
			core.MergeTokens(result, current.Tokens())
			current = result
		}
	case 3:
		{
			result := p.ParseAction()
			core.MergeTokens(result, current.Tokens())
			current = result
		}
	case 4:
		{
			token := p.state.Consume(Keyword_LeftParen_Idx)
			core.AssignToken(current, token, ElementLeftParen_0)
		}
		{
			result := p.ParseAlternatives()
			core.MergeTokens(result, current.Tokens())
			current = result
		}
		{
			token := p.state.Consume(Keyword_RightParen_Idx)
			core.AssignToken(current, token, ElementRightParen_0)
		}
	}
	{
		if p.state.Lookahead(ElementLookahead10) == 0 {
			switch p.state.Lookahead(ElementCardinalityLookaheadOr2) {
			case 0:
				token := p.state.Consume(Keyword_Asterisk_Idx)
				core.AssignToken(current, token, ElementCardinalityAsterisk_0)
				if token != nil {
					current.WithCardinality(token)
				}
			case 1:
				token := p.state.Consume(Keyword_Plus_Idx)
				core.AssignToken(current, token, ElementCardinalityPlus_0)
				if token != nil {
					current.WithCardinality(token)
				}
			case 2:
				token := p.state.Consume(Keyword_Question_Idx)
				core.AssignToken(current, token, ElementCardinalityQuestion_0)
				if token != nil {
					current.WithCardinality(token)
				}
			}
		}
	}
	current.WithSegmentEndToken(p.state.LA(0))
	return current
}

func (p *Parser) ParseKeyword() Keyword {
	current := NewKeyword()
	current.WithSegmentStartToken(p.state.LA(1))
	{
		token := p.state.Consume(Token_String_Idx)
		core.AssignToken(current, token, KeywordValueString_0)
		if token != nil {
			current.WithValue(token)
		}
	}
	current.WithSegmentEndToken(p.state.LA(0))
	return current
}

func (p *Parser) ParseAssignment() Assignment {
	current := NewAssignment()
	current.WithSegmentStartToken(p.state.LA(1))
	{
		token := p.state.Consume(Token_ID_Idx)
		core.AssignToken(current, token, AssignmentPropertyID_0)
		if token != nil {
			current.WithProperty(token)
		}
	}
	{
		switch p.state.Lookahead(AssignmentOperatorLookaheadOr3) {
		case 0:
			token := p.state.Consume(Keyword_PlusEquals_Idx)
			core.AssignToken(current, token, AssignmentOperatorPlusEquals_0)
			if token != nil {
				current.WithOperator(token)
			}
		case 1:
			token := p.state.Consume(Keyword_Equals_Idx)
			core.AssignToken(current, token, AssignmentOperatorEquals_0)
			if token != nil {
				current.WithOperator(token)
			}
		case 2:
			token := p.state.Consume(Keyword_QuestionEquals_Idx)
			core.AssignToken(current, token, AssignmentOperatorQuestionEquals_0)
			if token != nil {
				current.WithOperator(token)
			}
		}
	}
	{
		result := p.ParseAssignable()
		if result != nil {
			current.WithValue(result)
		}
	}
	current.WithSegmentEndToken(p.state.LA(0))
	return current
}

func (p *Parser) ParseAssignable() Assignable {
	current := NewAssignable()
	current.WithSegmentStartToken(p.state.LA(1))
	switch p.state.Lookahead(AssignableLookaheadOr4) {
	case 0:
		{
			result := p.ParseKeyword()
			core.MergeTokens(result, current.Tokens())
			current = result
		}
	case 1:
		{
			result := p.ParseRuleCall()
			core.MergeTokens(result, current.Tokens())
			current = result
		}
	case 2:
		{
			result := p.ParseCrossRef()
			core.MergeTokens(result, current.Tokens())
			current = result
		}
	case 3:
		{
			token := p.state.Consume(Keyword_LeftParen_Idx)
			core.AssignToken(current, token, AssignableLeftParen_0)
		}
		{
			result := p.ParseAssignableAlternatives()
			core.MergeTokens(result, current.Tokens())
			current = result
		}
		{
			token := p.state.Consume(Keyword_RightParen_Idx)
			core.AssignToken(current, token, AssignableRightParen_0)
		}
	}
	current.WithSegmentEndToken(p.state.LA(0))
	return current
}

func (p *Parser) ParseAssignableWithoutAlts() Assignable {
	current := NewAssignable()
	current.WithSegmentStartToken(p.state.LA(1))
	switch p.state.Lookahead(AssignableWithoutAltsLookaheadOr5) {
	case 0:
		{
			result := p.ParseKeyword()
			core.MergeTokens(result, current.Tokens())
			current = result
		}
	case 1:
		{
			result := p.ParseRuleCall()
			core.MergeTokens(result, current.Tokens())
			current = result
		}
	case 2:
		{
			result := p.ParseCrossRef()
			core.MergeTokens(result, current.Tokens())
			current = result
		}
	}
	current.WithSegmentEndToken(p.state.LA(0))
	return current
}

func (p *Parser) ParseAssignableAlternatives() Assignable {
	current := NewAssignable()
	current.WithSegmentStartToken(p.state.LA(1))
	{
		result := p.ParseAssignableWithoutAlts()
		core.MergeTokens(result, current.Tokens())
		current = result
	}
	if p.state.Lookahead(AssignableAlternativesLookahead11) == 0 {
		{
			result := NewAlternatives()
			result.WithSegment(current.Segment())
			result.WithAltsItem(current)
			current.WithSegmentEndToken(p.state.LA(0))
			current = result
		}
		current := current.(Alternatives)
		for ok := true; ok; ok = p.state.Lookahead(AssignableAlternativesLookahead12) == 0 {
			{
				token := p.state.Consume(Keyword_Pipe_Idx)
				core.AssignToken(current, token, AssignableAlternativesPipe_0)
			}
			{
				result := p.ParseAssignableWithoutAlts()
				if result != nil {
					current.WithAltsItem(result)
				}
			}
		}
	}
	current.WithSegmentEndToken(p.state.LA(0))
	return current
}

func (p *Parser) ParseCrossRef() CrossRef {
	current := NewCrossRef()
	current.WithSegmentStartToken(p.state.LA(1))
	{
		token := p.state.Consume(Keyword_LeftBracket_Idx)
		core.AssignToken(current, token, CrossRefLeftBracket_0)
	}
	{
		token := p.state.Consume(Token_ID_Idx)
		core.AssignToken(current, token, CrossRefTypeID_0)
		if token != nil {
			current.WithType(token)
		}
	}
	if p.state.Lookahead(CrossRefLookahead13) == 0 {
		{
			token := p.state.Consume(Keyword_Colon_Idx)
			core.AssignToken(current, token, CrossRefColon_0)
		}
		{
			result := p.ParseRuleCall()
			if result != nil {
				current.WithRule(result)
			}
		}
	}
	{
		token := p.state.Consume(Keyword_RightBracket_Idx)
		core.AssignToken(current, token, CrossRefRightBracket_0)
	}
	current.WithSegmentEndToken(p.state.LA(0))
	return current
}

func (p *Parser) ParseRuleCall() RuleCall {
	current := NewRuleCall()
	current.WithSegmentStartToken(p.state.LA(1))
	{
		token := p.state.Consume(Token_ID_Idx)
		core.AssignToken(current, token, RuleCallRuleID_0)
		if token != nil {
			current.WithRule(token)
		}
	}
	current.WithSegmentEndToken(p.state.LA(0))
	return current
}

func (p *Parser) ParseAction() Action {
	current := NewAction()
	current.WithSegmentStartToken(p.state.LA(1))
	{
		token := p.state.Consume(Keyword_LeftBrace_Idx)
		core.AssignToken(current, token, ActionLeftBrace_0)
	}
	{
		token := p.state.Consume(Token_ID_Idx)
		core.AssignToken(current, token, ActionTypeID_0)
		if token != nil {
			current.WithType(token)
		}
	}
	if p.state.Lookahead(ActionLookahead14) == 0 {
		{
			token := p.state.Consume(Keyword_Dot_Idx)
			core.AssignToken(current, token, ActionDot_0)
		}
		{
			token := p.state.Consume(Token_ID_Idx)
			core.AssignToken(current, token, ActionPropertyID_0)
			if token != nil {
				current.WithProperty(token)
			}
		}
		{
			switch p.state.Lookahead(ActionOperatorLookaheadOr6) {
			case 0:
				token := p.state.Consume(Keyword_PlusEquals_Idx)
				core.AssignToken(current, token, ActionOperatorPlusEquals_0)
				if token != nil {
					current.WithOperator(token)
				}
			case 1:
				token := p.state.Consume(Keyword_Equals_Idx)
				core.AssignToken(current, token, ActionOperatorEquals_0)
				if token != nil {
					current.WithOperator(token)
				}
			}
		}
		{
			token := p.state.Consume(Keyword_current_Idx)
			core.AssignToken(current, token, Actioncurrent_0)
		}
	}
	{
		token := p.state.Consume(Keyword_RightBrace_Idx)
		core.AssignToken(current, token, ActionRightBrace_0)
	}
	current.WithSegmentEndToken(p.state.LA(0))
	return current
}
