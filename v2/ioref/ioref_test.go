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
	"github.com/IBM/fp-go/v2/readerio"
	"github.com/stretchr/testify/assert"
)

func TestMakeIORef(t *testing.T) {
	t.Run("creates IORef with integer value", func(t *testing.T) {
		ref := MakeIORef(42)()
		assert.NotNil(t, ref)
		assert.Equal(t, 42, Read(ref)())
	})

	t.Run("creates IORef with string value", func(t *testing.T) {
		ref := MakeIORef("hello")()
		assert.NotNil(t, ref)
		assert.Equal(t, "hello", Read(ref)())
	})

	t.Run("creates IORef with slice value", func(t *testing.T) {
		slice := []int{1, 2, 3}
		ref := MakeIORef(slice)()
		assert.NotNil(t, ref)
		assert.Equal(t, slice, Read(ref)())
	})

	t.Run("creates IORef with struct value", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		person := Person{Name: "Alice", Age: 30}
		ref := MakeIORef(person)()
		assert.NotNil(t, ref)
		assert.Equal(t, person, Read(ref)())
	})

	t.Run("creates IORef with zero value", func(t *testing.T) {
		ref := MakeIORef(0)()
		assert.NotNil(t, ref)
		assert.Equal(t, 0, Read(ref)())
	})

	t.Run("creates IORef with nil pointer", func(t *testing.T) {
		var ptr *int
		ref := MakeIORef(ptr)()
		assert.NotNil(t, ref)
		assert.Nil(t, Read(ref)())
	})

	t.Run("multiple IORefs are independent", func(t *testing.T) {
		ref1 := MakeIORef(10)()
		ref2 := MakeIORef(20)()

		assert.Equal(t, 10, Read(ref1)())
		assert.Equal(t, 20, Read(ref2)())

		Write(30)(ref1)()
		assert.Equal(t, 30, Read(ref1)())
		assert.Equal(t, 20, Read(ref2)()) // ref2 unchanged
	})
}

func TestRead(t *testing.T) {
	t.Run("reads initial value", func(t *testing.T) {
		ref := MakeIORef(42)()
		value := Read(ref)()
		assert.Equal(t, 42, value)
	})

	t.Run("reads updated value", func(t *testing.T) {
		ref := MakeIORef(10)()
		Write(20)(ref)()
		value := Read(ref)()
		assert.Equal(t, 20, value)
	})

	t.Run("multiple reads return same value", func(t *testing.T) {
		ref := MakeIORef(100)()
		value1 := Read(ref)()
		value2 := Read(ref)()
		value3 := Read(ref)()
		assert.Equal(t, 100, value1)
		assert.Equal(t, 100, value2)
		assert.Equal(t, 100, value3)
	})

	t.Run("concurrent reads are thread-safe", func(t *testing.T) {
		ref := MakeIORef(42)()
		var wg sync.WaitGroup
		iterations := 100
		results := make([]int, iterations)

		for i := 0; i < iterations; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				results[idx] = Read(ref)()
			}(i)
		}

		wg.Wait()

		// All reads should return the same value
		for _, v := range results {
			assert.Equal(t, 42, v)
		}
	})

	t.Run("reads during concurrent writes", func(t *testing.T) {
		ref := MakeIORef(0)()
		var wg sync.WaitGroup
		iterations := 50

		// Start concurrent writes
		for i := 0; i < iterations; i++ {
			wg.Add(1)
			go func(val int) {
				defer wg.Done()
				Write(val)(ref)()
			}(i)
		}

		// Start concurrent reads
		for i := 0; i < iterations; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				value := Read(ref)()
				// Value should be valid (between 0 and iterations-1)
				assert.GreaterOrEqual(t, value, 0)
				assert.Less(t, value, iterations)
			}()
		}

		wg.Wait()
	})
}

func TestWrite(t *testing.T) {
	t.Run("writes new value", func(t *testing.T) {
		ref := MakeIORef(42)()
		result := Write(100)(ref)()
		assert.Equal(t, 100, result)
		assert.Equal(t, 100, Read(ref)())
	})

	t.Run("overwrites existing value", func(t *testing.T) {
		ref := MakeIORef(10)()
		Write(20)(ref)()
		Write(30)(ref)()
		assert.Equal(t, 30, Read(ref)())
	})

	t.Run("returns written value", func(t *testing.T) {
		ref := MakeIORef(0)()
		result := Write(42)(ref)()
		assert.Equal(t, 42, result)
	})

	t.Run("writes string value", func(t *testing.T) {
		ref := MakeIORef("hello")()
		result := Write("world")(ref)()
		assert.Equal(t, "world", result)
		assert.Equal(t, "world", Read(ref)())
	})

	t.Run("chained writes", func(t *testing.T) {
		ref := MakeIORef(1)()
		Write(2)(ref)()
		Write(3)(ref)()
		result := Write(4)(ref)()
		assert.Equal(t, 4, result)
		assert.Equal(t, 4, Read(ref)())
	})

	t.Run("concurrent writes are thread-safe", func(t *testing.T) {
		ref := MakeIORef(0)()
		var wg sync.WaitGroup
		iterations := 100

		for i := 0; i < iterations; i++ {
			wg.Add(1)
			go func(val int) {
				defer wg.Done()
				Write(val)(ref)()
			}(i)
		}

		wg.Wait()

		// Final value should be one of the written values
		finalValue := Read(ref)()
		assert.GreaterOrEqual(t, finalValue, 0)
		assert.Less(t, finalValue, iterations)
	})

	t.Run("write with zero value", func(t *testing.T) {
		ref := MakeIORef(42)()
		Write(0)(ref)()
		assert.Equal(t, 0, Read(ref)())
	})
}

func TestModify(t *testing.T) {
	t.Run("modifies value with simple function", func(t *testing.T) {
		ref := MakeIORef(10)()
		result := Modify(func(x int) int { return x * 2 })(ref)()
		assert.Equal(t, 20, result)
		assert.Equal(t, 20, Read(ref)())
	})

	t.Run("modifies with addition", func(t *testing.T) {
		ref := MakeIORef(5)()
		Modify(func(x int) int { return x + 10 })(ref)()
		assert.Equal(t, 15, Read(ref)())
	})

	t.Run("modifies string value", func(t *testing.T) {
		ref := MakeIORef("hello")()
		result := Modify(func(s string) string { return s + " world" })(ref)()
		assert.Equal(t, "hello world", result)
		assert.Equal(t, "hello world", Read(ref)())
	})

	t.Run("chained modifications", func(t *testing.T) {
		ref := MakeIORef(2)()
		Modify(func(x int) int { return x * 3 })(ref)() // 6
		Modify(func(x int) int { return x + 4 })(ref)() // 10
		result := Modify(func(x int) int { return x / 2 })(ref)()
		assert.Equal(t, 5, result)
		assert.Equal(t, 5, Read(ref)())
	})

	t.Run("concurrent modifications are thread-safe", func(t *testing.T) {
		ref := MakeIORef(0)()
		var wg sync.WaitGroup
		iterations := 100

		for i := 0; i < iterations; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				Modify(func(x int) int { return x + 1 })(ref)()
			}()
		}

		wg.Wait()
		assert.Equal(t, iterations, Read(ref)())
	})

	t.Run("modify with identity function", func(t *testing.T) {
		ref := MakeIORef(42)()
		result := Modify(func(x int) int { return x })(ref)()
		assert.Equal(t, 42, result)
		assert.Equal(t, 42, Read(ref)())
	})

	t.Run("modify returns new value", func(t *testing.T) {
		ref := MakeIORef(100)()
		result := Modify(func(x int) int { return x - 50 })(ref)()
		assert.Equal(t, 50, result)
	})
}

func TestModifyWithResult(t *testing.T) {
	t.Run("modifies and returns old value", func(t *testing.T) {
		ref := MakeIORef(42)()
		oldValue := ModifyWithResult(func(x int) pair.Pair[int, int] {
			return pair.MakePair(x+1, x)
		})(ref)()
		assert.Equal(t, 42, oldValue)
		assert.Equal(t, 43, Read(ref)())
	})

	t.Run("swaps value and returns old", func(t *testing.T) {
		ref := MakeIORef(100)()
		oldValue := ModifyWithResult(func(x int) pair.Pair[int, int] {
			return pair.MakePair(200, x)
		})(ref)()
		assert.Equal(t, 100, oldValue)
		assert.Equal(t, 200, Read(ref)())
	})

	t.Run("returns different type", func(t *testing.T) {
		ref := MakeIORef(42)()
		message := ModifyWithResult(func(x int) pair.Pair[int, string] {
			return pair.MakePair(x*2, fmt.Sprintf("doubled from %d", x))
		})(ref)()
		assert.Equal(t, "doubled from 42", message)
		assert.Equal(t, 84, Read(ref)())
	})

	t.Run("computes result based on old value", func(t *testing.T) {
		ref := MakeIORef(10)()
		wasPositive := ModifyWithResult(func(x int) pair.Pair[int, bool] {
			return pair.MakePair(x+5, x > 0)
		})(ref)()
		assert.True(t, wasPositive)
		assert.Equal(t, 15, Read(ref)())
	})

	t.Run("chained modifications with results", func(t *testing.T) {
		ref := MakeIORef(5)()
		result1 := ModifyWithResult(func(x int) pair.Pair[int, int] {
			return pair.MakePair(x*2, x)
		})(ref)()
		result2 := ModifyWithResult(func(x int) pair.Pair[int, int] {
			return pair.MakePair(x+10, x)
		})(ref)()
		assert.Equal(t, 5, result1)
		assert.Equal(t, 10, result2)
		assert.Equal(t, 20, Read(ref)())
	})

	t.Run("concurrent modifications with results are thread-safe", func(t *testing.T) {
		ref := MakeIORef(0)()
		var wg sync.WaitGroup
		iterations := 100
		results := make([]int, iterations)

		for i := 0; i < iterations; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				oldValue := ModifyWithResult(func(x int) pair.Pair[int, int] {
					return pair.MakePair(x+1, x)
				})(ref)()
				results[idx] = oldValue
			}(i)
		}

		wg.Wait()
		assert.Equal(t, iterations, Read(ref)())

		// All old values should be unique
		seen := make(map[int]bool)
		for _, v := range results {
			assert.False(t, seen[v])
			seen[v] = true
		}
	})

	t.Run("extract and replace pattern", func(t *testing.T) {
		ref := MakeIORef([]int{1, 2, 3})()
		first := ModifyWithResult(func(xs []int) pair.Pair[[]int, int] {
			if len(xs) == 0 {
				return pair.MakePair(xs, 0)
			}
			return pair.MakePair(xs[1:], xs[0])
		})(ref)()
		assert.Equal(t, 1, first)
		assert.Equal(t, []int{2, 3}, Read(ref)())
	})
}

func TestModifyReaderIOK(t *testing.T) {
	type Config struct {
		multiplier int
	}

	t.Run("modifies with environment", func(t *testing.T) {
		ref := MakeIORef(10)()
		config := Config{multiplier: 5}

		result := ModifyReaderIOK(func(x int) readerio.ReaderIO[Config, int] {
			return func(cfg Config) io.IO[int] {
				return io.Of(x * cfg.multiplier)
			}
		})(ref)(config)()

		assert.Equal(t, 50, result)
		assert.Equal(t, 50, Read(ref)())
	})

	t.Run("uses environment for computation", func(t *testing.T) {
		ref := MakeIORef(100)()
		config := Config{multiplier: 2}

		result := ModifyReaderIOK(func(x int) readerio.ReaderIO[Config, int] {
			return func(cfg Config) io.IO[int] {
				return func() int {
					return x / cfg.multiplier
				}
			}
		})(ref)(config)()

		assert.Equal(t, 50, result)
		assert.Equal(t, 50, Read(ref)())
	})

	t.Run("chained modifications with different configs", func(t *testing.T) {
		ref := MakeIORef(10)()
		config1 := Config{multiplier: 2}
		config2 := Config{multiplier: 3}

		ModifyReaderIOK(func(x int) readerio.ReaderIO[Config, int] {
			return func(cfg Config) io.IO[int] {
				return io.Of(x * cfg.multiplier)
			}
		})(ref)(config1)()

		result := ModifyReaderIOK(func(x int) readerio.ReaderIO[Config, int] {
			return func(cfg Config) io.IO[int] {
				return io.Of(x + cfg.multiplier)
			}
		})(ref)(config2)()

		assert.Equal(t, 23, result) // (10 * 2) + 3
		assert.Equal(t, 23, Read(ref)())
	})

	t.Run("concurrent modifications with environment are thread-safe", func(t *testing.T) {
		ref := MakeIORef(0)()
		config := Config{multiplier: 1}
		var wg sync.WaitGroup
		iterations := 100

		for i := 0; i < iterations; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				ModifyReaderIOK(func(x int) readerio.ReaderIO[Config, int] {
					return func(cfg Config) io.IO[int] {
						return io.Of(x + cfg.multiplier)
					}
				})(ref)(config)()
			}()
		}

		wg.Wait()
		assert.Equal(t, iterations, Read(ref)())
	})

	t.Run("environment provides configuration", func(t *testing.T) {
		type Settings struct {
			prefix string
		}
		ref := MakeIORef("world")()
		settings := Settings{prefix: "hello "}

		result := ModifyReaderIOK(func(s string) readerio.ReaderIO[Settings, string] {
			return func(cfg Settings) io.IO[string] {
				return io.Of(cfg.prefix + s)
			}
		})(ref)(settings)()

		assert.Equal(t, "hello world", result)
		assert.Equal(t, "hello world", Read(ref)())
	})
}

func TestModifyReaderIOKWithResult(t *testing.T) {
	type Config struct {
		logEnabled bool
		multiplier int
	}

	t.Run("modifies with environment and returns result", func(t *testing.T) {
		ref := MakeIORef(42)()
		config := Config{logEnabled: false, multiplier: 2}

		oldValue := ModifyReaderIOKWithResult(func(x int) readerio.ReaderIO[Config, pair.Pair[int, int]] {
			return func(cfg Config) io.IO[pair.Pair[int, int]] {
				return io.Of(pair.MakePair(x*cfg.multiplier, x))
			}
		})(ref)(config)()

		assert.Equal(t, 42, oldValue)
		assert.Equal(t, 84, Read(ref)())
	})

	t.Run("returns different type based on environment", func(t *testing.T) {
		ref := MakeIORef(10)()
		config := Config{logEnabled: true, multiplier: 3}

		message := ModifyReaderIOKWithResult(func(x int) readerio.ReaderIO[Config, pair.Pair[int, string]] {
			return func(cfg Config) io.IO[pair.Pair[int, string]] {
				return func() pair.Pair[int, string] {
					newVal := x * cfg.multiplier
					msg := fmt.Sprintf("multiplied %d by %d", x, cfg.multiplier)
					return pair.MakePair(newVal, msg)
				}
			}
		})(ref)(config)()

		assert.Equal(t, "multiplied 10 by 3", message)
		assert.Equal(t, 30, Read(ref)())
	})

	t.Run("conditional logic based on environment", func(t *testing.T) {
		ref := MakeIORef(-10)()
		config := Config{logEnabled: true, multiplier: 2}

		message := ModifyReaderIOKWithResult(func(x int) readerio.ReaderIO[Config, pair.Pair[int, string]] {
			return func(cfg Config) io.IO[pair.Pair[int, string]] {
				return func() pair.Pair[int, string] {
					if x < 0 {
						return pair.MakePair(0, "reset negative value")
					}
					return pair.MakePair(x*cfg.multiplier, "multiplied positive value")
				}
			}
		})(ref)(config)()

		assert.Equal(t, "reset negative value", message)
		assert.Equal(t, 0, Read(ref)())
	})

	t.Run("chained modifications with results", func(t *testing.T) {
		ref := MakeIORef(5)()
		config := Config{logEnabled: false, multiplier: 2}

		result1 := ModifyReaderIOKWithResult(func(x int) readerio.ReaderIO[Config, pair.Pair[int, int]] {
			return func(cfg Config) io.IO[pair.Pair[int, int]] {
				return io.Of(pair.MakePair(x*cfg.multiplier, x))
			}
		})(ref)(config)()

		result2 := ModifyReaderIOKWithResult(func(x int) readerio.ReaderIO[Config, pair.Pair[int, int]] {
			return func(cfg Config) io.IO[pair.Pair[int, int]] {
				return io.Of(pair.MakePair(x+cfg.multiplier, x))
			}
		})(ref)(config)()

		assert.Equal(t, 5, result1)
		assert.Equal(t, 10, result2)
		assert.Equal(t, 12, Read(ref)())
	})

	t.Run("concurrent modifications with environment are thread-safe", func(t *testing.T) {
		ref := MakeIORef(0)()
		config := Config{logEnabled: false, multiplier: 1}
		var wg sync.WaitGroup
		iterations := 100
		results := make([]int, iterations)

		for i := 0; i < iterations; i++ {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				oldValue := ModifyReaderIOKWithResult(func(x int) readerio.ReaderIO[Config, pair.Pair[int, int]] {
					return func(cfg Config) io.IO[pair.Pair[int, int]] {
						return io.Of(pair.MakePair(x+cfg.multiplier, x))
					}
				})(ref)(config)()
				results[idx] = oldValue
			}(i)
		}

		wg.Wait()
		assert.Equal(t, iterations, Read(ref)())

		// All old values should be unique
		seen := make(map[int]bool)
		for _, v := range results {
			assert.False(t, seen[v])
			seen[v] = true
		}
	})

	t.Run("environment provides validation rules", func(t *testing.T) {
		type ValidationConfig struct {
			maxValue int
		}
		ref := MakeIORef(100)()
		config := ValidationConfig{maxValue: 50}

		message := ModifyReaderIOKWithResult(func(x int) readerio.ReaderIO[ValidationConfig, pair.Pair[int, string]] {
			return func(cfg ValidationConfig) io.IO[pair.Pair[int, string]] {
				return func() pair.Pair[int, string] {
					if x > cfg.maxValue {
						return pair.MakePair(cfg.maxValue, fmt.Sprintf("capped at %d", cfg.maxValue))
					}
					return pair.MakePair(x, "value within limits")
				}
			}
		})(ref)(config)()

		assert.Equal(t, "capped at 50", message)
		assert.Equal(t, 50, Read(ref)())
	})
}

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
