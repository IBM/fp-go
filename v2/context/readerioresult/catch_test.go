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

package readerioresult

import (
	"context"
	"errors"
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// notFoundError implements error with a pointer receiver.
type notFoundError struct{ id int }

func (e *notFoundError) Error() string { return fmt.Sprintf("%d not found", e.id) }

// codeError implements error with a value receiver.
type codeError struct{ code int }

func (e codeError) Error() string { return fmt.Sprintf("code %d", e.code) }

// temporary is an interface error type.
type temporary interface {
	error
	Temporary() bool
}

type tempError struct{}

func (tempError) Error() string   { return "temporary" }
func (tempError) Temporary() bool { return true }

var errSentinel = errors.New("sentinel")

func recoverNotFound(e *notFoundError) ReaderIOResult[int] { return Of(e.id) }

func recoverCode(e codeError) ReaderIOResult[int] { return Of(e.code) }

func TestCatchAs_Success(t *testing.T) {
	ctx := t.Context()

	t.Run("recovers from a matching error", func(t *testing.T) {
		res := CatchAs(recoverNotFound)(Left[int](&notFoundError{id: 7}))(ctx)()
		assert.Equal(t, result.Of(7), res)
	})

	t.Run("recovers from a wrapped error", func(t *testing.T) {
		err := fmt.Errorf("lookup: %w", &notFoundError{id: 8})
		res := CatchAs(recoverNotFound)(Left[int](err))(ctx)()
		assert.Equal(t, result.Of(8), res)
	})

	t.Run("recovers from a joined error", func(t *testing.T) {
		err := errors.Join(errSentinel, &notFoundError{id: 9})
		res := CatchAs(recoverNotFound)(Left[int](err))(ctx)()
		assert.Equal(t, result.Of(9), res)
	})

	t.Run("recovers from an error with a value receiver", func(t *testing.T) {
		res := CatchAs(recoverCode)(Left[int](codeError{code: 404}))(ctx)()
		assert.Equal(t, result.Of(404), res)
	})

	t.Run("recovers from an interface error type", func(t *testing.T) {
		recoverTemp := func(temporary) ReaderIOResult[int] { return Of(1) }
		res := CatchAs(recoverTemp)(Left[int](tempError{}))(ctx)()
		assert.Equal(t, result.Of(1), res)
	})

	t.Run("the handler can fail with a new error", func(t *testing.T) {
		replace := func(*notFoundError) ReaderIOResult[int] { return Left[int](errSentinel) }
		res := CatchAs(replace)(Left[int](&notFoundError{id: 7}))(ctx)()
		assert.Equal(t, result.Left[int](errSentinel), res)
	})

	t.Run("passes a success through without calling the handler", func(t *testing.T) {
		called := false
		handler := func(*notFoundError) ReaderIOResult[int] { called = true; return Of(0) }
		assert.Equal(t, result.Of(42), CatchAs(handler)(Of(42))(ctx)())
		assert.False(t, called)
	})
}

func TestCatchAs_Failure(t *testing.T) {
	ctx := t.Context()

	t.Run("passes an error of another type through", func(t *testing.T) {
		res := CatchAs(recoverNotFound)(Left[int](errSentinel))(ctx)()
		assert.Equal(t, result.Left[int](errSentinel), res)
	})

	t.Run("pointer type does not match a value receiver error", func(t *testing.T) {
		recoverPtr := func(e *codeError) ReaderIOResult[int] { return Of(e.code) }
		err := codeError{code: 500}
		assert.Equal(t, result.Left[int](error(err)), CatchAs(recoverPtr)(Left[int](err))(ctx)())
	})

	t.Run("value type does not match a pointer to a value receiver error", func(t *testing.T) {
		err := &codeError{code: 500}
		assert.Equal(t, result.Left[int](error(err)), CatchAs(recoverCode)(Left[int](err))(ctx)())
	})
}

func TestCatchAs_EdgeCases(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	t.Run("a non-matching error passes through a cancelled context unchanged", func(t *testing.T) {
		res := CatchAs(recoverNotFound)(Left[int](errSentinel))(ctx)()
		assert.Equal(t, result.Left[int](errSentinel), res)
	})

	t.Run("the handler is not run with a cancelled context", func(t *testing.T) {
		res := CatchAs(recoverNotFound)(Left[int](&notFoundError{id: 7}))(ctx)()
		assert.Equal(t, result.Left[int](context.Canceled), res)
	})
}

func TestCatchIs_Success(t *testing.T) {
	ctx := t.Context()
	fallback := func(error) ReaderIOResult[int] { return Of(0) }

	t.Run("recovers from the sentinel", func(t *testing.T) {
		assert.Equal(t, result.Of(0), CatchIs(errSentinel, fallback)(Left[int](errSentinel))(ctx)())
	})

	t.Run("recovers from a wrapped sentinel", func(t *testing.T) {
		err := fmt.Errorf("read: %w", errSentinel)
		assert.Equal(t, result.Of(0), CatchIs(errSentinel, fallback)(Left[int](err))(ctx)())
	})

	t.Run("recovers from a joined sentinel", func(t *testing.T) {
		err := errors.Join(&notFoundError{id: 1}, errSentinel)
		assert.Equal(t, result.Of(0), CatchIs(errSentinel, fallback)(Left[int](err))(ctx)())
	})

	t.Run("passes the original error to the handler", func(t *testing.T) {
		err := fmt.Errorf("read: %w", errSentinel)
		var got error
		handler := func(e error) ReaderIOResult[int] { got = e; return Of(0) }
		CatchIs(errSentinel, handler)(Left[int](err))(ctx)()
		assert.Same(t, err, got)
	})

	t.Run("passes a success through", func(t *testing.T) {
		assert.Equal(t, result.Of(42), CatchIs(errSentinel, fallback)(Of(42))(ctx)())
	})
}

func TestCatchIs_Failure(t *testing.T) {
	other := errors.New("other")
	res := CatchIs(errSentinel, func(error) ReaderIOResult[int] { return Of(0) })(Left[int](other))(t.Context())()
	assert.Equal(t, result.Left[int](other), res)
}

func TestMapErrorAs_Success(t *testing.T) {
	ctx := t.Context()
	toCode := func(e *notFoundError) error { return codeError{code: 404 + e.id} }

	t.Run("transforms a matching error", func(t *testing.T) {
		res := MapErrorAs[*notFoundError, int](toCode)(Left[int](&notFoundError{id: 0}))(ctx)()
		assert.Equal(t, result.Left[int](codeError{code: 404}), res)
	})

	t.Run("transforms a wrapped error", func(t *testing.T) {
		err := fmt.Errorf("lookup: %w", &notFoundError{id: 1})
		res := MapErrorAs[*notFoundError, int](toCode)(Left[int](err))(ctx)()
		assert.Equal(t, result.Left[int](codeError{code: 405}), res)
	})

	t.Run("passes a success through", func(t *testing.T) {
		assert.Equal(t, result.Of(42), MapErrorAs[*notFoundError, int](toCode)(Of(42))(ctx)())
	})
}

func TestMapErrorAs_Failure(t *testing.T) {
	toCode := func(*notFoundError) error { return codeError{code: 404} }
	res := MapErrorAs[*notFoundError, int](toCode)(Left[int](errSentinel))(t.Context())()
	assert.Equal(t, result.Left[int](errSentinel), res)
}

func TestTapErrorAs_Success(t *testing.T) {
	ctx := t.Context()

	t.Run("observes a matching error and keeps it", func(t *testing.T) {
		var seen int
		observe := func(e *notFoundError) ReaderIOResult[Void] {
			return FromIO(func() Void { seen = e.id; return F.VOID })
		}
		err := fmt.Errorf("lookup: %w", &notFoundError{id: 3})
		res := TapErrorAs[*notFoundError, int](observe)(Left[int](err))(ctx)()
		assert.Equal(t, result.Left[int](err), res)
		assert.Equal(t, 3, seen)
	})

	t.Run("a failing side effect does not replace the error", func(t *testing.T) {
		fail := func(*notFoundError) ReaderIOResult[Void] { return Left[Void](errSentinel) }
		err := &notFoundError{id: 3}
		res := TapErrorAs[*notFoundError, int](fail)(Left[int](err))(ctx)()
		assert.Equal(t, result.Left[int](error(err)), res)
	})
}

func TestTapErrorAs_Failure(t *testing.T) {
	ctx := t.Context()
	called := false
	observe := func(*notFoundError) ReaderIOResult[Void] {
		return FromIO(func() Void { called = true; return F.VOID })
	}

	t.Run("does not observe an error of another type", func(t *testing.T) {
		res := TapErrorAs[*notFoundError, int](observe)(Left[int](errSentinel))(ctx)()
		assert.Equal(t, result.Left[int](errSentinel), res)
		assert.False(t, called)
	})

	t.Run("does not observe a success", func(t *testing.T) {
		assert.Equal(t, result.Of(42), TapErrorAs[*notFoundError, int](observe)(Of(42))(ctx)())
		assert.False(t, called)
	})
}

func TestCatch_Integration(t *testing.T) {
	ctx := t.Context()
	lookup := func(id int) ReaderIOResult[int] {
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
		CatchAs(func(*notFoundError) ReaderIOResult[int] { return Of(-1) }),
		CatchIs(errSentinel, func(error) ReaderIOResult[int] { return Of(-2) }),
	)

	assert.Equal(t, result.Of(-1), handle(0)(ctx)())
	assert.Equal(t, result.Of(-2), handle(1)(ctx)())
	assert.Equal(t, result.Of(20), handle(2)(ctx)())
}
