package services

import (
	"typefox.dev/fastbelt/internal/grammar/generated"
	"typefox.dev/fastbelt/linking"
	"typefox.dev/fastbelt/server"
	"typefox.dev/fastbelt/textdoc"
	"typefox.dev/fastbelt/workspace"
)

type GrammarSrv struct {
	textdoc.TextdocSrvContBlock
	workspace.GeneratedSrvContBlock
	workspace.WorkspaceSrvContBlock
	linking.LinkingSrvContBlock
	server.ServerSrvContBlock
	generated.FastbeltLinkingSrvContBlock
}

func CreateServices() *GrammarSrv {
	srv := &GrammarSrv{}
	textdoc.CreateDefaultServices(srv)
	workspace.CreateDefaultServices(srv)
	server.CreateDefaultServices(srv)
	linking.CreateDefaultServices(srv)
	generated.CreateDefaultServices(srv)
	return srv
}
