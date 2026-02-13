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

package assert

import (
	"iter"
	"testing"

	"github.com/IBM/fp-go/v2/context/readerio"
	"github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/optics/optional"
	"github.com/IBM/fp-go/v2/optics/prism"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/result"
)

type (
	// Result represents a computation that may fail with an error.
	//
	// This is an alias for [result.Result][T], which encapsulates either a successful
	// value of type T or an error. It's commonly used in test assertions to represent
	// operations that might fail, allowing for functional error handling without exceptions.
	//
	// A Result can be in one of two states:
	//   - Success: Contains a value of type T
	//   - Failure: Contains an error
	//
	// This type is particularly useful in testing scenarios where you need to:
	//   - Test functions that return results
	//   - Chain operations that might fail
	//   - Handle errors functionally
	//
	// Example:
	//
	//	func TestResultHandling(t *testing.T) {
	//	    successResult := result.Of[int](42)
	//	    assert.Success(successResult)(t)  // Passes
	//
	//	    failureResult := result.Error[int](errors.New("failed"))
	//	    assert.Failure(failureResult)(t)  // Passes
	//	}
	//
	// See also:
	//   - [Success]: Asserts a Result is successful
	//   - [Failure]: Asserts a Result contains an error
	//   - [result.Result]: The underlying Result type
	Result[T any] = result.Result[T]

	// Reader represents a test assertion that depends on a [testing.T] context and returns a boolean.
	//
	// This is the core type for all assertions in this package. It's an alias for
	// [reader.Reader][*testing.T, bool], which is a function that takes a testing context
	// and produces a boolean result indicating whether the assertion passed.
	//
	// The Reader pattern enables:
	//   - Composable assertions that can be combined using functional operators
	//   - Deferred execution - assertions are defined but not executed until applied to a test
	//   - Reusable assertion logic that can be applied to multiple tests
	//   - Functional composition of complex test conditions
	//
	// All assertion functions in this package return a Reader, which must be applied
	// to a *testing.T to execute the assertion:
	//
	//	assertion := assert.Equal(42)(result)  // Creates a Reader
	//	assertion(t)                            // Executes the assertion
	//
	// Readers can be composed using functions like [AllOf], [ApplicativeMonoid], or
	// functional operators from the reader package.
	//
	// Example:
	//
	//	func TestReaderComposition(t *testing.T) {
	//	    // Create individual assertions
	//	    assertion1 := assert.Equal(42)(42)
	//	    assertion2 := assert.StringNotEmpty("hello")
	//
	//	    // Combine them
	//	    combined := assert.AllOf([]assert.Reader{assertion1, assertion2})
	//
	//	    // Execute the combined assertion
	//	    combined(t)
	//	}
	//
	// See also:
	//   - [Kleisli]: Function that produces a Reader from a value
	//   - [AllOf]: Combines multiple Readers
	//   - [ApplicativeMonoid]: Monoid for combining Readers
	Reader = reader.Reader[*testing.T, bool]

	// Kleisli represents a function that produces a test assertion [Reader] from a value of type T.
	//
	// This is an alias for [reader.Reader][T, Reader], which is a function that takes a value
	// of type T and returns a Reader (test assertion). This pattern is fundamental to the
	// "data last" principle used throughout this package.
	//
	// Kleisli functions enable:
	//   - Partial application of assertions - configure the expected value first, apply actual value later
	//   - Reusable assertion builders that can be applied to different values
	//   - Functional composition of assertion pipelines
	//   - Point-free style programming with assertions
	//
	// Most assertion functions in this package return a Kleisli, which must be applied
	// to the actual value being tested, and then to a *testing.T:
	//
	//	kleisli := assert.Equal(42)  // Kleisli[int] - expects an int
	//	reader := kleisli(result)     // Reader - assertion ready to execute
	//	reader(t)                     // Execute the assertion
	//
	// Or more concisely:
	//
	//	assert.Equal(42)(result)(t)
	//
	// Example:
	//
	//	func TestKleisliPattern(t *testing.T) {
	//	    // Create a reusable assertion for positive numbers
	//	    isPositive := assert.That(func(n int) bool { return n > 0 })
	//
	//	    // Apply it to different values
	//	    isPositive(42)(t)   // Passes
	//	    isPositive(100)(t)  // Passes
	//	    // isPositive(-5)(t) would fail
	//
	//	    // Can be used with Local for property testing
	//	    type User struct { Age int }
	//	    checkAge := assert.Local(func(u User) int { return u.Age })(isPositive)
	//	    checkAge(User{Age: 25})(t)  // Passes
	//	}
	//
	// See also:
	//   - [Reader]: The assertion type produced by Kleisli
	//   - [Local]: Focuses a Kleisli on a property of a larger structure
	Kleisli[T any] = reader.Reader[T, Reader]

	// Predicate represents a function that tests a value of type T and returns a boolean.
	//
	// This is an alias for [predicate.Predicate][T], which is a simple function that
	// takes a value and returns true or false based on some condition. Predicates are
	// used with the [That] function to create custom assertions.
	//
	// Predicates enable:
	//   - Custom validation logic for any type
	//   - Reusable test conditions
	//   - Composition of complex validation rules
	//   - Integration with functional programming patterns
	//
	// Example:
	//
	//	func TestPredicates(t *testing.T) {
	//	    // Simple predicate
	//	    isEven := func(n int) bool { return n%2 == 0 }
	//	    assert.That(isEven)(42)(t)  // Passes
	//
	//	    // String predicate
	//	    hasPrefix := func(s string) bool { return strings.HasPrefix(s, "test") }
	//	    assert.That(hasPrefix)("test_file.go")(t)  // Passes
	//
	//	    // Complex predicate
	//	    isValidEmail := func(s string) bool {
	//	        return strings.Contains(s, "@") && strings.Contains(s, ".")
	//	    }
	//	    assert.That(isValidEmail)("user@example.com")(t)  // Passes
	//	}
	//
	// See also:
	//   - [That]: Creates an assertion from a Predicate
	//   - [predicate.Predicate]: The underlying predicate type
	Predicate[T any] = predicate.Predicate[T]

	// Lens is a functional reference to a subpart of a data structure.
	//
	// This is an alias for [lens.Lens][S, T], which provides a composable way to focus
	// on a specific field within a larger structure. Lenses enable getting and setting
	// values in nested data structures in a functional, immutable way.
	//
	// In the context of testing, lenses are used with [LocalL] to focus assertions
	// on specific properties of complex objects without manually extracting those properties.
	//
	// A Lens[S, T] focuses on a value of type T within a structure of type S.
	//
	// Example:
	//
	//	func TestLensUsage(t *testing.T) {
	//	    type Address struct { City string }
	//	    type User struct { Name string; Address Address }
	//
	//	    // Define lenses (typically generated)
	//	    addressLens := lens.Lens[User, Address]{...}
	//	    cityLens := lens.Lens[Address, string]{...}
	//
	//	    // Compose lenses to focus on nested field
	//	    userCityLens := lens.Compose(addressLens, cityLens)
	//
	//	    // Use with LocalL to assert on nested property
	//	    user := User{Name: "Alice", Address: Address{City: "NYC"}}
	//	    assert.LocalL(userCityLens)(assert.Equal("NYC"))(user)(t)
	//	}
	//
	// See also:
	//   - [LocalL]: Uses a Lens to focus assertions on a property
	//   - [lens.Lens]: The underlying lens type
	//   - [Optional]: Similar but for values that may not exist
	Lens[S, T any] = lens.Lens[S, T]

	// Optional is an optic that focuses on a value that may or may not be present.
	//
	// This is an alias for [optional.Optional][S, T], which is similar to a [Lens] but
	// handles cases where the focused value might not exist. Optionals are useful for
	// working with nullable fields, optional properties, or values that might be absent.
	//
	// In testing, Optionals are used with [FromOptional] to create assertions that
	// verify whether an optional value is present and, if so, whether it satisfies
	// certain conditions.
	//
	// An Optional[S, T] focuses on an optional value of type T within a structure of type S.
	//
	// Example:
	//
	//	func TestOptionalUsage(t *testing.T) {
	//	    type Config struct { Timeout *int }
	//
	//	    // Define optional (typically generated)
	//	    timeoutOptional := optional.Optional[Config, int]{...}
	//
	//	    // Test when value is present
	//	    config1 := Config{Timeout: ptr(30)}
	//	    assert.FromOptional(timeoutOptional)(
	//	        assert.Equal(30),
	//	    )(config1)(t)  // Passes
	//
	//	    // Test when value is absent
	//	    config2 := Config{Timeout: nil}
	//	    // FromOptional would fail because value is not present
	//	}
	//
	// See also:
	//   - [FromOptional]: Creates assertions for optional values
	//   - [optional.Optional]: The underlying optional type
	//   - [Lens]: Similar but for values that always exist
	Optional[S, T any] = optional.Optional[S, T]

	// Prism is an optic that focuses on a case of a sum type.
	//
	// This is an alias for [prism.Prism][S, T], which provides a way to focus on one
	// variant of a sum type (like Result, Option, Either, etc.). Prisms enable pattern
	// matching and extraction of values from sum types in a functional way.
	//
	// In testing, Prisms are used with [FromPrism] to create assertions that verify
	// whether a value matches a specific case and, if so, whether the contained value
	// satisfies certain conditions.
	//
	// A Prism[S, T] focuses on a value of type T that may be contained within a sum type S.
	//
	// Example:
	//
	//	func TestPrismUsage(t *testing.T) {
	//	    // Prism for extracting success value from Result
	//	    successPrism := prism.Success[int]()
	//
	//	    // Test successful result
	//	    successResult := result.Of[int](42)
	//	    assert.FromPrism(successPrism)(
	//	        assert.Equal(42),
	//	    )(successResult)(t)  // Passes
	//
	//	    // Prism for extracting error from Result
	//	    failurePrism := prism.Failure[int]()
	//
	//	    // Test failed result
	//	    failureResult := result.Error[int](errors.New("failed"))
	//	    assert.FromPrism(failurePrism)(
	//	        assert.Error,
	//	    )(failureResult)(t)  // Passes
	//	}
	//
	// See also:
	//   - [FromPrism]: Creates assertions for prism-focused values
	//   - [prism.Prism]: The underlying prism type
	//   - [Optional]: Similar but for optional values
	Prism[S, T any] = prism.Prism[S, T]

	// ReaderIOResult represents a context-aware, IO-based computation that may fail.
	//
	// This is an alias for [readerioresult.ReaderIOResult][A], which combines three
	// computational effects:
	//   - Reader: Depends on a context (like context.Context)
	//   - IO: Performs side effects (like file I/O, network calls)
	//   - Result: May fail with an error
	//
	// In testing, ReaderIOResult is used with [FromReaderIOResult] to convert
	// context-aware, effectful computations into test assertions. This is useful
	// when your test assertions need to:
	//   - Access a context for cancellation or deadlines
	//   - Perform IO operations (database queries, API calls, file access)
	//   - Handle potential errors gracefully
	//
	// Example:
	//
	//	func TestReaderIOResult(t *testing.T) {
	//	    // Create a ReaderIOResult that performs IO and may fail
	//	    checkDatabase := func(ctx context.Context) func() result.Result[assert.Reader] {
	//	        return func() result.Result[assert.Reader] {
	//	            // Perform database check with context
	//	            if err := db.PingContext(ctx); err != nil {
	//	                return result.Error[assert.Reader](err)
	//	            }
	//	            return result.Of[assert.Reader](assert.NoError(nil))
	//	        }
	//	    }
	//
	//	    // Convert to Reader and execute
	//	    assertion := assert.FromReaderIOResult(checkDatabase)
	//	    assertion(t)
	//	}
	//
	// See also:
	//   - [FromReaderIOResult]: Converts ReaderIOResult to Reader
	//   - [ReaderIO]: Similar but without error handling
	//   - [readerioresult.ReaderIOResult]: The underlying type
	ReaderIOResult[A any] = readerioresult.ReaderIOResult[A]

	// ReaderIO represents a context-aware, IO-based computation.
	//
	// This is an alias for [readerio.ReaderIO][A], which combines two computational effects:
	//   - Reader: Depends on a context (like context.Context)
	//   - IO: Performs side effects (like logging, metrics)
	//
	// In testing, ReaderIO is used with [FromReaderIO] to convert context-aware,
	// effectful computations into test assertions. This is useful when your test
	// assertions need to:
	//   - Access a context for cancellation or deadlines
	//   - Perform IO operations that don't fail (or handle failures internally)
	//   - Integrate with context-aware utilities
	//
	// Example:
	//
	//	func TestReaderIO(t *testing.T) {
	//	    // Create a ReaderIO that performs IO
	//	    logAndCheck := func(ctx context.Context) func() assert.Reader {
	//	        return func() assert.Reader {
	//	            // Log with context
	//	            logger.InfoContext(ctx, "Running test")
	//	            // Return assertion
	//	            return assert.Equal(42)(computeValue())
	//	        }
	//	    }
	//
	//	    // Convert to Reader and execute
	//	    assertion := assert.FromReaderIO(logAndCheck)
	//	    assertion(t)
	//	}
	//
	// See also:
	//   - [FromReaderIO]: Converts ReaderIO to Reader
	//   - [ReaderIOResult]: Similar but with error handling
	//   - [readerio.ReaderIO]: The underlying type
	ReaderIO[A any] = readerio.ReaderIO[A]

	// Seq2 represents a Go iterator that yields key-value pairs.
	//
	// This is an alias for [iter.Seq2][K, A], which is Go's standard iterator type
	// introduced in Go 1.23. It represents a sequence of key-value pairs that can be
	// iterated over using a for-range loop.
	//
	// In testing, Seq2 is used with [SequenceSeq2] to execute a sequence of named
	// test cases provided as an iterator. This enables:
	//   - Lazy evaluation of test cases
	//   - Memory-efficient testing of large test suites
	//   - Integration with Go's iterator patterns
	//   - Dynamic generation of test cases
	//
	// Example:
	//
	//	func TestSeq2Usage(t *testing.T) {
	//	    // Create an iterator of test cases
	//	    testCases := func(yield func(string, assert.Reader) bool) {
	//	        if !yield("test_addition", assert.Equal(4)(2+2)) {
	//	            return
	//	        }
	//	        if !yield("test_multiplication", assert.Equal(6)(2*3)) {
	//	            return
	//	        }
	//	    }
	//
	//	    // Execute all test cases
	//	    assert.SequenceSeq2[assert.Reader](testCases)(t)
	//	}
	//
	// See also:
	//   - [SequenceSeq2]: Executes a Seq2 of test cases
	//   - [TraverseArray]: Similar but for arrays
	//   - [iter.Seq2]: The underlying iterator type
	Seq2[K, A any] = iter.Seq2[K, A]

	// Pair represents a tuple of two values with potentially different types.
	//
	// This is an alias for [pair.Pair][L, R], which holds two values: a "head" (or "left")
	// of type L and a "tail" (or "right") of type R. Pairs are useful for grouping
	// related values together without defining a custom struct.
	//
	// In testing, Pairs are used with [TraverseArray] to associate test names with
	// their corresponding assertions. Each element in the array is transformed into
	// a Pair[string, Reader] where the string is the test name and the Reader is
	// the assertion to execute.
	//
	// Example:
	//
	//	func TestPairUsage(t *testing.T) {
	//	    type TestCase struct {
	//	        Input    int
	//	        Expected int
	//	    }
	//
	//	    testCases := []TestCase{
	//	        {Input: 2, Expected: 4},
	//	        {Input: 3, Expected: 9},
	//	    }
	//
	//	    // Transform each test case into a named assertion
	//	    traverse := assert.TraverseArray(func(tc TestCase) assert.Pair[string, assert.Reader] {
	//	        name := fmt.Sprintf("square(%d)=%d", tc.Input, tc.Expected)
	//	        assertion := assert.Equal(tc.Expected)(tc.Input * tc.Input)
	//	        return pair.MakePair(name, assertion)
	//	    })
	//
	//	    traverse(testCases)(t)
	//	}
	//
	// See also:
	//   - [TraverseArray]: Uses Pairs to create named test cases
	//   - [pair.Pair]: The underlying pair type
	//   - [pair.MakePair]: Creates a Pair
	//   - [pair.Head]: Extracts the first value
	//   - [pair.Tail]: Extracts the second value
	Pair[L, R any] = pair.Pair[L, R]

	// Void represents the absence of a meaningful value, similar to unit type in functional programming.
	//
	// This is an alias for [function.Void], which is used to represent operations that don't
	// return a meaningful value but may perform side effects. In the context of testing, Void
	// is used with IO operations that perform actions without producing a result.
	//
	// Void is conceptually similar to:
	//   - Unit type in functional languages (Haskell's (), Scala's Unit)
	//   - void in languages like C/Java (but as a value, not just a type)
	//   - Empty struct{} in Go (but with clearer semantic meaning)
	//
	// Example:
	//
	//	func TestWithSideEffect(t *testing.T) {
	//	    // An IO operation that logs but returns Void
	//	    logOperation := func() function.Void {
	//	        log.Println("Test executed")
	//	        return function.Void{}
	//	    }
	//
	//	    // Execute the operation
	//	    logOperation()
	//	}
	//
	// See also:
	//   - [IO]: Wraps side-effecting operations
	//   - [function.Void]: The underlying void type
	Void = function.Void

	// IO represents a side-effecting computation that produces a value of type A.
	//
	// This is an alias for [io.IO][A], which encapsulates operations that perform side effects
	// (like I/O operations, logging, or state mutations) and return a value. IO is a lazy
	// computation - it describes an effect but doesn't execute it until explicitly run.
	//
	// In testing, IO is used to:
	//   - Defer execution of side effects until needed
	//   - Compose multiple side-effecting operations
	//   - Maintain referential transparency in test setup
	//   - Separate effect description from effect execution
	//
	// An IO[A] is essentially a function `func() A` that:
	//   - Encapsulates a side effect
	//   - Returns a value of type A when executed
	//   - Can be composed with other IO operations
	//
	// Example:
	//
	//	func TestIOOperation(t *testing.T) {
	//	    // Define an IO operation that reads a file
	//	    readConfig := func() io.IO[string] {
	//	        return func() string {
	//	            data, _ := os.ReadFile("config.txt")
	//	            return string(data)
	//	        }
	//	    }
	//
	//	    // The IO is not executed yet - it's just a description
	//	    configIO := readConfig()
	//
	//	    // Execute the IO to get the result
	//	    config := configIO()
	//	    assert.StringNotEmpty(config)(t)
	//	}
	//
	// Example with composition:
	//
	//	func TestIOComposition(t *testing.T) {
	//	    // Chain multiple IO operations
	//	    pipeline := io.Map(
	//	        func(s string) int { return len(s) },
	//	    )(readFileIO)
	//
	//	    // Execute the composed operation
	//	    length := pipeline()
	//	    assert.That(func(n int) bool { return n > 0 })(length)(t)
	//	}
	//
	// See also:
	//   - [ReaderIO]: Combines Reader and IO effects
	//   - [ReaderIOResult]: Adds error handling to ReaderIO
	//   - [io.IO]: The underlying IO type
	//   - [Void]: Represents operations without meaningful return values
	IO[A any] = io.IO[A]
)
