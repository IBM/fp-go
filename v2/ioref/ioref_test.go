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

package ioref

import (
	"fmt"
	"sync"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/stretchr/testify/assert"
)

func TestModifyIOK(t *testing.T) {
	t.Run("basic modification with IO effect", func(t *testing.T) {
		ref := MakeIORef(42)()

		// Double the value using ModifyIOK
		newValue := ModifyIOK(func(x int) io.IO[int] {
			return io.Of(x * 2)
		})(ref)()

		assert.Equal(t, 84, newValue)
		assert.Equal(t, 84, Read(ref)())
	})

	t.Run("modification with side effects", func(t *testing.T) {
		ref := MakeIORef(10)()
		var sideEffect int

		// Modify with a side effect
		newValue := ModifyIOK(func(x int) io.IO[int] {
			return func() int {
				sideEffect = x // Capture old value
				return x + 5
			}
		})(ref)()

		assert.Equal(t, 15, newValue)
		assert.Equal(t, 10, sideEffect)
		assert.Equal(t, 15, Read(ref)())
	})

	t.Run("chained modifications", func(t *testing.T) {
		ref := MakeIORef(5)()

		// First modification: add 10
		ModifyIOK(func(x int) io.IO[int] {
			return io.Of(x + 10)
		})(ref)()

		// Second modification: multiply by 2
		result := ModifyIOK(func(x int) io.IO[int] {
			return io.Of(x * 2)
		})(ref)()

		assert.Equal(t, 30, result)
		assert.Equal(t, 30, Read(ref)())
	})

	t.Run("concurrent modifications are thread-safe", func(t *testing.T) {
		ref := MakeIORef(0)()
		var wg sync.WaitGroup
		iterations := 100

		// Increment concurrently
		for i := 0; i < iterations; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				ModifyIOK(func(x int) io.IO[int] {
					return io.Of(x + 1)
				})(ref)()
			}()
		}

		wg.Wait()
		assert.Equal(t, iterations, Read(ref)())
	})

	t.Run("modification with string type", func(t *testing.T) {
		ref := MakeIORef("hello")()

		newValue := ModifyIOK(func(s string) io.IO[string] {
			return io.Of(s + " world")
		})(ref)()

		assert.Equal(t, "hello world", newValue)
		assert.Equal(t, "hello world", Read(ref)())
	})

	t.Run("modification returns new value", func(t *testing.T) {
		ref := MakeIORef(100)()

		result := ModifyIOK(func(x int) io.IO[int] {
			return io.Of(x / 2)
		})(ref)()

		// ModifyIOK returns the new value
		assert.Equal(t, 50, result)
		assert.Equal(t, 50, Read(ref)())
	})

	t.Run("modification with complex IO computation", func(t *testing.T) {
		ref := MakeIORef(3)()

		// Use a more complex IO computation
		newValue := ModifyIOK(func(x int) io.IO[int] {
			return F.Pipe1(
				io.Of(x),
				io.Map(func(n int) int { return n * n }),
			)
		})(ref)()

		assert.Equal(t, 9, newValue)
		assert.Equal(t, 9, Read(ref)())
	})
}

func TestModifyIOKWithResult(t *testing.T) {
	t.Run("basic modification with result", func(t *testing.T) {
		ref := MakeIORef(42)()

		// Increment and return old value
		oldValue := ModifyIOKWithResult(func(x int) io.IO[pair.Pair[int, int]] {
			return io.Of(pair.MakePair(x+1, x))
		})(ref)()

		assert.Equal(t, 42, oldValue)
		assert.Equal(t, 43, Read(ref)())
	})

	t.Run("swap and return old value", func(t *testing.T) {
		ref := MakeIORef(100)()

		oldValue := ModifyIOKWithResult(func(x int) io.IO[pair.Pair[int, int]] {
			return io.Of(pair.MakePair(200, x))
		})(ref)()

		assert.Equal(t, 100, oldValue)
		assert.Equal(t, 200, Read(ref)())
	})

	t.Run("modification with different result type", func(t *testing.T) {
		ref := MakeIORef(42)()

		// Double the value and return a message
		message := ModifyIOKWithResult(func(x int) io.IO[pair.Pair[int, string]] {
			return io.Of(pair.MakePair(x*2, fmt.Sprintf("doubled from %d", x)))
		})(ref)()

		assert.Equal(t, "doubled from 42", message)
		assert.Equal(t, 84, Read(ref)())
	})

	t.Run("modification with side effects in IO", func(t *testing.T) {
		ref := MakeIORef(10)()
		var sideEffect string

		result := ModifyIOKWithResult(func(x int) io.IO[pair.Pair[int, bool]] {
			return func() pair.Pair[int, bool] {
				sideEffect = fmt.Sprintf("processing %d", x)
				return pair.MakePair(x+5, x > 5)
			}
		})(ref)()

		assert.True(t, result)
		assert.Equal(t, "processing 10", sideEffect)
		assert.Equal(t, 15, Read(ref)())
	})

	t.Run("chained modifications with results", func(t *testing.T) {
		ref := MakeIORef(5)()

		// First modification
		result1 := ModifyIOKWithResult(func(x int) io.IO[pair.Pair[int, int]] {
			return io.Of(pair.MakePair(x*2, x))
		})(ref)()

		// Second modification
		result2 := ModifyIOKWithResult(func(x int) io.IO[pair.Pair[int, int]] {
			return io.Of(pair.MakePair(x+10, x))
		})(ref)()

		assert.Equal(t, 5, result1)      // Original value
		assert.Equal(t, 10, result2)     // After first modification
		assert.Equal(t, 20, Read(ref)()) // After both modifications
	})

	t.Run("concurrent modifications with results are thread-safe", func(t *testing.T) {
		ref := MakeIORef(0)()
		var wg sync.WaitGroup
		iterations := 100
		results := make([]int, iterations)

		// Increment concurrently and collect old values
		for i := 0; i < iterations; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				oldValue := ModifyIOKWithResult(func(x int) io.IO[pair.Pair[int, int]] {
					return io.Of(pair.MakePair(x+1, x))
				})(ref)()
				results[idx] = oldValue
			}(i)
		}

		wg.Wait()

		// Final value should be iterations
		assert.Equal(t, iterations, Read(ref)())

		// All old values should be unique and in range [0, iterations)
		seen := make(map[int]bool)
		for _, v := range results {
			assert.False(t, seen[v], "duplicate old value: %d", v)
			assert.GreaterOrEqual(t, v, 0)
			assert.Less(t, v, iterations)
			seen[v] = true
		}
	})

	t.Run("modification with string types", func(t *testing.T) {
		ref := MakeIORef("hello")()

		length := ModifyIOKWithResult(func(s string) io.IO[pair.Pair[string, int]] {
			return io.Of(pair.MakePair(s+" world", len(s)))
		})(ref)()

		assert.Equal(t, 5, length)
		assert.Equal(t, "hello world", Read(ref)())
	})

	t.Run("modification with validation logic", func(t *testing.T) {
		ref := MakeIORef(-10)()

		message := ModifyIOKWithResult(func(x int) io.IO[pair.Pair[int, string]] {
			return func() pair.Pair[int, string] {
				if x < 0 {
					return pair.MakePair(0, "reset negative value")
				}
				return pair.MakePair(x*2, "doubled positive value")
			}
		})(ref)()

		assert.Equal(t, "reset negative value", message)
		assert.Equal(t, 0, Read(ref)())
	})

	t.Run("modification with complex IO computation", func(t *testing.T) {
		ref := MakeIORef(5)()

		result := ModifyIOKWithResult(func(x int) io.IO[pair.Pair[int, string]] {
			return F.Pipe1(
				io.Of(x),
				io.Map(func(n int) pair.Pair[int, string] {
					squared := n * n
					return pair.MakePair(squared, fmt.Sprintf("%d squared is %d", n, squared))
				}),
			)
		})(ref)()

		assert.Equal(t, "5 squared is 25", result)
		assert.Equal(t, 25, Read(ref)())
	})

	t.Run("extract and replace pattern", func(t *testing.T) {
		ref := MakeIORef([]int{1, 2, 3})()

		// Extract first element and remove it from the slice
		first := ModifyIOKWithResult(func(xs []int) io.IO[pair.Pair[[]int, int]] {
			return func() pair.Pair[[]int, int] {
				if len(xs) == 0 {
					return pair.MakePair(xs, 0)
				}
				return pair.MakePair(xs[1:], xs[0])
			}
		})(ref)()

		assert.Equal(t, 1, first)
		assert.Equal(t, []int{2, 3}, Read(ref)())
	})
}

func TestModifyIOKIntegration(t *testing.T) {
	t.Run("ModifyIOK integrates with Modify", func(t *testing.T) {
		ref := MakeIORef(10)()

		// Use Modify (which internally uses ModifyIOK)
		result1 := Modify(N.Mul(2))(ref)()

		assert.Equal(t, 20, result1)

		// Use ModifyIOK directly
		result2 := ModifyIOK(func(x int) io.IO[int] {
			return io.Of(x + 5)
		})(ref)()

		assert.Equal(t, 25, result2)
		assert.Equal(t, 25, Read(ref)())
	})
}

func TestModifyIOKWithResultIntegration(t *testing.T) {
	t.Run("ModifyIOKWithResult integrates with ModifyWithResult", func(t *testing.T) {
		ref := MakeIORef(10)()

		// Use ModifyWithResult (which internally uses ModifyIOKWithResult)
		result1 := ModifyWithResult(func(x int) pair.Pair[int, int] {
			return pair.MakePair(x*2, x)
		})(ref)()

		assert.Equal(t, 10, result1)
		assert.Equal(t, 20, Read(ref)())

		// Use ModifyIOKWithResult directly
		result2 := ModifyIOKWithResult(func(x int) io.IO[pair.Pair[int, string]] {
			return io.Of(pair.MakePair(x+5, fmt.Sprintf("was %d", x)))
		})(ref)()

		assert.Equal(t, "was 20", result2)
		assert.Equal(t, 25, Read(ref)())
	})
}
