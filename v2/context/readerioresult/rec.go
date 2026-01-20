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

package readerioresult

import (
	F "github.com/IBM/fp-go/v2/function"
	RIOR "github.com/IBM/fp-go/v2/readerioresult"
)

// TailRec implements stack-safe tail recursion for the context-aware ReaderIOResult monad.
//
// This function enables recursive computations that combine four powerful concepts:
//   - Context awareness: Automatic cancellation checking via [context.Context]
//   - Environment dependency (Reader aspect): Access to configuration, context, or dependencies
//   - Side effects (IO aspect): Logging, file I/O, network calls, etc.
//   - Error handling (Either aspect): Computations that can fail with an error
//
// The function uses an iterative loop to execute the recursion, making it safe for deep
// or unbounded recursion without risking stack overflow. Additionally, it integrates
// context cancellation checking through [WithContext], ensuring that recursive computations
// can be cancelled gracefully.
//
// # How It Works
//
// TailRec takes a Kleisli arrow that returns Trampoline[A, B]:
//   - Bounce(A): Continue recursion with the new state A
//   - Land(B): Terminate recursion successfully and return the final result B
//
// The function wraps each iteration with [WithContext] to ensure context cancellation
// is checked before each recursive step. If the context is cancelled, the recursion
// terminates early with a context cancellation error.
//
// # Type Parameters
//
//   - A: The state type that changes during recursion
//   - B: The final result type when recursion terminates successfully
//
// # Parameters
//
//   - f: A Kleisli arrow (A => ReaderIOResult[Trampoline[A, B]]) that:
//   - Takes the current state A
//   - Returns a ReaderIOResult that depends on [context.Context]
//   - Can fail with error (Left in the outer Either)
//   - Produces Trampoline[A, B] to control recursion flow (Right in the outer Either)
//
// # Returns
//
// A Kleisli arrow (A => ReaderIOResult[B]) that:
//   - Takes an initial state A
//   - Returns a ReaderIOResult that requires [context.Context]
//   - Can fail with error or context cancellation
//   - Produces the final result B after recursion completes
//
// # Context Cancellation
//
// Unlike the base [readerioresult.TailRec], this version automatically integrates
// context cancellation checking:
//   - Each recursive iteration checks if the context is cancelled
//   - If cancelled, recursion terminates immediately with a cancellation error
//   - This prevents runaway recursive computations in cancelled contexts
//   - Enables responsive cancellation for long-running recursive operations
//
// # Use Cases
//
//  1. Cancellable recursive algorithms:
//     - Tree traversals that can be cancelled mid-operation
//     - Graph algorithms with timeout requirements
//     - Recursive parsers that respect cancellation
//
//  2. Long-running recursive computations:
//     - File system traversals with cancellation support
//     - Network operations with timeout handling
//     - Database operations with connection timeout awareness
//
//  3. Interactive recursive operations:
//     - User-initiated operations that can be cancelled
//     - Background tasks with cancellation support
//     - Streaming operations with graceful shutdown
//
// # Example: Cancellable Countdown
//
//	countdownStep := func(n int) readerioresult.ReaderIOResult[tailrec.Trampoline[int, string]] {
//	    return func(ctx context.Context) ioeither.IOEither[error, tailrec.Trampoline[int, string]] {
//	        return func() either.Either[error, tailrec.Trampoline[int, string]] {
//	            if n <= 0 {
//	                return either.Right[error](tailrec.Land[int]("Done!"))
//	            }
//	            // Simulate some work
//	            time.Sleep(100 * time.Millisecond)
//	            return either.Right[error](tailrec.Bounce[string](n - 1))
//	        }
//	    }
//	}
//
//	countdown := readerioresult.TailRec(countdownStep)
//
//	// With cancellation
//	ctx, cancel := context.WithTimeout(t.Context(), 500*time.Millisecond)
//	defer cancel()
//	result := countdown(10)(ctx)() // Will be cancelled after ~500ms
//
// # Example: Cancellable File Processing
//
//	type ProcessState struct {
//	    files     []string
//	    processed []string
//	}
//
//	processStep := func(state ProcessState) readerioresult.ReaderIOResult[tailrec.Trampoline[ProcessState, []string]] {
//	    return func(ctx context.Context) ioeither.IOEither[error, tailrec.Trampoline[ProcessState, []string]] {
//	        return func() either.Either[error, tailrec.Trampoline[ProcessState, []string]] {
//	            if len(state.files) == 0 {
//	                return either.Right[error](tailrec.Land[ProcessState](state.processed))
//	            }
//
//	            file := state.files[0]
//	            // Process file (this could be cancelled via context)
//	            if err := processFileWithContext(ctx, file); err != nil {
//	                return either.Left[tailrec.Trampoline[ProcessState, []string]](err)
//	            }
//
//	            return either.Right[error](tailrec.Bounce[[]string](ProcessState{
//	                files:     state.files[1:],
//	                processed: append(state.processed, file),
//	            }))
//	        }
//	    }
//	}
//
//	processFiles := readerioresult.TailRec(processStep)
//	ctx, cancel := context.WithCancel(t.Context())
//
//	// Can be cancelled at any point during processing
//	go func() {
//	    time.Sleep(2 * time.Second)
//	    cancel() // Cancel after 2 seconds
//	}()
//
//	result := processFiles(ProcessState{files: manyFiles})(ctx)()
//
// # Stack Safety
//
// The iterative implementation ensures that even deeply recursive computations
// (thousands or millions of iterations) will not cause stack overflow, while
// still respecting context cancellation:
//
//	// Safe for very large inputs with cancellation support
//	largeCountdown := readerioresult.TailRec(countdownStep)
//	ctx := t.Context()
//	result := largeCountdown(1000000)(ctx)() // Safe, no stack overflow
//
// # Performance Considerations
//
//   - Each iteration includes context cancellation checking overhead
//   - Context checking happens before each recursive step
//   - For performance-critical code, consider the cancellation checking cost
//   - The [WithContext] wrapper adds minimal overhead for cancellation safety
//
// # See Also
//
//   - [readerioresult.TailRec]: Base tail recursion without automatic context checking
//   - [WithContext]: Context cancellation wrapper used internally
//   - [Chain]: For sequencing ReaderIOResult computations
//   - [Ask]: For accessing the context
//   - [Left]/[Right]: For creating error/success values
//
//go:inline
func TailRec[A, B any](f Kleisli[A, Trampoline[A, B]]) Kleisli[A, B] {
	return RIOR.TailRec(F.Flow2(f, WithContext))
}
