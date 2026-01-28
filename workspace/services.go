// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"typefox.dev/fastbelt/lexer"
	"typefox.dev/fastbelt/linking"
	"typefox.dev/fastbelt/parser"
	"typefox.dev/fastbelt/textdoc"
)

// WorkspaceSrvCont is an interface for service containers which include the workspace services.
type WorkspaceSrvCont interface {
	textdoc.TextdocSrvCont
	linking.LinkingSrvCont
	GeneratedSrvCont
	Workspace() *WorkspaceSrv
}

// WorkspaceSrvContBlock is used to define a service container satisfying WorkspaceSrvCont.
type WorkspaceSrvContBlock struct {
	workspace WorkspaceSrv
}

func (b *WorkspaceSrvContBlock) Workspace() *WorkspaceSrv {
	return &b.workspace
}

// WorkspaceSrv contains the services for the workspace package.
type WorkspaceSrv struct {
	Builder        Builder
	DocumentParser DocumentParser
}

// CreateDefaultServices creates the default services for the workspace package.
// If the services are already set, they are not overwritten.
// Package dependencies: textdoc, generated
func CreateDefaultServices(c WorkspaceSrvCont) {
	s := c.Workspace()
	if s.Builder == nil {
		s.Builder = NewDefaultBuilder(c)
	}
	if s.DocumentParser == nil {
		s.DocumentParser = NewDefaultDocumentParser(c)
	}
}

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
	Lexer  lexer.Lexer
	Parser parser.Parser
}
