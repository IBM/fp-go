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

// Package ioresult provides functional programming combinators for working with IO operations
// that can fail with errors, following Go's idiomatic (value, error) tuple pattern.
//
// # Overview
//
// IOResult[A] represents a computation that performs IO and returns either a value of type A
// or an error. It is defined as:
//
//	type IOResult[A any] = func() (A, error)
//
// This is the idiomatic Go version of IOEither, using Go's standard error handling pattern
// instead of the Either monad. It combines:
//   - IO effects (functions that perform side effects)
//   - Error handling via Go's (value, error) tuple return pattern
//
// # Why Parameterless Functions Represent IO Operations
//
// The key insight behind IOResult is that a function returning a value without taking any
// input can only produce that value through side effects. Consider:
//
//	func() int { return 42 }              // Pure: always returns 42
//	func() int { return readFromFile() }  // Impure: result depends on external state
//
// When a parameterless function returns different values on different calls, or when it
// interacts with the outside world (filesystem, network, random number generator, current
// time, database), it is performing a side effect - an observable interaction with state
// outside the function's scope.
//
// # Lazy Evaluation and Referential Transparency
//
// IOResult provides two critical benefits:
//
//  1. **Lazy Evaluation**: The side effect doesn't happen when you create the IOResult,
//     only when you call it (execute it). This allows you to build complex computations
//     as pure data structures and defer execution until needed.
//
//     // This doesn't read the file yet, just describes how to read it
//     readConfig := func() (Config, error) { return os.ReadFile("config.json") }
//
//     // Still hasn't read the file, just composed operations
//     parsed := Map(parseJSON)(readConfig)
//
//     // NOW it reads the file and parses it
//     config, err := parsed()
//
//  2. **Referential Transparency of the Description**: While the IO operation itself has
//     side effects, the IOResult value (the function) is referentially transparent. You can
//     pass it around, compose it, and reason about it without triggering the side effect.
//     The side effect only occurs when you explicitly call the function.
//
// # Distinguishing Pure from Impure Operations
//
// The type system makes the distinction clear:
//
//	// Pure function: always returns the same output for the same input
//	func double(x int) int { return x * 2 }
//
//	// Impure operation: encapsulated in IOResult
//	func readFile(path string) IOResult[[]byte] {
//	    return func() ([]byte, error) {
//	        return os.ReadFile(path)  // Side effect: file system access
//	    }
//	}
//
// The IOResult type explicitly marks operations as having side effects, making the
// distinction between pure and impure code visible in the type system. This allows
// developers to:
//   - Identify which parts of the code interact with external state
//   - Test pure logic separately from IO operations
//   - Compose IO operations while keeping them lazy
//   - Control when and where side effects occur
//
// # Examples of Side Effects Captured by IOResult
//
// IOResult is appropriate for operations that:
//   - Read from or write to files, databases, or network
//   - Generate random numbers
//   - Read the current time
//   - Modify mutable state
//   - Interact with external APIs
//   - Execute system commands
//   - Acquire or release resources
//
// Example:
//
//	// Each call potentially returns a different value
//	getCurrentTime := func() (time.Time, error) {
//	    return time.Now(), nil  // Side effect: reads system clock
//	}
//
//	// Each call reads from external state
//	readDatabase := func() (User, error) {
//	    return db.Query("SELECT * FROM users WHERE id = ?", 1)
//	}
//
//	// Composes multiple IO operations
//	pipeline := F.Pipe2(
//	    getCurrentTime,
//	    Chain(func(t time.Time) IOResult[string] {
//	        return func() (string, error) {
//	            return fmt.Sprintf("Time: %v", t), nil
//	        }
//	    }),
//	)
//
// # Core Concepts
//
// IOResult follows functional programming principles:
//   - Functor: Transform successful values with Map
//   - Applicative: Combine multiple IOResults with Ap, ApS
//   - Monad: Chain dependent computations with Chain, Bind
//   - Error recovery: Handle errors with ChainLeft, Alt
//
// # Basic Usage
//
// Creating IOResult values:
//
//	success := Of(42)                          // Right value
//	failure := Left[int](errors.New("error"))  // Left (error) value
//
// Transforming values:
//
//	doubled := Map(N.Mul(2))(success)
//
// Chaining computations:
//
//	result := Chain(func(x int) IOResult[string] {
//	    return Of(fmt.Sprintf("%d", x))
//	})(success)
//
// # Do Notation
//
// The package supports do-notation style composition for building complex computations:
//
//	result := F.Pipe5(
//	    Of("John"),
//	    BindTo(T.Of[string]),
//	    ApS(T.Push1[string, int], Of(42)),
//	    Bind(T.Push2[string, int, string], func(t T.Tuple2[string, int]) IOResult[string] {
//	        return Of(fmt.Sprintf("%s: %d", t.F1, t.F2))
//	    }),
//	    Map(transform),
//	)
//
// # Error Handling
//
// IOResult provides several ways to handle errors:
//   - ChainLeft: Transform error values into new computations
//   - Alt: Provide alternative computations when an error occurs
//   - GetOrElse: Extract values with a default for errors
//   - Fold: Handle both success and error cases explicitly
//
// # Concurrency
//
// IOResult supports both sequential and parallel execution:
//   - ApSeq, TraverseArraySeq: Sequential execution
//   - ApPar, TraverseArrayPar: Parallel execution (default)
//   - Ap, TraverseArray: Defaults to parallel execution
//
// # Resource Management
//
// The package provides resource management utilities:
//   - Bracket: Acquire, use, and release resources safely
//   - WithResource: Scoped resource management
//   - WithLock: Execute operations within a lock scope
//
// # Conversion Functions
//
// IOResult interoperates with other types:
//   - FromEither: Convert Either to IOResult
//   - FromResult: Convert (value, error) tuple to IOResult
//   - FromOption: Convert Option to IOResult
//   - FromIO: Convert pure IO to IOResult (always succeeds)
//
// # Examples
//
// See the example tests for detailed usage patterns:
//   - examples_create_test.go: Creating IOResult values
//   - examples_do_test.go: Using do-notation
//   - examples_extract_test.go: Extracting values from IOResult
package ioresult

//go:generate go run .. ioeither --count 10 --filename gen.go
