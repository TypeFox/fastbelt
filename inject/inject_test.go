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

// CircularServiceA depends on CircularServiceB
type CircularServiceA struct {
	ServiceB *CircularServiceB
	Value    string
}

// Inject implements the Injectable interface for CircularServiceA
func (s *CircularServiceA) Inject(container *ServiceContainer) error {
	key := NewServiceKey[*CircularServiceB]("circular-service-b")
	serviceB, err := Get(key, container)
	if err != nil {
		return err
	}
	s.ServiceB = serviceB
	return nil
}

// CircularServiceB depends on CircularServiceA
type CircularServiceB struct {
	ServiceA *CircularServiceA
	Value    string
}

// Inject implements the Injectable interface for CircularServiceB
func (s *CircularServiceB) Inject(container *ServiceContainer) error {
	key := NewServiceKey[*CircularServiceA]("circular-service-a")
	serviceA, err := Get(key, container)
	if err != nil {
		return err
	}
	s.ServiceA = serviceA
	return nil
}

func TestCircularDependency(t *testing.T) {
	// Create service keys
	serviceAKey := NewServiceKey[*CircularServiceA]("circular-service-a")
	serviceBKey := NewServiceKey[*CircularServiceB]("circular-service-b")

	// Create service container
	container := NewServiceContainer()

	// Create services
	serviceA := &CircularServiceA{Value: "service-a"}
	serviceB := &CircularServiceB{Value: "service-b"}

	// Register services
	if err := Register(serviceAKey, serviceA, container); err != nil {
		t.Fatalf("Failed to register service A: %v", err)
	}

	if err := Register(serviceBKey, serviceB, container); err != nil {
		t.Fatalf("Failed to register service B: %v", err)
	}

	// Verify dependencies are not set before injection
	if serviceA.ServiceB != nil {
		t.Error("ServiceA.ServiceB should be nil before injection")
	}
	if serviceB.ServiceA != nil {
		t.Error("ServiceB.ServiceA should be nil before injection")
	}

	// Inject dependencies - this should work despite circular dependency
	if err := InjectAll(container); err != nil {
		t.Fatalf("Failed to inject circular dependencies: %v", err)
	}

	// Verify circular dependency injection worked
	if serviceA.ServiceB == nil {
		t.Error("ServiceA.ServiceB should be set after injection")
	}
	if serviceB.ServiceA == nil {
		t.Error("ServiceB.ServiceA should be set after injection")
	}

	// Verify the services reference each other correctly
	if serviceA.ServiceB != serviceB {
		t.Error("ServiceA.ServiceB should reference the same instance as serviceB")
	}
	if serviceB.ServiceA != serviceA {
		t.Error("ServiceB.ServiceA should reference the same instance as serviceA")
	}

	// Verify original values are preserved
	if serviceA.Value != "service-a" {
		t.Errorf("Expected serviceA value 'service-a', got '%s'", serviceA.Value)
	}
	if serviceB.Value != "service-b" {
		t.Errorf("Expected serviceB value 'service-b', got '%s'", serviceB.Value)
	}

	// Verify we can retrieve services from container
	retrievedA, err := Get(serviceAKey, container)
	if err != nil {
		t.Fatalf("Failed to retrieve service A: %v", err)
	}
	if retrievedA.Value != "service-a" {
		t.Errorf("Expected retrieved serviceA value 'service-a', got '%s'", retrievedA.Value)
	}

	retrievedB, err := Get(serviceBKey, container)
	if err != nil {
		t.Fatalf("Failed to retrieve service B: %v", err)
	}
	if retrievedB.Value != "service-b" {
		t.Errorf("Expected retrieved serviceB value 'service-b', got '%s'", retrievedB.Value)
	}

	// Verify the circular references are maintained in retrieved services
	if retrievedA.ServiceB != retrievedB {
		t.Error("Retrieved serviceA.ServiceB should reference the same instance as retrieved serviceB")
	}
	if retrievedB.ServiceA != retrievedA {
		t.Error("Retrieved serviceB.ServiceA should reference the same instance as retrieved serviceA")
	}
}
