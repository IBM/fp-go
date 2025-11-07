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

// Package option provides utilities for working with lenses that focus on optional values.
//
// This package extends the lens optics pattern to handle Option types, enabling safe
// manipulation of potentially absent values in nested data structures. It provides
// functions for creating, composing, and transforming lenses that work with optional
// fields.
//
// # Core Concepts
//
// A LensO[S, A] is a Lens[S, Option[A]] - a lens that focuses on an optional value A
// within a structure S. This is particularly useful when dealing with nullable pointers,
// optional fields, or values that may not always be present.
//
// # Key Functions
//
// Creating Lenses from Optional Values:
//   - FromNillable: Creates a lens from a nullable pointer field
//   - FromNillableRef: Pointer-based version of FromNillable
//   - FromPredicate: Creates a lens based on a predicate function
//   - FromPredicateRef: Pointer-based version of FromPredicate
//   - FromOption: Converts an optional lens to a definite lens with a default value
//   - FromOptionRef: Pointer-based version of FromOption
//   - FromNullableProp: Creates a lens with a default value for nullable properties
//   - FromNullablePropRef: Pointer-based version of FromNullableProp
//
// Composing Lenses:
//   - ComposeOption: Composes a lens returning Option[A] with a lens returning B
//   - ComposeOptions: Composes two lenses that both return optional values
//
// Conversions:
//   - AsTraversal: Converts a lens to a traversal for use with traversal operations
//
// # Usage Examples
//
// Working with nullable pointers:
//
//	type Config struct {
//	    Database *DatabaseConfig
//	}
//
//	type DatabaseConfig struct {
//	    Host string
//	    Port int
//	}
//
//	// Create a lens for the optional database config
//	dbLens := lens.FromNillable(lens.MakeLens(
//	    func(c Config) *DatabaseConfig { return c.Database },
//	    func(c Config, db *DatabaseConfig) Config { c.Database = db; return c },
//	))
//
//	// Access the optional value
//	config := Config{Database: nil}
//	dbOpt := dbLens.Get(config) // Returns None[*DatabaseConfig]
//
//	// Set a value
//	newDB := &DatabaseConfig{Host: "localhost", Port: 5432}
//	updated := dbLens.Set(O.Some(newDB))(config)
//
// Composing optional lenses:
//
//	// Lens to access port through optional database
//	portLens := lens.MakeLensRef(
//	    func(db *DatabaseConfig) int { return db.Port },
//	    func(db *DatabaseConfig, port int) *DatabaseConfig { db.Port = port; return db },
//	)
//
//	defaultDB := &DatabaseConfig{Host: "localhost", Port: 5432}
//	configPortLens := F.Pipe1(dbLens,
//	    lens.ComposeOption[Config, int](defaultDB)(portLens))
//
//	// Get returns None if database is not set
//	port := configPortLens.Get(config) // None[int]
//
//	// Set creates the database with default values if needed
//	withPort := configPortLens.Set(O.Some(3306))(config)
//	// withPort.Database.Port == 3306, Host == "localhost"
//
// Working with predicates:
//
//	type Person struct {
//	    Age int
//	}
//
//	ageLens := lens.MakeLensRef(
//	    func(p *Person) int { return p.Age },
//	    func(p *Person, age int) *Person { p.Age = age; return p },
//	)
//
//	// Only consider adults (age >= 18)
//	adultLens := lens.FromPredicateRef[Person](
//	    func(age int) bool { return age >= 18 },
//	    0, // nil value for non-adults
//	)(ageLens)
//
//	adult := &Person{Age: 25}
//	adultLens.Get(adult) // Some(25)
//
//	minor := &Person{Age: 15}
//	adultLens.Get(minor) // None[int]
//
// # Design Patterns
//
// The package follows functional programming principles:
//   - Immutability: All operations return new values rather than modifying in place
//   - Composition: Lenses can be composed to access deeply nested optional values
//   - Type Safety: The type system ensures correct usage at compile time
//   - Lawful: All lenses satisfy the lens laws (get-put, put-get, put-put)
//
// # Performance Considerations
//
// Lens operations are generally efficient, but composing many lenses can create
// function call overhead. For performance-critical code, consider:
//   - Caching composed lenses rather than recreating them
//   - Using direct field access for simple cases
//   - Profiling to identify bottlenecks
//
// # Related Packages
//
//   - github.com/IBM/fp-go/v2/optics/lens: Core lens functionality
//   - github.com/IBM/fp-go/v2/option: Option type and operations
//   - github.com/IBM/fp-go/v2/optics/traversal/option: Traversals for optional values
package option
