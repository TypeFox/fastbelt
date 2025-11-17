package generated

import (
	"regexp"
	"strings"

	"github.com/TypeFox/langium-to-go/core"
	"github.com/TypeFox/langium-to-go/lexer"
)

const Keyword_LeftParen_Idx = 1

var Keyword_LeftParen = core.NewTokenType(
	Keyword_LeftParen_Idx,
	"(",
	"(",
	0,
	0,
	false,
	func(text string, offset int) int {
		if strings.HasPrefix(text[offset:], "(") {
			return 1
		}
		return 0
	},
	[]rune{
		'(',
	},
)

const Keyword_RightParen_Idx = 2

var Keyword_RightParen = core.NewTokenType(
	Keyword_RightParen_Idx,
	")",
	")",
	0,
	0,
	false,
	func(text string, offset int) int {
		if strings.HasPrefix(text[offset:], ")") {
			return 1
		}
		return 0
	},
	[]rune{
		')',
	},
)

const Keyword_Comma_Idx = 3

var Keyword_Comma = core.NewTokenType(
	Keyword_Comma_Idx,
	",",
	",",
	0,
	0,
	false,
	func(text string, offset int) int {
		if strings.HasPrefix(text[offset:], ",") {
			return 1
		}
		return 0
	},
	[]rune{
		',',
	},
)

const Keyword_Dot_Idx = 4

var Keyword_Dot = core.NewTokenType(
	Keyword_Dot_Idx,
	".",
	".",
	0,
	0,
	false,
	func(text string, offset int) int {
		if strings.HasPrefix(text[offset:], ".") {
			return 1
		}
		return 0
	},
	[]rune{
		'.',
	},
)

const Keyword_Colon_Idx = 5

var Keyword_Colon = core.NewTokenType(
	Keyword_Colon_Idx,
	":",
	":",
	0,
	0,
	false,
	func(text string, offset int) int {
		if strings.HasPrefix(text[offset:], ":") {
			return 1
		}
		return 0
	},
	[]rune{
		':',
	},
)

const Keyword_Semicolon_Idx = 6

var Keyword_Semicolon = core.NewTokenType(
	Keyword_Semicolon_Idx,
	";",
	";",
	0,
	0,
	false,
	func(text string, offset int) int {
		if strings.HasPrefix(text[offset:], ";") {
			return 1
		}
		return 0
	},
	[]rune{
		';',
	},
)

const Keyword_LeftBracket_Idx = 7

var Keyword_LeftBracket = core.NewTokenType(
	Keyword_LeftBracket_Idx,
	"[",
	"[",
	0,
	0,
	false,
	func(text string, offset int) int {
		if strings.HasPrefix(text[offset:], "[") {
			return 1
		}
		return 0
	},
	[]rune{
		'[',
	},
)

const Keyword_RightBracket_Idx = 8

var Keyword_RightBracket = core.NewTokenType(
	Keyword_RightBracket_Idx,
	"]",
	"]",
	0,
	0,
	false,
	func(text string, offset int) int {
		if strings.HasPrefix(text[offset:], "]") {
			return 1
		}
		return 0
	},
	[]rune{
		']',
	},
)

const Keyword_current_Idx = 9

var Keyword_current = core.NewTokenType(
	Keyword_current_Idx,
	"current",
	"current",
	0,
	0,
	false,
	func(text string, offset int) int {
		if strings.HasPrefix(text[offset:], "current") {
			return 7
		}
		return 0
	},
	[]rune{
		'c',
	},
)

const Keyword_extends_Idx = 10

var Keyword_extends = core.NewTokenType(
	Keyword_extends_Idx,
	"extends",
	"extends",
	0,
	0,
	false,
	func(text string, offset int) int {
		if strings.HasPrefix(text[offset:], "extends") {
			return 7
		}
		return 0
	},
	[]rune{
		'e',
	},
)

const Keyword_grammar_Idx = 11

var Keyword_grammar = core.NewTokenType(
	Keyword_grammar_Idx,
	"grammar",
	"grammar",
	0,
	0,
	false,
	func(text string, offset int) int {
		if strings.HasPrefix(text[offset:], "grammar") {
			return 7
		}
		return 0
	},
	[]rune{
		'g',
	},
)

const Keyword_hidden_Idx = 12

var Keyword_hidden = core.NewTokenType(
	Keyword_hidden_Idx,
	"hidden",
	"hidden",
	0,
	0,
	false,
	func(text string, offset int) int {
		if strings.HasPrefix(text[offset:], "hidden") {
			return 6
		}
		return 0
	},
	[]rune{
		'h',
	},
)

const Keyword_interface_Idx = 13

var Keyword_interface = core.NewTokenType(
	Keyword_interface_Idx,
	"interface",
	"interface",
	0,
	0,
	false,
	func(text string, offset int) int {
		if strings.HasPrefix(text[offset:], "interface") {
			return 9
		}
		return 0
	},
	[]rune{
		'i',
	},
)

const Keyword_returns_Idx = 14

var Keyword_returns = core.NewTokenType(
	Keyword_returns_Idx,
	"returns",
	"returns",
	0,
	0,
	false,
	func(text string, offset int) int {
		if strings.HasPrefix(text[offset:], "returns") {
			return 7
		}
		return 0
	},
	[]rune{
		'r',
	},
)

const Keyword_token_Idx = 15

var Keyword_token = core.NewTokenType(
	Keyword_token_Idx,
	"token",
	"token",
	0,
	0,
	false,
	func(text string, offset int) int {
		if strings.HasPrefix(text[offset:], "token") {
			return 5
		}
		return 0
	},
	[]rune{
		't',
	},
)

const Keyword_LeftBrace_Idx = 16

var Keyword_LeftBrace = core.NewTokenType(
	Keyword_LeftBrace_Idx,
	"{",
	"{",
	0,
	0,
	false,
	func(text string, offset int) int {
		if strings.HasPrefix(text[offset:], "{") {
			return 1
		}
		return 0
	},
	[]rune{
		'{',
	},
)

const Keyword_Pipe_Idx = 17

var Keyword_Pipe = core.NewTokenType(
	Keyword_Pipe_Idx,
	"|",
	"|",
	0,
	0,
	false,
	func(text string, offset int) int {
		if strings.HasPrefix(text[offset:], "|") {
			return 1
		}
		return 0
	},
	[]rune{
		'|',
	},
)

const Keyword_RightBrace_Idx = 18

var Keyword_RightBrace = core.NewTokenType(
	Keyword_RightBrace_Idx,
	"}",
	"}",
	0,
	0,
	false,
	func(text string, offset int) int {
		if strings.HasPrefix(text[offset:], "}") {
			return 1
		}
		return 0
	},
	[]rune{
		'}',
	},
)

const Token_AssignmentOperator_Idx = 19

var Token_AssignmentOperator_Regexp = regexp.MustCompile(`^(=|\+=|\?=)`)
var Token_AssignmentOperator = core.NewTokenType(
	Token_AssignmentOperator_Idx,
	"AssignmentOperator",
	"AssignmentOperator",
	0,
	0,
	false,
	func(text string, offset int) int {
		matches := Token_AssignmentOperator_Regexp.FindStringIndex(text[offset:])
		if matches != nil {
			return matches[1]
		}
		return 0
	},
	[]rune{
		'+', '=', '?',
	},
)

const Token_Cardinality_Idx = 20

var Token_Cardinality_Regexp = regexp.MustCompile(`^([\*\+\?])`)
var Token_Cardinality = core.NewTokenType(
	Token_Cardinality_Idx,
	"Cardinality",
	"Cardinality",
	0,
	0,
	false,
	func(text string, offset int) int {
		matches := Token_Cardinality_Regexp.FindStringIndex(text[offset:])
		if matches != nil {
			return matches[1]
		}
		return 0
	},
	[]rune{
		'*', '+', '?',
	},
)

const Token_String_Idx = 21

var Token_String_Regexp = regexp.MustCompile(`^("[^"]+")`)
var Token_String = core.NewTokenType(
	Token_String_Idx,
	"String",
	"String",
	0,
	0,
	false,
	func(text string, offset int) int {
		matches := Token_String_Regexp.FindStringIndex(text[offset:])
		if matches != nil {
			return matches[1]
		}
		return 0
	},
	[]rune{
		'"',
	},
)

const Token_ID_Idx = 22

var Token_ID_Regexp = regexp.MustCompile(`^([A-Z_a-z][0-9A-Z_a-z]*)`)
var Token_ID = core.NewTokenType(
	Token_ID_Idx,
	"ID",
	"ID",
	0,
	0,
	false,
	func(text string, offset int) int {
		matches := Token_ID_Regexp.FindStringIndex(text[offset:])
		if matches != nil {
			return matches[1]
		}
		return 0
	},
	[]rune{
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
		'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T',
		'U', 'V', 'W', 'X', 'Y', 'Z', '_', 'a', 'b', 'c',
		'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
		'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w',
		'x', 'y', 'z',
	},
)

const Token_RegexLiteral_Idx = 23

var Token_RegexLiteral_Regexp = regexp.MustCompile(`^((?-s:/([^\n\r/\[\\]|\\.|\[([^\n\r\\\]]|\\.)*\])+/))`)
var Token_RegexLiteral = core.NewTokenType(
	Token_RegexLiteral_Idx,
	"RegexLiteral",
	"RegexLiteral",
	0,
	0,
	false,
	func(text string, offset int) int {
		matches := Token_RegexLiteral_Regexp.FindStringIndex(text[offset:])
		if matches != nil {
			return matches[1]
		}
		return 0
	},
	[]rune{
		'/',
	},
)

const Token_WS_Idx = 24

var Token_WS_Regexp = regexp.MustCompile(`^([\t\n\r ]+)`)
var Token_WS = core.NewTokenType(
	Token_WS_Idx,
	"WS",
	"WS",
	-1,
	0,
	false,
	func(text string, offset int) int {
		matches := Token_WS_Regexp.FindStringIndex(text[offset:])
		if matches != nil {
			return matches[1]
		}
		return 0
	},
	[]rune{
		9, 10, 13, ' ',
	},
)

func NewLexer() lexer.Lexer {
	return lexer.NewLexer(
		Keyword_LeftParen,
		Keyword_RightParen,
		Keyword_Comma,
		Keyword_Dot,
		Keyword_Colon,
		Keyword_Semicolon,
		Keyword_LeftBracket,
		Keyword_RightBracket,
		Keyword_current,
		Keyword_extends,
		Keyword_grammar,
		Keyword_hidden,
		Keyword_interface,
		Keyword_returns,
		Keyword_token,
		Keyword_LeftBrace,
		Keyword_Pipe,
		Keyword_RightBrace,
		Token_AssignmentOperator,
		Token_Cardinality,
		Token_String,
		Token_ID,
		Token_RegexLiteral,
		Token_WS,
	)
}
