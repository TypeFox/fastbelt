// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"slices"
	"strconv"
	"strings"

	"github.com/TypeFox/langium-to-go/core"
	"github.com/TypeFox/langium-to-go/generator"
	"github.com/TypeFox/langium-to-go/internal/generated"
)

type ParserGeneratorContext struct {
	grammar            generated.Grammar
	accessNamesCounter map[string]int
	accessNames        map[core.AstNode]string
	lookaheads         map[core.AstNode]LookaheadValue
	orLookaheads       map[core.AstNode]LookaheadValue
}

type LookaheadValue struct {
	name string
	llk  LLkLookahead
}

func (context *ParserGeneratorContext) SetAccessName(node core.AstNode, name string) {
	index := context.accessNamesCounter[name]
	context.accessNames[node] = name + "_" + strconv.Itoa(index)
	index++
	context.accessNamesCounter[name] = index
}

func GenerateParser(grammar generated.Grammar) string {
	context := &ParserGeneratorContext{
		grammar:            grammar,
		accessNamesCounter: make(map[string]int),
		accessNames:        make(map[core.AstNode]string),
		lookaheads:         make(map[core.AstNode]LookaheadValue),
		orLookaheads:       make(map[core.AstNode]LookaheadValue),
	}
	populateContext(context)
	node := generator.NewNode()
	node.AppendLine("package generated")
	node.AppendLine()
	node.AppendLine("import (")
	node.Indent(func(n generator.Node) {
		n.AppendLine("\"github.com/TypeFox/langium-to-go/core\"")
		n.AppendLine("\"github.com/TypeFox/langium-to-go/parser\"")
	})
	node.AppendLine(")")
	node.AppendLine()

	node.AppendLine("type Parser struct {")
	node.Indent(func(n generator.Node) {
		n.AppendLine("state *parser.ParserState")
	})
	node.AppendLine("}")
	node.AppendLine()
	firstRule := grammar.Rules()[0]
	node.AppendLine("func (p *Parser) Parse(tokens []*core.Token) ", firstRule.Name(), " {")
	node.Indent(func(n generator.Node) {
		n.AppendLine("p.state = parser.NewParserState(tokens)")
		n.AppendLine("result := p.Parse", firstRule.Name(), "()")
		n.AppendLine("core.AssignContainers(result)")
		n.AppendLine("return result")
	})
	node.AppendLine("}")
	node.AppendLine()
	node.AppendLine("func NewParser() *Parser {")
	node.Indent(func(n generator.Node) {
		n.AppendLine("return &Parser{}")
	})
	node.AppendLine("}")
	node.AppendLine()

	node.AppendLine("const (")
	accessIota := true
	accessNames := make([]string, 0, len(context.accessNames))
	node.Indent(func(access generator.Node) {
		for _, name := range context.accessNames {
			accessNames = append(accessNames, name)
		}
		slices.Sort(accessNames)
		for _, name := range accessNames {
			if accessIota {
				access.AppendLine(name + " = iota + 1")
				accessIota = false
			} else {
				access.AppendLine(name)
			}
		}
	})
	node.AppendLine(")")
	node.AppendLine()

	lookaheads := make([]LookaheadValue, 0, len(context.lookaheads))

	for _, lookahead := range context.lookaheads {
		lookaheads = append(lookaheads, lookahead)
	}
	for _, lookahead := range context.orLookaheads {
		lookaheads = append(lookaheads, lookahead)
	}
	slices.SortFunc(lookaheads, func(a, b LookaheadValue) int {
		return strings.Compare(a.name, b.name)
	})
	for _, lookahead := range lookaheads {
		generateLLkLookahead(node, lookahead.name, lookahead.llk)
		node.AppendLine()
	}

	for _, rule := range grammar.Rules() {
		generateParseFunction(node, context, rule)
	}

	return formatIfPossible(node.String())
}

func populateContext(context *ParserGeneratorContext) {
	for _, rule := range context.grammar.Rules() {
		ruleName := rule.Name()
		populateContextWithNode(context, ruleName, rule.Body())
	}
}

func populateContextWithNode(context *ParserGeneratorContext, prefix string, node core.AstNode) {
	if _, exists := context.accessNames[node]; exists {
		return
	}
	switch n := node.(type) {
	case generated.Alternatives:
		if len(n.Alts()) > 1 {
			name := prefix + "LookaheadOr" + strconv.Itoa(len(context.orLookaheads))
			context.orLookaheads[n] = LookaheadValue{name: name, llk: GetLLkLookaheadOr(context.grammar, n)}
		}
		if n.Cardinality() != CardinalityOne {
			name := prefix + "Lookahead" + strconv.Itoa(len(context.lookaheads))
			context.lookaheads[n] = LookaheadValue{name: name, llk: GetLLkLookaheadOpt(context.grammar, n)}
		}
		for _, alt := range n.Alts() {
			populateContextWithNode(context, prefix, alt)
		}
	case generated.Group:
		if n.Cardinality() != CardinalityOne {
			name := prefix + "Lookahead" + strconv.Itoa(len(context.lookaheads))
			context.lookaheads[n] = LookaheadValue{name: name, llk: GetLLkLookaheadOpt(context.grammar, n)}
		}
		for _, element := range n.Elements() {
			populateContextWithNode(context, prefix, element)
		}
	case generated.Keyword:
		if n.Cardinality() != CardinalityOne {
			name := prefix + "Lookahead" + strconv.Itoa(len(context.lookaheads))
			context.lookaheads[n] = LookaheadValue{name: name, llk: GetLLkLookaheadOpt(context.grammar, n)}
		}
		name := prefix + KeywordName(n)
		context.SetAccessName(node, name)
	case generated.Assignment:
		if n.Cardinality() != CardinalityOne {
			name := prefix + "Lookahead" + strconv.Itoa(len(context.lookaheads))
			context.lookaheads[n] = LookaheadValue{name: name, llk: GetLLkLookaheadOpt(context.grammar, n)}
		}
		name := prefix + n.Property()
		populateContextWithNode(context, name, n.Value())
	case generated.CrossRef:
		populateContextWithNode(context, prefix, n.Rule())
	case generated.RuleCall:
		if n.Cardinality() != CardinalityOne {
			name := prefix + "Lookahead" + strconv.Itoa(len(context.lookaheads))
			context.lookaheads[n] = LookaheadValue{name: name, llk: GetLLkLookaheadOpt(context.grammar, n)}
		}
		token := getTokenWithName(context.grammar, n.Rule())
		if token != nil {
			name := prefix + token.Name()
			context.SetAccessName(node, name)
		}
	}
}

func generateParseFunction(node generator.Node, context *ParserGeneratorContext, rule generated.ParserRule) {
	node.AppendLine("func (p *Parser) Parse", rule.Name(), "() ", rule.ReturnType(), " {")
	node.Indent(func(n generator.Node) {
		n.AppendLine("node := New", rule.ReturnType(), "()")
		n.AppendLine("node.WithSegmentStartToken(p.state.LA(1))")
		generateAbstractElementParser(n, context, rule.Body())
		n.AppendLine("node.WithSegmentEndToken(p.state.LA(0))")
		n.AppendLine("return node")
	})
	node.AppendLine("}")
	node.AppendLine()
}

func generateAbstractElementParser(node generator.Node, context *ParserGeneratorContext, element generated.Element) {
	if alts, ok := element.(generated.Alternatives); ok {
		generateAlternativesParser(node, context, alts)
	} else if group, ok := element.(generated.Group); ok {
		generateGroupParser(node, context, group)
	} else if action, ok := element.(generated.Action); ok {
		node.AppendLine("{")
		node.Indent(func(n generator.Node) {
			n.AppendLine("result := New", action.Type(), "()")
			if action.Property() != "" {
				if action.Operator() == "+=" {
					n.AppendLine("result.With", action.Property(), "Item(node)")
				} else {
					n.AppendLine("result.With", action.Property(), "(node)")
				}
				n.AppendLine("node = result")
			} else {
				n.AppendLine("core.AssignTokens(result, node.Tokens())")
				n.AppendLine("node = result")
			}
		})
		node.AppendLine("}")
		node.AppendLine("node := node.(", action.Type(), ")")
	} else {
		node.AppendLine("{")
		node.Indent(func(indent generator.Node) {
			if keyword, ok := element.(generated.Keyword); ok {
				generateKeywordParser(indent, context, keyword)
			} else if ruleCall, ok := element.(generated.RuleCall); ok {
				resultName := generateRuleCallParser(indent, context, ruleCall)
				if resultName == "result" {
					indent.AppendLine("core.MergeTokens(result, node.Tokens())")
					indent.AppendLine("node = result")
				}
			} else if assignment, ok := element.(generated.Assignment); ok {
				generateCardinality(indent, func(n generator.Node) {
					generateAssignable(n, context, assignment.Value(), func(n2 generator.Node, resultName string) {
						n2.AppendLine("if ", resultName, " != nil {")
						n2.Indent(func(in generator.Node) {
							switch assignment.Operator() {
							case "+=":
								// Append to slice
								in.AppendLine("node.With", assignment.Property(), "Item(", resultName, ")")
							default:
								// Single assignment
								in.AppendLine("node.With", assignment.Property(), "(", resultName, ")")
							}
						})
						n2.AppendLine("}")
					})
				}, func(n generator.Node) {
					lookaheadName := context.lookaheads[element].name
					lookahead := generateLookaheadString(lookaheadName)
					n.Append(lookahead)
				}, element.Cardinality())
			}
		})
		node.AppendLine("}")
	}
}

func generateAssignable(node generator.Node, context *ParserGeneratorContext, assignable generated.Assignable, cb func(node generator.Node, resultName string)) {
	if crossRef, ok := assignable.(generated.CrossRef); ok {
		resultName := generateCrossReferenceParser(node, context, crossRef)
		cb(node, resultName)
	} else if keyword, ok := assignable.(generated.Keyword); ok {
		resultName := generateKeywordParser(node, context, keyword)
		cb(node, resultName)
	} else if ruleCall, ok := assignable.(generated.RuleCall); ok {
		resultName := generateRuleCallParser(node, context, ruleCall)
		cb(node, resultName)
	} else if alts, ok := assignable.(generated.Alternatives); ok {
		generateAssignableAlternatives(node, context, alts, cb)
	} else {
		panic("Unresolved assignment assignable")
	}
}

func generateAssignableAlternatives(node generator.Node, context *ParserGeneratorContext, alts generated.Alternatives, cb func(node generator.Node, resultName string)) {
	lookaheadName := context.orLookaheads[alts].name
	node.AppendLine("switch p.state.Lookahead(", lookaheadName, ") {")
	for i, alt := range alts.Alts() {
		node.AppendLine("case ", strconv.Itoa(i), ":")
		node.Indent(func(in generator.Node) {
			if assignable, ok := alt.(generated.Assignable); ok {
				generateAssignable(in, context, assignable, cb)
			}
		})
	}
	node.AppendLine("}")
}

func generateCrossReferenceParser(node generator.Node, context *ParserGeneratorContext, crossRef generated.CrossRef) string {
	return generateRuleCallParser(node, context, crossRef.Rule())
}

func generateGroupParser(node generator.Node, context *ParserGeneratorContext, group generated.Group) {
	lookaheadName := context.lookaheads[group].name
	lookahead := generateLookaheadString(lookaheadName)
	generateCardinality(node, func(n generator.Node) {
		for _, element := range group.Elements() {
			generateAbstractElementParser(n, context, element)
		}
	}, func(n generator.Node) {
		n.Append(lookahead)
	}, group.Cardinality())
}

func generateKeywordParser(node generator.Node, context *ParserGeneratorContext, keyword generated.Keyword) string {
	lookahead := "p.state.LA(1) == " + GeneratedKeywordIdxName(keyword)
	generateCardinality(node, func(n generator.Node) {
		n.AppendLine("token := p.state.Consume(", GeneratedKeywordIdxName(keyword), ")")
		n.AppendLine("core.AssignToken(node, token, ", context.accessNames[keyword], ")")
	}, func(n generator.Node) { n.Append(lookahead) }, keyword.Cardinality())
	return "token"
}

func generateRuleCallParser(node generator.Node, context *ParserGeneratorContext, ruleCall generated.RuleCall) string {
	token := getTokenWithName(context.grammar, ruleCall.Rule())
	rule := getRuleWithName(context.grammar, ruleCall.Rule())
	lookaheadName := context.lookaheads[ruleCall].name
	lookahead := generateLookaheadString(lookaheadName)
	var result string
	if token != nil {
		result = "token"
	} else if rule != nil {
		result = "result"
	} else {
		panic("Unresolved rule call: " + ruleCall.Rule())
	}
	first := true
	generateCardinality(node, func(n generator.Node) {
		eq := "="
		if first {
			eq = ":="
			first = false
		}
		if token != nil {
			n.AppendLine("token ", eq, " p.state.Consume(", GeneratedTokenIdxName(token), ")")
			n.AppendLine("core.AssignToken(node, token, ", context.accessNames[ruleCall], ")")
		} else if rule != nil {
			n.AppendLine("result ", eq, " p.Parse", rule.Name(), "()")
		}
	}, func(n generator.Node) { n.Append(lookahead) }, ruleCall.Cardinality())
	return result
}

func generateLookaheadString(name string) string {
	return "p.state.Lookahead(" + name + ") == 0"
}

func generateAlternativesParser(node generator.Node, context *ParserGeneratorContext, alts generated.Alternatives) {
	allLookaheadName := context.lookaheads[alts].name
	allLookahead := generateLookaheadString(allLookaheadName)
	lookaheadName := context.orLookaheads[alts].name
	generateCardinality(node, func(n generator.Node) {
		n.AppendLine("switch p.state.Lookahead(", lookaheadName, ") {")
		for i, alt := range alts.Alts() {
			n.AppendLine("case ", strconv.Itoa(i), ":")
			n.Indent(func(in generator.Node) {
				generateAbstractElementParser(in, context, alt)
			})
		}
		n.AppendLine("}")
	}, func(n generator.Node) {
		n.Append(allLookahead)
	}, alts.Cardinality())
}

func generateCardinality(node generator.Node, element, lookahead generator.Callback, cardinality string) {
	switch cardinality {
	case CardinalityOne:
		element(node)
	case CardinalityOptional:
		node.Append("if ")
		lookahead(node)
		node.AppendLine(" {")
		node.Indent(element)
		node.AppendLine("}")
	case CardinalityZeroOrMore:
		node.Append("for ")
		lookahead(node)
		node.AppendLine(" {")
		node.Indent(element)
		node.AppendLine("}")
	case CardinalityOneOrMore:
		node.Append("for ok := true; ok; ok = ")
		lookahead(node)
		node.AppendLine(" {")
		node.Indent(element)
		node.AppendLine("}")
	}
}

type PartialPathAndSuffixes struct {
	PartialPath []string
	Follow      []generated.Element
}

func remainingPathWith(nextDef []generated.Element, targetDef []generated.Element, i int) []generated.Element {
	targetSlice := targetDef[i+1:]
	arr := make([]generated.Element, 0, len(nextDef)+len(targetSlice))
	arr = append(arr, nextDef...)
	arr = append(arr, targetSlice...)
	return arr
}

func getAlternativesFor(grammar generated.Grammar, result []PartialPathAndSuffixes, elements []generated.Element, maxLength int, currPath []string) []PartialPathAndSuffixes {
	return append(result, possiblePathsFrom(grammar, maxLength, elements, currPath)...)
}

func IsOptionalCardinality(cardinality string) bool {
	return cardinality == CardinalityOptional || cardinality == CardinalityZeroOrMore
}

func possiblePathsFrom(grammar generated.Grammar, maxLength int, elements []generated.Element, currPath []string) []PartialPathAndSuffixes {
	result := []PartialPathAndSuffixes{}
	// Make a copy of currPath
	currPath = append([]string{}, currPath...)
	i := 0

	for len(currPath) < maxLength && i < len(elements) {
		element := elements[i]
		if IsOptionalCardinality(element.Cardinality()) {
			// Add path without this element
			result = getAlternativesFor(grammar, result, elements[i+1:], maxLength, currPath)
		}
		if group, ok := element.(generated.Group); ok {
			remain := remainingPathWith(group.Elements(), elements, i)
			return getAlternativesFor(grammar, result, remain, maxLength, currPath)
		} else if keyword, ok := element.(generated.Keyword); ok {
			currPath = append(currPath, GeneratedKeywordIdxName(keyword))
		} else if alts, ok := element.(generated.Alternatives); ok {
			for _, alt := range alts.Alts() {
				remain := remainingPathWith([]generated.Element{alt}, elements, i)
				result = getAlternativesFor(grammar, result, remain, maxLength, currPath)
			}
			return result
		} else if assignment, ok := element.(generated.Assignment); ok {
			if keyword, ok := assignment.Value().(generated.Keyword); ok {
				currPath = append(currPath, GeneratedKeywordIdxName(keyword))
			} else if ruleCall, ok := assignment.Value().(generated.RuleCall); ok {
				token := getTokenWithName(grammar, ruleCall.Rule())
				rule := getRuleWithName(grammar, ruleCall.Rule())
				if token != nil {
					currPath = append(currPath, GeneratedTokenIdxName(token))
				} else if rule != nil {
					remain := remainingPathWith([]generated.Element{rule.Body()}, elements, i)
					result = getAlternativesFor(grammar, result, remain, maxLength, currPath)
					return result
				}
			} else if crossRef, ok := assignment.Value().(generated.CrossRef); ok {
				token := getTokenWithName(grammar, crossRef.Rule().Rule())
				rule := getRuleWithName(grammar, crossRef.Rule().Rule())
				if token != nil {
					currPath = append(currPath, GeneratedTokenIdxName(token))
				} else if rule != nil {
					remain := remainingPathWith([]generated.Element{rule.Body()}, elements, i)
					result = getAlternativesFor(grammar, result, remain, maxLength, currPath)
					return result
				}
			} else if alts, ok := assignment.Value().(generated.Alternatives); ok {
				for _, alt := range alts.Alts() {
					remain := remainingPathWith([]generated.Element{alt}, elements, i)
					result = getAlternativesFor(grammar, result, remain, maxLength, currPath)
				}
				return result
			}
		} else if ruleCall, ok := element.(generated.RuleCall); ok {
			token := getTokenWithName(grammar, ruleCall.Rule())
			rule := getRuleWithName(grammar, ruleCall.Rule())
			if token != nil {
				currPath = append(currPath, GeneratedTokenIdxName(token))
			} else if rule != nil {
				remain := remainingPathWith([]generated.Element{rule.Body()}, elements, i)
				result = getAlternativesFor(grammar, result, remain, maxLength, currPath)
				return result
			}
		}
		i++
	}

	result = append(result, PartialPathAndSuffixes{
		PartialPath: currPath,
		Follow:      elements[i:],
	})
	return result
}

type LookaheadPath []string
type LookaheadOption []LookaheadPath
type LLkLookahead []LookaheadOption

func GetLLkLookaheadOr(grammar generated.Grammar, element generated.Alternatives) LLkLookahead {
	return generateCommonLLkLookahead(grammar, element.Alts())
}

func generateCommonLLkLookahead(grammar generated.Grammar, elements []generated.Element) LLkLookahead {
	lookahead := LLkLookahead{}
	for i := range 3 {
		lookahead = LLkLookahead{}
		for _, alt := range elements {
			option := generateLookaheadOption(grammar, alt, i+1)
			lookahead = append(lookahead, option)
		}
		if isUniqueLookahead(lookahead) {
			break
		}
	}
	return lookahead
}

func generateLLkLookahead(node generator.Node, name string, lookahead LLkLookahead) {
	node.AppendLine("var ", name, " = parser.LLkLookahead{")
	node.Indent(func(n generator.Node) {
		for _, option := range lookahead {
			n.AppendLine("parser.LookaheadOption{")
			n.Indent(func(in generator.Node) {
				for _, path := range option {
					in.AppendLine("parser.LookaheadPath{", strings.Join(path, ", "), "},")
				}
			})
			n.AppendLine("},")
		}
	})
	node.AppendLine("}")
}

func generateLookaheadOption(grammar generated.Grammar, element generated.Element, depth int) LookaheadOption {
	lookaheadOption := LookaheadOption{}
	partialPaths := possiblePathsFrom(grammar, depth, []generated.Element{element}, []string{})
	for _, partialPath := range partialPaths {
		if len(partialPath.PartialPath) > 0 {
			lookaheadOption = append(lookaheadOption, partialPath.PartialPath)
		}
	}
	return lookaheadOption
}

func GetLLkLookaheadOpt(grammar generated.Grammar, element generated.Element) LLkLookahead {
	// lookahead := LLkLookahead{}
	// for i := range 3 {
	// 	lookahead = LLkLookahead{}
	// 	for _, alt := range element.Alts {
	// 		option := generateLookaheadOption(grammar, alt, 1)
	// 		lookahead = append(lookahead, option)
	// 	}
	// 	if isUniqueLookahead(lookahead) {
	// 		break
	// 	}
	// }
	// Generate LL(1) decisions for options
	return LLkLookahead{
		generateLookaheadOption(grammar, element, 1),
	}
}

func isUniqueLookahead(lookahead LLkLookahead) bool {
	seen := make(map[string]bool)
	for _, option := range lookahead {
		localSeen := make(map[string]bool)
		for _, path := range option {
			key := strings.Join(path, ",")
			localSeen[key] = true
			if seen[key] {
				return false
			}
		}
		for key := range localSeen {
			seen[key] = true
		}
	}
	return true
}

func getTokenWithName(grammar generated.Grammar, name string) generated.Token {
	for _, t := range grammar.Terminals() {
		if t.Name() == name {
			return t
		}
	}
	return nil
}

func getRuleWithName(grammar generated.Grammar, name string) generated.ParserRule {
	for _, r := range grammar.Rules() {
		if r.Name() == name {
			return r
		}
	}
	return nil
}
