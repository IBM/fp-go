// Copyright (c) 2025 IBM Corp.
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

package readereither

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/reader"
)

// Promap is the profunctor map operation that transforms both the input and output of a ReaderEither.
// It applies f to the input environment (contravariantly) and g to the output value (covariantly).
//
// See: https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#profunctor
//
// This operation allows you to:
//   - Adapt the environment type before passing it to the ReaderEither (via f)
//   - Transform the success value after the computation completes (via g)
//
// The error type E remains unchanged through the transformation.
//
// Type Parameters:
//   - R: The original environment type expected by the ReaderEither
//   - E: The error type (unchanged)
//   - A: The original success type produced by the ReaderEither
//   - D: The new input environment type
//   - B: The new output success type
//
// Parameters:
//   - f: Function to transform the input environment from D to R (contravariant)
//   - g: Function to transform the output success value from A to B (covariant)
//
// Returns:
//   - A Kleisli arrow that takes a ReaderEither[R, E, A] and returns a ReaderEither[D, E, B]
//
//go:inline
func Promap[R, E, A, D, B any](f func(D) R, g func(A) B) Kleisli[D, E, ReaderEither[R, E, A], B] {
	return reader.Promap(f, either.Map[E](g))
}

// Contramap changes the value of the local environment during the execution of a ReaderEither.
// This is the contravariant functor operation that transforms the input environment.
//
// See: https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#profunctor
//
// Contramap is useful for adapting a ReaderEither to work with a different environment type
// by providing a function that converts the new environment to the expected one.
//
// Type Parameters:
//   - E: The error type (unchanged)
//   - A: The success type (unchanged)
//   - R2: The new input environment type
//   - R1: The original environment type expected by the ReaderEither
//
// Parameters:
//   - f: Function to transform the environment from R2 to R1
//
// Returns:
//   - A Kleisli arrow that takes a ReaderEither[R1, E, A] and returns a ReaderEither[R2, E, A]
//
//go:inline
func Contramap[E, A, R1, R2 any](f func(R2) R1) Kleisli[R2, E, ReaderEither[R1, E, A], A] {
	return reader.Contramap[Either[E, A]](f)
}
