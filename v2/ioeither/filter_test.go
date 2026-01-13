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

package ioeither

import (
	"fmt"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	N "github.com/IBM/fp-go/v2/number"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

func TestFilterOrElse_PredicateTrue(t *testing.T) {
	// Test that when predicate returns true, Right value passes through
	isPositive := N.MoreThan(0)
	onFalse := S.Format[int]("%d is not positive")

	filter := FilterOrElse(isPositive, onFalse)
	result := filter(Right[string](42))()

	assert.Equal(t, E.Right[string](42), result)
}

func TestFilterOrElse_PredicateFalse(t *testing.T) {
	// Test that when predicate returns false, Right value becomes Left
	isPositive := N.MoreThan(0)
	onFalse := S.Format[int]("%d is not positive")

	filter := FilterOrElse(isPositive, onFalse)
	result := filter(Right[string](-5))()

	assert.Equal(t, E.Left[int]("-5 is not positive"), result)
}

func TestFilterOrElse_LeftPassesThrough(t *testing.T) {
	// Test that Left values pass through unchanged
	isPositive := N.MoreThan(0)
	onFalse := S.Format[int]("%d is not positive")

	filter := FilterOrElse(isPositive, onFalse)
	result := filter(Left[int]("original error"))()

	assert.Equal(t, E.Left[int]("original error"), result)
}

func TestFilterOrElse_ZeroValue(t *testing.T) {
	// Test filtering with zero value
	isNonZero := func(n int) bool { return n != 0 }
	onZero := func(n int) string { return "value is zero" }

	filter := FilterOrElse(isNonZero, onZero)
	result := filter(Right[string](0))()

	assert.Equal(t, E.Left[int]("value is zero"), result)
}

func TestFilterOrElse_StringValidation(t *testing.T) {
	// Test with string validation
	isNonEmpty := S.IsNonEmpty
	onEmpty := func(s string) error { return fmt.Errorf("string is empty") }

	filter := FilterOrElse(isNonEmpty, onEmpty)

	// Non-empty string passes
	result1 := filter(Right[error]("hello"))()
	assert.Equal(t, E.Right[error]("hello"), result1)

	// Empty string becomes error
	result2 := filter(Right[error](""))()
	assert.True(t, E.IsLeft(result2))
	assert.Equal(t, "string is empty", E.ToError(result2).Error())
}

func TestFilterOrElse_ComplexPredicate(t *testing.T) {
	// Test with more complex predicate
	type User struct {
		Name string
		Age  int
	}

	isAdult := func(u User) bool { return u.Age >= 18 }
	onMinor := func(u User) string {
		return fmt.Sprintf("%s is only %d years old", u.Name, u.Age)
	}

	filter := FilterOrElse(isAdult, onMinor)

	// Adult user passes
	adult := User{Name: "Alice", Age: 25}
	result1 := filter(Right[string](adult))()
	assert.Equal(t, E.Right[string](adult), result1)

	// Minor becomes error
	minor := User{Name: "Bob", Age: 16}
	result2 := filter(Right[string](minor))()
	assert.Equal(t, E.Left[User]("Bob is only 16 years old"), result2)
}

func TestFilterOrElse_ChainedFilters(t *testing.T) {
	// Test chaining multiple filters
	isPositive := N.MoreThan(0)
	onNegative := func(n int) string { return "not positive" }

	isEven := func(n int) bool { return n%2 == 0 }
	onOdd := func(n int) string { return "not even" }

	filter1 := FilterOrElse(isPositive, onNegative)
	filter2 := FilterOrElse(isEven, onOdd)

	// Chain filters - apply filter1 first, then filter2
	result := filter2(filter1(Right[string](4)))()
	assert.Equal(t, E.Right[string](4), result)

	// Fails first filter
	result2 := filter2(filter1(Right[string](-2)))()
	assert.Equal(t, E.Left[int]("not positive"), result2)

	// Passes first but fails second
	result3 := filter2(filter1(Right[string](3)))()
	assert.Equal(t, E.Left[int]("not even"), result3)
}

func TestFilterOrElse_WithMap(t *testing.T) {
	// Test FilterOrElse combined with Map
	isPositive := N.MoreThan(0)
	onNegative := func(n int) string { return "negative number" }

	filter := FilterOrElse(isPositive, onNegative)
	double := Map[string](N.Mul(2))

	// Compose: filter then double
	result1 := double(filter(Right[string](5)))()
	assert.Equal(t, E.Right[string](10), result1)

	// Negative value filtered out
	result2 := double(filter(Right[string](-3)))()
	assert.Equal(t, E.Left[int]("negative number"), result2)
}

func TestFilterOrElse_BoundaryConditions(t *testing.T) {
	// Test boundary conditions
	isInRange := func(n int) bool { return n >= 0 && n <= 100 }
	onOutOfRange := func(n int) string {
		return fmt.Sprintf("%d is out of range [0, 100]", n)
	}

	filter := FilterOrElse(isInRange, onOutOfRange)

	// Lower boundary
	result1 := filter(Right[string](0))()
	assert.Equal(t, E.Right[string](0), result1)

	// Upper boundary
	result2 := filter(Right[string](100))()
	assert.Equal(t, E.Right[string](100), result2)

	// Below lower boundary
	result3 := filter(Right[string](-1))()
	assert.Equal(t, E.Left[int]("-1 is out of range [0, 100]"), result3)

	// Above upper boundary
	result4 := filter(Right[string](101))()
	assert.Equal(t, E.Left[int]("101 is out of range [0, 100]"), result4)
}

func TestFilterOrElse_AlwaysTrue(t *testing.T) {
	// Test with predicate that always returns true
	alwaysTrue := func(n int) bool { return true }
	onFalse := func(n int) string { return "never happens" }

	filter := FilterOrElse(alwaysTrue, onFalse)

	result1 := filter(Right[string](42))()
	assert.Equal(t, E.Right[string](42), result1)

	result2 := filter(Right[string](-42))()
	assert.Equal(t, E.Right[string](-42), result2)
}

func TestFilterOrElse_AlwaysFalse(t *testing.T) {
	// Test with predicate that always returns false
	alwaysFalse := func(n int) bool { return false }
	onFalse := S.Format[int]("rejected: %d")

	filter := FilterOrElse(alwaysFalse, onFalse)

	result := filter(Right[string](42))()
	assert.Equal(t, E.Left[int]("rejected: 42"), result)
}

func TestFilterOrElse_NilPointerValidation(t *testing.T) {
	// Test filtering nil pointers
	isNonNil := func(p *int) bool { return p != nil }
	onNil := func(p *int) string { return "pointer is nil" }

	filter := FilterOrElse(isNonNil, onNil)

	// Non-nil pointer passes
	value := 42
	result1 := filter(Right[string](&value))()
	assert.True(t, E.IsRight(result1))

	// Nil pointer becomes error
	result2 := filter(Right[string]((*int)(nil)))()
	assert.Equal(t, E.Left[*int]("pointer is nil"), result2)
}
