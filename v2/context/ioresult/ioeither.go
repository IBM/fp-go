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

package ioresult

import (
	"context"

	"github.com/IBM/fp-go/v2/result"
)

// WithContext wraps an IOResult and performs a context check for cancellation before executing.
// This ensures that if the context is already cancelled, the computation short-circuits immediately
// without executing the wrapped computation.
//
// This is useful for adding cancellation awareness to computations that might not check the context themselves.
//
// Type Parameters:
//   - A: The type of the success value
//
// Parameters:
//   - ctx: The context to check for cancellation
//   - ma: The IOResult to wrap with context checking
//
// Returns:
//   - An IOResult that checks for cancellation before executing
//
// Example:
//
//	computation := func() Result[string] {
//	    // Long-running operation
//	    return result.Of("done")
//	}
//
//	ctx, cancel := context.WithCancel(t.Context())
//	cancel() // Cancel immediately
//
//	wrapped := WithContext(ctx, computation)
//	result := wrapped() // Returns Left with context.Canceled error
func WithContext[A any](ctx context.Context, ma IOResult[A]) IOResult[A] {
	return func() Result[A] {
		if ctx.Err() != nil {
			return result.Left[A](context.Cause(ctx))
		}
		return ma()
	}
}
