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
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/result"
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

	// Result represents a computation that may fail with an error.
	// It's an alias for Either[error, T].
	Result[T any] = result.Result[T]

	// Endomorphism represents a function from a type to itself (T → T).
	Endomorphism[T any] = endomorphism.Endomorphism[T]

	// Reader represents a computation that depends on an environment R and produces a value T.
	Reader[R, T any] = reader.Reader[R, T]

	// Kleisli represents a function that takes a value of type A and returns a Prism[S, B].
	// This is commonly used for composing prisms in a monadic style.
	//
	// Type Parameters:
	//   - S: The source type of the resulting prism
	//   - A: The input type to the function
	//   - B: The focus type of the resulting prism
	Kleisli[S, A, B any] = func(A) Prism[S, B]

	// Operator represents a function that transforms one prism into another.
	// It takes a Prism[S, A] and returns a Prism[S, B], allowing for prism transformations.
	//
	// Type Parameters:
	//   - S: The source type (remains constant)
	//   - A: The original focus type
	//   - B: The new focus type
	Operator[S, A, B any] = func(Prism[S, A]) Prism[S, B]

	Predicate[A any] = predicate.Predicate[A]

	Lens[S, A any] = lens.Lens[S, A]
)
