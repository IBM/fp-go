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

package reader

import (
	"context"

	R "github.com/IBM/fp-go/v2/reader"
)

type (
	// Reader represents a computation that depends on a [context.Context] and produces a value of type A.
	//
	// This is a specialization of the generic Reader monad where the environment type is fixed
	// to [context.Context]. This is particularly useful for Go applications that need to thread
	// context through computations for cancellation, deadlines, and request-scoped values.
	//
	// Type Parameters:
	//   - A: The result type produced by the computation
	//
	// Reader[A] is equivalent to func(context.Context) A
	//
	// The Reader monad enables:
	//   - Dependency injection using context values
	//   - Cancellation and timeout handling
	//   - Request-scoped data propagation
	//   - Avoiding explicit context parameter threading
	//
	// Example:
	//
	//	// A Reader that extracts a user ID from context
	//	getUserID := func(ctx context.Context) string {
	//	    if userID, ok := ctx.Value("userID").(string); ok {
	//	        return userID
	//	    }
	//	    return "anonymous"
	//	}
	//
	//	// A Reader that checks if context is cancelled
	//	isCancelled := func(ctx context.Context) bool {
	//	    select {
	//	    case <-ctx.Done():
	//	        return true
	//	    default:
	//	        return false
	//	    }
	//	}
	//
	//	// Use the readers with a context
	//	ctx := context.WithValue(context.Background(), "userID", "user123")
	//	userID := getUserID(ctx)      // "user123"
	//	cancelled := isCancelled(ctx) // false
	Reader[A any] = R.Reader[context.Context, A]

	// Kleisli represents a Kleisli arrow for the context-based Reader monad.
	//
	// It's a function from A to Reader[B], used for composing Reader computations
	// that all depend on the same [context.Context].
	//
	// Type Parameters:
	//   - A: The input type
	//   - B: The output type wrapped in Reader
	//
	// Kleisli[A, B] is equivalent to func(A) func(context.Context) B
	//
	// Kleisli arrows are fundamental for monadic composition, allowing you to chain
	// operations that depend on context without explicitly passing the context through
	// each function call.
	//
	// Example:
	//
	//	// A Kleisli arrow that creates a greeting Reader from a name
	//	greet := func(name string) Reader[string] {
	//	    return func(ctx context.Context) string {
	//	        if deadline, ok := ctx.Deadline(); ok {
	//	            return fmt.Sprintf("Hello %s (deadline: %v)", name, deadline)
	//	        }
	//	        return fmt.Sprintf("Hello %s", name)
	//	    }
	//	}
	//
	//	// Use the Kleisli arrow
	//	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//	defer cancel()
	//	greeting := greet("Alice")(ctx) // "Hello Alice (deadline: ...)"
	Kleisli[A, B any] = R.Reader[A, Reader[B]]

	// Operator represents a transformation from one Reader to another.
	//
	// It takes a Reader[A] and produces a Reader[B], where both readers depend on
	// the same [context.Context]. This type is commonly used for operations like
	// Map, Chain, and other transformations that convert readers while preserving
	// the context dependency.
	//
	// Type Parameters:
	//   - A: The input Reader's result type
	//   - B: The output Reader's result type
	//
	// Operator[A, B] is equivalent to func(Reader[A]) func(context.Context) B
	//
	// Operators enable building pipelines of context-dependent computations where
	// each step can transform the result of the previous computation while maintaining
	// access to the shared context.
	//
	// Example:
	//
	//	// An operator that transforms int readers to string readers
	//	intToString := func(r Reader[int]) Reader[string] {
	//	    return func(ctx context.Context) string {
	//	        value := r(ctx)
	//	        return strconv.Itoa(value)
	//	    }
	//	}
	//
	//	// A Reader that extracts a timeout value from context
	//	getTimeout := func(ctx context.Context) int {
	//	    if deadline, ok := ctx.Deadline(); ok {
	//	        return int(time.Until(deadline).Seconds())
	//	    }
	//	    return 0
	//	}
	//
	//	// Transform the Reader
	//	getTimeoutStr := intToString(getTimeout)
	//	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	//	defer cancel()
	//	result := getTimeoutStr(ctx) // "30" (approximately)
	Operator[A, B any] = Kleisli[Reader[A], B]
)
