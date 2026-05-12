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

package iterresult

import (
	"errors"
	"strconv"
	"testing"

	"github.com/IBM/fp-go/v2/iterator/iter"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// Helper function to collect results into a slice
func collectResults[T any](seq SeqResult[T]) []Result[T] {
	var results []Result[T]
	for r := range seq {
		results = append(results, r)
	}
	return results
}

// TestTakeUntilError_Success tests basic TakeUntilError functionality
func TestTakeUntilError_Success(t *testing.T) {
	t.Run("takes all success values until first error", func(t *testing.T) {
		seq := iter.From(
			result.Of(1),
			result.Of(2),
			result.Of(3),
			result.Left[int](errors.New("error")),
			result.Of(4),
			result.Of(5),
		)
		collected := collectResults(TakeUntilError(seq))

		assert.Len(t, collected, 4)
		assert.True(t, result.IsRight(collected[0]))
		assert.True(t, result.IsRight(collected[1]))
		assert.True(t, result.IsRight(collected[2]))
		assert.True(t, result.IsLeft(collected[3]))
	})

	t.Run("takes all elements when no error exists", func(t *testing.T) {
		seq := iter.From(
			result.Of(1),
			result.Of(2),
			result.Of(3),
		)
		collected := collectResults(TakeUntilError(seq))

		assert.Len(t, collected, 3)
		for _, r := range collected {
			assert.True(t, result.IsRight(r))
		}
	})

	t.Run("stops immediately when first element is error", func(t *testing.T) {
		seq := iter.From(
			result.Left[int](errors.New("immediate error")),
			result.Of(1),
			result.Of(2),
		)
		collected := collectResults(TakeUntilError(seq))

		assert.Len(t, collected, 1)
		assert.True(t, result.IsLeft(collected[0]))
	})

	t.Run("handles single success element", func(t *testing.T) {
		seq := iter.From(result.Of(42))
		collected := collectResults(TakeUntilError(seq))

		assert.Len(t, collected, 1)
		assert.True(t, result.IsRight(collected[0]))
	})

	t.Run("handles single error element", func(t *testing.T) {
		seq := iter.From(result.Left[int](errors.New("error")))
		collected := collectResults(TakeUntilError(seq))

		assert.Len(t, collected, 1)
		assert.True(t, result.IsLeft(collected[0]))
	})
}

// TestTakeUntilError_Empty tests TakeUntilError with empty sequences
func TestTakeUntilError_Empty(t *testing.T) {
	t.Run("returns empty from empty sequence", func(t *testing.T) {
		seq := iter.Empty[Result[int]]()
		collected := collectResults(TakeUntilError(seq))
		assert.Empty(t, collected)
	})
}

// TestTakeUntilError_WithComplexTypes tests TakeUntilError with complex types
func TestTakeUntilError_WithComplexTypes(t *testing.T) {
	type User struct {
		ID   int
		Name string
	}

	t.Run("works with struct types", func(t *testing.T) {
		seq := iter.From(
			result.Of(User{1, "Alice"}),
			result.Of(User{2, "Bob"}),
			result.Left[User](errors.New("database error")),
			result.Of(User{3, "Charlie"}),
		)
		collected := collectResults(TakeUntilError(seq))

		assert.Len(t, collected, 3)
		assert.True(t, result.IsRight(collected[0]))
		assert.True(t, result.IsRight(collected[1]))
		assert.True(t, result.IsLeft(collected[2]))
	})

	t.Run("works with string types", func(t *testing.T) {
		seq := iter.From(
			result.Of("hello"),
			result.Of("world"),
			result.Left[string](errors.New("parse error")),
			result.Of("ignored"),
		)
		collected := collectResults(TakeUntilError(seq))

		assert.Len(t, collected, 3)
		assert.True(t, result.IsRight(collected[0]))
		assert.True(t, result.IsRight(collected[1]))
		assert.True(t, result.IsLeft(collected[2]))
	})

	t.Run("works with slice types", func(t *testing.T) {
		seq := iter.From(
			result.Of([]int{1, 2}),
			result.Of([]int{3, 4}),
			result.Left[[]int](errors.New("validation error")),
			result.Of([]int{5, 6}),
		)
		collected := collectResults(TakeUntilError(seq))

		assert.Len(t, collected, 3)
		assert.True(t, result.IsRight(collected[0]))
		assert.True(t, result.IsRight(collected[1]))
		assert.True(t, result.IsLeft(collected[2]))
	})
}

// TestTakeUntilError_WithChainedOperations tests TakeUntilError with other operations
func TestTakeUntilError_WithChainedOperations(t *testing.T) {
	t.Run("chains with Map", func(t *testing.T) {
		seq := iter.From(
			result.Of(1),
			result.Of(2),
			result.Left[int](errors.New("error")),
			result.Of(3),
		)
		collected := collectResults(
			TakeUntilError(
				MonadMap(seq, func(n int) int { return n * 2 }),
			),
		)

		assert.Len(t, collected, 3)
		assert.True(t, result.IsRight(collected[0]))
		assert.True(t, result.IsRight(collected[1]))
		assert.True(t, result.IsLeft(collected[2]))
	})

	t.Run("chains with multiple TakeUntilError", func(t *testing.T) {
		seq := iter.From(
			result.Of(1),
			result.Of(2),
			result.Left[int](errors.New("first error")),
			result.Of(3),
			result.Left[int](errors.New("second error")),
		)
		// First TakeUntilError stops at "first error"
		// Second TakeUntilError would also stop at "first error" (same result)
		collected := collectResults(TakeUntilError(TakeUntilError(seq)))

		assert.Len(t, collected, 3)
		assert.True(t, result.IsRight(collected[0]))
		assert.True(t, result.IsRight(collected[1]))
		assert.True(t, result.IsLeft(collected[2]))
	})
}

// TestTakeUntilError_WithRealWorldScenarios tests practical use cases
func TestTakeUntilError_WithRealWorldScenarios(t *testing.T) {
	t.Run("parsing numbers until error", func(t *testing.T) {
		inputs := []string{"1", "2", "3", "invalid", "4", "5"}
		seq := iter.MonadMap(
			iter.From(inputs...),
			result.Eitherize1(strconv.Atoi),
		)
		collected := collectResults(TakeUntilError(seq))

		assert.Len(t, collected, 4) // 1, 2, 3, and the error
		assert.True(t, result.IsRight(collected[0]))
		assert.True(t, result.IsRight(collected[1]))
		assert.True(t, result.IsRight(collected[2]))
		assert.True(t, result.IsLeft(collected[3]))
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

		validate := func(r Record) Result[Record] {
			if r.Valid {
				return result.Of(r)
			}
			return result.Left[Record](errors.New("invalid record"))
		}

		seq := iter.MonadMap(iter.From(records...), validate)
		collected := collectResults(TakeUntilError(seq))

		assert.Len(t, collected, 3)
		assert.True(t, result.IsRight(collected[0]))
		assert.True(t, result.IsRight(collected[1]))
		assert.True(t, result.IsLeft(collected[2]))
	})
}

// TestTakeUntilError_DoesNotConsumeEntireSequence verifies lazy evaluation
func TestTakeUntilError_DoesNotConsumeEntireSequence(t *testing.T) {
	t.Run("stops consuming after error", func(t *testing.T) {
		consumed := 0
		seq := func(yield func(Result[int]) bool) {
			for i := range 10 {
				consumed++
				var r Result[int]
				if i == 3 {
					r = result.Left[int](errors.New("error at 3"))
				} else {
					r = result.Of(i)
				}
				if !yield(r) {
					return
				}
			}
		}

		_ = collectResults(TakeUntilError(seq))

		// Should consume: 0, 1, 2, 3 (the error), then stop
		assert.Equal(t, 4, consumed)
	})

	t.Run("consumes all when no error", func(t *testing.T) {
		consumed := 0
		seq := func(yield func(Result[int]) bool) {
			for i := range 5 {
				consumed++
				if !yield(result.Of(i)) {
					return
				}
			}
		}

		_ = collectResults(TakeUntilError(seq))

		assert.Equal(t, 5, consumed)
	})
}

// TestTakeUntilError_EdgeCases tests edge cases
func TestTakeUntilError_EdgeCases(t *testing.T) {
	t.Run("handles multiple consecutive errors", func(t *testing.T) {
		seq := iter.From(
			result.Of(1),
			result.Left[int](errors.New("first error")),
			result.Left[int](errors.New("second error")),
			result.Of(2),
		)
		collected := collectResults(TakeUntilError(seq))

		assert.Len(t, collected, 2)
		assert.True(t, result.IsRight(collected[0]))
		assert.True(t, result.IsLeft(collected[1]))
	})

	t.Run("handles zero values", func(t *testing.T) {
		seq := iter.From(
			result.Of(0),
			result.Of(0),
			result.Left[int](errors.New("error")),
			result.Of(0),
		)
		collected := collectResults(TakeUntilError(seq))

		assert.Len(t, collected, 3)
		assert.True(t, result.IsRight(collected[0]))
		assert.True(t, result.IsRight(collected[1]))
		assert.True(t, result.IsLeft(collected[2]))
	})

	t.Run("handles nil error values", func(t *testing.T) {
		var nilErr error
		seq := iter.From(
			result.Of(1),
			result.Of(2),
			result.Left[int](nilErr),
			result.Of(3),
		)
		collected := collectResults(TakeUntilError(seq))

		assert.Len(t, collected, 3)
		assert.True(t, result.IsRight(collected[0]))
		assert.True(t, result.IsRight(collected[1]))
		assert.True(t, result.IsLeft(collected[2]))
	})
}

// TestStopOnError_Success tests that StopOnError is an alias for TakeUntilError
func TestStopOnError_Success(t *testing.T) {
	t.Run("behaves identically to TakeUntilError", func(t *testing.T) {
		seq := iter.From(
			result.Of(1),
			result.Of(2),
			result.Left[int](errors.New("error")),
			result.Of(3),
		)

		resultTake := collectResults(TakeUntilError(seq))
		resultStop := collectResults(StopOnError(seq))

		assert.Equal(t, len(resultTake), len(resultStop))
		for i := range resultTake {
			assert.Equal(t, result.IsRight(resultTake[i]), result.IsRight(resultStop[i]))
			assert.Equal(t, result.IsLeft(resultTake[i]), result.IsLeft(resultStop[i]))
		}
	})

	t.Run("stops on first error", func(t *testing.T) {
		seq := iter.From(
			result.Of(1),
			result.Of(2),
			result.Left[int](errors.New("error")),
			result.Of(3),
		)
		collected := collectResults(StopOnError(seq))

		assert.Len(t, collected, 3)
		assert.True(t, result.IsRight(collected[0]))
		assert.True(t, result.IsRight(collected[1]))
		assert.True(t, result.IsLeft(collected[2]))
	})

	t.Run("passes through all success values", func(t *testing.T) {
		seq := iter.From(
			result.Of(1),
			result.Of(2),
			result.Of(3),
		)
		collected := collectResults(StopOnError(seq))

		assert.Len(t, collected, 3)
		for _, r := range collected {
			assert.True(t, result.IsRight(r))
		}
	})
}

// BenchmarkTakeUntilError benchmarks basic TakeUntilError operation
func BenchmarkTakeUntilError(b *testing.B) {
	seq := iter.From(
		result.Of(1),
		result.Of(2),
		result.Of(3),
		result.Left[int](errors.New("error")),
		result.Of(4),
	)

	b.ResetTimer()
	for range b.N {
		_ = collectResults(TakeUntilError(seq))
	}
}

// BenchmarkTakeUntilError_LargeSequence benchmarks with a large sequence
func BenchmarkTakeUntilError_LargeSequence(b *testing.B) {
	makeSeq := func() SeqResult[int] {
		return func(yield func(Result[int]) bool) {
			for i := range 1000 {
				var r Result[int]
				if i == 500 {
					r = result.Left[int](errors.New("error at 500"))
				} else {
					r = result.Of(i)
				}
				if !yield(r) {
					return
				}
			}
		}
	}

	b.ResetTimer()
	for range b.N {
		_ = collectResults(TakeUntilError(makeSeq()))
	}
}

// BenchmarkTakeUntilError_AllSuccess benchmarks when all elements are success
func BenchmarkTakeUntilError_AllSuccess(b *testing.B) {
	makeSeq := func() SeqResult[int] {
		return func(yield func(Result[int]) bool) {
			for i := range 100 {
				if !yield(result.Of(i)) {
					return
				}
			}
		}
	}

	b.ResetTimer()
	for range b.N {
		_ = collectResults(TakeUntilError(makeSeq()))
	}
}

// BenchmarkTakeUntilError_ImmediateError benchmarks when first element is error
func BenchmarkTakeUntilError_ImmediateError(b *testing.B) {
	makeSeq := func() SeqResult[int] {
		return func(yield func(Result[int]) bool) {
			if !yield(result.Left[int](errors.New("immediate error"))) {
				return
			}
			for i := range 100 {
				if !yield(result.Of(i)) {
					return
				}
			}
		}
	}

	b.ResetTimer()
	for range b.N {
		_ = collectResults(TakeUntilError(makeSeq()))
	}
}

// BenchmarkStopOnError benchmarks StopOnError operation
func BenchmarkStopOnError(b *testing.B) {
	seq := iter.From(
		result.Of(1),
		result.Of(2),
		result.Of(3),
		result.Left[int](errors.New("error")),
		result.Of(4),
	)

	b.ResetTimer()
	for range b.N {
		_ = collectResults(StopOnError(seq))
	}
}

// Made with Bob
