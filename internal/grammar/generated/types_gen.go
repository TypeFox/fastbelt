package generated

import (
	core "typefox.dev/fastbelt"
)

type Grammar interface {
	core.AstNode

	IsGrammar()
	Name() string
	NameToken() *core.Token
	SetName(value *core.Token)
	Rules() []ParserRule
	SetRulesItem(item ParserRule)
	Terminals() []Token
	SetTerminalsItem(item Token)
	Interfaces() []Interface
	SetInterfacesItem(item Interface)
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

func (i *GrammarData) ForEachReference(fn func(core.UntypedReference)) {
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

func (i *GrammarData) SetName(value *core.Token) {
	i.name = value
}

func (i *GrammarData) Rules() []ParserRule {
	return i.rules
}

func (i *GrammarData) SetRulesItem(item ParserRule) {
	i.rules = append(i.rules, item)
}

func (i *GrammarData) Terminals() []Token {
	return i.terminals
}

func (i *GrammarData) SetTerminalsItem(item Token) {
	i.terminals = append(i.terminals, item)
}

func (i *GrammarData) Interfaces() []Interface {
	return i.interfaces
}

func (i *GrammarData) SetInterfacesItem(item Interface) {
	i.interfaces = append(i.interfaces, item)
}

type GrammarImpl struct {
	core.AstNodeBase
	GrammarData
}

func (i *GrammarImpl) ForEachNode(fn func(core.AstNode)) {
	i.GrammarData.ForEachNode(fn)
}

func (i *GrammarImpl) ForEachReference(fn func(core.UntypedReference)) {
	i.GrammarData.ForEachReference(fn)
}

type Interface interface {
	core.AstNode

	IsInterface()
	Name() string
	NameToken() *core.Token
	SetName(value *core.Token)
	Extends() []*core.Token
	SetExtendsItem(item *core.Token)
	Fields() []Field
	SetFieldsItem(item Field)
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

func (i *InterfaceData) ForEachReference(fn func(core.UntypedReference)) {
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

func (i *InterfaceData) SetName(value *core.Token) {
	i.name = value
}

func (i *InterfaceData) Extends() []*core.Token {
	return i.extends
}

func (i *InterfaceData) SetExtendsItem(item *core.Token) {
	i.extends = append(i.extends, item)
}

func (i *InterfaceData) Fields() []Field {
	return i.fields
}

func (i *InterfaceData) SetFieldsItem(item Field) {
	i.fields = append(i.fields, item)
}

type InterfaceImpl struct {
	core.AstNodeBase
	InterfaceData
}

func (i *InterfaceImpl) ForEachNode(fn func(core.AstNode)) {
	i.InterfaceData.ForEachNode(fn)
}

func (i *InterfaceImpl) ForEachReference(fn func(core.UntypedReference)) {
	i.InterfaceData.ForEachReference(fn)
}

type Field interface {
	core.AstNode

	IsField()
	Name() string
	NameToken() *core.Token
	SetName(value *core.Token)
	Type() FieldType
	SetType(value FieldType)
}

func NewField() Field {
	return &FieldImpl{
		AstNodeBase: core.NewAstNode(),
		FieldData:   NewFieldData(),
	}
}

type FieldData struct {
	name  *core.Token
	_Type FieldType
}

func NewFieldData() FieldData {
	return FieldData{}
}

func (i *FieldData) IsField() {}

func (i *FieldData) ForEachNode(fn func(core.AstNode)) {
	if i._Type != nil {
		fn(i._Type)
	}
}

func (i *FieldData) ForEachReference(fn func(core.UntypedReference)) {
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

func (i *FieldData) SetName(value *core.Token) {
	i.name = value
}

func (i *FieldData) Type() FieldType {
	if i != nil && i._Type != nil {
		return i._Type
	} else {
		return nil
	}
}

func (i *FieldData) SetType(value FieldType) {
	i._Type = value
}

type FieldImpl struct {
	core.AstNodeBase
	FieldData
}

func (i *FieldImpl) ForEachNode(fn func(core.AstNode)) {
	i.FieldData.ForEachNode(fn)
}

func (i *FieldImpl) ForEachReference(fn func(core.UntypedReference)) {
	i.FieldData.ForEachReference(fn)
}

type FieldType interface {
	core.AstNode

	IsFieldType()
}

func NewFieldType() FieldType {
	return &FieldTypeImpl{
		AstNodeBase:   core.NewAstNode(),
		FieldTypeData: NewFieldTypeData(),
	}
}

type FieldTypeData struct {
}

func NewFieldTypeData() FieldTypeData {
	return FieldTypeData{}
}

func (i *FieldTypeData) IsFieldType() {}

func (i *FieldTypeData) ForEachNode(fn func(core.AstNode)) {
}

func (i *FieldTypeData) ForEachReference(fn func(core.UntypedReference)) {
}

type FieldTypeImpl struct {
	core.AstNodeBase
	FieldTypeData
}

func (i *FieldTypeImpl) ForEachNode(fn func(core.AstNode)) {
	i.FieldTypeData.ForEachNode(fn)
}

func (i *FieldTypeImpl) ForEachReference(fn func(core.UntypedReference)) {
	i.FieldTypeData.ForEachReference(fn)
}

type ArrayType interface {
	core.AstNode
	FieldType

	IsArrayType()
	InternalType() FieldType
	SetInternalType(value FieldType)
}

func NewArrayType() ArrayType {
	return &ArrayTypeImpl{
		AstNodeBase:   core.NewAstNode(),
		FieldTypeData: NewFieldTypeData(),
		ArrayTypeData: NewArrayTypeData(),
	}
}

type ArrayTypeData struct {
	internalType FieldType
}

func NewArrayTypeData() ArrayTypeData {
	return ArrayTypeData{}
}

func (i *ArrayTypeData) IsArrayType() {}

func (i *ArrayTypeData) ForEachNode(fn func(core.AstNode)) {
	if i.internalType != nil {
		fn(i.internalType)
	}
}

func (i *ArrayTypeData) ForEachReference(fn func(core.UntypedReference)) {
}

func (i *ArrayTypeData) InternalType() FieldType {
	if i != nil && i.internalType != nil {
		return i.internalType
	} else {
		return nil
	}
}

func (i *ArrayTypeData) SetInternalType(value FieldType) {
	i.internalType = value
}

type ArrayTypeImpl struct {
	core.AstNodeBase
	FieldTypeData
	ArrayTypeData
}

func (i *ArrayTypeImpl) ForEachNode(fn func(core.AstNode)) {
	i.FieldTypeData.ForEachNode(fn)
	i.ArrayTypeData.ForEachNode(fn)
}

func (i *ArrayTypeImpl) ForEachReference(fn func(core.UntypedReference)) {
	i.FieldTypeData.ForEachReference(fn)
	i.ArrayTypeData.ForEachReference(fn)
}

type ReferenceType interface {
	core.AstNode
	FieldType

	IsReferenceType()
	Type() string
	TypeToken() *core.Token
	SetType(value *core.Token)
}

func NewReferenceType() ReferenceType {
	return &ReferenceTypeImpl{
		AstNodeBase:       core.NewAstNode(),
		FieldTypeData:     NewFieldTypeData(),
		ReferenceTypeData: NewReferenceTypeData(),
	}
}

type ReferenceTypeData struct {
	_Type *core.Token
}

func NewReferenceTypeData() ReferenceTypeData {
	return ReferenceTypeData{}
}

func (i *ReferenceTypeData) IsReferenceType() {}

func (i *ReferenceTypeData) ForEachNode(fn func(core.AstNode)) {
}

func (i *ReferenceTypeData) ForEachReference(fn func(core.UntypedReference)) {
}

func (i *ReferenceTypeData) Type() string {
	if i != nil && i._Type != nil {
		return i._Type.Image
	} else {
		return ""
	}
}

func (i *ReferenceTypeData) TypeToken() *core.Token {
	return i._Type
}

func (i *ReferenceTypeData) SetType(value *core.Token) {
	i._Type = value
}

type ReferenceTypeImpl struct {
	core.AstNodeBase
	FieldTypeData
	ReferenceTypeData
}

func (i *ReferenceTypeImpl) ForEachNode(fn func(core.AstNode)) {
	i.FieldTypeData.ForEachNode(fn)
	i.ReferenceTypeData.ForEachNode(fn)
}

func (i *ReferenceTypeImpl) ForEachReference(fn func(core.UntypedReference)) {
	i.FieldTypeData.ForEachReference(fn)
	i.ReferenceTypeData.ForEachReference(fn)
}

type SimpleType interface {
	core.AstNode
	FieldType

	IsSimpleType()
	Type() string
	TypeToken() *core.Token
	SetType(value *core.Token)
}

func NewSimpleType() SimpleType {
	return &SimpleTypeImpl{
		AstNodeBase:    core.NewAstNode(),
		FieldTypeData:  NewFieldTypeData(),
		SimpleTypeData: NewSimpleTypeData(),
	}
}

type SimpleTypeData struct {
	_Type *core.Token
}

func NewSimpleTypeData() SimpleTypeData {
	return SimpleTypeData{}
}

func (i *SimpleTypeData) IsSimpleType() {}

func (i *SimpleTypeData) ForEachNode(fn func(core.AstNode)) {
}

func (i *SimpleTypeData) ForEachReference(fn func(core.UntypedReference)) {
}

func (i *SimpleTypeData) Type() string {
	if i != nil && i._Type != nil {
		return i._Type.Image
	} else {
		return ""
	}
}

func (i *SimpleTypeData) TypeToken() *core.Token {
	return i._Type
}

func (i *SimpleTypeData) SetType(value *core.Token) {
	i._Type = value
}

type SimpleTypeImpl struct {
	core.AstNodeBase
	FieldTypeData
	SimpleTypeData
}

func (i *SimpleTypeImpl) ForEachNode(fn func(core.AstNode)) {
	i.FieldTypeData.ForEachNode(fn)
	i.SimpleTypeData.ForEachNode(fn)
}

func (i *SimpleTypeImpl) ForEachReference(fn func(core.UntypedReference)) {
	i.FieldTypeData.ForEachReference(fn)
	i.SimpleTypeData.ForEachReference(fn)
}

type ParserRule interface {
	core.AstNode

	IsParserRule()
	Name() string
	NameToken() *core.Token
	SetName(value *core.Token)
	ReturnType() *core.Reference[Interface]
	SetReturnType(value *core.Reference[Interface])
	Body() Element
	SetBody(value Element)
}

func NewParserRule() ParserRule {
	return &ParserRuleImpl{
		AstNodeBase:    core.NewAstNode(),
		ParserRuleData: NewParserRuleData(),
	}
}

type ParserRuleData struct {
	name       *core.Token
	returnType *core.Reference[Interface]
	body       Element
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

func (i *ParserRuleData) ForEachReference(fn func(core.UntypedReference)) {
	if i.returnType != nil {
		fn(i.returnType)
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

func (i *ParserRuleData) SetName(value *core.Token) {
	i.name = value
}

func (i *ParserRuleData) ReturnType() *core.Reference[Interface] {
	if i != nil && i.returnType != nil {
		return i.returnType
	} else {
		return nil
	}
}

func (i *ParserRuleData) SetReturnType(value *core.Reference[Interface]) {
	i.returnType = value
}

func (i *ParserRuleData) Body() Element {
	if i != nil && i.body != nil {
		return i.body
	} else {
		return nil
	}
}

func (i *ParserRuleData) SetBody(value Element) {
	i.body = value
}

type ParserRuleImpl struct {
	core.AstNodeBase
	ParserRuleData
}

func (i *ParserRuleImpl) ForEachNode(fn func(core.AstNode)) {
	i.ParserRuleData.ForEachNode(fn)
}

func (i *ParserRuleImpl) ForEachReference(fn func(core.UntypedReference)) {
	i.ParserRuleData.ForEachReference(fn)
}

type Token interface {
	core.AstNode

	IsToken()
	Type() string
	TypeToken() *core.Token
	SetType(value *core.Token)
	Name() string
	NameToken() *core.Token
	SetName(value *core.Token)
	Regexp() string
	RegexpToken() *core.Token
	SetRegexp(value *core.Token)
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

func (i *TokenData) ForEachReference(fn func(core.UntypedReference)) {
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

func (i *TokenData) SetType(value *core.Token) {
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

func (i *TokenData) SetName(value *core.Token) {
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

func (i *TokenData) SetRegexp(value *core.Token) {
	i.regexp = value
}

type TokenImpl struct {
	core.AstNodeBase
	TokenData
}

func (i *TokenImpl) ForEachNode(fn func(core.AstNode)) {
	i.TokenData.ForEachNode(fn)
}

func (i *TokenImpl) ForEachReference(fn func(core.UntypedReference)) {
	i.TokenData.ForEachReference(fn)
}

type Element interface {
	core.AstNode

	IsElement()
	Cardinality() string
	CardinalityToken() *core.Token
	SetCardinality(value *core.Token)
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

func (i *ElementData) ForEachReference(fn func(core.UntypedReference)) {
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

func (i *ElementData) SetCardinality(value *core.Token) {
	i.cardinality = value
}

type ElementImpl struct {
	core.AstNodeBase
	ElementData
}

func (i *ElementImpl) ForEachNode(fn func(core.AstNode)) {
	i.ElementData.ForEachNode(fn)
}

func (i *ElementImpl) ForEachReference(fn func(core.UntypedReference)) {
	i.ElementData.ForEachReference(fn)
}

type Alternatives interface {
	core.AstNode
	Assignable

	IsAlternatives()
	Alts() []Element
	SetAltsItem(item Element)
}

func NewAlternatives() Alternatives {
	return &AlternativesImpl{
		AstNodeBase:      core.NewAstNode(),
		AssignableData:   NewAssignableData(),
		ElementData:      NewElementData(),
		AlternativesData: NewAlternativesData(),
	}
}

type AlternativesData struct {
	alts []Element
}

func NewAlternativesData() AlternativesData {
	return AlternativesData{
		alts: []Element{},
	}
}

func (i *AlternativesData) IsAlternatives() {}

func (i *AlternativesData) ForEachNode(fn func(core.AstNode)) {
	for _, item := range i.alts {
		fn(item)
	}
}

func (i *AlternativesData) ForEachReference(fn func(core.UntypedReference)) {
}

func (i *AlternativesData) Alts() []Element {
	return i.alts
}

func (i *AlternativesData) SetAltsItem(item Element) {
	i.alts = append(i.alts, item)
}

type AlternativesImpl struct {
	core.AstNodeBase
	AssignableData
	ElementData
	AlternativesData
}

func (i *AlternativesImpl) ForEachNode(fn func(core.AstNode)) {
	i.AssignableData.ForEachNode(fn)
	i.ElementData.ForEachNode(fn)
	i.AlternativesData.ForEachNode(fn)
}

func (i *AlternativesImpl) ForEachReference(fn func(core.UntypedReference)) {
	i.AssignableData.ForEachReference(fn)
	i.ElementData.ForEachReference(fn)
	i.AlternativesData.ForEachReference(fn)
}

type Group interface {
	core.AstNode
	Element

	IsGroup()
	Elements() []Element
	SetElementsItem(item Element)
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

func (i *GroupData) ForEachReference(fn func(core.UntypedReference)) {
}

func (i *GroupData) Elements() []Element {
	return i.elements
}

func (i *GroupData) SetElementsItem(item Element) {
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

func (i *GroupImpl) ForEachReference(fn func(core.UntypedReference)) {
	i.ElementData.ForEachReference(fn)
	i.GroupData.ForEachReference(fn)
}

type Keyword interface {
	core.AstNode
	Assignable

	IsKeyword()
	Value() string
	ValueToken() *core.Token
	SetValue(value *core.Token)
}

func NewKeyword() Keyword {
	return &KeywordImpl{
		AstNodeBase:    core.NewAstNode(),
		AssignableData: NewAssignableData(),
		ElementData:    NewElementData(),
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

func (i *KeywordData) ForEachReference(fn func(core.UntypedReference)) {
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

func (i *KeywordData) SetValue(value *core.Token) {
	i.value = value
}

type KeywordImpl struct {
	core.AstNodeBase
	AssignableData
	ElementData
	KeywordData
}

func (i *KeywordImpl) ForEachNode(fn func(core.AstNode)) {
	i.AssignableData.ForEachNode(fn)
	i.ElementData.ForEachNode(fn)
	i.KeywordData.ForEachNode(fn)
}

func (i *KeywordImpl) ForEachReference(fn func(core.UntypedReference)) {
	i.AssignableData.ForEachReference(fn)
	i.ElementData.ForEachReference(fn)
	i.KeywordData.ForEachReference(fn)
}

type Assignment interface {
	core.AstNode
	Element

	IsAssignment()
	Property() string
	PropertyToken() *core.Token
	SetProperty(value *core.Token)
	Operator() string
	OperatorToken() *core.Token
	SetOperator(value *core.Token)
	Value() Assignable
	SetValue(value Assignable)
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

func (i *AssignmentData) ForEachReference(fn func(core.UntypedReference)) {
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

func (i *AssignmentData) SetProperty(value *core.Token) {
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

func (i *AssignmentData) SetOperator(value *core.Token) {
	i.operator = value
}

func (i *AssignmentData) Value() Assignable {
	if i != nil && i.value != nil {
		return i.value
	} else {
		return nil
	}
}

func (i *AssignmentData) SetValue(value Assignable) {
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

func (i *AssignmentImpl) ForEachReference(fn func(core.UntypedReference)) {
	i.ElementData.ForEachReference(fn)
	i.AssignmentData.ForEachReference(fn)
}

type Assignable interface {
	core.AstNode
	Element

	IsAssignable()
}

func NewAssignable() Assignable {
	return &AssignableImpl{
		AstNodeBase:    core.NewAstNode(),
		ElementData:    NewElementData(),
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

func (i *AssignableData) ForEachReference(fn func(core.UntypedReference)) {
}

type AssignableImpl struct {
	core.AstNodeBase
	ElementData
	AssignableData
}

func (i *AssignableImpl) ForEachNode(fn func(core.AstNode)) {
	i.ElementData.ForEachNode(fn)
	i.AssignableData.ForEachNode(fn)
}

func (i *AssignableImpl) ForEachReference(fn func(core.UntypedReference)) {
	i.ElementData.ForEachReference(fn)
	i.AssignableData.ForEachReference(fn)
}

type CrossRef interface {
	core.AstNode
	Assignable

	IsCrossRef()
	Type() string
	TypeToken() *core.Token
	SetType(value *core.Token)
	Rule() RuleCall
	SetRule(value RuleCall)
}

func NewCrossRef() CrossRef {
	return &CrossRefImpl{
		AstNodeBase:    core.NewAstNode(),
		AssignableData: NewAssignableData(),
		ElementData:    NewElementData(),
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

func (i *CrossRefData) ForEachReference(fn func(core.UntypedReference)) {
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

func (i *CrossRefData) SetType(value *core.Token) {
	i._Type = value
}

func (i *CrossRefData) Rule() RuleCall {
	if i != nil && i.rule != nil {
		return i.rule
	} else {
		return nil
	}
}

func (i *CrossRefData) SetRule(value RuleCall) {
	i.rule = value
}

type CrossRefImpl struct {
	core.AstNodeBase
	AssignableData
	ElementData
	CrossRefData
}

func (i *CrossRefImpl) ForEachNode(fn func(core.AstNode)) {
	i.AssignableData.ForEachNode(fn)
	i.ElementData.ForEachNode(fn)
	i.CrossRefData.ForEachNode(fn)
}

func (i *CrossRefImpl) ForEachReference(fn func(core.UntypedReference)) {
	i.AssignableData.ForEachReference(fn)
	i.ElementData.ForEachReference(fn)
	i.CrossRefData.ForEachReference(fn)
}

type RuleCall interface {
	core.AstNode
	Assignable

	IsRuleCall()
	Rule() string
	RuleToken() *core.Token
	SetRule(value *core.Token)
}

func NewRuleCall() RuleCall {
	return &RuleCallImpl{
		AstNodeBase:    core.NewAstNode(),
		AssignableData: NewAssignableData(),
		ElementData:    NewElementData(),
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

func (i *RuleCallData) ForEachReference(fn func(core.UntypedReference)) {
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

func (i *RuleCallData) SetRule(value *core.Token) {
	i.rule = value
}

type RuleCallImpl struct {
	core.AstNodeBase
	AssignableData
	ElementData
	RuleCallData
}

func (i *RuleCallImpl) ForEachNode(fn func(core.AstNode)) {
	i.AssignableData.ForEachNode(fn)
	i.ElementData.ForEachNode(fn)
	i.RuleCallData.ForEachNode(fn)
}

func (i *RuleCallImpl) ForEachReference(fn func(core.UntypedReference)) {
	i.AssignableData.ForEachReference(fn)
	i.ElementData.ForEachReference(fn)
	i.RuleCallData.ForEachReference(fn)
}

type Action interface {
	core.AstNode
	Element

	IsAction()
	Type() string
	TypeToken() *core.Token
	SetType(value *core.Token)
	Operator() string
	OperatorToken() *core.Token
	SetOperator(value *core.Token)
	Property() string
	PropertyToken() *core.Token
	SetProperty(value *core.Token)
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
	operator *core.Token
	property *core.Token
}

func NewActionData() ActionData {
	return ActionData{}
}

func (i *ActionData) IsAction() {}

func (i *ActionData) ForEachNode(fn func(core.AstNode)) {
}

func (i *ActionData) ForEachReference(fn func(core.UntypedReference)) {
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

func (i *ActionData) SetType(value *core.Token) {
	i._Type = value
}

func (i *ActionData) Operator() string {
	if i != nil && i.operator != nil {
		return i.operator.Image
	} else {
		return ""
	}
}

func (i *ActionData) OperatorToken() *core.Token {
	return i.operator
}

func (i *ActionData) SetOperator(value *core.Token) {
	i.operator = value
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

func (i *ActionData) SetProperty(value *core.Token) {
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

func (i *ActionImpl) ForEachReference(fn func(core.UntypedReference)) {
	i.ElementData.ForEachReference(fn)
	i.ActionData.ForEachReference(fn)
}
