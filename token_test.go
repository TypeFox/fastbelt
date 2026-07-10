// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package fastbelt

import "testing"

func TestSearchOffset2(t *testing.T) {
	tt := &TokenType{Id: 1, Name: "t"}
	tokens := TokenSlice{
		NewToken(tt, "aa", 2, 4),
		NewToken(tt, "bb", 6, 8),
		NewToken(tt, "cc", 8, 10),
	}

	cases := []struct {
		name     string
		offset   int
		wantPrev int // -1 for nil, else index into tokens
		wantNext int
	}{
		{"inside first token", 3, 0, -1},
		{"at start of first token", 2, 0, -1},
		{"at the end of first token", 4, 0, -1},
		{"boundary between token 1 and 2", 8, 1, 2},
		{"inside last token", 9, 2, -1},
		{"at end of last token", 10, 2, -1},
		{"past end of last token", 11, -1, -1},
		{"before start of first token", 1, -1, -1},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			prev, next := tokens.SearchOffset2(c.offset)

			var gotPrev, gotNext = -1, -1
			for i := range tokens {
				if prev == &tokens[i] {
					gotPrev = i
				}
				if next == &tokens[i] {
					gotNext = i
				}
			}

			if gotPrev != c.wantPrev || gotNext != c.wantNext {
				t.Errorf("SearchOffset2(%d) = (%d, %d), want (%d, %d)",
					c.offset, gotPrev, gotNext, c.wantPrev, c.wantNext)
			}
		})
	}
}
