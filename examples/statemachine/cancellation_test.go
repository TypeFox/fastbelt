package statemachine

import (
	"context"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"typefox.dev/fastbelt"
	"typefox.dev/fastbelt/lexer"
	"typefox.dev/fastbelt/parser"
	"typefox.dev/fastbelt/util/service"
)

func TestParseCancelledContextTerminates(t *testing.T) {
	content, _ := generateStatemachineContent(0)
	srv := CreateServices()
	doc, err := fastbelt.NewDocumentFromString("file:///workspace/statemachine_0.statemachine", "statemachine", content)
	if err != nil {
		t.Fatal(err)
	}
	lexer := service.MustGet[lexer.Lexer](srv)
	parser := service.MustGet[parser.Parser](srv)
	doc.Tokens = lexer.Exec(t.Context(), content).Tokens
	ctx, cancel := context.WithCancel(t.Context())
	// Cancel ahead of time
	cancel()
	parser.Parse(ctx, doc)
}

// hugeStateCount yields a ~5 MB document, large enough that timer and
// scheduling jitter stay small relative to a full lexer or parser run.
const hugeStateCount = 50_000

func TestLexerCancellationStopsEarly(t *testing.T) {
	content, _ := generateStatemachineContentWithStates(0, hugeStateCount)
	lexer := service.MustGet[lexer.Lexer](CreateServices())
	assertCancellationStopsEarly(t, func(ctx context.Context) { lexer.Exec(ctx, content) })
}

func TestParserCancellationStopsEarly(t *testing.T) {
	content, _ := generateStatemachineContentWithStates(0, hugeStateCount)
	srv := CreateServices()
	doc, err := fastbelt.NewDocumentFromString("file:///workspace/statemachine_0.statemachine", "statemachine", content)
	if err != nil {
		t.Fatal(err)
	}
	doc.Tokens = service.MustGet[lexer.Lexer](srv).Exec(t.Context(), content).Tokens
	parser := service.MustGet[parser.Parser](srv)
	assertCancellationStopsEarly(t, func(ctx context.Context) { parser.Parse(ctx, doc) })
}

// assertCancellationStopsEarly times a full run, then reruns with
// cancellation fired after 20% of that time and asserts the cancelled run
// finishes within 30% of it.
func assertCancellationStopsEarly(t *testing.T, run func(ctx context.Context)) {
	t.Helper()
	run(t.Context()) // warm up

	runtime.GC()
	start := time.Now()
	run(t.Context())
	full := time.Since(start)

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()
	runtime.GC()
	start = time.Now()
	timer := time.AfterFunc(full/5, cancel)
	defer timer.Stop()
	run(ctx)
	cancelled := time.Since(start)

	t.Logf("full: %v, cancelled: %v", full, cancelled)
	assert.Less(t, cancelled, full*3/10)
}
