// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package service

import (
	"fmt"
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

func panicMessage(recovered any) string {
	switch v := recovered.(type) {
	case error:
		return v.Error()
	case string:
		return v
	default:
		return fmt.Sprint(v)
	}
}

func TestNewContainer_IsInitiallyEmptyAndNotSealed(t *testing.T) {
	container := NewContainer()

	if container.Has[Greeter]() {
		t.Fatal("expected no services in a new container")
	}

	_, err := container.Get[Greeter]()
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
			t.Fatal("expected Put to panic for a sealed container")
		}
		if panicMessage(recovered) != "container is sealed" {
			t.Fatalf("unexpected panic: %v", recovered)
		}
	}()

	container.Put[Greeter](&EnglishGreeter{value: "hello"})
}

func TestPut_PanicsWhenContainerSealed(t *testing.T) {
	container := NewContainer()
	container.Seal()

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("expected Put to panic for a sealed container")
		}
		if panicMessage(recovered) != "container is sealed" {
			t.Fatalf("unexpected panic: %v", recovered)
		}
	}()

	container.Put[Greeter](&EnglishGreeter{value: "hello"})
}

func TestPut_PanicsWhenServiceAlreadyExists(t *testing.T) {
	container := NewContainer()
	container.Put[Greeter](&EnglishGreeter{value: "hello"})

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("expected Put to panic for duplicate service")
		}
		if panicMessage(recovered) != "service Greeter already exists" {
			t.Fatalf("unexpected panic: %v", recovered)
		}
	}()

	container.Put[Greeter](&EnglishGreeter{value: "hi"})
}

func TestHas_ReportsPresenceOfService(t *testing.T) {
	container := NewContainer()
	service := &EnglishGreeter{value: "hello"}
	container.Put[Greeter](service)

	if !container.Has[Greeter]() {
		t.Fatal("expected Has to report true for inserted service")
	}
	if container.Has[Counter]() {
		t.Fatal("expected Has to report false for missing service type")
	}
}

func TestGet_ReturnsServiceWhenSealedAndPresent(t *testing.T) {
	container := NewContainer()
	service := &EnglishGreeter{value: "hello"}
	container.Put[Greeter](service)
	container.Seal()

	got, err := container.Get[Greeter]()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != service {
		t.Fatal("expected retrieved service to match inserted service")
	}
}

func TestOverride_PanicsWhenContainerSealed(t *testing.T) {
	container := NewContainer()
	container.Put[Greeter](&EnglishGreeter{value: "hello"})
	container.Seal()

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("expected Override to panic for a sealed container")
		}
		if panicMessage(recovered) != "container is sealed" {
			t.Fatalf("unexpected panic: %v", recovered)
		}
	}()

	container.Override[Greeter](&EnglishGreeter{value: "hi"})
}

func TestOverride_PanicsWhenServiceDoesNotExist(t *testing.T) {
	container := NewContainer()

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("expected Override to panic for missing service")
		}
		if panicMessage(recovered) != "service Greeter does not exist" {
			t.Fatalf("unexpected panic: %v", recovered)
		}
	}()

	container.Override[Greeter](&EnglishGreeter{value: "hello"})
}

func TestOverride_ReplacesExistingService(t *testing.T) {
	container := NewContainer()
	first := &EnglishGreeter{value: "hello"}
	second := &EnglishGreeter{value: "hi"}
	container.Put[Greeter](first)

	container.Override[Greeter](second)
	container.Seal()
	got := container.MustGet[Greeter]()
	if got != second {
		t.Fatal("expected Override to replace the existing service")
	}
}

func TestGet_ReturnsErrorWhenServiceNotFound(t *testing.T) {
	container := NewContainer()
	container.Seal()

	_, err := container.Get[Greeter]()
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

	_ = container.MustGet[Greeter]()
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

	_ = container.MustGet[Greeter]()
}

func TestMustGet_ReturnsServiceWhenPresent(t *testing.T) {
	container := NewContainer()
	service := &StaticCounter{value: 42}
	container.Put[Counter](service)
	container.Seal()

	got := container.MustGet[Counter]()
	if got != service {
		t.Fatal("expected MustGet to return inserted service")
	}
}

func TestOverride_PanicsWhenContainerSealed_MustBehavior(t *testing.T) {
	container := NewContainer()
	container.Put[Greeter](&EnglishGreeter{value: "hello"})
	container.Seal()

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("expected Override to panic for sealed container")
		}
		if panicMessage(recovered) != "container is sealed" {
			t.Fatalf("unexpected panic: %v", recovered)
		}
	}()

	container.Override[Greeter](&EnglishGreeter{value: "hi"})
}

func TestOverride_PanicsWhenServiceDoesNotExist_MustBehavior(t *testing.T) {
	container := NewContainer()

	defer func() {
		recovered := recover()
		if recovered == nil {
			t.Fatal("expected Override to panic for missing service")
		}
		if panicMessage(recovered) != "service Greeter does not exist" {
			t.Fatalf("unexpected panic: %v", recovered)
		}
	}()

	container.Override[Greeter](&EnglishGreeter{value: "hello"})
}
