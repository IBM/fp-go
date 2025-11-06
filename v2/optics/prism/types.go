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

package prism

import (
	"github.com/IBM/fp-go/v2/either"
	O "github.com/IBM/fp-go/v2/option"
)

type (
	// Option is a type alias for O.Option[T], representing an optional value.
	// It is re-exported here for convenience when working with prisms.
	//
	// An Option[T] can be either:
	//   - Some(value): Contains a value of type T
	//   - None: Represents the absence of a value
	//
	// This type is commonly used in prism operations, particularly in the
	// GetOption method which returns Option[A] to indicate whether a value
	// could be extracted from the source type.
	//
	// Type Parameters:
	//   - T: The type of the value that may or may not be present
	//
	// Example:
	//
	//	// A prism's GetOption returns an Option
	//	prism := MakePrism(...)
	//	result := prism.GetOption(value)  // Returns Option[A]
	//
	//	// Check if the value was extracted successfully
	//	if O.IsSome(result) {
	//	    // Value was found
	//	} else {
	//	    // Value was not found (None)
	//	}
	//
	// See also:
	//   - github.com/IBM/fp-go/v2/option for the full Option API
	//   - Prism.GetOption for the primary use case within this package
	Option[T any] = O.Option[T]

	// Either is a type alias for either.Either[E, T], representing a value that can be one of two types.
	// It is re-exported here for convenience when working with prisms that handle error cases.
	//
	// An Either[E, T] can be either:
	//   - Left(error): Contains an error value of type E
	//   - Right(value): Contains a success value of type T
	//
	// This type is commonly used in prism operations for error handling, particularly with
	// the FromEither prism which extracts Right values and returns None for Left values.
	//
	// Type Parameters:
	//   - E: The type of the error/left value
	//   - T: The type of the success/right value
	//
	// Example:
	//
	//	// Using FromEither prism to extract success values
	//	prism := FromEither[error, int]()
	//
	//	// Extract from a Right value
	//	success := either.Right[error](42)
	//	result := prism.GetOption(success)  // Returns Some(42)
	//
	//	// Extract from a Left value
	//	failure := either.Left[int](errors.New("failed"))
	//	result = prism.GetOption(failure)   // Returns None
	//
	//	// ReverseGet wraps a value into Right
	//	wrapped := prism.ReverseGet(100)    // Returns Right(100)
	//
	// Common Use Cases:
	//   - Error handling in functional pipelines
	//   - Representing computations that may fail
	//   - Composing prisms that work with Either types
	//
	// See also:
	//   - github.com/IBM/fp-go/v2/either for the full Either API
	//   - FromEither for creating prisms that work with Either types
	//   - Prism composition for building complex error-handling pipelines
	Either[E, T any] = either.Either[E, T]
)
