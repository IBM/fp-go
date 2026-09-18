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
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

type valueTestKey string

const (
	userKey    valueTestKey = "user"
	requestKey valueTestKey = "request"
)

func TestAskValue(t *testing.T) {
	t.Run("returns Some when the key holds a value of the expected type", func(t *testing.T) {
		ctx := context.WithValue(t.Context(), userKey, "Alice")
		assert.Equal(t, E.Of[error](O.Some("Alice")), AskValue[string](userKey)(ctx))
	})

	t.Run("returns None when the key is absent", func(t *testing.T) {
		assert.Equal(t, E.Of[error](O.None[string]()), AskValue[string](userKey)(t.Context()))
	})

	t.Run("returns None when the stored value has a different type", func(t *testing.T) {
		ctx := context.WithValue(t.Context(), userKey, 42)
		assert.Equal(t, E.Of[error](O.None[string]()), AskValue[string](userKey)(ctx))
	})
}

func TestWithValue(t *testing.T) {
	askUser := F.Pipe1(
		AskValue[string](userKey),
		Map(O.GetOrElse(F.Constant("anonymous"))),
	)

	t.Run("makes the value visible to the wrapped computation", func(t *testing.T) {
		res := F.Pipe1(askUser, WithValue[string](userKey, "Alice"))
		assert.Equal(t, E.Of[error]("Alice"), res(t.Context()))
	})

	t.Run("does not modify the caller's context", func(t *testing.T) {
		outer := t.Context()
		F.Pipe1(askUser, WithValue[string](userKey, "Alice"))(outer)
		assert.Nil(t, outer.Value(userKey))
		assert.Equal(t, E.Of[error]("anonymous"), askUser(outer))
	})

	t.Run("inner values shadow outer values", func(t *testing.T) {
		res := F.Pipe2(
			askUser,
			WithValue[string](userKey, "inner"),
			WithValue[string](userKey, "outer"),
		)
		assert.Equal(t, E.Of[error]("inner"), res(t.Context()))
	})

	t.Run("stacks values under different keys", func(t *testing.T) {
		both := F.Pipe1(
			Ask(),
			Map(func(ctx context.Context) string {
				return fmt.Sprintf("%v/%v", ctx.Value(userKey), ctx.Value(requestKey))
			}),
		)
		res := F.Pipe2(
			both,
			WithValue[string](userKey, "Alice"),
			WithValue[string](requestKey, 7),
		)
		assert.Equal(t, E.Of[error]("Alice/7"), res(t.Context()))
	})

	t.Run("propagates errors of the wrapped computation", func(t *testing.T) {
		err := errors.New("boom")
		res := F.Pipe1(Left[int](err), WithValue[int](userKey, "Alice"))
		assert.Equal(t, E.Left[int](err), res(t.Context()))
	})
}

// waitForDone blocks until the context is done or the fallback elapses
func waitForDone(fallback time.Duration) ReaderResult[string] {
	return func(ctx context.Context) Result[string] {
		select {
		case <-ctx.Done():
			return E.Left[string](ctx.Err())
		case <-time.After(fallback):
			return E.Of[error]("completed")
		}
	}
}

func TestWithTimeout(t *testing.T) {
	t.Run("cancels a computation that exceeds the timeout", func(t *testing.T) {
		res := F.Pipe1(waitForDone(time.Second), WithTimeout[string](10*time.Millisecond))
		assert.Equal(t, E.Left[string](context.DeadlineExceeded), res(t.Context()))
	})

	t.Run("passes through a computation that finishes in time", func(t *testing.T) {
		res := F.Pipe1(waitForDone(time.Millisecond), WithTimeout[string](time.Second))
		assert.Equal(t, E.Of[error]("completed"), res(t.Context()))
	})

	t.Run("sets a deadline on the derived context", func(t *testing.T) {
		var deadline time.Time
		var ok bool
		probe := func(ctx context.Context) Result[int] {
			deadline, ok = ctx.Deadline()
			return E.Of[error](1)
		}
		before := time.Now()
		F.Pipe1(probe, WithTimeout[int](time.Minute))(t.Context())

		assert.True(t, ok)
		assert.WithinDuration(t, before.Add(time.Minute), deadline, time.Second)
		_, parentHasDeadline := t.Context().Deadline()
		assert.False(t, parentHasDeadline)
	})

	t.Run("releases the derived context after completion", func(t *testing.T) {
		var seen context.Context
		probe := func(ctx context.Context) Result[int] {
			seen = ctx
			return E.Of[error](1)
		}
		F.Pipe1(probe, WithTimeout[int](time.Minute))(t.Context())
		assert.ErrorIs(t, seen.Err(), context.Canceled)
	})
}

func TestWithDeadline(t *testing.T) {
	t.Run("cancels a computation that runs past the deadline", func(t *testing.T) {
		res := F.Pipe1(waitForDone(time.Second), WithDeadline[string](time.Now().Add(10*time.Millisecond)))
		assert.Equal(t, E.Left[string](context.DeadlineExceeded), res(t.Context()))
	})

	t.Run("passes through a computation that finishes in time", func(t *testing.T) {
		res := F.Pipe1(waitForDone(time.Millisecond), WithDeadline[string](time.Now().Add(time.Second)))
		assert.Equal(t, E.Of[error]("completed"), res(t.Context()))
	})

	t.Run("a deadline in the past cancels immediately", func(t *testing.T) {
		res := F.Pipe1(waitForDone(time.Second), WithDeadline[string](time.Now().Add(-time.Second)))
		assert.Equal(t, E.Left[string](context.DeadlineExceeded), res(t.Context()))
	})

	t.Run("an earlier parent deadline takes precedence", func(t *testing.T) {
		parent, cancel := context.WithTimeout(t.Context(), 10*time.Millisecond)
		defer cancel()
		res := F.Pipe1(waitForDone(time.Second), WithDeadline[string](time.Now().Add(time.Hour)))
		assert.Equal(t, E.Left[string](context.DeadlineExceeded), res(parent))
	})
}
