// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	ctx "context"
	"strconv"
	"strings"

	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/util/codegen"
)

// mustExpandInfixRules desugars all infix rules of the grammar (idempotent).
// Every generator entry point calls it so the synthesized rule bodies and
// operator token groups exist regardless of the order generators run in.
func mustExpandInfixRules(grammr grammar.Grammar) {
	if err := grammar.ExpandInfixRules(grammr); err != nil {
		panic(err.Error())
	}
}

// generateInfixParseFunction emits the Parse function of an infix rule.
//
// The completion parser reuses the generic element emitter over the
// synthesized flat body. The real parser feeds operands and operator tokens
// (a single token-group consume per operator, no alternatives decision) into
// the rule's parser.InfixFold field, which folds them into a binary tree as
// they arrive without any per-call slices.
func generateInfixParseFunction(node codegen.Node, context *ParserGeneratorContext, rule grammar.InfixRule, groupMembers map[string][]string) {
	returnType := grammar.FindReturnType(rule, ctx.Background())
	if returnType == nil {
		panic("Unable to find return type for infix rule: " + rule.Name())
	}
	context.inCompositeRule = false
	receiverType := context.receiverType()

	if context.completion {
		node.AppendLine("func (p *", receiverType, ") Parse", rule.Name(), "() {")
		node.Indent(func(n codegen.Node) {
			n.AppendLine("p.cp.EnterRule(", strconv.Quote(rule.Name()), ", ", context.atnData.ruleStart(rule), ")")
			n.AppendLine("defer p.cp.ExitRule()")
			generateAbstractElementParser(n, context, rule.Body())
		})
		node.AppendLine("}")
		node.AppendLine()
		return
	}

	// Destructure the body synthesized by grammar.ExpandInfixRules:
	// Group [ RuleCall(operand), Group* [ RuleCall(operatorGroup), RuleCall(operand) ] ]
	outer := rule.Body().(grammar.Group)
	firstCall := outer.Elements()[0].(grammar.RuleCall)
	inner := outer.Elements()[1].(grammar.Group)
	operatorCall := inner.Elements()[0].(grammar.RuleCall)
	secondCall := inner.Elements()[1].(grammar.RuleCall)
	operand := firstCall.Rule().Ref(ctx.Background()).(grammar.ParserRule)
	operatorGroup := operatorCall.Rule().Ref(ctx.Background()).(grammar.AbstractTokenRule)

	precedenceVar := rule.Name() + "Precedence"
	generateInfixPrecedenceVar(node, precedenceVar, rule, groupMembers)

	syncCall := buildSyncCall(context, inner)
	operatorState := context.atnData.elementStateName(operatorCall)
	fold := "p." + infixFoldFieldName(rule)
	build := "p.build" + rule.Name()
	emitOperand := func(n codegen.Node, call grammar.RuleCall) {
		n.AppendLine("{")
		n.Indent(func(in codegen.Node) {
			in.AppendLine("p.state.EnterRule(", context.atnData.followState(call), ")")
			in.AppendLine("result := p.Parse", operand.Name(), "()")
			in.AppendLine("p.state.ExitRule()")
			in.AppendLine(fold, ".Operand(result)")
		})
		n.AppendLine("}")
	}

	// build<Rule> constructs one binary node; it is handed to the fold as a
	// method value, which does not escape and therefore does not allocate.
	node.AppendLine("func (p *", receiverType, ") build", rule.Name(), "(left ", returnType.Name(), ", operator *core.Token, right ", returnType.Name(), ") ", returnType.Name(), " {")
	node.Indent(func(n codegen.Node) {
		n.AppendLine("result := New", rule.Name(), "()")
		n.AppendLine("if left != nil {")
		n.Indent(func(in codegen.Node) {
			in.AppendLine("result.SetLeft(left)")
			in.AppendLine("result.SetTextRangeStart(left.TextRange().Start)")
		})
		n.AppendLine("}")
		n.AppendLine("if operator != nil {")
		n.Indent(func(in codegen.Node) {
			in.AppendLine("core.AssignToken(result, operator, ", operatorState, ")")
			in.AppendLine("result.SetOperator(operator)")
			in.AppendLine("result.SetTextRangeEnd(operator.Range.End)")
		})
		n.AppendLine("}")
		n.AppendLine("if right != nil {")
		n.Indent(func(in codegen.Node) {
			in.AppendLine("result.SetRight(right)")
			in.AppendLine("result.SetTextRangeEnd(right.TextRange().End)")
		})
		n.AppendLine("}")
		n.AppendLine("return result")
	})
	node.AppendLine("}")
	node.AppendLine()

	node.AppendLine("func (p *", receiverType, ") Parse", rule.Name(), "() ", returnType.Name(), " {")
	node.Indent(func(n codegen.Node) {
		n.AppendLine("frame := ", fold, ".Begin()")
		emitOperand(n, firstCall)
		if syncCall != "" {
			n.AppendLine(syncCall)
		}
		n.AppendLine("for ", guardCall(context, inner), " {")
		n.Indent(func(in codegen.Node) {
			in.AppendLine("token := p.state.Consume(", GeneratedTokenName(operatorGroup), ")")
			in.AppendLine(fold, ".Operator(frame, token, ", precedenceVar, ", ", build, ")")
			emitOperand(in, secondCall)
			if syncCall != "" {
				in.AppendLine(syncCall)
			}
		})
		n.AppendLine("}")
		n.AppendLine("return ", fold, ".End(frame, ", build, ")")
	})
	node.AppendLine("}")
	node.AppendLine()
}

// infixFoldFieldName returns the name of the generated Parser field holding the
// parser.InfixFold of an infix rule.
func infixFoldFieldName(rule grammar.InfixRule) string {
	name := rule.Name()
	return strings.ToLower(name[:1]) + name[1:] + "Fold"
}

// generateInfixFoldFields emits one parser.InfixFold field per infix rule into
// the generated Parser struct.
func generateInfixFoldFields(node codegen.Node, grammr grammar.Grammar) {
	for _, rule := range grammr.InfixRules() {
		returnType := grammar.FindReturnType(rule, ctx.Background())
		if returnType == nil {
			panic("Unable to find return type for infix rule: " + rule.Name())
		}
		node.AppendLine(infixFoldFieldName(rule), " parser.InfixFold[", returnType.Name(), "]")
	}
}

// generateInfixPrecedenceVar emits the precedence table of an infix rule as a
// slice indexed by leaf token type id. Token-group operators are expanded to
// their leaf members so the runtime lookup is a single bounds-checked index on
// Token.Type.Id. Every operator the generated parser can consume (it consumes
// the rule's token group) has an entry.
func generateInfixPrecedenceVar(node codegen.Node, name string, rule grammar.InfixRule, groupMembers map[string][]string) {
	node.AppendLine("var ", name, " = []parser.InfixPrecedence{")
	node.Indent(func(n codegen.Node) {
		seen := map[string]bool{}
		for level, group := range rule.Groups() {
			right := ""
			if group.Associativity() == "right" {
				right = ", Right: true"
			}
			for _, operator := range group.Operators() {
				var operatorVar string
				switch op := operator.(type) {
				case grammar.Keyword:
					operatorVar = GeneratedTokenName(op)
				case grammar.RuleCall:
					operatorVar = GeneratedTokenName(op.Rule().Ref(ctx.Background()))
				default:
					continue
				}
				for _, leaf := range leafTokenVarNames(operatorVar, groupMembers, seen) {
					n.AppendLine(leaf, "_Idx: {Level: ", strconv.Itoa(level), right, "},")
				}
			}
		}
	})
	node.AppendLine("}")
	node.AppendLine()
}

// leafTokenVarNames recursively expands a token var name to its leaf member
// var names, resolving token groups through groupMembers. Names already in
// seen are skipped (and recorded), so duplicate operators produce a single
// map entry.
func leafTokenVarNames(varName string, groupMembers map[string][]string, seen map[string]bool) []string {
	if members, ok := groupMembers[varName]; ok {
		var result []string
		for _, member := range members {
			result = append(result, leafTokenVarNames(member, groupMembers, seen)...)
		}
		return result
	}
	if seen[varName] {
		return nil
	}
	seen[varName] = true
	return []string{varName}
}
