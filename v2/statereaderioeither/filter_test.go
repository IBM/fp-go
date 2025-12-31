// Copyright (c) 2024 - 2025 IBM Corp.
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

package statereaderioeither

import (
	"fmt"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	N "github.com/IBM/fp-go/v2/number"
	P "github.com/IBM/fp-go/v2/pair"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

type AppState struct {
	Counter int
}

type Config struct {
	MaxValue int
}

func TestFilterOrElse_PredicateTrue(t *testing.T) {
	// Test that when predicate returns true, Right value passes through
	isPositive := N.MoreThan(0)
	onFalse := S.Format[int]("%d is not positive")

	filter := FilterOrElse[AppState, Config](isPositive, onFalse)
	result := filter(Right[AppState, Config, string](42))(AppState{Counter: 0})(Config{MaxValue: 100})()

	assert.True(t, E.IsRight(result))
	E.Map[string](func(p P.Pair[AppState, int]) P.Pair[AppState, int] {
		assert.Equal(t, 42, P.Tail(p))
		return p
	})(result)
}

func TestFilterOrElse_PredicateFalse(t *testing.T) {
	// Test that when predicate returns false, Right value becomes Left
	isPositive := N.MoreThan(0)
	onFalse := S.Format[int]("%d is not positive")

	filter := FilterOrElse[AppState, Config](isPositive, onFalse)
	result := filter(Right[AppState, Config, string](-5))(AppState{Counter: 0})(Config{MaxValue: 100})()

	assert.True(t, E.IsLeft(result))
	assert.Equal(t, E.Left[P.Pair[AppState, int]]("-5 is not positive"), result)
}

func TestFilterOrElse_LeftPassesThrough(t *testing.T) {
	// Test that Left values pass through unchanged
	isPositive := N.MoreThan(0)
	onFalse := S.Format[int]("%d is not positive")

	filter := FilterOrElse[AppState, Config](isPositive, onFalse)
	result := filter(Left[AppState, Config, int]("original error"))(AppState{Counter: 0})(Config{MaxValue: 100})()

	assert.True(t, E.IsLeft(result))
	assert.Equal(t, E.Left[P.Pair[AppState, int]]("original error"), result)
}

func TestFilterOrElse_StatePreserved(t *testing.T) {
	// Test that state is preserved through filtering
	isEven := func(n int) bool { return n%2 == 0 }
	onOdd := func(n int) string { return "not even" }

	filter := FilterOrElse[AppState, Config](isEven, onOdd)
	initialState := AppState{Counter: 5}
	result := filter(Right[AppState, Config, string](42))(initialState)(Config{MaxValue: 100})()

	assert.True(t, E.IsRight(result))
	E.Map[string](func(p P.Pair[AppState, int]) P.Pair[AppState, int] {
		assert.Equal(t, 42, P.Tail(p))
		assert.Equal(t, 5, P.Head(p).Counter) // State unchanged
		return p
	})(result)
}

func TestFilterOrElse_WithContextValidation(t *testing.T) {
	// Test filtering with context-dependent validation
	cfg := Config{MaxValue: 100}

	isInRange := func(n int) bool { return n <= cfg.MaxValue }
	onOutOfRange := func(n int) string {
		return fmt.Sprintf("%d exceeds max %d", n, cfg.MaxValue)
	}

	filter := FilterOrElse[AppState, Config](isInRange, onOutOfRange)

	// Within range
	result1 := filter(Right[AppState, Config, string](50))(AppState{Counter: 0})(cfg)()
	assert.True(t, E.IsRight(result1))

	// Above range
	result2 := filter(Right[AppState, Config, string](150))(AppState{Counter: 0})(cfg)()
	assert.True(t, E.IsLeft(result2))
	assert.Equal(t, E.Left[P.Pair[AppState, int]]("150 exceeds max 100"), result2)
}

func TestFilterOrElse_ChainedFilters(t *testing.T) {
	// Test chaining multiple filters
	isPositive := N.MoreThan(0)
	onNegative := func(n int) string { return "not positive" }

	isEven := func(n int) bool { return n%2 == 0 }
	onOdd := func(n int) string { return "not even" }

	filter1 := FilterOrElse[AppState, Config](isPositive, onNegative)
	filter2 := FilterOrElse[AppState, Config](isEven, onOdd)

	state := AppState{Counter: 0}
	cfg := Config{MaxValue: 100}

	// Both filters pass
	result := filter2(filter1(Right[AppState, Config, string](4)))(state)(cfg)()
	assert.True(t, E.IsRight(result))

	// Fails first filter
	result2 := filter2(filter1(Right[AppState, Config, string](-2)))(state)(cfg)()
	assert.True(t, E.IsLeft(result2))
	assert.Equal(t, E.Left[P.Pair[AppState, int]]("not positive"), result2)

	// Passes first but fails second
	result3 := filter2(filter1(Right[AppState, Config, string](3)))(state)(cfg)()
	assert.True(t, E.IsLeft(result3))
	assert.Equal(t, E.Left[P.Pair[AppState, int]]("not even"), result3)
}

func TestFilterOrElse_BoundaryConditions(t *testing.T) {
	// Test boundary conditions
	isInRange := func(n int) bool { return n >= 0 && n <= 100 }
	onOutOfRange := func(n int) string {
		return fmt.Sprintf("%d is out of range [0, 100]", n)
	}

	filter := FilterOrElse[AppState, Config](isInRange, onOutOfRange)
	state := AppState{Counter: 0}
	cfg := Config{MaxValue: 100}

	// Lower boundary
	result1 := filter(Right[AppState, Config, string](0))(state)(cfg)()
	assert.True(t, E.IsRight(result1))

	// Upper boundary
	result2 := filter(Right[AppState, Config, string](100))(state)(cfg)()
	assert.True(t, E.IsRight(result2))

	// Below lower boundary
	result3 := filter(Right[AppState, Config, string](-1))(state)(cfg)()
	assert.True(t, E.IsLeft(result3))

	// Above upper boundary
	result4 := filter(Right[AppState, Config, string](101))(state)(cfg)()
	assert.True(t, E.IsLeft(result4))
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

	filter := FilterOrElse[AppState, Config](isAdult, onMinor)
	state := AppState{Counter: 0}
	cfg := Config{MaxValue: 100}

	// Adult user passes
	adult := User{Name: "Alice", Age: 25}
	result1 := filter(Right[AppState, Config, string](adult))(state)(cfg)()
	assert.True(t, E.IsRight(result1))

	// Minor becomes error
	minor := User{Name: "Bob", Age: 16}
	result2 := filter(Right[AppState, Config, string](minor))(state)(cfg)()
	assert.True(t, E.IsLeft(result2))
	assert.Equal(t, E.Left[P.Pair[AppState, User]]("Bob is only 16 years old"), result2)
}

func TestFilterOrElse_AlwaysTrue(t *testing.T) {
	// Test with predicate that always returns true
	alwaysTrue := func(n int) bool { return true }
	onFalse := func(n int) string { return "never happens" }

	filter := FilterOrElse[AppState, Config](alwaysTrue, onFalse)
	state := AppState{Counter: 0}
	cfg := Config{MaxValue: 100}

	result1 := filter(Right[AppState, Config, string](42))(state)(cfg)()
	assert.True(t, E.IsRight(result1))

	result2 := filter(Right[AppState, Config, string](-42))(state)(cfg)()
	assert.True(t, E.IsRight(result2))
}

func TestFilterOrElse_AlwaysFalse(t *testing.T) {
	// Test with predicate that always returns false
	alwaysFalse := func(n int) bool { return false }
	onFalse := S.Format[int]("rejected: %d")

	filter := FilterOrElse[AppState, Config](alwaysFalse, onFalse)
	state := AppState{Counter: 0}
	cfg := Config{MaxValue: 100}

	result := filter(Right[AppState, Config, string](42))(state)(cfg)()
	assert.True(t, E.IsLeft(result))
	assert.Equal(t, E.Left[P.Pair[AppState, int]]("rejected: 42"), result)
}
