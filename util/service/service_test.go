// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package service

import (
	"testing"
)

type Greeter interface {
	Greet() string
}

type EnglishGreeter struct {
	value string
}

func (g *EnglishGreeter) Greet() string {
	return g.value
}

type Counter interface {
	Count() int
}

type StaticCounter struct {
	value int
}

func (c *StaticCounter) Count() int {
	return c.value
}

func TestNewContainer_IsInitiallyEmptyAndNotSealed(t *testing.T) {
	container := NewContainer()

	if Has[Greeter](container) {
		t.Fatal("expected no services in a new container")
	}

	_, err := Get[Greeter](container)
	if err == nil {
		t.Fatal("expected error when getting from an unsealed container")
	}
	if err.Error() != "container is not sealed" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestSeal_PreventsFurtherServicesFromBeingAdded(t *testing.T) {
	container := NewContainer()
	container.Seal()

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("expected MustPut to panic for a sealed container")
		}
		if recovered != "container is sealed" {
			t.Fatalf("unexpected panic: %v", recovered)
		}
	}()

	MustPut[Greeter](container, &EnglishGreeter{value: "hello"})
}

func TestHas_ReportsPresenceOfService(t *testing.T) {
	container := NewContainer()
	service := &EnglishGreeter{value: "hello"}
	MustPut[Greeter](container, service)

	if !Has[Greeter](container) {
		t.Fatal("expected Has to report true for inserted service")
	}
	if Has[Counter](container) {
		t.Fatal("expected Has to report false for missing service type")
	}
}

func TestGet_ReturnsServiceWhenSealedAndPresent(t *testing.T) {
	container := NewContainer()
	service := &EnglishGreeter{value: "hello"}
	MustPut[Greeter](container, service)
	container.Seal()

	got, err := Get[Greeter](container)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != service {
		t.Fatal("expected retrieved service to match inserted service")
	}
}

func TestGet_ReturnsErrorWhenServiceNotFound(t *testing.T) {
	container := NewContainer()
	container.Seal()

	_, err := Get[Greeter](container)
	if err == nil {
		t.Fatal("expected error for missing service")
	}
	if err.Error() != "service Greeter not found" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMustGet_PanicsWhenContainerNotSealed(t *testing.T) {
	container := NewContainer()

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("expected MustGet to panic for unsealed container")
		}
	}()

	_ = MustGet[Greeter](container)
}

func TestMustGet_PanicsWhenServiceNotFound(t *testing.T) {
	container := NewContainer()
	container.Seal()

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("expected MustGet to panic for missing service")
		}
	}()

	_ = MustGet[Greeter](container)
}

func TestMustGet_ReturnsServiceWhenPresent(t *testing.T) {
	container := NewContainer()
	service := &StaticCounter{value: 42}
	MustPut[Counter](container, service)
	container.Seal()

	got := MustGet[Counter](container)
	if got != service {
		t.Fatal("expected MustGet to return inserted service")
	}
}
