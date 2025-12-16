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

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/idiomatic/result"
)

// WithContext wraps an existing ReaderResult and performs a context check for cancellation before delegating
// to the underlying computation.
//
// If the context has been cancelled (ctx.Err() != nil), it immediately returns an error containing the
// cancellation cause without executing the wrapped computation. Otherwise, it delegates to the original
// ReaderResult.
//
// This is useful for adding cancellation checks to computations that may not check the context themselves,
// ensuring that cancelled operations fail fast.
//
// Example:
//
//	// A computation that might take a long time
//	slowComputation := func(ctx context.Context) (int, error) {
//	    time.Sleep(5 * time.Second)
//	    return 42, nil
//	}
//
//	// Wrap it to check for cancellation before execution
//	safeSlow := readerresult.WithContext(slowComputation)
//
//	// If context is already cancelled, this returns immediately
//	ctx, cancel := context.WithCancel(context.Background())
//	cancel() // Cancel immediately
//	result, err := safeSlow(ctx) // Returns error immediately without sleeping
func WithContext[A any](ma ReaderResult[A]) ReaderResult[A] {
	return func(ctx context.Context) (A, error) {
		if ctx.Err() != nil {
			return result.Left[A](context.Cause(ctx))
		}
		return ma(ctx)
	}
}

// WithContextK wraps a Kleisli arrow (a function that returns a ReaderResult) with context cancellation checking.
//
// This is the Kleisli arrow version of WithContext. It takes a function A -> ReaderResult[B] and returns
// a new function that performs the same transformation but with an added context cancellation check before
// executing the resulting ReaderResult.
//
// A Kleisli arrow is a function that takes a value and returns a monadic computation. In this case,
// Kleisli[A, B] = func(A) ReaderResult[B], which represents a function from A to a context-dependent
// computation that may fail.
//
// WithContextK is particularly useful when composing operations with Chain/Bind, as it ensures that
// each step in the composition checks for cancellation before proceeding.
//
// Parameters:
//   - f: A Kleisli arrow (function from A to ReaderResult[B])
//
// Returns:
//   - A new Kleisli arrow that wraps the result of f with context cancellation checking
//
// Example:
//
//	// A function that fetches user details
//	fetchUserDetails := func(userID int) readerresult.ReaderResult[UserDetails] {
//	    return func(ctx context.Context) (UserDetails, error) {
//	        // Fetch from database...
//	        return details, nil
//	    }
//	}
//
//	// Wrap it to ensure cancellation is checked before each execution
//	safeFetchDetails := readerresult.WithContextK(fetchUserDetails)
//
//	// Use in a composition chain
//	pipeline := F.Pipe2(
//	    getUser(42),
//	    readerresult.Chain(safeFetchDetails), // Checks cancellation before fetching details
//	)
//
//	// If context is cancelled between getUser and fetchUserDetails,
//	// the details fetch will not execute
//	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
//	defer cancel()
//	result, err := pipeline(ctx)
//
// Use Cases:
//   - Adding cancellation checks to composed operations
//   - Ensuring long-running pipelines respect context cancellation
//   - Wrapping third-party functions that don't check context themselves
//   - Creating fail-fast behavior in complex operation chains
func WithContextK[A, B any](f Kleisli[A, B]) Kleisli[A, B] {
	return F.Flow2(
		f,
		WithContext,
	)
}
