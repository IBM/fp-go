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
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/catch/catchtest"
	R "github.com/IBM/fp-go/v2/result"
)

// lazy defers a (value, error) pair, so that the shared catch test suite can treat it as
// a single value.
type lazy[A any] = func() (A, error)

// strict turns a handler that returns a lazy pair into a Kleisli arrow.
func strict[E, A any](f func(E) lazy[A]) Kleisli[E, A] {
	return func(e E) (A, error) { return f(e)() }
}

// lazyOp applies an operator to a lazy pair.
func lazyOp(op Operator[int, int]) func(lazy[int]) lazy[int] {
	return func(m lazy[int]) lazy[int] {
		return func() (int, error) { return op(m()) }
	}
}

// lift adapts a combinator that takes a Kleisli arrow to lazy pairs.
func lift[E, B any](op func(Kleisli[E, B]) Operator[int, int]) func(func(E) lazy[B]) func(lazy[int]) lazy[int] {
	return func(f func(E) lazy[B]) func(lazy[int]) lazy[int] { return lazyOp(op(strict(f))) }
}

// catchMonad describes the monad of this package for the shared catch test suite.
func catchMonad() catchtest.Monad[lazy[int], lazy[F.Void], R.Result[int]] {
	return catchtest.Monad[lazy[int], lazy[F.Void], R.Result[int]]{
		Of:       func(a int) lazy[int] { return func() (int, error) { return Of(a) } },
		Left:     func(err error) lazy[int] { return func() (int, error) { return Left[int](err) } },
		LeftVoid: func(err error) lazy[F.Void] { return func() (F.Void, error) { return Left[F.Void](err) } },
		Observe:  func(record func()) lazy[F.Void] { return func() (F.Void, error) { record(); return Of(F.VOID) } },
		Run:      func(m lazy[int]) R.Result[int] { return R.TryCatchError(m()) },
		Want:     F.Identity[R.Result[int]],
	}
}

func TestCatchByType(t *testing.T) {
	catchtest.TestCatch(t, catchtest.CatchOps[lazy[int], lazy[F.Void], R.Result[int]]{
		Monad:            catchMonad(),
		CatchAsNotFound:  lift(CatchAs[*catchtest.NotFoundError, int]),
		CatchAsCode:      lift(CatchAs[catchtest.CodeError, int]),
		CatchAsCodePtr:   lift(CatchAs[*catchtest.CodeError, int]),
		CatchAsTemporary: lift(CatchAs[catchtest.Temporary, int]),
		CatchIs: func(target error, f func(error) lazy[int]) func(lazy[int]) lazy[int] {
			return lazyOp(CatchIs(target, strict(f)))
		},
		MapErrorAsNotFound: func(f func(*catchtest.NotFoundError) error) func(lazy[int]) lazy[int] {
			return lazyOp(MapErrorAs[*catchtest.NotFoundError, int](f))
		},
		TapErrorAsNotFound: lift(TapErrorAs[*catchtest.NotFoundError, int, F.Void]),
	})
}
