package generated

import (
	"slices"
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

const Keyword_bool_Idx = 15

var Keyword_bool = core.NewTokenType(
	Keyword_bool_Idx,
	"bool",
	"bool",
	0,
	0,
	false,
	func(text string, offset int) int {
		if strings.HasPrefix(text[offset:], "bool") {
			return 4
		}
		return 0
	},
	[]rune{'b'},
)

const Keyword_current_Idx = 16

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

const Keyword_extends_Idx = 17

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

const Keyword_grammar_Idx = 18

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

const Keyword_hidden_Idx = 19

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

const Keyword_interface_Idx = 20

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

const Keyword_returns_Idx = 21

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

const Keyword_string_Idx = 22

var Keyword_string = core.NewTokenType(
	Keyword_string_Idx,
	"string",
	"string",
	0,
	0,
	false,
	func(text string, offset int) int {
		if strings.HasPrefix(text[offset:], "string") {
			return 6
		}
		return 0
	},
	[]rune{'s'},
)

const Keyword_token_Idx = 23

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

const Keyword_LeftBrace_Idx = 24

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

const Keyword_Pipe_Idx = 25

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

const Keyword_RightBrace_Idx = 26

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

const Token_String_Idx = 27

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
		acceptedIndex := 0
		index := 0
	loop:
		for index < length {
			r, runeSize := utf8.DecodeRuneInString(input[index:])
			switch state {
			case 0:
				nextState := -1
				next := Token_String_Next[0]
				lookup := Token_String_Lookup[0]
				searchIndex := slices.IndexFunc(lookup, func(lowHigh int64) bool {
					lo := rune(lowHigh & 0xFFFFFFFF)
					hi := rune(lowHigh >> 32)
					return lo <= r && r <= hi
				})
				if searchIndex > -1 {
					nextState = next[searchIndex]
				}
				if nextState > -1 {
					state = nextState
				} else {
					break loop
				}
			case 1:
				nextState := -1
				next := Token_String_Next[1]
				lookup := Token_String_Lookup[1]
				searchIndex := slices.IndexFunc(lookup, func(lowHigh int64) bool {
					lo := rune(lowHigh & 0xFFFFFFFF)
					hi := rune(lowHigh >> 32)
					return lo <= r && r <= hi
				})
				if searchIndex > -1 {
					nextState = next[searchIndex]
				}
				if nextState > -1 {
					state = nextState
				} else {
					break loop
				}
			case 2:
				nextState := -1
				next := Token_String_Next[2]
				lookup := Token_String_Lookup[2]
				searchIndex := slices.IndexFunc(lookup, func(lowHigh int64) bool {
					lo := rune(lowHigh & 0xFFFFFFFF)
					hi := rune(lowHigh >> 32)
					return lo <= r && r <= hi
				})
				if searchIndex > -1 {
					nextState = next[searchIndex]
				}
				if nextState > -1 {
					state = nextState
				} else {
					break loop
				}
			case 3:
				nextState := -1
				next := Token_String_Next[3]
				lookup := Token_String_Lookup[3]
				searchIndex := slices.IndexFunc(lookup, func(lowHigh int64) bool {
					lo := rune(lowHigh & 0xFFFFFFFF)
					hi := rune(lowHigh >> 32)
					return lo <= r && r <= hi
				})
				if searchIndex > -1 {
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
	[]rune{'"'},
)
var Token_String_Lookup = [][]int64{
	{0x0000002200000022},
	{0x0000002100000000, 0x0010FFFF00000023},
	{0x0000002100000000, 0x0000002200000022, 0x0010FFFF00000023},
	{},
}
var Token_String_Next = [][]int{
	{1},
	{2, 2},
	{2, 3, 2},
	{},
}

const Token_ID_Idx = 28

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
		acceptedIndex := 0
		index := 0
	loop:
		for index < length {
			r, runeSize := utf8.DecodeRuneInString(input[index:])
			switch state {
			case 0:
				nextState := -1
				next := Token_ID_Next[0]
				lookup := Token_ID_Lookup[0]
				searchIndex := slices.IndexFunc(lookup, func(lowHigh int64) bool {
					lo := rune(lowHigh & 0xFFFFFFFF)
					hi := rune(lowHigh >> 32)
					return lo <= r && r <= hi
				})
				if searchIndex > -1 {
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
				searchIndex := slices.IndexFunc(lookup, func(lowHigh int64) bool {
					lo := rune(lowHigh & 0xFFFFFFFF)
					hi := rune(lowHigh >> 32)
					return lo <= r && r <= hi
				})
				if searchIndex > -1 {
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
var Token_ID_Lookup = [][]int64{
	{0x0000005A00000041, 0x0000005F0000005F, 0x0000007A00000061},
	{0x0000003900000030, 0x0000005A00000041, 0x0000005F0000005F, 0x0000007A00000061},
}
var Token_ID_Next = [][]int{
	{1, 1, 1},
	{1, 1, 1, 1},
}

const Token_RegexLiteral_Idx = 29

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
		acceptedIndex := 0
		index := 0
	loop:
		for index < length {
			r, runeSize := utf8.DecodeRuneInString(input[index:])
			switch state {
			case 0:
				nextState := -1
				next := Token_RegexLiteral_Next[0]
				lookup := Token_RegexLiteral_Lookup[0]
				searchIndex := slices.IndexFunc(lookup, func(lowHigh int64) bool {
					lo := rune(lowHigh & 0xFFFFFFFF)
					hi := rune(lowHigh >> 32)
					return lo <= r && r <= hi
				})
				if searchIndex > -1 {
					nextState = next[searchIndex]
				}
				if nextState > -1 {
					state = nextState
				} else {
					break loop
				}
			case 1:
				nextState := -1
				next := Token_RegexLiteral_Next[1]
				lookup := Token_RegexLiteral_Lookup[1]
				searchIndex := slices.IndexFunc(lookup, func(lowHigh int64) bool {
					lo := rune(lowHigh & 0xFFFFFFFF)
					hi := rune(lowHigh >> 32)
					return lo <= r && r <= hi
				})
				if searchIndex > -1 {
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
				searchIndex := slices.IndexFunc(lookup, func(lowHigh int64) bool {
					lo := rune(lowHigh & 0xFFFFFFFF)
					hi := rune(lowHigh >> 32)
					return lo <= r && r <= hi
				})
				if searchIndex > -1 {
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
				searchIndex := slices.IndexFunc(lookup, func(lowHigh int64) bool {
					lo := rune(lowHigh & 0xFFFFFFFF)
					hi := rune(lowHigh >> 32)
					return lo <= r && r <= hi
				})
				if searchIndex > -1 {
					nextState = next[searchIndex]
				}
				if nextState > -1 {
					state = nextState
				} else {
					break loop
				}
			case 4:
				nextState := -1
				next := Token_RegexLiteral_Next[4]
				lookup := Token_RegexLiteral_Lookup[4]
				searchIndex := slices.IndexFunc(lookup, func(lowHigh int64) bool {
					lo := rune(lowHigh & 0xFFFFFFFF)
					hi := rune(lowHigh >> 32)
					return lo <= r && r <= hi
				})
				if searchIndex > -1 {
					nextState = next[searchIndex]
				}
				if nextState > -1 {
					state = nextState
				} else {
					break loop
				}
			case 5:
				nextState := -1
				next := Token_RegexLiteral_Next[5]
				lookup := Token_RegexLiteral_Lookup[5]
				searchIndex := slices.IndexFunc(lookup, func(lowHigh int64) bool {
					lo := rune(lowHigh & 0xFFFFFFFF)
					hi := rune(lowHigh >> 32)
					return lo <= r && r <= hi
				})
				if searchIndex > -1 {
					nextState = next[searchIndex]
				}
				if nextState > -1 {
					state = nextState
				} else {
					break loop
				}
			case 6:
				nextState := -1
				next := Token_RegexLiteral_Next[6]
				lookup := Token_RegexLiteral_Lookup[6]
				searchIndex := slices.IndexFunc(lookup, func(lowHigh int64) bool {
					lo := rune(lowHigh & 0xFFFFFFFF)
					hi := rune(lowHigh >> 32)
					return lo <= r && r <= hi
				})
				if searchIndex > -1 {
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
	[]rune{'/'},
)
var Token_RegexLiteral_Lookup = [][]int64{
	{0x0000002F0000002F},
	{0x0000000900000000, 0x0000000C0000000B, 0x0000002E0000000E, 0x0000005A00000030, 0x0000005B0000005B, 0x0000005C0000005C, 0x0010FFFF0000005D},
	{0x0000000900000000, 0x0000000C0000000B, 0x0000002E0000000E, 0x0000002F0000002F, 0x0000005A00000030, 0x0000005B0000005B, 0x0000005C0000005C, 0x0010FFFF0000005D},
	{0x0000000900000000, 0x0000000C0000000B, 0x0000005B0000000E, 0x0000005C0000005C, 0x0000005D0000005D, 0x0010FFFF0000005E},
	{0x0000000900000000, 0x0010FFFF0000000B},
	{},
	{0x0000000900000000, 0x0010FFFF0000000B},
}
var Token_RegexLiteral_Next = [][]int{
	{1},
	{2, 2, 2, 2, 3, 4, 2},
	{2, 2, 2, 5, 2, 3, 4, 2},
	{3, 3, 3, 6, 2, 3},
	{2, 2},
	{},
	{3, 3},
}

const Token_WS_Idx = 30

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
		acceptedIndex := 0
		index := 0
	loop:
		for index < length {
			r, runeSize := utf8.DecodeRuneInString(input[index:])
			switch state {
			case 0:
				nextState := -1
				next := Token_WS_Next[0]
				lookup := Token_WS_Lookup[0]
				searchIndex := slices.IndexFunc(lookup, func(lowHigh int64) bool {
					lo := rune(lowHigh & 0xFFFFFFFF)
					hi := rune(lowHigh >> 32)
					return lo <= r && r <= hi
				})
				if searchIndex > -1 {
					nextState = next[searchIndex]
				}
				if nextState > -1 {
					state = nextState
				} else {
					break loop
				}
			case 1:
				nextState := -1
				next := Token_WS_Next[1]
				lookup := Token_WS_Lookup[1]
				searchIndex := slices.IndexFunc(lookup, func(lowHigh int64) bool {
					lo := rune(lowHigh & 0xFFFFFFFF)
					hi := rune(lowHigh >> 32)
					return lo <= r && r <= hi
				})
				if searchIndex > -1 {
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
	[]rune{'\u0009', '\u000A', '\u000D', ' '},
)
var Token_WS_Lookup = [][]int64{
	{0x0000000A00000009, 0x0000000D0000000D, 0x0000002000000020},
	{0x0000000A00000009, 0x0000000D0000000D, 0x0000002000000020},
}
var Token_WS_Next = [][]int{
	{1, 1, 1},
	{1, 1, 1},
}

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
		Keyword_bool,
		Keyword_current,
		Keyword_extends,
		Keyword_grammar,
		Keyword_hidden,
		Keyword_interface,
		Keyword_returns,
		Keyword_string,
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
