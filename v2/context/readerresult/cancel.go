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

package readerresult

import (
	"context"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
)

// WithContext wraps an existing ReaderResult and performs a context check for cancellation
// before delegating to the wrapped computation. This provides early cancellation detection,
// allowing computations to fail fast when the context has been cancelled or has exceeded
// its deadline.
//
// IMPORTANT: This function checks for context cancellation BEFORE executing the wrapped
// ReaderResult. If the context is already cancelled or has exceeded its deadline, the
// computation returns immediately with the cancellation error without executing the
// wrapped ReaderResult.
//
// The function uses context.Cause(ctx) to extract the cancellation reason, which may be:
//   - context.Canceled: The context was explicitly cancelled
//   - context.DeadlineExceeded: The context's deadline was exceeded
//   - A custom error: If the context was cancelled with a cause (Go 1.20+)
//
// Type Parameters:
//   - A: The success type of the ReaderResult
//
// Parameters:
//   - ma: The ReaderResult to wrap with cancellation checking
//
// Returns:
//   - A ReaderResult that checks for cancellation before executing ma
//
// Example:
//
//	// Create a long-running computation
//	longComputation := func(ctx context.Context) result.Result[int] {
//	    time.Sleep(5 * time.Second)
//	    return result.Of(42)
//	}
//
//	// Wrap with cancellation check
//	safeLongComputation := readerresult.WithContext(longComputation)
//
//	// Cancel the context before execution
//	ctx, cancel := context.WithCancel(t.Context())
//	cancel()
//
//	// The computation returns immediately with cancellation error
//	result := safeLongComputation(ctx)
//	// result is Left(context.Canceled) - longComputation never executes
//
// Example with timeout:
//
//	fetchData := func(ctx context.Context) result.Result[string] {
//	    // Simulate slow operation
//	    time.Sleep(2 * time.Second)
//	    return result.Of("data")
//	}
//
//	safeFetch := readerresult.WithContext(fetchData)
//
//	// Context with 1 second timeout
//	ctx, cancel := context.WithTimeout(t.Context(), 1*time.Second)
//	defer cancel()
//
//	time.Sleep(1500 * time.Millisecond) // Wait for timeout
//
//	result := safeFetch(ctx)
//	// result is Left(context.DeadlineExceeded) - fetchData never executes
//
// Use cases:
//   - Wrapping expensive computations to enable early cancellation
//   - Preventing unnecessary work when context is already cancelled
//   - Implementing timeout-aware operations
//   - Building cancellation-aware pipelines
//
//go:inline
func WithContext[A any](ma ReaderResult[A]) ReaderResult[A] {
	return func(ctx context.Context) E.Either[error, A] {
		if ctx.Err() != nil {
			return E.Left[A](context.Cause(ctx))
		}
		return ma(ctx)
	}
}

// WithContextK wraps a Kleisli arrow with context cancellation checking.
// This is a higher-order function that takes a Kleisli arrow and returns a new
// Kleisli arrow that checks for context cancellation before executing.
//
// IMPORTANT: This function composes the Kleisli arrow with WithContext, ensuring
// that the resulting ReaderResult checks for cancellation before execution. This
// is particularly useful when building pipelines of Kleisli arrows where you want
// cancellation checking at each step.
//
// Type Parameters:
//   - A: The input type of the Kleisli arrow
//   - B: The output type of the Kleisli arrow
//
// Parameters:
//   - f: The Kleisli arrow to wrap with cancellation checking
//
// Returns:
//   - A new Kleisli arrow that checks for cancellation before executing f
//
// Example:
//
//	// Define a Kleisli arrow
//	processUser := func(id int) readerresult.ReaderResult[User] {
//	    return func(ctx context.Context) result.Result[User] {
//	        // Expensive database operation
//	        return fetchUserFromDB(ctx, id)
//	    }
//	}
//
//	// Wrap with cancellation checking
//	safeProcessUser := readerresult.WithContextK(processUser)
//
//	// Use in a pipeline
//	pipeline := F.Pipe1(
//	    readerresult.Of(123),
//	    readerresult.Chain(safeProcessUser),
//	)
//
//	// If context is cancelled, processUser never executes
//	ctx, cancel := context.WithCancel(t.Context())
//	cancel()
//	result := pipeline(ctx) // Left(context.Canceled)
//
// Example with multiple steps:
//
//	getUserK := readerresult.WithContextK(func(id int) readerresult.ReaderResult[User] {
//	    return func(ctx context.Context) result.Result[User] {
//	        return fetchUser(ctx, id)
//	    }
//	})
//
//	getOrdersK := readerresult.WithContextK(func(user User) readerresult.ReaderResult[[]Order] {
//	    return func(ctx context.Context) result.Result[[]Order] {
//	        return fetchOrders(ctx, user.ID)
//	    }
//	})
//
//	// Each step checks for cancellation
//	pipeline := F.Pipe2(
//	    readerresult.Of(123),
//	    readerresult.Chain(getUserK),
//	    readerresult.Chain(getOrdersK),
//	)
//
//	// If context is cancelled at any point, remaining steps don't execute
//	ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
//	defer cancel()
//	result := pipeline(ctx)
//
// Use cases:
//   - Building cancellation-aware pipelines
//   - Ensuring each step in a chain respects cancellation
//   - Implementing timeout-aware multi-step operations
//   - Preventing cascading failures in long pipelines
//
//go:inline
func WithContextK[A, B any](f Kleisli[A, B]) Kleisli[A, B] {
	return F.Flow2(
		f,
		WithContext,
	)
}
