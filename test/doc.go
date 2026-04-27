// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Package test provides utilities for testing Fastbelt language implementations.
//
// Create a [Fixture] from your language's service container, then call [Fixture.Parse]
// or [Fixture.ParseAll] to build documents. Use [Doc] assertion methods and the
// generic Find/MustFind functions to inspect results.
//
// # Markers
//
// Content strings may embed position markers that are stripped before parsing.
// Positions are recorded relative to the cleaned text.
//
//   - Range marker:  <|label:text|>  — spans "text"; label identifies the range.
//   - Range shorthand: <|label|>     — equivalent to <|label:label|>.
//   - Empty range:   <|label:|>      — zero-length span at the marker position.
//   - Index marker:  <|label>        — single position (no text content).
//
// Range markers are tried before index markers when both share the same opening
// delimiter (the default). Marker delimiters are configurable via [TestMarking].
package test
