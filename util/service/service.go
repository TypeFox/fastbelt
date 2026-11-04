// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package service

import (
	"fmt"
	"iter"
	"maps"
	"reflect"

	"typefox.dev/fastbelt/util/extiter"
)

// A Container holds a collection of services keyed by Go type.
// Register services with [Container.Put] and [Container.Override] during setup, call [Container.Seal]
// when wiring is complete, then retrieve services with [Container.Get] or [Container.MustGet].
// The zero value is not usable; create containers with [NewContainer].
type Container struct {
	services map[reflect.Type]any
	sealed   bool
}

// NewContainer creates a new container.
// The container is initially empty and not sealed.
func NewContainer() *Container {
	return &Container{
		services: make(map[reflect.Type]any),
		sealed:   false,
	}
}

// Seal prevents further services from being added to the container.
// This must be called before the container is used.
func (c *Container) Seal() {
	c.sealed = true
}

// Has checks if a service is in the container.
func (c *Container) Has[T any]() bool {
	t := reflect.TypeFor[T]()
	_, ok := c.services[t]
	return ok
}

// Get retrieves a service from the container.
// It returns an error if the container is not sealed or the service is not found.
func (c *Container) Get[T any]() (T, error) {
	if !c.sealed {
		var zero T
		return zero, fmt.Errorf("container is not sealed")
	}
	t := reflect.TypeFor[T]()
	service, ok := c.services[t]
	if !ok {
		var zero T
		return zero, fmt.Errorf("service %s not found", t.Name())
	}
	return service.(T), nil
}

// GetAll retrieves all services of type T from the container.
// Will return an empty sequence if the container is not sealed.
func (c *Container) GetAll[T any]() iter.Seq[T] {
	if !c.sealed {
		return extiter.Empty[T]()
	}
	return func(yield func(T) bool) {
		for service := range c.All() {
			if t, ok := service.(T); ok && !yield(t) {
				return
			}
		}
	}
}

// MustGet retrieves a service from the container.
// It panics if the container is not sealed or the service is not found.
func (c *Container) MustGet[T any]() T {
	service, err := c.Get[T]()
	if err != nil {
		panic(err)
	}
	return service
}

// Put puts a service into the container.
// It panics if the container is sealed or the service already exists.
func (c *Container) Put[T any](service T) {
	if c.sealed {
		panic("container is sealed")
	}
	t := reflect.TypeFor[T]()
	if _, ok := c.services[t]; ok {
		panic(fmt.Sprintf("service %s already exists", t.Name()))
	}
	c.services[t] = service
}

// Override overrides a service in the container.
// It panics if the container is sealed or the service does not exist.
func (c *Container) Override[T any](service T) {
	if c.sealed {
		panic("container is sealed")
	}
	t := reflect.TypeFor[T]()
	if _, ok := c.services[t]; !ok {
		panic(fmt.Sprintf("service %s does not exist", t.Name()))
	}
	c.services[t] = service
}

// All returns all services in the container in arbitrary order.
func (c *Container) All() iter.Seq[any] {
	return maps.Values(c.services)
}
