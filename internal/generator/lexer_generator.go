// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"context"
	"fmt"
	"maps"
	goregex "regexp"
	"slices"
	"sort"
	"strconv"
	"unicode/utf8"

	"typefox.dev/fastbelt/internal/automatons"
	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/internal/regexp"
	"typefox.dev/fastbelt/util/codegen"
)

type GenerateTokenTypesResult struct {
	Tokens            []grammar.Token
	Keywords          []grammar.Keyword
	Imports           map[string]bool
	KeywordsCode      []codegen.Node
	TokensCode        []codegen.Node
	TokenGroupCode    []codegen.Node
	TokenTypeVarNames []string
	TokenTypeNames    []string
	TokenTypeIds      map[string]int
}

func GenerateTokenTypes(grammr grammar.Grammar) GenerateTokenTypesResult {
	tokens := grammr.Terminals()
	tokenGroups := grammr.TokenGroups()
	keywords := GetAllKeywords(grammr)
	keywordsCount := len(keywords)
	result := GenerateTokenTypesResult{
		Tokens:            tokens,
		Keywords:          keywords,
		KeywordsCode:      make([]codegen.Node, keywordsCount),
		TokensCode:        make([]codegen.Node, len(tokens)),
		TokenGroupCode:    make([]codegen.Node, len(tokenGroups)),
		TokenTypeNames:    make([]string, keywordsCount+len(tokens)+len(tokenGroups)),
		TokenTypeVarNames: make([]string, keywordsCount+len(tokens)+len(tokenGroups)),
		TokenTypeIds:      make(map[string]int),
		Imports:           map[string]bool{},
	}
	// Starting with 1 - prevent clash with EOF (index 0)
	tokenIndex := 1
	for index, keyword := range keywords {
		result.KeywordsCode[index] = generateKeywordTokenType(keyword, tokenIndex)
		varName := GeneratedTokenName(keyword)
		kwName := keyword.Value()
		result.TokenTypeVarNames[index] = varName
		result.TokenTypeNames[index] = kwName
		result.TokenTypeIds[kwName] = index
		result.Imports["strings"] = true
		tokenIndex++
	}
	for index, token := range tokens {
		tokenType := generateTokenType(token, tokenIndex)
		result.TokensCode[index] = tokenType.Code
		for imp := range tokenType.Imports {
			result.Imports[imp] = true
		}
		varName := GeneratedTokenName(token)
		tokName := token.Name()
		result.TokenTypeVarNames[keywordsCount+index] = varName
		result.TokenTypeNames[keywordsCount+index] = tokName
		result.TokenTypeIds[tokName] = keywordsCount + index
		tokenIndex++
	}
	tokenGroupMembers := map[string][]string{}
	for _, tokenGroup := range tokenGroups {
		tokenGroupMembers[tokenGroup.Name()] = getAllTokenGroupMembers(tokenGroup, keywords)
	}
	// Token groups need to be topologically sorted, so that nested groups appear after their members
	for index, tokenGroup := range sortTokenGroups(tokenGroups, tokenGroupMembers) {
		result.TokenGroupCode[index] = generateTokenGroupType(tokenGroup, tokenGroupMembers, tokenIndex)
		varName := GeneratedTokenName(tokenGroup)
		tokName := tokenGroup.Name()
		result.TokenTypeVarNames[keywordsCount+len(tokens)+index] = varName
		result.TokenTypeNames[keywordsCount+len(tokens)+index] = tokName
		result.TokenTypeIds[tokName] = keywordsCount + len(tokens) + index
		tokenIndex++
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

	generateMainLexerFunction(node, tokenTypes.Tokens, tokenTypes.Keywords)
	return FormatIfPossible(node.String())
}

func generateMainLexerFunction(node codegen.Node, tokens []grammar.Token, keywords []grammar.Keyword) {
	node.AppendLine("func NewLexer() lexer.Lexer {")
	node.Indent(func(n codegen.Node) {
		n.AppendLine("return lexer.NewDefaultLexer(")
		n.Indent(func(nn codegen.Node) {
			for _, keyword := range keywords {
				nn.AppendLine(GeneratedTokenName(keyword), ",")
			}
			for _, token := range tokens {
				nn.AppendLine(GeneratedTokenName(token), ",")
			}
		})
		n.AppendLine(")")
	})
	node.AppendLine("}")
}

func generateKeywordTokenType(keyword grammar.Keyword, id int) codegen.Node {
	code := codegen.NewNode()
	keywordValue := grammar.KeywordValue(keyword)
	code.AppendLine("const ", GeneratedTokenIdxName(keyword), " = ", strconv.Itoa(id))
	code.AppendLine()
	code.AppendLine("var ", GeneratedTokenName(keyword), " = core.NewTokenType(")
	code.Indent(func(n codegen.Node) {
		n.AppendLine(GeneratedTokenIdxName(keyword), ",")
		n.AppendLine("\"", keywordValue, "\",")
		n.AppendLine("\"", keywordValue, "\",")
		n.AppendLine("0,")
		n.AppendLine("core.TokenKindKeyword,")
		n.AppendLine("0,")
		n.AppendLine("false,")
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
	return code
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

func getAllTokenGroupMembers(tokenGroup grammar.TokenGroup, keywords []grammar.Keyword) []string {
	members := map[string]bool{}
	for _, tokenRef := range tokenGroup.TokenRefs() {
		tokenRule := tokenRef.Ref(context.Background())
		if tokenRule != nil {
			name := GeneratedTokenName(tokenRule)
			members[name] = true
		}
	}
	for _, keyword := range tokenGroup.Keywords() {
		name := GeneratedTokenName(keyword)
		members[name] = true
	}
	for _, regex := range tokenGroup.Regexps() {
		pattern := regex.Image[1 : len(regex.Image)-1]
		compiled, err := goregex.Compile(pattern)
		if err != nil {
			continue
		}
		for _, keyword := range keywords {
			keywordValue := KeywordValue(keyword)
			keywordName := GeneratedTokenName(keyword)
			matches := compiled.MatchString(keywordValue)
			if matches {
				members[keywordName] = true
			}
		}
	}
	slice := slices.Collect(maps.Keys(members))
	sort.Strings(slice)
	return slice
}

func generateTokenType(token grammar.Token, id int) GenerateLexerResult {
	var result regexp.GenerateRegExpResult
	imports := map[string]bool{}
	code := codegen.NewNode()
	regexPattern := token.Regexp()
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
		if token.Type() == "hidden" {
			n.AppendLine("core.SkippedGroup,")
		} else if token.Type() == "comment" {
			n.AppendLine("core.CommentGroup,")
		} else {
			n.AppendLine("0,")
		}
		n.AppendLine("core.TokenKindToken,")
		n.AppendLine("0,")
		n.AppendLine("false,")
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
