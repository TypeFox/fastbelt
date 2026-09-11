// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package json

// this language with its grammar interface definitions aims at checking the absence of any naming conflicts
//  between generated property inspired names and names hard-coded in generator of `json_gen.go` files

// therefore only `types_gen.go`, `json_gen.go`, and `doc.go` files are put into the repository,
//  while all other (generated) files are ignored via a local `.gitignore` configuration;

//go:generate go run ../../../cmd/fastbelt generate ./json.fb -v
