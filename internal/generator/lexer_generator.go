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

type TokenType struct {
	TokenIndex int
	VarName    string
	Name       string
	Code       codegen.Node
}

type TokenTypeUsage struct {
	GroupType string
	Command   grammar.TokenCommand
}

type TokenMode struct {
	Id               int
	VarName          string
	TokenTypeIndices []int
	TokenTypeUsages  map[int]TokenTypeUsage
}

type GenerateTokenTypesResult struct {
	Imports         map[string]bool
	Keywords        GetAllKeywordsResult
	Tokens          map[grammar.TokenDecl]int
	TokenGroups     map[grammar.TokenGroup]int
	TokenTypes      []TokenType
	TokenTypeUsages map[int]TokenTypeUsage
	TokenModes      map[string]*TokenMode
	TokenModeOrder  []string
}

func (r GenerateTokenTypesResult) TokenTypeIds() map[string]int {
	tokenTypeIds := map[string]int{}
	for _, tokenType := range r.TokenTypes {
		// TokenIndex is the runtime token id emitted by the lexer (the _Idx
		// constant). Keyword-backed tokens alias their keyword, so several
		// token types can share a TokenIndex. This must match la.TypeId at
		// runtime, so the ATN and lookahead tables agree with the lexer.
		tokenTypeIds[tokenType.Name] = tokenType.TokenIndex
	}
	return tokenTypeIds
}

// TokenTypeVarNamesByTokenIndex returns a slice mapping runtime token id
// (TokenIndex) to a token var name. Keyword-backed tokens alias their keyword
// and share its id; later entries (token declarations, groups) override the
// earlier keyword entry so the referenced token name wins.
func (r GenerateTokenTypesResult) TokenTypeVarNamesByTokenIndex() []string {
	maxId := 0
	for _, tokenType := range r.TokenTypes {
		if tokenType.TokenIndex > maxId {
			maxId = tokenType.TokenIndex
		}
	}
	names := make([]string, maxId+1)
	for _, tokenType := range r.TokenTypes {
		names[tokenType.TokenIndex] = tokenType.VarName
	}
	return names
}

func GenerateTokenTypes(grammr grammar.Grammar) GenerateTokenTypesResult {
	keywords := GetAllKeywords(grammr)
	tokens := grammr.Terminals()
	tokenGroups := grammr.TokenGroups()
	result := GenerateTokenTypesResult{
		Imports:         map[string]bool{},
		Keywords:        keywords,
		Tokens:          make(map[grammar.TokenDecl]int),
		TokenGroups:     make(map[grammar.TokenGroup]int),
		TokenTypes:      make([]TokenType, 0),
		TokenTypeUsages: make(map[int]TokenTypeUsage),
		TokenModes:      make(map[string]*TokenMode),
		TokenModeOrder:  []string{},
	}
	// Starting with 1 - prevent clash with EOF (index 0)
	tokenIndex := 1
	// keywordTokenIndex records the runtime token id (TokenIndex) assigned to
	// each keyword node, so keyword-backed tokens can alias to the correct id.
	keywordTokenIndex := map[string]int{}
	for _, keyword := range keywords.Keywords {
		code := generateKeywordTokenType(keyword, tokenIndex)
		mergeImports(&result.Imports, code.Imports)
		keywordTokenIndex[keyword.Value()] = tokenIndex
		result.TokenTypes = append(result.TokenTypes, TokenType{
			TokenIndex: tokenIndex,
			VarName:    GeneratedTokenName(keyword),
			Name:       keyword.Value(),
			Code:       code.Code,
		})
		tokenIndex++
	}
	for _, token := range tokens {
		varName := GeneratedTokenName(token)
		index := len(result.TokenTypes)
		switch element := token.Content().(type) {
		case grammar.KeywordTokenElement:
			keywordIndex := keywords.ByValue[element.Keyword().Value()]
			keyword := keywords.Keywords[keywordIndex]
			code := codegen.NewNode()
			code.AppendLine("const ", GeneratedTokenIdxName(token), " = ", GeneratedTokenIdxName(keyword))
			code.AppendLine()
			code.AppendLine("var ", varName, " = ", GeneratedTokenName(keyword))
			mergeImports(&result.Imports, map[string]bool{})
			keywordsTokenIndex := keywordTokenIndex[keyword.Value()]
			result.TokenTypes = append(result.TokenTypes, TokenType{
				TokenIndex: keywordsTokenIndex,
				VarName:    varName,
				Name:       token.Name(),
				Code:       code,
			})
		case grammar.RegexpTokenElement:
			lexerResult := generateRegexpTokenElement(token, element, tokenIndex)
			mergeImports(&result.Imports, lexerResult.Imports)
			result.TokenTypes = append(result.TokenTypes, TokenType{
				TokenIndex: tokenIndex,
				VarName:    varName,
				Name:       token.Name(),
				Code:       lexerResult.Code,
			})
			tokenIndex++
		}
		result.Tokens[token] = index
		if token.Type() != "" || token.Command() != nil {
			result.TokenTypeUsages[index] = TokenTypeUsage{
				GroupType: token.Type(),
				Command:   token.Command(),
			}
		}
	}
	tokenGroupMembers := map[string][]string{}
	for _, tokenGroup := range tokenGroups {
		tokenGroupMembers[tokenGroup.Name()] = getAllTokenGroupMembers(tokenGroup)
	}
	// Token groups need to be topologically sorted, so that nested groups appear after their members
	for _, tokenGroup := range sortTokenGroups(tokenGroups, tokenGroupMembers) {
		index := len(result.TokenTypes)
		lexerResult := generateTokenGroupType(tokenGroup, tokenGroupMembers, tokenIndex)
		mergeImports(&result.Imports, lexerResult.Imports)
		result.TokenTypes = append(result.TokenTypes, TokenType{
			TokenIndex: tokenIndex,
			VarName:    GeneratedTokenName(tokenGroup),
			Name:       tokenGroup.Name(),
			Code:       lexerResult.Code,
		})
		result.TokenGroups[tokenGroup] = index
		tokenIndex++
	}
	for index, tokenMode := range grammr.TokenModes() {
		modeName := "default"
		if !tokenMode.IsDefault() {
			modeName = tokenMode.Name()
		}
		current := TokenMode{
			Id:               index,
			VarName:          "TokenMode_" + modeName,
			TokenTypeIndices: make([]int, 0),
			TokenTypeUsages:  make(map[int]TokenTypeUsage),
		}
		result.TokenModes[modeName] = &current
		for _, member := range tokenMode.Members() {
			switch member := member.(type) {
			case grammar.TokenDeclUsage:
				//TODO
			case grammar.KeywordUsage:
				//TODO
			case grammar.TokenUsage:
				token := member.TokenRef().Ref(context.Background())
				var index int
				switch rule := token.(type) {
				case grammar.TokenDecl:
					index = result.Tokens[rule]
					current.TokenTypeIndices = append(current.TokenTypeIndices, index)
				case grammar.TokenGroup:
					index = result.TokenGroups[rule]
					current.TokenTypeIndices = append(current.TokenTypeIndices, index)
				}
				if member.Type() != "" || member.Command() != nil {
					current.TokenTypeUsages[index] = TokenTypeUsage{
						GroupType: member.Type(),
						Command:   member.Command(),
					}
				}
			case grammar.KeywordSelector:
				//TODO
			}
		}
		result.TokenModeOrder = append(result.TokenModeOrder, modeName)
	}
	if result.TokenModes["default"] == nil {
		//if token mode "default" is not defined, we need to create it
		result.TokenModes["default"] = &TokenMode{
			Id:               len(result.TokenModes),
			VarName:          "TokenMode_default",
			TokenTypeIndices: make([]int, 0),
			TokenTypeUsages:  make(map[int]TokenTypeUsage),
		}
		result.TokenModeOrder = append(result.TokenModeOrder, "default")
	}
	return result
}

func mergeImports(target *map[string]bool, source map[string]bool) {
	for imp := range source {
		(*target)[imp] = true
	}
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

	for _, tokenType := range tokenTypes.TokenTypes {
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
				for _, index := range tokenMode.TokenTypeIndices {
					tokenType := tokenTypes.TokenTypes[index]
					nn.Append("lexer.UseTokenType(", tokenType.VarName, ")")
					if usage, ok := tokenMode.TokenTypeUsages[index]; ok {
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
						if usage.GroupType == "comment" {
							nn.Append(".WithGroup(core.CommentGroup)")
						} else if usage.GroupType == "hidden" {
							nn.Append(".WithGroup(core.SkippedGroup)")
						}
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
