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
		{"Whitespace", WhitespacePattern(), generated.Whitespace, []string{" \t\r\nx", "x", "   ", "\u00a0"}},
		{"BlockComment", BlockCommentPattern(), generated.BlockComment, []string{
			"/* a */ b", "/** x **/", "/* unterminated", "/*/ x */", "/* a * / b */ c", "/* ü */",
		}},
		{"URL", URLPattern(), generated.URL, []string{"http://example.com/a?b=c#d", "https://x"}},
		{"Email", EmailPattern(), generated.Email, []string{"a.b@c.de", "a@b", "ä@b.c"}},
	}
	alphabet := []rune(`ab_1 "'\\*/.@-:` + "\t\n" + "üé€")
	rng := rand.New(rand.NewPCG(1, 2))
	for _, c := range cases {
		re := original.MustCompile(`^(?:` + c.pattern + `)`)
		re.Longest()
		inputs := append([]string{}, c.samples...)
		for range 3000 {
			var b strings.Builder
			for range rng.IntN(24) {
				b.WriteRune(alphabet[rng.IntN(len(alphabet))])
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
