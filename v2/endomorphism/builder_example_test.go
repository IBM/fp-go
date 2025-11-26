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

package endomorphism_test

import (
	"fmt"
	"time"

	A "github.com/IBM/fp-go/v2/array"
	"github.com/IBM/fp-go/v2/endomorphism"
	M "github.com/IBM/fp-go/v2/monoid"
	N "github.com/IBM/fp-go/v2/number"
)

// Example_build_basicUsage demonstrates basic usage of the Build function
// to construct a value from the zero value using endomorphisms.
func Example_build_basicUsage() {
	// Define simple endomorphisms
	addTen := N.Add(10)
	double := N.Mul(2)

	// Compose them using monoid (RIGHT-TO-LEFT execution)
	// double is applied first, then addTen
	builder := M.ConcatAll(endomorphism.Monoid[int]())(A.From(
		addTen,
		double,
	))

	// Build from zero value: 0 * 2 = 0, 0 + 10 = 10
	result := endomorphism.Build(builder)
	fmt.Println(result)
	// Output: 10
}

// Example_build_configBuilder demonstrates using Build as a configuration builder pattern.
func Example_build_configBuilder() {
	type Config struct {
		Host    string
		Port    int
		Timeout time.Duration
		Debug   bool
	}

	// Define builder functions as endomorphisms
	withHost := func(host string) endomorphism.Endomorphism[Config] {
		return func(c Config) Config {
			c.Host = host
			return c
		}
	}

	withPort := func(port int) endomorphism.Endomorphism[Config] {
		return func(c Config) Config {
			c.Port = port
			return c
		}
	}

	withTimeout := func(d time.Duration) endomorphism.Endomorphism[Config] {
		return func(c Config) Config {
			c.Timeout = d
			return c
		}
	}

	withDebug := func(debug bool) endomorphism.Endomorphism[Config] {
		return func(c Config) Config {
			c.Debug = debug
			return c
		}
	}

	// Compose builders using monoid
	configBuilder := M.ConcatAll(endomorphism.Monoid[Config]())([]endomorphism.Endomorphism[Config]{
		withHost("localhost"),
		withPort(8080),
		withTimeout(30 * time.Second),
		withDebug(true),
	})

	// Build the configuration from zero value
	config := endomorphism.Build(configBuilder)

	fmt.Printf("Host: %s\n", config.Host)
	fmt.Printf("Port: %d\n", config.Port)
	fmt.Printf("Timeout: %v\n", config.Timeout)
	fmt.Printf("Debug: %v\n", config.Debug)
	// Output:
	// Host: localhost
	// Port: 8080
	// Timeout: 30s
	// Debug: true
}

// Example_build_stringBuilder demonstrates building a string using endomorphisms.
func Example_build_stringBuilder() {
	// Define string transformation endomorphisms
	appendHello := func(s string) string { return s + "Hello" }
	appendSpace := func(s string) string { return s + " " }
	appendWorld := func(s string) string { return s + "World" }
	appendExclamation := func(s string) string { return s + "!" }

	// Compose transformations (RIGHT-TO-LEFT execution)
	stringBuilder := M.ConcatAll(endomorphism.Monoid[string]())([]endomorphism.Endomorphism[string]{
		appendHello,
		appendSpace,
		appendWorld,
		appendExclamation,
	})

	// Build the string from empty string
	result := endomorphism.Build(stringBuilder)
	fmt.Println(result)
	// Output: !World Hello
}

// Example_build_personBuilder demonstrates building a complex struct using the builder pattern.
func Example_build_personBuilder() {
	type Person struct {
		FirstName string
		LastName  string
		Age       int
		Email     string
	}

	// Define builder functions
	withFirstName := func(name string) endomorphism.Endomorphism[Person] {
		return func(p Person) Person {
			p.FirstName = name
			return p
		}
	}

	withLastName := func(name string) endomorphism.Endomorphism[Person] {
		return func(p Person) Person {
			p.LastName = name
			return p
		}
	}

	withAge := func(age int) endomorphism.Endomorphism[Person] {
		return func(p Person) Person {
			p.Age = age
			return p
		}
	}

	withEmail := func(email string) endomorphism.Endomorphism[Person] {
		return func(p Person) Person {
			p.Email = email
			return p
		}
	}

	// Build a person
	personBuilder := M.ConcatAll(endomorphism.Monoid[Person]())([]endomorphism.Endomorphism[Person]{
		withFirstName("Alice"),
		withLastName("Smith"),
		withAge(30),
		withEmail("alice.smith@example.com"),
	})

	person := endomorphism.Build(personBuilder)

	fmt.Printf("%s %s, Age: %d, Email: %s\n",
		person.FirstName, person.LastName, person.Age, person.Email)
	// Output: Alice Smith, Age: 30, Email: alice.smith@example.com
}

// Example_build_conditionalBuilder demonstrates conditional building using endomorphisms.
func Example_build_conditionalBuilder() {
	type Settings struct {
		Theme      string
		FontSize   int
		AutoSave   bool
		Animations bool
	}

	withTheme := func(theme string) endomorphism.Endomorphism[Settings] {
		return func(s Settings) Settings {
			s.Theme = theme
			return s
		}
	}

	withFontSize := func(size int) endomorphism.Endomorphism[Settings] {
		return func(s Settings) Settings {
			s.FontSize = size
			return s
		}
	}

	withAutoSave := func(enabled bool) endomorphism.Endomorphism[Settings] {
		return func(s Settings) Settings {
			s.AutoSave = enabled
			return s
		}
	}

	withAnimations := func(enabled bool) endomorphism.Endomorphism[Settings] {
		return func(s Settings) Settings {
			s.Animations = enabled
			return s
		}
	}

	// Build settings conditionally
	isDarkMode := true
	isAccessibilityMode := true

	// Note: Monoid executes RIGHT-TO-LEFT, so later items in the slice are applied first
	// We need to add items in reverse order for the desired effect
	builders := []endomorphism.Endomorphism[Settings]{}

	if isAccessibilityMode {
		builders = append(builders, withFontSize(18)) // Will be applied last (overrides)
		builders = append(builders, withAnimations(false))
	}

	if isDarkMode {
		builders = append(builders, withTheme("dark"))
	} else {
		builders = append(builders, withTheme("light"))
	}

	builders = append(builders, withAutoSave(true))
	builders = append(builders, withFontSize(14)) // Will be applied first

	settingsBuilder := M.ConcatAll(endomorphism.Monoid[Settings]())(builders)
	settings := endomorphism.Build(settingsBuilder)

	fmt.Printf("Theme: %s\n", settings.Theme)
	fmt.Printf("FontSize: %d\n", settings.FontSize)
	fmt.Printf("AutoSave: %v\n", settings.AutoSave)
	fmt.Printf("Animations: %v\n", settings.Animations)
	// Output:
	// Theme: dark
	// FontSize: 18
	// AutoSave: true
	// Animations: false
}
