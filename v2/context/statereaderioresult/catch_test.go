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

package statereaderioresult

import (
	"context"
	"errors"
	"fmt"
	"testing"

	RIORES "github.com/IBM/fp-go/v2/context/readerioresult"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/pair"
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

// counter is the state of the computations in the tests.
type counter int

// run returns the value of a computation started from the zero state.
func run[A any](m StateReaderIOResult[counter, A]) Result[A] {
	return result.Map(pair.Tail[counter, A])(m(counter(0))(context.Background())())
}

func recoverNotFound(e *notFoundError) StateReaderIOResult[counter, int] { return Of[counter](e.id) }

func recoverCode(e codeError) StateReaderIOResult[counter, int] { return Of[counter](e.code) }

func TestCatchAs_Success(t *testing.T) {
	t.Run("recovers from a matching error", func(t *testing.T) {
		assert.Equal(t, (result.Of(7)), run(CatchAs(recoverNotFound)(Left[counter, int](&notFoundError{id: 7}))))
	})

	t.Run("recovers from a wrapped error", func(t *testing.T) {
		err := fmt.Errorf("lookup: %w", &notFoundError{id: 8})
		assert.Equal(t, (result.Of(8)), run(CatchAs(recoverNotFound)(Left[counter, int](err))))
	})

	t.Run("recovers from a joined error", func(t *testing.T) {
		err := errors.Join(errSentinel, &notFoundError{id: 9})
		assert.Equal(t, (result.Of(9)), run(CatchAs(recoverNotFound)(Left[counter, int](err))))
	})

	t.Run("recovers from an error with a value receiver", func(t *testing.T) {
		assert.Equal(t, (result.Of(404)), run(CatchAs(recoverCode)(Left[counter, int](codeError{code: 404}))))
	})

	t.Run("recovers from an interface error type", func(t *testing.T) {
		recoverTemp := func(temporary) StateReaderIOResult[counter, int] { return Of[counter](1) }
		assert.Equal(t, (result.Of(1)), run(CatchAs(recoverTemp)(Left[counter, int](tempError{}))))
	})

	t.Run("the handler can fail with a new error", func(t *testing.T) {
		replace := func(*notFoundError) StateReaderIOResult[counter, int] { return Left[counter, int](errSentinel) }
		assert.Equal(t, (result.Left[int](errSentinel)), run(CatchAs(replace)(Left[counter, int](&notFoundError{id: 7}))))
	})

	t.Run("passes a success through without calling the handler", func(t *testing.T) {
		called := false
		handler := func(*notFoundError) StateReaderIOResult[counter, int] { called = true; return Of[counter](0) }
		assert.Equal(t, (result.Of(42)), run(CatchAs(handler)(Of[counter](42))))
		assert.False(t, called)
	})
}

func TestCatchAs_Failure(t *testing.T) {
	t.Run("passes an error of another type through", func(t *testing.T) {
		assert.Equal(t, (result.Left[int](errSentinel)), run(CatchAs(recoverNotFound)(Left[counter, int](errSentinel))))
	})

	t.Run("pointer type does not match a value receiver error", func(t *testing.T) {
		recoverPtr := func(e *codeError) StateReaderIOResult[counter, int] { return Of[counter](e.code) }
		err := codeError{code: 500}
		assert.Equal(t, (result.Left[int](error(err))), run(CatchAs(recoverPtr)(Left[counter, int](err))))
	})

	t.Run("value type does not match a pointer to a value receiver error", func(t *testing.T) {
		err := &codeError{code: 500}
		assert.Equal(t, (result.Left[int](error(err))), run(CatchAs(recoverCode)(Left[counter, int](err))))
	})
}

func TestCatchIs_Success(t *testing.T) {
	fallback := func(error) StateReaderIOResult[counter, int] { return Of[counter](0) }

	t.Run("recovers from the sentinel", func(t *testing.T) {
		assert.Equal(t, (result.Of(0)), run(CatchIs(errSentinel, fallback)(Left[counter, int](errSentinel))))
	})

	t.Run("recovers from a wrapped sentinel", func(t *testing.T) {
		err := fmt.Errorf("read: %w", errSentinel)
		assert.Equal(t, (result.Of(0)), run(CatchIs(errSentinel, fallback)(Left[counter, int](err))))
	})

	t.Run("recovers from a joined sentinel", func(t *testing.T) {
		err := errors.Join(&notFoundError{id: 1}, errSentinel)
		assert.Equal(t, (result.Of(0)), run(CatchIs(errSentinel, fallback)(Left[counter, int](err))))
	})

	t.Run("passes the original error to the handler", func(t *testing.T) {
		err := fmt.Errorf("read: %w", errSentinel)
		var got error
		handler := func(e error) StateReaderIOResult[counter, int] { got = e; return Of[counter](0) }
		run(CatchIs(errSentinel, handler)(Left[counter, int](err)))
		assert.Same(t, err, got)
	})

	t.Run("passes a success through", func(t *testing.T) {
		assert.Equal(t, (result.Of(42)), run(CatchIs(errSentinel, fallback)(Of[counter](42))))
	})
}

func TestCatchIs_Failure(t *testing.T) {
	other := errors.New("other")
	fallback := func(error) StateReaderIOResult[counter, int] { return Of[counter](0) }
	assert.Equal(t, (result.Left[int](other)), run(CatchIs(errSentinel, fallback)(Left[counter, int](other))))
}

func TestMapErrorAs_Success(t *testing.T) {
	toCode := func(e *notFoundError) error { return codeError{code: 404 + e.id} }

	t.Run("transforms a matching error", func(t *testing.T) {
		res := run(MapErrorAs[*notFoundError, counter, int](toCode)(Left[counter, int](&notFoundError{id: 0})))
		assert.Equal(t, (result.Left[int](codeError{code: 404})), res)
	})

	t.Run("transforms a wrapped error", func(t *testing.T) {
		err := fmt.Errorf("lookup: %w", &notFoundError{id: 1})
		res := run(MapErrorAs[*notFoundError, counter, int](toCode)(Left[counter, int](err)))
		assert.Equal(t, (result.Left[int](codeError{code: 405})), res)
	})

	t.Run("passes a success through", func(t *testing.T) {
		assert.Equal(t, (result.Of(42)), run(MapErrorAs[*notFoundError, counter, int](toCode)(Of[counter](42))))
	})
}

func TestMapErrorAs_Failure(t *testing.T) {
	toCode := func(*notFoundError) error { return codeError{code: 404} }
	res := run(MapErrorAs[*notFoundError, counter, int](toCode)(Left[counter, int](errSentinel)))
	assert.Equal(t, (result.Left[int](errSentinel)), res)
}

func TestCatch_Integration(t *testing.T) {
	lookup := func(id int) StateReaderIOResult[counter, int] {
		switch id {
		case 0:
			return Left[counter, int](&notFoundError{id: id})
		case 1:
			return Left[counter, int](fmt.Errorf("db: %w", errSentinel))
		default:
			return Of[counter](id * 10)
		}
	}
	handle := func(id int) StateReaderIOResult[counter, int] {
		return CatchIs(errSentinel, func(error) StateReaderIOResult[counter, int] { return Of[counter](-2) })(
			CatchAs(func(*notFoundError) StateReaderIOResult[counter, int] { return Of[counter](-1) })(lookup(id)))
	}

	assert.Equal(t, (result.Of(-1)), run(handle(0)))
	assert.Equal(t, (result.Of(-2)), run(handle(1)))
	assert.Equal(t, (result.Of(20)), run(handle(2)))
}

func TestTapErrorAs_Success(t *testing.T) {
	t.Run("observes a matching error and keeps it", func(t *testing.T) {
		var seen int
		observe := func(e *notFoundError) StateReaderIOResult[counter, F.Void] {
			seen = e.id
			return Of[counter](F.VOID)
		}
		err := fmt.Errorf("lookup: %w", &notFoundError{id: 3})
		assert.Equal(t, (result.Left[int](err)), run(TapErrorAs[*notFoundError, counter, int](observe)(Left[counter, int](err))))
		assert.Equal(t, 3, seen)
	})

	t.Run("a failing side effect does not replace the error", func(t *testing.T) {
		fail := func(*notFoundError) StateReaderIOResult[counter, F.Void] { return Left[counter, F.Void](errSentinel) }
		err := &notFoundError{id: 3}
		assert.Equal(t, (result.Left[int](error(err))), run(TapErrorAs[*notFoundError, counter, int](fail)(Left[counter, int](err))))
	})
}

func TestTapErrorAs_Failure(t *testing.T) {
	called := false
	observe := func(*notFoundError) StateReaderIOResult[counter, F.Void] {
		called = true
		return Of[counter](F.VOID)
	}

	t.Run("does not observe an error of another type", func(t *testing.T) {
		assert.Equal(t, (result.Left[int](errSentinel)), run(TapErrorAs[*notFoundError, counter, int](observe)(Left[counter, int](errSentinel))))
		assert.False(t, called)
	})

	t.Run("does not observe a success", func(t *testing.T) {
		assert.Equal(t, (result.Of(42)), run(TapErrorAs[*notFoundError, counter, int](observe)(Of[counter](42))))
		assert.False(t, called)
	})
}

// countFrom recovers with the current state as the value and increments the state.
func countFrom(*notFoundError) StateReaderIOResult[counter, int] {
	return func(s counter) ReaderIOResult[Pair[counter, int]] {
		return RIORES.Of(pair.MakePair(s+1, int(s)))
	}
}

func TestCatchAs_EdgeCases(t *testing.T) {
	t.Run("the handler starts from the state of the failed computation", func(t *testing.T) {
		res := CatchAs(countFrom)(Left[counter, int](&notFoundError{id: 1}))(counter(5))(t.Context())()
		assert.Equal(t, result.Of(pair.MakePair(counter(6), 5)), res)
	})

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	t.Run("a non-matching error passes through a cancelled context unchanged", func(t *testing.T) {
		res := CatchAs(countFrom)(Left[counter, int](errSentinel))(counter(0))(ctx)()
		assert.Equal(t, result.Left[Pair[counter, int]](errSentinel), res)
	})

	t.Run("the handler is not run with a cancelled context", func(t *testing.T) {
		res := CatchAs(countFrom)(Left[counter, int](&notFoundError{id: 7}))(counter(0))(ctx)()
		assert.Equal(t, result.Left[Pair[counter, int]](context.Canceled), res)
	})
}
