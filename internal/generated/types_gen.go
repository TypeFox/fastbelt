package generated

import (
	"github.com/TypeFox/langium-to-go/core"
)

type Grammar interface {
	core.AstNode

	IsGrammar()
	Name() string
	NameToken() *core.Token
	WithName(value *core.Token)
	Rules() []ParserRule
	WithRulesItem(item ParserRule)
	Terminals() []Token
	WithTerminalsItem(item Token)
	Interfaces() []Interface
	WithInterfacesItem(item Interface)
}

func NewGrammar() Grammar {
	return &GrammarImpl{
		AstNodeBase: core.NewAstNode(),
		GrammarData: NewGrammarData(),
	}
}

type GrammarData struct {
	name       *core.Token
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

func (i *GrammarData) ForEachNode(fn func(core.AstNode)) {
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

func (i *GrammarData) NameToken() *core.Token {
	return i.name
}

func (i *GrammarData) WithName(value *core.Token) {
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
	core.AstNodeBase
	GrammarData
}

func (i *GrammarImpl) ForEachNode(fn func(core.AstNode)) {
	i.GrammarData.ForEachNode(fn)
}

type Interface interface {
	core.AstNode

	IsInterface()
	Name() string
	NameToken() *core.Token
	WithName(value *core.Token)
	Extends() []*core.Token
	WithExtendsItem(item *core.Token)
	Fields() []Field
	WithFieldsItem(item Field)
}

func NewInterface() Interface {
	return &InterfaceImpl{
		AstNodeBase:   core.NewAstNode(),
		InterfaceData: NewInterfaceData(),
	}
}

type InterfaceData struct {
	name    *core.Token
	extends []*core.Token
	fields  []Field
}

func NewInterfaceData() InterfaceData {
	return InterfaceData{
		extends: []*core.Token{},
		fields:  []Field{},
	}
}

func (i *InterfaceData) IsInterface() {}

func (i *InterfaceData) ForEachNode(fn func(core.AstNode)) {
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

func (i *InterfaceData) NameToken() *core.Token {
	return i.name
}

func (i *InterfaceData) WithName(value *core.Token) {
	i.name = value
}

func (i *InterfaceData) Extends() []*core.Token {
	return i.extends
}

func (i *InterfaceData) WithExtendsItem(item *core.Token) {
	i.extends = append(i.extends, item)
}

func (i *InterfaceData) Fields() []Field {
	return i.fields
}

func (i *InterfaceData) WithFieldsItem(item Field) {
	i.fields = append(i.fields, item)
}

type InterfaceImpl struct {
	core.AstNodeBase
	InterfaceData
}

func (i *InterfaceImpl) ForEachNode(fn func(core.AstNode)) {
	i.InterfaceData.ForEachNode(fn)
}

type Field interface {
	core.AstNode

	IsField()
	Name() string
	NameToken() *core.Token
	WithName(value *core.Token)
	IsArray() bool
	ArrayToken() *core.Token
	WithArray(value *core.Token)
	Type() string
	TypeToken() *core.Token
	WithType(value *core.Token)
}

func NewField() Field {
	return &FieldImpl{
		AstNodeBase: core.NewAstNode(),
		FieldData:   NewFieldData(),
	}
}

type FieldData struct {
	name  *core.Token
	array *core.Token
	_Type *core.Token
}

func NewFieldData() FieldData {
	return FieldData{}
}

func (i *FieldData) IsField() {}

func (i *FieldData) ForEachNode(fn func(core.AstNode)) {
}

func (i *FieldData) Name() string {
	if i != nil && i.name != nil {
		return i.name.Image
	} else {
		return ""
	}
}

func (i *FieldData) NameToken() *core.Token {
	return i.name
}

func (i *FieldData) WithName(value *core.Token) {
	i.name = value
}

func (i *FieldData) IsArray() bool {
	return i != nil && i.array != nil
}

func (i *FieldData) ArrayToken() *core.Token {
	return i.array
}

func (i *FieldData) WithArray(value *core.Token) {
	i.array = value
}

func (i *FieldData) Type() string {
	if i != nil && i._Type != nil {
		return i._Type.Image
	} else {
		return ""
	}
}

func (i *FieldData) TypeToken() *core.Token {
	return i._Type
}

func (i *FieldData) WithType(value *core.Token) {
	i._Type = value
}

type FieldImpl struct {
	core.AstNodeBase
	FieldData
}

func (i *FieldImpl) ForEachNode(fn func(core.AstNode)) {
	i.FieldData.ForEachNode(fn)
}

type ParserRule interface {
	core.AstNode

	IsParserRule()
	Name() string
	NameToken() *core.Token
	WithName(value *core.Token)
	ReturnType() string
	ReturnTypeToken() *core.Token
	WithReturnType(value *core.Token)
	Body() Alternatives
	WithBody(value Alternatives)
}

func NewParserRule() ParserRule {
	return &ParserRuleImpl{
		AstNodeBase:    core.NewAstNode(),
		ParserRuleData: NewParserRuleData(),
	}
}

type ParserRuleData struct {
	name       *core.Token
	returnType *core.Token
	body       Alternatives
}

func NewParserRuleData() ParserRuleData {
	return ParserRuleData{}
}

func (i *ParserRuleData) IsParserRule() {}

func (i *ParserRuleData) ForEachNode(fn func(core.AstNode)) {
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

func (i *ParserRuleData) NameToken() *core.Token {
	return i.name
}

func (i *ParserRuleData) WithName(value *core.Token) {
	i.name = value
}

func (i *ParserRuleData) ReturnType() string {
	if i != nil && i.returnType != nil {
		return i.returnType.Image
	} else {
		return ""
	}
}

func (i *ParserRuleData) ReturnTypeToken() *core.Token {
	return i.returnType
}

func (i *ParserRuleData) WithReturnType(value *core.Token) {
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
	core.AstNodeBase
	ParserRuleData
}

func (i *ParserRuleImpl) ForEachNode(fn func(core.AstNode)) {
	i.ParserRuleData.ForEachNode(fn)
}

type Token interface {
	core.AstNode

	IsToken()
	Type() string
	TypeToken() *core.Token
	WithType(value *core.Token)
	Name() string
	NameToken() *core.Token
	WithName(value *core.Token)
	Regexp() string
	RegexpToken() *core.Token
	WithRegexp(value *core.Token)
}

func NewToken() Token {
	return &TokenImpl{
		AstNodeBase: core.NewAstNode(),
		TokenData:   NewTokenData(),
	}
}

type TokenData struct {
	_Type  *core.Token
	name   *core.Token
	regexp *core.Token
}

func NewTokenData() TokenData {
	return TokenData{}
}

func (i *TokenData) IsToken() {}

func (i *TokenData) ForEachNode(fn func(core.AstNode)) {
}

func (i *TokenData) Type() string {
	if i != nil && i._Type != nil {
		return i._Type.Image
	} else {
		return ""
	}
}

func (i *TokenData) TypeToken() *core.Token {
	return i._Type
}

func (i *TokenData) WithType(value *core.Token) {
	i._Type = value
}

func (i *TokenData) Name() string {
	if i != nil && i.name != nil {
		return i.name.Image
	} else {
		return ""
	}
}

func (i *TokenData) NameToken() *core.Token {
	return i.name
}

func (i *TokenData) WithName(value *core.Token) {
	i.name = value
}

func (i *TokenData) Regexp() string {
	if i != nil && i.regexp != nil {
		return i.regexp.Image
	} else {
		return ""
	}
}

func (i *TokenData) RegexpToken() *core.Token {
	return i.regexp
}

func (i *TokenData) WithRegexp(value *core.Token) {
	i.regexp = value
}

type TokenImpl struct {
	core.AstNodeBase
	TokenData
}

func (i *TokenImpl) ForEachNode(fn func(core.AstNode)) {
	i.TokenData.ForEachNode(fn)
}

type Element interface {
	core.AstNode

	IsElement()
	Cardinality() string
	CardinalityToken() *core.Token
	WithCardinality(value *core.Token)
}

func NewElement() Element {
	return &ElementImpl{
		AstNodeBase: core.NewAstNode(),
		ElementData: NewElementData(),
	}
}

type ElementData struct {
	cardinality *core.Token
}

func NewElementData() ElementData {
	return ElementData{}
}

func (i *ElementData) IsElement() {}

func (i *ElementData) ForEachNode(fn func(core.AstNode)) {
}

func (i *ElementData) Cardinality() string {
	if i != nil && i.cardinality != nil {
		return i.cardinality.Image
	} else {
		return ""
	}
}

func (i *ElementData) CardinalityToken() *core.Token {
	return i.cardinality
}

func (i *ElementData) WithCardinality(value *core.Token) {
	i.cardinality = value
}

type ElementImpl struct {
	core.AstNodeBase
	ElementData
}

func (i *ElementImpl) ForEachNode(fn func(core.AstNode)) {
	i.ElementData.ForEachNode(fn)
}

type Alternatives interface {
	core.AstNode
	Element

	IsAlternatives()
	Alts() []Group
	WithAltsItem(item Group)
}

func NewAlternatives() Alternatives {
	return &AlternativesImpl{
		AstNodeBase:      core.NewAstNode(),
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

func (i *AlternativesData) ForEachNode(fn func(core.AstNode)) {
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
	core.AstNodeBase
	ElementData
	AlternativesData
}

func (i *AlternativesImpl) ForEachNode(fn func(core.AstNode)) {
	i.ElementData.ForEachNode(fn)
	i.AlternativesData.ForEachNode(fn)
}

type Group interface {
	core.AstNode
	Element

	IsGroup()
	Elements() []Element
	WithElementsItem(item Element)
}

func NewGroup() Group {
	return &GroupImpl{
		AstNodeBase: core.NewAstNode(),
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

func (i *GroupData) ForEachNode(fn func(core.AstNode)) {
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
	core.AstNodeBase
	ElementData
	GroupData
}

func (i *GroupImpl) ForEachNode(fn func(core.AstNode)) {
	i.ElementData.ForEachNode(fn)
	i.GroupData.ForEachNode(fn)
}

type Keyword interface {
	core.AstNode
	Element
	Assignable

	IsKeyword()
	Value() string
	ValueToken() *core.Token
	WithValue(value *core.Token)
}

func NewKeyword() Keyword {
	return &KeywordImpl{
		AstNodeBase:    core.NewAstNode(),
		ElementData:    NewElementData(),
		AssignableData: NewAssignableData(),
		KeywordData:    NewKeywordData(),
	}
}

type KeywordData struct {
	value *core.Token
}

func NewKeywordData() KeywordData {
	return KeywordData{}
}

func (i *KeywordData) IsKeyword() {}

func (i *KeywordData) ForEachNode(fn func(core.AstNode)) {
}

func (i *KeywordData) Value() string {
	if i != nil && i.value != nil {
		return i.value.Image
	} else {
		return ""
	}
}

func (i *KeywordData) ValueToken() *core.Token {
	return i.value
}

func (i *KeywordData) WithValue(value *core.Token) {
	i.value = value
}

type KeywordImpl struct {
	core.AstNodeBase
	ElementData
	AssignableData
	KeywordData
}

func (i *KeywordImpl) ForEachNode(fn func(core.AstNode)) {
	i.ElementData.ForEachNode(fn)
	i.AssignableData.ForEachNode(fn)
	i.KeywordData.ForEachNode(fn)
}

type Assignment interface {
	core.AstNode
	Element

	IsAssignment()
	Property() string
	PropertyToken() *core.Token
	WithProperty(value *core.Token)
	Operator() string
	OperatorToken() *core.Token
	WithOperator(value *core.Token)
	Value() Assignable
	WithValue(value Assignable)
}

func NewAssignment() Assignment {
	return &AssignmentImpl{
		AstNodeBase:    core.NewAstNode(),
		ElementData:    NewElementData(),
		AssignmentData: NewAssignmentData(),
	}
}

type AssignmentData struct {
	property *core.Token
	operator *core.Token
	value    Assignable
}

func NewAssignmentData() AssignmentData {
	return AssignmentData{}
}

func (i *AssignmentData) IsAssignment() {}

func (i *AssignmentData) ForEachNode(fn func(core.AstNode)) {
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

func (i *AssignmentData) PropertyToken() *core.Token {
	return i.property
}

func (i *AssignmentData) WithProperty(value *core.Token) {
	i.property = value
}

func (i *AssignmentData) Operator() string {
	if i != nil && i.operator != nil {
		return i.operator.Image
	} else {
		return ""
	}
}

func (i *AssignmentData) OperatorToken() *core.Token {
	return i.operator
}

func (i *AssignmentData) WithOperator(value *core.Token) {
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
	core.AstNodeBase
	ElementData
	AssignmentData
}

func (i *AssignmentImpl) ForEachNode(fn func(core.AstNode)) {
	i.ElementData.ForEachNode(fn)
	i.AssignmentData.ForEachNode(fn)
}

type Assignable interface {
	core.AstNode

	IsAssignable()
}

func NewAssignable() Assignable {
	return &AssignableImpl{
		AstNodeBase:    core.NewAstNode(),
		AssignableData: NewAssignableData(),
	}
}

type AssignableData struct {
}

func NewAssignableData() AssignableData {
	return AssignableData{}
}

func (i *AssignableData) IsAssignable() {}

func (i *AssignableData) ForEachNode(fn func(core.AstNode)) {
}

type AssignableImpl struct {
	core.AstNodeBase
	AssignableData
}

func (i *AssignableImpl) ForEachNode(fn func(core.AstNode)) {
	i.AssignableData.ForEachNode(fn)
}

type CrossRef interface {
	core.AstNode
	Assignable

	IsCrossRef()
	Type() string
	TypeToken() *core.Token
	WithType(value *core.Token)
	Rule() RuleCall
	WithRule(value RuleCall)
}

func NewCrossRef() CrossRef {
	return &CrossRefImpl{
		AstNodeBase:    core.NewAstNode(),
		AssignableData: NewAssignableData(),
		CrossRefData:   NewCrossRefData(),
	}
}

type CrossRefData struct {
	_Type *core.Token
	rule  RuleCall
}

func NewCrossRefData() CrossRefData {
	return CrossRefData{}
}

func (i *CrossRefData) IsCrossRef() {}

func (i *CrossRefData) ForEachNode(fn func(core.AstNode)) {
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

func (i *CrossRefData) TypeToken() *core.Token {
	return i._Type
}

func (i *CrossRefData) WithType(value *core.Token) {
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
	core.AstNodeBase
	AssignableData
	CrossRefData
}

func (i *CrossRefImpl) ForEachNode(fn func(core.AstNode)) {
	i.AssignableData.ForEachNode(fn)
	i.CrossRefData.ForEachNode(fn)
}

type RuleCall interface {
	core.AstNode
	Element
	Assignable

	IsRuleCall()
	Rule() string
	RuleToken() *core.Token
	WithRule(value *core.Token)
}

func NewRuleCall() RuleCall {
	return &RuleCallImpl{
		AstNodeBase:    core.NewAstNode(),
		ElementData:    NewElementData(),
		AssignableData: NewAssignableData(),
		RuleCallData:   NewRuleCallData(),
	}
}

type RuleCallData struct {
	rule *core.Token
}

func NewRuleCallData() RuleCallData {
	return RuleCallData{}
}

func (i *RuleCallData) IsRuleCall() {}

func (i *RuleCallData) ForEachNode(fn func(core.AstNode)) {
}

func (i *RuleCallData) Rule() string {
	if i != nil && i.rule != nil {
		return i.rule.Image
	} else {
		return ""
	}
}

func (i *RuleCallData) RuleToken() *core.Token {
	return i.rule
}

func (i *RuleCallData) WithRule(value *core.Token) {
	i.rule = value
}

type RuleCallImpl struct {
	core.AstNodeBase
	ElementData
	AssignableData
	RuleCallData
}

func (i *RuleCallImpl) ForEachNode(fn func(core.AstNode)) {
	i.ElementData.ForEachNode(fn)
	i.AssignableData.ForEachNode(fn)
	i.RuleCallData.ForEachNode(fn)
}

type Action interface {
	core.AstNode
	Element

	IsAction()
	Type() string
	TypeToken() *core.Token
	WithType(value *core.Token)
	Property() string
	PropertyToken() *core.Token
	WithProperty(value *core.Token)
}

func NewAction() Action {
	return &ActionImpl{
		AstNodeBase: core.NewAstNode(),
		ElementData: NewElementData(),
		ActionData:  NewActionData(),
	}
}

type ActionData struct {
	_Type    *core.Token
	property *core.Token
}

func NewActionData() ActionData {
	return ActionData{}
}

func (i *ActionData) IsAction() {}

func (i *ActionData) ForEachNode(fn func(core.AstNode)) {
}

func (i *ActionData) Type() string {
	if i != nil && i._Type != nil {
		return i._Type.Image
	} else {
		return ""
	}
}

func (i *ActionData) TypeToken() *core.Token {
	return i._Type
}

func (i *ActionData) WithType(value *core.Token) {
	i._Type = value
}

func (i *ActionData) Property() string {
	if i != nil && i.property != nil {
		return i.property.Image
	} else {
		return ""
	}
}

func (i *ActionData) PropertyToken() *core.Token {
	return i.property
}

func (i *ActionData) WithProperty(value *core.Token) {
	i.property = value
}

type ActionImpl struct {
	core.AstNodeBase
	ElementData
	ActionData
}

func (i *ActionImpl) ForEachNode(fn func(core.AstNode)) {
	i.ElementData.ForEachNode(fn)
	i.ActionData.ForEachNode(fn)
}
