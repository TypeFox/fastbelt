// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Package multilang is an example of a single Fastbelt language server that
// serves two languages ("greeting", *.hello and "farewell", *.bye) from one
// grammar with two entry rules. Documents are routed to the matching entry
// rule by the registered core.LanguageSelector.
package multilang

//go:generate go run ./gen
