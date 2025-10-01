// Copyright 2025 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

package inject

import (
	"testing"
)

// TestServices extends Services with test-specific services using embedding
type TestServices struct {
	*Services
	TestService          *TestService
	TestDependentService *TestDependentService
	CircularServiceA     *CircularServiceA
	CircularServiceB     *CircularServiceB
}

// NewTestServices creates a new test services container
func NewTestServices() *TestServices {
	return &TestServices{
		Services: NewServices(),
	}
}

// TestService is a simple service that doesn't need other dependencies
type TestService struct {
	Name string
}

// TestDependentService is a service that depends on TestService
type TestDependentService struct {
	services *TestServices
	Value    string
}

// NewTestDependentService creates a new TestDependentService with access to services
func NewTestDependentService(services *TestServices, value string) *TestDependentService {
	return &TestDependentService{
		services: services,
		Value:    value,
	}
}

// GetDependency accesses the dependency directly from services
func (s *TestDependentService) GetDependency() *TestService {
	return s.services.TestService
}

func TestDependencyInjection(t *testing.T) {
	// Create service container with embedded base services
	services := NewTestServices()

	// Create and set services
	dependency := &TestService{Name: "dependency"}
	services.TestService = dependency

	dependent := NewTestDependentService(services, "dependent-value")
	services.TestDependentService = dependent

	// Verify dependency can be accessed
	retrievedDependency := dependent.GetDependency()
	if retrievedDependency == nil {
		t.Error("Dependency should be accessible")
	}

	if retrievedDependency.Name != "dependency" {
		t.Errorf("Expected dependency name 'dependency', got '%s'", retrievedDependency.Name)
	}

	if dependent.Value != "dependent-value" {
		t.Errorf("Expected dependent value 'dependent-value', got '%s'", dependent.Value)
	}

	// Verify we can access services from container
	retrievedFromContainer := services.TestService
	if retrievedFromContainer.Name != "dependency" {
		t.Errorf("Expected retrieved dependency name 'dependency', got '%s'", retrievedFromContainer.Name)
	}

	retrievedDependentFromContainer := services.TestDependentService
	if retrievedDependentFromContainer.Value != "dependent-value" {
		t.Errorf("Expected retrieved dependent value 'dependent-value', got '%s'", retrievedDependentFromContainer.Value)
	}

	// Verify the dependency is the same instance
	if retrievedDependency != retrievedFromContainer {
		t.Error("Dependency should be the same instance")
	}
}

func TestOverride(t *testing.T) {
	// Create service container
	services := NewTestServices()

	// Create initial service
	originalService := &TestService{Name: "original"}

	// Set the service
	services.TestService = originalService

	// Verify original service is set
	retrieved := services.TestService
	if retrieved.Name != "original" {
		t.Errorf("Expected original name 'original', got '%s'", retrieved.Name)
	}

	// Create replacement service
	replacementService := &TestService{Name: "replacement"}

	// Override the service
	services.TestService = replacementService

	// Verify service was overridden
	retrieved = services.TestService
	if retrieved.Name != "replacement" {
		t.Errorf("Expected overridden name 'replacement', got '%s'", retrieved.Name)
	}

	// Verify it's the same instance
	if retrieved != replacementService {
		t.Error("Retrieved service should be the same instance as the replacement")
	}
}

// CircularServiceA depends on CircularServiceB
type CircularServiceA struct {
	services *TestServices
	Value    string
}

// NewCircularServiceA creates a new CircularServiceA
func NewCircularServiceA(services *TestServices, value string) *CircularServiceA {
	return &CircularServiceA{
		services: services,
		Value:    value,
	}
}

// GetServiceB accesses ServiceB from services
func (s *CircularServiceA) GetServiceB() *CircularServiceB {
	return s.services.CircularServiceB
}

// CircularServiceB depends on CircularServiceA
type CircularServiceB struct {
	services *TestServices
	Value    string
}

// NewCircularServiceB creates a new CircularServiceB
func NewCircularServiceB(services *TestServices, value string) *CircularServiceB {
	return &CircularServiceB{
		services: services,
		Value:    value,
	}
}

// GetServiceA accesses ServiceA from services
func (s *CircularServiceB) GetServiceA() *CircularServiceA {
	return s.services.CircularServiceA
}

func TestCircularDependency(t *testing.T) {
	// Create service container
	services := NewTestServices()

	// Create services
	serviceA := NewCircularServiceA(services, "service-a")
	serviceB := NewCircularServiceB(services, "service-b")

	// Set services
	services.CircularServiceA = serviceA
	services.CircularServiceB = serviceB

	// Verify circular dependency works
	retrievedB := serviceA.GetServiceB()
	if retrievedB == nil {
		t.Error("ServiceA should be able to access ServiceB")
	}

	retrievedA := serviceB.GetServiceA()
	if retrievedA == nil {
		t.Error("ServiceB should be able to access ServiceA")
	}

	// Verify the services reference each other correctly
	if retrievedB != serviceB {
		t.Error("ServiceA.GetServiceB() should return the same instance as serviceB")
	}
	if retrievedA != serviceA {
		t.Error("ServiceB.GetServiceA() should return the same instance as serviceA")
	}

	// Verify original values are preserved
	if serviceA.Value != "service-a" {
		t.Errorf("Expected serviceA value 'service-a', got '%s'", serviceA.Value)
	}
	if serviceB.Value != "service-b" {
		t.Errorf("Expected serviceB value 'service-b', got '%s'", serviceB.Value)
	}

	// Verify we can retrieve services from container
	retrievedAFromContainer := services.CircularServiceA
	if retrievedAFromContainer.Value != "service-a" {
		t.Errorf("Expected retrieved serviceA value 'service-a', got '%s'", retrievedAFromContainer.Value)
	}

	retrievedBFromContainer := services.CircularServiceB
	if retrievedBFromContainer.Value != "service-b" {
		t.Errorf("Expected retrieved serviceB value 'service-b', got '%s'", retrievedBFromContainer.Value)
	}

	// Verify the circular references are maintained in retrieved services
	if retrievedAFromContainer.GetServiceB() != retrievedBFromContainer {
		t.Error("Retrieved serviceA.GetServiceB() should reference the same instance as retrieved serviceB")
	}
	if retrievedBFromContainer.GetServiceA() != retrievedAFromContainer {
		t.Error("Retrieved serviceB.GetServiceA() should reference the same instance as retrieved serviceA")
	}
}
