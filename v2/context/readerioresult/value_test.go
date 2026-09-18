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
	"fmt"
	"testing"

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
		assert.Equal(t, E.Of[error](O.Some("Alice")), AskValue[string](userKey)(ctx)())
	})

	t.Run("returns None when the key is absent", func(t *testing.T) {
		assert.Equal(t, E.Of[error](O.None[string]()), AskValue[string](userKey)(t.Context())())
	})

	t.Run("returns None when the stored value has a different type", func(t *testing.T) {
		ctx := context.WithValue(t.Context(), userKey, 42)
		assert.Equal(t, E.Of[error](O.None[string]()), AskValue[string](userKey)(ctx)())
	})

	t.Run("distinguishes keys of different types with the same underlying value", func(t *testing.T) {
		ctx := context.WithValue(t.Context(), "user", "plain")
		assert.Equal(t, E.Of[error](O.None[string]()), AskValue[string](userKey)(ctx)())
	})
}

func TestWithValue(t *testing.T) {
	askUser := F.Pipe1(
		AskValue[string](userKey),
		Map(O.GetOrElse(F.Constant("anonymous"))),
	)

	t.Run("makes the value visible to the wrapped computation", func(t *testing.T) {
		res := F.Pipe1(askUser, WithValue[string](userKey, "Alice"))
		assert.Equal(t, E.Of[error]("Alice"), res(t.Context())())
	})

	t.Run("does not modify the caller's context", func(t *testing.T) {
		var seen context.Context
		capture := func(ctx context.Context) IOResult[string] {
			seen = ctx
			return askUser(ctx)
		}
		outer := t.Context()
		F.Pipe1(capture, WithValue[string](userKey, "Alice"))(outer)()

		assert.NotSame(t, outer, seen)
		assert.Nil(t, outer.Value(userKey))
		assert.Equal(t, "Alice", seen.Value(userKey))
	})

	t.Run("inner values shadow outer values", func(t *testing.T) {
		res := F.Pipe2(
			askUser,
			WithValue[string](userKey, "inner"),
			WithValue[string](userKey, "outer"),
		)
		assert.Equal(t, E.Of[error]("inner"), res(t.Context())())
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
		assert.Equal(t, E.Of[error]("Alice/7"), res(t.Context())())
	})

	t.Run("does not cancel the derived context", func(t *testing.T) {
		var seen context.Context
		capture := func(ctx context.Context) IOResult[int] {
			seen = ctx
			return Of(1)(ctx)
		}
		F.Pipe1(capture, WithValue[int](userKey, "Alice"))(t.Context())()
		assert.NoError(t, seen.Err())
	})

	t.Run("propagates errors of the wrapped computation", func(t *testing.T) {
		err := errors.New("boom")
		res := F.Pipe1(Left[int](err), WithValue[int](userKey, "Alice"))
		assert.Equal(t, E.Left[int](err), res(t.Context())())
	})
}
