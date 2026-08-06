// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"context"

	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/server"
)

var legendProvider = server.NewExtendableSemanticTokensLegendProvider()

type GrammarTokenHighlightingStrategy struct{}

func NewGrammarTokenHighlightingStrategy() server.TokenHighlightingStrategy {
	return &GrammarTokenHighlightingStrategy{}
}

func (s *GrammarTokenHighlightingStrategy) Highlight(ctx context.Context, token core.Token, accept server.TokenHighlightingStrategyAcceptor) {
	switch token.Kind {
	case Grammar_Name_ID:
		accept(legendProvider.Namespace(), 0)
	case Interface_Name_ID,
		Interface_Extends_ID_0,
		Interface_Extends_ID_1:
		accept(legendProvider.Interface(), 0)
	case ParserRule_Name_ID,
		Token_Name_ID,
		CompositeRule_Name_ID,
		RuleCall_Rule_ID,
		TokenGroup_Name_ID,
		TokenGroup_TokenRefs_ID:
		accept(legendProvider.Function(), 0)
	case Field_Name_ID,
		Assignment_Property_ID,
		Action_Property_ID:
		accept(legendProvider.Property(), 0)
	case PrimitiveType_Type_bool,
		PrimitiveType_Type_composite,
		PrimitiveType_Type_string,
		SimpleType_Type_ID,
		ReferenceType_Type_ID,
		CrossRef_Type_ID,
		ParserRule_ReturnType_ID,
		Action_Type_ID,
		Action_current:
		accept(legendProvider.Type(), 0)
	case Token_Type_comment,
		Token_Type_hidden:
		accept(legendProvider.Modifier(), 0)
	}
}
