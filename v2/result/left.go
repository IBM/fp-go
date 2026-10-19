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

import (
	"github.com/IBM/fp-go/v2/identity"
	"github.com/IBM/fp-go/v2/internal/eithert"
)

// MonadChainFirstLeft runs f on the error of a failed Result for its side effect and keeps
// the original error.
//
// The value that f returns, success or failure, is discarded. A successful Result passes
// through without calling f.
//
// Type Parameters:
//   - A: The type of the success value
//   - B: The result type of the side effect (discarded)
//
// Parameters:
//   - ma: The Result to observe
//   - f: The side effect for the error
//
// Returns:
//   - Result[A]: ma unchanged
//
// See Also:
//   - ChainFirstLeft: The curried version
//   - MonadChainLeft: Recovers from the error instead of observing it
//
//go:inline
func MonadChainFirstLeft[A, B any](ma Result[A], f Kleisli[error, B]) Result[A] {
	return eithert.MonadChainFirstLeft(
		identity.MonadChain[Result[A], Result[A]],
		identity.MonadMap[Result[B], Result[A]],
		identity.Of[Result[A]],
		ma,
		f,
	)
}

// ChainFirstLeft runs f on the error of a failed Result for its side effect and keeps the
// original error.
//
// The value that f returns, success or failure, is discarded. A successful Result passes
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
//   - MonadChainFirstLeft: The uncurried version
//   - TapLeft: An alias
//   - ChainLeft: Recovers from the error instead of observing it
//
//go:inline
func ChainFirstLeft[A, B any](f Kleisli[error, B]) Operator[A, A] {
	return eithert.ChainFirstLeft(
		identity.Chain[Result[A], Result[A]],
		identity.Map[Result[B], Result[A]],
		identity.Of[Result[A]],
		f,
	)
}

// MonadTapLeft is an alias for MonadChainFirstLeft.
//
// Type Parameters:
//   - A: The type of the success value
//   - B: The result type of the side effect (discarded)
//
// Parameters:
//   - ma: The Result to observe
//   - f: The side effect for the error
//
// Returns:
//   - Result[A]: ma unchanged
//
// See Also:
//   - MonadChainFirstLeft: The function this is an alias for
//
//go:inline
func MonadTapLeft[A, B any](ma Result[A], f Kleisli[error, B]) Result[A] {
	return MonadChainFirstLeft(ma, f)
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
