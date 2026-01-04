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

// Package ioref provides mutable references in the IO monad.
//
// # Overview
//
// IORef represents a mutable reference that can be read and written within IO computations.
// It provides thread-safe access to shared mutable state using read-write locks, making it
// safe to use across multiple goroutines.
//
// This package is inspired by Haskell's Data.IORef module and provides a functional approach
// to managing mutable state with explicit IO effects.
//
// # Core Operations
//
// The package provides four main operations:
//
//   - MakeIORef: Creates a new IORef with an initial value
//   - Read: Atomically reads the current value from an IORef
//   - Write: Atomically writes a new value to an IORef
//   - Modify: Atomically modifies the value using a transformation function
//   - ModifyWithResult: Atomically modifies the value and returns a computed result
//
// # Thread Safety
//
// All operations on IORef are thread-safe:
//
//   - Read operations use read locks, allowing multiple concurrent readers
//   - Write and Modify operations use write locks, ensuring exclusive access
//   - The underlying sync.RWMutex ensures proper synchronization
//
// # Basic Usage
//
// Creating and using an IORef:
//
//	import (
//	    "github.com/IBM/fp-go/v2/ioref"
//	)
//
//	// Create a new IORef
//	ref := ioref.MakeIORef(42)()
//
//	// Read the current value
//	value := ioref.Read(ref)()  // 42
//
//	// Write a new value
//	ioref.Write(100)(ref)()
//
//	// Read the updated value
//	newValue := ioref.Read(ref)()  // 100
//
// # Modifying Values
//
// Use Modify to transform the value in place:
//
//	ref := ioref.MakeIORef(10)()
//
//	// Double the value
//	ioref.Modify(func(x int) int { return x * 2 })(ref)()
//
//	// Chain multiple modifications
//	ioref.Modify(func(x int) int { return x + 5 })(ref)()
//	ioref.Modify(func(x int) int { return x * 3 })(ref)()
//
//	result := ioref.Read(ref)()  // (10 * 2 + 5) * 3 = 75
//
// # Atomic Modify with Result
//
// Use ModifyWithResult when you need to both transform the value and compute a result
// from the old value in a single atomic operation:
//
//	ref := ioref.MakeIORef(42)()
//
//	// Increment and return the old value
//	oldValue := ioref.ModifyWithResult(func(x int) pair.Pair[int, int] {
//	    return pair.MakePair(x+1, x)
//	})(ref)()
//
//	// oldValue is 42, ref now contains 43
//
// This is particularly useful for implementing counters, swapping values, or any operation
// where you need to know the previous state.
//
// # Concurrent Usage
//
// IORef is safe to use across multiple goroutines:
//
//	ref := ioref.MakeIORef(0)()
//
//	// Multiple goroutines can safely modify the same IORef
//	var wg sync.WaitGroup
//	for i := 0; i < 100; i++ {
//	    wg.Add(1)
//	    go func() {
//	        defer wg.Done()
//	        ioref.Modify(func(x int) int { return x + 1 })(ref)()
//	    }()
//	}
//	wg.Wait()
//
//	result := ioref.Read(ref)()  // 100
//
// # Comparison with Haskell's IORef
//
// This implementation provides the following Haskell IORef operations:
//
//   - newIORef → MakeIORef
//   - readIORef → Read
//   - writeIORef → Write
//   - modifyIORef → Modify
//   - atomicModifyIORef → ModifyWithResult
//
// The main difference is that Go's implementation uses explicit locking (sync.RWMutex)
// rather than relying on the runtime's STM (Software Transactional Memory) as Haskell does.
//
// # Performance Considerations
//
// IORef operations are highly optimized:
//
//   - Read operations are very fast (~5ns) and allow concurrent access
//   - Write and Modify operations are slightly slower (~7-8ns) due to exclusive locking
//   - ModifyWithResult is marginally slower (~9ns) due to tuple creation
//   - All operations have zero allocations in the common case
//
// For high-contention scenarios, consider:
//
//   - Using multiple IORefs to reduce lock contention
//   - Batching modifications when possible
//   - Using Read locks for read-heavy workloads
//
// # Examples
//
// Counter with atomic increment:
//
//	counter := ioref.MakeIORef(0)()
//
//	increment := func() int {
//	    return ioref.ModifyWithResult(func(x int) pair.Pair[int, int] {
//	        return pair.MakePair(x+1, x+1)
//	    })(counter)()
//	}
//
//	id1 := increment()  // 1
//	id2 := increment()  // 2
//	id3 := increment()  // 3
//
// Shared configuration:
//
//	type Config struct {
//	    MaxRetries int
//	    Timeout    time.Duration
//	}
//
//	configRef := ioref.MakeIORef(Config{
//	    MaxRetries: 3,
//	    Timeout:    5 * time.Second,
//	})()
//
//	// Update configuration
//	ioref.Modify(func(c Config) Config {
//	    c.MaxRetries = 5
//	    return c
//	})(configRef)()
//
//	// Read configuration
//	config := ioref.Read(configRef)()
//
// Stack implementation:
//
//	type Stack []int
//
//	stackRef := ioref.MakeIORef(Stack{})()
//
//	push := func(value int) {
//	    ioref.Modify(func(s Stack) Stack {
//	        return append(s, value)
//	    })(stackRef)()
//	}
//
//	pop := func() option.Option[int] {
//	    return ioref.ModifyWithResult(func(s Stack) pair.Pair[Stack, option.Option[int]] {
//	        if len(s) == 0 {
//	            return pair.MakePair(s, option.None[int]())
//	        }
//	        return pair.MakePair(s[:len(s)-1], option.Some(s[len(s)-1]))
//	    })(stackRef)()
//	}
package ioref
