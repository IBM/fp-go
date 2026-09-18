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
	IC "github.com/IBM/fp-go/v2/internal/context"
)

// AskValue creates a ReaderIOResult that reads a typed value stored in the context under key.
//
// The value is retrieved with [context.Context.Value] and type-asserted to V.
// The computation always succeeds with Some(v) if the key is present and holds a V,
// or None if the key is absent or the stored value has a different type.
//
// AskValue is the read counterpart of [WithValue]. To treat a missing value as an
// error, combine it with a Chain over the Option.
//
// Type Parameters:
//   - V: The expected type of the stored value (must be given explicitly)
//   - K: The type of the key (inferred from the argument)
//
// Parameters:
//   - key: The context key to look up
//
// Returns:
//   - A ReaderIOResult[Option[V]] that yields the value, if present
//
// Example:
//
//	type contextKey string
//	const userKey contextKey = "user"
//
//	getUser := readerioresult.AskValue[User](userKey)
//	user := getUser(ctx)() // Some(User{...}) or None
//
// See Also:
//   - WithValue: Runs a computation with a value added to the context
//   - Ask: Retrieves the whole context
func AskValue[V, K any](key K) ReaderIOResult[Option[V]] {
	return FromReader(IC.AskValue[V](key))
}

// WithValue runs a ReaderIOResult computation with an additional key → value
// pair stored in its context.
//
// This is a convenience wrapper around Local that uses [context.WithValue].
// Only the wrapped computation observes the value; the caller's context is not
// modified. Since no resources are allocated, the associated cancel function
// is a no-op.
//
// Type Parameters:
//   - A: The value type of the ReaderIOResult (must be given explicitly)
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
//	type contextKey string
//	const userKey contextKey = "user"
//
//	greet := F.Pipe1(
//	    readerioresult.AskValue[string](userKey),
//	    readerioresult.Map(O.GetOrElse(F.Constant("anonymous"))),
//	)
//
//	result := F.Pipe1(
//	    greet,
//	    readerioresult.WithValue[string](userKey, "Alice"),
//	)
//	result(context.Background())() // Right("Alice")
//
// See Also:
//   - AskValue: Reads a value from the context
//   - Local: The general mechanism for transforming the context
//   - WithTimeout: Runs a computation with a timeout
func WithValue[A, K, V any](key K, value V) Operator[A, A] {
	return Local[A](IC.WithValueNopCancel(key, value))
}
