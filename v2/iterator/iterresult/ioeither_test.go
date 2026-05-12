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
	"slices"
	"strings"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/iterator/iter"
	R "github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

func TestOf(t *testing.T) {
	t.Run("creates Right value", func(t *testing.T) {
		result := Of(42)
		collected := slices.Collect(F.Pipe1(result, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{42}, collected)
	})
}

func TestLeft(t *testing.T) {
	t.Run("creates Left value", func(t *testing.T) {
		result := Left[int](errors.New("test error"))
		var err error
		for e := range result {
			err = R.MonadFold(e,
				F.Identity[error],
				func(v int) error { t.Fatal("expected Left"); return nil },
			)
			break
		}
		assert.Error(t, err)
		assert.Equal(t, "test error", err.Error())
	})
}

func TestFromIO_Success(t *testing.T) {
	t.Run("converts IO computation to single-element success sequence", func(t *testing.T) {
		io := func() int { return 42 }
		seq := FromIO(io)
		collected := slices.Collect(F.Pipe1(seq, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{42}, collected)
	})

	t.Run("executes IO computation when sequence is consumed", func(t *testing.T) {
		executed := false
		io := func() string {
			executed = true
			return "hello"
		}
		seq := FromIO(io)

		// IO should not be executed yet
		assert.False(t, executed)

		// Consume the sequence
		collected := slices.Collect(F.Pipe1(seq, GetOrElse(
			func(e error) Seq[string] { t.Fatal(e); return nil },
		)))

		// Now IO should be executed
		assert.True(t, executed)
		assert.Equal(t, []string{"hello"}, collected)
	})

	t.Run("works with different types", func(t *testing.T) {
		t.Run("string", func(t *testing.T) {
			io := func() string { return "test" }
			seq := FromIO(io)
			collected := slices.Collect(F.Pipe1(seq, GetOrElse(
				func(e error) Seq[string] { t.Fatal(e); return nil },
			)))
			assert.Equal(t, []string{"test"}, collected)
		})

		t.Run("struct", func(t *testing.T) {
			type Person struct {
				Name string
				Age  int
			}
			io := func() Person { return Person{Name: "Alice", Age: 30} }
			seq := FromIO(io)
			collected := slices.Collect(F.Pipe1(seq, GetOrElse(
				func(e error) Seq[Person] { t.Fatal(e); return nil },
			)))
			assert.Equal(t, []Person{{Name: "Alice", Age: 30}}, collected)
		})

		t.Run("pointer", func(t *testing.T) {
			value := 100
			io := func() *int { return &value }
			seq := FromIO(io)
			collected := slices.Collect(F.Pipe1(seq, GetOrElse(
				func(e error) Seq[*int] { t.Fatal(e); return nil },
			)))
			assert.Len(t, collected, 1)
			assert.Equal(t, 100, *collected[0])
		})
	})

	t.Run("can be composed with other operations", func(t *testing.T) {
		io := func() int { return 10 }
		seq := F.Pipe1(
			FromIO(io),
			Map(func(x int) int { return x * 2 }),
		)
		collected := slices.Collect(F.Pipe1(seq, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{20}, collected)
	})

	t.Run("can be used in chain operations", func(t *testing.T) {
		io := func() int { return 3 }
		seq := F.Pipe1(
			FromIO(io),
			Chain(func(n int) SeqResult[int] {
				return FromSeq(iter.Replicate(n, n))
			}),
		)
		collected := slices.Collect(F.Pipe1(seq, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{3, 3, 3}, collected)
	})
}

func TestFromIO_EdgeCases(t *testing.T) {
	t.Run("handles zero value", func(t *testing.T) {
		io := func() int { return 0 }
		seq := FromIO(io)
		collected := slices.Collect(F.Pipe1(seq, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{0}, collected)
	})

	t.Run("handles empty string", func(t *testing.T) {
		io := func() string { return "" }
		seq := FromIO(io)
		collected := slices.Collect(F.Pipe1(seq, GetOrElse(
			func(e error) Seq[string] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []string{""}, collected)
	})

	t.Run("handles nil pointer", func(t *testing.T) {
		io := func() *int { return nil }
		seq := FromIO(io)
		collected := slices.Collect(F.Pipe1(seq, GetOrElse(
			func(e error) Seq[*int] { t.Fatal(e); return nil },
		)))
		assert.Len(t, collected, 1)
		assert.Nil(t, collected[0])
	})
}

func TestFromIO_Integration(t *testing.T) {
	t.Run("multiple iterations execute IO multiple times", func(t *testing.T) {
		counter := 0
		io := func() int {
			counter++
			return counter
		}
		seq := FromIO(io)

		// First iteration
		result1 := slices.Collect(F.Pipe1(seq, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{1}, result1)

		// Second iteration - IO executes again
		result2 := slices.Collect(F.Pipe1(seq, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{2}, result2)
	})

	t.Run("works with MonadMap", func(t *testing.T) {
		io := func() int { return 5 }
		seq := MonadMap(FromIO(io), func(n int) string { return strings.Repeat("*", n) })
		collected := slices.Collect(F.Pipe1(seq, GetOrElse(
			func(e error) Seq[string] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []string{"*****"}, collected)
	})

	t.Run("works with MonadChain", func(t *testing.T) {
		io := func() int { return 2 }
		seq := MonadChain(
			FromIO(io),
			func(n int) SeqResult[int] {
				return FromSeq(iter.Replicate(n, n))
			},
		)
		collected := slices.Collect(F.Pipe1(seq, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{2, 2}, collected)
	})
}

func TestFromLazy_Success(t *testing.T) {
	t.Run("converts Lazy computation to single-element success sequence", func(t *testing.T) {
		lazy := func() int { return 42 }
		seq := FromLazy(lazy)
		collected := slices.Collect(F.Pipe1(seq, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{42}, collected)
	})

	t.Run("defers computation until sequence is consumed", func(t *testing.T) {
		executed := false
		lazy := func() string {
			executed = true
			return "lazy value"
		}
		seq := FromLazy(lazy)

		// Lazy computation should not be executed yet
		assert.False(t, executed)

		// Consume the sequence
		collected := slices.Collect(F.Pipe1(seq, GetOrElse(
			func(e error) Seq[string] { t.Fatal(e); return nil },
		)))

		// Now lazy computation should be executed
		assert.True(t, executed)
		assert.Equal(t, []string{"lazy value"}, collected)
	})

	t.Run("works with different types", func(t *testing.T) {
		t.Run("bool", func(t *testing.T) {
			lazy := func() bool { return true }
			seq := FromLazy(lazy)
			collected := slices.Collect(F.Pipe1(seq, GetOrElse(
				func(e error) Seq[bool] { t.Fatal(e); return nil },
			)))
			assert.Equal(t, []bool{true}, collected)
		})

		t.Run("slice", func(t *testing.T) {
			lazy := func() []int { return []int{1, 2, 3} }
			seq := FromLazy(lazy)
			collected := slices.Collect(F.Pipe1(seq, GetOrElse(
				func(e error) Seq[[]int] { t.Fatal(e); return nil },
			)))
			assert.Len(t, collected, 1)
			assert.Equal(t, []int{1, 2, 3}, collected[0])
		})

		t.Run("map", func(t *testing.T) {
			lazy := func() map[string]int {
				return map[string]int{"a": 1, "b": 2}
			}
			seq := FromLazy(lazy)
			collected := slices.Collect(F.Pipe1(seq, GetOrElse(
				func(e error) Seq[map[string]int] { t.Fatal(e); return nil },
			)))
			assert.Len(t, collected, 1)
			assert.Equal(t, map[string]int{"a": 1, "b": 2}, collected[0])
		})
	})

	t.Run("can be composed with other operations", func(t *testing.T) {
		lazy := func() int { return 5 }
		seq := F.Pipe1(
			FromLazy(lazy),
			Map(func(n int) int { return n * n }),
		)
		collected := slices.Collect(F.Pipe1(seq, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{25}, collected)
	})

	t.Run("can be used with MonadChain", func(t *testing.T) {
		lazy := func() int { return 3 }
		seq := MonadChain(
			FromLazy(lazy),
			func(n int) SeqResult[string] {
				return FromSeq(iter.Replicate(n, strings.Repeat("x", n)))
			},
		)
		collected := slices.Collect(F.Pipe1(seq, GetOrElse(
			func(e error) Seq[string] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []string{"xxx", "xxx", "xxx"}, collected)
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
		seq := FromLazy(lazy)
		collected := slices.Collect(F.Pipe1(seq, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{499500}, collected)
	})

	t.Run("handles function returning function", func(t *testing.T) {
		lazy := func() func(int) int {
			return func(x int) int { return x * 2 }
		}
		seq := FromLazy(lazy)
		collected := slices.Collect(F.Pipe1(seq, GetOrElse(
			func(e error) Seq[func(int) int] { t.Fatal(e); return nil },
		)))
		assert.Len(t, collected, 1)
		assert.Equal(t, 10, collected[0](5))
	})
}

func TestFromLazy_Integration(t *testing.T) {
	t.Run("multiple iterations execute lazy computation multiple times", func(t *testing.T) {
		counter := 0
		lazy := func() int {
			counter++
			return counter * 10
		}
		seq := FromLazy(lazy)

		// First iteration
		result1 := slices.Collect(F.Pipe1(seq, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{10}, result1)

		// Second iteration - lazy computation executes again
		result2 := slices.Collect(F.Pipe1(seq, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{20}, result2)
	})

	t.Run("works with MonadMapLeft", func(t *testing.T) {
		lazy := func() int { return 7 }
		// Even though we map left, the success value passes through
		seq := MonadMapLeft(FromLazy(lazy), func(e error) error { return errors.New("mapped: " + e.Error()) })
		collected := slices.Collect(F.Pipe1(seq, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{7}, collected)
	})

	t.Run("works with Fold", func(t *testing.T) {
		lazy := func() int { return 5 }
		seq := FromLazy(lazy)
		folded := MonadFold(
			seq,
			func(e error) Seq[string] { return iter.Of("error: " + e.Error()) },
			func(n int) Seq[string] { return iter.Of(strings.Repeat("*", n)) },
		)
		result := slices.Collect(folded)
		assert.Equal(t, []string{"*****"}, result)
	})
}

func TestFromIOResult_Success(t *testing.T) {
	t.Run("converts IOResult success to single-element success sequence", func(t *testing.T) {
		ior := func() R.Result[int] { return R.Of(42) }
		seq := FromIOResult(ior)
		collected := slices.Collect(F.Pipe1(seq, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{42}, collected)
	})

	t.Run("converts IOResult error to single-element error sequence", func(t *testing.T) {
		ior := func() R.Result[int] { return R.Left[int](errors.New("test error")) }
		seq := FromIOResult(ior)
		var err error
		for r := range seq {
			err = R.MonadFold(r,
				F.Identity[error],
				func(v int) error { t.Fatal("expected error"); return nil },
			)
			break
		}
		assert.Error(t, err)
		assert.Equal(t, "test error", err.Error())
	})

	t.Run("executes IOResult when sequence is consumed", func(t *testing.T) {
		executed := false
		ior := func() R.Result[string] {
			executed = true
			return R.Of("result")
		}
		seq := FromIOResult(ior)

		// IOResult should not be executed yet
		assert.False(t, executed)

		// Consume the sequence
		collected := slices.Collect(F.Pipe1(seq, GetOrElse(
			func(e error) Seq[string] { t.Fatal(e); return nil },
		)))

		// Now IOResult should be executed
		assert.True(t, executed)
		assert.Equal(t, []string{"result"}, collected)
	})

	t.Run("works with different error scenarios", func(t *testing.T) {
		t.Run("specific error", func(t *testing.T) {
			ior := func() R.Result[int] {
				return R.Left[int](errors.New("validation failed"))
			}
			seq := FromIOResult(ior)
			var err error
			for r := range seq {
				err = R.MonadFold(r,
					F.Identity[error],
					func(v int) error { t.Fatal("expected error"); return nil },
				)
				break
			}
			assert.Error(t, err)
			assert.Equal(t, "validation failed", err.Error())
		})

		t.Run("wrapped error", func(t *testing.T) {
			baseErr := errors.New("base error")
			ior := func() R.Result[int] {
				return R.Left[int](fmt.Errorf("wrapped: %w", baseErr))
			}
			seq := FromIOResult(ior)
			var err error
			for r := range seq {
				err = R.MonadFold(r,
					F.Identity[error],
					func(v int) error { t.Fatal("expected error"); return nil },
				)
				break
			}
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "wrapped")
			assert.ErrorIs(t, err, baseErr)
		})
	})

	t.Run("can be composed with other operations", func(t *testing.T) {
		ior := func() R.Result[int] { return R.Of(10) }
		seq := F.Pipe1(
			FromIOResult(ior),
			Map(func(x int) int { return x * 2 }),
		)
		collected := slices.Collect(F.Pipe1(seq, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{20}, collected)
	})

	t.Run("error values stop chain operations", func(t *testing.T) {
		ior := func() R.Result[int] { return R.Left[int](errors.New("error")) }
		seq := F.Pipe1(
			FromIOResult(ior),
			Map(func(x int) int { return x * 2 }),
		)
		var err error
		for r := range seq {
			err = R.MonadFold(r,
				F.Identity[error],
				func(v int) error { t.Fatal("expected error"); return nil },
			)
			break
		}
		assert.Error(t, err)
	})
}

func TestFromIOResult_EdgeCases(t *testing.T) {
	t.Run("handles zero value in success", func(t *testing.T) {
		ior := func() R.Result[int] { return R.Of(0) }
		seq := FromIOResult(ior)
		collected := slices.Collect(F.Pipe1(seq, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{0}, collected)
	})

	t.Run("handles nil pointer in success", func(t *testing.T) {
		ior := func() R.Result[*int] {
			var nilPtr *int
			return R.Of(nilPtr)
		}
		seq := FromIOResult(ior)
		collected := slices.Collect(F.Pipe1(seq, GetOrElse(
			func(e error) Seq[*int] { t.Fatal(e); return nil },
		)))
		assert.Len(t, collected, 1)
		assert.Nil(t, collected[0])
	})
}

func TestFromIOResult_Integration(t *testing.T) {
	t.Run("multiple iterations execute IOResult multiple times", func(t *testing.T) {
		counter := 0
		ior := func() R.Result[int] {
			counter++
			if counter%2 == 0 {
				return R.Left[int](errors.New("even"))
			}
			return R.Of(counter)
		}
		seq := FromIOResult(ior)

		// First iteration - odd
		result1 := slices.Collect(F.Pipe1(seq, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{1}, result1)

		// Second iteration - even (error)
		var err error
		for r := range seq {
			err = R.MonadFold(r,
				F.Identity[error],
				func(v int) error { t.Fatal("expected error"); return nil },
			)
			break
		}
		assert.Error(t, err)
		assert.Equal(t, "even", err.Error())

		// Third iteration - odd
		result3 := slices.Collect(F.Pipe1(seq, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{3}, result3)
	})

	t.Run("works with MonadChain", func(t *testing.T) {
		ior := func() R.Result[int] { return R.Of(2) }
		seq := MonadChain(
			FromIOResult(ior),
			func(n int) SeqResult[int] {
				return FromSeq(iter.From(n, n*2))
			},
		)
		collected := slices.Collect(F.Pipe1(seq, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{2, 4}, collected)
	})

	t.Run("works with GetOrElse", func(t *testing.T) {
		t.Run("success value", func(t *testing.T) {
			ior := func() R.Result[int] { return R.Of(42) }
			seq := FromIOResult(ior)
			result := slices.Collect(GetOrElse(func(e error) Seq[int] {
				return iter.Of(0)
			})(seq))
			assert.Equal(t, []int{42}, result)
		})

		t.Run("error value", func(t *testing.T) {
			ior := func() R.Result[int] { return R.Left[int](errors.New("error")) }
			seq := FromIOResult(ior)
			result := slices.Collect(GetOrElse(func(e error) Seq[int] {
				return iter.Of(-1)
			})(seq))
			assert.Equal(t, []int{-1}, result)
		})
	})

	t.Run("works with OrElse for error recovery", func(t *testing.T) {
		ior := func() R.Result[int] { return R.Left[int](errors.New("not found")) }
		recover := OrElse(func(err error) SeqResult[int] {
			if err.Error() == "not found" {
				return Of(0) // default value
			}
			return Left[int](err)
		})
		seq := recover(FromIOResult(ior))
		collected := slices.Collect(F.Pipe1(seq, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{0}, collected)
	})
}

func TestFromIO_vs_FromLazy_vs_FromIOResult(t *testing.T) {
	t.Run("FromLazy delegates to FromIO", func(t *testing.T) {
		value := 42
		io := func() int { return value }
		lazy := func() int { return value }

		seqIO := FromIO(io)
		seqLazy := FromLazy(lazy)

		resultIO := slices.Collect(F.Pipe1(seqIO, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		resultLazy := slices.Collect(F.Pipe1(seqLazy, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))

		assert.Equal(t, resultIO, resultLazy)
	})

	t.Run("FromIO and FromLazy always produce success", func(t *testing.T) {
		io := func() int { return 42 }
		lazy := func() int { return 42 }

		seqIO := FromIO(io)
		seqLazy := FromLazy(lazy)

		// Both should always be success
		for r := range seqIO {
			assert.True(t, R.IsRight(r))
			break
		}
		for r := range seqLazy {
			assert.True(t, R.IsRight(r))
			break
		}
	})

	t.Run("FromIOResult can produce error or success", func(t *testing.T) {
		iorSuccess := func() R.Result[int] { return R.Of(42) }
		iorError := func() R.Result[int] { return R.Left[int](errors.New("error")) }

		seqSuccess := FromIOResult(iorSuccess)
		seqError := FromIOResult(iorError)

		for r := range seqSuccess {
			assert.True(t, R.IsRight(r))
			break
		}
		for r := range seqError {
			assert.True(t, R.IsLeft(r))
			break
		}
	})

	t.Run("all handle side effects similarly", func(t *testing.T) {
		counterIO := 0
		counterLazy := 0
		counterIOR := 0

		io := func() int {
			counterIO++
			return counterIO
		}
		lazy := func() int {
			counterLazy++
			return counterLazy
		}
		ior := func() R.Result[int] {
			counterIOR++
			return R.Of(counterIOR)
		}

		seqIO := FromIO(io)
		seqLazy := FromLazy(lazy)
		seqIOR := FromIOResult(ior)

		// All should execute on consumption
		resultIO := slices.Collect(F.Pipe1(seqIO, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		resultLazy := slices.Collect(F.Pipe1(seqLazy, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		resultIOR := slices.Collect(F.Pipe1(seqIOR, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))

		assert.Equal(t, []int{1}, resultIO)
		assert.Equal(t, []int{1}, resultLazy)
		assert.Equal(t, []int{1}, resultIOR)
		assert.Equal(t, 1, counterIO)
		assert.Equal(t, 1, counterLazy)
		assert.Equal(t, 1, counterIOR)

	})
}

func TestRight(t *testing.T) {
	t.Run("creates Right value", func(t *testing.T) {
		result := Right(42)
		collected := slices.Collect(F.Pipe1(result, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{42}, collected)
	})
}

func TestMonadMap(t *testing.T) {
	t.Run("maps over Right values", func(t *testing.T) {
		input := func(yield func(R.Result[int]) bool) {
			if !yield(R.Right(1)) {
				return
			}
			if !yield(R.Right(2)) {
				return
			}
			yield(R.Right(3))
		}
		result := MonadMap(input, func(x int) int { return x * 2 })
		collected := slices.Collect(F.Pipe1(result, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{2, 4, 6}, collected)
	})

	t.Run("preserves Left values", func(t *testing.T) {
		testErr := errors.New("test error")
		input := func(yield func(R.Result[int]) bool) {
			if !yield(R.Right(1)) {
				return
			}
			if !yield(R.Left[int](testErr)) {
				return
			}
			yield(R.Right(3))
		}
		result := MonadMap(input, func(x int) int { return x * 2 })

		var foundError bool
		for e := range result {
			if R.IsLeft(e) {
				foundError = true
				err := R.MonadFold(e, F.Identity[error], func(int) error { return nil })
				assert.Equal(t, testErr, err)
			}
		}
		assert.True(t, foundError)
	})
}

func TestMap(t *testing.T) {
	t.Run("curried map function", func(t *testing.T) {
		input := func(yield func(R.Result[int]) bool) {
			if !yield(R.Right(1)) {
				return
			}
			yield(R.Right(2))
		}
		double := Map(func(x int) int { return x * 2 })
		result := double(input)
		collected := slices.Collect(F.Pipe1(result, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{2, 4}, collected)
	})
}

func TestMonadChain(t *testing.T) {
	t.Run("chains successful computations", func(t *testing.T) {
		input := func(yield func(R.Result[int]) bool) {
			if !yield(R.Right(1)) {
				return
			}
			yield(R.Right(2))
		}

		expand := func(x int) SeqResult[int] {
			return func(yield func(R.Result[int]) bool) {
				if !yield(R.Right(x)) {
					return
				}
				yield(R.Right(x * 10))
			}
		}

		result := MonadChain(input, expand)
		collected := slices.Collect(F.Pipe1(result, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{1, 10, 2, 20}, collected)
	})

	t.Run("stops on Left", func(t *testing.T) {
		testErr := errors.New("test error")
		input := func(yield func(R.Result[int]) bool) {
			if !yield(R.Right(1)) {
				return
			}
			yield(R.Left[int](testErr))
		}

		expand := func(x int) SeqResult[int] {
			return func(yield func(R.Result[int]) bool) {
				if !yield(R.Right(x)) {
					return
				}
				yield(R.Right(x * 10))
			}
		}

		result := MonadChain(input, expand)

		var foundError bool
		for e := range result {
			if R.IsLeft(e) {
				foundError = true
				err := R.MonadFold(e, F.Identity[error], func(int) error { return nil })
				assert.Equal(t, testErr, err)
			}
		}
		assert.True(t, foundError)
	})
}

func TestChain(t *testing.T) {
	t.Run("curried chain function", func(t *testing.T) {
		input := func(yield func(R.Result[int]) bool) {
			if !yield(R.Right(1)) {
				return
			}
			yield(R.Right(2))
		}

		expand := Chain(func(x int) SeqResult[int] {
			return func(yield func(R.Result[int]) bool) {
				if !yield(R.Right(x)) {
					return
				}
				yield(R.Right(x * 10))
			}
		})

		result := expand(input)
		collected := slices.Collect(F.Pipe1(result, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{1, 10, 2, 20}, collected)
	})
}

func TestMonadAlt(t *testing.T) {
	t.Run("uses alternative on Left", func(t *testing.T) {
		first := func(yield func(R.Result[int]) bool) {
			yield(R.Left[int](errors.New("error")))
		}

		second := func() SeqResult[int] {
			return func(yield func(R.Result[int]) bool) {
				yield(R.Right(42))
			}
		}

		result := MonadAlt(first, second)
		collected := slices.Collect(F.Pipe1(result, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{42}, collected)
	})

	t.Run("keeps Right values", func(t *testing.T) {
		first := func(yield func(R.Result[int]) bool) {
			if !yield(R.Right(1)) {
				return
			}
			yield(R.Right(2))
		}

		second := func() SeqResult[int] {
			return func(yield func(R.Result[int]) bool) {
				yield(R.Right(99))
			}
		}

		result := MonadAlt(first, second)
		collected := slices.Collect(F.Pipe1(result, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{1, 2}, collected)
	})
}

func TestMonadReduce(t *testing.T) {
	t.Run("reduces Right values", func(t *testing.T) {
		input := func(yield func(R.Result[int]) bool) {
			if !yield(R.Right(1)) {
				return
			}
			if !yield(R.Right(2)) {
				return
			}
			yield(R.Right(3))
		}

		result := MonadReduce(input, func(acc, x int) int { return acc + x }, 0)
		value := R.MonadFold(result(),
			func(e error) int { t.Fatal(e); return 0 },
			F.Identity[int],
		)
		assert.Equal(t, 6, value)
	})

	t.Run("stops on Left", func(t *testing.T) {
		testErr := errors.New("test error")
		input := func(yield func(R.Result[int]) bool) {
			if !yield(R.Right(1)) {
				return
			}
			if !yield(R.Left[int](testErr)) {
				return
			}
			yield(R.Right(3))
		}

		result := MonadReduce(input, func(acc, x int) int { return acc + x }, 0)
		err := R.MonadFold(result(),
			F.Identity[error],
			func(int) error { t.Fatal("expected Left"); return nil },
		)
		assert.Equal(t, testErr, err)
	})
}

func TestReduce(t *testing.T) {
	t.Run("curried reduce function", func(t *testing.T) {
		input := func(yield func(R.Result[int]) bool) {
			if !yield(R.Right(1)) {
				return
			}
			if !yield(R.Right(2)) {
				return
			}
			yield(R.Right(3))
		}

		sum := Reduce(func(acc, x int) int { return acc + x }, 0)
		result := sum(input)
		value := R.MonadFold(result(),
			func(e error) int { t.Fatal(e); return 0 },
			F.Identity[int],
		)
		assert.Equal(t, 6, value)
	})
}

func TestOrElse(t *testing.T) {
	t.Run("recovers from Left", func(t *testing.T) {
		input := func(yield func(R.Result[int]) bool) {
			yield(R.Left[int](errors.New("not found")))
		}

		recover := OrElse(func(err error) SeqResult[int] {
			if err.Error() == "not found" {
				return Right(0)
			}
			return Left[int](err)
		})

		result := recover(input)
		collected := slices.Collect(F.Pipe1(result, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{0}, collected)
	})

	t.Run("keeps Right values", func(t *testing.T) {
		input := func(yield func(R.Result[int]) bool) {
			yield(R.Right(42))
		}

		recover := OrElse(func(err error) SeqResult[int] {
			return Right(0)
		})

		result := recover(input)
		collected := slices.Collect(F.Pipe1(result, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{42}, collected)
	})
}

func TestFromEither(t *testing.T) {
	t.Run("lifts Right value", func(t *testing.T) {
		either := R.Right(42)
		result := FromEither(either)
		collected := slices.Collect(F.Pipe1(result, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{42}, collected)
	})

	t.Run("lifts Left value", func(t *testing.T) {
		testErr := errors.New("test error")
		either := R.Left[int](testErr)
		result := FromEither(either)

		var err error
		for e := range result {
			err = R.MonadFold(e,
				F.Identity[error],
				func(int) error { t.Fatal("expected Left"); return nil },
			)
			break
		}
		assert.Equal(t, testErr, err)
	})
}

func TestMapTo(t *testing.T) {
	t.Run("replaces Right values with constant", func(t *testing.T) {
		input := func(yield func(R.Result[int]) bool) {
			if !yield(R.Right(1)) {
				return
			}
			yield(R.Right(2))
		}

		result := F.Pipe1(input, MapTo[int]("constant"))
		collected := slices.Collect(F.Pipe1(result, GetOrElse(
			func(e error) Seq[string] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []string{"constant", "constant"}, collected)
	})
}

func TestMonadBiMap(t *testing.T) {
	t.Run("maps both Left and Right", func(t *testing.T) {
		input := func(yield func(R.Result[int]) bool) {
			if !yield(R.Right(5)) {
				return
			}
			yield(R.Left[int](errors.New("err")))
		}

		result := MonadBiMap(input,
			func(e error) error { return errors.New("mapped: " + e.Error()) },
			func(x int) int { return x * 2 },
		)

		var values []int
		var errs []string
		for e := range result {
			R.MonadFold(e,
				func(err error) int { errs = append(errs, err.Error()); return 0 },
				func(v int) int { values = append(values, v); return 0 },
			)
		}

		assert.Equal(t, []int{10}, values)
		assert.Equal(t, []string{"mapped: err"}, errs)
	})
}

func TestFlatten(t *testing.T) {
	t.Run("flattens nested SeqResult", func(t *testing.T) {
		inner1 := func(yield func(R.Result[int]) bool) {
			if !yield(R.Right(1)) {
				return
			}
			yield(R.Right(2))
		}
		inner2 := func(yield func(R.Result[int]) bool) {
			if !yield(R.Right(3)) {
				return
			}
			yield(R.Right(4))
		}

		outer := func(yield func(R.Result[SeqResult[int]]) bool) {
			if !yield(R.Right[SeqResult[int]](inner1)) {
				return
			}
			yield(R.Right[SeqResult[int]](inner2))
		}

		result := Flatten(outer)
		collected := slices.Collect(F.Pipe1(result, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{1, 2, 3, 4}, collected)
	})
}
