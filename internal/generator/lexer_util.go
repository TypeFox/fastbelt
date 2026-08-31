// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"context"
	"maps"
	"regexp"
	"slices"
	"sort"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/util/codegen"
)

// A helper type to manage generator-related information about token types.
// A token type is either originated from
// - a keyword (standalone keyword or token type with keyword content)
// - a token declaration (token type with regexp content)
// - a token group (token type that groups other token types)
type TokenType struct {
	// Token index is the unique ID at runtime, used as a constant in
	// the generated lexer and parser code. It is also used as the index
	// into the token type lookup table.
	TokenIndex int
	// Name of the token type, used to reflect usage of grammar rules.
	// No extra prefix!
	Name string
	// VarName is the name of the generated variable containing the token type.
	VarName string
	// Generated code to declare the token type variable.
	Code codegen.Node
}

// Describes how a [TokenType] is used in a [TokenMode].
// Mind the same naming in the modes package: [modes.TokenTypeUsage]. The concept is the same. Implementation is slightly different.
type TokenTypeUsage struct {
	// comment or hidden, look at it like public, protected, private in object-oriented programming.
	TokenModifier string
	// Command to manipulate the mode stack when this token type is matched. Can be nil.
	// Distinguishes between:
	// - push(mmode): adds a new mode onto the stack
	// - set(mode): replace the current mode with a new one
	// - pop: pops the current mode from the stack
	Command grammar.TokenCommand
}

// Belongs to a [TokenMode] and describes how a [TokenType] is used in that mode.
// Actually it would be sufficient to store one array of token type indices, but
// we split this set by category: keywords and regexp tokens. Keywords are matched first,
// so they have priority over regexp tokens.
type TokenTypeIndices struct {
	// token indices of keywords that are part of this mode
	Keywords []int
	// token indices of regexp tokens that are part of this mode
	Tokens []int
}

// A token mode ist an ordered set of token types that are active in a given mode.
type TokenMode struct {
	// Unique ID of the token mode, used as a constant.
	Id int
	// Variable name of the generated token mode constant.
	VarName string
	// Holds the token indices of the token types that are part of this mode.
	TokenTypeIndices TokenTypeIndices
	// Holds the token type usage information (modifier and command) for each
	// token type in this mode. The key is the token index.
	TokenTypeUsages map[int]TokenTypeUsage
}

// TokenIndexSource is an enum and describes for one token index, from
// which source it was generated.
type TokenIndexSource int

const (
	// was created by a token declaration (regexp or keyword)
	SourceTokenDecl TokenIndexSource = iota
	// was created by a keyword (standalone and no token declaration with keyword content)
	SourceKeyword
	// was created by a token group (so it has members!)
	SourceGroup
)

// TokenIndexLookup is a helper type to manage the mapping between token
// indices and their sources (keyword, token, token group, keyword).
type TokenIndexLookup struct {
	// Keyword value to token index mapping.
	ByKeyword map[string]int
	// Token declaration to token index mapping.
	ByToken map[grammar.TokenDecl]int
	// Token group to token index mapping.
	ByTokenGroup map[grammar.TokenGroup]int
	// Token group indices of the members of a token group. The key is the parent token group.
	ByTokenGroupParent map[grammar.TokenGroup][]int
	// Token index to source type mapping. So that you can find out from which source a token index was generated (keyword, token, token group).
	SourceType map[int]TokenIndexSource
}

// TokenTypeLookup is a helper type to manage the mapping between token types and their token indices.
type TokenTypeLookup struct {
	// Get all token types that were created (it can happen that multiple
	// token types share the same token index, for example when a keyword
	// was declared standalone and in a token declaration).
	All []TokenType
	// Token index to token type mapping. So that you can find out which token type was generated for a given token index.
	ByTokenIndex map[int]*TokenType
}

// Result of GenerateTokenTypes. It contains all the information about the generated token types, token modes, and their usage in the grammar.
type GenerateTokenTypesResult struct {
	// All keywords that were found in the grammar. This includes standalone keywords and keywords that are part of token declarations.
	Keywords GetAllKeywordsResult
	// All token declarations that were found in the grammar. This includes top-level and mode-level token declarations.
	TokenDecls GetAllTokenDeclsResult
	// All token groups that were found in the grammar. This includes top-level and mode-level token groups.
	TokenGroups GetAllTokenGroupsResult
	// All imports that are required for the generated code. Each entry is a package path.
	Imports map[string]bool
	// Lookup tables for token indices: What token index has X? Where X is a keyword, a token declaration, or a token group.
	TokenIndex TokenIndexLookup
	// Lookup tables for token types: What token type has X? Where X is a token index.
	TokenTypes TokenTypeLookup
	// Lookup tables for token modes: What token mode has X? Where X is a mode name (always contains "default").
	TokenModes map[string]*TokenMode
	// The order of token modes inside the grammar. Used mostly to keep the order of the generated code stable.
	TokenModeOrder []string
}

func (r GenerateTokenTypesResult) TokenTypeIds() map[string]int {
	tokenTypeIds := map[string]int{}
	for _, tokenType := range r.TokenTypes.All {
		// TokenIndex is the runtime token id emitted by the lexer (the _Idx
		// constant). Stores the token index for each token type name
		// (which is something like the token rule name). Especially for
		// keywords which can be standalone or part of a token declaration,
		// this is important to know which token index is used at runtime
		// (so: two entries can have the same token index).
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

// TokenTypeNamesByTokenIndex returns a slice mapping runtime token id
// (TokenIndex) to a token name. Keyword-backed tokens alias their keyword
// and share its id; later entries (token declarations, groups) override the
// earlier keyword entry so the referenced token name wins.
func (r GenerateTokenTypesResult) TokenTypeNamesByTokenIndex() []string {
	maxId := 0
	for _, tokenType := range r.TokenTypes.All {
		if tokenType.TokenIndex > maxId {
			maxId = tokenType.TokenIndex
		}
	}
	names := make([]string, maxId+1)
	for _, tokenType := range r.TokenTypes.All {
		names[tokenType.TokenIndex] = tokenType.Name
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
			ByKeyword:          make(map[string]int),
			ByToken:            make(map[grammar.TokenDecl]int),
			ByTokenGroup:       make(map[grammar.TokenGroup]int),
			ByTokenGroupParent: make(map[grammar.TokenGroup][]int),
			SourceType:         make(map[int]TokenIndexSource),
		},
		TokenTypes: TokenTypeLookup{
			All:          make([]TokenType, 0),
			ByTokenIndex: make(map[int]*TokenType),
		},
		TokenModes:     make(map[string]*TokenMode),
		TokenModeOrder: []string{},
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
			Id:      index,
			VarName: "TokenMode_" + modeName,
			TokenTypeIndices: TokenTypeIndices{
				Keywords: make([]int, 0),
				Tokens:   make([]int, 0),
			},
			TokenTypeUsages: make(map[int]TokenTypeUsage),
		}
		result.TokenModes[modeName] = &current
		// Tracked across all members of the mode: a keyword selector must not
		// register a keyword that an earlier member already added.
		alreadyAdded := map[int]bool{}
		for _, member := range tokenMode.Members() {
			pushTokenTypeUsage := func(tokenIndex int, tokenModifier string, command grammar.TokenCommand) {
				if _, ok := alreadyAdded[tokenIndex]; !ok {
					if source, ok := result.TokenIndex.SourceType[tokenIndex]; ok && source != SourceGroup {
						if source == SourceKeyword {
							current.TokenTypeIndices.Keywords = append(current.TokenTypeIndices.Keywords, tokenIndex)
						} else {
							current.TokenTypeIndices.Tokens = append(current.TokenTypeIndices.Tokens, tokenIndex)
						}
						alreadyAdded[tokenIndex] = true
					} else {
						return
					}
				}
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
				for _, tokenIndex := range result.TokenIndex.ByTokenGroupParent[tokenGroup] {
					pushTokenTypeUsage(tokenIndex, tokenGroup.Modifier(), tokenGroup.Command())
				}
			case grammar.KeywordUsage:
				tokenIndex := result.TokenIndex.ByKeyword[member.Keyword().Value()]
				pushTokenTypeUsage(tokenIndex, member.Modifier(), member.Command())
			case grammar.TokenUsage:
				token := member.TokenRef().Ref(context.Background())

				overrideUsage := func(tokenIndex int) {
					if member.Modifier() != "" || member.Command() != nil {
						current.TokenTypeUsages[tokenIndex] = TokenTypeUsage{
							TokenModifier: member.Modifier(),
							Command:       member.Command(),
						}
					}
				}

				var tokenIndex int
				switch rule := token.(type) {
				case grammar.TokenDecl:
					tokenIndex = result.TokenIndex.ByToken[rule]
					pushTokenTypeUsage(tokenIndex, rule.Modifier(), rule.Command())
					overrideUsage(tokenIndex)
				case grammar.TokenGroup:
					for _, tokenIndex := range result.TokenIndex.ByTokenGroupParent[rule] {
						pushTokenTypeUsage(tokenIndex, rule.Modifier(), rule.Command())
						overrideUsage(tokenIndex)
					}
				}
			case grammar.KeywordSelector:
				pattern := regexp.MustCompile(grammar.RegexpValue(member.Selector()))
				for _, keyword := range keywords.Keywords {
					tokenIndex := result.TokenIndex.ByKeyword[keyword.Value()]
					value := grammar.KeywordValue(keyword)
					if !alreadyAdded[tokenIndex] && pattern.MatchString(value) {
						current.TokenTypeIndices.Keywords = append(current.TokenTypeIndices.Keywords, tokenIndex)
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
			Id:      len(result.TokenModes),
			VarName: "TokenMode_default",
			TokenTypeIndices: TokenTypeIndices{
				Keywords: make([]int, 0),
				Tokens:   make([]int, 0),
			},
			TokenTypeUsages: make(map[int]TokenTypeUsage),
		}
		result.TokenModes["default"] = &defaultMode
		result.TokenModeOrder = append(result.TokenModeOrder, "default")

		for _, tokenGroup := range tokenGroups.TopLevel {
			for _, tokenIndex := range result.TokenIndex.ByTokenGroupParent[tokenGroup] {
				if source, ok := result.TokenIndex.SourceType[tokenIndex]; ok && source != SourceGroup {
					if source == SourceKeyword {
						if slices.Contains(defaultMode.TokenTypeIndices.Keywords, tokenIndex) {
							continue
						}
						defaultMode.TokenTypeIndices.Keywords = append(defaultMode.TokenTypeIndices.Keywords, tokenIndex)
					} else {
						if slices.Contains(defaultMode.TokenTypeIndices.Tokens, tokenIndex) {
							continue
						}
						defaultMode.TokenTypeIndices.Tokens = append(defaultMode.TokenTypeIndices.Tokens, tokenIndex)
					}
					if tokenGroup.Modifier() != "" || tokenGroup.Command() != nil {
						defaultMode.TokenTypeUsages[tokenIndex] = TokenTypeUsage{
							TokenModifier: tokenGroup.Modifier(),
							Command:       tokenGroup.Command(),
						}
					}
				}
			}
		}

		for _, keyword := range keywords.Keywords {
			tokenIndex := result.TokenIndex.ByKeyword[keyword.Value()]
			if slices.Contains(defaultMode.TokenTypeIndices.Keywords, tokenIndex) {
				continue
			}
			defaultMode.TokenTypeIndices.Keywords = append(defaultMode.TokenTypeIndices.Keywords, tokenIndex)
			//keywords don't have type or command, so we don't need to add anything to TokenTypeUsages
		}

		for _, token := range tokens.TopLevel {
			tokenIndex := result.TokenIndex.ByToken[token]
			if result.TokenIndex.SourceType[tokenIndex] != SourceTokenDecl {
				continue
			}
			if !slices.Contains(defaultMode.TokenTypeIndices.Tokens, tokenIndex) {
				defaultMode.TokenTypeIndices.Tokens = append(defaultMode.TokenTypeIndices.Tokens, tokenIndex)
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
		mergeImports(result.Imports, code.Imports)
		tokenType := TokenType{
			TokenIndex: tokenIndex,
			VarName:    GeneratedTokenName(keyword),
			Name:       keyword.Value(),
			Code:       code.Code,
		}
		result.TokenTypes.All = append(result.TokenTypes.All, tokenType)
		result.TokenTypes.ByTokenIndex[tokenIndex] = &tokenType
		result.TokenIndex.ByKeyword[keyword.Value()] = tokenIndex
		result.TokenIndex.SourceType[tokenIndex] = SourceKeyword
		tokenIndex++
	}
	for _, token := range tokens.All {
		varName := GeneratedTokenName(token)
		var currentTokenIndex int
		switch element := token.Content().(type) {
		case grammar.KeywordTokenContent:
			keywordIndex := keywords.ByValue[element.Keyword().Value()]
			keyword := keywords.Keywords[keywordIndex]
			code := codegen.NewNode()
			code.AppendLine("const ", GeneratedTokenIdxName(token), " = ", GeneratedTokenIdxName(keyword))
			code.AppendLine()
			code.AppendLine("var ", varName, " = ", GeneratedTokenName(keyword))
			mergeImports(result.Imports, map[string]bool{})
			currentTokenIndex = result.TokenIndex.ByKeyword[keyword.Value()]
			result.TokenIndex.SourceType[currentTokenIndex] = SourceTokenDecl
			tokenType := TokenType{
				TokenIndex: currentTokenIndex,
				VarName:    varName,
				Name:       token.Name(),
				Code:       code,
			}
			result.TokenTypes.All = append(result.TokenTypes.All, tokenType)
			result.TokenTypes.ByTokenIndex[currentTokenIndex] = &tokenType
			result.TokenIndex.ByToken[token] = currentTokenIndex
		case grammar.RegexpTokenContent:
			lexerResult := generateRegexpTokenElement(token, element, tokenIndex)
			mergeImports(result.Imports, lexerResult.Imports)
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
			result.TokenIndex.SourceType[currentTokenIndex] = SourceTokenDecl
			tokenIndex++
		}
	}
	tokenGroupMembers := map[string][]string{}
	for _, tokenGroup := range tokenGroups.All {
		tokenGroupMembers[tokenGroup.Name()] = getAllTokenGroupMembers(tokenGroup, keywords)
	}
	// Token groups need to be topologically sorted, so that nested groups appear after their members
	sortedTokenGroups := sortTokenGroups(tokenGroups.All, tokenGroupMembers)
	for _, tokenGroup := range sortedTokenGroups {
		lexerResult := generateTokenGroupType(tokenGroup, tokenGroupMembers, tokenIndex)
		mergeImports(result.Imports, lexerResult.Imports)
		tokenType := TokenType{
			TokenIndex: tokenIndex,
			VarName:    GeneratedTokenName(tokenGroup),
			Name:       tokenGroup.Name(),
			Code:       lexerResult.Code,
		}
		result.TokenTypes.All = append(result.TokenTypes.All, tokenType)
		result.TokenTypes.ByTokenIndex[tokenIndex] = &tokenType
		result.TokenIndex.ByTokenGroup[tokenGroup] = tokenIndex
		result.TokenIndex.SourceType[tokenIndex] = SourceGroup
		tokenIndex++
	}
	result.TokenIndex.ByTokenGroupParent = getAllTokenGroupMemberTokenIndices(sortedTokenGroups, keywords, result.TokenIndex)
}

func mergeImports(target map[string]bool, source map[string]bool) {
	maps.Copy(target, source)
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

func getAllTokenGroupMemberTokenIndices(sortedTokenGroups []grammar.TokenGroup, keywords GetAllKeywordsResult, lookup TokenIndexLookup) map[grammar.TokenGroup][]int {
	tokenGroupMembers := map[grammar.TokenGroup][]int{}
	for _, tokenGroup := range sortedTokenGroups {
		tokenGroupMembers[tokenGroup] = []int{}
	}
	for _, tokenGroup := range sortedTokenGroups {
		for _, tokenRef := range tokenGroup.TokenRefs() {
			tokenRule := tokenRef.Ref(context.Background())
			if tokenRule != nil {
				if subTokenGroup, ok := tokenRule.(grammar.TokenGroup); ok {
					tokenGroupMembers[tokenGroup] = append(tokenGroupMembers[tokenGroup], tokenGroupMembers[subTokenGroup]...)
				} else if tokenDecl, ok := tokenRule.(grammar.TokenDecl); ok {
					if idx, ok := lookup.ByToken[tokenDecl]; ok {
						tokenGroupMembers[tokenGroup] = append(tokenGroupMembers[tokenGroup], idx)
					}
				}
			}
		}
		for _, selector := range tokenGroup.KeywordSelectors() {
			pattern := regexp.MustCompile(grammar.RegexpValue(selector.Image))
			for _, keyword := range keywords.Keywords {
				if idx, ok := lookup.ByKeyword[keyword.Value()]; ok {
					if pattern.MatchString(keyword.Value()) {
						tokenGroupMembers[tokenGroup] = append(tokenGroupMembers[tokenGroup], idx)
					}
				}
			}
		}
		for _, keyword := range tokenGroup.Keywords() {
			if idx, ok := lookup.ByKeyword[keyword.Value()]; ok {
				tokenGroupMembers[tokenGroup] = append(tokenGroupMembers[tokenGroup], idx)
			}
		}
	}
	return tokenGroupMembers
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
