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
// IORef represents a mutable reference that can be read and written within IO computations.
// It provides thread-safe access to shared mutable state using read-write locks.
//
// This is inspired by Haskell's Data.IORef module and provides a functional approach
// to managing mutable state with explicit IO effects.
//
// Example usage:
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
//	// Modify the value
//	ioref.Modify(func(x int) int { return x * 2 })(ref)()
//
//	// Read the modified value
//	newValue := ioref.Read(ref)()  // 200
package ioref

import (
	"sync"

	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/readerio"
)

type (
	// ioRef is the internal implementation of a mutable reference.
	// It uses a read-write mutex to ensure thread-safe access to the stored value.
	//
	// The mutex allows multiple concurrent readers (using RLock) but ensures
	// exclusive access for writers (using Lock), preventing race conditions
	// when reading or modifying the stored value.
	//
	// This type is not exported; users interact with it through the IORef type alias.
	ioRef[A any] struct {
		mu sync.RWMutex // Protects concurrent access to the stored value
		a  A            // The stored value
	}

	// IO represents a synchronous computation that may have side effects.
	// It's a function that takes no arguments and returns a value of type A.
	//
	// IO computations are lazy - they don't execute until explicitly invoked
	// by calling the function. This allows for composing and chaining effects
	// before execution.
	//
	// Example:
	//
	//	// Define an IO computation
	//	computation := func() int {
	//	    fmt.Println("Computing...")
	//	    return 42
	//	}
	//
	//	// Nothing happens yet - the computation is lazy
	//	result := computation()  // Now it executes and prints "Computing..."
	IO[A any] = io.IO[A]

	// ReaderIO represents a computation that requires an environment of type R
	// and produces an IO effect that yields a value of type A.
	//
	// This combines the Reader pattern (dependency injection) with IO effects,
	// allowing computations to access shared configuration or context while
	// performing side effects.
	//
	// Example:
	//
	//	type Config struct {
	//	    multiplier int
	//	}
	//
	//	// A ReaderIO that uses config to compute a value
	//	computation := func(cfg Config) io.IO[int] {
	//	    return func() int {
	//	        return 42 * cfg.multiplier
	//	    }
	//	}
	//
	//	// Execute with specific config
	//	result := computation(Config{multiplier: 2})()  // Returns 84
	ReaderIO[R, A any] = readerio.ReaderIO[R, A]

	// IORef represents a mutable reference to a value of type A.
	// Operations on IORef are thread-safe and performed within the IO monad.
	//
	// IORef provides a way to work with mutable state in a functional style,
	// where mutations are explicit and contained within IO computations.
	// This makes side effects visible in the type system and allows for
	// better reasoning about code that uses mutable state.
	//
	// All operations on IORef (Read, Write, Modify, etc.) are atomic and
	// thread-safe, making it safe to share IORefs across goroutines.
	//
	// Example:
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
	//	// Modify the value atomically
	//	ioref.Modify(func(x int) int { return x * 2 })(ref)()
	IORef[A any] = *ioRef[A]

	// Endomorphism represents a function from A to A.
	// It's commonly used with Modify to transform the value in an IORef.
	//
	// An endomorphism is a morphism (structure-preserving map) from a
	// mathematical object to itself. In programming terms, it's simply
	// a function that takes a value and returns a value of the same type.
	//
	// Example:
	//
	//	// An endomorphism that doubles an integer
	//	double := func(x int) int { return x * 2 }
	//
	//	// An endomorphism that uppercases a string
	//	upper := func(s string) string { return strings.ToUpper(s) }
	//
	//	// Use with IORef
	//	ref := ioref.MakeIORef(21)()
	//	ioref.Modify(double)(ref)()  // ref now contains 42
	Endomorphism[A any] = endomorphism.Endomorphism[A]

	// Pair represents a tuple of two values of types A and B.
	// It's used with ModifyWithResult and ModifyIOKWithResult to return both
	// a new value for the IORef (head) and a computed result (tail).
	//
	// The head of the pair contains the new value to store in the IORef,
	// while the tail contains the result to return from the operation.
	// This allows atomic operations that both update the reference and
	// compute a result based on the old value.
	//
	// Example:
	//
	//	// Create a pair where head is the new value and tail is the old value
	//	p := pair.MakePair(newValue, oldValue)
	//
	//	// Extract values
	//	newVal := pair.Head(p)  // Gets the head (new value)
	//	oldVal := pair.Tail(p)  // Gets the tail (old value)
	//
	//	// Use with ModifyWithResult to swap and return old value
	//	ref := ioref.MakeIORef(42)()
	//	oldValue := ioref.ModifyWithResult(func(x int) pair.Pair[int, int] {
	//	    return pair.MakePair(100, x)  // Store 100, return old value
	//	})(ref)()  // oldValue is 42, ref now contains 100
	Pair[A, B any] = pair.Pair[A, B]
)
