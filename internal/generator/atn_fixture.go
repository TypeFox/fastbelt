package generator

import (
	"testing"

	"github.com/stretchr/testify/require"
	"typefox.dev/fastbelt/internal/grammar"
	"typefox.dev/fastbelt/test"
)

func GrammarTemplate(rules string) string {
	return "grammar Test;\n" + rules + `
token ID: /[a-zA-Z_][a-zA-Z0-9_]*/;
hidden token WS: /[ \n\r\t]+/;
`
}

func FixtureATN(t *testing.T, grammarStr string) (*ATN, map[string]*grammar.ParserRule, map[string]TokenInfo) {
	f := test.New(t, grammar.CreateServices())
	doc := f.Parse(grammarStr).AssertNoErrors()
	grammr := doc.Root().(grammar.Grammar)
	return CreateATN(grammr)
}

func RequireATNRecognizes(t *testing.T, atn *ATN, rules map[string]*grammar.ParserRule, tokenTypes map[string]TokenInfo, start string, inputTokenTypes []string, expectedDecisions [][]int) {
	startRule := *rules[start]
	inputTokenTypeIds := make([]int, len(inputTokenTypes))
	for i, tokenType := range inputTokenTypes {
		info := tokenTypes[tokenType]
		inputTokenTypeIds[i] = info.ID
	}
	actual := recognizeATN(atn, startRule, inputTokenTypeIds)
	require.ElementsMatch(t, expectedDecisions, actual)
}

type recognitionContext struct {
	inputTokenTypeIds []int
	position          int
	output            []int
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
			output:            ctx.output,
		},
		atnState: ctx.atnState,
	}
}

func (ctx parserContext) WriteOutput(decision int) parserContext {
	return parserContext{
		recognitionContext: recognitionContext{
			inputTokenTypeIds: ctx.inputTokenTypeIds,
			position:          ctx.position,
			output:            append(ctx.output, decision),
		},
		atnState: ctx.atnState,
	}
}

func (ctx parserContext) Next(atnState *ATNState) parserContext {
	return parserContext{
		recognitionContext: recognitionContext{
			inputTokenTypeIds: ctx.inputTokenTypeIds,
			position:          ctx.position,
			output:            ctx.output,
		},
		atnState: atnState,
	}
}

func recognizeATN(atn *ATN, startRule grammar.ParserRule, inputTokenTypeIds []int) [][]int {
	recognitionContext := recognitionContext{
		inputTokenTypeIds: inputTokenTypeIds,
		position:          0,
		output:            []int{},
	}
	contexts := recognizeParserRule(atn, startRule, recognitionContext)
	decisionChains := [][]int{}
	for _, context := range contexts {
		decisionChains = append(decisionChains, context.output)
	}
	return decisionChains
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
		decision := context.atnState.Decision
		for _, transition := range context.atnState.Transitions {
			if atomTransition, ok := transition.(*AtomTransition); ok {
				if context.position >= len(context.inputTokenTypeIds) {
					continue //reject
				}
				lookahead := context.inputTokenTypeIds[context.position]
				if atomTransition.TokenTypeID == lookahead {
					contextQueue = append(contextQueue, context.ReadInput().WriteOutput(decision).Next(atomTransition.Target()))
				}
			} else if ruleTransition, ok := transition.(*RuleTransition); ok {
				newContext := context.WriteOutput(decision)
				nextContext := recognizeParserRule(atn, *ruleTransition.Rule, newContext.recognitionContext)
				for _, nextCtx := range nextContext {
					contextQueue = append(contextQueue, nextCtx.Next(ruleTransition.Target()))
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
