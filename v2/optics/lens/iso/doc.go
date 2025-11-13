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
Package iso provides utilities for composing lenses with isomorphisms.

# Overview

This package bridges lenses and isomorphisms, allowing you to transform the focus type
of a lens using an isomorphism. It provides functions to compose lenses with isomorphisms
and to create isomorphisms for common patterns like nullable pointers.

The key insight is that if you have a Lens[S, A] and an Iso[A, B], you can create a
Lens[S, B] by composing them. This allows you to work with transformed views of your
data without changing the underlying structure.

# Core Functions

## FromNillable

Creates an isomorphism between a nullable pointer and an Option type:

	type Config struct {
		Timeout *int
	}

	// Create isomorphism: *int ↔ Option[int]
	timeoutIso := iso.FromNillable[int]()

	// nil → None, &value → Some(value)
	opt := timeoutIso.Get(nil)  // None[int]
	num := 42
	opt = timeoutIso.Get(&num)  // Some(42)

	// None → nil, Some(value) → &value
	ptr := timeoutIso.ReverseGet(O.None[int]())  // nil
	ptr = timeoutIso.ReverseGet(O.Some(42))      // &42

## Compose

Composes a lens with an isomorphism to transform the focus type:

	type Person struct {
		Name string
		Age  int
	}

	type Celsius float64
	type Fahrenheit float64

	type Weather struct {
		Temperature Celsius
	}

	// Lens to access temperature
	tempLens := L.MakeLens(
		func(w Weather) Celsius { return w.Temperature },
		func(w Weather, t Celsius) Weather {
			w.Temperature = t
			return w
		},
	)

	// Isomorphism: Celsius ↔ Fahrenheit
	celsiusToFahrenheit := I.MakeIso(
		func(c Celsius) Fahrenheit { return Fahrenheit(c*9/5 + 32) },
		func(f Fahrenheit) Celsius { return Celsius((f - 32) * 5 / 9) },
	)

	// Compose to work with Fahrenheit
	tempFahrenheitLens := F.Pipe1(
		tempLens,
		iso.Compose[Weather, Celsius, Fahrenheit](celsiusToFahrenheit),
	)

	weather := Weather{Temperature: 20} // 20°C
	tempF := tempFahrenheitLens.Get(weather)        // 68°F
	updated := tempFahrenheitLens.Set(86)(weather)  // Set to 86°F (30°C)

# Use Cases

## Working with Nullable Fields

Convert between nullable pointers and Option types:

	type DatabaseConfig struct {
		Host     string
		Port     int
		Username string
		Password *string // Nullable
	}

	type AppConfig struct {
		Database *DatabaseConfig
	}

	// Lens to database config
	dbLens := L.MakeLens(
		func(c AppConfig) *DatabaseConfig { return c.Database },
		func(c AppConfig, db *DatabaseConfig) AppConfig {
			c.Database = db
			return c
		},
	)

	// Isomorphism for nullable pointer
	dbIso := iso.FromNillable[DatabaseConfig]()

	// Compose to work with Option
	dbOptLens := F.Pipe1(
		dbLens,
		iso.Compose[AppConfig, *DatabaseConfig, O.Option[DatabaseConfig]](dbIso),
	)

	config := AppConfig{Database: nil}
	dbOpt := dbOptLens.Get(config) // None[DatabaseConfig]

	// Set with Some
	newDB := DatabaseConfig{Host: "localhost", Port: 5432}
	updated := dbOptLens.Set(O.Some(newDB))(config)

## Unit Conversions

Work with different units of measurement:

	type Distance struct {
		Meters float64
	}

	type Kilometers float64
	type Miles float64

	// Lens to meters
	metersLens := L.MakeLens(
		func(d Distance) float64 { return d.Meters },
		func(d Distance, m float64) Distance {
			d.Meters = m
			return d
		},
	)

	// Isomorphism: meters ↔ kilometers
	metersToKm := I.MakeIso(
		func(m float64) Kilometers { return Kilometers(m / 1000) },
		func(km Kilometers) float64 { return float64(km * 1000) },
	)

	// Compose to work with kilometers
	kmLens := F.Pipe1(
		metersLens,
		iso.Compose[Distance, float64, Kilometers](metersToKm),
	)

	distance := Distance{Meters: 5000}
	km := kmLens.Get(distance)           // 5 km
	updated := kmLens.Set(Kilometers(10))(distance) // 10000 meters

## Type Wrappers

Work with newtype wrappers:

	type UserId int
	type User struct {
		ID   UserId
		Name string
	}

	// Lens to user ID
	idLens := L.MakeLens(
		func(u User) UserId { return u.ID },
		func(u User, id UserId) User {
			u.ID = id
			return u
		},
	)

	// Isomorphism: UserId ↔ int
	userIdIso := I.MakeIso(
		func(id UserId) int { return int(id) },
		func(i int) UserId { return UserId(i) },
	)

	// Compose to work with raw int
	idIntLens := F.Pipe1(
		idLens,
		iso.Compose[User, UserId, int](userIdIso),
	)

	user := User{ID: 42, Name: "Alice"}
	rawId := idIntLens.Get(user)      // 42 (int)
	updated := idIntLens.Set(100)(user) // UserId(100)

## Nested Nullable Fields

Safely navigate through nullable nested structures:

	type Address struct {
		Street string
		City   string
	}

	type Person struct {
		Name    string
		Address *Address
	}

	type Company struct {
		Name string
		CEO  *Person
	}

	// Lens to CEO
	ceoLens := L.MakeLens(
		func(c Company) *Person { return c.CEO },
		func(c Company, p *Person) Company {
			c.CEO = p
			return c
		},
	)

	// Isomorphism for nullable person
	personIso := iso.FromNillable[Person]()

	// Compose to work with Option[Person]
	ceoOptLens := F.Pipe1(
		ceoLens,
		iso.Compose[Company, *Person, O.Option[Person]](personIso),
	)

	company := Company{Name: "Acme Corp", CEO: nil}
	ceo := ceoOptLens.Get(company) // None[Person]

	// Set CEO
	newCEO := Person{Name: "Alice", Address: nil}
	updated := ceoOptLens.Set(O.Some(newCEO))(company)

# Composition Patterns

## Chaining Multiple Isomorphisms

	type Meters float64
	type Kilometers float64
	type Miles float64

	type Journey struct {
		Distance Meters
	}

	// Lens to distance
	distLens := L.MakeLens(
		func(j Journey) Meters { return j.Distance },
		func(j Journey, d Meters) Journey {
			j.Distance = d
			return j
		},
	)

	// Isomorphisms
	metersToKm := I.MakeIso(
		func(m Meters) Kilometers { return Kilometers(m / 1000) },
		func(km Kilometers) Meters { return Meters(km * 1000) },
	)

	kmToMiles := I.MakeIso(
		func(km Kilometers) Miles { return Miles(km * 0.621371) },
		func(mi Miles) Kilometers { return Kilometers(mi / 0.621371) },
	)

	// Compose lens with chained isomorphisms
	milesLens := F.Pipe2(
		distLens,
		iso.Compose[Journey, Meters, Kilometers](metersToKm),
		iso.Compose[Journey, Kilometers, Miles](kmToMiles),
	)

	journey := Journey{Distance: 5000} // 5000 meters
	miles := milesLens.Get(journey)    // ~3.11 miles

## Combining with Optional Lenses

	type Config struct {
		Database *DatabaseConfig
	}

	type DatabaseConfig struct {
		Port int
	}

	// Lens to database (nullable)
	dbLens := L.MakeLens(
		func(c Config) *DatabaseConfig { return c.Database },
		func(c Config, db *DatabaseConfig) Config {
			c.Database = db
			return c
		},
	)

	// Convert to Option lens
	dbIso := iso.FromNillable[DatabaseConfig]()
	dbOptLens := F.Pipe1(
		dbLens,
		iso.Compose[Config, *DatabaseConfig, O.Option[DatabaseConfig]](dbIso),
	)

	// Now compose with lens to port
	portLens := L.MakeLens(
		func(db DatabaseConfig) int { return db.Port },
		func(db DatabaseConfig, port int) DatabaseConfig {
			db.Port = port
			return db
		},
	)

	// Use ComposeOption to handle the Option
	defaultDB := DatabaseConfig{Port: 5432}
	configPortLens := F.Pipe1(
		dbOptLens,
		L.ComposeOption[Config, int](defaultDB)(portLens),
	)

# Performance Considerations

Composing lenses with isomorphisms is efficient:
  - No additional allocations beyond the lens and iso structures
  - Composition creates function closures but is still performant
  - The isomorphism transformations are applied on-demand
  - Consider caching composed lenses for frequently used paths

# Type Safety

All operations are fully type-safe:
  - Compile-time type checking ensures correct composition
  - Generic type parameters prevent type mismatches
  - No runtime type assertions needed
  - The compiler enforces that isomorphisms are properly reversible

# Related Packages

  - github.com/IBM/fp-go/v2/optics/lens: Core lens functionality
  - github.com/IBM/fp-go/v2/optics/iso: Core isomorphism functionality
  - github.com/IBM/fp-go/v2/optics/iso/lens: Convert isomorphisms to lenses
  - github.com/IBM/fp-go/v2/option: Option type and operations
  - github.com/IBM/fp-go/v2/function: Function composition utilities

# See Also

For more information on lenses and isomorphisms:
  - optics/lens package documentation
  - optics/iso package documentation
  - optics package overview
*/
package iso
