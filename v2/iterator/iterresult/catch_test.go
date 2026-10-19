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

func run[A any](seq SeqResult[A]) []result.Result[A] { return slices.Collect(seq) }

func one[A any](r result.Result[A]) []result.Result[A] { return []result.Result[A]{r} }

func recoverNotFound(e *notFoundError) SeqResult[int] { return Of(e.id) }

func recoverCode(e codeError) SeqResult[int] { return Of(e.code) }

func TestCatchAs_Success(t *testing.T) {
	t.Run("recovers from a matching error", func(t *testing.T) {
		assert.Equal(t, one(result.Of(7)), run(CatchAs(recoverNotFound)(Left[int](&notFoundError{id: 7}))))
	})

	t.Run("recovers from a wrapped error", func(t *testing.T) {
		err := fmt.Errorf("lookup: %w", &notFoundError{id: 8})
		assert.Equal(t, one(result.Of(8)), run(CatchAs(recoverNotFound)(Left[int](err))))
	})

	t.Run("recovers from a joined error", func(t *testing.T) {
		err := errors.Join(errSentinel, &notFoundError{id: 9})
		assert.Equal(t, one(result.Of(9)), run(CatchAs(recoverNotFound)(Left[int](err))))
	})

	t.Run("recovers from an error with a value receiver", func(t *testing.T) {
		assert.Equal(t, one(result.Of(404)), run(CatchAs(recoverCode)(Left[int](codeError{code: 404}))))
	})

	t.Run("recovers from an interface error type", func(t *testing.T) {
		recoverTemp := func(temporary) SeqResult[int] { return Of(1) }
		assert.Equal(t, one(result.Of(1)), run(CatchAs(recoverTemp)(Left[int](tempError{}))))
	})

	t.Run("the handler can fail with a new error", func(t *testing.T) {
		replace := func(*notFoundError) SeqResult[int] { return Left[int](errSentinel) }
		assert.Equal(t, one(result.Left[int](errSentinel)), run(CatchAs(replace)(Left[int](&notFoundError{id: 7}))))
	})

	t.Run("passes a success through without calling the handler", func(t *testing.T) {
		called := false
		handler := func(*notFoundError) SeqResult[int] { called = true; return Of(0) }
		assert.Equal(t, one(result.Of(42)), run(CatchAs(handler)(Of(42))))
		assert.False(t, called)
	})
}

func TestCatchAs_Failure(t *testing.T) {
	t.Run("passes an error of another type through", func(t *testing.T) {
		assert.Equal(t, one(result.Left[int](errSentinel)), run(CatchAs(recoverNotFound)(Left[int](errSentinel))))
	})

	t.Run("pointer type does not match a value receiver error", func(t *testing.T) {
		recoverPtr := func(e *codeError) SeqResult[int] { return Of(e.code) }
		err := codeError{code: 500}
		assert.Equal(t, one(result.Left[int](error(err))), run(CatchAs(recoverPtr)(Left[int](err))))
	})

	t.Run("value type does not match a pointer to a value receiver error", func(t *testing.T) {
		err := &codeError{code: 500}
		assert.Equal(t, one(result.Left[int](error(err))), run(CatchAs(recoverCode)(Left[int](err))))
	})
}

func TestCatchIs_Success(t *testing.T) {
	fallback := func(error) SeqResult[int] { return Of(0) }

	t.Run("recovers from the sentinel", func(t *testing.T) {
		assert.Equal(t, one(result.Of(0)), run(CatchIs(errSentinel, fallback)(Left[int](errSentinel))))
	})

	t.Run("recovers from a wrapped sentinel", func(t *testing.T) {
		err := fmt.Errorf("read: %w", errSentinel)
		assert.Equal(t, one(result.Of(0)), run(CatchIs(errSentinel, fallback)(Left[int](err))))
	})

	t.Run("recovers from a joined sentinel", func(t *testing.T) {
		err := errors.Join(&notFoundError{id: 1}, errSentinel)
		assert.Equal(t, one(result.Of(0)), run(CatchIs(errSentinel, fallback)(Left[int](err))))
	})

	t.Run("passes the original error to the handler", func(t *testing.T) {
		err := fmt.Errorf("read: %w", errSentinel)
		var got error
		handler := func(e error) SeqResult[int] { got = e; return Of(0) }
		run(CatchIs(errSentinel, handler)(Left[int](err)))
		assert.Same(t, err, got)
	})

	t.Run("passes a success through", func(t *testing.T) {
		assert.Equal(t, one(result.Of(42)), run(CatchIs(errSentinel, fallback)(Of(42))))
	})
}

func TestCatchIs_Failure(t *testing.T) {
	other := errors.New("other")
	fallback := func(error) SeqResult[int] { return Of(0) }
	assert.Equal(t, one(result.Left[int](other)), run(CatchIs(errSentinel, fallback)(Left[int](other))))
}

func TestMapErrorAs_Success(t *testing.T) {
	toCode := func(e *notFoundError) error { return codeError{code: 404 + e.id} }

	t.Run("transforms a matching error", func(t *testing.T) {
		res := run(MapErrorAs[*notFoundError, int](toCode)(Left[int](&notFoundError{id: 0})))
		assert.Equal(t, one(result.Left[int](codeError{code: 404})), res)
	})

	t.Run("transforms a wrapped error", func(t *testing.T) {
		err := fmt.Errorf("lookup: %w", &notFoundError{id: 1})
		res := run(MapErrorAs[*notFoundError, int](toCode)(Left[int](err)))
		assert.Equal(t, one(result.Left[int](codeError{code: 405})), res)
	})

	t.Run("passes a success through", func(t *testing.T) {
		assert.Equal(t, one(result.Of(42)), run(MapErrorAs[*notFoundError, int](toCode)(Of(42))))
	})
}

func TestMapErrorAs_Failure(t *testing.T) {
	toCode := func(*notFoundError) error { return codeError{code: 404} }
	res := run(MapErrorAs[*notFoundError, int](toCode)(Left[int](errSentinel)))
	assert.Equal(t, one(result.Left[int](errSentinel)), res)
}

func TestCatch_Integration(t *testing.T) {
	lookup := func(id int) SeqResult[int] {
		switch id {
		case 0:
			return Left[int](&notFoundError{id: id})
		case 1:
			return Left[int](fmt.Errorf("db: %w", errSentinel))
		default:
			return Of(id * 10)
		}
	}
	handle := func(id int) SeqResult[int] {
		return CatchIs(errSentinel, func(error) SeqResult[int] { return Of(-2) })(
			CatchAs(func(*notFoundError) SeqResult[int] { return Of(-1) })(lookup(id)))
	}

	assert.Equal(t, one(result.Of(-1)), run(handle(0)))
	assert.Equal(t, one(result.Of(-2)), run(handle(1)))
	assert.Equal(t, one(result.Of(20)), run(handle(2)))
}

func TestTapErrorAs_Success(t *testing.T) {
	t.Run("observes a matching error and keeps it", func(t *testing.T) {
		var seen int
		observe := func(e *notFoundError) SeqResult[Void] {
			return FromIO(func() Void { seen = e.id; return F.VOID })
		}
		err := fmt.Errorf("lookup: %w", &notFoundError{id: 3})
		assert.Equal(t, one(result.Left[int](err)), run(TapErrorAs[*notFoundError, int](observe)(Left[int](err))))
		assert.Equal(t, 3, seen)
	})

	t.Run("a failing side effect does not replace the error", func(t *testing.T) {
		fail := func(*notFoundError) SeqResult[Void] { return Left[Void](errSentinel) }
		err := &notFoundError{id: 3}
		assert.Equal(t, one(result.Left[int](error(err))), run(TapErrorAs[*notFoundError, int](fail)(Left[int](err))))
	})
}

func TestTapErrorAs_Failure(t *testing.T) {
	called := false
	observe := func(*notFoundError) SeqResult[Void] {
		return FromIO(func() Void { called = true; return F.VOID })
	}

	t.Run("does not observe an error of another type", func(t *testing.T) {
		assert.Equal(t, one(result.Left[int](errSentinel)), run(TapErrorAs[*notFoundError, int](observe)(Left[int](errSentinel))))
		assert.False(t, called)
	})

	t.Run("does not observe a success", func(t *testing.T) {
		assert.Equal(t, one(result.Of(42)), run(TapErrorAs[*notFoundError, int](observe)(Of(42))))
		assert.False(t, called)
	})
}

func TestCatchAs_EdgeCases(t *testing.T) {
	t.Run("the handler can recover with several values", func(t *testing.T) {
		expand := func(e *notFoundError) SeqResult[int] { return FromSeq(slices.Values([]int{e.id, e.id + 1})) }
		res := run(CatchAs(expand)(Left[int](&notFoundError{id: 1})))
		assert.Equal(t, []result.Result[int]{result.Of(1), result.Of(2)}, res)
	})

	t.Run("recovers each matching element of the sequence", func(t *testing.T) {
		input := func(yield func(result.Result[int]) bool) {
			_ = yield(result.Of(1)) && yield(result.Left[int](&notFoundError{id: 2})) &&
				yield(result.Left[int](errSentinel)) && yield(result.Left[int](&notFoundError{id: 4}))
		}
		res := run(CatchAs(recoverNotFound)(input))
		assert.Equal(t, []result.Result[int]{result.Of(1), result.Of(2), result.Left[int](errSentinel), result.Of(4)}, res)
	})
}
