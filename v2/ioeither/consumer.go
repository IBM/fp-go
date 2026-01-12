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

package ioeither

import "github.com/IBM/fp-go/v2/io"

// ChainConsumer converts a Consumer into an IOEither operator that executes the consumer
// as a side effect on successful (Right) values and returns an empty struct.
//
// This function bridges the gap between pure consumers (functions that consume values
// without returning anything) and the IOEither monad. It takes a Consumer[A] and returns
// an Operator that:
//  1. If the IOEither is Right, executes the consumer with the value as a side effect
//  2. If the IOEither is Left, propagates the error without calling the consumer
//  3. Returns IOEither[E, struct{}] to maintain the monadic chain
//
// The consumer is only executed for successful (Right) values. Errors (Left values) are
// propagated unchanged. This is useful for operations like logging successful results,
// collecting metrics, or updating external state within an IOEither pipeline.
//
// Type Parameters:
//   - E: The error type of the IOEither
//   - A: The type of value consumed by the consumer
//
// Parameters:
//   - c: A Consumer[A] that performs side effects on values of type A
//
// Returns:
//   - An Operator[E, A, struct{}] that executes the consumer on Right values and returns an empty struct
//
// Example:
//
//	// Create a consumer that logs successful values
//	logger := func(x int) {
//	    fmt.Printf("Success: %d\n", x)
//	}
//
//	// Convert it to an IOEither operator
//	logOp := ioeither.ChainConsumer[error](logger)
//
//	// Use it in an IOEither pipeline
//	result := F.Pipe2(
//	    ioeither.Right[error](42),
//	    logOp,                    // Logs "Success: 42"
//	    ioeither.Map[error](func(struct{}) string { return "done" }),
//	)
//	result() // Returns Right("done") after logging
//
//	// Errors are propagated without calling the consumer
//	errorResult := F.Pipe2(
//	    ioeither.Left[int](errors.New("failed")),
//	    logOp,                    // Consumer NOT called
//	    ioeither.Map[error](func(struct{}) string { return "done" }),
//	)
//	errorResult() // Returns Left(error) without logging
//
//	// Example with data collection
//	var successfulValues []int
//	collector := func(x int) {
//	    successfulValues = append(successfulValues, x)
//	}
//
//	pipeline := F.Pipe2(
//	    ioeither.Right[error](100),
//	    ioeither.ChainConsumer[error](collector),  // Collects the value
//	    ioeither.Map[error](func(struct{}) int { return len(successfulValues) }),
//	)
//	count := pipeline() // Returns Right(1), successfulValues contains [100]
//
//go:inline
func ChainConsumer[E, A any](c Consumer[A]) Operator[E, A, struct{}] {
	return ChainIOK[E](io.FromConsumer(c))
}

//go:inline
func ChainFirstConsumer[E, A any](c Consumer[A]) Operator[E, A, A] {
	return ChainFirstIOK[E](io.FromConsumer(c))
}
