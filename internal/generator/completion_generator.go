// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"context"
	"sort"

	"typefox.dev/fastbelt"
	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/util/codegen"
)

func GenerateCompletion(grammr grammar.Grammar, packageName string) string {
	ctx := generateContext(grammr)
	node := NewRootNode()
	node.AppendLine("package ", packageName)
	node.AppendLine()
	node.AppendLine("import (")
	node.Indent(func(n codegen.Node) {
		n.AppendLine("\"context\"")
		n.AppendLine("\"iter\"")
		n.AppendLine()
		n.AppendLine("core \"typefox.dev/fastbelt\"")
		n.AppendLine("\"typefox.dev/fastbelt/parser\"")
		n.AppendLine("\"typefox.dev/fastbelt/util/service\"")
	})
	node.AppendLine(")")
	node.AppendLine()

	actions := collectActions(grammr)

	node.AppendNode(generateCompletionProvider(ctx))
	node.AppendNode(generateSyntheticFactories(ctx, grammr))
	node.AppendNode(generateCompletionDispatch(ctx))
	node.AppendNode(generateLspAdapter(ctx, actions))

	return FormatIfPossible(node.String())
}

// actionEntry captures the metadata needed to emit one ApplyAction case:
// "create an X and assign value to its field". Whether the field is single
// or array (and thus whether SetX or SetXItem is emitted) is derived from
// the grammar Field's type, not from the action's operator.
type actionEntry struct {
	targetType string // grammar interface name, e.g. "MemberCall"
	property   string // field on targetType, e.g. "Previous"
	isArray    bool   // true if the field is []T (use SetXItem), false for single T (use SetX)
	valueType  string // element type to cast value to, e.g. "MemberCall"
}

func collectActions(grammr grammar.Grammar) []actionEntry {
	seen := map[string]actionEntry{}
	for node := range fastbelt.AllNodes(grammr) {
		action, ok := node.(grammar.Action)
		if !ok {
			continue
		}
		typeRef := action.Type()
		propRef := action.Property()
		if typeRef == nil || propRef == nil {
			continue
		}
		iface := typeRef.Ref(context.TODO())
		if iface == nil {
			continue
		}
		field := propRef.Ref(context.TODO())
		if field == nil {
			continue
		}
		_, isArray := field.Type().(grammar.ArrayType)
		valueType := getTypeName(field.Type())
		if valueType == "" {
			continue
		}
		entry := actionEntry{
			targetType: iface.Name(),
			property:   field.Name(),
			isArray:    isArray,
			valueType:  valueType,
		}
		key := entry.targetType + "/" + entry.property
		seen[key] = entry
	}
	out := make([]actionEntry, 0, len(seen))
	for _, e := range seen {
		out = append(out, e)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].targetType != out[j].targetType {
			return out[i].targetType < out[j].targetType
		}
		return out[i].property < out[j].property
	})
	return out
}

// generateCompletionProvider emits the per-language Filter interface and its
// identity default implementation.
func generateCompletionProvider(ctx *LinkerGeneratorContext) codegen.Node {
	node := codegen.NewNode()
	name := ctx.grammar.Name()

	node.AppendLine("// ", name, "CompletionFilter lets adopters refine the candidates produced by")
	node.AppendLine("// the existing ", name, "ScopeProvider when completing a cross-reference.")
	node.AppendLine("type ", name, "CompletionFilter interface {")
	node.Indent(func(n codegen.Node) {
		for _, field := range ctx.fields {
			n.AppendLine("Filter", field.typeName, field.name,
				"(ctx context.Context, reference *core.Reference[", field.target,
				"], in iter.Seq[*core.SymbolDescription]) iter.Seq[*core.SymbolDescription]")
		}
	})
	node.AppendLine("}")
	node.AppendLine()

	node.AppendLine("type Default", name, "CompletionFilter struct{}")
	node.AppendLine()
	node.AppendLine("func NewDefault", name, "CompletionFilter() ", name, "CompletionFilter {")
	node.AppendLine(" return &Default", name, "CompletionFilter{}")
	node.AppendLine("}")
	node.AppendLine()

	for _, field := range ctx.fields {
		node.AppendLine("func (*Default", name, "CompletionFilter) Filter", field.typeName, field.name,
			"(_ context.Context, _ *core.Reference[", field.target,
			"], in iter.Seq[*core.SymbolDescription]) iter.Seq[*core.SymbolDescription] {")
		node.AppendLine(" return in")
		node.AppendLine("}")
		node.AppendLine()
	}
	return node
}

// Synthetic factories are required to produce the expected AST shape for completion
// when the user is completing an element which hasn't been produced by the parser yet.
func generateSyntheticFactories(ctx *LinkerGeneratorContext, grammr grammar.Grammar) codegen.Node {
	node := codegen.NewNode()
	name := ctx.grammar.Name()

	// One entry per parser rule (composite rules don't carry their own AST node).
	type ruleEntry struct {
		ruleName string
		typeName string
	}
	entries := []ruleEntry{}
	for _, rule := range grammr.Rules() {
		if returnType, ok := ruleReturnTypeName(rule); ok {
			entries = append(entries, ruleEntry{ruleName: rule.Name(), typeName: returnType})
		}
	}
	// Sort for stable output.
	sort.Slice(entries, func(i, j int) bool { return entries[i].ruleName < entries[j].ruleName })

	node.AppendLine("var ", name, "SyntheticFactories = map[string]func() core.AstNode{")
	node.Indent(func(n codegen.Node) {
		for _, e := range entries {
			n.AppendLine("\"", e.ruleName, "\": func() core.AstNode { return New", e.typeName, "() },")
		}
	})
	node.AppendLine("}")
	node.AppendLine()
	return node
}

// generateCompletionDispatch emits the hint-field -> dispatcher map. Each
// dispatcher resolves scope candidates for the in-progress reference and runs
// them through the language's CompletionProvider filter for that field.
func generateCompletionDispatch(ctx *LinkerGeneratorContext) codegen.Node {
	node := codegen.NewNode()
	name := ctx.grammar.Name()

	if len(ctx.fields) == 0 {
		node.AppendLine("var ", name, "CompletionDispatch = map[string]", name, "CompletionDispatchFunc{}")
		node.AppendLine()
		node.AppendLine("type ", name, "CompletionDispatchFunc func(ctx context.Context, sc *core.AstNode) iter.Seq[*core.SymbolDescription]")
		node.AppendLine()
		return node
	}

	node.AppendLine("type ", name, "CompletionDispatchFunc func(")
	node.AppendLine(" ctx context.Context,")
	node.AppendLine(" sc *service.Container,")
	node.AppendLine(" owner core.AstNode,")
	node.AppendLine(") iter.Seq[*core.SymbolDescription]")
	node.AppendLine()

	node.AppendLine("var ", name, "CompletionDispatch = map[string]", name, "CompletionDispatchFunc{")
	node.Indent(func(n codegen.Node) {
		for _, field := range ctx.fields {
			key := field.typeName + "." + field.name
			n.AppendLine("\"", key, "\": func(ctx context.Context, sc *service.Container, owner core.AstNode) iter.Seq[*core.SymbolDescription] {")
			n.Indent(func(in codegen.Node) {
				in.AppendLine("typedOwner, ok := owner.(", field.typeName, ")")
				in.AppendLine("if !ok {")
				in.AppendLine("	return func(yield func(*core.SymbolDescription) bool) {}")
				in.AppendLine("}")
				in.AppendLine("refs := service.MustGet[", name, "ReferencesConstructor](sc)")
				in.AppendLine("scopes := service.MustGet[", name, "ScopeProvider](sc)")
				in.AppendLine("filter := service.MustGet[", name, "CompletionFilter](sc)")
				in.AppendLine("ref := refs.", field.typeName, field.name, "(typedOwner, nil)")
				in.AppendLine("candidates := scopes.Scope", field.typeName, field.name, "(ctx, ref).AllElements()")
				in.AppendLine("return filter.Filter", field.typeName, field.name, "(ctx, ref, candidates)")
			})
			n.AppendLine("},")
		}
	})
	node.AppendLine("}")
	node.AppendLine()
	return node
}

func ruleReturnTypeName(rule grammar.ParserRule) (string, bool) {
	rt := grammar.FindReturnType(rule, context.Background())
	if rt == nil {
		return "", false
	}
	return rt.Name(), true
}

func generateLspAdapter(ctx *LinkerGeneratorContext, actions []actionEntry) codegen.Node {
	node := codegen.NewNode()
	name := ctx.grammar.Name()

	node.AppendLine("type ", name, "CompletionAdapter struct {")
	node.AppendLine(" sc *service.Container")
	node.AppendLine("}")
	node.AppendLine()

	node.AppendLine("func New", name, "CompletionAdapter(sc *service.Container) *", name, "CompletionAdapter {")
	node.AppendLine(" return &", name, "CompletionAdapter{sc: sc}")
	node.AppendLine("}")
	node.AppendLine()

	node.AppendLine("func (a *", name, "CompletionAdapter) Parse(tokens []core.Token) *parser.CompletionParseResult {")
	node.AppendLine(" return NewCompletionParser(a.sc).Parse(tokens)")
	node.AppendLine("}")
	node.AppendLine()

	node.AppendLine("func (a *", name, "CompletionAdapter) ATN() *parser.RuntimeATN {")
	node.AppendLine(" return BuildATN()")
	node.AppendLine("}")
	node.AppendLine()

	node.AppendLine("func (a *", name, "CompletionAdapter) SyntheticOwnerFor(ruleKey string) (core.AstNode, bool) {")
	node.AppendLine(" factory, ok := ", name, "SyntheticFactories[ruleKey]")
	node.AppendLine(" if !ok {")
	node.AppendLine("  return nil, false")
	node.AppendLine(" }")
	node.AppendLine("return factory(), true")
	node.AppendLine("}")
	node.AppendLine()

	node.AppendLine("func (a *", name, "CompletionAdapter) DispatchCompletion(ctx context.Context, field string, owner core.AstNode) (iter.Seq[*core.SymbolDescription], bool) {")
	if len(ctx.fields) == 0 {
		node.AppendLine("_ = field")
		node.AppendLine("_ = owner")
		node.AppendLine("_ = ctx")
		node.AppendLine("return nil, false")
	} else {
		node.AppendLine("dispatch, ok := ", name, "CompletionDispatch[field]")
		node.AppendLine("if !ok {")
		node.AppendLine(" return nil, false")
		node.AppendLine("}")
		node.AppendLine("return dispatch(ctx, a.sc, owner), true")
	}
	node.AppendLine("}")
	node.AppendLine()

	node.AppendNode(generateHasAssignment(ctx))
	node.AppendNode(generateApplyAction(ctx, actions))

	return node
}

// generateHasAssignment emits a typed-switch over every (interface, property)
// pair that can carry a cross-reference assignment in this grammar. The
// completion engine calls it to ask "is the hint's slot already filled on
// this owner?" - if yes, the PrecedingAction logic wraps the owner in a
// fresh synthetic parent; if no, the existing AST is the in-progress node
// and is used directly.
func generateHasAssignment(ctx *LinkerGeneratorContext) codegen.Node {
	node := codegen.NewNode()
	name := ctx.grammar.Name()

	// Index CR fields by interface so each type case emits one property switch.
	byIface := map[string][]LinkerField{}
	for _, f := range ctx.fields {
		byIface[f.typeName] = append(byIface[f.typeName], f)
	}
	ifaceNames := make([]string, 0, len(byIface))
	for n := range byIface {
		ifaceNames = append(ifaceNames, n)
	}
	sort.Strings(ifaceNames)

	node.AppendLine("func (a *", name, "CompletionAdapter) HasAssignment(node core.AstNode, property string) bool {")
	if len(ifaceNames) == 0 {
		node.AppendLine("_ = node")
		node.AppendLine("_ = property")
		node.AppendLine("return false")
		node.AppendLine("}")
		node.AppendLine()
		return node
	}
	node.AppendLine("switch n := node.(type) {")
	for _, iface := range ifaceNames {
		flds := byIface[iface]
		sort.Slice(flds, func(i, j int) bool { return flds[i].name < flds[j].name })
		node.AppendLine("case ", iface, ":")
		node.AppendLine("switch property {")
		for _, f := range flds {
			node.AppendLine("case \"", f.name, "\":")
			node.AppendLine("  return n.", f.name, "() != nil")
		}
		node.AppendLine("}")
	}
	node.AppendLine("}")
	node.AppendLine("return false")
	node.AppendLine("}")
	node.AppendLine()
	return node
}

// generateApplyAction emits a switch over every (actionType, property,
// operator) triple that appears in a `{Type.prop=current}` action anywhere
// in the grammar. The completion engine calls it to materialise the action
// when the main parser couldn't (e.g. the action's trigger token wasn't
// typed yet) - it allocates a new node and wires `value` into the named
// slot, mirroring what the main parser would have done on the next
// iteration.
func generateApplyAction(ctx *LinkerGeneratorContext, actions []actionEntry) codegen.Node {
	node := codegen.NewNode()
	name := ctx.grammar.Name()

	node.AppendLine("func (a *", name, "CompletionAdapter) ApplyAction(actionType, property string, value core.AstNode) core.AstNode {")
	if len(actions) == 0 {
		node.AppendLine("_ = actionType")
		node.AppendLine("_ = property")
		node.AppendLine("_ = value")
		node.AppendLine("return nil")
		node.AppendLine("}")
		node.AppendLine()
		return node
	}
	node.AppendLine("switch actionType {")
	// Group by target type so each case emits a single New<T>() and a
	// property switch.
	byTarget := map[string][]actionEntry{}
	for _, e := range actions {
		byTarget[e.targetType] = append(byTarget[e.targetType], e)
	}
	targets := make([]string, 0, len(byTarget))
	for t := range byTarget {
		targets = append(targets, t)
	}
	sort.Strings(targets)
	for _, target := range targets {
		node.AppendLine("case \"", target, "\":")
		node.AppendLine(" node := New", target, "()")
		node.AppendLine(" switch property {")
		entries := byTarget[target]
		sort.Slice(entries, func(i, j int) bool { return entries[i].property < entries[j].property })
		for _, e := range entries {
			setter := "Set" + e.property
			if e.isArray {
				setter += "Item"
			}
			node.AppendLine("case \"", e.property, "\":")
			node.AppendLine(" if v, ok := value.(", e.valueType, "); ok {")
			node.AppendLine("  node.", setter, "(v)")
			node.AppendLine("}")
		}
		node.AppendLine("}")
		node.AppendLine("return node")
	}
	node.AppendLine("}")
	node.AppendLine("return nil")
	node.AppendLine("}")
	node.AppendLine()
	return node
}
