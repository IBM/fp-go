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

package readerresult

import (
	"context"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/catch/catchtest"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

func TestErrorChannel(t *testing.T) {
	catchtest.TestLeft(t, catchtest.LeftOps[ReaderResult[int], ReaderResult[F.Void], result.Result[int]]{
		Monad:               catchMonad(),
		ChainFirstLeft:      ChainFirstLeft[int, F.Void],
		TapLeft:             TapLeft[int, F.Void],
		MonadChainFirstLeft: MonadChainFirstLeft[int, F.Void],
		MonadTapLeft:        MonadTapLeft[int, F.Void],
	})
}

func TestChainFirstLeft_EdgeCases(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	called := false
	observe := func(error) ReaderResult[F.Void] {
		return func(context.Context) (F.Void, error) { called = true; return F.VOID, nil }
	}
	res := result.TryCatchError(ChainFirstLeft[int](observe)(Left[int](catchtest.ErrSentinel))(ctx))
	assert.Equal(t, result.Left[int](catchtest.ErrSentinel), res)
	assert.False(t, called)
}
