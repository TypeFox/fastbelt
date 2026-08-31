// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"unicode"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/collections"
)

const (
	ValidateUniqueRuleName                   = "uniqueRuleName"
	ValidateUniqueRuleNameInTokenMode        = "uniqueRuleNameInTokenMode"
	ValidateUniqueInterfaceName              = "uniqueInterfaceName"
	ValidateUniqueTokenModeName              = "uniqueTokenModeName"
	ValidateEmptyToken                       = "emptyTerminalRule"
	ValidateEmptyKeyword                     = "emptyKeyword"
	ValidateWhitespaceOnlyKeyword            = "whitespaceOnlyKeyword"
	ValidateKeywordWithWhitespace            = "keywordWithWhitespace"
	ValidateRuleReturnType                   = "ruleReturnType"
	ValidateInterfaceExtends                 = "interfaceExtends"
	ValidateRuleCallReturnType               = "ruleCallReturnType"
	ValidateRuleCallPosition                 = "ruleCallPosition"
	ValidateActionAssignmentType             = "actionAssignmentType"
	ValidateActionPropertyType               = "actionPropertyType"
	ValidateAssignmentType                   = "assignmentType"
	ValidateRecursiveTokenGroup              = "recursiveTokenGroup"
	ValidateInvalidTokenInGroup              = "invalidTokenInGroup"
	ValidateInvalidTokenInCrossRef           = "invalidTokenInCrossRef"
	ValidateMissingCrossRefTerminal          = "missingCrossRefTerminal"
	ValidateUniqueFieldName                  = "uniqueFieldName"
	ValidateFieldNameCapitalLetter           = "fieldNameCapitalLetter"
	ValidateReservedFieldName                = "reservedFieldName"
	ValidateNestedArrayType                  = "nestedArrayType"
	ValidateDefaultTokenModeRequired         = "defaultTokenModeRequired"
	ValidateTokenCommandMode                 = "tokenCommandMode"
	ValidateEmptyTokenMode                   = "emptyTokenMode"
	ValidateUnreachableTokenMode             = "unreachableTokenMode"
	ValidateKeywordNotInTokenMode            = "keywordNotInTokenMode"
	ValidateTokenNotInTokenMode              = "tokenNotInTokenMode"
	ValidateNonDefaultTokenModeNoPop         = "nonDefaultTokenModeNoPop"
	ValidateTerminalNotCoveredByParserRule   = "terminalNotCoveredByParserRule"
	ValidateTokenGroupNotCoveredByParserRule = "tokenGroupNotCoveredByParserRule"
	ValidateInvalidRegExpLiteral             = "invalidRegExpLiteral"
)

// defaultTokenModeName is the name under which the mode marked
// `token mode default` is registered. It is the mode the lexer starts in.
const defaultTokenModeName = "default"

// reservedFieldNames lists field names that must not be used because they
// conflict with [core.AstNode] methods. Keep in sync with ast.go.
var reservedFieldNames = map[string]string{
	"Document":         "AstNode.Document",
	"Container":        "AstNode.Container",
	"ContainmentData":  "AstNode.ContainmentData",
	"Tokens":           "AstNode.Tokens",
	"TextRange":        "AstNode.TextRange",
	"Text":             "AstNode.Text",
	"ForEachNode":      "AstNode.ForEachNode",
	"ForEachReference": "AstNode.ForEachReference",
	"Resolve":          "AstNode.Resolve",
}

// GrammarImpl.Validate checks grammar-level constraints
func (g *GrammarImpl) Validate(ctx context.Context, _ string, accept core.ValidationAcceptor) {
	checkUniqueRuleNames(g, accept)
	checkUniqueInterfaceNames(g, accept)
	checkUniqueTokenModeNames(g, accept)
	checkIfDefaultTokenModeIsRequired(g, accept)
	checkTokenModesAreReachable(g, ctx, accept)
	checkTokenModesCoverParserTokens(g, ctx, accept)
	checkParserRulesCoverVisibleTokens(g, ctx, accept)
	checkIfNonDefaultTokenModesHasNoExit(g, ctx, accept)
}

// tokenModeName returns the name a token mode is registered under. The mode
// marked `default` has no name of its own.
func tokenModeName(mode TokenMode) string {
	if mode.IsDefault() {
		return defaultTokenModeName
	}
	return mode.Name()
}

// tokenModeNameToken returns the token to anchor diagnostics about mode itself
// on, which is either its name or the `default` marker.
func tokenModeNameToken(mode TokenMode) *core.Token {
	if mode.IsDefault() {
		return mode.DefaultToken()
	}
	return mode.NameToken()
}

// checkTokenModesAreReachable reports token modes that no command switches to.
// The lexer starts in the default mode and can only leave it through a `push` or
// `mode` command, so a mode nothing targets is dead weight.
func checkTokenModesAreReachable(g Grammar, ctx context.Context, accept core.ValidationAcceptor) {
	if len(g.TokenModes()) == 0 {
		return
	}
	var entryMode TokenMode = nil
	transitions := map[TokenMode][]TokenMode{}
	for _, mode := range g.TokenModes() {
		if mode.IsDefault() {
			entryMode = mode
		}
		for _, member := range mode.Members() {
			command := getCommand(member)
			if command != nil && command.Mode() != nil {
				transitions[mode] = append(transitions[mode], command.Mode().Ref(ctx))
			}
		}
	}
	visited := collections.NewSet[TokenMode]()
	queue := []TokenMode{}
	if entryMode != nil {
		queue = append(queue, entryMode)
	}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		visited.Add(current)
		for _, next := range transitions[current] {
			if !visited.Has(next) {
				queue = append(queue, next)
			}
		}
	}
	for _, mode := range g.TokenModes() {
		if visited.Has(mode) {
			continue
		}
		accept(core.NewDiagnostic(
			core.SeverityWarning,
			fmt.Sprintf("The token mode '%s' is never entered from the default mode. Add a 'push(%s)' or 'mode(%s)' command to a token.", mode.Name(), mode.Name(), mode.Name()),
			mode,
			core.WithToken(mode.NameToken()),
			core.WithCode(ValidateUnreachableTokenMode),
		))
	}
}

func checkIfNonDefaultTokenModesHasNoExit(g Grammar, ctx context.Context, accept core.ValidationAcceptor) {
	if len(g.TokenModes()) == 0 {
		return
	}
	for _, mode := range g.TokenModes() {
		if mode.IsDefault() {
			continue
		}
		found := false
		for _, member := range mode.Members() {
			command := getCommand(member)
			if command != nil && command.Type() != tokenCommandPush {
				found = true
			}
		}
		if !found {
			accept(core.NewDiagnostic(
				core.SeverityWarning,
				fmt.Sprintf("The non-default token mode '%s' has no exit command ('mode' or 'pop').", mode.Name()),
				mode,
				core.WithToken(mode.NameToken()),
				core.WithCode(ValidateNonDefaultTokenModeNoPop),
			))
		}
	}
}

func getCommand(member TokenModeMember) TokenCommand {
	switch casted := member.(type) {
	case TokenDeclUsage:
		return casted.Declaration().Command()
	case TokenGroupUsage:
		return casted.Group().Command()
	case TokenUsage:
		command := casted.Command()
		if command == nil {
			ref := casted.TokenRef().Ref(context.Background())
			if ref != nil {
				command = ref.Command()
			}
		}
		return command
	case KeywordUsage:
		return casted.Command()
	case KeywordSelector:
		return nil
	}
	return nil
}

func checkIfDefaultTokenModeIsRequired(g Grammar, accept core.ValidationAcceptor) {
	if len(g.TokenModes()) > 0 {
		hasDefault := false
		var nonDefaultTokenMode TokenMode = nil
		for _, mode := range g.TokenModes() {
			if mode.IsDefault() {
				hasDefault = true
				break
			} else if nonDefaultTokenMode == nil {
				//mark only the first non-default token mode
				//one diagnostic is enough to indicate that a default token mode is required
				nonDefaultTokenMode = mode
			}
		}
		if !hasDefault {
			accept(core.NewDiagnostic(
				core.SeverityError,
				"At least one token mode must be marked as default.",
				nonDefaultTokenMode,
				core.WithToken(nonDefaultTokenMode.NameToken()),
				core.WithCode(ValidateDefaultTokenModeRequired),
			))
		}
	}
}

func checkUniqueTokenModeNames(g Grammar, accept core.ValidationAcceptor) {
	seen := map[string][]TokenMode{}
	for _, mode := range g.TokenModes() {
		name := tokenModeName(mode)
		seen[name] = append(seen[name], mode)
	}
	for name, modes := range seen {
		if len(modes) > 1 {
			for _, mode := range modes {
				token := tokenModeNameToken(mode)
				accept(core.NewDiagnostic(
					core.SeverityError,
					fmt.Sprintf("A token mode's name has to be unique. '%s' is used multiple times.", name),
					mode,
					core.WithToken(token),
					core.WithCode(ValidateUniqueTokenModeName),
				))
			}
		}
	}
}

// Token command types as they appear after the `->` arrow.
const (
	tokenCommandPush = "push"
	tokenCommandPop  = "pop"
	tokenCommandMode = "mode"
)

// TokenCommandImpl.Validate checks lexer mode switch constraints:
//   - `push` and `mode` need a target mode, otherwise there is nothing to
//     switch to and the command is dropped during code generation.
//   - `pop` returns to the mode below the current one on the stack, so a target
//     mode cannot be honored.
func (c *TokenCommandImpl) Validate(_ context.Context, _ string, accept core.ValidationAcceptor) {
	checkTokenCommandMode(c, accept)
}

func checkTokenCommandMode(c TokenCommand, accept core.ValidationAcceptor) {
	// An unresolvable mode reference is still a target; the linker reports it.
	hasMode := c.IsDefault() || c.Mode() != nil
	switch c.Type() {
	case tokenCommandPop:
		if !hasMode {
			return
		}
		options := []core.DiagnosticOption{core.WithCode(ValidateTokenCommandMode)}
		if c.Mode() != nil {
			options = append(options, core.WithReference(c.Mode()))
		} else {
			options = append(options, core.WithToken(c.DefaultToken()))
		}
		accept(core.NewDiagnostic(
			core.SeverityError,
			"The 'pop' command returns to the previous token mode and cannot take a target mode.",
			c,
			options...,
		))
	case tokenCommandPush, tokenCommandMode:
		if hasMode {
			return
		}
		accept(core.NewDiagnostic(
			core.SeverityError,
			fmt.Sprintf("The '%s' command requires a target token mode, for example '%s(default)'.", c.Type(), c.Type()),
			c,
			core.WithToken(c.TypeToken()),
			core.WithCode(ValidateTokenCommandMode),
		))
	}
}

// TokenModeImpl.Validate checks token mode constraints:
//   - A mode without members leaves the lexer with nothing to match.
func (m *TokenModeImpl) Validate(_ context.Context, _ string, accept core.ValidationAcceptor) {
	checkTokenModeNotEmpty(m, accept)
	checkTokenModeMembersAreUnique(m, accept)
}

func checkTokenModeNotEmpty(mode TokenMode, accept core.ValidationAcceptor) {
	if len(mode.Members()) > 0 {
		return
	}
	accept(core.NewDiagnostic(
		core.SeverityWarning,
		fmt.Sprintf("The token mode '%s' does not declare any token, so the lexer cannot match anything while it is active.", tokenModeName(mode)),
		mode,
		core.WithToken(tokenModeNameToken(mode)),
		core.WithCode(ValidateEmptyTokenMode),
	))
}

func getTokenNeverReferencedMessage(tokenTypeName string, tokenValue string) string {
	return fmt.Sprintf("The %s '%s' is never referenced in a parser rule, so the lexer can never produce it.", tokenTypeName, tokenValue)
}

func checkParserRulesCoverVisibleTokens(g Grammar, ctx context.Context, accept core.ValidationAcceptor) {
	severity := core.SeverityWarning
	seen := collections.NewSet[string]()
	queue := []core.AstNode{}
	for _, composite := range g.Composites() {
		for node := range core.AllChildren(composite) {
			queue = append(queue, node)
		}
	}
	for _, rule := range g.Rules() {
		for node := range core.AllChildren(rule) {
			queue = append(queue, node)
		}
	}
	for _, group := range g.TokenGroups() {
		for node := range core.AllChildren(group) {
			queue = append(queue, node)
		}
	}
	for _, node := range queue {
		switch casted := node.(type) {
		case Keyword:
			seen.Add(casted.Value())
		case RuleCall:
			rule, ok := casted.Rule().Ref(ctx).(AbstractTokenRule)
			if ok {
				seen.Add(rule.Name())
			}
		}
	}

	if len(g.TokenModes()) == 0 {
		for _, terminal := range g.Terminals() {
			if terminal.Modifier() != "" {
				continue
			}
			if !seen.Has(terminal.Name()) {
				accept(core.NewDiagnostic(
					severity,
					getTokenNeverReferencedMessage("token", terminal.Name()),
					terminal,
					core.WithToken(terminal.NameToken()),
					core.WithCode(ValidateTerminalNotCoveredByParserRule),
				))
			}
		}
		for _, group := range g.TokenGroups() {
			if group.Modifier() != "" {
				continue
			}
			if !seen.Has(group.Name()) {
				accept(core.NewDiagnostic(
					severity,
					getTokenNeverReferencedMessage("token group", group.Name()),
					group,
					core.WithToken(group.NameToken()),
					core.WithCode(ValidateTokenGroupNotCoveredByParserRule),
				))
			}
		}
	} else {
		for _, mode := range g.TokenModes() {
			for _, member := range mode.Members() {
				switch member := member.(type) {
				case KeywordUsage:
					if !seen.Has(member.Keyword().Value()) {
						accept(core.NewDiagnostic(
							severity,
							getTokenNeverReferencedMessage("keyword", member.Keyword().Value()),
							member.Keyword(),
							core.WithToken(member.Keyword().ValueToken()),
							core.WithCode(ValidateTerminalNotCoveredByParserRule),
						))
					}
				case TokenUsage:
					if tokenRef := member.TokenRef().Ref(context.Background()); tokenRef != nil {
						if tokenRef.Modifier() != "" {
							continue
						}
						if !seen.Has(tokenRef.Name()) {
							accept(core.NewDiagnostic(
								severity,
								getTokenNeverReferencedMessage("token", tokenRef.Name()),
								member,
								core.WithCode(ValidateTerminalNotCoveredByParserRule),
							))
						}
					}
				case TokenDeclUsage:
					if member.Declaration().Modifier() != "" {
						continue
					}
					if !seen.Has(member.Declaration().Name()) {
						accept(core.NewDiagnostic(
							severity,
							getTokenNeverReferencedMessage("token", member.Declaration().Name()),
							member.Declaration(),
							core.WithToken(member.Declaration().NameToken()),
							core.WithCode(ValidateTerminalNotCoveredByParserRule),
						))
					}
				case TokenGroupUsage:
					if member.Group().Modifier() != "" {
						continue
					}
					if !seen.Has(member.Group().Name()) {
						accept(core.NewDiagnostic(
							severity,
							getTokenNeverReferencedMessage("token group", member.Group().Name()),
							member.Group(),
							core.WithToken(member.Group().NameToken()),
							core.WithCode(ValidateTokenGroupNotCoveredByParserRule),
						))
					}
				}
			}
		}
	}
}

// tokenModeCoverage records the keywords and token rules that the declared token
// modes register with the lexer.
type tokenModeCoverage struct {
	keywords collections.Set[string]
	rules    collections.Set[string]
}

// checkTokenModesCoverParserTokens reports keywords and tokens that the parser
// can demand but that no token mode registers. Declaring token modes takes over
// token registration from the grammar as a whole, so anything left out of every
// mode can never be produced by the lexer and makes the rule using it dead.
//
// Only the first occurrence of each keyword or token is reported: the fix is a
// single entry in a token mode, not one per use site.
func checkTokenModesCoverParserTokens(g Grammar, ctx context.Context, accept core.ValidationAcceptor) {
	if len(g.TokenModes()) == 0 {
		// Without explicit token modes every keyword and token is registered
		// automatically, so nothing can be missing.
		return
	}
	coverage := collectTokenModeCoverage(g, ctx)
	reported := collections.NewSet[string]()
outerLoop:
	for node := range core.AllChildren(g) {
		switch node := node.(type) {
		case Keyword:
			value := node.Value()
			if !insideParserRule(node) || coverage.keywords.Has(value) || !reported.Add(value) {
				continue
			}
			accept(core.NewDiagnostic(
				core.SeverityError,
				fmt.Sprintf("The keyword %s is not registered in any token mode, so the lexer can never produce it. List it in a token mode or cover it with a 'keywords' selector.", value),
				node,
				core.WithToken(node.ValueToken()),
				core.WithCode(ValidateKeywordNotInTokenMode),
			))
		case RuleCall:
			rule, ok := node.Rule().Ref(ctx).(AbstractTokenRule)
			if !ok || !insideParserRule(node) || coverage.covers(rule) || !reported.Add(rule.Name()) {
				continue
			}
			if tokenGroup, ok := rule.(TokenGroup); ok {
				//Ignore token groups FTM, since they are meant to be
				//used as a grouped token type during parsing, not during lexing
				//Instead test the individual token types in the group for coverage
				//If all token types in the group are not covered, then the individual token group will be reported as missing
				for _, member := range tokenGroup.TokenRefs() {
					if member := member.Ref(ctx); member != nil && coverage.covers(member) {
						continue outerLoop
					}
				}
				for _, member := range tokenGroup.Keywords() {
					value := KeywordValue(member)
					if coverage.keywords.Has(value) {
						continue outerLoop
					}
				}
				for _, member := range tokenGroup.KeywordSelectors() {
					pattern, err := regexp.Compile(RegexpValue(member.Image))
					if err != nil {
						//error is already reported by the grammar validator, so ignore it here
						continue
					}
					for _, keyword := range allGrammarKeywords(g) {
						value := KeywordValue(keyword)
						if pattern.MatchString(value) && coverage.keywords.Has(value) {
							continue outerLoop
						}
					}
				}
			}
			// Deliberately a warning, not an error: unlike a keyword, a token
			// rule may be registered with the lexer by hand-written code, and a
			// grammar that leaves one out still generates and builds.
			accept(core.NewDiagnostic(
				core.SeverityWarning,
				fmt.Sprintf("The token '%s' is not registered in any token mode, so the lexer can never produce it. List it in a token mode.", rule.Name()),
				rule,
				core.WithToken(rule.NameToken()),
				core.WithCode(ValidateTokenNotInTokenMode),
			))
		}
	}
}

// insideParserRule reports whether node contributes to what the parser matches.
// Keywords and rule calls nested in a token declaration, token group or token
// mode describe the lexer instead and are therefore not parser input.
func insideParserRule(node core.AstNode) bool {
	for container := node.Container(); container != nil; container = container.Container() {
		switch container.(type) {
		case ParserRule, CompositeRule:
			return true
		case TokenDecl, TokenGroup, TokenMode:
			return false
		}
	}
	return false
}

func collectTokenModeCoverage(g Grammar, ctx context.Context) tokenModeCoverage {
	coverage := tokenModeCoverage{
		keywords: collections.NewSet[string](),
		rules:    collections.NewSet[string](),
	}
	allKeywords := allGrammarKeywords(g)
	// Shared across modes: coverage is grammar-wide, so a group already expanded
	// for one mode needs no second pass, and the set breaks reference cycles.
	visitedGroups := collections.NewSet[string]()
	for _, mode := range g.TokenModes() {
		for _, member := range mode.Members() {
			switch member := member.(type) {
			case TokenDeclUsage:
				coverage.addTokenDecl(member.Declaration())
			case TokenGroupUsage:
				coverage.addTokenGroup(member.Group(), ctx, allKeywords, visitedGroups)
			case TokenUsage:
				coverage.addTokenRule(member.TokenRef().Ref(ctx), ctx, allKeywords, visitedGroups)
			case KeywordUsage:
				coverage.addKeyword(member.Keyword())
			case KeywordSelector:
				coverage.addKeywordsMatching(member.Selector(), allKeywords)
			}
		}
	}
	return coverage
}

// covers reports whether the lexer can produce a token for rule.
func (c tokenModeCoverage) covers(rule AbstractTokenRule) bool {
	if c.rules.Has(rule.Name()) {
		return true
	}
	// A token declared as a single keyword shares that keyword's token id, so
	// covering the keyword - through a 'keywords' selector, for instance -
	// covers the token as well.
	if decl, ok := rule.(TokenDecl); ok {
		if content, ok := decl.Content().(KeywordTokenContent); ok && content.Keyword() != nil {
			return c.keywords.Has(content.Keyword().Value())
		}
	}
	return false
}

func (c tokenModeCoverage) addTokenRule(rule AbstractTokenRule, ctx context.Context, allKeywords []Keyword, visited collections.Set[string]) {
	switch rule := rule.(type) {
	case TokenDecl:
		c.addTokenDecl(rule)
	case TokenGroup:
		c.addTokenGroup(rule, ctx, allKeywords, visited)
	}
}

func (c tokenModeCoverage) addTokenDecl(decl TokenDecl) {
	if decl == nil {
		return
	}
	c.rules.Add(decl.Name())
	// `token LBRACE: "{"` registers the keyword under the token's name, so
	// listing the token covers the keyword as well.
	if content, ok := decl.Content().(KeywordTokenContent); ok {
		c.addKeyword(content.Keyword())
	}
}

func (c tokenModeCoverage) addTokenGroup(group TokenGroup, ctx context.Context, allKeywords []Keyword, visited collections.Set[string]) {
	if group == nil || !visited.Add(group.Name()) {
		return
	}
	c.rules.Add(group.Name())
	for _, keyword := range group.Keywords() {
		c.addKeyword(keyword)
	}
	for _, selector := range group.KeywordSelectors() {
		c.addKeywordsMatching(selector.Image, allKeywords)
	}
	for _, tokenRef := range group.TokenRefs() {
		c.addTokenRule(tokenRef.Ref(ctx), ctx, allKeywords, visited)
	}
}

func (c tokenModeCoverage) addKeyword(keyword Keyword) {
	if keyword == nil || keyword.Value() == "" {
		return
	}
	c.keywords.Add(keyword.Value())
}

func (c tokenModeCoverage) addKeywordsMatching(selector string, allKeywords []Keyword) {
	if len(selector) < 2 {
		return
	}
	pattern, err := regexp.Compile(RegexpValue(selector))
	if err != nil {
		return
	}
	for _, keyword := range allKeywords {
		value, err := convertString(keyword)
		if err != nil {
			continue
		}
		if pattern.MatchString(value) {
			c.keywords.Add(keyword.Value())
		}
	}
}

// allGrammarKeywords returns every distinct keyword of the grammar, wherever it
// is declared. A `keywords` selector is matched against all of them.
func allGrammarKeywords(g Grammar) []Keyword {
	seen := collections.NewSet[string]()
	keywords := []Keyword{}
	for node := range core.AllChildren(g) {
		if keyword, ok := node.(Keyword); ok && seen.Add(keyword.Value()) {
			keywords = append(keywords, keyword)
		}
	}
	return keywords
}

func checkUniqueRuleNames(g Grammar, accept core.ValidationAcceptor) {
	seen := map[string][]core.NamedTokenNode{}
	for _, rule := range g.Rules() {
		if rule.Name() != "" {
			seen[rule.Name()] = append(seen[rule.Name()], rule)
		}
	}
	for _, terminal := range g.Terminals() {
		if terminal.Name() != "" {
			seen[terminal.Name()] = append(seen[terminal.Name()], terminal)
		}
	}
	for _, tokenGroup := range g.TokenGroups() {
		if tokenGroup.Name() != "" {
			seen[tokenGroup.Name()] = append(seen[tokenGroup.Name()], tokenGroup)
		}
	}
	for _, tokenMode := range g.TokenModes() {
		for _, member := range tokenMode.Members() {
			if usage, ok := member.(TokenDeclUsage); ok {
				decl := usage.Declaration()
				if decl.Name() != "" {
					seen[decl.Name()] = append(seen[decl.Name()], decl)
				}
			} else if group, ok := member.(TokenGroupUsage); ok {
				decl := group.Group()
				if decl.Name() != "" {
					seen[decl.Name()] = append(seen[decl.Name()], decl)
				}
			}
		}
	}
	for name, nodes := range seen {
		if len(nodes) > 1 {
			for _, node := range nodes {
				accept(core.NewDiagnostic(
					core.SeverityError,
					fmt.Sprintf("A rule's name has to be unique. '%s' is used multiple times.", name),
					node,
					core.WithToken(node.NameToken()),
					core.WithCode(ValidateUniqueRuleName),
				))
			}
		}
	}
}

func checkUniqueInterfaceNames(g Grammar, accept core.ValidationAcceptor) {
	seen := map[string][]Interface{}
	for _, iface := range g.Interfaces() {
		if iface.Name() != "" {
			seen[iface.Name()] = append(seen[iface.Name()], iface)
		}
	}
	for name, ifaces := range seen {
		if len(ifaces) > 1 {
			for _, iface := range ifaces {
				accept(core.NewDiagnostic(
					core.SeverityError,
					fmt.Sprintf("An interface name has to be unique. '%s' is used multiple times.", name),
					iface,
					core.WithToken(iface.NameToken()),
					core.WithCode(ValidateUniqueInterfaceName),
				))
			}
		}
	}
}

// TokenImpl.Validate checks terminal rule constraints:
//   - The regular expression should not match the empty string.
func (t *TokenDeclImpl) Validate(_ context.Context, _ string, accept core.ValidationAcceptor) {
	checkEmptyTerminalRule(t, accept)
}

func checkEmptyTerminalRule(t TokenDecl, accept core.ValidationAcceptor) {
	var canBeEmpty bool
	switch content := t.Content().(type) {
	case KeywordTokenContent:
		raw := KeywordValue(content.Keyword())
		canBeEmpty = raw == ""
	case RegexpTokenContent:
		pattern := RegexpValue(content.Regexp())
		re, err := regexp.Compile(pattern)
		if err != nil {
			return
		}
		canBeEmpty = re.MatchString("")
	}
	if canBeEmpty {
		accept(core.NewDiagnostic(
			core.SeverityError,
			"This terminal could match an empty string.",
			t,
			core.WithToken(t.NameToken()),
			core.WithCode(ValidateEmptyToken),
		))
	}
}

// KeywordImpl.Validate checks keyword constraints:
//   - Keywords cannot be empty.
//   - Keywords cannot consist only of whitespace.
//   - Keywords should not contain whitespace characters (warning).
func (k *KeywordImpl) Validate(_ context.Context, _ string, accept core.ValidationAcceptor) {
	checkKeyword(k, accept)
}

func checkKeyword(k Keyword, accept core.ValidationAcceptor) {
	value, err := convertString(k)
	if err != nil {
		return
	}
	if value == "" {
		accept(core.NewDiagnostic(
			core.SeverityError,
			"Keywords cannot be empty.",
			k,
			core.WithToken(k.ValueToken()),
			core.WithCode(ValidateEmptyKeyword),
		))
	} else if strings.TrimSpace(value) == "" {
		accept(core.NewDiagnostic(
			core.SeverityError,
			"Keywords cannot only consist of whitespace characters.",
			k,
			core.WithToken(k.ValueToken()),
			core.WithCode(ValidateWhitespaceOnlyKeyword),
		))
	} else if strings.ContainsAny(value, " \t\n\r") {
		accept(core.NewDiagnostic(
			core.SeverityWarning,
			"Keywords should not contain whitespace characters.",
			k,
			core.WithToken(k.ValueToken()),
			core.WithCode(ValidateKeywordWithWhitespace),
		))
	}
}

func (rule *ParserRuleImpl) Validate(ctx context.Context, _ string, accept core.ValidationAcceptor) {
	checkRuleReturnType(rule, ctx, accept)
}

func checkRuleReturnType(rule ParserRule, _ context.Context, accept core.ValidationAcceptor) {
	// Only search if not explicitly provided
	if rule.ReturnType() == nil && rule.Name() != "" {
		grammar, ok := rule.Container().(Grammar)
		if !ok || grammar == nil {
			return
		}
		returnType := FindInterfaceByName(grammar, rule.Name())
		if returnType == nil {
			accept(
				core.NewDiagnostic(
					core.SeverityError,
					fmt.Sprintf("Unable to find return type for rule '%s'. Either define an interface with the same name as the rule or explicitly specify the return type.", rule.Name()),
					rule,
					core.WithToken(rule.NameToken()),
					core.WithCode(ValidateRuleReturnType),
				),
			)
		}
	}
}

func (i *InterfaceImpl) Validate(ctx context.Context, _ string, accept core.ValidationAcceptor) {
	checkInterfaceExtends(i, ctx, accept)
	checkInterfaceFieldNames(i, ctx, accept)
	checkInterfaceFieldTypes(i, accept)
}

func collectInheritedFieldNames(iface Interface, ctx context.Context, collected map[string]Interface, visited collections.Set[string]) {
	if !visited.Add(iface.Name()) {
		return
	}
	for _, ext := range iface.Extends() {
		extType := ext.Ref(ctx)
		if extType == nil {
			continue
		}
		// Recurse into the parent's ancestry first so the deepest (originating)
		// declarant wins in `collected` when the same name appears at multiple levels.
		collectInheritedFieldNames(extType, ctx, collected, visited)
		for _, field := range extType.Fields() {
			name := field.Name()
			if name == "" {
				continue
			}
			lower := strings.ToLower(name)
			if _, exists := collected[lower]; !exists {
				collected[lower] = extType
			}
		}
	}
}

func checkInterfaceFieldNames(iface Interface, ctx context.Context, accept core.ValidationAcceptor) {
	inherited := map[string]Interface{}
	collectInheritedFieldNames(iface, ctx, inherited, collections.NewSet[string]())

	allFields := map[string]Interface{}
	for lower, declaringIface := range inherited {
		allFields[lower] = declaringIface
	}
	for _, field := range iface.Fields() {
		name := field.Name()
		if name == "" {
			continue
		}
		lower := strings.ToLower(name)
		if _, exists := allFields[lower]; !exists {
			allFields[lower] = iface
		}
	}

	seen := collections.NewSet[string]()
	for _, field := range iface.Fields() {
		name := field.Name()
		if name == "" {
			continue
		}
		if !unicode.IsUpper(rune(name[0])) {
			accept(core.NewDiagnostic(
				core.SeverityError,
				"Field names must start with a capital letter.",
				field,
				core.WithToken(field.NameToken()),
				core.WithCode(ValidateFieldNameCapitalLetter),
			))
		}
		checkReservedFieldName(field, iface, allFields, accept)
		lower := strings.ToLower(name)
		if seen.Has(lower) {
			accept(core.NewDiagnostic(
				core.SeverityError,
				fmt.Sprintf("A field's name has to be unique (case-insensitively). '%s' is already used above.", name),
				field,
				core.WithToken(field.NameToken()),
				core.WithCode(ValidateUniqueFieldName),
			))
		} else {
			seen.Add(lower)
			if declaringIface, dup := inherited[lower]; dup {
				accept(core.NewDiagnostic(
					core.SeverityError,
					fmt.Sprintf("A field's name has to be unique (case-insensitively). '%s' is already declared in '%s'.", name, declaringIface.Name()),
					field,
					core.WithToken(field.NameToken()),
					core.WithCode(ValidateUniqueFieldName),
				))
			}
		}
	}
}

func checkInterfaceFieldTypes(iface Interface, accept core.ValidationAcceptor) {
	for _, field := range iface.Fields() {
		fieldType := field.Type()
		if isNestedArrayType(fieldType) {
			accept(core.NewDiagnostic(
				core.SeverityError,
				"Nested array types are not supported.",
				fieldType,
				core.WithCode(ValidateNestedArrayType),
			))
		}
	}
}

func isNestedArrayType(fieldType FieldType) bool {
	arrayType, ok := fieldType.(ArrayType)
	if !ok {
		return false
	}
	_, ok = arrayType.InternalType().(ArrayType)
	return ok
}

func checkReservedFieldName(field Field, iface Interface, allFields map[string]Interface, accept core.ValidationAcceptor) {
	name := field.Name()
	if name == "" {
		return
	}

	if strings.HasPrefix(name, "Set") {
		accept(core.NewDiagnostic(
			core.SeverityError,
			"Field names must not start with 'Set' because the framework generates Set{Name}() setter methods.",
			field,
			core.WithToken(field.NameToken()),
			core.WithCode(ValidateReservedFieldName),
		))
		return
	}
	if strings.HasPrefix(name, "Is") {
		accept(core.NewDiagnostic(
			core.SeverityError,
			"Field names must not start with 'Is' because the framework generates Is{Name}() methods for boolean fields.",
			field,
			core.WithToken(field.NameToken()),
			core.WithCode(ValidateReservedFieldName),
		))
		return
	}
	if name == iface.Name() {
		accept(core.NewDiagnostic(
			core.SeverityError,
			"A field's name cannot be the same as the interface name due to potential conflicts with generated methods.",
			field,
			core.WithToken(field.NameToken()),
			core.WithCode(ValidateReservedFieldName),
		))
		return
	}
	if reserved, ok := reservedFieldNames[name]; ok {
		accept(core.NewDiagnostic(
			core.SeverityError,
			fmt.Sprintf("The field name '%s' is reserved by the framework and cannot be used because it would conflict with %s.", name, reserved),
			field,
			core.WithToken(field.NameToken()),
			core.WithCode(ValidateReservedFieldName),
		))
		return
	}
	if base, _, ok := tokenOrNodeSuffixBase(name); ok {
		lower := strings.ToLower(base)
		if _, exists := allFields[lower]; exists && !strings.EqualFold(base, name) {
			accept(core.NewDiagnostic(
				core.SeverityError,
				fmt.Sprintf("The field name '%s' conflicts with '%s' due to potential conflicts with generated methods.", name, base),
				field,
				core.WithToken(field.NameToken()),
				core.WithCode(ValidateReservedFieldName),
			))
		}
	}
}

func tokenOrNodeSuffixBase(name string) (base, suffix string, ok bool) {
	for _, suffix := range []string{"Token", "Node"} {
		if strings.HasSuffix(name, suffix) && len(name) > len(suffix) {
			return name[:len(name)-len(suffix)], suffix, true
		}
	}
	return "", "", false
}

func checkInterfaceExtends(iface Interface, ctx context.Context, accept core.ValidationAcceptor) {
	for _, ext := range iface.Extends() {
		extType := ext.Ref(ctx)
		if extType == nil {
			continue
		}
		if appearsInExtends(iface, extType, ctx, collections.NewSet[string]()) {
			accept(core.NewDiagnostic(
				core.SeverityError,
				"An interface cannot extend itself, neither directly nor indirectly.",
				iface,
				core.WithReference(ext),
				core.WithCode(ValidateInterfaceExtends),
			))
		}
	}
}

func appearsInExtends(target Interface, current Interface, ctx context.Context, visited collections.Set[string]) bool {
	if current.Name() == target.Name() {
		return true
	}
	if !visited.Add(current.Name()) {
		return false
	}
	for _, ext := range current.Extends() {
		extType := ext.Ref(ctx)
		if extType == nil {
			continue
		}
		if appearsInExtends(target, extType, ctx, visited) {
			return true
		}
	}
	return false
}

func (r *RuleCallImpl) Validate(ctx context.Context, _ string, accept core.ValidationAcceptor) {
	assignment := core.ContainerOfType[Assignment](r)
	if assignment == nil {
		// Some validations only apply to unassigned rule calls
		checkRuleCallReturnType(r, ctx, accept)
		checkRuleCallPosition(r, ctx, accept)
	}
}

func checkRuleCallReturnType(call RuleCall, ctx context.Context, accept core.ValidationAcceptor) {
	ownRule := core.ContainerOfType[ParserRule](call)
	ownType := FindReturnType(ownRule, ctx)
	if ownType == nil {
		return
	}
	// Unassigned rule call
	targetRule := call.Rule().Ref(ctx)
	if targetRule == nil {
		return
	}
	if parserRule, ok := targetRule.(ParserRule); ok {
		targetType := FindReturnType(parserRule, ctx)
		if targetType == nil {
			return
		}
		if !interfaceIsAssignableTo(targetType, ownType) {
			accept(core.NewDiagnostic(
				core.SeverityError,
				fmt.Sprintf("The return type '%s' of the called rule is not assignable to the return type '%s' of the current rule.", targetType.Name(), ownType.Name()),
				call,
				core.WithCode(ValidateRuleCallReturnType),
			))
		}
	}
}

func checkRuleCallPosition(call RuleCall, ctx context.Context, accept core.ValidationAcceptor) {
	rule := call.Rule().Ref(ctx)
	if _, ok := rule.(ParserRule); !ok {
		// Only parser rules can cause information loss, so we only check those
		return
	}
	// An unassigned rule call cannot be preceded by an action or assignment
	// This would lead to information loss, as the result of the rule call overrides the current AST node
	var node core.AstNode = call
	for node != nil {
		container := node.Container()
		if _, ok := container.(ParserRule); ok {
			break
		}
		if group, ok := container.(Group); ok {
			for _, elem := range group.Elements() {
				if elem == node {
					break
				}
				if action, ok := elem.(Action); ok && action.Property() != nil {
					accept(core.NewDiagnostic(
						core.SeverityError,
						"An unassigned rule call cannot be preceded by an assigned action.",
						call,
						core.WithCode(ValidateRuleCallPosition),
					))
					return
				}
				if _, ok := elem.(Assignment); ok {
					accept(core.NewDiagnostic(
						core.SeverityError,
						"An unassigned rule call cannot be preceded by an assignment.",
						call,
						core.WithCode(ValidateRuleCallPosition),
					))
					return
				}
			}
		}
		node = container
	}
}

func (a *ActionImpl) Validate(ctx context.Context, _ string, accept core.ValidationAcceptor) {
	checkActionAssignmentType(a, ctx, accept)
	checkActionPropertyType(a, ctx, accept)
}

func checkActionAssignmentType(a Action, ctx context.Context, accept core.ValidationAcceptor) {
	targetType := a.Type().Ref(ctx)
	if targetType == nil {
		return
	}
	rule := core.ContainerOfType[ParserRule](a)
	if rule == nil {
		return
	}
	returnType := FindReturnType(rule, ctx)
	if returnType == nil {
		return
	}
	if !interfaceIsAssignableTo(targetType, returnType) {
		accept(core.NewDiagnostic(
			core.SeverityError,
			fmt.Sprintf("The type '%s' of the action is not assignable to the rule's return type '%s'.", targetType.Name(), returnType.Name()),
			a,
			core.WithReference(a.Type()),
			core.WithCode(ValidateActionAssignmentType),
		))
	}
}

func checkActionPropertyType(a Action, ctx context.Context, accept core.ValidationAcceptor) {
	targetField := a.Property().Ref(ctx)
	if targetField == nil {
		return
	}
	currentType := getCurrentType(ctx, a)
	if currentType == nil {
		return
	}
	targetType := targetField.Type()
	if a.Operator() == "+=" {
		if arrayType, ok := targetType.(ArrayType); ok {
			// Reassign target type to the array internal type
			targetType = arrayType.InternalType()
		} else {
			accept(core.NewDiagnostic(
				core.SeverityError,
				"The '+=' operator can only be used on array fields.",
				a,
				core.WithToken(a.OperatorToken()),
				core.WithCode(ValidateActionPropertyType),
			))
			return
		}
	}

	if simpleType, ok := targetType.(SimpleType); ok {
		targetInterface := simpleType.Type().Ref(ctx)
		if targetInterface == nil {
			return
		}
		if !interfaceIsAssignableTo(currentType, targetInterface) {
			accept(core.NewDiagnostic(
				core.SeverityError,
				fmt.Sprintf("The local type '%s' is not assignable to the target field type '%s'.", currentType.Name(), targetInterface.Name()),
				a,
				core.WithReference(a.Property()),
				core.WithCode(ValidateActionPropertyType),
			))
		}
	} else {
		accept(core.NewDiagnostic(
			core.SeverityError,
			"Cannot assign a parser rule to a non-interface field.",
			a,
			core.WithReference(a.Property()),
			core.WithCode(ValidateActionPropertyType),
		))
	}
}

func (a *AssignmentImpl) Validate(ctx context.Context, _ string, accept core.ValidationAcceptor) {
	checkAssignmentType(a, ctx, accept)
}

func checkAssignmentType(a Assignment, ctx context.Context, accept core.ValidationAcceptor) {
	propRef := a.Property()
	if propRef == nil {
		return
	}
	field := propRef.Ref(ctx)
	if field == nil {
		return
	}
	fieldType := field.Type()
	if fieldType == nil {
		return
	}

	var effectiveFieldType FieldType

	switch a.Operator() {
	case "?=":
		pt, ok := fieldType.(PrimitiveType)
		if !ok || pt.Type() != "bool" {
			accept(core.NewDiagnostic(
				core.SeverityError,
				"The '?=' operator can only be used on boolean fields.",
				a,
				core.WithToken(a.OperatorToken()),
				core.WithCode(ValidateAssignmentType),
			))
		}
		return
	case "+=":
		at, ok := fieldType.(ArrayType)
		if !ok {
			accept(core.NewDiagnostic(
				core.SeverityError,
				"The '+=' operator can only be used on array fields.",
				a,
				core.WithToken(a.OperatorToken()),
				core.WithCode(ValidateAssignmentType),
			))
			return
		}
		effectiveFieldType = at.InternalType()
	default:
		effectiveFieldType = fieldType
	}

	value := a.Value()
	if value != nil {
		isAssignableTo(ctx, value, effectiveFieldType, accept)
	}
}

func isAssignableTo(ctx context.Context, source Assignable, fieldType FieldType, accept core.ValidationAcceptor) {
	switch v := source.(type) {
	case CrossRef:
		if refType, ok := fieldType.(ReferenceType); ok {
			toType := refType.Type().Ref(ctx)
			if toType == nil {
				return
			}
			fromType := v.Type().Ref(ctx)
			if fromType == nil {
				return
			}
			if !interfaceIsAssignableTo(fromType, toType) {
				accept(core.NewDiagnostic(
					core.SeverityError,
					fmt.Sprintf("The type '%s' of the cross-reference value is not assignable to the target field type '%s'.", fromType.Name(), toType.Name()),
					v,
					core.WithReference(v.Type()),
					core.WithCode(ValidateAssignmentType),
				))
			}
		} else {
			accept(core.NewDiagnostic(
				core.SeverityError,
				"Cannot assign a cross-reference value to a non-reference field.",
				v,
				core.WithCode(ValidateAssignmentType),
			))
		}
	case RuleCall:
		resolvedRule := v.Rule().Ref(ctx)
		if resolvedRule == nil {
			return
		}
		switch rule := resolvedRule.(type) {
		case TokenDecl:
			if primitiveType, ok := fieldType.(PrimitiveType); !ok || primitiveType.Type() != "string" {
				accept(core.NewDiagnostic(
					core.SeverityError,
					"Cannot assign a token to a non-string field.",
					v,
					core.WithCode(ValidateAssignmentType),
				))
			}
		case ParserRule:
			if simpleType, ok := fieldType.(SimpleType); ok {
				ruleType := FindReturnType(rule, ctx)
				if ruleType == nil {
					return
				}
				targetType := simpleType.Type().Ref(ctx)
				if targetType == nil {
					return
				}
				if !interfaceIsAssignableTo(ruleType, targetType) {
					accept(core.NewDiagnostic(
						core.SeverityError,
						fmt.Sprintf("The return type '%s' of the called rule is not assignable to the target field type '%s'.", ruleType.Name(), targetType.Name()),
						v,
						core.WithCode(ValidateAssignmentType),
					))
				}
			} else {
				accept(core.NewDiagnostic(
					core.SeverityError,
					"Cannot assign a parser rule to a non-interface field.",
					v,
					core.WithCode(ValidateAssignmentType),
				))
			}
		case CompositeRule:
			if primitiveType, ok := fieldType.(PrimitiveType); !ok || primitiveType.Type() != "composite" {
				if primitiveType, ok := fieldType.(PrimitiveType); ok && primitiveType.Type() == "string" {
					accept(core.NewDiagnostic(
						core.SeverityError,
						"Cannot assign a composite rule to a string field. Use 'composite' as the field type instead.",
						v,
						core.WithCode(ValidateAssignmentType),
					))
				} else {
					accept(core.NewDiagnostic(
						core.SeverityError,
						"Cannot assign a composite rule to a non-composite field.",
						v,
						core.WithCode(ValidateAssignmentType),
					))
				}
			}
		}
	case Keyword:
		if primitiveType, ok := fieldType.(PrimitiveType); !ok || primitiveType.Type() != "string" {
			accept(core.NewDiagnostic(
				core.SeverityError,
				"Cannot assign a keyword value to a non-string field.",
				v,
				core.WithCode(ValidateAssignmentType),
			))
		}
	case Alternatives:
		for _, option := range v.Alts() {
			if assignableOption, ok := option.(Assignable); ok {
				isAssignableTo(ctx, assignableOption, fieldType, accept)
			}
		}
	}
}

func interfaceIsAssignableTo(source Interface, target Interface) bool {
	return doInterfaceIsAssignableTo(source, target, collections.NewSet[string]())
}

func doInterfaceIsAssignableTo(source Interface, target Interface, visited collections.Set[string]) bool {
	if source.Name() == target.Name() {
		return true
	}
	if !visited.Add(source.Name()) {
		return false
	}
	for _, ext := range source.Extends() {
		extType := ext.Ref(context.Background())
		if extType == nil {
			continue
		}
		if doInterfaceIsAssignableTo(extType, target, visited) {
			return true
		}
	}
	return false
}

func (tg *TokenGroupImpl) Validate(_ context.Context, _ string, accept core.ValidationAcceptor) {
	checkRecursiveTokenGroup(tg, accept)
	checkTokenGroupContainsOnlyValidTokens(tg, accept)
	for _, selector := range tg.KeywordSelectors() {
		checkRegExpIsValid(selector, accept)
	}
}

func checkRecursiveTokenGroup(tg TokenGroup, accept core.ValidationAcceptor) {
	if appearsInTokenGroup(tg, tg, context.Background(), collections.NewSet[string]()) {
		accept(core.NewDiagnostic(
			core.SeverityError,
			"A token group cannot contain itself, neither directly nor indirectly.",
			tg,
			core.WithToken(tg.NameToken()),
			core.WithCode(ValidateRecursiveTokenGroup),
		))
	}
}

func checkTokenGroupContainsOnlyValidTokens(tg TokenGroup, accept core.ValidationAcceptor) {
	for _, tokenRef := range tg.TokenRefs() {
		token := tokenRef.Ref(context.Background())
		if token == nil {
			continue
		}
		if description, special := hiddenOrCommentTokenDescription(token); special {
			accept(core.NewDiagnostic(
				core.SeverityError,
				fmt.Sprintf("The token '%s' cannot be used in a token group because it is %s.", token.Name(), description),
				tg,
				core.WithReference(tokenRef),
				core.WithCode(ValidateInvalidTokenInGroup),
			))
		}
	}
}

func appearsInTokenGroup(target TokenGroup, current TokenGroup, ctx context.Context, visited collections.Set[string]) bool {
	if !visited.Add(current.Name()) {
		return false
	}
	for _, ext := range current.TokenRefs() {
		token := ext.Ref(ctx)
		if token == nil {
			continue
		}
		if tokenGroup, ok := token.(TokenGroup); ok &&
			(target.Name() == tokenGroup.Name() || appearsInTokenGroup(target, tokenGroup, ctx, visited)) {
			return true
		}
	}
	return false
}

func hiddenOrCommentTokenDescription(tokenDecl AbstractTokenRule) (description string, ok bool) {
	switch tokenDecl.Modifier() {
	case "hidden":
		return "hidden", true
	case "comment":
		return "a comment", true
	default:
		return "", false
	}
}

func (cr *CrossRefImpl) Validate(ctx context.Context, _ string, accept core.ValidationAcceptor) {
	checkCrossRefHasTerminal(cr, ctx, accept)
	checkCrossRefToken(cr, ctx, accept)
}

func checkCrossRefHasTerminal(cr CrossRef, ctx context.Context, accept core.ValidationAcceptor) {
	if cr.Rule() != nil {
		return
	}
	resolvedType := cr.Type().Ref(ctx)
	if resolvedType == nil {
		return
	}
	accept(core.NewDiagnostic(
		core.SeverityError,
		fmt.Sprintf("A cross-reference must specify a token or composite rule after ':' (e.g. [%s:ID]).", resolvedType.Name()),
		cr,
		core.WithReference(cr.Type()),
		core.WithCode(ValidateMissingCrossRefTerminal),
	))
}

func checkCrossRefToken(cr CrossRef, ctx context.Context, accept core.ValidationAcceptor) {
	ruleCall := cr.Rule()
	if ruleCall == nil {
		return
	}
	resolved := ruleCall.Rule().Ref(ctx)
	if resolved == nil {
		return
	}
	tokenDecl, ok := resolved.(TokenDecl)
	if !ok {
		return
	}
	// Hidden/comment tokens are not allowed in cross-references because they
	// are not stored in the token slice and cannot identify named elements.
	if description, special := hiddenOrCommentTokenDescription(tokenDecl); special {
		accept(core.NewDiagnostic(
			core.SeverityError,
			fmt.Sprintf("The token '%s' cannot be used in a cross-reference because it is %s.", tokenDecl.Name(), description),
			cr,
			core.WithReference(ruleCall.Rule()),
			core.WithCode(ValidateInvalidTokenInCrossRef),
		))
	}
}

func checkTokenModeMembersAreUnique(tm TokenMode, accept core.ValidationAcceptor) {
	seen := collections.NewSet[string]()
	for _, member := range tm.Members() {
		var name string
		var textRange core.TextRange
		var token *core.Token = nil
		switch member := member.(type) {
		case TokenDeclUsage:
			continue // TokenDeclUsage is already checked in checkUniqueRuleNames
		case TokenGroupUsage:
			continue // TokenGroupUsage is already checked in checkUniqueRuleNames
		case TokenUsage:
			decl := member.TokenRef().Ref(context.Background())
			textRange = member.TextRange()
			if decl == nil {
				continue
			}
			name = decl.Name()
		case KeywordUsage:
			name = member.Keyword().Value()
			token = member.Keyword().ValueToken()
		default:
			continue
		}
		if token != nil {
			textRange = token.TextRange()
		}
		if !seen.Add(name) {
			accept(core.NewDiagnostic(
				core.SeverityError,
				fmt.Sprintf("The token mode '%s' contains multiple members with the same name '%s'.", tm.Name(), name),
				member,
				core.WithTextRange(textRange),
				core.WithCode(ValidateUniqueRuleNameInTokenMode),
			))
		}
	}
}

func (m *KeywordSelectorImpl) Validate(_ context.Context, _ string, accept core.ValidationAcceptor) {
	checkRegExpIsValid(m.SelectorToken(), accept)
}

func (m *RegexpTokenContentImpl) Validate(_ context.Context, _ string, accept core.ValidationAcceptor) {
	checkRegExpIsValid(m.RegexpToken(), accept)
}

func checkRegExpIsValid(patternToken *core.Token, accept core.ValidationAcceptor) {
	_, err := regexp.Compile(RegexpValue(patternToken.Image))
	if err != nil {
		accept(core.NewDiagnostic(
			core.SeverityError,
			fmt.Sprintf("The keyword selector '%s' is not a valid regular expression: %s", patternToken.Image, err.Error()),
			patternToken.Element,
			core.WithToken(patternToken),
			core.WithCode(ValidateInvalidRegExpLiteral),
		))
	}
}
