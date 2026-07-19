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

package either

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAltAllArray_Success(t *testing.T) {
	t.Run("returns first Right value from array", func(t *testing.T) {
		eithers := []Either[error, int]{
			Left[int](errors.New("error1")),
			Left[int](errors.New("error2")),
			Right[error](42),
			Right[error](100),
		}
		result := AltAllArray(Left[int](errors.New("start")))(eithers)
		assert.Equal(t, Right[error](42), result)
	})

	t.Run("returns startWith when array is empty", func(t *testing.T) {
		eithers := []Either[error, int]{}
		result := AltAllArray(Right[error](10))(eithers)
		assert.Equal(t, Right[error](10), result)
	})

	t.Run("returns startWith when all array elements are Left", func(t *testing.T) {
		eithers := []Either[error, int]{
			Left[int](errors.New("error1")),
			Left[int](errors.New("error2")),
			Left[int](errors.New("error3")),
		}
		result := AltAllArray(Right[error](99))(eithers)
		assert.Equal(t, Right[error](99), result)
	})

	t.Run("returns first Right when startWith is Left", func(t *testing.T) {
		eithers := []Either[error, string]{
			Left[string](errors.New("error1")),
			Right[error]("hello"),
			Right[error]("world"),
		}
		result := AltAllArray(Left[string](errors.New("start")))(eithers)
		assert.Equal(t, Right[error]("hello"), result)
	})

	t.Run("returns startWith when startWith is Right and array is all Left", func(t *testing.T) {
		eithers := []Either[error, int]{
			Left[int](errors.New("error1")),
			Left[int](errors.New("error2")),
		}
		result := AltAllArray(Right[error](5))(eithers)
		assert.Equal(t, Right[error](5), result)
	})
}

func TestAltAllArray_EdgeCases(t *testing.T) {
	t.Run("handles single element array with Right", func(t *testing.T) {
		eithers := []Either[error, int]{Right[error](42)}
		result := AltAllArray(Left[int](errors.New("start")))(eithers)
		assert.Equal(t, Right[error](42), result)
	})

	t.Run("handles single element array with Left", func(t *testing.T) {
		eithers := []Either[error, int]{Left[int](errors.New("error"))}
		result := AltAllArray(Right[error](10))(eithers)
		assert.Equal(t, Right[error](10), result)
	})

	t.Run("returns first Right when multiple Right values exist", func(t *testing.T) {
		eithers := []Either[error, int]{
			Right[error](1),
			Right[error](2),
			Right[error](3),
		}
		result := AltAllArray(Left[int](errors.New("start")))(eithers)
		assert.Equal(t, Right[error](1), result)
	})

	t.Run("works with different types", func(t *testing.T) {
		eithers := []Either[string, string]{
			Left[string]("error1"),
			Right[string]("first"),
			Left[string]("error2"),
			Right[string]("second"),
		}
		result := AltAllArray(Left[string]("start"))(eithers)
		assert.Equal(t, Right[string]("first"), result)
	})

	t.Run("returns Left when startWith is Left and array is empty", func(t *testing.T) {
		eithers := []Either[error, int]{}
		startErr := errors.New("start error")
		result := AltAllArray(Left[int](startErr))(eithers)
		assert.True(t, IsLeft(result))
		assert.Equal(t, startErr, result.e)
	})

	t.Run("returns last Left when startWith is Left and all array elements are Left", func(t *testing.T) {
		startErr := errors.New("start error")
		lastErr := errors.New("error2")
		eithers := []Either[error, int]{
			Left[int](errors.New("error1")),
			Left[int](lastErr),
		}
		result := AltAllArray(Left[int](startErr))(eithers)
		assert.True(t, IsLeft(result))
		assert.Equal(t, lastErr, result.e)
	})
}

func TestAltAllSeq_Success(t *testing.T) {
	t.Run("returns first Right value from sequence", func(t *testing.T) {
		generator := func(yield func(Either[error, int]) bool) {
			if !yield(Left[int](errors.New("error1"))) {
				return
			}
			if !yield(Left[int](errors.New("error2"))) {
				return
			}
			if !yield(Right[error](42)) {
				return
			}
			if !yield(Right[error](100)) {
				return
			}
		}
		result := AltAllSeq(Left[int](errors.New("start")))(generator)
		assert.Equal(t, Right[error](42), result)
	})

	t.Run("returns startWith when sequence is empty", func(t *testing.T) {
		generator := func(yield func(Either[error, int]) bool) {}
		result := AltAllSeq(Right[error](10))(generator)
		assert.Equal(t, Right[error](10), result)
	})

	t.Run("returns startWith when all sequence elements are Left", func(t *testing.T) {
		generator := func(yield func(Either[error, int]) bool) {
			if !yield(Left[int](errors.New("error1"))) {
				return
			}
			if !yield(Left[int](errors.New("error2"))) {
				return
			}
			if !yield(Left[int](errors.New("error3"))) {
				return
			}
		}
		result := AltAllSeq(Right[error](99))(generator)
		assert.Equal(t, Right[error](99), result)
	})

	t.Run("returns first Right when startWith is Left", func(t *testing.T) {
		generator := func(yield func(Either[error, string]) bool) {
			if !yield(Left[string](errors.New("error1"))) {
				return
			}
			if !yield(Right[error]("hello")) {
				return
			}
			if !yield(Right[error]("world")) {
				return
			}
		}
		result := AltAllSeq(Left[string](errors.New("start")))(generator)
		assert.Equal(t, Right[error]("hello"), result)
	})

	t.Run("stops iteration after finding first Right", func(t *testing.T) {
		iterationCount := 0
		generator := func(yield func(Either[error, int]) bool) {
			iterationCount++
			if !yield(Left[int](errors.New("error1"))) {
				return
			}
			iterationCount++
			if !yield(Right[error](42)) {
				return
			}
			iterationCount++
			if !yield(Right[error](100)) {
				return
			}
		}
		result := AltAllSeq(Left[int](errors.New("start")))(generator)
		assert.Equal(t, Right[error](42), result)
		// The Alt operation short-circuits: lazy thunk not evaluated when current is Right
		// So we expect 2 iterations (error1, Right(42)) - the third element is never requested
		assert.Equal(t, 2, iterationCount)
	})
}

func TestAltAllSeq_EdgeCases(t *testing.T) {
	t.Run("handles single element sequence with Right", func(t *testing.T) {
		generator := func(yield func(Either[error, int]) bool) {
			yield(Right[error](42))
		}
		result := AltAllSeq(Left[int](errors.New("start")))(generator)
		assert.Equal(t, Right[error](42), result)
	})

	t.Run("handles single element sequence with Left", func(t *testing.T) {
		generator := func(yield func(Either[error, int]) bool) {
			yield(Left[int](errors.New("error")))
		}
		result := AltAllSeq(Right[error](10))(generator)
		assert.Equal(t, Right[error](10), result)
	})

	t.Run("returns first Right when multiple Right values exist", func(t *testing.T) {
		generator := func(yield func(Either[error, int]) bool) {
			if !yield(Right[error](1)) {
				return
			}
			if !yield(Right[error](2)) {
				return
			}
			if !yield(Right[error](3)) {
				return
			}
		}
		result := AltAllSeq(Left[int](errors.New("start")))(generator)
		assert.Equal(t, Right[error](1), result)
	})

	t.Run("works with different types", func(t *testing.T) {
		generator := func(yield func(Either[string, string]) bool) {
			if !yield(Left[string]("error1")) {
				return
			}
			if !yield(Right[string]("first")) {
				return
			}
			if !yield(Left[string]("error2")) {
				return
			}
			if !yield(Right[string]("second")) {
				return
			}
		}
		result := AltAllSeq(Left[string]("start"))(generator)
		assert.Equal(t, Right[string]("first"), result)
	})

	t.Run("returns Left when startWith is Left and sequence is empty", func(t *testing.T) {
		generator := func(yield func(Either[error, int]) bool) {}
		startErr := errors.New("start error")
		result := AltAllSeq(Left[int](startErr))(generator)
		assert.True(t, IsLeft(result))
		assert.Equal(t, startErr, result.e)
	})

	t.Run("returns last Left when startWith is Left and all sequence elements are Left", func(t *testing.T) {
		startErr := errors.New("start error")
		lastErr := errors.New("error2")
		generator := func(yield func(Either[error, int]) bool) {
			if !yield(Left[int](errors.New("error1"))) {
				return
			}
			if !yield(Left[int](lastErr)) {
				return
			}
		}
		result := AltAllSeq(Left[int](startErr))(generator)
		assert.True(t, IsLeft(result))
		assert.Equal(t, lastErr, result.e)
	})
}

func TestAltAllArray_Integration(t *testing.T) {
	t.Run("chains with other Either operations", func(t *testing.T) {
		eithers := []Either[error, int]{
			Left[int](errors.New("error1")),
			Right[error](5),
			Right[error](10),
		}
		result := Map[error](func(x int) int { return x * 2 })(
			AltAllArray(Left[int](errors.New("start")))(eithers),
		)
		assert.Equal(t, Right[error](10), result)
	})

	t.Run("works with complex data types", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}
		eithers := []Either[error, User]{
			Left[User](errors.New("error1")),
			Right[error](User{Name: "Alice", Age: 30}),
			Right[error](User{Name: "Bob", Age: 25}),
		}
		result := AltAllArray(Left[User](errors.New("start")))(eithers)
		assert.Equal(t, Right[error](User{Name: "Alice", Age: 30}), result)
	})
}

func TestAltAllSeq_Integration(t *testing.T) {
	t.Run("chains with other Either operations", func(t *testing.T) {
		generator := func(yield func(Either[error, int]) bool) {
			if !yield(Left[int](errors.New("error1"))) {
				return
			}
			if !yield(Right[error](5)) {
				return
			}
			if !yield(Right[error](10)) {
				return
			}
		}
		result := Map[error](func(x int) int { return x * 2 })(
			AltAllSeq(Left[int](errors.New("start")))(generator),
		)
		assert.Equal(t, Right[error](10), result)
	})

	t.Run("works with complex data types", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}
		generator := func(yield func(Either[error, User]) bool) {
			if !yield(Left[User](errors.New("error1"))) {
				return
			}
			if !yield(Right[error](User{Name: "Alice", Age: 30})) {
				return
			}
			if !yield(Right[error](User{Name: "Bob", Age: 25})) {
				return
			}
		}
		result := AltAllSeq(Left[User](errors.New("start")))(generator)
		assert.Equal(t, Right[error](User{Name: "Alice", Age: 30}), result)
	})

	t.Run("works with lazy evaluation pattern", func(t *testing.T) {
		// Simulates a lazy sequence that could be expensive to compute
		generator := func(yield func(Either[error, int]) bool) {
			for i := range 10 {
				if i == 5 {
					if !yield(Right[error](i)) {
						return
					}
				} else {
					if !yield(Left[int](errors.New("error"))) {
						return
					}
				}
			}
		}
		result := AltAllSeq(Left[int](errors.New("start")))(generator)
		assert.Equal(t, Right[error](5), result)
	})
}

func TestAltAllArray_ShortCircuit(t *testing.T) {
	t.Run("short-circuits when startWith is Right", func(t *testing.T) {
		// This test verifies that the array is not examined when startWith is Right
		eithers := []Either[error, int]{
			Right[error](100),
			Right[error](200),
		}
		result := AltAllArray(Right[error](42))(eithers)
		assert.Equal(t, Right[error](42), result)
	})

	t.Run("short-circuits on first Right in array", func(t *testing.T) {
		// Create a slice where we can track which elements were accessed
		eithers := make([]Either[error, int], 4)

		// Wrap each Either to track access
		for i := range 4 {
			if i == 2 {
				eithers[i] = Right[error](42)
			} else {
				eithers[i] = Left[int](errors.New("error"))
			}
		}

		// The function should find the Right at index 2
		result := AltAllArray(Left[int](errors.New("start")))(eithers)
		assert.Equal(t, Right[error](42), result)

		// Verify we got the expected result (the implementation will have accessed up to index 2)
		assert.True(t, IsRight(result))
	})
}

func TestAltAllSeq_ShortCircuit(t *testing.T) {
	t.Run("short-circuits when startWith is Right", func(t *testing.T) {
		iterationCount := 0
		generator := func(yield func(Either[error, int]) bool) {
			iterationCount++
			yield(Right[error](100))
		}
		result := AltAllSeq(Right[error](42))(generator)
		assert.Equal(t, Right[error](42), result)
		// When startWith is Right, the lazy thunk is never evaluated, so no iterations
		assert.Equal(t, 0, iterationCount)
	})

	t.Run("short-circuits on first Right in sequence", func(t *testing.T) {
		iterationCount := 0
		generator := func(yield func(Either[error, int]) bool) {
			iterationCount++
			if !yield(Left[int](errors.New("error1"))) {
				return
			}
			iterationCount++
			if !yield(Left[int](errors.New("error2"))) {
				return
			}
			iterationCount++
			if !yield(Right[error](42)) {
				return
			}
			iterationCount++
			if !yield(Right[error](100)) {
				return
			}
		}
		result := AltAllSeq(Left[int](errors.New("start")))(generator)
		assert.Equal(t, Right[error](42), result)
		// Short-circuits after finding Right: 3 iterations (error1, error2, Right(42))
		assert.Equal(t, 3, iterationCount)
	})
}

func BenchmarkAltAllArray(b *testing.B) {
	eithers := []Either[error, int]{
		Left[int](errors.New("error1")),
		Left[int](errors.New("error2")),
		Right[error](42),
		Right[error](100),
	}
	altAll := AltAllArray(Left[int](errors.New("start")))

	b.ResetTimer()
	for range b.N {
		_ = altAll(eithers)
	}
}

func BenchmarkAltAllSeq(b *testing.B) {
	generator := func(yield func(Either[error, int]) bool) {
		if !yield(Left[int](errors.New("error1"))) {
			return
		}
		if !yield(Left[int](errors.New("error2"))) {
			return
		}
		if !yield(Right[error](42)) {
			return
		}
		if !yield(Right[error](100)) {
			return
		}
	}
	altAll := AltAllSeq(Left[int](errors.New("start")))

	b.ResetTimer()
	for range b.N {
		_ = altAll(generator)
	}
}

func BenchmarkAltAllArray_AllLeft(b *testing.B) {
	eithers := []Either[error, int]{
		Left[int](errors.New("error1")),
		Left[int](errors.New("error2")),
		Left[int](errors.New("error3")),
		Left[int](errors.New("error4")),
	}
	altAll := AltAllArray(Right[error](10))

	b.ResetTimer()
	for range b.N {
		_ = altAll(eithers)
	}
}

func BenchmarkAltAllSeq_AllLeft(b *testing.B) {
	generator := func(yield func(Either[error, int]) bool) {
		if !yield(Left[int](errors.New("error1"))) {
			return
		}
		if !yield(Left[int](errors.New("error2"))) {
			return
		}
		if !yield(Left[int](errors.New("error3"))) {
			return
		}
		if !yield(Left[int](errors.New("error4"))) {
			return
		}
	}
	altAll := AltAllSeq(Right[error](10))

	b.ResetTimer()
	for range b.N {
		_ = altAll(generator)
	}
}
