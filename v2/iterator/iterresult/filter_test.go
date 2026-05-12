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
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/iterator/iter"
	N "github.com/IBM/fp-go/v2/number"
	R "github.com/IBM/fp-go/v2/result"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

func TestFilterOrElse_PredicateTrue(t *testing.T) {
	isPositive := N.MoreThan(0)
	onFalse := func(n int) error { return fmt.Errorf("%d is not positive", n) }

	filter := FilterOrElse(isPositive, onFalse)
	seq := iter.From(R.Of(42))
	result := collectResults(filter(seq))

	assert.Equal(t, []R.Result[int]{R.Of(42)}, result)
}

func TestFilterOrElse_PredicateFalse(t *testing.T) {
	isPositive := N.MoreThan(0)
	onFalse := func(n int) error { return fmt.Errorf("%d is not positive", n) }

	filter := FilterOrElse(isPositive, onFalse)
	seq := iter.From(R.Of(-5))
	result := collectResults(filter(seq))

	assert.Len(t, result, 1)
	assert.True(t, R.IsLeft(result[0]))
	err := F.Pipe1(result[0], R.Fold(
		F.Identity[error],
		func(_ int) error { t.Fatal("expected error"); return nil },
	))
	assert.Equal(t, "-5 is not positive", err.Error())
}

func TestFilterOrElse_ErrorPassesThrough(t *testing.T) {
	isPositive := N.MoreThan(0)
	onFalse := func(n int) error { return fmt.Errorf("%d is not positive", n) }

	filter := FilterOrElse(isPositive, onFalse)
	originalErr := errors.New("original error")
	seq := iter.From(R.Left[int](originalErr))
	result := collectResults(filter(seq))

	assert.Len(t, result, 1)
	assert.True(t, R.IsLeft(result[0]))
	err := F.Pipe1(result[0], R.Fold(
		F.Identity[error],
		func(_ int) error { t.Fatal("expected error"); return nil },
	))
	assert.Equal(t, originalErr, err)
}

func TestFilterOrElse_ZeroValue(t *testing.T) {
	isNonZero := func(n int) bool { return n != 0 }
	onZero := func(n int) error { return errors.New("value is zero") }

	filter := FilterOrElse(isNonZero, onZero)
	seq := iter.From(R.Of(0))
	result := collectResults(filter(seq))

	assert.Len(t, result, 1)
	assert.True(t, R.IsLeft(result[0]))
}

func TestFilterOrElse_StringValidation(t *testing.T) {
	isNonEmpty := S.IsNonEmpty
	onEmpty := func(s string) error { return fmt.Errorf("string is empty") }

	filter := FilterOrElse(isNonEmpty, onEmpty)

	t.Run("non-empty string passes", func(t *testing.T) {
		seq := iter.From(R.Of("hello"))
		result := collectResults(filter(seq))
		assert.Equal(t, []R.Result[string]{R.Of("hello")}, result)
	})

	t.Run("empty string fails", func(t *testing.T) {
		seq := iter.From(R.Of(""))
		result := collectResults(filter(seq))
		assert.Len(t, result, 1)
		assert.True(t, R.IsLeft(result[0]))
	})
}

func TestFilterOrElse_MultipleValues(t *testing.T) {
	isEven := func(n int) bool { return n%2 == 0 }
	onOdd := func(n int) error { return fmt.Errorf("%d is odd", n) }

	filter := FilterOrElse(isEven, onOdd)
	seq := iter.From(
		R.Of(2),
		R.Of(3),
		R.Of(4),
		R.Left[int](errors.New("error")),
		R.Of(5),
	)
	result := collectResults(filter(seq))

	assert.Len(t, result, 5)
	assert.True(t, R.IsRight(result[0]))
	assert.True(t, R.IsLeft(result[1]))
	assert.True(t, R.IsRight(result[2]))
	assert.True(t, R.IsLeft(result[3]))
	assert.True(t, R.IsLeft(result[4]))
}

func TestFilterOrElse_InPipeline(t *testing.T) {
	isPositive := N.MoreThan(0)
	onNegative := func(n int) error { return fmt.Errorf("%d is not positive", n) }

	result := F.Pipe2(
		iter.From(1, -2, 3, -4, 5),
		FromSeq,
		FilterOrElse(isPositive, onNegative),
	)

	collected := collectResults(result)
	assert.Len(t, collected, 5)
	assert.True(t, R.IsRight(collected[0]))
	assert.True(t, R.IsLeft(collected[1]))
	assert.True(t, R.IsRight(collected[2]))
	assert.True(t, R.IsLeft(collected[3]))
	assert.True(t, R.IsRight(collected[4]))
}

func TestFilterOrElse_WithChain(t *testing.T) {
	isPositive := N.MoreThan(0)
	onNegative := func(n int) error { return fmt.Errorf("%d is not positive", n) }

	result := F.Pipe3(
		iter.From(1, 2, 3),
		FromSeq,
		Map(func(n int) int { return n - 2 }),
		FilterOrElse(isPositive, onNegative),
	)

	collected := collectResults(result)
	assert.Len(t, collected, 3)
	assert.True(t, R.IsLeft(collected[0]))
	assert.True(t, R.IsLeft(collected[1]))
	assert.True(t, R.IsRight(collected[2]))
}

func TestFilterOrElse_CustomPredicate(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	isAdult := func(u User) bool { return u.Age >= 18 }
	onMinor := func(u User) error {
		return fmt.Errorf("%s is only %d years old", u.Name, u.Age)
	}

	filter := FilterOrElse(isAdult, onMinor)
	seq := iter.From(
		R.Of(User{"Alice", 25}),
		R.Of(User{"Bob", 16}),
		R.Of(User{"Charlie", 30}),
	)

	result := collectResults(filter(seq))
	assert.Len(t, result, 3)
	assert.True(t, R.IsRight(result[0]))
	assert.True(t, R.IsLeft(result[1]))
	assert.True(t, R.IsRight(result[2]))
}

// TestMonadFilter_Success tests basic MonadFilter functionality
func TestMonadFilter_Success(t *testing.T) {
	t.Run("filters Ok values based on predicate", func(t *testing.T) {
		seq := iter.From(
			R.Of(1),
			R.Of(2),
			R.Of(3),
			R.Of(4),
			R.Of(5),
		)
		isEven := func(n int) bool { return n%2 == 0 }
		result := collectResults(MonadFilter(seq, isEven))

		expected := []R.Result[int]{
			R.Of(2),
			R.Of(4),
		}
		assert.Equal(t, expected, result)
	})

	t.Run("keeps all Ok values when all satisfy predicate", func(t *testing.T) {
		seq := iter.From(
			R.Of(2),
			R.Of(4),
			R.Of(6),
		)
		isEven := func(n int) bool { return n%2 == 0 }
		result := collectResults(MonadFilter(seq, isEven))

		expected := []R.Result[int]{
			R.Of(2),
			R.Of(4),
			R.Of(6),
		}
		assert.Equal(t, expected, result)
	})

	t.Run("filters out all Ok values when none satisfy predicate", func(t *testing.T) {
		seq := iter.From(
			R.Of(1),
			R.Of(3),
			R.Of(5),
		)
		isEven := func(n int) bool { return n%2 == 0 }
		result := collectResults(MonadFilter(seq, isEven))

		assert.Empty(t, result)
	})

	t.Run("passes through error values unchanged", func(t *testing.T) {
		err1 := errors.New("error1")
		err2 := errors.New("error2")
		seq := iter.From(
			R.Of(1),
			R.Left[int](err1),
			R.Of(2),
			R.Left[int](err2),
			R.Of(3),
		)
		isEven := func(n int) bool { return n%2 == 0 }
		result := collectResults(MonadFilter(seq, isEven))

		assert.Len(t, result, 3)
		assert.True(t, R.IsLeft(result[0]))
		assert.True(t, R.IsRight(result[1]))
		assert.True(t, R.IsLeft(result[2]))
	})
}

// TestMonadFilter_Empty tests MonadFilter with empty sequences
func TestMonadFilter_Empty(t *testing.T) {
	t.Run("returns empty from empty sequence", func(t *testing.T) {
		seq := iter.Empty[R.Result[int]]()
		isEven := func(n int) bool { return n%2 == 0 }
		result := collectResults(MonadFilter(seq, isEven))

		assert.Empty(t, result)
	})
}

// TestMonadFilter_EdgeCases tests edge cases
func TestMonadFilter_EdgeCases(t *testing.T) {
	t.Run("handles zero values", func(t *testing.T) {
		seq := iter.From(
			R.Of(0),
			R.Of(1),
			R.Of(2),
		)
		isZero := func(n int) bool { return n == 0 }
		result := collectResults(MonadFilter(seq, isZero))

		expected := []R.Result[int]{R.Of(0)}
		assert.Equal(t, expected, result)
	})

	t.Run("handles single Ok element that passes", func(t *testing.T) {
		seq := iter.From(R.Of(2))
		isEven := func(n int) bool { return n%2 == 0 }
		result := collectResults(MonadFilter(seq, isEven))

		expected := []R.Result[int]{R.Of(2)}
		assert.Equal(t, expected, result)
	})

	t.Run("handles single Ok element that fails", func(t *testing.T) {
		seq := iter.From(R.Of(1))
		isEven := func(n int) bool { return n%2 == 0 }
		result := collectResults(MonadFilter(seq, isEven))

		assert.Empty(t, result)
	})

	t.Run("handles single error element", func(t *testing.T) {
		err := errors.New("error")
		seq := iter.From(R.Left[int](err))
		isEven := func(n int) bool { return n%2 == 0 }
		result := collectResults(MonadFilter(seq, isEven))

		assert.Len(t, result, 1)
		assert.True(t, R.IsLeft(result[0]))
	})

	t.Run("handles sequence of only error values", func(t *testing.T) {
		err1 := errors.New("error1")
		err2 := errors.New("error2")
		err3 := errors.New("error3")
		seq := iter.From(
			R.Left[int](err1),
			R.Left[int](err2),
			R.Left[int](err3),
		)
		isEven := func(n int) bool { return n%2 == 0 }
		result := collectResults(MonadFilter(seq, isEven))

		assert.Len(t, result, 3)
		assert.True(t, R.IsLeft(result[0]))
		assert.True(t, R.IsLeft(result[1]))
		assert.True(t, R.IsLeft(result[2]))
	})
}

// TestMonadFilter_WithComplexTypes tests MonadFilter with complex types
func TestMonadFilter_WithComplexTypes(t *testing.T) {
	type User struct {
		ID   int
		Name string
		Age  int
	}

	t.Run("filters struct types", func(t *testing.T) {
		seq := iter.From(
			R.Of(User{1, "Alice", 25}),
			R.Of(User{2, "Bob", 17}),
			R.Of(User{3, "Charlie", 30}),
			R.Left[User](errors.New("database error")),
			R.Of(User{4, "David", 16}),
		)
		isAdult := func(u User) bool { return u.Age >= 18 }
		result := collectResults(MonadFilter(seq, isAdult))

		assert.Len(t, result, 3)
		assert.True(t, R.IsRight(result[0]))
		assert.True(t, R.IsRight(result[1]))
		assert.True(t, R.IsLeft(result[2]))
	})

	t.Run("filters string types", func(t *testing.T) {
		seq := iter.From(
			R.Of("hello"),
			R.Of(""),
			R.Of("world"),
			R.Left[string](errors.New("error")),
			R.Of("test"),
		)
		isNonEmpty := S.IsNonEmpty
		result := collectResults(MonadFilter(seq, isNonEmpty))

		assert.Len(t, result, 4)
		assert.True(t, R.IsRight(result[0]))
		assert.True(t, R.IsRight(result[1]))
		assert.True(t, R.IsLeft(result[2]))
		assert.True(t, R.IsRight(result[3]))
	})
}

// TestFilter_Success tests the curried Filter function
func TestFilter_Success(t *testing.T) {
	t.Run("creates reusable filter function", func(t *testing.T) {
		isEven := func(n int) bool { return n%2 == 0 }
		evens := Filter(isEven)

		seq1 := iter.From(
			R.Of(1),
			R.Of(2),
			R.Of(3),
		)
		result1 := collectResults(evens(seq1))
		expected1 := []R.Result[int]{R.Of(2)}
		assert.Equal(t, expected1, result1)

		seq2 := iter.From(
			R.Of(4),
			R.Left[int](errors.New("error")),
			R.Of(5),
			R.Of(6),
		)
		result2 := collectResults(evens(seq2))
		assert.Len(t, result2, 3)
		assert.True(t, R.IsRight(result2[0]))
		assert.True(t, R.IsLeft(result2[1]))
		assert.True(t, R.IsRight(result2[2]))
	})

	t.Run("works in pipeline", func(t *testing.T) {
		isPositive := N.MoreThan(0)
		result := F.Pipe2(
			iter.From(
				R.Of(-1),
				R.Of(2),
				R.Left[int](errors.New("error")),
				R.Of(3),
				R.Of(-4),
			),
			Filter(isPositive),
			collectResults[int],
		)

		assert.Len(t, result, 3)
		assert.True(t, R.IsRight(result[0]))
		assert.True(t, R.IsLeft(result[1]))
		assert.True(t, R.IsRight(result[2]))
	})

	t.Run("chains with Map", func(t *testing.T) {
		isEven := func(n int) bool { return n%2 == 0 }
		result := F.Pipe3(
			iter.From(1, 2, 3, 4, 5),
			FromSeq,
			Map(func(n int) int { return n * 2 }),
			Filter(isEven),
		)

		collected := collectResults(result)
		expected := []R.Result[int]{
			R.Of(2),
			R.Of(4),
			R.Of(6),
			R.Of(8),
			R.Of(10),
		}
		assert.Equal(t, expected, collected)
	})
}

// TestFilter_VsFilterOrElse compares Filter and FilterOrElse behavior
func TestFilter_VsFilterOrElse(t *testing.T) {
	t.Run("Filter removes failing values", func(t *testing.T) {
		seq := iter.From(
			R.Of(1),
			R.Of(2),
			R.Of(3),
		)
		isEven := func(n int) bool { return n%2 == 0 }

		// Filter removes odd numbers
		filterResult := collectResults(Filter(isEven)(seq))
		assert.Len(t, filterResult, 1)
		assert.Equal(t, R.Of(2), filterResult[0])
	})

	t.Run("FilterOrElse converts failing values to errors", func(t *testing.T) {
		seq := iter.From(
			R.Of(1),
			R.Of(2),
			R.Of(3),
		)
		isEven := func(n int) bool { return n%2 == 0 }
		onOdd := func(n int) error { return fmt.Errorf("%d is odd", n) }

		// FilterOrElse converts odd numbers to errors
		filterOrElseResult := collectResults(FilterOrElse(isEven, onOdd)(seq))
		assert.Len(t, filterOrElseResult, 3)
		assert.True(t, R.IsLeft(filterOrElseResult[0]))
		assert.True(t, R.IsRight(filterOrElseResult[1]))
		assert.True(t, R.IsLeft(filterOrElseResult[2]))
	})
}

// BenchmarkMonadFilter benchmarks basic MonadFilter operation
func BenchmarkMonadFilter(b *testing.B) {
	seq := iter.From(
		R.Of(1),
		R.Of(2),
		R.Of(3),
		R.Of(4),
		R.Of(5),
	)
	isEven := func(n int) bool { return n%2 == 0 }

	b.ResetTimer()
	for range b.N {
		_ = collectResults(MonadFilter(seq, isEven))
	}
}

// BenchmarkFilter benchmarks curried Filter operation
func BenchmarkFilter(b *testing.B) {
	seq := iter.From(
		R.Of(1),
		R.Of(2),
		R.Of(3),
		R.Of(4),
		R.Of(5),
	)
	isEven := func(n int) bool { return n%2 == 0 }
	evens := Filter(isEven)

	b.ResetTimer()
	for range b.N {
		_ = collectResults(evens(seq))
	}
}

// BenchmarkMonadFilter_LargeSequence benchmarks with a large sequence
func BenchmarkMonadFilter_LargeSequence(b *testing.B) {
	makeSeq := func() SeqResult[int] {
		return func(yield func(R.Result[int]) bool) {
			for i := range 1000 {
				if !yield(R.Of(i)) {
					return
				}
			}
		}
	}
	isEven := func(n int) bool { return n%2 == 0 }

	b.ResetTimer()
	for range b.N {
		_ = collectResults(MonadFilter(makeSeq(), isEven))
	}
}

// Made with Bob
