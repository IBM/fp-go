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
	"errors"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// countingIOResult returns an IOResult Kleisli that counts its executions and
// either doubles its input or fails with err, if err is not nil.
func countingIOResult(calls *int, err error) ioresult.Kleisli[int, int] {
	return func(n int) IOResult[int] {
		return func() Result[int] {
			*calls++
			if err != nil {
				return result.Left[int](err)
			}
			return result.Of(n * 2)
		}
	}
}

func TestChainIOResultK_Success(t *testing.T) {
	t.Run("chains a successful IOResult", func(t *testing.T) {
		calls := 0
		computation := F.Pipe1(Of[AppConfig](21), ChainIOResultK[AppConfig](countingIOResult(&calls, nil)))
		assert.Equal(t, result.Of(42), computation(defaultConfig)(t.Context())())
		assert.Equal(t, 1, calls)
	})

	t.Run("is lazy until executed", func(t *testing.T) {
		calls := 0
		computation := F.Pipe1(Of[AppConfig](21), ChainIOResultK[AppConfig](countingIOResult(&calls, nil)))
		assert.Equal(t, 0, calls)
		computation(defaultConfig)(t.Context())()
		computation(defaultConfig)(t.Context())()
		assert.Equal(t, 2, calls)
	})

	t.Run("monadic version", func(t *testing.T) {
		calls := 0
		computation := MonadChainIOResultK(Of[AppConfig](21), countingIOResult(&calls, nil))
		assert.Equal(t, result.Of(42), computation(defaultConfig)(t.Context())())
		assert.Equal(t, 1, calls)
	})
}

func TestChainIOResultK_Failure(t *testing.T) {
	t.Run("propagates error of the IOResult", func(t *testing.T) {
		err := errors.New("io failed")
		calls := 0
		computation := F.Pipe1(Of[AppConfig](21), ChainIOResultK[AppConfig](countingIOResult(&calls, err)))
		assert.Equal(t, result.Left[int](err), computation(defaultConfig)(t.Context())())
		assert.Equal(t, 1, calls)
	})

	t.Run("skips the IOResult on upstream failure", func(t *testing.T) {
		err := errors.New("upstream")
		calls := 0
		computation := MonadChainIOResultK(Left[AppConfig, int](err), countingIOResult(&calls, nil))
		assert.Equal(t, result.Left[int](err), computation(defaultConfig)(t.Context())())
		assert.Equal(t, 0, calls)
	})
}

func TestChainFirstIOResultK_Success(t *testing.T) {
	t.Run("keeps the original value", func(t *testing.T) {
		calls := 0
		computation := F.Pipe1(Of[AppConfig](21), ChainFirstIOResultK[AppConfig](countingIOResult(&calls, nil)))
		assert.Equal(t, result.Of(21), computation(defaultConfig)(t.Context())())
		assert.Equal(t, 1, calls)
	})

	t.Run("monadic version and tap aliases", func(t *testing.T) {
		calls := 0
		f := countingIOResult(&calls, nil)
		for _, computation := range []ReaderReaderIOResult[AppConfig, int]{
			MonadChainFirstIOResultK(Of[AppConfig](21), f),
			MonadTapIOResultK(Of[AppConfig](21), f),
			TapIOResultK[AppConfig](f)(Of[AppConfig](21)),
		} {
			assert.Equal(t, result.Of(21), computation(defaultConfig)(t.Context())())
		}
		assert.Equal(t, 3, calls)
	})
}

func TestChainFirstIOResultK_Failure(t *testing.T) {
	t.Run("error of the side effect replaces the value", func(t *testing.T) {
		err := errors.New("io failed")
		calls := 0
		f := countingIOResult(&calls, err)
		for _, computation := range []ReaderReaderIOResult[AppConfig, int]{
			ChainFirstIOResultK[AppConfig](f)(Of[AppConfig](21)),
			MonadChainFirstIOResultK(Of[AppConfig](21), f),
			MonadTapIOResultK(Of[AppConfig](21), f),
			TapIOResultK[AppConfig](f)(Of[AppConfig](21)),
		} {
			assert.Equal(t, result.Left[int](err), computation(defaultConfig)(t.Context())())
		}
		assert.Equal(t, 4, calls)
	})

	t.Run("skips the side effect on upstream failure", func(t *testing.T) {
		err := errors.New("upstream")
		calls := 0
		computation := F.Pipe1(Left[AppConfig, int](err), TapIOResultK[AppConfig](countingIOResult(&calls, nil)))
		assert.Equal(t, result.Left[int](err), computation(defaultConfig)(t.Context())())
		assert.Equal(t, 0, calls)
	})
}

func TestChainFirstIOEitherK_Success(t *testing.T) {
	t.Run("keeps the original value", func(t *testing.T) {
		calls := 0
		f := countingIOResult(&calls, nil)
		for _, computation := range []ReaderReaderIOResult[AppConfig, int]{
			ChainFirstIOEitherK[AppConfig](f)(Of[AppConfig](21)),
			MonadChainFirstIOEitherK(Of[AppConfig](21), f),
			MonadTapIOEitherK(Of[AppConfig](21), f),
			TapIOEitherK[AppConfig](f)(Of[AppConfig](21)),
		} {
			assert.Equal(t, result.Of(21), computation(defaultConfig)(t.Context())())
		}
		assert.Equal(t, 4, calls)
	})
}

func TestChainFirstIOEitherK_Failure(t *testing.T) {
	t.Run("error of the side effect replaces the value", func(t *testing.T) {
		err := errors.New("io failed")
		calls := 0
		f := countingIOResult(&calls, err)
		for _, computation := range []ReaderReaderIOResult[AppConfig, int]{
			ChainFirstIOEitherK[AppConfig](f)(Of[AppConfig](21)),
			MonadChainFirstIOEitherK(Of[AppConfig](21), f),
			MonadTapIOEitherK(Of[AppConfig](21), f),
			TapIOEitherK[AppConfig](f)(Of[AppConfig](21)),
		} {
			assert.Equal(t, result.Left[int](err), computation(defaultConfig)(t.Context())())
		}
		assert.Equal(t, 4, calls)
	})

	t.Run("skips the side effect on upstream failure", func(t *testing.T) {
		err := errors.New("upstream")
		calls := 0
		computation := MonadChainFirstIOEitherK(Left[AppConfig, int](err), countingIOResult(&calls, nil))
		assert.Equal(t, result.Left[int](err), computation(defaultConfig)(t.Context())())
		assert.Equal(t, 0, calls)
	})
}

func TestIOResultK_AgreesWithIOEitherK(t *testing.T) {
	err := errors.New("io failed")
	for _, f := range []func(n int) IOResult[int]{
		countingIOResult(new(int), nil),
		countingIOResult(new(int), err),
	} {
		for _, start := range []ReaderReaderIOResult[AppConfig, int]{Of[AppConfig](21), Left[AppConfig, int](err)} {
			run := func(m ReaderReaderIOResult[AppConfig, int]) Result[int] { return m(defaultConfig)(t.Context())() }
			assert.Equal(t, run(ChainIOEitherK[AppConfig](f)(start)), run(ChainIOResultK[AppConfig](f)(start)))
			assert.Equal(t, run(ChainFirstIOEitherK[AppConfig](f)(start)), run(ChainFirstIOResultK[AppConfig](f)(start)))
		}
	}
}

// positiveResult fails for numbers that are not positive.
func positiveResult(n int) Result[int] {
	if n > 0 {
		return result.Of(n)
	}
	return result.Left[int](errors.New("not positive"))
}

func TestMonadChainResultK(t *testing.T) {
	run := func(m ReaderReaderIOResult[AppConfig, int]) Result[int] { return m(defaultConfig)(t.Context())() }
	assert.Equal(t, result.Of(5), run(MonadChainResultK(Of[AppConfig](5), positiveResult)))
	assert.Equal(t, run(ChainResultK[AppConfig](positiveResult)(Of[AppConfig](-5))), run(MonadChainResultK(Of[AppConfig](-5), positiveResult)))
}

func TestChainFirstResultK(t *testing.T) {
	run := func(m ReaderReaderIOResult[AppConfig, int]) Result[int] { return m(defaultConfig)(t.Context())() }
	all := func(start ReaderReaderIOResult[AppConfig, int]) []ReaderReaderIOResult[AppConfig, int] {
		return []ReaderReaderIOResult[AppConfig, int]{
			ChainFirstResultK[AppConfig](positiveResult)(start),
			TapResultK[AppConfig](positiveResult)(start),
			MonadChainFirstResultK(start, positiveResult),
			MonadTapResultK(start, positiveResult),
		}
	}
	t.Run("keeps the value", func(t *testing.T) {
		for _, m := range all(Of[AppConfig](7)) {
			assert.Equal(t, result.Of(7), run(m))
		}
	})
	t.Run("failing check replaces the value", func(t *testing.T) {
		for _, m := range all(Of[AppConfig](-7)) {
			assert.Equal(t, result.Left[int](errors.New("not positive")), run(m))
		}
	})
	t.Run("upstream error is preserved", func(t *testing.T) {
		err := errors.New("upstream")
		for _, m := range all(Left[AppConfig, int](err)) {
			assert.Equal(t, result.Left[int](err), run(m))
		}
	})
}
