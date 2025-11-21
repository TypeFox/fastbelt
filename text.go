// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

type TextIndex int32
type TextLine int32
type TextColumn int32

type TextIndexRange struct {
	Start TextIndex
	End   TextIndex
}

type TextLocation struct {
	Line   TextLine
	Column TextColumn
}

type TextRange struct {
	Start TextLocation
	End   TextLocation
}

type TextSegment struct {
	Indices TextIndexRange
	Range   TextRange
}
