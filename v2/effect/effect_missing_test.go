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
	"strconv"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

func TestRead(t *testing.T) {
	t.Run("provides context to effect", func(t *testing.T) {
		ctx := TestContext{Value: "test-context"}
		eff := Of[TestContext](42)

		thunk := Read[int](ctx)(eff)
		ioResult := thunk(context.Background())
		res := ioResult()

		assert.True(t, result.IsRight(res))
		value, err := result.Unwrap(res)
		assert.NoError(t, err)
		assert.Equal(t, 42, value)
	})

	t.Run("provides context to failing effect", func(t *testing.T) {
		expectedErr := errors.New("read error")
		ctx := TestContext{Value: "test"}
		eff := Fail[TestContext, string](expectedErr)

		thunk := Read[string](ctx)(eff)
		ioResult := thunk(context.Background())
		res := ioResult()

		assert.True(t, result.IsLeft(res))
		_, err := result.Unwrap(res)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("provides context to chained effects", func(t *testing.T) {
		ctx := TestContext{Value: "base"}
		eff := Chain(func(x int) Effect[TestContext, string] {
			return Of[TestContext](strconv.Itoa(x * 2))
		})(Of[TestContext](21))

		thunk := Read[string](ctx)(eff)
		ioResult := thunk(context.Background())
		res := ioResult()

		assert.True(t, result.IsRight(res))
		value, err := result.Unwrap(res)
		assert.NoError(t, err)
		assert.Equal(t, "42", value)
	})

	t.Run("works with different context types", func(t *testing.T) {
		type CustomContext struct {
			ID   int
			Name string
		}

		ctx := CustomContext{ID: 100, Name: "custom"}
		eff := Of[CustomContext]("result")

		thunk := Read[string](ctx)(eff)
		ioResult := thunk(context.Background())
		res := ioResult()

		assert.True(t, result.IsRight(res))
		value, err := result.Unwrap(res)
		assert.NoError(t, err)
		assert.Equal(t, "result", value)
	})

	t.Run("can be composed with RunSync", func(t *testing.T) {
		ctx := TestContext{Value: "test"}
		eff := Of[TestContext](100)

		thunk := Read[int](ctx)(eff)
		readerResult := RunSync(thunk)
		value, err := readerResult(context.Background())

		assert.NoError(t, err)
		assert.Equal(t, 100, value)
	})
}

func TestChainResultK(t *testing.T) {
	t.Run("chains successful Result function", func(t *testing.T) {
		parseIntResult := result.Eitherize1(strconv.Atoi)
		eff := Of[TestContext]("42")
		chained := ChainResultK[TestContext](parseIntResult)(eff)

		result, err := runEffect(chained, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, 42, result)
	})

	t.Run("chains failing Result function", func(t *testing.T) {
		parseIntResult := result.Eitherize1(strconv.Atoi)
		eff := Of[TestContext]("not-a-number")
		chained := ChainResultK[TestContext](parseIntResult)(eff)

		_, err := runEffect(chained, TestContext{Value: "test"})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid syntax")
	})

	t.Run("propagates error from original effect", func(t *testing.T) {
		expectedErr := errors.New("original error")
		parseIntResult := result.Eitherize1(strconv.Atoi)
		eff := Fail[TestContext, string](expectedErr)
		chained := ChainResultK[TestContext](parseIntResult)(eff)

		_, err := runEffect(chained, TestContext{Value: "test"})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("chains multiple Result functions", func(t *testing.T) {
		parseIntResult := result.Eitherize1(strconv.Atoi)
		formatResult := func(x int) result.Result[string] {
			return result.Of("value: " + strconv.Itoa(x))
		}

		eff := Of[TestContext]("42")
		chained := ChainResultK[TestContext](formatResult)(
			ChainResultK[TestContext](parseIntResult)(eff),
		)

		result, err := runEffect(chained, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "value: 42", result)
	})

	t.Run("integrates with other effect operations", func(t *testing.T) {
		parseIntResult := result.Eitherize1(strconv.Atoi)

		eff := Map[TestContext](func(x int) string {
			return "final: " + strconv.Itoa(x)
		})(ChainResultK[TestContext](parseIntResult)(Of[TestContext]("100")))

		result, err := runEffect(eff, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "final: 100", result)
	})

	t.Run("works with custom Result functions", func(t *testing.T) {
		validatePositive := func(x int) result.Result[int] {
			if x > 0 {
				return result.Of(x)
			}
			return result.Left[int](errors.New("must be positive"))
		}

		parseIntResult := result.Eitherize1(strconv.Atoi)

		// Test with positive number
		eff1 := ChainResultK[TestContext](validatePositive)(
			ChainResultK[TestContext](parseIntResult)(Of[TestContext]("42")),
		)
		result1, err1 := runEffect(eff1, TestContext{Value: "test"})
		assert.NoError(t, err1)
		assert.Equal(t, 42, result1)

		// Test with negative number
		eff2 := ChainResultK[TestContext](validatePositive)(
			ChainResultK[TestContext](parseIntResult)(Of[TestContext]("-5")),
		)
		_, err2 := runEffect(eff2, TestContext{Value: "test"})
		assert.Error(t, err2)
		assert.Contains(t, err2.Error(), "must be positive")
	})

	t.Run("preserves error context", func(t *testing.T) {
		customError := errors.New("custom validation error")
		validateFunc := func(s string) result.Result[string] {
			if len(s) > 0 {
				return result.Of(s)
			}
			return result.Left[string](customError)
		}

		eff := ChainResultK[TestContext](validateFunc)(Of[TestContext](""))
		_, err := runEffect(eff, TestContext{Value: "test"})

		assert.Error(t, err)
		assert.Equal(t, customError, err)
	})
}

func TestFromReader_Success(t *testing.T) {
	t.Run("lifts a Reader that returns a constant", func(t *testing.T) {
		r := func(_ TestContext) int { return 42 }
		eff := FromReader(r)

		value, err := runEffect(eff, TestContext{Value: "ignored"})

		assert.NoError(t, err)
		assert.Equal(t, 42, value)
	})

	t.Run("lifts a Reader that reads from the environment", func(t *testing.T) {
		r := func(ctx TestContext) string { return ctx.Value }
		eff := FromReader(r)

		value, err := runEffect(eff, TestContext{Value: "hello"})

		assert.NoError(t, err)
		assert.Equal(t, "hello", value)
	})

	t.Run("lifts a Reader that computes a derived value", func(t *testing.T) {
		r := func(cfg TestConfig) int { return cfg.Multiplier * 10 }
		eff := FromReader(r)

		value, err := runEffect(eff, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, 30, value)
	})
}

func TestFromReader_EdgeCases(t *testing.T) {
	t.Run("handles zero value environment", func(t *testing.T) {
		r := func(ctx TestContext) string { return ctx.Value }
		eff := FromReader(r)

		value, err := runEffect(eff, TestContext{})

		assert.NoError(t, err)
		assert.Equal(t, "", value)
	})

	t.Run("handles boolean result", func(t *testing.T) {
		r := func(cfg TestConfig) bool { return cfg.Multiplier > 0 }
		eff := FromReader(r)

		value, err := runEffect(eff, testConfig)

		assert.NoError(t, err)
		assert.True(t, value)
	})
}

func TestFromReader_Integration(t *testing.T) {
	t.Run("composes with Map", func(t *testing.T) {
		r := func(cfg TestConfig) int { return cfg.Multiplier }
		eff := Map[TestConfig](func(n int) string { return strconv.Itoa(n) })(
			FromReader(r),
		)

		value, err := runEffect(eff, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, "3", value)
	})

	t.Run("composes with Chain", func(t *testing.T) {
		r := func(cfg TestConfig) int { return cfg.Multiplier }
		eff := Chain(func(n int) Effect[TestConfig, int] {
			return Of[TestConfig](n * 2)
		})(FromReader(r))

		value, err := runEffect(eff, testConfig)

		assert.NoError(t, err)
		assert.Equal(t, 6, value)
	})

	t.Run("produces same result as Asks for identical Reader", func(t *testing.T) {
		r := func(cfg TestConfig) string { return cfg.Prefix }

		effFromReader := FromReader(r)
		effAsks := Asks(r)

		v1, err1 := runEffect(effFromReader, testConfig)
		v2, err2 := runEffect(effAsks, testConfig)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, v1, v2)
	})
}

func TestFromReaderResult_Success(t *testing.T) {
	t.Run("lifts a ReaderResult returning a Right", func(t *testing.T) {
		rr := func(_ TestContext) result.Result[int] { return result.Of(42) }
		eff := FromReaderResult(rr)

		value, err := runEffect(eff, TestContext{Value: "env"})

		assert.NoError(t, err)
		assert.Equal(t, 42, value)
	})

	t.Run("reads from the environment before producing a Right", func(t *testing.T) {
		rr := func(ctx TestContext) result.Result[string] {
			return result.Of("hello " + ctx.Value)
		}
		eff := FromReaderResult(rr)

		value, err := runEffect(eff, TestContext{Value: "world"})

		assert.NoError(t, err)
		assert.Equal(t, "hello world", value)
	})

	t.Run("lifts a ReaderResult built from result.Eitherize1", func(t *testing.T) {
		parseIntRR := func(_ TestContext) result.Result[int] {
			return result.Eitherize1(strconv.Atoi)("99")
		}
		eff := FromReaderResult(parseIntRR)

		value, err := runEffect(eff, TestContext{})

		assert.NoError(t, err)
		assert.Equal(t, 99, value)
	})
}

func TestFromReaderResult_Failure(t *testing.T) {
	t.Run("propagates a Left as a failed effect", func(t *testing.T) {
		expectedErr := errors.New("something went wrong")
		rr := func(_ TestContext) result.Result[int] {
			return result.Left[int](expectedErr)
		}
		eff := FromReaderResult(rr)

		_, err := runEffect(eff, TestContext{})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("failure from a reader-based parse error", func(t *testing.T) {
		rr := func(_ TestContext) result.Result[int] {
			return result.Eitherize1(strconv.Atoi)("not-a-number")
		}
		eff := FromReaderResult(rr)

		_, err := runEffect(eff, TestContext{})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid syntax")
	})

	t.Run("error message is preserved exactly", func(t *testing.T) {
		customErr := errors.New("custom validation error")
		rr := func(_ TestContext) result.Result[string] {
			return result.Left[string](customErr)
		}
		eff := FromReaderResult(rr)

		_, err := runEffect(eff, TestContext{})

		assert.Equal(t, customErr, err)
	})
}

func TestFromReaderResult_EdgeCases(t *testing.T) {
	t.Run("zero-value environment", func(t *testing.T) {
		rr := func(ctx TestContext) result.Result[string] {
			return result.Of(ctx.Value)
		}
		eff := FromReaderResult(rr)

		value, err := runEffect(eff, TestContext{})

		assert.NoError(t, err)
		assert.Equal(t, "", value)
	})

	t.Run("environment is passed exactly once", func(t *testing.T) {
		callCount := 0
		rr := func(_ TestContext) result.Result[int] {
			callCount++
			return result.Of(callCount)
		}
		eff := FromReaderResult(rr)

		value, err := runEffect(eff, TestContext{})

		assert.NoError(t, err)
		assert.Equal(t, 1, value)
		assert.Equal(t, 1, callCount)
	})
}

func TestFromReaderResult_Integration(t *testing.T) {
	t.Run("composes with Map on success", func(t *testing.T) {
		rr := func(_ TestContext) result.Result[int] { return result.Of(21) }
		eff := F.Pipe1(
			FromReaderResult(rr),
			Map[TestContext](func(n int) int { return n * 2 }),
		)

		value, err := runEffect(eff, TestContext{})

		assert.NoError(t, err)
		assert.Equal(t, 42, value)
	})

	t.Run("composes with Chain on success", func(t *testing.T) {
		rr := func(ctx TestContext) result.Result[string] {
			return result.Of(ctx.Value)
		}
		eff := F.Pipe1(
			FromReaderResult(rr),
			Chain(func(s string) Effect[TestContext, int] {
				n, _ := strconv.Atoi(s)
				return Of[TestContext](n)
			}),
		)

		value, err := runEffect(eff, TestContext{Value: "7"})

		assert.NoError(t, err)
		assert.Equal(t, 7, value)
	})

	t.Run("Map is skipped when ReaderResult is a Left", func(t *testing.T) {
		rr := func(_ TestContext) result.Result[int] {
			return result.Left[int](errors.New("upstream error"))
		}
		eff := F.Pipe1(
			FromReaderResult(rr),
			Map[TestContext](func(n int) int { return n * 2 }),
		)

		_, err := runEffect(eff, TestContext{})

		assert.Error(t, err)
		assert.EqualError(t, err, "upstream error")
	})

	t.Run("composes with FromReader for the same environment", func(t *testing.T) {
		plainReader := func(ctx TestContext) int { return len(ctx.Value) }
		rrReader := func(ctx TestContext) result.Result[int] {
			if ctx.Value == "" {
				return result.Left[int](errors.New("empty value"))
			}
			return result.Of(len(ctx.Value))
		}

		effFromReader := FromReader(plainReader)
		effFromReaderResult := FromReaderResult(rrReader)

		v1, err1 := runEffect(effFromReader, TestContext{Value: "abc"})
		v2, err2 := runEffect(effFromReaderResult, TestContext{Value: "abc"})

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, v1, v2)
	})
}

// ExampleFromReaderResult demonstrates lifting a ReaderResult into an Effect.
func ExampleFromReaderResult() {
	type Config struct {
		Input string
	}

	parseIntRR := func(cfg Config) result.Result[int] {
		return result.Eitherize1(strconv.Atoi)(cfg.Input)
	}
	eff := FromReaderResult(parseIntRR)

	res := eff(Config{Input: "42"})(context.Background())()
	value, _ := result.Unwrap(res)
	fmt.Println(value)
	// Output: 42
}

// ExampleFromReaderResult_failure demonstrates that a Left result becomes a failed Effect.
func ExampleFromReaderResult_failure() {
	type Config struct {
		Input string
	}

	parseIntRR := func(cfg Config) result.Result[int] {
		return result.Eitherize1(strconv.Atoi)(cfg.Input)
	}
	eff := FromReaderResult(parseIntRR)

	res := eff(Config{Input: "bad"})(context.Background())()
	_, err := result.Unwrap(res)
	fmt.Println(result.IsLeft(res))
	fmt.Println(err != nil)
	// Output:
	// true
	// true
}

// ExampleFromReaderResult_composition demonstrates composing FromReaderResult with Map.
func ExampleFromReaderResult_composition() {
	type Config struct {
		Base int
	}

	rr := func(cfg Config) result.Result[int] { return result.Of(cfg.Base) }
	eff := F.Pipe1(
		FromReaderResult(rr),
		Map[Config](func(n int) string { return "value=" + strconv.Itoa(n*2) }),
	)

	res := eff(Config{Base: 21})(context.Background())()
	value, _ := result.Unwrap(res)
	fmt.Println(value)
	// Output: value=42
}
