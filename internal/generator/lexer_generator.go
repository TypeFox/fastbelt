// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"fmt"
	"regexp/syntax"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/TypeFox/langium-to-go/internal/generated"
)

func GenerateLexer(grammar generated.Grammar) string {
	sb := &strings.Builder{}
	WriteSB(
		sb,
		"package generated",
		EOL,
		EOL,
		"import (",
		EOLIndent(1),
		"\"regexp\"",
		EOLIndent(1),
		"\"strings\"",
		EOL,
		EOLIndent(1),
		"\"github.com/TypeFox/langium-to-go/core\"",
		EOLIndent(1),
		"\"github.com/TypeFox/langium-to-go/lexer\"",
		EOL,
		")",
		EOL,
		EOL,
	)
	// keywords := GetAllKeywords(grammar)
	tokens := grammar.Terminals()
	keywords := GetAllKeywords(grammar)
	id := 1
	for _, keyword := range keywords {
		generateKeywordTokenType(sb, keyword, id)
		id++
	}
	for _, token := range tokens {
		generateTokenType(sb, token, id)
		id++
	}
	generateMainLexerFunction(sb, tokens, keywords)
	return formatIfPossible(sb.String())
}

func generateMainLexerFunction(sb *strings.Builder, tokens []generated.Token, keywords []generated.Keyword) {
	WriteSB(sb, "func NewLexer() lexer.Lexer {", EOL)
	WriteSB(sb, Indent, "return lexer.NewLexer(", EOL)
	for _, keyword := range keywords {
		WriteSB(sb, Indent, Indent, GeneratedKeywordName(keyword), ",", EOL)
	}
	for _, token := range tokens {
		WriteSB(sb, Indent, Indent, GeneratedTokenName(token), ",", EOL)
	}
	WriteSB(sb, Indent, ")", EOL)
	WriteSB(sb, "}", EOL)
}

func generateKeywordTokenType(sb *strings.Builder, keyword generated.Keyword, id int) {
	WriteSB(sb, "const ", GeneratedKeywordIdxName(keyword), " = ", strconv.Itoa(id), EOL, EOL)
	WriteSB(sb, "var ", GeneratedKeywordName(keyword), " = core.NewTokenType(", EOL)
	WriteSB(sb, Indent, GeneratedKeywordIdxName(keyword), ",", EOL)
	WriteSB(sb, Indent, "\"", KeywordValue(keyword), "\",", EOL)
	WriteSB(sb, Indent, "\"", KeywordValue(keyword), "\",", EOL)
	WriteSB(sb, Indent, "0,", EOL)
	WriteSB(sb, Indent, "0,", EOL)
	WriteSB(sb, Indent, "false,", EOL)
	WriteSB(sb, Indent, "func (text string, offset int) int {", EOL)
	WriteSB(sb, Indent, Indent, "if strings.HasPrefix(text[offset:], \"", KeywordValue(keyword), "\") {", EOL)
	WriteSB(sb, Indent, Indent, Indent, "return ", strconv.Itoa(utf8.RuneCountInString(KeywordValue(keyword))), EOL)
	WriteSB(sb, Indent, Indent, "}", EOL)
	WriteSB(sb, Indent, Indent, "return 0", EOL)
	WriteSB(sb, Indent, "},", EOL)
	WriteSB(sb, Indent, "[]rune{", EOLIndent(2))
	runeArrayToString([]rune{rune(KeywordValue(keyword)[0])}, sb)
	WriteSB(sb, EOLIndent(1), "},", EOL)
	WriteSB(sb, ")", EOL)
}

func generateTokenType(sb *strings.Builder, token generated.Token, id int) {
	regexPattern := token.Regexp()
	regexPattern = regexPattern[1 : len(regexPattern)-1] // remove leading and trailing backticks
	regex, err := syntax.Parse(regexPattern, syntax.Perl)
	if err != nil {
		panic(err)
	}
	regex = regex.Simplify()
	WriteSB(sb, "const ", GeneratedTokenIdxName(token), " = ", strconv.Itoa(id), EOL)
	WriteSB(sb, "var ", GeneratedTokenName(token), "_Regexp = regexp.MustCompile(`^(", regex.String(), ")`) ", EOL)
	WriteSB(sb, "var ", GeneratedTokenName(token), " = core.NewTokenType(", EOL)
	WriteSB(sb, Indent, GeneratedTokenIdxName(token), ",", EOL)
	WriteSB(sb, Indent, "\"", token.Name(), "\",", EOL)
	WriteSB(sb, Indent, "\"", token.Name(), "\",", EOL)
	if token.Type() == "hidden" {
		WriteSB(sb, Indent, "-1,", EOL)
	} else {
		WriteSB(sb, Indent, "0,", EOL)
	}
	WriteSB(sb, Indent, "0,", EOL)
	WriteSB(sb, Indent, "false,", EOL)
	WriteSB(sb, Indent, "func (text string, offset int) int {", EOL)
	WriteSB(sb, Indent, Indent, "matches := ", GeneratedTokenName(token), "_Regexp.FindStringIndex(text[offset:])", EOL)
	WriteSB(sb, Indent, Indent, "if matches != nil {", EOL)
	WriteSB(sb, Indent, Indent, Indent, "return matches[1]", EOL)
	WriteSB(sb, Indent, Indent, "}", EOL)
	WriteSB(sb, Indent, Indent, "return 0", EOL)
	WriteSB(sb, Indent, "},", EOL)
	WriteSB(sb, Indent, "[]rune{", EOLIndent(2))
	runeArrayToString(getStartChars(regexPattern), sb)
	WriteSB(sb, EOLIndent(1), "},", EOL)
	WriteSB(sb, ")", EOL)
}

type RuneSlice []rune

func (x RuneSlice) Len() int           { return len(x) }
func (x RuneSlice) Less(i, j int) bool { return x[i] < x[j] }
func (x RuneSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// Sort is a convenience method: x.Sort() calls Sort(x).
func (x RuneSlice) Sort() { sort.Sort(x) }

func runeArrayToString(runes RuneSlice, sb *strings.Builder) {
	runes.Sort()
	for i, r := range runes {
		if r == '\'' {
			sb.WriteString("'\\''")
		} else if r == '\\' {
			sb.WriteString("'\\\\'")
		} else if r >= 32 && r <= 126 {
			WriteSB(sb, "'", string(r), "'")
		} else {
			WriteSB(sb, fmt.Sprint(int64(r)))
		}
		sb.WriteString(",")
		if (i+1)%10 == 0 && i < len(runes)-1 {
			sb.WriteString(EOLIndent(2))
		} else if i < len(runes)-1 {
			sb.WriteString(" ")
		}
	}
}

func getStartChars(pattern string) []rune {
	runes := map[rune]bool{}
	regex, err := syntax.Parse(pattern, syntax.POSIX)
	if err != nil {
		panic(err)
	}
	regex = regex.Simplify()
	startChars, _ := getStartCharsFromRegex(regex)
	for _, r := range startChars {
		runes[r] = true
	}
	returnValue := []rune{}
	for k := range runes {
		returnValue = append(returnValue, k)
	}
	return returnValue
}

func getStartCharsFromRegex(regex *syntax.Regexp) ([]rune, bool) {
	switch regex.Op {
	case syntax.OpConcat:
		runes := []rune{}
		for _, sub := range regex.Sub {
			subRunes, isOptional := getStartCharsFromRegex(sub)
			runes = append(runes, subRunes...)
			if !isOptional {
				return runes, false
			}
		}
		return runes, true
	case syntax.OpLiteral:
		return []rune{regex.Rune[0]}, false
	case syntax.OpCharClass:
		runes := []rune{}
		// The rune slice contains pairs of runes that form ranges
		// e.g. [a-zA-Z] is represented as []rune{'a', 'z', 'A', 'Z'}
		// so we need to iterate over the slice in steps of 2
		for i := 0; i < len(regex.Rune); i += 2 {
			for r := regex.Rune[i]; r <= regex.Rune[i+1]; r++ {
				runes = append(runes, r)
			}
		}
		return runes, false
	case syntax.OpAnyCharNotNL:
	case syntax.OpAnyChar:
		size := 256
		runes := make([]rune, size)
		for r := range size {
			runes[r] = rune(r)
		}
		return runes, false
	case syntax.OpCapture:
		return getStartCharsFromRegex(regex.Sub[0])
	case syntax.OpStar:
	case syntax.OpQuest:
		runes, _ := getStartCharsFromRegex(regex.Sub[0])
		return runes, true
	case syntax.OpPlus:
		runes, _ := getStartCharsFromRegex(regex.Sub[0])
		return runes, false
	case syntax.OpAlternate:
		runes := []rune{}
		for _, sub := range regex.Sub {
			subRunes, _ := getStartCharsFromRegex(sub)
			runes = append(runes, subRunes...)
		}
		return runes, false
	case syntax.OpEmptyMatch:
		return []rune{}, false
	}
	return []rune{}, false
}
