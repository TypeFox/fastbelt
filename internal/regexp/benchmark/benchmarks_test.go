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
	for n := 0; n < b.N; n++ {
		generated.Email(email, 0)
	}
}

func BenchmarkEmailOriginal(b *testing.B) {
	regexp := original.MustCompile(EmailPattern())
	for n := 0; n < b.N; n++ {
		regexp.FindStringIndex(email)
	}
}

func BenchmarkURLCustom(b *testing.B) {
	for n := 0; n < b.N; n++ {
		generated.URL(url, 0)
	}
}

func BenchmarkURLOriginal(b *testing.B) {
	regexp := original.MustCompile(URLPattern())
	for n := 0; n < b.N; n++ {
		regexp.FindStringIndex(url)
	}
}

func BenchmarkIPv4Custom(b *testing.B) {
	for n := 0; n < b.N; n++ {
		generated.IPv4(ipv4, 0)
	}
}

func BenchmarkIPv4Original(b *testing.B) {
	regexp := original.MustCompile(IPv4Pattern())
	for n := 0; n < b.N; n++ {
		regexp.FindStringIndex(ipv4)
	}
}
