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
Package optics provides functional optics for composable data access and manipulation.

# Overview

Optics are first-class, composable references to parts of data structures. They provide
a uniform interface for reading, writing, and transforming nested immutable data without
verbose boilerplate code.

The optics package family includes several types of optics, each suited for different
data structure patterns:

  - Lens: Focus on a field within a product type (struct)
  - Prism: Focus on a variant within a sum type (union/Either)
  - Iso: Bidirectional transformation between equivalent types
  - Optional: Focus on a value that may not exist
  - Traversal: Focus on multiple values simultaneously

# Why Optics?

Working with deeply nested immutable data structures in Go can be verbose:

	// Without optics - updating nested data
	updated := Person{
		Name: person.Name,
		Age:  person.Age,
		Address: Address{
			Street: person.Address.Street,
			City:   "New York", // Only this changed!
			Zip:    person.Address.Zip,
		},
	}

With optics, this becomes:

	// With optics - clean and composable
	updated := cityLens.Set("New York")(person)

# Core Optics Types

## Lens - Product Types (Structs)

A Lens focuses on a single field within a struct. It provides get and set operations
that maintain immutability.

	type Person struct {
		Name string
		Age  int
	}

	nameLens := lens.MakeLens(
		func(p Person) string { return p.Name },
		func(p Person, name string) Person {
			p.Name = name
			return p
		},
	)

	person := Person{Name: "Alice", Age: 30}
	name := nameLens.Get(person)           // "Alice"
	updated := nameLens.Set("Bob")(person) // Person{Name: "Bob", Age: 30}

**Use lenses when:**
  - Working with struct fields
  - The field always exists
  - You need both read and write access

## Prism - Sum Types (Variants)

A Prism focuses on one variant of a sum type. It provides optional get (the variant
may not match) and definite set operations.

	type Result interface{ isResult() }
	type Success struct{ Value int }
	type Failure struct{ Error string }

	successPrism := prism.MakePrism(
		func(r Result) option.Option[int] {
			if s, ok := r.(Success); ok {
				return option.Some(s.Value)
			}
			return option.None[int]()
		},
		func(v int) Result { return Success{Value: v} },
	)

	result := Success{Value: 42}
	value := successPrism.GetOption(result) // Some(42)

**Use prisms when:**
  - Working with sum types (Either, Result, etc.)
  - The value may not be the expected variant
  - You need to match on specific cases

## Iso - Isomorphisms

An Iso represents a bidirectional transformation between two equivalent types with
no information loss.

	celsiusToFahrenheit := iso.MakeIso(
		func(c float64) float64 { return c*9/5 + 32 },
		func(f float64) float64 { return (f - 32) * 5 / 9 },
	)

	fahrenheit := celsiusToFahrenheit.Get(20.0)        // 68.0
	celsius := celsiusToFahrenheit.ReverseGet(68.0)    // 20.0

**Use isos when:**
  - Converting between equivalent representations
  - Wrapping/unwrapping newtypes
  - Encoding/decoding data

## Optional - Maybe Values

An Optional focuses on a value that may or may not exist, similar to Option[A].

**Use optionals when:**
  - Working with nullable fields
  - The value may be absent
  - You need to handle the None case

## Traversal - Multiple Values

A Traversal focuses on multiple values simultaneously, allowing batch operations.

**Use traversals when:**
  - Working with collections
  - Updating multiple fields at once
  - Applying transformations to all matching elements

# Composition

The real power of optics comes from composition. Optics of the same or compatible
types can be composed to create more complex accessors:

	type Company struct {
		Name    string
		Address Address
	}

	type Address struct {
		Street string
		City   string
	}

	// Individual lenses
	addressLens := lens.MakeLens(
		func(c Company) Address { return c.Address },
		func(c Company, a Address) Company {
			c.Address = a
			return c
		},
	)

	cityLens := lens.MakeLens(
		func(a Address) string { return a.City },
		func(a Address, city string) Address {
			a.City = city
			return a
		},
	)

	// Compose to access city directly from company
	companyCityLens := F.Pipe1(
		addressLens,
		lens.Compose[Company](cityLens),
	)

	company := Company{
		Name: "Acme Corp",
		Address: Address{Street: "Main St", City: "NYC"},
	}

	city := companyCityLens.Get(company)           // "NYC"
	updated := companyCityLens.Set("Boston")(company)

# Optics Hierarchy

Optics form a hierarchy where more specific optics can be converted to more general ones:

	Iso[S, A]
	    ↓
	Lens[S, A]
	    ↓
	Optional[S, A]
	    ↓
	Traversal[S, A]

	Prism[S, A]
	    ↓
	Optional[S, A]
	    ↓
	Traversal[S, A]

This means:
  - Every Iso is a Lens
  - Every Lens is an Optional
  - Every Prism is an Optional
  - Every Optional is a Traversal

# Laws

Each optic type must satisfy specific laws to ensure correct behavior:

**Lens Laws:**
 1. GetSet: lens.Set(lens.Get(s))(s) == s
 2. SetGet: lens.Get(lens.Set(a)(s)) == a
 3. SetSet: lens.Set(a2)(lens.Set(a1)(s)) == lens.Set(a2)(s)

**Prism Laws:**
 1. GetOptionReverseGet: prism.GetOption(prism.ReverseGet(a)) == Some(a)
 2. ReverseGetGetOption: if GetOption(s) == Some(a), then ReverseGet(a) == s

**Iso Laws:**
 1. RoundTrip1: iso.ReverseGet(iso.Get(s)) == s
 2. RoundTrip2: iso.Get(iso.ReverseGet(a)) == a

# Real-World Example: Configuration Management

	type DatabaseConfig struct {
		Host     string
		Port     int
		Username string
		Password string
	}

	type CacheConfig struct {
		TTL     int
		MaxSize int
	}

	type AppConfig struct {
		Database *DatabaseConfig
		Cache    *CacheConfig
		Debug    bool
	}

	// Create lenses for nested access
	dbLens := lens.FromNillable(lens.MakeLens(
		func(c AppConfig) *DatabaseConfig { return c.Database },
		func(c AppConfig, db *DatabaseConfig) AppConfig {
			c.Database = db
			return c
		},
	))

	dbHostLens := lens.MakeLensRef(
		func(db *DatabaseConfig) string { return db.Host },
		func(db *DatabaseConfig, host string) *DatabaseConfig {
			db.Host = host
			return db
		},
	)

	defaultDB := &DatabaseConfig{
		Host:     "localhost",
		Port:     5432,
		Username: "admin",
		Password: "",
	}

	// Compose to access database host from app config
	appDbHostLens := F.Pipe1(
		dbLens,
		lens.ComposeOption[AppConfig, string](defaultDB)(dbHostLens),
	)

	config := AppConfig{Database: nil, Debug: true}

	// Get returns None when database is not configured
	host := appDbHostLens.Get(config) // None[string]

	// Set creates database with default values
	updated := appDbHostLens.Set(option.Some("prod.example.com"))(config)
	// updated.Database.Host == "prod.example.com"
	// updated.Database.Port == 5432 (from default)

# Package Structure

The optics package is organized into subpackages:

  - optics/lens: Lenses for product types
  - optics/prism: Prisms for sum types
  - optics/iso: Isomorphisms for equivalent types
  - optics/optional: Optional optics for maybe values
  - optics/traversal: Traversals for multiple values

Each subpackage may have additional specialized subpackages for common patterns:
  - array: Optics for array/slice operations
  - either: Optics for Either types
  - option: Optics for Option types
  - record: Optics for record/map types

# Performance Considerations

Optics are designed to be efficient:
  - No reflection - all operations are type-safe at compile time
  - Minimal allocations - optics themselves are lightweight
  - Composition is efficient - creates function closures
  - Immutability ensures thread safety

For performance-critical code:
  - Cache composed optics rather than recomposing
  - Use pointer-based lenses (MakeLensRef) for large structs
  - Consider batch operations with traversals

# Type Safety

All optics are fully type-safe:
  - Compile-time type checking
  - No runtime type assertions
  - Generic type parameters ensure correctness
  - Composition maintains type relationships

# Getting Started

1. Choose the right optic for your data structure
2. Create basic optics for your types
3. Compose optics for nested access
4. Use Modify for transformations
5. Leverage the optics hierarchy when needed

# Further Reading

For detailed documentation on each optic type, see:
  - github.com/IBM/fp-go/v2/optics/lens
  - github.com/IBM/fp-go/v2/optics/prism
  - github.com/IBM/fp-go/v2/optics/iso
  - github.com/IBM/fp-go/v2/optics/optional
  - github.com/IBM/fp-go/v2/optics/traversal

For related functional programming concepts:
  - github.com/IBM/fp-go/v2/option: Optional values
  - github.com/IBM/fp-go/v2/either: Sum types
  - github.com/IBM/fp-go/v2/function: Function composition
*/
package optics
