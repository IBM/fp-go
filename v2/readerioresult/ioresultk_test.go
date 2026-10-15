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
	"fmt"
	"testing"

	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

type iorKConfig struct{}

// iorKCounter counts its executions and doubles its input or fails with err, if err is not nil.
func iorKCounter(calls *int, err error) ioresult.Kleisli[int, int] {
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

func TestMonadChainFirstIOResultK_Success(t *testing.T) {
	calls := 0
	f := iorKCounter(&calls, nil)
	for _, m := range []ReaderIOResult[iorKConfig, int]{
		MonadChainFirstIOEitherK(Of[iorKConfig](21), f),
		MonadTapIOEitherK(Of[iorKConfig](21), f),
		MonadChainFirstIOResultK(Of[iorKConfig](21), f),
		MonadTapIOResultK(Of[iorKConfig](21), f),
	} {
		assert.Equal(t, result.Of(21), m(iorKConfig{})())
	}
	assert.Equal(t, 4, calls)
}

func TestMonadChainFirstIOResultK_Failure(t *testing.T) {
	t.Run("error of the side effect replaces the value", func(t *testing.T) {
		err := errors.New("io failed")
		f := iorKCounter(new(int), err)
		for _, m := range []ReaderIOResult[iorKConfig, int]{
			MonadChainFirstIOEitherK(Of[iorKConfig](21), f),
			MonadTapIOEitherK(Of[iorKConfig](21), f),
			MonadChainFirstIOResultK(Of[iorKConfig](21), f),
			MonadTapIOResultK(Of[iorKConfig](21), f),
		} {
			assert.Equal(t, result.Left[int](err), m(iorKConfig{})())
		}
	})

	t.Run("skips the side effect on upstream failure", func(t *testing.T) {
		err := errors.New("upstream")
		calls := 0
		m := MonadTapIOResultK(Left[iorKConfig, int](err), iorKCounter(&calls, nil))
		assert.Equal(t, result.Left[int](err), m(iorKConfig{})())
		assert.Equal(t, 0, calls)
	})
}

// ExampleMonadChainFirstIOResultK runs a fallible side effect and keeps the original value.
func ExampleMonadChainFirstIOResultK() {
	fmt.Println(result.Unwrap(MonadChainFirstIOResultK(Of[iorKConfig](21), iorKCounter(new(int), nil))(iorKConfig{})()))
	// Output: 21 <nil>
}

// ExampleMonadTapIOResultK is an alias for MonadChainFirstIOResultK; a failing
// side effect replaces the value.
func ExampleMonadTapIOResultK() {
	fmt.Println(result.Unwrap(MonadTapIOResultK(Of[iorKConfig](21), iorKCounter(new(int), errors.New("io failed")))(iorKConfig{})()))
	// Output: 0 io failed
}

// ExampleMonadChainFirstIOEitherK is an alias for MonadChainFirstIOResultK.
func ExampleMonadChainFirstIOEitherK() {
	fmt.Println(result.Unwrap(MonadChainFirstIOEitherK(Of[iorKConfig](21), iorKCounter(new(int), nil))(iorKConfig{})()))
	// Output: 21 <nil>
}

// ExampleMonadTapIOEitherK is an alias for MonadTapIOResultK.
func ExampleMonadTapIOEitherK() {
	fmt.Println(result.Unwrap(MonadTapIOEitherK(Of[iorKConfig](21), iorKCounter(new(int), nil))(iorKConfig{})()))
	// Output: 21 <nil>
}
