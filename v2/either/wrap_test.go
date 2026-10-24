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

package either

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/stretchr/testify/assert"
)

func TestWrap_Success(t *testing.T) {
	t.Run("wraps a value with nil error into Right", func(t *testing.T) {
		assert.Equal(t, Right[error](42), Wrap(42, nil))
	})

	t.Run("wraps a zero value with nil error into Right", func(t *testing.T) {
		assert.Equal(t, Right[error](0), Wrap(0, nil))
		assert.Equal(t, Right[error](""), Wrap("", nil))
	})

	t.Run("wraps a nil pointer with nil error into Right", func(t *testing.T) {
		result := Wrap[*int](nil, nil)
		assert.True(t, IsRight(result))
	})
}

func TestWrap_Failure(t *testing.T) {
	t.Run("returns Left when err is non-nil", func(t *testing.T) {
		err := errors.New("boom")
		assert.Equal(t, Left[int](err), Wrap(0, err))
	})

	t.Run("discards the value when err is non-nil", func(t *testing.T) {
		err := errors.New("boom")
		assert.Equal(t, Left[int](err), Wrap(42, err))
		assert.Equal(t, Left[string](err), Wrap("ignored", err))
	})

	t.Run("preserves the error identity and chain", func(t *testing.T) {
		cause := errors.New("cause")
		wrapped := fmt.Errorf("context: %w", cause)
		_, err := UnwrapError(Wrap(0, wrapped))
		assert.Same(t, wrapped, err)
		assert.ErrorIs(t, err, cause)
	})
}

func TestWrap_EdgeCases(t *testing.T) {
	t.Run("works with struct types", func(t *testing.T) {
		type point struct{ X, Y int }
		err := errors.New("boom")
		assert.Equal(t, Right[error](point{1, 2}), Wrap(point{1, 2}, nil))
		assert.Equal(t, Left[point](err), Wrap(point{1, 2}, err))
	})

	t.Run("works with a typed error", func(t *testing.T) {
		var pathErr error = &os.PathError{Op: "open", Path: "x", Err: os.ErrNotExist}
		result := Wrap(0, pathErr)
		assert.True(t, IsLeft(result))
		_, err := UnwrapError(result)
		assert.ErrorIs(t, err, os.ErrNotExist)
	})
}

func TestWrap_Integration(t *testing.T) {
	t.Run("is an alias for TryCatchError", func(t *testing.T) {
		err := errors.New("boom")
		assert.Equal(t, TryCatchError(42, nil), Wrap(42, nil))
		assert.Equal(t, TryCatchError(42, err), Wrap(42, err))
	})

	t.Run("is the inverse of UnwrapError", func(t *testing.T) {
		err := errors.New("boom")
		assert.Equal(t, Right[error](42), Wrap(UnwrapError(Right[error](42))))
		assert.Equal(t, Left[int](err), Wrap(UnwrapError(Left[int](err))))
	})

	t.Run("UnwrapError is the inverse of Wrap", func(t *testing.T) {
		val, err := UnwrapError(Wrap(42, nil))
		assert.NoError(t, err)
		assert.Equal(t, 42, val)

		boom := errors.New("boom")
		val, err = UnwrapError(Wrap(42, boom))
		assert.Equal(t, boom, err)
		assert.Equal(t, 0, val)
	})

	t.Run("accepts a (value, error) call directly", func(t *testing.T) {
		assert.Equal(t, Right[error](42), Wrap(strconv.Atoi("42")))
		assert.True(t, IsLeft(Wrap(strconv.Atoi("abc"))))
	})

	t.Run("agrees with Eitherize1", func(t *testing.T) {
		for _, s := range []string{"42", "abc"} {
			assert.Equal(t, Eitherize1(strconv.Atoi)(s), Wrap(strconv.Atoi(s)))
		}
	})

	t.Run("composes with Map", func(t *testing.T) {
		result := F.Pipe1(Wrap(strconv.Atoi("21")), Map[error](N.Mul(2)))
		assert.Equal(t, Right[error](42), result)
	})
}

func BenchmarkWrap(b *testing.B) {
	err := errors.New("boom")
	for i := range b.N {
		if i%2 == 0 {
			_ = Wrap(i, nil)
		} else {
			_ = Wrap(i, err)
		}
	}
}

// ExampleWrap demonstrates creating a value from an idiomatic (value, error) tuple.
func ExampleWrap() {
	fmt.Println(Wrap(42, nil))
	fmt.Println(Wrap(0, errors.New("boom")))
	// Output:
	// Right[int](42)
	// Left[*errors.errorString](boom)
}

// ExampleWrap_call demonstrates applying Wrap directly to a function that
// returns (A, error).
func ExampleWrap_call() {
	fmt.Println(Wrap(strconv.Atoi("42")))
	fmt.Println(Wrap(strconv.Atoi("abc")))
	// Output:
	// Right[int](42)
	// Left[*strconv.NumError](strconv.Atoi: parsing "abc": invalid syntax)
}

// ExampleWrap_unwrap demonstrates that Wrap is the inverse of UnwrapError.
func ExampleWrap_unwrap() {
	fmt.Println(Wrap(UnwrapError(Right[error]("value"))))
	fmt.Println(Wrap(UnwrapError(Left[string](errors.New("boom")))))
	// Output:
	// Right[string](value)
	// Left[*errors.errorString](boom)
}
