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

// Package io provides the IO monad, representing synchronous computations that cannot fail.
//
// IO is a lazy computation that encapsulates side effects and ensures referential transparency.
// Unlike functions that execute immediately, IO values describe computations that will be
// executed when explicitly invoked.
//
// # Fantasy Land Specification
//
// This implementation corresponds to the Fantasy Land IO type:
// https://github.com/fantasyland/fantasy-land
//
// Implemented Fantasy Land algebras:
//   - Functor: https://github.com/fantasyland/fantasy-land#functor
//   - Apply: https://github.com/fantasyland/fantasy-land#apply
//   - Applicative: https://github.com/fantasyland/fantasy-land#applicative
//   - Chain: https://github.com/fantasyland/fantasy-land#chain
//   - Monad: https://github.com/fantasyland/fantasy-land#monad
//
// # Core Concepts
//
// The IO type is defined as a function that takes no arguments and returns a value:
//
//	type IO[A any] = func() A
//
// This simple definition provides powerful guarantees:
//   - Lazy evaluation: computations are not executed until explicitly called
//   - Composability: IO operations can be combined without executing them
//   - Referential transparency: IO values can be safely reused and passed around
//
// # Basic Usage
//
//	// Creating IO values
//	greeting := io.Of("Hello, World!")
//	timestamp := io.Now
//
//	// Transforming values with Map
//	upper := io.Map(strings.ToUpper)(greeting)
//
//	// Chaining computations with Chain
//	result := io.Chain(func(s string) io.IO[int] {
//	    return io.Of(len(s))
//	})(greeting)
//
//	// Executing the computation
//	value := result() // Only now does the computation run
//
// # Monadic Operations
//
// IO implements the Monad interface, providing:
//   - Of: Wrap a pure value in IO
//   - Map: Transform the result of a computation
//   - Chain (FlatMap): Sequence computations that return IO
//   - Ap: Apply a function wrapped in IO to a value wrapped in IO
//
// # Parallel vs Sequential Execution
//
// IO supports both parallel and sequential execution of applicative operations:
//   - Ap/MonadAp: Parallel execution (default)
//   - ApSeq/MonadApSeq: Sequential execution
//   - ApPar/MonadApPar: Explicit parallel execution
//
// # Time-based Operations
//
//	// Delay execution
//	delayed := io.Delay(time.Second)(computation)
//
//	// Execute after a specific time
//	scheduled := io.After(timestamp)(computation)
//
//	// Measure execution time
//	withDuration := io.WithDuration(computation)
//	withTime := io.WithTime(computation)
//
// # Resource Management
//
// IO provides utilities for safe resource management:
//
//	// Bracket ensures cleanup
//	result := io.Bracket(
//	    acquire,
//	    use,
//	    release,
//	)
//
//	// WithResource simplifies resource patterns
//	withFile := io.WithResource(openFile, closeFile)
//	result := withFile(func(f *os.File) io.IO[Data] {
//	    return readData(f)
//	})
//
// # Retry Logic
//
//	// Retry with exponential backoff
//	result := io.Retrying(
//	    retry.ExponentialBackoff(time.Second, 5),
//	    func(status retry.RetryStatus) io.IO[Result] {
//	        return fetchData()
//	    },
//	    func(r Result) bool { return r.ShouldRetry },
//	)
//
// # Traversal Operations
//
// IO provides utilities for working with collections:
//   - TraverseArray: Apply IO-returning function to array elements
//   - TraverseRecord: Apply IO-returning function to map values
//   - SequenceArray: Convert []IO[A] to IO[[]A]
//   - SequenceRecord: Convert map[K]IO[A] to IO[map[K]A]
//
// Both parallel and sequential variants are available (e.g., TraverseArraySeq).
//
// # Do Notation
//
// IO supports do-notation style composition for imperative-looking code:
//
//	result := pipe.Pipe3(
//	    io.Do(State{}),
//	    io.Bind("user", func(s State) io.IO[User] {
//	        return fetchUser(s.userId)
//	    }),
//	    io.Bind("posts", func(s State) io.IO[[]Post] {
//	        return fetchPosts(s.user.Id)
//	    }),
//	    io.Map(func(s State) Result {
//	        return formatResult(s.user, s.posts)
//	    }),
//	)
//
// # Logging and Debugging
//
//	// Log values during computation
//	logged := io.ChainFirst(io.Logger()("Fetched user"))(fetchUser)
//
//	// Printf-style logging
//	logged := io.ChainFirst(io.Printf("User: %+v"))(fetchUser)
//
// # Subpackages
//
//   - io/file: File system operations returning IO
//   - io/generic: Generic IO utilities and type classes
//   - io/testing: Testing utilities for IO laws
//
// # Relationship to Other Monads
//
// IO is the simplest effect monad in the fp-go library:
//   - IOEither: IO that can fail (combines IO with Either)
//   - IOOption: IO that may not return a value (combines IO with Option)
//   - ReaderIO: IO with dependency injection (combines Reader with IO)
//   - ReaderIOEither: Full effect system with DI and error handling
package io

//go:generate go run .. io --count 10 --filename gen.go
