// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package generated

import (
	core "typefox.dev/fastbelt"
	"typefox.dev/fastbelt/lexer"
	"typefox.dev/fastbelt/parser"
)

// GeneratedSrvCont is an interface for service containers which include the generated services.
type GeneratedSrvCont interface {
	Generated() *GeneratedSrv
}

// GeneratedSrvContBlock is used to define a service container satisfying GeneratedSrvCont.
type GeneratedSrvContBlock struct {
	generated GeneratedSrv
}

func (b *GeneratedSrvContBlock) Generated() *GeneratedSrv {
	return &b.generated
}

// GeneratedSrv contains the generated services for a specific language.
type GeneratedSrv struct {
	Lexer            lexer.Lexer
	Parser           parser.Parser
	SymbolContainers core.SymbolContainers
}
