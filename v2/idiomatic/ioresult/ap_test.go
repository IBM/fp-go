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

package ioresult

import (
	"errors"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/stretchr/testify/assert"
)

func TestMonadApFirst(t *testing.T) {
	t.Run("Both Right - keeps first value", func(t *testing.T) {
		first := Of("first")
		second := Of("second")
		result := MonadApFirst(first, second)

		val, err := result()
		assert.NoError(t, err)
		assert.Equal(t, "first", val)
	})

	t.Run("First Right, Second Left - returns second's error", func(t *testing.T) {
		first := Of(42)
		second := Left[string](errors.New("second error"))
		result := MonadApFirst(first, second)

		_, err := result()
		assert.Error(t, err)
		assert.Equal(t, "second error", err.Error())
	})

	t.Run("First Left, Second Right - returns first's error", func(t *testing.T) {
		first := Left[int](errors.New("first error"))
		second := Of("success")
		result := MonadApFirst(first, second)

		_, err := result()
		assert.Error(t, err)
		assert.Equal(t, "first error", err.Error())
	})

	t.Run("Both Left - returns first error", func(t *testing.T) {
		first := Left[int](errors.New("first error"))
		second := Left[string](errors.New("second error"))
		result := MonadApFirst(first, second)

		_, err := result()
		assert.Error(t, err)
		assert.Equal(t, "first error", err.Error())
	})

	t.Run("Different types", func(t *testing.T) {
		first := Of(123)
		second := Of("text")
		result := MonadApFirst(first, second)

		val, err := result()
		assert.NoError(t, err)
		assert.Equal(t, 123, val)
	})

	t.Run("Side effects execute for both", func(t *testing.T) {
		firstCalled := false
		secondCalled := false

		first := func() (int, error) {
			firstCalled = true
			return 1, nil
		}
		second := func() (string, error) {
			secondCalled = true
			return "test", nil
		}

		result := MonadApFirst(first, second)
		val, err := result()

		assert.True(t, firstCalled)
		assert.True(t, secondCalled)
		assert.NoError(t, err)
		assert.Equal(t, 1, val)
	})
}

func TestApFirst(t *testing.T) {
	t.Run("Both Right - keeps first value", func(t *testing.T) {
		result := F.Pipe1(
			Of("first"),
			ApFirst[string](Of("second")),
		)

		val, err := result()
		assert.NoError(t, err)
		assert.Equal(t, "first", val)
	})

	t.Run("First Right, Second Left - returns error", func(t *testing.T) {
		result := F.Pipe1(
			Of(100),
			ApFirst[int](Left[string](errors.New("error"))),
		)

		_, err := result()
		assert.Error(t, err)
		assert.Equal(t, "error", err.Error())
	})

	t.Run("First Left, Second Right - returns error", func(t *testing.T) {
		result := F.Pipe1(
			Left[int](errors.New("first error")),
			ApFirst[int](Of("success")),
		)

		_, err := result()
		assert.Error(t, err)
		assert.Equal(t, "first error", err.Error())
	})

	t.Run("Chain multiple ApFirst", func(t *testing.T) {
		result := F.Pipe2(
			Of("value"),
			ApFirst[string](Of(1)),
			ApFirst[string](Of(true)),
		)

		val, err := result()
		assert.NoError(t, err)
		assert.Equal(t, "value", val)
	})

	t.Run("With Map composition", func(t *testing.T) {
		result := F.Pipe2(
			Of(5),
			ApFirst[int](Of("ignored")),
			Map(N.Mul(2)),
		)

		val, err := result()
		assert.NoError(t, err)
		assert.Equal(t, 10, val)
	})

	t.Run("Side effect in second executes", func(t *testing.T) {
		effectExecuted := false
		second := func() (string, error) {
			effectExecuted = true
			return "effect", nil
		}

		result := F.Pipe1(
			Of(42),
			ApFirst[int](second),
		)

		val, err := result()
		assert.True(t, effectExecuted)
		assert.NoError(t, err)
		assert.Equal(t, 42, val)
	})
}

func TestMonadApSecond(t *testing.T) {
	t.Run("Both Right - keeps second value", func(t *testing.T) {
		first := Of("first")
		second := Of("second")
		result := MonadApSecond(first, second)

		val, err := result()
		assert.NoError(t, err)
		assert.Equal(t, "second", val)
	})

	t.Run("First Right, Second Left - returns second's error", func(t *testing.T) {
		first := Of(42)
		second := Left[string](errors.New("second error"))
		result := MonadApSecond(first, second)

		_, err := result()
		assert.Error(t, err)
		assert.Equal(t, "second error", err.Error())
	})

	t.Run("First Left, Second Right - returns first's error", func(t *testing.T) {
		first := Left[int](errors.New("first error"))
		second := Of("success")
		result := MonadApSecond(first, second)

		_, err := result()
		assert.Error(t, err)
		assert.Equal(t, "first error", err.Error())
	})

	t.Run("Both Left - returns first error", func(t *testing.T) {
		first := Left[int](errors.New("first error"))
		second := Left[string](errors.New("second error"))
		result := MonadApSecond(first, second)

		_, err := result()
		assert.Error(t, err)
		assert.Equal(t, "first error", err.Error())
	})

	t.Run("Different types", func(t *testing.T) {
		first := Of(123)
		second := Of("text")
		result := MonadApSecond(first, second)

		val, err := result()
		assert.NoError(t, err)
		assert.Equal(t, "text", val)
	})

	t.Run("Side effects execute for both", func(t *testing.T) {
		firstCalled := false
		secondCalled := false

		first := func() (int, error) {
			firstCalled = true
			return 1, nil
		}
		second := func() (string, error) {
			secondCalled = true
			return "test", nil
		}

		result := MonadApSecond(first, second)
		val, err := result()

		assert.True(t, firstCalled)
		assert.True(t, secondCalled)
		assert.NoError(t, err)
		assert.Equal(t, "test", val)
	})
}

func TestApSecond(t *testing.T) {
	t.Run("Both Right - keeps second value", func(t *testing.T) {
		result := F.Pipe1(
			Of("first"),
			ApSecond[string](Of("second")),
		)

		val, err := result()
		assert.NoError(t, err)
		assert.Equal(t, "second", val)
	})

	t.Run("First Right, Second Left - returns error", func(t *testing.T) {
		result := F.Pipe1(
			Of(100),
			ApSecond[int](Left[string](errors.New("error"))),
		)

		_, err := result()
		assert.Error(t, err)
		assert.Equal(t, "error", err.Error())
	})

	t.Run("First Left, Second Right - returns error", func(t *testing.T) {
		result := F.Pipe1(
			Left[int](errors.New("first error")),
			ApSecond[int](Of("success")),
		)

		_, err := result()
		assert.Error(t, err)
		assert.Equal(t, "first error", err.Error())
	})

	t.Run("Chain multiple ApSecond", func(t *testing.T) {
		result := F.Pipe2(
			Of("initial"),
			ApSecond[string](Of("middle")),
			ApSecond[string](Of("final")),
		)

		val, err := result()
		assert.NoError(t, err)
		assert.Equal(t, "final", val)
	})

	t.Run("With Map composition", func(t *testing.T) {
		result := F.Pipe2(
			Of(1),
			ApSecond[int](Of(5)),
			Map(N.Mul(2)),
		)

		val, err := result()
		assert.NoError(t, err)
		assert.Equal(t, 10, val)
	})

	t.Run("Side effect in first executes", func(t *testing.T) {
		effectExecuted := false
		first := func() (int, error) {
			effectExecuted = true
			return 1, nil
		}

		result := F.Pipe1(
			first,
			ApSecond[int](Of("result")),
		)

		val, err := result()
		assert.True(t, effectExecuted)
		assert.NoError(t, err)
		assert.Equal(t, "result", val)
	})
}

func TestApFirstApSecondInteraction(t *testing.T) {
	t.Run("ApFirst then ApSecond", func(t *testing.T) {
		// ApFirst keeps "first", then ApSecond discards it for "second"
		result := F.Pipe2(
			Of("first"),
			ApFirst[string](Of("middle")),
			ApSecond[string](Of("second")),
		)

		val, err := result()
		assert.NoError(t, err)
		assert.Equal(t, "second", val)
	})

	t.Run("ApSecond then ApFirst", func(t *testing.T) {
		// ApSecond picks "middle", then ApFirst keeps that over "ignored"
		result := F.Pipe2(
			Of("first"),
			ApSecond[string](Of("middle")),
			ApFirst[string](Of("ignored")),
		)

		val, err := result()
		assert.NoError(t, err)
		assert.Equal(t, "middle", val)
	})

	t.Run("Error propagation in chain", func(t *testing.T) {
		result := F.Pipe2(
			Of("first"),
			ApFirst[string](Left[int](errors.New("error in middle"))),
			ApSecond[string](Of("second")),
		)

		_, err := result()
		assert.Error(t, err)
		assert.Equal(t, "error in middle", err.Error())
	})
}

func TestApSequencingBehavior(t *testing.T) {
	t.Run("MonadApFirst executes both operations", func(t *testing.T) {
		executionOrder := make([]string, 2)

		first := func() (int, error) {
			executionOrder[0] = "first"
			return 1, nil
		}
		second := func() (string, error) {
			executionOrder[1] = "second"
			return "test", nil
		}

		result := MonadApFirst(first, second)
		_, err := result()

		assert.NoError(t, err)
		// Note: execution order is second then first due to applicative implementation
		assert.Len(t, executionOrder, 2)
		assert.Contains(t, executionOrder, "first")
		assert.Contains(t, executionOrder, "second")
	})

	t.Run("MonadApSecond executes both operations", func(t *testing.T) {
		executionOrder := make([]string, 2)

		first := func() (int, error) {
			executionOrder[0] = "first"
			return 1, nil
		}
		second := func() (string, error) {
			executionOrder[1] = "second"
			return "test", nil
		}

		result := MonadApSecond(first, second)
		_, err := result()

		assert.NoError(t, err)
		// Note: execution order is second then first due to applicative implementation
		assert.Len(t, executionOrder, 2)
		assert.Contains(t, executionOrder, "first")
		assert.Contains(t, executionOrder, "second")
	})

	t.Run("Error in first stops second from affecting result in MonadApFirst", func(t *testing.T) {
		secondExecuted := false

		first := Left[int](errors.New("first error"))
		second := func() (string, error) {
			secondExecuted = true
			return "test", nil
		}

		result := MonadApFirst(first, second)
		_, err := result()

		// Second still executes but error is from first
		assert.True(t, secondExecuted)
		assert.Error(t, err)
		assert.Equal(t, "first error", err.Error())
	})
}
