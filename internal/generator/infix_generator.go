// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	ctx "context"
	"strconv"

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
// synthesized flat body. The real parser collects operands and operator
// tokens in flat slices with a single token-group consume per operator (no
// alternatives decision), then folds them into a binary tree with
// parser.BuildInfixTree.
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
	emitOperand := func(n codegen.Node, call grammar.RuleCall) {
		n.AppendLine("{")
		n.Indent(func(in codegen.Node) {
			in.AppendLine("p.state.EnterRule(", context.atnData.followState(call), ")")
			in.AppendLine("result := p.Parse", operand.Name(), "()")
			in.AppendLine("p.state.ExitRule()")
			in.AppendLine("parts = append(parts, result)")
		})
		n.AppendLine("}")
	}

	node.AppendLine("func (p *", receiverType, ") Parse", rule.Name(), "() ", returnType.Name(), " {")
	node.Indent(func(n codegen.Node) {
		n.AppendLine("var parts []", returnType.Name())
		n.AppendLine("var operators []*core.Token")
		emitOperand(n, firstCall)
		if syncCall != "" {
			n.AppendLine(syncCall)
		}
		n.AppendLine("for ", guardCall(context, inner), " {")
		n.Indent(func(in codegen.Node) {
			in.AppendLine("if token := p.state.Consume(", GeneratedTokenName(operatorGroup), "); token != nil {")
			in.Indent(func(iin codegen.Node) {
				iin.AppendLine("operators = append(operators, token)")
			})
			in.AppendLine("}")
			emitOperand(in, secondCall)
			if syncCall != "" {
				in.AppendLine(syncCall)
			}
		})
		n.AppendLine("}")
		n.AppendLine("if len(parts) == 1 {")
		n.Indent(func(in codegen.Node) {
			in.AppendLine("// A lone operand is returned unchanged; no binary node is created.")
			in.AppendLine("return parts[0]")
		})
		n.AppendLine("}")
		n.AppendLine("return parser.BuildInfixTree(parts, operators, ", precedenceVar, ", func(left ", returnType.Name(), ", operator *core.Token, right ", returnType.Name(), ") ", returnType.Name(), " {")
		n.Indent(func(in codegen.Node) {
			in.AppendLine("result := New", rule.Name(), "()")
			in.AppendLine("if left != nil {")
			in.Indent(func(iin codegen.Node) {
				iin.AppendLine("result.SetLeft(left)")
				iin.AppendLine("result.SetTextRangeStart(left.TextRange().Start)")
			})
			in.AppendLine("}")
			in.AppendLine("if operator != nil {")
			in.Indent(func(iin codegen.Node) {
				iin.AppendLine("core.AssignToken(result, operator, ", operatorState, ")")
				iin.AppendLine("result.SetOperator(operator)")
				iin.AppendLine("result.SetTextRangeEnd(operator.Range.End)")
			})
			in.AppendLine("}")
			in.AppendLine("if right != nil {")
			in.Indent(func(iin codegen.Node) {
				iin.AppendLine("result.SetRight(right)")
				iin.AppendLine("result.SetTextRangeEnd(right.TextRange().End)")
			})
			in.AppendLine("}")
			in.AppendLine("return result")
		})
		n.AppendLine("})")
	})
	node.AppendLine("}")
	node.AppendLine()
}

// generateInfixPrecedenceVar emits the precedence table of an infix rule,
// keyed by leaf token type id. Token-group operators are expanded to their
// leaf members so the runtime lookup is a single map access on Token.TypeId.
func generateInfixPrecedenceVar(node codegen.Node, name string, rule grammar.InfixRule, groupMembers map[string][]string) {
	node.AppendLine("var ", name, " = map[int]parser.InfixPrecedence{")
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
