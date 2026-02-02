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
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/readerio"
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

// Write atomically writes a new value to an IORef and returns the written value.
//
// This function returns a Kleisli arrow that takes an IORef and produces an IO
// computation that writes the given value to the reference. The write operation
// is atomic and thread-safe, using a write lock to ensure exclusive access.
//
// Parameters:
//   - a: The new value to write to the IORef
//
// Returns:
//   - A Kleisli arrow from IORef[A] to IO[A] that writes the value and returns it
//
// Example:
//
//	ref := ioref.MakeIORef(42)()
//
//	// Write a new value
//	newValue := ioref.Write(100)(ref)()  // Returns 100, ref now contains 100
//
//	// Chain writes
//	pipe.Pipe2(
//	    ref,
//	    ioref.Write(50),
//	    io.Chain(ioref.Write(75)),
//	)()  // ref now contains 75
//
//go:inline
func Write[A any](a A) io.Kleisli[IORef[A], A] {
	return func(ref IORef[A]) IO[A] {
		return func() A {
			ref.mu.Lock()
			defer ref.mu.Unlock()

			ref.a = a
			return a
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
func Modify[A any](f Endomorphism[A]) io.Kleisli[IORef[A], A] {
	return ModifyIOK(function.Flow2(f, io.Of))
}

// ModifyIOK atomically modifies the value in an IORef using an IO-based transformation.
//
// This is a more powerful version of Modify that allows the transformation function
// to perform IO effects. The function takes a Kleisli arrow (a function from A to IO[A])
// and returns a Kleisli arrow that modifies the IORef atomically.
//
// The modification is atomic and thread-safe, using a write lock to ensure exclusive
// access during the read-modify-write cycle. The IO effect in the transformation
// function is executed while holding the lock.
//
// Parameters:
//   - f: A Kleisli arrow (io.Kleisli[A, A]) that transforms the current value with IO effects
//
// Returns:
//   - A Kleisli arrow from IORef[A] to IO[A] that returns the new value
//
// Example:
//
//	ref := ioref.MakeIORef(42)()
//
//	// Modify with an IO effect (e.g., logging)
//	modifyWithLog := ioref.ModifyIOK(func(x int) io.IO[int] {
//	    return func() int {
//	        fmt.Printf("Old value: %d\n", x)
//	        return x * 2
//	    }
//	})
//	newValue := modifyWithLog(ref)()  // Logs and returns 84
//
//	// Chain multiple IO-based modifications
//	pipe.Pipe2(
//	    ref,
//	    ioref.ModifyIOK(func(x int) io.IO[int] {
//	        return io.Of(x + 10)
//	    }),
//	    io.Chain(ioref.ModifyIOK(func(x int) io.IO[int] {
//	        return io.Of(x * 2)
//	    })),
//	)()
func ModifyIOK[A any](f io.Kleisli[A, A]) io.Kleisli[IORef[A], A] {
	return func(ref IORef[A]) IO[A] {
		return func() A {
			ref.mu.Lock()
			defer ref.mu.Unlock()

			ref.a = f(ref.a)()
			return ref.a
		}
	}
}

// ModifyReaderIOK atomically modifies the value in an IORef using a ReaderIO-based transformation.
//
// This is a variant of ModifyIOK that works with ReaderIO computations, allowing the
// transformation function to access an environment of type R while performing IO effects.
// This is useful when the modification logic needs access to configuration, context,
// or other shared resources.
//
// The modification is atomic and thread-safe, using a write lock to ensure exclusive
// access during the read-modify-write cycle. The ReaderIO effect in the transformation
// function is executed while holding the lock.
//
// Parameters:
//   - f: A ReaderIO Kleisli arrow (readerio.Kleisli[R, A, A]) that takes the current value
//     and an environment R, and returns an IO computation producing the new value
//
// Returns:
//   - A ReaderIO Kleisli arrow from IORef[A] to ReaderIO[R, A] that returns the new value
//
// Example:
//
//	type Config struct {
//	    multiplier int
//	}
//
//	ref := ioref.MakeIORef(10)()
//
//	// Modify using environment
//	modifyWithConfig := ioref.ModifyReaderIOK(func(x int) readerio.ReaderIO[Config, int] {
//	    return func(cfg Config) io.IO[int] {
//	        return func() int {
//	            return x * cfg.multiplier
//	        }
//	    }
//	})
//
//	config := Config{multiplier: 5}
//	newValue := modifyWithConfig(ref)(config)()  // Returns 50, ref now contains 50
func ModifyReaderIOK[R, A any](f readerio.Kleisli[R, A, A]) readerio.Kleisli[R, IORef[A], A] {
	return func(ref IORef[A]) ReaderIO[R, A] {
		return func(r R) readerio.IO[A] {
			return func() A {
				ref.mu.Lock()
				defer ref.mu.Unlock()

				ref.a = f(ref.a)(r)()
				return ref.a
			}
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
	return ModifyIOKWithResult(function.Flow2(f, io.Of))
}

// ModifyIOKWithResult atomically modifies the value in an IORef and returns a result,
// using an IO-based transformation function.
//
// This is a more powerful version of ModifyWithResult that allows the transformation
// function to perform IO effects. The function takes a Kleisli arrow that transforms
// the old value into an IO computation producing a Pair of (new value, result).
//
// This is useful when you need to:
//   - Both transform the stored value and compute some result based on the old value
//   - Perform IO effects during the transformation (e.g., logging, validation)
//   - Ensure atomicity of the entire read-transform-write-compute cycle
//
// The modification is atomic and thread-safe, using a write lock to ensure exclusive
// access. The IO effect in the transformation function is executed while holding the lock.
//
// Parameters:
//   - f: A Kleisli arrow (io.Kleisli[A, Pair[A, B]]) that takes the old value and
//     returns an IO computation producing a Pair of (new value, result)
//
// Returns:
//   - A Kleisli arrow from IORef[A] to IO[B] that produces the result
//
// Example:
//
//	ref := ioref.MakeIORef(42)()
//
//	// Increment with IO effect and return old value
//	incrementWithLog := ioref.ModifyIOKWithResult(func(x int) io.IO[pair.Pair[int, int]] {
//	    return func() pair.Pair[int, int] {
//	        fmt.Printf("Incrementing from %d\n", x)
//	        return pair.MakePair(x+1, x)
//	    }
//	})
//	oldValue := incrementWithLog(ref)()  // Logs and returns 42, ref now contains 43
//
//	// Swap with validation
//	swapWithValidation := ioref.ModifyIOKWithResult(func(old int) io.IO[pair.Pair[int, string]] {
//	    return func() pair.Pair[int, string] {
//	        if old < 0 {
//	            return pair.MakePair(0, "reset negative")
//	        }
//	        return pair.MakePair(100, fmt.Sprintf("swapped %d", old))
//	    }
//	})
//	message := swapWithValidation(ref)()
func ModifyIOKWithResult[A, B any](f io.Kleisli[A, Pair[A, B]]) io.Kleisli[IORef[A], B] {
	return func(ref IORef[A]) IO[B] {
		return func() B {
			ref.mu.Lock()
			defer ref.mu.Unlock()

			result := f(ref.a)()
			ref.a = pair.Head(result)
			return pair.Tail(result)
		}
	}
}

// ModifyReaderIOKWithResult atomically modifies the value in an IORef and returns a result,
// using a ReaderIO-based transformation function.
//
// This combines the capabilities of ModifyIOKWithResult and ModifyReaderIOK, allowing the
// transformation function to:
//   - Access an environment of type R (like configuration or context)
//   - Perform IO effects during the transformation
//   - Both update the stored value and compute a result based on the old value
//   - Ensure atomicity of the entire read-transform-write-compute cycle
//
// The modification is atomic and thread-safe, using a write lock to ensure exclusive
// access. The ReaderIO effect in the transformation function is executed while holding the lock.
//
// Parameters:
//   - f: A ReaderIO Kleisli arrow (readerio.Kleisli[R, A, Pair[A, B]]) that takes the old value
//     and an environment R, and returns an IO computation producing a Pair of (new value, result)
//
// Returns:
//   - A ReaderIO Kleisli arrow from IORef[A] to ReaderIO[R, B] that produces the result
//
// Example:
//
//	type Config struct {
//	    logEnabled bool
//	}
//
//	ref := ioref.MakeIORef(42)()
//
//	// Increment with conditional logging, return old value
//	incrementWithLog := ioref.ModifyReaderIOKWithResult(
//	    func(x int) readerio.ReaderIO[Config, pair.Pair[int, int]] {
//	        return func(cfg Config) io.IO[pair.Pair[int, int]] {
//	            return func() pair.Pair[int, int] {
//	                if cfg.logEnabled {
//	                    fmt.Printf("Incrementing from %d\n", x)
//	                }
//	                return pair.MakePair(x+1, x)
//	            }
//	        }
//	    },
//	)
//
//	config := Config{logEnabled: true}
//	oldValue := incrementWithLog(ref)(config)()  // Logs and returns 42, ref now contains 43
func ModifyReaderIOKWithResult[R, A, B any](f readerio.Kleisli[R, A, Pair[A, B]]) readerio.Kleisli[R, IORef[A], B] {
	return func(ref IORef[A]) ReaderIO[R, B] {
		return func(r R) readerio.IO[B] {
			return func() B {
				ref.mu.Lock()
				defer ref.mu.Unlock()

				result := f(ref.a)(r)()
				ref.a = pair.Head(result)
				return pair.Tail(result)
			}
		}
	}
}
