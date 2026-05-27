// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"typefox.dev/fastbelt/internal/scaffold"
)

type scaffoldOptions struct {
	modulePath        string
	packagePath       string
	language          string
	noVSCodeExtension bool
}

func runScaffoldCLI(opts scaffoldOptions) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	packageRel := strings.TrimSpace(opts.packagePath)
	if packageRel == "" {
		packageRel = "."
	}

	scaffolder := &scaffold.Scaffolder{
		CreateVSCodeExtension: !opts.noVSCodeExtension,
		Language:              opts.language,
		CreateModule:          opts.modulePath != "",
	}

	moduleArg := strings.TrimSpace(opts.modulePath)
	if moduleArg != "" {
		importPath := path.Clean(moduleArg)
		dirBase := path.Base(importPath)
		if dirBase == "." || dirBase == "/" {
			return fmt.Errorf("invalid -module path %q (cannot determine output directory)", opts.modulePath)
		}
		moduleDir := filepath.Join(wd, dirBase)
		err = scaffolder.PopulateDirectoriesFromModuleDir(moduleDir, importPath, packageRel)
		if err != nil {
			return err
		}
	} else {
		err = scaffolder.PopulateDirectorysFromWorkDir(wd, packageRel)
		if err != nil {
			return err
		}
	}

	err = scaffolder.Run()
	if err != nil {
		debug, _ := json.MarshalIndent(scaffolder, "", "  ")
		return fmt.Errorf("scaffold: %w; scaffolder: %s", err, string(debug))
	}
	return nil
}
