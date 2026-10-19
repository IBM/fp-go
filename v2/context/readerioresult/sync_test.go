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
	"context"
	"errors"
	"testing"
	"time"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

func TestWithSemaphore_Success(t *testing.T) {
	t.Run("limits the number of concurrent computations", func(t *testing.T) {
		probe := &concurrencyProbe{}
		sem := make(chan Void, 2)
		guarded := F.Flow2(probe.track(5*time.Millisecond), WithSemaphore[int](sem))

		res := TraverseArrayPar(guarded)(make([]int, 10))(t.Context())()

		assert.True(t, result.IsRight(res))
		assert.LessOrEqual(t, probe.peak.Load(), int32(2))
		assert.Empty(t, sem)
	})

	t.Run("passes the value through", func(t *testing.T) {
		res := WithSemaphore[int](make(chan Void, 1))(Of(42))(t.Context())()
		assert.Equal(t, result.Of(42), res)
	})
}

func TestWithSemaphore_Failure(t *testing.T) {
	t.Run("returns the slot when the computation fails", func(t *testing.T) {
		boom := errors.New("boom")
		sem := make(chan Void, 1)

		res := WithSemaphore[int](sem)(Left[int](boom))(t.Context())()

		assert.Equal(t, result.Left[int](boom), res)
		assert.Empty(t, sem)
	})

	t.Run("fails when the context is cancelled while waiting", func(t *testing.T) {
		sem := make(chan Void, 1)
		sem <- F.VOID // the only slot is taken
		ctx, cancel := context.WithTimeout(t.Context(), 10*time.Millisecond)
		defer cancel()
		called := false

		res := WithSemaphore[int](sem)(FromLazy(func() int {
			called = true
			return 1
		}))(ctx)()

		assert.Equal(t, result.Left[int](context.DeadlineExceeded), res)
		assert.False(t, called)
		assert.Len(t, sem, 1)
	})
}

func TestWithSemaphore_EdgeCases(t *testing.T) {
	t.Run("already cancelled context does not take a slot", func(t *testing.T) {
		sem := make(chan Void, 1)
		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		res := WithSemaphore[int](sem)(Of(1))(ctx)()

		assert.Equal(t, result.Left[int](context.Canceled), res)
		assert.Empty(t, sem)
	})
}
