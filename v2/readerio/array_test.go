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

package readerio

import (
	"context"
	"fmt"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/stretchr/testify/assert"
)

func TestTraverseArray(t *testing.T) {
	ctx := context.Background()

	t.Run("empty array", func(t *testing.T) {
		f := TraverseArray(func(a string) ReaderIO[context.Context, string] {
			return Of[context.Context](a + a)
		})
		assert.Equal(t, A.Empty[string](), F.Pipe1(A.Empty[string](), f)(ctx)())
	})

	t.Run("successful transformation", func(t *testing.T) {
		f := TraverseArray(func(a string) ReaderIO[context.Context, string] {
			return Of[context.Context](a + a)
		})
		assert.Equal(t, []string{"aa", "bb"}, F.Pipe1([]string{"a", "b"}, f)(ctx)())
	})

	t.Run("transformation with context", func(t *testing.T) {
		type Config struct {
			prefix string
		}
		cfg := Config{prefix: "pre_"}

		f := TraverseArray(func(a string) ReaderIO[Config, string] {
			return func(c Config) func() string {
				return func() string { return c.prefix + a }
			}
		})
		result := f([]string{"a", "b"})(cfg)()
		assert.Equal(t, []string{"pre_a", "pre_b"}, result)
	})

	t.Run("large array performance", func(t *testing.T) {
		input := A.MakeBy(1000, func(i int) int { return i })
		f := TraverseArray(func(a int) ReaderIO[context.Context, int] {
			return Of[context.Context](a * 2)
		})
		result := f(input)(ctx)()
		assert.Equal(t, 1000, len(result))
		assert.Equal(t, 0, result[0])
		assert.Equal(t, 1998, result[999])
	})
}

func TestTraverseArrayWithIndex(t *testing.T) {
	ctx := context.Background()

	t.Run("empty array", func(t *testing.T) {
		f := TraverseArrayWithIndex(func(i int, a string) ReaderIO[context.Context, string] {
			return Of[context.Context](fmt.Sprintf("%d:%s", i, a))
		})
		result := f([]string{})(ctx)()
		assert.Equal(t, []string{}, result)
	})

	t.Run("transformation with index", func(t *testing.T) {
		f := TraverseArrayWithIndex(func(i int, a string) ReaderIO[context.Context, string] {
			return Of[context.Context](fmt.Sprintf("%d:%s", i, a))
		})
		result := f([]string{"a", "b", "c"})(ctx)()
		assert.Equal(t, []string{"0:a", "1:b", "2:c"}, result)
	})

	t.Run("index-dependent transformation", func(t *testing.T) {
		f := TraverseArrayWithIndex(func(i int, a int) ReaderIO[context.Context, int] {
			return Of[context.Context](a + i)
		})
		result := f([]int{10, 20, 30})(ctx)()
		assert.Equal(t, []int{10, 21, 32}, result)
	})
}

func TestSequenceArray(t *testing.T) {
	ctx := context.Background()

	t.Run("empty array", func(t *testing.T) {
		computations := []ReaderIO[context.Context, int]{}
		result := SequenceArray(computations)(ctx)()
		assert.Equal(t, []int{}, result)
	})

	t.Run("multiple computations", func(t *testing.T) {
		computations := []ReaderIO[context.Context, int]{
			Of[context.Context](1),
			Of[context.Context](2),
			Of[context.Context](3),
		}
		result := SequenceArray(computations)(ctx)()
		assert.Equal(t, []int{1, 2, 3}, result)
	})

	t.Run("computations with context", func(t *testing.T) {
		type Config struct {
			multiplier int
		}
		cfg := Config{multiplier: 10}

		computations := []ReaderIO[Config, int]{
			func(c Config) func() int { return func() int { return 1 * c.multiplier } },
			func(c Config) func() int { return func() int { return 2 * c.multiplier } },
			func(c Config) func() int { return func() int { return 3 * c.multiplier } },
		}
		result := SequenceArray(computations)(cfg)()
		assert.Equal(t, []int{10, 20, 30}, result)
	})
}
