// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package statemachine

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"typefox.dev/fastbelt"
	"typefox.dev/fastbelt/lexer"
	"typefox.dev/fastbelt/parser"
	"typefox.dev/fastbelt/util/service"
)

func parseDocument(tb testing.TB, srv *service.Container, content string) *fastbelt.Document {
	tb.Helper()
	lexer := service.MustGet[lexer.Lexer](srv)
	parser := service.MustGet[parser.Parser](srv)
	doc := fastbelt.NewDocumentFromString("file:///workspace/test.statemachine", "statemachine", content)
	if err := lexer.Exec(context.Background(), doc); err != nil {
		tb.Fatal(err)
	}
	if err := parser.Parse(context.Background(), doc); err != nil {
		tb.Fatal(err)
	}
	return doc
}

func BenchmarkTraverseContentSeq(b *testing.B) {
	content, _ := generateStatemachineContent(0)
	srv := CreateServices()
	doc := parseDocument(b, srv, content)

	for b.Loop() {
		count := 0
		for range fastbelt.AllChildren(doc.Root) {
			count++
		}
		_ = count
	}
}

// BenchmarkTraverseContentSeqHalf breaks out of the loop after visiting half of the
// nodes. With short-circuiting traversal this should take roughly half the time of
// BenchmarkTraverseContentSeq.
func BenchmarkTraverseContentSeqHalf(b *testing.B) {
	content, elementCount := generateStatemachineContent(0)
	srv := CreateServices()
	doc := parseDocument(b, srv, content)
	half := elementCount / 2

	for b.Loop() {
		count := 0
		for range fastbelt.AllChildren(doc.Root) {
			count++
			if count == half {
				break
			}
		}
		_ = count
	}
}

func TestAllNodesEquivalence(t *testing.T) {
	content, elementCount := generateStatemachineContent(0)
	srv := CreateServices()
	doc := parseDocument(t, srv, content)
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
	doc := parseDocument(t, srv, content)
	childCount := 0
	for range fastbelt.AllChildren(doc.Root) {
		childCount++
	}
	assert.Equal(t, totalCount, childCount, "AllChildren should iterate all child nodes in the document")
}

func TestAllNodesEarlyExitShortCircuits(t *testing.T) {
	content, elementCount := generateStatemachineContent(0)
	srv := CreateServices()
	doc := parseDocument(t, srv, content)
	visited := 0
	for range fastbelt.AllNodes(doc.Root) {
		visited++
		if visited == 3 {
			break
		}
	}
	assert.Equal(t, 3, visited)
	assert.Greater(t, elementCount, 3)
}
