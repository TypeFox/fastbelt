package generator

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/test"
)

func GrammarTemplate(rules string) string {
	return "grammar Test;\n" + rules + `
token ID: /[a-zA-Z_][a-zA-Z0-9_]*/;
token NUMBER: /[0-9]+/;
hidden token WS: /[ \n\r\t]+/;
`
}

func FixtureATN(t *testing.T, grammarStr string) (*ATN, map[string]*grammar.ParserRule, map[string]TokenInfo) {
	f := test.New(t, grammar.CreateServices())
	doc := f.Parse(grammarStr).AssertNoErrors()
	grammr := doc.Root().(grammar.Grammar)
	atn, rules, tokenTypes := CreateATN(grammr)
	rtn := BuildRuntimeATN(atn)
	tokenTypeNames := make(map[int]string, len(tokenTypes))
	for name, info := range tokenTypes {
		tokenTypeNames[info.ID] = name
	}
	node := EmitMarkdownSource("test", rtn, tokenTypeNames)
	content := node.String()
	os.WriteFile("atn.test.md", []byte(content), 0644)
	return atn, rules, tokenTypes
}

func RequireATNRecognizes(t *testing.T, atn *ATN, rules map[string]*grammar.ParserRule, tokenTypes map[string]TokenInfo, start string, inputTokenTypes []string, expected int) {
	startRule := *rules[start]
	inputTokenTypeIds := make([]int, len(inputTokenTypes))
	for i, tokenType := range inputTokenTypes {
		info := tokenTypes[tokenType]
		inputTokenTypeIds[i] = info.ID
	}
	actual := recognizeATN(atn, startRule, inputTokenTypeIds)
	require.Equal(t, expected, actual)
}

type recognitionContext struct {
	inputTokenTypeIds []int
	position          int
}

type parserContext struct {
	recognitionContext
	atnState *ATNState
}

func (ctx parserContext) ReadInput() parserContext {
	return parserContext{
		recognitionContext: recognitionContext{
			inputTokenTypeIds: ctx.inputTokenTypeIds,
			position:          ctx.position + 1,
		},
		atnState: ctx.atnState,
	}
}

func (ctx parserContext) Next(atnState *ATNState) parserContext {
	return parserContext{
		recognitionContext: recognitionContext{
			inputTokenTypeIds: ctx.inputTokenTypeIds,
			position:          ctx.position,
		},
		atnState: atnState,
	}
}

func recognizeATN(atn *ATN, startRule grammar.ParserRule, inputTokenTypeIds []int) int {
	recognitionContext := recognitionContext{
		inputTokenTypeIds: inputTokenTypeIds,
		position:          0,
	}
	contexts := recognizeParserRule(atn, startRule, recognitionContext)
	finished := make([]parserContext, 0)
	for _, ctx := range contexts {
		if ctx.position == len(ctx.inputTokenTypeIds) {
			finished = append(finished, ctx)
		}
	}
	return len(finished)
}

func recognizeParserRule(atn *ATN, rule grammar.ParserRule, context recognitionContext) []parserContext {
	END := atn.RuleToStopState[rule]

	contextQueue := []parserContext{
		parserContext{
			recognitionContext: context,
			atnState:           atn.RuleToStartState[rule],
		},
	}

	result := []parserContext{}
	for len(contextQueue) > 0 {
		context := contextQueue[0]
		contextQueue = contextQueue[1:]
		for _, transition := range context.atnState.Transitions {
			if atomTransition, ok := transition.(*AtomTransition); ok {
				if context.position >= len(context.inputTokenTypeIds) {
					continue //reject
				}
				lookahead := context.inputTokenTypeIds[context.position]
				if atomTransition.TokenTypeID == lookahead {
					contextQueue = append(contextQueue, context.ReadInput().Next(atomTransition.Target()))
				}
			} else if ruleTransition, ok := transition.(*RuleTransition); ok {
				nextContext := recognizeParserRule(atn, *ruleTransition.Rule, context.recognitionContext)
				for _, nextCtx := range nextContext {
					contextQueue = append(contextQueue, nextCtx.Next(ruleTransition.FollowState))
				}
			} else if epsilonTransition, ok := transition.(*EpsilonTransition); ok {
				contextQueue = append(contextQueue, context.Next(epsilonTransition.Target()))
			} else {
				panic("unsupported transition type")
			}
		}
		if context.atnState.StateNumber == END.StateNumber {
			result = append(result, context)
		}
	}
	return result
}
