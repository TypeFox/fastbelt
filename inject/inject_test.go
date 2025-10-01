// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package inject

import (
	"testing"
)

// TestService is a simple service that doesn't implement Injectable
type TestService struct {
	Name string
}

// TestInjectableService is a service that implements Injectable and depends on TestService
type TestInjectableService struct {
	Dependency *TestService
	Value      string
}

// Inject implements the Injectable interface
func (s *TestInjectableService) Inject(container *ServiceContainer) error {
	key := NewServiceKey[*TestService]("test-service")
	dependency, err := Get(key, container)
	if err != nil {
		return err
	}
	s.Dependency = dependency
	return nil
}

func TestDependencyInjection(t *testing.T) {
	// Create service keys
	testServiceKey := NewServiceKey[*TestService]("test-service")
	injectableServiceKey := NewServiceKey[*TestInjectableService]("injectable-service")

	// Create service container
	container := NewServiceContainer()

	// Create services
	dependency := &TestService{Name: "dependency"}
	injectable := &TestInjectableService{Value: "injectable"}

	// Register services
	if err := Register(testServiceKey, dependency, container); err != nil {
		t.Fatalf("Failed to register dependency: %v", err)
	}

	if err := Register(injectableServiceKey, injectable, container); err != nil {
		t.Fatalf("Failed to register injectable service: %v", err)
	}

	// Verify dependency is not set before injection
	if injectable.Dependency != nil {
		t.Error("Dependency should be nil before injection")
	}

	// Inject dependencies
	if err := InjectAll(container); err != nil {
		t.Fatalf("Failed to inject dependencies: %v", err)
	}

	// Verify dependency injection worked
	if injectable.Dependency == nil {
		t.Error("Dependency should be set after injection")
	}

	if injectable.Dependency.Name != "dependency" {
		t.Errorf("Expected dependency name 'dependency', got '%s'", injectable.Dependency.Name)
	}

	if injectable.Value != "injectable" {
		t.Errorf("Expected injectable value 'injectable', got '%s'", injectable.Value)
	}

	// Verify we can retrieve services from container
	retrievedDependency, err := Get(testServiceKey, container)
	if err != nil {
		t.Fatalf("Failed to retrieve dependency: %v", err)
	}

	if retrievedDependency.Name != "dependency" {
		t.Errorf("Expected retrieved dependency name 'dependency', got '%s'", retrievedDependency.Name)
	}

	retrievedInjectable, err := Get(injectableServiceKey, container)
	if err != nil {
		t.Fatalf("Failed to retrieve injectable service: %v", err)
	}

	if retrievedInjectable.Value != "injectable" {
		t.Errorf("Expected retrieved injectable value 'injectable', got '%s'", retrievedInjectable.Value)
	}

	// Verify the dependency is the same instance
	if retrievedInjectable.Dependency != retrievedDependency {
		t.Error("Dependency should be the same instance")
	}
}

func TestOverride(t *testing.T) {
	// Create service key
	serviceKey := NewServiceKey[*TestService]("test-service")

	// Create service container
	container := NewServiceContainer()

	// Create initial service
	originalService := &TestService{Name: "original"}

	// Register the service
	if err := Register(serviceKey, originalService, container); err != nil {
		t.Fatalf("Failed to register original service: %v", err)
	}

	// Verify original service is registered
	retrieved, err := Get(serviceKey, container)
	if err != nil {
		t.Fatalf("Failed to retrieve original service: %v", err)
	}
	if retrieved.Name != "original" {
		t.Errorf("Expected original name 'original', got '%s'", retrieved.Name)
	}

	// Create replacement service
	replacementService := &TestService{Name: "replacement"}

	// Override the service
	if err := Override(serviceKey, replacementService, container); err != nil {
		t.Fatalf("Failed to override service: %v", err)
	}

	// Verify service was overridden
	retrieved, err = Get(serviceKey, container)
	if err != nil {
		t.Fatalf("Failed to retrieve overridden service: %v", err)
	}
	if retrieved.Name != "replacement" {
		t.Errorf("Expected overridden name 'replacement', got '%s'", retrieved.Name)
	}

	// Verify it's the same instance
	if retrieved != replacementService {
		t.Error("Retrieved service should be the same instance as the replacement")
	}
}

func TestOverrideNonExistentService(t *testing.T) {
	// Create service key
	serviceKey := NewServiceKey[*TestService]("non-existent-service")

	// Create service container
	container := NewServiceContainer()

	// Create service to override with
	service := &TestService{Name: "test"}

	// Try to override non-existent service - should fail
	err := Override(serviceKey, service, container)
	if err == nil {
		t.Error("Expected error when overriding non-existent service")
	}

	expectedError := "service not registered, cannot override: non-existent-service"
	if err.Error() != expectedError {
		t.Errorf("Expected error '%s', got '%s'", expectedError, err.Error())
	}
}
