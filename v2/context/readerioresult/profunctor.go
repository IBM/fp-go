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

package readerioresult

import (
	"context"

	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/result"
)

// Promap is the profunctor map operation that transforms both the input and output of a context-based ReaderIOResult.
// It applies f to the input context (contravariantly) and g to the output value (covariantly).
//
// See: https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#profunctor
//
// This operation allows you to:
//   - Modify the context before passing it to the ReaderIOResult (via f)
//   - Transform the success value after the IO effect completes (via g)
//
// The function f returns both a new context and a CancelFunc that should be called to release resources.
// The error type is fixed as error and remains unchanged through the transformation.
//
// Type Parameters:
//   - A: The original success type produced by the ReaderIOResult
//   - B: The new output success type
//
// Parameters:
//   - f: Function to transform the input context (contravariant)
//   - g: Function to transform the output success value from A to B (covariant)
//
// Returns:
//   - An Operator that takes a ReaderIOResult[A] and returns a ReaderIOResult[B]
//
//go:inline
func Promap[A, B any](f pair.Kleisli[context.CancelFunc, context.Context, context.Context], g func(A) B) Operator[A, B] {
	return function.Flow2(
		Local[A](f),
		Map(g),
	)
}

// Contramap changes the context during the execution of a ReaderIOResult.
// This is the contravariant functor operation that transforms the input context.
//
// See: https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#profunctor
//
// Contramap is an alias for Local and is useful for adapting a ReaderIOResult to work with
// a modified context by providing a function that transforms the context.
//
// Type Parameters:
//   - A: The success type (unchanged)
//
// Parameters:
//   - f: Function to transform the context, returning a new context and CancelFunc
//
// Returns:
//   - An Operator that takes a ReaderIOResult[A] and returns a ReaderIOResult[A]
//
//go:inline
func Contramap[A any](f pair.Kleisli[context.CancelFunc, context.Context, context.Context]) Operator[A, A] {
	return Local[A](f)
}

func ContramapIOK[A any](f io.Kleisli[context.Context, ContextCancel]) Operator[A, A] {
	return LocalIOK[A](f)
}

// LocalIOK transforms the context using an IO-based function before passing it to a ReaderIOResult.
// This is similar to Local but the context transformation itself is wrapped in an IO effect.
//
// The function f takes a context and returns an IO effect that produces a ContextCancel
// (a pair of CancelFunc and the new Context). This allows the context transformation to
// perform side effects.
//
// # Use Cases
//
// This function is useful for sharing information via the Context that is computed through
// side effects that cannot fail, such as:
//   - Generating unique request IDs or trace IDs
//   - Recording timestamps or metrics
//   - Logging context information
//   - Computing derived values from existing context data
//
// The side effect is executed during the context transformation, and the resulting data is
// stored in the context for downstream computations to access.
//
// # Type Parameters
//
//   - A: The success type (unchanged through the transformation)
//
// # Parameters
//
//   - f: An IO-based Kleisli function that transforms the context
//
// # Returns
//
//   - An Operator that applies the context transformation before executing the ReaderIOResult
//
// # Example Usage
//
//	// Generate a request ID via side effect and add to context
//	addRequestID := func(ctx context.Context) io.IO[ContextCancel] {
//	    return func() ContextCancel {
//	        // Side effect: generate unique ID
//	        requestID := uuid.New().String()
//	        // Share the ID via context
//	        newCtx := context.WithValue(ctx, "requestID", requestID)
//	        return pair.MakePair(func() {}, newCtx)
//	    }
//	}
//	adapted := LocalIOK[int](addRequestID)(computation)
//
// # See Also
//
//   - Local: For pure context transformations
//   - LocalIOResultK: For context transformations that can fail
//
//go:inline
func LocalIOK[A any](f io.Kleisli[context.Context, ContextCancel]) Operator[A, A] {
	return LocalIOResultK[A](function.Flow2(f, ioresult.FromIO))
}

// LocalIOResultK transforms the context using an IOResult-based function before passing it to a ReaderIOResult.
// This is similar to Local but the context transformation can fail with an error.
//
// The function f takes a context and returns an IOResult that produces either an error or a ContextCancel
// (a pair of CancelFunc and the new Context). If the transformation fails, the error is propagated
// and the original ReaderIOResult is not executed.
//
// # Use Cases
//
// This function is particularly useful for sharing information via the Context that is computed
// through side effects, such as:
//   - Loading configuration from a file or database
//   - Fetching authentication tokens from an external service
//   - Computing derived values that require I/O operations
//   - Validating and enriching context with data from external sources
//
// The side effect is executed during the context transformation, and the resulting data is
// stored in the context for downstream computations to access.
//
// # Type Parameters
//
//   - A: The success type (unchanged through the transformation)
//
// # Parameters
//
//   - f: An IOResult-based Kleisli function that transforms the context and may fail
//
// # Returns
//
//   - An Operator that applies the context transformation before executing the ReaderIOResult
//
// # Example Usage
//
//	// Load configuration via side effect and add to context
//	loadConfig := func(ctx context.Context) ioresult.IOResult[ContextCancel] {
//	    return func() result.Result[ContextCancel] {
//	        // Side effect: read from file system
//	        config, err := os.ReadFile("config.json")
//	        if err != nil {
//	            return result.Left[ContextCancel](err)
//	        }
//	        // Share the loaded config via context
//	        newCtx := context.WithValue(ctx, "config", config)
//	        return result.Of(pair.MakePair(func() {}, newCtx))
//	    }
//	}
//	adapted := LocalIOResultK[int](loadConfig)(computation)
//
// # See Also
//
//   - Local: For pure context transformations
//   - LocalIOK: For context transformations with side effects that cannot fail
//
//go:inline
func LocalIOResultK[A any](f ioresult.Kleisli[context.Context, ContextCancel]) Operator[A, A] {
	return func(rr ReaderIOResult[A]) ReaderIOResult[A] {
		return func(ctx context.Context) IOResult[A] {
			return func() Result[A] {
				if ctx.Err() != nil {
					return result.Left[A](context.Cause(ctx))
				}
				p, err := result.Unwrap(f(ctx)())
				if err != nil {
					return result.Left[A](err)
				}
				// unwrap
				otherCancel, otherCtx := pair.Unpack(p)
				defer otherCancel()
				return rr(otherCtx)()
			}
		}
	}
}
