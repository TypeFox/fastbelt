package benchmark

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"typefox.dev/fastbelt/generator"
	custom "typefox.dev/fastbelt/internal/regexp"
)

func writeRegexpFile(name string, pattern string) {
	regexp := custom.MustCompile(pattern)
	root := generator.NewNode()
	result := regexp.(*custom.RegexpImpl).GenerateRegExp(name, name)
	root.AppendLine("package benchmarkGenerated")
	root.AppendLine()
	root.AppendLine("import (")
	for imp := range result.Imports {
		root.AppendLine(fmt.Sprintf(`	"%s"`, imp))
	}
	root.AppendLine(")")
	root.AppendLine()
	root.AppendNode(result.Lookup)
	root.AppendNode(result.Next)
	root.AppendNode(result.Code)
	err := os.WriteFile("generated/"+strings.ToLower(name)+".go", []byte(root.String()), 0644)
	if err != nil {
		panic(err)
	}
}

func TestGenerateCustomRegExpFiles(t *testing.T) {
	writeRegexpFile("URL", URLPattern())
	writeRegexpFile("Email", EmailPattern())
	writeRegexpFile("IPv4", IPv4Pattern())
}
