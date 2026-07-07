// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"sort"
	"strconv"
	"unicode/utf8"

	"typefox.dev/fastbelt/internal/automatons"
	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/internal/regexp"
	"typefox.dev/fastbelt/util/codegen"
)

type TokenMode struct {
	Id               int
	VarName          string
	TokenDecls       []grammar.TokenDecl
	TokenUsages      []grammar.TokenUsage
	Keywords         []grammar.KeywordUsage
	KeywordSelectors []string
}

type GenerateTokenTypesResult struct {
	Tokens            []grammar.TokenDecl
	Keywords          []grammar.Keyword
	Imports           map[string]bool
	KeywordsCode      []codegen.Node
	TokensCode        []codegen.Node
	TokenGroupCode    []codegen.Node
	TokenTypeVarNames []string
	TokenTypeNames    []string
	TokenTypeIds      map[string]int
	TokenModes        map[string]*TokenMode
	TokenModeOrder    []string
}

func GenerateTokenTypes(grammr grammar.Grammar) GenerateTokenTypesResult {
	tokens := grammr.Terminals()
	tokenGroups := grammr.TokenGroups()
	result := GenerateTokenTypesResult{
		Tokens:            tokens,
		TokensCode:        make([]codegen.Node, len(tokens)),
		TokenGroupCode:    make([]codegen.Node, len(tokenGroups)),
		TokenTypeNames:    make([]string, len(tokens)+len(tokenGroups)),
		TokenTypeVarNames: make([]string, len(tokens)+len(tokenGroups)),
		TokenTypeIds:      make(map[string]int),
		Imports:           map[string]bool{},
		TokenModes:        map[string]*TokenMode{},
	}
	// Starting with 1 - prevent clash with EOF (index 0)
	tokenIndex := 1
	for index, token := range tokens {
		tokenType := generateTokenType(token, tokenIndex)
		result.TokensCode[index] = tokenType.Code
		for imp := range tokenType.Imports {
			result.Imports[imp] = true
		}
		varName := GeneratedTokenName(token)
		tokName := token.Name()
		result.TokenTypeVarNames[index] = varName
		result.TokenTypeNames[index] = tokName
		result.TokenTypeIds[tokName] = index
		tokenIndex++
	}
	tokenGroupMembers := map[string][]string{}
	for _, tokenGroup := range tokenGroups {
		tokenGroupMembers[tokenGroup.Name()] = getAllTokenGroupMembers(tokenGroup)
	}
	// Token groups need to be topologically sorted, so that nested groups appear after their members
	for index, tokenGroup := range sortTokenGroups(tokenGroups, tokenGroupMembers) {
		result.TokenGroupCode[index] = generateTokenGroupType(tokenGroup, tokenGroupMembers, tokenIndex)
		varName := GeneratedTokenName(tokenGroup)
		tokName := tokenGroup.Name()
		result.TokenTypeVarNames[len(tokens)+index] = varName
		result.TokenTypeNames[len(tokens)+index] = tokName
		result.TokenTypeIds[tokName] = len(tokens) + index
		tokenIndex++
	}
	for index, tokenMode := range grammr.TokenModes() {
		modeName := "default"
		if !tokenMode.IsDefault() {
			modeName = tokenMode.Name()
		}
		result.TokenModes[modeName] = &TokenMode{
			Id:          index,
			VarName:     "TokenMode_" + modeName,
			TokenUsages: tokenMode.TokenRefs(),
		}
		result.TokenModeOrder = append(result.TokenModeOrder, modeName)
	}
	if result.TokenModes["default"] == nil {
		result.TokenModes["default"] = &TokenMode{
			Id:      len(result.TokenModes),
			VarName: "TokenMode_default",
			//TODO
			TokenDecls:       nil,
			TokenUsages:      nil,
			Keywords:         nil,
			KeywordSelectors: nil,
		}
		result.TokenModeOrder = append(result.TokenModeOrder, "default")
	}
	return result
}

func sortTokenGroups(tokenGroups []grammar.TokenGroup, members map[string][]string) []grammar.TokenGroup {
	names := make([]string, len(tokenGroups))
	for i, tg := range tokenGroups {
		names[i] = tg.Name()
	}
	sort.Strings(names)
	topoSorted := topoSort(names, members, false)
	sortedGroups := make([]grammar.TokenGroup, len(tokenGroups))
	for i, name := range topoSorted {
		for _, tg := range tokenGroups {
			if tg.Name() == name {
				sortedGroups[i] = tg
				break
			}
		}
	}
	return sortedGroups
}

func GenerateLexer(grammr grammar.Grammar, packageName string, tokenTypes GenerateTokenTypesResult) string {
	nodes := []codegen.Node{}

	imports := map[string]bool{}
	maps.Copy(imports, tokenTypes.Imports)
	nodes = append(nodes, tokenTypes.KeywordsCode...)
	nodes = append(nodes, tokenTypes.TokensCode...)
	nodes = append(nodes, tokenTypes.TokenGroupCode...)

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
			tokenUsages := tokenTypes.TokenModes[modeName]
			varName := "modes[" + tokenTypes.TokenModes[modeName].VarName + "]"
			n.AppendLine(varName, " = lexer.NewTokenMode(\"", modeName, "\",")
			n.Indent(func(nn codegen.Node) {
				for _, tokenUsage := range tokenUsages.TokenUsages {
					nn.Append("lexer.UseTokenType(", GeneratedTokenName(tokenUsage.TokenRef().Ref(context)), ")")
					if cmd := tokenUsage.Command(); cmd != nil {
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
					if tokenUsage.Type() == "comment" {
						nn.Append(".WithGroup(core.CommentGroup)")
					} else if tokenUsage.Type() == "hidden" {
						nn.Append(".WithGroup(core.SkippedGroup)")
					}
					nn.AppendLine(",")
				}
			})
			n.AppendLine(")")
		}
		n.AppendLine("return lexer.NewDefaultLexer(" + tokenTypes.TokenModes["default"].VarName + ", modes...)")
	})
	node.AppendLine("}")
}

func generateTokenType(token grammar.TokenDecl, id int) GenerateLexerResult {
	switch element := token.Content().(type) {
	case grammar.KeywordTokenElement:
		return generateKeywordTokenElement(token, element, id)
	case grammar.RegexpTokenElement:
		return generateRegexpTokenElement(token, element, id)
	default:
		panic(fmt.Sprintf("Unknown token element type: %T", element))
	}
}

func generateKeywordTokenElement(token grammar.TokenDecl, keywordElement grammar.KeywordTokenElement, id int) GenerateLexerResult {
	code := codegen.NewNode()
	keyword := keywordElement.Keyword()
	keywordToken := keyword.Value()
	keywordValue := keywordToken[1 : len(keywordToken)-1]
	code.AppendLine("const ", GeneratedTokenIdxName(token), " = ", strconv.Itoa(id))
	code.AppendLine()
	code.AppendLine("var ", GeneratedTokenName(token), " = core.NewTokenType(")
	code.Indent(func(n codegen.Node) {
		n.AppendLine(GeneratedTokenIdxName(token), ",")
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

func generateTokenGroupType(tokenGroup grammar.TokenGroup, tokenGroupMembers map[string][]string, id int) codegen.Node {
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
	return code
}

func getAllTokenGroupMembers(tokenGroup grammar.TokenGroup) []string {
	members := map[string]bool{}
	for _, tokenRef := range tokenGroup.TokenRefs() {
		tokenRule := tokenRef.Ref(context.Background())
		if tokenRule != nil {
			name := GeneratedTokenName(tokenRule)
			members[name] = true
		}
	}
	slice := slices.Collect(maps.Keys(members))
	sort.Strings(slice)
	return slice
}

func generateRegexpTokenElement(token grammar.TokenDecl, regexpTokenElement grammar.RegexpTokenElement, id int) GenerateLexerResult {
	var result regexp.GenerateRegExpResult
	imports := map[string]bool{}
	code := codegen.NewNode()
	regexPattern := regexpTokenElement.Regexp()
	regexPattern = regexPattern[1 : len(regexPattern)-1] // remove leading and trailing backticks
	regex, err := regexp.Compile(regexPattern)
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
		impl := regex.(*regexp.RegexpImpl)
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
