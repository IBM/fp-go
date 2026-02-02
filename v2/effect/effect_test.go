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
	"fmt"
	"testing"

	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

type TestContext struct {
	Value string
}

func runEffect[A any](eff Effect[TestContext, A], ctx TestContext) (A, error) {
	ioResult := Provide[TestContext, A](ctx)(eff)
	readerResult := RunSync(ioResult)
	return readerResult(context.Background())
}

func TestSucceed(t *testing.T) {
	t.Run("creates successful effect with value", func(t *testing.T) {
		eff := Succeed[TestContext](42)
		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, 42, result)
	})

	t.Run("creates successful effect with string", func(t *testing.T) {
		eff := Succeed[TestContext]("hello")
		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "hello", result)
	})

	t.Run("creates successful effect with struct", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}
		user := User{Name: "Alice", Age: 30}
		eff := Succeed[TestContext](user)
		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, user, result)
	})
}

func TestFail(t *testing.T) {
	t.Run("creates failed effect with error", func(t *testing.T) {
		expectedErr := errors.New("test error")
		eff := Fail[TestContext, int](expectedErr)
		_, err := runEffect(eff, TestContext{Value: "test"})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("creates failed effect with custom error", func(t *testing.T) {
		expectedErr := fmt.Errorf("custom error: %s", "details")
		eff := Fail[TestContext, string](expectedErr)
		_, err := runEffect(eff, TestContext{Value: "test"})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestOf(t *testing.T) {
	t.Run("lifts value into effect", func(t *testing.T) {
		eff := Of[TestContext](100)
		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, 100, result)
	})

	t.Run("is equivalent to Succeed", func(t *testing.T) {
		value := "test value"
		eff1 := Of[TestContext](value)
		eff2 := Succeed[TestContext](value)

		result1, err1 := runEffect(eff1, TestContext{Value: "test"})
		result2, err2 := runEffect(eff2, TestContext{Value: "test"})

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, result1, result2)
	})
}

func TestMap(t *testing.T) {
	t.Run("maps over successful effect", func(t *testing.T) {
		eff := Of[TestContext](10)
		mapped := Map[TestContext](func(x int) int {
			return x * 2
		})(eff)

		result, err := runEffect(mapped, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, 20, result)
	})

	t.Run("maps to different type", func(t *testing.T) {
		eff := Of[TestContext](42)
		mapped := Map[TestContext](func(x int) string {
			return fmt.Sprintf("value: %d", x)
		})(eff)

		result, err := runEffect(mapped, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "value: 42", result)
	})

	t.Run("preserves error in failed effect", func(t *testing.T) {
		expectedErr := errors.New("original error")
		eff := Fail[TestContext, int](expectedErr)
		mapped := Map[TestContext](func(x int) int {
			return x * 2
		})(eff)

		_, err := runEffect(mapped, TestContext{Value: "test"})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("chains multiple maps", func(t *testing.T) {
		eff := Of[TestContext](5)
		result := Map[TestContext](func(x int) int {
			return x + 1
		})(Map[TestContext](func(x int) int {
			return x * 2
		})(eff))

		value, err := runEffect(result, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, 11, value) // (5 * 2) + 1
	})
}

func TestChain(t *testing.T) {
	t.Run("chains successful effects", func(t *testing.T) {
		eff := Of[TestContext](10)
		chained := Chain(func(x int) Effect[TestContext, int] {
			return Of[TestContext](x * 2)
		})(eff)

		result, err := runEffect(chained, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, 20, result)
	})

	t.Run("chains to different type", func(t *testing.T) {
		eff := Of[TestContext](42)
		chained := Chain(func(x int) Effect[TestContext, string] {
			return Of[TestContext](fmt.Sprintf("number: %d", x))
		})(eff)

		result, err := runEffect(chained, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "number: 42", result)
	})

	t.Run("propagates first error", func(t *testing.T) {
		expectedErr := errors.New("first error")
		eff := Fail[TestContext, int](expectedErr)
		chained := Chain(func(x int) Effect[TestContext, int] {
			return Of[TestContext](x * 2)
		})(eff)

		_, err := runEffect(chained, TestContext{Value: "test"})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("propagates second error", func(t *testing.T) {
		expectedErr := errors.New("second error")
		eff := Of[TestContext](10)
		chained := Chain(func(x int) Effect[TestContext, int] {
			return Fail[TestContext, int](expectedErr)
		})(eff)

		_, err := runEffect(chained, TestContext{Value: "test"})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("chains multiple operations", func(t *testing.T) {
		eff := Of[TestContext](5)
		result := Chain(func(x int) Effect[TestContext, int] {
			return Of[TestContext](x + 10)
		})(Chain(func(x int) Effect[TestContext, int] {
			return Of[TestContext](x * 2)
		})(eff))

		value, err := runEffect(result, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, 20, value) // (5 * 2) + 10
	})
}

func TestAp(t *testing.T) {
	t.Run("applies function effect to value effect", func(t *testing.T) {
		fn := Of[TestContext](func(x int) int {
			return x * 2
		})
		value := Of[TestContext](21)

		result := Ap[int](value)(fn)
		val, err := runEffect(result, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, 42, val)
	})

	t.Run("applies function to different type", func(t *testing.T) {
		fn := Of[TestContext](func(x int) string {
			return fmt.Sprintf("value: %d", x)
		})
		value := Of[TestContext](42)

		result := Ap[string](value)(fn)
		val, err := runEffect(result, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "value: 42", val)
	})

	t.Run("propagates error from function effect", func(t *testing.T) {
		expectedErr := errors.New("function error")
		fn := Fail[TestContext, func(int) int](expectedErr)
		value := Of[TestContext](42)

		result := Ap[int](value)(fn)
		_, err := runEffect(result, TestContext{Value: "test"})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("propagates error from value effect", func(t *testing.T) {
		expectedErr := errors.New("value error")
		fn := Of[TestContext](func(x int) int {
			return x * 2
		})
		value := Fail[TestContext, int](expectedErr)

		result := Ap[int](value)(fn)
		_, err := runEffect(result, TestContext{Value: "test"})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestSuspend(t *testing.T) {
	t.Run("suspends effect computation", func(t *testing.T) {
		callCount := 0
		eff := Suspend(func() Effect[TestContext, int] {
			callCount++
			return Of[TestContext](42)
		})

		// Effect not executed yet
		assert.Equal(t, 0, callCount)

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, 42, result)
		assert.Equal(t, 1, callCount)
	})

	t.Run("suspends failing effect", func(t *testing.T) {
		expectedErr := errors.New("suspended error")
		eff := Suspend(func() Effect[TestContext, int] {
			return Fail[TestContext, int](expectedErr)
		})

		_, err := runEffect(eff, TestContext{Value: "test"})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("allows lazy evaluation", func(t *testing.T) {
		var value int
		eff := Suspend(func() Effect[TestContext, int] {
			return Of[TestContext](value)
		})

		value = 10
		result1, err1 := runEffect(eff, TestContext{Value: "test"})

		value = 20
		result2, err2 := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, 10, result1)
		assert.Equal(t, 20, result2)
	})
}

func TestTap(t *testing.T) {
	t.Run("executes side effect without changing value", func(t *testing.T) {
		sideEffectValue := 0
		eff := Of[TestContext](42)
		tapped := Tap(func(x int) Effect[TestContext, any] {
			sideEffectValue = x * 2
			return Of[TestContext, any](nil)
		})(eff)

		result, err := runEffect(tapped, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, 42, result)
		assert.Equal(t, 84, sideEffectValue)
	})

	t.Run("propagates original error", func(t *testing.T) {
		expectedErr := errors.New("original error")
		eff := Fail[TestContext, int](expectedErr)
		tapped := Tap(func(x int) Effect[TestContext, any] {
			return Of[TestContext, any](nil)
		})(eff)

		_, err := runEffect(tapped, TestContext{Value: "test"})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("propagates tap error", func(t *testing.T) {
		expectedErr := errors.New("tap error")
		eff := Of[TestContext](42)
		tapped := Tap(func(x int) Effect[TestContext, any] {
			return Fail[TestContext, any](expectedErr)
		})(eff)

		_, err := runEffect(tapped, TestContext{Value: "test"})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("chains multiple taps", func(t *testing.T) {
		values := []int{}
		eff := Of[TestContext](10)
		result := Tap(func(x int) Effect[TestContext, any] {
			values = append(values, x+2)
			return Of[TestContext, any](nil)
		})(Tap(func(x int) Effect[TestContext, any] {
			values = append(values, x+1)
			return Of[TestContext, any](nil)
		})(eff))

		value, err := runEffect(result, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, 10, value)
		assert.Equal(t, []int{11, 12}, values)
	})
}

func TestTernary(t *testing.T) {
	t.Run("executes onTrue when predicate is true", func(t *testing.T) {
		kleisli := Ternary(
			func(x int) bool { return x > 10 },
			func(x int) Effect[TestContext, string] {
				return Of[TestContext]("greater")
			},
			func(x int) Effect[TestContext, string] {
				return Of[TestContext]("less or equal")
			},
		)

		result, err := runEffect(kleisli(15), TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "greater", result)
	})

	t.Run("executes onFalse when predicate is false", func(t *testing.T) {
		kleisli := Ternary(
			func(x int) bool { return x > 10 },
			func(x int) Effect[TestContext, string] {
				return Of[TestContext]("greater")
			},
			func(x int) Effect[TestContext, string] {
				return Of[TestContext]("less or equal")
			},
		)

		result, err := runEffect(kleisli(5), TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "less or equal", result)
	})

	t.Run("handles errors in onTrue branch", func(t *testing.T) {
		expectedErr := errors.New("true branch error")
		kleisli := Ternary(
			func(x int) bool { return x > 10 },
			func(x int) Effect[TestContext, string] {
				return Fail[TestContext, string](expectedErr)
			},
			func(x int) Effect[TestContext, string] {
				return Of[TestContext]("less or equal")
			},
		)

		_, err := runEffect(kleisli(15), TestContext{Value: "test"})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("handles errors in onFalse branch", func(t *testing.T) {
		expectedErr := errors.New("false branch error")
		kleisli := Ternary(
			func(x int) bool { return x > 10 },
			func(x int) Effect[TestContext, string] {
				return Of[TestContext]("greater")
			},
			func(x int) Effect[TestContext, string] {
				return Fail[TestContext, string](expectedErr)
			},
		)

		_, err := runEffect(kleisli(5), TestContext{Value: "test"})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})
}

func TestEffectComposition(t *testing.T) {
	t.Run("composes Map and Chain", func(t *testing.T) {
		eff := Of[TestContext](5)
		result := Chain(func(x int) Effect[TestContext, string] {
			return Of[TestContext](fmt.Sprintf("result: %d", x))
		})(Map[TestContext](func(x int) int {
			return x * 2
		})(eff))

		value, err := runEffect(result, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "result: 10", value)
	})

	t.Run("composes Chain and Tap", func(t *testing.T) {
		sideEffect := 0
		eff := Of[TestContext](10)
		result := Tap(func(x int) Effect[TestContext, any] {
			sideEffect = x
			return Of[TestContext, any](nil)
		})(Chain(func(x int) Effect[TestContext, int] {
			return Of[TestContext](x * 2)
		})(eff))

		value, err := runEffect(result, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, 20, value)
		assert.Equal(t, 20, sideEffect)
	})
}

func TestEffectWithResult(t *testing.T) {
	t.Run("converts result to effect", func(t *testing.T) {
		res := result.Of(42)
		// This demonstrates integration with result package
		assert.True(t, result.IsRight(res))
	})
}
