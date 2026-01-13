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

package readerioeither

import (
	"context"
	"fmt"
	"testing"

	N "github.com/IBM/fp-go/v2/number"
	S "github.com/IBM/fp-go/v2/string"

	E "github.com/IBM/fp-go/v2/either"
	"github.com/stretchr/testify/assert"
)

type Config struct {
	MaxValue int
	MinValue int
}

func TestFilterOrElse_PredicateTrue(t *testing.T) {
	// Test that when predicate returns true, Right value passes through
	isPositive := N.MoreThan(0)
	onFalse := S.Format[int]("%d is not positive")

	filter := FilterOrElse[context.Context](isPositive, onFalse)
	result := filter(Right[context.Context, string](42))(context.Background())()

	assert.Equal(t, E.Right[string](42), result)
}

func TestFilterOrElse_PredicateFalse(t *testing.T) {
	// Test that when predicate returns false, Right value becomes Left
	isPositive := N.MoreThan(0)
	onFalse := S.Format[int]("%d is not positive")

	filter := FilterOrElse[context.Context](isPositive, onFalse)
	result := filter(Right[context.Context, string](-5))(context.Background())()

	assert.Equal(t, E.Left[int]("-5 is not positive"), result)
}

func TestFilterOrElse_LeftPassesThrough(t *testing.T) {
	// Test that Left values pass through unchanged
	isPositive := N.MoreThan(0)
	onFalse := S.Format[int]("%d is not positive")

	filter := FilterOrElse[context.Context](isPositive, onFalse)
	result := filter(Left[context.Context, int]("original error"))(context.Background())()

	assert.Equal(t, E.Left[int]("original error"), result)
}

func TestFilterOrElse_WithContext(t *testing.T) {
	// Test filtering with context-dependent validation
	cfg := Config{MaxValue: 100, MinValue: 0}

	isInRange := func(n int) bool { return n >= cfg.MinValue && n <= cfg.MaxValue }
	onOutOfRange := func(n int) string {
		return fmt.Sprintf("%d is out of range [%d, %d]", n, cfg.MinValue, cfg.MaxValue)
	}

	filter := FilterOrElse[Config](isInRange, onOutOfRange)

	// Within range
	result1 := filter(Right[Config, string](50))(cfg)()
	assert.Equal(t, E.Right[string](50), result1)

	// Below range
	result2 := filter(Right[Config, string](-10))(cfg)()
	assert.Equal(t, E.Left[int]("-10 is out of range [0, 100]"), result2)

	// Above range
	result3 := filter(Right[Config, string](150))(cfg)()
	assert.Equal(t, E.Left[int]("150 is out of range [0, 100]"), result3)
}

func TestFilterOrElse_ZeroValue(t *testing.T) {
	// Test filtering with zero value
	isNonZero := func(n int) bool { return n != 0 }
	onZero := func(n int) string { return "value is zero" }

	filter := FilterOrElse[context.Context](isNonZero, onZero)
	result := filter(Right[context.Context, string](0))(context.Background())()

	assert.Equal(t, E.Left[int]("value is zero"), result)
}

func TestFilterOrElse_StringValidation(t *testing.T) {
	// Test with string validation
	isNonEmpty := S.IsNonEmpty
	onEmpty := func(s string) error { return fmt.Errorf("string is empty") }

	filter := FilterOrElse[context.Context](isNonEmpty, onEmpty)

	// Non-empty string passes
	result1 := filter(Right[context.Context, error]("hello"))(context.Background())()
	assert.Equal(t, E.Right[error]("hello"), result1)

	// Empty string becomes error
	result2 := filter(Right[context.Context, error](""))(context.Background())()
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

	filter := FilterOrElse[context.Context](isAdult, onMinor)

	// Adult user passes
	adult := User{Name: "Alice", Age: 25}
	result1 := filter(Right[context.Context, string](adult))(context.Background())()
	assert.Equal(t, E.Right[string](adult), result1)

	// Minor becomes error
	minor := User{Name: "Bob", Age: 16}
	result2 := filter(Right[context.Context, string](minor))(context.Background())()
	assert.Equal(t, E.Left[User]("Bob is only 16 years old"), result2)
}

func TestFilterOrElse_ChainedFilters(t *testing.T) {
	// Test chaining multiple filters
	isPositive := N.MoreThan(0)
	onNegative := func(n int) string { return "not positive" }

	isEven := func(n int) bool { return n%2 == 0 }
	onOdd := func(n int) string { return "not even" }

	filter1 := FilterOrElse[context.Context](isPositive, onNegative)
	filter2 := FilterOrElse[context.Context](isEven, onOdd)

	ctx := context.Background()

	// Chain filters - apply filter1 first, then filter2
	result := filter2(filter1(Right[context.Context, string](4)))(ctx)()
	assert.Equal(t, E.Right[string](4), result)

	// Fails first filter
	result2 := filter2(filter1(Right[context.Context, string](-2)))(ctx)()
	assert.Equal(t, E.Left[int]("not positive"), result2)

	// Passes first but fails second
	result3 := filter2(filter1(Right[context.Context, string](3)))(ctx)()
	assert.Equal(t, E.Left[int]("not even"), result3)
}

func TestFilterOrElse_WithMap(t *testing.T) {
	// Test FilterOrElse combined with Map
	isPositive := N.MoreThan(0)
	onNegative := func(n int) string { return "negative number" }

	filter := FilterOrElse[context.Context](isPositive, onNegative)
	double := Map[context.Context, string](N.Mul(2))

	ctx := context.Background()

	// Compose: filter then double
	result1 := double(filter(Right[context.Context, string](5)))(ctx)()
	assert.Equal(t, E.Right[string](10), result1)

	// Negative value filtered out
	result2 := double(filter(Right[context.Context, string](-3)))(ctx)()
	assert.Equal(t, E.Left[int]("negative number"), result2)
}

func TestFilterOrElse_BoundaryConditions(t *testing.T) {
	// Test boundary conditions
	isInRange := func(n int) bool { return n >= 0 && n <= 100 }
	onOutOfRange := func(n int) string {
		return fmt.Sprintf("%d is out of range [0, 100]", n)
	}

	filter := FilterOrElse[context.Context](isInRange, onOutOfRange)
	ctx := context.Background()

	// Lower boundary
	result1 := filter(Right[context.Context, string](0))(ctx)()
	assert.Equal(t, E.Right[string](0), result1)

	// Upper boundary
	result2 := filter(Right[context.Context, string](100))(ctx)()
	assert.Equal(t, E.Right[string](100), result2)

	// Below lower boundary
	result3 := filter(Right[context.Context, string](-1))(ctx)()
	assert.Equal(t, E.Left[int]("-1 is out of range [0, 100]"), result3)

	// Above upper boundary
	result4 := filter(Right[context.Context, string](101))(ctx)()
	assert.Equal(t, E.Left[int]("101 is out of range [0, 100]"), result4)
}

func TestFilterOrElse_AlwaysTrue(t *testing.T) {
	// Test with predicate that always returns true
	alwaysTrue := func(n int) bool { return true }
	onFalse := func(n int) string { return "never happens" }

	filter := FilterOrElse[context.Context](alwaysTrue, onFalse)
	ctx := context.Background()

	result1 := filter(Right[context.Context, string](42))(ctx)()
	assert.Equal(t, E.Right[string](42), result1)

	result2 := filter(Right[context.Context, string](-42))(ctx)()
	assert.Equal(t, E.Right[string](-42), result2)
}

func TestFilterOrElse_AlwaysFalse(t *testing.T) {
	// Test with predicate that always returns false
	alwaysFalse := func(n int) bool { return false }
	onFalse := S.Format[int]("rejected: %d")

	filter := FilterOrElse[context.Context](alwaysFalse, onFalse)
	ctx := context.Background()

	result := filter(Right[context.Context, string](42))(ctx)()
	assert.Equal(t, E.Left[int]("rejected: 42"), result)
}

func TestFilterOrElse_NilPointerValidation(t *testing.T) {
	// Test filtering nil pointers
	isNonNil := func(p *int) bool { return p != nil }
	onNil := func(p *int) string { return "pointer is nil" }

	filter := FilterOrElse[context.Context](isNonNil, onNil)
	ctx := context.Background()

	// Non-nil pointer passes
	value := 42
	result1 := filter(Right[context.Context, string](&value))(ctx)()
	assert.True(t, E.IsRight(result1))

	// Nil pointer becomes error
	result2 := filter(Right[context.Context, string]((*int)(nil)))(ctx)()
	assert.Equal(t, E.Left[*int]("pointer is nil"), result2)
}

func TestFilterOrElse_ContextPropagation(t *testing.T) {
	// Test that context is properly propagated
	type ctxKey string
	const key ctxKey = "test-key"

	ctx := context.WithValue(context.Background(), key, "test-value")

	isPositive := N.MoreThan(0)
	onNegative := func(n int) string { return "negative" }

	filter := FilterOrElse[context.Context](isPositive, onNegative)

	// The context should be available when the computation runs
	result := filter(Right[context.Context, string](42))(ctx)()
	assert.Equal(t, E.Right[string](42), result)
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

	result := filter(Right[AppConfig, error]("1.0.0"))(cfg)()
	assert.Equal(t, E.Right[error]("1.0.0"), result)
}
