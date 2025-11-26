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

package endomorphism

import (
	"github.com/IBM/fp-go/v2/function"
	A "github.com/IBM/fp-go/v2/internal/array"
)

// Build applies an endomorphism to the zero value of type A, effectively using
// the endomorphism as a builder pattern.
//
// # Endomorphism as Builder Pattern
//
// An endomorphism (a function from type A to type A) can be viewed as a builder pattern
// because it transforms a value of a type into another value of the same type. When you
// compose multiple endomorphisms together, you create a pipeline of transformations that
// build up a final value step by step.
//
// The Build function starts with the zero value of type A and applies the endomorphism
// to it, making it particularly useful for building complex values from scratch using
// a functional composition of transformations.
//
// # Builder Pattern Characteristics
//
// Traditional builder patterns have these characteristics:
//  1. Start with an initial (often empty) state
//  2. Apply a series of transformations/configurations
//  3. Return the final built object
//
// Endomorphisms provide the same pattern functionally:
//  1. Start with zero value: var a A
//  2. Apply composed endomorphisms: e(a)
//  3. Return the transformed value
//
// # Type Parameters
//
//   - A: The type being built/transformed
//
// # Parameters
//
//   - e: An endomorphism (or composition of endomorphisms) that transforms type A
//
// # Returns
//
//	The result of applying the endomorphism to the zero value of type A
//
// # Example - Building a Configuration Object
//
//	type Config struct {
//	    Host     string
//	    Port     int
//	    Timeout  time.Duration
//	    Debug    bool
//	}
//
//	// Define builder functions as endomorphisms
//	withHost := func(host string) Endomorphism[Config] {
//	    return func(c Config) Config {
//	        c.Host = host
//	        return c
//	    }
//	}
//
//	withPort := func(port int) Endomorphism[Config] {
//	    return func(c Config) Config {
//	        c.Port = port
//	        return c
//	    }
//	}
//
//	withTimeout := func(d time.Duration) Endomorphism[Config] {
//	    return func(c Config) Config {
//	        c.Timeout = d
//	        return c
//	    }
//	}
//
//	withDebug := func(debug bool) Endomorphism[Config] {
//	    return func(c Config) Config {
//	        c.Debug = debug
//	        return c
//	    }
//	}
//
//	// Compose builders using monoid operations
//	import M "github.com/IBM/fp-go/v2/monoid"
//
//	configBuilder := M.ConcatAll(Monoid[Config]())(
//	    withHost("localhost"),
//	    withPort(8080),
//	    withTimeout(30 * time.Second),
//	    withDebug(true),
//	)
//
//	// Build the final configuration
//	config := Build(configBuilder)
//	// Result: Config{Host: "localhost", Port: 8080, Timeout: 30s, Debug: true}
//
// # Example - Building a String with Transformations
//
//	import (
//	    "strings"
//	    M "github.com/IBM/fp-go/v2/monoid"
//	)
//
//	// Define string transformation endomorphisms
//	appendHello := func(s string) string { return s + "Hello" }
//	appendSpace := func(s string) string { return s + " " }
//	appendWorld := func(s string) string { return s + "World" }
//	toUpper := strings.ToUpper
//
//	// Compose transformations
//	stringBuilder := M.ConcatAll(Monoid[string]())(
//	    appendHello,
//	    appendSpace,
//	    appendWorld,
//	    toUpper,
//	)
//
//	// Build the final string from empty string
//	result := Build(stringBuilder)
//	// Result: "HELLO WORLD"
//
// # Example - Building a Slice with Operations
//
//	type IntSlice []int
//
//	appendValue := func(v int) Endomorphism[IntSlice] {
//	    return func(s IntSlice) IntSlice {
//	        return append(s, v)
//	    }
//	}
//
//	sortSlice := func(s IntSlice) IntSlice {
//	    sorted := make(IntSlice, len(s))
//	    copy(sorted, s)
//	    sort.Ints(sorted)
//	    return sorted
//	}
//
//	// Build a sorted slice
//	sliceBuilder := M.ConcatAll(Monoid[IntSlice]())(
//	    appendValue(5),
//	    appendValue(2),
//	    appendValue(8),
//	    appendValue(1),
//	    sortSlice,
//	)
//
//	result := Build(sliceBuilder)
//	// Result: IntSlice{1, 2, 5, 8}
//
// # Advantages of Endomorphism Builder Pattern
//
//  1. **Composability**: Builders can be composed using monoid operations
//  2. **Immutability**: Each transformation returns a new value (if implemented immutably)
//  3. **Type Safety**: The type system ensures all transformations work on the same type
//  4. **Reusability**: Individual builder functions can be reused and combined differently
//  5. **Testability**: Each transformation can be tested independently
//  6. **Declarative**: The composition clearly expresses the building process
//
// # Comparison with Traditional Builder Pattern
//
// Traditional OOP Builder:
//
//	config := NewConfigBuilder().
//	    WithHost("localhost").
//	    WithPort(8080).
//	    WithTimeout(30 * time.Second).
//	    Build()
//
// Endomorphism Builder:
//
//	config := Build(M.ConcatAll(Monoid[Config]())(
//	    withHost("localhost"),
//	    withPort(8080),
//	    withTimeout(30 * time.Second),
//	))
//
// Both achieve the same goal, but the endomorphism approach:
//   - Uses pure functions instead of methods
//   - Leverages algebraic properties (monoid) for composition
//   - Allows for more flexible composition patterns
//   - Integrates naturally with other functional programming constructs
func Build[A any](e Endomorphism[A]) A {
	var a A
	return e(a)
}

// ConcatAll combines multiple endomorphisms into a single endomorphism using composition.
//
// This function takes a slice of endomorphisms and combines them using the monoid's
// concat operation (which is composition). The resulting endomorphism, when applied,
// will execute all the input endomorphisms in RIGHT-TO-LEFT order (mathematical composition order).
//
// IMPORTANT: Execution order is RIGHT-TO-LEFT:
//   - ConcatAll([]Endomorphism{f, g, h}) creates an endomorphism that applies h, then g, then f
//   - This is equivalent to f ∘ g ∘ h in mathematical notation
//   - The last endomorphism in the slice is applied first
//
// If the slice is empty, returns the identity endomorphism.
//
// # Type Parameters
//
//   - T: The type that the endomorphisms operate on
//
// # Parameters
//
//   - es: A slice of endomorphisms to combine
//
// # Returns
//
//	A single endomorphism that represents the composition of all input endomorphisms
//
// # Example - Basic Composition
//
//	double := N.Mul(2)
//	increment := N.Add(1)
//	square := func(x int) int { return x * x }
//
//	// Combine endomorphisms (RIGHT-TO-LEFT execution)
//	combined := ConcatAll([]Endomorphism[int]{double, increment, square})
//	result := combined(5)
//	// Execution: square(5) = 25, increment(25) = 26, double(26) = 52
//	// Result: 52
//
// # Example - Building with ConcatAll
//
//	type Config struct {
//	    Host string
//	    Port int
//	}
//
//	withHost := func(host string) Endomorphism[Config] {
//	    return func(c Config) Config {
//	        c.Host = host
//	        return c
//	    }
//	}
//
//	withPort := func(port int) Endomorphism[Config] {
//	    return func(c Config) Config {
//	        c.Port = port
//	        return c
//	    }
//	}
//
//	// Combine configuration builders
//	configBuilder := ConcatAll([]Endomorphism[Config]{
//	    withHost("localhost"),
//	    withPort(8080),
//	})
//
//	// Apply to zero value
//	config := Build(configBuilder)
//	// Result: Config{Host: "localhost", Port: 8080}
//
// # Example - Empty Slice
//
//	// Empty slice returns identity
//	identity := ConcatAll([]Endomorphism[int]{})
//	result := identity(42) // Returns: 42
//
// # Relationship to Monoid
//
// ConcatAll is equivalent to using M.ConcatAll with the endomorphism Monoid:
//
//	import M "github.com/IBM/fp-go/v2/monoid"
//
//	// These are equivalent:
//	result1 := ConcatAll(endomorphisms)
//	result2 := M.ConcatAll(Monoid[T]())(endomorphisms)
//
// # Use Cases
//
//  1. **Pipeline Construction**: Build transformation pipelines from individual steps
//  2. **Configuration Building**: Combine multiple configuration setters
//  3. **Data Transformation**: Chain multiple data transformations
//  4. **Middleware Composition**: Combine middleware functions
//  5. **Validation Chains**: Compose multiple validation functions
func ConcatAll[T any](es []Endomorphism[T]) Endomorphism[T] {
	return A.Reduce(es, MonadCompose[T], function.Identity[T])
}

// Reduce applies a slice of endomorphisms to the zero value of type T in LEFT-TO-RIGHT order.
//
// This function is a convenience wrapper that:
//  1. Starts with the zero value of type T
//  2. Applies each endomorphism in the slice from left to right
//  3. Returns the final transformed value
//
// IMPORTANT: Execution order is LEFT-TO-RIGHT:
//   - Reduce([]Endomorphism{f, g, h}) applies f first, then g, then h
//   - This is the opposite of ConcatAll's RIGHT-TO-LEFT order
//   - Each endomorphism receives the result of the previous one
//
// This is equivalent to: Build(ConcatAll(reverse(es))) but more efficient and clearer
// for left-to-right sequential application.
//
// # Type Parameters
//
//   - T: The type being transformed
//
// # Parameters
//
//   - es: A slice of endomorphisms to apply sequentially
//
// # Returns
//
//	The final value after applying all endomorphisms to the zero value
//
// # Example - Sequential Transformations
//
//	double := N.Mul(2)
//	increment := N.Add(1)
//	square := func(x int) int { return x * x }
//
//	// Apply transformations LEFT-TO-RIGHT
//	result := Reduce([]Endomorphism[int]{double, increment, square})
//	// Execution: 0 -> double(0) = 0 -> increment(0) = 1 -> square(1) = 1
//	// Result: 1
//
//	// With a non-zero starting point, use a custom initial value:
//	addTen := N.Add(10)
//	result2 := Reduce([]Endomorphism[int]{addTen, double, increment})
//	// Execution: 0 -> addTen(0) = 10 -> double(10) = 20 -> increment(20) = 21
//	// Result: 21
//
// # Example - Building a String
//
//	appendHello := func(s string) string { return s + "Hello" }
//	appendSpace := func(s string) string { return s + " " }
//	appendWorld := func(s string) string { return s + "World" }
//
//	// Build string LEFT-TO-RIGHT
//	result := Reduce([]Endomorphism[string]{
//	    appendHello,
//	    appendSpace,
//	    appendWorld,
//	})
//	// Execution: "" -> "Hello" -> "Hello " -> "Hello World"
//	// Result: "Hello World"
//
// # Example - Configuration Building
//
//	type Settings struct {
//	    Theme    string
//	    FontSize int
//	}
//
//	withTheme := func(theme string) Endomorphism[Settings] {
//	    return func(s Settings) Settings {
//	        s.Theme = theme
//	        return s
//	    }
//	}
//
//	withFontSize := func(size int) Endomorphism[Settings] {
//	    return func(s Settings) Settings {
//	        s.FontSize = size
//	        return s
//	    }
//	}
//
//	// Build settings LEFT-TO-RIGHT
//	settings := Reduce([]Endomorphism[Settings]{
//	    withTheme("dark"),
//	    withFontSize(14),
//	})
//	// Result: Settings{Theme: "dark", FontSize: 14}
//
// # Comparison with ConcatAll
//
//	// ConcatAll: RIGHT-TO-LEFT composition, returns endomorphism
//	endo := ConcatAll([]Endomorphism[int]{f, g, h})
//	result1 := endo(value) // Applies h, then g, then f
//
//	// Reduce: LEFT-TO-RIGHT application, returns final value
//	result2 := Reduce([]Endomorphism[int]{f, g, h})
//	// Applies f to zero, then g, then h
//
// # Use Cases
//
//  1. **Sequential Processing**: Apply transformations in order
//  2. **Pipeline Execution**: Execute a pipeline from start to finish
//  3. **Builder Pattern**: Build objects step by step
//  4. **State Machines**: Apply state transitions in sequence
//  5. **Data Flow**: Transform data through multiple stages
func Reduce[T any](es []Endomorphism[T]) T {
	var t T
	return A.Reduce(es, func(t T, e Endomorphism[T]) T { return e(t) }, t)
}
