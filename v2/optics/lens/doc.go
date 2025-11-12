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
Package lens provides functional optics for zooming into and modifying nested data structures.

# Overview

A Lens is a first-class reference to a subpart of a data structure. It provides a composable
way to focus on a particular field within a nested structure, allowing you to get and set
values in an immutable, functional manner.

Lenses are particularly useful when working with deeply nested immutable data structures,
as they eliminate the need for verbose copying and updating code.

# Mathematical Foundation

A Lens[S, A] is defined by two operations:
  - Get: S → A (extract a value of type A from a structure of type S)
  - Set: A → S → S (update the value of type A in structure S, returning a new S)

Lenses must satisfy the lens laws:
 1. GetSet: lens.Set(lens.Get(s))(s) == s
 2. SetGet: lens.Get(lens.Set(a)(s)) == a
 3. SetSet: lens.Set(a2)(lens.Set(a1)(s)) == lens.Set(a2)(s)

# Basic Usage

Creating a lens for a struct field:

	type Person struct {
		Name string
		Age  int
	}

	// Create a lens for the Name field
	nameLens := lens.MakeLens(
		func(p Person) string { return p.Name },
		func(p Person, name string) Person {
			p.Name = name
			return p
		},
	)

	person := Person{Name: "Alice", Age: 30}

	// Get the name
	name := nameLens.Get(person) // "Alice"

	// Set a new name (returns a new Person)
	updated := nameLens.Set("Bob")(person)
	// person.Name is still "Alice", updated.Name is "Bob"

# Composing Lenses

Lenses can be composed to focus on deeply nested fields:

	type Address struct {
		Street string
		City   string
	}

	type Person struct {
		Name    string
		Address Address
	}

	addressLens := lens.MakeLens(
		func(p Person) Address { return p.Address },
		func(p Person, a Address) Person {
			p.Address = a
			return p
		},
	)

	streetLens := lens.MakeLens(
		func(a Address) string { return a.Street },
		func(a Address, s string) Address {
			a.Street = s
			return a
		},
	)

	// Compose to access street directly from person
	personStreetLens := F.Pipe1(
		addressLens,
		lens.Compose[Person](streetLens),
	)

	person := Person{
		Name: "Alice",
		Address: Address{Street: "Main St", City: "NYC"},
	}

	street := personStreetLens.Get(person) // "Main St"
	updated := personStreetLens.Set("Oak Ave")(person)

# Working with Pointers

For pointer-based structures, use MakeLensRef which handles copying automatically:

	type Person struct {
		Name string
		Age  int
	}

	func (p *Person) GetName() string {
		return p.Name
	}

	func (p *Person) SetName(name string) *Person {
		p.Name = name
		return p
	}

	// MakeLensRef handles pointer copying
	nameLens := lens.MakeLensRef(
		(*Person).GetName,
		(*Person).SetName,
	)

	person := &Person{Name: "Alice", Age: 30}
	updated := nameLens.Set("Bob")(person)
	// person.Name is still "Alice", updated is a new pointer

# Optional Values

Lenses can work with optional values using Option types:

	type Config struct {
		Port    *int
		Timeout *int
	}

	portLens := lens.MakeLens(
		func(c Config) *int { return c.Port },
		func(c Config, p *int) Config {
			c.Port = p
			return c
		},
	)

	// Convert to optional lens
	optPortLens := lens.FromNillable(portLens)

	config := Config{Port: nil}

	// Get returns None for nil
	port := optPortLens.Get(config) // None[*int]

	// Set with Some updates the value
	newPort := 8080
	updated := optPortLens.Set(O.Some(&newPort))(config)

	// Set with None removes the value
	cleared := optPortLens.Set(O.None[*int]())(updated)

# Composing with Optional Values

ComposeOption allows composing a lens returning an optional value with a regular lens:

	type Database struct {
		Host string
		Port int
	}

	type Config struct {
		Database *Database
	}

	dbLens := lens.FromNillable(lens.MakeLens(
		func(c Config) *Database { return c.Database },
		func(c Config, db *Database) Config {
			c.Database = db
			return c
		},
	))

	portLens := lens.MakeLensRef(
		func(db *Database) int { return db.Port },
		func(db *Database, port int) *Database {
			db.Port = port
			return db
		},
	)

	defaultDB := &Database{Host: "localhost", Port: 5432}

	// Compose with default value
	configPortLens := F.Pipe1(
		dbLens,
		lens.ComposeOption[Config, int](defaultDB)(portLens),
	)

	config := Config{Database: nil}

	// Get returns None when database is nil
	port := configPortLens.Get(config) // None[int]

	// Set creates database with default values
	updated := configPortLens.Set(O.Some(3306))(config)
	// updated.Database.Port == 3306, Host == "localhost"

# Modifying Values

Use Modify to transform a value through a lens:

	type Counter struct {
		Value int
	}

	valueLens := lens.MakeLens(
		func(c Counter) int { return c.Value },
		func(c Counter, v int) Counter {
			c.Value = v
			return c
		},
	)

	counter := Counter{Value: 5}

	// Increment the counter
	incremented := F.Pipe2(
		counter,
		valueLens,
		lens.Modify[Counter](func(v int) int { return v + 1 }),
	)
	// incremented.Value == 6

# Identity Lens

The identity lens focuses on the entire structure:

	idLens := lens.Id[Person]()

	person := Person{Name: "Alice", Age: 30}
	same := idLens.Get(person) // returns person
	updated := idLens.Set(Person{Name: "Bob", Age: 25})(person)

# Isomorphic Mapping

IMap allows you to transform the focus type of a lens:

	type Celsius float64
	type Fahrenheit float64

	celsiusToFahrenheit := func(c Celsius) Fahrenheit {
		return Fahrenheit(c*9/5 + 32)
	}

	fahrenheitToCelsius := func(f Fahrenheit) Celsius {
		return Celsius((f - 32) * 5 / 9)
	}

	type Weather struct {
		Temperature Celsius
	}

	tempCelsiusLens := lens.MakeLens(
		func(w Weather) Celsius { return w.Temperature },
		func(w Weather, t Celsius) Weather {
			w.Temperature = t
			return w
		},
	)

	// Create a lens that works with Fahrenheit
	tempFahrenheitLens := F.Pipe1(
		tempCelsiusLens,
		lens.IMap[Weather](celsiusToFahrenheit, fahrenheitToCelsius),
	)

	weather := Weather{Temperature: 20} // 20°C
	tempF := tempFahrenheitLens.Get(weather) // 68°F
	updated := tempFahrenheitLens.Set(86)(weather) // Set to 86°F (30°C)

# Nullable Properties

FromNullableProp creates a lens that provides a default value for nullable properties:

	type Config struct {
		Timeout *int
	}

	timeoutLens := lens.MakeLens(
		func(c Config) *int { return c.Timeout },
		func(c Config, t *int) Config {
			c.Timeout = t
			return c
		},
	)

	// Provide default value of 30 for nil timeout
	safeTimeoutLens := F.Pipe1(
		timeoutLens,
		lens.FromNullableProp[Config](
			O.FromNillable[int],
			func() *int { v := 30; return &v }(),
		),
	)

	config := Config{Timeout: nil}
	timeout := safeTimeoutLens.Get(config) // returns pointer to 30

# FromOption

FromOption creates a lens from an Option property, providing a default value:

	type Settings struct {
		MaxRetries O.Option[int]
	}

	retriesLens := lens.MakeLens(
		func(s Settings) O.Option[int] { return s.MaxRetries },
		func(s Settings, r O.Option[int]) Settings {
			s.MaxRetries = r
			return s
		},
	)

	// Provide default of 3 retries
	safeRetriesLens := F.Pipe1(
		retriesLens,
		lens.FromOption[Settings](3),
	)

	settings := Settings{MaxRetries: O.None[int]()}
	retries := safeRetriesLens.Get(settings) // returns 3
	updated := safeRetriesLens.Set(5)(settings) // sets to Some(5)

# Predicate-Based Lenses

FromPredicate creates an optional lens based on a predicate:

	type User struct {
		Age int
	}

	ageLens := lens.MakeLens(
		func(u User) int { return u.Age },
		func(u User, age int) User {
			u.Age = age
			return u
		},
	)

	// Only consider valid ages (18+)
	adultAgeLens := F.Pipe1(
		ageLens,
		lens.FromPredicate[User](func(age int) bool {
			return age >= 18
		}, 0),
	)

	user := User{Age: 25}
	age := adultAgeLens.Get(user) // Some(25)

	minor := User{Age: 15}
	minorAge := adultAgeLens.Get(minor) // None[int]

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

	// Create lenses for each level
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
	updated := appDbHostLens.Set(O.Some("prod.example.com"))(config)
	// updated.Database.Host == "prod.example.com"
	// updated.Database.Port == 5432 (from default)

# Performance Considerations

Lenses create new copies of data structures on each Set operation. For deeply nested
structures, this can be expensive. Consider:

1. Using pointer-based structures with MakeLensRef for better performance
2. Batching multiple updates using Modify
3. Using specialized lenses for common patterns (arrays, records, etc.)

# Type Safety

Lenses are fully type-safe. The compiler ensures that:
- Get returns the correct type
- Set accepts the correct type
- Composed lenses maintain type relationships

# Function Reference

Core Lens Creation:
  - MakeLens: Create a lens from getter and setter functions
  - MakeLensCurried: Create a lens with curried setter
  - MakeLensRef: Create a lens for pointer-based structures
  - MakeLensRefCurried: Create a lens for pointers with curried setter
  - MakeLensWithEq: Create a lens with equality optimization for pointer structures
  - MakeLensStrict: Create a lens with strict equality optimization for pointer structures
  - Id: Create an identity lens
  - IdRef: Create an identity lens for pointers

Composition:
  - Compose: Compose two lenses
  - ComposeRef: Compose lenses for pointer structures
  - ComposeOption: Compose lens returning Option with regular lens
  - ComposeOptions: Compose two lenses returning Options

Transformation:
  - Modify: Transform a value through a lens
  - IMap: Transform the focus type of a lens

Optional Value Handling:
  - FromNillable: Create optional lens from nullable pointer
  - FromNillableRef: Create optional lens from nullable pointer (ref version)
  - FromNullableProp: Create lens with default for nullable property
  - FromNullablePropRef: Create lens with default for nullable property (ref version)
  - FromOption: Create lens from Option property with default
  - FromOptionRef: Create lens from Option property with default (ref version)
  - FromPredicate: Create optional lens based on predicate
  - FromPredicateRef: Create optional lens based on predicate (ref version)

# Related Packages

  - github.com/IBM/fp-go/v2/optics/iso: Isomorphisms (bidirectional transformations)
  - github.com/IBM/fp-go/v2/optics/prism: Prisms (focus on sum types)
  - github.com/IBM/fp-go/v2/optics/optional: Optional optics
  - github.com/IBM/fp-go/v2/optics/traversal: Traversals (focus on multiple values)
  - github.com/IBM/fp-go/v2/option: Optional values
  - github.com/IBM/fp-go/v2/endomorphism: Endomorphisms (A → A functions)
*/
package lens
