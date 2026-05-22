// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Package grammar documents the Fastbelt grammar language.
//
// A grammar definition file (extension .fb) describes the concrete syntax
// and type structure of a language. Fastbelt reads .fb files and generates
// Go lexers, parsers, and AST types from them. The grammar language is
// itself implemented as a Fastbelt grammar.
//
// See [typefox.dev/fastbelt] for the toolchain overview and the
// [typefox.dev/fastbelt/cmd/fastbelt] command for code generation.
//
// # Language Declaration
//
// Every grammar file begins with the grammar keyword, the language name,
// and a semicolon:
//
//	grammar MyLanguage;
//
// The name is used by Fastbelt when naming the generated package and
// language server.
//
// # Interface Declarations
//
// Every type that a parser rule creates must be declared explicitly as an
// interface. An interface lists the fields of that type along with their
// types:
//
//	interface Person {
//	    Name string
//	    Age  string
//	}
//
// The following field types are available:
//
//   - string — a string value set by a = or += assignment to a token rule.
//   - bool — a boolean value set by a ?= assignment.
//   - composite — a value set by a = or += assignment to a composite rule.
//   - TypeName — an embedded sub-object whose type is the named interface.
//   - *TypeName — a cross-reference to an object of the named interface type.
//   - []FieldType — an array of any of the above types.
//
// An interface can extend one or more other interfaces to inherit their
// fields:
//
//	interface NamedElement {
//	    Name string
//	}
//
//	interface Person extends NamedElement {
//	    Age string
//	}
//
// # Generated Go Types for Interfaces
//
// For every grammar interface, Fastbelt generates a set of Go types in
// `types_gen.go`:
//
//   - a public Go interface that embeds `core.AstNode` and declares
//     accessors/setters for each grammar field,
//   - a `<TypeName>Data` struct that stores the field data and implements
//     the accessor methods, and
//   - a `<TypeName>Impl` struct that embeds `core.AstNodeBase` plus
//     `<TypeName>Data`.
//
// The split between `Data` and `Impl` keeps type hierarchies composable.
// `Data` holds one layer of fields and behavior, while `Impl` embeds the
// data structs for all parent interfaces plus its own. This avoids copying
// inherited fields into every concrete type and keeps generator output
// regular as hierarchies grow.
//
// For example, this grammar interface:
//
//	interface Person {
//	    Name string
//	}
//
// generates a Go interface like:
//
//	type Person interface {
//	    core.AstNode
//	    IsPerson()
//	    Name() string
//	    NameToken() *core.Token
//	    SetName(value *core.Token)
//	}
//
// Field types in grammar interfaces map to generated Go API types as follows:
//
//   - string -> getter returns `string`, plus `FieldToken() *core.Token`.
//   - bool -> getter is named `IsField()` and returns `bool`.
//   - composite -> getter returns `string`; scalar composite fields also get a
//     `FieldNode() core.CompositeNode` accessor.
//   - TypeName -> getter returns `TypeName`.
//   - *TypeName -> getter returns `*core.Reference[TypeName]`.
//   - []FieldType -> getter returns a slice of the mapped element type.
//
// Slice fields use append-style setters (`Set<Field>Item`) in generated code.
// Scalar fields use `Set<Field>(value ...)`.
//
// # Token Rules
//
// Lexing transforms a stream of characters into a stream of tokens.
// A token rule matches a character sequence with a regular expression and
// assigns it a name. Token rule names are conventionally written in upper
// case:
//
//	token ID:  /[_a-zA-Z][\w_]*/;
//	token INT: /[0-9]+/;
//
// The order in which token rules are declared is significant: the lexer
// returns the first match, so more specific patterns should appear before
// more general ones.
//
// # Hidden Tokens
//
// Tokens that should be silently consumed — whitespace, comments — are
// declared with the hidden modifier:
//
//	hidden token WS:         /\s+/;
//	hidden token ML_COMMENT: /\/\*[\s\S]*?\*\//;
//	hidden token SL_COMMENT: /\/\/[^\r\n]*/;
//
// Hidden tokens are global and apply to the entire grammar.
//
// # Comment Tokens
//
// Documentation comments attached to the grammar element that follows them
// are declared with the comment modifier:
//
//	comment token SL_COMMENT: /\/\/[^\r\n]*/;
//
// Comment tokens are moved in front of the element they precede during
// parsing, making them available for tooling such as hover documentation in
// a language server.
//
// # Parser Rules
//
// Parser rules define what sequences of tokens are valid and how to populate
// the fields of the AST nodes they create. The first parser rule in the file
// is the entry rule: the starting point of the parse.
//
// A parser rule starts with its name, an optional returns clause naming the
// interface type it creates, a colon, the rule body, and a semicolon:
//
//	Person returns Person:
//	    "person" Name=ID;
//
// When the rule name matches a declared interface the returns clause may be
// omitted and Fastbelt resolves the type by name.
//
// # Cardinalities
//
// A cardinality suffix controls how many times an element may appear:
//
//   - (none) — exactly once
//   - ? — zero or one time
//   - * — zero or more times
//   - + — one or more times
//
// # Groups
//
// Elements in sequence form a group and must appear in the declared order:
//
//	Person returns Person:
//	    "person" Name=ID Address=Address;
//
// Parentheses create a sub-group that can carry its own cardinality:
//
//	State returns State:
//	    "state" Name=ID
//	        ("actions" "{" Actions+=[Command:ID]+ "}")?
//	    "end";
//
// # Alternatives
//
// The pipe operator | matches one of several options. Alternatives inside
// parentheses can carry a cardinality:
//
//	Model returns Model:
//	    (Persons+=Person | Greetings+=Greeting)*;
//
// # Keywords
//
// A keyword is a literal string in double quotes. Keywords guide the parser
// and provide visible structure to the language. They must not be empty and
// must not contain whitespace:
//
//	Person returns Person:
//	    "person" Name=ID "age" Age=INT;
//
// Keywords help the parser disambiguate between alternatives that would
// otherwise be identical:
//
//	interface Student { Name string }
//	interface Teacher { Name string }
//
//	Student returns Student:
//	    "student" Name=ID;
//
//	Teacher returns Teacher:
//	    "teacher" Name=ID;
//
//	Person:
//	    Student | Teacher;
//
// Without the "student" and "teacher" keywords the grammar would be
// ambiguous and the parser could not distinguish the two rules.
//
// # Assignments
//
// Assignments populate fields on the object being built by the surrounding
// rule. The left side names a field declared on the return type; the right
// side names what to parse. There are three forms:
//
// Single-value assignment (=) stores one parsed value in the field:
//
//	Person returns Person:
//	    "person" Name=ID;
//
// Array assignment (+=) appends each matched value to a slice field:
//
//	Model returns Model:
//	    Events+=Event*;
//
// Boolean assignment (?=) sets a bool field to true when the right side is
// consumed; the field remains false otherwise:
//
//	Employee returns Employee:
//	    "employee" Name=ID (Remote?="remote")?;
//
// Assignments with cardinality + or * form a contiguous group: the sequence
// of matched values must not be interrupted by elements belonging to a
// different assignment before the group is complete. Hidden tokens such as
// whitespace and comments may appear anywhere without interrupting a group.
//
// # Cross-References
//
// A cross-reference reads an identifying token from the input and resolves it
// to an existing object rather than creating a new one. The syntax is:
//
//	property=[Type:TOKEN]
//
// Type is the name of an interface and TOKEN is the name of a token rule that
// identifies objects of that type. If TOKEN is omitted, Fastbelt uses the
// token matched by the Name field assignment of the referenced type:
//
//	interface State {
//	    Name string
//	}
//
//	interface Transition {
//	    Event *Event
//	    State *State
//	}
//
//	Transition returns Transition:
//	    Event=[Event:ID] "=>" State=[State:ID];
//
// The linker resolves cross-references after parsing. If no object matching
// the token value is found in scope, a diagnostic error is reported.
//
// # Unassigned Rule Calls
//
// A rule call without an assignment delegates parsing to another rule without
// creating a new object in the current rule. The called rule is responsible
// for producing the object:
//
//	AbstractDefinition:
//	    Definition | DeclaredParameter;
//
// The parser rule AbstractDefinition does not create an object of its own.
// Instead it calls either Definition or DeclaredParameter, and whichever
// rule matches creates the object. This pattern is the standard way to write
// rules that match one of several concrete types.
//
// In contrast, an assigned rule call such as parameter=DeclaredParameter
// creates an object in the current rule and assigns the result of the called
// rule to the named property.
//
// # Actions
//
// Actions explicitly set the type of the object being built at the point
// where they appear in the rule body. A simple action creates a new object
// of the named type:
//
//	interface TypeOne { Name string }
//	interface TypeTwo extends TypeOne {}
//
//	RuleOne returns TypeOne:
//	    "one" Name=ID | {TypeTwo} "two" Name=ID;
//
// A tree-rewriting action creates a new object of the named type and assigns
// the object built so far to one of its properties. This technique handles
// structures that would require left recursion if written directly. The
// current object is referred to by the keyword current:
//
//	Addition returns Expression:
//	    SimpleExpr ({Addition.Left=current} "+" Right=SimpleExpr)*;
//
// When the "+" keyword is found, a new Addition object is created, the
// object parsed so far is stored in its Left property, and that new Addition
// becomes the current object. The operator += is also valid for tree-rewriting
// actions on slice properties.
//
// # Composite Rules
//
// A composite rule matches a structured token value such as a qualified name
// or dotted path. Unlike parser rules, composite rules do not create AST
// objects; they yield an opaque composite value that is stored in a field of
// type composite.
//
// Composite rules support keywords, rule calls, parenthesized alternatives,
// and cardinalities, but not assignments or cross-references:
//
//	composite QualifiedName: ID ("." ID)*;
//
// Declare the receiving field as composite in the interface, then use a
// normal = or += assignment on the right side of a parser rule:
//
//	interface TypeRef {
//	    Name composite
//	}
//
//	TypeRef returns TypeRef:
//	    Name=QualifiedName;
package grammar
