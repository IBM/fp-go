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

	N "github.com/IBM/fp-go/v2/number"
	"github.com/stretchr/testify/assert"
)

// TestBiMap tests mapping over both error and value channels
func TestBiMap(t *testing.T) {
	wrapError := func(e error) error {
		return fmt.Errorf("wrapped: %w", e)
	}
	double := N.Mul(2)

	t.Run("BiMap on Right", func(t *testing.T) {
		val, err := BiMap(wrapError, double)(Right(21))
		AssertEq(Right(42))(val, err)(t)
	})

	t.Run("BiMap on Left", func(t *testing.T) {
		originalErr := errors.New("original")
		val, err := BiMap(wrapError, double)(Left[int](originalErr))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "wrapped")
		assert.Contains(t, err.Error(), "original")
		assert.Equal(t, 0, val)
	})
}

// TestMapTo tests mapping to a constant value
func TestMapTo(t *testing.T) {
	t.Run("MapTo on Right", func(t *testing.T) {
		val, err := MapTo[int]("constant")(Right(42))
		AssertEq(Right("constant"))(val, err)(t)
	})

	t.Run("MapTo on Left", func(t *testing.T) {
		originalErr := errors.New("error")
		val, err := MapTo[int]("constant")(Left[int](originalErr))
		assert.Error(t, err)
		assert.Equal(t, originalErr, err)
		// MapTo still applies the constant value even for Left
		assert.Equal(t, "constant", val)
	})
}

// TestChainTo tests chaining to a constant value
func TestChainTo(t *testing.T) {
	t.Run("ChainTo Right to Right", func(t *testing.T) {
		val, err := ChainTo[int]("success", nil)(Right(42))
		AssertEq(Right("success"))(val, err)(t)
	})

	t.Run("ChainTo Right to Left", func(t *testing.T) {
		targetErr := errors.New("target error")
		val, err := ChainTo[int]("", targetErr)(Right(42))
		assert.Error(t, err)
		assert.Equal(t, targetErr, err)
		assert.Equal(t, "", val)
	})

	t.Run("ChainTo Left", func(t *testing.T) {
		sourceErr := errors.New("source error")
		val, err := ChainTo[int]("success", nil)(Left[int](sourceErr))
		assert.Error(t, err)
		assert.Equal(t, sourceErr, err)
		assert.Equal(t, "", val)
	})
}

// TestReduce tests the reduce/fold operation
func TestReduce(t *testing.T) {
	sum := func(acc, val int) int {
		return acc + val
	}

	t.Run("Reduce on Right", func(t *testing.T) {
		result := Reduce(sum, 10)(Right(32))
		assert.Equal(t, 42, result)
	})

	t.Run("Reduce on Left", func(t *testing.T) {
		result := Reduce(sum, 10)(Left[int](errors.New("error")))
		assert.Equal(t, 10, result) // Returns initial value
	})
}

// TestFromPredicate tests creating Result from a predicate
func TestFromPredicate(t *testing.T) {
	isPositive := FromPredicate(
		N.MoreThan(0),
		func(x int) error { return fmt.Errorf("%d is not positive", x) },
	)

	t.Run("Predicate passes", func(t *testing.T) {
		val, err := isPositive(42)
		AssertEq(Right(42))(val, err)(t)
	})

	t.Run("Predicate fails", func(t *testing.T) {
		val, err := isPositive(-5)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not positive")
		// FromPredicate returns zero value on Left
		assert.Equal(t, 0, val)
	})

	t.Run("Predicate with zero", func(t *testing.T) {
		val, err := isPositive(0)
		assert.Error(t, err)
		assert.Equal(t, 0, val)
	})
}

// TestFromNillable tests creating Result from nullable pointers
func TestFromNillable(t *testing.T) {
	nilErr := errors.New("value is nil")
	fromPtr := FromNillable[int](nilErr)

	t.Run("Non-nil pointer", func(t *testing.T) {
		value := 42
		val, err := fromPtr(&value)
		assert.NoError(t, err)
		assert.NotNil(t, val)
		assert.Equal(t, 42, *val)
	})

	t.Run("Nil pointer", func(t *testing.T) {
		val, err := fromPtr(nil)
		assert.Error(t, err)
		assert.Equal(t, nilErr, err)
		assert.Nil(t, val)
	})
}

// TestToType tests type conversion
func TestToType(t *testing.T) {
	toInt := ToType[int](func(v any) error {
		return fmt.Errorf("cannot convert %T to int", v)
	})

	t.Run("Correct type", func(t *testing.T) {
		val, err := toInt(42)
		AssertEq(Right(42))(val, err)(t)
	})

	t.Run("Wrong type", func(t *testing.T) {
		val, err := toInt("string")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot convert")
		assert.Equal(t, 0, val)
	})

	t.Run("Nil value", func(t *testing.T) {
		val, err := toInt(nil)
		assert.Error(t, err)
		assert.Equal(t, 0, val)
	})
}

// TestMemoize tests that Memoize returns the value unchanged
func TestMemoize(t *testing.T) {
	t.Run("Memoize Right", func(t *testing.T) {
		val, err := Memoize(Right(42))
		AssertEq(Right(42))(val, err)(t)
	})

	t.Run("Memoize Left", func(t *testing.T) {
		originalErr := errors.New("error")
		val, err := Memoize(Left[int](originalErr))
		assert.Error(t, err)
		assert.Equal(t, originalErr, err)
		assert.Equal(t, 0, val)
	})
}

// TestFlap tests reverse application
func TestFlap(t *testing.T) {
	t.Run("Flap with Right function", func(t *testing.T) {
		double := N.Mul(2)
		val, err := Flap[int](21)(Right(double))
		AssertEq(Right(42))(val, err)(t)
	})

	t.Run("Flap with Left function", func(t *testing.T) {
		fnErr := errors.New("function error")
		val, err := Flap[int](21)(Left[func(int) int](fnErr))
		assert.Error(t, err)
		assert.Equal(t, fnErr, err)
		assert.Equal(t, 0, val)
	})
}

// TestToError tests extracting error from Result
func TestToError(t *testing.T) {
	t.Run("ToError from Right", func(t *testing.T) {
		err := ToError(Right(42))
		assert.NoError(t, err)
	})

	t.Run("ToError from Left", func(t *testing.T) {
		originalErr := errors.New("error")
		err := ToError(Left[int](originalErr))
		assert.Error(t, err)
		assert.Equal(t, originalErr, err)
	})
}

// TestLet tests the Let operation for do-notation
func TestLet(t *testing.T) {
	type State struct {
		value int
	}

	t.Run("Let with Right", func(t *testing.T) {
		val, err := Pipe2(
			State{value: 10},
			Right,
			Let(
				func(v int) func(State) State {
					return func(s State) State { return State{value: s.value + v} }
				},
				func(s State) int { return 32 },
			),
		)
		assert.NoError(t, err)
		assert.Equal(t, 42, val.value)
	})

	t.Run("Let with Left", func(t *testing.T) {
		originalErr := errors.New("error")
		val, err := Pipe2(
			originalErr,
			Left[State],
			Let(
				func(v int) func(State) State {
					return func(s State) State { return State{value: v} }
				},
				func(s State) int { return 42 },
			),
		)
		assert.Error(t, err)
		assert.Equal(t, originalErr, err)
		assert.Equal(t, State{}, val)
	})
}

// TestLetTo tests the LetTo operation
func TestLetTo(t *testing.T) {
	type State struct {
		name string
	}

	t.Run("LetTo with Right", func(t *testing.T) {
		val, err := Pipe2(
			State{},
			Right,
			LetTo(
				func(n string) func(State) State {
					return func(s State) State { return State{name: n} }
				},
				"Alice",
			),
		)
		assert.NoError(t, err)
		assert.Equal(t, "Alice", val.name)
	})

	t.Run("LetTo with Left", func(t *testing.T) {
		originalErr := errors.New("error")
		val, err := Pipe2(
			originalErr,
			Left[State],
			LetTo(
				func(n string) func(State) State {
					return func(s State) State { return State{name: n} }
				},
				"Bob",
			),
		)
		assert.Error(t, err)
		assert.Equal(t, State{}, val)
	})
}

// TestBindTo tests the BindTo operation
func TestBindTo(t *testing.T) {
	type State struct {
		value int
	}

	t.Run("BindTo with Right", func(t *testing.T) {
		val, err := Pipe2(
			42,
			Right,
			BindTo(func(v int) State { return State{value: v} }),
		)
		assert.NoError(t, err)
		assert.Equal(t, 42, val.value)
	})

	t.Run("BindTo with Left", func(t *testing.T) {
		originalErr := errors.New("error")
		val, err := Pipe2(
			originalErr,
			Left[int],
			BindTo(func(v int) State { return State{value: v} }),
		)
		assert.Error(t, err)
		assert.Equal(t, State{}, val)
	})
}

// TestMapLeft tests mapping over the error channel
func TestMapLeft(t *testing.T) {
	wrapError := func(e error) error {
		return fmt.Errorf("wrapped: %w", e)
	}

	t.Run("MapLeft on Right", func(t *testing.T) {
		val, err := MapLeft[int](wrapError)(Right(42))
		AssertEq(Right(42))(val, err)(t)
	})

	t.Run("MapLeft on Left", func(t *testing.T) {
		originalErr := errors.New("original")
		_, err := MapLeft[int](wrapError)(Left[int](originalErr))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "wrapped")
		assert.Contains(t, err.Error(), "original")
	})
}

// TestOrElse tests recovery from error
func TestOrElse(t *testing.T) {
	recover := OrElse(func(e error) (int, error) {
		return Right(0) // default value
	})

	t.Run("OrElse on Right", func(t *testing.T) {
		val, err := recover(Right(42))
		AssertEq(Right(42))(val, err)(t)
	})

	t.Run("OrElse on Left recovers", func(t *testing.T) {
		val, err := recover(Left[int](errors.New("error")))
		AssertEq(Right(0))(val, err)(t)
	})
}

// TestDo tests the Do operation
func TestDo(t *testing.T) {
	type State struct {
		x int
		y int
	}

	result, err := Do(State{})
	assert.NoError(t, err)
	assert.Equal(t, State{}, result)
}

// TestOf tests the Of/pure operation
func TestOf(t *testing.T) {
	val, err := Of(42)
	AssertEq(Right(42))(val, err)(t)
}

// TestToString tests string representation
func TestToString(t *testing.T) {
	t.Run("ToString Right", func(t *testing.T) {
		str := ToString(Right(42))
		assert.Equal(t, "Right[int](42)", str)
	})

	t.Run("ToString Left", func(t *testing.T) {
		str := ToString(Left[int](errors.New("error")))
		assert.Contains(t, str, "Left(")
		assert.Contains(t, str, "error")
	})
}

// TestToOption tests conversion to Option
func TestToOption(t *testing.T) {
	t.Run("ToOption from Right", func(t *testing.T) {
		val, ok := ToOption(Right(42))
		assert.True(t, ok)
		assert.Equal(t, 42, val)
	})

	t.Run("ToOption from Left", func(t *testing.T) {
		val, ok := ToOption(Left[int](errors.New("error")))
		assert.False(t, ok)
		assert.Equal(t, 0, val)
	})
}

// TestFromError tests creating Result from error-returning function
func TestFromError(t *testing.T) {
	validate := func(x int) error {
		if x < 0 {
			return errors.New("negative")
		}
		return nil
	}

	toResult := FromError(validate)

	t.Run("FromError with valid value", func(t *testing.T) {
		val, err := toResult(42)
		AssertEq(Right(42))(val, err)(t)
	})

	t.Run("FromError with invalid value", func(t *testing.T) {
		val, err := toResult(-5)
		assert.Error(t, err)
		assert.Equal(t, "negative", err.Error())
		assert.Equal(t, -5, val)
	})
}

// TestGetOrElse tests extracting value with default
func TestGetOrElse(t *testing.T) {
	defaultValue := func(error) int { return 0 }

	t.Run("GetOrElse on Right", func(t *testing.T) {
		val := GetOrElse(defaultValue)(Right(42))
		assert.Equal(t, 42, val)
	})

	t.Run("GetOrElse on Left", func(t *testing.T) {
		val := GetOrElse(defaultValue)(Left[int](errors.New("error")))
		assert.Equal(t, 0, val)
	})
}
