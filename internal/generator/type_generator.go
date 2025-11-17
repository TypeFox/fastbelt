// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"strings"

	"github.com/TypeFox/langium-to-go/generator"
	"github.com/TypeFox/langium-to-go/internal/generated"
)

const TOKEN_TYPE = "*core.Token"

func GenerateTypes(grammar generated.Grammar) string {
	node := generator.NewNode()
	node.AppendLine("package generated")
	node.AppendLine()
	node.AppendLine("import (")
	node.Indent(func(n generator.Node) {
		n.AppendLine("\"github.com/TypeFox/langium-to-go/core\"")
	})
	node.AppendLine(")")
	node.AppendLine()
	for _, iface := range grammar.Interfaces() {
		generateInterface(node, grammar, iface)
	}
	return formatIfPossible(node.String())
}

var reservedKeywords = map[string]bool{
	"break":       true,
	"case":        true,
	"chan":        true,
	"const":       true,
	"continue":    true,
	"default":     true,
	"defer":       true,
	"else":        true,
	"fallthrough": true,
	"for":         true,
	"func":        true,
	"go":          true,
	"goto":        true,
	"if":          true,
	"import":      true,
	"interface":   true,
	"map":         true,
	"package":     true,
	"range":       true,
	"return":      true,
	"select":      true,
	"struct":      true,
	"switch":      true,
	"type":        true,
	"var":         true,
}

type FieldInfo struct {
	Name string
	// Private name, used to avoid conflicts with reserved keywords
	PName string

	Array          bool
	Boolean        bool
	Type           string
	HasTokenGetter bool
	// Generated type, e.g. *core.Token or an interface type.
	// Also used for setter methods.
	GType string
}

func getFieldInfo(field generated.Field) FieldInfo {
	name := field.Name()
	pname := strings.ToLower(name[0:1]) + name[1:]
	if reservedKeywords[pname] {
		pname = "_" + name
	}
	array := field.IsArray()
	typ := field.Type()
	gtype := typ
	hasTokenGetter := false
	boolean := false
	if typ == "string" || typ == "bool" {
		gtype = TOKEN_TYPE
		if !array {
			hasTokenGetter = true
		}
		if typ == "bool" {
			boolean = true
		}
	}
	return FieldInfo{
		Name:           name,
		PName:          pname,
		Array:          array,
		Type:           typ,
		HasTokenGetter: hasTokenGetter,
		Boolean:        boolean,
		GType:          gtype,
	}
}

func generateInterface(node generator.Node, grammar generated.Grammar, iface generated.Interface) {
	fields := []FieldInfo{}
	for _, field := range iface.Fields() {
		fields = append(fields, getFieldInfo(field))
	}
	node.AppendLine("type ", iface.Name(), " interface {")
	node.Indent(func(n generator.Node) {
		n.AppendLine("core.AstNode")
		for _, extends := range iface.Extends() {
			n.AppendLine(extends.Image)
		}
		n.AppendLine()
		n.AppendLine("Is", iface.Name(), "()")
		for _, field := range fields {
			var typeStr string
			prefix := ""
			if field.Array {
				typeStr = "[]" + field.GType
			} else {
				typeStr = field.Type
				if field.Boolean {
					prefix = "Is"
				}
			}
			// Getter
			n.AppendLine(prefix, field.Name, "() ", typeStr)
			// Token getter
			if field.HasTokenGetter {
				n.AppendLine(field.Name, "Token() ", TOKEN_TYPE)
			}
			// Setter
			if field.Array {
				n.AppendLine("With", field.Name, "Item(item ", field.GType, ")")
			} else {
				n.AppendLine("With", field.Name, "(value ", field.GType, ")")
			}
		}
	})
	node.AppendLine("}")
	node.AppendLine()
	node.AppendLine("func New", iface.Name(), "() ", iface.Name(), " {")
	node.Indent(func(n generator.Node) {
		n.AppendLine("return &", iface.Name(), "Impl{")
		n.Indent(func(n2 generator.Node) {
			n2.AppendLine("AstNodeBase: core.NewAstNode(),")
			for _, extend := range iface.Extends() {
				n2.AppendLine(extend.Image, "Data: New", extend.Image, "Data(),")
			}
			n2.AppendLine(iface.Name(), "Data: New", iface.Name(), "Data(),")
		})
		n.AppendLine("}")
	})
	node.AppendLine("}")
	node.AppendLine()
	generateDataStruct(node, iface, fields)
	generateImplStruct(node, iface, fields)
}

func generateImplStruct(node generator.Node, iface generated.Interface, fields []FieldInfo) {
	node.AppendLine("type ", iface.Name(), "Impl struct {")
	node.Indent(func(n generator.Node) {
		n.AppendLine("core.AstNodeBase")
		for _, extends := range iface.Extends() {
			n.AppendLine(extends.Image, "Data")
		}
		n.AppendLine(iface.Name(), "Data")
	})
	node.AppendLine("}")
	node.AppendLine()

	node.AppendLine("func (i *", iface.Name(), "Impl) ForEachNode(fn func(core.AstNode)) {")
	node.Indent(func(n generator.Node) {
		for _, extends := range iface.Extends() {
			n.AppendLine("i.", extends.Image, "Data.ForEachNode(fn)")
		}
		n.AppendLine("i.", iface.Name(), "Data.ForEachNode(fn)")
	})
	node.AppendLine("}")
	node.AppendLine()
}

func generateDataStruct(node generator.Node, iface generated.Interface, fields []FieldInfo) {
	node.AppendLine("type ", iface.Name(), "Data struct {")
	node.Indent(func(n generator.Node) {
		for _, field := range fields {
			var typeStr string
			if field.Array {
				typeStr = "[]" + field.GType
			} else {
				typeStr = field.GType
			}
			n.AppendLine(field.PName + " " + typeStr)
		}
	})
	node.AppendLine("}")
	node.AppendLine()
	node.AppendLine("func New", iface.Name(), "Data() ", iface.Name(), "Data {")
	node.Indent(func(n generator.Node) {
		n.AppendLine("return ", iface.Name(), "Data{")
		n.Indent(func(n2 generator.Node) {
			for _, field := range fields {
				if field.Array {
					n2.AppendLine(field.PName, ": []", field.GType, "{},")
				}
			}
		})
		n.AppendLine("}")
	})
	node.AppendLine("}")
	node.AppendLine()
	node.AppendLine("func (i *", iface.Name(), "Data) Is", iface.Name(), "() {}")
	node.AppendLine()
	node.AppendLine("func (i *", iface.Name(), "Data) ForEachNode(fn func(core.AstNode)) {")
	node.Indent(func(n generator.Node) {
		for _, field := range fields {
			name := field.PName
			if field.GType == TOKEN_TYPE {
				continue
			}
			if field.Array {
				n.AppendLine("for _, item := range i.", name, " {")
				n.Indent(func(n2 generator.Node) {
					n2.AppendLine("fn(item)")
				})
				n.AppendLine("}")
			} else {
				n.AppendLine("if i.", field.PName, " != nil {")
				n.Indent(func(n2 generator.Node) {
					n2.AppendLine("fn(i.", field.PName, ")")
				})
				n.AppendLine("}")
			}
		}
	})
	node.AppendLine("}")
	node.AppendLine()
	for _, field := range fields {
		// Getter
		getterName := field.Name
		if field.Boolean && !field.Array {
			getterName = "Is" + getterName
		}
		returnType := field.Type
		if field.Array {
			returnType = "[]" + field.GType
		}
		node.AppendLine("func (i *", iface.Name(), "Data) ", getterName, "() ", returnType, " {")
		node.Indent(func(n generator.Node) {
			if field.Array {
				// Arrays are always initialized
				n.AppendLine("return i.", field.PName)
			} else if field.Boolean {
				// Boolean fields return true if their token is present
				n.AppendLine("return i != nil && i.", field.PName, " != nil")
			} else {
				n.AppendLine("if i != nil && i.", field.PName, " != nil {")
				n.Indent(func(n2 generator.Node) {
					n2.Append("return i.", field.PName)
					if field.GType == TOKEN_TYPE {
						n2.Append(".Image")
					}
					n2.AppendLine()
				})
				n.AppendLine("} else {")
				n.Indent(func(n2 generator.Node) {
					defaultReturn := "nil"
					switch field.Type {
					case "string":
						defaultReturn = "\"\""
					case "bool":
						defaultReturn = "false"
					}
					n2.AppendLine("return ", defaultReturn)
				})
				n.AppendLine("}")
			}
		})
		node.AppendLine("}")
		node.AppendLine()

		// Token getter
		if field.HasTokenGetter {
			node.AppendLine("func (i *", iface.Name(), "Data) ", field.Name, "Token() ", TOKEN_TYPE, " {")
			node.Indent(func(n generator.Node) {
				n.AppendLine("return i.", field.PName)
			})
			node.AppendLine("}")
			node.AppendLine()
		}

		// Setter
		if field.Array {
			node.AppendLine("func (i *", iface.Name(), "Data) With", field.Name, "Item(item ", field.GType, ") {")
			node.Indent(func(n generator.Node) {
				n.AppendLine("i.", field.PName, " = append(i.", field.PName, ", item)")
			})
		} else {
			node.AppendLine("func (i *", iface.Name(), "Data) With", field.Name, "(value ", field.GType, ") {")
			node.Indent(func(n generator.Node) {
				n.AppendLine("i.", field.PName, " = value")
			})
		}
		node.AppendLine("}")
		node.AppendLine()
	}
}
