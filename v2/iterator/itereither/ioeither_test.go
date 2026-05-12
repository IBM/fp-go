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
	"slices"
	"strings"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	"github.com/IBM/fp-go/v2/iterator/iter"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

func sum(a, b int) int {
	return a + b
}

func collectEithers[E, A any](seq SeqEither[E, A]) []Either[E, A] {
	return slices.Collect(seq)
}

func TestLeft(t *testing.T) {
	result := collectEithers(Left[int]("error"))
	assert.Equal(t, []Either[string, int]{E.Left[int]("error")}, result)
}

func TestRight(t *testing.T) {
	result := collectEithers(Right[string](42))
	assert.Equal(t, []Either[string, int]{E.Right[string](42)}, result)
}

func TestOf(t *testing.T) {
	result := collectEithers(Of[string](42))
	assert.Equal(t, []Either[string, int]{E.Right[string](42)}, result)
}

func TestMonadOf(t *testing.T) {
	result := collectEithers(MonadOf[string](42))
	assert.Equal(t, []Either[string, int]{E.Right[string](42)}, result)
}

func TestFromEither(t *testing.T) {
	t.Run("from Right", func(t *testing.T) {
		result := collectEithers(FromEither(E.Right[string](42)))
		assert.Equal(t, []Either[string, int]{E.Right[string](42)}, result)
	})

	t.Run("from Left", func(t *testing.T) {
		result := collectEithers(FromEither(E.Left[int]("error")))
		assert.Equal(t, []Either[string, int]{E.Left[int]("error")}, result)
	})
}

func TestFromOption(t *testing.T) {
	onNone := F.Constant("none")

	t.Run("from Some", func(t *testing.T) {
		result := collectEithers(FromOption[int](onNone)(O.Some(42)))
		assert.Equal(t, []Either[string, int]{E.Right[string](42)}, result)
	})

	t.Run("from None", func(t *testing.T) {
		result := collectEithers(FromOption[int](onNone)(O.None[int]()))
		assert.Equal(t, []Either[string, int]{E.Left[int]("none")}, result)
	})
}

func TestFromSeq(t *testing.T) {
	seq := iter.From(1, 2, 3)
	result := collectEithers(FromSeq[string](seq))
	expected := []Either[string, int]{
		E.Right[string](1),
		E.Right[string](2),
		E.Right[string](3),
	}
	assert.Equal(t, expected, result)
}
func TestFromIO_Success(t *testing.T) {
	t.Run("converts IO computation to single-element Right sequence", func(t *testing.T) {
		io := func() int { return 42 }
		seq := FromIO[string](io)
		result := collectEithers(seq)
		assert.Equal(t, []Either[string, int]{E.Right[string](42)}, result)
	})

	t.Run("executes IO computation when sequence is consumed", func(t *testing.T) {
		executed := false
		io := func() string {
			executed = true
			return "hello"
		}
		seq := FromIO[error](io)

		// IO should not be executed yet
		assert.False(t, executed)

		// Consume the sequence
		result := collectEithers(seq)

		// Now IO should be executed
		assert.True(t, executed)
		assert.Equal(t, []Either[error, string]{E.Right[error]("hello")}, result)
	})

	t.Run("works with different types", func(t *testing.T) {
		t.Run("string", func(t *testing.T) {
			io := func() string { return "test" }
			seq := FromIO[error](io)
			result := collectEithers(seq)
			assert.Equal(t, []Either[error, string]{E.Right[error]("test")}, result)
		})

		t.Run("struct", func(t *testing.T) {
			type Person struct {
				Name string
				Age  int
			}
			io := func() Person { return Person{Name: "Alice", Age: 30} }
			seq := FromIO[string](io)
			result := collectEithers(seq)
			assert.Equal(t, []Either[string, Person]{E.Right[string](Person{Name: "Alice", Age: 30})}, result)
		})

		t.Run("pointer", func(t *testing.T) {
			value := 100
			io := func() *int { return &value }
			seq := FromIO[error](io)
			result := collectEithers(seq)
			assert.Len(t, result, 1)
			assert.True(t, E.IsRight(result[0]))
			ptr, _ := E.Unwrap(result[0])
			assert.Equal(t, 100, *ptr)
		})
	})

	t.Run("can be composed with other operations", func(t *testing.T) {
		io := func() int { return 10 }
		seq := F.Pipe1(
			FromIO[string](io),
			Map[string](func(x int) int { return x * 2 }),
		)
		result := collectEithers(seq)
		assert.Equal(t, []Either[string, int]{E.Right[string](20)}, result)
	})

	t.Run("can be used in chain operations", func(t *testing.T) {
		io := func() int { return 3 }
		seq := F.Pipe1(
			FromIO[string](io),
			Chain(func(n int) SeqEither[string, int] {
				return FromSeq[string](iter.From(n, n*2, n*3))
			}),
		)
		result := collectEithers(seq)
		expected := []Either[string, int]{
			E.Right[string](3),
			E.Right[string](6),
			E.Right[string](9),
		}
		assert.Equal(t, expected, result)
	})
}

func TestFromIO_EdgeCases(t *testing.T) {
	t.Run("handles zero value", func(t *testing.T) {
		io := func() int { return 0 }
		seq := FromIO[string](io)
		result := collectEithers(seq)
		assert.Equal(t, []Either[string, int]{E.Right[string](0)}, result)
	})

	t.Run("handles empty string", func(t *testing.T) {
		io := func() string { return "" }
		seq := FromIO[error](io)
		result := collectEithers(seq)
		assert.Equal(t, []Either[error, string]{E.Right[error]("")}, result)
	})

	t.Run("handles nil pointer", func(t *testing.T) {
		io := func() *int { return nil }
		seq := FromIO[error](io)
		result := collectEithers(seq)
		assert.Len(t, result, 1)
		assert.True(t, E.IsRight(result[0]))
		ptr, _ := E.Unwrap(result[0])
		assert.Nil(t, ptr)
	})
}

func TestFromIO_Integration(t *testing.T) {
	t.Run("multiple iterations execute IO multiple times", func(t *testing.T) {
		counter := 0
		io := func() int {
			counter++
			return counter
		}
		seq := FromIO[error](io)

		// First iteration
		result1 := collectEithers(seq)
		assert.Equal(t, []Either[error, int]{E.Right[error](1)}, result1)

		// Second iteration - IO executes again
		result2 := collectEithers(seq)
		assert.Equal(t, []Either[error, int]{E.Right[error](2)}, result2)
	})

	t.Run("works with MonadMap", func(t *testing.T) {
		io := func() int { return 5 }
		seq := MonadMap(FromIO[string](io), func(n int) string { return strings.Repeat("*", n) })
		result := collectEithers(seq)
		assert.Equal(t, []Either[string, string]{E.Right[string]("*****")}, result)
	})

	t.Run("works with MonadChain", func(t *testing.T) {
		io := func() int { return 2 }
		seq := MonadChain(
			FromIO[string](io),
			func(n int) SeqEither[string, int] {
				return FromSeq[string](iter.Replicate(n, n))
			},
		)
		result := collectEithers(seq)
		expected := []Either[string, int]{
			E.Right[string](2),
			E.Right[string](2),
		}
		assert.Equal(t, expected, result)
	})
}

func TestFromLazy_Success(t *testing.T) {
	t.Run("converts Lazy computation to single-element Right sequence", func(t *testing.T) {
		lazy := func() int { return 42 }
		seq := FromLazy[string](lazy)
		result := collectEithers(seq)
		assert.Equal(t, []Either[string, int]{E.Right[string](42)}, result)
	})

	t.Run("defers computation until sequence is consumed", func(t *testing.T) {
		executed := false
		lazy := func() string {
			executed = true
			return "lazy value"
		}
		seq := FromLazy[error](lazy)

		// Lazy computation should not be executed yet
		assert.False(t, executed)

		// Consume the sequence
		result := collectEithers(seq)

		// Now lazy computation should be executed
		assert.True(t, executed)
		assert.Equal(t, []Either[error, string]{E.Right[error]("lazy value")}, result)
	})

	t.Run("works with different types", func(t *testing.T) {
		t.Run("bool", func(t *testing.T) {
			lazy := func() bool { return true }
			seq := FromLazy[string](lazy)
			result := collectEithers(seq)
			assert.Equal(t, []Either[string, bool]{E.Right[string](true)}, result)
		})

		t.Run("slice", func(t *testing.T) {
			lazy := func() []int { return []int{1, 2, 3} }
			seq := FromLazy[error](lazy)
			result := collectEithers(seq)
			assert.Len(t, result, 1)
			assert.True(t, E.IsRight(result[0]))
			slice, _ := E.Unwrap(result[0])
			assert.Equal(t, []int{1, 2, 3}, slice)
		})

		t.Run("map", func(t *testing.T) {
			lazy := func() map[string]int {
				return map[string]int{"a": 1, "b": 2}
			}
			seq := FromLazy[error](lazy)
			result := collectEithers(seq)
			assert.Len(t, result, 1)
			assert.True(t, E.IsRight(result[0]))
			m, _ := E.Unwrap(result[0])
			assert.Equal(t, map[string]int{"a": 1, "b": 2}, m)
		})
	})

	t.Run("can be composed with other operations", func(t *testing.T) {
		lazy := func() int { return 5 }
		seq := F.Pipe1(
			FromLazy[string](lazy),
			Map[string](func(n int) int { return n * n }),
		)
		result := collectEithers(seq)
		assert.Equal(t, []Either[string, int]{E.Right[string](25)}, result)
	})

	t.Run("can be used with MonadChain", func(t *testing.T) {
		lazy := func() int { return 3 }
		seq := MonadChain(
			FromLazy[string](lazy),
			func(n int) SeqEither[string, string] {
				return FromSeq[string](iter.Replicate(n, strings.Repeat("x", n)))
			},
		)
		result := collectEithers(seq)
		expected := []Either[string, string]{
			E.Right[string]("xxx"),
			E.Right[string]("xxx"),
			E.Right[string]("xxx"),
		}
		assert.Equal(t, expected, result)
	})
}

func TestFromLazy_EdgeCases(t *testing.T) {
	t.Run("handles expensive computation", func(t *testing.T) {
		lazy := func() int {
			// Simulate expensive computation
			sum := 0
			for i := range 1000 {
				sum += i
			}
			return sum
		}
		seq := FromLazy[error](lazy)
		result := collectEithers(seq)
		assert.Equal(t, []Either[error, int]{E.Right[error](499500)}, result)
	})

	t.Run("handles function returning function", func(t *testing.T) {
		lazy := func() func(int) int {
			return func(x int) int { return x * 2 }
		}
		seq := FromLazy[error](lazy)
		result := collectEithers(seq)
		assert.Len(t, result, 1)
		assert.True(t, E.IsRight(result[0]))
		fn, _ := E.Unwrap(result[0])
		assert.Equal(t, 10, fn(5))
	})
}

func TestFromLazy_Integration(t *testing.T) {
	t.Run("multiple iterations execute lazy computation multiple times", func(t *testing.T) {
		counter := 0
		lazy := func() int {
			counter++
			return counter * 10
		}
		seq := FromLazy[error](lazy)

		// First iteration
		result1 := collectEithers(seq)
		assert.Equal(t, []Either[error, int]{E.Right[error](10)}, result1)

		// Second iteration - lazy computation executes again
		result2 := collectEithers(seq)
		assert.Equal(t, []Either[error, int]{E.Right[error](20)}, result2)
	})

	t.Run("works with MonadMapLeft", func(t *testing.T) {
		lazy := func() int { return 7 }
		// Even though we map left, the Right value passes through
		seq := MonadMapLeft(FromLazy[string](lazy), func(s string) int { return len(s) })
		result := collectEithers(seq)
		assert.Equal(t, []Either[int, int]{E.Right[int](7)}, result)
	})

	t.Run("works with Fold", func(t *testing.T) {
		lazy := func() int { return 42 }
		seq := FromLazy[string](lazy)
		folded := MonadFold(
			seq,
			func(s string) iter.Seq[string] { return iter.Of("error: " + s) },
			func(n int) iter.Seq[string] { return iter.Of(strings.Repeat("*", n)) },
		)
		result := slices.Collect(folded)
		assert.Equal(t, []string{strings.Repeat("*", 42)}, result)
	})
}

func TestFromIOEither_Success(t *testing.T) {
	t.Run("converts IOEither Right to single-element Right sequence", func(t *testing.T) {
		ioe := func() Either[string, int] { return E.Right[string](42) }
		seq := FromIOEither(ioe)
		result := collectEithers(seq)
		assert.Equal(t, []Either[string, int]{E.Right[string](42)}, result)
	})

	t.Run("converts IOEither Left to single-element Left sequence", func(t *testing.T) {
		ioe := func() Either[string, int] { return E.Left[int]("error") }
		seq := FromIOEither(ioe)
		result := collectEithers(seq)
		assert.Equal(t, []Either[string, int]{E.Left[int]("error")}, result)
	})

	t.Run("executes IOEither when sequence is consumed", func(t *testing.T) {
		executed := false
		ioe := func() Either[error, string] {
			executed = true
			return E.Right[error]("result")
		}
		seq := FromIOEither(ioe)

		// IOEither should not be executed yet
		assert.False(t, executed)

		// Consume the sequence
		result := collectEithers(seq)

		// Now IOEither should be executed
		assert.True(t, executed)
		assert.Equal(t, []Either[error, string]{E.Right[error]("result")}, result)
	})

	t.Run("multiple iterations execute IOEither multiple times", func(t *testing.T) {
		counter := 0
		ioe := func() Either[string, int] {
			counter++
			return E.Right[string](counter)
		}
		seq := FromIOEither(ioe)

		result1 := collectEithers(seq)
		assert.Equal(t, []Either[string, int]{E.Right[string](1)}, result1)

		result2 := collectEithers(seq)
		assert.Equal(t, []Either[string, int]{E.Right[string](2)}, result2)
	})

	t.Run("works with different error types", func(t *testing.T) {
		t.Run("string error", func(t *testing.T) {
			ioe := func() Either[string, int] { return E.Left[int]("validation failed") }
			seq := FromIOEither(ioe)
			result := collectEithers(seq)
			assert.Equal(t, []Either[string, int]{E.Left[int]("validation failed")}, result)
		})

		t.Run("error type", func(t *testing.T) {
			ioe := func() Either[error, int] {
				return E.Left[int](assert.AnError)
			}
			seq := FromIOEither(ioe)
			result := collectEithers(seq)
			assert.Len(t, result, 1)
			assert.True(t, E.IsLeft(result[0]))
			_, err := E.Unwrap(result[0])
			assert.Equal(t, assert.AnError, err)
		})

		t.Run("custom error struct", func(t *testing.T) {
			type ValidationError struct {
				Field   string
				Message string
			}
			ioe := func() Either[ValidationError, int] {
				return E.Left[int](ValidationError{Field: "age", Message: "must be positive"})
			}
			seq := FromIOEither(ioe)
			result := collectEithers(seq)
			assert.Len(t, result, 1)
			assert.True(t, E.IsLeft(result[0]))
			_, verr := E.Unwrap(result[0])
			assert.Equal(t, "age", verr.Field)
			assert.Equal(t, "must be positive", verr.Message)
		})
	})

	t.Run("can be composed with other operations", func(t *testing.T) {
		ioe := func() Either[string, int] { return E.Right[string](10) }
		seq := F.Pipe1(
			FromIOEither(ioe),
			Map[string](func(x int) int { return x * 2 }),
		)
		result := collectEithers(seq)
		assert.Equal(t, []Either[string, int]{E.Right[string](20)}, result)
	})

	t.Run("Left values stop chain operations", func(t *testing.T) {
		ioe := func() Either[string, int] { return E.Left[int]("error") }
		seq := F.Pipe1(
			FromIOEither(ioe),
			Map[string](func(x int) int { return x * 2 }),
		)
		result := collectEithers(seq)
		assert.Equal(t, []Either[string, int]{E.Left[int]("error")}, result)
	})
}

func TestFromIOEither_EdgeCases(t *testing.T) {
	t.Run("handles zero value in Right", func(t *testing.T) {
		ioe := func() Either[string, int] { return E.Right[string](0) }
		seq := FromIOEither(ioe)
		result := collectEithers(seq)
		assert.Equal(t, []Either[string, int]{E.Right[string](0)}, result)
	})

	t.Run("handles empty string in Left", func(t *testing.T) {
		ioe := func() Either[string, int] { return E.Left[int]("") }
		seq := FromIOEither(ioe)
		result := collectEithers(seq)
		assert.Equal(t, []Either[string, int]{E.Left[int]("")}, result)
	})

	t.Run("handles nil pointer in Right", func(t *testing.T) {
		ioe := func() Either[error, *int] {
			var nilPtr *int
			return E.Right[error](nilPtr)
		}
		seq := FromIOEither(ioe)
		result := collectEithers(seq)
		assert.Len(t, result, 1)
		assert.True(t, E.IsRight(result[0]))
		ptr, _ := E.Unwrap(result[0])
		assert.Nil(t, ptr)
	})
}

func TestFromIOEither_Integration(t *testing.T) {
	t.Run("multiple iterations execute IOEither multiple times", func(t *testing.T) {
		counter := 0
		ioe := func() Either[string, int] {
			counter++
			if counter%2 == 0 {
				return E.Left[int]("even")
			}
			return E.Right[string](counter)
		}
		seq := FromIOEither(ioe)

		// First iteration - odd
		result1 := collectEithers(seq)
		assert.Equal(t, []Either[string, int]{E.Right[string](1)}, result1)

		// Second iteration - even
		result2 := collectEithers(seq)
		assert.Equal(t, []Either[string, int]{E.Left[int]("even")}, result2)

		// Third iteration - odd
		result3 := collectEithers(seq)
		assert.Equal(t, []Either[string, int]{E.Right[string](3)}, result3)
	})

	t.Run("works with MonadChain", func(t *testing.T) {
		ioe := func() Either[string, int] { return E.Right[string](2) }
		seq := MonadChain(
			FromIOEither(ioe),
			func(n int) SeqEither[string, int] {
				return FromSeq[string](iter.From(n, n*2))
			},
		)
		result := collectEithers(seq)
		expected := []Either[string, int]{
			E.Right[string](2),
			E.Right[string](4),
		}
		assert.Equal(t, expected, result)
	})

	t.Run("works with GetOrElse", func(t *testing.T) {
		t.Run("Right value", func(t *testing.T) {
			ioe := func() Either[string, int] { return E.Right[string](42) }
			seq := FromIOEither(ioe)
			result := slices.Collect(GetOrElse(func(s string) iter.Seq[int] {
				return iter.Of(0)
			})(seq))
			assert.Equal(t, []int{42}, result)
		})

		t.Run("Left value", func(t *testing.T) {
			ioe := func() Either[string, int] { return E.Left[int]("error") }
			seq := FromIOEither(ioe)
			result := slices.Collect(GetOrElse(func(s string) iter.Seq[int] {
				return iter.Of(-1)
			})(seq))
			assert.Equal(t, []int{-1}, result)
		})
	})

	t.Run("works with OrElse for error recovery", func(t *testing.T) {
		ioe := func() Either[string, int] { return E.Left[int]("not found") }
		recover := OrElse(func(err string) SeqEither[string, int] {
			if err == "not found" {
				return Right[string](0) // default value
			}
			return Left[int](err)
		})
		seq := recover(FromIOEither(ioe))
		result := collectEithers(seq)
		assert.Equal(t, []Either[string, int]{E.Right[string](0)}, result)
	})
}

func TestFromIO_vs_FromLazy_vs_FromIOEither(t *testing.T) {
	t.Run("FromLazy delegates to FromIO", func(t *testing.T) {
		value := 42
		io := func() int { return value }
		lazy := func() int { return value }

		seqIO := FromIO[string](io)
		seqLazy := FromLazy[string](lazy)

		resultIO := collectEithers(seqIO)
		resultLazy := collectEithers(seqLazy)

		assert.Equal(t, resultIO, resultLazy)
	})

	t.Run("FromIO and FromLazy always produce Right", func(t *testing.T) {
		io := func() int { return 42 }
		lazy := func() int { return 42 }

		seqIO := FromIO[string](io)
		seqLazy := FromLazy[string](lazy)

		resultIO := collectEithers(seqIO)
		resultLazy := collectEithers(seqLazy)

		// Both should always be Right
		assert.True(t, E.IsRight(resultIO[0]))
		assert.True(t, E.IsRight(resultLazy[0]))
	})

	t.Run("FromIOEither can produce Left or Right", func(t *testing.T) {
		ioeRight := func() Either[string, int] { return E.Right[string](42) }
		ioeLeft := func() Either[string, int] { return E.Left[int]("error") }

		seqRight := FromIOEither(ioeRight)
		seqLeft := FromIOEither(ioeLeft)

		resultRight := collectEithers(seqRight)
		resultLeft := collectEithers(seqLeft)

		assert.True(t, E.IsRight(resultRight[0]))
		assert.True(t, E.IsLeft(resultLeft[0]))
	})

	t.Run("all handle side effects similarly", func(t *testing.T) {
		counterIO := 0
		counterLazy := 0
		counterIOE := 0

		io := func() int {
			counterIO++
			return counterIO
		}
		lazy := func() int {
			counterLazy++
			return counterLazy
		}
		ioe := func() Either[string, int] {
			counterIOE++
			return E.Right[string](counterIOE)
		}

		seqIO := FromIO[string](io)
		seqLazy := FromLazy[string](lazy)
		seqIOE := FromIOEither(ioe)

		// All should execute on consumption
		resultIO := collectEithers(seqIO)
		resultLazy := collectEithers(seqLazy)
		resultIOE := collectEithers(seqIOE)

		assert.Equal(t, []Either[string, int]{E.Right[string](1)}, resultIO)
		assert.Equal(t, []Either[string, int]{E.Right[string](1)}, resultLazy)
		assert.Equal(t, []Either[string, int]{E.Right[string](1)}, resultIOE)
		assert.Equal(t, 1, counterIO)
		assert.Equal(t, 1, counterLazy)
		assert.Equal(t, 1, counterIOE)
	})
}

func TestMonadMap(t *testing.T) {
	t.Run("maps Right values", func(t *testing.T) {
		seq := iter.From(E.Right[string](1), E.Right[string](2))
		result := collectEithers(MonadMap(seq, utils.Double))
		expected := []Either[string, int]{
			E.Right[string](2),
			E.Right[string](4),
		}
		assert.Equal(t, expected, result)
	})

	t.Run("preserves Left values", func(t *testing.T) {
		seq := iter.From(E.Left[int]("error"), E.Right[string](2))
		result := collectEithers(MonadMap(seq, utils.Double))
		expected := []Either[string, int]{
			E.Left[int]("error"),
			E.Right[string](4),
		}
		assert.Equal(t, expected, result)
	})
}

func TestMap(t *testing.T) {
	seq := iter.From(E.Right[string](1), E.Right[string](2))
	result := F.Pipe1(seq, Map[string](utils.Double))
	expected := []Either[string, int]{
		E.Right[string](2),
		E.Right[string](4),
	}
	assert.Equal(t, expected, collectEithers(result))
}

func TestMonadMapTo(t *testing.T) {
	seq := iter.From(E.Right[string](1), E.Right[string](2))
	result := collectEithers(MonadMapTo(seq, 99))
	expected := []Either[string, int]{
		E.Right[string](99),
		E.Right[string](99),
	}
	assert.Equal(t, expected, result)
}

func TestMapTo(t *testing.T) {
	seq := iter.From(E.Right[string](1), E.Right[string](2))
	result := F.Pipe1(seq, MapTo[string, int, int](99))
	expected := []Either[string, int]{
		E.Right[string](99),
		E.Right[string](99),
	}
	assert.Equal(t, expected, collectEithers(result))
}

func TestMonadChain(t *testing.T) {
	t.Run("chains Right values", func(t *testing.T) {
		seq := iter.From(E.Right[string](1), E.Right[string](2))
		f := func(n int) SeqEither[string, int] {
			return iter.From(E.Right[string](n*10), E.Right[string](n*100))
		}
		result := collectEithers(MonadChain(seq, f))
		expected := []E.Either[string, int]{
			E.Right[string](10),
			E.Right[string](100),
			E.Right[string](20),
			E.Right[string](200),
		}
		assert.Equal(t, expected, result)
	})

	t.Run("stops on Left", func(t *testing.T) {
		seq := iter.From(E.Right[string](1), E.Left[int]("error"))
		f := func(n int) SeqEither[string, int] {
			return iter.From(E.Right[string](n * 10))
		}
		result := collectEithers(MonadChain(seq, f))
		expected := []E.Either[string, int]{
			E.Right[string](10),
			E.Left[int]("error"),
		}
		assert.Equal(t, expected, result)
	})
}

func TestChain(t *testing.T) {
	seq := iter.From(E.Right[string](1), E.Right[string](2))
	f := func(n int) SeqEither[string, int] {
		return iter.From(E.Right[string](n * 10))
	}
	result := F.Pipe1(seq, Chain(f))
	expected := []E.Either[string, int]{
		E.Right[string](10),
		E.Right[string](20),
	}
	assert.Equal(t, expected, collectEithers(result))
}

func TestChainEitherK(t *testing.T) {
	f := ChainEitherK(func(n int) E.Either[string, int] {
		if n > 0 {
			return E.Right[string](n * 2)
		}
		return E.Left[int]("negative")
	})

	t.Run("chains successful Either", func(t *testing.T) {
		seq := iter.From(E.Right[string](5))
		result := collectEithers(f(seq))
		assert.Equal(t, []E.Either[string, int]{E.Right[string](10)}, result)
	})

	t.Run("chains failing Either", func(t *testing.T) {
		seq := iter.From(E.Right[string](-5))
		result := collectEithers(f(seq))
		assert.Equal(t, []E.Either[string, int]{E.Left[int]("negative")}, result)
	})

	t.Run("preserves Left", func(t *testing.T) {
		seq := iter.From(E.Left[int]("original"))
		result := collectEithers(f(seq))
		assert.Equal(t, []E.Either[string, int]{E.Left[int]("original")}, result)
	})
}

func TestChainOptionK(t *testing.T) {
	f := ChainOptionK[int, int](F.Constant("none"))(func(n int) O.Option[int] {
		if n > 0 {
			return O.Some(n * 2)
		}
		return O.None[int]()
	})

	t.Run("chains Some", func(t *testing.T) {
		seq := iter.From(E.Right[string](5))
		result := collectEithers(f(seq))
		assert.Equal(t, []E.Either[string, int]{E.Right[string](10)}, result)
	})

	t.Run("chains None", func(t *testing.T) {
		seq := iter.From(E.Right[string](-5))
		result := collectEithers(f(seq))
		assert.Equal(t, []E.Either[string, int]{E.Left[int]("none")}, result)
	})
}

func TestChainSeqK(t *testing.T) {
	f := ChainSeqK[string](func(n int) iter.Seq[int] {
		return iter.From(n*10, n*100)
	})

	seq := iter.From(E.Right[string](1), E.Right[string](2))
	result := collectEithers(f(seq))
	expected := []E.Either[string, int]{
		E.Right[string](10),
		E.Right[string](100),
		E.Right[string](20),
		E.Right[string](200),
	}
	assert.Equal(t, expected, result)
}

func TestFlatten(t *testing.T) {
	inner := iter.From(E.Right[string](1), E.Right[string](2))
	outer := iter.From(E.Right[string](inner))
	result := collectEithers(Flatten(outer))
	expected := []E.Either[string, int]{
		E.Right[string](1),
		E.Right[string](2),
	}
	assert.Equal(t, expected, result)
}

func TestMonadMapLeft(t *testing.T) {
	seq := iter.From(E.Left[int]("error"), E.Right[string](42))
	result := collectEithers(MonadMapLeft(seq, func(s string) int { return len(s) }))
	expected := []E.Either[int, int]{
		E.Left[int](5),
		E.Right[int](42),
	}
	assert.Equal(t, expected, result)
}

func TestMapLeft(t *testing.T) {
	seq := iter.From(E.Left[int]("error"), E.Right[string](42))
	result := F.Pipe1(seq, MapLeft[int](func(s string) int { return len(s) }))
	expected := []E.Either[int, int]{
		E.Left[int](5),
		E.Right[int](42),
	}
	assert.Equal(t, expected, collectEithers(result))
}

func TestMonadBiMap(t *testing.T) {
	seq := iter.From(E.Left[int]("error"), E.Right[string](42))
	result := collectEithers(MonadBiMap(
		seq,
		func(s string) int { return len(s) },
		utils.Double,
	))
	expected := []E.Either[int, int]{
		E.Left[int](5),
		E.Right[int](84),
	}
	assert.Equal(t, expected, result)
}

func TestBiMap(t *testing.T) {
	seq := iter.From(E.Left[int]("error"), E.Right[string](42))
	result := F.Pipe1(
		seq,
		BiMap(func(s string) int { return len(s) }, utils.Double),
	)
	expected := []E.Either[int, int]{
		E.Left[int](5),
		E.Right[int](84),
	}
	assert.Equal(t, expected, collectEithers(result))
}

func TestSwap(t *testing.T) {
	seq := iter.From(E.Left[int]("error"), E.Right[string](42))
	result := collectEithers(Swap(seq))
	expected := []E.Either[int, string]{
		E.Right[int]("error"),
		E.Left[string](42),
	}
	assert.Equal(t, expected, result)
}

func TestMonadAlt(t *testing.T) {
	t.Run("Right stays Right", func(t *testing.T) {
		first := iter.From(E.Right[string](1))
		second := func() SeqEither[string, int] {
			return iter.From(E.Right[string](2))
		}
		result := collectEithers(MonadAlt(first, second))
		assert.Equal(t, []E.Either[string, int]{E.Right[string](1)}, result)
	})

	t.Run("Left uses alternative", func(t *testing.T) {
		first := iter.From(E.Left[int]("error"))
		second := func() SeqEither[string, int] {
			return iter.From(E.Right[string](2))
		}
		result := collectEithers(MonadAlt(first, second))
		assert.Equal(t, []E.Either[string, int]{E.Right[string](2)}, result)
	})
}

func TestAlt(t *testing.T) {
	first := iter.From(E.Left[int]("error"))
	second := func() SeqEither[string, int] {
		return iter.From(E.Right[string](2))
	}
	result := F.Pipe1(first, Alt(second))
	assert.Equal(t, []E.Either[string, int]{E.Right[string](2)}, collectEithers(result))
}

func TestMonadChainLeft(t *testing.T) {
	t.Run("transforms Left", func(t *testing.T) {
		seq := iter.From(E.Left[int]("error"))
		f := func(s string) SeqEither[int, int] {
			return iter.From(E.Left[int](len(s)))
		}
		result := collectEithers(MonadChainLeft(seq, f))
		assert.Equal(t, []E.Either[int, int]{E.Left[int](5)}, result)
	})

	t.Run("preserves Right", func(t *testing.T) {
		seq := iter.From(E.Right[string](42))
		f := func(s string) SeqEither[int, int] {
			return iter.From(E.Left[int](len(s)))
		}
		result := collectEithers(MonadChainLeft(seq, f))
		assert.Equal(t, []E.Either[int, int]{E.Right[int](42)}, result)
	})
}

func TestChainLeft(t *testing.T) {
	seq := iter.From(E.Left[int]("error"))
	f := func(s string) SeqEither[int, int] {
		return iter.From(E.Left[int](len(s)))
	}
	result := F.Pipe1(seq, ChainLeft(f))
	assert.Equal(t, []E.Either[int, int]{E.Left[int](5)}, collectEithers(result))
}

func TestOrElse(t *testing.T) {
	t.Run("recovers from Left", func(t *testing.T) {
		onLeft := func(s string) SeqEither[string, int] {
			if s == "recoverable" {
				return iter.From(E.Right[string](0))
			}
			return iter.From(E.Left[int](s))
		}
		seq := iter.From(E.Left[int]("recoverable"))
		result := collectEithers(OrElse(onLeft)(seq))
		assert.Equal(t, []Either[string, int]{E.Right[string](0)}, result)
	})

	t.Run("preserves Right", func(t *testing.T) {
		onLeft := func(s string) SeqEither[string, int] {
			return iter.From(E.Left[int]("fallback"))
		}
		seq := iter.From(E.Right[string](42))
		result := collectEithers(OrElse(onLeft)(seq))
		assert.Equal(t, []Either[string, int]{E.Right[string](42)}, result)
	})
}

func TestMonadChainFirstLeft(t *testing.T) {
	t.Run("executes side effect but preserves Left", func(t *testing.T) {
		var sideEffect string
		seq := iter.From(E.Left[int]("error"))
		f := func(s string) SeqEither[string, int] {
			sideEffect = "logged: " + s
			return iter.From(E.Right[string](999))
		}
		result := collectEithers(MonadChainFirstLeft(seq, f))
		assert.Equal(t, []E.Either[string, int]{E.Left[int]("error")}, result)
		assert.Equal(t, "logged: error", sideEffect)
	})

	t.Run("preserves Right without calling function", func(t *testing.T) {
		var called bool
		seq := iter.From(E.Right[string](42))
		f := func(s string) SeqEither[string, int] {
			called = true
			return iter.From(E.Right[string](999))
		}
		result := collectEithers(MonadChainFirstLeft(seq, f))
		assert.Equal(t, []E.Either[string, int]{E.Right[string](42)}, result)
		assert.False(t, called)
	})
}

func TestChainFirstLeft(t *testing.T) {
	var sideEffect string
	seq := iter.From(E.Left[int]("error"))
	f := func(s string) SeqEither[string, int] {
		sideEffect = "logged: " + s
		return iter.From(E.Right[string](999))
	}
	result := F.Pipe1(seq, ChainFirstLeft[int](f))
	assert.Equal(t, []E.Either[string, int]{E.Left[int]("error")}, collectEithers(result))
	assert.Equal(t, "logged: error", sideEffect)
}

func TestMonadFlap(t *testing.T) {
	fab := iter.From(E.Right[string](utils.Double))
	result := collectEithers(MonadFlap(fab, 21))
	assert.Equal(t, []E.Either[string, int]{E.Right[string](42)}, result)
}

func TestFlap(t *testing.T) {
	fab := iter.From(E.Right[string](utils.Double))
	result := F.Pipe1(fab, Flap[string, int](21))
	assert.Equal(t, []E.Either[string, int]{E.Right[string](42)}, collectEithers(result))
}

func TestMonadChainTo(t *testing.T) {
	first := iter.From(E.Right[string](1))
	second := iter.From(E.Right[string](2))
	result := collectEithers(MonadChainTo(first, second))
	assert.Equal(t, []E.Either[string, int]{E.Right[string](2)}, result)
}

func TestChainTo(t *testing.T) {
	first := iter.From(E.Right[string](1))
	second := iter.From(E.Right[string](2))
	result := F.Pipe1(first, ChainTo[int](second))
	assert.Equal(t, []E.Either[string, int]{E.Right[string](2)}, collectEithers(result))
}

func TestMonadChainFirst(t *testing.T) {
	t.Run("executes side effect and returns original", func(t *testing.T) {
		var sideEffect int
		seq := iter.From(E.Right[string](42))
		f := func(n int) SeqEither[string, string] {
			sideEffect = n * 2
			return iter.From(E.Right[string]("ignored"))
		}
		result := collectEithers(MonadChainFirst(seq, f))
		assert.Equal(t, []E.Either[string, int]{E.Right[string](42)}, result)
		assert.Equal(t, 84, sideEffect)
	})
}

func TestChainFirst(t *testing.T) {
	var sideEffect int
	seq := iter.From(E.Right[string](42))
	f := func(n int) SeqEither[string, string] {
		sideEffect = n * 2
		return iter.From(E.Right[string]("ignored"))
	}
	result := F.Pipe1(seq, ChainFirst[string](f))
	assert.Equal(t, []E.Either[string, int]{E.Right[string](42)}, collectEithers(result))
	assert.Equal(t, 84, sideEffect)
}

func TestMonadTap(t *testing.T) {
	var sideEffect int
	seq := iter.From(E.Right[string](42))
	f := func(n int) SeqEither[string, string] {
		sideEffect = n * 2
		return iter.From(E.Right[string]("ignored"))
	}
	result := collectEithers(MonadTap(seq, f))
	assert.Equal(t, []E.Either[string, int]{E.Right[string](42)}, result)
	assert.Equal(t, 84, sideEffect)
}

func TestTap(t *testing.T) {
	var sideEffect int
	seq := iter.From(E.Right[string](42))
	f := func(n int) SeqEither[string, string] {
		sideEffect = n * 2
		return iter.From(E.Right[string]("ignored"))
	}
	result := F.Pipe1(seq, Tap[string](f))
	assert.Equal(t, []E.Either[string, int]{E.Right[string](42)}, collectEithers(result))
	assert.Equal(t, 84, sideEffect)
}

func TestGetOrElse(t *testing.T) {
	t.Run("extracts Right value", func(t *testing.T) {
		seq := iter.From(E.Right[string](42))
		onLeft := func(s string) iter.Seq[int] {
			return iter.From(0)
		}
		result := slices.Collect(GetOrElse(onLeft)(seq))
		assert.Equal(t, []int{42}, result)
	})

	t.Run("uses default for Left", func(t *testing.T) {
		seq := iter.From(E.Left[int]("error"))
		onLeft := func(s string) iter.Seq[int] {
			return iter.From(0)
		}
		result := slices.Collect(GetOrElse(onLeft)(seq))
		assert.Equal(t, []int{0}, result)
	})
}

func TestGetOrElseOf(t *testing.T) {
	t.Run("extracts Right value", func(t *testing.T) {
		seq := iter.From(E.Right[string](42))
		onLeft := func(s string) int { return 0 }
		result := slices.Collect(GetOrElseOf(onLeft)(seq))
		assert.Equal(t, []int{42}, result)
	})

	t.Run("uses default for Left", func(t *testing.T) {
		seq := iter.From(E.Left[int]("error"))
		onLeft := func(s string) int { return 0 }
		result := slices.Collect(GetOrElseOf(onLeft)(seq))
		assert.Equal(t, []int{0}, result)
	})
}

func TestWithResource(t *testing.T) {
	var released bool
	onCreate := iter.From(E.Right[error]("resource"))
	onRelease := func(r string) SeqEither[error, F.Void] {
		released = true
		return iter.From(E.Right[error](F.VOID))
	}
	use := func(r string) SeqEither[error, int] {
		return iter.From(E.Right[error](len(r)))
	}

	withRes := WithResource[int, error, string, F.Void](onCreate, onRelease)
	result := collectEithers(withRes(use))
	assert.Equal(t, []Either[error, int]{E.Right[error](8)}, result)
	assert.True(t, released)
}

func TestMonadChainFirstEitherK(t *testing.T) {
	var sideEffect int
	seq := iter.From(E.Right[string](42))
	f := func(n int) E.Either[string, string] {
		sideEffect = n * 2
		return E.Right[string]("ignored")
	}
	result := collectEithers(MonadChainFirstEitherK(seq, f))
	assert.Equal(t, []E.Either[string, int]{E.Right[string](42)}, result)
	assert.Equal(t, 84, sideEffect)
}

func TestChainFirstEitherK(t *testing.T) {
	var sideEffect int
	seq := iter.From(E.Right[string](42))
	f := func(n int) E.Either[string, string] {
		sideEffect = n * 2
		return E.Right[string]("ignored")
	}
	result := F.Pipe1(seq, ChainFirstEitherK[int](f))
	assert.Equal(t, []E.Either[string, int]{E.Right[string](42)}, collectEithers(result))
	assert.Equal(t, 84, sideEffect)
}

func TestMonadTapEitherK(t *testing.T) {
	var sideEffect int
	seq := iter.From(E.Right[string](42))
	f := func(n int) E.Either[string, string] {
		sideEffect = n * 2
		return E.Right[string]("ignored")
	}
	result := collectEithers(MonadTapEitherK(seq, f))
	assert.Equal(t, []E.Either[string, int]{E.Right[string](42)}, result)
	assert.Equal(t, 84, sideEffect)
}

func TestTapEitherK(t *testing.T) {
	var sideEffect int
	seq := iter.From(E.Right[string](42))
	f := func(n int) E.Either[string, string] {
		sideEffect = n * 2
		return E.Right[string]("ignored")
	}
	result := F.Pipe1(seq, TapEitherK[int](f))
	assert.Equal(t, []E.Either[string, int]{E.Right[string](42)}, collectEithers(result))
	assert.Equal(t, 84, sideEffect)
}

func TestMonadFold(t *testing.T) {
	t.Run("folds Right", func(t *testing.T) {
		seq := iter.From(E.Right[string](42))
		onLeft := func(s string) iter.Seq[int] { return iter.From(-1) }
		onRight := func(n int) iter.Seq[int] { return iter.From(n * 2) }
		result := slices.Collect(MonadFold(seq, onLeft, onRight))
		assert.Equal(t, []int{84}, result)
	})

	t.Run("folds Left", func(t *testing.T) {
		seq := iter.From(E.Left[int]("error"))
		onLeft := func(s string) iter.Seq[int] { return iter.From(-1) }
		onRight := func(n int) iter.Seq[int] { return iter.From(n * 2) }
		result := slices.Collect(MonadFold(seq, onLeft, onRight))
		assert.Equal(t, []int{-1}, result)
	})
}

func TestFold(t *testing.T) {
	seq := iter.From(E.Right[string](42), E.Left[int]("error"))
	onLeft := func(s string) iter.Seq[int] { return iter.From(-1) }
	onRight := func(n int) iter.Seq[int] { return iter.From(n * 2) }
	result := slices.Collect(Fold(onLeft, onRight)(seq))
	assert.Equal(t, []int{84, -1}, result)
}

func TestLeftSeq(t *testing.T) {
	seq := iter.From("error1", "error2")
	result := collectEithers(LeftSeq[int](seq))
	expected := []E.Either[string, int]{
		E.Left[int]("error1"),
		E.Left[int]("error2"),
	}
	assert.Equal(t, expected, result)
}

func TestRightSeq(t *testing.T) {
	seq := iter.From(1, 2, 3)
	result := collectEithers(RightSeq[string](seq))
	expected := []E.Either[string, int]{
		E.Right[string](1),
		E.Right[string](2),
		E.Right[string](3),
	}
	assert.Equal(t, expected, result)
}

func TestMonadMergeMap(t *testing.T) {
	seq := iter.From(E.Right[string](1), E.Right[string](2))
	f := func(n int) SeqEither[string, int] {
		return iter.From(E.Right[string](n*10), E.Right[string](n*100))
	}
	result := collectEithers(MonadMergeMap(seq, f))
	// MergeMap should interleave results
	assert.Len(t, result, 4)
	assert.Contains(t, result, E.Right[string](10))
	assert.Contains(t, result, E.Right[string](100))
	assert.Contains(t, result, E.Right[string](20))
	assert.Contains(t, result, E.Right[string](200))
}

func TestMergeMap(t *testing.T) {
	seq := iter.From(E.Right[string](1), E.Right[string](2))
	f := func(n int) SeqEither[string, int] {
		return iter.From(E.Right[string](n * 10))
	}
	result := F.Pipe1(seq, MergeMap(f))
	assert.Len(t, collectEithers(result), 2)
}

func TestMonadChainSeqK(t *testing.T) {
	seq := iter.From(E.Right[string](1), E.Right[string](2))
	f := func(n int) iter.Seq[int] {
		return iter.From(n*10, n*100)
	}
	result := collectEithers(MonadChainSeqK(seq, f))
	expected := []E.Either[string, int]{
		E.Right[string](10),
		E.Right[string](100),
		E.Right[string](20),
		E.Right[string](200),
	}
	assert.Equal(t, expected, result)
}

func TestMonadMergeMapSeqK(t *testing.T) {
	seq := iter.From(E.Right[string](1), E.Right[string](2))
	f := func(n int) iter.Seq[int] {
		return iter.From(n * 10)
	}
	result := collectEithers(MonadMergeMapSeqK(seq, f))
	assert.Len(t, result, 2)
}

func TestMergeMapSeqK(t *testing.T) {
	seq := iter.From(E.Right[string](1), E.Right[string](2))
	f := func(n int) iter.Seq[int] {
		return iter.From(n * 10)
	}
	result := F.Pipe1(seq, MergeMapSeqK[string](f))
	assert.Len(t, collectEithers(result), 2)
}

func TestMonadChainToSeq(t *testing.T) {
	first := iter.From(E.Right[string](1))
	second := iter.From(2, 3)
	result := collectEithers(MonadChainToSeq(first, second))
	expected := []E.Either[string, int]{
		E.Right[string](2),
		E.Right[string](3),
	}
	assert.Equal(t, expected, result)
}

func TestChainToSeq(t *testing.T) {
	first := iter.From(E.Right[string](1))
	second := iter.From(2, 3)
	result := F.Pipe1(first, ChainToSeq[string, int](second))
	expected := []E.Either[string, int]{
		E.Right[string](2),
		E.Right[string](3),
	}
	assert.Equal(t, expected, collectEithers(result))
}

func TestMonadTapLeft(t *testing.T) {
	var sideEffect string
	seq := iter.From(E.Left[int]("error"))
	f := func(s string) SeqEither[string, int] {
		sideEffect = "logged: " + s
		return iter.From(E.Right[string](999))
	}
	result := collectEithers(MonadTapLeft(seq, f))
	assert.Equal(t, []E.Either[string, int]{E.Left[int]("error")}, result)
	assert.Equal(t, "logged: error", sideEffect)
}

func TestTapLeft(t *testing.T) {
	var sideEffect string
	seq := iter.From(E.Left[int]("error"))
	f := func(s string) SeqEither[string, int] {
		sideEffect = "logged: " + s
		return iter.From(E.Right[string](999))
	}
	result := F.Pipe1(seq, TapLeft[int](f))
	assert.Equal(t, []E.Either[string, int]{E.Left[int]("error")}, collectEithers(result))
	assert.Equal(t, "logged: error", sideEffect)
}

func TestMonadChainEitherK(t *testing.T) {
	seq := iter.From(E.Right[string](5))
	f := func(n int) E.Either[string, int] {
		if n > 0 {
			return E.Right[string](n * 2)
		}
		return E.Left[int]("negative")
	}
	result := collectEithers(MonadChainEitherK(seq, f))
	assert.Equal(t, []E.Either[string, int]{E.Right[string](10)}, result)
}

func TestEdgeCases(t *testing.T) {
	t.Run("empty sequence", func(t *testing.T) {
		seq := iter.From[E.Either[string, int]]()
		result := collectEithers(seq)
		assert.Empty(t, result)
	})

	t.Run("mixed Left and Right", func(t *testing.T) {
		seq := iter.From(
			E.Right[string](1),
			E.Left[int]("error1"),
			E.Right[string](2),
			E.Left[int]("error2"),
		)
		result := collectEithers(MonadMap(seq, utils.Double))
		expected := []E.Either[string, int]{
			E.Right[string](2),
			E.Left[int]("error1"),
			E.Right[string](4),
			E.Left[int]("error2"),
		}
		assert.Equal(t, expected, result)
	})

	t.Run("error propagation in chain", func(t *testing.T) {
		seq := iter.From(E.Right[string](1), E.Left[int]("error"), E.Right[string](2))
		f := func(n int) SeqEither[string, int] {
			return iter.From(E.Right[string](n * 10))
		}
		result := collectEithers(MonadChain(seq, f))
		expected := []E.Either[string, int]{
			E.Right[string](10),
			E.Left[int]("error"),
			E.Right[string](20),
		}
		assert.Equal(t, expected, result)
	})
}

func TestComplexPipeline(t *testing.T) {
	// Test a complex pipeline combining multiple operations
	result := F.Pipe3(
		iter.From(1, 2, 3),
		FromSeq[string],
		Map[string](utils.Double),
		Chain(func(n int) SeqEither[string, int] {
			if n > 5 {
				return iter.From(E.Left[int]("too large"))
			}
			return iter.From(E.Right[string](n * 10))
		}),
	)

	collected := collectEithers(result)
	expected := []E.Either[string, int]{
		E.Right[string](20),
		E.Right[string](40),
		E.Left[int]("too large"),
	}
	assert.Equal(t, expected, collected)
}

func TestErrorRecovery(t *testing.T) {
	// Test error recovery with OrElse
	result := F.Pipe2(
		iter.From(E.Left[int]("recoverable"), E.Right[string](42)),
		OrElse(func(s string) SeqEither[string, int] {
			if s == "recoverable" {
				return iter.From(E.Right[string](0))
			}
			return iter.From(E.Left[int](s))
		}),
		Map[string](utils.Double),
	)

	collected := collectEithers(result)
	expected := []E.Either[string, int]{
		E.Right[string](0),
		E.Right[string](84),
	}
	assert.Equal(t, expected, collected)
}
func TestMonadReduce_Success(t *testing.T) {
	t.Run("reduces all Right values", func(t *testing.T) {
		seq := iter.From(E.Right[string](1), E.Right[string](2), E.Right[string](3), E.Right[string](4), E.Right[string](5))
		result := MonadReduce(seq, sum, 0)
		assert.Equal(t, E.Right[string](15), result())
	})

	t.Run("reduces with multiplication", func(t *testing.T) {
		seq := iter.From(E.Right[string](2), E.Right[string](3), E.Right[string](4))
		result := MonadReduce(seq, func(acc, x int) int { return acc * x }, 1)
		assert.Equal(t, E.Right[string](24), result())
	})

	t.Run("reduces with string concatenation", func(t *testing.T) {
		seq := iter.From(E.Right[string]("a"), E.Right[string]("b"), E.Right[string]("c"))
		result := MonadReduce(seq, func(acc, x string) string { return acc + x }, "")
		assert.Equal(t, E.Right[string]("abc"), result())
	})

	t.Run("reduces to different type", func(t *testing.T) {
		seq := iter.From(E.Right[string](1), E.Right[string](2), E.Right[string](3))
		result := MonadReduce(seq, func(acc []int, x int) []int {
			return append(acc, x)
		}, []int{})
		assert.Equal(t, E.Right[string]([]int{1, 2, 3}), result())
	})

	t.Run("empty sequence returns initial value", func(t *testing.T) {
		seq := iter.From[E.Either[string, int]]()
		result := MonadReduce(seq, sum, 42)
		assert.Equal(t, E.Right[string](42), result())
	})

	t.Run("single Right value", func(t *testing.T) {
		seq := iter.From(E.Right[string](10))
		result := MonadReduce(seq, sum, 5)
		assert.Equal(t, E.Right[string](15), result())
	})
}

func TestMonadReduce_Failure(t *testing.T) {
	t.Run("stops at first Left", func(t *testing.T) {
		seq := iter.From(E.Right[string](1), E.Right[string](2), E.Left[int]("error"), E.Right[string](4))
		result := MonadReduce(seq, sum, 0)
		assert.Equal(t, E.Left[int]("error"), result())
	})

	t.Run("Left at beginning", func(t *testing.T) {
		seq := iter.From(E.Left[int]("first error"), E.Right[string](1), E.Right[string](2))
		result := MonadReduce(seq, sum, 0)
		assert.Equal(t, E.Left[int]("first error"), result())
	})

	t.Run("Left at end", func(t *testing.T) {
		seq := iter.From(E.Right[string](1), E.Right[string](2), E.Left[int]("last error"))
		result := MonadReduce(seq, sum, 0)
		assert.Equal(t, E.Left[int]("last error"), result())
	})

	t.Run("only Left values", func(t *testing.T) {
		seq := iter.From(E.Left[int]("error1"), E.Left[int]("error2"))
		result := MonadReduce(seq, sum, 0)
		assert.Equal(t, E.Left[int]("error1"), result())
	})

	t.Run("preserves error type", func(t *testing.T) {
		seq := iter.From(E.Right[error](1), E.Left[int](assert.AnError))
		result := MonadReduce(seq, sum, 0)
		assert.Equal(t, E.Left[int](assert.AnError), result())
	})
}

func TestMonadReduce_EdgeCases(t *testing.T) {
	t.Run("accumulator with complex state", func(t *testing.T) {
		type State struct {
			sum   int
			count int
		}
		seq := iter.From(E.Right[string](1), E.Right[string](2), E.Right[string](3))
		result := MonadReduce(seq, func(acc State, x int) State {
			return State{sum: acc.sum + x, count: acc.count + 1}
		}, State{sum: 0, count: 0})
		expected := E.Right[string](State{sum: 6, count: 3})
		assert.Equal(t, expected, result())
	})

	t.Run("reducer with side effects", func(t *testing.T) {
		var sideEffects []int
		seq := iter.From(E.Right[string](1), E.Right[string](2), E.Right[string](3))
		result := MonadReduce(seq, func(acc, x int) int {
			sideEffects = append(sideEffects, x)
			return acc + x
		}, 0)
		assert.Equal(t, E.Right[string](6), result())
		assert.Equal(t, []int{1, 2, 3}, sideEffects)
	})

	t.Run("zero initial value", func(t *testing.T) {
		seq := iter.From(E.Right[string](5), E.Right[string](10))
		result := MonadReduce(seq, sum, 0)
		assert.Equal(t, E.Right[string](15), result())
	})
}

func TestReduce_Success(t *testing.T) {
	t.Run("curried version reduces all Right values", func(t *testing.T) {
		sum := Reduce[string](sum, 0)
		seq := iter.From(E.Right[string](1), E.Right[string](2), E.Right[string](3))
		result := sum(seq)
		assert.Equal(t, E.Right[string](6), result())
	})

	t.Run("reusable reducer function", func(t *testing.T) {
		multiply := Reduce[string](func(acc, x int) int { return acc * x }, 1)

		seq1 := iter.From(E.Right[string](2), E.Right[string](3))
		result1 := multiply(seq1)
		assert.Equal(t, E.Right[string](6), result1())

		seq2 := iter.From(E.Right[string](4), E.Right[string](5))
		result2 := multiply(seq2)
		assert.Equal(t, E.Right[string](20), result2())
	})

	t.Run("used in pipeline", func(t *testing.T) {
		result := F.Pipe1(
			iter.From(E.Right[string](1), E.Right[string](2), E.Right[string](3)),
			Reduce[string](sum, 0),
		)
		assert.Equal(t, E.Right[string](6), result())
	})

	t.Run("complex pipeline with map and reduce", func(t *testing.T) {
		result := F.Pipe2(
			iter.From(E.Right[string](1), E.Right[string](2), E.Right[string](3)),
			Map[string](utils.Double),
			Reduce[string](sum, 0),
		)
		assert.Equal(t, E.Right[string](12), result())
	})
}

func TestReduce_Failure(t *testing.T) {
	t.Run("curried version stops at Left", func(t *testing.T) {
		sum := Reduce[string](sum, 0)
		seq := iter.From(E.Right[string](1), E.Left[int]("error"), E.Right[string](3))
		result := sum(seq)
		assert.Equal(t, E.Left[int]("error"), result())
	})

	t.Run("error in pipeline", func(t *testing.T) {
		result := F.Pipe2(
			iter.From(E.Right[string](1), E.Left[int]("error"), E.Right[string](3)),
			Map[string](utils.Double),
			Reduce[string](sum, 0),
		)
		assert.Equal(t, E.Left[int]("error"), result())
	})
}

func TestReduce_EdgeCases(t *testing.T) {
	t.Run("empty sequence with curried version", func(t *testing.T) {
		sum := Reduce[string](sum, 100)
		seq := iter.From[E.Either[string, int]]()
		result := sum(seq)
		assert.Equal(t, E.Right[string](100), result())
	})

	t.Run("type transformation in pipeline", func(t *testing.T) {
		result := F.Pipe1(
			iter.From(E.Right[string](1), E.Right[string](2), E.Right[string](3)),
			Reduce[string](func(acc []int, x int) []int {
				return append(acc, x)
			}, []int{}),
		)
		assert.Equal(t, E.Right[string]([]int{1, 2, 3}), result())
	})
}

func TestReduce_Integration(t *testing.T) {
	t.Run("reduce after chain", func(t *testing.T) {
		result := F.Pipe2(
			iter.From(E.Right[string](1), E.Right[string](2)),
			Chain(func(n int) SeqEither[string, int] {
				return iter.From(E.Right[string](n), E.Right[string](n*10))
			}),
			Reduce[string](sum, 0),
		)
		assert.Equal(t, E.Right[string](33), result()) // 1 + 10 + 2 + 20
	})

	t.Run("reduce with filter-like behavior", func(t *testing.T) {
		result := F.Pipe1(
			iter.From(E.Right[string](1), E.Right[string](2), E.Right[string](3), E.Right[string](4)),
			Reduce[string](func(acc, x int) int {
				if x%2 == 0 {
					return acc + x
				}
				return acc
			}, 0),
		)
		assert.Equal(t, E.Right[string](6), result()) // 2 + 4
	})

	t.Run("reduce to find max", func(t *testing.T) {
		resultIO := F.Pipe1(
			iter.From(E.Right[string](3), E.Right[string](7), E.Right[string](2), E.Right[string](9), E.Right[string](1)),
			Reduce[string](func(acc, x int) int {
				if x > acc {
					return x
				}
				return acc
			}, 0),
		)
		assert.Equal(t, E.Right[string](9), resultIO())
	})

	t.Run("reduce to count elements", func(t *testing.T) {
		resultIO := F.Pipe1(
			iter.From(E.Right[string]("a"), E.Right[string]("b"), E.Right[string]("c")),
			Reduce[string](func(acc int, _ string) int { return acc + 1 }, 0),
		)
		assert.Equal(t, E.Right[string](3), resultIO())
	})
}

// TestCollect_Success tests basic Collect functionality
func TestCollect_Success(t *testing.T) {
	t.Run("collects all Right values into slice", func(t *testing.T) {
		seq := iter.From(
			E.Right[string](1),
			E.Right[string](2),
			E.Right[string](3),
		)
		result := Collect(seq)
		assert.Equal(t, E.Right[string]([]int{1, 2, 3}), result())
	})

	t.Run("stops on first Left and returns error", func(t *testing.T) {
		seq := iter.From(
			E.Right[string](1),
			E.Right[string](2),
			E.Left[int]("error"),
			E.Right[string](4),
		)
		result := Collect(seq)
		assert.Equal(t, E.Left[[]int]("error"), result())
	})

	t.Run("collects empty sequence into empty slice", func(t *testing.T) {
		seq := iter.From[E.Either[string, int]]()
		result := Collect(seq)
		assert.Equal(t, E.Right[string, []int](nil), result())
	})

	t.Run("collects single Right value", func(t *testing.T) {
		seq := iter.From(E.Right[string](42))
		result := Collect(seq)
		assert.Equal(t, E.Right[string]([]int{42}), result())
	})

	t.Run("returns first Left immediately", func(t *testing.T) {
		seq := iter.From(
			E.Left[int]("first error"),
			E.Left[int]("second error"),
		)
		result := Collect(seq)
		assert.Equal(t, E.Left[[]int]("first error"), result())
	})

	t.Run("works with different types", func(t *testing.T) {
		t.Run("string", func(t *testing.T) {
			seq := iter.From(
				E.Right[int]("hello"),
				E.Right[int]("world"),
			)
			result := Collect(seq)
			assert.Equal(t, E.Right[int]([]string{"hello", "world"}), result())
		})

		t.Run("struct", func(t *testing.T) {
			type Person struct {
				Name string
				Age  int
			}
			seq := iter.From(
				E.Right[string](Person{"Alice", 30}),
				E.Right[string](Person{"Bob", 25}),
			)
			result := Collect(seq)
			expected := E.Right[string]([]Person{{"Alice", 30}, {"Bob", 25}})
			assert.Equal(t, expected, result())
		})
	})
}

// TestCollect_EdgeCases tests edge cases
func TestCollect_EdgeCases(t *testing.T) {
	t.Run("handles zero values", func(t *testing.T) {
		seq := iter.From(
			E.Right[string](0),
			E.Right[string](0),
		)
		result := Collect(seq)
		assert.Equal(t, E.Right[string]([]int{0, 0}), result())
	})

	t.Run("handles empty strings", func(t *testing.T) {
		seq := iter.From(
			E.Right[int](""),
			E.Right[int]("a"),
			E.Right[int](""),
		)
		result := Collect(seq)
		assert.Equal(t, E.Right[int]([]string{"", "a", ""}), result())
	})

	t.Run("handles nil pointers", func(t *testing.T) {
		var nilPtr *int
		val := 42
		seq := iter.From(
			E.Right[string](nilPtr),
			E.Right[string](&val),
		)
		result := Collect(seq)
		assert.Equal(t, E.Right[string]([]*int{nilPtr, &val}), result())
	})

	t.Run("handles large sequence", func(t *testing.T) {
		// Create a large sequence of Right values
		eithers := make([]E.Either[string, int], 1000)
		for i := range 1000 {
			eithers[i] = E.Right[string](i)
		}
		seq := iter.From(eithers...)
		resultIO := Collect(seq)
		result := resultIO()

		assert.True(t, E.IsRight(result))
		collected := F.Pipe1(result, E.Fold(
			func(e string) []int { t.Fatal(e); return nil },
			F.Identity[[]int],
		))
		assert.Len(t, collected, 1000)
		assert.Equal(t, 0, collected[0])
		assert.Equal(t, 999, collected[999])
	})
}

// TestCollect_Integration tests integration with other operations
func TestCollect_Integration(t *testing.T) {
	t.Run("works with MonadMap", func(t *testing.T) {
		seq := iter.From(
			E.Right[string](1),
			E.Right[string](2),
			E.Right[string](3),
		)
		mapped := MonadMap(seq, func(x int) int { return x * 2 })
		result := Collect(mapped)
		assert.Equal(t, E.Right[string]([]int{2, 4, 6}), result())
	})

	t.Run("works with MonadChain", func(t *testing.T) {
		seq := iter.From(
			E.Right[string](1),
			E.Right[string](2),
		)
		chained := MonadChain(seq, func(x int) SeqEither[string, int] {
			return iter.From(E.Right[string](x), E.Right[string](x*10))
		})
		result := Collect(chained)
		assert.Equal(t, E.Right[string]([]int{1, 10, 2, 20}), result())
	})

	t.Run("works with Flatten", func(t *testing.T) {
		nested := iter.From(
			E.Right[string](iter.From(E.Right[string](1), E.Right[string](2))),
			E.Right[string](iter.From(E.Right[string](3))),
		)
		flattened := Flatten(nested)
		result := Collect(flattened)
		assert.Equal(t, E.Right[string]([]int{1, 2, 3}), result())
	})

	t.Run("works with FromIO", func(t *testing.T) {
		counter := 0
		io := func() int {
			counter++
			return counter * 10
		}
		seq := FromIO[string](io)
		result := Collect(seq)
		assert.Equal(t, E.Right[string]([]int{10}), result())
		assert.Equal(t, 1, counter)
	})

	t.Run("works with FromLazy", func(t *testing.T) {
		lazy := func() string { return "computed" }
		seq := FromLazy[int](lazy)
		result := Collect(seq)
		assert.Equal(t, E.Right[int]([]string{"computed"}), result())
	})

	t.Run("works with FromIOEither", func(t *testing.T) {
		t.Run("success case", func(t *testing.T) {
			ioEither := func() E.Either[string, int] {
				return E.Right[string](42)
			}
			seq := FromIOEither(ioEither)
			result := Collect(seq)
			assert.Equal(t, E.Right[string]([]int{42}), result())
		})

		t.Run("error case", func(t *testing.T) {
			ioEither := func() E.Either[string, int] {
				return E.Left[int]("io error")
			}
			seq := FromIOEither(ioEither)
			result := Collect(seq)
			assert.Equal(t, E.Left[[]int]("io error"), result())
		})
	})

	t.Run("short-circuits on error in chain", func(t *testing.T) {
		executionCount := 0
		seq := iter.From(
			E.Right[string](1),
			E.Right[string](2),
			E.Left[int]("error"),
			E.Right[string](4),
		)
		mapped := MonadMap(seq, func(x int) int {
			executionCount++
			return x * 2
		})
		result := Collect(mapped)

		assert.Equal(t, E.Left[[]int]("error"), result())
		// Map should only execute for the first two Right values
		assert.Equal(t, 2, executionCount)
	})
}

// TestCollect_Comparison tests Collect behavior
func TestCollect_Comparison(t *testing.T) {
	t.Run("Collect is equivalent to MonadReduce with append", func(t *testing.T) {
		seq := iter.From(
			E.Right[string](1),
			E.Right[string](2),
			E.Right[string](3),
		)

		result1 := Collect(seq)
		result2 := MonadReduce(seq, func(acc []int, x int) []int {
			return append(acc, x)
		}, []int{})

		assert.Equal(t, result1(), result2())
	})

	t.Run("Collect stops on first error like MonadReduce", func(t *testing.T) {
		seq := iter.From(
			E.Right[string](1),
			E.Left[int]("error"),
			E.Right[string](3),
		)

		result1 := Collect(seq)
		result2 := MonadReduce(seq, func(acc []int, x int) []int {
			return append(acc, x)
		}, []int{})

		assert.Equal(t, result1(), result2())
		assert.True(t, E.IsLeft(result1()))
	})
}
