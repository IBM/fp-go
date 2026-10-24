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

package tailrec

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetters_Land(t *testing.T) {
	tramp := Land[int]("done")

	t.Run("methods", func(t *testing.T) {
		assert.True(t, tramp.IsLand())
		assert.False(t, tramp.IsBounce())

		l, ok := tramp.GetLand()
		assert.True(t, ok)
		assert.Equal(t, "done", l)

		b, ok := tramp.GetBounce()
		assert.False(t, ok)
		assert.Equal(t, 0, b)
	})

	t.Run("standalone functions", func(t *testing.T) {
		assert.True(t, IsLand(tramp))
		assert.False(t, IsBounce(tramp))

		l, ok := GetLand(tramp)
		assert.True(t, ok)
		assert.Equal(t, "done", l)

		b, ok := GetBounce(tramp)
		assert.False(t, ok)
		assert.Equal(t, 0, b)
	})
}

func TestGetters_Bounce(t *testing.T) {
	tramp := Bounce[string](42)

	t.Run("methods", func(t *testing.T) {
		assert.False(t, tramp.IsLand())
		assert.True(t, tramp.IsBounce())

		b, ok := tramp.GetBounce()
		assert.True(t, ok)
		assert.Equal(t, 42, b)

		l, ok := tramp.GetLand()
		assert.False(t, ok)
		assert.Equal(t, "", l)
	})

	t.Run("standalone functions", func(t *testing.T) {
		assert.False(t, IsLand(tramp))
		assert.True(t, IsBounce(tramp))

		b, ok := GetBounce(tramp)
		assert.True(t, ok)
		assert.Equal(t, 42, b)

		l, ok := GetLand(tramp)
		assert.False(t, ok)
		assert.Equal(t, "", l)
	})
}

func TestGetters_EdgeCases(t *testing.T) {
	t.Run("zero value Trampoline bounces with zero state", func(t *testing.T) {
		var tramp Trampoline[int, string]
		assert.True(t, IsBounce(tramp))
		b, ok := GetBounce(tramp)
		assert.True(t, ok)
		assert.Equal(t, 0, b)
	})

	t.Run("landed with zero value is still landed", func(t *testing.T) {
		l, ok := GetLand(Land[int](0))
		assert.True(t, ok)
		assert.Equal(t, 0, l)
	})

	t.Run("inactive field values are not exposed", func(t *testing.T) {
		l, ok := GetLand(Trampoline[int, string]{Land: "stale", Bounce: 7})
		assert.False(t, ok)
		assert.Equal(t, "", l)

		b, ok := GetBounce(Trampoline[int, string]{Land: "x", Bounce: 7, Landed: true})
		assert.False(t, ok)
		assert.Equal(t, 0, b)
	})
}

func TestGetters_Integration(t *testing.T) {
	t.Run("drive a factorial trampoline with the getters", func(t *testing.T) {
		type state struct{ n, acc int }
		step := func(s state) Trampoline[state, int] {
			if s.n <= 1 {
				return Land[state](s.acc)
			}
			return Bounce[int](state{s.n - 1, s.acc * s.n})
		}

		current := step(state{5, 1})
		for current.IsBounce() {
			s, _ := current.GetBounce()
			current = step(s)
		}
		result, ok := current.GetLand()
		assert.True(t, ok)
		assert.Equal(t, 120, result)
	})
}
