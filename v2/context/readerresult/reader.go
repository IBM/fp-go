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

	"github.com/IBM/fp-go/v2/function"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/fromio"
	"github.com/IBM/fp-go/v2/internal/fromioeither"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readereither"
	"github.com/IBM/fp-go/v2/result"
)

// FromIO lifts a pure IO computation into a ReaderResult.
// The IO computation is executed when the ReaderResult is run, ignoring the context.
//
// IMPORTANT: While IO represents a side-effectful computation, combining it with ReaderResult
// makes sense because ReaderResult already has side effects due to its context.Context dependency.
// The context can be cancelled, has deadlines, and carries values - all side effects. Therefore,
// adding IO operations (which also have side effects) is a natural fit.
//
// Type Parameters:
//   - A: The success type of the IO computation
//
// Parameters:
//   - t: The IO computation to lift
//
// Returns:
//   - A ReaderResult that executes the IO computation and wraps the result in Right
//
// Example:
//
//	ioOp := func() int { return 42 }
//	rr := readerresult.FromIO(ioOp)
//	result := rr(t.Context()) // Right(42)
//
//go:inline
func FromIO[A any](t io.IO[A]) ReaderResult[A] {
	return func(_ context.Context) Result[A] {
		return result.Of(t())
	}
}

// FromIOResult lifts an IOResult computation into a ReaderResult.
// The IOResult computation is executed when the ReaderResult is run, ignoring the context.
//
// IMPORTANT: Combining IOResult with ReaderResult makes sense because both represent side-effectful
// computations. ReaderResult has side effects from context.Context (cancellation, deadlines, values),
// and IOResult has side effects from IO operations. This combination allows you to work with
// context-aware error handling while performing IO operations.
//
// Type Parameters:
//   - A: The success type of the IOResult computation
//
// Parameters:
//   - t: The IOResult computation to lift
//
// Returns:
//   - A ReaderResult that executes the IOResult computation
//
// Example:
//
//	ioResultOp := func() result.Result[int] {
//	    return result.Of(42)
//	}
//	rr := readerresult.FromIOResult(ioResultOp)
//	result := rr(t.Context()) // Right(42)
//
//go:inline
func FromIOResult[A any](t ioresult.IOResult[A]) ReaderResult[A] {
	return func(_ context.Context) Result[A] {
		return t()
	}
}

// FromReader lifts a Reader computation into a ReaderResult.
// The Reader computation receives the context and its result is wrapped in Right.
//
// Type Parameters:
//   - A: The success type of the Reader computation
//
// Parameters:
//   - r: The Reader computation to lift
//
// Returns:
//   - A ReaderResult that executes the Reader and wraps the result in Right
//
// Example:
//
//	reader := func(ctx context.Context) int {
//	    return 42
//	}
//	rr := readerresult.FromReader(reader)
//	result := rr(t.Context()) // Right(42)
//
//go:inline
func FromReader[A any](r Reader[context.Context, A]) ReaderResult[A] {
	return readereither.FromReader[error](r)
}

// FromEither lifts an Either value into a ReaderResult.
// The Either value is returned as-is, ignoring the context.
//
// Type Parameters:
//   - A: The success type of the Either value
//
// Parameters:
//   - e: The Either value to lift
//
// Returns:
//   - A ReaderResult that returns the Either value
//
// Example:
//
//	either := result.Of[error](42)
//	rr := readerresult.FromEither(either)
//	result := rr(t.Context()) // Right(42)
//
//go:inline
func FromEither[A any](e Either[A]) ReaderResult[A] {
	return readereither.FromEither[context.Context](e)
}

// Left creates a ReaderResult that always returns a Left (error) value.
//
// Type Parameters:
//   - A: The success type (not used, as this always returns an error)
//
// Parameters:
//   - l: The error value to return
//
// Returns:
//   - A ReaderResult that always returns Left(l)
//
// Example:
//
//	rr := readerresult.Left[int](errors.New("failed"))
//	result := rr(t.Context()) // Left(error("failed"))
//
//go:inline
func Left[A any](l error) ReaderResult[A] {
	return readereither.Left[context.Context, A](l)
}

// Right creates a ReaderResult that always returns a Right (success) value.
// This is an alias for Of.
//
// Type Parameters:
//   - A: The success type
//
// Parameters:
//   - r: The success value to return
//
// Returns:
//   - A ReaderResult that always returns Right(r)
//
// Example:
//
//	rr := readerresult.Right(42)
//	result := rr(t.Context()) // Right(42)
//
//go:inline
func Right[A any](r A) ReaderResult[A] {
	return readereither.Right[context.Context, error](r)
}

// MonadMap applies a function to the success value of a ReaderResult.
// This is the monadic version that takes both the ReaderResult and the function as parameters.
//
// Type Parameters:
//   - A: The input success type
//   - B: The output success type
//
// Parameters:
//   - fa: The ReaderResult to map over
//   - f: The function to apply to the success value
//
// Returns:
//   - A ReaderResult with the function applied to the success value
//
// Example:
//
//	rr := readerresult.Of(42)
//	mapped := readerresult.MonadMap(rr, func(x int) string {
//	    return fmt.Sprintf("value: %d", x)
//	})
//	result := mapped(t.Context()) // Right("value: 42")
//
//go:inline
func MonadMap[A, B any](fa ReaderResult[A], f func(A) B) ReaderResult[B] {
	return readereither.MonadMap(fa, f)
}

// Map creates an operator that applies a function to the success value of a ReaderResult.
// This is the curried version where the function is provided first.
//
// Type Parameters:
//   - A: The input success type
//   - B: The output success type
//
// Parameters:
//   - f: The function to apply to the success value
//
// Returns:
//   - An Operator that applies the function to a ReaderResult
//
// Example:
//
//	toString := readerresult.Map(func(x int) string {
//	    return fmt.Sprintf("value: %d", x)
//	})
//	rr := readerresult.Of(42)
//	result := toString(rr)(t.Context()) // Right("value: 42")
//
//go:inline
func Map[A, B any](f func(A) B) Operator[A, B] {
	return readereither.Map[context.Context, error](f)
}

// MonadChain sequences two ReaderResult computations, passing the success value from the first
// to the second. This is the monadic version that takes both the ReaderResult and the Kleisli
// function as parameters.
//
// Type Parameters:
//   - A: The success type of the first ReaderResult
//   - B: The success type of the second ReaderResult
//
// Parameters:
//   - ma: The first ReaderResult to execute
//   - f: The Kleisli function that takes the success value and returns a new ReaderResult
//
// Returns:
//   - A ReaderResult that sequences both computations
//
// Example:
//
//	rr := readerresult.Of(42)
//	chained := readerresult.MonadChain(rr, func(x int) readerresult.ReaderResult[string] {
//	    return readerresult.Of(fmt.Sprintf("value: %d", x))
//	})
//	result := chained(t.Context()) // Right("value: 42")
//
//go:inline
func MonadChain[A, B any](ma ReaderResult[A], f Kleisli[A, B]) ReaderResult[B] {
	return readereither.MonadChain(ma, F.Flow2(f, WithContext))
}

// Chain creates an operator that sequences two ReaderResult computations.
// This is the curried version where the Kleisli function is provided first.
//
// Type Parameters:
//   - A: The success type of the first ReaderResult
//   - B: The success type of the second ReaderResult
//
// Parameters:
//   - f: The Kleisli function that takes the success value and returns a new ReaderResult
//
// Returns:
//   - An Operator that sequences the computations
//
// Example:
//
//	toUpper := readerresult.Chain(func(s string) readerresult.ReaderResult[string] {
//	    return readerresult.Of(strings.ToUpper(s))
//	})
//	rr := readerresult.Of("hello")
//	result := toUpper(rr)(t.Context()) // Right("HELLO")
//
//go:inline
func Chain[A, B any](f Kleisli[A, B]) Operator[A, B] {
	return readereither.Chain(F.Flow2(f, WithContext))
}

// Of creates a ReaderResult that always returns a Right (success) value.
// This is the pointed functor constructor.
//
// Type Parameters:
//   - A: The success type
//
// Parameters:
//   - a: The success value to return
//
// Returns:
//   - A ReaderResult that always returns Right(a)
//
// Example:
//
//	rr := readerresult.Of(42)
//	result := rr(t.Context()) // Right(42)
//
//go:inline
func Of[A any](a A) ReaderResult[A] {
	return readereither.Of[context.Context, error](a)
}

// MonadAp applies a ReaderResult containing a function to a ReaderResult containing a value.
// This is the monadic version that takes both ReaderResults as parameters.
//
// Type Parameters:
//   - A: The input type
//   - B: The output type
//
// Parameters:
//   - fab: The ReaderResult containing the function
//   - fa: The ReaderResult containing the value
//
// Returns:
//   - A ReaderResult with the function applied to the value
//
// Example:
//
//	fabRR := readerresult.Of(func(x int) string {
//	    return fmt.Sprintf("value: %d", x)
//	})
//	faRR := readerresult.Of(42)
//	result := readerresult.MonadAp(fabRR, faRR)(t.Context()) // Right("value: 42")
//
//go:inline
func MonadAp[A, B any](fab ReaderResult[func(A) B], fa ReaderResult[A]) ReaderResult[B] {
	return readereither.MonadAp(fab, fa)
}

// Ap creates a function that applies a ReaderResult containing a function to a ReaderResult
// containing a value. This is the curried version where the value ReaderResult is provided first.
//
// Type Parameters:
//   - A: The input type
//   - B: The output type
//
// Parameters:
//   - fa: The ReaderResult containing the value
//
// Returns:
//   - A function that takes a ReaderResult containing a function and returns the result
//
// Example:
//
//	faRR := readerresult.Of(42)
//	applyTo42 := readerresult.Ap[int, string](faRR)
//	fabRR := readerresult.Of(func(x int) string {
//	    return fmt.Sprintf("value: %d", x)
//	})
//	result := applyTo42(fabRR)(t.Context()) // Right("value: 42")
//
//go:inline
func Ap[A, B any](fa ReaderResult[A]) func(ReaderResult[func(A) B]) ReaderResult[B] {
	return readereither.Ap[B](fa)
}

// FromPredicate creates a Kleisli arrow that validates a value using a predicate.
// If the predicate returns true, the value is wrapped in Right.
// If the predicate returns false, the onFalse function is called to generate an error.
//
// Type Parameters:
//   - A: The type of the value to validate
//
// Parameters:
//   - pred: The predicate function to test the value
//   - onFalse: The function to generate an error when the predicate fails
//
// Returns:
//   - A Kleisli arrow that validates the value
//
// Example:
//
//	isPositive := readerresult.FromPredicate(
//	    func(x int) bool { return x > 0 },
//	    func(x int) error { return fmt.Errorf("%d is not positive", x) },
//	)
//	result1 := isPositive(42)(t.Context()) // Right(42)
//	result2 := isPositive(-1)(t.Context()) // Left(error("-1 is not positive"))
//
//go:inline
func FromPredicate[A any](pred func(A) bool, onFalse func(A) error) Kleisli[A, A] {
	return readereither.FromPredicate[context.Context](pred, onFalse)
}

// OrElse recovers from a Left (error) by providing an alternative computation with access to context.Context.
// If the ReaderResult is Right, it returns the value unchanged.
// If the ReaderResult is Left, it applies the provided function to the error value,
// which returns a new ReaderResult that replaces the original.
//
// This is useful for error recovery, fallback logic, or chaining alternative computations
// that need access to the context (for cancellation, deadlines, or values).
//
// Example:
//
//	// Recover with context-aware fallback
//	recover := readerresult.OrElse(func(err error) readerresult.ReaderResult[int] {
//	    if err.Error() == "not found" {
//	        return func(ctx context.Context) result.Result[int] {
//	            // Could check ctx.Err() here
//	            return result.Of(42)
//	        }
//	    }
//	    return readerresult.Left[int](err)
//	})
//
// OrElse recovers from a Left (error) by providing an alternative computation.
// If the ReaderResult is Right, it returns the value unchanged.
// If the ReaderResult is Left, it applies the provided function to the error value,
// which returns a new ReaderResult that replaces the original.
//
// This is useful for error recovery, fallback logic, or chaining alternative computations
// in the context of Reader computations with context.Context.
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
func OrElse[A any](onLeft Kleisli[error, A]) Kleisli[ReaderResult[A], A] {
	return readereither.OrElse(F.Flow2(onLeft, WithContext))
}

// Ask returns a ReaderResult that provides access to the context.Context.
// This allows you to read the context within a ReaderResult computation.
//
// Returns:
//   - A ReaderResult that returns the context.Context as its success value
//
// Example:
//
//	rr := readerresult.Ask()
//	result := rr(t.Context()) // Right(t.Context())
//
//	// Use in a chain to access context
//	pipeline := F.Pipe2(
//	    readerresult.Ask(),
//	    readerresult.Chain(func(ctx context.Context) readerresult.ReaderResult[string] {
//	        if deadline, ok := ctx.Deadline(); ok {
//	            return readerresult.Of(fmt.Sprintf("deadline: %v", deadline))
//	        }
//	        return readerresult.Of("no deadline")
//	    }),
//	)
//
//go:inline
func Ask() ReaderResult[context.Context] {
	return readereither.Ask[context.Context, error]()
}

// MonadChainEitherK sequences a ReaderResult with a function that returns an Either.
// This is the monadic version that takes both the ReaderResult and the function as parameters.
//
// Type Parameters:
//   - A: The success type of the input ReaderResult
//   - B: The success type of the output Either
//
// Parameters:
//   - ma: The ReaderResult to execute
//   - f: The function that takes the success value and returns an Either
//
// Returns:
//   - A ReaderResult that sequences both computations
//
// Example:
//
//	rr := readerresult.Of(42)
//	chained := readerresult.MonadChainEitherK(rr, func(x int) result.Result[string] {
//	    if x > 0 {
//	        return result.Of(fmt.Sprintf("positive: %d", x))
//	    }
//	    return result.Error[string](errors.New("not positive"))
//	})
//	result := chained(t.Context()) // Right("positive: 42")
//
//go:inline
func MonadChainEitherK[A, B any](ma ReaderResult[A], f func(A) Either[B]) ReaderResult[B] {
	return readereither.MonadChainEitherK(ma, f)
}

// ChainEitherK creates an operator that sequences a ReaderResult with a function that returns an Either.
// This is the curried version where the function is provided first.
//
// Type Parameters:
//   - A: The success type of the input ReaderResult
//   - B: The success type of the output Either
//
// Parameters:
//   - f: The function that takes the success value and returns an Either
//
// Returns:
//   - An Operator that sequences the computations
//
// Example:
//
//	validate := readerresult.ChainEitherK(func(x int) result.Result[int] {
//	    if x > 0 {
//	        return result.Of(x)
//	    }
//	    return result.Error[int](errors.New("must be positive"))
//	})
//	rr := readerresult.Of(42)
//	result := validate(rr)(t.Context()) // Right(42)
//
//go:inline
func ChainEitherK[A, B any](f func(A) Either[B]) func(ma ReaderResult[A]) ReaderResult[B] {
	return readereither.ChainEitherK[context.Context](f)
}

// ChainOptionK creates an operator that sequences a ReaderResult with a function that returns an Option.
// If the Option is None, the onNone function is called to generate an error.
//
// Type Parameters:
//   - A: The success type of the input ReaderResult
//   - B: The success type of the output Option
//
// Parameters:
//   - onNone: The function to generate an error when the Option is None
//
// Returns:
//   - A function that takes an Option Kleisli and returns an Operator
//
// Example:
//
//	chainOpt := readerresult.ChainOptionK[int, string](func() error {
//	    return errors.New("value not found")
//	})
//	optKleisli := func(x int) option.Option[string] {
//	    if x > 0 {
//	        return option.Some(fmt.Sprintf("value: %d", x))
//	    }
//	    return option.None[string]()
//	}
//	operator := chainOpt(optKleisli)
//	result := operator(readerresult.Of(42))(t.Context()) // Right("value: 42")
//
//go:inline
func ChainOptionK[A, B any](onNone func() error) func(option.Kleisli[A, B]) Operator[A, B] {
	return readereither.ChainOptionK[context.Context, A, B](onNone)
}

// MonadFlap applies a value to a ReaderResult containing a function.
// This is the monadic version that takes both the ReaderResult and the value as parameters.
// Flap is the reverse of Ap - instead of applying a function to a value, it applies a value to a function.
//
// Type Parameters:
//   - B: The output type
//   - A: The input type
//
// Parameters:
//   - fab: The ReaderResult containing the function
//   - a: The value to apply to the function
//
// Returns:
//   - A ReaderResult with the value applied to the function
//
// Example:
//
//	fabRR := readerresult.Of(func(x int) string {
//	    return fmt.Sprintf("value: %d", x)
//	})
//	result := readerresult.MonadFlap(fabRR, 42)(t.Context()) // Right("value: 42")
//
//go:inline
func MonadFlap[B, A any](fab ReaderResult[func(A) B], a A) ReaderResult[B] {
	return readereither.MonadFlap(fab, a)
}

// Flap creates an operator that applies a value to a ReaderResult containing a function.
// This is the curried version where the value is provided first.
// Flap is the reverse of Ap - instead of applying a function to a value, it applies a value to a function.
//
// Type Parameters:
//   - B: The output type
//   - A: The input type
//
// Parameters:
//   - a: The value to apply to the function
//
// Returns:
//   - An Operator that applies the value to a ReaderResult containing a function
//
// Example:
//
//	applyTo42 := readerresult.Flap[string](42)
//	fabRR := readerresult.Of(func(x int) string {
//	    return fmt.Sprintf("value: %d", x)
//	})
//	result := applyTo42(fabRR)(t.Context()) // Right("value: 42")
//
//go:inline
func Flap[B, A any](a A) Operator[func(A) B, B] {
	return readereither.Flap[context.Context, error, B](a)
}

// Read executes a ReaderResult by providing it with a context.Context value.
// This function "runs" the ReaderResult computation with the given context.
//
// Type Parameters:
//   - A: The success type of the ReaderResult
//
// Parameters:
//   - r: The context.Context to provide to the ReaderResult
//
// Returns:
//   - A function that takes a ReaderResult and returns its Result
//
// Example:
//
//	rr := readerresult.Of(42)
//	ctx := t.Context()
//	runWithCtx := readerresult.Read[int](ctx)
//	result := runWithCtx(rr) // Right(42)
//
//go:inline
func Read[A any](r context.Context) func(ReaderResult[A]) Result[A] {
	return readereither.Read[error, A](r)
}

// ReadEither executes a ReaderResult by providing it with a Result[context.Context].
// If the Result contains an error, that error is returned immediately.
// If the Result contains a context, the ReaderResult is executed with that context.
//
// Type Parameters:
//   - A: The success type of the ReaderResult
//
// Parameters:
//   - r: The Result[context.Context] to provide to the ReaderResult
//
// Returns:
//   - A function that takes a ReaderResult and returns its Result
//
// Example:
//
//	rr := readerresult.Of(42)
//	ctxResult := result.Of[error](t.Context())
//	runWithCtxResult := readerresult.ReadEither[int](ctxResult)
//	result := runWithCtxResult(rr) // Right(42)
//
//go:inline
func ReadEither[A any](r Result[context.Context]) func(ReaderResult[A]) Result[A] {
	return readereither.ReadEither[error, A](r)
}

// ReadResult executes a ReaderResult by providing it with a Result[context.Context].
// This is an alias for ReadEither.
//
// Type Parameters:
//   - A: The success type of the ReaderResult
//
// Parameters:
//   - r: The Result[context.Context] to provide to the ReaderResult
//
// Returns:
//   - A function that takes a ReaderResult and returns its Result
//
// Example:
//
//	rr := readerresult.Of(42)
//	ctxResult := result.Of[error](t.Context())
//	runWithCtxResult := readerresult.ReadResult[int](ctxResult)
//	result := runWithCtxResult(rr) // Right(42)
//
//go:inline
func ReadResult[A any](r Result[context.Context]) func(ReaderResult[A]) Result[A] {
	return readereither.ReadEither[error, A](r)
}

// MonadMapTo executes a ReaderResult computation, discards its success value, and returns a constant value.
// This is the monadic version that takes both the ReaderResult and the constant value as parameters.
//
// IMPORTANT: ReaderResult represents a side-effectful computation because it depends on context.Context,
// which is effectful (can be cancelled, has deadlines, carries values). For this reason, MonadMapTo WILL
// execute the original ReaderResult to allow any side effects to occur, then discard the success result
// and return the constant value. If the original computation fails, the error is preserved.
//
// Type Parameters:
//   - A: The success type of the first ReaderResult (will be discarded if successful)
//   - B: The type of the constant value to return on success
//
// Parameters:
//   - ma: The ReaderResult to execute (side effects will occur, success value discarded)
//   - b: The constant value to return if ma succeeds
//
// Returns:
//   - A ReaderResult that executes ma, preserves errors, but replaces success values with b
//
// Example:
//
//	type Config struct { Counter int }
//	increment := func(ctx context.Context) result.Result[int] {
//	    // Side effect: log the operation
//	    fmt.Println("incrementing")
//	    return result.Of(5)
//	}
//	r := readerresult.MonadMapTo(increment, "done")
//	result := r(t.Context()) // Prints "incrementing", returns Right("done")
//
//go:inline
func MonadMapTo[A, B any](ma ReaderResult[A], b B) ReaderResult[B] {
	return MonadMap(ma, reader.Of[A](b))
}

// MapTo creates an operator that executes a ReaderResult computation, discards its success value,
// and returns a constant value. This is the curried version where the constant value is provided first,
// returning a function that can be applied to any ReaderResult.
//
// IMPORTANT: ReaderResult represents a side-effectful computation because it depends on context.Context,
// which is effectful (can be cancelled, has deadlines, carries values). For this reason, MapTo WILL
// execute the input ReaderResult to allow any side effects to occur, then discard the success result
// and return the constant value. If the computation fails, the error is preserved.
//
// Type Parameters:
//   - A: The success type of the input ReaderResult (will be discarded if successful)
//   - B: The type of the constant value to return on success
//
// Parameters:
//   - b: The constant value to return on success
//
// Returns:
//   - An Operator that executes a ReaderResult[A], preserves errors, but replaces success with b
//
// Example:
//
//	logStep := func(ctx context.Context) result.Result[int] {
//	    fmt.Println("step executed")
//	    return result.Of(42)
//	}
//	toDone := readerresult.MapTo[int, string]("done")
//	pipeline := toDone(logStep)
//	result := pipeline(t.Context()) // Prints "step executed", returns Right("done")
//
// Example - In a functional pipeline:
//
//	step1 := func(ctx context.Context) result.Result[int] {
//	    fmt.Println("processing")
//	    return result.Of(1)
//	}
//	pipeline := F.Pipe1(
//	    step1,
//	    readerresult.MapTo[int, string]("complete"),
//	)
//	output := pipeline(t.Context()) // Prints "processing", returns Right("complete")
//
//go:inline
func MapTo[A, B any](b B) Operator[A, B] {
	return Map(reader.Of[A](b))
}

// MonadChainTo sequences two ReaderResult computations where the second ignores the first's success value.
// This is the monadic version that takes both ReaderResults as parameters.
//
// IMPORTANT: ReaderResult represents a side-effectful computation because it depends on context.Context,
// which is effectful (can be cancelled, has deadlines, carries values). For this reason, MonadChainTo WILL
// execute the first ReaderResult to allow any side effects to occur, then discard the success result and
// execute the second ReaderResult with the same context. If the first computation fails, the error is
// returned immediately without executing the second computation.
//
// Type Parameters:
//   - A: The success type of the first ReaderResult (will be discarded if successful)
//   - B: The success type of the second ReaderResult
//
// Parameters:
//   - ma: The first ReaderResult to execute (side effects will occur, success value discarded)
//   - b: The second ReaderResult to execute if ma succeeds
//
// Returns:
//   - A ReaderResult that executes ma, then b if ma succeeds, returning b's result
//
// Example:
//
//	logStart := func(ctx context.Context) result.Result[int] {
//	    fmt.Println("starting")
//	    return result.Of(1)
//	}
//	logEnd := func(ctx context.Context) result.Result[string] {
//	    fmt.Println("ending")
//	    return result.Of("done")
//	}
//	r := readerresult.MonadChainTo(logStart, logEnd)
//	result := r(t.Context()) // Prints "starting" then "ending", returns Right("done")
//
//go:inline
func MonadChainTo[A, B any](ma ReaderResult[A], b ReaderResult[B]) ReaderResult[B] {
	return MonadChain(ma, reader.Of[A](b))
}

// ChainTo creates an operator that sequences two ReaderResult computations where the second ignores
// the first's success value. This is the curried version where the second ReaderResult is provided first,
// returning a function that can be applied to any first ReaderResult.
//
// IMPORTANT: ReaderResult represents a side-effectful computation because it depends on context.Context,
// which is effectful (can be cancelled, has deadlines, carries values). For this reason, ChainTo WILL
// execute the first ReaderResult to allow any side effects to occur, then discard the success result and
// execute the second ReaderResult with the same context. If the first computation fails, the error is
// returned immediately without executing the second computation.
//
// Type Parameters:
//   - A: The success type of the first ReaderResult (will be discarded if successful)
//   - B: The success type of the second ReaderResult
//
// Parameters:
//   - b: The second ReaderResult to execute after the first succeeds
//
// Returns:
//   - An Operator that executes the first ReaderResult, then b if successful
//
// Example:
//
//	logEnd := func(ctx context.Context) result.Result[string] {
//	    fmt.Println("ending")
//	    return result.Of("done")
//	}
//	thenLogEnd := readerresult.ChainTo[int, string](logEnd)
//
//	logStart := func(ctx context.Context) result.Result[int] {
//	    fmt.Println("starting")
//	    return result.Of(1)
//	}
//	pipeline := thenLogEnd(logStart)
//	result := pipeline(t.Context()) // Prints "starting" then "ending", returns Right("done")
//
// Example - In a functional pipeline:
//
//	step1 := func(ctx context.Context) result.Result[int] {
//	    fmt.Println("step 1")
//	    return result.Of(1)
//	}
//	step2 := func(ctx context.Context) result.Result[string] {
//	    fmt.Println("step 2")
//	    return result.Of("complete")
//	}
//	pipeline := F.Pipe1(
//	    step1,
//	    readerresult.ChainTo[int, string](step2),
//	)
//	output := pipeline(t.Context()) // Prints "step 1" then "step 2", returns Right("complete")
//
//go:inline
func ChainTo[A, B any](b ReaderResult[B]) Operator[A, B] {
	return Chain(reader.Of[A](b))
}

// MonadChainFirst sequences two ReaderResult computations, executing the second for its side effects
// but returning the value from the first. This is the monadic version that takes both the ReaderResult
// and the Kleisli function as parameters.
//
// IMPORTANT: Combining with IO operations makes sense because ReaderResult already has side effects
// due to context.Context (cancellation, deadlines, values). ChainFirst executes both computations
// for their side effects, which is natural when working with effectful computations.
//
// Type Parameters:
//   - A: The success type of the first ReaderResult (returned value)
//   - B: The success type of the second ReaderResult (discarded)
//
// Parameters:
//   - ma: The first ReaderResult to execute
//   - f: The Kleisli function that takes the success value and returns a second ReaderResult
//
// Returns:
//   - A ReaderResult that executes both computations but returns the first's value
//
// Example:
//
//	rr := readerresult.Of(42)
//	withLogging := readerresult.MonadChainFirst(rr, func(x int) readerresult.ReaderResult[string] {
//	    return func(ctx context.Context) result.Result[string] {
//	        fmt.Printf("Value: %d\n", x)
//	        return result.Of("logged")
//	    }
//	})
//	result := withLogging(t.Context()) // Prints "Value: 42", returns Right(42)
//
//go:inline
func MonadChainFirst[A, B any](ma ReaderResult[A], f Kleisli[A, B]) ReaderResult[A] {
	return chain.MonadChainFirst(
		MonadChain,
		MonadMap,
		ma,
		F.Flow2(f, WithContext),
	)
}

// ChainFirst creates an operator that sequences two ReaderResult computations, executing the second
// for its side effects but returning the value from the first. This is the curried version where
// the Kleisli function is provided first.
//
// IMPORTANT: Combining with IO operations makes sense because ReaderResult already has side effects
// due to context.Context (cancellation, deadlines, values). ChainFirst executes both computations
// for their side effects, which is natural when working with effectful computations.
//
// Type Parameters:
//   - A: The success type of the first ReaderResult (returned value)
//   - B: The success type of the second ReaderResult (discarded)
//
// Parameters:
//   - f: The Kleisli function that takes the success value and returns a second ReaderResult
//
// Returns:
//   - An Operator that executes both computations but returns the first's value
//
// Example:
//
//	logValue := readerresult.ChainFirst(func(x int) readerresult.ReaderResult[string] {
//	    return func(ctx context.Context) result.Result[string] {
//	        fmt.Printf("Value: %d\n", x)
//	        return result.Of("logged")
//	    }
//	})
//	result := logValue(readerresult.Of(42))(t.Context()) // Prints "Value: 42", returns Right(42)
//
//go:inline
func ChainFirst[A, B any](f Kleisli[A, B]) Operator[A, A] {
	return chain.ChainFirst(
		Chain,
		Map,
		F.Flow2(f, WithContext),
	)
}

// MonadChainIOK sequences a ReaderResult with an IO computation, lifting the IO into ReaderResult.
// This is the monadic version that takes both the ReaderResult and the IO Kleisli function as parameters.
//
// IMPORTANT: Combining IO with ReaderResult makes sense because both represent side-effectful computations.
// ReaderResult has side effects from context.Context (cancellation, deadlines, values), and IO has side
// effects from IO operations. This combination allows context-aware error handling with IO operations.
//
// Type Parameters:
//   - A: The success type of the input ReaderResult
//   - B: The success type of the IO computation
//
// Parameters:
//   - ma: The ReaderResult to execute
//   - f: The IO Kleisli function that takes the success value and returns an IO computation
//
// Returns:
//   - A ReaderResult that sequences both computations
//
// Example:
//
//	rr := readerresult.Of(42)
//	withIO := readerresult.MonadChainIOK(rr, func(x int) func() string {
//	    return func() string {
//	        fmt.Printf("Value: %d\n", x)
//	        return "done"
//	    }
//	})
//	result := withIO(t.Context()) // Prints "Value: 42", returns Right("done")
//
//go:inline
func MonadChainIOK[A, B any](ma ReaderResult[A], f io.Kleisli[A, B]) ReaderResult[B] {
	return fromio.MonadChainIOK(
		MonadChain[A, B],
		FromIO[B],
		ma,
		f,
	)
}

// ChainIOK creates an operator that sequences a ReaderResult with an IO computation.
// This is the curried version where the IO Kleisli function is provided first.
//
// IMPORTANT: Combining IO with ReaderResult makes sense because both represent side-effectful computations.
// ReaderResult has side effects from context.Context (cancellation, deadlines, values), and IO has side
// effects from IO operations. This combination allows context-aware error handling with IO operations.
//
// Type Parameters:
//   - A: The success type of the input ReaderResult
//   - B: The success type of the IO computation
//
// Parameters:
//   - f: The IO Kleisli function that takes the success value and returns an IO computation
//
// Returns:
//   - An Operator that sequences the computations
//
// Example:
//
//	logIO := readerresult.ChainIOK(func(x int) func() string {
//	    return func() string {
//	        fmt.Printf("Value: %d\n", x)
//	        return "logged"
//	    }
//	})
//	result := logIO(readerresult.Of(42))(t.Context()) // Prints "Value: 42", returns Right("logged")
//
//go:inline
func ChainIOK[A, B any](f io.Kleisli[A, B]) Operator[A, B] {
	return fromio.ChainIOK(
		Chain[A, B],
		FromIO[B],
		f,
	)
}

// MonadChainFirstIOK sequences a ReaderResult with an IO computation for its side effects,
// but returns the original value. This is the monadic version.
//
// IMPORTANT: Combining IO with ReaderResult makes sense because both represent side-effectful computations.
// This function executes the IO operation for its side effects (like logging or metrics) while preserving
// the original value, which is natural when working with effectful computations.
//
// Type Parameters:
//   - A: The success type of the ReaderResult (returned value)
//   - B: The success type of the IO computation (discarded)
//
// Parameters:
//   - ma: The ReaderResult to execute
//   - f: The IO Kleisli function for side effects
//
// Returns:
//   - A ReaderResult that executes both but returns the original value
//
// Example:
//
//	rr := readerresult.Of(42)
//	withLog := readerresult.MonadChainFirstIOK(rr, func(x int) func() string {
//	    return func() string {
//	        fmt.Printf("Processing: %d\n", x)
//	        return "logged"
//	    }
//	})
//	result := withLog(t.Context()) // Prints "Processing: 42", returns Right(42)
//
//go:inline
func MonadChainFirstIOK[A, B any](ma ReaderResult[A], f io.Kleisli[A, B]) ReaderResult[A] {
	return fromio.MonadChainFirstIOK(
		MonadChain[A, A],
		MonadMap[B, A],
		FromIO[B],
		ma,
		f,
	)
}

// MonadTapIOK is an alias for MonadChainFirstIOK. It sequences a ReaderResult with an IO computation
// for its side effects, but returns the original value.
//
// IMPORTANT: Combining IO with ReaderResult makes sense because both represent side-effectful computations.
// Tap executes the IO operation for its side effects while preserving the original value.
//
// Type Parameters:
//   - A: The success type of the ReaderResult (returned value)
//   - B: The success type of the IO computation (discarded)
//
// Parameters:
//   - ma: The ReaderResult to execute
//   - f: The IO Kleisli function for side effects
//
// Returns:
//   - A ReaderResult that executes both but returns the original value
//
// Example:
//
//	rr := readerresult.Of(42)
//	withLog := readerresult.MonadTapIOK(rr, func(x int) func() string {
//	    return func() string {
//	        fmt.Printf("Tapping: %d\n", x)
//	        return "logged"
//	    }
//	})
//	result := withLog(t.Context()) // Prints "Tapping: 42", returns Right(42)
//
//go:inline
func MonadTapIOK[A, B any](ma ReaderResult[A], f io.Kleisli[A, B]) ReaderResult[A] {
	return fromio.MonadChainFirstIOK(
		MonadChain[A, A],
		MonadMap[B, A],
		FromIO[B],
		ma,
		f,
	)
}

// ChainFirstIOK creates an operator that sequences a ReaderResult with an IO computation for its
// side effects, but returns the original value. This is the curried version.
//
// IMPORTANT: Combining IO with ReaderResult makes sense because both represent side-effectful computations.
// This function executes the IO operation for its side effects while preserving the original value.
//
// Type Parameters:
//   - A: The success type of the ReaderResult (returned value)
//   - B: The success type of the IO computation (discarded)
//
// Parameters:
//   - f: The IO Kleisli function for side effects
//
// Returns:
//   - An Operator that executes both but returns the original value
//
// Example:
//
//	logIO := readerresult.ChainFirstIOK(func(x int) func() string {
//	    return func() string {
//	        fmt.Printf("Processing: %d\n", x)
//	        return "logged"
//	    }
//	})
//	result := logIO(readerresult.Of(42))(t.Context()) // Prints "Processing: 42", returns Right(42)
//
//go:inline
func ChainFirstIOK[A, B any](f io.Kleisli[A, B]) Operator[A, A] {
	return fromio.ChainFirstIOK(
		Chain[A, A],
		Map[B, A],
		FromIO[B],
		f,
	)
}

// TapIOK is an alias for ChainFirstIOK. It creates an operator that sequences a ReaderResult with
// an IO computation for its side effects, but returns the original value.
//
// IMPORTANT: Combining IO with ReaderResult makes sense because both represent side-effectful computations.
// Tap executes the IO operation for its side effects while preserving the original value.
//
// Type Parameters:
//   - A: The success type of the ReaderResult (returned value)
//   - B: The success type of the IO computation (discarded)
//
// Parameters:
//   - f: The IO Kleisli function for side effects
//
// Returns:
//   - An Operator that executes both but returns the original value
//
// Example:
//
//	tapLog := readerresult.TapIOK(func(x int) func() string {
//	    return func() string {
//	        fmt.Printf("Tapping: %d\n", x)
//	        return "logged"
//	    }
//	})
//	result := tapLog(readerresult.Of(42))(t.Context()) // Prints "Tapping: 42", returns Right(42)
//
//go:inline
func TapIOK[A, B any](f io.Kleisli[A, B]) Operator[A, A] {
	return fromio.ChainFirstIOK(
		Chain[A, A],
		Map[B, A],
		FromIO[B],
		f,
	)
}

// ChainIOEitherK creates an operator that sequences a ReaderResult with an IOResult computation.
//
// IMPORTANT: Combining IOResult with ReaderResult makes sense because both represent side-effectful
// computations. ReaderResult has side effects from context.Context, and IOResult has side effects
// from IO operations with error handling. This combination provides context-aware error handling
// with IO operations.
//
// Type Parameters:
//   - A: The success type of the input ReaderResult
//   - B: The success type of the IOResult computation
//
// Parameters:
//   - f: The IOResult Kleisli function
//
// Returns:
//   - An Operator that sequences the computations
//
// Example:
//
//	ioResultOp := readerresult.ChainIOEitherK(func(x int) func() result.Result[string] {
//	    return func() result.Result[string] {
//	        if x > 0 {
//	            return result.Of(fmt.Sprintf("positive: %d", x))
//	        }
//	        return result.Error[string](errors.New("not positive"))
//	    }
//	})
//	result := ioResultOp(readerresult.Of(42))(t.Context()) // Right("positive: 42")
//
//go:inline
func ChainIOEitherK[A, B any](f ioresult.Kleisli[A, B]) Operator[A, B] {
	return fromioeither.ChainIOEitherK(
		Chain[A, B],
		FromIOResult[B],
		f,
	)
}

// ChainIOResultK is an alias for ChainIOEitherK. It creates an operator that sequences a ReaderResult
// with an IOResult computation.
//
// IMPORTANT: Combining IOResult with ReaderResult makes sense because both represent side-effectful
// computations. This provides context-aware error handling with IO operations.
//
// Type Parameters:
//   - A: The success type of the input ReaderResult
//   - B: The success type of the IOResult computation
//
// Parameters:
//   - f: The IOResult Kleisli function
//
// Returns:
//   - An Operator that sequences the computations
//
// Example:
//
//	ioResultOp := readerresult.ChainIOResultK(func(x int) func() result.Result[string] {
//	    return func() result.Result[string] {
//	        return result.Of(fmt.Sprintf("value: %d", x))
//	    }
//	})
//	result := ioResultOp(readerresult.Of(42))(t.Context()) // Right("value: 42")
//
//go:inline
func ChainIOResultK[A, B any](f ioresult.Kleisli[A, B]) Operator[A, B] {
	return fromioeither.ChainIOEitherK(
		Chain[A, B],
		FromIOResult[B],
		f,
	)
}

// ReadIO executes a ReaderResult by providing it with a context obtained from an IO computation.
//
// IMPORTANT: Combining IO with ReaderResult makes sense because both represent side-effectful computations.
// This allows the context itself to be obtained through an IO operation.
//
// Type Parameters:
//   - A: The success type of the ReaderResult
//
// Parameters:
//   - r: The IO computation that produces a context.Context
//
// Returns:
//   - A function that takes a ReaderResult and returns an IOResult
//
// Example:
//
//	getCtx := func() context.Context { return t.Context() }
//	rr := readerresult.Of(42)
//	runWithIO := readerresult.ReadIO[int](getCtx)
//	ioResult := runWithIO(rr)
//	result := ioResult() // Right(42)
//
//go:inline
func ReadIO[A any](r IO[context.Context]) func(ReaderResult[A]) IOResult[A] {
	return func(rr ReaderResult[A]) IOResult[A] {
		return func() Result[A] {
			return rr(r())
		}
	}
}

// ReadIOEither executes a ReaderResult by providing it with a context obtained from an IOResult computation.
// If the IOResult contains an error, that error is returned immediately.
//
// IMPORTANT: Combining IOResult with ReaderResult makes sense because both represent side-effectful
// computations with error handling. This allows the context itself to be obtained through an IO operation
// that may fail.
//
// Type Parameters:
//   - A: The success type of the ReaderResult
//
// Parameters:
//   - r: The IOResult computation that produces a context.Context
//
// Returns:
//   - A function that takes a ReaderResult and returns an IOResult
//
// Example:
//
//	getCtx := func() result.Result[context.Context] {
//	    return result.Of[error](t.Context())
//	}
//	rr := readerresult.Of(42)
//	runWithIOResult := readerresult.ReadIOEither[int](getCtx)
//	ioResult := runWithIOResult(rr)
//	result := ioResult() // Right(42)
//
//go:inline
func ReadIOEither[A any](r IOResult[context.Context]) func(ReaderResult[A]) IOResult[A] {
	return func(rr ReaderResult[A]) IOResult[A] {
		return F.Pipe1(
			r,
			ioresult.ChainResultK(rr),
		)
	}
}

// ReadIOResult is an alias for ReadIOEither. It executes a ReaderResult by providing it with a context
// obtained from an IOResult computation.
//
// IMPORTANT: Combining IOResult with ReaderResult makes sense because both represent side-effectful
// computations with error handling.
//
// Type Parameters:
//   - A: The success type of the ReaderResult
//
// Parameters:
//   - r: The IOResult computation that produces a context.Context
//
// Returns:
//   - A function that takes a ReaderResult and returns an IOResult
//
// Example:
//
//	getCtx := func() result.Result[context.Context] {
//	    return result.Of[error](t.Context())
//	}
//	rr := readerresult.Of(42)
//	runWithIOResult := readerresult.ReadIOResult[int](getCtx)
//	ioResult := runWithIOResult(rr)
//	result := ioResult() // Right(42)
//
//go:inline
func ReadIOResult[A any](r IOResult[context.Context]) func(ReaderResult[A]) IOResult[A] {
	return ReadIOEither[A](r)
}

// ChainFirstLeft executes a computation on the Left (error) value for its side effects,
// but preserves the original error. This is useful for error logging or metrics.
//
// Type Parameters:
//   - A: The success type of the ReaderResult
//   - B: The success type of the error handler (discarded)
//
// Parameters:
//   - f: The Kleisli function that handles the error
//
// Returns:
//   - An Operator that executes the error handler but preserves the original error
//
// Example:
//
//	logError := readerresult.ChainFirstLeft[int](func(err error) readerresult.ReaderResult[string] {
//	    return func(ctx context.Context) result.Result[string] {
//	        fmt.Printf("Error occurred: %v\n", err)
//	        return result.Of("logged")
//	    }
//	})
//	rr := readerresult.Left[int](errors.New("failed"))
//	result := logError(rr)(t.Context()) // Prints "Error occurred: failed", returns Left(error("failed"))
//
//go:inline
func ChainFirstLeft[A, B any](f Kleisli[error, B]) Operator[A, A] {
	return readereither.ChainFirstLeft[A](f)
}

// ChainFirstLeftIOK executes an IO computation on the Left (error) value for its side effects,
// but preserves the original error. This is useful for error logging or metrics with IO operations.
//
// IMPORTANT: Combining IO with ReaderResult makes sense because both represent side-effectful computations.
// This allows error handling with side effects like logging to external systems.
//
// Type Parameters:
//   - A: The success type of the ReaderResult
//   - B: The success type of the IO error handler (discarded)
//
// Parameters:
//   - f: The IO Kleisli function that handles the error
//
// Returns:
//   - An Operator that executes the IO error handler but preserves the original error
//
// Example:
//
//	logErrorIO := readerresult.ChainFirstLeftIOK[int](func(err error) func() string {
//	    return func() string {
//	        fmt.Printf("Error: %v\n", err)
//	        return "logged"
//	    }
//	})
//	rr := readerresult.Left[int](errors.New("failed"))
//	result := logErrorIO(rr)(t.Context()) // Prints "Error: failed", returns Left(error("failed"))
//
//go:inline
func ChainFirstLeftIOK[A, B any](f io.Kleisli[error, B]) Operator[A, A] {
	return ChainFirstLeft[A](function.Flow2(
		f,
		FromIO,
	))
}

// TapLeftIOK is an alias for ChainFirstLeftIOK. It executes an IO computation on the Left (error)
// value for its side effects, but preserves the original error.
//
// IMPORTANT: Combining IO with ReaderResult makes sense because both represent side-effectful computations.
// Tap allows error handling with side effects while preserving the error.
//
// Type Parameters:
//   - A: The success type of the ReaderResult
//   - B: The success type of the IO error handler (discarded)
//
// Parameters:
//   - f: The IO Kleisli function that handles the error
//
// Returns:
//   - An Operator that executes the IO error handler but preserves the original error
//
// Example:
//
//	tapErrorIO := readerresult.TapLeftIOK[int](func(err error) func() string {
//	    return func() string {
//	        fmt.Printf("Tapping error: %v\n", err)
//	        return "logged"
//	    }
//	})
//	rr := readerresult.Left[int](errors.New("failed"))
//	result := tapErrorIO(rr)(t.Context()) // Prints "Tapping error: failed", returns Left(error("failed"))
//
//go:inline
func TapLeftIOK[A, B any](f io.Kleisli[error, B]) Operator[A, A] {
	return ChainFirstLeftIOK[A](f)
}
