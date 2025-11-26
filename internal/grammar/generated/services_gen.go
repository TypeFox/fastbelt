package generated

import (
	"typefox.dev/fastbelt/workspace"
)

func CreateDefaultServices(c workspace.GeneratedSrvCont) {
	s := c.Generated()
	if s.Lexer == nil {
		s.Lexer = NewLexer()
	}
	if s.Parser == nil {
		s.Parser = NewParser()
	}
}
