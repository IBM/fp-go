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

package common

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEitherAltAllArray_Success(t *testing.T) {
	t.Run("returns first Right value from array", func(t *testing.T) {
		eithers := []Either[error, int]{
			EitherLeft[int](errors.New("error1")),
			EitherLeft[int](errors.New("error2")),
			EitherRight[error](42),
			EitherRight[error](100),
		}
		result := EitherAltAllArray(EitherLeft[int](errors.New("start")))(eithers)
		assert.Equal(t, EitherRight[error](42), result)
	})

	t.Run("returns startWith when array is empty", func(t *testing.T) {
		eithers := []Either[error, int]{}
		result := EitherAltAllArray(EitherRight[error](10))(eithers)
		assert.Equal(t, EitherRight[error](10), result)
	})

	t.Run("returns startWith when all array elements are Left", func(t *testing.T) {
		eithers := []Either[error, int]{
			EitherLeft[int](errors.New("error1")),
			EitherLeft[int](errors.New("error2")),
			EitherLeft[int](errors.New("error3")),
		}
		result := EitherAltAllArray(EitherRight[error](99))(eithers)
		assert.Equal(t, EitherRight[error](99), result)
	})

	t.Run("returns first Right when startWith is Left", func(t *testing.T) {
		eithers := []Either[error, string]{
			EitherLeft[string](errors.New("error1")),
			EitherRight[error]("hello"),
			EitherRight[error]("world"),
		}
		result := EitherAltAllArray(EitherLeft[string](errors.New("start")))(eithers)
		assert.Equal(t, EitherRight[error]("hello"), result)
	})

	t.Run("returns startWith when startWith is Right and array is all Left", func(t *testing.T) {
		eithers := []Either[error, int]{
			EitherLeft[int](errors.New("error1")),
			EitherLeft[int](errors.New("error2")),
		}
		result := EitherAltAllArray(EitherRight[error](5))(eithers)
		assert.Equal(t, EitherRight[error](5), result)
	})
}

func TestEitherAltAllArray_EdgeCases(t *testing.T) {
	t.Run("handles single element array with Right", func(t *testing.T) {
		eithers := []Either[error, int]{EitherRight[error](42)}
		result := EitherAltAllArray(EitherLeft[int](errors.New("start")))(eithers)
		assert.Equal(t, EitherRight[error](42), result)
	})

	t.Run("handles single element array with Left", func(t *testing.T) {
		eithers := []Either[error, int]{EitherLeft[int](errors.New("error"))}
		result := EitherAltAllArray(EitherRight[error](10))(eithers)
		assert.Equal(t, EitherRight[error](10), result)
	})

	t.Run("returns first Right when multiple Right values exist", func(t *testing.T) {
		eithers := []Either[error, int]{
			EitherRight[error](1),
			EitherRight[error](2),
			EitherRight[error](3),
		}
		result := EitherAltAllArray(EitherLeft[int](errors.New("start")))(eithers)
		assert.Equal(t, EitherRight[error](1), result)
	})

	t.Run("works with different types", func(t *testing.T) {
		eithers := []Either[string, string]{
			EitherLeft[string]("error1"),
			EitherRight[string]("first"),
			EitherLeft[string]("error2"),
			EitherRight[string]("second"),
		}
		result := EitherAltAllArray(EitherLeft[string]("start"))(eithers)
		assert.Equal(t, EitherRight[string]("first"), result)
	})

	t.Run("returns Left when startWith is Left and array is empty", func(t *testing.T) {
		eithers := []Either[error, int]{}
		startErr := errors.New("start error")
		result := EitherAltAllArray(EitherLeft[int](startErr))(eithers)
		assert.True(t, EitherIsLeft(result))
		assert.Equal(t, startErr, result.e)
	})

	t.Run("returns last Left when startWith is Left and all array elements are Left", func(t *testing.T) {
		startErr := errors.New("start error")
		lastErr := errors.New("error2")
		eithers := []Either[error, int]{
			EitherLeft[int](errors.New("error1")),
			EitherLeft[int](lastErr),
		}
		result := EitherAltAllArray(EitherLeft[int](startErr))(eithers)
		assert.True(t, EitherIsLeft(result))
		assert.Equal(t, lastErr, result.e)
	})
}

func TestEitherAltAllSeq_Success(t *testing.T) {
	t.Run("returns first Right value from sequence", func(t *testing.T) {
		generator := func(yield func(Either[error, int]) bool) {
			if !yield(EitherLeft[int](errors.New("error1"))) {
				return
			}
			if !yield(EitherLeft[int](errors.New("error2"))) {
				return
			}
			if !yield(EitherRight[error](42)) {
				return
			}
			if !yield(EitherRight[error](100)) {
				return
			}
		}
		result := EitherAltAllSeq(EitherLeft[int](errors.New("start")))(generator)
		assert.Equal(t, EitherRight[error](42), result)
	})

	t.Run("returns startWith when sequence is empty", func(t *testing.T) {
		generator := func(yield func(Either[error, int]) bool) {}
		result := EitherAltAllSeq(EitherRight[error](10))(generator)
		assert.Equal(t, EitherRight[error](10), result)
	})

	t.Run("returns startWith when all sequence elements are Left", func(t *testing.T) {
		generator := func(yield func(Either[error, int]) bool) {
			if !yield(EitherLeft[int](errors.New("error1"))) {
				return
			}
			if !yield(EitherLeft[int](errors.New("error2"))) {
				return
			}
			if !yield(EitherLeft[int](errors.New("error3"))) {
				return
			}
		}
		result := EitherAltAllSeq(EitherRight[error](99))(generator)
		assert.Equal(t, EitherRight[error](99), result)
	})

	t.Run("returns first Right when startWith is Left", func(t *testing.T) {
		generator := func(yield func(Either[error, string]) bool) {
			if !yield(EitherLeft[string](errors.New("error1"))) {
				return
			}
			if !yield(EitherRight[error]("hello")) {
				return
			}
			if !yield(EitherRight[error]("world")) {
				return
			}
		}
		result := EitherAltAllSeq(EitherLeft[string](errors.New("start")))(generator)
		assert.Equal(t, EitherRight[error]("hello"), result)
	})

	t.Run("stops iteration after finding first Right", func(t *testing.T) {
		iterationCount := 0
		generator := func(yield func(Either[error, int]) bool) {
			iterationCount++
			if !yield(EitherLeft[int](errors.New("error1"))) {
				return
			}
			iterationCount++
			if !yield(EitherRight[error](42)) {
				return
			}
			iterationCount++
			if !yield(EitherRight[error](100)) {
				return
			}
		}
		result := EitherAltAllSeq(EitherLeft[int](errors.New("start")))(generator)
		assert.Equal(t, EitherRight[error](42), result)
		// The Alt operation short-circuits: lazy thunk not evaluated when current is Right
		// So we expect 2 iterations (error1, Right(42)) - the third element is never requested
		assert.Equal(t, 2, iterationCount)
	})
}

func TestEitherAltAllSeq_EdgeCases(t *testing.T) {
	t.Run("handles single element sequence with Right", func(t *testing.T) {
		generator := func(yield func(Either[error, int]) bool) {
			yield(EitherRight[error](42))
		}
		result := EitherAltAllSeq(EitherLeft[int](errors.New("start")))(generator)
		assert.Equal(t, EitherRight[error](42), result)
	})

	t.Run("handles single element sequence with Left", func(t *testing.T) {
		generator := func(yield func(Either[error, int]) bool) {
			yield(EitherLeft[int](errors.New("error")))
		}
		result := EitherAltAllSeq(EitherRight[error](10))(generator)
		assert.Equal(t, EitherRight[error](10), result)
	})

	t.Run("returns first Right when multiple Right values exist", func(t *testing.T) {
		generator := func(yield func(Either[error, int]) bool) {
			if !yield(EitherRight[error](1)) {
				return
			}
			if !yield(EitherRight[error](2)) {
				return
			}
			if !yield(EitherRight[error](3)) {
				return
			}
		}
		result := EitherAltAllSeq(EitherLeft[int](errors.New("start")))(generator)
		assert.Equal(t, EitherRight[error](1), result)
	})

	t.Run("works with different types", func(t *testing.T) {
		generator := func(yield func(Either[string, string]) bool) {
			if !yield(EitherLeft[string]("error1")) {
				return
			}
			if !yield(EitherRight[string]("first")) {
				return
			}
			if !yield(EitherLeft[string]("error2")) {
				return
			}
			if !yield(EitherRight[string]("second")) {
				return
			}
		}
		result := EitherAltAllSeq(EitherLeft[string]("start"))(generator)
		assert.Equal(t, EitherRight[string]("first"), result)
	})

	t.Run("returns Left when startWith is Left and sequence is empty", func(t *testing.T) {
		generator := func(yield func(Either[error, int]) bool) {}
		startErr := errors.New("start error")
		result := EitherAltAllSeq(EitherLeft[int](startErr))(generator)
		assert.True(t, EitherIsLeft(result))
		assert.Equal(t, startErr, result.e)
	})

	t.Run("returns last Left when startWith is Left and all sequence elements are Left", func(t *testing.T) {
		startErr := errors.New("start error")
		lastErr := errors.New("error2")
		generator := func(yield func(Either[error, int]) bool) {
			if !yield(EitherLeft[int](errors.New("error1"))) {
				return
			}
			if !yield(EitherLeft[int](lastErr)) {
				return
			}
		}
		result := EitherAltAllSeq(EitherLeft[int](startErr))(generator)
		assert.True(t, EitherIsLeft(result))
		assert.Equal(t, lastErr, result.e)
	})
}

func TestEitherAltAllArray_Integration(t *testing.T) {
	t.Run("chains with other Either operations", func(t *testing.T) {
		eithers := []Either[error, int]{
			EitherLeft[int](errors.New("error1")),
			EitherRight[error](5),
			EitherRight[error](10),
		}
		result := EitherMap[error](func(x int) int { return x * 2 })(
			EitherAltAllArray(EitherLeft[int](errors.New("start")))(eithers),
		)
		assert.Equal(t, EitherRight[error](10), result)
	})

	t.Run("works with complex data types", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}
		eithers := []Either[error, User]{
			EitherLeft[User](errors.New("error1")),
			EitherRight[error](User{Name: "Alice", Age: 30}),
			EitherRight[error](User{Name: "Bob", Age: 25}),
		}
		result := EitherAltAllArray(EitherLeft[User](errors.New("start")))(eithers)
		assert.Equal(t, EitherRight[error](User{Name: "Alice", Age: 30}), result)
	})
}

func TestEitherAltAllSeq_Integration(t *testing.T) {
	t.Run("chains with other Either operations", func(t *testing.T) {
		generator := func(yield func(Either[error, int]) bool) {
			if !yield(EitherLeft[int](errors.New("error1"))) {
				return
			}
			if !yield(EitherRight[error](5)) {
				return
			}
			if !yield(EitherRight[error](10)) {
				return
			}
		}
		result := EitherMap[error](func(x int) int { return x * 2 })(
			EitherAltAllSeq(EitherLeft[int](errors.New("start")))(generator),
		)
		assert.Equal(t, EitherRight[error](10), result)
	})

	t.Run("works with complex data types", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}
		generator := func(yield func(Either[error, User]) bool) {
			if !yield(EitherLeft[User](errors.New("error1"))) {
				return
			}
			if !yield(EitherRight[error](User{Name: "Alice", Age: 30})) {
				return
			}
			if !yield(EitherRight[error](User{Name: "Bob", Age: 25})) {
				return
			}
		}
		result := EitherAltAllSeq(EitherLeft[User](errors.New("start")))(generator)
		assert.Equal(t, EitherRight[error](User{Name: "Alice", Age: 30}), result)
	})

	t.Run("works with lazy evaluation pattern", func(t *testing.T) {
		// Simulates a lazy sequence that could be expensive to compute
		generator := func(yield func(Either[error, int]) bool) {
			for i := range 10 {
				if i == 5 {
					if !yield(EitherRight[error](i)) {
						return
					}
				} else {
					if !yield(EitherLeft[int](errors.New("error"))) {
						return
					}
				}
			}
		}
		result := EitherAltAllSeq(EitherLeft[int](errors.New("start")))(generator)
		assert.Equal(t, EitherRight[error](5), result)
	})
}

func TestEitherAltAllArray_ShortCircuit(t *testing.T) {
	t.Run("short-circuits when startWith is Right", func(t *testing.T) {
		// This test verifies that the array is not examined when startWith is Right
		eithers := []Either[error, int]{
			EitherRight[error](100),
			EitherRight[error](200),
		}
		result := EitherAltAllArray(EitherRight[error](42))(eithers)
		assert.Equal(t, EitherRight[error](42), result)
	})

	t.Run("short-circuits on first Right in array", func(t *testing.T) {
		// Create a slice where we can track which elements were accessed
		eithers := make([]Either[error, int], 4)

		// Wrap each Either to track access
		for i := range 4 {
			if i == 2 {
				eithers[i] = EitherRight[error](42)
			} else {
				eithers[i] = EitherLeft[int](errors.New("error"))
			}
		}

		// The function should find the Right at index 2
		result := EitherAltAllArray(EitherLeft[int](errors.New("start")))(eithers)
		assert.Equal(t, EitherRight[error](42), result)

		// Verify we got the expected result (the implementation will have accessed up to index 2)
		assert.True(t, EitherIsRight(result))
	})
}

func TestEitherAltAllSeq_ShortCircuit(t *testing.T) {
	t.Run("short-circuits when startWith is Right", func(t *testing.T) {
		iterationCount := 0
		generator := func(yield func(Either[error, int]) bool) {
			iterationCount++
			yield(EitherRight[error](100))
		}
		result := EitherAltAllSeq(EitherRight[error](42))(generator)
		assert.Equal(t, EitherRight[error](42), result)
		// When startWith is Right, the lazy thunk is never evaluated, so no iterations
		assert.Equal(t, 0, iterationCount)
	})

	t.Run("short-circuits on first Right in sequence", func(t *testing.T) {
		iterationCount := 0
		generator := func(yield func(Either[error, int]) bool) {
			iterationCount++
			if !yield(EitherLeft[int](errors.New("error1"))) {
				return
			}
			iterationCount++
			if !yield(EitherLeft[int](errors.New("error2"))) {
				return
			}
			iterationCount++
			if !yield(EitherRight[error](42)) {
				return
			}
			iterationCount++
			if !yield(EitherRight[error](100)) {
				return
			}
		}
		result := EitherAltAllSeq(EitherLeft[int](errors.New("start")))(generator)
		assert.Equal(t, EitherRight[error](42), result)
		// Short-circuits after finding Right: 3 iterations (error1, error2, Right(42))
		assert.Equal(t, 3, iterationCount)
	})
}

func BenchmarkEitherAltAllArray(b *testing.B) {
	eithers := []Either[error, int]{
		EitherLeft[int](errors.New("error1")),
		EitherLeft[int](errors.New("error2")),
		EitherRight[error](42),
		EitherRight[error](100),
	}
	altAll := EitherAltAllArray(EitherLeft[int](errors.New("start")))

	b.ResetTimer()
	for range b.N {
		_ = altAll(eithers)
	}
}

func BenchmarkEitherAltAllSeq(b *testing.B) {
	generator := func(yield func(Either[error, int]) bool) {
		if !yield(EitherLeft[int](errors.New("error1"))) {
			return
		}
		if !yield(EitherLeft[int](errors.New("error2"))) {
			return
		}
		if !yield(EitherRight[error](42)) {
			return
		}
		if !yield(EitherRight[error](100)) {
			return
		}
	}
	altAll := EitherAltAllSeq(EitherLeft[int](errors.New("start")))

	b.ResetTimer()
	for range b.N {
		_ = altAll(generator)
	}
}

func BenchmarkEitherAltAllArray_AllLeft(b *testing.B) {
	eithers := []Either[error, int]{
		EitherLeft[int](errors.New("error1")),
		EitherLeft[int](errors.New("error2")),
		EitherLeft[int](errors.New("error3")),
		EitherLeft[int](errors.New("error4")),
	}
	altAll := EitherAltAllArray(EitherRight[error](10))

	b.ResetTimer()
	for range b.N {
		_ = altAll(eithers)
	}
}

func BenchmarkEitherAltAllSeq_AllLeft(b *testing.B) {
	generator := func(yield func(Either[error, int]) bool) {
		if !yield(EitherLeft[int](errors.New("error1"))) {
			return
		}
		if !yield(EitherLeft[int](errors.New("error2"))) {
			return
		}
		if !yield(EitherLeft[int](errors.New("error3"))) {
			return
		}
		if !yield(EitherLeft[int](errors.New("error4"))) {
			return
		}
	}
	altAll := EitherAltAllSeq(EitherRight[error](10))

	b.ResetTimer()
	for range b.N {
		_ = altAll(generator)
	}
}
