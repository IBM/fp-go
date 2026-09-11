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

package effect

import (
	"context"
	"errors"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// TestChainFirstLeft_Contract pins down the documented semantics shared by
// ChainFirstLeft, TapLeft, MonadChainFirstLeft and MonadTapLeft: the handler runs
// only on failure, sees the error and the effect's context, and never changes
// the outcome — neither by succeeding nor by failing.
func TestChainFirstLeft_Contract(t *testing.T) {
	originalErr := errors.New("original")
	handlerErr := errors.New("handler")

	type handler = Kleisli[TestContext, error, string]

	variants := map[string]func(handler) Operator[TestContext, int, int]{
		"ChainFirstLeft": ChainFirstLeft[int, TestContext, string],
		"TapLeft":        TapLeft[int, TestContext, string],
		"MonadChainFirstLeft": func(f handler) Operator[TestContext, int, int] {
			return func(ma Effect[TestContext, int]) Effect[TestContext, int] {
				return MonadChainFirstLeft(ma, f)
			}
		},
		"MonadTapLeft": func(f handler) Operator[TestContext, int, int] {
			return func(ma Effect[TestContext, int]) Effect[TestContext, int] {
				return MonadTapLeft(ma, f)
			}
		},
	}

	for name, op := range variants {
		t.Run(name, func(t *testing.T) {
			t.Run("original error survives a failing handler", func(t *testing.T) {
				failing := func(error) Effect[TestContext, string] {
					return Fail[TestContext, string](handlerErr)
				}

				_, err := runEffect(op(failing)(Fail[TestContext, int](originalErr)), TestContext{})

				assert.Equal(t, originalErr, err)
			})

			t.Run("succeeding handler does not recover", func(t *testing.T) {
				succeeding := func(error) Effect[TestContext, string] {
					return Of[TestContext]("ignored")
				}

				_, err := runEffect(op(succeeding)(Fail[TestContext, int](originalErr)), TestContext{})

				assert.Equal(t, originalErr, err)
			})

			t.Run("handler sees the error and the effect's context", func(t *testing.T) {
				var seenErr error
				var seenValue string
				observing := func(err error) Effect[TestContext, string] {
					seenErr = err
					return Asks(func(c TestContext) string {
						seenValue = c.Value
						return c.Value
					})
				}

				_, err := runEffect(op(observing)(Fail[TestContext, int](originalErr)), TestContext{Value: "cfg"})

				assert.Equal(t, originalErr, err)
				assert.Equal(t, originalErr, seenErr)
				assert.Equal(t, "cfg", seenValue)
			})

			t.Run("handler is skipped on success", func(t *testing.T) {
				called := false
				observing := func(error) Effect[TestContext, string] {
					called = true
					return Of[TestContext]("")
				}

				value, err := runEffect(op(observing)(Of[TestContext](42)), TestContext{})

				assert.NoError(t, err)
				assert.Equal(t, 42, value)
				assert.False(t, called)
			})
		})
	}
}

// pairedTestKey is a context.Context key used to check that Paired forwards the
// head of the pair as the runtime context.
type pairedTestKey struct{}

// TestPaired checks that Paired routes the pair's head to the runtime
// context.Context and its tail to the effect's dependency.
func TestPaired(t *testing.T) {
	t.Run("head is the runtime context, tail is the dependency", func(t *testing.T) {
		requestID := FromThunk[TestContext](func(ctx context.Context) IOResult[string] {
			return func() Result[string] {
				return result.Of(ctx.Value(pairedTestKey{}).(string))
			}
		})
		eff := F.Pipe1(
			requestID,
			Chain(func(id string) Effect[TestContext, string] {
				return Asks(func(c TestContext) string { return id + "/" + c.Value })
			}),
		)

		ctx := context.WithValue(context.Background(), pairedTestKey{}, "req-1")
		res := Paired(eff)(pair.MakePair(ctx, TestContext{Value: "cfg"}))()

		assert.Equal(t, result.Of("req-1/cfg"), res)
	})

	t.Run("failure is returned as Left", func(t *testing.T) {
		err := errors.New("boom")

		res := Paired(Fail[TestContext, int](err))(pair.MakePair(context.Background(), TestContext{}))()

		assert.Equal(t, result.Left[int](err), res)
	})

	t.Run("does nothing until the IOResult is executed", func(t *testing.T) {
		executed := false
		eff := FromIO[TestContext](func() int {
			executed = true
			return 1
		})

		io := Paired(eff)(pair.MakePair(context.Background(), TestContext{}))
		assert.False(t, executed)

		io()
		assert.True(t, executed)
	})
}
