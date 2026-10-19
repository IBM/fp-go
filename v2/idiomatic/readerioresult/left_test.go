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
	"errors"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

func TestChainFirstLeft_Success(t *testing.T) {
	t.Run("observes the error and keeps it", func(t *testing.T) {
		var seen error
		observe := func(err error) ReaderIOResult[config, F.Void] { seen = err; return Of[config](F.VOID) }
		assert.Equal(t, (result.Left[int](errSentinel)), run(ChainFirstLeft[int](observe)(Left[config, int](errSentinel))))
		assert.Equal(t, errSentinel, seen)
	})

	t.Run("a failing side effect does not replace the error", func(t *testing.T) {
		fail := func(error) ReaderIOResult[config, F.Void] { return Left[config, F.Void](errors.New("side effect")) }
		assert.Equal(t, (result.Left[int](errSentinel)), run(ChainFirstLeft[int](fail)(Left[config, int](errSentinel))))
	})

	t.Run("TapLeft is the same as ChainFirstLeft", func(t *testing.T) {
		var seen error
		observe := func(err error) ReaderIOResult[config, F.Void] { seen = err; return Of[config](F.VOID) }
		assert.Equal(t, (result.Left[int](errSentinel)), run(TapLeft[int](observe)(Left[config, int](errSentinel))))
		assert.Equal(t, errSentinel, seen)
	})

	t.Run("the monadic forms behave the same", func(t *testing.T) {
		count := 0
		observe := func(error) ReaderIOResult[config, F.Void] { count++; return Of[config](F.VOID) }
		assert.Equal(t, (result.Left[int](errSentinel)), run(MonadChainFirstLeft(Left[config, int](errSentinel), observe)))
		assert.Equal(t, (result.Left[int](errSentinel)), run(MonadTapLeft(Left[config, int](errSentinel), observe)))
		assert.Equal(t, 2, count)
	})
}

func TestChainFirstLeft_Failure(t *testing.T) {
	called := false
	observe := func(error) ReaderIOResult[config, F.Void] { called = true; return Of[config](F.VOID) }
	assert.Equal(t, (result.Of(42)), run(ChainFirstLeft[int](observe)(Of[config](42))))
	assert.False(t, called)
}

func TestChainLeft_Success(t *testing.T) {
	recoverAll := func(error) ReaderIOResult[config, int] { return Of[config](1) }

	t.Run("recovers from the error", func(t *testing.T) {
		assert.Equal(t, (result.Of(1)), run(ChainLeft(recoverAll)(Left[config, int](errSentinel))))
	})

	t.Run("OrElse and MonadChainLeft behave the same", func(t *testing.T) {
		assert.Equal(t, (result.Of(1)), run(OrElse(recoverAll)(Left[config, int](errSentinel))))
		assert.Equal(t, (result.Of(1)), run(MonadChainLeft(Left[config, int](errSentinel), recoverAll)))
	})

	t.Run("the handler can fail with a new error", func(t *testing.T) {
		other := errors.New("other")
		replace := func(error) ReaderIOResult[config, int] { return Left[config, int](other) }
		assert.Equal(t, (result.Left[int](other)), run(ChainLeft(replace)(Left[config, int](errSentinel))))
	})
}

func TestChainLeft_Failure(t *testing.T) {
	called := false
	handler := func(error) ReaderIOResult[config, int] { called = true; return Of[config](0) }
	assert.Equal(t, (result.Of(42)), run(ChainLeft(handler)(Of[config](42))))
	assert.False(t, called)
}

func TestMapLeft_Success(t *testing.T) {
	toCode := func(error) error { return codeError{code: 1} }

	t.Run("transforms the error", func(t *testing.T) {
		assert.Equal(t, (result.Left[int](codeError{code: 1})), run(MapLeft[config, int](toCode)(Left[config, int](errSentinel))))
	})

	t.Run("MonadMapLeft behaves the same", func(t *testing.T) {
		assert.Equal(t, (result.Left[int](codeError{code: 1})), run(MonadMapLeft(Left[config, int](errSentinel), toCode)))
	})

	t.Run("passes a success through", func(t *testing.T) {
		assert.Equal(t, (result.Of(42)), run(MapLeft[config, int](toCode)(Of[config](42))))
	})
}
