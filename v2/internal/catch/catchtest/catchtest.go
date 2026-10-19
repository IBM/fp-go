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

// Package catchtest is a shared test suite for the error channel of the monads whose error
// type is error: the catch-by-type combinators (CatchAs, CatchIs, MapErrorAs, TapErrorAs)
// and the primitives they build on (ChainLeft, OrElse, MapLeft, ChainFirstLeft, TapLeft).
//
// A package describes its monad with a Monad value and passes the functions under test,
// instantiated for the fixture types of this package, in a CatchOps or LeftOps value. Go
// cannot pass generic functions without instantiating them, so the suite fixes the success
// type to int and the side effect type to Void.
//
// The suite only depends on either, so the tests of every package above it can use it.
package catchtest

import (
	"errors"
	"fmt"
	"testing"

	ET "github.com/IBM/fp-go/v2/either"
	"github.com/stretchr/testify/assert"
)

// NotFoundError implements error with a pointer receiver.
type NotFoundError struct{ ID int }

func (e *NotFoundError) Error() string { return fmt.Sprintf("%d not found", e.ID) }

// CodeError implements error with a value receiver.
type CodeError struct{ Code int }

func (e CodeError) Error() string { return fmt.Sprintf("code %d", e.Code) }

// Temporary is an interface error type.
type Temporary interface {
	error
	Temporary() bool
}

// TempError implements Temporary.
type TempError struct{}

func (TempError) Error() string   { return "temporary" }
func (TempError) Temporary() bool { return true }

// ErrSentinel is a sentinel error for errors.Is.
var ErrSentinel = errors.New("sentinel")

// Monad describes how to build and run the computations of a monad.
//
// Type Parameters:
//   - HKTInt: The computation that produces an int
//   - HKTVoid: The computation that produces a side effect
//   - R: The value a computation is run into, compared with assert.Equal
type Monad[HKTInt, HKTVoid, R any] struct {
	// Of creates a successful computation.
	Of func(int) HKTInt
	// Left creates a failed computation.
	Left func(error) HKTInt
	// LeftVoid creates a failed side effect.
	LeftVoid func(error) HKTVoid
	// Observe creates a successful side effect that calls record when it runs.
	Observe func(record func()) HKTVoid
	// Run runs a computation.
	Run func(HKTInt) R
	// Want is what Run returns for a computation with the given outcome.
	Want func(ET.Either[error, int]) R
}

// CatchOps holds the catch-by-type combinators of a monad, instantiated for the fixtures.
type CatchOps[HKTInt, HKTVoid, R any] struct {
	Monad[HKTInt, HKTVoid, R]

	CatchAsNotFound    func(func(*NotFoundError) HKTInt) func(HKTInt) HKTInt
	CatchAsCode        func(func(CodeError) HKTInt) func(HKTInt) HKTInt
	CatchAsCodePtr     func(func(*CodeError) HKTInt) func(HKTInt) HKTInt
	CatchAsTemporary   func(func(Temporary) HKTInt) func(HKTInt) HKTInt
	CatchIs            func(error, func(error) HKTInt) func(HKTInt) HKTInt
	MapErrorAsNotFound func(func(*NotFoundError) error) func(HKTInt) HKTInt
	TapErrorAsNotFound func(func(*NotFoundError) HKTVoid) func(HKTInt) HKTInt
}

// LeftOps holds the error channel primitives of a monad. Primitives the monad does not
// have are nil, and their tests are skipped.
type LeftOps[HKTInt, HKTVoid, R any] struct {
	Monad[HKTInt, HKTVoid, R]

	ChainFirstLeft      func(func(error) HKTVoid) func(HKTInt) HKTInt
	TapLeft             func(func(error) HKTVoid) func(HKTInt) HKTInt
	MonadChainFirstLeft func(HKTInt, func(error) HKTVoid) HKTInt
	MonadTapLeft        func(HKTInt, func(error) HKTVoid) HKTInt

	ChainLeft      func(func(error) HKTInt) func(HKTInt) HKTInt
	OrElse         func(func(error) HKTInt) func(HKTInt) HKTInt
	MonadChainLeft func(HKTInt, func(error) HKTInt) HKTInt
	MapLeft        func(func(error) error) func(HKTInt) HKTInt
	MonadMapLeft   func(HKTInt, func(error) error) HKTInt
}

func right(a int) ET.Either[error, int] { return ET.Right[error](a) }

func left(err error) ET.Either[error, int] { return ET.Left[int](err) }

// TestCatch runs the tests of CatchAs, CatchIs, MapErrorAs and TapErrorAs.
func TestCatch[HKTInt, HKTVoid, R any](t *testing.T, ops CatchOps[HKTInt, HKTVoid, R]) {
	m := ops.Monad
	recoverNotFound := func(e *NotFoundError) HKTInt { return m.Of(e.ID) }
	recoverCode := func(e CodeError) HKTInt { return m.Of(e.Code) }
	check := func(t *testing.T, want ET.Either[error, int], got HKTInt) {
		t.Helper()
		assert.Equal(t, m.Want(want), m.Run(got))
	}

	t.Run("CatchAs", func(t *testing.T) {
		t.Run("recovers from a matching error", func(t *testing.T) {
			check(t, right(7), ops.CatchAsNotFound(recoverNotFound)(m.Left(&NotFoundError{ID: 7})))
		})

		t.Run("recovers from a wrapped error", func(t *testing.T) {
			err := fmt.Errorf("lookup: %w", &NotFoundError{ID: 8})
			check(t, right(8), ops.CatchAsNotFound(recoverNotFound)(m.Left(err)))
		})

		t.Run("recovers from a joined error", func(t *testing.T) {
			err := errors.Join(ErrSentinel, &NotFoundError{ID: 9})
			check(t, right(9), ops.CatchAsNotFound(recoverNotFound)(m.Left(err)))
		})

		t.Run("recovers from an error with a value receiver", func(t *testing.T) {
			check(t, right(404), ops.CatchAsCode(recoverCode)(m.Left(CodeError{Code: 404})))
		})

		t.Run("recovers from an interface error type", func(t *testing.T) {
			recoverTemp := func(Temporary) HKTInt { return m.Of(1) }
			check(t, right(1), ops.CatchAsTemporary(recoverTemp)(m.Left(TempError{})))
		})

		t.Run("the handler can fail with a new error", func(t *testing.T) {
			replace := func(*NotFoundError) HKTInt { return m.Left(ErrSentinel) }
			check(t, left(ErrSentinel), ops.CatchAsNotFound(replace)(m.Left(&NotFoundError{ID: 7})))
		})

		t.Run("passes a success through without calling the handler", func(t *testing.T) {
			called := false
			handler := func(*NotFoundError) HKTInt { called = true; return m.Of(0) }
			check(t, right(42), ops.CatchAsNotFound(handler)(m.Of(42)))
			assert.False(t, called)
		})

		t.Run("passes an error of another type through", func(t *testing.T) {
			check(t, left(ErrSentinel), ops.CatchAsNotFound(recoverNotFound)(m.Left(ErrSentinel)))
		})

		t.Run("pointer type does not match a value receiver error", func(t *testing.T) {
			recoverPtr := func(e *CodeError) HKTInt { return m.Of(e.Code) }
			err := CodeError{Code: 500}
			check(t, left(err), ops.CatchAsCodePtr(recoverPtr)(m.Left(err)))
		})

		t.Run("value type does not match a pointer to a value receiver error", func(t *testing.T) {
			err := &CodeError{Code: 500}
			check(t, left(err), ops.CatchAsCode(recoverCode)(m.Left(err)))
		})
	})

	t.Run("CatchIs", func(t *testing.T) {
		fallback := func(error) HKTInt { return m.Of(0) }

		t.Run("recovers from the sentinel", func(t *testing.T) {
			check(t, right(0), ops.CatchIs(ErrSentinel, fallback)(m.Left(ErrSentinel)))
		})

		t.Run("recovers from a wrapped sentinel", func(t *testing.T) {
			err := fmt.Errorf("read: %w", ErrSentinel)
			check(t, right(0), ops.CatchIs(ErrSentinel, fallback)(m.Left(err)))
		})

		t.Run("recovers from a joined sentinel", func(t *testing.T) {
			err := errors.Join(&NotFoundError{ID: 1}, ErrSentinel)
			check(t, right(0), ops.CatchIs(ErrSentinel, fallback)(m.Left(err)))
		})

		t.Run("passes the original error to the handler", func(t *testing.T) {
			err := fmt.Errorf("read: %w", ErrSentinel)
			var got error
			handler := func(e error) HKTInt { got = e; return m.Of(0) }
			m.Run(ops.CatchIs(ErrSentinel, handler)(m.Left(err)))
			assert.Same(t, err, got)
		})

		t.Run("passes a success through", func(t *testing.T) {
			check(t, right(42), ops.CatchIs(ErrSentinel, fallback)(m.Of(42)))
		})

		t.Run("passes another error through", func(t *testing.T) {
			other := errors.New("other")
			check(t, left(other), ops.CatchIs(ErrSentinel, fallback)(m.Left(other)))
		})
	})

	t.Run("MapErrorAs", func(t *testing.T) {
		toCode := func(e *NotFoundError) error { return CodeError{Code: 404 + e.ID} }

		t.Run("transforms a matching error", func(t *testing.T) {
			check(t, left(CodeError{Code: 404}), ops.MapErrorAsNotFound(toCode)(m.Left(&NotFoundError{ID: 0})))
		})

		t.Run("transforms a wrapped error", func(t *testing.T) {
			err := fmt.Errorf("lookup: %w", &NotFoundError{ID: 1})
			check(t, left(CodeError{Code: 405}), ops.MapErrorAsNotFound(toCode)(m.Left(err)))
		})

		t.Run("passes a success through", func(t *testing.T) {
			check(t, right(42), ops.MapErrorAsNotFound(toCode)(m.Of(42)))
		})

		t.Run("passes an error of another type through", func(t *testing.T) {
			check(t, left(ErrSentinel), ops.MapErrorAsNotFound(toCode)(m.Left(ErrSentinel)))
		})
	})

	t.Run("TapErrorAs", func(t *testing.T) {
		t.Run("observes a matching error and keeps it", func(t *testing.T) {
			var seen int
			observe := func(e *NotFoundError) HKTVoid { return m.Observe(func() { seen = e.ID }) }
			err := fmt.Errorf("lookup: %w", &NotFoundError{ID: 3})
			check(t, left(err), ops.TapErrorAsNotFound(observe)(m.Left(err)))
			assert.Equal(t, 3, seen)
		})

		t.Run("a failing side effect does not replace the error", func(t *testing.T) {
			fail := func(*NotFoundError) HKTVoid { return m.LeftVoid(ErrSentinel) }
			err := &NotFoundError{ID: 3}
			check(t, left(err), ops.TapErrorAsNotFound(fail)(m.Left(err)))
		})

		t.Run("does not observe an error of another type or a success", func(t *testing.T) {
			called := false
			observe := func(*NotFoundError) HKTVoid { return m.Observe(func() { called = true }) }
			check(t, left(ErrSentinel), ops.TapErrorAsNotFound(observe)(m.Left(ErrSentinel)))
			check(t, right(42), ops.TapErrorAsNotFound(observe)(m.Of(42)))
			assert.False(t, called)
		})
	})

	t.Run("CatchAs and CatchIs compose", func(t *testing.T) {
		lookup := func(id int) HKTInt {
			switch id {
			case 0:
				return m.Left(&NotFoundError{ID: id})
			case 1:
				return m.Left(fmt.Errorf("db: %w", ErrSentinel))
			default:
				return m.Of(id * 10)
			}
		}
		handle := func(id int) HKTInt {
			return ops.CatchIs(ErrSentinel, func(error) HKTInt { return m.Of(-2) })(
				ops.CatchAsNotFound(func(*NotFoundError) HKTInt { return m.Of(-1) })(lookup(id)))
		}

		check(t, right(-1), handle(0))
		check(t, right(-2), handle(1))
		check(t, right(20), handle(2))
	})
}

// TestLeft runs the tests of the error channel primitives that ops provides.
func TestLeft[HKTInt, HKTVoid, R any](t *testing.T, ops LeftOps[HKTInt, HKTVoid, R]) {
	m := ops.Monad
	check := func(t *testing.T, want ET.Either[error, int], got HKTInt) {
		t.Helper()
		assert.Equal(t, m.Want(want), m.Run(got))
	}

	taps := map[string]func(func(error) HKTVoid) func(HKTInt) HKTInt{
		"ChainFirstLeft": ops.ChainFirstLeft,
		"TapLeft":        ops.TapLeft,
	}
	for name, monadic := range map[string]func(HKTInt, func(error) HKTVoid) HKTInt{
		"MonadChainFirstLeft": ops.MonadChainFirstLeft,
		"MonadTapLeft":        ops.MonadTapLeft,
	} {
		if monadic != nil {
			taps[name] = func(f func(error) HKTVoid) func(HKTInt) HKTInt {
				return func(ma HKTInt) HKTInt { return monadic(ma, f) }
			}
		}
	}
	for name, tap := range taps {
		if tap == nil {
			continue
		}
		t.Run(name, func(t *testing.T) {
			t.Run("observes the error and keeps it", func(t *testing.T) {
				var seen error
				observe := func(err error) HKTVoid { return m.Observe(func() { seen = err }) }
				check(t, left(ErrSentinel), tap(observe)(m.Left(ErrSentinel)))
				assert.Equal(t, ErrSentinel, seen)
			})

			t.Run("a failing side effect does not replace the error", func(t *testing.T) {
				fail := func(error) HKTVoid { return m.LeftVoid(errors.New("side effect")) }
				check(t, left(ErrSentinel), tap(fail)(m.Left(ErrSentinel)))
			})

			t.Run("does not observe a success", func(t *testing.T) {
				called := false
				observe := func(error) HKTVoid { return m.Observe(func() { called = true }) }
				check(t, right(42), tap(observe)(m.Of(42)))
				assert.False(t, called)
			})
		})
	}

	chains := map[string]func(func(error) HKTInt) func(HKTInt) HKTInt{
		"ChainLeft": ops.ChainLeft,
		"OrElse":    ops.OrElse,
	}
	if ops.MonadChainLeft != nil {
		chains["MonadChainLeft"] = func(f func(error) HKTInt) func(HKTInt) HKTInt {
			return func(ma HKTInt) HKTInt { return ops.MonadChainLeft(ma, f) }
		}
	}
	for name, chain := range chains {
		if chain == nil {
			continue
		}
		t.Run(name, func(t *testing.T) {
			t.Run("recovers from the error", func(t *testing.T) {
				check(t, right(1), chain(func(error) HKTInt { return m.Of(1) })(m.Left(ErrSentinel)))
			})

			t.Run("the handler can fail with a new error", func(t *testing.T) {
				other := errors.New("other")
				check(t, left(other), chain(func(error) HKTInt { return m.Left(other) })(m.Left(ErrSentinel)))
			})

			t.Run("passes a success through without calling the handler", func(t *testing.T) {
				called := false
				handler := func(error) HKTInt { called = true; return m.Of(0) }
				check(t, right(42), chain(handler)(m.Of(42)))
				assert.False(t, called)
			})
		})
	}

	maps := map[string]func(func(error) error) func(HKTInt) HKTInt{"MapLeft": ops.MapLeft}
	if ops.MonadMapLeft != nil {
		maps["MonadMapLeft"] = func(f func(error) error) func(HKTInt) HKTInt {
			return func(ma HKTInt) HKTInt { return ops.MonadMapLeft(ma, f) }
		}
	}
	for name, mapLeft := range maps {
		if mapLeft == nil {
			continue
		}
		t.Run(name, func(t *testing.T) {
			toCode := func(error) error { return CodeError{Code: 1} }

			t.Run("transforms the error", func(t *testing.T) {
				check(t, left(CodeError{Code: 1}), mapLeft(toCode)(m.Left(ErrSentinel)))
			})

			t.Run("passes a success through", func(t *testing.T) {
				check(t, right(42), mapLeft(toCode)(m.Of(42)))
			})
		})
	}
}
