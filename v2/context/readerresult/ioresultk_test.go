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
	"errors"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// iorKCounter counts its executions and doubles its input or fails with err, if err is not nil.
func iorKCounter(calls *int, err error) ioresult.Kleisli[int, int] {
	return func(n int) ioresult.IOResult[int] {
		return func() result.Result[int] {
			*calls++
			if err != nil {
				return result.Left[int](err)
			}
			return result.Of(n * 2)
		}
	}
}

func TestChainIOResultK_Success(t *testing.T) {
	calls := 0
	f := iorKCounter(&calls, nil)
	for _, m := range []ReaderResult[int]{
		F.Pipe1(Of(21), ChainIOEitherK(f)),
		F.Pipe1(Of(21), ChainIOResultK(f)),
		MonadChainIOEitherK(Of(21), f),
		MonadChainIOResultK(Of(21), f),
	} {
		assert.Equal(t, result.Of(42), m(t.Context()))
	}
	assert.Equal(t, 4, calls)
}

func TestChainIOResultK_Failure(t *testing.T) {
	t.Run("propagates the error of the IOResult", func(t *testing.T) {
		err := errors.New("io failed")
		f := iorKCounter(new(int), err)
		for _, m := range []ReaderResult[int]{
			F.Pipe1(Of(21), ChainIOResultK(f)),
			MonadChainIOEitherK(Of(21), f),
			MonadChainIOResultK(Of(21), f),
		} {
			assert.Equal(t, result.Left[int](err), m(t.Context()))
		}
	})

	t.Run("skips the IOResult on upstream failure", func(t *testing.T) {
		err := errors.New("upstream")
		calls := 0
		m := MonadChainIOResultK(Left[int](err), iorKCounter(&calls, nil))
		assert.Equal(t, result.Left[int](err), m(t.Context()))
		assert.Equal(t, 0, calls)
	})
}

func TestChainFirstIOResultK_Success(t *testing.T) {
	calls := 0
	f := iorKCounter(&calls, nil)
	for _, m := range []ReaderResult[int]{
		F.Pipe1(Of(21), ChainFirstIOEitherK(f)),
		F.Pipe1(Of(21), TapIOEitherK(f)),
		F.Pipe1(Of(21), ChainFirstIOResultK(f)),
		F.Pipe1(Of(21), TapIOResultK(f)),
		MonadChainFirstIOEitherK(Of(21), f),
		MonadTapIOEitherK(Of(21), f),
		MonadChainFirstIOResultK(Of(21), f),
		MonadTapIOResultK(Of(21), f),
	} {
		assert.Equal(t, result.Of(21), m(t.Context()))
	}
	assert.Equal(t, 8, calls)
}

func TestChainFirstIOResultK_Failure(t *testing.T) {
	t.Run("error of the side effect replaces the value", func(t *testing.T) {
		err := errors.New("io failed")
		f := iorKCounter(new(int), err)
		for _, m := range []ReaderResult[int]{
			F.Pipe1(Of(21), ChainFirstIOEitherK(f)),
			F.Pipe1(Of(21), TapIOResultK(f)),
			MonadChainFirstIOResultK(Of(21), f),
			MonadTapIOEitherK(Of(21), f),
		} {
			assert.Equal(t, result.Left[int](err), m(t.Context()))
		}
	})

	t.Run("skips the side effect on upstream failure", func(t *testing.T) {
		err := errors.New("upstream")
		calls := 0
		f := iorKCounter(&calls, nil)
		for _, m := range []ReaderResult[int]{
			F.Pipe1(Left[int](err), ChainFirstIOResultK(f)),
			MonadTapIOResultK(Left[int](err), f),
		} {
			assert.Equal(t, result.Left[int](err), m(t.Context()))
		}
		assert.Equal(t, 0, calls)
	})
}

// positiveResult fails for numbers that are not positive.
func positiveResult(n int) result.Result[int] {
	if n > 0 {
		return result.Of(n)
	}
	return result.Left[int](errors.New("not positive"))
}

func TestFromResult(t *testing.T) {
	assert.Equal(t, result.Of(1), FromResult(result.Of(1))(t.Context()))
	assert.Equal(t, FromEither(positiveResult(0))(t.Context()), FromResult(positiveResult(0))(t.Context()))
}

func TestMonadChainResultK(t *testing.T) {
	assert.Equal(t, result.Of(5), MonadChainResultK(Of(5), positiveResult)(t.Context()))
	assert.Equal(t, result.Of(5), ChainResultK(positiveResult)(Of(5))(t.Context()))
	assert.Equal(t, ChainEitherK(positiveResult)(Of(-5))(t.Context()), MonadChainResultK(Of(-5), positiveResult)(t.Context()))
}

func TestChainFirstEitherK(t *testing.T) {
	all := func(start ReaderResult[int]) []ReaderResult[int] {
		return []ReaderResult[int]{
			ChainFirstEitherK(positiveResult)(start),
			TapEitherK(positiveResult)(start),
			MonadChainFirstEitherK(start, positiveResult),
			MonadTapEitherK(start, positiveResult),
			ChainFirstResultK(positiveResult)(start),
			TapResultK(positiveResult)(start),
			MonadChainFirstResultK(start, positiveResult),
			MonadTapResultK(start, positiveResult),
		}
	}
	t.Run("keeps the value", func(t *testing.T) {
		for _, m := range all(Of(7)) {
			assert.Equal(t, result.Of(7), m(t.Context()))
		}
	})
	t.Run("failing check replaces the value", func(t *testing.T) {
		for _, m := range all(Of(-7)) {
			assert.Equal(t, result.Left[int](errors.New("not positive")), m(t.Context()))
		}
	})
	t.Run("upstream error is preserved", func(t *testing.T) {
		err := errors.New("upstream")
		for _, m := range all(Left[int](err)) {
			assert.Equal(t, result.Left[int](err), m(t.Context()))
		}
	})
}
