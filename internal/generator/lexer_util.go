// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"context"
	"regexp"
	"slices"
	"sort"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/util/codegen"
)

type TokenType struct {
	TokenIndex int
	VarName    string
	Name       string
	Code       codegen.Node
}

type TokenTypeUsage struct {
	TokenModifier string
	Command       grammar.TokenCommand
}

type TokenMode struct {
	Id               int
	VarName          string
	TokenTypeIndices []int
	TokenTypeUsages  map[int]TokenTypeUsage
}

type TokenIndexLookup struct {
	ByKeyword    map[string]int
	ByToken      map[grammar.TokenDecl]int
	ByTokenGroup map[grammar.TokenGroup]int
}

type TokenTypeLookup struct {
	All          []TokenType
	ByTokenIndex map[int]*TokenType
}

type GenerateTokenTypesResult struct {
	Keywords        GetAllKeywordsResult
	TokenDecls      GetAllTokenDeclsResult
	TokenGroups     GetAllTokenGroupsResult
	Imports         map[string]bool
	TokenIndex      TokenIndexLookup
	TokenTypes      TokenTypeLookup
	TokenTypeUsages map[int]TokenTypeUsage
	TokenModes      map[string]*TokenMode
	TokenModeOrder  []string
}

func (r GenerateTokenTypesResult) TokenTypeIds() map[string]int {
	tokenTypeIds := map[string]int{}
	for _, tokenType := range r.TokenTypes.All {
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
	for _, tokenType := range r.TokenTypes.All {
		if tokenType.TokenIndex > maxId {
			maxId = tokenType.TokenIndex
		}
	}
	names := make([]string, maxId+1)
	for _, tokenType := range r.TokenTypes.All {
		names[tokenType.TokenIndex] = tokenType.VarName
	}
	return names
}

func GenerateTokenTypes(grammr grammar.Grammar) GenerateTokenTypesResult {
	keywords := GetAllKeywords(grammr)
	tokenDecls := GetAllTokenDecls(grammr)
	tokenGroups := GetAllTokenGroups(grammr)
	result := GenerateTokenTypesResult{
		Keywords:    keywords,
		TokenDecls:  tokenDecls,
		TokenGroups: tokenGroups,
		Imports:     map[string]bool{},
		TokenIndex: TokenIndexLookup{
			ByKeyword:    make(map[string]int),
			ByToken:      make(map[grammar.TokenDecl]int),
			ByTokenGroup: make(map[grammar.TokenGroup]int),
		},
		TokenTypes: TokenTypeLookup{
			All:          make([]TokenType, 0),
			ByTokenIndex: make(map[int]*TokenType),
		},
		TokenTypeUsages: make(map[int]TokenTypeUsage),
		TokenModes:      make(map[string]*TokenMode),
		TokenModeOrder:  []string{},
	}
	populateTokenTypes(&result)
	populateTokenModes(&result, grammr.TokenModes())
	return result
}

func populateTokenModes(result *GenerateTokenTypesResult, tokenModes []grammar.TokenMode) {
	keywords := result.Keywords
	tokenGroups := result.TokenGroups
	tokens := result.TokenDecls
	for index, tokenMode := range tokenModes {
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
		// Tracked across all members of the mode: a keyword selector must not
		// register a keyword that an earlier member already added.
		alreadyAdded := map[int]bool{}
		for _, member := range tokenMode.Members() {
			pushTokenTypeUsage := func(tokenIndex int, tokenModifier string, command grammar.TokenCommand) {
				alreadyAdded[tokenIndex] = true
				current.TokenTypeIndices = append(current.TokenTypeIndices, tokenIndex)
				if tokenModifier != "" || command != nil {
					current.TokenTypeUsages[tokenIndex] = TokenTypeUsage{
						TokenModifier: tokenModifier,
						Command:       command,
					}
				}
			}

			switch member := member.(type) {
			case grammar.TokenDeclUsage:
				token := member.Declaration()
				tokenIndex := result.TokenIndex.ByToken[token]
				pushTokenTypeUsage(tokenIndex, token.Modifier(), token.Command())
			case grammar.TokenGroupUsage:
				tokenGroup := member.Group()
				tokenIndex := result.TokenIndex.ByTokenGroup[tokenGroup]
				pushTokenTypeUsage(tokenIndex, tokenGroup.Modifier(), tokenGroup.Command())
			case grammar.KeywordUsage:
				tokenIndex := result.TokenIndex.ByKeyword[member.Keyword().Value()]
				pushTokenTypeUsage(tokenIndex, member.Modifier(), member.Command())
			case grammar.TokenUsage:
				token := member.TokenRef().Ref(context.Background())
				var tokenIndex int
				switch rule := token.(type) {
				case grammar.TokenDecl:
					tokenIndex = result.TokenIndex.ByToken[rule]
					pushTokenTypeUsage(tokenIndex, rule.Modifier(), rule.Command())
				case grammar.TokenGroup:
					tokenIndex = result.TokenIndex.ByTokenGroup[rule]
					pushTokenTypeUsage(tokenIndex, rule.Modifier(), rule.Command())
				}
				//overwrite usage if defined on token mode level
				if member.Modifier() != "" || member.Command() != nil {
					current.TokenTypeUsages[tokenIndex] = TokenTypeUsage{
						TokenModifier: member.Modifier(),
						Command:       member.Command(),
					}
				}
			case grammar.KeywordSelector:
				pattern := regexp.MustCompile(grammar.RegexpValue(member.Selector()))
				for _, keyword := range keywords.Keywords {
					tokenIndex := result.TokenIndex.ByKeyword[keyword.Value()]
					value := grammar.KeywordValue(keyword)
					if !alreadyAdded[tokenIndex] && pattern.MatchString(value) {
						current.TokenTypeIndices = append(current.TokenTypeIndices, tokenIndex)
						alreadyAdded[tokenIndex] = true
					}
				}
			}
		}
		result.TokenModeOrder = append(result.TokenModeOrder, modeName)
	}
	if result.TokenModes["default"] == nil {
		//if token mode "default" is not defined, we need to create it
		defaultMode := TokenMode{
			Id:               len(result.TokenModes),
			VarName:          "TokenMode_default",
			TokenTypeIndices: make([]int, 0),
			TokenTypeUsages:  make(map[int]TokenTypeUsage),
		}
		result.TokenModes["default"] = &defaultMode
		result.TokenModeOrder = append(result.TokenModeOrder, "default")

		for _, tokenGroup := range tokenGroups.TopLevel {
			tokenIndex := result.TokenIndex.ByTokenGroup[tokenGroup]
			defaultMode.TokenTypeIndices = append(defaultMode.TokenTypeIndices, tokenIndex)
			if tokenGroup.Modifier() != "" || tokenGroup.Command() != nil {
				defaultMode.TokenTypeUsages[tokenIndex] = TokenTypeUsage{
					TokenModifier: tokenGroup.Modifier(),
					Command:       tokenGroup.Command(),
				}
			}
		}

		for _, keyword := range keywords.Keywords {
			tokenIndex := result.TokenIndex.ByKeyword[keyword.Value()]
			defaultMode.TokenTypeIndices = append(defaultMode.TokenTypeIndices, tokenIndex)
			//keywords don't have type or command, so we don't need to add anything to TokenTypeUsages
		}

		for _, token := range tokens.TopLevel {
			tokenIndex := result.TokenIndex.ByToken[token]
			if !slices.Contains(defaultMode.TokenTypeIndices, tokenIndex) {
				defaultMode.TokenTypeIndices = append(defaultMode.TokenTypeIndices, tokenIndex)
			}
			if token.Modifier() != "" || token.Command() != nil {
				defaultMode.TokenTypeUsages[tokenIndex] = TokenTypeUsage{
					TokenModifier: token.Modifier(),
					Command:       token.Command(),
				}
			}
		}
	}
}

func populateTokenTypes(result *GenerateTokenTypesResult) {
	keywords := result.Keywords
	tokens := result.TokenDecls
	tokenGroups := result.TokenGroups

	// Starting with 1 - prevent clash with EOF (index 0)
	tokenIndex := 1
	for _, keyword := range keywords.Keywords {
		code := generateKeywordTokenType(keyword, tokenIndex)
		mergeImports(&result.Imports, code.Imports)
		tokenType := TokenType{
			TokenIndex: tokenIndex,
			VarName:    GeneratedTokenName(keyword),
			Name:       keyword.Value(),
			Code:       code.Code,
		}
		result.TokenTypes.All = append(result.TokenTypes.All, tokenType)
		result.TokenTypes.ByTokenIndex[tokenIndex] = &tokenType
		result.TokenIndex.ByKeyword[keyword.Value()] = tokenIndex
		tokenIndex++
	}
	for _, token := range tokens.All {
		varName := GeneratedTokenName(token)
		var currentTokenIndex int
		switch element := token.Content().(type) {
		case grammar.KeywordTokenElement:
			keywordIndex := keywords.ByValue[element.Keyword().Value()]
			keyword := keywords.Keywords[keywordIndex]
			code := codegen.NewNode()
			code.AppendLine("const ", GeneratedTokenIdxName(token), " = ", GeneratedTokenIdxName(keyword))
			code.AppendLine()
			code.AppendLine("var ", varName, " = ", GeneratedTokenName(keyword))
			mergeImports(&result.Imports, map[string]bool{})
			currentTokenIndex := result.TokenIndex.ByKeyword[keyword.Value()]
			tokenType := TokenType{
				TokenIndex: currentTokenIndex,
				VarName:    varName,
				Name:       token.Name(),
				Code:       code,
			}
			result.TokenTypes.All = append(result.TokenTypes.All, tokenType)
			result.TokenTypes.ByTokenIndex[currentTokenIndex] = &tokenType
			result.TokenIndex.ByToken[token] = currentTokenIndex
		case grammar.RegexpTokenElement:
			lexerResult := generateRegexpTokenElement(token, element, tokenIndex)
			mergeImports(&result.Imports, lexerResult.Imports)
			currentTokenIndex = tokenIndex
			tokenType := TokenType{
				TokenIndex: currentTokenIndex,
				VarName:    varName,
				Name:       token.Name(),
				Code:       lexerResult.Code,
			}
			result.TokenTypes.All = append(result.TokenTypes.All, tokenType)
			result.TokenTypes.ByTokenIndex[currentTokenIndex] = &tokenType
			result.TokenIndex.ByToken[token] = currentTokenIndex
			tokenIndex++
		}
		if token.Modifier() != "" || token.Command() != nil {
			result.TokenTypeUsages[currentTokenIndex] = TokenTypeUsage{
				TokenModifier: token.Modifier(),
				Command:       token.Command(),
			}
		}
	}
	tokenGroupMembers := map[string][]string{}
	for _, tokenGroup := range tokenGroups.All {
		tokenGroupMembers[tokenGroup.Name()] = getAllTokenGroupMembers(tokenGroup, keywords)
	}
	// Token groups need to be topologically sorted, so that nested groups appear after their members
	for _, tokenGroup := range sortTokenGroups(tokenGroups.All, tokenGroupMembers) {
		lexerResult := generateTokenGroupType(tokenGroup, tokenGroupMembers, tokenIndex)
		mergeImports(&result.Imports, lexerResult.Imports)
		tokenType := TokenType{
			TokenIndex: tokenIndex,
			VarName:    GeneratedTokenName(tokenGroup),
			Name:       tokenGroup.Name(),
			Code:       lexerResult.Code,
		}
		result.TokenTypes.All = append(result.TokenTypes.All, tokenType)
		result.TokenTypes.ByTokenIndex[tokenIndex] = &tokenType
		result.TokenIndex.ByTokenGroup[tokenGroup] = tokenIndex
		tokenIndex++
	}
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

type GetAllKeywordsResult struct {
	Keywords []grammar.Keyword
	ByValue  map[string]int
}

func GetAllKeywords(grammr grammar.Grammar) GetAllKeywordsResult {
	keywords := map[string]grammar.Keyword{}
	allNodes := []grammar.Keyword{}
	for node := range core.AllChildren(grammr) {
		if keyword, ok := node.(grammar.Keyword); ok {
			keywords[keyword.Value()] = keyword
			allNodes = append(allNodes, keyword)
		}
	}
	return keysFromMap(keywords, allNodes)
}

type GetAllTokenDeclsResult struct {
	TopLevel  []grammar.TokenDecl
	ModeLevel []grammar.TokenDecl
	All       []grammar.TokenDecl
}

func GetAllTokenDecls(grammr grammar.Grammar) GetAllTokenDeclsResult {
	topLevel := []grammar.TokenDecl{}
	modeLevel := []grammar.TokenDecl{}
	for node := range core.AllChildren(grammr) {
		if tokenDecl, ok := node.(grammar.TokenDecl); ok {
			_, ok := tokenDecl.Container().(grammar.TokenDeclUsage)
			if ok {
				modeLevel = append(modeLevel, tokenDecl)
			} else {
				topLevel = append(topLevel, tokenDecl)
			}
		}
	}
	length := len(topLevel) + len(modeLevel)
	all := make([]grammar.TokenDecl, length)
	copy(all, topLevel)
	copy(all[len(topLevel):], modeLevel)
	return GetAllTokenDeclsResult{
		All:       all,
		TopLevel:  topLevel,
		ModeLevel: modeLevel,
	}
}

type GetAllTokenGroupsResult struct {
	All       []grammar.TokenGroup
	TopLevel  []grammar.TokenGroup
	ModeLevel []grammar.TokenGroup
}

func GetAllTokenGroups(grammr grammar.Grammar) GetAllTokenGroupsResult {
	topLevel := []grammar.TokenGroup{}
	modeLevel := []grammar.TokenGroup{}

	for node := range core.AllChildren(grammr) {
		if tokenGroup, ok := node.(grammar.TokenGroup); ok {
			_, ok := tokenGroup.Container().(grammar.TokenGroupUsage)
			if ok {
				modeLevel = append(modeLevel, tokenGroup)
			} else {
				topLevel = append(topLevel, tokenGroup)
			}
		}
	}

	length := len(topLevel) + len(modeLevel)
	all := make([]grammar.TokenGroup, length)
	copy(all, topLevel)
	copy(all[len(topLevel):], modeLevel)
	return GetAllTokenGroupsResult{
		All:       all,
		TopLevel:  topLevel,
		ModeLevel: modeLevel,
	}
}

func keysFromMap(m map[string]grammar.Keyword, allNodes []grammar.Keyword) GetAllKeywordsResult {
	keywords := []grammar.Keyword{}
	for _, v := range m {
		keywords = append(keywords, v)
	}
	sort.Slice(keywords, func(i, j int) bool {
		return keywords[i].Value() < keywords[j].Value()
	})
	byValue := map[string]int{}
	for i, k := range keywords {
		byValue[k.Value()] = i
	}
	return GetAllKeywordsResult{
		Keywords: keywords,
		ByValue:  byValue,
	}
}

func GeneratedTokenName(t core.AstNode) string {
	switch t := t.(type) {
	case grammar.TokenDecl:
		return "Token_" + t.Name()
	case grammar.TokenGroup:
		return "TokenGroup_" + t.Name()
	case grammar.Keyword:
		return "Keyword_" + grammar.KeywordName(t)
	default:
		panic("unexpected type")
	}
}

func GeneratedTokenIdxName(t core.AstNode) string {
	return GeneratedTokenName(t) + "_Idx"
}
