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

package statereaderioresult

import (
	"time"

	RIORES "github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/option"
)

// fromContextOperator lifts a context-scoping operator of the underlying
// ReaderIOResult into an Operator on StateReaderIOResult. The state is threaded
// through unchanged; only the context seen by the computation is affected.
func fromContextOperator[S, A any](op RIORES.Operator[Pair[S, A], Pair[S, A]]) Operator[S, A, A] {
	return func(ma StateReaderIOResult[S, A]) StateReaderIOResult[S, A] {
		return function.Flow2(ma, op)
	}
}

// AskValue creates a StateReaderIOResult that reads a typed value stored in the context under key.
//
// The value is retrieved with [context.Context.Value] and type-asserted to V.
// The computation always succeeds with Some(v) if the key is present and holds a V,
// or None if the key is absent or the stored value has a different type.
// The state is passed through unchanged.
//
// Type Parameters:
//   - S: The state type (must be given explicitly)
//   - V: The expected type of the stored value (must be given explicitly)
//   - K: The type of the key (inferred from the argument)
//
// Parameters:
//   - key: The context key to look up
//
// Returns:
//   - A StateReaderIOResult[S, Option[V]] that yields the value, if present
//
// Example:
//
//	type ctxKey string
//	const userIDKey ctxKey = "userID"
//
//	getUserID := statereaderioresult.AskValue[AppState, string](userIDKey)
//
// See Also:
//   - WithValue: Runs a computation with a value added to the context
//   - Asks: Derives a computation from the whole context
func AskValue[S, V, K any](key K) StateReaderIOResult[S, option.Option[V]] {
	return FromReaderIOResult[S](RIORES.AskValue[V](key))
}

// WithValue runs a StateReaderIOResult computation with an additional key → value
// pair stored in its context.
//
// This is a convenience wrapper around Local that uses [context.WithValue].
// Only the wrapped computation observes the value; neither the caller's context
// nor the state is modified.
//
// Type Parameters:
//   - S: The state type (must be given explicitly)
//   - A: The value type (must be given explicitly)
//   - K: The type of the key (inferred from the argument)
//   - V: The type of the stored value (inferred from the argument)
//
// Parameters:
//   - key: The context key
//   - value: The value to associate with key
//
// Returns:
//   - An Operator that runs the computation with the value in its context
//
// Example:
//
//	result := F.Pipe1(
//	    computation,
//	    statereaderioresult.WithValue[AppState, Data](userIDKey, "user-123"),
//	)
//
// See Also:
//   - AskValue: Reads a value from the context
//   - Local: The general mechanism for transforming the context
func WithValue[S, A, K, V any](key K, value V) Operator[S, A, A] {
	return fromContextOperator[S](RIORES.WithValue[Pair[S, A]](key, value))
}

// WithTimeout adds a timeout to the context for a StateReaderIOResult computation.
//
// This is a convenience wrapper around Local that uses [context.WithTimeout].
// The timeout is relative to when the computation is executed. The cancel
// function is automatically called when the computation completes.
//
// Type Parameters:
//   - S: The state type
//   - A: The value type
//
// Parameters:
//   - timeout: The maximum duration for the computation
//
// Returns:
//   - An Operator that runs the computation with a timeout
//
// Example:
//
//	result := F.Pipe1(
//	    computation,
//	    statereaderioresult.WithTimeout[AppState, Data](5*time.Second),
//	)
func WithTimeout[S, A any](timeout time.Duration) Operator[S, A, A] {
	return fromContextOperator[S](RIORES.WithTimeout[Pair[S, A]](timeout))
}

// WithDeadline adds an absolute deadline to the context for a StateReaderIOResult computation.
//
// This is a convenience wrapper around Local that uses [context.WithDeadline].
// If the parent context already has an earlier deadline, that one takes precedence.
// The cancel function is automatically called when the computation completes.
//
// Type Parameters:
//   - S: The state type
//   - A: The value type
//
// Parameters:
//   - deadline: The absolute time by which the computation must complete
//
// Returns:
//   - An Operator that runs the computation with a deadline
//
// Example:
//
//	result := F.Pipe1(
//	    computation,
//	    statereaderioresult.WithDeadline[AppState, Data](time.Now().Add(time.Minute)),
//	)
func WithDeadline[S, A any](deadline time.Time) Operator[S, A, A] {
	return fromContextOperator[S](RIORES.WithDeadline[Pair[S, A]](deadline))
}
