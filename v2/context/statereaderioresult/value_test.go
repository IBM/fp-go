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

package statereaderioresult

import (
	"context"
	"testing"
	"time"

	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

type valueTestKey string

const valueUserKey valueTestKey = "user"

// countingProbe increments the state and returns the context it ran with
func countingProbe(seen *context.Context) StateReaderIOResult[int, string] {
	return func(s int) ReaderIOResult[Pair[int, string]] {
		return func(ctx context.Context) IOResult[Pair[int, string]] {
			return func() Result[Pair[int, string]] {
				*seen = ctx
				user, _ := ctx.Value(valueUserKey).(string)
				return result.Of(pair.MakePair(s+1, user))
			}
		}
	}
}

// waitForDone blocks until the context is done or the fallback elapses
func waitForDone(fallback time.Duration) StateReaderIOResult[int, string] {
	return func(s int) ReaderIOResult[Pair[int, string]] {
		return func(ctx context.Context) IOResult[Pair[int, string]] {
			return func() Result[Pair[int, string]] {
				select {
				case <-ctx.Done():
					return result.Left[Pair[int, string]](ctx.Err())
				case <-time.After(fallback):
					return result.Of(pair.MakePair(s, "completed"))
				}
			}
		}
	}
}

func TestAskValue(t *testing.T) {
	t.Run("returns Some and preserves the state", func(t *testing.T) {
		ctx := context.WithValue(t.Context(), valueUserKey, "Alice")
		res := AskValue[int, string](valueUserKey)(7)(ctx)()
		assert.Equal(t, result.Of(pair.MakePair(7, O.Some("Alice"))), res)
	})

	t.Run("returns None when the key is absent", func(t *testing.T) {
		res := AskValue[int, string](valueUserKey)(7)(t.Context())()
		assert.Equal(t, result.Of(pair.MakePair(7, O.None[string]())), res)
	})

	t.Run("returns None when the stored value has a different type", func(t *testing.T) {
		ctx := context.WithValue(t.Context(), valueUserKey, 42)
		res := AskValue[int, string](valueUserKey)(7)(ctx)()
		assert.Equal(t, result.Of(pair.MakePair(7, O.None[string]())), res)
	})
}

func TestWithValue(t *testing.T) {
	t.Run("makes the value visible and threads the state", func(t *testing.T) {
		var seen context.Context
		res := F.Pipe1(
			countingProbe(&seen),
			WithValue[int, string](valueUserKey, "Alice"),
		)(1)(t.Context())()

		assert.Equal(t, result.Of(pair.MakePair(2, "Alice")), res)
		assert.Nil(t, t.Context().Value(valueUserKey))
		assert.NoError(t, seen.Err())
	})

	t.Run("round-trips with AskValue", func(t *testing.T) {
		res := F.Pipe1(
			AskValue[int, string](valueUserKey),
			WithValue[int, O.Option[string]](valueUserKey, "Alice"),
		)(0)(t.Context())()
		assert.Equal(t, result.Of(pair.MakePair(0, O.Some("Alice"))), res)
	})
}

func TestWithTimeout(t *testing.T) {
	t.Run("cancels a computation that exceeds the timeout", func(t *testing.T) {
		res := F.Pipe1(waitForDone(time.Second), WithTimeout[int, string](10*time.Millisecond))(0)(t.Context())()
		assert.Equal(t, result.Left[Pair[int, string]](context.DeadlineExceeded), res)
	})

	t.Run("passes through a computation that finishes in time", func(t *testing.T) {
		res := F.Pipe1(waitForDone(time.Millisecond), WithTimeout[int, string](time.Second))(3)(t.Context())()
		assert.Equal(t, result.Of(pair.MakePair(3, "completed")), res)
	})

	t.Run("releases the derived context after completion", func(t *testing.T) {
		var seen context.Context
		F.Pipe1(countingProbe(&seen), WithTimeout[int, string](time.Minute))(0)(t.Context())()
		assert.ErrorIs(t, seen.Err(), context.Canceled)
	})
}

func TestWithDeadline(t *testing.T) {
	t.Run("cancels a computation that runs past the deadline", func(t *testing.T) {
		res := F.Pipe1(waitForDone(time.Second), WithDeadline[int, string](time.Now().Add(10*time.Millisecond)))(0)(t.Context())()
		assert.Equal(t, result.Left[Pair[int, string]](context.DeadlineExceeded), res)
	})

	t.Run("sets the deadline on the derived context", func(t *testing.T) {
		var seen context.Context
		deadline := time.Now().Add(time.Hour)
		F.Pipe1(countingProbe(&seen), WithDeadline[int, string](deadline))(0)(t.Context())()

		actual, ok := seen.Deadline()
		assert.True(t, ok)
		assert.True(t, deadline.Equal(actual))
	})
}
