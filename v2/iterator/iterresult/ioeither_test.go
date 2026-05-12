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
	"slices"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
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
			yield(R.Right(1))
			yield(R.Right(2))
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
			yield(R.Right(1))
			yield(R.Left[int](testErr))
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
			yield(R.Right(1))
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
			yield(R.Right(1))
			yield(R.Right(2))
		}

		expand := func(x int) SeqResult[int] {
			return func(yield func(R.Result[int]) bool) {
				yield(R.Right(x))
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
			yield(R.Right(1))
			yield(R.Left[int](testErr))
		}

		expand := func(x int) SeqResult[int] {
			return func(yield func(R.Result[int]) bool) {
				yield(R.Right(x))
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
			yield(R.Right(1))
			yield(R.Right(2))
		}

		expand := Chain(func(x int) SeqResult[int] {
			return func(yield func(R.Result[int]) bool) {
				yield(R.Right(x))
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
			yield(R.Right(1))
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
			yield(R.Right(1))
			yield(R.Right(2))
			yield(R.Right(3))
		}

		result := MonadReduce(input, func(acc, x int) int { return acc + x }, 0)
		value := R.MonadFold(result,
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
		err := R.MonadFold(result,
			F.Identity[error],
			func(int) error { t.Fatal("expected Left"); return nil },
		)
		assert.Equal(t, testErr, err)
	})
}

func TestReduce(t *testing.T) {
	t.Run("curried reduce function", func(t *testing.T) {
		input := func(yield func(R.Result[int]) bool) {
			yield(R.Right(1))
			yield(R.Right(2))
			yield(R.Right(3))
		}

		sum := Reduce(func(acc, x int) int { return acc + x }, 0)
		result := sum(input)
		value := R.MonadFold(result,
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
			yield(R.Right(1))
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
			yield(R.Right(5))
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
			yield(R.Right(1))
			yield(R.Right(2))
		}
		inner2 := func(yield func(R.Result[int]) bool) {
			yield(R.Right(3))
			yield(R.Right(4))
		}

		outer := func(yield func(R.Result[SeqResult[int]]) bool) {
			yield(R.Right[SeqResult[int]](inner1))
			yield(R.Right[SeqResult[int]](inner2))
		}

		result := Flatten(outer)
		collected := slices.Collect(F.Pipe1(result, GetOrElse(
			func(e error) Seq[int] { t.Fatal(e); return nil },
		)))
		assert.Equal(t, []int{1, 2, 3, 4}, collected)
	})
}
