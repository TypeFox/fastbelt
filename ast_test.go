// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

import (
	"math/rand/v2"
	"testing"
)

func BenchmarkCompositeNodeSequential(b *testing.B) {
	const (
		nodeCount        = 10_000
		minTokensPerNode = 4
		maxTokensPerNode = 8
	)

	nodes := generateRandomNodes(nodeCount, minTokensPerNode, maxTokensPerNode)

	b.ResetTimer()
	// Sequential access should mostly hit the warm cache path after the first iteration
	// and thus serve as a baseline for the best-case performance of the string cache.
	for b.Loop() {
		for i := range nodes {
			_ = nodes[i].String()
		}
	}
}

// BenchmarkCompositeNodeConcurrent benchmarks concurrent access to the
// atomic-pointer string cache on CompositeNodeBase.
//
// The setup creates a large pool of nodes whose tokens have nil TokenType
// (only Image is used by stringSlow). Each goroutine picks nodes pseudo-randomly
// so that some nodes are accessed while their cache is still cold (races to
// populate the atomic pointer) and others are accessed after the cache is warm.
// The access pattern is deliberately non-uniform to stress both paths.
func BenchmarkCompositeNodeConcurrent(b *testing.B) {
	const (
		nodeCount        = 10_000
		minTokensPerNode = 4
		maxTokensPerNode = 8
	)

	nodes := generateRandomNodes(nodeCount, minTokensPerNode, maxTokensPerNode)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		// Per-goroutine RNG avoids inter-goroutine synchronisation on the
		// RNG itself, keeping contention squarely on the atomic pointer.
		rng := rand.New(rand.NewPCG(rand.Uint64(), rand.Uint64()))

		for pb.Next() {
			// Pick a random starting index, then access a small cluster of
			// nearby nodes. This creates "hot spots" that exercise the warm-
			// cache path while the tail of the pool stays cold longer.
			base := rng.IntN(nodeCount)
			accesses := rng.IntN(4) + 1
			for range accesses {
				// Spread within a 64-node window around the base index so
				// that goroutines frequently race on the same nodes.
				idx := (base + rng.IntN(64)) % nodeCount
				_ = nodes[idx].String()
			}
		}
	})
}

func generateRandomNodes(nodeCount, minTokensPerNode, maxTokensPerNode int) []CompositeNode {
	// A small, fixed word list keeps token images cheap to allocate while
	// still producing non-trivial composite strings.
	words := [16]string{
		"alpha", "beta", "gamma", "delta",
		"epsilon", "zeta", "eta", "theta",
		"iota", "kappa", "lambda", "mu",
		"nu", "xi", "omicron", "pi",
	}

	// Build the node pool using a deterministic seed so the benchmark is
	// reproducible, but varied enough to defeat any trivial optimisation.
	src := rand.New(rand.NewPCG(0xdeadbeef, 0xcafebabe))
	nodes := make([]CompositeNode, nodeCount)
	for i := range nodes {
		node := NewCompositeNode()
		tokenCount := src.IntN(maxTokensPerNode-minTokensPerNode+1) + minTokensPerNode
		for range tokenCount {
			// String only reads Token.Image.
			node.SetToken(&Token{Image: words[src.IntN(len(words))]})
		}
		nodes[i] = node
	}
	return nodes
}
