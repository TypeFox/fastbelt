package benchmark

import (
	"testing"

	"github.com/stretchr/testify/assert"
	generated "typefox.dev/fastbelt/internal/regexp/benchmark/generated"
)

func TestURLPattern(t *testing.T) {
	assert.Equal(t, 18, generated.URL("http://example.com", 0))
	assert.Equal(t, 32, generated.URL("http://example.com/sub/page.html", 0))
	assert.Equal(t, 18, generated.URL("https://vercel.app", 0))
	assert.Equal(t, -1, generated.URL("möp", 0))
}

func TestEmailPattern(t *testing.T) {
	assert.Equal(t, 8, generated.Email("a.b@c.de", 0))
	assert.Equal(t, 15, generated.Email("abc.def@ghi.jkl", 0))
	assert.Equal(t, -1, generated.Email("möp", 0))
}

func TestIPv4Pattern(t *testing.T) {
	assert.Equal(t, 13, generated.IPv4("123.21.123.21", 0))
	assert.Equal(t, 13, generated.IPv4("255.255.255.1", 0))
	assert.Equal(t, -1, generated.IPv4("möp", 0))
}
