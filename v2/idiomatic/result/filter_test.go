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
	"errors"
	"fmt"
	"testing"

	N "github.com/IBM/fp-go/v2/number"
	"github.com/stretchr/testify/assert"
)

func TestFilterOrElse(t *testing.T) {
	// Test with positive predicate
	isPositive := N.MoreThan(0)
	onNegative := func(x int) error { return fmt.Errorf("%d is not positive", x) }

	// Test value that passes predicate
	AssertEq(Right(5))(Pipe2(5, Of, FilterOrElse(isPositive, onNegative)))(t)

	// Test value that fails predicate
	_, err := Pipe2(-3, Of, FilterOrElse(isPositive, onNegative))
	assert.Error(t, err)
	assert.Equal(t, "-3 is not positive", err.Error())

	// Test value at boundary (zero)
	_, err = Pipe2(0, Of, FilterOrElse(isPositive, onNegative))
	assert.Error(t, err)

	// Test error value (should pass through unchanged)
	originalError := errors.New("original error")
	_, err = Pipe2(originalError, Left[int], FilterOrElse(isPositive, onNegative))
	assert.Error(t, err)
	assert.Equal(t, originalError, err)
}

func TestFilterOrElse_StringValidation(t *testing.T) {
	// Test with string length validation
	isNotEmpty := func(s string) bool { return len(s) > 0 }
	onEmpty := func(s string) error { return errors.New("string is empty") }

	// Test non-empty string
	AssertEq(Right("hello"))(Pipe2("hello", Of, FilterOrElse(isNotEmpty, onEmpty)))(t)

	// Test empty string
	_, err := Pipe2("", Of, FilterOrElse(isNotEmpty, onEmpty))
	assert.Error(t, err)
	assert.Equal(t, "string is empty", err.Error())

	// Test error value
	originalError := errors.New("validation error")
	_, err = Pipe2(originalError, Left[string], FilterOrElse(isNotEmpty, onEmpty))
	assert.Error(t, err)
	assert.Equal(t, originalError, err)
}

func TestFilterOrElse_ComplexPredicate(t *testing.T) {
	// Test with range validation
	inRange := func(x int) bool { return x >= 10 && x <= 100 }
	outOfRange := func(x int) error { return fmt.Errorf("%d is out of range [10, 100]", x) }

	// Test value in range
	AssertEq(Right(50))(Pipe2(50, Of, FilterOrElse(inRange, outOfRange)))(t)

	// Test value below range
	_, err := Pipe2(5, Of, FilterOrElse(inRange, outOfRange))
	assert.Error(t, err)

	// Test value above range
	_, err = Pipe2(150, Of, FilterOrElse(inRange, outOfRange))
	assert.Error(t, err)

	// Test boundary values
	AssertEq(Right(10))(Pipe2(10, Of, FilterOrElse(inRange, outOfRange)))(t)
	AssertEq(Right(100))(Pipe2(100, Of, FilterOrElse(inRange, outOfRange)))(t)
}

func TestFilterOrElse_ChainedFilters(t *testing.T) {
	// Test chaining multiple filters
	isPositive := N.MoreThan(0)
	onNegative := func(x int) error { return fmt.Errorf("%d is not positive", x) }

	isEven := func(x int) bool { return x%2 == 0 }
	onOdd := func(x int) error { return fmt.Errorf("%d is not even", x) }

	// Test value that passes both filters
	AssertEq(Right(4))(Pipe3(4, Of, FilterOrElse(isPositive, onNegative), FilterOrElse(isEven, onOdd)))(t)

	// Test value that fails first filter
	_, err := Pipe3(-2, Of, FilterOrElse(isPositive, onNegative), FilterOrElse(isEven, onOdd))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not positive")

	// Test value that passes first but fails second filter
	_, err = Pipe3(3, Of, FilterOrElse(isPositive, onNegative), FilterOrElse(isEven, onOdd))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not even")
}

func TestFilterOrElse_WithStructs(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	// Test with struct validation
	isAdult := func(u User) bool { return u.Age >= 18 }
	onMinor := func(u User) error { return fmt.Errorf("%s is not an adult (age: %d)", u.Name, u.Age) }

	// Test adult user
	adult := User{Name: "Alice", Age: 25}
	AssertEq(Right(adult))(Pipe2(adult, Of, FilterOrElse(isAdult, onMinor)))(t)

	// Test minor user
	minor := User{Name: "Bob", Age: 16}
	_, err := Pipe2(minor, Of, FilterOrElse(isAdult, onMinor))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Bob is not an adult")
}

func TestFilterOrElse_WithChain(t *testing.T) {
	// Test FilterOrElse in a chain with other operations
	isPositive := N.MoreThan(0)
	onNegative := func(x int) error { return fmt.Errorf("%d is not positive", x) }

	double := func(x int) (int, error) { return x * 2, nil }

	// Test successful chain
	AssertEq(Right(10))(Pipe3(5, Of, FilterOrElse(isPositive, onNegative), Chain(double)))(t)

	// Test chain with filter failure
	_, err := Pipe3(-5, Of, FilterOrElse(isPositive, onNegative), Chain(double))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not positive")
}

func TestFilterOrElse_ErrorMessages(t *testing.T) {
	// Test that error messages are properly propagated
	isPositive := N.MoreThan(0)
	onNegative := func(x int) error { return fmt.Errorf("value %d is not positive", x) }

	result, err := Pipe2(-5, Of, FilterOrElse(isPositive, onNegative))
	assert.Error(t, err)
	assert.Equal(t, "value -5 is not positive", err.Error())
	assert.Equal(t, 0, result) // default value for int
}
