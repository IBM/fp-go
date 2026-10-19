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
	"slices"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/catch/catchtest"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// catchMonad describes the monad of this package for the shared catch test suite.
func catchMonad() catchtest.Monad[SeqResult[int], SeqResult[F.Void], []result.Result[int]] {
	return catchtest.Monad[SeqResult[int], SeqResult[F.Void], []result.Result[int]]{
		Of:       Of[int],
		Left:     Left[int],
		LeftVoid: Left[F.Void],
		Observe:  func(record func()) SeqResult[F.Void] { return FromIO(func() F.Void { record(); return F.VOID }) },
		Run:      func(s SeqResult[int]) []result.Result[int] { return slices.Collect(s) },
		Want:     func(r result.Result[int]) []result.Result[int] { return []result.Result[int]{r} },
	}
}

func TestCatchByType(t *testing.T) {
	catchtest.TestCatch(t, catchtest.CatchOps[SeqResult[int], SeqResult[F.Void], []result.Result[int]]{
		Monad:              catchMonad(),
		CatchAsNotFound:    CatchAs[*catchtest.NotFoundError, int],
		CatchAsCode:        CatchAs[catchtest.CodeError, int],
		CatchAsCodePtr:     CatchAs[*catchtest.CodeError, int],
		CatchAsTemporary:   CatchAs[catchtest.Temporary, int],
		CatchIs:            CatchIs[int],
		MapErrorAsNotFound: MapErrorAs[*catchtest.NotFoundError, int],
		TapErrorAsNotFound: TapErrorAs[*catchtest.NotFoundError, int, F.Void],
	})
}

func TestCatchAs_EdgeCases(t *testing.T) {
	recoverNotFound := func(e *catchtest.NotFoundError) SeqResult[int] { return Of(e.ID) }

	t.Run("the handler can recover with several values", func(t *testing.T) {
		expand := func(e *catchtest.NotFoundError) SeqResult[int] { return FromSeq(slices.Values([]int{e.ID, e.ID + 1})) }
		res := slices.Collect(CatchAs(expand)(Left[int](&catchtest.NotFoundError{ID: 1})))
		assert.Equal(t, []result.Result[int]{result.Of(1), result.Of(2)}, res)
	})

	t.Run("recovers each matching element of the sequence", func(t *testing.T) {
		input := func(yield func(result.Result[int]) bool) {
			_ = yield(result.Of(1)) && yield(result.Left[int](&catchtest.NotFoundError{ID: 2})) &&
				yield(result.Left[int](catchtest.ErrSentinel)) && yield(result.Left[int](&catchtest.NotFoundError{ID: 4}))
		}
		res := slices.Collect(CatchAs(recoverNotFound)(input))
		assert.Equal(t, []result.Result[int]{result.Of(1), result.Of(2), result.Left[int](catchtest.ErrSentinel), result.Of(4)}, res)
	})
}
