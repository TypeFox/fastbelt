// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	core "typefox.dev/fastbelt"
)

// GrammarImpl.Validate checks grammar-level constraints:
//   - Rule names must be unique within the grammar.
//   - Interface names must be unique within the grammar.
func (g *GrammarImpl) Validate(_ context.Context, _ string, accept core.ValidationAcceptor) {
	checkUniqueRuleNames(g, accept)
	checkUniqueInterfaceNames(g, accept)
}

func checkUniqueRuleNames(g Grammar, accept core.ValidationAcceptor) {
	seen := map[string][]core.NamedNode{}
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
	for name, nodes := range seen {
		if len(nodes) > 1 {
			for _, node := range nodes {
				accept(core.NewDiagnostic(
					core.SeverityError,
					fmt.Sprintf("A rule's name has to be unique. '%s' is used multiple times.", name),
					node,
					core.WithToken(node.NameToken()),
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
				))
			}
		}
	}
}

// TokenImpl.Validate checks terminal rule constraints:
//   - The regular expression should not match the empty string.
func (t *TokenImpl) Validate(_ context.Context, _ string, accept core.ValidationAcceptor) {
	checkEmptyTerminalRule(t, accept)
}

func checkEmptyTerminalRule(t Token, accept core.ValidationAcceptor) {
	raw := t.Regexp()
	if raw == "" {
		return
	}
	// Strip surrounding slashes from the regex literal
	pattern := raw
	if len(pattern) >= 2 && pattern[0] == '/' && pattern[len(pattern)-1] == '/' {
		pattern = pattern[1 : len(pattern)-1]
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		return
	}
	if re.MatchString("") {
		accept(core.NewDiagnostic(
			core.SeverityError,
			"This terminal could match an empty string.",
			t,
			core.WithToken(t.NameToken()),
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
		))
	} else if strings.TrimSpace(value) == "" {
		accept(core.NewDiagnostic(
			core.SeverityError,
			"Keywords cannot only consist of whitespace characters.",
			k,
			core.WithToken(k.ValueToken()),
		))
	} else if strings.ContainsAny(value, " \t\n\r") {
		accept(core.NewDiagnostic(
			core.SeverityWarning,
			"Keywords should not contain whitespace characters.",
			k,
			core.WithToken(k.ValueToken()),
		))
	}
}

func (rule *ParserRuleImpl) Validate(ctx context.Context, _ string, accept core.ValidationAcceptor) {
	checkRuleReturnType(rule, ctx, accept)
}

func checkRuleReturnType(rule ParserRule, _ context.Context, accept core.ValidationAcceptor) {
	// Only search if not explicitly provided
	if rule.ReturnType() == nil && rule.Name() != "" {
		grammar := rule.Container().(Grammar)
		if grammar == nil {
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
				),
			)
		}
	}
}

func (i *InterfaceImpl) Validate(ctx context.Context, _ string, accept core.ValidationAcceptor) {
	checkInterfaceExtends(i, ctx, accept)
}

func checkInterfaceExtends(iface Interface, ctx context.Context, accept core.ValidationAcceptor) {
	for _, ext := range iface.Extends() {
		extType := ext.Ref(ctx)
		if extType == nil {
			continue
		}
		if appearsInExtends(iface, extType, ctx, map[string]struct{}{}) {
			accept(core.NewDiagnostic(
				core.SeverityError,
				"An interface cannot extend itself, neither directly nor indirectly.",
				iface,
				core.WithToken(ext.Token()),
			))
		}
	}
}

func appearsInExtends(target Interface, current Interface, ctx context.Context, visited map[string]struct{}) bool {
	if current.Name() == target.Name() {
		return true
	}
	if _, alreadyVisited := visited[current.Name()]; alreadyVisited {
		return false
	}
	visited[current.Name()] = struct{}{}
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
			))
		}
	}
}

func checkRuleCallPosition(call RuleCall, _ context.Context, accept core.ValidationAcceptor) {
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
					))
					return
				}
				if _, ok := elem.(Assignment); ok {
					accept(core.NewDiagnostic(
						core.SeverityError,
						"An unassigned rule call cannot be preceded by an assignment.",
						call,
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
			core.WithToken(a.Type().Token()),
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
				core.WithToken(a.Property().Token()),
			))
		}
	} else {
		accept(core.NewDiagnostic(
			core.SeverityError,
			"Cannot assign a parser rule to a non-interface field.",
			a,
			core.WithToken(a.Property().Token()),
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
					core.WithToken(v.Type().Token()),
				))
			}
		} else {
			accept(core.NewDiagnostic(
				core.SeverityError,
				"Cannot assign a cross-reference value to a non-reference field.",
				v,
			))
		}
	case RuleCall:
		resolvedRule := v.Rule().Ref(ctx)
		if resolvedRule == nil {
			return
		}
		switch rule := resolvedRule.(type) {
		case Token:
			if primitiveType, ok := fieldType.(PrimitiveType); !ok || primitiveType.Type() != "string" {
				accept(core.NewDiagnostic(
					core.SeverityError,
					"Cannot assign a token to a non-string field.",
					v,
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
					))
				}
			} else {
				accept(core.NewDiagnostic(
					core.SeverityError,
					"Cannot assign a parser rule to a non-interface field.",
					v,
				))
			}
		}
	case Keyword:
		if primitiveType, ok := fieldType.(PrimitiveType); !ok || primitiveType.Type() != "string" {
			accept(core.NewDiagnostic(
				core.SeverityError,
				"Cannot assign a keyword value to a non-string field.",
				v,
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
	return doInterfaceIsAssignableTo(source, target, map[string]struct{}{})
}

func doInterfaceIsAssignableTo(source Interface, target Interface, visited map[string]struct{}) bool {
	if source.Name() == target.Name() {
		return true
	}
	if _, alreadyVisited := visited[source.Name()]; alreadyVisited {
		return false
	}
	visited[source.Name()] = struct{}{}
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
