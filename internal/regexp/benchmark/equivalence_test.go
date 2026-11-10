package benchmark

import (
	"math/rand/v2"
	original "regexp"
	"strings"
	"testing"

	generated "typefox.dev/fastbelt/internal/regexp/benchmark/generated"
)

// TestGeneratedMatchesGoRegexp checks the generated matchers against Go's
// regexp package (anchored, leftmost-longest) on hand-picked and random inputs.
func TestGeneratedMatchesGoRegexp(t *testing.T) {
	cases := []struct {
		name    string
		pattern string
		match   func(string, int) int
		samples []string
	}{
		{"String", StringPattern(), generated.String, []string{
			`"abc"`, `'abc'`, `"a\"b"`, `'it\'s'`, `"unterminated`, `"a\\"b"`, `"ünï\ncode"`, `""`, `"\`, `"x"y"`,
		}},
		{"ID", IDPattern(), generated.ID, []string{"abc", "_a1 b", "a-b", "1abc", "ab€c", "x"}},
		{"UnicodeID", UnicodeIDPattern(), generated.UnicodeID, []string{
			"abcädef", "äbc", "abcä", "ab€c", "日本語abc", "a𝒜b c", "ab\xffcd", "a\u00a0b", "_1é2",
		}},
		{"Whitespace", WhitespacePattern(), generated.Whitespace, []string{" \t\r\nx", "x", "   ", "\u00a0"}},
		{"BlockComment", BlockCommentPattern(), generated.BlockComment, []string{
			"/* a */ b", "/** x **/", "/* unterminated", "/*/ x */", "/* a * / b */ c", "/* ü */",
		}},
		{"URL", URLPattern(), generated.URL, []string{"http://example.com/a?b=c#d", "https://x"}},
		{"Email", EmailPattern(), generated.Email, []string{"a.b@c.de", "a@b", "ä@b.c"}},
	}
	// ASCII bytes of each fast path and its exits, plus 2-, 3- and 4-byte
	// runes (letters and non-letters), and an invalid UTF-8 byte below.
	alphabet := []string{"a", "b", "_", "1", " ", `"`, "'", `\`, "*", "/", ".", "@", "-", ":", "\t", "\n",
		"ü", "é", "€", "日", "𝒜", "😀", "\u00a0", "\xff"}
	rng := rand.New(rand.NewPCG(1, 2))
	for _, c := range cases {
		re := original.MustCompile(`^(?:` + c.pattern + `)`)
		re.Longest()
		inputs := append([]string{}, c.samples...)
		for range 3000 {
			var b strings.Builder
			for range rng.IntN(24) {
				b.WriteString(alphabet[rng.IntN(len(alphabet))])
			}
			inputs = append(inputs, b.String())
		}
		for _, in := range inputs {
			want := 0
			if loc := re.FindStringIndex(in); loc != nil {
				want = loc[1]
			}
			if got := c.match(in, 0); got != want {
				t.Errorf("%s(%q) = %d, want %d", c.name, in, got, want)
			}
		}
	}
}

// TestSelfLoopAcrossNonASCII covers identifiers and strings in which the
// ASCII fast path is interrupted by multi-byte runes and resumes after them.
func TestSelfLoopAcrossNonASCII(t *testing.T) {
	cases := []struct {
		name  string
		match func(string, int) int
		input string
		want  int
	}{
		{"UnicodeID", generated.UnicodeID, "abcädef", len("abcädef")},
		{"UnicodeID", generated.UnicodeID, "abcädef ghi", len("abcädef")},
		{"UnicodeID", generated.UnicodeID, "äöü", len("äöü")},
		{"UnicodeID", generated.UnicodeID, "x日本語y1", len("x日本語y1")},
		{"UnicodeID", generated.UnicodeID, "a𝒜b", len("a𝒜b")},     // 4-byte letter
		{"UnicodeID", generated.UnicodeID, "ab€cd", len("ab")},    // € is no letter
		{"UnicodeID", generated.UnicodeID, "ab\xffcd", len("ab")}, // invalid UTF-8
		{"UnicodeID", generated.UnicodeID, "€abc", 0},
		{"ID", generated.ID, "abcädef", len("abc")}, // \w is ASCII-only
		{"String", generated.String, `"abcädef"x`, len(`"abcädef"`)},
		{"String", generated.String, `"日本\"語😀"`, len(`"日本\"語😀"`)},
		{"String", generated.String, "\"a\xffb\"", len("\"a\xffb\"")},
		{"String", generated.String, `"abcäd`, 0},
		{"BlockComment", generated.BlockComment, "/* ä * ö */x", len("/* ä * ö */")},
	}
	for _, c := range cases {
		if got := c.match(c.input, 0); got != c.want {
			t.Errorf("%s(%q) = %d, want %d", c.name, c.input, got, c.want)
		}
		// Also at an offset, as the lexer calls it.
		if got := c.match("..."+c.input, 3); got != c.want {
			t.Errorf("%s(%q, offset 3) = %d, want %d", c.name, "..."+c.input, got, c.want)
		}
	}
}
