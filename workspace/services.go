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
	DocumentManager DocumentManager
	DocumentUpdater DocumentUpdater
	Builder         Builder
	DocumentParser  DocumentParser
	Initializer     Initializer
}

// CreateDefaultServices creates the default services for the workspace package.
// If the services are already set, they are not overwritten.
// Package dependencies: textdoc, linking, generated
func CreateDefaultServices(c WorkspaceSrvCont) {
	s := c.Workspace()
	if s.DocumentManager == nil {
		s.DocumentManager = NewDefaultDocumentManager()
	}
	if s.DocumentUpdater == nil {
		s.DocumentUpdater = NewDefaultDocumentUpdater(c)
	}
	if s.Builder == nil {
		s.Builder = NewDefaultBuilder(c)
	}
	if s.DocumentParser == nil {
		s.DocumentParser = NewDefaultDocumentParser(c)
	}
	if s.Initializer == nil {
		s.Initializer = NewDefaultInitializer(c)
	}
}

// GeneratedSrvCont is an interface for service containers which include the generated services.
// TODO Move this stuff to the core package?
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
