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

package readerio

import (
	"context"

	"github.com/IBM/fp-go/v2/function"
)

// Promap is the profunctor map operation that transforms both the input and output of a context-based ReaderIO.
// It applies f to the input context (contravariantly) and g to the output value (covariantly).
//
// See: https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#profunctor
//
// This operation allows you to:
//   - Modify the context before passing it to the ReaderIO (via f)
//   - Transform the result value after the IO effect completes (via g)
//
// The function f returns both a new context and a CancelFunc that should be called to release resources.
//
// Type Parameters:
//   - A: The original result type produced by the ReaderIO
//   - B: The new output result type
//
// Parameters:
//   - f: Function to transform the input context (contravariant)
//   - g: Function to transform the output value from A to B (covariant)
//
// Returns:
//   - An Operator that takes a ReaderIO[A] and returns a ReaderIO[B]
//
//go:inline
func Promap[A, B any](f func(context.Context) (context.Context, context.CancelFunc), g func(A) B) Operator[A, B] {
	return function.Flow2(
		Local[A](f),
		Map(g),
	)
}

// Contramap changes the context during the execution of a ReaderIO.
// This is the contravariant functor operation that transforms the input context.
//
// See: https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#profunctor
//
// Contramap is an alias for Local and is useful for adapting a ReaderIO to work with
// a modified context by providing a function that transforms the context.
//
// Type Parameters:
//   - A: The result type (unchanged)
//
// Parameters:
//   - f: Function to transform the context, returning a new context and CancelFunc
//
// Returns:
//   - An Operator that takes a ReaderIO[A] and returns a ReaderIO[A]
//
//go:inline
func Contramap[A any](f func(context.Context) (context.Context, context.CancelFunc)) Operator[A, A] {
	return Local[A](f)
}
