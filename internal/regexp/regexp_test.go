package regexp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	automatons "typefox.dev/fastbelt/internal/automatons"
)

func checkRegexp(t *testing.T, regexp Regexp, input string, expected []int) {
	loc := regexp.FindStringIndex(input)
	if (loc == nil && expected != nil) || (loc != nil && expected == nil) {
		panic("Location mismatch")
	}
	if loc != nil && expected != nil {
		assert.Equal(t, len(loc), len(expected))
		for i := range loc {
			assert.Equal(t, loc[i], expected[i])
		}
	}
}

func TestSimple(t *testing.T) {
	regexp := MustCompile("a+")
	checkRegexp(t, regexp, "aaab", []int{0, 3})
	checkRegexp(t, regexp, "ab", []int{0, 1})
	checkRegexp(t, regexp, "aace", []int{0, 2})
}

func TestEmpty(t *testing.T) {
	regexp := MustCompile("")
	checkRegexp(t, regexp, "", []int{0, 0})
	checkRegexp(t, regexp, "123", []int{0, 0})
}

func TestStar(t *testing.T) {
	regexp := MustCompile("a*")
	checkRegexp(t, regexp, "", []int{0, 0})
	checkRegexp(t, regexp, "a", []int{0, 1})
	checkRegexp(t, regexp, "aaaabbb", []int{0, 4})
}

func TestEmail(t *testing.T) {
	regexp := MustCompile(`[\w\.+-]+@[\w\.-]+\.[\w\.-]+`)
	checkRegexp(t, regexp, "a.b@c.de", []int{0, 8})
}

func TestIP(t *testing.T) {
	regexp := MustCompile(`(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9])\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9])`)
	checkRegexp(t, regexp, "255.241.123.10", []int{0, 14})
}

func TestWhitespace(t *testing.T) {
	regexp := MustCompile(`[ \n\r\t]+`)
	checkRegexp(t, regexp, " ", []int{0, 1})
}

func TestRegexpLiteral(t *testing.T) {
	regexp := MustCompile("/([^\\r\\n\\[\\/\\\\]|\\\\.|\\[([^\\r\\n\\]\\\\]|\\\\.)*\\])+/")
	checkRegexp(t, regexp, "/github.com/", []int{0, 12})
}

func TestStartChars_SingleRune(t *testing.T) {
	regexp := MustCompile("a[bc]d")
	startChars := regexp.(*RegexpImpl).GetStartChars()
	expectedRunes := automatons.NewRuneSetRune('a')
	assert.True(t, expectedRunes.Equals(*startChars))
}

func TestStartChars_RuneRange(t *testing.T) {
	regexp := MustCompile("(a|b|c)")
	startChars := regexp.(*RegexpImpl).GetStartChars()
	expectedRunes := automatons.NewRuneSetOneOf([]rune{'a', 'b', 'c'})
	assert.True(t, expectedRunes.Equals(*startChars))
}

func TestStartChars_RuneNonAscii(t *testing.T) {
	regexp := MustCompile("🔥")
	startChars := regexp.(*RegexpImpl).GetStartChars()
	expectedRunes := automatons.NewRuneSetRune('🔥' % 0x100)
	assert.True(t, expectedRunes.Equals(*startChars))
}

func TestStartChars_RuneNonAsciiBigRange(t *testing.T) {
	regexp := MustCompile("[\U0001F525-\U0001F625]")
	startChars := regexp.(*RegexpImpl).GetStartChars()
	expectedRunes := automatons.NewRuneSetRange(0, 0xFF)
	assert.True(t, expectedRunes.Equals(*startChars))
}
