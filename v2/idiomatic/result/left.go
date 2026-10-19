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

package result

// ChainFirstLeft runs f on the error for its side effect and keeps the original error.
//
// The value that f returns, success or failure, is discarded. A successful value passes
// through without calling f.
//
// Type Parameters:
//   - A: The type of the success value
//   - B: The result type of the side effect (discarded)
//
// Parameters:
//   - f: The side effect for the error
//
// Returns:
//   - Operator[A, A]: An operator that preserves the original value or error
//
// See Also:
//   - TapLeft: An alias
//   - OrElse: Recovers from the error instead of observing it
func ChainFirstLeft[A, B any](f Kleisli[error, B]) Operator[A, A] {
	return func(a A, err error) (A, error) {
		if err != nil {
			_, _ = f(err)
			return Left[A](err)
		}
		return Of(a)
	}
}

// TapLeft is an alias for ChainFirstLeft.
//
// Type Parameters:
//   - A: The type of the success value
//   - B: The result type of the side effect (discarded)
//
// Parameters:
//   - f: The side effect for the error
//
// Returns:
//   - Operator[A, A]: An operator that preserves the original value or error
//
// See Also:
//   - ChainFirstLeft: The function this is an alias for
//   - TapErrorAs: Observes only errors of a given type
//
//go:inline
func TapLeft[A, B any](f Kleisli[error, B]) Operator[A, A] {
	return ChainFirstLeft[A](f)
}
