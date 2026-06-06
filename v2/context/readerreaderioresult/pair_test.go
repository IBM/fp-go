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

package readerreaderioresult

import (
	"context"
	"errors"
	"testing"

	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

func TestPaired(t *testing.T) {
	t.Run("applies outer environment from pair tail", func(t *testing.T) {
		cfg := AppConfig{DatabaseURL: "postgres://test", LogLevel: "debug"}
		ctx := t.Context()

		var f ReaderReaderIOResult[AppConfig, string] = func(c AppConfig) func(context.Context) IOResult[string] {
			return func(context.Context) IOResult[string] {
				return func() Result[string] {
					return result.Of(c.DatabaseURL)
				}
			}
		}

		p := pair.MakePair[context.Context, AppConfig](ctx, cfg)
		outcome := Paired(f)(p)()

		assert.Equal(t, result.Of("postgres://test"), outcome)
	})

	t.Run("passes context from pair head to inner reader", func(t *testing.T) {
		cfg := AppConfig{}
		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		var capturedCtx context.Context
		var f ReaderReaderIOResult[AppConfig, string] = func(AppConfig) func(context.Context) IOResult[string] {
			return func(c context.Context) IOResult[string] {
				capturedCtx = c
				return func() Result[string] {
					return result.Of("done")
				}
			}
		}

		p := pair.MakePair[context.Context, AppConfig](ctx, cfg)
		_ = Paired(f)(p)()

		assert.ErrorIs(t, capturedCtx.Err(), context.Canceled)
	})

	t.Run("propagates error from computation", func(t *testing.T) {
		cfg := AppConfig{}
		testErr := errors.New("computation failed")

		var f ReaderReaderIOResult[AppConfig, string] = func(AppConfig) func(context.Context) IOResult[string] {
			return func(context.Context) IOResult[string] {
				return func() Result[string] {
					return result.Left[string](testErr)
				}
			}
		}

		p := pair.MakePair[context.Context, AppConfig](t.Context(), cfg)
		outcome := Paired(f)(p)()

		assert.Equal(t, result.Left[string](testErr), outcome)
	})

	t.Run("is equivalent to direct curried application", func(t *testing.T) {
		cfg := AppConfig{DatabaseURL: "postgres://direct", LogLevel: "warn"}
		ctx := t.Context()

		f := Of[AppConfig]("value")

		direct := f(cfg)(ctx)()
		p := pair.MakePair[context.Context, AppConfig](ctx, cfg)
		viaPair := Paired(f)(p)()

		assert.Equal(t, direct, viaPair)
	})

	t.Run("maintains referential transparency", func(t *testing.T) {
		cfg := AppConfig{DatabaseURL: "postgres://ref", LogLevel: "info"}
		ctx := t.Context()

		var f ReaderReaderIOResult[AppConfig, int] = func(c AppConfig) func(context.Context) IOResult[int] {
			return func(context.Context) IOResult[int] {
				return func() Result[int] {
					return result.Of(len(c.DatabaseURL))
				}
			}
		}

		paired := Paired(f)
		p := pair.MakePair[context.Context, AppConfig](ctx, cfg)

		first := paired(p)()
		second := paired(p)()

		assert.Equal(t, first, second)
		assert.Equal(t, result.Of(len("postgres://ref")), first)
	})

	t.Run("works with different environment types", func(t *testing.T) {
		type Env struct {
			Multiplier int
		}
		env := Env{Multiplier: 7}
		ctx := t.Context()

		var f ReaderReaderIOResult[Env, int] = func(e Env) func(context.Context) IOResult[int] {
			return func(context.Context) IOResult[int] {
				return func() Result[int] {
					return result.Of(e.Multiplier * 6)
				}
			}
		}

		p := pair.MakePair[context.Context, Env](ctx, env)
		outcome := Paired(f)(p)()

		assert.Equal(t, result.Of(42), outcome)
	})
}
