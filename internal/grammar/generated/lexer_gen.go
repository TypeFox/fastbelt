package generated

import (
	"sort"
	"strings"
	"unicode/utf8"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/lexer"
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
	[]rune{'('},
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
	[]rune{')'},
)

const Keyword_Asterisk_Idx = 3

var Keyword_Asterisk = core.NewTokenType(
	Keyword_Asterisk_Idx,
	"*",
	"*",
	0,
	0,
	false,
	func(text string, offset int) int {
		if strings.HasPrefix(text[offset:], "*") {
			return 1
		}
		return 0
	},
	[]rune{'*'},
)

const Keyword_Plus_Idx = 4

var Keyword_Plus = core.NewTokenType(
	Keyword_Plus_Idx,
	"+",
	"+",
	0,
	0,
	false,
	func(text string, offset int) int {
		if strings.HasPrefix(text[offset:], "+") {
			return 1
		}
		return 0
	},
	[]rune{'+'},
)

const Keyword_PlusEquals_Idx = 5

var Keyword_PlusEquals = core.NewTokenType(
	Keyword_PlusEquals_Idx,
	"+=",
	"+=",
	0,
	0,
	false,
	func(text string, offset int) int {
		if strings.HasPrefix(text[offset:], "+=") {
			return 2
		}
		return 0
	},
	[]rune{'+'},
)

const Keyword_Comma_Idx = 6

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
	[]rune{','},
)

const Keyword_Dot_Idx = 7

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
	[]rune{'.'},
)

const Keyword_Colon_Idx = 8

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
	[]rune{':'},
)

const Keyword_Semicolon_Idx = 9

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
	[]rune{';'},
)

const Keyword_Equals_Idx = 10

var Keyword_Equals = core.NewTokenType(
	Keyword_Equals_Idx,
	"=",
	"=",
	0,
	0,
	false,
	func(text string, offset int) int {
		if strings.HasPrefix(text[offset:], "=") {
			return 1
		}
		return 0
	},
	[]rune{'='},
)

const Keyword_Question_Idx = 11

var Keyword_Question = core.NewTokenType(
	Keyword_Question_Idx,
	"?",
	"?",
	0,
	0,
	false,
	func(text string, offset int) int {
		if strings.HasPrefix(text[offset:], "?") {
			return 1
		}
		return 0
	},
	[]rune{'?'},
)

const Keyword_QuestionEquals_Idx = 12

var Keyword_QuestionEquals = core.NewTokenType(
	Keyword_QuestionEquals_Idx,
	"?=",
	"?=",
	0,
	0,
	false,
	func(text string, offset int) int {
		if strings.HasPrefix(text[offset:], "?=") {
			return 2
		}
		return 0
	},
	[]rune{'?'},
)

const Keyword_LeftBracket_Idx = 13

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
	[]rune{'['},
)

const Keyword_RightBracket_Idx = 14

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
	[]rune{']'},
)

const Keyword_current_Idx = 15

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
	[]rune{'c'},
)

const Keyword_extends_Idx = 16

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
	[]rune{'e'},
)

const Keyword_grammar_Idx = 17

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
	[]rune{'g'},
)

const Keyword_hidden_Idx = 18

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
	[]rune{'h'},
)

const Keyword_interface_Idx = 19

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
	[]rune{'i'},
)

const Keyword_returns_Idx = 20

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
	[]rune{'r'},
)

const Keyword_token_Idx = 21

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
	[]rune{'t'},
)

const Keyword_LeftBrace_Idx = 22

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
	[]rune{'{'},
)

const Keyword_Pipe_Idx = 23

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
	[]rune{'|'},
)

const Keyword_RightBrace_Idx = 24

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
	[]rune{'}'},
)

const Token_String_Idx = 25

var Token_String = core.NewTokenType(
	Token_String_Idx,
	"String",
	"String",
	0,
	0,
	false,
	func(s string, offset int) int {
		input := s[offset:]
		length := len(input)
		accepted := map[int]bool{3: true}
		state := 0
		acceptedIndex := -1
		if accepted[state] {
			acceptedIndex = 0
		}
		index := 0
	loop:
		for index < length {
			r, runeSize := utf8.DecodeRuneInString(input[index:])
			switch state {
			case 0:
				if r == 0x22 { // '"',
					state = 1
				} else {
					break loop
				}
			case 1:
				if r >= 0x00 && r <= 0x21 || r >= 0x23 && r <= 0x10FFFF { // '\u0000'..'!', '#'..'\U0010FFFF',
					state = 2
				} else {
					break loop
				}
			case 2:
				if r >= 0x00 && r <= 0x21 || r >= 0x23 && r <= 0x10FFFF { // '\u0000'..'!', '#'..'\U0010FFFF',
					state = 2
				} else if r == 0x22 { // '"',
					state = 3
				} else {
					break loop
				}
			default:
				break loop
			}
			index += runeSize
			if accepted[state] {
				acceptedIndex = index
			}
		}
		return acceptedIndex
	},
	[]rune{'"'},
)
var Token_String_Lookup = map[int][]rune{}
var Token_String_Next = map[int][]int{}

const Token_ID_Idx = 26

var Token_ID = core.NewTokenType(
	Token_ID_Idx,
	"ID",
	"ID",
	0,
	0,
	false,
	func(s string, offset int) int {
		input := s[offset:]
		length := len(input)
		accepted := map[int]bool{1: true}
		state := 0
		acceptedIndex := -1
		if accepted[state] {
			acceptedIndex = 0
		}
		index := 0
	loop:
		for index < length {
			r, runeSize := utf8.DecodeRuneInString(input[index:])
			switch state {
			case 0:
				nextState := -1
				next := Token_ID_Next[0]
				lookup := Token_ID_Lookup[0]
				searchIndex := sort.Search(len(next), func(i int) bool {
					return lookup[i*2] > r
				}) - 1
				if searchIndex > -1 && lookup[searchIndex*2] <= r && r <= lookup[searchIndex*2+1] {
					nextState = next[searchIndex]
				}
				if nextState > -1 {
					state = nextState
				} else {
					break loop
				}
			case 1:
				nextState := -1
				next := Token_ID_Next[1]
				lookup := Token_ID_Lookup[1]
				searchIndex := sort.Search(len(next), func(i int) bool {
					return lookup[i*2] > r
				}) - 1
				if searchIndex > -1 && lookup[searchIndex*2] <= r && r <= lookup[searchIndex*2+1] {
					nextState = next[searchIndex]
				}
				if nextState > -1 {
					state = nextState
				} else {
					break loop
				}
			default:
				break loop
			}
			index += runeSize
			if accepted[state] {
				acceptedIndex = index
			}
		}
		return acceptedIndex
	},
	[]rune{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '_', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'},
)
var Token_ID_Lookup = map[int][]rune{
	0: {0x41, 0x5A, 0x5F, 0x5F, 0x61, 0x7A},
	1: {0x30, 0x39, 0x41, 0x5A, 0x5F, 0x5F, 0x61, 0x7A},
}
var Token_ID_Next = map[int][]int{
	0: {1, 1, 1},
	1: {1, 1, 1, 1},
}

const Token_RegexLiteral_Idx = 27

var Token_RegexLiteral = core.NewTokenType(
	Token_RegexLiteral_Idx,
	"RegexLiteral",
	"RegexLiteral",
	0,
	0,
	false,
	func(s string, offset int) int {
		input := s[offset:]
		length := len(input)
		accepted := map[int]bool{5: true}
		state := 0
		acceptedIndex := -1
		if accepted[state] {
			acceptedIndex = 0
		}
		index := 0
	loop:
		for index < length {
			r, runeSize := utf8.DecodeRuneInString(input[index:])
			switch state {
			case 0:
				if r == 0x2F { // '/',
					state = 1
				} else {
					break loop
				}
			case 1:
				nextState := -1
				next := Token_RegexLiteral_Next[1]
				lookup := Token_RegexLiteral_Lookup[1]
				searchIndex := sort.Search(len(next), func(i int) bool {
					return lookup[i*2] > r
				}) - 1
				if searchIndex > -1 && lookup[searchIndex*2] <= r && r <= lookup[searchIndex*2+1] {
					nextState = next[searchIndex]
				}
				if nextState > -1 {
					state = nextState
				} else {
					break loop
				}
			case 2:
				nextState := -1
				next := Token_RegexLiteral_Next[2]
				lookup := Token_RegexLiteral_Lookup[2]
				searchIndex := sort.Search(len(next), func(i int) bool {
					return lookup[i*2] > r
				}) - 1
				if searchIndex > -1 && lookup[searchIndex*2] <= r && r <= lookup[searchIndex*2+1] {
					nextState = next[searchIndex]
				}
				if nextState > -1 {
					state = nextState
				} else {
					break loop
				}
			case 3:
				nextState := -1
				next := Token_RegexLiteral_Next[3]
				lookup := Token_RegexLiteral_Lookup[3]
				searchIndex := sort.Search(len(next), func(i int) bool {
					return lookup[i*2] > r
				}) - 1
				if searchIndex > -1 && lookup[searchIndex*2] <= r && r <= lookup[searchIndex*2+1] {
					nextState = next[searchIndex]
				}
				if nextState > -1 {
					state = nextState
				} else {
					break loop
				}
			case 4:
				if r >= 0x00 && r <= 0x09 || r >= 0x0B && r <= 0x10FFFF { // '\u0000'..'\u0009', '\u000B'..'\U0010FFFF',
					state = 2
				} else {
					break loop
				}
			case 6:
				if r >= 0x00 && r <= 0x09 || r >= 0x0B && r <= 0x10FFFF { // '\u0000'..'\u0009', '\u000B'..'\U0010FFFF',
					state = 3
				} else {
					break loop
				}
			default:
				break loop
			}
			index += runeSize
			if accepted[state] {
				acceptedIndex = index
			}
		}
		return acceptedIndex
	},
	[]rune{'/'},
)
var Token_RegexLiteral_Lookup = map[int][]rune{
	1: {0x00, 0x09, 0x0B, 0x0C, 0x0E, 0x2E, 0x30, 0x5A, 0x5B, 0x5B, 0x5C, 0x5C, 0x5D, 0x10FFFF},
	2: {0x00, 0x09, 0x0B, 0x0C, 0x0E, 0x2E, 0x2F, 0x2F, 0x30, 0x5A, 0x5B, 0x5B, 0x5C, 0x5C, 0x5D, 0x10FFFF},
	3: {0x00, 0x09, 0x0B, 0x0C, 0x0E, 0x5B, 0x5C, 0x5C, 0x5D, 0x5D, 0x5E, 0x10FFFF},
}
var Token_RegexLiteral_Next = map[int][]int{
	1: {2, 2, 2, 2, 3, 4, 2},
	2: {2, 2, 2, 5, 2, 3, 4, 2},
	3: {3, 3, 3, 6, 2, 3},
}

const Token_WS_Idx = 28

var Token_WS = core.NewTokenType(
	Token_WS_Idx,
	"WS",
	"WS",
	-1,
	0,
	false,
	func(s string, offset int) int {
		input := s[offset:]
		length := len(input)
		accepted := map[int]bool{1: true}
		state := 0
		acceptedIndex := -1
		if accepted[state] {
			acceptedIndex = 0
		}
		index := 0
	loop:
		for index < length {
			r, runeSize := utf8.DecodeRuneInString(input[index:])
			switch state {
			case 0:
				switch r {
				case 0x09:
					fallthrough
				case 0x0A:
					fallthrough
				case 0x0D:
					fallthrough
				case 0x20:
					state = 1
				default:
					break loop
				}
			case 1:
				switch r {
				case 0x09:
					fallthrough
				case 0x0A:
					fallthrough
				case 0x0D:
					fallthrough
				case 0x20:
					state = 1
				default:
					break loop
				}
			default:
				break loop
			}
			index += runeSize
			if accepted[state] {
				acceptedIndex = index
			}
		}
		return acceptedIndex
	},
	[]rune{'\u0009', '\u000A', '\u000D', ' '},
)
var Token_WS_Lookup = map[int][]rune{}
var Token_WS_Next = map[int][]int{}

func NewLexer() lexer.Lexer {
	return lexer.NewDefaultLexer(
		Keyword_LeftParen,
		Keyword_RightParen,
		Keyword_Asterisk,
		Keyword_Plus,
		Keyword_PlusEquals,
		Keyword_Comma,
		Keyword_Dot,
		Keyword_Colon,
		Keyword_Semicolon,
		Keyword_Equals,
		Keyword_Question,
		Keyword_QuestionEquals,
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
		Token_String,
		Token_ID,
		Token_RegexLiteral,
		Token_WS,
	)
}
