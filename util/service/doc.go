// Copyright 2026 TypeFox GmbH
// This program and the accompanying materials are made available under the
// terms of the MIT License, which is available in the project root.

// Package service provides a typed dependency injection container used across
// Fastbelt. Framework packages and language implementations register shared
// components (parsers, linkers, workspace builders, LSP handlers, etc.) in a
// [Container] keyed by Go type, then retrieve them with [Get] or [MustGet]
// after the container is sealed.
//
// # Container lifecycle
//
// DI containers have two phases. During setup, register services with [Put]
// and replace existing registrations with [Override]. Call [Container.Seal]
// when wiring is complete; after that, [Put] and [Override] panic and [Get]
// becomes available. [Has] reports whether a type is registered and may be
// called before sealing.
//
// A typical factory helper looks like:
//
//	func CreateServices() *service.Container {
//	    sc := service.NewContainer()
//	    SetupServices(sc)
//	    sc.Seal()
//	    return sc
//	}
//
// # SetupDefaultServices in framework packages
//
// Each key package of Fastbelt exposes a SetupDefaultServices function that
// registers its default implementations into a container:
//
//   - [typefox.dev/fastbelt/textdoc.SetupDefaultServices]
//   - [typefox.dev/fastbelt/linking.SetupDefaultServices]
//   - [typefox.dev/fastbelt/workspace.SetupDefaultServices]
//   - [typefox.dev/fastbelt/server.SetupDefaultServices]
//
// These functions are idempotent: they call [Has] before [Put] and skip types
// that are already registered. A language can pre-register custom
// implementations before calling them, or call [Override] afterward to replace
// specific defaults while still using the framework defaults for everything
// else.
//
// # SetupServices and CreateServices in language implementations
//
// Each language project defines SetupServices to wire a complete container
// for that language. By convention, SetupServices:
//
//   - registers language-specific configuration such as
//     [typefox.dev/fastbelt/workspace.LanguageID] and
//     [typefox.dev/fastbelt/workspace.FileExtensions],
//   - calls the framework SetupDefaultServices functions listed above,
//   - calls SetupGeneratedServices from generated code (registers the lexer,
//     parser, scope provider, and other grammar-specific services), and
//   - optionally replaces defaults with [Override] for language-specific
//     behavior.
//
// CreateServices is the usual entry point for CLI tools and tests: it creates
// a container with [NewContainer], calls SetupServices, seals the container,
// and returns it. Language servers follow the same pattern but also call
// SetupGeneratedServerServices and [typefox.dev/fastbelt/server.SetupDefaultServices]
// before sealing; many projects expose a separate CreateLspServices helper for
// that path.
//
// Scaffolded language projects include a generated services.go that follows
// this layout and can be extended for customization.
package service
