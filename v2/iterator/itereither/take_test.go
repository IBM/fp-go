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

package itereither

import (
	"errors"
	"strconv"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/iterator/iter"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// TestTakeUntilLeft_Success tests basic TakeUntilLeft functionality with successful sequences
func TestTakeUntilLeft_Success(t *testing.T) {
	t.Run("takes all Right values until first Left", func(t *testing.T) {
		seq := iter.From(
			E.Right[string](1),
			E.Right[string](2),
			E.Right[string](3),
			E.Left[int]("error"),
			E.Right[string](4),
			E.Right[string](5),
		)
		result := collectEithers(TakeUntilLeft(seq))
		expected := []Either[string, int]{
			E.Right[string](1),
			E.Right[string](2),
			E.Right[string](3),
			E.Left[int]("error"),
		}
		assert.Equal(t, expected, result)
	})

	t.Run("takes all elements when no Left exists", func(t *testing.T) {
		seq := iter.From(
			E.Right[string](1),
			E.Right[string](2),
			E.Right[string](3),
		)
		result := collectEithers(TakeUntilLeft(seq))
		expected := []Either[string, int]{
			E.Right[string](1),
			E.Right[string](2),
			E.Right[string](3),
		}
		assert.Equal(t, expected, result)
	})

	t.Run("stops immediately when first element is Left", func(t *testing.T) {
		seq := iter.From(
			E.Left[int]("immediate error"),
			E.Right[string](1),
			E.Right[string](2),
		)
		result := collectEithers(TakeUntilLeft(seq))
		expected := []Either[string, int]{
			E.Left[int]("immediate error"),
		}
		assert.Equal(t, expected, result)
	})

	t.Run("handles single Right element", func(t *testing.T) {
		seq := iter.From(E.Right[string](42))
		result := collectEithers(TakeUntilLeft(seq))
		expected := []Either[string, int]{E.Right[string](42)}
		assert.Equal(t, expected, result)
	})

	t.Run("handles single Left element", func(t *testing.T) {
		seq := iter.From(E.Left[int]("error"))
		result := collectEithers(TakeUntilLeft(seq))
		expected := []Either[string, int]{E.Left[int]("error")}
		assert.Equal(t, expected, result)
	})
}

// TestTakeUntilLeft_Empty tests TakeUntilLeft with empty sequences
func TestTakeUntilLeft_Empty(t *testing.T) {
	t.Run("returns empty from empty sequence", func(t *testing.T) {
		seq := iter.Empty[Either[string, int]]()
		result := collectEithers(TakeUntilLeft(seq))
		assert.Empty(t, result)
	})
}

// TestTakeUntilLeft_WithComplexTypes tests TakeUntilLeft with complex types
func TestTakeUntilLeft_WithComplexTypes(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}

	t.Run("works with struct types", func(t *testing.T) {
		seq := iter.From(
			E.Right[error](User{1, "Alice"}),
			E.Right[error](User{2, "Bob"}),
			E.Left[User](errors.New("database error")),
			E.Right[error](User{3, "Charlie"}),
		)
		result := collectEithers(TakeUntilLeft(seq))
		expected := []Either[error, User]{
			E.Right[error](User{1, "Alice"}),
			E.Right[error](User{2, "Bob"}),
			E.Left[User](errors.New("database error")),
		}
		assert.Equal(t, expected, result)
	})

	t.Run("works with string errors", func(t *testing.T) {
		seq := iter.From(
			E.Right[string]("hello"),
			E.Right[string]("world"),
			E.Left[string]("parse error"),
			E.Right[string]("ignored"),
		)
		result := collectEithers(TakeUntilLeft(seq))
		expected := []Either[string, string]{
			E.Right[string]("hello"),
			E.Right[string]("world"),
			E.Left[string]("parse error"),
		}
		assert.Equal(t, expected, result)
	})

	t.Run("works with slice types", func(t *testing.T) {
		seq := iter.From(
			E.Right[string]([]int{1, 2}),
			E.Right[string]([]int{3, 4}),
			E.Left[[]int]("validation error"),
			E.Right[string]([]int{5, 6}),
		)
		result := collectEithers(TakeUntilLeft(seq))
		expected := []Either[string, []int]{
			E.Right[string]([]int{1, 2}),
			E.Right[string]([]int{3, 4}),
			E.Left[[]int]("validation error"),
		}
		assert.Equal(t, expected, result)
	})
}

// TestTakeUntilLeft_WithChainedOperations tests TakeUntilLeft with other operations
func TestTakeUntilLeft_WithChainedOperations(t *testing.T) {
	t.Run("chains with Map", func(t *testing.T) {
		seq := iter.From(
			E.Right[string](1),
			E.Right[string](2),
			E.Left[int]("error"),
			E.Right[string](3),
		)
		result := collectEithers(
			TakeUntilLeft(
				MonadMap(seq, func(n int) int { return n * 2 }),
			),
		)
		expected := []Either[string, int]{
			E.Right[string](2),
			E.Right[string](4),
			E.Left[int]("error"),
		}
		assert.Equal(t, expected, result)
	})

	t.Run("chains with FilterOrElse", func(t *testing.T) {
		seq := iter.From(
			E.Right[string](1),
			E.Right[string](2),
			E.Right[string](3),
			E.Right[string](4),
		)
		isEven := func(n int) bool { return n%2 == 0 }
		onOdd := func(n int) string { return "odd number" }

		result := collectEithers(
			TakeUntilLeft(
				FilterOrElse(isEven, onOdd)(seq),
			),
		)
		expected := []Either[string, int]{
			E.Left[int]("odd number"), // 1 is odd, creates Left, stops here
		}
		assert.Equal(t, expected, result)
	})

	t.Run("chains with multiple TakeUntilLeft", func(t *testing.T) {
		seq := iter.From(
			E.Right[string](1),
			E.Right[string](2),
			E.Left[int]("first error"),
			E.Right[string](3),
			E.Left[int]("second error"),
		)
		// First TakeUntilLeft stops at "first error"
		// Second TakeUntilLeft would also stop at "first error" (same result)
		result := collectEithers(TakeUntilLeft(TakeUntilLeft(seq)))
		expected := []Either[string, int]{
			E.Right[string](1),
			E.Right[string](2),
			E.Left[int]("first error"),
		}
		assert.Equal(t, expected, result)
	})
}

// TestTakeUntilLeft_WithRealWorldScenarios tests practical use cases
func TestTakeUntilLeft_WithRealWorldScenarios(t *testing.T) {
	t.Run("parsing numbers until error", func(t *testing.T) {
		inputs := []string{"1", "2", "3", "invalid", "4", "5"}
		seq := iter.MonadMap(
			iter.From(inputs...),
			result.Eitherize1(strconv.Atoi),
		)
		result := collectEithers(TakeUntilLeft(seq))

		assert.Len(t, result, 4) // 1, 2, 3, and the error
		assert.True(t, E.IsRight(result[0]))
		assert.True(t, E.IsRight(result[1]))
		assert.True(t, E.IsRight(result[2]))
		assert.True(t, E.IsLeft(result[3]))
	})

	t.Run("processing until validation failure", func(t *testing.T) {
		type Record struct {
			ID    int
			Valid bool
		}

		records := []Record{
			{1, true},
			{2, true},
			{3, false},
			{4, true},
		}

		validate := func(r Record) Either[string, Record] {
			if r.Valid {
				return E.Right[string](r)
			}
			return E.Left[Record]("invalid record")
		}

		seq := iter.MonadMap(iter.From(records...), validate)
		result := collectEithers(TakeUntilLeft(seq))

		expected := []Either[string, Record]{
			E.Right[string](Record{1, true}),
			E.Right[string](Record{2, true}),
			E.Left[Record]("invalid record"),
		}
		assert.Equal(t, expected, result)
	})
}

// TestTakeUntilLeft_DoesNotConsumeEntireSequence verifies lazy evaluation
func TestTakeUntilLeft_DoesNotConsumeEntireSequence(t *testing.T) {
	t.Run("stops consuming after Left", func(t *testing.T) {
		consumed := 0
		seq := func(yield func(Either[string, int]) bool) {
			for i := range 10 {
				consumed++
				var e Either[string, int]
				if i == 3 {
					e = E.Left[int]("error at 3")
				} else {
					e = E.Right[string](i)
				}
				if !yield(e) {
					return
				}
			}
		}

		_ = collectEithers(TakeUntilLeft(seq))

		// Should consume: 0, 1, 2, 3 (the Left), then stop
		assert.Equal(t, 4, consumed)
	})

	t.Run("consumes all when no Left", func(t *testing.T) {
		consumed := 0
		seq := func(yield func(Either[string, int]) bool) {
			for i := range 5 {
				consumed++
				if !yield(E.Right[string](i)) {
					return
				}
			}
		}

		_ = collectEithers(TakeUntilLeft(seq))

		assert.Equal(t, 5, consumed)
	})
}

// TestTakeUntilLeft_EdgeCases tests edge cases
func TestTakeUntilLeft_EdgeCases(t *testing.T) {
	t.Run("handles multiple consecutive Lefts", func(t *testing.T) {
		seq := iter.From(
			E.Right[string](1),
			E.Left[int]("first error"),
			E.Left[int]("second error"),
			E.Right[string](2),
		)
		result := collectEithers(TakeUntilLeft(seq))
		expected := []Either[string, int]{
			E.Right[string](1),
			E.Left[int]("first error"),
		}
		assert.Equal(t, expected, result)
	})

	t.Run("handles zero values", func(t *testing.T) {
		seq := iter.From(
			E.Right[string](0),
			E.Right[string](0),
			E.Left[int](""),
			E.Right[string](0),
		)
		result := collectEithers(TakeUntilLeft(seq))
		expected := []Either[string, int]{
			E.Right[string](0),
			E.Right[string](0),
			E.Left[int](""),
		}
		assert.Equal(t, expected, result)
	})

	t.Run("handles nil error values", func(t *testing.T) {
		var nilErr error
		seq := iter.From(
			E.Right[error](1),
			E.Right[error](2),
			E.Left[int](nilErr),
			E.Right[error](3),
		)
		result := collectEithers(TakeUntilLeft(seq))
		expected := []Either[error, int]{
			E.Right[error](1),
			E.Right[error](2),
			E.Left[int](nilErr),
		}
		assert.Equal(t, expected, result)
	})
}

// BenchmarkTakeUntilLeft benchmarks basic TakeUntilLeft operation
func BenchmarkTakeUntilLeft(b *testing.B) {
	seq := iter.From(
		E.Right[string](1),
		E.Right[string](2),
		E.Right[string](3),
		E.Left[int]("error"),
		E.Right[string](4),
	)

	b.ResetTimer()
	for range b.N {
		_ = collectEithers(TakeUntilLeft(seq))
	}
}

// BenchmarkTakeUntilLeft_LargeSequence benchmarks with a large sequence
func BenchmarkTakeUntilLeft_LargeSequence(b *testing.B) {
	makeSeq := func() SeqEither[string, int] {
		return func(yield func(Either[string, int]) bool) {
			for i := range 1000 {
				var e Either[string, int]
				if i == 500 {
					e = E.Left[int]("error at 500")
				} else {
					e = E.Right[string](i)
				}
				if !yield(e) {
					return
				}
			}
		}
	}

	b.ResetTimer()
	for range b.N {
		_ = collectEithers(TakeUntilLeft(makeSeq()))
	}
}

// BenchmarkTakeUntilLeft_AllRight benchmarks when all elements are Right
func BenchmarkTakeUntilLeft_AllRight(b *testing.B) {
	makeSeq := func() SeqEither[string, int] {
		return func(yield func(Either[string, int]) bool) {
			for i := range 100 {
				if !yield(E.Right[string](i)) {
					return
				}
			}
		}
	}

	b.ResetTimer()
	for range b.N {
		_ = collectEithers(TakeUntilLeft(makeSeq()))
	}
}

// BenchmarkTakeUntilLeft_ImmediateLeft benchmarks when first element is Left
func BenchmarkTakeUntilLeft_ImmediateLeft(b *testing.B) {
	makeSeq := func() SeqEither[string, int] {
		return func(yield func(Either[string, int]) bool) {
			if !yield(E.Left[int]("immediate error")) {
				return
			}
			for i := range 100 {
				if !yield(E.Right[string](i)) {
					return
				}
			}
		}
	}

	b.ResetTimer()
	for range b.N {
		_ = collectEithers(TakeUntilLeft(makeSeq()))
	}
}
