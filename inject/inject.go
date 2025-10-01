// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package inject

import (
	"fmt"
)

// ServiceContainer is a container for services used for dependency injection.
type ServiceContainer struct {
	services map[string]any
}

// NewServiceContainer creates a new DI service container.
func NewServiceContainer() *ServiceContainer {
	return &ServiceContainer{
		services: make(map[string]any),
	}
}

// Injectable is for services that retrieve their dependencies from a DI container.
type Injectable interface {
	// Inject is called by `InjectAll` to inject the service's dependencies.
	// In service implementations, you can use `Get` to retrieve dependencies from the container.
	// However, you cannot assume that a retrieved dependency is already injected while injecting the service.
	Inject(container *ServiceContainer) error
}

// ServiceKey holds the map string key and the service type.
type ServiceKey[T any] struct {
	key string
}

// NewServiceKey creates a new ServiceKey with the given key string.
func NewServiceKey[T any](key string) ServiceKey[T] {
	return ServiceKey[T]{key: key}
}

// Register puts the service into the container map, checking if it's already registered.
func Register[T any](key ServiceKey[T], service T, container *ServiceContainer) error {
	if container.services[key.key] != nil {
		return fmt.Errorf("service already registered: %s", key.key)
	}

	container.services[key.key] = service
	return nil
}

// Override replaces an existing service in the container map, expecting the key to be already set.
func Override[T any](key ServiceKey[T], service T, container *ServiceContainer) error {
	if container.services[key.key] == nil {
		return fmt.Errorf("service not registered, cannot override: %s", key.key)
	}

	container.services[key.key] = service
	return nil
}

// InjectAll runs over the contents of the service container and calls the Inject method where applicable.
func InjectAll(container *ServiceContainer) error {
	for key, service := range container.services {
		if i, ok := service.(Injectable); ok {
			if err := i.Inject(container); err != nil {
				return fmt.Errorf("failed to inject service %s: %w", key, err)
			}
		}
	}
	return nil
}

// Get retrieves a service from the container and casts it to the type T.
func Get[T any](key ServiceKey[T], container *ServiceContainer) (T, error) {
	var s T
	mapValue := container.services[key.key]
	if mapValue == nil {
		return s, fmt.Errorf("service not injected: %s", key.key)
	}

	if s, ok := mapValue.(T); ok {
		return s, nil
	} else {
		return s, fmt.Errorf("service %s is not of type %T", key.key, *new(T))
	}
}
