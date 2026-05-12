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
	"fmt"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/iterator/iter"
	N "github.com/IBM/fp-go/v2/number"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

func TestFilterOrElse_PredicateTrue(t *testing.T) {
	isPositive := N.MoreThan(0)
	onFalse := S.Format[int]("%d is not positive")

	filter := FilterOrElse(isPositive, onFalse)
	seq := iter.From(E.Right[string](42))
	result := collectEithers(filter(seq))

	assert.Equal(t, []Either[string, int]{E.Right[string](42)}, result)
}

func TestFilterOrElse_PredicateFalse(t *testing.T) {
	isPositive := N.MoreThan(0)
	onFalse := S.Format[int]("%d is not positive")

	filter := FilterOrElse(isPositive, onFalse)
	seq := iter.From(E.Right[string](-5))
	result := collectEithers(filter(seq))

	assert.Equal(t, []Either[string, int]{E.Left[int]("-5 is not positive")}, result)
}

func TestFilterOrElse_LeftPassesThrough(t *testing.T) {
	isPositive := N.MoreThan(0)
	onFalse := S.Format[int]("%d is not positive")

	filter := FilterOrElse(isPositive, onFalse)
	seq := iter.From(E.Left[int]("original error"))
	result := collectEithers(filter(seq))

	assert.Equal(t, []Either[string, int]{E.Left[int]("original error")}, result)
}

func TestFilterOrElse_ZeroValue(t *testing.T) {
	isNonZero := func(n int) bool { return n != 0 }
	onZero := func(n int) string { return "value is zero" }

	filter := FilterOrElse(isNonZero, onZero)
	seq := iter.From(E.Right[string](0))
	result := collectEithers(filter(seq))

	assert.Equal(t, []Either[string, int]{E.Left[int]("value is zero")}, result)
}

func TestFilterOrElse_StringValidation(t *testing.T) {
	isNonEmpty := S.IsNonEmpty
	onEmpty := func(s string) error { return fmt.Errorf("string is empty") }

	filter := FilterOrElse(isNonEmpty, onEmpty)

	t.Run("non-empty string passes", func(t *testing.T) {
		seq := iter.From(E.Right[error]("hello"))
		result := collectEithers(filter(seq))
		assert.Equal(t, []Either[error, string]{E.Right[error]("hello")}, result)
	})

	t.Run("empty string fails", func(t *testing.T) {
		seq := iter.From(E.Right[error](""))
		result := collectEithers(filter(seq))
		assert.Len(t, result, 1)
		assert.True(t, E.IsLeft(result[0]))
	})
}

func TestFilterOrElse_MultipleValues(t *testing.T) {
	isEven := func(n int) bool { return n%2 == 0 }
	onOdd := S.Format[int]("%d is odd")

	filter := FilterOrElse(isEven, onOdd)
	seq := iter.From(
		E.Right[string](2),
		E.Right[string](3),
		E.Right[string](4),
		E.Left[int]("error"),
		E.Right[string](5),
	)
	result := collectEithers(filter(seq))

	expected := []Either[string, int]{
		E.Right[string](2),
		E.Left[int]("3 is odd"),
		E.Right[string](4),
		E.Left[int]("error"),
		E.Left[int]("5 is odd"),
	}
	assert.Equal(t, expected, result)
}

func TestFilterOrElse_InPipeline(t *testing.T) {
	isPositive := N.MoreThan(0)
	onNegative := S.Format[int]("%d is not positive")

	result := F.Pipe2(
		iter.From(1, -2, 3, -4, 5),
		FromSeq[string],
		FilterOrElse(isPositive, onNegative),
	)

	collected := collectEithers(result)
	expected := []Either[string, int]{
		E.Right[string](1),
		E.Left[int]("-2 is not positive"),
		E.Right[string](3),
		E.Left[int]("-4 is not positive"),
		E.Right[string](5),
	}
	assert.Equal(t, expected, collected)
}

func TestFilterOrElse_WithChain(t *testing.T) {
	isPositive := N.MoreThan(0)
	onNegative := S.Format[int]("%d is not positive")

	result := F.Pipe3(
		iter.From(1, 2, 3),
		FromSeq[string],
		Map[string](func(n int) int { return n - 2 }),
		FilterOrElse(isPositive, onNegative),
	)

	collected := collectEithers(result)
	expected := []Either[string, int]{
		E.Left[int]("-1 is not positive"),
		E.Left[int]("0 is not positive"),
		E.Right[string](1),
	}
	assert.Equal(t, expected, collected)
}

func TestFilterOrElse_CustomPredicate(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	isAdult := func(u User) bool { return u.Age >= 18 }
	onMinor := func(u User) string {
		return fmt.Sprintf("%s is only %d years old", u.Name, u.Age)
	}

	filter := FilterOrElse(isAdult, onMinor)
	seq := iter.From(
		E.Right[string](User{"Alice", 25}),
		E.Right[string](User{"Bob", 16}),
		E.Right[string](User{"Charlie", 30}),
	)

	result := collectEithers(filter(seq))
	assert.Len(t, result, 3)
	assert.True(t, E.IsRight(result[0]))
	assert.True(t, E.IsLeft(result[1]))
	assert.True(t, E.IsRight(result[2]))
}

// TestMonadFilter_Success tests basic MonadFilter functionality
func TestMonadFilter_Success(t *testing.T) {
	t.Run("filters Right values based on predicate", func(t *testing.T) {
		seq := iter.From(
			E.Right[string](1),
			E.Right[string](2),
			E.Right[string](3),
			E.Right[string](4),
			E.Right[string](5),
		)
		isEven := func(n int) bool { return n%2 == 0 }
		result := collectEithers(MonadFilter(seq, isEven))

		expected := []Either[string, int]{
			E.Right[string](2),
			E.Right[string](4),
		}
		assert.Equal(t, expected, result)
	})

	t.Run("keeps all Right values when all satisfy predicate", func(t *testing.T) {
		seq := iter.From(
			E.Right[string](2),
			E.Right[string](4),
			E.Right[string](6),
		)
		isEven := func(n int) bool { return n%2 == 0 }
		result := collectEithers(MonadFilter(seq, isEven))

		expected := []Either[string, int]{
			E.Right[string](2),
			E.Right[string](4),
			E.Right[string](6),
		}
		assert.Equal(t, expected, result)
	})

	t.Run("filters out all Right values when none satisfy predicate", func(t *testing.T) {
		seq := iter.From(
			E.Right[string](1),
			E.Right[string](3),
			E.Right[string](5),
		)
		isEven := func(n int) bool { return n%2 == 0 }
		result := collectEithers(MonadFilter(seq, isEven))

		assert.Empty(t, result)
	})

	t.Run("passes through Left values unchanged", func(t *testing.T) {
		seq := iter.From(
			E.Right[string](1),
			E.Left[int]("error1"),
			E.Right[string](2),
			E.Left[int]("error2"),
			E.Right[string](3),
		)
		isEven := func(n int) bool { return n%2 == 0 }
		result := collectEithers(MonadFilter(seq, isEven))

		expected := []Either[string, int]{
			E.Left[int]("error1"),
			E.Right[string](2),
			E.Left[int]("error2"),
		}
		assert.Equal(t, expected, result)
	})
}

// TestMonadFilter_Empty tests MonadFilter with empty sequences
func TestMonadFilter_Empty(t *testing.T) {
	t.Run("returns empty from empty sequence", func(t *testing.T) {
		seq := iter.Empty[Either[string, int]]()
		isEven := func(n int) bool { return n%2 == 0 }
		result := collectEithers(MonadFilter(seq, isEven))

		assert.Empty(t, result)
	})
}

// TestMonadFilter_EdgeCases tests edge cases
func TestMonadFilter_EdgeCases(t *testing.T) {
	t.Run("handles zero values", func(t *testing.T) {
		seq := iter.From(
			E.Right[string](0),
			E.Right[string](1),
			E.Right[string](2),
		)
		isZero := func(n int) bool { return n == 0 }
		result := collectEithers(MonadFilter(seq, isZero))

		expected := []Either[string, int]{E.Right[string](0)}
		assert.Equal(t, expected, result)
	})

	t.Run("handles single Right element that passes", func(t *testing.T) {
		seq := iter.From(E.Right[string](2))
		isEven := func(n int) bool { return n%2 == 0 }
		result := collectEithers(MonadFilter(seq, isEven))

		expected := []Either[string, int]{E.Right[string](2)}
		assert.Equal(t, expected, result)
	})

	t.Run("handles single Right element that fails", func(t *testing.T) {
		seq := iter.From(E.Right[string](1))
		isEven := func(n int) bool { return n%2 == 0 }
		result := collectEithers(MonadFilter(seq, isEven))

		assert.Empty(t, result)
	})

	t.Run("handles single Left element", func(t *testing.T) {
		seq := iter.From(E.Left[int]("error"))
		isEven := func(n int) bool { return n%2 == 0 }
		result := collectEithers(MonadFilter(seq, isEven))

		expected := []Either[string, int]{E.Left[int]("error")}
		assert.Equal(t, expected, result)
	})

	t.Run("handles sequence of only Left values", func(t *testing.T) {
		seq := iter.From(
			E.Left[int]("error1"),
			E.Left[int]("error2"),
			E.Left[int]("error3"),
		)
		isEven := func(n int) bool { return n%2 == 0 }
		result := collectEithers(MonadFilter(seq, isEven))

		expected := []Either[string, int]{
			E.Left[int]("error1"),
			E.Left[int]("error2"),
			E.Left[int]("error3"),
		}
		assert.Equal(t, expected, result)
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
			E.Right[error](User{1, "Alice", 25}),
			E.Right[error](User{2, "Bob", 17}),
			E.Right[error](User{3, "Charlie", 30}),
			E.Left[User](errors.New("database error")),
			E.Right[error](User{4, "David", 16}),
		)
		isAdult := func(u User) bool { return u.Age >= 18 }
		result := collectEithers(MonadFilter(seq, isAdult))

		expected := []Either[error, User]{
			E.Right[error](User{1, "Alice", 25}),
			E.Right[error](User{3, "Charlie", 30}),
			E.Left[User](errors.New("database error")),
		}
		assert.Equal(t, expected, result)
	})

	t.Run("filters string types", func(t *testing.T) {
		seq := iter.From(
			E.Right[string]("hello"),
			E.Right[string](""),
			E.Right[string]("world"),
			E.Left[string]("error"),
			E.Right[string]("test"),
		)
		isNonEmpty := S.IsNonEmpty
		result := collectEithers(MonadFilter(seq, isNonEmpty))

		expected := []Either[string, string]{
			E.Right[string]("hello"),
			E.Right[string]("world"),
			E.Left[string]("error"),
			E.Right[string]("test"),
		}
		assert.Equal(t, expected, result)
	})
}

// TestFilter_Success tests the curried Filter function
func TestFilter_Success(t *testing.T) {
	t.Run("creates reusable filter function", func(t *testing.T) {
		isEven := func(n int) bool { return n%2 == 0 }
		evens := Filter[string](isEven)

		seq1 := iter.From(
			E.Right[string](1),
			E.Right[string](2),
			E.Right[string](3),
		)
		result1 := collectEithers(evens(seq1))
		expected1 := []Either[string, int]{E.Right[string](2)}
		assert.Equal(t, expected1, result1)

		seq2 := iter.From(
			E.Right[string](4),
			E.Left[int]("error"),
			E.Right[string](5),
			E.Right[string](6),
		)
		result2 := collectEithers(evens(seq2))
		expected2 := []Either[string, int]{
			E.Right[string](4),
			E.Left[int]("error"),
			E.Right[string](6),
		}
		assert.Equal(t, expected2, result2)
	})

	t.Run("works in pipeline", func(t *testing.T) {
		isPositive := N.MoreThan(0)
		result := F.Pipe2(
			iter.From(
				E.Right[string](-1),
				E.Right[string](2),
				E.Left[int]("error"),
				E.Right[string](3),
				E.Right[string](-4),
			),
			Filter[string](isPositive),
			collectEithers[string, int],
		)

		expected := []Either[string, int]{
			E.Right[string](2),
			E.Left[int]("error"),
			E.Right[string](3),
		}
		assert.Equal(t, expected, result)
	})

	t.Run("chains with Map", func(t *testing.T) {
		isEven := func(n int) bool { return n%2 == 0 }
		result := F.Pipe3(
			iter.From(1, 2, 3, 4, 5),
			FromSeq[string],
			Map[string](func(n int) int { return n * 2 }),
			Filter[string](isEven),
		)

		collected := collectEithers(result)
		expected := []Either[string, int]{
			E.Right[string](2),
			E.Right[string](4),
			E.Right[string](6),
			E.Right[string](8),
			E.Right[string](10),
		}
		assert.Equal(t, expected, collected)
	})
}

// TestFilter_VsFilterOrElse compares Filter and FilterOrElse behavior
func TestFilter_VsFilterOrElse(t *testing.T) {
	t.Run("Filter removes failing values", func(t *testing.T) {
		seq := iter.From(
			E.Right[string](1),
			E.Right[string](2),
			E.Right[string](3),
		)
		isEven := func(n int) bool { return n%2 == 0 }

		// Filter removes odd numbers
		filterResult := collectEithers(Filter[string](isEven)(seq))
		assert.Len(t, filterResult, 1)
		assert.Equal(t, E.Right[string](2), filterResult[0])
	})

	t.Run("FilterOrElse converts failing values to Left", func(t *testing.T) {
		seq := iter.From(
			E.Right[string](1),
			E.Right[string](2),
			E.Right[string](3),
		)
		isEven := func(n int) bool { return n%2 == 0 }
		onOdd := S.Format[int]("%d is odd")

		// FilterOrElse converts odd numbers to Left
		filterOrElseResult := collectEithers(FilterOrElse(isEven, onOdd)(seq))
		assert.Len(t, filterOrElseResult, 3)
		assert.True(t, E.IsLeft(filterOrElseResult[0]))
		assert.True(t, E.IsRight(filterOrElseResult[1]))
		assert.True(t, E.IsLeft(filterOrElseResult[2]))
	})
}

// BenchmarkMonadFilter benchmarks basic MonadFilter operation
func BenchmarkMonadFilter(b *testing.B) {
	seq := iter.From(
		E.Right[string](1),
		E.Right[string](2),
		E.Right[string](3),
		E.Right[string](4),
		E.Right[string](5),
	)
	isEven := func(n int) bool { return n%2 == 0 }

	b.ResetTimer()
	for range b.N {
		_ = collectEithers(MonadFilter(seq, isEven))
	}
}

// BenchmarkFilter benchmarks curried Filter operation
func BenchmarkFilter(b *testing.B) {
	seq := iter.From(
		E.Right[string](1),
		E.Right[string](2),
		E.Right[string](3),
		E.Right[string](4),
		E.Right[string](5),
	)
	isEven := func(n int) bool { return n%2 == 0 }
	evens := Filter[string](isEven)

	b.ResetTimer()
	for range b.N {
		_ = collectEithers(evens(seq))
	}
}

// BenchmarkMonadFilter_LargeSequence benchmarks with a large sequence
func BenchmarkMonadFilter_LargeSequence(b *testing.B) {
	makeSeq := func() SeqEither[string, int] {
		return func(yield func(Either[string, int]) bool) {
			for i := range 1000 {
				if !yield(E.Right[string](i)) {
					return
				}
			}
		}
	}
	isEven := func(n int) bool { return n%2 == 0 }

	b.ResetTimer()
	for range b.N {
		_ = collectEithers(MonadFilter(makeSeq(), isEven))
	}
}
