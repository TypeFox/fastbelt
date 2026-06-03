// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Package codegen provides a small builder for assembling multi-line source
// text with indentation. Fastbelt's own code generators use it to construct
// Go (and other) source before passing the result to a formatter. Fastbelt
// adopters can use the same API when writing custom code generators for their
// own languages.
//
// Create a root with [NewNode], write lines with [Node.Append] and
// [Node.AppendLine], increase nesting with [Node.Indent], splice completed
// fragments with [Node.AppendNode], and read the result from [Node.String].
// Line breaks in the output use [EOL], which matches the host platform.
package codegen
