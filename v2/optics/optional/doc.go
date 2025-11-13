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
Package optional provides optional optics for focusing on values that may not exist.

# Overview

An Optional is an optic that focuses on a subpart of a data structure that may or may not
be present. Unlike lenses which always focus on an existing field, optionals handle cases
where the target value might be absent, returning Option[A] instead of A.

Optionals are the bridge between lenses (which always succeed) and prisms (which may fail
to match). They combine aspects of both:
  - Like lenses: Focus on a specific location in a structure
  - Like prisms: The value at that location may not exist

Optionals are essential for:
  - Working with nullable fields (pointers that may be nil)
  - Accessing nested optional values
  - Conditional updates based on value presence
  - Safe navigation through potentially missing data

# Mathematical Foundation

An Optional[S, A] consists of two operations:
  - GetOption: S → Option[A] (try to extract A from S, may return None)
  - Set: A → S → S (update A in S, may be a no-op if value doesn't exist)

Optionals must satisfy the optional laws:
 1. GetOptionSet: if GetOption(s) == Some(a), then GetOption(Set(a)(s)) == Some(a)
 2. SetGetOption: if GetOption(s) == Some(a), then Set(a)(s) preserves other parts of s
 3. SetSet: Set(a2)(Set(a1)(s)) == Set(a2)(s)

# Basic Usage

Creating an optional for a nullable field:

	type Config struct {
		Timeout *int
		MaxSize *int
	}

	timeoutOptional := optional.MakeOptional(
		func(c Config) option.Option[*int] {
			return option.FromNillable(c.Timeout)
		},
		func(c Config, t *int) Config {
			c.Timeout = t
			return c
		},
	)

	config := Config{Timeout: nil, MaxSize: ptr(100)}

	// Get returns None for nil
	timeout := timeoutOptional.GetOption(config) // None[*int]

	// Set updates the value
	newTimeout := 30
	updated := timeoutOptional.Set(&newTimeout)(config)
	// updated.Timeout points to 30

# Working with Pointers

For pointer-based structures, use MakeOptionalRef which handles copying automatically:

	type Database struct {
		Host string
		Port int
	}

	type Config struct {
		Database *Database
	}

	dbOptional := optional.MakeOptionalRef(
		func(c *Config) option.Option[*Database] {
			return option.FromNillable(c.Database)
		},
		func(c *Config, db *Database) *Config {
			c.Database = db
			return c
		},
	)

	config := &Config{Database: nil}

	// Get returns None when database is nil
	db := dbOptional.GetOption(config) // None[*Database]

	// Set creates a new config with the database
	newDB := &Database{Host: "localhost", Port: 5432}
	updated := dbOptional.Set(newDB)(config)
	// config.Database is still nil, updated.Database points to newDB

# Identity Optional

The identity optional focuses on the entire structure:

	idOpt := optional.Id[Config]()

	config := Config{Timeout: ptr(30)}
	value := idOpt.GetOption(config) // Some(config)
	updated := idOpt.Set(Config{Timeout: ptr(60)})(config)

# Composing Optionals

Optionals can be composed to navigate through nested optional structures:

	type Address struct {
		Street string
		City   string
	}

	type Person struct {
		Name    string
		Address *Address
	}

	addressOpt := optional.MakeOptional(
		func(p Person) option.Option[*Address] {
			return option.FromNillable(p.Address)
		},
		func(p Person, a *Address) Person {
			p.Address = a
			return p
		},
	)

	cityOpt := optional.MakeOptionalRef(
		func(a *Address) option.Option[string] {
			if a == nil {
				return option.None[string]()
			}
			return option.Some(a.City)
		},
		func(a *Address, city string) *Address {
			a.City = city
			return a
		},
	)

	// Compose to access city from person
	personCityOpt := F.Pipe1(
		addressOpt,
		optional.Compose[Person, *Address, string](cityOpt),
	)

	person := Person{Name: "Alice", Address: nil}

	// Get returns None when address is nil
	city := personCityOpt.GetOption(person) // None[string]

	// Set updates the city if address exists
	withAddress := Person{
		Name:    "Alice",
		Address: &Address{Street: "Main St", City: "NYC"},
	}
	updated := personCityOpt.Set("Boston")(withAddress)
	// updated.Address.City == "Boston"

# From Predicate

Create an optional that only focuses on values satisfying a predicate:

	type User struct {
		Age int
	}

	ageOpt := optional.FromPredicate[User, int](
		func(age int) bool { return age >= 18 },
	)(
		func(u User) int { return u.Age },
		func(u User, age int) User {
			u.Age = age
			return u
		},
	)

	adult := User{Age: 25}
	age := ageOpt.GetOption(adult) // Some(25)

	minor := User{Age: 15}
	minorAge := ageOpt.GetOption(minor) // None[int]

	// Set only works if predicate is satisfied
	updated := ageOpt.Set(30)(adult) // Age becomes 30
	unchanged := ageOpt.Set(30)(minor) // Age stays 15 (predicate fails)

# Modifying Values

Use ModifyOption to transform values that exist:

	type Counter struct {
		Value *int
	}

	valueOpt := optional.MakeOptional(
		func(c Counter) option.Option[*int] {
			return option.FromNillable(c.Value)
		},
		func(c Counter, v *int) Counter {
			c.Value = v
			return c
		},
	)

	counter := Counter{Value: ptr(5)}

	// Increment if value exists
	incremented := F.Pipe3(
		counter,
		valueOpt,
		optional.ModifyOption[Counter, *int](func(v *int) *int {
			newVal := *v + 1
			return &newVal
		}),
		option.GetOrElse(F.Constant(counter)),
	)
	// incremented.Value points to 6

	// No change if value is nil
	nilCounter := Counter{Value: nil}
	result := F.Pipe3(
		nilCounter,
		valueOpt,
		optional.ModifyOption[Counter, *int](func(v *int) *int {
			newVal := *v + 1
			return &newVal
		}),
		option.GetOrElse(F.Constant(nilCounter)),
	)
	// result.Value is still nil

# Bidirectional Mapping

Transform the focus type of an optional:

	type Celsius float64
	type Fahrenheit float64

	type Weather struct {
		Temperature *Celsius
	}

	tempCelsiusOpt := optional.MakeOptional(
		func(w Weather) option.Option[*Celsius] {
			return option.FromNillable(w.Temperature)
		},
		func(w Weather, t *Celsius) Weather {
			w.Temperature = t
			return w
		},
	)

	// Create optional that works with Fahrenheit
	tempFahrenheitOpt := F.Pipe1(
		tempCelsiusOpt,
		optional.IMap[Weather, *Celsius, *Fahrenheit](
			func(c *Celsius) *Fahrenheit {
				f := Fahrenheit(*c*9/5 + 32)
				return &f
			},
			func(f *Fahrenheit) *Celsius {
				c := Celsius((*f - 32) * 5 / 9)
				return &c
			},
		),
	)

	celsius := Celsius(20)
	weather := Weather{Temperature: &celsius}

	tempF := tempFahrenheitOpt.GetOption(weather) // Some(68°F)

# Real-World Example: Configuration with Defaults

	type DatabaseConfig struct {
		Host     string
		Port     int
		Username string
		Password string
	}

	type AppConfig struct {
		Database *DatabaseConfig
		Debug    bool
	}

	dbOpt := optional.MakeOptional(
		func(c AppConfig) option.Option[*DatabaseConfig] {
			return option.FromNillable(c.Database)
		},
		func(c AppConfig, db *DatabaseConfig) AppConfig {
			c.Database = db
			return c
		},
	)

	dbHostOpt := optional.MakeOptionalRef(
		func(db *DatabaseConfig) option.Option[string] {
			if db == nil {
				return option.None[string]()
			}
			return option.Some(db.Host)
		},
		func(db *DatabaseConfig, host string) *DatabaseConfig {
			db.Host = host
			return db
		},
	)

	// Compose to access database host
	appDbHostOpt := F.Pipe1(
		dbOpt,
		optional.Compose[AppConfig, *DatabaseConfig, string](dbHostOpt),
	)

	config := AppConfig{Database: nil, Debug: true}

	// Get returns None when database is not configured
	host := appDbHostOpt.GetOption(config) // None[string]

	// Set creates database if needed
	withDB := AppConfig{
		Database: &DatabaseConfig{Host: "localhost", Port: 5432},
		Debug:    true,
	}
	updated := appDbHostOpt.Set("prod.example.com")(withDB)
	// updated.Database.Host == "prod.example.com"

# Real-World Example: Safe Navigation

	type Company struct {
		Name string
		CEO  *Person
	}

	type Person struct {
		Name    string
		Address *Address
	}

	type Address struct {
		City string
	}

	ceoOpt := optional.MakeOptional(
		func(c Company) option.Option[*Person] {
			return option.FromNillable(c.CEO)
		},
		func(c Company, p *Person) Company {
			c.CEO = p
			return c
		},
	)

	addressOpt := optional.MakeOptionalRef(
		func(p *Person) option.Option[*Address] {
			return option.FromNillable(p.Address)
		},
		func(p *Person, a *Address) *Person {
			p.Address = a
			return p
		},
	)

	cityOpt := optional.MakeOptionalRef(
		func(a *Address) option.Option[string] {
			if a == nil {
				return option.None[string]()
			}
			return option.Some(a.City)
		},
		func(a *Address, city string) *Address {
			a.City = city
			return a
		},
	)

	// Compose all optionals for safe navigation
	ceoCityOpt := F.Pipe2(
		ceoOpt,
		optional.Compose[Company, *Person, *Address](addressOpt),
		optional.Compose[Company, *Address, string](cityOpt),
	)

	company := Company{Name: "Acme Corp", CEO: nil}

	// Safe navigation returns None at any missing level
	city := ceoCityOpt.GetOption(company) // None[string]

# Optionals in the Optics Hierarchy

Optionals sit between lenses and traversals in the optics hierarchy:

	Lens[S, A]
	    ↓
	Optional[S, A]
	    ↓
	Traversal[S, A]

	Prism[S, A]
	    ↓
	Optional[S, A]

This means:
  - Every Lens can be converted to an Optional (value always exists)
  - Every Prism can be converted to an Optional (variant may not match)
  - Every Optional can be converted to a Traversal (0 or 1 values)

# Performance Considerations

Optionals are efficient:
  - No reflection - all operations are type-safe at compile time
  - Minimal allocations - optionals themselves are lightweight
  - GetOption short-circuits on None
  - Set operations create new copies (immutability)

For best performance:
  - Use MakeOptionalRef for pointer structures to ensure proper copying
  - Cache composed optionals rather than recomposing
  - Consider batch operations when updating multiple optional values

# Type Safety

Optionals are fully type-safe:
  - Compile-time type checking
  - No runtime type assertions
  - Generic type parameters ensure correctness
  - Composition maintains type relationships

# Function Reference

Core Optional Creation:
  - MakeOptional: Create an optional from getter and setter functions
  - MakeOptionalRef: Create an optional for pointer-based structures
  - Id: Create an identity optional
  - IdRef: Create an identity optional for pointers

Composition:
  - Compose: Compose two optionals
  - ComposeRef: Compose optionals for pointer structures

Transformation:
  - ModifyOption: Transform a value through an optional (returns Option[S])
  - SetOption: Set a value through an optional (returns Option[S])
  - IMap: Bidirectionally map an optional
  - IChain: Bidirectionally map with optional results
  - IChainAny: Map to/from any type

Predicate-Based:
  - FromPredicate: Create optional from predicate
  - FromPredicateRef: Create optional from predicate (ref version)

# Related Packages

  - github.com/IBM/fp-go/v2/optics/lens: Lenses for fields that always exist
  - github.com/IBM/fp-go/v2/optics/prism: Prisms for sum types
  - github.com/IBM/fp-go/v2/optics/traversal: Traversals for multiple values
  - github.com/IBM/fp-go/v2/option: Optional values
  - github.com/IBM/fp-go/v2/endomorphism: Endomorphisms (A → A functions)
*/
package optional
