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
	"context"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/catch/catchtest"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// catchMonad describes the monad of this package for the shared catch test suite.
func catchMonad() catchtest.Monad[ReaderReaderIOResult[int, int], ReaderReaderIOResult[int, F.Void], Result[int]] {
	return catchtest.Monad[ReaderReaderIOResult[int, int], ReaderReaderIOResult[int, F.Void], Result[int]]{
		Of:       Of[int, int],
		Left:     Left[int, int],
		LeftVoid: Left[int, F.Void],
		Observe: func(record func()) ReaderReaderIOResult[int, F.Void] {
			return FromIO[int](func() F.Void { record(); return F.VOID })
		},
		Run:  func(m ReaderReaderIOResult[int, int]) Result[int] { return m(0)(context.Background())() },
		Want: F.Identity[Result[int]],
	}
}

func TestCatchByType(t *testing.T) {
	catchtest.TestCatch(t, catchtest.CatchOps[ReaderReaderIOResult[int, int], ReaderReaderIOResult[int, F.Void], Result[int]]{
		Monad:              catchMonad(),
		CatchAsNotFound:    CatchAs[*catchtest.NotFoundError, int, int],
		CatchAsCode:        CatchAs[catchtest.CodeError, int, int],
		CatchAsCodePtr:     CatchAs[*catchtest.CodeError, int, int],
		CatchAsTemporary:   CatchAs[catchtest.Temporary, int, int],
		CatchIs:            CatchIs[int, int],
		MapErrorAsNotFound: MapErrorAs[*catchtest.NotFoundError, int, int],
		TapErrorAsNotFound: TapErrorAs[*catchtest.NotFoundError, int, int, F.Void],
	})
}

func TestCatchAs_EdgeCases(t *testing.T) {
	t.Run("recovers from a matching error using the environment", func(t *testing.T) {
		recoverNotFound := func(e *catchtest.NotFoundError) ReaderReaderIOResult[int, int] {
			return Asks(func(offset int) int { return offset + e.ID })
		}
		res := CatchAs(recoverNotFound)(Left[int, int](&catchtest.NotFoundError{ID: 7}))(100)(t.Context())()
		assert.Equal(t, result.Of(107), res)
	})
}
