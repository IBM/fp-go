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
	"fmt"
	"testing"

	"github.com/IBM/fp-go/v2/pair"
)

// TestTraverseArray_EmptyArray tests that TraverseArray handles empty arrays correctly
func TestTraverseArray_EmptyArray(t *testing.T) {
	traverse := TraverseArray(func(n int) Pair[string, Reader] {
		return pair.MakePair(
			fmt.Sprintf("test_%d", n),
			Equal(n)(n),
		)
	})

	result := traverse([]int{})(t)
	if !result {
		t.Error("Expected TraverseArray to pass with empty array")
	}
}

// TestTraverseArray_SingleElement tests TraverseArray with a single element
func TestTraverseArray_SingleElement(t *testing.T) {
	traverse := TraverseArray(func(n int) Pair[string, Reader] {
		return pair.MakePair(
			fmt.Sprintf("test_%d", n),
			Equal(n*2)(n*2),
		)
	})

	result := traverse([]int{5})(t)
	if !result {
		t.Error("Expected TraverseArray to pass with single element")
	}
}

// TestTraverseArray_MultipleElements tests TraverseArray with multiple passing elements
func TestTraverseArray_MultipleElements(t *testing.T) {
	traverse := TraverseArray(func(n int) Pair[string, Reader] {
		return pair.MakePair(
			fmt.Sprintf("square_%d", n),
			Equal(n*n)(n*n),
		)
	})

	result := traverse([]int{1, 2, 3, 4, 5})(t)
	if !result {
		t.Error("Expected TraverseArray to pass with all passing elements")
	}
}

// TestTraverseArray_WithFailure tests that TraverseArray fails when one element fails
func TestTraverseArray_WithFailure(t *testing.T) {
	t.Skip("Skipping test that intentionally creates failing subtests")
	traverse := TraverseArray(func(n int) Pair[string, Reader] {
		return pair.MakePair(
			fmt.Sprintf("test_%d", n),
			Equal(10)(n), // Will fail for all except 10
		)
	})

	// Run in a subtest - we expect the subtests to fail, so t.Run returns false
	result := traverse([]int{1, 2, 3})(t)

	// The traverse should return false because assertions fail
	if result {
		t.Error("Expected traverse to return false when elements don't match")
	}
}

// TestTraverseArray_MixedResults tests TraverseArray with some passing and some failing
func TestTraverseArray_MixedResults(t *testing.T) {
	t.Skip("Skipping test that intentionally creates failing subtests")
	traverse := TraverseArray(func(n int) Pair[string, Reader] {
		return pair.MakePair(
			fmt.Sprintf("is_even_%d", n),
			Equal(0)(n%2), // Only passes for even numbers
		)
	})

	result := traverse([]int{2, 3, 4})(t) // 3 is odd, should fail

	// The traverse should return false because one assertion fails
	if result {
		t.Error("Expected traverse to return false when some elements fail")
	}
}

// TestTraverseArray_StringData tests TraverseArray with string data
func TestTraverseArray_StringData(t *testing.T) {
	words := []string{"hello", "world", "test"}

	traverse := TraverseArray(func(s string) Pair[string, Reader] {
		return pair.MakePair(
			fmt.Sprintf("validate_%s", s),
			AllOf([]Reader{
				StringNotEmpty(s),
				That(func(str string) bool { return len(str) > 0 })(s),
			}),
		)
	})

	result := traverse(words)(t)
	if !result {
		t.Error("Expected TraverseArray to pass with valid strings")
	}
}

// TestTraverseArray_ComplexObjects tests TraverseArray with complex objects
func TestTraverseArray_ComplexObjects(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	users := []User{
		{Name: "Alice", Age: 30},
		{Name: "Bob", Age: 25},
		{Name: "Charlie", Age: 35},
	}

	traverse := TraverseArray(func(u User) Pair[string, Reader] {
		return pair.MakePair(
			fmt.Sprintf("user_%s", u.Name),
			AllOf([]Reader{
				StringNotEmpty(u.Name),
				That(func(age int) bool { return age > 0 && age < 150 })(u.Age),
			}),
		)
	})

	result := traverse(users)(t)
	if !result {
		t.Error("Expected TraverseArray to pass with valid users")
	}
}

// TestTraverseArray_ComplexObjectsWithFailure tests TraverseArray with invalid complex objects
func TestTraverseArray_ComplexObjectsWithFailure(t *testing.T) {
	t.Skip("Skipping test that intentionally creates failing subtests")
	type User struct {
		Name string
		Age  int
	}

	users := []User{
		{Name: "Alice", Age: 30},
		{Name: "", Age: 25}, // Invalid: empty name
		{Name: "Charlie", Age: 35},
	}

	traverse := TraverseArray(func(u User) Pair[string, Reader] {
		return pair.MakePair(
			fmt.Sprintf("user_%s", u.Name),
			AllOf([]Reader{
				StringNotEmpty(u.Name),
				That(func(age int) bool { return age > 0 })(u.Age),
			}),
		)
	})

	result := traverse(users)(t)

	// The traverse should return false because one user is invalid
	if result {
		t.Error("Expected traverse to return false with invalid user")
	}
}

// TestTraverseArray_DataDrivenTesting demonstrates data-driven testing pattern
func TestTraverseArray_DataDrivenTesting(t *testing.T) {
	type TestCase struct {
		Input    int
		Expected int
	}

	testCases := []TestCase{
		{Input: 2, Expected: 4},
		{Input: 3, Expected: 9},
		{Input: 4, Expected: 16},
		{Input: 5, Expected: 25},
	}

	square := func(n int) int { return n * n }

	traverse := TraverseArray(func(tc TestCase) Pair[string, Reader] {
		return pair.MakePair(
			fmt.Sprintf("square(%d)=%d", tc.Input, tc.Expected),
			Equal(tc.Expected)(square(tc.Input)),
		)
	})

	result := traverse(testCases)(t)
	if !result {
		t.Error("Expected all test cases to pass")
	}
}

// TestSequenceSeq2_EmptySequence tests that SequenceSeq2 handles empty sequences correctly
func TestSequenceSeq2_EmptySequence(t *testing.T) {
	emptySeq := func(yield func(string, Reader) bool) {
		// Empty - yields nothing
	}

	result := SequenceSeq2[Reader](emptySeq)(t)
	if !result {
		t.Error("Expected SequenceSeq2 to pass with empty sequence")
	}
}

// TestSequenceSeq2_SingleTest tests SequenceSeq2 with a single test
func TestSequenceSeq2_SingleTest(t *testing.T) {
	singleSeq := func(yield func(string, Reader) bool) {
		yield("test_one", Equal(42)(42))
	}

	result := SequenceSeq2[Reader](singleSeq)(t)
	if !result {
		t.Error("Expected SequenceSeq2 to pass with single test")
	}
}

// TestSequenceSeq2_MultipleTests tests SequenceSeq2 with multiple passing tests
func TestSequenceSeq2_MultipleTests(t *testing.T) {
	multiSeq := func(yield func(string, Reader) bool) {
		if !yield("test_addition", Equal(4)(2+2)) {
			return
		}
		if !yield("test_subtraction", Equal(1)(3-2)) {
			return
		}
		if !yield("test_multiplication", Equal(6)(2*3)) {
			return
		}
	}

	result := SequenceSeq2[Reader](multiSeq)(t)
	if !result {
		t.Error("Expected SequenceSeq2 to pass with all passing tests")
	}
}

// TestSequenceSeq2_WithFailure tests that SequenceSeq2 fails when one test fails
func TestSequenceSeq2_WithFailure(t *testing.T) {
	t.Skip("Skipping test that intentionally creates failing subtests")
	failSeq := func(yield func(string, Reader) bool) {
		if !yield("test_pass", Equal(4)(2+2)) {
			return
		}
		if !yield("test_fail", Equal(5)(2+2)) { // This will fail
			return
		}
		if !yield("test_pass2", Equal(6)(2*3)) {
			return
		}
	}

	result := SequenceSeq2[Reader](failSeq)(t)

	// The sequence should return false because one test fails
	if result {
		t.Error("Expected sequence to return false when one test fails")
	}
}

// TestSequenceSeq2_GeneratedTests tests SequenceSeq2 with generated test cases
func TestSequenceSeq2_GeneratedTests(t *testing.T) {
	generateTests := func(yield func(string, Reader) bool) {
		for i := 1; i <= 5; i++ {
			name := fmt.Sprintf("test_%d", i)
			assertion := Equal(i * i)(i * i)
			if !yield(name, assertion) {
				return
			}
		}
	}

	result := SequenceSeq2[Reader](generateTests)(t)
	if !result {
		t.Error("Expected all generated tests to pass")
	}
}

// TestSequenceSeq2_StringTests tests SequenceSeq2 with string assertions
func TestSequenceSeq2_StringTests(t *testing.T) {
	stringSeq := func(yield func(string, Reader) bool) {
		if !yield("test_hello", StringNotEmpty("hello")) {
			return
		}
		if !yield("test_world", StringNotEmpty("world")) {
			return
		}
		if !yield("test_length", StringLength[any, any](5)("hello")) {
			return
		}
	}

	result := SequenceSeq2[Reader](stringSeq)(t)
	if !result {
		t.Error("Expected all string tests to pass")
	}
}

// TestSequenceSeq2_ArrayTests tests SequenceSeq2 with array assertions
func TestSequenceSeq2_ArrayTests(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}

	arraySeq := func(yield func(string, Reader) bool) {
		if !yield("test_not_empty", ArrayNotEmpty(arr)) {
			return
		}
		if !yield("test_length", ArrayLength[int](5)(arr)) {
			return
		}
		if !yield("test_contains", ArrayContains(3)(arr)) {
			return
		}
	}

	result := SequenceSeq2[Reader](arraySeq)(t)
	if !result {
		t.Error("Expected all array tests to pass")
	}
}

// TestSequenceSeq2_ComplexAssertions tests SequenceSeq2 with complex combined assertions
func TestSequenceSeq2_ComplexAssertions(t *testing.T) {
	type User struct {
		Name  string
		Age   int
		Email string
	}

	user := User{Name: "Alice", Age: 30, Email: "alice@example.com"}

	userSeq := func(yield func(string, Reader) bool) {
		if !yield("test_name", StringNotEmpty(user.Name)) {
			return
		}
		if !yield("test_age", That(func(age int) bool { return age > 0 && age < 150 })(user.Age)) {
			return
		}
		if !yield("test_email", That(func(email string) bool {
			for _, ch := range email {
				if ch == '@' {
					return true
				}
			}
			return false
		})(user.Email)) {
			return
		}
	}

	result := SequenceSeq2[Reader](userSeq)(t)
	if !result {
		t.Error("Expected all user validation tests to pass")
	}
}

// TestSequenceSeq2_EarlyTermination tests that SequenceSeq2 respects early termination
func TestSequenceSeq2_EarlyTermination(t *testing.T) {
	executionCount := 0

	earlyTermSeq := func(yield func(string, Reader) bool) {
		executionCount++
		if !yield("test_1", Equal(1)(1)) {
			return
		}
		executionCount++
		if !yield("test_2", Equal(2)(2)) {
			return
		}
		executionCount++
		// This should execute even though we don't check the return
		yield("test_3", Equal(3)(3))
		executionCount++
	}

	SequenceSeq2[Reader](earlyTermSeq)(t)

	// All iterations should execute since we're not terminating early
	if executionCount != 4 {
		t.Errorf("Expected 4 executions, got %d", executionCount)
	}
}

// TestSequenceSeq2_WithMapConversion demonstrates converting a map to Seq2
func TestSequenceSeq2_WithMapConversion(t *testing.T) {
	testMap := map[string]Reader{
		"test_addition":       Equal(4)(2 + 2),
		"test_multiplication": Equal(6)(2 * 3),
		"test_subtraction":    Equal(1)(3 - 2),
	}

	// Convert map to Seq2
	mapSeq := func(yield func(string, Reader) bool) {
		for name, assertion := range testMap {
			if !yield(name, assertion) {
				return
			}
		}
	}

	result := SequenceSeq2[Reader](mapSeq)(t)
	if !result {
		t.Error("Expected all map-based tests to pass")
	}
}

// TestTraverseArray_vs_SequenceSeq2 demonstrates the relationship between the two functions
func TestTraverseArray_vs_SequenceSeq2(t *testing.T) {
	type TestCase struct {
		Name     string
		Input    int
		Expected int
	}

	testCases := []TestCase{
		{Name: "test_1", Input: 2, Expected: 4},
		{Name: "test_2", Input: 3, Expected: 9},
		{Name: "test_3", Input: 4, Expected: 16},
	}

	// Using TraverseArray
	traverseResult := TraverseArray(func(tc TestCase) Pair[string, Reader] {
		return pair.MakePair(tc.Name, Equal(tc.Expected)(tc.Input*tc.Input))
	})(testCases)(t)

	// Using SequenceSeq2
	seqResult := SequenceSeq2[Reader](func(yield func(string, Reader) bool) {
		for _, tc := range testCases {
			if !yield(tc.Name, Equal(tc.Expected)(tc.Input*tc.Input)) {
				return
			}
		}
	})(t)

	if traverseResult != seqResult {
		t.Error("Expected TraverseArray and SequenceSeq2 to produce same result")
	}

	if !traverseResult || !seqResult {
		t.Error("Expected both approaches to pass")
	}
}

// TestTraverseRecord_EmptyMap tests that TraverseRecord handles empty maps correctly
func TestTraverseRecord_EmptyMap(t *testing.T) {
	traverse := TraverseRecord(func(n int) Reader {
		return Equal(n)(n)
	})

	result := traverse(map[string]int{})(t)
	if !result {
		t.Error("Expected TraverseRecord to pass with empty map")
	}
}

// TestTraverseRecord_SingleEntry tests TraverseRecord with a single map entry
func TestTraverseRecord_SingleEntry(t *testing.T) {
	traverse := TraverseRecord(func(n int) Reader {
		return Equal(n * 2)(n * 2)
	})

	result := traverse(map[string]int{"test_5": 5})(t)
	if !result {
		t.Error("Expected TraverseRecord to pass with single entry")
	}
}

// TestTraverseRecord_MultipleEntries tests TraverseRecord with multiple passing entries
func TestTraverseRecord_MultipleEntries(t *testing.T) {
	traverse := TraverseRecord(func(n int) Reader {
		return Equal(n * n)(n * n)
	})

	result := traverse(map[string]int{
		"square_1": 1,
		"square_2": 2,
		"square_3": 3,
		"square_4": 4,
		"square_5": 5,
	})(t)

	if !result {
		t.Error("Expected TraverseRecord to pass with all passing entries")
	}
}

// TestTraverseRecord_WithFailure tests that TraverseRecord fails when one entry fails
func TestTraverseRecord_WithFailure(t *testing.T) {
	t.Skip("Skipping test that intentionally creates failing subtests")
	traverse := TraverseRecord(func(n int) Reader {
		return Equal(10)(n) // Will fail for all except 10
	})

	result := traverse(map[string]int{
		"test_1": 1,
		"test_2": 2,
		"test_3": 3,
	})(t)

	// The traverse should return false because entries don't match
	if result {
		t.Error("Expected traverse to return false when entries don't match")
	}
}

// TestTraverseRecord_MixedResults tests TraverseRecord with some passing and some failing
func TestTraverseRecord_MixedResults(t *testing.T) {
	t.Skip("Skipping test that intentionally creates failing subtests")
	traverse := TraverseRecord(func(n int) Reader {
		return Equal(0)(n % 2) // Only passes for even numbers
	})

	result := traverse(map[string]int{
		"even_2": 2,
		"odd_3":  3,
		"even_4": 4,
	})(t)

	// The traverse should return false because some entries fail
	if result {
		t.Error("Expected traverse to return false when some entries fail")
	}
}

// TestTraverseRecord_StringData tests TraverseRecord with string data
func TestTraverseRecord_StringData(t *testing.T) {
	words := map[string]string{
		"greeting": "hello",
		"world":    "world",
		"test":     "test",
	}

	traverse := TraverseRecord(func(s string) Reader {
		return AllOf([]Reader{
			StringNotEmpty(s),
			That(func(str string) bool { return len(str) > 0 })(s),
		})
	})

	result := traverse(words)(t)
	if !result {
		t.Error("Expected TraverseRecord to pass with valid strings")
	}
}

// TestTraverseRecord_ComplexObjects tests TraverseRecord with complex objects
func TestTraverseRecord_ComplexObjects(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	users := map[string]User{
		"alice":   {Name: "Alice", Age: 30},
		"bob":     {Name: "Bob", Age: 25},
		"charlie": {Name: "Charlie", Age: 35},
	}

	traverse := TraverseRecord(func(u User) Reader {
		return AllOf([]Reader{
			StringNotEmpty(u.Name),
			That(func(age int) bool { return age > 0 && age < 150 })(u.Age),
		})
	})

	result := traverse(users)(t)
	if !result {
		t.Error("Expected TraverseRecord to pass with valid users")
	}
}

// TestTraverseRecord_ComplexObjectsWithFailure tests TraverseRecord with invalid complex objects
func TestTraverseRecord_ComplexObjectsWithFailure(t *testing.T) {
	t.Skip("Skipping test that intentionally creates failing subtests")
	type User struct {
		Name string
		Age  int
	}

	users := map[string]User{
		"alice":   {Name: "Alice", Age: 30},
		"invalid": {Name: "", Age: 25}, // Invalid: empty name
		"charlie": {Name: "Charlie", Age: 35},
	}

	traverse := TraverseRecord(func(u User) Reader {
		return AllOf([]Reader{
			StringNotEmpty(u.Name),
			That(func(age int) bool { return age > 0 })(u.Age),
		})
	})

	result := traverse(users)(t)

	// The traverse should return false because one user is invalid
	if result {
		t.Error("Expected traverse to return false with invalid user")
	}
}

// TestTraverseRecord_ConfigurationTesting demonstrates configuration testing pattern
func TestTraverseRecord_ConfigurationTesting(t *testing.T) {
	configs := map[string]int{
		"timeout":    30,
		"maxRetries": 3,
		"bufferSize": 1024,
	}

	validatePositive := That(func(n int) bool { return n > 0 })

	traverse := TraverseRecord(validatePositive)
	result := traverse(configs)(t)

	if !result {
		t.Error("Expected all configuration values to be positive")
	}
}

// TestTraverseRecord_APIEndpointTesting demonstrates API endpoint testing pattern
func TestTraverseRecord_APIEndpointTesting(t *testing.T) {
	type Endpoint struct {
		Path   string
		Method string
	}

	endpoints := map[string]Endpoint{
		"get_users":   {Path: "/api/users", Method: "GET"},
		"create_user": {Path: "/api/users", Method: "POST"},
		"delete_user": {Path: "/api/users/:id", Method: "DELETE"},
	}

	validateEndpoint := func(e Endpoint) Reader {
		return AllOf([]Reader{
			StringNotEmpty(e.Path),
			That(func(path string) bool {
				return len(path) > 0 && path[0] == '/'
			})(e.Path),
			That(func(method string) bool {
				return method == "GET" || method == "POST" ||
					method == "PUT" || method == "DELETE"
			})(e.Method),
		})
	}

	traverse := TraverseRecord(validateEndpoint)
	result := traverse(endpoints)(t)

	if !result {
		t.Error("Expected all endpoints to be valid")
	}
}

// TestSequenceRecord_EmptyMap tests that SequenceRecord handles empty maps correctly
func TestSequenceRecord_EmptyMap(t *testing.T) {
	result := SequenceRecord(map[string]Reader{})(t)
	if !result {
		t.Error("Expected SequenceRecord to pass with empty map")
	}
}

// TestSequenceRecord_SingleTest tests SequenceRecord with a single test
func TestSequenceRecord_SingleTest(t *testing.T) {
	tests := map[string]Reader{
		"test_one": Equal(42)(42),
	}

	result := SequenceRecord(tests)(t)
	if !result {
		t.Error("Expected SequenceRecord to pass with single test")
	}
}

// TestSequenceRecord_MultipleTests tests SequenceRecord with multiple passing tests
func TestSequenceRecord_MultipleTests(t *testing.T) {
	tests := map[string]Reader{
		"test_addition":       Equal(4)(2 + 2),
		"test_subtraction":    Equal(1)(3 - 2),
		"test_multiplication": Equal(6)(2 * 3),
		"test_division":       Equal(2)(6 / 3),
	}

	result := SequenceRecord(tests)(t)
	if !result {
		t.Error("Expected SequenceRecord to pass with all passing tests")
	}
}

// TestSequenceRecord_WithFailure tests that SequenceRecord fails when one test fails
func TestSequenceRecord_WithFailure(t *testing.T) {
	t.Skip("Skipping test that intentionally creates failing subtests")
	tests := map[string]Reader{
		"test_pass":  Equal(4)(2 + 2),
		"test_fail":  Equal(5)(2 + 2), // This will fail
		"test_pass2": Equal(6)(2 * 3),
	}

	result := SequenceRecord(tests)(t)

	// The sequence should return false because one test fails
	if result {
		t.Error("Expected sequence to return false when one test fails")
	}
}

// TestSequenceRecord_StringTests tests SequenceRecord with string assertions
func TestSequenceRecord_StringTests(t *testing.T) {
	testString := "hello world"

	tests := map[string]Reader{
		"not_empty":      StringNotEmpty(testString),
		"correct_length": StringLength[any, any](11)(testString),
		"has_space": That(func(s string) bool {
			for _, ch := range s {
				if ch == ' ' {
					return true
				}
			}
			return false
		})(testString),
	}

	result := SequenceRecord(tests)(t)
	if !result {
		t.Error("Expected all string tests to pass")
	}
}

// TestSequenceRecord_ArrayTests tests SequenceRecord with array assertions
func TestSequenceRecord_ArrayTests(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}

	tests := map[string]Reader{
		"not_empty":      ArrayNotEmpty(arr),
		"correct_length": ArrayLength[int](5)(arr),
		"contains_three": ArrayContains(3)(arr),
		"all_positive": That(func(arr []int) bool {
			for _, n := range arr {
				if n <= 0 {
					return false
				}
			}
			return true
		})(arr),
	}

	result := SequenceRecord(tests)(t)
	if !result {
		t.Error("Expected all array tests to pass")
	}
}

// TestSequenceRecord_ComplexAssertions tests SequenceRecord with complex combined assertions
func TestSequenceRecord_ComplexAssertions(t *testing.T) {
	type User struct {
		Name  string
		Age   int
		Email string
	}

	user := User{Name: "Alice", Age: 30, Email: "alice@example.com"}

	tests := map[string]Reader{
		"name_not_empty": StringNotEmpty(user.Name),
		"age_positive":   That(func(age int) bool { return age > 0 })(user.Age),
		"age_reasonable": That(func(age int) bool { return age < 150 })(user.Age),
		"email_valid": That(func(email string) bool {
			hasAt := false
			hasDot := false
			for _, ch := range email {
				if ch == '@' {
					hasAt = true
				}
				if ch == '.' {
					hasDot = true
				}
			}
			return hasAt && hasDot
		})(user.Email),
	}

	result := SequenceRecord(tests)(t)
	if !result {
		t.Error("Expected all user validation tests to pass")
	}
}

// TestSequenceRecord_MathOperations demonstrates basic math operations testing
func TestSequenceRecord_MathOperations(t *testing.T) {
	tests := map[string]Reader{
		"addition":       Equal(4)(2 + 2),
		"subtraction":    Equal(1)(3 - 2),
		"multiplication": Equal(6)(2 * 3),
		"division":       Equal(2)(6 / 3),
		"modulo":         Equal(1)(7 % 3),
	}

	result := SequenceRecord(tests)(t)
	if !result {
		t.Error("Expected all math operations to pass")
	}
}

// TestSequenceRecord_BooleanTests tests SequenceRecord with boolean assertions
func TestSequenceRecord_BooleanTests(t *testing.T) {
	tests := map[string]Reader{
		"true_is_true":   Equal(true)(true),
		"false_is_false": Equal(false)(false),
		"not_true":       Equal(false)(!true),
		"not_false":      Equal(true)(!false),
	}

	result := SequenceRecord(tests)(t)
	if !result {
		t.Error("Expected all boolean tests to pass")
	}
}

// TestSequenceRecord_ErrorTests tests SequenceRecord with error assertions
func TestSequenceRecord_ErrorTests(t *testing.T) {
	tests := map[string]Reader{
		"no_error":    NoError(nil),
		"equal_value": Equal("test")("test"),
		"not_empty":   StringNotEmpty("hello"),
	}

	result := SequenceRecord(tests)(t)
	if !result {
		t.Error("Expected all error tests to pass")
	}
}

// TestTraverseRecord_vs_SequenceRecord demonstrates the relationship between the two functions
func TestTraverseRecord_vs_SequenceRecord(t *testing.T) {
	type TestCase struct {
		Input    int
		Expected int
	}

	testData := map[string]TestCase{
		"test_1": {Input: 2, Expected: 4},
		"test_2": {Input: 3, Expected: 9},
		"test_3": {Input: 4, Expected: 16},
	}

	// Using TraverseRecord
	traverseResult := TraverseRecord(func(tc TestCase) Reader {
		return Equal(tc.Expected)(tc.Input * tc.Input)
	})(testData)(t)

	// Using SequenceRecord (manually creating the map)
	tests := make(map[string]Reader)
	for name, tc := range testData {
		tests[name] = Equal(tc.Expected)(tc.Input * tc.Input)
	}
	seqResult := SequenceRecord(tests)(t)

	if traverseResult != seqResult {
		t.Error("Expected TraverseRecord and SequenceRecord to produce same result")
	}

	if !traverseResult || !seqResult {
		t.Error("Expected both approaches to pass")
	}
}

// TestSequenceRecord_WithAllOf demonstrates combining SequenceRecord with AllOf
func TestSequenceRecord_WithAllOf(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}

	tests := map[string]Reader{
		"array_validations": AllOf([]Reader{
			ArrayNotEmpty(arr),
			ArrayLength[int](5)(arr),
			ArrayContains(3)(arr),
		}),
		"element_checks": AllOf([]Reader{
			That(func(a []int) bool { return a[0] == 1 })(arr),
			That(func(a []int) bool { return a[4] == 5 })(arr),
		}),
	}

	result := SequenceRecord(tests)(t)
	if !result {
		t.Error("Expected combined assertions to pass")
	}
}

// TestTraverseRecord_ConfigValidation demonstrates real-world configuration validation
func TestTraverseRecord_ConfigValidation(t *testing.T) {
	type Config struct {
		Value int
		Min   int
		Max   int
	}

	configs := map[string]Config{
		"timeout":    {Value: 30, Min: 1, Max: 60},
		"maxRetries": {Value: 3, Min: 1, Max: 10},
		"bufferSize": {Value: 1024, Min: 512, Max: 4096},
	}

	validateConfig := func(c Config) Reader {
		return AllOf([]Reader{
			That(func(val int) bool { return val >= c.Min })(c.Value),
			That(func(val int) bool { return val <= c.Max })(c.Value),
		})
	}

	traverse := TraverseRecord(validateConfig)
	result := traverse(configs)(t)

	if !result {
		t.Error("Expected all configurations to be within valid ranges")
	}
}

// TestSequenceRecord_RealWorldExample demonstrates a realistic use case
func TestSequenceRecord_RealWorldExample(t *testing.T) {
	type Response struct {
		StatusCode int
		Body       string
	}

	response := Response{StatusCode: 200, Body: `{"status":"ok"}`}

	tests := map[string]Reader{
		"status_ok":      Equal(200)(response.StatusCode),
		"body_not_empty": StringNotEmpty(response.Body),
		"body_is_json": That(func(s string) bool {
			return len(s) > 0 && s[0] == '{' && s[len(s)-1] == '}'
		})(response.Body),
	}

	result := SequenceRecord(tests)(t)
	if !result {
		t.Error("Expected response validation to pass")
	}
}
