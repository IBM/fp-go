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

// Package consumer provides types and utilities for functions that consume values without returning results.
//
// A Consumer represents a side-effecting operation that accepts a value but produces no output.
// This is useful for operations like logging, printing, updating state, or any action where
// the return value is not needed.
package consumer

type (
	// Consumer represents a function that accepts a value of type A and performs a side effect.
	// It does not return any value, making it useful for operations where only the side effect matters,
	// such as logging, printing, or updating external state.
	//
	// This is a fundamental concept in functional programming for handling side effects in a
	// controlled manner. Consumers can be composed, chained, or used in higher-order functions
	// to build complex side-effecting behaviors.
	//
	// Type Parameters:
	//   - A: The type of value consumed by the function
	//
	// Example:
	//
	//	// A simple consumer that prints values
	//	var printInt Consumer[int] = func(x int) {
	//	    fmt.Println(x)
	//	}
	//	printInt(42) // Prints: 42
	//
	//	// A consumer that logs messages
	//	var logger Consumer[string] = func(msg string) {
	//	    log.Println(msg)
	//	}
	//	logger("Hello, World!") // Logs: Hello, World!
	//
	//	// Consumers can be used in functional pipelines
	//	var saveToDatabase Consumer[User] = func(user User) {
	//	    db.Save(user)
	//	}
	Consumer[A any] = func(A)

	// Operator represents a function that transforms a Consumer[A] into a Consumer[B].
	// This is useful for composing and adapting consumers to work with different types.
	Operator[A, B any] = func(Consumer[A]) Consumer[B]
)
