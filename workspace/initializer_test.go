// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package workspace

import (
	"context"
	"testing"

	"typefox.dev/fastbelt/util/service"
)

func TestDefaultInitializerReadyClosesAfterInitialize(t *testing.T) {
	init := NewDefaultInitializer(service.NewContainer())
	select {
	case <-init.Ready():
		t.Fatal("Ready closed before Initialize")
	default:
	}
	if err := init.Initialize(context.Background(), nil); err != nil {
		t.Fatal(err)
	}
	select {
	case <-init.Ready():
	default:
		t.Fatal("Ready not closed after Initialize")
	}
}

func TestNoopInitializerReadyIsClosed(t *testing.T) {
	select {
	case <-NewNoopInitializer().Ready():
	default:
		t.Fatal("Ready not closed")
	}
}
