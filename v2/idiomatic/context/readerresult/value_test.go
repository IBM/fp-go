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

	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

type valueTestKey string

const (
	valueUserKey    valueTestKey = "user"
	valueRequestKey valueTestKey = "request"
)

func TestAskValue(t *testing.T) {
	t.Run("returns Some when the key holds a value of the expected type", func(t *testing.T) {
		ctx := context.WithValue(t.Context(), valueUserKey, "Alice")
		v, err := AskValue[string](valueUserKey)(ctx)
		assert.NoError(t, err)
		assert.Equal(t, O.Some("Alice"), v)
	})

	t.Run("returns None when the key is absent", func(t *testing.T) {
		v, err := AskValue[string](valueUserKey)(t.Context())
		assert.NoError(t, err)
		assert.Equal(t, O.None[string](), v)
	})

	t.Run("returns None when the stored value has a different type", func(t *testing.T) {
		ctx := context.WithValue(t.Context(), valueUserKey, 42)
		v, err := AskValue[string](valueUserKey)(ctx)
		assert.NoError(t, err)
		assert.Equal(t, O.None[string](), v)
	})
}

func TestWithValue(t *testing.T) {
	askUser := F.Pipe1(
		AskValue[string](valueUserKey),
		Map(O.GetOrElse(F.Constant("anonymous"))),
	)

	t.Run("makes the value visible to the wrapped computation", func(t *testing.T) {
		v, err := F.Pipe1(askUser, WithValue[string](valueUserKey, "Alice"))(t.Context())
		assert.NoError(t, err)
		assert.Equal(t, "Alice", v)
	})

	t.Run("does not modify the caller's context", func(t *testing.T) {
		outer := t.Context()
		_, _ = F.Pipe1(askUser, WithValue[string](valueUserKey, "Alice"))(outer)
		assert.Nil(t, outer.Value(valueUserKey))

		v, err := askUser(outer)
		assert.NoError(t, err)
		assert.Equal(t, "anonymous", v)
	})

	t.Run("inner values shadow outer values", func(t *testing.T) {
		v, err := F.Pipe2(
			askUser,
			WithValue[string](valueUserKey, "inner"),
			WithValue[string](valueUserKey, "outer"),
		)(t.Context())
		assert.NoError(t, err)
		assert.Equal(t, "inner", v)
	})

	t.Run("stacks values under different keys", func(t *testing.T) {
		both := F.Pipe1(
			Ask(),
			Map(func(ctx context.Context) string {
				return fmt.Sprintf("%v/%v", ctx.Value(valueUserKey), ctx.Value(valueRequestKey))
			}),
		)
		v, err := F.Pipe2(
			both,
			WithValue[string](valueUserKey, "Alice"),
			WithValue[string](valueRequestKey, 7),
		)(t.Context())
		assert.NoError(t, err)
		assert.Equal(t, "Alice/7", v)
	})

	t.Run("propagates errors of the wrapped computation", func(t *testing.T) {
		boom := errors.New("boom")
		_, err := F.Pipe1(Left[int](boom), WithValue[int](valueUserKey, "Alice"))(t.Context())
		assert.ErrorIs(t, err, boom)
	})

	t.Run("short-circuits on an already cancelled context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		executed := false
		probe := func(context.Context) (int, error) {
			executed = true
			return 1, nil
		}
		_, err := F.Pipe1(probe, WithValue[int](valueUserKey, "Alice"))(ctx)
		assert.ErrorIs(t, err, context.Canceled)
		assert.False(t, executed)
	})
}
