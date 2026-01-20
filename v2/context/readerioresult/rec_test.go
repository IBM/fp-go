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
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/tailrec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTailRec_BasicRecursion(t *testing.T) {
	// Test basic countdown recursion
	countdownStep := func(n int) ReaderIOResult[Trampoline[int, string]] {
		return func(ctx context.Context) IOEither[Trampoline[int, string]] {
			return func() Either[Trampoline[int, string]] {
				if n <= 0 {
					return E.Right[error](tailrec.Land[int]("Done!"))
				}
				return E.Right[error](tailrec.Bounce[string](n - 1))
			}
		}
	}

	countdown := TailRec(countdownStep)
	result := countdown(5)(t.Context())()

	assert.Equal(t, E.Of[error]("Done!"), result)
}

func TestTailRec_FactorialRecursion(t *testing.T) {
	// Test factorial computation using tail recursion
	type FactorialState struct {
		n   int
		acc int
	}

	factorialStep := func(state FactorialState) ReaderIOResult[Trampoline[FactorialState, int]] {
		return func(ctx context.Context) IOEither[Trampoline[FactorialState, int]] {
			return func() Either[Trampoline[FactorialState, int]] {
				if state.n <= 1 {
					return E.Right[error](tailrec.Land[FactorialState](state.acc))
				}
				return E.Right[error](tailrec.Bounce[int](FactorialState{
					n:   state.n - 1,
					acc: state.acc * state.n,
				}))
			}
		}
	}

	factorial := TailRec(factorialStep)
	result := factorial(FactorialState{n: 5, acc: 1})(t.Context())()

	assert.Equal(t, E.Of[error](120), result) // 5! = 120
}

func TestTailRec_ErrorHandling(t *testing.T) {
	// Test that errors are properly propagated
	testErr := errors.New("computation error")

	errorStep := func(n int) ReaderIOResult[Trampoline[int, string]] {
		return func(ctx context.Context) IOEither[Trampoline[int, string]] {
			return func() Either[Trampoline[int, string]] {
				if n == 3 {
					return E.Left[Trampoline[int, string]](testErr)
				}
				if n <= 0 {
					return E.Right[error](tailrec.Land[int]("Done!"))
				}
				return E.Right[error](tailrec.Bounce[string](n - 1))
			}
		}
	}

	errorRecursion := TailRec(errorStep)
	result := errorRecursion(5)(t.Context())()

	assert.True(t, E.IsLeft(result))
	err := E.ToError(result)
	assert.Equal(t, testErr, err)
}

func TestTailRec_ContextCancellation(t *testing.T) {
	// Test that recursion gets cancelled early when context is canceled
	var iterationCount int32

	slowStep := func(n int) ReaderIOResult[Trampoline[int, string]] {
		return func(ctx context.Context) IOEither[Trampoline[int, string]] {
			return func() Either[Trampoline[int, string]] {
				atomic.AddInt32(&iterationCount, 1)

				// Simulate some work
				time.Sleep(50 * time.Millisecond)

				if n <= 0 {
					return E.Right[error](tailrec.Land[int]("Done!"))
				}
				return E.Right[error](tailrec.Bounce[string](n - 1))
			}
		}
	}

	slowRecursion := TailRec(slowStep)

	// Create a context that will be cancelled after 100ms
	ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
	defer cancel()

	start := time.Now()
	result := slowRecursion(10)(ctx)()
	elapsed := time.Since(start)

	// Should be cancelled and return an error
	assert.True(t, E.IsLeft(result))

	// Should complete quickly due to cancellation (much less than 10 * 50ms = 500ms)
	assert.Less(t, elapsed, 200*time.Millisecond)

	// Should have executed only a few iterations before cancellation
	iterations := atomic.LoadInt32(&iterationCount)
	assert.Less(t, iterations, int32(5), "Should have been cancelled before completing all iterations")
}

func TestTailRec_ImmediateCancellation(t *testing.T) {
	// Test with an already cancelled context
	countdownStep := func(n int) ReaderIOResult[Trampoline[int, string]] {
		return func(ctx context.Context) IOEither[Trampoline[int, string]] {
			return func() Either[Trampoline[int, string]] {
				if n <= 0 {
					return E.Right[error](tailrec.Land[int]("Done!"))
				}
				return E.Right[error](tailrec.Bounce[string](n - 1))
			}
		}
	}

	countdown := TailRec(countdownStep)

	// Create an already cancelled context
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	result := countdown(5)(ctx)()

	// Should immediately return a cancellation error
	assert.True(t, E.IsLeft(result))
	err := E.ToError(result)
	assert.Equal(t, context.Canceled, err)
}

func TestTailRec_StackSafety(t *testing.T) {
	// Test that deep recursion doesn't cause stack overflow
	const largeN = 10000

	countdownStep := func(n int) ReaderIOResult[Trampoline[int, int]] {
		return func(ctx context.Context) IOEither[Trampoline[int, int]] {
			return func() Either[Trampoline[int, int]] {
				if n <= 0 {
					return E.Right[error](tailrec.Land[int](0))
				}
				return E.Right[error](tailrec.Bounce[int](n - 1))
			}
		}
	}

	countdown := TailRec(countdownStep)
	result := countdown(largeN)(t.Context())()

	assert.Equal(t, E.Of[error](0), result)
}

func TestTailRec_StackSafetyWithCancellation(t *testing.T) {
	// Test stack safety with cancellation after many iterations
	const largeN = 100000
	var iterationCount int32

	countdownStep := func(n int) ReaderIOResult[Trampoline[int, int]] {
		return func(ctx context.Context) IOEither[Trampoline[int, int]] {
			return func() Either[Trampoline[int, int]] {
				atomic.AddInt32(&iterationCount, 1)

				// Add a small delay every 1000 iterations to make cancellation more likely
				if n%1000 == 0 {
					time.Sleep(1 * time.Millisecond)
				}

				if n <= 0 {
					return E.Right[error](tailrec.Land[int](0))
				}
				return E.Right[error](tailrec.Bounce[int](n - 1))
			}
		}
	}

	countdown := TailRec(countdownStep)

	// Cancel after 50ms to allow some iterations but not all
	ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
	defer cancel()

	result := countdown(largeN)(ctx)()

	// Should be cancelled (or completed if very fast)
	// The key is that it doesn't cause a stack overflow
	iterations := atomic.LoadInt32(&iterationCount)
	assert.Greater(t, iterations, int32(0))

	// If it was cancelled, verify it didn't complete all iterations
	if E.IsLeft(result) {
		assert.Less(t, iterations, int32(largeN))
	}
}

func TestTailRec_ComplexState(t *testing.T) {
	// Test with more complex state management
	type ProcessState struct {
		items     []string
		processed []string
		errors    []error
	}

	processStep := func(state ProcessState) ReaderIOResult[Trampoline[ProcessState, []string]] {
		return func(ctx context.Context) IOEither[Trampoline[ProcessState, []string]] {
			return func() Either[Trampoline[ProcessState, []string]] {
				if A.IsEmpty(state.items) {
					return E.Right[error](tailrec.Land[ProcessState](state.processed))
				}

				item := state.items[0]

				// Simulate processing that might fail for certain items
				if item == "error-item" {
					return E.Left[Trampoline[ProcessState, []string]](
						fmt.Errorf("failed to process item: %s", item))
				}

				return E.Right[error](tailrec.Bounce[[]string](ProcessState{
					items:     state.items[1:],
					processed: append(state.processed, item),
					errors:    state.errors,
				}))
			}
		}
	}

	processItems := TailRec(processStep)

	t.Run("successful processing", func(t *testing.T) {
		initialState := ProcessState{
			items:     []string{"item1", "item2", "item3"},
			processed: []string{},
			errors:    []error{},
		}

		result := processItems(initialState)(t.Context())()

		assert.Equal(t, E.Of[error]([]string{"item1", "item2", "item3"}), result)
	})

	t.Run("processing with error", func(t *testing.T) {
		initialState := ProcessState{
			items:     []string{"item1", "error-item", "item3"},
			processed: []string{},
			errors:    []error{},
		}

		result := processItems(initialState)(t.Context())()

		assert.True(t, E.IsLeft(result))
		err := E.ToError(result)
		assert.Contains(t, err.Error(), "failed to process item: error-item")
	})
}

func TestTailRec_CancellationDuringProcessing(t *testing.T) {
	// Test cancellation during a realistic processing scenario
	type FileProcessState struct {
		files     []string
		processed int
	}

	var processedCount int32

	processFileStep := func(state FileProcessState) ReaderIOResult[Trampoline[FileProcessState, int]] {
		return func(ctx context.Context) IOEither[Trampoline[FileProcessState, int]] {
			return func() Either[Trampoline[FileProcessState, int]] {
				if A.IsEmpty(state.files) {
					return E.Right[error](tailrec.Land[FileProcessState](state.processed))
				}

				// Simulate file processing time
				time.Sleep(20 * time.Millisecond)
				atomic.AddInt32(&processedCount, 1)

				return E.Right[error](tailrec.Bounce[int](FileProcessState{
					files:     state.files[1:],
					processed: state.processed + 1,
				}))
			}
		}
	}

	processFiles := TailRec(processFileStep)

	// Create many files to process
	files := make([]string, 20)
	for i := range files {
		files[i] = fmt.Sprintf("file%d.txt", i)
	}

	initialState := FileProcessState{
		files:     files,
		processed: 0,
	}

	// Cancel after 100ms (should allow ~5 files to be processed)
	ctx, cancel := context.WithTimeout(t.Context(), 100*time.Millisecond)
	defer cancel()

	start := time.Now()
	result := processFiles(initialState)(ctx)()
	elapsed := time.Since(start)

	// Should be cancelled
	assert.True(t, E.IsLeft(result))

	// Should complete quickly due to cancellation
	assert.Less(t, elapsed, 150*time.Millisecond)

	// Should have processed some but not all files
	processed := atomic.LoadInt32(&processedCount)
	assert.Greater(t, processed, int32(0))
	assert.Less(t, processed, int32(20))
}

func TestTailRec_ZeroIterations(t *testing.T) {
	// Test case where recursion terminates immediately
	immediateStep := func(n int) ReaderIOResult[Trampoline[int, string]] {
		return func(ctx context.Context) IOEither[Trampoline[int, string]] {
			return func() Either[Trampoline[int, string]] {
				return E.Right[error](tailrec.Land[int]("immediate"))
			}
		}
	}

	immediate := TailRec(immediateStep)
	result := immediate(100)(t.Context())()

	assert.Equal(t, E.Of[error]("immediate"), result)
}

func TestTailRec_ContextWithDeadline(t *testing.T) {
	// Test with context deadline
	var iterationCount int32

	slowStep := func(n int) ReaderIOResult[Trampoline[int, string]] {
		return func(ctx context.Context) IOEither[Trampoline[int, string]] {
			return func() Either[Trampoline[int, string]] {
				atomic.AddInt32(&iterationCount, 1)
				time.Sleep(30 * time.Millisecond)

				if n <= 0 {
					return E.Right[error](tailrec.Land[int]("Done!"))
				}
				return E.Right[error](tailrec.Bounce[string](n - 1))
			}
		}
	}

	slowRecursion := TailRec(slowStep)

	// Set deadline 80ms from now
	ctx, cancel := context.WithDeadline(t.Context(), time.Now().Add(80*time.Millisecond))
	defer cancel()

	result := slowRecursion(10)(ctx)()

	// Should be cancelled due to deadline
	assert.True(t, E.IsLeft(result))

	// Should have executed only a few iterations
	iterations := atomic.LoadInt32(&iterationCount)
	assert.Greater(t, iterations, int32(0))
	assert.Less(t, iterations, int32(5))
}

func TestTailRec_ContextWithValue(t *testing.T) {
	// Test that context values are preserved through recursion
	type contextKey string
	const testKey contextKey = "test"

	valueStep := func(n int) ReaderIOResult[Trampoline[int, string]] {
		return func(ctx context.Context) IOEither[Trampoline[int, string]] {
			return func() Either[Trampoline[int, string]] {
				value := ctx.Value(testKey)
				require.NotNil(t, value)
				assert.Equal(t, "test-value", value.(string))

				if n <= 0 {
					return E.Right[error](tailrec.Land[int]("Done!"))
				}
				return E.Right[error](tailrec.Bounce[string](n - 1))
			}
		}
	}

	valueRecursion := TailRec(valueStep)
	ctx := context.WithValue(t.Context(), testKey, "test-value")
	result := valueRecursion(3)(ctx)()

	assert.Equal(t, E.Of[error]("Done!"), result)
}
