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

package boolean

import (
	"fmt"
	"testing"

	"github.com/IBM/fp-go/v2/lazy"
	"github.com/stretchr/testify/assert"
)

func TestFold_True(t *testing.T) {
	t.Run("returns onTrue result when given true", func(t *testing.T) {
		fold := Fold(lazy.Of("no"), lazy.Of("yes"))
		assert.Equal(t, "yes", fold(true))
	})

	t.Run("returns onTrue int when given true", func(t *testing.T) {
		fold := Fold(lazy.Of(0), lazy.Of(1))
		assert.Equal(t, 1, fold(true))
	})
}

func TestFold_False(t *testing.T) {
	t.Run("returns onFalse result when given false", func(t *testing.T) {
		fold := Fold(lazy.Of("no"), lazy.Of("yes"))
		assert.Equal(t, "no", fold(false))
	})

	t.Run("returns onFalse int when given false", func(t *testing.T) {
		fold := Fold(lazy.Of(0), lazy.Of(1))
		assert.Equal(t, 0, fold(false))
	})
}

func TestFold_LazyEvaluation(t *testing.T) {
	t.Run("onTrue thunk is not called when input is false", func(t *testing.T) {
		called := false
		onTrue := func() string { called = true; return "yes" }
		Fold(lazy.Of("no"), onTrue)(false)
		assert.False(t, called, "onTrue should not have been called")
	})

	t.Run("onFalse thunk is not called when input is true", func(t *testing.T) {
		called := false
		onFalse := func() string { called = true; return "no" }
		Fold(onFalse, lazy.Of("yes"))(true)
		assert.False(t, called, "onFalse should not have been called")
	})
}

func TestFold_ReusableCurried(t *testing.T) {
	t.Run("same fold function applied to both values", func(t *testing.T) {
		toLabel := Fold(lazy.Of("disabled"), lazy.Of("enabled"))
		assert.Equal(t, "enabled", toLabel(true))
		assert.Equal(t, "disabled", toLabel(false))
	})
}

func TestFold_EdgeCases(t *testing.T) {
	t.Run("works with struct result type", func(t *testing.T) {
		type Status struct{ Active bool }
		fold := Fold(lazy.Of(Status{false}), lazy.Of(Status{true}))
		assert.Equal(t, Status{true}, fold(true))
		assert.Equal(t, Status{false}, fold(false))
	})

	t.Run("works with slice result type", func(t *testing.T) {
		fold := Fold(lazy.Of([]int{0}), lazy.Of([]int{1, 2}))
		assert.Equal(t, []int{1, 2}, fold(true))
		assert.Equal(t, []int{0}, fold(false))
	})
}

// ExampleFold demonstrates selecting a string based on a boolean condition.
func ExampleFold_string() {
	label := Fold(lazy.Of("inactive"), lazy.Of("active"))
	fmt.Println(label(true))
	fmt.Println(label(false))
	// Output:
	// active
	// inactive
}

// ExampleFold_int demonstrates mapping a boolean to integer values.
func ExampleFold_int() {
	score := Fold(lazy.Of(0), lazy.Of(100))
	fmt.Println(score(true))
	fmt.Println(score(false))
	// Output:
	// 100
	// 0
}

// ExampleFold_lazy demonstrates that thunks are evaluated lazily — the
// unchosen branch is never called.
func ExampleFold_lazy() {
	expensive := func() string {
		// In real code this might be a costly computation; it is only
		// called when the boolean is false.
		return "fallback"
	}
	result := Fold(expensive, lazy.Of("fast-path"))(true)
	fmt.Println(result)
	// Output:
	// fast-path
}
