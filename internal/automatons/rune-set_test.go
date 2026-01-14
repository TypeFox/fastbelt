package automatons

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRuneSet_Add(t *testing.T) {
	t.Run("EmptySet_AddSingleCharacter_Success", func(t *testing.T) {
		c := rune(123)
		runeSet := NewRuneSetEmpty()
		runeSet.AddRune(c)
		assert.Equal(t, 1, runeSet.Length())
		assert.Equal(t, true, runeSet.IncludesRune(c))
	})

	t.Run("EmptySet_AddCharacterRange_Success", func(t *testing.T) {
		cFrom := rune(123)
		cTo := rune(125)
		runeSet := NewRuneSetEmpty()
		runeSet.AddRange(cFrom, cTo)
		assert.Equal(t, 3, runeSet.Length())
		assert.Equal(t, true, runeSet.IncludesRange(cFrom, cTo))
		assert.Equal(t, true, runeSet.IncludesRune(cFrom))
		assert.Equal(t, true, runeSet.IncludesRune(cFrom+1))
		assert.Equal(t, true, runeSet.IncludesRune(cFrom+2))
	})

	t.Run("EmptySet_AddBadCharacterRange_Fail", func(t *testing.T) {
		assert.Panics(t, func() {
			runeSet := NewRuneSetEmpty()
			runeSet.AddRange(125, 123)
		})
	})

	t.Run("OneCharSet_AddSingleCharacterSame_Success", func(t *testing.T) {
		c := rune(123)
		runeSet := NewRuneSetEmpty()
		runeSet.AddRune(c)
		runeSet.AddRune(c)
		assert.Equal(t, 1, runeSet.Length())
		assert.Equal(t, true, runeSet.IncludesRune(c))
	})

	t.Run("OneCharSet_AddSingleCharacterBeside_Success", func(t *testing.T) {
		c1 := rune(123)
		c2 := rune(124)
		runeSet := NewRuneSetEmpty()
		runeSet.AddRune(c1)
		runeSet.AddRune(c2)
		assert.Equal(t, 2, runeSet.Length())
		assert.Equal(t, true, runeSet.IncludesRange(c1, c2))
	})

	t.Run("OneCharSet_AddSingleCharacterWithGap_Success", func(t *testing.T) {
		c1 := rune(123)
		c2 := rune(125)
		runeSet := NewRuneSetEmpty()
		runeSet.AddRune(c1)
		runeSet.AddRune(c2)
		assert.Equal(t, 2, runeSet.Length())
		assert.Equal(t, false, runeSet.IncludesRange(c1, c2))
		assert.Equal(t, true, runeSet.IncludesRune(c1))
		assert.Equal(t, true, runeSet.IncludesRune(c2))
	})

	t.Run("OneCharSet_AddCharacterRangeWithin_Success", func(t *testing.T) {
		cFrom1 := rune(1)
		cTo1 := rune(3)
		cFrom2 := rune(4)
		cTo2 := rune(6)
		runeSet := NewRuneSetEmpty()
		runeSet.AddRange(cFrom1, cTo1)
		runeSet.AddRange(cFrom2, cTo2)
		assert.Equal(t, 6, runeSet.Length())
		assert.Equal(t, true, runeSet.IncludesRange(cFrom1, cTo1))
	})

	t.Run("OneCharSet_AddCharacterRangeOverlapping_Success", func(t *testing.T) {
		cFrom1 := rune(1)
		cTo1 := rune(3)
		cFrom2 := rune(5)
		cTo2 := rune(7)
		runeSet := NewRuneSetEmpty()
		runeSet.AddRange(cFrom1, cTo1)
		runeSet.AddRange(cFrom2, cTo2)
		assert.Equal(t, 6, runeSet.Length())
		assert.Equal(t, true, runeSet.IncludesRange(cFrom1, cTo1))
		assert.Equal(t, true, runeSet.IncludesRange(cFrom2, cTo2))
	})

	t.Run("ThreeCharsSet_AddSingleCharacterAlreadyExists_Success", func(t *testing.T) {
		runeSet := NewRuneSetEmpty()
		runeSet.AddRange('a', 'c')
		runeSet.AddRune('b')
		assert.Equal(t, 3, runeSet.Length())
		assert.Equal(t, true, runeSet.IncludesRange('a', 'c'))
	})

	t.Run("ThreeCharsSet_AddSingleCharacterBeside_Success", func(t *testing.T) {
		runeSet := NewRuneSetEmpty()
		runeSet.AddRange('a', 'c')
		runeSet.AddRune('d')
		assert.Equal(t, 4, runeSet.Length())
		assert.Equal(t, true, runeSet.IncludesRange('a', 'd'))
	})

	t.Run("ThreeCharsSet_AddSingleCharacterWithGap_Success", func(t *testing.T) {
		runeSet := NewRuneSetEmpty()
		runeSet.AddRange('a', 'c')
		runeSet.AddRune('e')
		assert.Equal(t, 4, runeSet.Length())
		assert.Equal(t, false, runeSet.IncludesRange('a', 'e'))
		assert.Equal(t, true, runeSet.IncludesRune('e'))
	})

	t.Run("ThreeCharsSet_AddCharacterRangeWithin_Success", func(t *testing.T) {
		runeSet := NewRuneSetEmpty()
		runeSet.AddRange('a', 'c')
		runeSet.AddRange('b', 'c')
		assert.Equal(t, 3, runeSet.Length())
		assert.Equal(t, true, runeSet.IncludesRange('a', 'c'))
	})

	t.Run("ThreeCharsSet_AddCharacterRangeOverlapping_Success", func(t *testing.T) {
		runeSet := NewRuneSetEmpty()
		runeSet.AddRange('a', 'c')
		runeSet.AddRange('b', 'd')
		assert.Equal(t, 4, runeSet.Length())
		assert.Equal(t, true, runeSet.IncludesRange('a', 'd'))
	})

	t.Run("ThreeCharsSet_AddCharacterRangeBadRange_Fail", func(t *testing.T) {
		assert.Panics(t, func() {
			runeSet := NewRuneSetEmpty()
			runeSet.AddRange('a', 'c')
			runeSet.AddRange('d', 'b')
		})
	})

	t.Run("ThreeCharsSet_AddCharacterRangeOverlappingBothSides_Success", func(t *testing.T) {
		runeSet := NewRuneSetEmpty()
		runeSet.AddRange('d', 'f')
		runeSet.AddRange('b', 'h')
		assert.Equal(t, 7, runeSet.Length())
		assert.Equal(t, true, runeSet.IncludesRange('b', 'h'))
	})

	t.Run("ThreeCharsSet_AddCharacterRangeCoveringAll_Success", func(t *testing.T) {
		runeSet := NewRuneSetEmpty()
		runeSet.AddRange('d', 'f')
		runeSet.AddRange('a', 'z')
		assert.Equal(t, 26, runeSet.Length())
		assert.Equal(t, true, runeSet.IncludesRange('a', 'z'))
	})

	t.Run("ThreeCharsSet_AddCharacterRangeBeside_Success", func(t *testing.T) {
		runeSet := NewRuneSetEmpty()
		runeSet.AddRange('a', 'c')
		runeSet.AddRange('d', 'f')
		assert.Equal(t, 6, runeSet.Length())
		assert.Equal(t, true, runeSet.IncludesRange('a', 'f'))
	})

	t.Run("ThreeCharsSet_AddCharacterRangeWithGap_Success", func(t *testing.T) {
		runeSet := NewRuneSetEmpty()
		runeSet.AddRange('a', 'c')
		runeSet.AddRange('e', 'g')
		assert.Equal(t, 6, runeSet.Length())
		assert.Equal(t, false, runeSet.IncludesRange('a', 'g'))
		assert.Equal(t, true, runeSet.IncludesRange('e', 'g'))
	})

	t.Run("TwoThreeGroupsCharsSet_AddCharacterFillGap_Success", func(t *testing.T) {
		runeSet := NewRuneSetEmpty()
		runeSet.AddRange('a', 'c')
		runeSet.AddRange('e', 'g')
		runeSet.AddRune('d')
		assert.Equal(t, 7, runeSet.Length())
		assert.Equal(t, true, runeSet.IncludesRange('a', 'g'))
	})

	t.Run("TwoThreeGroupsCharsSet_AddCharacterRangeFillGap_Success", func(t *testing.T) {
		runeSet := NewRuneSetEmpty()
		runeSet.AddRange('a', 'c')
		runeSet.AddRange('g', 'i')
		runeSet.AddRange('d', 'f')
		assert.Equal(t, 9, runeSet.Length())
		assert.Equal(t, true, runeSet.IncludesRange('a', 'i'))
	})
}

func TestRuneSet_Remove(t *testing.T) {
	t.Run("EmptyCharSet_RemoveSingleCharacter_Success", func(t *testing.T) {
		runeSet := NewRuneSetEmpty()
		runeSet.RemoveRune('a')
		assert.Equal(t, 0, runeSet.Length())
	})

	t.Run("EmptyCharSet_RemoveCharacterRange_Success", func(t *testing.T) {
		runeSet := NewRuneSetEmpty()
		runeSet.RemoveRange('a', 'c')
		assert.Equal(t, 0, runeSet.Length())
	})

	t.Run("OneCharSet_RemoveSingleCharacter_Success", func(t *testing.T) {
		runeSet := NewRuneSetRune('a')
		runeSet.RemoveRune('a')
		assert.Equal(t, 0, runeSet.Length())
	})

	t.Run("OneCharSet_RemoveCharacterRange_Success", func(t *testing.T) {
		runeSet := NewRuneSetRune('b')
		runeSet.RemoveRange('a', 'c')
		assert.Equal(t, 0, runeSet.Length())
	})

	t.Run("ThreeCharsSet_RemoveSingleCharacterMiddle_Success", func(t *testing.T) {
		runeSet := NewRuneSetRange('a', 'c')
		runeSet.RemoveRune('b')
		assert.Equal(t, 2, runeSet.Length())
		assert.Equal(t, true, runeSet.IncludesRune('a'))
		assert.Equal(t, true, runeSet.IncludesRune('c'))
	})

	t.Run("ThreeCharsSet_RemoveSingleCharacterLeft_Success", func(t *testing.T) {
		runeSet := NewRuneSetRange('a', 'c')
		runeSet.RemoveRune('a')
		assert.Equal(t, 2, runeSet.Length())
		assert.Equal(t, true, runeSet.IncludesRune('b'))
		assert.Equal(t, true, runeSet.IncludesRune('c'))
	})

	t.Run("ThreeCharsSet_RemoveSingleCharacterRight_Success", func(t *testing.T) {
		runeSet := NewRuneSetRange('a', 'c')
		runeSet.RemoveRune('c')
		assert.Equal(t, 2, runeSet.Length())
		assert.Equal(t, true, runeSet.IncludesRune('a'))
		assert.Equal(t, true, runeSet.IncludesRune('b'))
	})

	t.Run("ThreeCharsSet_RemoveCharacterRangeExact_Success", func(t *testing.T) {
		runeSet := NewRuneSetRange('a', 'c')
		runeSet.RemoveRange('a', 'c')
		assert.Equal(t, 0, runeSet.Length())
	})

	t.Run("ThreeCharsSet_RemoveCharacterRangeOutside_Success", func(t *testing.T) {
		runeSet := NewRuneSetRange('a', 'c')
		runeSet.RemoveRange('d', 'f')
		assert.Equal(t, 3, runeSet.Length())
		assert.Equal(t, true, runeSet.IncludesRange('a', 'c'))
	})

	t.Run("ThreeCharsSet_RemoveCharacterRangeOverlapRight_Success", func(t *testing.T) {
		runeSet := NewRuneSetRange('a', 'c')
		runeSet.RemoveRange('b', 'd')
		assert.Equal(t, 1, runeSet.Length())
		assert.Equal(t, true, runeSet.IncludesRune('a'))
	})

	t.Run("ThreeCharsSet_RemoveCharacterRangeOverlapLeft_Success", func(t *testing.T) {
		runeSet := NewRuneSetRange('b', 'd')
		runeSet.RemoveRange('a', 'c')
		assert.Equal(t, 1, runeSet.Length())
		assert.Equal(t, true, runeSet.IncludesRune('d'))
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
		assert.Equal(t, 21, runeSet.Length())
		assert.Equal(t, true, runeSet.IncludesRange(0, 20))
		assert.Equal(t, true, runeSet.ExcludesRange(21, MaxRune))
		assert.Equal(t, true, runeSet.ExcludesRune(22))
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

		assert.Equal(t, 2, len(includedRanges))

		if len(includedRanges) >= 1 {
			left := includedRanges[0]
			assert.Equal(t, rune(10), left.Start)
			assert.Equal(t, rune(14), left.End)
		}

		if len(includedRanges) >= 2 {
			right := includedRanges[1]
			assert.Equal(t, rune(19), right.Start)
			assert.Equal(t, rune(20), right.End)
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
		assert.Equal(t, expected, runeSet.String())
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

		assert.Equal(t, true, runeSet1.Equals(*runeSet2))

		runeSet3 := &RuneSet{
			Ranges: []RuneRange{
				*NewRuneRange(10, 14, true),
				*NewRuneRange(16, 18, false),
				*NewRuneRange(19, 20, true),
			},
		}
		assert.Equal(t, false, runeSet1.Equals(*runeSet3))
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
		assert.Equal(t, false, runeSet1.Equals(*runeSet2))
	})

	t.Run("Factories", func(t *testing.T) {
		t.Run("RuneSet_OneOf", func(t *testing.T) {
			chars := []rune{'a', 'b', 'c'}
			runeSet := NewRuneSetOneOf(chars)
			assert.Equal(t, 3, runeSet.Length())
			assert.Equal(t, true, runeSet.IncludesRange('a', 'c'))
		})

		t.Run("RuneSet_Char", func(t *testing.T) {
			runeSet := NewRuneSetRune('a')
			assert.Equal(t, 1, runeSet.Length())
			assert.Equal(t, true, runeSet.IncludesRune('a'))
		})

		t.Run("RuneSet_Range", func(t *testing.T) {
			runeSet := NewRuneSetRange('a', 'c')
			assert.Equal(t, 3, runeSet.Length())
			assert.Equal(t, true, runeSet.IncludesRange('a', 'c'))
		})

		t.Run("RuneSet_Empty", func(t *testing.T) {
			runeSet := NewRuneSetEmpty()
			assert.Equal(t, 0, runeSet.Length())
		})

		t.Run("RuneSet_Full", func(t *testing.T) {
			runeSet := NewRuneSetFull()
			expectedLength := int(MaxRune) + 1
			assert.Equal(t, expectedLength, runeSet.Length())
			assert.Equal(t, true, runeSet.IncludesRange(0, MaxRune))
		})
	})
}
