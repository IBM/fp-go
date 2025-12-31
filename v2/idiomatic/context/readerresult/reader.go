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

package readerresult

import (
	"context"
	"sync"
	"time"

	RS "github.com/IBM/fp-go/v2/context/readerresult"
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/idiomatic/option"
	RR "github.com/IBM/fp-go/v2/idiomatic/readerresult"
	"github.com/IBM/fp-go/v2/idiomatic/result"
	"github.com/IBM/fp-go/v2/reader"
	RES "github.com/IBM/fp-go/v2/result"
)

// FromEither lifts a Result (Either[error, A]) into a ReaderResult.
//
// The resulting ReaderResult ignores the context.Context environment and simply
// returns the Result value. This is useful for converting existing Result values
// into the ReaderResult monad for composition with other ReaderResult operations.
//
// Type Parameters:
//   - A: The success value type
//
// Parameters:
//   - e: A Result[A] (Either[error, A]) to lift
//
// Returns:
//   - A ReaderResult[A] that ignores the context and returns the Result
//
//go:inline
func FromEither[A any](e Result[A]) ReaderResult[A] {
	return RR.FromEither[context.Context](e)
}

// FromResult creates a ReaderResult from a Go-style (value, error) tuple.
//
// This is a convenience function for converting standard Go error handling
// into the ReaderResult monad. The resulting ReaderResult ignores the context.
//
// Type Parameters:
//   - A: The value type
//
// Parameters:
//   - a: The value
//   - err: The error (nil for success)
//
// Returns:
//   - A ReaderResult[A] that returns the given value and error
//
//go:inline
func FromResult[A any](a A, err error) ReaderResult[A] {
	return RR.FromResult[context.Context](a, err)
}

//go:inline
func RightReader[A any](rdr Reader[context.Context, A]) ReaderResult[A] {
	return RR.RightReader(rdr)
}

//go:inline
func LeftReader[A, R any](l Reader[context.Context, error]) ReaderResult[A] {
	return RR.LeftReader[A](l)
}

// Left creates a ReaderResult that always fails with the given error.
//
// This is the error constructor for ReaderResult, analogous to Either's Left.
// The resulting computation ignores the context and immediately returns the error.
//
// Type Parameters:
//   - A: The success type (for type inference)
//
// Parameters:
//   - err: The error to return
//
// Returns:
//   - A ReaderResult[A] that always fails with the given error
//
//go:inline
func Left[A any](err error) ReaderResult[A] {
	return RR.Left[context.Context, A](err)
}

// Right creates a ReaderResult that always succeeds with the given value.
//
// This is the success constructor for ReaderResult, analogous to Either's Right.
// The resulting computation ignores the context and immediately returns the value.
//
// Type Parameters:
//   - A: The value type
//
// Parameters:
//   - a: The value to return
//
// Returns:
//   - A ReaderResult[A] that always succeeds with the given value
//
//go:inline
func Right[A any](a A) ReaderResult[A] {
	return RR.Right[context.Context](a)
}

// FromReader lifts a Reader into a ReaderResult that always succeeds.
//
// The Reader computation is executed and its result is wrapped in a successful Result.
// This is useful for incorporating Reader computations into ReaderResult pipelines.
//
// Type Parameters:
//   - A: The value type
//
// Parameters:
//   - r: A Reader[context.Context, A] to lift
//
// Returns:
//   - A ReaderResult[A] that executes the Reader and always succeeds
//
//go:inline
func FromReader[A any](r Reader[context.Context, A]) ReaderResult[A] {
	return RR.FromReader(r)
}

//go:inline
func FromReaderResult[A any](r RS.ReaderResult[A]) ReaderResult[A] {
	return func(ctx context.Context) (A, error) {
		return either.Unwrap(r(ctx))
	}
}

//go:inline
func ToReaderResult[A any](r ReaderResult[A]) RS.ReaderResult[A] {
	return func(ctx context.Context) Result[A] {
		return either.TryCatchError(r(ctx))
	}
}

// MonadMap transforms the success value of a ReaderResult using the given function.
//
// If the ReaderResult fails, the error is propagated unchanged. This is the
// Functor's map operation for ReaderResult.
//
// Type Parameters:
//   - A: The input value type
//   - B: The output value type
//
// Parameters:
//   - fa: The ReaderResult to transform
//   - f: The transformation function
//
// Returns:
//   - A ReaderResult[B] with the transformed value
//
//go:inline
func MonadMap[A, B any](fa ReaderResult[A], f func(A) B) ReaderResult[B] {
	return RR.MonadMap(fa, f)
}

// Map is the curried version of MonadMap, useful for function composition.
//
// It returns an Operator that can be used in pipelines with F.Pipe.
//
// Type Parameters:
//   - A: The input value type
//   - B: The output value type
//
// Parameters:
//   - f: The transformation function
//
// Returns:
//   - An Operator that transforms ReaderResult[A] to ReaderResult[B]
//
//go:inline
func Map[A, B any](f func(A) B) Operator[A, B] {
	return RR.Map[context.Context](f)
}

// MonadChain sequences two ReaderResult computations where the second depends on the first.
//
// This is the monadic bind operation (flatMap). If the first computation fails,
// the error is propagated and the second computation is not executed. Both
// computations share the same context.Context environment.
//
// Type Parameters:
//   - A: The input value type
//   - B: The output value type
//
// Parameters:
//   - ma: The first ReaderResult computation
//   - f: A Kleisli arrow that produces the second computation based on the first's result
//
// Returns:
//   - A ReaderResult[B] representing the sequenced computation
//
//go:inline
func MonadChain[A, B any](ma ReaderResult[A], f Kleisli[A, B]) ReaderResult[B] {
	return RR.MonadChain(ma, WithContextK(f))
}

// Chain is the curried version of MonadChain, useful for function composition.
//
// It returns an Operator that can be used in pipelines with F.Pipe.
//
// Type Parameters:
//   - A: The input value type
//   - B: The output value type
//
// Parameters:
//   - f: A Kleisli arrow for the second computation
//
// Returns:
//   - An Operator that chains ReaderResult computations
//
//go:inline
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B] {
	return RR.Chain(WithContextK(f))
}

// Of creates a ReaderResult that always succeeds with the given value.
//
// This is an alias for Right and represents the Applicative's pure/return operation.
// The resulting computation ignores the context and immediately returns the value.
//
// Type Parameters:
//   - A: The value type
//
// Parameters:
//   - a: The value to wrap
//
// Returns:
//   - A ReaderResult[A] that always succeeds with the given value
//
//go:inline
func Of[A any](a A) ReaderResult[A] {
	return RR.Of[context.Context](a)
}

// MonadAp applies a function wrapped in a ReaderResult to a value wrapped in a ReaderResult.
//
// This is the Applicative's ap operation. Both computations are executed concurrently
// using goroutines, and the context is shared between them. If either computation fails,
// the entire operation fails. If the context is cancelled, the operation is aborted.
//
// The concurrent execution allows for parallel independent computations, which can
// improve performance when both operations involve I/O or other blocking operations.
//
// Type Parameters:
//   - B: The result type after applying the function
//   - A: The input type to the function
//
// Parameters:
//   - fab: A ReaderResult containing a function from A to B
//   - fa: A ReaderResult containing a value of type A
//
// Returns:
//   - A ReaderResult[B] that applies the function to the value
//
// Example:
//
//	// Create a function wrapped in ReaderResult
//	addTen := readerresult.Right(func(n int) int {
//	    return n + 10
//	})
//
//	// Create a value wrapped in ReaderResult
//	value := readerresult.Right(32)
//
//	// Apply the function to the value
//	result := readerresult.MonadAp(addTen, value)
//	output, err := result(ctx)  // Returns (42, nil)
//
// Error Handling:
//
//	// If the function fails
//	failedFn := readerresult.Left[func(int) int](errors.New("function error"))
//	result := readerresult.MonadAp(failedFn, value)
//	_, err := result(ctx)  // Returns function error
//
//	// If the value fails
//	failedValue := readerresult.Left[int](errors.New("value error"))
//	result := readerresult.MonadAp(addTen, failedValue)
//	_, err := result(ctx)  // Returns value error
//
// Context Cancellation:
//
//	ctx, cancel := context.WithCancel(context.Background())
//	cancel()  // Cancel immediately
//	result := readerresult.MonadAp(addTen, value)
//	_, err := result(ctx)  // Returns context cancellation error
func MonadAp[B, A any](fab ReaderResult[func(A) B], fa ReaderResult[A]) ReaderResult[B] {
	return func(ctx context.Context) (B, error) {

		if ctx.Err() != nil {
			return result.Left[B](context.Cause(ctx))
		}

		var wg sync.WaitGroup
		wg.Add(1)

		cancelCtx, cancelFct := context.WithCancel(ctx)
		defer cancelFct()

		var a A
		var aerr error

		go func() {
			defer wg.Done()
			a, aerr = fa(cancelCtx)
			if aerr != nil {
				cancelFct()
			}
		}()

		ab, aberr := fab(cancelCtx)
		if aberr != nil {
			cancelFct()
			wg.Wait()
			return result.Left[B](aberr)
		}

		wg.Wait()

		if aerr != nil {
			return result.Left[B](aerr)
		}

		return result.Of(ab(a))
	}
}

// Ap is the curried version of MonadAp, useful for function composition.
//
// It fixes the value argument and returns an Operator that can be applied
// to a ReaderResult containing a function. This is particularly useful in
// pipelines where you want to apply a fixed value to various functions.
//
// Type Parameters:
//   - B: The result type after applying the function
//   - A: The input type to the function
//
// Parameters:
//   - fa: A ReaderResult containing a value of type A
//
// Returns:
//   - An Operator that applies the value to a function wrapped in ReaderResult
//
// Example:
//
//	import F "github.com/IBM/fp-go/v2/function"
//
//	value := readerresult.Right(32)
//	addTen := readerresult.Right(N.Add(10))
//
//	result := F.Pipe1(
//	    addTen,
//	    readerresult.Ap[int](value),
//	)
//	output, err := result(ctx)  // Returns (42, nil)
//
//go:inline
func Ap[B, A any](fa ReaderResult[A]) Operator[func(A) B, B] {
	return function.Bind2nd(MonadAp[B, A], fa)
}

//go:inline
func FromPredicate[A any](pred func(A) bool, onFalse func(A) error) Kleisli[A, A] {
	return WithContextK(RR.FromPredicate[context.Context](pred, onFalse))
}

//go:inline
func Fold[A, B any](onLeft reader.Kleisli[context.Context, error, B], onRight reader.Kleisli[context.Context, A, B]) func(ReaderResult[A]) Reader[context.Context, B] {
	return RR.Fold(onLeft, onRight)
}

//go:inline
func GetOrElse[A any](onLeft reader.Kleisli[context.Context, error, A]) func(ReaderResult[A]) Reader[context.Context, A] {
	return RR.GetOrElse(onLeft)
}

// OrElse recovers from a Left (error) by providing an alternative computation.
// If the ReaderResult is Right, it returns the value unchanged.
// If the ReaderResult is Left, it applies the provided function to the error value,
// which returns a new ReaderResult that replaces the original.
//
// This is the idiomatic version that works with context.Context-based ReaderResult.
// This is useful for error recovery, fallback logic, or chaining alternative computations.
//
// Example:
//
//	// Recover from specific errors with fallback values
//	recover := readerresult.OrElse(func(err error) readerresult.ReaderResult[int] {
//	    if err.Error() == "not found" {
//	        return readerresult.Of[int](0) // default value
//	    }
//	    return readerresult.Left[int](err) // propagate other errors
//	})
//	result := recover(readerresult.Left[int](errors.New("not found")))(ctx) // Right(0)
//	result := recover(readerresult.Of(42))(ctx) // Right(42) - unchanged
//
//go:inline
func OrElse[A any](onLeft Kleisli[error, A]) Operator[A, A] {
	return RR.OrElse(WithContextK(onLeft))
}

//go:inline
func OrLeft[A any](onLeft reader.Kleisli[context.Context, error, error]) Operator[A, A] {
	return RR.OrLeft[A](onLeft)
}

// Ask retrieves the current context.Context environment.
//
// This is the Reader's ask operation, which provides access to the environment.
// It always succeeds and returns the context that was passed in.
//
// Returns:
//   - A ReaderResult[context.Context] that returns the environment
//
//go:inline
func Ask() ReaderResult[context.Context] {
	return RR.Ask[context.Context]()
}

// Asks extracts a value from the context.Context environment using a Reader function.
//
// This is useful for accessing specific parts of the environment. The Reader
// function is applied to the context, and the result is wrapped in a successful ReaderResult.
//
// Type Parameters:
//   - A: The extracted value type
//
// Parameters:
//   - r: A Reader function that extracts a value from the context
//
// Returns:
//   - A ReaderResult[A] that extracts and returns the value
//
//go:inline
func Asks[A any](r Reader[context.Context, A]) ReaderResult[A] {
	return RR.Asks(r)
}

//go:inline
func MonadChainEitherK[A, B any](ma ReaderResult[A], f RES.Kleisli[A, B]) ReaderResult[B] {
	return RR.MonadChainEitherK(ma, f)
}

//go:inline
func ChainEitherK[A, B any](f RES.Kleisli[A, B]) Operator[A, B] {
	return RR.ChainEitherK[context.Context](f)
}

//go:inline
func MonadChainReaderK[A, B any](ma ReaderResult[A], f result.Kleisli[A, B]) ReaderResult[B] {
	return RR.MonadChainReaderK(ma, f)
}

//go:inline
func ChainReaderK[A, B any](f result.Kleisli[A, B]) Operator[A, B] {
	return RR.ChainReaderK[context.Context](f)
}

//go:inline
func ChainOptionK[A, B any](onNone Lazy[error]) func(option.Kleisli[A, B]) Operator[A, B] {
	return RR.ChainOptionK[context.Context, A, B](onNone)
}

// Flatten removes one level of ReaderResult nesting.
//
// This is equivalent to Chain with the identity function. It's useful when you have
// a ReaderResult that produces another ReaderResult and want to collapse them into one.
//
// Type Parameters:
//   - A: The inner value type
//
// Parameters:
//   - mma: A nested ReaderResult[ReaderResult[A]]
//
// Returns:
//   - A flattened ReaderResult[A]
//
//go:inline
func Flatten[A any](mma ReaderResult[ReaderResult[A]]) ReaderResult[A] {
	return RR.Flatten(mma)
}

//go:inline
func MonadBiMap[A, B any](fa ReaderResult[A], f Endomorphism[error], g func(A) B) ReaderResult[B] {
	return RR.MonadBiMap(fa, f, g)
}

//go:inline
func BiMap[A, B any](f Endomorphism[error], g func(A) B) Operator[A, B] {
	return RR.BiMap[context.Context](f, g)
}

// Read executes a ReaderResult by providing it with a context.Context.
//
// This is the elimination form for ReaderResult - it "runs" the computation
// by supplying the required environment, producing a (value, error) tuple.
//
// Type Parameters:
//   - A: The result value type
//
// Parameters:
//   - ctx: The context.Context environment to provide
//
// Returns:
//   - A function that executes a ReaderResult[A] and returns (A, error)
//
//go:inline
func Read[A any](ctx context.Context) func(ReaderResult[A]) (A, error) {
	return RR.Read[A](ctx)
}

//go:inline
func MonadFlap[A, B any](fab ReaderResult[func(A) B], a A) ReaderResult[B] {
	return RR.MonadFlap(fab, a)
}

//go:inline
func Flap[B, A any](a A) Operator[func(A) B, B] {
	return RR.Flap[context.Context, B](a)
}

//go:inline
func MonadMapLeft[A any](fa ReaderResult[A], f Endomorphism[error]) ReaderResult[A] {
	return RR.MonadMapLeft(fa, f)
}

//go:inline
func MapLeft[A any](f Endomorphism[error]) Operator[A, A] {
	return RR.MapLeft[context.Context, A](f)
}

//go:inline
func MonadAlt[A any](first ReaderResult[A], second Lazy[ReaderResult[A]]) ReaderResult[A] {
	return RR.MonadAlt(first, second)
}

//go:inline
func Alt[A any](second Lazy[ReaderResult[A]]) Operator[A, A] {
	return RR.Alt(second)
}

// Local transforms the context.Context environment before passing it to a ReaderResult computation.
//
// This is the Reader's local operation, which allows you to modify the environment
// for a specific computation without affecting the outer context. The transformation
// function receives the current context and returns a new context along with a
// cancel function. The cancel function is automatically called when the computation
// completes (via defer), ensuring proper cleanup of resources.
//
// This is useful for:
//   - Adding timeouts or deadlines to specific operations
//   - Adding context values for nested computations
//   - Creating isolated context scopes
//   - Implementing context-based dependency injection
//
// Type Parameters:
//   - A: The value type of the ReaderResult
//
// Parameters:
//   - f: A function that transforms the context and returns a cancel function
//
// Returns:
//   - An Operator that runs the computation with the transformed context
//
// Example:
//
//	import F "github.com/IBM/fp-go/v2/function"
//
//	// Add a custom value to the context
//	type key int
//	const userKey key = 0
//
//	addUser := readerresult.Local[string](func(ctx context.Context) (context.Context, context.CancelFunc) {
//	    newCtx := context.WithValue(ctx, userKey, "Alice")
//	    return newCtx, func() {} // No-op cancel
//	})
//
//	getUser := readerresult.Asks(func(ctx context.Context) string {
//	    return ctx.Value(userKey).(string)
//	})
//
//	result := F.Pipe1(
//	    getUser,
//	    addUser,
//	)
//	user, err := result(context.Background())  // Returns ("Alice", nil)
//
// Timeout Example:
//
//	// Add a 5-second timeout to a specific operation
//	withTimeout := readerresult.Local[Data](func(ctx context.Context) (context.Context, context.CancelFunc) {
//	    return context.WithTimeout(ctx, 5*time.Second)
//	})
//
//	result := F.Pipe1(
//	    fetchData,
//	    withTimeout,
//	)
func Local[A any](f func(context.Context) (context.Context, context.CancelFunc)) Operator[A, A] {
	return func(rr ReaderResult[A]) ReaderResult[A] {
		return func(ctx context.Context) (A, error) {
			if ctx.Err() != nil {
				return result.Left[A](context.Cause(ctx))
			}
			otherCtx, otherCancel := f(ctx)
			defer otherCancel()
			return rr(otherCtx)
		}
	}
}

// WithTimeout adds a timeout to the context for a ReaderResult computation.
//
// This is a convenience wrapper around Local that uses context.WithTimeout.
// The computation must complete within the specified duration, or it will be
// cancelled. This is useful for ensuring operations don't run indefinitely
// and for implementing timeout-based error handling.
//
// The timeout is relative to when the ReaderResult is executed, not when
// WithTimeout is called. The cancel function is automatically called when
// the computation completes, ensuring proper cleanup.
//
// Type Parameters:
//   - A: The value type of the ReaderResult
//
// Parameters:
//   - timeout: The maximum duration for the computation
//
// Returns:
//   - An Operator that runs the computation with a timeout
//
// Example:
//
//	import (
//	    "time"
//	    F "github.com/IBM/fp-go/v2/function"
//	)
//
//	// Fetch data with a 5-second timeout
//	fetchData := readerresult.FromReader(func(ctx context.Context) Data {
//	    // Simulate slow operation
//	    select {
//	    case <-time.After(10 * time.Second):
//	        return Data{Value: "slow"}
//	    case <-ctx.Done():
//	        return Data{}
//	    }
//	})
//
//	result := F.Pipe1(
//	    fetchData,
//	    readerresult.WithTimeout[Data](5*time.Second),
//	)
//	_, err := result(context.Background())  // Returns context.DeadlineExceeded after 5s
//
// Successful Example:
//
//	quickFetch := readerresult.Right(Data{Value: "quick"})
//	result := F.Pipe1(
//	    quickFetch,
//	    readerresult.WithTimeout[Data](5*time.Second),
//	)
//	data, err := result(context.Background())  // Returns (Data{Value: "quick"}, nil)
func WithTimeout[A any](timeout time.Duration) Operator[A, A] {
	return Local[A](func(ctx context.Context) (context.Context, context.CancelFunc) {
		return context.WithTimeout(ctx, timeout)
	})
}

// WithDeadline adds an absolute deadline to the context for a ReaderResult computation.
//
// This is a convenience wrapper around Local that uses context.WithDeadline.
// The computation must complete before the specified time, or it will be
// cancelled. This is useful for coordinating operations that must finish
// by a specific time, such as request deadlines or scheduled tasks.
//
// The deadline is an absolute time, unlike WithTimeout which uses a relative
// duration. The cancel function is automatically called when the computation
// completes, ensuring proper cleanup.
//
// Type Parameters:
//   - A: The value type of the ReaderResult
//
// Parameters:
//   - deadline: The absolute time by which the computation must complete
//
// Returns:
//   - An Operator that runs the computation with a deadline
func WithDeadline[A any](deadline time.Time) Operator[A, A] {
	return Local[A](func(ctx context.Context) (context.Context, context.CancelFunc) {
		return context.WithDeadline(ctx, deadline)
	})
}
