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

	"github.com/stretchr/testify/assert"
)

func TestProvide(t *testing.T) {
	t.Run("provides context to effect", func(t *testing.T) {
		ctx := TestContext{Value: "test-value"}
		eff := Of[TestContext]("result")

		ioResult := Provide[TestContext, string](ctx)(eff)
		readerResult := RunSync(ioResult)
		result, err := readerResult(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "result", result)
	})

	t.Run("provides context with specific values", func(t *testing.T) {
		type Config struct {
			Host string
			Port int
		}

		cfg := Config{Host: "localhost", Port: 8080}
		eff := Of[Config]("connected")

		ioResult := Provide[Config, string](cfg)(eff)
		readerResult := RunSync(ioResult)
		result, err := readerResult(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "connected", result)
	})

	t.Run("propagates errors", func(t *testing.T) {
		expectedErr := errors.New("provide error")
		ctx := TestContext{Value: "test"}
		eff := Fail[TestContext, string](expectedErr)

		ioResult := Provide[TestContext, string](ctx)(eff)
		readerResult := RunSync(ioResult)
		_, err := readerResult(context.Background())

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("works with different context types", func(t *testing.T) {
		type SimpleContext struct {
			ID int
		}

		ctx := SimpleContext{ID: 42}
		eff := Of[SimpleContext](100)

		ioResult := Provide[SimpleContext, int](ctx)(eff)
		readerResult := RunSync(ioResult)
		result, err := readerResult(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, 100, result)
	})

	t.Run("provides context to chained effects", func(t *testing.T) {
		ctx := TestContext{Value: "base"}

		eff := Chain(func(x int) Effect[TestContext, string] {
			return Of[TestContext]("result")
		})(Of[TestContext](42))

		ioResult := Provide[TestContext, string](ctx)(eff)
		readerResult := RunSync(ioResult)
		result, err := readerResult(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "result", result)
	})

	t.Run("provides context to mapped effects", func(t *testing.T) {
		ctx := TestContext{Value: "test"}

		eff := Map[TestContext](func(x int) string {
			return "mapped"
		})(Of[TestContext](42))

		ioResult := Provide[TestContext, string](ctx)(eff)
		readerResult := RunSync(ioResult)
		result, err := readerResult(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "mapped", result)
	})
}

func TestRunSync(t *testing.T) {
	t.Run("runs effect synchronously", func(t *testing.T) {
		ctx := TestContext{Value: "test"}
		eff := Of[TestContext](42)

		ioResult := Provide[TestContext, int](ctx)(eff)
		readerResult := RunSync(ioResult)
		result, err := readerResult(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, 42, result)
	})

	t.Run("runs effect with context.Context", func(t *testing.T) {
		ctx := TestContext{Value: "test"}
		eff := Of[TestContext]("hello")

		ioResult := Provide[TestContext, string](ctx)(eff)
		readerResult := RunSync(ioResult)

		bgCtx := context.Background()
		result, err := readerResult(bgCtx)

		assert.NoError(t, err)
		assert.Equal(t, "hello", result)
	})

	t.Run("propagates errors synchronously", func(t *testing.T) {
		expectedErr := errors.New("sync error")
		ctx := TestContext{Value: "test"}
		eff := Fail[TestContext, int](expectedErr)

		ioResult := Provide[TestContext, int](ctx)(eff)
		readerResult := RunSync(ioResult)
		_, err := readerResult(context.Background())

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("runs complex effect chains", func(t *testing.T) {
		ctx := TestContext{Value: "test"}

		eff := Chain(func(x int) Effect[TestContext, int] {
			return Of[TestContext](x * 2)
		})(Chain(func(x int) Effect[TestContext, int] {
			return Of[TestContext](x + 10)
		})(Of[TestContext](5)))

		ioResult := Provide[TestContext, int](ctx)(eff)
		readerResult := RunSync(ioResult)
		result, err := readerResult(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, 30, result) // (5 + 10) * 2
	})

	t.Run("handles multiple sequential runs", func(t *testing.T) {
		ctx := TestContext{Value: "test"}
		eff := Of[TestContext](42)

		ioResult := Provide[TestContext, int](ctx)(eff)
		readerResult := RunSync(ioResult)

		// Run multiple times
		result1, err1 := readerResult(context.Background())
		result2, err2 := readerResult(context.Background())
		result3, err3 := readerResult(context.Background())

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NoError(t, err3)
		assert.Equal(t, 42, result1)
		assert.Equal(t, 42, result2)
		assert.Equal(t, 42, result3)
	})

	t.Run("works with different result types", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}

		ctx := TestContext{Value: "test"}
		user := User{Name: "Alice", Age: 30}
		eff := Of[TestContext](user)

		ioResult := Provide[TestContext, User](ctx)(eff)
		readerResult := RunSync(ioResult)
		result, err := readerResult(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, user, result)
	})
}

func TestProvideAndRunSyncIntegration(t *testing.T) {
	t.Run("complete workflow with success", func(t *testing.T) {
		type AppConfig struct {
			APIKey  string
			Timeout int
		}

		cfg := AppConfig{APIKey: "secret", Timeout: 30}

		// Create an effect that uses the config
		eff := Of[AppConfig]("API call successful")

		// Provide config and run
		result, err := RunSync(Provide[AppConfig, string](cfg)(eff))(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "API call successful", result)
	})

	t.Run("complete workflow with error", func(t *testing.T) {
		type AppConfig struct {
			APIKey string
		}

		expectedErr := errors.New("API error")
		cfg := AppConfig{APIKey: "secret"}

		eff := Fail[AppConfig, string](expectedErr)

		_, err := RunSync(Provide[AppConfig, string](cfg)(eff))(context.Background())

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("workflow with transformations", func(t *testing.T) {
		ctx := TestContext{Value: "test"}

		eff := Map[TestContext](func(x int) string {
			return "final"
		})(Chain(func(x int) Effect[TestContext, int] {
			return Of[TestContext](x * 2)
		})(Of[TestContext](21)))

		result, err := RunSync(Provide[TestContext, string](ctx)(eff))(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "final", result)
	})

	t.Run("workflow with bind operations", func(t *testing.T) {
		type State struct {
			X int
			Y int
		}

		ctx := TestContext{Value: "test"}

		eff := Bind(
			func(y int) func(State) State {
				return func(s State) State {
					s.Y = y
					return s
				}
			},
			func(s State) Effect[TestContext, int] {
				return Of[TestContext](s.X * 2)
			},
		)(BindTo[TestContext](func(x int) State {
			return State{X: x}
		})(Of[TestContext](10)))

		result, err := RunSync(Provide[TestContext, State](ctx)(eff))(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, 10, result.X)
		assert.Equal(t, 20, result.Y)
	})

	t.Run("workflow with context transformation", func(t *testing.T) {
		type OuterCtx struct {
			Value string
		}
		type InnerCtx struct {
			Data string
		}

		outerCtx := OuterCtx{Value: "outer"}
		innerEff := Of[InnerCtx]("inner result")

		// Transform context
		transformedEff := Local[OuterCtx, InnerCtx, string](func(outer OuterCtx) InnerCtx {
			return InnerCtx{Data: outer.Value + "-transformed"}
		})(innerEff)

		result, err := RunSync(Provide[OuterCtx, string](outerCtx)(transformedEff))(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, "inner result", result)
	})

	t.Run("workflow with array traversal", func(t *testing.T) {
		ctx := TestContext{Value: "test"}
		input := []int{1, 2, 3, 4, 5}

		eff := TraverseArray(func(x int) Effect[TestContext, int] {
			return Of[TestContext](x * 2)
		})(input)

		result, err := RunSync(Provide[TestContext, []int](ctx)(eff))(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, []int{2, 4, 6, 8, 10}, result)
	})
}
