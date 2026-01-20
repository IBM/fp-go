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

// package readerioresult provides a specialized version of [readerioeither.ReaderIOResult] that uses
// [context.Context] as the context type and [error] as the left (error) type. This package is designed
// for typical Go applications where context-aware, effectful computations with error handling are needed.
//
// # Core Concept
//
// ReaderIOResult[A] represents a computation that:
//   - Depends on a [context.Context] (Reader aspect)
//   - Performs side effects (IO aspect)
//   - Can fail with an [error] (Either aspect)
//   - Produces a value of type A on success
//
// The type is defined as:
//
//	ReaderIOResult[A] = func(context.Context) func() Either[error, A]
//
// This combines three powerful functional programming concepts:
//   - Reader: Dependency injection via context
//   - IO: Deferred side effects
//   - Either: Explicit error handling
//
// # Key Features
//
//   - Context-aware operations with automatic cancellation support
//   - Parallel and sequential execution modes for applicative operations
//   - Resource management with automatic cleanup (RAII pattern)
//   - Thread-safe operations with lock management
//   - Comprehensive error handling combinators
//   - Conversion utilities for standard Go error-returning functions
//
// # Core Operations
//
// Construction:
//   - [Of], [Right]: Create successful computations
//   - [Left]: Create failed computations
//   - [FromEither], [FromIO], [FromIOEither]: Convert from other types
//   - [TryCatch]: Wrap error-returning functions
//   - [Eitherize0-10]: Convert standard Go functions to ReaderIOResult
//
// Transformation:
//   - [Map]: Transform success values
//   - [MapLeft]: Transform error values
//   - [BiMap]: Transform both success and error values
//   - [Chain]: Sequence dependent computations
//   - [ChainFirst]: Sequence computations, keeping first result
//
// Combination:
//   - [Ap], [ApSeq], [ApPar]: Apply functions to values (sequential/parallel)
//   - [SequenceT2-10]: Combine multiple computations into tuples
//   - [TraverseArray], [TraverseRecord]: Transform collections
//
// Error Handling:
//   - [Fold]: Handle both success and error cases
//   - [GetOrElse]: Provide default value on error
//   - [OrElse]: Try alternative computation on error
//   - [Alt]: Alternative computation with lazy evaluation
//
// Context Operations:
//   - [Ask]: Access the context
//   - [Asks]: Access and transform the context
//   - [WithContext]: Add context cancellation checks
//   - [Never]: Create a computation that waits for context cancellation
//
// Resource Management:
//   - [Bracket]: Ensure resource cleanup (acquire/use/release pattern)
//   - [WithResource]: Manage resource lifecycle with automatic cleanup
//   - [WithLock]: Execute operations within lock scope
//
// Timing:
//   - [Delay]: Delay execution by duration
//   - [Timer]: Return current time after delay
//
// # Usage Example
//
//	import (
//	    "context"
//	    "fmt"
//	    RIOE "github.com/IBM/fp-go/v2/context/readerioresult"
//	    F "github.com/IBM/fp-go/v2/function"
//	)
//
//	// Define a computation that reads from context and may fail
//	func fetchUser(id string) RIOE.ReaderIOResult[User] {
//	    return F.Pipe2(
//	        RIOE.Ask(),
//	        RIOE.Chain(func(ctx context.Context) RIOE.ReaderIOResult[User] {
//	            return RIOE.TryCatch(func(ctx context.Context) func() (User, error) {
//	                return func() (User, error) {
//	                    return userService.Get(ctx, id)
//	                }
//	            })
//	        }),
//	        RIOE.Map(func(user User) User {
//	            // Transform the user
//	            return user
//	        }),
//	    )
//	}
//
//	// Execute the computation
//	ctx := t.Context()
//	result := fetchUser("123")(ctx)()
//	// result is Either[error, User]
//
// # Parallel vs Sequential Execution
//
// The package supports both parallel and sequential execution for applicative operations:
//
//	// Sequential execution (default for Ap)
//	result := RIOE.ApSeq(value)(function)
//
//	// Parallel execution with automatic cancellation
//	result := RIOE.ApPar(value)(function)
//
// When using parallel execution, if any operation fails, all other operations are automatically
// cancelled via context cancellation.
//
// # Resource Management Example
//
//	// Automatic resource cleanup with WithResource
//	result := F.Pipe1(
//	    RIOE.WithResource(
//	        openFile("data.txt"),
//	        closeFile,
//	    ),
//	    func(use func(func(*os.File) RIOE.ReaderIOResult[string]) RIOE.ReaderIOResult[string]) RIOE.ReaderIOResult[string] {
//	        return use(func(file *os.File) RIOE.ReaderIOResult[string] {
//	            return readContent(file)
//	        })
//	    },
//	)
//
// # Bind Operations
//
// The package provides do-notation style operations for building complex computations:
//
//	result := F.Pipe3(
//	    RIOE.Do(State{}),
//	    RIOE.Bind(setState, getFirstValue),
//	    RIOE.Bind(setState2, getSecondValue),
//	    RIOE.Map(finalTransform),
//	)
//
// # Context Cancellation
//
// All operations respect context cancellation. When a context is cancelled, operations
// will return an error containing the cancellation cause:
//
//	ctx, cancel := context.WithCancelCause(t.Context())
//	cancel(errors.New("operation cancelled"))
//	result := computation(ctx)() // Returns Left with cancellation error
//
//go:generate go run ../.. contextreaderioeither --count 10 --filename gen.go
package readerioresult
