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
	"github.com/stretchr/testify/assert"
)

// notFoundError implements error with a pointer receiver.
type notFoundError struct{ id int }

func (e *notFoundError) Error() string { return fmt.Sprintf("%d not found", e.id) }

// temporary is an interface error type.
type temporary interface {
	error
	Temporary() bool
}

type tempError struct{}

func (tempError) Error() string   { return "temporary" }
func (tempError) Temporary() bool { return true }

func catchNotFound(e *notFoundError) Result[int] { return Of(e.id) }

func recoverCode(e codeError) Result[int] { return Of(e.code) }

func TestCatchAs_Success(t *testing.T) {
	t.Run("recovers from a matching error", func(t *testing.T) {
		res := CatchAs(catchNotFound)(Left[int](&notFoundError{id: 7}))
		assert.Equal(t, Of(7), res)
	})

	t.Run("recovers from a wrapped error", func(t *testing.T) {
		err := fmt.Errorf("lookup: %w", &notFoundError{id: 8})
		res := CatchAs(catchNotFound)(Left[int](err))
		assert.Equal(t, Of(8), res)
	})

	t.Run("recovers from a joined error", func(t *testing.T) {
		err := errors.Join(errSentinel, &notFoundError{id: 9})
		res := CatchAs(catchNotFound)(Left[int](err))
		assert.Equal(t, Of(9), res)
	})

	t.Run("recovers from an error with a value receiver", func(t *testing.T) {
		res := CatchAs(recoverCode)(Left[int](codeError{code: 404}))
		assert.Equal(t, Of(404), res)
	})

	t.Run("recovers from an interface error type", func(t *testing.T) {
		recoverTemp := func(temporary) Result[int] { return Of(1) }
		res := CatchAs(recoverTemp)(Left[int](tempError{}))
		assert.Equal(t, Of(1), res)
	})

	t.Run("the handler can fail with a new error", func(t *testing.T) {
		replace := func(*notFoundError) Result[int] { return Left[int](errSentinel) }
		res := CatchAs(replace)(Left[int](&notFoundError{id: 7}))
		assert.Equal(t, Left[int](errSentinel), res)
	})

	t.Run("passes a success through without calling the handler", func(t *testing.T) {
		called := false
		handler := func(*notFoundError) Result[int] { called = true; return Of(0) }
		assert.Equal(t, Of(42), CatchAs(handler)(Of(42)))
		assert.False(t, called)
	})
}

func TestCatchAs_Failure(t *testing.T) {
	t.Run("passes an error of another type through", func(t *testing.T) {
		res := CatchAs(catchNotFound)(Left[int](errSentinel))
		assert.Equal(t, Left[int](errSentinel), res)
	})

	t.Run("pointer type does not match a value receiver error", func(t *testing.T) {
		recoverPtr := func(e *codeError) Result[int] { return Of(e.code) }
		err := codeError{code: 500}
		assert.Equal(t, Left[int](error(err)), CatchAs(recoverPtr)(Left[int](err)))
	})

	t.Run("value type does not match a pointer to a value receiver error", func(t *testing.T) {
		err := &codeError{code: 500}
		assert.Equal(t, Left[int](error(err)), CatchAs(recoverCode)(Left[int](err)))
	})
}

func TestCatchIs_Success(t *testing.T) {
	fallback := func(error) Result[int] { return Of(0) }

	t.Run("recovers from the sentinel", func(t *testing.T) {
		assert.Equal(t, Of(0), CatchIs(errSentinel, fallback)(Left[int](errSentinel)))
	})

	t.Run("recovers from a wrapped sentinel", func(t *testing.T) {
		err := fmt.Errorf("read: %w", errSentinel)
		assert.Equal(t, Of(0), CatchIs(errSentinel, fallback)(Left[int](err)))
	})

	t.Run("recovers from a joined sentinel", func(t *testing.T) {
		err := errors.Join(&notFoundError{id: 1}, errSentinel)
		assert.Equal(t, Of(0), CatchIs(errSentinel, fallback)(Left[int](err)))
	})

	t.Run("passes the original error to the handler", func(t *testing.T) {
		err := fmt.Errorf("read: %w", errSentinel)
		var got error
		handler := func(e error) Result[int] { got = e; return Of(0) }
		CatchIs(errSentinel, handler)(Left[int](err))
		assert.Same(t, err, got)
	})

	t.Run("passes a success through", func(t *testing.T) {
		assert.Equal(t, Of(42), CatchIs(errSentinel, fallback)(Of(42)))
	})
}

func TestCatchIs_Failure(t *testing.T) {
	other := errors.New("other")
	res := CatchIs(errSentinel, func(error) Result[int] { return Of(0) })(Left[int](other))
	assert.Equal(t, Left[int](other), res)
}

func TestMapErrorAs_Success(t *testing.T) {
	toCode := func(e *notFoundError) error { return codeError{code: 404 + e.id} }

	t.Run("transforms a matching error", func(t *testing.T) {
		res := MapErrorAs[*notFoundError, int](toCode)(Left[int](&notFoundError{id: 0}))
		assert.Equal(t, Left[int](codeError{code: 404}), res)
	})

	t.Run("transforms a wrapped error", func(t *testing.T) {
		err := fmt.Errorf("lookup: %w", &notFoundError{id: 1})
		res := MapErrorAs[*notFoundError, int](toCode)(Left[int](err))
		assert.Equal(t, Left[int](codeError{code: 405}), res)
	})

	t.Run("passes a success through", func(t *testing.T) {
		assert.Equal(t, Of(42), MapErrorAs[*notFoundError, int](toCode)(Of(42)))
	})
}

func TestMapErrorAs_Failure(t *testing.T) {
	toCode := func(*notFoundError) error { return codeError{code: 404} }
	res := MapErrorAs[*notFoundError, int](toCode)(Left[int](errSentinel))
	assert.Equal(t, Left[int](errSentinel), res)
}

func TestCatch_Integration(t *testing.T) {
	lookup := func(id int) Result[int] {
		switch id {
		case 0:
			return Left[int](&notFoundError{id: id})
		case 1:
			return Left[int](fmt.Errorf("db: %w", errSentinel))
		default:
			return Of(id * 10)
		}
	}
	handle := F.Flow3(
		lookup,
		CatchAs(func(*notFoundError) Result[int] { return Of(-1) }),
		CatchIs(errSentinel, func(error) Result[int] { return Of(-2) }),
	)

	assert.Equal(t, Of(-1), handle(0))
	assert.Equal(t, Of(-2), handle(1))
	assert.Equal(t, Of(20), handle(2))
}

func TestTapErrorAs_Success(t *testing.T) {
	t.Run("observes a matching error and keeps it", func(t *testing.T) {
		var seen int
		observe := func(e *notFoundError) Result[F.Void] { seen = e.id; return Of(F.VOID) }
		err := fmt.Errorf("lookup: %w", &notFoundError{id: 3})
		assert.Equal(t, Left[int](err), TapErrorAs[*notFoundError, int](observe)(Left[int](err)))
		assert.Equal(t, 3, seen)
	})

	t.Run("a failing side effect does not replace the error", func(t *testing.T) {
		fail := func(*notFoundError) Result[F.Void] { return Left[F.Void](errSentinel) }
		err := &notFoundError{id: 3}
		assert.Equal(t, Left[int](error(err)), TapErrorAs[*notFoundError, int](fail)(Left[int](err)))
	})
}

func TestTapErrorAs_Failure(t *testing.T) {
	called := false
	observe := func(*notFoundError) Result[F.Void] { called = true; return Of(F.VOID) }

	t.Run("does not observe an error of another type", func(t *testing.T) {
		assert.Equal(t, Left[int](errSentinel), TapErrorAs[*notFoundError, int](observe)(Left[int](errSentinel)))
		assert.False(t, called)
	})

	t.Run("does not observe a success", func(t *testing.T) {
		assert.Equal(t, Of(42), TapErrorAs[*notFoundError, int](observe)(Of(42)))
		assert.False(t, called)
	})
}
