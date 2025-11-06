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

// Package json provides functional wrappers around Go's encoding/json package using Either for error handling.
//
// This package wraps JSON marshaling and unmarshaling operations in Either monads, making it easier
// to compose JSON operations with other functional code and handle errors in a functional style.
//
// # Core Concepts
//
// The json package provides type-safe JSON operations that return Either[error, A] instead of
// the traditional (value, error) tuple pattern. This allows for better composition with other
// functional operations and eliminates the need for explicit error checking at each step.
//
// # Basic Usage
//
//	// Unmarshaling JSON data
//	type Person struct {
//	    Name string `json:"name"`
//	    Age  int    `json:"age"`
//	}
//
//	data := []byte(`{"name":"Alice","age":30}`)
//	result := json.Unmarshal[Person](data)
//	// result is Either[error, Person]
//
//	// Marshaling to JSON
//	person := Person{Name: "Bob", Age: 25}
//	jsonBytes := json.Marshal(person)
//	// jsonBytes is Either[error, []byte]
//
// # Chaining Operations
//
// Since Marshal and Unmarshal return Either values, they can be easily composed:
//
//	result := function.Pipe2(
//	    person,
//	    json.Marshal[Person],
//	    either.Chain(json.Unmarshal[Person]),
//	)
//
// # Type Conversion
//
// The package provides utilities for converting between types using JSON as an intermediate format:
//
//	// Convert from one type to another via JSON (returns Either)
//	type Source struct { Value int }
//	type Target struct { Value int }
//
//	src := Source{Value: 42}
//	result := json.ToTypeE[Target](src)
//	// result is Either[error, Target]
//
//	// Convert with Option (discards error details)
//	maybeTarget := json.ToTypeO[Target](src)
//	// maybeTarget is Option[Target]
//
// # Error Handling
//
// All operations return Either[error, A], allowing you to handle errors functionally:
//
//	result := function.Pipe1(
//	    json.Unmarshal[Person](data),
//	    either.Fold(
//	        func(err error) string { return "Failed: " + err.Error() },
//	        func(p Person) string { return "Success: " + p.Name },
//	    ),
//	)
//
// # Type Aliases
//
// The package defines convenient type aliases:
//   - Either[A] = either.Either[error, A]
//   - Option[A] = option.Option[A]
//
// These aliases simplify type signatures and make the code more readable.
package json
