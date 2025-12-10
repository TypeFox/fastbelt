package automatons

import (
	"testing"
)

func TestRuneSet_Add(t *testing.T) {
	t.Run("EmptySet_AddSingleCharacter_Success", func(t *testing.T) {
		c := rune(123)
		runeSet := NewRuneSet_Empty()
		runeSet.AddRune(c)
		if runeSet.Length() != 1 {
			t.Errorf("Expected length 1, got %d", runeSet.Length())
		}
		if !runeSet.IncludesRune(c) {
			t.Errorf("Expected rune %d to be included", c)
		}
	})

	t.Run("EmptySet_AddCharacterRange_Success", func(t *testing.T) {
		cFrom := rune(123)
		cTo := rune(125)
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange(cFrom, cTo)
		if runeSet.Length() != 3 {
			t.Errorf("Expected length 3, got %d", runeSet.Length())
		}
		if !runeSet.IncludesRange(cFrom, cTo) {
			t.Errorf("Expected range %d-%d to be included", cFrom, cTo)
		}
		if !runeSet.IncludesRune(cFrom) {
			t.Errorf("Expected rune %d to be included", cFrom)
		}
		if !runeSet.IncludesRune(cFrom + 1) {
			t.Errorf("Expected rune %d to be included", cFrom+1)
		}
		if !runeSet.IncludesRune(cFrom + 2) {
			t.Errorf("Expected rune %d to be included", cFrom+2)
		}
	})

	t.Run("EmptySet_AddBadCharacterRange_Fail", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for invalid range order")
			}
		}()
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange(125, 123)
	})

	t.Run("OneCharSet_AddSingleCharacterSame_Success", func(t *testing.T) {
		c := rune(123)
		runeSet := NewRuneSet_Empty()
		runeSet.AddRune(c)
		runeSet.AddRune(c)
		if runeSet.Length() != 1 {
			t.Errorf("Expected length 1, got %d", runeSet.Length())
		}
		if !runeSet.IncludesRune(c) {
			t.Errorf("Expected rune %d to be included", c)
		}
	})

	t.Run("OneCharSet_AddSingleCharacterBeside_Success", func(t *testing.T) {
		c1 := rune(123)
		c2 := rune(124)
		runeSet := NewRuneSet_Empty()
		runeSet.AddRune(c1)
		runeSet.AddRune(c2)
		if runeSet.Length() != 2 {
			t.Errorf("Expected length 2, got %d", runeSet.Length())
		}
		if !runeSet.IncludesRange(c1, c2) {
			t.Errorf("Expected range %d-%d to be included", c1, c2)
		}
	})

	t.Run("OneCharSet_AddSingleCharacterWithGap_Success", func(t *testing.T) {
		c1 := rune(123)
		c2 := rune(125)
		runeSet := NewRuneSet_Empty()
		runeSet.AddRune(c1)
		runeSet.AddRune(c2)
		if runeSet.Length() != 2 {
			t.Errorf("Expected length 2, got %d", runeSet.Length())
		}
		if runeSet.IncludesRange(c1, c2) {
			t.Errorf("Expected range %d-%d NOT to be included (has gap)", c1, c2)
		}
		if !runeSet.IncludesRune(c1) {
			t.Errorf("Expected rune %d to be included", c1)
		}
		if !runeSet.IncludesRune(c2) {
			t.Errorf("Expected rune %d to be included", c2)
		}
	})

	t.Run("OneCharSet_AddCharacterRangeWithin_Success", func(t *testing.T) {
		cFrom1 := rune(1)
		cTo1 := rune(3)
		cFrom2 := rune(4)
		cTo2 := rune(6)
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange(cFrom1, cTo1)
		runeSet.AddRange(cFrom2, cTo2)
		if runeSet.Length() != 6 {
			t.Errorf("Expected length 6, got %d", runeSet.Length())
		}
		if !runeSet.IncludesRange(cFrom1, cTo1) {
			t.Errorf("Expected range %d-%d to be included", cFrom1, cTo1)
		}
	})

	t.Run("OneCharSet_AddCharacterRangeOverlapping_Success", func(t *testing.T) {
		cFrom1 := rune(1)
		cTo1 := rune(3)
		cFrom2 := rune(5)
		cTo2 := rune(7)
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange(cFrom1, cTo1)
		runeSet.AddRange(cFrom2, cTo2)
		if runeSet.Length() != 6 {
			t.Errorf("Expected length 6, got %d", runeSet.Length())
		}
		if runeSet.IncludesRange(cFrom1, cTo2) {
			t.Errorf("Expected range %d-%d NOT to be included (has gap)", cFrom1, cTo2)
		}
		if !runeSet.IncludesRange(cFrom2, cTo2) {
			t.Errorf("Expected range %d-%d to be included", cFrom2, cTo2)
		}
	})

	t.Run("ThreeCharsSet_AddSingleCharacterAlreadyExists_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange('a', 'c')
		runeSet.AddRune('b')
		if runeSet.Length() != 3 {
			t.Errorf("Expected length 3, got %d", runeSet.Length())
		}
		if !runeSet.IncludesRange('a', 'c') {
			t.Errorf("Expected range 'a'-'c' to be included")
		}
	})

	t.Run("ThreeCharsSet_AddSingleCharacterBeside_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange('a', 'c')
		runeSet.AddRune('d')
		if runeSet.Length() != 4 {
			t.Errorf("Expected length 4, got %d", runeSet.Length())
		}
		if !runeSet.IncludesRange('a', 'd') {
			t.Errorf("Expected range 'a'-'d' to be included")
		}
	})

	t.Run("ThreeCharsSet_AddSingleCharacterWithGap_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange('a', 'c')
		runeSet.AddRune('e')
		if runeSet.Length() != 4 {
			t.Errorf("Expected length 4, got %d", runeSet.Length())
		}
		if runeSet.IncludesRange('a', 'e') {
			t.Errorf("Expected range 'a'-'e' NOT to be included (has gap)")
		}
		if !runeSet.IncludesRune('e') {
			t.Errorf("Expected rune 'e' to be included")
		}
	})

	t.Run("ThreeCharsSet_AddCharacterRangeWithin_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange('a', 'c')
		runeSet.AddRange('b', 'c')
		if runeSet.Length() != 3 {
			t.Errorf("Expected length 3, got %d", runeSet.Length())
		}
		if !runeSet.IncludesRange('a', 'c') {
			t.Errorf("Expected range 'a'-'c' to be included")
		}
	})

	t.Run("ThreeCharsSet_AddCharacterRangeOverlapping_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange('a', 'c')
		runeSet.AddRange('b', 'd')
		if runeSet.Length() != 4 {
			t.Errorf("Expected length 4, got %d", runeSet.Length())
		}
		if !runeSet.IncludesRange('a', 'd') {
			t.Errorf("Expected range 'a'-'d' to be included")
		}
	})

	t.Run("ThreeCharsSet_AddCharacterRangeBadRange_Fail", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic for invalid range order")
			}
		}()
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange('a', 'c')
		runeSet.AddRange('d', 'b')
	})

	t.Run("ThreeCharsSet_AddCharacterRangeOverlappingBothSides_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange('d', 'f')
		runeSet.AddRange('b', 'h')
		if runeSet.Length() != 7 {
			t.Errorf("Expected length 7, got %d", runeSet.Length())
		}
		if !runeSet.IncludesRange('b', 'h') {
			t.Errorf("Expected range 'b'-'h' to be included")
		}
	})

	t.Run("ThreeCharsSet_AddCharacterRangeCoveringAll_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange('d', 'f')
		runeSet.AddRange('a', 'z')
		if runeSet.Length() != 26 {
			t.Errorf("Expected length 26, got %d", runeSet.Length())
		}
		if !runeSet.IncludesRange('a', 'z') {
			t.Errorf("Expected range 'a'-'z' to be included")
		}
	})

	t.Run("ThreeCharsSet_AddCharacterRangeBeside_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange('a', 'c')
		runeSet.AddRange('d', 'f')
		if runeSet.Length() != 6 {
			t.Errorf("Expected length 6, got %d", runeSet.Length())
		}
		if !runeSet.IncludesRange('a', 'f') {
			t.Errorf("Expected range 'a'-'f' to be included")
		}
	})

	t.Run("ThreeCharsSet_AddCharacterRangeWithGap_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange('a', 'c')
		runeSet.AddRange('e', 'g')
		if runeSet.Length() != 6 {
			t.Errorf("Expected length 6, got %d", runeSet.Length())
		}
		if runeSet.IncludesRange('a', 'g') {
			t.Errorf("Expected range 'a'-'g' NOT to be included (has gap)")
		}
		if !runeSet.IncludesRange('e', 'g') {
			t.Errorf("Expected range 'e'-'g' to be included")
		}
	})

	t.Run("TwoThreeGroupsCharsSet_AddCharacterFillGap_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange('a', 'c')
		runeSet.AddRange('e', 'g')
		runeSet.AddRune('d')
		if runeSet.Length() != 7 {
			t.Errorf("Expected length 7, got %d", runeSet.Length())
		}
		if !runeSet.IncludesRange('a', 'g') {
			t.Errorf("Expected range 'a'-'g' to be included")
		}
	})

	t.Run("TwoThreeGroupsCharsSet_AddCharacterRangeFillGap_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange('a', 'c')
		runeSet.AddRange('g', 'i')
		runeSet.AddRange('d', 'f')
		if runeSet.Length() != 9 {
			t.Errorf("Expected length 9, got %d", runeSet.Length())
		}
		if !runeSet.IncludesRange('a', 'i') {
			t.Errorf("Expected range 'a'-'i' to be included")
		}
	})
}

func TestRuneSet_Remove(t *testing.T) {
	t.Run("EmptyCharSet_RemoveSingleCharacter_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Empty()
		runeSet.RemoveRune('a')
		if runeSet.Length() != 0 {
			t.Errorf("Expected length 0, got %d", runeSet.Length())
		}
	})

	t.Run("EmptyCharSet_RemoveCharacterRange_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Empty()
		runeSet.RemoveRange('a', 'c')
		if runeSet.Length() != 0 {
			t.Errorf("Expected length 0, got %d", runeSet.Length())
		}
	})

	t.Run("OneCharSet_RemoveSingleCharacter_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Rune('a')
		runeSet.RemoveRune('a')
		if runeSet.Length() != 0 {
			t.Errorf("Expected length 0, got %d", runeSet.Length())
		}
	})

	t.Run("OneCharSet_RemoveCharacterRange_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Rune('b')
		runeSet.RemoveRange('a', 'c')
		if runeSet.Length() != 0 {
			t.Errorf("Expected length 0, got %d", runeSet.Length())
		}
	})

	t.Run("ThreeCharsSet_RemoveSingleCharacterMiddle_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Range('a', 'c')
		runeSet.RemoveRune('b')
		if runeSet.Length() != 2 {
			t.Errorf("Expected length 2, got %d", runeSet.Length())
		}
		if !runeSet.IncludesRune('a') {
			t.Errorf("Expected rune 'a' to be included")
		}
		if !runeSet.IncludesRune('c') {
			t.Errorf("Expected rune 'c' to be included")
		}
	})

	t.Run("ThreeCharsSet_RemoveSingleCharacterLeft_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Range('a', 'c')
		runeSet.RemoveRune('a')
		if runeSet.Length() != 2 {
			t.Errorf("Expected length 2, got %d", runeSet.Length())
		}
		if !runeSet.IncludesRune('b') {
			t.Errorf("Expected rune 'b' to be included")
		}
		if !runeSet.IncludesRune('c') {
			t.Errorf("Expected rune 'c' to be included")
		}
	})

	t.Run("ThreeCharsSet_RemoveSingleCharacterRight_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Range('a', 'c')
		runeSet.RemoveRune('c')
		if runeSet.Length() != 2 {
			t.Errorf("Expected length 2, got %d", runeSet.Length())
		}
		if !runeSet.IncludesRune('a') {
			t.Errorf("Expected rune 'a' to be included")
		}
		if !runeSet.IncludesRune('b') {
			t.Errorf("Expected rune 'b' to be included")
		}
	})

	t.Run("ThreeCharsSet_RemoveCharacterRangeExact_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Range('a', 'c')
		runeSet.RemoveRange('a', 'c')
		if runeSet.Length() != 0 {
			t.Errorf("Expected length 0, got %d", runeSet.Length())
		}
	})

	t.Run("ThreeCharsSet_RemoveCharacterRangeOutside_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Range('a', 'c')
		runeSet.RemoveRange('d', 'f')
		if runeSet.Length() != 3 {
			t.Errorf("Expected length 3, got %d", runeSet.Length())
		}
		if !runeSet.IncludesRange('a', 'c') {
			t.Errorf("Expected range 'a'-'c' to be included")
		}
	})

	t.Run("ThreeCharsSet_RemoveCharacterRangeOverlapRight_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Range('a', 'c')
		runeSet.RemoveRange('b', 'd')
		if runeSet.Length() != 1 {
			t.Errorf("Expected length 1, got %d", runeSet.Length())
		}
		if !runeSet.IncludesRune('a') {
			t.Errorf("Expected rune 'a' to be included")
		}
	})

	t.Run("ThreeCharsSet_RemoveCharacterRangeOverlapLeft_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Range('b', 'd')
		runeSet.RemoveRange('a', 'c')
		if runeSet.Length() != 1 {
			t.Errorf("Expected length 1, got %d", runeSet.Length())
		}
		if !runeSet.IncludesRune('d') {
			t.Errorf("Expected rune 'd' to be included")
		}
	})
}

func TestRuneSet_Coverage(t *testing.T) {
	t.Run("CharSet_ConstructorProvidesCharRanges", func(t *testing.T) {
		// Create a RuneSet with included ranges from 0-20 (similar to the constructor test)
		runeSet := &RuneSet{
			Ranges: []RuneRange{
				*NewRuneRange(0, 20, true),
				*NewRuneRange(21, MaxRune, false),
			},
		}
		if !runeSet.IncludesRange(0, 20) {
			t.Errorf("Expected range 0-20 to be included")
		}
		if !runeSet.ExcludesRange(21, MaxRune) {
			t.Errorf("Expected range 21-0xFFFF to be excluded")
		}
		if !runeSet.ExcludesRune(22) {
			t.Errorf("Expected rune 22 to be excluded")
		}
	})

	t.Run("CharSet_UseIterator", func(t *testing.T) {
		// Create a RuneSet similar to the TypeScript test: 10-20 included, 15-18 excluded
		runeSet := &RuneSet{
			Ranges: []RuneRange{
				*NewRuneRange(10, 14, true),
				*NewRuneRange(15, 18, false),
				*NewRuneRange(19, 20, true),
			},
		}

		// Find the included ranges
		var includedRanges []RuneRange
		for _, r := range runeSet.Ranges {
			if r.Includes {
				includedRanges = append(includedRanges, r)
			}
		}

		if len(includedRanges) != 2 {
			t.Errorf("Expected 2 included ranges, got %d", len(includedRanges))
		}

		if len(includedRanges) >= 1 {
			left := includedRanges[0]
			if left.Start != 10 || left.End != 14 {
				t.Errorf("Expected left range 10-14, got %d-%d", left.Start, left.End)
			}
		}

		if len(includedRanges) >= 2 {
			right := includedRanges[1]
			if right.Start != 19 || right.End != 20 {
				t.Errorf("Expected right range 19-20, got %d-%d", right.Start, right.End)
			}
		}
	})

	t.Run("RuneSet_ToString", func(t *testing.T) {
		runeSet := &RuneSet{
			Ranges: []RuneRange{
				*NewRuneRange(10, 14, true),
				*NewRuneRange(15, 18, false),
				*NewRuneRange(19, 20, true),
			},
		}
		expected := "[" + string(rune(10)) + "-" + string(rune(14)) + "],[" + string(rune(19)) + "-" + string(rune(20)) + "]"
		if runeSet.String() != expected {
			t.Errorf("Expected string %s, got %s", expected, runeSet.String())
		}
	})

	t.Run("RuneSet_Equals", func(t *testing.T) {
		runeSet1 := &RuneSet{
			Ranges: []RuneRange{
				*NewRuneRange(10, 14, true),
				*NewRuneRange(15, 18, false),
				*NewRuneRange(19, 20, true),
			},
		}
		runeSet2 := &RuneSet{
			Ranges: []RuneRange{
				*NewRuneRange(10, 14, true),
				*NewRuneRange(15, 18, false),
				*NewRuneRange(19, 20, true),
			},
		}
		if !runeSet1.Equals(*runeSet2) {
			t.Errorf("Expected rune sets to be equal")
		}

		runeSet3 := &RuneSet{
			Ranges: []RuneRange{
				*NewRuneRange(10, 14, true),
				*NewRuneRange(16, 18, false),
				*NewRuneRange(19, 20, true),
			},
		}
		if runeSet1.Equals(*runeSet3) {
			t.Errorf("Expected rune sets to be different")
		}
	})

	t.Run("RuneSet_Equals_DifferentLength_False", func(t *testing.T) {
		runeSet1 := &RuneSet{
			Ranges: []RuneRange{
				*NewRuneRange(10, 14, true),
				*NewRuneRange(15, 18, false),
				*NewRuneRange(19, 20, true),
			},
		}
		runeSet2 := &RuneSet{
			Ranges: []RuneRange{
				*NewRuneRange(10, 14, true),
				*NewRuneRange(19, 21, true),
			},
		}
		if runeSet1.Equals(*runeSet2) {
			t.Errorf("Expected rune sets to be different (different lengths)")
		}
	})

	t.Run("Factories", func(t *testing.T) {
		t.Run("RuneSet_OneOf", func(t *testing.T) {
			chars := []rune{'a', 'b', 'c'}
			runeSet := NewRuneSet_OneOf(chars)
			if runeSet.Length() != 3 {
				t.Errorf("Expected length 3, got %d", runeSet.Length())
			}
			if !runeSet.IncludesRange('a', 'c') {
				t.Errorf("Expected range 'a'-'c' to be included")
			}
		})

		t.Run("RuneSet_Char", func(t *testing.T) {
			runeSet := NewRuneSet_Rune('a')
			if runeSet.Length() != 1 {
				t.Errorf("Expected length 1, got %d", runeSet.Length())
			}
			if !runeSet.IncludesRune('a') {
				t.Errorf("Expected rune 'a' to be included")
			}
		})

		t.Run("RuneSet_Range", func(t *testing.T) {
			runeSet := NewRuneSet_Range('a', 'c')
			if runeSet.Length() != 3 {
				t.Errorf("Expected length 3, got %d", runeSet.Length())
			}
			if !runeSet.IncludesRange('a', 'c') {
				t.Errorf("Expected range 'a'-'c' to be included")
			}
		})

		t.Run("RuneSet_Empty", func(t *testing.T) {
			runeSet := NewRuneSet_Empty()
			if runeSet.Length() != 0 {
				t.Errorf("Expected length 0, got %d", runeSet.Length())
			}
		})

		t.Run("RuneSet_Full", func(t *testing.T) {
			runeSet := NewRuneSet_Full()
			expectedLength := int(MaxRune) + 1
			if runeSet.Length() != expectedLength {
				t.Errorf("Expected length %d, got %d", expectedLength, runeSet.Length())
			}
			if !runeSet.IncludesRange(0, MaxRune) {
				t.Errorf("Expected range 0-0xFFFF to be included")
			}
		})
	})
}
