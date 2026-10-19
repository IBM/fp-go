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
	"errors"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

func TestChainFirstLeft_Success(t *testing.T) {
	t.Run("observes the error and keeps it", func(t *testing.T) {
		var seen error
		observe := func(err error) Result[F.Void] { seen = err; return Of(F.VOID) }
		assert.Equal(t, Left[int](errSentinel), ChainFirstLeft[int](observe)(Left[int](errSentinel)))
		assert.Equal(t, errSentinel, seen)
	})

	t.Run("a failing side effect does not replace the error", func(t *testing.T) {
		fail := func(error) Result[F.Void] { return Left[F.Void](errors.New("side effect")) }
		assert.Equal(t, Left[int](errSentinel), ChainFirstLeft[int](fail)(Left[int](errSentinel)))
	})

	t.Run("the aliases and monadic forms behave the same", func(t *testing.T) {
		count := 0
		observe := func(error) Result[F.Void] { count++; return Of(F.VOID) }
		assert.Equal(t, Left[int](errSentinel), TapLeft[int](observe)(Left[int](errSentinel)))
		assert.Equal(t, Left[int](errSentinel), MonadChainFirstLeft(Left[int](errSentinel), observe))
		assert.Equal(t, Left[int](errSentinel), MonadTapLeft(Left[int](errSentinel), observe))
		assert.Equal(t, 3, count)
	})
}

func TestChainFirstLeft_Failure(t *testing.T) {
	called := false
	observe := func(error) Result[F.Void] { called = true; return Of(F.VOID) }
	assert.Equal(t, Of(42), ChainFirstLeft[int](observe)(Of(42)))
	assert.False(t, called)
}
