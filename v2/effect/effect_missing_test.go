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
	"strconv"
	"testing"

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
