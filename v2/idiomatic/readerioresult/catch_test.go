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
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/catch/catchtest"
	"github.com/IBM/fp-go/v2/result"
)

// config is the environment of the computations in the tests.
type config struct{}

// catchMonad describes the monad of this package for the shared catch test suite.
func catchMonad() catchtest.Monad[ReaderIOResult[config, int], ReaderIOResult[config, F.Void], result.Result[int]] {
	return catchtest.Monad[ReaderIOResult[config, int], ReaderIOResult[config, F.Void], result.Result[int]]{
		Of:       Of[config, int],
		Left:     Left[config, int],
		LeftVoid: Left[config, F.Void],
		Observe: func(record func()) ReaderIOResult[config, F.Void] {
			return func(config) IOResult[F.Void] {
				return func() (F.Void, error) { record(); return F.VOID, nil }
			}
		},
		Run:  func(m ReaderIOResult[config, int]) result.Result[int] { return result.TryCatchError(m(config{})()) },
		Want: F.Identity[result.Result[int]],
	}
}

func TestCatchByType(t *testing.T) {
	catchtest.TestCatch(t, catchtest.CatchOps[ReaderIOResult[config, int], ReaderIOResult[config, F.Void], result.Result[int]]{
		Monad:              catchMonad(),
		CatchAsNotFound:    CatchAs[*catchtest.NotFoundError, config, int],
		CatchAsCode:        CatchAs[catchtest.CodeError, config, int],
		CatchAsCodePtr:     CatchAs[*catchtest.CodeError, config, int],
		CatchAsTemporary:   CatchAs[catchtest.Temporary, config, int],
		CatchIs:            CatchIs[config, int],
		MapErrorAsNotFound: MapErrorAs[*catchtest.NotFoundError, config, int],
		TapErrorAsNotFound: TapErrorAs[*catchtest.NotFoundError, config, int, F.Void],
	})
}
