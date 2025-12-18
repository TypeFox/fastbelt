package automatons

import (
	"testing"
)

func TestRuneSet_Add(t *testing.T) {
	t.Run("EmptySet_AddSingleCharacter_Success", func(t *testing.T) {
		c := rune(123)
		runeSet := NewRuneSet_Empty()
		runeSet.AddRune(c)
		Expect(runeSet.Length()).ToEqual(1)
		Expect(runeSet.IncludesRune(c)).ToEqual(true)
	})

	t.Run("EmptySet_AddCharacterRange_Success", func(t *testing.T) {
		cFrom := rune(123)
		cTo := rune(125)
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange(cFrom, cTo)
		Expect(runeSet.Length()).ToEqual(3)
		Expect(runeSet.IncludesRange(cFrom, cTo)).ToEqual(true)
		Expect(runeSet.IncludesRune(cFrom)).ToEqual(true)
		Expect(runeSet.IncludesRune(cFrom + 1)).ToEqual(true)
		Expect(runeSet.IncludesRune(cFrom + 2)).ToEqual(true)
	})

	t.Run("EmptySet_AddBadCharacterRange_Fail", func(t *testing.T) {
		Expect(func() {
			runeSet := NewRuneSet_Empty()
			runeSet.AddRange(125, 123)
		}).ToPanic()
	})

	t.Run("OneCharSet_AddSingleCharacterSame_Success", func(t *testing.T) {
		c := rune(123)
		runeSet := NewRuneSet_Empty()
		runeSet.AddRune(c)
		runeSet.AddRune(c)
		Expect(runeSet.Length()).ToEqual(1)
		Expect(runeSet.IncludesRune(c)).ToEqual(true)
	})

	t.Run("OneCharSet_AddSingleCharacterBeside_Success", func(t *testing.T) {
		c1 := rune(123)
		c2 := rune(124)
		runeSet := NewRuneSet_Empty()
		runeSet.AddRune(c1)
		runeSet.AddRune(c2)
		Expect(runeSet.Length()).ToEqual(2)
		Expect(runeSet.IncludesRange(c1, c2)).ToEqual(true)
	})

	t.Run("OneCharSet_AddSingleCharacterWithGap_Success", func(t *testing.T) {
		c1 := rune(123)
		c2 := rune(125)
		runeSet := NewRuneSet_Empty()
		runeSet.AddRune(c1)
		runeSet.AddRune(c2)
		Expect(runeSet.Length()).ToEqual(2)
		Expect(runeSet.IncludesRange(c1, c2)).ToEqual(false)
		Expect(runeSet.IncludesRune(c1)).ToEqual(true)
		Expect(runeSet.IncludesRune(c2)).ToEqual(true)
	})

	t.Run("OneCharSet_AddCharacterRangeWithin_Success", func(t *testing.T) {
		cFrom1 := rune(1)
		cTo1 := rune(3)
		cFrom2 := rune(4)
		cTo2 := rune(6)
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange(cFrom1, cTo1)
		runeSet.AddRange(cFrom2, cTo2)
		Expect(runeSet.Length()).ToEqual(6)
		Expect(runeSet.IncludesRange(cFrom1, cTo1)).ToEqual(true)
	})

	t.Run("OneCharSet_AddCharacterRangeOverlapping_Success", func(t *testing.T) {
		cFrom1 := rune(1)
		cTo1 := rune(3)
		cFrom2 := rune(5)
		cTo2 := rune(7)
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange(cFrom1, cTo1)
		runeSet.AddRange(cFrom2, cTo2)
		Expect(runeSet.Length()).ToEqual(6)
		Expect(runeSet.IncludesRange(cFrom1, cTo1)).ToEqual(true)
		Expect(runeSet.IncludesRange(cFrom2, cTo2)).ToEqual(true)
	})

	t.Run("ThreeCharsSet_AddSingleCharacterAlreadyExists_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange('a', 'c')
		runeSet.AddRune('b')
		Expect(runeSet.Length()).ToEqual(3)
		Expect(runeSet.IncludesRange('a', 'c')).ToEqual(true)
	})

	t.Run("ThreeCharsSet_AddSingleCharacterBeside_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange('a', 'c')
		runeSet.AddRune('d')
		Expect(runeSet.Length()).ToEqual(4)
		Expect(runeSet.IncludesRange('a', 'd')).ToEqual(true)
	})

	t.Run("ThreeCharsSet_AddSingleCharacterWithGap_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange('a', 'c')
		runeSet.AddRune('e')
		Expect(runeSet.Length()).ToEqual(4)
		Expect(runeSet.IncludesRange('a', 'e')).ToEqual(false)
		Expect(runeSet.IncludesRune('e')).ToEqual(true)
	})

	t.Run("ThreeCharsSet_AddCharacterRangeWithin_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange('a', 'c')
		runeSet.AddRange('b', 'c')
		Expect(runeSet.Length()).ToEqual(3)
		Expect(runeSet.IncludesRange('a', 'c')).ToEqual(true)
	})

	t.Run("ThreeCharsSet_AddCharacterRangeOverlapping_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange('a', 'c')
		runeSet.AddRange('b', 'd')
		Expect(runeSet.Length()).ToEqual(4)
		Expect(runeSet.IncludesRange('a', 'd')).ToEqual(true)
	})

	t.Run("ThreeCharsSet_AddCharacterRangeBadRange_Fail", func(t *testing.T) {
		Expect(func() {
			runeSet := NewRuneSet_Empty()
			runeSet.AddRange('a', 'c')
			runeSet.AddRange('d', 'b')
		}).ToPanic()
	})

	t.Run("ThreeCharsSet_AddCharacterRangeOverlappingBothSides_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange('d', 'f')
		runeSet.AddRange('b', 'h')
		Expect(runeSet.Length()).ToEqual(7)
		Expect(runeSet.IncludesRange('b', 'h')).ToEqual(true)
	})

	t.Run("ThreeCharsSet_AddCharacterRangeCoveringAll_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange('d', 'f')
		runeSet.AddRange('a', 'z')
		Expect(runeSet.Length()).ToEqual(26)
		Expect(runeSet.IncludesRange('a', 'z')).ToEqual(true)
	})

	t.Run("ThreeCharsSet_AddCharacterRangeBeside_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange('a', 'c')
		runeSet.AddRange('d', 'f')
		Expect(runeSet.Length()).ToEqual(6)
		Expect(runeSet.IncludesRange('a', 'f')).ToEqual(true)
	})

	t.Run("ThreeCharsSet_AddCharacterRangeWithGap_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange('a', 'c')
		runeSet.AddRange('e', 'g')
		Expect(runeSet.Length()).ToEqual(6)
		Expect(runeSet.IncludesRange('a', 'g')).ToEqual(false)
		Expect(runeSet.IncludesRange('e', 'g')).ToEqual(true)
	})

	t.Run("TwoThreeGroupsCharsSet_AddCharacterFillGap_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange('a', 'c')
		runeSet.AddRange('e', 'g')
		runeSet.AddRune('d')
		Expect(runeSet.Length()).ToEqual(7)
		Expect(runeSet.IncludesRange('a', 'g')).ToEqual(true)
	})

	t.Run("TwoThreeGroupsCharsSet_AddCharacterRangeFillGap_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Empty()
		runeSet.AddRange('a', 'c')
		runeSet.AddRange('g', 'i')
		runeSet.AddRange('d', 'f')
		Expect(runeSet.Length()).ToEqual(9)
		Expect(runeSet.IncludesRange('a', 'i')).ToEqual(true)
	})
}

func TestRuneSet_Remove(t *testing.T) {
	t.Run("EmptyCharSet_RemoveSingleCharacter_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Empty()
		runeSet.RemoveRune('a')
		Expect(runeSet.Length()).ToEqual(0)
	})

	t.Run("EmptyCharSet_RemoveCharacterRange_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Empty()
		runeSet.RemoveRange('a', 'c')
		Expect(runeSet.Length()).ToEqual(0)
	})

	t.Run("OneCharSet_RemoveSingleCharacter_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Rune('a')
		runeSet.RemoveRune('a')
		Expect(runeSet.Length()).ToEqual(0)
	})

	t.Run("OneCharSet_RemoveCharacterRange_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Rune('b')
		runeSet.RemoveRange('a', 'c')
		Expect(runeSet.Length()).ToEqual(0)
	})

	t.Run("ThreeCharsSet_RemoveSingleCharacterMiddle_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Range('a', 'c')
		runeSet.RemoveRune('b')
		Expect(runeSet.Length()).ToEqual(2)
		Expect(runeSet.IncludesRune('a')).ToEqual(true)
		Expect(runeSet.IncludesRune('c')).ToEqual(true)
	})

	t.Run("ThreeCharsSet_RemoveSingleCharacterLeft_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Range('a', 'c')
		runeSet.RemoveRune('a')
		Expect(runeSet.Length()).ToEqual(2)
		Expect(runeSet.IncludesRune('b')).ToEqual(true)
		Expect(runeSet.IncludesRune('c')).ToEqual(true)
	})

	t.Run("ThreeCharsSet_RemoveSingleCharacterRight_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Range('a', 'c')
		runeSet.RemoveRune('c')
		Expect(runeSet.Length()).ToEqual(2)
		Expect(runeSet.IncludesRune('a')).ToEqual(true)
		Expect(runeSet.IncludesRune('b')).ToEqual(true)
	})

	t.Run("ThreeCharsSet_RemoveCharacterRangeExact_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Range('a', 'c')
		runeSet.RemoveRange('a', 'c')
		Expect(runeSet.Length()).ToEqual(0)
	})

	t.Run("ThreeCharsSet_RemoveCharacterRangeOutside_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Range('a', 'c')
		runeSet.RemoveRange('d', 'f')
		Expect(runeSet.Length()).ToEqual(3)
		Expect(runeSet.IncludesRange('a', 'c')).ToEqual(true)
	})

	t.Run("ThreeCharsSet_RemoveCharacterRangeOverlapRight_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Range('a', 'c')
		runeSet.RemoveRange('b', 'd')
		Expect(runeSet.Length()).ToEqual(1)
		Expect(runeSet.IncludesRune('a')).ToEqual(true)
	})

	t.Run("ThreeCharsSet_RemoveCharacterRangeOverlapLeft_Success", func(t *testing.T) {
		runeSet := NewRuneSet_Range('b', 'd')
		runeSet.RemoveRange('a', 'c')
		Expect(runeSet.Length()).ToEqual(1)
		Expect(runeSet.IncludesRune('d')).ToEqual(true)
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
		Expect(runeSet.Length()).ToEqual(21)
		Expect(runeSet.IncludesRange(0, 20)).ToEqual(true)
		Expect(runeSet.ExcludesRange(21, MaxRune)).ToEqual(true)
		Expect(runeSet.ExcludesRune(22)).ToEqual(true)
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

		Expect(len(includedRanges)).ToEqual(2)

		if len(includedRanges) >= 1 {
			left := includedRanges[0]
			Expect(left.Start).ToEqual(rune(10))
			Expect(left.End).ToEqual(rune(14))
		}

		if len(includedRanges) >= 2 {
			right := includedRanges[1]
			Expect(right.Start).ToEqual(rune(19))
			Expect(right.End).ToEqual(rune(20))
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
		Expect(runeSet.String()).ToEqual(expected)
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

		Expect(runeSet1.Equals(*runeSet2)).ToEqual(true)

		runeSet3 := &RuneSet{
			Ranges: []RuneRange{
				*NewRuneRange(10, 14, true),
				*NewRuneRange(16, 18, false),
				*NewRuneRange(19, 20, true),
			},
		}
		Expect(runeSet1.Equals(*runeSet3)).ToEqual(false)
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
		Expect(runeSet1.Equals(*runeSet2)).ToEqual(false)
	})

	t.Run("Factories", func(t *testing.T) {
		t.Run("RuneSet_OneOf", func(t *testing.T) {
			chars := []rune{'a', 'b', 'c'}
			runeSet := NewRuneSet_OneOf(chars)
			Expect(runeSet.Length()).ToEqual(3)
			Expect(runeSet.IncludesRange('a', 'c')).ToEqual(true)
		})

		t.Run("RuneSet_Char", func(t *testing.T) {
			runeSet := NewRuneSet_Rune('a')
			Expect(runeSet.Length()).ToEqual(1)
			Expect(runeSet.IncludesRune('a')).ToEqual(true)
		})

		t.Run("RuneSet_Range", func(t *testing.T) {
			runeSet := NewRuneSet_Range('a', 'c')
			Expect(runeSet.Length()).ToEqual(3)
			Expect(runeSet.IncludesRange('a', 'c')).ToEqual(true)
		})

		t.Run("RuneSet_Empty", func(t *testing.T) {
			runeSet := NewRuneSet_Empty()
			Expect(runeSet.Length()).ToEqual(0)
		})

		t.Run("RuneSet_Full", func(t *testing.T) {
			runeSet := NewRuneSet_Full()
			expectedLength := int(MaxRune) + 1
			Expect(runeSet.Length()).ToEqual(expectedLength)
			Expect(runeSet.IncludesRange(0, MaxRune)).ToEqual(true)
		})
	})
}
