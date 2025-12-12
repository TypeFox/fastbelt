package generated

import (
	"strings"

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
			r := rune(input[index])
			switch state {
			case 0:
				if r == 34 { // '"',
					state = 1
				} else {
					break loop
				}
			case 1:
				if r >= 0 && r <= 33 || r >= 35 && r <= 1114111 { // '\u0000'..'!', '#'..'\u10FFFF',
					state = 2
				} else {
					break loop
				}
			case 2:
				if r >= 0 && r <= 33 || r >= 35 && r <= 1114111 { // '\u0000'..'!', '#'..'\u10FFFF',
					state = 2
				} else if r == 34 { // '"',
					state = 3
				} else {
					break loop
				}
			}
			if accepted[state] {
				acceptedIndex = index
			}
			index++
		}
		return acceptedIndex
	},
	[]rune{'"'},
)

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
			r := rune(input[index])
			switch state {
			case 0:
				if r >= 65 && r <= 90 || r == 95 || r >= 97 && r <= 122 { // 'A'..'Z', '_', 'a'..'z',
					state = 1
				} else {
					break loop
				}
			case 1:
				if r >= 48 && r <= 57 || r >= 65 && r <= 90 || r == 95 || r >= 97 && r <= 122 { // '0'..'9', 'A'..'Z', '_', 'a'..'z',
					state = 1
				} else {
					break loop
				}
			}
			if accepted[state] {
				acceptedIndex = index
			}
			index++
		}
		return acceptedIndex
	},
	[]rune{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '_', 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'},
)

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
			r := rune(input[index])
			switch state {
			case 0:
				if r == 47 { // '/',
					state = 1
				} else {
					break loop
				}
			case 1:
				if r >= 0 && r <= 9 || r >= 11 && r <= 12 || r >= 14 && r <= 46 || r >= 48 && r <= 90 || r >= 93 && r <= 1114111 { // '\u0000'..'\u0009', '\u000B'..'\u000C', '\u000E'..'.', '0'..'Z', ']'..'\u10FFFF',
					state = 2
				} else if r == 91 { // '[',
					state = 3
				} else if r == 92 { // '\\',
					state = 4
				} else {
					break loop
				}
			case 2:
				if r >= 0 && r <= 9 || r >= 11 && r <= 12 || r >= 14 && r <= 46 || r >= 48 && r <= 90 || r >= 93 && r <= 1114111 { // '\u0000'..'\u0009', '\u000B'..'\u000C', '\u000E'..'.', '0'..'Z', ']'..'\u10FFFF',
					state = 2
				} else if r == 47 { // '/',
					state = 5
				} else if r == 91 { // '[',
					state = 3
				} else if r == 92 { // '\\',
					state = 4
				} else {
					break loop
				}
			case 3:
				if r == 92 { // '\\',
					state = 6
				} else if r == 93 { // ']',
					state = 2
				} else if r >= 0 && r <= 9 || r >= 11 && r <= 12 || r >= 14 && r <= 91 || r >= 94 && r <= 1114111 { // '\u0000'..'\u0009', '\u000B'..'\u000C', '\u000E'..'[', '^'..'\u10FFFF',
					state = 3
				} else {
					break loop
				}
			case 4:
				if r >= 0 && r <= 9 || r >= 11 && r <= 1114111 { // '\u0000'..'\u0009', '\u000B'..'\u10FFFF',
					state = 2
				} else {
					break loop
				}
			case 6:
				if r >= 0 && r <= 9 || r >= 11 && r <= 1114111 { // '\u0000'..'\u0009', '\u000B'..'\u10FFFF',
					state = 3
				} else {
					break loop
				}
			}
			if accepted[state] {
				acceptedIndex = index
			}
			index++
		}
		return acceptedIndex
	},
	[]rune{'/'},
)

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
			r := rune(input[index])
			switch state {
			case 0:
				if r >= 9 && r <= 10 || r == 13 || r == 32 { // '\u0009'..'\u000A', '\u000D', ' ',
					state = 1
				} else {
					break loop
				}
			case 1:
				if r >= 9 && r <= 10 || r == 13 || r == 32 { // '\u0009'..'\u000A', '\u000D', ' ',
					state = 1
				} else {
					break loop
				}
			}
			if accepted[state] {
				acceptedIndex = index
			}
			index++
		}
		return acceptedIndex
	},
	[]rune{'\u0009', '\u000A', '\u000D', ' '},
)

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
