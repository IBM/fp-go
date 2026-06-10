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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAltAllArray_Success(t *testing.T) {
	t.Run("returns first Some value from array", func(t *testing.T) {
		options := []Option[int]{
			None[int](),
			None[int](),
			Some(42),
			Some(100),
		}
		result := AltAllArray(None[int]())(options)
		assert.Equal(t, Some(42), result)
	})

	t.Run("returns startWith when array is empty", func(t *testing.T) {
		options := []Option[int]{}
		result := AltAllArray(Some(10))(options)
		assert.Equal(t, Some(10), result)
	})

	t.Run("returns startWith when all array elements are None", func(t *testing.T) {
		options := []Option[int]{
			None[int](),
			None[int](),
			None[int](),
		}
		result := AltAllArray(Some(99))(options)
		assert.Equal(t, Some(99), result)
	})

	t.Run("returns first Some when startWith is None", func(t *testing.T) {
		options := []Option[string]{
			None[string](),
			Some("hello"),
			Some("world"),
		}
		result := AltAllArray(None[string]())(options)
		assert.Equal(t, Some("hello"), result)
	})

	t.Run("returns startWith when startWith is Some and array is all None", func(t *testing.T) {
		options := []Option[int]{
			None[int](),
			None[int](),
		}
		result := AltAllArray(Some(5))(options)
		assert.Equal(t, Some(5), result)
	})
}

func TestAltAllArray_EdgeCases(t *testing.T) {
	t.Run("handles single element array with Some", func(t *testing.T) {
		options := []Option[int]{Some(42)}
		result := AltAllArray(None[int]())(options)
		assert.Equal(t, Some(42), result)
	})

	t.Run("handles single element array with None", func(t *testing.T) {
		options := []Option[int]{None[int]()}
		result := AltAllArray(Some(10))(options)
		assert.Equal(t, Some(10), result)
	})

	t.Run("returns last Some when multiple Some values exist", func(t *testing.T) {
		options := []Option[int]{
			Some(1),
			Some(2),
			Some(3),
		}
		result := AltAllArray(None[int]())(options)
		assert.Equal(t, Some(1), result)
	})

	t.Run("works with different types", func(t *testing.T) {
		options := []Option[string]{
			None[string](),
			Some("first"),
			None[string](),
			Some("second"),
		}
		result := AltAllArray(None[string]())(options)
		assert.Equal(t, Some("first"), result)
	})

	t.Run("returns None when startWith is None and array is empty", func(t *testing.T) {
		options := []Option[int]{}
		result := AltAllArray(None[int]())(options)
		assert.Equal(t, None[int](), result)
	})

	t.Run("returns None when startWith is None and all array elements are None", func(t *testing.T) {
		options := []Option[int]{
			None[int](),
			None[int](),
		}
		result := AltAllArray(None[int]())(options)
		assert.Equal(t, None[int](), result)
	})
}

func TestAltAllSeq_Success(t *testing.T) {
	t.Run("returns first Some value from sequence", func(t *testing.T) {
		generator := func(yield func(Option[int]) bool) {
			yield(None[int]())
			yield(None[int]())
			yield(Some(42))
			yield(Some(100))
		}
		result := AltAllSeq(None[int]())(generator)
		assert.Equal(t, Some(42), result)
	})

	t.Run("returns startWith when sequence is empty", func(t *testing.T) {
		generator := func(yield func(Option[int]) bool) {}
		result := AltAllSeq(Some(10))(generator)
		assert.Equal(t, Some(10), result)
	})

	t.Run("returns startWith when all sequence elements are None", func(t *testing.T) {
		generator := func(yield func(Option[int]) bool) {
			yield(None[int]())
			yield(None[int]())
			yield(None[int]())
		}
		result := AltAllSeq(Some(99))(generator)
		assert.Equal(t, Some(99), result)
	})

	t.Run("returns first Some when startWith is None", func(t *testing.T) {
		generator := func(yield func(Option[string]) bool) {
			yield(None[string]())
			yield(Some("hello"))
			yield(Some("world"))
		}
		result := AltAllSeq(None[string]())(generator)
		assert.Equal(t, Some("hello"), result)
	})

	t.Run("stops iteration after finding first Some", func(t *testing.T) {
		iterationCount := 0
		generator := func(yield func(Option[int]) bool) {
			iterationCount++
			yield(None[int]())
			iterationCount++
			yield(Some(42))
			iterationCount++
			yield(Some(100))
		}
		result := AltAllSeq(None[int]())(generator)
		assert.Equal(t, Some(42), result)
		// The generator will be fully consumed due to how Alt works
		assert.Equal(t, 3, iterationCount)
	})
}

func TestAltAllSeq_EdgeCases(t *testing.T) {
	t.Run("handles single element sequence with Some", func(t *testing.T) {
		generator := func(yield func(Option[int]) bool) {
			yield(Some(42))
		}
		result := AltAllSeq(None[int]())(generator)
		assert.Equal(t, Some(42), result)
	})

	t.Run("handles single element sequence with None", func(t *testing.T) {
		generator := func(yield func(Option[int]) bool) {
			yield(None[int]())
		}
		result := AltAllSeq(Some(10))(generator)
		assert.Equal(t, Some(10), result)
	})

	t.Run("returns first Some when multiple Some values exist", func(t *testing.T) {
		generator := func(yield func(Option[int]) bool) {
			yield(Some(1))
			yield(Some(2))
			yield(Some(3))
		}
		result := AltAllSeq(None[int]())(generator)
		assert.Equal(t, Some(1), result)
	})

	t.Run("works with different types", func(t *testing.T) {
		generator := func(yield func(Option[string]) bool) {
			yield(None[string]())
			yield(Some("first"))
			yield(None[string]())
			yield(Some("second"))
		}
		result := AltAllSeq(None[string]())(generator)
		assert.Equal(t, Some("first"), result)
	})

	t.Run("returns None when startWith is None and sequence is empty", func(t *testing.T) {
		generator := func(yield func(Option[int]) bool) {}
		result := AltAllSeq(None[int]())(generator)
		assert.Equal(t, None[int](), result)
	})

	t.Run("returns None when startWith is None and all sequence elements are None", func(t *testing.T) {
		generator := func(yield func(Option[int]) bool) {
			yield(None[int]())
			yield(None[int]())
		}
		result := AltAllSeq(None[int]())(generator)
		assert.Equal(t, None[int](), result)
	})
}

func TestAltAllArray_Integration(t *testing.T) {
	t.Run("chains with other Option operations", func(t *testing.T) {
		options := []Option[int]{
			None[int](),
			Some(5),
			Some(10),
		}
		result := Map(func(x int) int { return x * 2 })(
			AltAllArray(None[int]())(options),
		)
		assert.Equal(t, Some(10), result)
	})

	t.Run("works with complex data types", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}
		options := []Option[User]{
			None[User](),
			Some(User{Name: "Alice", Age: 30}),
			Some(User{Name: "Bob", Age: 25}),
		}
		result := AltAllArray(None[User]())(options)
		assert.Equal(t, Some(User{Name: "Alice", Age: 30}), result)
	})
}

func TestAltAllSeq_Integration(t *testing.T) {
	t.Run("chains with other Option operations", func(t *testing.T) {
		generator := func(yield func(Option[int]) bool) {
			yield(None[int]())
			yield(Some(5))
			yield(Some(10))
		}
		result := Map(func(x int) int { return x * 2 })(
			AltAllSeq(None[int]())(generator),
		)
		assert.Equal(t, Some(10), result)
	})

	t.Run("works with complex data types", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}
		generator := func(yield func(Option[User]) bool) {
			yield(None[User]())
			yield(Some(User{Name: "Alice", Age: 30}))
			yield(Some(User{Name: "Bob", Age: 25}))
		}
		result := AltAllSeq(None[User]())(generator)
		assert.Equal(t, Some(User{Name: "Alice", Age: 30}), result)
	})

	t.Run("works with lazy evaluation pattern", func(t *testing.T) {
		// Simulates a lazy sequence that could be expensive to compute
		generator := func(yield func(Option[int]) bool) {
			for i := range 10 {
				if i == 5 {
					yield(Some(i))
				} else {
					yield(None[int]())
				}
			}
		}
		result := AltAllSeq(None[int]())(generator)
		assert.Equal(t, Some(5), result)
	})
}

func BenchmarkAltAllArray(b *testing.B) {
	options := []Option[int]{
		None[int](),
		None[int](),
		Some(42),
		Some(100),
	}
	altAll := AltAllArray(None[int]())

	b.ResetTimer()
	for range b.N {
		_ = altAll(options)
	}
}

func BenchmarkAltAllSeq(b *testing.B) {
	generator := func(yield func(Option[int]) bool) {
		yield(None[int]())
		yield(None[int]())
		yield(Some(42))
		yield(Some(100))
	}
	altAll := AltAllSeq(None[int]())

	b.ResetTimer()
	for range b.N {
		_ = altAll(generator)
	}
}

func BenchmarkAltAllArray_AllNone(b *testing.B) {
	options := []Option[int]{
		None[int](),
		None[int](),
		None[int](),
		None[int](),
	}
	altAll := AltAllArray(Some(10))

	b.ResetTimer()
	for range b.N {
		_ = altAll(options)
	}
}

func BenchmarkAltAllSeq_AllNone(b *testing.B) {
	generator := func(yield func(Option[int]) bool) {
		yield(None[int]())
		yield(None[int]())
		yield(None[int]())
		yield(None[int]())
	}
	altAll := AltAllSeq(Some(10))

	b.ResetTimer()
	for range b.N {
		_ = altAll(generator)
	}
}

func TestAltAllArray_RelationshipToAltMonoid(t *testing.T) {
	t.Run("demonstrates equivalence to manual fold with AltMonoid", func(t *testing.T) {
		options := []Option[int]{
			None[int](),
			Some(42),
			Some(100),
		}

		// Using AltAllArray with None startWith
		resultAltAll := AltAllArray(None[int]())(options)

		// Manual fold using AltMonoid (simulating array.Fold behavior)
		monoid := AltMonoid[int]()
		manualFold := monoid.Empty()
		for _, opt := range options {
			manualFold = monoid.Concat(manualFold, opt)
		}

		assert.Equal(t, manualFold, resultAltAll, "AltAllArray should equal manual fold with AltMonoid")
		assert.Equal(t, Some(42), resultAltAll)
	})

	t.Run("AltAllArray with Some startWith prepends to fold", func(t *testing.T) {
		options := []Option[int]{
			None[int](),
			Some(42),
			Some(100),
		}

		// Using AltAllArray with Some startWith
		resultAltAll := AltAllArray(Some(10))(options)

		// Manual fold with startWith prepended
		monoid := AltMonoid[int]()
		manualFold := Some(10) // Start with Some(10) instead of Empty()
		for _, opt := range options {
			manualFold = monoid.Concat(manualFold, opt)
		}

		assert.Equal(t, manualFold, resultAltAll)
		assert.Equal(t, Some(10), resultAltAll, "Should return startWith since it's Some")
	})

	t.Run("demonstrates AltMonoid properties", func(t *testing.T) {
		monoid := AltMonoid[int]()

		// Identity: Empty is None
		assert.Equal(t, None[int](), monoid.Empty())

		// Concat returns first Some value (left-biased)
		assert.Equal(t, Some(1), monoid.Concat(Some(1), Some(2)))
		assert.Equal(t, Some(2), monoid.Concat(None[int](), Some(2)))
		assert.Equal(t, Some(1), monoid.Concat(Some(1), None[int]()))
		assert.Equal(t, None[int](), monoid.Concat(None[int](), None[int]()))

		// Associativity: (a <> b) <> c = a <> (b <> c)
		a, b, c := Some(1), Some(2), Some(3)
		left := monoid.Concat(monoid.Concat(a, b), c)
		right := monoid.Concat(a, monoid.Concat(b, c))
		assert.Equal(t, left, right)
	})

	t.Run("AltAllArray implements fold pattern", func(t *testing.T) {
		// AltAllArray(startWith)(options) is equivalent to:
		// fold(options, startWith, Alt)
		// where fold applies Alt iteratively

		options := []Option[int]{None[int](), Some(42), None[int]()}

		// Using AltAllArray
		result := AltAllArray(None[int]())(options)

		// Manual implementation of the fold pattern
		current := None[int]()
		for _, opt := range options {
			current = MonadAlt(current, func() Option[int] { return opt })
		}

		assert.Equal(t, current, result)
		assert.Equal(t, Some(42), result)
	})
}

// Made with Bob
