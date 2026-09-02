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

package either

import C "github.com/IBM/fp-go/v2/internal/common"

// Exists creates a predicate that tests whether an Either value is Right and its value satisfies the given predicate.
// It returns a function that takes an Either[E, T] and returns true only if the Either is Right and the predicate p
// returns true for the Right value.
//
// This function is useful for checking if an Either contains a successful value that meets certain criteria,
// commonly used in filtering operations, validation chains, or conditional logic where you need to verify
// both the success state and a property of the success value.
//
// The behavior is as follows:
//   - If the input is Left, returns false (regardless of the predicate)
//   - If the input is Right and p returns true for the Right value, returns true
//   - If the input is Right and p returns false for the Right value, returns false
//
// Type Parameters:
//   - E: The type of the Left value (error type)
//   - T: The type of the Right value (success type)
//
// Parameters:
//   - p: A predicate function that tests values of type T
//
// Returns:
//
//	A Predicate function that takes an Either[E, T] and returns true if it's Right and satisfies p
//
// See Also:
//   - ExistsLeft: Tests if Either is Left and satisfies a predicate
//   - Filter: Converts Right values that fail a predicate to Left
func Exists[E, T any](p Predicate[T]) Predicate[Either[E, T]] {
	return C.EitherExists[E, T](p)
}

// ExistsLeft creates a predicate that tests whether an Either value is Left and its value satisfies the given predicate.
// It returns a function that takes an Either[E, T] and returns true only if the Either is Left and the predicate p
// returns true for the Left value.
//
// The behavior is as follows:
//   - If the input is Right, returns false (regardless of the predicate)
//   - If the input is Left and p returns true for the Left value, returns true
//   - If the input is Left and p returns false for the Left value, returns false
//
// Type Parameters:
//   - T: The type of the Right value (success type)
//   - E: The type of the Left value (error type)
//
// Parameters:
//   - p: A predicate function that tests values of type E
//
// Returns:
//
//	A Predicate function that takes an Either[E, T] and returns true if it's Left and satisfies p
//
// See Also:
//   - Exists: Tests if Either is Right and satisfies a predicate
//   - IsLeft: Tests if Either is Left without checking the value
func ExistsLeft[T, E any](p Predicate[E]) Predicate[Either[E, T]] {
	return C.EitherExistsLeft[T, E](p)
}

// ForAll creates a predicate that tests whether an Either value is Left or its Right value satisfies the given predicate.
// It returns a function that takes an Either[E, T] and returns true if the Either is Left (regardless of its value)
// or if it's Right and the predicate p returns true for the Right value.
//
// The behavior is as follows:
//   - If the input is Left, returns true (vacuous truth - predicate holds for empty case)
//   - If the input is Right and p returns true for the Right value, returns true
//   - If the input is Right and p returns false for the Right value, returns false
//
// Type Parameters:
//   - E: The type of the Left value (error type)
//   - T: The type of the Right value (success type)
//
// Parameters:
//   - p: A predicate function that tests values of type T
//
// Returns:
//
//	A Predicate function that takes an Either[E, T] and returns true if it's Left or Right with p satisfied
//
// See Also:
//   - Exists: Tests if Either is Right and satisfies a predicate (existential quantification)
//   - ExistsLeft: Tests if Either is Left and satisfies a predicate
//   - Filter: Converts Right values that fail a predicate to Left
func ForAll[E, T any](p Predicate[T]) Predicate[Either[E, T]] {
	return C.EitherForAll[E, T](p)
}
