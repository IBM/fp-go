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

package string

import (
	"testing"

	M "github.com/IBM/fp-go/v2/monoid/testing"
	"github.com/stretchr/testify/assert"
)

func TestMonoid(t *testing.T) {
	M.AssertLaws(t, Monoid)([]string{"", "a", "some value"})
}

func TestMonoidConcat(t *testing.T) {
	t.Run("basic concatenation", func(t *testing.T) {
		result := Monoid.Concat("hello", " world")
		assert.Equal(t, "hello world", result)
	})

	t.Run("empty identity", func(t *testing.T) {
		empty := Monoid.Empty()
		assert.Equal(t, "", empty)
	})

	t.Run("left identity", func(t *testing.T) {
		// empty • a = a
		result := Monoid.Concat(Monoid.Empty(), "test")
		assert.Equal(t, "test", result)
	})

	t.Run("right identity", func(t *testing.T) {
		// a • empty = a
		result := Monoid.Concat("test", Monoid.Empty())
		assert.Equal(t, "test", result)
	})

	t.Run("associativity", func(t *testing.T) {
		// (a • b) • c = a • (b • c)
		a, b, c := "foo", "bar", "baz"
		left := Monoid.Concat(Monoid.Concat(a, b), c)
		right := Monoid.Concat(a, Monoid.Concat(b, c))
		assert.Equal(t, left, right)
		assert.Equal(t, "foobarbaz", left)
	})
}

func TestIntersperseMonoid(t *testing.T) {
	// Test with comma separator
	commaMonoid := IntersperseMonoid(", ")
	M.AssertLaws(t, commaMonoid)([]string{"", "a", "b", "some value"})

	// Test with dash separator
	dashMonoid := IntersperseMonoid("-")
	M.AssertLaws(t, dashMonoid)([]string{"", "x", "y", "test"})
}

func TestIntersperseMonoidConcat(t *testing.T) {
	t.Run("comma separator", func(t *testing.T) {
		commaMonoid := IntersperseMonoid(", ")
		result := commaMonoid.Concat("a", "b")
		assert.Equal(t, "a, b", result)
	})

	t.Run("empty identity", func(t *testing.T) {
		commaMonoid := IntersperseMonoid(", ")
		empty := commaMonoid.Empty()
		assert.Equal(t, "", empty)
	})

	t.Run("left identity with separator", func(t *testing.T) {
		// empty • a = a (no separator added)
		commaMonoid := IntersperseMonoid(", ")
		result := commaMonoid.Concat(commaMonoid.Empty(), "test")
		assert.Equal(t, "test", result)
	})

	t.Run("right identity with separator", func(t *testing.T) {
		// a • empty = a (no separator added)
		commaMonoid := IntersperseMonoid(", ")
		result := commaMonoid.Concat("test", commaMonoid.Empty())
		assert.Equal(t, "test", result)
	})

	t.Run("associativity with separator", func(t *testing.T) {
		// (a • b) • c = a • (b • c)
		commaMonoid := IntersperseMonoid(", ")
		a, b, c := "x", "y", "z"
		left := commaMonoid.Concat(commaMonoid.Concat(a, b), c)
		right := commaMonoid.Concat(a, commaMonoid.Concat(b, c))
		assert.Equal(t, left, right)
		assert.Equal(t, "x, y, z", left)
	})

	t.Run("multiple separators", func(t *testing.T) {
		dashMonoid := IntersperseMonoid("-")
		result := dashMonoid.Concat("foo", "bar")
		assert.Equal(t, "foo-bar", result)

		spaceMonoid := IntersperseMonoid(" ")
		result = spaceMonoid.Concat("hello", "world")
		assert.Equal(t, "hello world", result)
	})
}
