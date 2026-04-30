// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package service

import (
	"fmt"
	"reflect"
)

// Container holds a collection of services to be used in the framework.
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
func Has[T any](container *Container) bool {
	t := reflect.TypeFor[T]()
	_, ok := container.services[t]
	return ok
}

// Get retrieves a service from the container.
// It returns an error if the container is not sealed or the service is not found.
//
// TODO make this a method once generic methods are supported in Go.
func Get[T any](container *Container) (T, error) {
	if !container.sealed {
		var zero T
		return zero, fmt.Errorf("container is not sealed")
	}
	t := reflect.TypeFor[T]()
	service, ok := container.services[t]
	if !ok {
		var zero T
		return zero, fmt.Errorf("service %s not found", t.Name())
	}
	return service.(T), nil
}

// MustGet retrieves a service from the container.
// It panics if the container is not sealed or the service is not found.
//
// TODO make this a method once generic methods are supported in Go.
func MustGet[T any](container *Container) T {
	service, err := Get[T](container)
	if err != nil {
		panic(err)
	}
	return service
}

// MustPut puts a service into the container.
// It panics if the container is sealed or the service already exists.
//
// TODO make this a method once generic methods are supported in Go.
func MustPut[T any](container *Container, service T) {
	if container.sealed {
		panic("container is sealed")
	}
	t := reflect.TypeFor[T]()
	if _, ok := container.services[t]; ok {
		panic(fmt.Sprintf("service %s already exists", t.Name()))
	}
	container.services[t] = service
}

// MustOverride overrides a service in the container.
// It panics if the container is sealed or the service does not exist.
//
// TODO make this a method once generic methods are supported in Go.
func MustOverride[T any](container *Container, service T) {
	if container.sealed {
		panic("container is sealed")
	}
	t := reflect.TypeFor[T]()
	if _, ok := container.services[t]; !ok {
		panic(fmt.Sprintf("service %s does not exist", t.Name()))
	}
	container.services[t] = service
}
