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
	"errors"
	"testing"

	RIORES "github.com/IBM/fp-go/v2/context/readerioresult"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

func TestChainFirstLeft_Success(t *testing.T) {
	t.Run("observes the error and keeps it", func(t *testing.T) {
		var seen error
		observe := func(err error) StateReaderIOResult[counter, F.Void] { seen = err; return Of[counter](F.VOID) }
		assert.Equal(t, (result.Left[int](errSentinel)), run(ChainFirstLeft[int](observe)(Left[counter, int](errSentinel))))
		assert.Equal(t, errSentinel, seen)
	})

	t.Run("a failing side effect does not replace the error", func(t *testing.T) {
		fail := func(error) StateReaderIOResult[counter, F.Void] {
			return Left[counter, F.Void](errors.New("side effect"))
		}
		assert.Equal(t, (result.Left[int](errSentinel)), run(ChainFirstLeft[int](fail)(Left[counter, int](errSentinel))))
	})

	t.Run("TapLeft is the same as ChainFirstLeft", func(t *testing.T) {
		var seen error
		observe := func(err error) StateReaderIOResult[counter, F.Void] { seen = err; return Of[counter](F.VOID) }
		assert.Equal(t, (result.Left[int](errSentinel)), run(TapLeft[int](observe)(Left[counter, int](errSentinel))))
		assert.Equal(t, errSentinel, seen)
	})

	t.Run("the monadic forms behave the same", func(t *testing.T) {
		count := 0
		observe := func(error) StateReaderIOResult[counter, F.Void] { count++; return Of[counter](F.VOID) }
		assert.Equal(t, (result.Left[int](errSentinel)), run(MonadChainFirstLeft(Left[counter, int](errSentinel), observe)))
		assert.Equal(t, (result.Left[int](errSentinel)), run(MonadTapLeft(Left[counter, int](errSentinel), observe)))
		assert.Equal(t, 2, count)
	})
}

func TestChainFirstLeft_Failure(t *testing.T) {
	called := false
	observe := func(error) StateReaderIOResult[counter, F.Void] { called = true; return Of[counter](F.VOID) }
	assert.Equal(t, (result.Of(42)), run(ChainFirstLeft[int](observe)(Of[counter](42))))
	assert.False(t, called)
}

func TestChainLeft_Success(t *testing.T) {
	recoverAll := func(error) StateReaderIOResult[counter, int] { return Of[counter](1) }

	t.Run("recovers from the error", func(t *testing.T) {
		assert.Equal(t, (result.Of(1)), run(ChainLeft(recoverAll)(Left[counter, int](errSentinel))))
	})

	t.Run("OrElse and MonadChainLeft behave the same", func(t *testing.T) {
		assert.Equal(t, (result.Of(1)), run(OrElse(recoverAll)(Left[counter, int](errSentinel))))
		assert.Equal(t, (result.Of(1)), run(MonadChainLeft(Left[counter, int](errSentinel), recoverAll)))
	})

	t.Run("the handler can fail with a new error", func(t *testing.T) {
		other := errors.New("other")
		replace := func(error) StateReaderIOResult[counter, int] { return Left[counter, int](other) }
		assert.Equal(t, (result.Left[int](other)), run(ChainLeft(replace)(Left[counter, int](errSentinel))))
	})
}

func TestChainLeft_Failure(t *testing.T) {
	called := false
	handler := func(error) StateReaderIOResult[counter, int] { called = true; return Of[counter](0) }
	assert.Equal(t, (result.Of(42)), run(ChainLeft(handler)(Of[counter](42))))
	assert.False(t, called)
}

func TestMapLeft_Success(t *testing.T) {
	toCode := func(error) error { return codeError{code: 1} }

	t.Run("transforms the error", func(t *testing.T) {
		assert.Equal(t, (result.Left[int](codeError{code: 1})), run(MapLeft[counter, int](toCode)(Left[counter, int](errSentinel))))
	})

	t.Run("MonadMapLeft behaves the same", func(t *testing.T) {
		assert.Equal(t, (result.Left[int](codeError{code: 1})), run(MonadMapLeft(Left[counter, int](errSentinel), toCode)))
	})

	t.Run("passes a success through", func(t *testing.T) {
		assert.Equal(t, (result.Of(42)), run(MapLeft[counter, int](toCode)(Of[counter](42))))
	})
}

func TestChainLeft_EdgeCases(t *testing.T) {
	// fromState recovers with the current state as the value and increments the state.
	fromState := func(error) StateReaderIOResult[counter, int] {
		return func(s counter) ReaderIOResult[Pair[counter, int]] {
			return RIORES.Of(pair.MakePair(s+1, int(s)))
		}
	}

	t.Run("the handler starts from the state of the failed computation", func(t *testing.T) {
		res := ChainLeft(fromState)(Left[counter, int](errSentinel))(counter(5))(t.Context())()
		assert.Equal(t, result.Of(pair.MakePair(counter(6), 5)), res)
	})

	t.Run("the state of a side effect is discarded", func(t *testing.T) {
		res := TapLeft[int](fromState)(Left[counter, int](errSentinel))(counter(5))(t.Context())()
		assert.Equal(t, result.Left[Pair[counter, int]](errSentinel), res)
	})

	t.Run("the handler is not run with a cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()
		res := ChainLeft(fromState)(Left[counter, int](errSentinel))(counter(5))(ctx)()
		assert.Equal(t, result.Left[Pair[counter, int]](context.Canceled), res)
	})
}
