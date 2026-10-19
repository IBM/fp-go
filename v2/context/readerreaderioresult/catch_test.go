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

package readerreaderioresult

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

var errSentinel = errors.New("sentinel")

// recoverNotFound recovers with the id of the error plus the offset from the environment.
func recoverNotFound(e *notFoundError) ReaderReaderIOResult[int, int] {
	return Asks(func(offset int) int { return offset + e.id })
}

func recoverCode(e codeError) ReaderReaderIOResult[int, int] { return Of[int](e.code) }

func TestCatchAs_Success(t *testing.T) {
	run := func(eff ReaderReaderIOResult[int, int]) Result[int] { return eff(100)(t.Context())() }

	t.Run("recovers from a matching error using the environment", func(t *testing.T) {
		assert.Equal(t, result.Of(107), run(CatchAs(recoverNotFound)(Left[int, int](&notFoundError{id: 7}))))
	})

	t.Run("recovers from a wrapped error", func(t *testing.T) {
		err := fmt.Errorf("lookup: %w", &notFoundError{id: 8})
		assert.Equal(t, result.Of(108), run(CatchAs(recoverNotFound)(Left[int, int](err))))
	})

	t.Run("recovers from a joined error", func(t *testing.T) {
		err := errors.Join(errSentinel, &notFoundError{id: 9})
		assert.Equal(t, result.Of(109), run(CatchAs(recoverNotFound)(Left[int, int](err))))
	})

	t.Run("recovers from an error with a value receiver", func(t *testing.T) {
		assert.Equal(t, result.Of(404), run(CatchAs(recoverCode)(Left[int, int](codeError{code: 404}))))
	})

	t.Run("passes a success through", func(t *testing.T) {
		assert.Equal(t, result.Of(42), run(CatchAs(recoverNotFound)(Of[int](42))))
	})
}

func TestCatchAs_Failure(t *testing.T) {
	run := func(eff ReaderReaderIOResult[int, int]) Result[int] { return eff(0)(t.Context())() }

	t.Run("passes an error of another type through", func(t *testing.T) {
		assert.Equal(t, result.Left[int](errSentinel), run(CatchAs(recoverNotFound)(Left[int, int](errSentinel))))
	})

	t.Run("pointer type does not match a value receiver error", func(t *testing.T) {
		recoverPtr := func(e *codeError) ReaderReaderIOResult[int, int] { return Of[int](e.code) }
		err := codeError{code: 500}
		assert.Equal(t, result.Left[int](error(err)), run(CatchAs(recoverPtr)(Left[int, int](err))))
	})

	t.Run("value type does not match a pointer to a value receiver error", func(t *testing.T) {
		err := &codeError{code: 500}
		assert.Equal(t, result.Left[int](error(err)), run(CatchAs(recoverCode)(Left[int, int](err))))
	})
}

func TestCatchIs_Success(t *testing.T) {
	run := func(eff ReaderReaderIOResult[int, int]) Result[int] { return eff(0)(t.Context())() }
	fallback := func(error) ReaderReaderIOResult[int, int] { return Of[int](0) }

	t.Run("recovers from a wrapped sentinel", func(t *testing.T) {
		err := fmt.Errorf("read: %w", errSentinel)
		assert.Equal(t, result.Of(0), run(CatchIs(errSentinel, fallback)(Left[int, int](err))))
	})

	t.Run("recovers from a joined sentinel", func(t *testing.T) {
		err := errors.Join(&notFoundError{id: 1}, errSentinel)
		assert.Equal(t, result.Of(0), run(CatchIs(errSentinel, fallback)(Left[int, int](err))))
	})

	t.Run("passes the original error to the handler", func(t *testing.T) {
		err := fmt.Errorf("read: %w", errSentinel)
		var got error
		handler := func(e error) ReaderReaderIOResult[int, int] { got = e; return Of[int](0) }
		run(CatchIs(errSentinel, handler)(Left[int, int](err)))
		assert.Same(t, err, got)
	})
}

func TestCatchIs_Failure(t *testing.T) {
	other := errors.New("other")
	eff := CatchIs(errSentinel, func(error) ReaderReaderIOResult[int, int] { return Of[int](0) })(Left[int, int](other))
	assert.Equal(t, result.Left[int](other), eff(0)(t.Context())())
}

func TestMapErrorAs_Success(t *testing.T) {
	toCode := func(e *notFoundError) error { return codeError{code: 404 + e.id} }
	err := fmt.Errorf("lookup: %w", &notFoundError{id: 1})

	eff := MapErrorAs[*notFoundError, int, int](toCode)(Left[int, int](err))
	assert.Equal(t, result.Left[int](codeError{code: 405}), eff(0)(t.Context())())
}

func TestMapErrorAs_Failure(t *testing.T) {
	toCode := func(*notFoundError) error { return codeError{code: 404} }
	mapErr := MapErrorAs[*notFoundError, int, int](toCode)

	assert.Equal(t, result.Left[int](errSentinel), mapErr(Left[int, int](errSentinel))(0)(t.Context())())
	assert.Equal(t, result.Of(42), mapErr(Of[int](42))(0)(t.Context())())
}

func TestTapErrorAs_Success(t *testing.T) {
	var seen int
	observe := func(e *notFoundError) ReaderReaderIOResult[int, Void] {
		return FromIO[int](func() Void { seen = e.id; return F.VOID })
	}
	err := fmt.Errorf("lookup: %w", &notFoundError{id: 3})

	eff := TapErrorAs[*notFoundError, int, int](observe)(Left[int, int](err))
	assert.Equal(t, result.Left[int](err), eff(0)(t.Context())())
	assert.Equal(t, 3, seen)
}

func TestTapErrorAs_Failure(t *testing.T) {
	called := false
	observe := func(*notFoundError) ReaderReaderIOResult[int, Void] {
		return FromIO[int](func() Void { called = true; return F.VOID })
	}
	tap := TapErrorAs[*notFoundError, int, int](observe)

	assert.Equal(t, result.Left[int](errSentinel), tap(Left[int, int](errSentinel))(0)(t.Context())())
	assert.Equal(t, result.Of(42), tap(Of[int](42))(0)(t.Context())())
	assert.False(t, called)
}

func TestCatch_Integration(t *testing.T) {
	lookup := func(id int) ReaderReaderIOResult[int, int] {
		switch id {
		case 0:
			return Left[int, int](&notFoundError{id: id})
		case 1:
			return Left[int, int](fmt.Errorf("db: %w", errSentinel))
		default:
			return Of[int](id * 10)
		}
	}
	handle := F.Flow3(
		lookup,
		CatchAs(func(*notFoundError) ReaderReaderIOResult[int, int] { return Of[int](-1) }),
		CatchIs(errSentinel, func(error) ReaderReaderIOResult[int, int] { return Of[int](-2) }),
	)

	assert.Equal(t, result.Of(-1), handle(0)(0)(t.Context())())
	assert.Equal(t, result.Of(-2), handle(1)(0)(t.Context())())
	assert.Equal(t, result.Of(20), handle(2)(0)(t.Context())())
}
