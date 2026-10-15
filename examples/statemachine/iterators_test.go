package statemachine

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"typefox.dev/fastbelt"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
)

func BenchmarkTraverseContentSeq(b *testing.B) {
	content, _ := generateStatemachineContent(0)
	srv := CreateServices()
	documentParser := service.MustGet[workspace.DocumentParser](srv)
	doc, err := fastbelt.NewDocumentFromString("file:///workspace/statemachine_0.statemachine", "statemachine", content)
	if err != nil {
		b.Fatal(err)
	}
	documentParser.Parse(b.Context(), doc)

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
	documentParser := service.MustGet[workspace.DocumentParser](srv)
	doc, err := fastbelt.NewDocumentFromString("file:///workspace/statemachine_0.statemachine", "statemachine", content)
	if err != nil {
		b.Fatal(err)
	}
	documentParser.Parse(b.Context(), doc)
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
	documentParser := service.MustGet[workspace.DocumentParser](srv)
	doc, err := fastbelt.NewDocumentFromString("file:///workspace/statemachine_0.statemachine", "statemachine", content)
	if err != nil {
		t.Fatal(err)
	}
	documentParser.Parse(t.Context(), doc)
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
	documentParser.Parse(t.Context(), doc)
	childCount := 0
	for range fastbelt.AllChildren(doc.Root) {
		childCount++
	}
	assert.Equal(t, totalCount, childCount, "AllChildren should iterate all child nodes in the document")
}

func TestAllNodesEarlyExitShortCircuits(t *testing.T) {
	content, elementCount := generateStatemachineContent(0)
	srv := CreateServices()
	documentParser := service.MustGet[workspace.DocumentParser](srv)
	doc, err := fastbelt.NewDocumentFromString("file:///workspace/statemachine_0.statemachine", "statemachine", content)
	if err != nil {
		t.Fatal(err)
	}
	documentParser.Parse(t.Context(), doc)
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
