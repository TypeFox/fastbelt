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
	ValidateUniqueRuleName          = "uniqueRuleName"
	ValidateUniqueInterfaceName     = "uniqueInterfaceName"
	ValidateEmptyToken              = "emptyTerminalRule"
	ValidateEmptyKeyword            = "emptyKeyword"
	ValidateWhitespaceOnlyKeyword   = "whitespaceOnlyKeyword"
	ValidateKeywordWithWhitespace   = "keywordWithWhitespace"
	ValidateRuleReturnType          = "ruleReturnType"
	ValidateInterfaceExtends        = "interfaceExtends"
	ValidateRuleCallReturnType      = "ruleCallReturnType"
	ValidateRuleCallPosition        = "ruleCallPosition"
	ValidateActionAssignmentType    = "actionAssignmentType"
	ValidateActionPropertyType      = "actionPropertyType"
	ValidateAssignmentType          = "assignmentType"
	ValidateRecursiveTokenGroup     = "recursiveTokenGroup"
	ValidateInvalidTokenInGroup     = "invalidTokenInGroup"
	ValidateInvalidTokenInCrossRef  = "invalidTokenInCrossRef"
	ValidateMissingCrossRefTerminal = "missingCrossRefTerminal"
	ValidateUniqueFieldName         = "uniqueFieldName"
	ValidateFieldNameCapitalLetter  = "fieldNameCapitalLetter"
	ValidateReservedFieldName       = "reservedFieldName"
	ValidateNestedArrayType         = "nestedArrayType"
)

// reservedFieldNames lists field names that must not be used because they
// conflict with [core.AstNode] methods. Keep in sync with ast.go.
var reservedFieldNames = map[string]string{
	"Document":         "AstNode.Document",
	"Container":        "AstNode.Container",
	"Tokens":           "AstNode.Tokens",
	"Segment":          "AstNode.Segment",
	"Text":             "AstNode.Text",
	"ForEachNode":      "AstNode.ForEachNode",
	"ForEachReference": "AstNode.ForEachReference",
}

// GrammarImpl.Validate checks grammar-level constraints:
//   - Rule names must be unique within the grammar.
//   - Interface names must be unique within the grammar.
func (g *GrammarImpl) Validate(_ context.Context, _ string, accept core.ValidationAcceptor) {
	checkUniqueRuleNames(g, accept)
	checkUniqueInterfaceNames(g, accept)
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
		case Token:
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
	checkInvalidTokensInGroup(tg, accept)
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

func hiddenOrCommentTokenDescription(token Token) (description string, ok bool) {
	switch token.Type() {
	case "hidden":
		return "hidden", true
	case "comment":
		return "a comment", true
	default:
		return "", false
	}
}

func checkInvalidTokensInGroup(tg TokenGroup, accept core.ValidationAcceptor) {
	for _, ext := range tg.TokenRefs() {
		abstractToken := ext.Ref(context.Background())
		if abstractToken == nil {
			continue
		}
		if token, ok := abstractToken.(Token); ok {
			// Hidden/comment tokens are not allowed in token groups.
			// They are not meant to be consumed in parser rules,
			// and do not appear in the token slice.
			if description, special := hiddenOrCommentTokenDescription(token); special {
				accept(core.NewDiagnostic(
					core.SeverityError,
					fmt.Sprintf("The token '%s' cannot be used in a token group because it is %s.", token.Name(), description),
					tg,
					core.WithReference(ext),
					core.WithCode(ValidateInvalidTokenInGroup),
				))
			}
		}
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
	token, ok := resolved.(Token)
	if !ok {
		return
	}
	// Hidden/comment tokens are not allowed in cross-references because they
	// are not stored in the token slice and cannot identify named elements.
	if description, special := hiddenOrCommentTokenDescription(token); special {
		accept(core.NewDiagnostic(
			core.SeverityError,
			fmt.Sprintf("The token '%s' cannot be used in a cross-reference because it is %s.", token.Name(), description),
			cr,
			core.WithReference(ruleCall.Rule()),
			core.WithCode(ValidateInvalidTokenInCrossRef),
		))
	}
}
