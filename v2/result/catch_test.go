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
)

// catchMonad describes the monad of this package for the shared catch test suite.
func catchMonad() catchtest.Monad[Result[int], Result[F.Void], Result[int]] {
	return catchtest.Monad[Result[int], Result[F.Void], Result[int]]{
		Of:       Of[int],
		Left:     Left[int],
		LeftVoid: Left[F.Void],
		Observe:  func(record func()) Result[F.Void] { record(); return Of(F.VOID) },
		Run:      F.Identity[Result[int]],
		Want:     F.Identity[Result[int]],
	}
}

func TestCatchByType(t *testing.T) {
	catchtest.TestCatch(t, catchtest.CatchOps[Result[int], Result[F.Void], Result[int]]{
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
