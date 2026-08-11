// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package statemachine

import (
	"context"
	"reflect"
	"testing"

	"typefox.dev/fastbelt"
	"typefox.dev/fastbelt/linking"
	"typefox.dev/fastbelt/util/extiter"
	"typefox.dev/fastbelt/util/service"
	"typefox.dev/fastbelt/workspace"
)

// Microbenchmarks for the linking primitives. Each benchmark isolates one layer
// of the reference resolution stack so regressions are attributable:
//
//	Reference.Resolve  = resolve machinery + scope construction + scope lookup
//	scope construction = LocalScopeOfType chain + GlobalScopeOfType
//	scope lookup       = SeqScope.ElementByName over the scope chain

// setupLinkedStatemachine parses and fully builds a single generated document
// and returns it together with a typed reference from a transition (a node
// nested two levels below the root, so its scope chain has realistic depth).
func setupLinkedStatemachine(b *testing.B) (*fastbelt.Document, *fastbelt.Reference[State]) {
	b.Helper()
	content, _ := generateStatemachineContent(0)
	srv := CreateServices()
	doc, err := fastbelt.NewDocumentFromString("file:///workspace/statemachine_0.statemachine", "statemachine", content)
	if err != nil {
		b.Fatal(err)
	}
	builder := service.MustGet[workspace.Builder](srv)
	if err := builder.Build(b.Context(), []*fastbelt.Document{doc}, nil); err != nil {
		b.Fatalf("build failed: %v", err)
	}
	for _, ref := range doc.References {
		if typed, ok := ref.(*fastbelt.Reference[State]); ok && typed.Owner().Container() != nil {
			return doc, typed
		}
	}
	b.Fatal("no transition state reference found")
	return nil, nil
}

// BenchmarkRefResolveMachinery measures the Reference[T] resolution machinery
// alone: atomic fast path check, mutex, context.WithValue cycle guard, and the
// desc.Node type assertion. The getter is a constant, so scope construction and
// lookup are excluded.
func BenchmarkRefResolveMachinery(b *testing.B) {
	_, real := setupLinkedStatemachine(b)
	desc := real.Description()
	if desc == nil {
		b.Fatal("reference did not resolve")
	}
	getter := func(ctx context.Context, r *fastbelt.Reference[State]) (*fastbelt.SymbolDescription, *fastbelt.ReferenceError) {
		return desc, nil
	}
	ref := fastbelt.NewReference(real.Owner(), real.Unit(), getter)
	ctx := b.Context()
	b.ResetTimer()
	for b.Loop() {
		ref.Reset()
		ref.Resolve(ctx)
	}
}

// BenchmarkRefResolveFull measures one complete resolution of a real reference:
// machinery + scope construction + lookup. The difference to
// BenchmarkRefResolveMachinery is the cost of the scoping/linking layers.
func BenchmarkRefResolveFull(b *testing.B) {
	_, ref := setupLinkedStatemachine(b)
	ctx := b.Context()
	b.ResetTimer()
	for b.Loop() {
		ref.Reset()
		ref.Resolve(ctx)
	}
}

// BenchmarkRefResolvedFastPath measures repeated access to an already resolved
// reference (the atomic.Bool fast path).
func BenchmarkRefResolvedFastPath(b *testing.B) {
	_, ref := setupLinkedStatemachine(b)
	ctx := b.Context()
	ref.Resolve(ctx)
	b.ResetTimer()
	for b.Loop() {
		_ = ref.Ref(ctx)
	}
}

// BenchmarkScopeConstruction measures building the scope chain for a reference
// owner (LocalScopeOfType per ancestor + GlobalScopeOfType), without any lookup.
func BenchmarkScopeConstruction(b *testing.B) {
	_, ref := setupLinkedStatemachine(b)
	owner := ref.Owner()
	b.ResetTimer()
	for b.Loop() {
		_ = linking.DefaultScopeOfType[State](owner)
	}
}

// BenchmarkScopeLookup measures ElementByName on a prebuilt scope chain,
// i.e. the linear SeqScope scan without scope construction.
func BenchmarkScopeLookup(b *testing.B) {
	_, ref := setupLinkedStatemachine(b)
	scope := linking.DefaultScopeOfType[State](ref.Owner())
	name := ref.Text()
	b.ResetTimer()
	for b.Loop() {
		if scope.ElementByName(name) == nil {
			b.Fatal("lookup failed")
		}
	}
}

// BenchmarkScopeLookupMiss measures ElementByName for a name that does not
// exist, forcing a full scan of the entire scope chain including the global
// scope. This is the worst case hit by every unresolvable reference.
func BenchmarkScopeLookupMiss(b *testing.B) {
	_, ref := setupLinkedStatemachine(b)
	scope := linking.DefaultScopeOfType[State](ref.Owner())
	b.ResetTimer()
	for b.Loop() {
		if scope.ElementByName("does_not_exist") != nil {
			b.Fatal("unexpected hit")
		}
	}
}

// BenchmarkSymbolContainerForType measures the generated SymbolContainer.ForType
// dispatch plus draining the resulting sequence.
func BenchmarkSymbolContainerForType(b *testing.B) {
	doc, ref := setupLinkedStatemachine(b)
	root := ref.Owner()
	for root.Container() != nil {
		root = root.Container()
	}
	container := doc.LocalSymbols.For(root)
	stateType := reflect.TypeFor[State]()
	b.ResetTimer()
	for b.Loop() {
		count := 0
		for range container.ForType(stateType) {
			count++
		}
		if count == 0 {
			b.Fatal("no symbols")
		}
	}
}

// BenchmarkIsEmptyOverContainer measures extiter.IsEmpty on a container-backed
// sequence — the check LocalScopeOfType performs at every ancestor level.
func BenchmarkIsEmptyOverContainer(b *testing.B) {
	doc, ref := setupLinkedStatemachine(b)
	root := ref.Owner()
	for root.Container() != nil {
		root = root.Container()
	}
	container := doc.LocalSymbols.For(root)
	stateType := reflect.TypeFor[State]()
	b.ResetTimer()
	for b.Loop() {
		if extiter.IsEmpty(container.ForType(stateType)) {
			b.Fatal("unexpectedly empty")
		}
	}
}

// BenchmarkSymbolDescriptionName measures SymbolDescription.Name.String(),
// which SeqScope.ElementByName calls once per candidate symbol.
func BenchmarkSymbolDescriptionName(b *testing.B) {
	_, ref := setupLinkedStatemachine(b)
	desc := ref.Description()
	if desc == nil {
		b.Fatal("reference did not resolve")
	}
	b.ResetTimer()
	for b.Loop() {
		_ = desc.Name.String()
	}
}
