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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMonoidAny(t *testing.T) {
	t.Run("identity element is false", func(t *testing.T) {
		assert.Equal(t, false, MonoidAny.Empty())
	})

	t.Run("false OR false = false", func(t *testing.T) {
		result := MonoidAny.Concat(false, false)
		assert.Equal(t, false, result)
	})

	t.Run("false OR true = true", func(t *testing.T) {
		result := MonoidAny.Concat(false, true)
		assert.Equal(t, true, result)
	})

	t.Run("true OR false = true", func(t *testing.T) {
		result := MonoidAny.Concat(true, false)
		assert.Equal(t, true, result)
	})

	t.Run("true OR true = true", func(t *testing.T) {
		result := MonoidAny.Concat(true, true)
		assert.Equal(t, true, result)
	})

	t.Run("left identity: empty OR x = x", func(t *testing.T) {
		assert.Equal(t, true, MonoidAny.Concat(MonoidAny.Empty(), true))
		assert.Equal(t, false, MonoidAny.Concat(MonoidAny.Empty(), false))
	})

	t.Run("right identity: x OR empty = x", func(t *testing.T) {
		assert.Equal(t, true, MonoidAny.Concat(true, MonoidAny.Empty()))
		assert.Equal(t, false, MonoidAny.Concat(false, MonoidAny.Empty()))
	})

	t.Run("associativity: (a OR b) OR c = a OR (b OR c)", func(t *testing.T) {
		a, b, c := true, false, true
		left := MonoidAny.Concat(MonoidAny.Concat(a, b), c)
		right := MonoidAny.Concat(a, MonoidAny.Concat(b, c))
		assert.Equal(t, left, right)
	})
}

func TestMonoidAll(t *testing.T) {
	t.Run("identity element is true", func(t *testing.T) {
		assert.Equal(t, true, MonoidAll.Empty())
	})

	t.Run("false AND false = false", func(t *testing.T) {
		result := MonoidAll.Concat(false, false)
		assert.Equal(t, false, result)
	})

	t.Run("false AND true = false", func(t *testing.T) {
		result := MonoidAll.Concat(false, true)
		assert.Equal(t, false, result)
	})

	t.Run("true AND false = false", func(t *testing.T) {
		result := MonoidAll.Concat(true, false)
		assert.Equal(t, false, result)
	})

	t.Run("true AND true = true", func(t *testing.T) {
		result := MonoidAll.Concat(true, true)
		assert.Equal(t, true, result)
	})

	t.Run("left identity: empty AND x = x", func(t *testing.T) {
		assert.Equal(t, true, MonoidAll.Concat(MonoidAll.Empty(), true))
		assert.Equal(t, false, MonoidAll.Concat(MonoidAll.Empty(), false))
	})

	t.Run("right identity: x AND empty = x", func(t *testing.T) {
		assert.Equal(t, true, MonoidAll.Concat(true, MonoidAll.Empty()))
		assert.Equal(t, false, MonoidAll.Concat(false, MonoidAll.Empty()))
	})

	t.Run("associativity: (a AND b) AND c = a AND (b AND c)", func(t *testing.T) {
		a, b, c := true, false, true
		left := MonoidAll.Concat(MonoidAll.Concat(a, b), c)
		right := MonoidAll.Concat(a, MonoidAll.Concat(b, c))
		assert.Equal(t, left, right)
	})
}

func TestEq(t *testing.T) {
	t.Run("true equals true", func(t *testing.T) {
		assert.True(t, Eq.Equals(true, true))
	})

	t.Run("false equals false", func(t *testing.T) {
		assert.True(t, Eq.Equals(false, false))
	})

	t.Run("true not equals false", func(t *testing.T) {
		assert.False(t, Eq.Equals(true, false))
	})

	t.Run("false not equals true", func(t *testing.T) {
		assert.False(t, Eq.Equals(false, true))
	})

	t.Run("reflexivity: x equals x", func(t *testing.T) {
		assert.True(t, Eq.Equals(true, true))
		assert.True(t, Eq.Equals(false, false))
	})

	t.Run("symmetry: if x equals y then y equals x", func(t *testing.T) {
		assert.Equal(t, Eq.Equals(true, false), Eq.Equals(false, true))
		assert.Equal(t, Eq.Equals(true, true), Eq.Equals(true, true))
	})

	t.Run("transitivity: if x equals y and y equals z then x equals z", func(t *testing.T) {
		x, y, z := true, true, true
		if Eq.Equals(x, y) && Eq.Equals(y, z) {
			assert.True(t, Eq.Equals(x, z))
		}
	})
}

func TestOrd(t *testing.T) {
	t.Run("false < true", func(t *testing.T) {
		result := Ord.Compare(false, true)
		assert.Equal(t, -1, result)
	})

	t.Run("true > false", func(t *testing.T) {
		result := Ord.Compare(true, false)
		assert.Equal(t, 1, result)
	})

	t.Run("true == true", func(t *testing.T) {
		result := Ord.Compare(true, true)
		assert.Equal(t, 0, result)
	})

	t.Run("false == false", func(t *testing.T) {
		result := Ord.Compare(false, false)
		assert.Equal(t, 0, result)
	})

	t.Run("Equals method works", func(t *testing.T) {
		assert.True(t, Ord.Equals(true, true))
		assert.True(t, Ord.Equals(false, false))
		assert.False(t, Ord.Equals(true, false))
		assert.False(t, Ord.Equals(false, true))
	})

	t.Run("reflexivity: x <= x", func(t *testing.T) {
		assert.True(t, Ord.Compare(true, true) == 0)
		assert.True(t, Ord.Compare(false, false) == 0)
	})

	t.Run("antisymmetry: if x <= y and y <= x then x == y", func(t *testing.T) {
		x, y := true, true
		if Ord.Compare(x, y) <= 0 && Ord.Compare(y, x) <= 0 {
			assert.Equal(t, 0, Ord.Compare(x, y))
		}
	})

	t.Run("transitivity: if x <= y and y <= z then x <= z", func(t *testing.T) {
		x, y, z := false, true, true
		if Ord.Compare(x, y) <= 0 && Ord.Compare(y, z) <= 0 {
			assert.True(t, Ord.Compare(x, z) <= 0)
		}
	})

	t.Run("totality: x <= y or y <= x", func(t *testing.T) {
		assert.True(t, Ord.Compare(true, false) >= 0 || Ord.Compare(false, true) >= 0)
		assert.True(t, Ord.Compare(false, true) <= 0 || Ord.Compare(true, false) <= 0)
	})
}

// Example tests that also serve as documentation
func ExampleMonoidAny() {
	// Combine booleans with OR
	result := MonoidAny.Concat(false, true)
	println(result) // true

	// Identity element
	identity := MonoidAny.Empty()
	println(identity) // false

	// Output:
}

func ExampleMonoidAll() {
	// Combine booleans with AND
	result := MonoidAll.Concat(true, true)
	println(result) // true

	// Identity element
	identity := MonoidAll.Empty()
	println(identity) // true

	// Output:
}

func ExampleEq() {
	// Check equality
	equal := Eq.Equals(true, true)
	println(equal) // true

	notEqual := Eq.Equals(true, false)
	println(notEqual) // false

	// Output:
}

func ExampleOrd() {
	// Compare booleans (false < true)
	cmp := Ord.Compare(false, true)
	println(cmp) // -1

	cmp2 := Ord.Compare(true, false)
	println(cmp2) // 1

	cmp3 := Ord.Compare(true, true)
	println(cmp3) // 0

	// Output:
}
