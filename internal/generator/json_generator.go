// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generator

import (
	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/util/codegen"
)

func GenerateJSON(grammar grammar.Grammar, packageName string) string {
	node := NewRootNode()
	node.AppendLine("package ", packageName)
	node.AppendLine()
	node.AppendLine("import (")
	node.Indent(func(n codegen.Node) {
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
	thisRef := "_this"
	thisDot := thisRef + "."

	node.AppendLine("func (", thisRef, " *", iface.Name(), "Impl) MarshalJSONTo(_encoder *jsontext.Encoder) error {")
	node.Indent(func(n codegen.Node) {
		// preprocess token/composite fields: they are stored as *core.Token / core.CompositeNode
		// internally but shall be serialized as plain strings; instead of implementing
		// 'MarshalJSON()' for *core.Token, a preprocessing loop/block calling '.String()' is added
		// for each of such fields. Non-array fields preprocess into a *string local instead of a
		// plain string, nil when the field's underlying token/composite is absent, so marshaling
		// can tell that apart from an explicitly present but empty string, matching how
		// UnmarshalJSONFrom revives them (array fields don't need this: a missing array decodes
		// to an empty slice, indistinguishable from - and treated the same as - an empty one).
		var stringListFields = map[string]string{}
		var stringFields = map[string]string{}
		for _, field := range fields {
			if field.Boolean || (field.GType != TOKEN_TYPE && field.GType != COMPOSITE_TYPE) {
				continue
			}
			varName := field.LName
			if field.Array {
				stringListFields[field.Name] = varName

				n.AppendLine(varName, " := make([]string, len(", thisDot, field.Name, "()))")
				n.AppendLine("for _j, _item := range ", thisDot, field.Name, "() {")
				n.Indent(func(n2 codegen.Node) {
					n2.AppendLine(varName, "[_j] = _item.String()")
				})
				n.AppendLine("}")
			} else {
				stringFields[field.Name] = varName

				presenceGetter := field.Name + "Token"
				if field.GType == COMPOSITE_TYPE {
					presenceGetter = field.Name + "Node"
				}
				n.AppendLine("var ", varName, " *string")
				n.AppendLine("if ", thisDot, presenceGetter, "() != nil {")
				n.Indent(func(n2 codegen.Node) {
					n2.AppendLine("_v := ", thisDot, field.Name, "()")
					n2.AppendLine(varName, " = &_v")
				})
				n.AppendLine("}")
			}
		}
		n.AppendLine("return json.MarshalEncode(_encoder, struct {")
		n.Indent(func(n2 codegen.Node) {
			n2.AppendLine("T__", " ", "string", " `json:\"$type\"`")
			for _, field := range fields {
				if field.Array {
					n2.AppendLine(field.Name, " []", field.Type, " `json:\"", field.JsonPropName, ",omitempty\"`")
				} else if _, present := stringFields[field.Name]; present {
					// omitempty also drops a non-nil pointer to an empty string (it checks the
					// pointee, not the pointer), which would undo the absent/empty distinction
					// this *string encoding exists for; omitzero only checks the pointer itself.
					n2.AppendLine(field.Name, " *", field.Type, " `json:\"", field.JsonPropName, ",omitzero\"`")
				} else {
					n2.AppendLine(field.Name, " ", field.Type, " `json:\"", field.JsonPropName, ",omitempty\"`")
				}
			}
		})
		n.AppendLine("}{")
		n.Indent(func(n2 codegen.Node) {
			n2.AppendLine("T__: ", "\"", iface.Name(), "\",")
			for _, field := range fields {
				if varName, present := stringListFields[field.Name]; present {
					n2.AppendLine(field.Name, ": ", varName, ",")
				} else if varName, present := stringFields[field.Name]; present {
					n2.AppendLine(field.Name, ": ", varName, ",")
				} else {
					getterName := field.Name
					if field.Boolean && !field.Array {
						getterName = "Is" + field.Name
					}
					n2.AppendLine(field.Name, ": ", thisDot, getterName, "(),")
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
		// non-array token/composite fields use *string so that an omitted JSON property
		// (nil) can be told apart from an explicitly present, empty string; array fields
		// don't need this, since a missing array property already decodes to a nil/empty
		// slice distinct from a present-but-empty one (both loops below produce no items)
		if !field.Array {
			typ = "*" + typ
		}
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
	fields := collectAllFields(iface, map[string]struct{}{})

	node.AppendLine("func (_this *", iface.Name(), "Impl) UnmarshalJSONFrom(_decoder *jsontext.Decoder) error {")
	thisRef := "_this"
	thisDotSet := thisRef + ".Set"
	errRef := "_err"
	genCreateNewToken := func(valueRef string) string {
		return "core.NewSyntheticToken(" + valueRef + ", " + thisRef + ")"
	}

	if len(fields) == 0 {
		node.Indent(func(n2 codegen.Node) {
			// json.UnmarshalerFrom requires reading exactly one value; with no fields to
			// populate there's nothing to decode into, but the value must still be consumed.
			n2.AppendLine("return _decoder.SkipValue()")
		})
		node.AppendLine("}")
		node.AppendLine()
		return
	}

	node.Indent(func(n codegen.Node) {
		n.AppendLine("aux := &struct {")
		n.Indent(func(n2 codegen.Node) {
			for _, field := range fields {
				n2.AppendLine(field.Name, " ", getAuxFieldType(field), " `json:\"", field.JsonPropName, "\"`")
			}
		})
		n.AppendLine("}{}")
		n.AppendLine("if ", errRef, " := json.UnmarshalDecode(_decoder, aux); ", errRef, " != nil {")
		genReturnErr(n, errRef)

		for _, field := range fields {
			if field.Array {
				n.AppendLine(thisRef, ".", field.PName, " = make([]", field.GType, ", 0, len(aux."+field.Name+"))")
				n.AppendLine("for _, item := range aux.", field.Name, " {")
				n.Indent(func(n2 codegen.Node) {
					switch field.GType {
					case TOKEN_TYPE:
						// Note: 'if field.Boolean' will never be true here, since we cannot parse lists of present test results ('?=')
						//  and there're no other ways of creating pure boolean values in the grammar, although '[]bool' is a valid type
						// all other primitive values are of type 'string' or 'composite'
						n2.AppendLine(thisDotSet, field.Name, "Item(", genCreateNewToken("item"), ")")
					case COMPOSITE_TYPE:
						n2.AppendLine("cn := core.NewCompositeNode()")
						n2.AppendLine("cn.AppendToken(", genCreateNewToken("item"), ")")
						n2.AppendLine(thisDotSet, field.Name, "Item(cn)")
					default:
						if field.Reference {
							genUnmarshalReference(n2, field, "item", "reference", thisRef, thisDotSet, errRef, true)
						} else {
							genUnmarshalChild(n2, field, "item", "node", thisDotSet, errRef, true)
						}
					}
				})
				n.AppendLine("}")
			} else if field.Boolean {
				n.AppendLine("if aux.", field.Name, "{")
				n.Indent(func(n2 codegen.Node) {
					n2.AppendLine(thisDotSet, field.Name, "(", genCreateNewToken("\"\""), ")")
				})
				n.AppendLine("}")
			} else if field.GType == TOKEN_TYPE {
				n.AppendLine("if aux.", field.Name, " != nil {")
				n.Indent(func(n2 codegen.Node) {
					n2.AppendLine(thisDotSet, field.Name, "(", genCreateNewToken("*aux."+field.Name), ")")
				})
				n.AppendLine("}")
			} else if field.GType == COMPOSITE_TYPE {
				n.AppendLine("if aux.", field.Name, " != nil {")
				n.Indent(func(n2 codegen.Node) {
					n2.AppendLine("cn := core.NewCompositeNode()")
					n2.AppendLine("cn.AppendToken(", genCreateNewToken("*aux."+field.Name), ")")
					n2.AppendLine(thisDotSet, field.Name, "(cn)")
				})
				n.AppendLine("}")
			} else if field.Reference {
				genUnmarshalReference(n, field, "aux."+field.Name, field.PName, thisRef, thisDotSet, errRef, false)
			} else {
				genUnmarshalChild(n, field, "aux."+field.Name, field.PName, thisDotSet, errRef, false)
			}
		}
		n.AppendLine("return nil")
	})
	node.AppendLine("}")
	node.AppendLine()
}

func genUnmarshalReference(
	node codegen.Node,
	field FieldInfo,
	srcName string,
	targetName string,
	thisRef string,
	thisDotSet string,
	errRef string,
	loopItem bool,
) {
	genUnmarshalFieldContent(node, field, srcName, targetName, thisDotSet, loopItem, func(body codegen.Node) {
		// for the sake simplicity and performance we call 'target.UnmarshalJSON()' directly instead of taking the route
		// via json.Unmarshal(...), since Reference implements that method
		// note: the generic impl has special handling for "RawMessage" being equal "null", sets the target pointer to "nil"
		body.AppendLine(targetName, ", ", errRef, " := util.UnmarshalReference[", field.RefTypeArg, "](", thisRef, ", ", srcName, ")")
		body.AppendLine("if ", errRef, " != nil {")
		genReturnErr(body, errRef)
	})
}

func genUnmarshalChild(
	node codegen.Node,
	field FieldInfo,
	srcName string,
	targetName string,
	thisDotSet string,
	errRef string,
	loopItem bool,
) {
	genUnmarshalFieldContent(node, field, srcName, targetName, thisDotSet, loopItem, func(body codegen.Node) {
		body.AppendLine(targetName, ", ", errRef, " := UnmarshalValue[", field.Type, "](", srcName, ")")
		body.AppendLine("if ", errRef, " != nil {")
		genReturnErr(body, errRef)
	})
}

func genUnmarshalFieldContent(
	node codegen.Node,
	field FieldInfo,
	srcName string,
	targetName string,
	thisDotSet string,
	loopItem bool,
	unmarshalBody func(codegen.Node),
) {
	if loopItem {
		unmarshalBody(node)
		node.AppendLine("if ", targetName, " != nil {")
		node.Indent(func(n2 codegen.Node) {
			node.AppendLine(thisDotSet, field.Name, "Item(", targetName, ")")
		})
		node.AppendLine("}")
	} else {
		node.AppendLine("if ", srcName, " != nil {")
		node.Indent(func(n2 codegen.Node) {
			unmarshalBody(node)
			node.AppendLine(thisDotSet, field.Name, "(", targetName, ")")
		})
		node.AppendLine("}")
	}
}

func genReturnErr(node codegen.Node, errRef string) {
	node.Indent(func(n2 codegen.Node) {
		n2.AppendLine("return ", errRef)
	})
	node.AppendLine("}")
}

func generateDispatchingUnmarshalFunc(node codegen.Node, g grammar.Grammar) {
	node.AppendLine(`// UnmarshalValue is a sugar method delegating to [util.UnmarshalValue]`)
	node.AppendLine(`// that decodes 'value' into an instance of type 'T' by reading the "$type" field,`)
	node.AppendLine(`// selecting a corresponding factory, creating an instance, and unmarshaling its content.`)
	node.AppendLine("func UnmarshalValue[T core.AstNode](value jsontext.Value) (T, error) {")
	node.Indent(func(n codegen.Node) {
		n.AppendLine("return util.UnmarshalValue[T](value, ", g.Name(), "SyntheticFactories)")
	})
	node.AppendLine("}")
	node.AppendLine()
}
