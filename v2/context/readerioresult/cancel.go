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

package readerioresult

import (
	"context"

	CIOE "github.com/IBM/fp-go/v2/context/ioresult"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/ioeither"
)

// WithContext wraps an existing [ReaderIOResult] and performs a context check for cancellation before delegating.
// This ensures that if the context is already canceled, the computation short-circuits immediately
// without executing the wrapped computation.
//
// This is useful for adding cancellation awareness to computations that might not check the context themselves.
//
// Parameters:
//   - ma: The ReaderIOResult to wrap with context checking
//
// Returns a ReaderIOResult that checks for cancellation before executing.
func WithContext[A any](ma ReaderIOResult[A]) ReaderIOResult[A] {
	return func(ctx context.Context) IOEither[A] {
		if ctx.Err() != nil {
			return ioeither.Left[A](context.Cause(ctx))
		}
		return CIOE.WithContext(ctx, ma(ctx))
	}
}

// WithContextK wraps a Kleisli arrow with context cancellation checking.
// This ensures that the computation checks for context cancellation before executing,
// providing a convenient way to add cancellation awareness to Kleisli arrows.
//
// This is particularly useful when composing multiple Kleisli arrows where each step
// should respect context cancellation.
//
// Type Parameters:
//   - A: The input type of the Kleisli arrow
//   - B: The output type of the Kleisli arrow
//
// Parameters:
//   - f: The Kleisli arrow to wrap with context checking
//
// Returns:
//   - A Kleisli arrow that checks for cancellation before executing
//
// Example:
//
//	fetchUser := func(id int) ReaderIOResult[User] {
//	    return func(ctx context.Context) IOResult[User] {
//	        return func() Result[User] {
//	            // Long-running operation
//	            return result.Of(User{ID: id})
//	        }
//	    }
//	}
//
//	// Wrap with context checking
//	safeFetch := WithContextK(fetchUser)
//
//	// If context is cancelled, returns immediately without executing fetchUser
//	ctx, cancel := context.WithCancel(t.Context())
//	cancel() // Cancel immediately
//	result := safeFetch(123)(ctx)() // Returns context.Canceled error
//
//go:inline
func WithContextK[A, B any](f Kleisli[A, B]) Kleisli[A, B] {
	return F.Flow2(
		f,
		WithContext,
	)
}
