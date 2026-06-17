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
	"sync"
	"sync/atomic"
	"testing"

	"github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/stretchr/testify/assert"
)

type memoizeTestContext struct {
	Key   string
	Value string
}

func TestMemoize_Success(t *testing.T) {
	t.Run("computes once per comparable environment", func(t *testing.T) {
		var calls atomic.Int32

		eff := Memoize(func(cfg memoizeTestContext) ReaderIOResult[int] {
			calls.Add(1)
			return readerioresult.Of(len(cfg.Value))
		})

		result1, err1 := runEffect(eff, memoizeTestContext{Key: "alpha", Value: "first"})
		result2, err2 := runEffect(eff, memoizeTestContext{Key: "alpha", Value: "first"})
		result3, err3 := runEffect(eff, memoizeTestContext{Key: "beta", Value: "second"})

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NoError(t, err3)
		assert.Equal(t, 5, result1)
		assert.Equal(t, 5, result2)
		assert.Equal(t, 6, result3)
		assert.Equal(t, int32(2), calls.Load())
	})

	t.Run("reuses cached successful result across repeated runs", func(t *testing.T) {
		var calls atomic.Int32

		eff := Memoize(func(cfg memoizeTestContext) ReaderIOResult[string] {
			calls.Add(1)
			return readerioresult.Of(cfg.Value + "-computed")
		})

		result1, err1 := runEffect(eff, memoizeTestContext{Key: "same", Value: "payload"})
		result2, err2 := runEffect(eff, memoizeTestContext{Key: "same", Value: "payload"})

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, "payload-computed", result1)
		assert.Equal(t, "payload-computed", result2)
		assert.Equal(t, int32(1), calls.Load())
	})
}

func TestMemoize_Failure(t *testing.T) {
	t.Run("memoizes failures per environment", func(t *testing.T) {
		expectedErr := assert.AnError
		var calls atomic.Int32

		eff := Memoize(func(cfg memoizeTestContext) ReaderIOResult[int] {
			calls.Add(1)
			return readerioresult.Left[int](expectedErr)
		})

		_, err1 := runEffect(eff, memoizeTestContext{Key: "same", Value: "ignored"})
		_, err2 := runEffect(eff, memoizeTestContext{Key: "same", Value: "ignored"})

		assert.ErrorIs(t, err1, expectedErr)
		assert.ErrorIs(t, err2, expectedErr)
		assert.Equal(t, int32(1), calls.Load())
	})
}

func TestMemoize_EdgeCases(t *testing.T) {
	t.Run("treats distinct comparable environments as distinct cache keys", func(t *testing.T) {
		var calls atomic.Int32

		eff := Memoize(func(cfg memoizeTestContext) ReaderIOResult[string] {
			calls.Add(1)
			return readerioresult.Of(cfg.Key + ":" + cfg.Value)
		})

		result1, err1 := runEffect(eff, memoizeTestContext{Key: "shared", Value: "one"})
		result2, err2 := runEffect(eff, memoizeTestContext{Key: "shared", Value: "two"})

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, "shared:one", result1)
		assert.Equal(t, "shared:two", result2)
		assert.Equal(t, int32(2), calls.Load())
	})
}

func TestMemoize_Integration(t *testing.T) {
	t.Run("preserves memoized value through map", func(t *testing.T) {
		var calls atomic.Int32

		base := Memoize(func(cfg memoizeTestContext) ReaderIOResult[int] {
			calls.Add(1)
			return readerioresult.Of(len(cfg.Value))
		})

		eff := Map[memoizeTestContext](func(n int) string {
			return "size"
		})(base)

		result1, err1 := runEffect(eff, memoizeTestContext{Key: "same", Value: "hello"})
		result2, err2 := runEffect(eff, memoizeTestContext{Key: "same", Value: "hello"})

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, "size", result1)
		assert.Equal(t, "size", result2)
		assert.Equal(t, int32(1), calls.Load())
	})
}

func TestContramapMemoize_Success(t *testing.T) {
	t.Run("memoizes by derived key", func(t *testing.T) {
		var calls atomic.Int32

		eff := ContramapMemoize[int](func(cfg memoizeTestContext) string {
			return cfg.Key
		})(func(cfg memoizeTestContext) ReaderIOResult[int] {
			calls.Add(1)
			return readerioresult.Of(len(cfg.Value))
		})

		result1, err1 := runEffect(eff, memoizeTestContext{Key: "shared", Value: "first"})
		result2, err2 := runEffect(eff, memoizeTestContext{Key: "shared", Value: "second"})
		result3, err3 := runEffect(eff, memoizeTestContext{Key: "other", Value: "third"})

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NoError(t, err3)
		assert.Equal(t, 5, result1)
		assert.Equal(t, 5, result2)
		assert.Equal(t, 5, result3)
		assert.Equal(t, int32(2), calls.Load())
	})
}

func TestContramapMemoize_Failure(t *testing.T) {
	t.Run("memoizes failures by derived key", func(t *testing.T) {
		expectedErr := assert.AnError
		var calls atomic.Int32

		eff := ContramapMemoize[int](func(cfg memoizeTestContext) string {
			return cfg.Key
		})(func(cfg memoizeTestContext) ReaderIOResult[int] {
			calls.Add(1)
			return readerioresult.Left[int](expectedErr)
		})

		_, err1 := runEffect(eff, memoizeTestContext{Key: "shared", Value: "first"})
		_, err2 := runEffect(eff, memoizeTestContext{Key: "shared", Value: "second"})

		assert.ErrorIs(t, err1, expectedErr)
		assert.ErrorIs(t, err2, expectedErr)
		assert.Equal(t, int32(1), calls.Load())
	})
}

func TestContramapMemoize_EdgeCases(t *testing.T) {
	t.Run("supports non-comparable environments through comparable derived keys", func(t *testing.T) {
		type nonComparableContext struct {
			Key    string
			Values []string
		}

		var calls atomic.Int32

		eff := ContramapMemoize[int](func(cfg nonComparableContext) string {
			return cfg.Key
		})(func(cfg nonComparableContext) ReaderIOResult[int] {
			calls.Add(1)
			return readerioresult.Of(len(cfg.Values))
		})

		result1, err1 := runEffect(eff, nonComparableContext{Key: "shared", Values: []string{"a", "b"}})
		result2, err2 := runEffect(eff, nonComparableContext{Key: "shared", Values: []string{"x", "y", "z"}})
		result3, err3 := runEffect(eff, nonComparableContext{Key: "other", Values: []string{"q"}})

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.NoError(t, err3)
		assert.Equal(t, 2, result1)
		assert.Equal(t, 2, result2)
		assert.Equal(t, 1, result3)
		assert.Equal(t, int32(2), calls.Load())
	})
}

func TestContramapMemoize_Integration(t *testing.T) {
	t.Run("computes once under concurrent access for the same derived key", func(t *testing.T) {
		var calls atomic.Int32

		eff := ContramapMemoize[int](func(cfg memoizeTestContext) string {
			return cfg.Key
		})(func(cfg memoizeTestContext) ReaderIOResult[int] {
			calls.Add(1)
			return readerioresult.Of(len(cfg.Value))
		})

		const workers = 10
		results := make([]int, workers)
		errs := make([]error, workers)

		var wg sync.WaitGroup
		wg.Add(workers)

		for i := range workers {
			go func(idx int) {
				defer wg.Done()
				results[idx], errs[idx] = runEffect(eff, memoizeTestContext{Key: "shared", Value: "first"})
			}(i)
		}

		wg.Wait()

		for i := range workers {
			assert.NoError(t, errs[i])
			assert.Equal(t, 5, results[i])
		}
		assert.Equal(t, int32(1), calls.Load())
	})
}
