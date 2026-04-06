// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package scaffold

import (
	"runtime/debug"
	"strings"
)

const fastbeltModulePath = "typefox.dev/fastbelt"

// fastbeltModuleVersion returns a version string suitable for "go get M@V"
// (for example v1.2.3 or pseudo-versions), or "latest" when unknown.
func fastbeltModuleVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "latest"
	}
	if info.Main.Path == fastbeltModulePath+"/cmd/fastbelt" || info.Main.Path == fastbeltModulePath+"/cmd" || info.Main.Path == fastbeltModulePath {
		if info.Main.Version != "" && info.Main.Version != "(devel)" && !isUnsetModulePseudoVersion(info.Main.Version) {
			return info.Main.Version
		}
	}
	for _, m := range info.Deps {
		if m.Path == fastbeltModulePath && m.Version != "" {
			if isUnsetModulePseudoVersion(m.Version) {
				return "latest"
			}
			return m.Version
		}
	}
	return "latest"
}

func isUnsetModulePseudoVersion(v string) bool {
	return strings.Contains(v, "000000000000") || strings.HasPrefix(v, "v0.0.0-00010101000000")
}
