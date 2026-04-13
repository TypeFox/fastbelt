package grammar

import (
	"context"
	"testing"
	"fmt"
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/parser"
	allstar "typefox.dev/fastbelt/parser/allstar"
)

func TestDebugGroupOption(t *testing.T) {
	// Grammar: 'Person returns Person: "person" Name=ID Age=ID;'
	// After parsing '"person"' (Keyword), next token should be Name=ID
	// Group_Option_1 should return true (enter optional group)
	
	input := `grammar Test;
interface Person { Name string Age string }
Person returns Person: "person" Name=ID Age=ID;
token ID: /[a-zA-Z_][a-zA-Z0-9_]*/;
hidden token WS: /[ \n\r\t]+/;
`
	lexResult := NewLexer().Lex(input)
	fmt.Printf("Tokens count: %d\n", len(lexResult.Tokens))
	
	rtn := buildATN()
	
	// Check Group_Option_1 state
	groupOptState := rtn.DecisionMap["Group_Option_1"]
	if groupOptState == nil {
		t.Fatal("Group_Option_1 not found in decision map")
	}
	fmt.Printf("Group_Option_1: State#%d, Decision=%d\n", groupOptState.StateNumber, groupOptState.Decision)

	// Test prediction with a token stream that simulates being after '"person"'
	// In a real parse, after parsing Person rule's first Keyword '"person"',
	// the state is before 'Name=ID'
	// Let's tokenize just the rule body: "person" Name=ID Age=ID
	bodyTokens := lexResult.Tokens
	// Find position of '"person"' String token
	var stringIdx int
	for i, tok := range bodyTokens {
		if tok != nil && tok.TypeId == Token_String_Idx {
			stringIdx = i
			break
		}
	}
	fmt.Printf("Found String token at index %d\n", stringIdx)
	
	// Set state to be AFTER the String token (i.e., at Name=ID)
	state := parser.NewParserState(bodyTokens)
	// Advance to just after the String token
	for i := 0; i <= stringIdx; i++ {
		state.Consume(bodyTokens[i].TypeId)
	}
	
	la1 := state.LA(1)
	if la1 != nil {
		fmt.Printf("After String token, LA(1): TypeId=%d Image=%q\n", la1.TypeId, la1.Image)
	}
	
	strat := allstar.NewLLStarLookaheadFromRuntime(rtn, nil).AsParserStrategy()
	result := strat.PredictOpt(state, "Group_Option_1")
	fmt.Printf("PredictOpt('Group_Option_1'): %v\n", result)
}

func TestDebugAssignmentScope(t *testing.T) {
	doc, err := core.NewDocumentFromString("inmemory://debug.fb", "fastbelt", `
grammar Test;
interface Person { Name string Age string }
Person returns Person: "person" Name=ID Age=ID;
token ID: /[a-zA-Z_][a-zA-Z0-9_]*/;
hidden token WS: /[ \n\r\t]+/;
`)
	if err != nil {
		t.Fatalf("failed: %v", err)
	}
	
	srv := CreateServices()
	ctx := context.Background()
	err = srv.Workspace().Builder.Build(ctx, []*core.Document{doc})
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	
	grammar, ok := doc.Root.(Grammar)
	if !ok {
		t.Fatalf("Root is not Grammar: %T", doc.Root)
	}
	rules := grammar.Rules()
	fmt.Printf("Rules: %d, Interfaces: %d\n", len(rules), len(grammar.Interfaces()))
	if len(rules) > 0 {
		body := rules[0].Body()
		fmt.Printf("Body type: %T\n", body)
	}
	_ = ctx
}
