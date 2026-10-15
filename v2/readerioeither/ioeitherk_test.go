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

package readerioeither

import (
	"errors"
	"fmt"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	"github.com/stretchr/testify/assert"
)

type ioeKConfig struct{}

// ioeKCounter counts its executions and doubles its input or fails with err, if err is not nil.
func ioeKCounter(calls *int, err error) IOE.Kleisli[error, int, int] {
	return func(n int) IOEither[error, int] {
		return func() E.Either[error, int] {
			*calls++
			if err != nil {
				return E.Left[int](err)
			}
			return E.Right[error](n * 2)
		}
	}
}

func TestMonadChainFirstIOEitherK_Success(t *testing.T) {
	calls := 0
	f := ioeKCounter(&calls, nil)
	for _, m := range []ReaderIOEither[ioeKConfig, error, int]{
		MonadChainFirstIOEitherK(Of[ioeKConfig, error](21), f),
		MonadTapIOEitherK(Of[ioeKConfig, error](21), f),
	} {
		assert.Equal(t, E.Right[error](21), m(ioeKConfig{})())
	}
	assert.Equal(t, 2, calls)
}

func TestMonadChainFirstIOEitherK_Failure(t *testing.T) {
	t.Run("error of the side effect replaces the value", func(t *testing.T) {
		err := errors.New("io failed")
		f := ioeKCounter(new(int), err)
		for _, m := range []ReaderIOEither[ioeKConfig, error, int]{
			MonadChainFirstIOEitherK(Of[ioeKConfig, error](21), f),
			MonadTapIOEitherK(Of[ioeKConfig, error](21), f),
		} {
			assert.Equal(t, E.Left[int](err), m(ioeKConfig{})())
		}
	})

	t.Run("skips the side effect on upstream failure", func(t *testing.T) {
		err := errors.New("upstream")
		calls := 0
		m := MonadTapIOEitherK(Left[ioeKConfig, int](err), ioeKCounter(&calls, nil))
		assert.Equal(t, E.Left[int](err), m(ioeKConfig{})())
		assert.Equal(t, 0, calls)
	})
}

func TestMonadChainFirstIOEitherK_AgreesWithCurried(t *testing.T) {
	err := errors.New("io failed")
	for _, f := range []IOE.Kleisli[error, int, int]{ioeKCounter(new(int), nil), ioeKCounter(new(int), err)} {
		for _, start := range []ReaderIOEither[ioeKConfig, error, int]{Of[ioeKConfig, error](21), Left[ioeKConfig, int](err)} {
			assert.Equal(t, ChainFirstIOEitherK[ioeKConfig](f)(start)(ioeKConfig{})(), MonadChainFirstIOEitherK(start, f)(ioeKConfig{})())
			assert.Equal(t, TapIOEitherK[ioeKConfig](f)(start)(ioeKConfig{})(), MonadTapIOEitherK(start, f)(ioeKConfig{})())
		}
	}
}

// ExampleMonadChainFirstIOEitherK runs a fallible side effect and keeps the original value.
func ExampleMonadChainFirstIOEitherK() {
	fmt.Println(MonadChainFirstIOEitherK(Of[ioeKConfig, error](21), ioeKCounter(new(int), nil))(ioeKConfig{})())
	// Output: Right[int](21)
}

// ExampleMonadTapIOEitherK is an alias for MonadChainFirstIOEitherK; a failing
// side effect replaces the value.
func ExampleMonadTapIOEitherK() {
	fmt.Println(MonadTapIOEitherK(Of[ioeKConfig, error](21), ioeKCounter(new(int), errors.New("io failed")))(ioeKConfig{})())
	// Output: Left[*errors.errorString](io failed)
}
