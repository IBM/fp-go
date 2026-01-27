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

package result

import (
	"errors"
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	IO "github.com/IBM/fp-go/v2/io"
	O "github.com/IBM/fp-go/v2/option"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

func TestIsLeft(t *testing.T) {
	err := errors.New("Some error")
	withError := Left[string](err)

	assert.True(t, IsLeft(withError))
	assert.False(t, IsRight(withError))
}

func TestIsRight(t *testing.T) {
	noError := Right("Carsten")

	assert.True(t, IsRight(noError))
	assert.False(t, IsLeft(noError))
}

func TestMapEither(t *testing.T) {
	e := errors.New("s")
	assert.Equal(t, F.Pipe1(Right("abc"), Map(utils.StringLen)), Right(3))

	val2 := F.Pipe1(Left[string](e), Map(utils.StringLen))
	exp2 := Left[int](e)

	assert.Equal(t, val2, exp2)
}

func TestUnwrapError(t *testing.T) {
	a := ""
	err := errors.New("Some error")
	withError := Left[string](err)

	res, extracted := UnwrapError(withError)
	assert.Equal(t, a, res)
	assert.Equal(t, extracted, err)

}

func TestReduce(t *testing.T) {

	s := S.Semigroup

	assert.Equal(t, "foobar", F.Pipe1(Right("bar"), Reduce(s.Concat, "foo")))
	assert.Equal(t, "foo", F.Pipe1(Left[string](errors.New("bar")), Reduce(s.Concat, "foo")))

}
func TestAp(t *testing.T) {
	f := S.Size

	maError := errors.New("maError")
	mabError := errors.New("mabError")

	assert.Equal(t, Right(3), F.Pipe1(Right(f), Ap[int](Right("abc"))))
	assert.Equal(t, Left[int](maError), F.Pipe1(Right(f), Ap[int](Left[string](maError))))
	assert.Equal(t, Left[int](mabError), F.Pipe1(Left[func(string) int](mabError), Ap[int](Left[string](maError))))
}

func TestAlt(t *testing.T) {

	a := errors.New("a")
	b := errors.New("b")

	assert.Equal(t, Right(1), F.Pipe1(Right(1), Alt(F.Constant(Right(2)))))
	assert.Equal(t, Right(1), F.Pipe1(Right(1), Alt(F.Constant(Left[int](a)))))
	assert.Equal(t, Right(2), F.Pipe1(Left[int](b), Alt(F.Constant(Right(2)))))
	assert.Equal(t, Left[int](b), F.Pipe1(Left[int](a), Alt(F.Constant(Left[int](b)))))
}

func TestChainFirst(t *testing.T) {
	f := F.Flow2(S.Size, Right[int])
	maError := errors.New("maError")

	assert.Equal(t, Right("abc"), F.Pipe1(Right("abc"), ChainFirst(f)))
	assert.Equal(t, Left[string](maError), F.Pipe1(Left[string](maError), ChainFirst(f)))
}

func TestChainOptionK(t *testing.T) {
	a := errors.New("a")
	b := errors.New("b")

	f := ChainOptionK[int, int](F.Constant(a))(func(n int) Option[int] {
		if n > 0 {
			return O.Some(n)
		}
		return O.None[int]()
	})
	assert.Equal(t, Right(1), f(Right(1)))
	assert.Equal(t, Left[int](a), f(Right(-1)))
	assert.Equal(t, Left[int](b), f(Left[int](b)))
}

func TestFromOption(t *testing.T) {
	none := errors.New("none")

	assert.Equal(t, Left[int](none), FromOption[int](F.Constant(none))(O.None[int]()))
	assert.Equal(t, Right(1), FromOption[int](F.Constant(none))(O.Some(1)))
}

func TestStringer(t *testing.T) {
	e := Of("foo")
	exp := "Right[string](foo)"

	assert.Equal(t, exp, e.String())

	var s fmt.Stringer = &e
	assert.Equal(t, exp, s.String())
}

func TestFromIO(t *testing.T) {
	f := IO.Of("abc")
	e := FromIO(f)

	assert.Equal(t, Right("abc"), e)
}

// TestOrElse tests recovery from error
func TestOrElse(t *testing.T) {
	// Test basic recovery from Left
	recover := OrElse(func(err error) Result[int] {
		return Right(0) // default value
	})

	leftResult := Left[int](errors.New("fail"))
	assert.Equal(t, Right(0), recover(leftResult))

	// Test that Right values pass through unchanged
	rightResult := Right(42)
	assert.Equal(t, Right(42), recover(rightResult))

	// Test conditional recovery
	recoverSpecific := OrElse(func(err error) Result[int] {
		if err.Error() == "not found" {
			return Right(0) // default for not found
		}
		return Left[int](err) // propagate other errors
	})

	notFoundErr := errors.New("not found")
	assert.Equal(t, Right(0), recoverSpecific(Left[int](notFoundErr)))

	otherErr := errors.New("other error")
	assert.Equal(t, Left[int](otherErr), recoverSpecific(Left[int](otherErr)))
}

// TestZeroEqualsDefaultInitialization tests that Zero returns the same value as default initialization
func TestZeroEqualsDefaultInitialization(t *testing.T) {
	// Default initialization of Result
	var defaultInit Result[int]

	// Zero function
	zero := Zero[int]()

	// They should be equal
	assert.Equal(t, defaultInit, zero, "Zero should equal default initialization")
	assert.Equal(t, IsRight(defaultInit), IsRight(zero), "Both should be Right")
	assert.Equal(t, IsLeft(defaultInit), IsLeft(zero), "Both should not be Left")
}

// TestInstanceOf tests the InstanceOf function for type assertions
func TestInstanceOf(t *testing.T) {
	// Test successful type assertion with int
	t.Run("successful int assertion", func(t *testing.T) {
		var value any = 42
		result := InstanceOf[int](value)
		assert.True(t, IsRight(result))
		assert.Equal(t, Right(42), result)
	})

	// Test successful type assertion with string
	t.Run("successful string assertion", func(t *testing.T) {
		var value any = "hello"
		result := InstanceOf[string](value)
		assert.True(t, IsRight(result))
		assert.Equal(t, Right("hello"), result)
	})

	// Test successful type assertion with float64
	t.Run("successful float64 assertion", func(t *testing.T) {
		var value any = 3.14
		result := InstanceOf[float64](value)
		assert.True(t, IsRight(result))
		assert.Equal(t, Right(3.14), result)
	})

	// Test successful type assertion with pointer
	t.Run("successful pointer assertion", func(t *testing.T) {
		val := 42
		var value any = &val
		result := InstanceOf[*int](value)
		assert.True(t, IsRight(result))
		v, err := UnwrapError(result)
		assert.NoError(t, err)
		assert.Equal(t, 42, *v)
	})

	// Test successful type assertion with struct
	t.Run("successful struct assertion", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		var value any = Person{Name: "Alice", Age: 30}
		result := InstanceOf[Person](value)
		assert.True(t, IsRight(result))
		assert.Equal(t, Right(Person{Name: "Alice", Age: 30}), result)
	})

	// Test failed type assertion - int to string
	t.Run("failed int to string assertion", func(t *testing.T) {
		var value any = 42
		result := InstanceOf[string](value)
		assert.True(t, IsLeft(result))
		_, err := UnwrapError(result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "expected")
		assert.Contains(t, err.Error(), "got")
	})

	// Test failed type assertion - string to int
	t.Run("failed string to int assertion", func(t *testing.T) {
		var value any = "hello"
		result := InstanceOf[int](value)
		assert.True(t, IsLeft(result))
		_, err := UnwrapError(result)
		assert.Error(t, err)
	})

	// Test failed type assertion - int to float64
	t.Run("failed int to float64 assertion", func(t *testing.T) {
		var value any = 42
		result := InstanceOf[float64](value)
		assert.True(t, IsLeft(result))
	})

	// Test with nil value
	t.Run("nil value assertion", func(t *testing.T) {
		var value any = nil
		result := InstanceOf[string](value)
		assert.True(t, IsLeft(result))
	})

	// Test chaining with Map
	t.Run("chaining with Map", func(t *testing.T) {
		var value any = 10
		result := F.Pipe2(
			value,
			InstanceOf[int],
			Map(func(n int) int { return n * 2 }),
		)
		assert.Equal(t, Right(20), result)
	})

	// Test chaining with Map on failed assertion
	t.Run("chaining with Map on failed assertion", func(t *testing.T) {
		var value any = "not a number"
		result := F.Pipe2(
			value,
			InstanceOf[int],
			Map(func(n int) int { return n * 2 }),
		)
		assert.True(t, IsLeft(result))
	})

	// Test with Chain for dependent operations
	t.Run("chaining with Chain", func(t *testing.T) {
		var value any = 5
		result := F.Pipe2(
			value,
			InstanceOf[int],
			Chain(func(n int) Result[string] {
				if n > 0 {
					return Right(fmt.Sprintf("positive: %d", n))
				}
				return Left[string](errors.New("not positive"))
			}),
		)
		assert.Equal(t, Right("positive: 5"), result)
	})

	// Test with GetOrElse for default value
	t.Run("GetOrElse with failed assertion", func(t *testing.T) {
		var value any = "not an int"
		result := F.Pipe2(
			value,
			InstanceOf[int],
			GetOrElse(func(err error) int { return -1 }),
		)
		assert.Equal(t, -1, result)
	})

	// Test with GetOrElse for successful assertion
	t.Run("GetOrElse with successful assertion", func(t *testing.T) {
		var value any = 42
		result := F.Pipe2(
			value,
			InstanceOf[int],
			GetOrElse(func(err error) int { return -1 }),
		)
		assert.Equal(t, 42, result)
	})

	// Test with interface type
	t.Run("interface type assertion", func(t *testing.T) {
		var value any = errors.New("test error")
		result := InstanceOf[error](value)
		assert.True(t, IsRight(result))
		v, err := UnwrapError(result)
		assert.NoError(t, err)
		assert.Equal(t, "test error", v.Error())
	})

	// Test with slice type
	t.Run("slice type assertion", func(t *testing.T) {
		var value any = []int{1, 2, 3}
		result := InstanceOf[[]int](value)
		assert.True(t, IsRight(result))
		assert.Equal(t, Right([]int{1, 2, 3}), result)
	})

	// Test with map type
	t.Run("map type assertion", func(t *testing.T) {
		var value any = map[string]int{"a": 1, "b": 2}
		result := InstanceOf[map[string]int](value)
		assert.True(t, IsRight(result))
		v, err := UnwrapError(result)
		assert.NoError(t, err)
		assert.Equal(t, 1, v["a"])
		assert.Equal(t, 2, v["b"])
	})
}
