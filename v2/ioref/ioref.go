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
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/pair"
)

// MakeIORef creates a new IORef containing the given initial value.
//
// This function returns an IO computation that, when executed, creates a new
// mutable reference initialized with the provided value. The reference is
// thread-safe and can be safely shared across goroutines.
//
// Parameters:
//   - a: The initial value to store in the IORef
//
// Returns:
//   - An IO computation that produces a new IORef[A]
//
// Example:
//
//	// Create a new IORef with initial value 42
//	refIO := ioref.MakeIORef(42)
//	ref := refIO()  // Execute the IO to get the IORef
//
//	// Create an IORef with a string
//	strRefIO := ioref.MakeIORef("hello")
//	strRef := strRefIO()
//
//go:inline
func MakeIORef[A any](a A) IO[IORef[A]] {
	return func() IORef[A] {
		return &ioRef[A]{a: a}
	}
}

// Write atomically writes a new value to an IORef.
//
// This function returns a Kleisli arrow that takes an IORef and produces an IO
// computation that writes the new value. The write operation is atomic and
// thread-safe, using a write lock to ensure exclusive access.
//
// The function returns the IORef itself, allowing for easy chaining of operations.
//
// Parameters:
//   - a: The new value to write to the IORef
//
// Returns:
//   - A Kleisli arrow from IORef[A] to IO[IORef[A]]
//
// Example:
//
//	ref := ioref.MakeIORef(42)()
//
//	// Write a new value
//	ioref.Write(100)(ref)()
//
//	// Chain multiple writes
//	pipe.Pipe2(
//	    ref,
//	    ioref.Write(200),
//	    io.Chain(ioref.Write(300)),
//	)()
//
//go:inline
func Write[A any](a A) io.Kleisli[IORef[A], IORef[A]] {
	return func(ref IORef[A]) IO[IORef[A]] {
		return func() IORef[A] {
			ref.mu.Lock()
			defer ref.mu.Unlock()

			ref.a = a
			return ref
		}
	}
}

// Read atomically reads the current value from an IORef.
//
// This function returns an IO computation that reads the value stored in the
// IORef. The read operation is thread-safe, using a read lock that allows
// multiple concurrent readers but excludes writers.
//
// Parameters:
//   - ref: The IORef to read from
//
// Returns:
//   - An IO computation that produces the current value of type A
//
// Example:
//
//	ref := ioref.MakeIORef(42)()
//
//	// Read the current value
//	value := ioref.Read(ref)()  // 42
//
//	// Use in a pipeline
//	result := pipe.Pipe2(
//	    ref,
//	    ioref.Read[int],
//	    io.Map(func(x int) int { return x * 2 }),
//	)()
//
//go:inline
func Read[A any](ref IORef[A]) IO[A] {
	return func() A {
		ref.mu.RLock()
		defer ref.mu.RUnlock()

		return ref.a
	}
}

// Modify atomically modifies the value in an IORef using the given function.
//
// This function returns a Kleisli arrow that takes an IORef and produces an IO
// computation that applies the transformation function to the current value.
// The modification is atomic and thread-safe, using a write lock to ensure
// exclusive access during the read-modify-write cycle.
//
// Parameters:
//   - f: An endomorphism (function from A to A) that transforms the current value
//
// Returns:
//   - A Kleisli arrow from IORef[A] to IO[IORef[A]]
//
// Example:
//
//	ref := ioref.MakeIORef(42)()
//
//	// Double the value
//	ioref.Modify(func(x int) int { return x * 2 })(ref)()
//
//	// Chain multiple modifications
//	pipe.Pipe2(
//	    ref,
//	    ioref.Modify(func(x int) int { return x + 10 }),
//	    io.Chain(ioref.Modify(func(x int) int { return x * 2 })),
//	)()
//
//go:inline
func Modify[A any](f Endomorphism[A]) io.Kleisli[IORef[A], IORef[A]] {
	return func(ref IORef[A]) IO[IORef[A]] {
		return func() IORef[A] {
			ref.mu.Lock()
			defer ref.mu.Unlock()

			ref.a = f(ref.a)
			return ref
		}
	}
}

// ModifyWithResult atomically modifies the value in an IORef and returns both
// the new value and an additional result computed from the old value.
//
// This function is useful when you need to both transform the stored value and
// compute some result based on the old value in a single atomic operation.
// It's similar to Haskell's atomicModifyIORef.
//
// Parameters:
//   - f: A function that takes the old value and returns a Pair of (new value, result)
//
// Returns:
//   - A Kleisli arrow from IORef[A] to IO[B] that produces the result
//
// Example:
//
//	ref := ioref.MakeIORef(42)()
//
//	// Increment and return the old value
//	oldValue := ioref.ModifyWithResult(func(x int) pair.Pair[int, int] {
//	    return pair.MakePair(x+1, x)
//	})(ref)()  // Returns 42, ref now contains 43
//
//	// Swap and return the old value
//	old := ioref.ModifyWithResult(func(x int) pair.Pair[int, int] {
//	    return pair.MakePair(100, x)
//	})(ref)()  // Returns 43, ref now contains 100
//
//go:inline
func ModifyWithResult[A, B any](f func(A) Pair[A, B]) io.Kleisli[IORef[A], B] {
	return func(ref IORef[A]) IO[B] {
		return func() B {
			ref.mu.Lock()
			defer ref.mu.Unlock()

			result := f(ref.a)
			ref.a = pair.Head(result)
			return pair.Tail(result)
		}
	}
}
