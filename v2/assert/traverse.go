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
	"testing"

	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/reader"
)

// TraverseArray transforms an array of values into a test suite by applying a function
// that generates named test cases for each element.
//
// This function enables data-driven testing where you have a collection of test inputs
// and want to run a named subtest for each one. It follows the functional programming
// pattern of "traverse" - transforming a collection while preserving structure and
// accumulating effects (in this case, test execution).
//
// The function takes each element of the array, applies the provided function to generate
// a [Pair] of (test name, test assertion), and runs each as a separate subtest using
// Go's t.Run. All subtests must pass for the overall test to pass.
//
// # Parameters
//
//   - f: A function that takes a value of type T and returns a [Pair] containing:
//   - Head: The test name (string) for the subtest
//   - Tail: The test assertion ([Reader]) to execute
//
// # Returns
//
//   - A [Kleisli] function that takes an array of T and returns a [Reader] that:
//   - Executes each element as a named subtest
//   - Returns true only if all subtests pass
//   - Provides proper test isolation and reporting via t.Run
//
// # Use Cases
//
//   - Data-driven testing with multiple test cases
//   - Parameterized tests where each parameter gets its own subtest
//   - Testing collections where each element needs validation
//   - Property-based testing with generated test data
//
// # Example - Basic Data-Driven Testing
//
//	func TestMathOperations(t *testing.T) {
//	    type TestCase struct {
//	        Input    int
//	        Expected int
//	    }
//
//	    testCases := []TestCase{
//	        {Input: 2, Expected: 4},
//	        {Input: 3, Expected: 9},
//	        {Input: 4, Expected: 16},
//	    }
//
//	    square := func(n int) int { return n * n }
//
//	    traverse := assert.TraverseArray(func(tc TestCase) assert.Pair[string, assert.Reader] {
//	        name := fmt.Sprintf("square(%d)=%d", tc.Input, tc.Expected)
//	        assertion := assert.Equal(tc.Expected)(square(tc.Input))
//	        return pair.MakePair(name, assertion)
//	    })
//
//	    traverse(testCases)(t)
//	}
//
// # Example - String Validation
//
//	func TestStringValidation(t *testing.T) {
//	    inputs := []string{"hello", "world", "test"}
//
//	    traverse := assert.TraverseArray(func(s string) assert.Pair[string, assert.Reader] {
//	        return pair.MakePair(
//	            fmt.Sprintf("validate_%s", s),
//	            assert.AllOf([]assert.Reader{
//	                assert.StringNotEmpty(s),
//	                assert.That(func(str string) bool { return len(str) > 0 })(s),
//	            }),
//	        )
//	    })
//
//	    traverse(inputs)(t)
//	}
//
// # Example - Complex Object Testing
//
//	func TestUsers(t *testing.T) {
//	    type User struct {
//	        Name  string
//	        Age   int
//	        Email string
//	    }
//
//	    users := []User{
//	        {Name: "Alice", Age: 30, Email: "alice@example.com"},
//	        {Name: "Bob", Age: 25, Email: "bob@example.com"},
//	    }
//
//	    traverse := assert.TraverseArray(func(u User) assert.Pair[string, assert.Reader] {
//	        return pair.MakePair(
//	            fmt.Sprintf("user_%s", u.Name),
//	            assert.AllOf([]assert.Reader{
//	                assert.StringNotEmpty(u.Name),
//	                assert.That(func(age int) bool { return age > 0 })(u.Age),
//	                assert.That(func(email string) bool {
//	                    return len(email) > 0 && strings.Contains(email, "@")
//	                })(u.Email),
//	            }),
//	        )
//	    })
//
//	    traverse(users)(t)
//	}
//
// # Comparison with RunAll
//
// TraverseArray and [RunAll] serve similar purposes but differ in their approach:
//
//   - TraverseArray: Generates test cases from an array of data
//
//   - Input: Array of values + function to generate test cases
//
//   - Use when: You have test data and need to generate test cases from it
//
//   - RunAll: Executes pre-defined named test cases
//
//   - Input: Map of test names to assertions
//
//   - Use when: You have already defined test cases with names
//
// # Related Functions
//
//   - [SequenceSeq2]: Similar but works with Go iterators (Seq2) instead of arrays
//   - [RunAll]: Executes a map of named test cases
//   - [AllOf]: Combines multiple assertions without subtests
//
// # References
//
//   - Haskell traverse: https://hackage.haskell.org/package/base/docs/Data-Traversable.html#v:traverse
//   - Go subtests: https://go.dev/blog/subtests
func TraverseArray[T any](f func(T) Pair[string, Reader]) Kleisli[[]T] {
	return func(ts []T) Reader {
		return func(t *testing.T) bool {
			ok := true
			for _, src := range ts {
				test := f(src)
				res := t.Run(pair.Head(test), func(t *testing.T) {
					pair.Tail(test)(t)
				})
				ok = ok && res
			}
			return ok
		}
	}
}

// SequenceSeq2 executes a sequence of named test cases provided as a Go iterator.
//
// This function takes a [Seq2] iterator that yields (name, assertion) pairs and
// executes each as a separate subtest using Go's t.Run. It's similar to [TraverseArray]
// but works directly with Go's iterator protocol (introduced in Go 1.23) rather than
// requiring an array.
//
// The function iterates through all test cases, running each as a named subtest.
// All subtests must pass for the overall test to pass. This provides proper test
// isolation and clear reporting of which specific test cases fail.
//
// # Parameters
//
//   - s: A [Seq2] iterator that yields pairs of:
//   - Key: Test name (string) for the subtest
//   - Value: Test assertion ([Reader]) to execute
//
// # Returns
//
//   - A [Reader] that:
//   - Executes each test case as a named subtest
//   - Returns true only if all subtests pass
//   - Provides proper test isolation via t.Run
//
// # Use Cases
//
//   - Working with iterator-based test data
//   - Lazy evaluation of test cases
//   - Integration with Go 1.23+ iterator patterns
//   - Memory-efficient testing of large test suites
//
// # Example - Basic Usage with Iterator
//
//	func TestWithIterator(t *testing.T) {
//	    // Create an iterator of test cases
//	    testCases := func(yield func(string, assert.Reader) bool) {
//	        if !yield("test_addition", assert.Equal(4)(2+2)) {
//	            return
//	        }
//	        if !yield("test_subtraction", assert.Equal(1)(3-2)) {
//	            return
//	        }
//	        if !yield("test_multiplication", assert.Equal(6)(2*3)) {
//	            return
//	        }
//	    }
//
//	    assert.SequenceSeq2(testCases)(t)
//	}
//
// # Example - Generated Test Cases
//
//	func TestGeneratedCases(t *testing.T) {
//	    // Generate test cases on the fly
//	    generateTests := func(yield func(string, assert.Reader) bool) {
//	        for i := 1; i <= 5; i++ {
//	            name := fmt.Sprintf("test_%d", i)
//	            assertion := assert.Equal(i*i)(i * i)
//	            if !yield(name, assertion) {
//	                return
//	            }
//	        }
//	    }
//
//	    assert.SequenceSeq2(generateTests)(t)
//	}
//
// # Example - Filtering Test Cases
//
//	func TestFilteredCases(t *testing.T) {
//	    type TestCase struct {
//	        Name     string
//	        Input    int
//	        Expected int
//	        Skip     bool
//	    }
//
//	    allCases := []TestCase{
//	        {Name: "test1", Input: 2, Expected: 4, Skip: false},
//	        {Name: "test2", Input: 3, Expected: 9, Skip: true},
//	        {Name: "test3", Input: 4, Expected: 16, Skip: false},
//	    }
//
//	    // Create iterator that filters out skipped tests
//	    activeTests := func(yield func(string, assert.Reader) bool) {
//	        for _, tc := range allCases {
//	            if !tc.Skip {
//	                assertion := assert.Equal(tc.Expected)(tc.Input * tc.Input)
//	                if !yield(tc.Name, assertion) {
//	                    return
//	                }
//	            }
//	        }
//	    }
//
//	    assert.SequenceSeq2(activeTests)(t)
//	}
//
// # Comparison with TraverseArray
//
// SequenceSeq2 and [TraverseArray] serve similar purposes but differ in their input:
//
//   - SequenceSeq2: Works with iterators (Seq2)
//
//   - Input: Iterator yielding (name, assertion) pairs
//
//   - Use when: Working with Go 1.23+ iterators or lazy evaluation
//
//   - Memory: More efficient for large test suites (lazy evaluation)
//
//   - TraverseArray: Works with arrays
//
//   - Input: Array of values + transformation function
//
//   - Use when: You have an array of test data
//
//   - Memory: All test data must be in memory
//
// # Comparison with RunAll
//
// SequenceSeq2 and [RunAll] are very similar:
//
//   - SequenceSeq2: Takes an iterator (Seq2)
//   - RunAll: Takes a map[string]Reader
//
// Both execute named test cases as subtests. Choose based on your data structure:
// use SequenceSeq2 for iterators, RunAll for maps.
//
// # Related Functions
//
//   - [TraverseArray]: Similar but works with arrays instead of iterators
//   - [RunAll]: Executes a map of named test cases
//   - [AllOf]: Combines multiple assertions without subtests
//
// # References
//
//   - Go iterators: https://go.dev/blog/range-functions
//   - Go subtests: https://go.dev/blog/subtests
//   - Haskell sequence: https://hackage.haskell.org/package/base/docs/Data-Traversable.html#v:sequence
func SequenceSeq2[T any](s Seq2[string, Reader]) Reader {
	return func(t *testing.T) bool {
		ok := true
		for name, test := range s {
			res := t.Run(name, func(t *testing.T) {
				test(t)
			})
			ok = ok && res
		}
		return ok
	}
}

// TraverseRecord transforms a map of values into a test suite by applying a function
// that generates test assertions for each map entry.
//
// This function enables data-driven testing where you have a map of test data and want
// to run a named subtest for each entry. The map keys become test names, and the function
// transforms each value into a test assertion. It follows the functional programming
// pattern of "traverse" - transforming a collection while preserving structure and
// accumulating effects (in this case, test execution).
//
// The function takes each key-value pair from the map, applies the provided function to
// generate a [Reader] assertion, and runs each as a separate subtest using Go's t.Run.
// All subtests must pass for the overall test to pass.
//
// # Parameters
//
//   - f: A [Kleisli] function that takes a value of type T and returns a [Reader] assertion
//
// # Returns
//
//   - A [Kleisli] function that takes a map[string]T and returns a [Reader] that:
//   - Executes each map entry as a named subtest (using the key as the test name)
//   - Returns true only if all subtests pass
//   - Provides proper test isolation and reporting via t.Run
//
// # Use Cases
//
//   - Data-driven testing with named test cases in a map
//   - Testing configuration maps where keys are meaningful names
//   - Validating collections where natural keys exist
//   - Property-based testing with named scenarios
//
// # Example - Basic Configuration Testing
//
//	func TestConfigurations(t *testing.T) {
//	    configs := map[string]int{
//	        "timeout":     30,
//	        "maxRetries":  3,
//	        "bufferSize":  1024,
//	    }
//
//	    validatePositive := assert.That(func(n int) bool { return n > 0 })
//
//	    traverse := assert.TraverseRecord(validatePositive)
//	    traverse(configs)(t)
//	}
//
// # Example - User Validation
//
//	func TestUserMap(t *testing.T) {
//	    type User struct {
//	        Name string
//	        Age  int
//	    }
//
//	    users := map[string]User{
//	        "alice": {Name: "Alice", Age: 30},
//	        "bob":   {Name: "Bob", Age: 25},
//	        "carol": {Name: "Carol", Age: 35},
//	    }
//
//	    validateUser := func(u User) assert.Reader {
//	        return assert.AllOf([]assert.Reader{
//	            assert.StringNotEmpty(u.Name),
//	            assert.That(func(age int) bool { return age > 0 && age < 150 })(u.Age),
//	        })
//	    }
//
//	    traverse := assert.TraverseRecord(validateUser)
//	    traverse(users)(t)
//	}
//
// # Example - API Endpoint Testing
//
//	func TestEndpoints(t *testing.T) {
//	    type Endpoint struct {
//	        Path   string
//	        Method string
//	    }
//
//	    endpoints := map[string]Endpoint{
//	        "get_users":    {Path: "/api/users", Method: "GET"},
//	        "create_user":  {Path: "/api/users", Method: "POST"},
//	        "delete_user":  {Path: "/api/users/:id", Method: "DELETE"},
//	    }
//
//	    validateEndpoint := func(e Endpoint) assert.Reader {
//	        return assert.AllOf([]assert.Reader{
//	            assert.StringNotEmpty(e.Path),
//	            assert.That(func(path string) bool {
//	                return strings.HasPrefix(path, "/api/")
//	            })(e.Path),
//	            assert.That(func(method string) bool {
//	                return method == "GET" || method == "POST" ||
//	                       method == "PUT" || method == "DELETE"
//	            })(e.Method),
//	        })
//	    }
//
//	    traverse := assert.TraverseRecord(validateEndpoint)
//	    traverse(endpoints)(t)
//	}
//
// # Comparison with TraverseArray
//
// TraverseRecord and [TraverseArray] serve similar purposes but differ in their input:
//
//   - TraverseRecord: Works with maps (records)
//
//   - Input: Map with string keys + transformation function
//
//   - Use when: You have named test data in a map
//
//   - Test names: Derived from map keys
//
//   - TraverseArray: Works with arrays
//
//   - Input: Array of values + function that generates names and assertions
//
//   - Use when: You have sequential test data
//
//   - Test names: Generated by the transformation function
//
// # Comparison with SequenceRecord
//
// TraverseRecord and [SequenceRecord] are closely related:
//
//   - TraverseRecord: Transforms values into assertions
//
//   - Input: map[string]T + function T -> Reader
//
//   - Use when: You need to transform data before asserting
//
//   - SequenceRecord: Executes pre-defined assertions
//
//   - Input: map[string]Reader
//
//   - Use when: Assertions are already defined
//
// # Related Functions
//
//   - [SequenceRecord]: Similar but takes pre-defined assertions
//   - [TraverseArray]: Similar but works with arrays
//   - [RunAll]: Alias for SequenceRecord
//
// # References
//
//   - Haskell traverse: https://hackage.haskell.org/package/base/docs/Data-Traversable.html#v:traverse
//   - Go subtests: https://go.dev/blog/subtests
func TraverseRecord[T any](f Kleisli[T]) Kleisli[map[string]T] {
	return func(m map[string]T) Reader {
		return func(t *testing.T) bool {
			ok := true
			for name, src := range m {
				res := t.Run(name, func(t *testing.T) {
					f(src)(t)
				})
				ok = ok && res
			}
			return ok
		}
	}
}

// SequenceRecord executes a map of named test cases as subtests.
//
// This function takes a map where keys are test names and values are test assertions
// ([Reader]), and executes each as a separate subtest using Go's t.Run. It's the
// record (map) equivalent of [SequenceSeq2] and is actually aliased as [RunAll] for
// convenience.
//
// The function iterates through all map entries, running each as a named subtest.
// All subtests must pass for the overall test to pass. This provides proper test
// isolation and clear reporting of which specific test cases fail.
//
// # Parameters
//
//   - m: A map[string]Reader where:
//   - Keys: Test names (strings) for the subtests
//   - Values: Test assertions ([Reader]) to execute
//
// # Returns
//
//   - A [Reader] that:
//   - Executes each map entry as a named subtest
//   - Returns true only if all subtests pass
//   - Provides proper test isolation via t.Run
//
// # Use Cases
//
//   - Executing a collection of pre-defined named test cases
//   - Organizing related tests in a map structure
//   - Running multiple assertions with descriptive names
//   - Building test suites programmatically
//
// # Example - Basic Named Tests
//
//	func TestMathOperations(t *testing.T) {
//	    tests := map[string]assert.Reader{
//	        "addition":       assert.Equal(4)(2 + 2),
//	        "subtraction":    assert.Equal(1)(3 - 2),
//	        "multiplication": assert.Equal(6)(2 * 3),
//	        "division":       assert.Equal(2)(6 / 3),
//	    }
//
//	    assert.SequenceRecord(tests)(t)
//	}
//
// # Example - String Validation Suite
//
//	func TestStringValidations(t *testing.T) {
//	    testString := "hello world"
//
//	    tests := map[string]assert.Reader{
//	        "not_empty":      assert.StringNotEmpty(testString),
//	        "correct_length": assert.StringLength[any, any](11)(testString),
//	        "has_space":      assert.That(func(s string) bool {
//	            return strings.Contains(s, " ")
//	        })(testString),
//	        "lowercase":      assert.That(func(s string) bool {
//	            return s == strings.ToLower(s)
//	        })(testString),
//	    }
//
//	    assert.SequenceRecord(tests)(t)
//	}
//
// # Example - Complex Object Validation
//
//	func TestUserValidation(t *testing.T) {
//	    type User struct {
//	        Name  string
//	        Age   int
//	        Email string
//	    }
//
//	    user := User{Name: "Alice", Age: 30, Email: "alice@example.com"}
//
//	    tests := map[string]assert.Reader{
//	        "name_not_empty": assert.StringNotEmpty(user.Name),
//	        "age_positive":   assert.That(func(age int) bool { return age > 0 })(user.Age),
//	        "age_reasonable": assert.That(func(age int) bool { return age < 150 })(user.Age),
//	        "email_valid":    assert.That(func(email string) bool {
//	            return strings.Contains(email, "@") && strings.Contains(email, ".")
//	        })(user.Email),
//	    }
//
//	    assert.SequenceRecord(tests)(t)
//	}
//
// # Example - Array Validation Suite
//
//	func TestArrayValidations(t *testing.T) {
//	    numbers := []int{1, 2, 3, 4, 5}
//
//	    tests := map[string]assert.Reader{
//	        "not_empty":      assert.ArrayNotEmpty(numbers),
//	        "correct_length": assert.ArrayLength[int](5)(numbers),
//	        "contains_three": assert.ArrayContains(3)(numbers),
//	        "all_positive":   assert.That(func(arr []int) bool {
//	            for _, n := range arr {
//	                if n <= 0 {
//	                    return false
//	                }
//	            }
//	            return true
//	        })(numbers),
//	    }
//
//	    assert.SequenceRecord(tests)(t)
//	}
//
// # Comparison with TraverseRecord
//
// SequenceRecord and [TraverseRecord] are closely related:
//
//   - SequenceRecord: Executes pre-defined assertions
//
//   - Input: map[string]Reader (assertions already created)
//
//   - Use when: You have already defined test cases with assertions
//
//   - TraverseRecord: Transforms values into assertions
//
//   - Input: map[string]T + function T -> Reader
//
//   - Use when: You need to transform data before asserting
//
// # Comparison with SequenceSeq2
//
// SequenceRecord and [SequenceSeq2] serve similar purposes but differ in their input:
//
//   - SequenceRecord: Works with maps
//
//   - Input: map[string]Reader
//
//   - Use when: You have named test cases in a map
//
//   - Iteration order: Non-deterministic (map iteration)
//
//   - SequenceSeq2: Works with iterators
//
//   - Input: Seq2[string, Reader]
//
//   - Use when: You have test cases in an iterator
//
//   - Iteration order: Deterministic (iterator order)
//
// # Note on Map Iteration Order
//
// Go maps have non-deterministic iteration order. If test execution order matters,
// consider using [SequenceSeq2] with an iterator that provides deterministic ordering,
// or use [TraverseArray] with a slice of test cases.
//
// # Related Functions
//
//   - [RunAll]: Alias for SequenceRecord
//   - [TraverseRecord]: Similar but transforms values into assertions
//   - [SequenceSeq2]: Similar but works with iterators
//   - [TraverseArray]: Similar but works with arrays
//
// # References
//
//   - Go subtests: https://go.dev/blog/subtests
//   - Haskell sequence: https://hackage.haskell.org/package/base/docs/Data-Traversable.html#v:sequence
func SequenceRecord(m map[string]Reader) Reader {
	return TraverseRecord(reader.Ask[Reader]())(m)
}
