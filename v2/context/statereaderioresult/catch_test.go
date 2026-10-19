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
	"testing"

	RIORES "github.com/IBM/fp-go/v2/context/readerioresult"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/catch/catchtest"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// counter is the state of the computations in the tests.
type counter int

// catchMonad describes the monad of this package for the shared catch test suite.
func catchMonad() catchtest.Monad[StateReaderIOResult[counter, int], StateReaderIOResult[counter, F.Void], Result[int]] {
	return catchtest.Monad[StateReaderIOResult[counter, int], StateReaderIOResult[counter, F.Void], Result[int]]{
		Of:       Of[counter, int],
		Left:     Left[counter, int],
		LeftVoid: Left[counter, F.Void],
		Observe: func(record func()) StateReaderIOResult[counter, F.Void] {
			return FromIO[counter](func() F.Void { record(); return F.VOID })
		},
		Run: func(m StateReaderIOResult[counter, int]) Result[int] {
			return result.Map(pair.Tail[counter, int])(m(counter(0))(context.Background())())
		},
		Want: F.Identity[Result[int]],
	}
}

func TestCatchByType(t *testing.T) {
	catchtest.TestCatch(t, catchtest.CatchOps[StateReaderIOResult[counter, int], StateReaderIOResult[counter, F.Void], Result[int]]{
		Monad:              catchMonad(),
		CatchAsNotFound:    CatchAs[*catchtest.NotFoundError, counter, int],
		CatchAsCode:        CatchAs[catchtest.CodeError, counter, int],
		CatchAsCodePtr:     CatchAs[*catchtest.CodeError, counter, int],
		CatchAsTemporary:   CatchAs[catchtest.Temporary, counter, int],
		CatchIs:            CatchIs[counter, int],
		MapErrorAsNotFound: MapErrorAs[*catchtest.NotFoundError, counter, int],
		TapErrorAsNotFound: TapErrorAs[*catchtest.NotFoundError, counter, int, F.Void],
	})
}

// countFrom recovers with the current state as the value and increments the state.
func countFrom[E any](E) StateReaderIOResult[counter, int] {
	return func(s counter) ReaderIOResult[Pair[counter, int]] {
		return RIORES.Of(pair.MakePair(s+1, int(s)))
	}
}

func TestCatchAs_EdgeCases(t *testing.T) {
	t.Run("the handler starts from the state of the failed computation", func(t *testing.T) {
		res := CatchAs(countFrom[*catchtest.NotFoundError])(Left[counter, int](&catchtest.NotFoundError{ID: 1}))(counter(5))(t.Context())()
		assert.Equal(t, result.Of(pair.MakePair(counter(6), 5)), res)
	})

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	t.Run("a non-matching error passes through a cancelled context unchanged", func(t *testing.T) {
		res := CatchAs(countFrom[*catchtest.NotFoundError])(Left[counter, int](catchtest.ErrSentinel))(counter(0))(ctx)()
		assert.Equal(t, result.Left[Pair[counter, int]](catchtest.ErrSentinel), res)
	})

	t.Run("the handler is not run with a cancelled context", func(t *testing.T) {
		res := CatchAs(countFrom[*catchtest.NotFoundError])(Left[counter, int](&catchtest.NotFoundError{ID: 7}))(counter(0))(ctx)()
		assert.Equal(t, result.Left[Pair[counter, int]](context.Canceled), res)
	})
}
