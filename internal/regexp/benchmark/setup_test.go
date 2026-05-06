package benchmark

import (
	"fmt"
	"os"
	"strings"
	"testing"

	format "typefox.dev/fastbelt/internal/generator"
	custom "typefox.dev/fastbelt/internal/regexp"
	"typefox.dev/fastbelt/util/codegen"
)

func writeRegexpFile(name string, pattern string) {
	regexp := custom.MustCompile(pattern)
	root := codegen.NewNode()
	result := regexp.(*custom.RegexpImpl).GenerateRegExp(name, name)
	root.AppendLine("package benchmarkGenerated")
	root.AppendLine()
	root.AppendLine("import (")
	for imp := range result.Imports {
		root.AppendLine(fmt.Sprintf(`	"%s"`, imp))
	}
	root.AppendLine(")")
	root.AppendLine()
	root.AppendNode(result.Vars)
	root.AppendNode(result.Code)
	err := os.WriteFile("generated/"+strings.ToLower(name)+".go", []byte(format.FormatIfPossible(root.String())), 0644)
	if err != nil {
		panic(err)
	}
}

func TestGenerateCustomRegExpFiles(t *testing.T) {
	writeRegexpFile("URL", URLPattern())
	writeRegexpFile("Email", EmailPattern())
	writeRegexpFile("IPv4", IPv4Pattern())
}
