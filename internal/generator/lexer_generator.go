// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"fmt"
	"regexp/syntax"
	"sort"
	"strconv"
	"unicode/utf8"

	gen "typefox.dev/fastbelt/generator"
	"typefox.dev/fastbelt/internal/grammar/generated"
	"typefox.dev/fastbelt/internal/regexp"
)

func GenerateLexer(grammar generated.Grammar) string {
	node := gen.NewNode()
	node.AppendLine("package generated")
	node.AppendLine()
	node.AppendLine("import (")
	node.Indent(func(n gen.Node) {
		n.AppendLine("\"regexp\"")
		n.AppendLine("\"strings\"")
		n.AppendLine()
		n.AppendLine("core \"typefox.dev/fastbelt\"")
		n.AppendLine("\"typefox.dev/fastbelt/lexer\"")
	})
	node.AppendLine(")")
	node.AppendLine()

	tokens := grammar.Terminals()
	keywords := GetAllKeywords(grammar)
	id := 1
	for _, keyword := range keywords {
		generateKeywordTokenType(node, keyword, id)
		id++
	}
	for _, token := range tokens {
		generateTokenType(node, token, id)
		id++
	}
	generateMainLexerFunction(node, tokens, keywords)
	return formatIfPossible(node.String())
}

func generateMainLexerFunction(node gen.Node, tokens []generated.Token, keywords []generated.Keyword) {
	node.AppendLine("func NewLexer() lexer.Lexer {")
	node.Indent(func(n gen.Node) {
		n.AppendLine("return lexer.NewDefaultLexer(")
		n.Indent(func(nn gen.Node) {
			for _, keyword := range keywords {
				nn.AppendLine(GeneratedKeywordName(keyword), ",")
			}
			for _, token := range tokens {
				nn.AppendLine(GeneratedTokenName(token), ",")
			}
		})
		n.AppendLine(")")
	})
	node.AppendLine("}")
}

func generateKeywordTokenType(node gen.Node, keyword generated.Keyword, id int) {
	keywordValue := KeywordValue(keyword)
	node.AppendLine("const ", GeneratedKeywordIdxName(keyword), " = ", strconv.Itoa(id))
	node.AppendLine()
	node.AppendLine("var ", GeneratedKeywordName(keyword), " = core.NewTokenType(")
	node.Indent(func(n gen.Node) {
		n.AppendLine(GeneratedKeywordIdxName(keyword), ",")
		n.AppendLine("\"", keywordValue, "\",")
		n.AppendLine("\"", keywordValue, "\",")
		n.AppendLine("0,")
		n.AppendLine("0,")
		n.AppendLine("false,")
		n.AppendLine("func (text string, offset int) int {")
		n.Indent(func(nn gen.Node) {
			nn.AppendLine("if strings.HasPrefix(text[offset:], \"", keywordValue, "\") {")
			nn.Indent(func(nnn gen.Node) {
				nnn.AppendLine("return ", strconv.Itoa(len(keywordValue)))
			})
			nn.AppendLine("}")
			nn.AppendLine("return 0")
		})
		n.AppendLine("},")
		n.Append("[]rune{")
		firstRune, _ := utf8.DecodeRune([]byte(keywordValue))
		n.Indent(func(nn gen.Node) {
			runeArrayToString([]rune{firstRune}, nn)
		})
		n.AppendLine("},")
	})
	node.AppendLine(")")
}

func generateTokenType(node gen.Node, token generated.Token, id int) {
	regexPattern := token.Regexp()
	regexPattern = regexPattern[1 : len(regexPattern)-1] // remove leading and trailing backticks
	regex, err := regexp.Compile(regexPattern)
	if err != nil {
		panic(err)
	}
	node.AppendLine("const ", GeneratedTokenIdxName(token), " = ", strconv.Itoa(id))
	node.AppendLine("var ", GeneratedTokenName(token), " = core.NewTokenType(")
	node.Indent(func(n gen.Node) {
		n.AppendLine(GeneratedTokenIdxName(token), ",")
		n.AppendLine("\"", token.Name(), "\",")
		n.AppendLine("\"", token.Name(), "\",")
		if token.Type() == "hidden" {
			n.AppendLine("-1,")
		} else {
			n.AppendLine("0,")
		}
		n.AppendLine("0,")
		n.AppendLine("false,")
		n.AppendNode(regex.(*regexp.RegexpImpl).GenerateLambda())
		n.Append("[]rune{")
		n.Indent(func(nn gen.Node) {
			runeArrayToString(getStartChars(regexPattern), nn)
		})
		n.AppendLine("},")
	})
	node.AppendLine(")")
}

type RuneSlice []rune

func (x RuneSlice) Len() int           { return len(x) }
func (x RuneSlice) Less(i, j int) bool { return x[i] < x[j] }
func (x RuneSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// Sort is a convenience method: x.Sort() calls Sort(x).
func (x RuneSlice) Sort() { sort.Sort(x) }

func runeArrayToString(runes RuneSlice, node gen.Node) {
	runes.Sort()
	for i, r := range runes {
		var runeStr string
		if r == '\'' {
			runeStr = "'\\'"
		} else if r == '\\' {
			runeStr = "'\\\\'"
		} else if r >= 32 && r <= 126 {
			runeStr = "'" + string(r) + "'"
		} else {
			runeStr = fmt.Sprint(int64(r))
		}

		if (i+1)%10 == 0 && i < len(runes)-1 {
			node.AppendLine(runeStr, ",")
		} else if i < len(runes)-1 {
			node.Append(runeStr, ", ")
		} else {
			node.Append(runeStr, ",")
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
