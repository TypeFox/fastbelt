package benchmark

import (
	original "regexp"
	"strings"
	"testing"

	generated "typefox.dev/fastbelt/internal/regexp/benchmark/generated"
)

const email = "abc.def@ghi.jkl"
const url = "https://www.example.com/path?query=string#fragment"
const ipv4 = "255.123.2.1"

func BenchmarkEmailCustom(b *testing.B) {
	for b.Loop() {
		generated.Email(email, 0)
	}
}

func BenchmarkEmailOriginal(b *testing.B) {
	regexp := original.MustCompile(EmailPattern())
	for b.Loop() {
		regexp.FindStringIndex(email)
	}
}

func BenchmarkURLCustom(b *testing.B) {
	for b.Loop() {
		generated.URL(url, 0)
	}
}

func BenchmarkURLOriginal(b *testing.B) {
	regexp := original.MustCompile(URLPattern())
	for b.Loop() {
		regexp.FindStringIndex(url)
	}
}

func BenchmarkIPv4Custom(b *testing.B) {
	for b.Loop() {
		generated.IPv4(ipv4, 0)
	}
}

func BenchmarkIPv4Original(b *testing.B) {
	regexp := original.MustCompile(IPv4Pattern())
	for b.Loop() {
		regexp.FindStringIndex(ipv4)
	}
}

var (
	longString  = `"` + strings.Repeat(`some text with an \"escaped\" quote, `, 8) + `"`
	longID      = strings.Repeat("identifier_", 8)
	longSpace   = strings.Repeat(" ", 64) + "\n"
	longComment = "/*" + strings.Repeat(" a comment line\n *", 8) + "/"
)

func BenchmarkStringCustom(b *testing.B) {
	for b.Loop() {
		generated.String(longString, 0)
	}
}

func BenchmarkIDCustom(b *testing.B) {
	for b.Loop() {
		generated.ID(longID, 0)
	}
}

func BenchmarkWhitespaceCustom(b *testing.B) {
	for b.Loop() {
		generated.Whitespace(longSpace, 0)
	}
}

func BenchmarkBlockCommentCustom(b *testing.B) {
	for b.Loop() {
		generated.BlockComment(longComment, 0)
	}
}
