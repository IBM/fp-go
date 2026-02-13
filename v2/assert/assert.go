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
// # Data Last Principle
//
// This package follows the "data last" functional programming principle, where
// the data being operated on comes as the last parameter in a chain of function
// applications. This design enables several powerful functional programming patterns:
//
//  1. **Partial Application**: You can create reusable assertion functions by providing
//     configuration parameters first, leaving the data and testing context for later.
//
//  2. **Function Composition**: Assertions can be composed and combined before being
//     applied to actual data.
//
//  3. **Point-Free Style**: You can pass assertion functions around without immediately
//     providing the data they operate on.
//
// The general pattern is:
//
//	assert.Function(config)(data)(testingContext)
//	               ↑        ↑     ↑
//	            expected  actual  *testing.T (always last)
//
// For single-parameter assertions:
//
//	assert.Function(data)(testingContext)
//	                ↑     ↑
//	              actual  *testing.T (always last)
//
// Examples of "data last" in action:
//
//	// Multi-parameter: expected value → actual value → testing context
//	assert.Equal(42)(result)(t)
//	assert.ArrayContains(3)(numbers)(t)
//
//	// Single-parameter: data → testing context
//	assert.NoError(err)(t)
//	assert.ArrayNotEmpty(arr)(t)
//
//	// Partial application - create reusable assertions
//	isPositive := assert.That(N.MoreThan(0))
//	// Later, apply to different values:
//	isPositive(42)(t)   // Passes
//	isPositive(-5)(t)   // Fails
//
//	// Composition - combine assertions before applying data
//	validateUser := func(u User) assert.Reader {
//	    return assert.AllOf([]assert.Reader{
//	        assert.Equal("Alice")(u.Name),
//	        assert.That(func(age int) bool { return age >= 18 })(u.Age),
//	    })
//	}
//	validateUser(user)(t)
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
	"github.com/IBM/fp-go/v2/option"
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

// NotEqual tests if the expected and the actual values are not equal.
//
// This function follows the "data last" principle - you provide the expected value first,
// then the actual value, and finally the testing.T context.
//
// Example:
//
//	func TestNotEqual(t *testing.T) {
//	    value := 42
//	    assert.NotEqual(10)(value)(t)  // Passes: 42 != 10
//	    assert.NotEqual(42)(value)(t)  // Fails: 42 == 42
//	}
func NotEqual[T any](expected T) Kleisli[T] {
	return wrap1(assert.NotEqual, expected)
}

// Equal tests if the expected and the actual values are equal.
//
// This is one of the most commonly used assertions. It follows the "data last" principle -
// you provide the expected value first, then the actual value, and finally the testing.T context.
//
// Example:
//
//	func TestEqual(t *testing.T) {
//	    result := 2 + 2
//	    assert.Equal(4)(result)(t)  // Passes
//
//	    name := "Alice"
//	    assert.Equal("Alice")(name)(t)  // Passes
//
//	    // Can be composed with other assertions
//	    user := User{Name: "Bob", Age: 30}
//	    assertions := assert.AllOf([]assert.Reader{
//	        assert.Equal("Bob")(user.Name),
//	        assert.Equal(30)(user.Age),
//	    })
//	    assertions(t)
//	}
func Equal[T any](expected T) Kleisli[T] {
	return wrap1(assert.Equal, expected)
}

// ArrayNotEmpty checks if an array is not empty.
//
// Example:
//
//	func TestArrayNotEmpty(t *testing.T) {
//	    numbers := []int{1, 2, 3}
//	    assert.ArrayNotEmpty(numbers)(t)  // Passes
//
//	    empty := []int{}
//	    assert.ArrayNotEmpty(empty)(t)  // Fails
//	}
func ArrayNotEmpty[T any](arr []T) Reader {
	return func(t *testing.T) bool {
		return assert.NotEmpty(t, arr)
	}
}

// ArrayEmpty checks if an array is empty.
//
// This is the complement of ArrayNotEmpty, asserting that a slice has no elements.
//
// Example:
//
//	func TestArrayEmpty(t *testing.T) {
//	    empty := []int{}
//	    assert.ArrayEmpty(empty)(t)  // Passes
//
//	    numbers := []int{1, 2, 3}
//	    assert.ArrayEmpty(numbers)(t)  // Fails
//	}
func ArrayEmpty[T any](arr []T) Reader {
	return func(t *testing.T) bool {
		return assert.Empty(t, arr)
	}
}

// RecordNotEmpty checks if a map is not empty.
//
// Example:
//
//	func TestRecordNotEmpty(t *testing.T) {
//	    config := map[string]int{"timeout": 30, "retries": 3}
//	    assert.RecordNotEmpty(config)(t)  // Passes
//
//	    empty := map[string]int{}
//	    assert.RecordNotEmpty(empty)(t)  // Fails
//	}
func RecordNotEmpty[K comparable, T any](mp map[K]T) Reader {
	return func(t *testing.T) bool {
		return assert.NotEmpty(t, mp)
	}
}

// RecordEmpty checks if a map is empty.
//
// This is the complement of RecordNotEmpty, asserting that a map has no key-value pairs.
//
// Example:
//
//	func TestRecordEmpty(t *testing.T) {
//	    empty := map[string]int{}
//	    assert.RecordEmpty(empty)(t)  // Passes
//
//	    config := map[string]int{"timeout": 30}
//	    assert.RecordEmpty(config)(t)  // Fails
//	}
func RecordEmpty[K comparable, T any](mp map[K]T) Reader {
	return func(t *testing.T) bool {
		return assert.Empty(t, mp)
	}
}

// StringNotEmpty checks if a string is not empty.
//
// Example:
//
//	func TestStringNotEmpty(t *testing.T) {
//	    message := "Hello, World!"
//	    assert.StringNotEmpty(message)(t)  // Passes
//
//	    empty := ""
//	    assert.StringNotEmpty(empty)(t)  // Fails
//	}
func StringNotEmpty(s string) Reader {
	return func(t *testing.T) bool {
		return assert.NotEmpty(t, s)
	}
}

// ArrayLength tests if an array has the expected length.
//
// Example:
//
//	func TestArrayLength(t *testing.T) {
//	    numbers := []int{1, 2, 3, 4, 5}
//	    assert.ArrayLength[int](5)(numbers)(t)  // Passes
//	    assert.ArrayLength[int](3)(numbers)(t)  // Fails
//	}
func ArrayLength[T any](expected int) Kleisli[[]T] {
	return func(actual []T) Reader {
		return func(t *testing.T) bool {
			return assert.Len(t, actual, expected)
		}
	}
}

// RecordLength tests if a map has the expected length.
//
// Example:
//
//	func TestRecordLength(t *testing.T) {
//	    config := map[string]string{"host": "localhost", "port": "8080"}
//	    assert.RecordLength[string, string](2)(config)(t)  // Passes
//	    assert.RecordLength[string, string](3)(config)(t)  // Fails
//	}
func RecordLength[K comparable, T any](expected int) Kleisli[map[K]T] {
	return func(actual map[K]T) Reader {
		return func(t *testing.T) bool {
			return assert.Len(t, actual, expected)
		}
	}
}

// StringLength tests if a string has the expected length.
//
// Example:
//
//	func TestStringLength(t *testing.T) {
//	    message := "Hello"
//	    assert.StringLength[any, any](5)(message)(t)  // Passes
//	    assert.StringLength[any, any](10)(message)(t)  // Fails
//	}
func StringLength[K comparable, T any](expected int) Kleisli[string] {
	return func(actual string) Reader {
		return func(t *testing.T) bool {
			return assert.Len(t, actual, expected)
		}
	}
}

// NoError validates that there is no error.
//
// This is commonly used to assert that operations complete successfully.
//
// Example:
//
//	func TestNoError(t *testing.T) {
//	    err := doSomething()
//	    assert.NoError(err)(t)  // Passes if err is nil
//
//	    // Can be used with result types
//	    result := result.TryCatch(func() (int, error) {
//	        return 42, nil
//	    })
//	    assert.Success(result)(t)  // Uses NoError internally
//	}
func NoError(err error) Reader {
	return func(t *testing.T) bool {
		return assert.NoError(t, err)
	}
}

// Error validates that there is an error.
//
// This is used to assert that operations fail as expected.
//
// Example:
//
//	func TestError(t *testing.T) {
//	    err := validateInput("")
//	    assert.Error(err)(t)  // Passes if err is not nil
//
//	    err2 := validateInput("valid")
//	    assert.Error(err2)(t)  // Fails if err2 is nil
//	}
func Error(err error) Reader {
	return func(t *testing.T) bool {
		return assert.Error(t, err)
	}
}

// Success checks if a [Result] represents success.
//
// This is a convenience function for testing Result types from the fp-go library.
//
// Example:
//
//	func TestSuccess(t *testing.T) {
//	    res := result.Of[int](42)
//	    assert.Success(res)(t)  // Passes
//
//	    failedRes := result.Error[int](errors.New("failed"))
//	    assert.Success(failedRes)(t)  // Fails
//	}
func Success[T any](res Result[T]) Reader {
	return NoError(result.ToError(res))
}

// Failure checks if a [Result] represents failure.
//
// This is a convenience function for testing Result types from the fp-go library.
//
// Example:
//
//	func TestFailure(t *testing.T) {
//	    res := result.Error[int](errors.New("something went wrong"))
//	    assert.Failure(res)(t)  // Passes
//
//	    successRes := result.Of[int](42)
//	    assert.Failure(successRes)(t)  // Fails
//	}
func Failure[T any](res Result[T]) Reader {
	return Error(result.ToError(res))
}

// ArrayContains tests if a value is contained in an array.
//
// Example:
//
//	func TestArrayContains(t *testing.T) {
//	    numbers := []int{1, 2, 3, 4, 5}
//	    assert.ArrayContains(3)(numbers)(t)  // Passes
//	    assert.ArrayContains(10)(numbers)(t)  // Fails
//
//	    names := []string{"Alice", "Bob", "Charlie"}
//	    assert.ArrayContains("Bob")(names)(t)  // Passes
//	}
func ArrayContains[T any](expected T) Kleisli[[]T] {
	return func(actual []T) Reader {
		return func(t *testing.T) bool {
			return assert.Contains(t, actual, expected)
		}
	}
}

// ContainsKey tests if a key is contained in a map.
//
// Example:
//
//	func TestContainsKey(t *testing.T) {
//	    config := map[string]int{"timeout": 30, "retries": 3}
//	    assert.ContainsKey[int]("timeout")(config)(t)  // Passes
//	    assert.ContainsKey[int]("maxSize")(config)(t)  // Fails
//	}
func ContainsKey[T any, K comparable](expected K) Kleisli[map[K]T] {
	return func(actual map[K]T) Reader {
		return func(t *testing.T) bool {
			return assert.Contains(t, actual, expected)
		}
	}
}

// NotContainsKey tests if a key is not contained in a map.
//
// Example:
//
//	func TestNotContainsKey(t *testing.T) {
//	    config := map[string]int{"timeout": 30, "retries": 3}
//	    assert.NotContainsKey[int]("maxSize")(config)(t)  // Passes
//	    assert.NotContainsKey[int]("timeout")(config)(t)  // Fails
//	}
func NotContainsKey[T any, K comparable](expected K) Kleisli[map[K]T] {
	return func(actual map[K]T) Reader {
		return func(t *testing.T) bool {
			return assert.NotContains(t, actual, expected)
		}
	}
}

// That asserts that a particular predicate matches.
//
// This is a powerful function that allows you to create custom assertions using predicates.
//
// Example:
//
//	func TestThat(t *testing.T) {
//	    // Test if a number is positive
//	    isPositive := N.MoreThan(0)
//	    assert.That(isPositive)(42)(t)  // Passes
//	    assert.That(isPositive)(-5)(t)  // Fails
//
//	    // Test if a string is uppercase
//	    isUppercase := func(s string) bool { return s == strings.ToUpper(s) }
//	    assert.That(isUppercase)("HELLO")(t)  // Passes
//	    assert.That(isUppercase)("Hello")(t)  // Fails
//
//	    // Can be combined with Local for property testing
//	    type User struct { Age int }
//	    ageIsAdult := assert.Local(func(u User) int { return u.Age })(
//	        assert.That(func(age int) bool { return age >= 18 }),
//	    )
//	    user := User{Age: 25}
//	    ageIsAdult(user)(t)  // Passes
//	}
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
	return SequenceRecord(testcases)
}

// Local transforms a Reader that works on type R1 into a Reader that works on type R2,
// by providing a function that converts R2 to R1. This allows you to focus a test on a
// specific property or subset of a larger data structure.
//
// See: https://github.com/fantasyland/fantasy-land?tab=readme-ov-file#profunctor
//
// This is particularly useful when you have an assertion that operates on a specific field
// or property, and you want to apply it to a complete object. Instead of extracting the
// property and then asserting on it, you can transform the assertion to work directly
// on the whole object.
//
// Parameters:
//   - f: A function that extracts or transforms R2 into R1
//
// Returns:
//   - A function that transforms a Reader[R1, Reader] into a Reader[R2, Reader]
//
// Example:
//
//	type User struct {
//	    Name string
//	    Age  int
//	}
//
//	// Create an assertion that checks if age is positive
//	ageIsPositive := assert.That(func(age int) bool { return age > 0 })
//
//	// Focus this assertion on the Age field of User
//	userAgeIsPositive := assert.Local(func(u User) int { return u.Age })(ageIsPositive)
//
//	// Now we can test the whole User object
//	user := User{Name: "Alice", Age: 30}
//	userAgeIsPositive(user)(t)
//
//go:inline
func Local[R1, R2 any](f func(R2) R1) func(Kleisli[R1]) Kleisli[R2] {
	return reader.Local[Reader](f)
}

// LocalL is similar to Local but uses a Lens to focus on a specific property.
// A Lens is a functional programming construct that provides a composable way to
// focus on a part of a data structure.
//
// This function is particularly useful when you want to focus a test on a specific
// field of a struct using a lens, making the code more declarative and composable.
// Lenses are often code-generated or predefined for common data structures.
//
// Parameters:
//   - l: A Lens that focuses from type S to type T
//
// Returns:
//   - A function that transforms a Reader[T, Reader] into a Reader[S, Reader]
//
// Example:
//
//	type Person struct {
//	    Name  string
//	    Email string
//	}
//
//	// Assume we have a lens that focuses on the Email field
//	var emailLens = lens.Prop[Person, string]("Email")
//
//	// Create an assertion for email format
//	validEmail := assert.That(func(email string) bool {
//	    return strings.Contains(email, "@")
//	})
//
//	// Focus this assertion on the Email property using a lens
//	validPersonEmail := assert.LocalL(emailLens)(validEmail)
//
//	// Test a Person object
//	person := Person{Name: "Bob", Email: "bob@example.com"}
//	validPersonEmail(person)(t)
//
//go:inline
func LocalL[S, T any](l Lens[S, T]) func(Kleisli[T]) Kleisli[S] {
	return reader.Local[Reader](l.Get)
}

// fromOptionalGetter is an internal helper that creates an assertion Reader from
// an optional getter function. It asserts that the optional value is present (Some).
func fromOptionalGetter[S, T any](getter func(S) option.Option[T], msgAndArgs ...any) Kleisli[S] {
	return func(s S) Reader {
		return func(t *testing.T) bool {
			return assert.True(t, option.IsSome(getter(s)), msgAndArgs...)
		}
	}
}

// FromOptional creates an assertion that checks if an Optional can successfully extract a value.
// An Optional is an optic that represents an optional reference to a subpart of a data structure.
//
// This function is useful when you have an Optional optic and want to assert that the optional
// value is present (Some) rather than absent (None). The assertion passes if the Optional's
// GetOption returns Some, and fails if it returns None.
//
// This enables property-focused testing where you verify that a particular optional field or
// sub-structure exists and is accessible.
//
// Parameters:
//   - opt: An Optional optic that focuses from type S to type T
//
// Returns:
//   - A Reader that asserts the optional value is present when applied to a value of type S
//
// Example:
//
//	type Config struct {
//	    Database *DatabaseConfig  // Optional field
//	}
//
//	type DatabaseConfig struct {
//	    Host string
//	    Port int
//	}
//
//	// Create an Optional that focuses on the Database field
//	dbOptional := optional.MakeOptional(
//	    func(c Config) option.Option[*DatabaseConfig] {
//	        if c.Database != nil {
//	            return option.Some(c.Database)
//	        }
//	        return option.None[*DatabaseConfig]()
//	    },
//	    func(c Config, db *DatabaseConfig) Config {
//	        c.Database = db
//	        return c
//	    },
//	)
//
//	// Assert that the database config is present
//	hasDatabaseConfig := assert.FromOptional(dbOptional)
//
//	config := Config{Database: &DatabaseConfig{Host: "localhost", Port: 5432}}
//	hasDatabaseConfig(config)(t)  // Passes
//
//	emptyConfig := Config{Database: nil}
//	hasDatabaseConfig(emptyConfig)(t)  // Fails
//
//go:inline
func FromOptional[S, T any](opt Optional[S, T]) reader.Reader[S, Reader] {
	return fromOptionalGetter(opt.GetOption, "Optional: %s", opt)
}

// FromPrism creates an assertion that checks if a Prism can successfully extract a value.
// A Prism is an optic used to select part of a sum type (tagged union or variant).
//
// This function is useful when you have a Prism optic and want to assert that a value
// matches a specific variant of a sum type. The assertion passes if the Prism's GetOption
// returns Some (meaning the value is of the expected variant), and fails if it returns None
// (meaning the value is a different variant).
//
// This enables variant-focused testing where you verify that a value is of a particular
// type or matches a specific condition within a sum type.
//
// Parameters:
//   - p: A Prism optic that focuses from type S to type T
//
// Returns:
//   - A Reader that asserts the prism successfully extracts when applied to a value of type S
//
// Example:
//
//	type Result interface{ isResult() }
//	type Success struct{ Value int }
//	type Failure struct{ Error string }
//
//	func (Success) isResult() {}
//	func (Failure) isResult() {}
//
//	// Create a Prism that focuses on Success variant
//	successPrism := prism.MakePrism(
//	    func(r Result) option.Option[int] {
//	        if s, ok := r.(Success); ok {
//	            return option.Some(s.Value)
//	        }
//	        return option.None[int]()
//	    },
//	    func(v int) Result { return Success{Value: v} },
//	)
//
//	// Assert that the result is a Success
//	isSuccess := assert.FromPrism(successPrism)
//
//	result1 := Success{Value: 42}
//	isSuccess(result1)(t)  // Passes
//
//	result2 := Failure{Error: "something went wrong"}
//	isSuccess(result2)(t)  // Fails
//
//go:inline
func FromPrism[S, T any](p Prism[S, T]) reader.Reader[S, Reader] {
	return fromOptionalGetter(p.GetOption, "Prism: %s", p)
}
