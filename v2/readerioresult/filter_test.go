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

package readerioresult

import (
	"context"
	"fmt"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	N "github.com/IBM/fp-go/v2/number"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

type Config struct {
	MaxValue int
	MinValue int
}

func TestFilterOrElse_PredicateTrue(t *testing.T) {
	// Test that when predicate returns true, Right value passes through
	isPositive := N.MoreThan(0)
	onFalse := func(n int) error { return fmt.Errorf("%d is not positive", n) }

	filter := FilterOrElse[context.Context](isPositive, onFalse)
	result := filter(Right[context.Context](42))(context.Background())()

	assert.Equal(t, E.Right[error](42), result)
}

func TestFilterOrElse_PredicateFalse(t *testing.T) {
	// Test that when predicate returns false, Right value becomes Left
	isPositive := N.MoreThan(0)
	onFalse := func(n int) error { return fmt.Errorf("%d is not positive", n) }

	filter := FilterOrElse[context.Context](isPositive, onFalse)
	result := filter(Right[context.Context](-5))(context.Background())()

	assert.True(t, E.IsLeft(result))
	assert.Equal(t, "-5 is not positive", E.ToError(result).Error())
}

func TestFilterOrElse_LeftPassesThrough(t *testing.T) {
	// Test that Left values pass through unchanged
	isPositive := N.MoreThan(0)
	onFalse := func(n int) error { return fmt.Errorf("%d is not positive", n) }

	filter := FilterOrElse[context.Context](isPositive, onFalse)
	originalErr := fmt.Errorf("original error")
	result := filter(Left[context.Context, int](originalErr))(context.Background())()

	assert.True(t, E.IsLeft(result))
	assert.Equal(t, "original error", E.ToError(result).Error())
}

func TestFilterOrElse_WithContext(t *testing.T) {
	// Test filtering with context-dependent validation
	cfg := Config{MaxValue: 100, MinValue: 0}

	isInRange := func(n int) bool { return n >= cfg.MinValue && n <= cfg.MaxValue }
	onOutOfRange := func(n int) error {
		return fmt.Errorf("%d is out of range [%d, %d]", n, cfg.MinValue, cfg.MaxValue)
	}

	filter := FilterOrElse[Config](isInRange, onOutOfRange)

	// Within range
	result1 := filter(Right[Config](50))(cfg)()
	assert.Equal(t, E.Right[error](50), result1)

	// Below range
	result2 := filter(Right[Config](-10))(cfg)()
	assert.True(t, E.IsLeft(result2))
	assert.Equal(t, "-10 is out of range [0, 100]", E.ToError(result2).Error())

	// Above range
	result3 := filter(Right[Config](150))(cfg)()
	assert.True(t, E.IsLeft(result3))
	assert.Equal(t, "150 is out of range [0, 100]", E.ToError(result3).Error())
}

func TestFilterOrElse_ZeroValue(t *testing.T) {
	// Test filtering with zero value
	isNonZero := func(n int) bool { return n != 0 }
	onZero := func(n int) error { return fmt.Errorf("value is zero") }

	filter := FilterOrElse[context.Context](isNonZero, onZero)
	result := filter(Right[context.Context](0))(context.Background())()

	assert.True(t, E.IsLeft(result))
	assert.Equal(t, "value is zero", E.ToError(result).Error())
}

func TestFilterOrElse_StringValidation(t *testing.T) {
	// Test with string validation
	isNonEmpty := S.IsNonEmpty
	onEmpty := func(s string) error { return fmt.Errorf("string is empty") }

	filter := FilterOrElse[context.Context](isNonEmpty, onEmpty)

	// Non-empty string passes
	result1 := filter(Right[context.Context]("hello"))(context.Background())()
	assert.Equal(t, E.Right[error]("hello"), result1)

	// Empty string becomes error
	result2 := filter(Right[context.Context](""))(context.Background())()
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
	onMinor := func(u User) error {
		return fmt.Errorf("%s is only %d years old", u.Name, u.Age)
	}

	filter := FilterOrElse[context.Context](isAdult, onMinor)

	// Adult user passes
	adult := User{Name: "Alice", Age: 25}
	result1 := filter(Right[context.Context](adult))(context.Background())()
	assert.Equal(t, E.Right[error](adult), result1)

	// Minor becomes error
	minor := User{Name: "Bob", Age: 16}
	result2 := filter(Right[context.Context](minor))(context.Background())()
	assert.True(t, E.IsLeft(result2))
	assert.Equal(t, "Bob is only 16 years old", E.ToError(result2).Error())
}

func TestFilterOrElse_ChainedFilters(t *testing.T) {
	// Test chaining multiple filters
	isPositive := N.MoreThan(0)
	onNegative := func(n int) error { return fmt.Errorf("not positive") }

	isEven := func(n int) bool { return n%2 == 0 }
	onOdd := func(n int) error { return fmt.Errorf("not even") }

	filter1 := FilterOrElse[context.Context](isPositive, onNegative)
	filter2 := FilterOrElse[context.Context](isEven, onOdd)

	ctx := context.Background()

	// Chain filters - apply filter1 first, then filter2
	result := filter2(filter1(Right[context.Context](4)))(ctx)()
	assert.Equal(t, E.Right[error](4), result)

	// Fails first filter
	result2 := filter2(filter1(Right[context.Context](-2)))(ctx)()
	assert.True(t, E.IsLeft(result2))
	assert.Equal(t, "not positive", E.ToError(result2).Error())

	// Passes first but fails second
	result3 := filter2(filter1(Right[context.Context](3)))(ctx)()
	assert.True(t, E.IsLeft(result3))
	assert.Equal(t, "not even", E.ToError(result3).Error())
}

func TestFilterOrElse_WithMap(t *testing.T) {
	// Test FilterOrElse combined with Map
	isPositive := N.MoreThan(0)
	onNegative := func(n int) error { return fmt.Errorf("negative number") }

	filter := FilterOrElse[context.Context](isPositive, onNegative)
	double := Map[context.Context](N.Mul(2))

	ctx := context.Background()

	// Compose: filter then double
	result1 := double(filter(Right[context.Context](5)))(ctx)()
	assert.Equal(t, E.Right[error](10), result1)

	// Negative value filtered out
	result2 := double(filter(Right[context.Context](-3)))(ctx)()
	assert.True(t, E.IsLeft(result2))
	assert.Equal(t, "negative number", E.ToError(result2).Error())
}

func TestFilterOrElse_BoundaryConditions(t *testing.T) {
	// Test boundary conditions
	isInRange := func(n int) bool { return n >= 0 && n <= 100 }
	onOutOfRange := func(n int) error {
		return fmt.Errorf("%d is out of range [0, 100]", n)
	}

	filter := FilterOrElse[context.Context](isInRange, onOutOfRange)
	ctx := context.Background()

	// Lower boundary
	result1 := filter(Right[context.Context](0))(ctx)()
	assert.Equal(t, E.Right[error](0), result1)

	// Upper boundary
	result2 := filter(Right[context.Context](100))(ctx)()
	assert.Equal(t, E.Right[error](100), result2)

	// Below lower boundary
	result3 := filter(Right[context.Context](-1))(ctx)()
	assert.True(t, E.IsLeft(result3))
	assert.Equal(t, "-1 is out of range [0, 100]", E.ToError(result3).Error())

	// Above upper boundary
	result4 := filter(Right[context.Context](101))(ctx)()
	assert.True(t, E.IsLeft(result4))
	assert.Equal(t, "101 is out of range [0, 100]", E.ToError(result4).Error())
}

func TestFilterOrElse_AlwaysTrue(t *testing.T) {
	// Test with predicate that always returns true
	alwaysTrue := func(n int) bool { return true }
	onFalse := func(n int) error { return fmt.Errorf("never happens") }

	filter := FilterOrElse[context.Context](alwaysTrue, onFalse)
	ctx := context.Background()

	result1 := filter(Right[context.Context](42))(ctx)()
	assert.Equal(t, E.Right[error](42), result1)

	result2 := filter(Right[context.Context](-42))(ctx)()
	assert.Equal(t, E.Right[error](-42), result2)
}

func TestFilterOrElse_AlwaysFalse(t *testing.T) {
	// Test with predicate that always returns false
	alwaysFalse := func(n int) bool { return false }
	onFalse := func(n int) error { return fmt.Errorf("rejected: %d", n) }

	filter := FilterOrElse[context.Context](alwaysFalse, onFalse)
	ctx := context.Background()

	result := filter(Right[context.Context](42))(ctx)()
	assert.True(t, E.IsLeft(result))
	assert.Equal(t, "rejected: 42", E.ToError(result).Error())
}

func TestFilterOrElse_NilPointerValidation(t *testing.T) {
	// Test filtering nil pointers
	isNonNil := func(p *int) bool { return p != nil }
	onNil := func(p *int) error { return fmt.Errorf("pointer is nil") }

	filter := FilterOrElse[context.Context](isNonNil, onNil)
	ctx := context.Background()

	// Non-nil pointer passes
	value := 42
	result1 := filter(Right[context.Context](&value))(ctx)()
	assert.True(t, E.IsRight(result1))

	// Nil pointer becomes error
	result2 := filter(Right[context.Context]((*int)(nil)))(ctx)()
	assert.True(t, E.IsLeft(result2))
	assert.Equal(t, "pointer is nil", E.ToError(result2).Error())
}

func TestFilterOrElse_ContextPropagation(t *testing.T) {
	// Test that context is properly propagated
	type ctxKey string
	const key ctxKey = "test-key"

	ctx := context.WithValue(context.Background(), key, "test-value")

	isPositive := N.MoreThan(0)
	onNegative := func(n int) error { return fmt.Errorf("negative") }

	filter := FilterOrElse[context.Context](isPositive, onNegative)

	// The context should be available when the computation runs
	result := filter(Right[context.Context](42))(ctx)()
	assert.Equal(t, E.Right[error](42), result)
}

func TestFilterOrElse_DifferentContextTypes(t *testing.T) {
	// Test with different context types
	type AppConfig struct {
		Name    string
		Version string
	}

	cfg := AppConfig{Name: "TestApp", Version: "1.0.0"}

	isValidVersion := func(v string) bool { return len(v) > 0 }
	onInvalid := func(v string) error { return fmt.Errorf("invalid version") }

	filter := FilterOrElse[AppConfig](isValidVersion, onInvalid)

	result := filter(Right[AppConfig]("1.0.0"))(cfg)()
	assert.Equal(t, E.Right[error]("1.0.0"), result)
}

func TestFilterOrElse_ErrorWrapping(t *testing.T) {
	// Test that errors are properly created with context
	isValid := func(n int) bool { return n >= 0 && n <= 100 }
	onInvalid := func(n int) error {
		return fmt.Errorf("validation failed: value %d is out of range", n)
	}

	filter := FilterOrElse[context.Context](isValid, onInvalid)
	ctx := context.Background()

	result := filter(Right[context.Context](150))(ctx)()
	assert.True(t, E.IsLeft(result))
	assert.Contains(t, E.ToError(result).Error(), "validation failed")
	assert.Contains(t, E.ToError(result).Error(), "150")
}
