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
	R "github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

func TestChainFirstLeft_Success(t *testing.T) {
	t.Run("observes the error and keeps it", func(t *testing.T) {
		var seen error
		observe := func(err error) (F.Void, error) { seen = err; return Of(F.VOID) }
		assert.Equal(t, (R.Left[int](errSentinel)), run(ChainFirstLeft[int](observe)(Left[int](errSentinel))))
		assert.Equal(t, errSentinel, seen)
	})

	t.Run("a failing side effect does not replace the error", func(t *testing.T) {
		fail := func(error) (F.Void, error) { return Left[F.Void](errors.New("side effect")) }
		assert.Equal(t, (R.Left[int](errSentinel)), run(ChainFirstLeft[int](fail)(Left[int](errSentinel))))
	})

	t.Run("TapLeft is the same as ChainFirstLeft", func(t *testing.T) {
		var seen error
		observe := func(err error) (F.Void, error) { seen = err; return Of(F.VOID) }
		assert.Equal(t, (R.Left[int](errSentinel)), run(TapLeft[int](observe)(Left[int](errSentinel))))
		assert.Equal(t, errSentinel, seen)
	})
}

func TestChainFirstLeft_Failure(t *testing.T) {
	called := false
	observe := func(error) (F.Void, error) { called = true; return Of(F.VOID) }
	assert.Equal(t, (R.Of(42)), run(ChainFirstLeft[int](observe)(Of(42))))
	assert.False(t, called)
}
