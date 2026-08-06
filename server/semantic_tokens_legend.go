// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"slices"

	"typefox.dev/lsp"
)

// SemanticTokensLegendProvider provides the legend for semantic tokens LSP requests.
// Must be registered together with a [SemanticTokensProvider] in the service container
// to enable semantic tokens support for the language server.
type SemanticTokensLegendProvider interface {
	Legend() lsp.SemanticTokensLegend
}

// ExtendableSemanticTokensLegendProvider is a type of [SemanticTokensLegendProvider] that allows
// adding new token types and modifiers at runtime, extending the default legend.
//
// Note that this interface is useful beyond the ability to add new token types and modifiers.
// It also provides a convenient way to retrieve the index of the default token types and the bit
// values of default token modifiers, which can be used when implementing custom semantic token provider.
//
//	// Declaring a new legend provider with custom token types and modifiers
//	var legendProvider = server.NewExtendableSemanticTokensLegendProvider()
//	// Returns the index of the new token type within the legend
//	var extraType = legendProvider.AddType("extraTokenType")
//	// Returns the bit value of the new token modifier within the legend
//	var extraModifier = legendProvider.AddModifier("extraTokenModifier")
//	// Register within the service container
//	service.Put[server.SemanticTokensLegendProvider](sc, legendProvider)
type ExtendableSemanticTokensLegendProvider interface {
	SemanticTokensLegendProvider

	// Type returns the index of the "type" token type within the legend.
	Type() uint32
	// Class returns the index of the "class" token type within the legend.
	Class() uint32
	// Enum returns the index of the "enum" token type within the legend.
	Enum() uint32
	// Interface returns the index of the "interface" token type within the legend.
	Interface() uint32
	// Struct returns the index of the "struct" token type within the legend.
	Struct() uint32
	// TypeParameter returns the index of the "typeParameter" token type within the legend.
	TypeParameter() uint32
	// Parameter returns the index of the "parameter" token type within the legend.
	Parameter() uint32
	// Variable returns the index of the "variable" token type within the legend.
	Variable() uint32
	// Property returns the index of the "property" token type within the legend.
	Property() uint32
	// EnumMember returns the index of the "enumMember" token type within the legend.
	EnumMember() uint32
	// Event returns the index of the "event" token type within the legend.
	Event() uint32
	// Function returns the index of the "function" token type within the legend.
	Function() uint32
	// Method returns the index of the "method" token type within the legend.
	Method() uint32
	// Macro returns the index of the "macro" token type within the legend.
	Macro() uint32
	// Keyword returns the index of the "keyword" token type within the legend.
	Keyword() uint32
	// Modifier returns the index of the "modifier" token type within the legend.
	Modifier() uint32
	// Comment returns the index of the "comment" token type within the legend.
	Comment() uint32
	// String returns the index of the "string" token type within the legend.
	String() uint32
	// Number returns the index of the "number" token type within the legend.
	Number() uint32
	// Regexp returns the index of the "regexp" token type within the legend.
	Regexp() uint32
	// Operator returns the index of the "operator" token type within the legend.
	Operator() uint32
	// Decorator returns the index of the "decorator" token type within the legend.
	Decorator() uint32
	// Label returns the index of the "label" token type within the legend.
	Label() uint32
	// Namespace returns the index of the "namespace" token type within the legend.
	Namespace() uint32

	// Declaration returns the bit value of the "declaration" token modifier within the legend.
	Declaration() uint32
	// Definition returns the bit value of the "definition" token modifier within the legend.
	Definition() uint32
	// Readonly returns the bit value of the "readonly" token modifier within the legend.
	Readonly() uint32
	// Static returns the bit value of the "static" token modifier within the legend.
	Static() uint32
	// Deprecated returns the bit value of the "deprecated" token modifier within the legend.
	Deprecated() uint32
	// Abstract returns the bit value of the "abstract" token modifier within the legend.
	Abstract() uint32
	// Async returns the bit value of the "async" token modifier within the legend.
	Async() uint32
	// Modification returns the bit value of the "modification" token modifier within the legend.
	Modification() uint32
	// Documentation returns the bit value of the "documentation" token modifier within the legend.
	Documentation() uint32
	// DefaultLibrary returns the bit value of the "defaultLibrary" token modifier within the legend.
	DefaultLibrary() uint32

	// AddType adds a new token type to the legend and returns its index.
	AddType(name string) uint32
	// AddModifier adds a new token modifier to the legend and returns its bit value.
	AddModifier(name string) uint32
}

// NewExtendableSemanticTokensLegendProvider creates a new instance of [ExtendableSemanticTokensLegendProvider].
func NewExtendableSemanticTokensLegendProvider() ExtendableSemanticTokensLegendProvider {
	return &extendableSemanticTokensLegendProvider{}
}

func add(name string, existing *[]string) uint32 {
	index := uint32(len(*existing))
	*existing = append(*existing, name)
	return index
}

var defaultTokenTypes []string

var _type = add("type", &defaultTokenTypes)
var _class = add("class", &defaultTokenTypes)
var _enum = add("enum", &defaultTokenTypes)
var _interface = add("interface", &defaultTokenTypes)
var _struct = add("struct", &defaultTokenTypes)
var _typeParameter = add("typeParameter", &defaultTokenTypes)
var _parameter = add("parameter", &defaultTokenTypes)
var _variable = add("variable", &defaultTokenTypes)
var _property = add("property", &defaultTokenTypes)
var _enumMember = add("enumMember", &defaultTokenTypes)
var _event = add("event", &defaultTokenTypes)
var _function = add("function", &defaultTokenTypes)
var _method = add("method", &defaultTokenTypes)
var _macro = add("macro", &defaultTokenTypes)
var _keyword = add("keyword", &defaultTokenTypes)
var _modifier = add("modifier", &defaultTokenTypes)
var _comment = add("comment", &defaultTokenTypes)
var _string = add("string", &defaultTokenTypes)
var _number = add("number", &defaultTokenTypes)
var _regexp = add("regexp", &defaultTokenTypes)
var _operator = add("operator", &defaultTokenTypes)
var _decorator = add("decorator", &defaultTokenTypes)
var _label = add("label", &defaultTokenTypes)
var _namespace = add("namespace", &defaultTokenTypes)

var defaultTokenModifiers []string

var _declaration uint32 = 1 << add("declaration", &defaultTokenModifiers)
var _definition uint32 = 1 << add("definition", &defaultTokenModifiers)
var _readonly uint32 = 1 << add("readonly", &defaultTokenModifiers)
var _static uint32 = 1 << add("static", &defaultTokenModifiers)
var _deprecated uint32 = 1 << add("deprecated", &defaultTokenModifiers)
var _abstract uint32 = 1 << add("abstract", &defaultTokenModifiers)
var _async uint32 = 1 << add("async", &defaultTokenModifiers)
var _modification uint32 = 1 << add("modification", &defaultTokenModifiers)
var _documentation uint32 = 1 << add("documentation", &defaultTokenModifiers)
var _defaultLibrary uint32 = 1 << add("defaultLibrary", &defaultTokenModifiers)

type extendableSemanticTokensLegendProvider struct {
	types     []string
	modifiers []string
}

func (d *extendableSemanticTokensLegendProvider) Type() uint32          { return _type }
func (d *extendableSemanticTokensLegendProvider) Class() uint32         { return _class }
func (d *extendableSemanticTokensLegendProvider) Enum() uint32          { return _enum }
func (d *extendableSemanticTokensLegendProvider) Interface() uint32     { return _interface }
func (d *extendableSemanticTokensLegendProvider) Struct() uint32        { return _struct }
func (d *extendableSemanticTokensLegendProvider) TypeParameter() uint32 { return _typeParameter }
func (d *extendableSemanticTokensLegendProvider) Parameter() uint32     { return _parameter }
func (d *extendableSemanticTokensLegendProvider) Variable() uint32      { return _variable }
func (d *extendableSemanticTokensLegendProvider) Property() uint32      { return _property }
func (d *extendableSemanticTokensLegendProvider) EnumMember() uint32    { return _enumMember }
func (d *extendableSemanticTokensLegendProvider) Event() uint32         { return _event }
func (d *extendableSemanticTokensLegendProvider) Function() uint32      { return _function }
func (d *extendableSemanticTokensLegendProvider) Method() uint32        { return _method }
func (d *extendableSemanticTokensLegendProvider) Macro() uint32         { return _macro }
func (d *extendableSemanticTokensLegendProvider) Keyword() uint32       { return _keyword }
func (d *extendableSemanticTokensLegendProvider) Modifier() uint32      { return _modifier }
func (d *extendableSemanticTokensLegendProvider) Comment() uint32       { return _comment }
func (d *extendableSemanticTokensLegendProvider) String() uint32        { return _string }
func (d *extendableSemanticTokensLegendProvider) Number() uint32        { return _number }
func (d *extendableSemanticTokensLegendProvider) Regexp() uint32        { return _regexp }
func (d *extendableSemanticTokensLegendProvider) Operator() uint32      { return _operator }
func (d *extendableSemanticTokensLegendProvider) Decorator() uint32     { return _decorator }
func (d *extendableSemanticTokensLegendProvider) Label() uint32         { return _label }
func (d *extendableSemanticTokensLegendProvider) Namespace() uint32     { return _namespace }

func (d *extendableSemanticTokensLegendProvider) Declaration() uint32    { return _declaration }
func (d *extendableSemanticTokensLegendProvider) Definition() uint32     { return _definition }
func (d *extendableSemanticTokensLegendProvider) Readonly() uint32       { return _readonly }
func (d *extendableSemanticTokensLegendProvider) Static() uint32         { return _static }
func (d *extendableSemanticTokensLegendProvider) Deprecated() uint32     { return _deprecated }
func (d *extendableSemanticTokensLegendProvider) Abstract() uint32       { return _abstract }
func (d *extendableSemanticTokensLegendProvider) Async() uint32          { return _async }
func (d *extendableSemanticTokensLegendProvider) Modification() uint32   { return _modification }
func (d *extendableSemanticTokensLegendProvider) Documentation() uint32  { return _documentation }
func (d *extendableSemanticTokensLegendProvider) DefaultLibrary() uint32 { return _defaultLibrary }

func (d *extendableSemanticTokensLegendProvider) AddType(name string) uint32 {
	if slices.Contains(defaultTokenTypes, name) {
		panic("Cannot add a token type that already exists in the default legend: " + name)
	} else if slices.Contains(d.types, name) {
		panic("Cannot add a token type that already exists in the legend: " + name)
	}
	return add(name, &d.types) + uint32(len(defaultTokenTypes))
}

func (d *extendableSemanticTokensLegendProvider) AddModifier(name string) uint32 {
	if slices.Contains(defaultTokenModifiers, name) {
		panic("Cannot add a token modifier that already exists in the default legend: " + name)
	} else if slices.Contains(d.modifiers, name) {
		panic("Cannot add a token modifier that already exists in the legend: " + name)
	} else if (len(d.modifiers) + len(defaultTokenModifiers)) >= 32 {
		panic("Cannot add a token modifier because the legend already has 32 modifiers")
	}
	return 1 << (add(name, &d.modifiers) + uint32(len(defaultTokenModifiers)))
}

func (d *extendableSemanticTokensLegendProvider) Legend() lsp.SemanticTokensLegend {
	tokenTypes := make([]string, len(defaultTokenTypes)+len(d.types))
	copy(tokenTypes, defaultTokenTypes)
	copy(tokenTypes[len(defaultTokenTypes):], d.types)
	tokenModifiers := make([]string, len(defaultTokenModifiers)+len(d.modifiers))
	copy(tokenModifiers, defaultTokenModifiers)
	copy(tokenModifiers[len(defaultTokenModifiers):], d.modifiers)
	return lsp.SemanticTokensLegend{
		TokenTypes:     tokenTypes,
		TokenModifiers: tokenModifiers,
	}
}
