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
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/pair"
	RIO "github.com/IBM/fp-go/v2/readerio"
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
//   - R: The input environment type that f transforms into context.Context
//   - A: The original result type produced by the ReaderIO
//   - B: The new output result type
//
// Parameters:
//   - f: Function to transform the input environment R into context.Context (contravariant)
//   - g: Function to transform the output value from A to B (covariant)
//
// Returns:
//   - A Kleisli arrow that takes a ReaderIO[A] and returns a function from R to B
//
// Note: When R is context.Context, this simplifies to an Operator[A, B]
//
//go:inline
func Promap[R, A, B any](f pair.Kleisli[context.CancelFunc, R, context.Context], g func(A) B) RIO.Kleisli[R, ReaderIO[A], B] {
	return function.Flow2(
		Local[A](f),
		RIO.Map[R](g),
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
//   - R: The input environment type that f transforms into context.Context
//
// Parameters:
//   - f: Function to transform the input environment R into context.Context, returning a new context and CancelFunc
//
// Returns:
//   - A Kleisli arrow that takes a ReaderIO[A] and returns a function from R to A
//
// Note: When R is context.Context, this simplifies to an Operator[A, A]
//
//go:inline
func Contramap[A, R any](f pair.Kleisli[context.CancelFunc, R, context.Context]) RIO.Kleisli[R, ReaderIO[A], A] {
	return Local[A](f)
}

// LocalIOK transforms the context using an IO effect before passing it to a ReaderIO computation.
//
// This is similar to Local, but the context transformation itself is wrapped in an IO effect,
// allowing for side-effectful context transformations. The transformation function receives
// the current context and returns an IO effect that produces a new context along with a
// cancel function. The cancel function is automatically called when the computation completes
// (via defer), ensuring proper cleanup of resources.
//
// This is useful for:
//   - Context transformations that require side effects (e.g., loading configuration)
//   - Lazy initialization of context values
//   - Context transformations that may fail or need to perform I/O
//   - Composing effectful context setup with computations
//
// Type Parameters:
//   - A: The value type of the ReaderIO
//
// Parameters:
//   - f: An IO Kleisli arrow that transforms the context with side effects
//
// Returns:
//   - An Operator that runs the computation with the effectfully transformed context
//
// Example:
//
//	import (
//	    "context"
//	    G "github.com/IBM/fp-go/v2/io"
//	    F "github.com/IBM/fp-go/v2/function"
//	)
//
//	// Context transformation with side effects (e.g., loading config)
//	loadConfig := func(ctx context.Context) G.IO[ContextCancel] {
//	    return func() ContextCancel {
//	        // Simulate loading configuration
//	        config := loadConfigFromFile()
//	        newCtx := context.WithValue(ctx, "config", config)
//	        return pair.MakePair[context.CancelFunc](func() {}, newCtx)
//	    }
//	}
//
//	getValue := readerio.FromReader(func(ctx context.Context) string {
//	    if cfg := ctx.Value("config"); cfg != nil {
//	        return cfg.(string)
//	    }
//	    return "default"
//	})
//
//	result := F.Pipe1(
//	    getValue,
//	    readerio.LocalIOK[string](loadConfig),
//	)
//	value := result(t.Context())()  // Loads config and uses it
//
// Comparison with Local:
//   - Local: Takes a pure function that transforms the context
//   - LocalIOK: Takes an IO effect that transforms the context, allowing side effects
func LocalIOK[A any](f io.Kleisli[context.Context, ContextCancel]) Operator[A, A] {
	return func(r ReaderIO[A]) ReaderIO[A] {
		return func(ctx context.Context) IO[A] {
			p := f(ctx)
			return func() A {
				otherCancel, otherCtx := pair.Unpack(p())
				defer otherCancel()
				return r(otherCtx)()
			}
		}
	}
}
