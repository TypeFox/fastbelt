// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package grammar

import (
	"testing"

	"typefox.dev/fastbelt/test"
)

func TestSemanticTokensIntegration(t *testing.T) {
	fixture := test.New(t, CreateServices())

	grammarText := `<|comment:// Grammar for semantic token testing|>
grammar <|namespace:Test|>;

interface <|interface:Expression|> {}
interface <|interface:BinaryExpression|> extends <|interface:Expression|> {
	<|property:Left|> <|type:Expression|>
	<|property:Operator|> <|type:string|>
	<|property:Right|> <|type:Expression|>
}

<|function:Addition|> returns <|type:Expression|>:
	<|function:Primary|>
	({<|type:BinaryExpression|>.<|property:Left|>=<|type:current|>}
		<|property:Operator|>=("+" | "-") <|property:Right|>=<|function:Primary|>)*
<|function:Primary|> returns <|type:Expression|>:
	<|property:Operator|>=<|function:ID|>

token <|function:ID|>: /[a-zA-Z_][a-zA-Z0-9_]*/;
<|modifier:hidden|> token <|function:WS|>: /[ \n\r\t]+/;
`

	doc := fixture.ParseURI(grammarText, "file:///semantic.fb")
	doc.AssertNoParseErrors().
		AssertSemanticTokens("namespace", legendProvider.Namespace(), 0).
		AssertSemanticTokens("interface", legendProvider.Interface(), 0).
		AssertSemanticTokens("function", legendProvider.Function(), 0).
		AssertSemanticTokens("property", legendProvider.Property(), 0).
		AssertSemanticTokens("type", legendProvider.Type(), 0).
		AssertSemanticTokens("modifier", legendProvider.Modifier(), 0).
		AssertSemanticTokens("comment", legendProvider.Comment(), 0)
}
