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

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/catch/catchtest"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

func TestErrorChannel(t *testing.T) {
	catchtest.TestLeft(t, catchtest.LeftOps[StateReaderIOResult[counter, int], StateReaderIOResult[counter, F.Void], Result[int]]{
		Monad:               catchMonad(),
		ChainFirstLeft:      ChainFirstLeft[int, counter, F.Void],
		TapLeft:             TapLeft[int, counter, F.Void],
		MonadChainFirstLeft: MonadChainFirstLeft[counter, int, F.Void],
		MonadTapLeft:        MonadTapLeft[counter, int, F.Void],
		ChainLeft:           ChainLeft[counter, int],
		OrElse:              OrElse[counter, int],
		MonadChainLeft:      MonadChainLeft[counter, int],
		MapLeft:             MapLeft[counter, int],
		MonadMapLeft:        MonadMapLeft[counter, int],
	})
}

func TestChainLeft_EdgeCases(t *testing.T) {
	t.Run("the handler starts from the state of the failed computation", func(t *testing.T) {
		res := ChainLeft(countFrom[error])(Left[counter, int](catchtest.ErrSentinel))(counter(5))(t.Context())()
		assert.Equal(t, result.Of(pair.MakePair(counter(6), 5)), res)
	})

	t.Run("the state of a side effect is discarded", func(t *testing.T) {
		res := TapLeft[int](countFrom[error])(Left[counter, int](catchtest.ErrSentinel))(counter(5))(t.Context())()
		assert.Equal(t, result.Left[Pair[counter, int]](catchtest.ErrSentinel), res)
	})

	t.Run("the handler is not run with a cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		res := ChainLeft(countFrom[error])(Left[counter, int](catchtest.ErrSentinel))(counter(5))(ctx)()
		assert.Equal(t, result.Left[Pair[counter, int]](context.Canceled), res)
	})
}
