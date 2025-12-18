package regexp

import (
	"testing"

	"typefox.dev/fastbelt/internal/automatons"
)

func checkRegexp(regexp Regexp, input string, expected []int) {
	loc := regexp.FindStringIndex(input)
	if (loc == nil && expected != nil) || (loc != nil && expected == nil) {
		panic("Location mismatch")
	}
	if loc != nil && expected != nil {
		automatons.Expect(len(loc)).ToEqual(len(expected))
		for i := range loc {
			automatons.Expect(loc[i]).ToEqual(expected[i])
		}
	}
}

func TestSimple(t *testing.T) {
	regexp := MustCompile("a+")
	checkRegexp(regexp, "aaab", []int{0, 3})
	checkRegexp(regexp, "ab", []int{0, 1})
	checkRegexp(regexp, "aace", []int{0, 2})
}

func TestEmpty(t *testing.T) {
	regexp := MustCompile("")
	checkRegexp(regexp, "", []int{0, 0})
	checkRegexp(regexp, "123", []int{0, 0})
}

func TestStar(t *testing.T) {
	regexp := MustCompile("a*")
	checkRegexp(regexp, "", []int{0, 0})
	checkRegexp(regexp, "a", []int{0, 1})
	checkRegexp(regexp, "aaaabbb", []int{0, 4})
}

func TestEmail(t *testing.T) {
	regexp := MustCompile(`[\w\.+-]+@[\w\.-]+\.[\w\.-]+`)
	checkRegexp(regexp, "a.b@c.de", []int{0, 8})
}

func TestIP(t *testing.T) {
	regexp := MustCompile(`(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9])\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9])`)
	checkRegexp(regexp, "255.241.123.10", []int{0, 14})
}

func TestWhitespace(t *testing.T) {
	regexp := MustCompile(`[ \n\r\t]+`)
	checkRegexp(regexp, " ", []int{0, 1})
}

func TestRegexpLiteral(t *testing.T) {
	regexp := MustCompile("/([^\\r\\n\\[\\/\\\\]|\\\\.|\\[([^\\r\\n\\]\\\\]|\\\\.)*\\])+/")
	checkRegexp(regexp, "/github.com/", []int{0, 12})
}
