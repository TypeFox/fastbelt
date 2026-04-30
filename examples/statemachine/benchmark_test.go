// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package statemachine

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"typefox.dev/fastbelt"
	"typefox.dev/fastbelt/lexer"
	"typefox.dev/fastbelt/parser"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
)

const resourceCount = 200

// BenchmarkWorkspaceCycle benchmarks a full workspace build cycle over
// resourceCount in-memory statemachine documents: parse, index, link, and validate.
func BenchmarkWorkspaceCycle(b *testing.B) {
	length := 0
	contents := make([]string, resourceCount)
	for i := range contents {
		contents[i], _ = generateStatemachineContent(i)
		length += len(contents[i])
	}
	srv := CreateServices()
	docManager := service.MustGet[workspace.DocumentManager](srv)
	lock := service.MustGet[workspace.Lock](srv)
	builder := service.MustGet[workspace.Builder](srv)

	var totalNs int64
	b.SetBytes(int64(length))
	b.ResetTimer()
	for range b.N {
		// Fresh document manager per cycle so each build starts from a clean state.
		docManager.Clear()

		docs := make([]*fastbelt.Document, resourceCount)
		for i, content := range contents {
			uri := fmt.Sprintf("file:///workspace/statemachine_%d.statemachine", i)
			doc, err := fastbelt.NewDocumentFromString(uri, "statemachine", content)
			if err != nil {
				b.Fatal(err)
			}
			docs[i] = doc
			docManager.Set(doc)
		}

		start := time.Now()
		lock.Write(context.Background(), func(ctx context.Context, downgrade func()) {
			if err := builder.Build(ctx, docs, downgrade); err != nil {
				b.Errorf("build failed: %v", err)
			}
		})
		totalNs += time.Since(start).Nanoseconds()
	}

	msPerOp := float64(totalNs) / float64(b.N) / 1e6
	msPerResource := msPerOp / resourceCount
	b.ReportMetric(msPerOp, "ms/op")
	b.ReportMetric(msPerResource, "ms/resource")
}

func BenchmarkTraverseContentSeq(b *testing.B) {
	content, _ := generateStatemachineContent(0)
	srv := CreateServices()
	documentParser := service.MustGet[workspace.DocumentParser](srv)
	doc, err := fastbelt.NewDocumentFromString("file:///workspace/statemachine_0.statemachine", "statemachine", content)
	if err != nil {
		b.Fatal(err)
	}
	documentParser.Parse(doc)

	for b.Loop() {
		count := 0
		for range fastbelt.AllChildren(doc.Root) {
			count++
		}
		_ = count
	}
}

func TestAllNodesEquivalence(t *testing.T) {
	content, elementCount := generateStatemachineContent(0)
	srv := CreateServices()
	documentParser := service.MustGet[workspace.DocumentParser](srv)
	doc, err := fastbelt.NewDocumentFromString("file:///workspace/statemachine_0.statemachine", "statemachine", content)
	if err != nil {
		t.Fatal(err)
	}
	documentParser.Parse(doc)
	nodeCount := 0
	for range fastbelt.AllNodes(doc.Root) {
		nodeCount++
	}
	assert.Equal(t, elementCount, nodeCount, "AllNodes should iterate all nodes in the document")
}

func TestAllChildrenEquivalence(t *testing.T) {
	content, elementCount := generateStatemachineContent(0)
	totalCount := elementCount - 1 // AllChildren does not include the root node, so we subtract 1 from the total count
	srv := CreateServices()
	documentParser := service.MustGet[workspace.DocumentParser](srv)
	doc, err := fastbelt.NewDocumentFromString("file:///workspace/statemachine_0.statemachine", "statemachine", content)
	if err != nil {
		t.Fatal(err)
	}
	documentParser.Parse(doc)
	childCount := 0
	for range fastbelt.AllChildren(doc.Root) {
		childCount++
	}
	assert.Equal(t, totalCount, childCount, "AllChildren should iterate all child nodes in the document")
}

// BenchmarkParser benchmarks parsing a single generated statemachine document,
// reusing the pre-lexed token slice every iteration.
func BenchmarkParser(b *testing.B) {
	content, _ := generateStatemachineContent(0)
	srv := CreateServices()
	lexerService := service.MustGet[lexer.Lexer](srv)
	parserService := service.MustGet[parser.Parser](srv)
	tokens := lexerService.Lex(content).Tokens
	doc, err := fastbelt.NewDocumentFromString("file:///workspace/statemachine_0.statemachine", "statemachine", content)
	if err != nil {
		b.Fatal(err)
	}
	doc.Tokens = tokens
	b.SetBytes(int64(len(content)))
	b.ResetTimer()
	for b.Loop() {
		result := parserService.Parse(doc)
		doc.Root = result.Node
	}
}

// BenchmarkLexer benchmarks tokenizing a single generated statemachine document.
func BenchmarkLexer(b *testing.B) {
	content, _ := generateStatemachineContent(0)
	l := NewLexer()
	b.SetBytes(int64(len(content)))
	b.ResetTimer()
	for b.Loop() {
		_ = l.Lex(content)
	}
}

// generateStatemachineContent generates a syntactically valid statemachine
// document for the given index. Each document contains:
//   - 4 events
//   - 3 commands
//   - 50 states, each with transitions that cycle through events/states
func generateStatemachineContent(index int) (string, int) {
	const numEvents = 4
	const numCommands = 3
	const numStates = 50
	elementCount := 1

	var sb strings.Builder

	fmt.Fprintf(&sb, "statemachine sm%d\n\n", index)

	// Events block
	sb.WriteString("events\n")
	for e := range numEvents {
		fmt.Fprintf(&sb, "  evt%d_%d\n", index, e)
		elementCount++
	}

	// Commands block
	sb.WriteString("commands\n")
	for c := range numCommands {
		fmt.Fprintf(&sb, "  cmd%d_%d\n", index, c)
		elementCount++
	}

	// Initial state
	fmt.Fprintf(&sb, "initialState s%d_0\n\n", index)

	// States
	for s := range numStates {
		fmt.Fprintf(&sb, "state s%d_%d\n", index, s)
		if s == 0 {
			fmt.Fprintf(&sb, "  actions { cmd%d_0 cmd%d_1 }\n", index, index)
		}
		for e := range numEvents {
			target := (s + e + 1) % numStates
			fmt.Fprintf(&sb, "  evt%d_%d => s%d_%d\n", index, e, index, target)
			elementCount++
		}
		sb.WriteString("end\n\n")
		elementCount++
	}

	return sb.String(), elementCount
}
