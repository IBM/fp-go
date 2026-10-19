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

// config is the environment of the computations in the tests.
type config struct{}

func run[A any](rio ReaderIOResult[config, A]) result.Result[A] {
	return result.TryCatchError(rio(config{})())
}

func recoverNotFound(e *notFoundError) ReaderIOResult[config, int] { return Of[config](e.id) }

func recoverCode(e codeError) ReaderIOResult[config, int] { return Of[config](e.code) }

func TestCatchAs_Success(t *testing.T) {
	t.Run("recovers from a matching error", func(t *testing.T) {
		assert.Equal(t, (result.Of(7)), run(CatchAs(recoverNotFound)(Left[config, int](&notFoundError{id: 7}))))
	})

	t.Run("recovers from a wrapped error", func(t *testing.T) {
		err := fmt.Errorf("lookup: %w", &notFoundError{id: 8})
		assert.Equal(t, (result.Of(8)), run(CatchAs(recoverNotFound)(Left[config, int](err))))
	})

	t.Run("recovers from a joined error", func(t *testing.T) {
		err := errors.Join(errSentinel, &notFoundError{id: 9})
		assert.Equal(t, (result.Of(9)), run(CatchAs(recoverNotFound)(Left[config, int](err))))
	})

	t.Run("recovers from an error with a value receiver", func(t *testing.T) {
		assert.Equal(t, (result.Of(404)), run(CatchAs(recoverCode)(Left[config, int](codeError{code: 404}))))
	})

	t.Run("recovers from an interface error type", func(t *testing.T) {
		recoverTemp := func(temporary) ReaderIOResult[config, int] { return Of[config](1) }
		assert.Equal(t, (result.Of(1)), run(CatchAs(recoverTemp)(Left[config, int](tempError{}))))
	})

	t.Run("the handler can fail with a new error", func(t *testing.T) {
		replace := func(*notFoundError) ReaderIOResult[config, int] { return Left[config, int](errSentinel) }
		assert.Equal(t, (result.Left[int](errSentinel)), run(CatchAs(replace)(Left[config, int](&notFoundError{id: 7}))))
	})

	t.Run("passes a success through without calling the handler", func(t *testing.T) {
		called := false
		handler := func(*notFoundError) ReaderIOResult[config, int] { called = true; return Of[config](0) }
		assert.Equal(t, (result.Of(42)), run(CatchAs(handler)(Of[config](42))))
		assert.False(t, called)
	})
}

func TestCatchAs_Failure(t *testing.T) {
	t.Run("passes an error of another type through", func(t *testing.T) {
		assert.Equal(t, (result.Left[int](errSentinel)), run(CatchAs(recoverNotFound)(Left[config, int](errSentinel))))
	})

	t.Run("pointer type does not match a value receiver error", func(t *testing.T) {
		recoverPtr := func(e *codeError) ReaderIOResult[config, int] { return Of[config](e.code) }
		err := codeError{code: 500}
		assert.Equal(t, (result.Left[int](error(err))), run(CatchAs(recoverPtr)(Left[config, int](err))))
	})

	t.Run("value type does not match a pointer to a value receiver error", func(t *testing.T) {
		err := &codeError{code: 500}
		assert.Equal(t, (result.Left[int](error(err))), run(CatchAs(recoverCode)(Left[config, int](err))))
	})
}

func TestCatchIs_Success(t *testing.T) {
	fallback := func(error) ReaderIOResult[config, int] { return Of[config](0) }

	t.Run("recovers from the sentinel", func(t *testing.T) {
		assert.Equal(t, (result.Of(0)), run(CatchIs(errSentinel, fallback)(Left[config, int](errSentinel))))
	})

	t.Run("recovers from a wrapped sentinel", func(t *testing.T) {
		err := fmt.Errorf("read: %w", errSentinel)
		assert.Equal(t, (result.Of(0)), run(CatchIs(errSentinel, fallback)(Left[config, int](err))))
	})

	t.Run("recovers from a joined sentinel", func(t *testing.T) {
		err := errors.Join(&notFoundError{id: 1}, errSentinel)
		assert.Equal(t, (result.Of(0)), run(CatchIs(errSentinel, fallback)(Left[config, int](err))))
	})

	t.Run("passes the original error to the handler", func(t *testing.T) {
		err := fmt.Errorf("read: %w", errSentinel)
		var got error
		handler := func(e error) ReaderIOResult[config, int] { got = e; return Of[config](0) }
		run(CatchIs(errSentinel, handler)(Left[config, int](err)))
		assert.Same(t, err, got)
	})

	t.Run("passes a success through", func(t *testing.T) {
		assert.Equal(t, (result.Of(42)), run(CatchIs(errSentinel, fallback)(Of[config](42))))
	})
}

func TestCatchIs_Failure(t *testing.T) {
	other := errors.New("other")
	fallback := func(error) ReaderIOResult[config, int] { return Of[config](0) }
	assert.Equal(t, (result.Left[int](other)), run(CatchIs(errSentinel, fallback)(Left[config, int](other))))
}

func TestMapErrorAs_Success(t *testing.T) {
	toCode := func(e *notFoundError) error { return codeError{code: 404 + e.id} }

	t.Run("transforms a matching error", func(t *testing.T) {
		res := run(MapErrorAs[*notFoundError, config, int](toCode)(Left[config, int](&notFoundError{id: 0})))
		assert.Equal(t, (result.Left[int](codeError{code: 404})), res)
	})

	t.Run("transforms a wrapped error", func(t *testing.T) {
		err := fmt.Errorf("lookup: %w", &notFoundError{id: 1})
		res := run(MapErrorAs[*notFoundError, config, int](toCode)(Left[config, int](err)))
		assert.Equal(t, (result.Left[int](codeError{code: 405})), res)
	})

	t.Run("passes a success through", func(t *testing.T) {
		assert.Equal(t, (result.Of(42)), run(MapErrorAs[*notFoundError, config, int](toCode)(Of[config](42))))
	})
}

func TestMapErrorAs_Failure(t *testing.T) {
	toCode := func(*notFoundError) error { return codeError{code: 404} }
	res := run(MapErrorAs[*notFoundError, config, int](toCode)(Left[config, int](errSentinel)))
	assert.Equal(t, (result.Left[int](errSentinel)), res)
}

func TestCatch_Integration(t *testing.T) {
	lookup := func(id int) ReaderIOResult[config, int] {
		switch id {
		case 0:
			return Left[config, int](&notFoundError{id: id})
		case 1:
			return Left[config, int](fmt.Errorf("db: %w", errSentinel))
		default:
			return Of[config](id * 10)
		}
	}
	handle := func(id int) ReaderIOResult[config, int] {
		return CatchIs(errSentinel, func(error) ReaderIOResult[config, int] { return Of[config](-2) })(
			CatchAs(func(*notFoundError) ReaderIOResult[config, int] { return Of[config](-1) })(lookup(id)))
	}

	assert.Equal(t, (result.Of(-1)), run(handle(0)))
	assert.Equal(t, (result.Of(-2)), run(handle(1)))
	assert.Equal(t, (result.Of(20)), run(handle(2)))
}

func TestTapErrorAs_Success(t *testing.T) {
	t.Run("observes a matching error and keeps it", func(t *testing.T) {
		var seen int
		observe := func(e *notFoundError) ReaderIOResult[config, F.Void] {
			seen = e.id
			return Of[config](F.VOID)
		}
		err := fmt.Errorf("lookup: %w", &notFoundError{id: 3})
		assert.Equal(t, (result.Left[int](err)), run(TapErrorAs[*notFoundError, config, int](observe)(Left[config, int](err))))
		assert.Equal(t, 3, seen)
	})

	t.Run("a failing side effect does not replace the error", func(t *testing.T) {
		fail := func(*notFoundError) ReaderIOResult[config, F.Void] { return Left[config, F.Void](errSentinel) }
		err := &notFoundError{id: 3}
		assert.Equal(t, (result.Left[int](error(err))), run(TapErrorAs[*notFoundError, config, int](fail)(Left[config, int](err))))
	})
}

func TestTapErrorAs_Failure(t *testing.T) {
	called := false
	observe := func(*notFoundError) ReaderIOResult[config, F.Void] {
		called = true
		return Of[config](F.VOID)
	}

	t.Run("does not observe an error of another type", func(t *testing.T) {
		assert.Equal(t, (result.Left[int](errSentinel)), run(TapErrorAs[*notFoundError, config, int](observe)(Left[config, int](errSentinel))))
		assert.False(t, called)
	})

	t.Run("does not observe a success", func(t *testing.T) {
		assert.Equal(t, (result.Of(42)), run(TapErrorAs[*notFoundError, config, int](observe)(Of[config](42))))
		assert.False(t, called)
	})
}
