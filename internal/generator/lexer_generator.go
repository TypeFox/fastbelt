// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"sort"
	"strconv"
	"unicode/utf8"

	gen "typefox.dev/fastbelt/generator"
	"typefox.dev/fastbelt/internal/automatons"
	"typefox.dev/fastbelt/internal/grammar/generated"
	"typefox.dev/fastbelt/internal/regexp"
)

func GenerateLexer(grammar generated.Grammar) string {
	node := gen.NewNode()
	node.AppendLine("package generated")
	node.AppendLine()
	node.AppendLine("import (")
	node.Indent(func(n gen.Node) {
		n.AppendLine("\"strings\"")
		n.AppendLine()
		n.AppendLine("core \"typefox.dev/fastbelt\"")
		n.AppendLine("\"typefox.dev/fastbelt/lexer\"")
		n.AppendLine("\"typefox.dev/fastbelt/internal/regexp\"")
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
		n.Append(automatons.FormatRune(firstRune))
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
		impl := regex.(*regexp.RegexpImpl)
		n.AppendNode(impl.GenerateRegExp())
		n.Append("[]rune{")
		startCharsSet := impl.GetStartChars()
		n.AppendNode(runeSetToNode(startCharsSet))
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

func runeSetToNode(set *automatons.RuneSet) gen.Node {
	root := gen.NewNode()
	for _, rng := range set.Ranges {
		if !rng.Includes {
			continue
		}
		if rng.Start == rng.End {
			root.Append(automatons.FormatRune(rng.Start), ", ")
		} else {
			for r := rng.Start; r <= rng.End; r++ {
				root.Append(automatons.FormatRune(r))
				root.Append(", ")
			}
		}
	}
	return root
}
