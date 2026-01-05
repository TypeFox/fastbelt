package benchmark

import (
	"os"
	"strings"
	"testing"

	"typefox.dev/fastbelt/generator"
	custom "typefox.dev/fastbelt/internal/regexp"
)

func writeRegexpFile(name string, pattern string) {
	regexp := custom.MustCompile(pattern)
	root := generator.NewNode()
	root.AppendLine("package benchmarkGenerated")
	root.AppendLine()
	root.AppendLine("import (")
	root.AppendLine(`	"sort"`)
	root.AppendLine(`	"unicode/utf8"`)
	root.AppendLine(")")
	root.AppendLine()
	root.AppendNode(regexp.(*custom.RegexpImpl).GenerateRegExp(name))
	os.WriteFile("generated/"+strings.ToLower(name)+".go", []byte(root.String()), 0644)
}

func TestGenerateCustomRegExpFiles(t *testing.T) {
	t.Skip()
	writeRegexpFile("URL", URLPattern())
	writeRegexpFile("Email", EmailPattern())
	writeRegexpFile("IPv4", IPv4Pattern())
}
