// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"strings"

	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/util/codegen"
)

func GenerateJSON(grammar grammar.Grammar, packageName string) string {
	node := NewRootNode()
	node.AppendLine("package ", packageName)
	node.AppendLine()
	node.AppendLine("import (")
	node.Indent(func(n codegen.Node) {
		n.AppendLine("\"bytes\"")
		n.AppendLine("\"encoding/json/jsontext\"")
		n.AppendLine("\"encoding/json/v2\"")
		n.AppendLine()
		n.AppendLine("core \"typefox.dev/fastbelt\"")
		n.AppendLine("\"typefox.dev/fastbelt/util\"")
	})
	node.AppendLine(")")
	node.AppendLine()

	generateFunctions(grammar, node)

	return FormatIfPossible(node.String())
}

func generateFunctions(grammar grammar.Grammar, node codegen.Node) codegen.Node {
	for _, iface := range grammar.Interfaces() {
		generateJSONMarshalTo(node, iface)
	}

	for _, iface := range grammar.Interfaces() {
		generateJSONUnmarshalFrom(node, iface)
	}

	generateDispatchingUnmarshalFunc(node, grammar)

	return node
}

func generateJSONMarshalTo(node codegen.Node, iface grammar.Interface) {
	fields := collectAllFields(iface, map[string]struct{}{})

	node.AppendLine("func (i *", iface.Name(), "Impl) MarshalJSONTo(_encoder *jsontext.Encoder) error {")
	node.Indent(func(n codegen.Node) {
		// preprocess []string fields: they are stored as []*core.Token internally but shall be serialized as plain strings;
		// instead of implementing 'MarshalJSON()' for *core.Token, a preprocessing loop calling '.String()' for each item
		// is added for each of such fields
		var stringListFields = map[string]string{}
		for _, field := range fields {
			if field.Array && (field.GType == TOKEN_TYPE || field.GType == COMPOSITE_TYPE) {
				varName := field.PName
				stringListFields[field.Name] = varName

				n.AppendLine(varName, " := make([]string, len(i.", field.Name, "()))")
				n.AppendLine("for j, item := range i.", field.Name, "() {")
				n.Indent(func(n2 codegen.Node) {
					n2.AppendLine(varName, "[j] = item.String()")
				})
				n.AppendLine("}")
			}
		}
		n.AppendLine("return json.MarshalEncode(_encoder, struct {")
		n.Indent(func(n2 codegen.Node) {
			n2.AppendLine("T__", " ", "string", " `json:\"$type\"`")
			for _, field := range fields {
				typeStr := field.Type
				if field.Array {
					typeStr = "[]" + typeStr
				}
				n2.AppendLine(field.Name, " ", typeStr, " `json:\"", field.JsonPropName, ",omitempty\"`")
			}
		})
		n.AppendLine("}{")
		n.Indent(func(n2 codegen.Node) {
			n2.AppendLine("T__: ", "\"", iface.Name(), "\",")
			for _, field := range fields {
				if varName, present := stringListFields[field.Name]; present {
					n2.AppendLine(field.Name, ": ", varName, ",")
				} else {
					getterName := field.Name
					if field.Boolean && !field.Array {
						getterName = "Is" + field.Name
					}
					n2.AppendLine(field.Name, ": i.", getterName, "(),")
				}
			}
		})
		n.AppendLine("})")
	})
	node.AppendLine("}")
	node.AppendLine()
}

func getAuxFieldType(field FieldInfo) string {
	var typ string
	if field.Boolean {
		typ = "bool"
	} else if field.GType == TOKEN_TYPE || field.GType == COMPOSITE_TYPE {
		typ = "string"
	} else {
		typ = "jsontext.Value"
	}

	if field.Array {
		return "[]" + typ
	} else {
		return typ
	}
}

func generateJSONUnmarshalFrom(node codegen.Node, iface grammar.Interface) {
	genCreateNewToken := func(valueRef string) string {
		return "token := core.NewToken(nil, " + valueRef + ", -1, -1)"
	}

	fields := collectAllFields(iface, map[string]struct{}{})

	node.AppendLine("func (i *", iface.Name(), "Impl) UnmarshalJSONFrom(_decoder *jsontext.Decoder) error {")

	if len(fields) == 0 {
		node.Indent(func(n2 codegen.Node) {
			n2.AppendLine("return nil")
		})
		node.AppendLine("}")
		node.AppendLine()
		return
	}

	node.Indent(func(n codegen.Node) {
		n.AppendLine("aux := &struct {")
		n.Indent(func(n2 codegen.Node) {
			n2.AppendLine("T__", " ", "string", " `json:\"$type\"`")
			for _, field := range fields {
				n2.AppendLine(field.Name, " ", getAuxFieldType(field), " `json:\"", field.JsonPropName, "\"`")
			}
		})
		n.AppendLine("}{}")
		n.AppendLine("if err := json.UnmarshalDecode(_decoder, aux); err != nil {")
		genReturnErr(n)

		for _, field := range fields {
			if field.Array {
				n.AppendLine("i.", field.PName, " = make([]", field.GType, ", 0, len(aux."+field.Name+"))")
				n.AppendLine("for _, item := range aux.", field.Name, " {")
				n.Indent(func(n2 codegen.Node) {
					switch field.GType {
					case TOKEN_TYPE:
						n2.AppendLine("{")
						n2.Indent(func(n3 codegen.Node) {
							if field.Boolean {
								n3.AppendLine(genCreateNewToken("\"\""))
							} else {
								n3.AppendLine(genCreateNewToken("item"))
							}
							n3.AppendLine("i.Set", field.Name, "Item(&token)")
						})
						n2.AppendLine("}")
					case COMPOSITE_TYPE:
						n2.AppendLine("{")
						n2.Indent(func(n3 codegen.Node) {
							n3.AppendLine(genCreateNewToken("item"))
							n3.AppendLine("cn := core.NewCompositeNode()")
							n3.AppendLine("cn.AppendToken(&token)")
							n3.AppendLine("i.Set", field.Name, "Item(cn)")
						})
						n2.AppendLine("}")

					default:
						if field.Reference {
							genUnmarshalReference(n2, field, "item", "node", true)
						} else {
							genUnmarshalChild(n2, field, "item", "node", true)
						}
					}
				})
				n.AppendLine("}")
			} else if field.Boolean {
				n.AppendLine("if aux.", field.Name, "{")
				n.Indent(func(n2 codegen.Node) {
					n2.AppendLine(genCreateNewToken("\"\""))
					n2.AppendLine("i.Set", field.Name, "(&token)")
				})
				n.AppendLine("}")
			} else if field.GType == TOKEN_TYPE {
				n.AppendLine("{")
				n.Indent(func(n2 codegen.Node) {
					n2.AppendLine(genCreateNewToken("aux." + field.Name))
					n2.AppendLine("i.Set", field.Name, "(&token)")
				})
				n.AppendLine("}")
			} else if field.GType == COMPOSITE_TYPE {
				n.AppendLine("{")
				n.Indent(func(n2 codegen.Node) {
					n2.AppendLine(genCreateNewToken("aux." + field.Name))
					n2.AppendLine("cn := core.NewCompositeNode()")
					n2.AppendLine("cn.AppendToken(&token)")
					n2.AppendLine("i.Set", field.Name, "(cn)")
				})
				n.AppendLine("}")
			} else if field.Reference {
				genUnmarshalReference(n, field, "aux."+field.Name, field.PName, false)
			} else {
				genUnmarshalChild(n, field, "aux."+field.Name, field.PName, false)
			}
		}
		n.AppendLine("return nil")
	})
	node.AppendLine("}")
	node.AppendLine()
}

func genUnmarshalReference(node codegen.Node, field FieldInfo, srcName string, targetName string, loopItem bool) {
	genUnmarshalFieldContent(node, field, srcName, targetName, loopItem, func(body codegen.Node) {
		body.AppendLine(targetName, " := core.NewReference[", strings.Split(field.Type, "[")[1], "(i, nil, nil)")
		// for the sake simplicity and performance we call 'target.UnmarshalJSON()' directly instead of taking the route
		// via json.Unmarshal(...), since Reference implements that method
		// note: the generic impl has special handling for "RawMessage" being equal "null", sets the target pointer to "nil"
		body.AppendLine("if err := ", targetName, ".UnmarshalJSON(", srcName, "); err != nil {")
		genReturnErr(body)
	})
}

func genUnmarshalChild(node codegen.Node, field FieldInfo, srcName string, targetName string, loopItem bool) {
	genUnmarshalFieldContent(node, field, srcName, targetName, loopItem, func(body codegen.Node) {
		body.AppendLine(targetName, ", err := UnmarshalValue[", field.Type, "](", srcName, ")")
		body.AppendLine("if err != nil {")
		genReturnErr(body)
	})
}

func genUnmarshalFieldContent(node codegen.Node, field FieldInfo, srcName string, targetName string, loopItem bool, unmarshalBody func(codegen.Node)) {
	if loopItem {
		unmarshalBody(node)
		node.AppendLine("i.Set", field.Name, "Item(", targetName, ")")
	} else {
		node.AppendLine("if ", srcName, " != nil {")
		node.Indent(func(n2 codegen.Node) {
			unmarshalBody(node)
			node.AppendLine("i.Set", field.Name, "(", targetName, ")")
		})
		node.AppendLine("}")
	}
}

func genReturnErr(node codegen.Node) {
	node.Indent(func(n2 codegen.Node) {
		n2.AppendLine("return err")
	})
	node.AppendLine("}")
}

func generateDispatchingUnmarshalFunc(node codegen.Node, g grammar.Grammar) {
	node.AppendLine(`// UnmarshalValue is a sugar method delegating to [util.UnmarshalDecode]`)
	node.AppendLine(`// that decodes 'value' into an instance of type 'T' by reading the "$type" field,`)
	node.AppendLine(`// selecting a corresponding factory, creating an instance, and unmarshaling its content.`)
	node.AppendLine("func UnmarshalValue[T core.AstNode](value jsontext.Value) (T, error) {")
	node.Indent(func(n codegen.Node) {
		n.AppendLine("return util.UnmarshalDecode[T](jsontext.NewDecoder(bytes.NewReader(value)), ", g.Name(), "SyntheticFactories)")
	})
	node.AppendLine("}")
	node.AppendLine()
}
