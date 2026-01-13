// Copyright (c) 2025 IBM Corp.
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
	"errors"
	"fmt"
	"strings"
	"testing"

	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/reader"
	S "github.com/IBM/fp-go/v2/semigroup"
	STR "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// Helper function to create a semigroup that concatenates error messages
func makeErrorConcatSemigroup() S.Semigroup[error] {
	return S.MakeSemigroup(func(e1, e2 error) error {
		return fmt.Errorf("%v; %v", e1, e2)
	})
}

// Helper function to create a semigroup that collects error messages in a slice
func makeErrorListSemigroup() S.Semigroup[error] {
	return S.MakeSemigroup(func(e1, e2 error) error {
		msg1 := e1.Error()
		msg2 := e2.Error()

		// Parse existing lists
		var msgs []string
		if strings.HasPrefix(msg1, "[") && strings.HasSuffix(msg1, "]") {
			trimmed := strings.Trim(msg1, "[]")
			if STR.IsNonEmpty(trimmed) {
				msgs = strings.Split(trimmed, ", ")
			}
		} else {
			msgs = []string{msg1}
		}

		if strings.HasPrefix(msg2, "[") && strings.HasSuffix(msg2, "]") {
			trimmed := strings.Trim(msg2, "[]")
			if STR.IsNonEmpty(trimmed) {
				msgs = append(msgs, strings.Split(trimmed, ", ")...)
			}
		} else {
			msgs = append(msgs, msg2)
		}

		return fmt.Errorf("[%s]", strings.Join(msgs, ", "))
	})
}

// TestApV_BothRight tests ApV when both the value and function are Right
func TestApV_BothRight(t *testing.T) {
	sg := makeErrorConcatSemigroup()
	apv := ApV[int, int](sg)

	double := N.Mul(2)

	value, verr := Right(5)
	fn, ferr := Right(double)

	result, err := apv(value, verr)(fn, ferr)

	assert.NoError(t, err)
	assert.Equal(t, 10, result)
}

// TestApV_ValueLeft_FunctionRight tests ApV when value is Left and function is Right
func TestApV_ValueLeft_FunctionRight(t *testing.T) {
	sg := makeErrorConcatSemigroup()
	apv := ApV[int, int](sg)

	double := N.Mul(2)

	valueError := errors.New("invalid value")
	value, verr := Left[int](valueError)
	fn, ferr := Right(double)

	result, err := apv(value, verr)(fn, ferr)

	assert.Error(t, err)
	assert.Equal(t, valueError, err)
	assert.Equal(t, 0, result) // zero value for int
}

// TestApV_ValueRight_FunctionLeft tests ApV when value is Right and function is Left
func TestApV_ValueRight_FunctionLeft(t *testing.T) {
	sg := makeErrorConcatSemigroup()
	apv := ApV[int, int](sg)

	fnError := errors.New("invalid function")
	value, verr := Right(5)
	fn, ferr := Left[func(int) int](fnError)

	result, err := apv(value, verr)(fn, ferr)

	assert.Error(t, err)
	assert.Equal(t, fnError, err)
	assert.Equal(t, 0, result) // zero value for int
}

// TestApV_BothLeft tests ApV when both value and function are Left
func TestApV_BothLeft(t *testing.T) {
	sg := makeErrorConcatSemigroup()
	apv := ApV[int, int](sg)

	valueError := errors.New("invalid value")
	fnError := errors.New("invalid function")

	value, verr := Left[int](valueError)
	fn, ferr := Left[func(int) int](fnError)

	result, err := apv(value, verr)(fn, ferr)

	assert.Error(t, err)
	assert.Equal(t, 0, result) // zero value for int

	// Verify the error message contains both errors
	expectedMsg := "invalid function; invalid value"
	assert.Equal(t, expectedMsg, err.Error())
}

// TestApV_BothLeft_WithListSemigroup tests error accumulation with a list semigroup
func TestApV_BothLeft_WithListSemigroup(t *testing.T) {
	sg := makeErrorListSemigroup()
	apv := ApV[string, string](sg)

	valueError := errors.New("error1")
	fnError := errors.New("error2")

	value, verr := Left[string](valueError)
	fn, ferr := Left[func(string) string](fnError)

	result, err := apv(value, verr)(fn, ferr)

	assert.Error(t, err)
	assert.Equal(t, "", result) // zero value for string

	// Verify both errors are in the list
	expectedMsg := "[error2, error1]"
	assert.Equal(t, expectedMsg, err.Error())
}

// TestApV_StringTransformation tests ApV with string transformation
func TestApV_StringTransformation(t *testing.T) {
	sg := makeErrorConcatSemigroup()
	apv := ApV[string, string](sg)

	toUpper := strings.ToUpper

	value, verr := Right("hello")
	fn, ferr := Right(toUpper)

	result, err := apv(value, verr)(fn, ferr)

	assert.NoError(t, err)
	assert.Equal(t, "HELLO", result)
}

// TestApV_DifferentTypes tests ApV with different input and output types
func TestApV_DifferentTypes(t *testing.T) {
	sg := makeErrorConcatSemigroup()
	apv := ApV[string, int](sg)

	intToString := func(x int) string { return fmt.Sprintf("Number: %d", x) }

	value, verr := Right(42)
	fn, ferr := Right(intToString)

	result, err := apv(value, verr)(fn, ferr)

	assert.NoError(t, err)
	assert.Equal(t, "Number: 42", result)
}

// TestApV_ComplexType tests ApV with complex types (structs)
func TestApV_ComplexType(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	sg := makeErrorConcatSemigroup()
	apv := ApV[string, Person](sg)

	getName := func(p Person) string { return p.Name }

	person := Person{Name: "Alice", Age: 30}
	value, verr := Right(person)
	fn, ferr := Right(getName)

	result, err := apv(value, verr)(fn, ferr)

	assert.NoError(t, err)
	assert.Equal(t, "Alice", result)
}

// TestApV_MultipleValidations demonstrates chaining multiple validations
func TestApV_MultipleValidations(t *testing.T) {
	_ = makeErrorListSemigroup() // Semigroup available for future use

	// Validation functions
	validatePositive := func(x int) (int, error) {
		if x > 0 {
			return Right(x)
		}
		return Left[int](errors.New("must be positive"))
	}

	validateEven := func(x int) (int, error) {
		if x%2 == 0 {
			return Right(x)
		}
		return Left[int](errors.New("must be even"))
	}

	// Test valid value (positive and even)
	t.Run("valid value", func(t *testing.T) {
		value, err := Right(4)
		validatedPositive, err1 := validatePositive(value)
		validatedEven, err2 := validateEven(validatedPositive)

		assert.NoError(t, err)
		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, 4, validatedEven)
	})

	// Test invalid value (negative)
	t.Run("negative value", func(t *testing.T) {
		value, _ := Right(-3)
		_, err := validatePositive(value)

		assert.Error(t, err)
		assert.Equal(t, "must be positive", err.Error())
	})

	// Test invalid value (odd)
	t.Run("odd value", func(t *testing.T) {
		value, _ := Right(3)
		validatedPositive, err1 := validatePositive(value)
		_, err2 := validateEven(validatedPositive)

		assert.NoError(t, err1)
		assert.Error(t, err2)
		assert.Equal(t, "must be even", err2.Error())
	})
}

// TestApV_ZeroValues tests ApV with zero values
func TestApV_ZeroValues(t *testing.T) {
	sg := makeErrorConcatSemigroup()
	apv := ApV[int, int](sg)

	identity := reader.Ask[int]()

	value, verr := Right(0)
	fn, ferr := Right(identity)

	result, err := apv(value, verr)(fn, ferr)

	assert.NoError(t, err)
	assert.Equal(t, 0, result)
}

// TestApV_NilError tests that nil errors are handled correctly
func TestApV_NilError(t *testing.T) {
	sg := makeErrorConcatSemigroup()
	apv := ApV[string, string](sg)

	identity := func(s string) string { return s }

	// Right is equivalent to (value, nil)
	value, verr := Right("test")
	fn, ferr := Right(identity)

	assert.Nil(t, verr)
	assert.Nil(t, ferr)

	result, err := apv(value, verr)(fn, ferr)

	assert.NoError(t, err)
	assert.Equal(t, "test", result)
}

// TestApV_SemigroupAssociativity tests that error combination is associative
func TestApV_SemigroupAssociativity(t *testing.T) {
	sg := makeErrorConcatSemigroup()

	e1 := errors.New("error1")
	e2 := errors.New("error2")
	e3 := errors.New("error3")

	// (e1 + e2) + e3
	left := sg.Concat(sg.Concat(e1, e2), e3)
	// e1 + (e2 + e3)
	right := sg.Concat(e1, sg.Concat(e2, e3))

	assert.Equal(t, left.Error(), right.Error())
}

// TestApV_CustomSemigroup tests ApV with a custom semigroup
func TestApV_CustomSemigroup(t *testing.T) {
	// Custom semigroup that counts errors
	type ErrorCount struct {
		count int
		msg   string
	}

	countSemigroup := S.MakeSemigroup(func(e1, e2 error) error {
		// Simple counter in error message
		return fmt.Errorf("combined: %v | %v", e1, e2)
	})

	apv := ApV[int, int](countSemigroup)

	e1 := errors.New("first")
	e2 := errors.New("second")

	value, verr := Left[int](e1)
	fn, ferr := Left[func(int) int](e2)

	_, err := apv(value, verr)(fn, ferr)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "first")
	assert.Contains(t, err.Error(), "second")
}

// BenchmarkApV_BothRight benchmarks the happy path
func BenchmarkApV_BothRight(b *testing.B) {
	sg := makeErrorConcatSemigroup()
	apv := ApV[int, int](sg)

	double := N.Mul(2)
	value, verr := Right(5)
	fn, ferr := Right(double)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		apv(value, verr)(fn, ferr)
	}
}

// BenchmarkApV_BothLeft benchmarks the error accumulation path
func BenchmarkApV_BothLeft(b *testing.B) {
	sg := makeErrorConcatSemigroup()
	apv := ApV[int, int](sg)

	valueError := errors.New("value error")
	fnError := errors.New("function error")

	value, verr := Left[int](valueError)
	fn, ferr := Left[func(int) int](fnError)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		apv(value, verr)(fn, ferr)
	}
}
