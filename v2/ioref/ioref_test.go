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
	"sync"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/stretchr/testify/assert"
)

func TestMakeIORef(t *testing.T) {
	t.Run("creates IORef with initial value", func(t *testing.T) {
		value := F.Pipe2(
			42,
			MakeIORef[int],
			io.Chain(Read[int]),
		)()
		assert.Equal(t, 42, value)
	})

	t.Run("creates IORef with string value", func(t *testing.T) {
		value := F.Pipe2(
			"hello",
			MakeIORef[string],
			io.Chain(Read[string]),
		)()
		assert.Equal(t, "hello", value)
	})

	t.Run("creates IORef with struct value", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		person := Person{Name: "Alice", Age: 30}
		value := F.Pipe2(
			person,
			MakeIORef[Person],
			io.Chain(Read[Person]),
		)()
		assert.Equal(t, person, value)
	})

	t.Run("creates IORef with slice value", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		value := F.Pipe2(
			slice,
			MakeIORef[[]int],
			io.Chain(Read[[]int]),
		)()
		assert.Equal(t, slice, value)
	})

	t.Run("creates IORef with zero value", func(t *testing.T) {
		value := F.Pipe2(
			0,
			MakeIORef[int],
			io.Chain(Read[int]),
		)()
		assert.Equal(t, 0, value)
	})
}

func TestRead(t *testing.T) {
	t.Run("reads the current value", func(t *testing.T) {
		value := F.Pipe2(
			100,
			MakeIORef[int],
			io.Chain(Read[int]),
		)()
		assert.Equal(t, 100, value)
	})

	t.Run("reads value multiple times", func(t *testing.T) {
		ref := MakeIORef(42)()
		value1 := Read(ref)()
		value2 := Read(ref)()
		value3 := Read(ref)()
		assert.Equal(t, 42, value1)
		assert.Equal(t, 42, value2)
		assert.Equal(t, 42, value3)
	})

	t.Run("reads updated value", func(t *testing.T) {
		value := F.Pipe3(
			10,
			MakeIORef[int],
			io.Chain(Write(20)),
			io.Chain(Read[int]),
		)()
		assert.Equal(t, 20, value)
	})
}

func TestWrite(t *testing.T) {
	t.Run("writes a new value", func(t *testing.T) {
		value := F.Pipe3(
			42,
			MakeIORef[int],
			io.Chain(Write(100)),
			io.Chain(Read[int]),
		)()
		assert.Equal(t, 100, value)
	})

	t.Run("overwrites previous value", func(t *testing.T) {
		value := F.Pipe5(
			1,
			MakeIORef[int],
			io.Chain(Write(2)),
			io.Chain(Write(3)),
			io.Chain(Write(4)),
			io.Chain(Read[int]),
		)()
		assert.Equal(t, 4, value)
	})

	t.Run("returns the IORef for chaining", func(t *testing.T) {
		ref := F.Pipe2(
			10,
			MakeIORef[int],
			io.Chain(Write(20)),
		)()
		assert.NotNil(t, ref)
		assert.Equal(t, 20, Read(ref)())
	})

	t.Run("writes different types", func(t *testing.T) {
		strValue := F.Pipe3(
			"initial",
			MakeIORef[string],
			io.Chain(Write("updated")),
			io.Chain(Read[string]),
		)()
		assert.Equal(t, "updated", strValue)

		boolValue := F.Pipe3(
			false,
			MakeIORef[bool],
			io.Chain(Write(true)),
			io.Chain(Read[bool]),
		)()
		assert.Equal(t, true, boolValue)
	})
}

func TestModify(t *testing.T) {
	t.Run("modifies value with function", func(t *testing.T) {
		value := F.Pipe3(
			10,
			MakeIORef[int],
			io.Chain(Modify(N.Mul(2))),
			io.Chain(Read[int]),
		)()
		assert.Equal(t, 20, value)
	})

	t.Run("chains multiple modifications", func(t *testing.T) {
		value := F.Pipe5(
			5,
			MakeIORef[int],
			io.Chain(Modify(N.Add(10))),
			io.Chain(Modify(N.Mul(2))),
			io.Chain(Modify(N.Sub(5))),
			io.Chain(Read[int]),
		)()
		assert.Equal(t, 25, value) // (5 + 10) * 2 - 5 = 25
	})

	t.Run("modifies string value", func(t *testing.T) {
		value := F.Pipe3(
			"hello",
			MakeIORef[string],
			io.Chain(Modify(func(s string) string { return s + " world" })),
			io.Chain(Read[string]),
		)()
		assert.Equal(t, "hello world", value)
	})

	t.Run("returns the IORef for chaining", func(t *testing.T) {
		ref := F.Pipe2(
			42,
			MakeIORef[int],
			io.Chain(Modify(N.Add(1))),
		)()
		assert.NotNil(t, ref)
		assert.Equal(t, 43, Read(ref)())
	})

	t.Run("modifies slice by appending", func(t *testing.T) {
		value := F.Pipe3(
			[]int{1, 2, 3},
			MakeIORef[[]int],
			io.Chain(Modify(func(s []int) []int { return append(s, 4, 5) })),
			io.Chain(Read[[]int]),
		)()
		assert.Equal(t, []int{1, 2, 3, 4, 5}, value)
	})
}

func TestModifyWithResult(t *testing.T) {
	t.Run("modifies and returns result", func(t *testing.T) {
		oldValue := F.Pipe2(
			42,
			MakeIORef[int],
			io.Chain(ModifyWithResult(func(x int) pair.Pair[int, int] {
				return pair.MakePair(x+1, x)
			})),
		)()

		assert.Equal(t, 42, oldValue)
	})

	t.Run("swaps value and returns old", func(t *testing.T) {
		old := F.Pipe2(
			100,
			MakeIORef[int],
			io.Chain(ModifyWithResult(func(x int) pair.Pair[int, int] {
				return pair.MakePair(200, x)
			})),
		)()

		assert.Equal(t, 100, old)
	})

	t.Run("computes result from old value", func(t *testing.T) {
		doubled := F.Pipe2(
			10,
			MakeIORef[int],
			io.Chain(ModifyWithResult(func(x int) pair.Pair[int, int] {
				return pair.MakePair(x+5, x*2)
			})),
		)()

		assert.Equal(t, 20, doubled) // old value * 2
	})

	t.Run("returns different type", func(t *testing.T) {
		message := F.Pipe2(
			42,
			MakeIORef[int],
			io.Chain(ModifyWithResult(func(x int) pair.Pair[int, string] {
				return pair.MakePair(x*2, "doubled")
			})),
		)()

		assert.Equal(t, "doubled", message)
	})

	t.Run("chains multiple ModifyWithResult calls", func(t *testing.T) {
		ref := MakeIORef(5)()

		result1 := ModifyWithResult(func(x int) pair.Pair[int, int] {
			return pair.MakePair(x+10, x)
		})(ref)()
		assert.Equal(t, 5, result1)

		result2 := ModifyWithResult(func(x int) pair.Pair[int, int] {
			return pair.MakePair(x*2, x)
		})(ref)()
		assert.Equal(t, 15, result2)

		finalValue := Read(ref)()
		assert.Equal(t, 30, finalValue)
	})
}

func TestConcurrency(t *testing.T) {
	t.Run("concurrent reads are safe", func(t *testing.T) {
		ref := MakeIORef(42)()
		var wg sync.WaitGroup

		for range 100 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				value := Read(ref)()
				assert.Equal(t, 42, value)
			}()
		}

		wg.Wait()
	})

	t.Run("concurrent writes are safe", func(t *testing.T) {
		ref := MakeIORef(0)()
		var wg sync.WaitGroup

		for i := range 100 {
			wg.Add(1)
			go func(val int) {
				defer wg.Done()
				Write(val)(ref)()
			}(i)
		}

		wg.Wait()
		// Final value should be one of the written values
		value := Read(ref)()
		assert.GreaterOrEqual(t, value, 0)
		assert.Less(t, value, 100)
	})

	t.Run("concurrent modifications are safe", func(t *testing.T) {
		ref := MakeIORef(0)()
		var wg sync.WaitGroup

		for range 100 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				Modify(N.Add(1))(ref)()
			}()
		}

		wg.Wait()
		value := Read(ref)()
		assert.Equal(t, 100, value)
	})

	t.Run("concurrent ModifyWithResult is safe", func(t *testing.T) {
		ref := MakeIORef(0)()
		var wg sync.WaitGroup
		results := make([]int, 100)

		for i := range 100 {
			wg.Add(1)
			go func(idx int) {
				defer wg.Done()
				old := ModifyWithResult(func(x int) pair.Pair[int, int] {
					return pair.MakePair(x+1, x)
				})(ref)()
				results[idx] = old
			}(i)
		}

		wg.Wait()

		// Final value should be 100
		assert.Equal(t, 100, Read(ref)())

		// All returned old values should be unique and in range [0, 99]
		seen := make(map[int]bool)
		for _, v := range results {
			assert.GreaterOrEqual(t, v, 0)
			assert.Less(t, v, 100)
			assert.False(t, seen[v], "duplicate old value: %d", v)
			seen[v] = true
		}
	})

	t.Run("mixed concurrent operations are safe", func(t *testing.T) {
		ref := MakeIORef(0)()
		var wg sync.WaitGroup

		// Concurrent reads
		for range 50 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				Read(ref)()
			}()
		}

		// Concurrent writes
		for i := range 25 {
			wg.Add(1)
			go func(val int) {
				defer wg.Done()
				Write(val)(ref)()
			}(i)
		}

		// Concurrent modifications
		for range 25 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				Modify(N.Add(1))(ref)()
			}()
		}

		wg.Wait()
		// Should complete without deadlock or race conditions
		value := Read(ref)()
		assert.NotNil(t, value)
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("IORef with nil pointer", func(t *testing.T) {
		var ptr *int
		value := F.Pipe2(
			ptr,
			MakeIORef[*int],
			io.Chain(Read[*int]),
		)()
		assert.Nil(t, value)

		newPtr := new(int)
		*newPtr = 42
		updatedValue := F.Pipe3(
			ptr,
			MakeIORef[*int],
			io.Chain(Write(newPtr)),
			io.Chain(Read[*int]),
		)()
		assert.NotNil(t, updatedValue)
		assert.Equal(t, 42, *updatedValue)
	})

	t.Run("IORef with empty slice", func(t *testing.T) {
		value := F.Pipe2(
			[]int{},
			MakeIORef[[]int],
			io.Chain(Read[[]int]),
		)()
		assert.Empty(t, value)

		updatedValue := F.Pipe3(
			[]int{},
			MakeIORef[[]int],
			io.Chain(Modify(func(s []int) []int { return append(s, 1) })),
			io.Chain(Read[[]int]),
		)()
		assert.Equal(t, []int{1}, updatedValue)
	})

	t.Run("IORef with empty string", func(t *testing.T) {
		value := F.Pipe2(
			"",
			MakeIORef[string],
			io.Chain(Read[string]),
		)()
		assert.Equal(t, "", value)

		updatedValue := F.Pipe3(
			"",
			MakeIORef[string],
			io.Chain(Write("not empty")),
			io.Chain(Read[string]),
		)()
		assert.Equal(t, "not empty", updatedValue)
	})

	t.Run("identity modification", func(t *testing.T) {
		value := F.Pipe3(
			42,
			MakeIORef[int],
			io.Chain(Modify(F.Identity[int])),
			io.Chain(Read[int]),
		)()
		assert.Equal(t, 42, value)
	})

	t.Run("ModifyWithResult with identity", func(t *testing.T) {
		ref := MakeIORef(42)()
		result := ModifyWithResult(pair.Of[int])(ref)()

		assert.Equal(t, 42, result)
		assert.Equal(t, 42, Read(ref)())
	})
}

func TestComplexTypes(t *testing.T) {
	t.Run("IORef with map", func(t *testing.T) {
		m := map[string]int{"a": 1, "b": 2}
		value := F.Pipe3(
			m,
			MakeIORef[map[string]int],
			io.Chain(Modify(func(m map[string]int) map[string]int {
				m["c"] = 3
				return m
			})),
			io.Chain(Read[map[string]int]),
		)()

		assert.Equal(t, 3, len(value))
		assert.Equal(t, 3, value["c"])
	})

	t.Run("IORef with channel", func(t *testing.T) {
		ch := make(chan int, 1)
		retrievedCh := F.Pipe2(
			ch,
			MakeIORef[chan int],
			io.Chain(Read[chan int]),
		)()

		retrievedCh <- 42
		value := <-retrievedCh
		assert.Equal(t, 42, value)
	})

	t.Run("IORef with function", func(t *testing.T) {
		fn := N.Mul(2)
		retrievedFn := F.Pipe2(
			fn,
			MakeIORef[func(int) int],
			io.Chain(Read[func(int) int]),
		)()

		result := retrievedFn(21)
		assert.Equal(t, 42, result)
	})
}

// Benchmark tests
func BenchmarkMakeIORef(b *testing.B) {
	for b.Loop() {
		MakeIORef(42)()
	}
}

func BenchmarkRead(b *testing.B) {
	ref := MakeIORef(42)()

	for b.Loop() {
		Read(ref)()
	}
}

func BenchmarkWrite(b *testing.B) {
	ref := MakeIORef(0)()

	for i := 0; b.Loop(); i++ {
		Write(i)(ref)()
	}
}

func BenchmarkModify(b *testing.B) {
	ref := MakeIORef(0)()

	for b.Loop() {
		Modify(N.Add(1))(ref)()
	}
}

func BenchmarkModifyWithResult(b *testing.B) {
	ref := MakeIORef(0)()

	for b.Loop() {
		ModifyWithResult(func(x int) pair.Pair[int, int] {
			return pair.MakePair(x+1, x)
		})(ref)()
	}
}

func BenchmarkConcurrentReads(b *testing.B) {
	ref := MakeIORef(42)()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Read(ref)()
		}
	})
}

func BenchmarkConcurrentWrites(b *testing.B) {
	ref := MakeIORef(0)()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Write(42)(ref)()
		}
	})
}

func BenchmarkConcurrentModify(b *testing.B) {
	ref := MakeIORef(0)()
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			Modify(N.Add(1))(ref)()
		}
	})
}

//
