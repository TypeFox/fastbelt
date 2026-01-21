package generator

import (
	"context"

	"typefox.dev/fastbelt/generator"
	"typefox.dev/fastbelt/internal/grammar/generated"
)

type LinkerGeneratorContext struct {
	grammar generated.Grammar
	fields  []LinkerField
}

type LinkerField struct {
	typeName string
	name     string
	target   string
}

func GenerateLinker(grammar generated.Grammar) string {
	node := generator.NewNode()
	node.AppendLine("package generated")
	node.AppendLine()
	node.AppendLine("import (")
	node.Indent(func(n generator.Node) {
		n.AppendLine("\"context\"")
		n.AppendLine()
		n.AppendLine("core \"typefox.dev/fastbelt\"")
		n.AppendLine("\"typefox.dev/fastbelt/linking\"")
	})
	node.AppendLine(")")
	node.AppendLine()

	context := generateContext(grammar)
	node.AppendNode(generateScopeProvider(context))
	node.AppendNode(generateLinker(context))
	node.AppendNode(generateReferenceGenerator(context))

	return formatIfPossible(node.String())
}

func generateContext(grammar generated.Grammar) *LinkerGeneratorContext {
	fields := []LinkerField{}
	for _, iface := range grammar.Interfaces() {
		for _, field := range iface.Fields() {
			if refType, ok := field.Type().(generated.ReferenceType); ok {
				fields = append(fields, LinkerField{
					typeName: iface.Name(),
					name:     field.Name(),
					target:   refType.Type().Ref(context.TODO()).Name(),
				})
			} else if arrayType, ok := field.Type().(generated.ArrayType); ok {
				if refType, ok := arrayType.InternalType().(generated.ReferenceType); ok {
					fields = append(fields, LinkerField{
						typeName: iface.Name(),
						name:     field.Name(),
						target:   refType.Type().Ref(context.TODO()).Name(),
					})
				}
			}
		}
	}
	return &LinkerGeneratorContext{
		grammar: grammar,
		fields:  fields,
	}
}

func generateScopeProvider(context *LinkerGeneratorContext) generator.Node {
	node := generator.NewNode()
	node.AppendLine("type ", context.grammar.Name(), "ScopeProvider interface {")
	node.Indent(func(n generator.Node) {
		for _, field := range context.fields {
			n.AppendLine("Scope", field.typeName, field.name, "(ctx context.Context, reference *core.Reference[", field.target, "]) core.Scope")
		}
	})
	node.AppendLine("}")
	node.AppendLine()

	node.AppendLine("type Default", context.grammar.Name(), "ScopeProvider struct {")
	node.AppendLine("	srv ", context.grammar.Name(), "LinkingSrvCont")
	node.AppendLine("}")
	node.AppendLine()

	node.AppendLine("func NewDefault", context.grammar.Name(), "ScopeProvider(srv ", context.grammar.Name(), "LinkingSrvCont) *Default", context.grammar.Name(), "ScopeProvider {")
	node.AppendLine("	return &Default", context.grammar.Name(), "ScopeProvider{srv: srv}")
	node.AppendLine("}")
	node.AppendLine()

	for _, field := range context.fields {
		node.AppendLine("func (s *Default", context.grammar.Name(), "ScopeProvider) Scope", field.typeName, field.name, "(ctx context.Context, reference *core.Reference[", field.target, "]) core.Scope {")
		node.AppendLine("	return linking.LocalScopeOfType[", field.target, "](reference.Owner, s.srv.Linking().SymbolTable.LocalSymbols)")
		node.AppendLine("}").AppendLine()
	}
	return node
}

func generateLinker(context *LinkerGeneratorContext) generator.Node {
	node := generator.NewNode()
	node.AppendLine("type ", context.grammar.Name(), "Linker interface {")
	node.Indent(func(n generator.Node) {
		for _, field := range context.fields {
			n.AppendLine("Link", field.typeName, field.name, "(ctx context.Context, reference *core.Reference[", field.target, "]) (*core.AstNodeDescription, *core.ReferenceError)")
		}
	})
	node.AppendLine("}")
	node.AppendLine()

	node.AppendLine("type Default", context.grammar.Name(), "Linker struct {")
	node.AppendLine("	srv ", context.grammar.Name(), "LinkingSrvCont")
	node.AppendLine("}")
	node.AppendLine()

	node.AppendLine("func NewDefault", context.grammar.Name(), "Linker(srv ", context.grammar.Name(), "LinkingSrvCont) *Default", context.grammar.Name(), "Linker {")
	node.AppendLine("	return &Default", context.grammar.Name(), "Linker{srv: srv}")
	node.AppendLine("}")
	node.AppendLine()

	for _, field := range context.fields {
		node.AppendLine("func (l *Default", context.grammar.Name(), "Linker) Link", field.typeName, field.name, "(ctx context.Context, reference *core.Reference[", field.target, "]) (*core.AstNodeDescription, *core.ReferenceError) {")
		node.AppendLine("    scope := l.srv.", context.grammar.Name(), "Linking().ScopeProvider.Scope", field.typeName, field.name, "(ctx, reference)")
		node.AppendLine("    return core.DefaultLink(scope, reference.Text)")
		node.AppendLine("}").AppendLine()
	}
	return node
}

func generateReferenceGenerator(context *LinkerGeneratorContext) generator.Node {
	node := generator.NewNode()
	node.AppendLine("type ", context.grammar.Name(), "ReferenceGenerator interface {")
	node.Indent(func(n generator.Node) {
		for _, field := range context.fields {
			n.AppendLine(field.typeName, field.name, "(owner core.AstNode, token *core.Token) *core.Reference[", field.target, "]")
		}
	})
	node.AppendLine("}")
	node.AppendLine()

	node.AppendLine("type Default", context.grammar.Name(), "ReferenceGenerator struct {")
	node.AppendLine("	srv ", context.grammar.Name(), "LinkingSrvCont")
	node.AppendLine("}")
	node.AppendLine()

	node.AppendLine("func NewDefault", context.grammar.Name(), "ReferenceGenerator(srv ", context.grammar.Name(), "LinkingSrvCont) *Default", context.grammar.Name(), "ReferenceGenerator {")
	node.AppendLine("	return &Default", context.grammar.Name(), "ReferenceGenerator{srv: srv}")
	node.AppendLine("}")
	node.AppendLine()

	for _, field := range context.fields {
		node.AppendLine("func (g *Default", context.grammar.Name(), "ReferenceGenerator) ", field.typeName, field.name, "(owner core.AstNode, token *core.Token) *core.Reference[", field.target, "] {")
		node.AppendLine("    fn := g.srv.", context.grammar.Name(), "Linking().Linker.Link", field.typeName, field.name)
		node.AppendLine("    return core.NewReference(owner, token, fn)")
		node.AppendLine("}").AppendLine()
	}
	return node
}
