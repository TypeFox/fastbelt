// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package server

import (
	"context"
	"fmt"

	"typefox.dev/fastbelt/util/service"
	"typefox.dev/lsp"
)

// CommandProvider is a service for handling LSP execute command requests.
//
// Usage:
//
//	type MyCommandProvider struct{ sc *service.Container }
//
//	func (p *MyCommandProvider) HandleExecuteCommandRequest(ctx context.Context, params *lsp.ExecuteCommandParams) (any, error) {
//	    switch params.Command {
//	    case "myLanguage.refactor":
//	        // Execute refactoring with params.Arguments
//	        return "Refactoring completed", nil
//	    default:
//	        return nil, fmt.Errorf("unknown command: %s", params.Command)
//	    }
//	}
type CommandProvider interface {
	HandleExecuteCommandRequest(ctx context.Context, params *lsp.ExecuteCommandParams) (any, error)
}

// DefaultCommandProvider returns an error for all commands.
type DefaultCommandProvider struct {
	sc *service.Container
}

func NewDefaultCommandProvider(sc *service.Container) CommandProvider {
	return &DefaultCommandProvider{sc: sc}
}

func (p *DefaultCommandProvider) HandleExecuteCommandRequest(ctx context.Context, params *lsp.ExecuteCommandParams) (any, error) {
	return nil, fmt.Errorf("command not found: %s", params.Command)
}
