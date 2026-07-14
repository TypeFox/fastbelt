// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"context"
	"fmt"
	"maps"
	"regexp"
	"slices"
	"sort"
	"strconv"
	"unicode/utf8"

	"typefox.dev/fastbelt/internal/automatons"
	"typefox.dev/fastbelt/internal/grammar"
	fbRegexp "typefox.dev/fastbelt/internal/regexp"
	"typefox.dev/fastbelt/util/codegen"
)

func GenerateLexer(grammr grammar.Grammar, packageName string, tokenTypes GenerateTokenTypesResult) string {
	nodes := []codegen.Node{}

	imports := map[string]bool{}
	maps.Copy(imports, tokenTypes.Imports)

	for _, tokenType := range tokenTypes.TokenTypes.All {
		nodes = append(nodes, tokenType.Code)
	}

	node := NewRootNode()
	node.AppendLine("package ", packageName)
	node.AppendLine()
	node.AppendLine("import (")
	node.Indent(func(n codegen.Node) {
		importList := make([]string, 0, len(imports))
		for imp := range imports {
			importList = append(importList, imp)
		}
		sort.Strings(importList)
		for _, imp := range importList {
			n.AppendLine(fmt.Sprintf(`"%s"`, imp))
		}
		n.AppendLine("core \"typefox.dev/fastbelt\"")
		n.AppendLine("\"typefox.dev/fastbelt/lexer\"")
	})
	node.AppendLine(")")
	node.AppendLine()

	for _, n := range nodes {
		node.AppendNode(n)
		node.AppendLine()
	}

	generateLexerModeEnums(node, tokenTypes)
	generateMainLexerFunction(context.Background(), node, tokenTypes)
	return FormatIfPossible(node.String())
}

func generateLexerModeEnums(node codegen.Node, tokenTypes GenerateTokenTypesResult) {
	node.AppendLine("const (")
	node.Indent(func(n codegen.Node) {
		for _, modeName := range tokenTypes.TokenModeOrder {
			modeId := tokenTypes.TokenModes[modeName].Id
			n.AppendLine("TokenMode_", modeName, " = ", strconv.Itoa(modeId))
		}
	})
	node.AppendLine(")")
	node.AppendLine()
}

func generateMainLexerFunction(context context.Context, node codegen.Node, tokenTypes GenerateTokenTypesResult) {
	node.AppendLine("func NewLexer() lexer.Lexer {")
	node.Indent(func(n codegen.Node) {
		count := strconv.Itoa(len(tokenTypes.TokenModes))
		n.AppendLine("modes := make([]*lexer.TokenMode, " + count + ", " + count + ")")
		for _, modeName := range tokenTypes.TokenModeOrder {
			tokenMode := tokenTypes.TokenModes[modeName]
			varName := "modes[" + tokenTypes.TokenModes[modeName].VarName + "]"
			n.AppendLine(varName, " = lexer.NewTokenMode(\"", modeName, "\",")
			n.Indent(func(nn codegen.Node) {
				for _, tokenIndex := range tokenMode.TokenTypeIndices {
					tokenType := tokenTypes.TokenTypes.ByTokenIndex[tokenIndex]
					generateTokenTypeUsage(context, nn, tokenType, tokenMode, tokenIndex, tokenTypes)
				}
			})
			n.AppendLine(")")
		}
		n.AppendLine("return lexer.NewDefaultLexer(" + tokenTypes.TokenModes["default"].VarName + ", modes...)")
	})
	node.AppendLine("}")
}

func generateTokenTypeUsage(context context.Context, nn codegen.Node, tokenType *TokenType, tokenMode *TokenMode, tokenIndex int, tokenTypes GenerateTokenTypesResult) {
	nn.Append("lexer.UseTokenType(", tokenType.VarName, ")")
	if usage, ok := tokenMode.TokenTypeUsages[tokenIndex]; ok {
		if cmd := usage.Command; cmd != nil {
			var cmdModeName string
			if mode := cmd.Mode().Ref(context); cmd.IsDefault() || mode != nil {
				if cmd.IsDefault() {
					cmdModeName = "default"
				} else {
					cmdModeName = cmd.Mode().Ref(context).Name()
				}
			}
			switch cmd.Type() {
			case "push":
				nn.Append(".WithPushMode(", tokenTypes.TokenModes[cmdModeName].VarName, ")")
			case "pop":
				nn.Append(".WithPopMode()")
			default: //"mode"
				nn.Append(".WithSetMode(", tokenTypes.TokenModes[cmdModeName].VarName, ")")
			}
		}
		switch usage.GroupType {
		case "comment":
			nn.Append(".WithGroup(core.CommentGroup)")
		case "hidden":
			nn.Append(".WithGroup(core.SkippedGroup)")
		}
	}
	nn.AppendLine(",")
}

func generateKeywordTokenType(keyword grammar.Keyword, id int) GenerateLexerResult {
	code := codegen.NewNode()
	keywordValue := grammar.KeywordValue(keyword)
	code.AppendLine("const ", GeneratedTokenIdxName(keyword), " = ", strconv.Itoa(id))
	code.AppendLine()
	code.AppendLine("var ", GeneratedTokenName(keyword), " = core.NewTokenType(")
	code.Indent(func(n codegen.Node) {
		n.AppendLine(GeneratedTokenIdxName(keyword), ",")
		n.AppendLine("\"", keywordValue, "\",")
		n.AppendLine("\"", keywordValue, "\",")
		n.AppendLine("core.TokenKindKeyword,")
		n.AppendLine("func (text string, offset int) int {")
		n.Indent(func(nn codegen.Node) {
			nn.AppendLine("if strings.HasPrefix(text[offset:], \"", keywordValue, "\") {")
			nn.Indent(func(nnn codegen.Node) {
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
	code.Append(")")
	return GenerateLexerResult{
		Imports: map[string]bool{
			"strings": true,
		},
		Code: code,
	}
}

type GenerateLexerResult struct {
	Imports map[string]bool
	Code    codegen.Node
}

func generateTokenGroupType(tokenGroup grammar.TokenGroup, tokenGroupMembers map[string][]string, id int) GenerateLexerResult {
	code := codegen.NewNode()
	code.AppendLine("const ", GeneratedTokenIdxName(tokenGroup), " = ", strconv.Itoa(id))
	code.AppendLine()
	code.AppendLine("var ", GeneratedTokenName(tokenGroup), " = core.NewTokenGroup(")
	code.Indent(func(n codegen.Node) {
		n.AppendLine(GeneratedTokenIdxName(tokenGroup), ",")
		n.AppendLine("\"", tokenGroup.Name(), "\",")
		n.AppendLine("\"", tokenGroup.Name(), "\",")
		n.AppendLine("[]*core.TokenType{")
		for _, member := range tokenGroupMembers[tokenGroup.Name()] {
			n.AppendLine(member, ",")
		}
		n.AppendLine("},")
	})
	code.Append(")")
	return GenerateLexerResult{
		Imports: map[string]bool{},
		Code:    code,
	}
}

func getAllTokenGroupMembers(tokenGroup grammar.TokenGroup, keywords GetAllKeywordsResult) []string {
	members := map[string]bool{}
	for _, tokenRef := range tokenGroup.TokenRefs() {
		tokenRule := tokenRef.Ref(context.Background())
		if tokenRule != nil {
			name := GeneratedTokenName(tokenRule)
			members[name] = true
		}
	}
	for _, selector := range tokenGroup.KeywordSelectors() {
		pattern := regexp.MustCompile(grammar.RegexpValue(selector.Image))
		for _, keyword := range keywords.Keywords {
			name := GeneratedTokenName(keyword)
			value := grammar.KeywordValue(keyword)
			if !members[name] && pattern.MatchString(value) {
				members[name] = true
			}
		}
	}
	for _, keyword := range tokenGroup.Keywords() {
		name := GeneratedTokenName(keyword)
		members[name] = true
	}
	slice := slices.Collect(maps.Keys(members))
	sort.Strings(slice)
	return slice
}

func generateRegexpTokenElement(token grammar.TokenDecl, regexpTokenElement grammar.RegexpTokenElement, id int) GenerateLexerResult {
	var result fbRegexp.GenerateRegExpResult
	imports := map[string]bool{}
	code := codegen.NewNode()
	regexPattern := regexpTokenElement.Regexp()
	regexPattern = grammar.RegexpValue(regexPattern)
	regex, err := fbRegexp.Compile(regexPattern)
	if err != nil {
		panic(err)
	}
	code.AppendLine("const ", GeneratedTokenIdxName(token), " = ", strconv.Itoa(id))
	code.AppendLine("var ", GeneratedTokenName(token), " = core.NewTokenType(")
	code.Indent(func(n codegen.Node) {
		n.AppendLine(GeneratedTokenIdxName(token), ",")
		n.AppendLine("\"", token.Name(), "\",")
		n.AppendLine("\"", token.Name(), "\",")
		n.AppendLine("core.TokenKindToken,")
		impl := regex.(*fbRegexp.RegexpImpl)
		result = impl.GenerateRegExp("", GeneratedTokenName(token))
		for imp := range result.Imports {
			imports[imp] = true
		}
		n.AppendNode(result.Code)
		n.AppendLine(",")
		n.Append("[]rune{")
		startCharsSet := impl.GetStartChars()
		n.AppendNode(runeSetToNode(startCharsSet))
		n.AppendLine("},")
	})
	code.AppendLine(")")
	code.AppendNode(result.Vars)
	return GenerateLexerResult{
		Imports: imports,
		Code:    code,
	}
}

type RuneSlice []rune

func (x RuneSlice) Len() int           { return len(x) }
func (x RuneSlice) Less(i, j int) bool { return x[i] < x[j] }
func (x RuneSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }

// Sort is a convenience method: x.Sort() calls Sort(x).
func (x RuneSlice) Sort() { sort.Sort(x) }

func runeSetToNode(set *automatons.RuneSet) codegen.Node {
	root := codegen.NewNode()
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
