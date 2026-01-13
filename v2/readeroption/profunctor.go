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

package readeroption

import (
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
)

// Promap is the profunctor map operation that transforms both the input and output of a ReaderOption.
// It applies f to the input environment (contravariantly) and g to the output value (covariantly).
//
// See: https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#profunctor
//
// This operation allows you to:
//   - Adapt the environment type before passing it to the ReaderOption (via f)
//   - Transform the Some value after the computation completes (via g)
//
// The None case remains unchanged through the transformation.
//
// Type Parameters:
//   - R: The original environment type expected by the ReaderOption
//   - A: The original value type produced by the ReaderOption
//   - D: The new input environment type
//   - B: The new output value type
//
// Parameters:
//   - f: Function to transform the input environment from D to R (contravariant)
//   - g: Function to transform the output Some value from A to B (covariant)
//
// Returns:
//   - A Kleisli arrow that takes a ReaderOption[R, A] and returns a ReaderOption[D, B]
//
//go:inline
func Promap[R, A, D, B any](f func(D) R, g func(A) B) Kleisli[D, ReaderOption[R, A], B] {
	return reader.Promap(f, option.Map(g))
}

// Contramap changes the value of the local environment during the execution of a ReaderOption.
// This is the contravariant functor operation that transforms the input environment.
//
// See: https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#profunctor
//
// Contramap is useful for adapting a ReaderOption to work with a different environment type
// by providing a function that converts the new environment to the expected one.
//
// Type Parameters:
//   - A: The value type (unchanged)
//   - R2: The new input environment type
//   - R1: The original environment type expected by the ReaderOption
//
// Parameters:
//   - f: Function to transform the environment from R2 to R1
//
// Returns:
//   - A Kleisli arrow that takes a ReaderOption[R1, A] and returns a ReaderOption[R2, A]
//
//go:inline
func Contramap[A, R1, R2 any](f func(R2) R1) Kleisli[R2, ReaderOption[R1, A], A] {
	return reader.Contramap[Option[A]](f)
}
