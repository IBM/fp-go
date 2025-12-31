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
	filter := FilterOrElse(isPositive, onNegative)

	// Test Right value that passes predicate
	result := filter(Right(5))
	assert.Equal(t, Right(5), result)

	// Test Right value that fails predicate
	result = filter(Right(-3))
	assert.True(t, IsLeft(result))

	// Test Right value at boundary (zero)
	result = filter(Right(0))
	assert.True(t, IsLeft(result))

	// Test Left value (should pass through unchanged)
	originalError := errors.New("original error")
	result = filter(Left[int](originalError))
	assert.Equal(t, Left[int](originalError), result)
}

func TestFilterOrElse_StringValidation(t *testing.T) {
	// Test with string length validation
	isNotEmpty := func(s string) bool { return len(s) > 0 }
	onEmpty := func(s string) error { return errors.New("string is empty") }
	filter := FilterOrElse(isNotEmpty, onEmpty)

	// Test non-empty string
	result := filter(Right("hello"))
	assert.Equal(t, Right("hello"), result)

	// Test empty string
	result = filter(Right(""))
	assert.True(t, IsLeft(result))

	// Test Left value
	originalError := errors.New("validation error")
	result = filter(Left[string](originalError))
	assert.Equal(t, Left[string](originalError), result)
}

func TestFilterOrElse_ComplexPredicate(t *testing.T) {
	// Test with range validation
	inRange := func(x int) bool { return x >= 10 && x <= 100 }
	outOfRange := func(x int) error { return fmt.Errorf("%d is out of range [10, 100]", x) }
	filter := FilterOrElse(inRange, outOfRange)

	// Test value in range
	result := filter(Right(50))
	assert.Equal(t, Right(50), result)

	// Test value below range
	result = filter(Right(5))
	assert.True(t, IsLeft(result))

	// Test value above range
	result = filter(Right(150))
	assert.True(t, IsLeft(result))

	// Test boundary values
	result = filter(Right(10))
	assert.Equal(t, Right(10), result)

	result = filter(Right(100))
	assert.Equal(t, Right(100), result)
}

func TestFilterOrElse_ChainedFilters(t *testing.T) {
	// Test chaining multiple filters
	isPositive := N.MoreThan(0)
	onNegative := func(x int) error { return fmt.Errorf("%d is not positive", x) }

	isEven := func(x int) bool { return x%2 == 0 }
	onOdd := func(x int) error { return fmt.Errorf("%d is not even", x) }

	filterPositive := FilterOrElse(isPositive, onNegative)
	filterEven := FilterOrElse(isEven, onOdd)

	// Test value that passes both filters
	result := filterEven(filterPositive(Right(4)))
	assert.Equal(t, Right(4), result)

	// Test value that fails first filter
	result = filterEven(filterPositive(Right(-2)))
	assert.True(t, IsLeft(result))

	// Test value that passes first but fails second filter
	result = filterEven(filterPositive(Right(3)))
	assert.True(t, IsLeft(result))
}

func TestFilterOrElse_WithStructs(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	// Test with struct validation
	isAdult := func(u User) bool { return u.Age >= 18 }
	onMinor := func(u User) error { return fmt.Errorf("%s is not an adult (age: %d)", u.Name, u.Age) }
	filter := FilterOrElse(isAdult, onMinor)

	// Test adult user
	adult := User{Name: "Alice", Age: 25}
	result := filter(Right(adult))
	assert.Equal(t, Right(adult), result)

	// Test minor user
	minor := User{Name: "Bob", Age: 16}
	result = filter(Right(minor))
	assert.True(t, IsLeft(result))
}

func TestFilterOrElse_WithOf(t *testing.T) {
	// Test using Of constructor
	isPositive := N.MoreThan(0)
	onNegative := func(x int) error { return fmt.Errorf("%d is not positive", x) }
	filter := FilterOrElse(isPositive, onNegative)

	// Test with Of
	result := filter(Of(5))
	assert.Equal(t, Of(5), result)

	result = filter(Of(-3))
	assert.True(t, IsLeft(result))
}

func TestFilterOrElse_ErrorMessages(t *testing.T) {
	// Test that error messages are properly propagated
	isPositive := N.MoreThan(0)
	onNegative := func(x int) error { return fmt.Errorf("value %d is not positive", x) }
	filter := FilterOrElse(isPositive, onNegative)

	result := filter(Right(-5))
	assert.True(t, IsLeft(result))

	_, err := UnwrapError(result)
	assert.Equal(t, "value -5 is not positive", err.Error())
}
