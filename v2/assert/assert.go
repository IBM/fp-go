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

// Package assert provides functional assertion helpers for testing.
//
// This package wraps testify/assert functions in a Reader monad pattern,
// allowing for composable and functional test assertions. Each assertion
// returns a Reader that takes a *testing.T and performs the assertion.
//
// The package supports:
//   - Equality and inequality assertions
//   - Collection assertions (arrays, maps, strings)
//   - Error handling assertions
//   - Result type assertions
//   - Custom predicate assertions
//   - Composable test suites
//
// Example:
//
//	func TestExample(t *testing.T) {
//	    value := 42
//	    assert.Equal(42)(value)(t)  // Curried style
//
//	    // Composing multiple assertions
//	    arr := []int{1, 2, 3}
//	    assertions := assert.AllOf([]assert.Reader{
//	        assert.ArrayNotEmpty(arr),
//	        assert.ArrayLength[int](3)(arr),
//	        assert.ArrayContains(2)(arr),
//	    })
//	    assertions(t)
//	}
package assert

import (
	"fmt"
	"testing"

	"github.com/IBM/fp-go/v2/boolean"
	"github.com/IBM/fp-go/v2/eq"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

var (
	// Eq is the equal predicate checking if objects are equal
	Eq = eq.FromEquals(assert.ObjectsAreEqual)
)

// wrap1 is an internal helper function that wraps testify assertion functions
// into the Reader monad pattern with curried parameters.
//
// It takes a testify assertion function and converts it into a curried function
// that first takes an expected value, then an actual value, and finally returns
// a Reader that performs the assertion when given a *testing.T.
//
// Parameters:
//   - wrapped: The testify assertion function to wrap
//   - expected: The expected value for comparison
//   - msgAndArgs: Optional message and arguments for assertion failure
//
// Returns:
//   - A Kleisli function that takes the actual value and returns a Reader
func wrap1[T any](wrapped func(t assert.TestingT, expected, actual any, msgAndArgs ...any) bool, expected T, msgAndArgs ...any) Kleisli[T] {
	return func(actual T) Reader {
		return func(t *testing.T) bool {
			return wrapped(t, expected, actual, msgAndArgs...)
		}
	}
}

// NotEqual tests if the expected and the actual values are not equal
func NotEqual[T any](expected T) Kleisli[T] {
	return wrap1(assert.NotEqual, expected)
}

// Equal tests if the expected and the actual values are equal
func Equal[T any](expected T) Kleisli[T] {
	return wrap1(assert.Equal, expected)
}

// ArrayNotEmpty checks if an array is not empty
func ArrayNotEmpty[T any](arr []T) Reader {
	return func(t *testing.T) bool {
		return assert.NotEmpty(t, arr)
	}
}

// RecordNotEmpty checks if an map is not empty
func RecordNotEmpty[K comparable, T any](mp map[K]T) Reader {
	return func(t *testing.T) bool {
		return assert.NotEmpty(t, mp)
	}
}

// ArrayLength tests if an array has the expected length
func ArrayLength[T any](expected int) Kleisli[[]T] {
	return func(actual []T) Reader {
		return func(t *testing.T) bool {
			return assert.Len(t, actual, expected)
		}
	}
}

// RecordLength tests if a map has the expected length
func RecordLength[K comparable, T any](expected int) Kleisli[map[K]T] {
	return func(actual map[K]T) Reader {
		return func(t *testing.T) bool {
			return assert.Len(t, actual, expected)
		}
	}
}

// StringLength tests if a string has the expected length
func StringLength[K comparable, T any](expected int) Kleisli[string] {
	return func(actual string) Reader {
		return func(t *testing.T) bool {
			return assert.Len(t, actual, expected)
		}
	}
}

// NoError validates that there is no error
func NoError(err error) Reader {
	return func(t *testing.T) bool {
		return assert.NoError(t, err)
	}
}

// Error validates that there is an error
func Error(err error) Reader {
	return func(t *testing.T) bool {
		return assert.Error(t, err)
	}
}

// Success checks if a [Result] represents success
func Success[T any](res Result[T]) Reader {
	return NoError(result.ToError(res))
}

// Failure checks if a [Result] represents failure
func Failure[T any](res Result[T]) Reader {
	return Error(result.ToError(res))
}

// ArrayContains tests if a value is contained in an array
func ArrayContains[T any](expected T) Kleisli[[]T] {
	return func(actual []T) Reader {
		return func(t *testing.T) bool {
			return assert.Contains(t, actual, expected)
		}
	}
}

// ContainsKey tests if a key is contained in a map
func ContainsKey[T any, K comparable](expected K) Kleisli[map[K]T] {
	return func(actual map[K]T) Reader {
		return func(t *testing.T) bool {
			return assert.Contains(t, actual, expected)
		}
	}
}

// NotContainsKey tests if a key is not contained in a map
func NotContainsKey[T any, K comparable](expected K) Kleisli[map[K]T] {
	return func(actual map[K]T) Reader {
		return func(t *testing.T) bool {
			return assert.NotContains(t, actual, expected)
		}
	}
}

// That asserts that a particular predicate matches
func That[T any](pred Predicate[T]) Kleisli[T] {
	return func(a T) Reader {
		return func(t *testing.T) bool {
			if pred(a) {
				return true
			}
			return assert.Fail(t, fmt.Sprintf("Preficate %v does not match value %v", pred, a))
		}
	}
}

// AllOf combines multiple assertion Readers into a single Reader that passes
// only if all assertions pass.
//
// This function uses boolean AND logic (MonoidAll) to combine the results of
// all assertions. If any assertion fails, the combined assertion fails.
//
// This is useful for grouping related assertions together and ensuring all
// conditions are met.
//
// Parameters:
//   - readers: Array of assertion Readers to combine
//
// Returns:
//   - A single Reader that performs all assertions and returns true only if all pass
//
// Example:
//
//	func TestUser(t *testing.T) {
//	    user := User{Name: "Alice", Age: 30, Active: true}
//	    assertions := assert.AllOf([]assert.Reader{
//	        assert.Equal("Alice")(user.Name),
//	        assert.Equal(30)(user.Age),
//	        assert.Equal(true)(user.Active),
//	    })
//	    assertions(t)
//	}
//
//go:inline
func AllOf(readers []Reader) Reader {
	return reader.MonadReduceArrayM(readers, boolean.MonoidAll)
}

// RunAll executes a map of named test cases, running each as a subtest.
//
// This function creates a Reader that runs multiple named test cases using
// Go's t.Run for proper test isolation and reporting. Each test case is
// executed as a separate subtest with its own name.
//
// The function returns true only if all subtests pass. This allows for
// better test organization and clearer test output.
//
// Parameters:
//   - testcases: Map of test names to assertion Readers
//
// Returns:
//   - A Reader that executes all named test cases and returns true if all pass
//
// Example:
//
//	func TestMathOperations(t *testing.T) {
//	    testcases := map[string]assert.Reader{
//	        "addition":       assert.Equal(4)(2 + 2),
//	        "multiplication": assert.Equal(6)(2 * 3),
//	        "subtraction":    assert.Equal(1)(3 - 2),
//	    }
//	    assert.RunAll(testcases)(t)
//	}
//
//go:inline
func RunAll(testcases map[string]Reader) Reader {
	return func(t *testing.T) bool {
		current := true
		for k, r := range testcases {
			current = current && t.Run(k, func(t1 *testing.T) {
				r(t1)
			})
		}
		return current
	}
}
