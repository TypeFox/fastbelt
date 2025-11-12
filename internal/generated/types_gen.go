package generated

import (
	"github.com/TypeFox/langium-to-go/ast"
	"github.com/TypeFox/langium-to-go/lexer"
)

type Grammar interface {
	ast.AstNode

	IsGrammar()
	Name() string
	NameToken() *lexer.Token
	WithName(value *lexer.Token)
	Rules() []ParserRule
	WithRulesItem(item ParserRule)
	Terminals() []Token
	WithTerminalsItem(item Token)
	Interfaces() []Interface
	WithInterfacesItem(item Interface)
}

func NewGrammar() Grammar {
	return &GrammarImpl{
		AstNodeBase: ast.NewAstNode(),
		GrammarData: NewGrammarData(),
	}
}

type GrammarData struct {
	name       *lexer.Token
	rules      []ParserRule
	terminals  []Token
	interfaces []Interface
}

func NewGrammarData() GrammarData {
	return GrammarData{
		rules:      []ParserRule{},
		terminals:  []Token{},
		interfaces: []Interface{},
	}
}

func (i *GrammarData) IsGrammar() {}

func (i *GrammarData) ForEachNode(fn func(ast.AstNode)) {
	for _, item := range i.rules {
		fn(item)
	}
	for _, item := range i.terminals {
		fn(item)
	}
	for _, item := range i.interfaces {
		fn(item)
	}
}

func (i *GrammarData) Name() string {
	if i != nil && i.name != nil {
		return i.name.Image
	} else {
		return ""
	}
}

func (i *GrammarData) NameToken() *lexer.Token {
	return i.name
}

func (i *GrammarData) WithName(value *lexer.Token) {
	i.name = value
}

func (i *GrammarData) Rules() []ParserRule {
	return i.rules
}

func (i *GrammarData) WithRulesItem(item ParserRule) {
	i.rules = append(i.rules, item)
}

func (i *GrammarData) Terminals() []Token {
	return i.terminals
}

func (i *GrammarData) WithTerminalsItem(item Token) {
	i.terminals = append(i.terminals, item)
}

func (i *GrammarData) Interfaces() []Interface {
	return i.interfaces
}

func (i *GrammarData) WithInterfacesItem(item Interface) {
	i.interfaces = append(i.interfaces, item)
}

type GrammarImpl struct {
	ast.AstNodeBase
	GrammarData
}

func (i *GrammarImpl) ForEachNode(fn func(ast.AstNode)) {
	i.GrammarData.ForEachNode(fn)
}

type Interface interface {
	ast.AstNode

	IsInterface()
	Name() string
	NameToken() *lexer.Token
	WithName(value *lexer.Token)
	Extends() []*lexer.Token
	WithExtendsItem(item *lexer.Token)
	Fields() []Field
	WithFieldsItem(item Field)
}

func NewInterface() Interface {
	return &InterfaceImpl{
		AstNodeBase:   ast.NewAstNode(),
		InterfaceData: NewInterfaceData(),
	}
}

type InterfaceData struct {
	name    *lexer.Token
	extends []*lexer.Token
	fields  []Field
}

func NewInterfaceData() InterfaceData {
	return InterfaceData{
		extends: []*lexer.Token{},
		fields:  []Field{},
	}
}

func (i *InterfaceData) IsInterface() {}

func (i *InterfaceData) ForEachNode(fn func(ast.AstNode)) {
	for _, item := range i.fields {
		fn(item)
	}
}

func (i *InterfaceData) Name() string {
	if i != nil && i.name != nil {
		return i.name.Image
	} else {
		return ""
	}
}

func (i *InterfaceData) NameToken() *lexer.Token {
	return i.name
}

func (i *InterfaceData) WithName(value *lexer.Token) {
	i.name = value
}

func (i *InterfaceData) Extends() []*lexer.Token {
	return i.extends
}

func (i *InterfaceData) WithExtendsItem(item *lexer.Token) {
	i.extends = append(i.extends, item)
}

func (i *InterfaceData) Fields() []Field {
	return i.fields
}

func (i *InterfaceData) WithFieldsItem(item Field) {
	i.fields = append(i.fields, item)
}

type InterfaceImpl struct {
	ast.AstNodeBase
	InterfaceData
}

func (i *InterfaceImpl) ForEachNode(fn func(ast.AstNode)) {
	i.InterfaceData.ForEachNode(fn)
}

type Field interface {
	ast.AstNode

	IsField()
	Name() string
	NameToken() *lexer.Token
	WithName(value *lexer.Token)
	IsArray() bool
	ArrayToken() *lexer.Token
	WithArray(value *lexer.Token)
	Type() string
	TypeToken() *lexer.Token
	WithType(value *lexer.Token)
}

func NewField() Field {
	return &FieldImpl{
		AstNodeBase: ast.NewAstNode(),
		FieldData:   NewFieldData(),
	}
}

type FieldData struct {
	name  *lexer.Token
	array *lexer.Token
	_Type *lexer.Token
}

func NewFieldData() FieldData {
	return FieldData{}
}

func (i *FieldData) IsField() {}

func (i *FieldData) ForEachNode(fn func(ast.AstNode)) {
}

func (i *FieldData) Name() string {
	if i != nil && i.name != nil {
		return i.name.Image
	} else {
		return ""
	}
}

func (i *FieldData) NameToken() *lexer.Token {
	return i.name
}

func (i *FieldData) WithName(value *lexer.Token) {
	i.name = value
}

func (i *FieldData) IsArray() bool {
	return i != nil && i.array != nil
}

func (i *FieldData) ArrayToken() *lexer.Token {
	return i.array
}

func (i *FieldData) WithArray(value *lexer.Token) {
	i.array = value
}

func (i *FieldData) Type() string {
	if i != nil && i._Type != nil {
		return i._Type.Image
	} else {
		return ""
	}
}

func (i *FieldData) TypeToken() *lexer.Token {
	return i._Type
}

func (i *FieldData) WithType(value *lexer.Token) {
	i._Type = value
}

type FieldImpl struct {
	ast.AstNodeBase
	FieldData
}

func (i *FieldImpl) ForEachNode(fn func(ast.AstNode)) {
	i.FieldData.ForEachNode(fn)
}

type ParserRule interface {
	ast.AstNode

	IsParserRule()
	Name() string
	NameToken() *lexer.Token
	WithName(value *lexer.Token)
	ReturnType() string
	ReturnTypeToken() *lexer.Token
	WithReturnType(value *lexer.Token)
	Body() Alternatives
	WithBody(value Alternatives)
}

func NewParserRule() ParserRule {
	return &ParserRuleImpl{
		AstNodeBase:    ast.NewAstNode(),
		ParserRuleData: NewParserRuleData(),
	}
}

type ParserRuleData struct {
	name       *lexer.Token
	returnType *lexer.Token
	body       Alternatives
}

func NewParserRuleData() ParserRuleData {
	return ParserRuleData{}
}

func (i *ParserRuleData) IsParserRule() {}

func (i *ParserRuleData) ForEachNode(fn func(ast.AstNode)) {
	if i.body != nil {
		fn(i.body)
	}
}

func (i *ParserRuleData) Name() string {
	if i != nil && i.name != nil {
		return i.name.Image
	} else {
		return ""
	}
}

func (i *ParserRuleData) NameToken() *lexer.Token {
	return i.name
}

func (i *ParserRuleData) WithName(value *lexer.Token) {
	i.name = value
}

func (i *ParserRuleData) ReturnType() string {
	if i != nil && i.returnType != nil {
		return i.returnType.Image
	} else {
		return ""
	}
}

func (i *ParserRuleData) ReturnTypeToken() *lexer.Token {
	return i.returnType
}

func (i *ParserRuleData) WithReturnType(value *lexer.Token) {
	i.returnType = value
}

func (i *ParserRuleData) Body() Alternatives {
	if i != nil && i.body != nil {
		return i.body
	} else {
		return nil
	}
}

func (i *ParserRuleData) WithBody(value Alternatives) {
	i.body = value
}

type ParserRuleImpl struct {
	ast.AstNodeBase
	ParserRuleData
}

func (i *ParserRuleImpl) ForEachNode(fn func(ast.AstNode)) {
	i.ParserRuleData.ForEachNode(fn)
}

type Token interface {
	ast.AstNode

	IsToken()
	Type() string
	TypeToken() *lexer.Token
	WithType(value *lexer.Token)
	Name() string
	NameToken() *lexer.Token
	WithName(value *lexer.Token)
	Regexp() string
	RegexpToken() *lexer.Token
	WithRegexp(value *lexer.Token)
}

func NewToken() Token {
	return &TokenImpl{
		AstNodeBase: ast.NewAstNode(),
		TokenData:   NewTokenData(),
	}
}

type TokenData struct {
	_Type  *lexer.Token
	name   *lexer.Token
	regexp *lexer.Token
}

func NewTokenData() TokenData {
	return TokenData{}
}

func (i *TokenData) IsToken() {}

func (i *TokenData) ForEachNode(fn func(ast.AstNode)) {
}

func (i *TokenData) Type() string {
	if i != nil && i._Type != nil {
		return i._Type.Image
	} else {
		return ""
	}
}

func (i *TokenData) TypeToken() *lexer.Token {
	return i._Type
}

func (i *TokenData) WithType(value *lexer.Token) {
	i._Type = value
}

func (i *TokenData) Name() string {
	if i != nil && i.name != nil {
		return i.name.Image
	} else {
		return ""
	}
}

func (i *TokenData) NameToken() *lexer.Token {
	return i.name
}

func (i *TokenData) WithName(value *lexer.Token) {
	i.name = value
}

func (i *TokenData) Regexp() string {
	if i != nil && i.regexp != nil {
		return i.regexp.Image
	} else {
		return ""
	}
}

func (i *TokenData) RegexpToken() *lexer.Token {
	return i.regexp
}

func (i *TokenData) WithRegexp(value *lexer.Token) {
	i.regexp = value
}

type TokenImpl struct {
	ast.AstNodeBase
	TokenData
}

func (i *TokenImpl) ForEachNode(fn func(ast.AstNode)) {
	i.TokenData.ForEachNode(fn)
}

type Element interface {
	ast.AstNode

	IsElement()
	Cardinality() string
	CardinalityToken() *lexer.Token
	WithCardinality(value *lexer.Token)
}

func NewElement() Element {
	return &ElementImpl{
		AstNodeBase: ast.NewAstNode(),
		ElementData: NewElementData(),
	}
}

type ElementData struct {
	cardinality *lexer.Token
}

func NewElementData() ElementData {
	return ElementData{}
}

func (i *ElementData) IsElement() {}

func (i *ElementData) ForEachNode(fn func(ast.AstNode)) {
}

func (i *ElementData) Cardinality() string {
	if i != nil && i.cardinality != nil {
		return i.cardinality.Image
	} else {
		return ""
	}
}

func (i *ElementData) CardinalityToken() *lexer.Token {
	return i.cardinality
}

func (i *ElementData) WithCardinality(value *lexer.Token) {
	i.cardinality = value
}

type ElementImpl struct {
	ast.AstNodeBase
	ElementData
}

func (i *ElementImpl) ForEachNode(fn func(ast.AstNode)) {
	i.ElementData.ForEachNode(fn)
}

type Alternatives interface {
	ast.AstNode
	Element

	IsAlternatives()
	Alts() []Group
	WithAltsItem(item Group)
}

func NewAlternatives() Alternatives {
	return &AlternativesImpl{
		AstNodeBase:      ast.NewAstNode(),
		ElementData:      NewElementData(),
		AlternativesData: NewAlternativesData(),
	}
}

type AlternativesData struct {
	alts []Group
}

func NewAlternativesData() AlternativesData {
	return AlternativesData{
		alts: []Group{},
	}
}

func (i *AlternativesData) IsAlternatives() {}

func (i *AlternativesData) ForEachNode(fn func(ast.AstNode)) {
	for _, item := range i.alts {
		fn(item)
	}
}

func (i *AlternativesData) Alts() []Group {
	return i.alts
}

func (i *AlternativesData) WithAltsItem(item Group) {
	i.alts = append(i.alts, item)
}

type AlternativesImpl struct {
	ast.AstNodeBase
	ElementData
	AlternativesData
}

func (i *AlternativesImpl) ForEachNode(fn func(ast.AstNode)) {
	i.ElementData.ForEachNode(fn)
	i.AlternativesData.ForEachNode(fn)
}

type Group interface {
	ast.AstNode
	Element

	IsGroup()
	Elements() []Element
	WithElementsItem(item Element)
}

func NewGroup() Group {
	return &GroupImpl{
		AstNodeBase: ast.NewAstNode(),
		ElementData: NewElementData(),
		GroupData:   NewGroupData(),
	}
}

type GroupData struct {
	elements []Element
}

func NewGroupData() GroupData {
	return GroupData{
		elements: []Element{},
	}
}

func (i *GroupData) IsGroup() {}

func (i *GroupData) ForEachNode(fn func(ast.AstNode)) {
	for _, item := range i.elements {
		fn(item)
	}
}

func (i *GroupData) Elements() []Element {
	return i.elements
}

func (i *GroupData) WithElementsItem(item Element) {
	i.elements = append(i.elements, item)
}

type GroupImpl struct {
	ast.AstNodeBase
	ElementData
	GroupData
}

func (i *GroupImpl) ForEachNode(fn func(ast.AstNode)) {
	i.ElementData.ForEachNode(fn)
	i.GroupData.ForEachNode(fn)
}

type Keyword interface {
	ast.AstNode
	Element
	Assignable

	IsKeyword()
	Value() string
	ValueToken() *lexer.Token
	WithValue(value *lexer.Token)
}

func NewKeyword() Keyword {
	return &KeywordImpl{
		AstNodeBase:    ast.NewAstNode(),
		ElementData:    NewElementData(),
		AssignableData: NewAssignableData(),
		KeywordData:    NewKeywordData(),
	}
}

type KeywordData struct {
	value *lexer.Token
}

func NewKeywordData() KeywordData {
	return KeywordData{}
}

func (i *KeywordData) IsKeyword() {}

func (i *KeywordData) ForEachNode(fn func(ast.AstNode)) {
}

func (i *KeywordData) Value() string {
	if i != nil && i.value != nil {
		return i.value.Image
	} else {
		return ""
	}
}

func (i *KeywordData) ValueToken() *lexer.Token {
	return i.value
}

func (i *KeywordData) WithValue(value *lexer.Token) {
	i.value = value
}

type KeywordImpl struct {
	ast.AstNodeBase
	ElementData
	AssignableData
	KeywordData
}

func (i *KeywordImpl) ForEachNode(fn func(ast.AstNode)) {
	i.ElementData.ForEachNode(fn)
	i.AssignableData.ForEachNode(fn)
	i.KeywordData.ForEachNode(fn)
}

type Assignment interface {
	ast.AstNode
	Element

	IsAssignment()
	Property() string
	PropertyToken() *lexer.Token
	WithProperty(value *lexer.Token)
	Operator() string
	OperatorToken() *lexer.Token
	WithOperator(value *lexer.Token)
	Value() Assignable
	WithValue(value Assignable)
}

func NewAssignment() Assignment {
	return &AssignmentImpl{
		AstNodeBase:    ast.NewAstNode(),
		ElementData:    NewElementData(),
		AssignmentData: NewAssignmentData(),
	}
}

type AssignmentData struct {
	property *lexer.Token
	operator *lexer.Token
	value    Assignable
}

func NewAssignmentData() AssignmentData {
	return AssignmentData{}
}

func (i *AssignmentData) IsAssignment() {}

func (i *AssignmentData) ForEachNode(fn func(ast.AstNode)) {
	if i.value != nil {
		fn(i.value)
	}
}

func (i *AssignmentData) Property() string {
	if i != nil && i.property != nil {
		return i.property.Image
	} else {
		return ""
	}
}

func (i *AssignmentData) PropertyToken() *lexer.Token {
	return i.property
}

func (i *AssignmentData) WithProperty(value *lexer.Token) {
	i.property = value
}

func (i *AssignmentData) Operator() string {
	if i != nil && i.operator != nil {
		return i.operator.Image
	} else {
		return ""
	}
}

func (i *AssignmentData) OperatorToken() *lexer.Token {
	return i.operator
}

func (i *AssignmentData) WithOperator(value *lexer.Token) {
	i.operator = value
}

func (i *AssignmentData) Value() Assignable {
	if i != nil && i.value != nil {
		return i.value
	} else {
		return nil
	}
}

func (i *AssignmentData) WithValue(value Assignable) {
	i.value = value
}

type AssignmentImpl struct {
	ast.AstNodeBase
	ElementData
	AssignmentData
}

func (i *AssignmentImpl) ForEachNode(fn func(ast.AstNode)) {
	i.ElementData.ForEachNode(fn)
	i.AssignmentData.ForEachNode(fn)
}

type Assignable interface {
	ast.AstNode

	IsAssignable()
}

func NewAssignable() Assignable {
	return &AssignableImpl{
		AstNodeBase:    ast.NewAstNode(),
		AssignableData: NewAssignableData(),
	}
}

type AssignableData struct {
}

func NewAssignableData() AssignableData {
	return AssignableData{}
}

func (i *AssignableData) IsAssignable() {}

func (i *AssignableData) ForEachNode(fn func(ast.AstNode)) {
}

type AssignableImpl struct {
	ast.AstNodeBase
	AssignableData
}

func (i *AssignableImpl) ForEachNode(fn func(ast.AstNode)) {
	i.AssignableData.ForEachNode(fn)
}

type CrossRef interface {
	ast.AstNode
	Assignable

	IsCrossRef()
	Type() string
	TypeToken() *lexer.Token
	WithType(value *lexer.Token)
	Rule() RuleCall
	WithRule(value RuleCall)
}

func NewCrossRef() CrossRef {
	return &CrossRefImpl{
		AstNodeBase:    ast.NewAstNode(),
		AssignableData: NewAssignableData(),
		CrossRefData:   NewCrossRefData(),
	}
}

type CrossRefData struct {
	_Type *lexer.Token
	rule  RuleCall
}

func NewCrossRefData() CrossRefData {
	return CrossRefData{}
}

func (i *CrossRefData) IsCrossRef() {}

func (i *CrossRefData) ForEachNode(fn func(ast.AstNode)) {
	if i.rule != nil {
		fn(i.rule)
	}
}

func (i *CrossRefData) Type() string {
	if i != nil && i._Type != nil {
		return i._Type.Image
	} else {
		return ""
	}
}

func (i *CrossRefData) TypeToken() *lexer.Token {
	return i._Type
}

func (i *CrossRefData) WithType(value *lexer.Token) {
	i._Type = value
}

func (i *CrossRefData) Rule() RuleCall {
	if i != nil && i.rule != nil {
		return i.rule
	} else {
		return nil
	}
}

func (i *CrossRefData) WithRule(value RuleCall) {
	i.rule = value
}

type CrossRefImpl struct {
	ast.AstNodeBase
	AssignableData
	CrossRefData
}

func (i *CrossRefImpl) ForEachNode(fn func(ast.AstNode)) {
	i.AssignableData.ForEachNode(fn)
	i.CrossRefData.ForEachNode(fn)
}

type RuleCall interface {
	ast.AstNode
	Element
	Assignable

	IsRuleCall()
	Rule() string
	RuleToken() *lexer.Token
	WithRule(value *lexer.Token)
}

func NewRuleCall() RuleCall {
	return &RuleCallImpl{
		AstNodeBase:    ast.NewAstNode(),
		ElementData:    NewElementData(),
		AssignableData: NewAssignableData(),
		RuleCallData:   NewRuleCallData(),
	}
}

type RuleCallData struct {
	rule *lexer.Token
}

func NewRuleCallData() RuleCallData {
	return RuleCallData{}
}

func (i *RuleCallData) IsRuleCall() {}

func (i *RuleCallData) ForEachNode(fn func(ast.AstNode)) {
}

func (i *RuleCallData) Rule() string {
	if i != nil && i.rule != nil {
		return i.rule.Image
	} else {
		return ""
	}
}

func (i *RuleCallData) RuleToken() *lexer.Token {
	return i.rule
}

func (i *RuleCallData) WithRule(value *lexer.Token) {
	i.rule = value
}

type RuleCallImpl struct {
	ast.AstNodeBase
	ElementData
	AssignableData
	RuleCallData
}

func (i *RuleCallImpl) ForEachNode(fn func(ast.AstNode)) {
	i.ElementData.ForEachNode(fn)
	i.AssignableData.ForEachNode(fn)
	i.RuleCallData.ForEachNode(fn)
}

type Action interface {
	ast.AstNode
	Element

	IsAction()
	Type() string
	TypeToken() *lexer.Token
	WithType(value *lexer.Token)
	Property() string
	PropertyToken() *lexer.Token
	WithProperty(value *lexer.Token)
}

func NewAction() Action {
	return &ActionImpl{
		AstNodeBase: ast.NewAstNode(),
		ElementData: NewElementData(),
		ActionData:  NewActionData(),
	}
}

type ActionData struct {
	_Type    *lexer.Token
	property *lexer.Token
}

func NewActionData() ActionData {
	return ActionData{}
}

func (i *ActionData) IsAction() {}

func (i *ActionData) ForEachNode(fn func(ast.AstNode)) {
}

func (i *ActionData) Type() string {
	if i != nil && i._Type != nil {
		return i._Type.Image
	} else {
		return ""
	}
}

func (i *ActionData) TypeToken() *lexer.Token {
	return i._Type
}

func (i *ActionData) WithType(value *lexer.Token) {
	i._Type = value
}

func (i *ActionData) Property() string {
	if i != nil && i.property != nil {
		return i.property.Image
	} else {
		return ""
	}
}

func (i *ActionData) PropertyToken() *lexer.Token {
	return i.property
}

func (i *ActionData) WithProperty(value *lexer.Token) {
	i.property = value
}

type ActionImpl struct {
	ast.AstNodeBase
	ElementData
	ActionData
}

func (i *ActionImpl) ForEachNode(fn func(ast.AstNode)) {
	i.ElementData.ForEachNode(fn)
	i.ActionData.ForEachNode(fn)
}
