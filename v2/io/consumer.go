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

package io

import "github.com/IBM/fp-go/v2/function"

// ChainConsumer converts a Consumer into an IO operator that executes the consumer
// as a side effect and returns an empty struct.
//
// This function bridges the gap between pure consumers (functions that consume values
// without returning anything) and the IO monad. It takes a Consumer[A] and returns
// an Operator that:
//  1. Executes the source IO[A] to get a value
//  2. Passes that value to the consumer for side effects
//  3. Returns IO[struct{}] to maintain the monadic chain
//
// The returned IO[struct{}] allows the operation to be composed with other IO operations
// while discarding the consumed value. This is useful for operations like logging,
// printing, or updating external state within an IO pipeline.
//
// Type Parameters:
//   - A: The type of value consumed by the consumer
//
// Parameters:
//   - c: A Consumer[A] that performs side effects on values of type A
//
// Returns:
//   - An Operator[A, struct{}] that executes the consumer and returns an empty struct
//
// Example:
//
//	// Create a consumer that logs values
//	logger := func(x int) {
//	    fmt.Printf("Value: %d\n", x)
//	}
//
//	// Convert it to an IO operator
//	logOp := io.ChainConsumer(logger)
//
//	// Use it in an IO pipeline
//	result := F.Pipe2(
//	    io.Of(42),
//	    logOp,                    // Logs "Value: 42"
//	    io.Map(func(struct{}) string { return "done" }),
//	)
//	result() // Returns "done" after logging
//
//	// Another example with multiple operations
//	var values []int
//	collector := func(x int) {
//	    values = append(values, x)
//	}
//
//	pipeline := F.Pipe2(
//	    io.Of(100),
//	    io.ChainConsumer(collector),  // Collects the value
//	    io.Map(func(struct{}) int { return len(values) }),
//	)
//	count := pipeline() // Returns 1, values contains [100]
func ChainConsumer[A any](c Consumer[A]) Operator[A, Void] {
	return Chain(FromConsumer(c))
}

// FromConsumer converts a Consumer into a Kleisli arrow that wraps the consumer
// in an IO context.
//
// This function lifts a Consumer[A] (a function that consumes a value and performs
// side effects) into a Kleisli[A, struct{}] (a function that takes a value and returns
// an IO computation that performs the side effect and returns an empty struct).
//
// The resulting Kleisli arrow can be used with Chain and other monadic operations
// to integrate consumers into IO pipelines. This is a lower-level function compared
// to ChainConsumer, which directly returns an Operator.
//
// Type Parameters:
//   - A: The type of value consumed by the consumer
//
// Parameters:
//   - c: A Consumer[A] that performs side effects on values of type A
//
// Returns:
//   - A Kleisli[A, struct{}] that wraps the consumer in an IO context
//
// Example:
//
//	// Create a consumer
//	logger := func(x int) {
//	    fmt.Printf("Logging: %d\n", x)
//	}
//
//	// Convert to Kleisli arrow
//	logKleisli := io.FromConsumer(logger)
//
//	// Use with Chain
//	result := F.Pipe2(
//	    io.Of(42),
//	    io.Chain(logKleisli),  // Logs "Logging: 42"
//	    io.Map(func(struct{}) string { return "completed" }),
//	)
//	result() // Returns "completed"
//
//	// Can also be used to build more complex operations
//	logAndCount := func(x int) io.IO[int] {
//	    return F.Pipe2(
//	        logKleisli(x),
//	        io.Map(func(struct{}) int { return 1 }),
//	    )
//	}
func FromConsumer[A any](c Consumer[A]) Kleisli[A, Void] {
	return func(a A) IO[Void] {
		return func() Void {
			c(a)
			return function.VOID
		}
	}
}
