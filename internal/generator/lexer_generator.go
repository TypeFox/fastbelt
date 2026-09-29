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
	"strings"
	"unicode/utf8"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/internal/automatons"
	"typefox.dev/fastbelt/internal/grammar"
	fbRegexp "typefox.dev/fastbelt/internal/regexp"
	"typefox.dev/fastbelt/util/codegen"
)

// GenerateLexer emits the lexer for grammr. entryRules lists one entry rule per
// language and startModes the token mode each language starts in (index
// aligned; nil or "" means the default mode). See lexerModes for the token
// modes that end up in the lexer.
func GenerateLexer(grammr grammar.Grammar, entryRules []grammar.ParserRule, startModes []string, packageName string, tokenTypes GenerateTokenTypesResult) string {
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
		n.AppendLine("\"typefox.dev/fastbelt/util/service\"")
	})
	node.AppendLine(")")
	node.AppendLine()

	for _, n := range nodes {
		node.AppendNode(n)
		node.AppendLine()
	}

	modes, starts := lexerModes(grammr, entryRules, startModes, tokenTypes)
	generateLexerModeEnums(node, modes)
	generateMainLexerFunction(context.Background(), node, modes, starts, tokenTypes)
	return FormatIfPossible(node.String())
}

// lexerMode is one token mode of the generated lexer.
type lexerMode struct {
	// name is the display name of the mode.
	name string
	mode *TokenMode
	// reachable limits the mode to the token types with these var names, plus
	// those carrying a modifier or mode command (hidden/comment tokens and mode
	// switches must be lexed regardless of the parser rules). Nil keeps every
	// token type of the mode.
	reachable map[string]bool
}

// lexerModes returns the token modes of the generated lexer and, per language,
// the var name of the mode its lexer starts in.
//
// Token modes declared in the grammar are emitted as they are: the grammar
// author decides which tokens a language sees by giving it a start mode. A
// grammar without token modes that serves several languages gets one synthetic
// mode per language instead, named after the entry rule and holding only the
// token types reachable from it, so that a keyword of one language stays an
// ordinary identifier in the others.
func lexerModes(grammr grammar.Grammar, entryRules []grammar.ParserRule, startModes []string, tokenTypes GenerateTokenTypesResult) ([]lexerMode, []string) {
	declared := len(grammr.TokenModes()) > 0
	if !declared && len(entryRules) > 1 {
		modes := make([]lexerMode, len(entryRules))
		starts := make([]string, len(entryRules))
		for i, entry := range entryRules {
			mode := *tokenTypes.TokenModes["default"]
			mode.Id = i
			mode.VarName = "TokenMode_" + entry.Name()
			modes[i] = lexerMode{name: entry.Name(), mode: &mode, reachable: reachableTokenNames(grammr, entry)}
			starts[i] = mode.VarName
		}
		return modes, starts
	}

	hasDefault := !declared
	for _, mode := range grammr.TokenModes() {
		hasDefault = hasDefault || mode.IsDefault()
	}
	modes := make([]lexerMode, 0, len(tokenTypes.TokenModeOrder))
	byName := map[string]*TokenMode{}
	for _, name := range tokenTypes.TokenModeOrder {
		if name == "default" && !hasDefault {
			// The implicit default mode is only a fallback for grammars
			// without token modes; here nothing could enter it.
			continue
		}
		mode := tokenTypes.TokenModes[name]
		modes = append(modes, lexerMode{name: name, mode: mode})
		byName[name] = mode
	}
	starts := make([]string, max(1, len(entryRules)))
	for i := range starts {
		name := "default"
		if i < len(startModes) && startModes[i] != "" {
			name = startModes[i]
		}
		mode := byName[name]
		if mode == nil {
			// Rejected by validation or the build; keep the output compilable.
			mode = modes[0].mode
		}
		starts[i] = mode.VarName
	}
	return modes, starts
}

func generateLexerModeEnums(node codegen.Node, modes []lexerMode) {
	node.AppendLine("const (")
	node.Indent(func(n codegen.Node) {
		for _, m := range modes {
			n.AppendLine(m.mode.VarName, " = ", strconv.Itoa(m.mode.Id))
		}
	})
	node.AppendLine(")")
	node.AppendLine()
}

func generateMainLexerFunction(context context.Context, node codegen.Node, modes []lexerMode, starts []string, tokenTypes GenerateTokenTypesResult) {
	node.AppendLine("func NewLexer(sc *service.Container) lexer.Lexer {")
	node.Indent(func(n codegen.Node) {
		n.AppendLine("modes := make([]*lexer.TokenMode, ", strconv.Itoa(len(modes)), ")")
		for _, m := range modes {
			generateTokenMode(context, n, m, modes, tokenTypes)
		}
		if len(starts) == 1 {
			n.AppendLine("return lexer.NewDefaultLexer(sc, ", starts[0], ", modes...)")
			return
		}
		// One start mode per language, index-aligned with the parser's entry
		// dispatch.
		n.AppendLine("return lexer.NewMultiLanguageLexer(sc, []int{", strings.Join(starts, ", "), "}, modes...)")
	})
	node.AppendLine("}")
}

// generateTokenMode emits the NewTokenMode call of m into the modes slice.
func generateTokenMode(context context.Context, n codegen.Node, m lexerMode, modes []lexerMode, tokenTypes GenerateTokenTypesResult) {
	// Mode commands resolve their target among the emitted modes. A synthetic
	// language mode stands in for the default mode of its language.
	targets := map[string]*TokenMode{}
	if m.reachable != nil {
		targets["default"] = m.mode
	} else {
		for _, other := range modes {
			targets[other.name] = other.mode
		}
	}
	tokenTypes.TokenModes = targets

	n.AppendLine("modes[", m.mode.VarName, "] = lexer.NewTokenMode(\"", m.name, "\",")
	n.Indent(func(nn codegen.Node) {
		for _, tokenIndex := range slices.Concat(m.mode.ModeTokenTypes.Keywords, m.mode.ModeTokenTypes.Tokens) {
			tokenType := tokenTypes.TokenTypes.ByTokenIndex[tokenIndex]
			if m.reachable == nil || m.reachable[tokenType.VarName] || m.mode.TokenTypeUsages[tokenIndex] != (tokenTypeUsage{}) {
				generateTokenTypeUsage(context, nn, tokenType, m.mode, tokenIndex, tokenTypes)
			}
		}
	})
	n.AppendLine(")")
}

// reachableTokenNames returns the generated var names (Keyword_*/Token_*) of
// all token types referenced by the language rooted at entry: every keyword
// and token declaration reachable through rule bodies, rule calls,
// cross-references and token groups.
func reachableTokenNames(grammr grammar.Grammar, entry grammar.ParserRule) map[string]bool {
	ctx := context.Background()
	reachable := map[string]bool{}
	visitedRules := map[string]bool{}
	visitedGroups := map[string]bool{}
	keywords := GetAllKeywords(grammr)

	var visitGroup func(tokenGroup grammar.TokenGroup)
	visitGroup = func(tokenGroup grammar.TokenGroup) {
		if visitedGroups[tokenGroup.Name()] {
			return
		}
		visitedGroups[tokenGroup.Name()] = true
		for _, tokenRef := range tokenGroup.TokenRefs() {
			switch target := tokenRef.Ref(ctx).(type) {
			case grammar.TokenGroup:
				visitGroup(target)
			case grammar.TokenDecl:
				reachable[GeneratedTokenName(target)] = true
			}
		}
		for _, name := range getAllTokenGroupMembers(tokenGroup, keywords) {
			reachable[name] = true
		}
	}

	var visitRule func(rule core.AstNode)
	visitRule = func(rule core.AstNode) {
		for node := range core.AllChildren(rule) {
			switch n := node.(type) {
			case grammar.Keyword:
				reachable[GeneratedTokenName(n)] = true
			case grammar.RuleCall:
				if n.Rule() == nil {
					continue
				}
				switch target := n.Rule().Ref(ctx).(type) {
				case grammar.TokenDecl:
					reachable[GeneratedTokenName(target)] = true
				case grammar.TokenGroup:
					visitGroup(target)
				case grammar.ParserRule, grammar.CompositeRule, grammar.InfixRule:
					// An infix rule's operators (keywords or token calls) and its
					// operand call are children of the rule node like any body.
					if !visitedRules[target.Name()] {
						visitedRules[target.Name()] = true
						visitRule(target)
					}
				}
			}
		}
	}
	visitedRules[entry.Name()] = true
	visitRule(entry)
	return reachable
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
			// A push/mode command without a resolvable target mode is reported
			// by validation; skip it here instead of emitting a broken lexer.
			targetMode := tokenTypes.TokenModes[cmdModeName]
			switch {
			case cmd.Type() == "pop":
				nn.Append(".WithPopMode()")
			case targetMode == nil:
				// no target mode to switch to
			case cmd.Type() == "push":
				nn.Append(".WithPushMode(", targetMode.VarName, ")")
			default: //"mode"
				nn.Append(".WithSetMode(", targetMode.VarName, ")")
			}
		}
		switch usage.TokenModifier {
		case "comment":
			nn.Append(".WithModifier(core.CommentModifier)")
		case "hidden":
			nn.Append(".WithModifier(core.SkippedModifier)")
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

func generateRegexpTokenElement(token grammar.TokenDecl, regexpTokenElement grammar.RegexpTokenContent, id int) GenerateLexerResult {
	var result fbRegexp.GenerateRegExpResult
	imports := map[string]bool{}
	code := codegen.NewNode()
	regexPattern := grammar.RegexpValue(regexpTokenElement.Regexp())
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
