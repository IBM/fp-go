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

package option

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

func TestWrap_Success(t *testing.T) {
	t.Run("wraps a present value into Some", func(t *testing.T) {
		assert.Equal(t, Some(42), Wrap(42, true))
	})

	t.Run("wraps a present zero value into Some", func(t *testing.T) {
		assert.Equal(t, Some(0), Wrap(0, true))
		assert.Equal(t, Some(""), Wrap("", true))
	})

	t.Run("wraps a present nil pointer into Some", func(t *testing.T) {
		result := Wrap[*int](nil, true)
		assert.True(t, IsSome(result))
	})
}

func TestWrap_Failure(t *testing.T) {
	t.Run("returns None when ok is false", func(t *testing.T) {
		assert.Equal(t, None[int](), Wrap(0, false))
	})

	t.Run("discards the value when ok is false", func(t *testing.T) {
		assert.Equal(t, None[int](), Wrap(42, false))
		assert.Equal(t, None[string](), Wrap("ignored", false))
	})
}

func TestWrap_EdgeCases(t *testing.T) {
	t.Run("works with struct types", func(t *testing.T) {
		type point struct{ X, Y int }
		assert.Equal(t, Some(point{1, 2}), Wrap(point{1, 2}, true))
		assert.Equal(t, None[point](), Wrap(point{1, 2}, false))
	})

	t.Run("works with nested options", func(t *testing.T) {
		assert.Equal(t, Some(None[int]()), Wrap(None[int](), true))
		assert.Equal(t, None[Option[int]](), Wrap(Some(1), false))
	})
}

func TestWrap_Integration(t *testing.T) {
	t.Run("is the inverse of Unwrap", func(t *testing.T) {
		assert.Equal(t, Some(42), Wrap(Unwrap(Some(42))))
		assert.Equal(t, None[int](), Wrap(Unwrap(None[int]())))
	})

	t.Run("Unwrap is the inverse of Wrap", func(t *testing.T) {
		val, ok := Unwrap(Wrap(42, true))
		assert.True(t, ok)
		assert.Equal(t, 42, val)

		val, ok = Unwrap(Wrap(42, false))
		assert.False(t, ok)
		assert.Equal(t, 0, val)
	})

	t.Run("accepts a comma-ok call directly", func(t *testing.T) {
		assert.Equal(t, Some("1.2.3"), Wrap(strings.CutPrefix("v1.2.3", "v")))
		assert.Equal(t, None[string](), Wrap(strings.CutPrefix("1.2.3", "v")))
	})

	t.Run("composes with map lookups", func(t *testing.T) {
		m := map[string]int{"a": 1}
		lookup := func(k string) Option[int] {
			v, ok := m[k]
			return Wrap(v, ok)
		}
		assert.Equal(t, Some(1), lookup("a"))
		assert.Equal(t, None[int](), lookup("b"))
	})

	t.Run("composes with Map", func(t *testing.T) {
		result := F.Pipe1(Wrap(42, true), Map(strconv.Itoa))
		assert.Equal(t, Some("42"), result)
	})

	t.Run("agrees with Optionize0", func(t *testing.T) {
		for _, ok := range []bool{true, false} {
			expected := Optionize0(func() (int, bool) { return 7, ok })()
			assert.Equal(t, expected, Wrap(7, ok))
		}
	})
}

func BenchmarkWrap(b *testing.B) {
	for i := range b.N {
		_ = Wrap(i, i%2 == 0)
	}
}

// ExampleWrap demonstrates creating an Option from a value and a presence flag.
func ExampleWrap() {
	fmt.Println(Wrap(42, true))
	fmt.Println(Wrap(42, false))
	// Output:
	// Some[int](42)
	// None[int]
}

// ExampleWrap_commaOk demonstrates applying Wrap directly to a function that
// returns the comma-ok tuple (A, bool).
func ExampleWrap_commaOk() {
	fmt.Println(Wrap(strings.CutPrefix("v1.2.3", "v")))
	fmt.Println(Wrap(strings.CutPrefix("1.2.3", "v")))
	// Output:
	// Some[string](1.2.3)
	// None[string]
}

// ExampleWrap_unwrap demonstrates that Wrap is the inverse of Unwrap.
func ExampleWrap_unwrap() {
	fmt.Println(Wrap(Unwrap(Some("value"))))
	fmt.Println(Wrap(Unwrap(None[string]())))
	// Output:
	// Some[string](value)
	// None[string]
}
