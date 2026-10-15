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

package readerreaderioeither

import (
	"errors"
	"fmt"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	"github.com/stretchr/testify/assert"
)

type ioeKOuter struct{}
type ioeKInner struct{}

type ioeKRRIOE = ReaderReaderIOEither[ioeKOuter, ioeKInner, error, int]

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

func ioeKRun(m ioeKRRIOE) E.Either[error, int] {
	return m(ioeKOuter{})(ioeKInner{})()
}

func TestChainFirstIOEitherK_Success(t *testing.T) {
	calls := 0
	f := ioeKCounter(&calls, nil)
	for _, m := range []ioeKRRIOE{
		F.Pipe1(Of[ioeKOuter, ioeKInner, error](21), ChainFirstIOEitherK[ioeKOuter, ioeKInner](f)),
		F.Pipe1(Of[ioeKOuter, ioeKInner, error](21), TapIOEitherK[ioeKOuter, ioeKInner](f)),
		MonadChainFirstIOEitherK(Of[ioeKOuter, ioeKInner, error](21), f),
		MonadTapIOEitherK(Of[ioeKOuter, ioeKInner, error](21), f),
	} {
		assert.Equal(t, E.Right[error](21), ioeKRun(m))
	}
	assert.Equal(t, 4, calls)
}

func TestChainFirstIOEitherK_Failure(t *testing.T) {
	t.Run("error of the side effect replaces the value", func(t *testing.T) {
		err := errors.New("io failed")
		f := ioeKCounter(new(int), err)
		for _, m := range []ioeKRRIOE{
			F.Pipe1(Of[ioeKOuter, ioeKInner, error](21), ChainFirstIOEitherK[ioeKOuter, ioeKInner](f)),
			F.Pipe1(Of[ioeKOuter, ioeKInner, error](21), TapIOEitherK[ioeKOuter, ioeKInner](f)),
			MonadChainFirstIOEitherK(Of[ioeKOuter, ioeKInner, error](21), f),
			MonadTapIOEitherK(Of[ioeKOuter, ioeKInner, error](21), f),
		} {
			assert.Equal(t, E.Left[int](err), ioeKRun(m))
		}
	})

	t.Run("skips the side effect on upstream failure", func(t *testing.T) {
		err := errors.New("upstream")
		calls := 0
		f := ioeKCounter(&calls, nil)
		for _, m := range []ioeKRRIOE{
			F.Pipe1(Left[ioeKOuter, ioeKInner, int](err), ChainFirstIOEitherK[ioeKOuter, ioeKInner](f)),
			MonadTapIOEitherK(Left[ioeKOuter, ioeKInner, int](err), f),
		} {
			assert.Equal(t, E.Left[int](err), ioeKRun(m))
		}
		assert.Equal(t, 0, calls)
	})
}

// ExampleChainFirstIOEitherK runs a fallible side effect and keeps the original value.
func ExampleChainFirstIOEitherK() {
	fmt.Println(ioeKRun(F.Pipe1(Of[ioeKOuter, ioeKInner, error](21), ChainFirstIOEitherK[ioeKOuter, ioeKInner](ioeKCounter(new(int), nil)))))
	// Output: Right[int](21)
}

// ExampleTapIOEitherK is an alias for ChainFirstIOEitherK; a failing side
// effect replaces the value.
func ExampleTapIOEitherK() {
	fmt.Println(ioeKRun(F.Pipe1(Of[ioeKOuter, ioeKInner, error](21), TapIOEitherK[ioeKOuter, ioeKInner](ioeKCounter(new(int), errors.New("io failed"))))))
	// Output: Left[*errors.errorString](io failed)
}

// ExampleMonadChainFirstIOEitherK is the non curried version of ChainFirstIOEitherK.
func ExampleMonadChainFirstIOEitherK() {
	fmt.Println(ioeKRun(MonadChainFirstIOEitherK(Of[ioeKOuter, ioeKInner, error](21), ioeKCounter(new(int), nil))))
	// Output: Right[int](21)
}

// ExampleMonadTapIOEitherK is the non curried version of TapIOEitherK.
func ExampleMonadTapIOEitherK() {
	fmt.Println(ioeKRun(MonadTapIOEitherK(Of[ioeKOuter, ioeKInner, error](21), ioeKCounter(new(int), nil))))
	// Output: Right[int](21)
}
