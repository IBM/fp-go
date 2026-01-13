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

package readerresult

import "github.com/IBM/fp-go/v2/idiomatic/result"

// Promap is the profunctor map operation that transforms both the input and output of a ReaderResult.
// It applies f to the input environment (contravariantly) and g to the output value (covariantly).
//
// See: https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#profunctor
//
// This operation allows you to:
//   - Adapt the environment type before passing it to the ReaderResult (via f)
//   - Transform the success value after the computation completes (via g)
//
// The error type is fixed as error and remains unchanged through the transformation.
//
// Type Parameters:
//   - E: The original environment type expected by the ReaderResult
//   - A: The original success type produced by the ReaderResult
//   - D: The new input environment type
//   - B: The new output success type
//
// Parameters:
//   - f: Function to transform the input environment from D to E (contravariant)
//   - g: Function to transform the output success value from A to B (covariant)
//
// Returns:
//   - A Kleisli arrow that takes a ReaderResult[E, A] and returns a ReaderResult[D, B]
//
//go:inline
func Promap[E, A, D, B any](f func(D) E, g func(A) B) Kleisli[D, ReaderResult[E, A], B] {
	mp := result.Map(g)
	return func(rr ReaderResult[E, A]) ReaderResult[D, B] {
		return func(d D) (B, error) {
			return mp(rr(f(d)))
		}
	}
}

// Contramap changes the value of the local environment during the execution of a ReaderResult.
// This is the contravariant functor operation that transforms the input environment.
//
// See: https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#profunctor
//
// Contramap is useful for adapting a ReaderResult to work with a different environment type
// by providing a function that converts the new environment to the expected one.
//
// Type Parameters:
//   - A: The success type (unchanged)
//   - R2: The new input environment type
//   - R1: The original environment type expected by the ReaderResult
//
// Parameters:
//   - f: Function to transform the environment from R2 to R1
//
// Returns:
//   - A Kleisli arrow that takes a ReaderResult[R1, A] and returns a ReaderResult[R2, A]
//
//go:inline
func Contramap[A, R1, R2 any](f func(R2) R1) Kleisli[R2, ReaderResult[R1, A], A] {
	return Local[A](f)
}
