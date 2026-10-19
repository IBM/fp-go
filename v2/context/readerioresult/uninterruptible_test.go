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

type uninterruptibleKey struct{}

// contextAwareRelease returns a release action that fails if its context is cancelled
// and otherwise records that it ran.
func contextAwareRelease(released *bool) ReaderIOResult[Void] {
	return F.Pipe1(
		FromIO(func() Void { *released = true; return F.VOID }),
		WithContext[Void],
	)
}

func TestUninterruptible_Success(t *testing.T) {
	t.Run("runs with a cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		var released bool
		assert.Equal(t, result.Of(F.VOID), Uninterruptible(contextAwareRelease(&released))(ctx)())
		assert.True(t, released)
	})

	t.Run("drops the deadline", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(t.Context(), time.Hour)
		defer cancel()

		hasDeadline := func(ctx context.Context) IOResult[bool] {
			_, ok := ctx.Deadline()
			return Of(ok)(ctx)
		}
		assert.Equal(t, result.Of(false), Uninterruptible(hasDeadline)(ctx)())
	})

	t.Run("keeps context values", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.WithValue(t.Context(), uninterruptibleKey{}, "req-1"))
		cancel()

		value := func(ctx context.Context) IOResult[any] {
			return Of(ctx.Value(uninterruptibleKey{}))(ctx)
		}
		assert.Equal(t, result.Of[any]("req-1"), Uninterruptible(value)(ctx)())
	})
}

func TestUninterruptible_Failure(t *testing.T) {
	boom := errors.New("boom")
	assert.Equal(t, result.Left[int](boom), Uninterruptible(Left[int](boom))(t.Context())())
}

func TestBracket_CancelledDuringUse(t *testing.T) {
	t.Run("runs a context-aware release", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		defer cancel()

		var released bool
		use := func(int) ReaderIOResult[int] {
			return FromIO(func() int { cancel(); return 2 })
		}
		release := func(int, Result[int]) ReaderIOResult[Void] {
			return contextAwareRelease(&released)
		}

		Bracket(Of(1), use, release)(ctx)()
		assert.True(t, released)
	})

	t.Run("does not run release when the context is cancelled before acquire", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		var released bool
		res := Bracket(
			WithContext(Of(1)),
			Of[int],
			func(int, Result[int]) ReaderIOResult[Void] { return contextAwareRelease(&released) },
		)(ctx)()

		assert.True(t, result.IsLeft(res))
		assert.False(t, released)
	})
}

func TestWithResource_CancelledDuringUse(t *testing.T) {
	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	var released bool
	withRes := WithResource[int](Of(1), func(int) ReaderIOResult[Void] { return contextAwareRelease(&released) })

	withRes(func(int) ReaderIOResult[int] {
		return FromIO(func() int { cancel(); return 2 })
	})(ctx)()
	assert.True(t, released)
}

func BenchmarkBracket(b *testing.B) {
	release := func(int, Result[int]) ReaderIOResult[Void] { return Of(F.VOID) }
	eff := Bracket(Of(1), Of[int], release)
	ctx := b.Context()
	for b.Loop() {
		eff(ctx)()
	}
}
