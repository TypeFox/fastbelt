// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package textdoc

import "typefox.dev/fastbelt/util/service"

// SetupDefaultServices sets up the default services for the textdoc package.
// If any service is already set, it's not overwritten.
func SetupDefaultServices(sc *service.Container) {
	if !service.Has[Store](sc) {
		service.MustPut(sc, NewDefaultStore())
	}
}
