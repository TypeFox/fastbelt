package benchmark

import (
	original "regexp"
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
