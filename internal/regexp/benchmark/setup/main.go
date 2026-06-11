package main

//go:generate go run main.go

import (
	"fmt"
	"os"
	"strings"

	format "typefox.dev/fastbelt/internal/generator"
	custom "typefox.dev/fastbelt/internal/regexp"
	"typefox.dev/fastbelt/internal/regexp/benchmark"
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
	err := os.Mkdir("../generated", 0755)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
	err = os.WriteFile("../generated/"+strings.ToLower(name)+".go", []byte(format.FormatIfPossible(root.String())), 0644)
	if err != nil {
		panic(err)
	}
}

func main() {
	writeRegexpFile("URL", benchmark.URLPattern())
	writeRegexpFile("Email", benchmark.EmailPattern())
	writeRegexpFile("IPv4", benchmark.IPv4Pattern())
}
