// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"go/format"
	"runtime"
	"sort"
	"strings"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/internal/grammar"
)

const CardinalityOne = ""
const CardinalityOptional = "?"
const CardinalityZeroOrMore = "*"
const CardinalityOneOrMore = "+"

func FormatIfPossible(text string) string {
	formatted, err := format.Source([]byte(text))
	if err != nil {
		return text
	}
	return string(formatted)
}

func eol() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

var EOL = eol()

const Indent = "    "

func EOLIndent(count int) string {
	return EOL + strings.Repeat(Indent, count)
}

func WriteSB(sb *strings.Builder, texts ...string) {
	for _, text := range texts {
		sb.WriteString(text)
	}
}

func GetAllKeywords(grammr grammar.Grammar) []grammar.Keyword {
	keywords := map[string]grammar.Keyword{}
	core.TraverseContent(grammr, func(node core.AstNode) {
		if keyword, ok := node.(grammar.Keyword); ok {
			keywords[keyword.Value()] = keyword
		}
	})
	return keysFromMap(keywords)
}

func keysFromMap(m map[string]grammar.Keyword) []grammar.Keyword {
	keys := []grammar.Keyword{}
	for _, v := range m {
		keys = append(keys, v)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].Value() < keys[j].Value()
	})
	return keys
}

func GeneratedTokenName(t grammar.Token) string {
	return "Token_" + t.Name()
}

func GeneratedTokenIdxName(t grammar.Token) string {
	return GeneratedTokenName(t) + "_Idx"
}

func KeywordValue(k grammar.Keyword) string {
	return k.Value()[1 : len(k.Value())-1]
}

func KeywordName(k grammar.Keyword) string {
	sb := &strings.Builder{}
	for _, r := range KeywordValue(k) {
		switch r {
		case '(':
			sb.WriteString("LeftParen")
		case ')':
			sb.WriteString("RightParen")
		case '{':
			sb.WriteString("LeftBrace")
		case '}':
			sb.WriteString("RightBrace")
		case '[':
			sb.WriteString("LeftBracket")
		case ']':
			sb.WriteString("RightBracket")
		case '_':
			sb.WriteString("Underscore")
		case '$':
			sb.WriteString("Dollar")
		case '%':
			sb.WriteString("Percent")
		case '#':
			sb.WriteString("Hash")
		case '@':
			sb.WriteString("At")
		case '!':
			sb.WriteString("Exclamation")
		case '^':
			sb.WriteString("Caret")
		case '&':
			sb.WriteString("Ampersand")
		case '*':
			sb.WriteString("Asterisk")
		case '-':
			sb.WriteString("Dash")
		case '+':
			sb.WriteString("Plus")
		case '=':
			sb.WriteString("Equals")
		case '<':
			sb.WriteString("LessThan")
		case '>':
			sb.WriteString("GreaterThan")
		case '?':
			sb.WriteString("Question")
		case '/':
			sb.WriteString("Slash")
		case '\\':
			sb.WriteString("Backslash")
		case '|':
			sb.WriteString("Pipe")
		case '~':
			sb.WriteString("Tilde")
		case '`':
			sb.WriteString("Backtick")
		case '.':
			sb.WriteString("Dot")
		case ',':
			sb.WriteString("Comma")
		case ':':
			sb.WriteString("Colon")
		case ';':
			sb.WriteString("Semicolon")
		case ' ':
			sb.WriteString("Space")
		case '\t':
			sb.WriteString("Tab")
		case '\n':
			sb.WriteString("Newline")
		case '\r':
			sb.WriteString("CarriageReturn")
		default:
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func GeneratedKeywordName(k grammar.Keyword) string {
	return "Keyword_" + KeywordName(k)
}

func GeneratedKeywordIdxName(k grammar.Keyword) string {
	return GeneratedKeywordName(k) + "_Idx"
}
