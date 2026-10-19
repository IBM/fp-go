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
	"github.com/IBM/fp-go/v2/internal/catch"
)

// ErrorAs selects the errors of type E.
//
// If errors.As finds an E in the chain of err, ErrorAs returns it as the Right value.
// Otherwise it returns err unchanged as the Left value, so a following Fold can recover
// from errors of type E and propagate all others.
//
// Type Parameters:
//   - E: The error type to select, a concrete error type or an interface that implements error
//
// Parameters:
//   - err: The error to inspect
//
// Returns:
//   - Result[E]: The E found in the chain of err, or err itself as the Left value
//
// See Also:
//   - ErrorIs: Selects the errors that match a sentinel value
//   - errors.As: The function from the fp-go errors package that returns an Option
//
//go:inline
func ErrorAs[E error](err error) Result[E] {
	return catch.ErrorAs[E](err)
}

// ErrorIs selects the errors that match target.
//
// The returned function returns err as the Right value if errors.Is(err, target) holds,
// and as the Left value otherwise.
//
// Parameters:
//   - target: The sentinel error to match with errors.Is
//
// Returns:
//   - Kleisli[error, error]: A function that selects the errors that match target
//
// See Also:
//   - ErrorAs: Selects the errors of a given type
//
//go:inline
func ErrorIs(target error) Kleisli[error, error] {
	return catch.ErrorIs(target)
}
