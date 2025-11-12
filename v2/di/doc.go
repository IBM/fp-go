// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*
Package di implements functions and data types supporting dependency injection patterns.

# Overview

The dependency injection (DI) framework provides a type-safe way to manage dependencies
between components in your application. It ensures that all dependencies are resolved
correctly and that instances are created as singletons.

# Core Concepts

Dependency - An abstract concept representing a service, function, or value with a specific type.
Dependencies can be:
  - Simple values (API URLs, configuration strings, numbers)
  - Complex objects (HTTP clients, database connections)
  - Service interfaces
  - Functions

InjectionToken - A unique identifier for a dependency that includes type information.
Created using MakeToken[T](name).

Provider - The implementation of a dependency. Providers specify:
  - Which dependency they provide (via an InjectionToken)
  - What other dependencies they need
  - How to create the instance

InjectableFactory - The container that manages all providers and resolves dependencies.
All resolved instances are singletons.

# Basic Usage

Creating and using dependencies:

	import (
		"github.com/IBM/fp-go/v2/di"
		IOE "github.com/IBM/fp-go/v2/ioeither"
	)

	// Define injection tokens
	var (
		ConfigToken = di.MakeToken[Config]("Config")
		DBToken     = di.MakeToken[Database]("Database")
		APIToken    = di.MakeToken[APIService]("APIService")
	)

	// Create providers
	configProvider := di.ConstProvider(ConfigToken, Config{Port: 8080})

	dbProvider := di.MakeProvider1(
		DBToken,
		ConfigToken.Identity(),
		func(cfg Config) IOResult[Database] {
			return ioresult.Of(NewDatabase(cfg))
		},
	)

	apiProvider := di.MakeProvider2(
		APIToken,
		ConfigToken.Identity(),
		DBToken.Identity(),
		func(cfg Config, db Database) IOResult[APIService] {
			return ioresult.Of(NewAPIService(cfg, db))
		},
	)

	// Create injector and resolve
	injector := DIE.MakeInjector([]DIE.Provider{
		configProvider,
		dbProvider,
		apiProvider,
	})

	// Resolve the API service
	resolver := di.Resolve(APIToken)
	result := resolver(injector)()

# Dependency Types

Identity (Required) - The dependency must be resolved, or the injection fails:

	token.Identity() // Returns Dependency[T]

Option (Optional) - The dependency is optional, returns Option[T]:

	token.Option() // Returns Dependency[Option[T]]

IOEither (Lazy Required) - Lazy evaluation, memoized singleton:

	token.IOEither() // Returns Dependency[IOEither[error, T]]

IOOption (Lazy Optional) - Lazy optional evaluation:

	token.IOOption() // Returns Dependency[IOOption[T]]

# Provider Creation

Providers are created using MakeProvider functions with suffixes indicating
the number of dependencies (0-15):

MakeProvider0 - No dependencies:

	provider := di.MakeProvider0(
		token,
		ioresult.Of(value),
	)

MakeProvider1 - One dependency:

	provider := di.MakeProvider1(
		resultToken,
		dep1Token.Identity(),
		func(dep1 Dep1Type) IOResult[ResultType] {
			return ioresult.Of(createResult(dep1))
		},
	)

MakeProvider2 - Two dependencies:

	provider := di.MakeProvider2(
		resultToken,
		dep1Token.Identity(),
		dep2Token.Identity(),
		func(dep1 Dep1Type, dep2 Dep2Type) IOResult[ResultType] {
			return ioresult.Of(createResult(dep1, dep2))
		},
	)

# Constant Providers

For simple constant values:

	provider := di.ConstProvider(token, value)

# Default Implementations

Tokens can have default implementations that are used when no explicit
provider is registered:

	token := di.MakeTokenWithDefault0(
		"ServiceName",
		ioresult.Of(defaultImplementation),
	)

	// Or with dependencies
	token := di.MakeTokenWithDefault2(
		"ServiceName",
		dep1Token.Identity(),
		dep2Token.Identity(),
		func(dep1 Dep1Type, dep2 Dep2Type) IOResult[ResultType] {
			return ioresult.Of(createDefault(dep1, dep2))
		},
	)

# Multi-Value Dependencies

For dependencies that can have multiple implementations:

	// Create a multi-token
	loggersToken := di.MakeMultiToken[Logger]("Loggers")

	// Provide multiple items
	consoleLogger := di.ConstProvider(loggersToken.Item(), ConsoleLogger{})
	fileLogger := di.ConstProvider(loggersToken.Item(), FileLogger{})

	// Resolve all items as an array
	resolver := di.Resolve(loggersToken.Container())
	loggers := resolver(injector)() // Returns []Logger

# Lazy vs Eager Resolution

Eager (Identity/Option) - Resolved immediately when the injector is created:

	dep1Token.Identity() // Resolved eagerly
	dep2Token.Option()   // Resolved eagerly

Lazy (IOEither/IOOption) - Resolved only when accessed:

	dep3Token.IOEither() // Resolved lazily
	dep4Token.IOOption() // Resolved lazily

Lazy dependencies are memoized, so they're only created once.

# Main Application Pattern

The framework provides a convenient pattern for running applications:

	import (
		"github.com/IBM/fp-go/v2/di"
		IOE "github.com/IBM/fp-go/v2/ioeither"
	)

	// Define your main application logic
	mainProvider := di.MakeProvider1(
		di.InjMain,
		APIToken.Identity(),
		func(api APIService) IOResult[any] {
			return ioresult.Of(api.Start())
		},
	)

	// Run the application
	err := di.RunMain([]DIE.Provider{
		configProvider,
		dbProvider,
		apiProvider,
		mainProvider,
	})()

# Practical Examples

Example 1: Configuration-based Service

	type Config struct {
		APIKey string
		Timeout int
	}

	type HTTPClient struct {
		config Config
	}

	var (
		ConfigToken = di.MakeToken[Config]("Config")
		ClientToken = di.MakeToken[HTTPClient]("HTTPClient")
	)

	configProvider := di.ConstProvider(ConfigToken, Config{
		APIKey: "secret",
		Timeout: 30,
	})

	clientProvider := di.MakeProvider1(
		ClientToken,
		ConfigToken.Identity(),
		func(cfg Config) IOResult[HTTPClient] {
			return ioresult.Of(HTTPClient{config: cfg})
		},
	)

Example 2: Optional Dependencies

	var (
		CacheToken  = di.MakeToken[Cache]("Cache")
		ServiceToken = di.MakeToken[Service]("Service")
	)

	// Service works with or without cache
	serviceProvider := di.MakeProvider1(
		ServiceToken,
		CacheToken.Option(), // Optional dependency
		func(cache Option[Cache]) IOResult[Service] {
			return ioresult.Of(NewService(cache))
		},
	)

Example 3: Lazy Dependencies

	var (
		DBToken      = di.MakeToken[Database]("Database")
		ReporterToken = di.MakeToken[Reporter]("Reporter")
	)

	// Reporter only connects to DB when needed
	reporterProvider := di.MakeProvider1(
		ReporterToken,
		DBToken.IOEither(), // Lazy dependency
		func(dbIO IOResult[Database]) IOResult[Reporter] {
			return ioresult.Of(NewReporter(dbIO))
		},
	)

# Function Reference

Token Creation:
  - MakeToken[T](name) InjectionToken[T] - Creates a unique injection token
  - MakeTokenWithDefault[T](name, factory) InjectionToken[T] - Token with default implementation
  - MakeTokenWithDefault0-15 - Token with default and N dependencies
  - MakeMultiToken[T](name) MultiInjectionToken[T] - Token for multiple implementations

Provider Creation:
  - ConstProvider[T](token, value) Provider - Simple constant provider
  - MakeProvider0[R](token, factory) Provider - Provider with no dependencies
  - MakeProvider1-15 - Providers with 1-15 dependencies
  - MakeProviderFactory0-15 - Lower-level factory creation

Resolution:
  - Resolve[T](token) ReaderIOEither[InjectableFactory, error, T] - Resolves a dependency

Application:
  - InjMain - Injection token for main application
  - Main - Resolver for main application
  - RunMain(providers) IO[error] - Runs the main application

Utility:
  - asDependency[T](t) Dependency - Converts to dependency interface

# Related Packages

  - github.com/IBM/fp-go/v2/di/erasure - Type-erased DI implementation
  - github.com/IBM/fp-go/v2/ioeither - IO operations with error handling
  - github.com/IBM/fp-go/v2/option - Optional values
  - github.com/IBM/fp-go/v2/either - Either monad for error handling

[Provider]: [github.com/IBM/fp-go/v2/di/erasure.Provider]
[InjectableFactory]: [github.com/IBM/fp-go/v2/di/erasure.InjectableFactory]
[MakeInjector]: [github.com/IBM/fp-go/v2/di/erasure.MakeInjector]
*/
package di

//go:generate go run .. di --count 15 --filename gen.go
