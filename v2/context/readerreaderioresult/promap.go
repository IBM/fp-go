package readerreaderioresult

import (
	"context"

	"github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioresult"
	RRIOE "github.com/IBM/fp-go/v2/readerreaderioeither"
	"github.com/IBM/fp-go/v2/result"
)

// Local modifies the outer environment before passing it to a computation.
// Useful for providing different configurations to sub-computations.
//
//go:inline
func Local[A, R1, R2 any](f func(R2) R1) func(ReaderReaderIOResult[R1, A]) ReaderReaderIOResult[R2, A] {
	return RRIOE.Local[context.Context, error, A](f)
}

// LocalIOK transforms the outer environment of a ReaderReaderIOResult using an IO-based Kleisli arrow.
// It allows you to modify the outer environment through an effectful computation before
// passing it to the ReaderReaderIOResult.
//
// This is useful when the outer environment transformation itself requires IO effects,
// such as reading from a file, making a network call, or accessing system resources,
// but these effects cannot fail (or failures are not relevant).
//
// The transformation happens in two stages:
//  1. The IO effect f is executed with the R2 environment to produce an R1 value
//  2. The resulting R1 value is passed as the outer environment to the ReaderReaderIOResult[R1, A]
//
// Type Parameters:
//   - A: The success type produced by the ReaderReaderIOResult
//   - R1: The original outer environment type expected by the ReaderReaderIOResult
//   - R2: The new input outer environment type
//
// Parameters:
//   - f: An IO Kleisli arrow that transforms R2 to R1 with IO effects
//
// Returns:
//   - A function that takes a ReaderReaderIOResult[R1, A] and returns a ReaderReaderIOResult[R2, A]
//
//go:inline
func LocalIOK[A, R1, R2 any](f io.Kleisli[R2, R1]) func(ReaderReaderIOResult[R1, A]) ReaderReaderIOResult[R2, A] {
	return RRIOE.LocalIOK[context.Context, error, A](f)
}

// LocalIOEitherK transforms the outer environment of a ReaderReaderIOResult using an IOResult-based Kleisli arrow.
// It allows you to modify the outer environment through an effectful computation that can fail before
// passing it to the ReaderReaderIOResult.
//
// This is useful when the outer environment transformation itself requires IO effects that can fail,
// such as reading from a file that might not exist, making a network call that might timeout,
// or parsing data that might be invalid.
//
// The transformation happens in two stages:
//  1. The IOResult effect f is executed with the R2 environment to produce Result[R1]
//  2. If successful (Ok), the R1 value is passed as the outer environment to the ReaderReaderIOResult[R1, A]
//  3. If failed (Err), the error is propagated without executing the ReaderReaderIOResult
//
// Type Parameters:
//   - A: The success type produced by the ReaderReaderIOResult
//   - R1: The original outer environment type expected by the ReaderReaderIOResult
//   - R2: The new input outer environment type
//
// Parameters:
//   - f: An IOResult Kleisli arrow that transforms R2 to R1 with IO effects that can fail
//
// Returns:
//   - A function that takes a ReaderReaderIOResult[R1, A] and returns a ReaderReaderIOResult[R2, A]
//
//go:inline
func LocalIOEitherK[A, R1, R2 any](f ioresult.Kleisli[R2, R1]) func(ReaderReaderIOResult[R1, A]) ReaderReaderIOResult[R2, A] {
	return RRIOE.LocalIOEitherK[context.Context, A](f)
}

// LocalIOResultK transforms the outer environment of a ReaderReaderIOResult using an IOResult-based Kleisli arrow.
// This is a type-safe alias for LocalIOEitherK specialized for Result types (which use error as the error type).
//
// It allows you to modify the outer environment through an effectful computation that can fail before
// passing it to the ReaderReaderIOResult.
//
// The transformation happens in two stages:
//  1. The IOResult effect f is executed with the R2 environment to produce Result[R1]
//  2. If successful (Ok), the R1 value is passed as the outer environment to the ReaderReaderIOResult[R1, A]
//  3. If failed (Err), the error is propagated without executing the ReaderReaderIOResult
//
// Type Parameters:
//   - A: The success type produced by the ReaderReaderIOResult
//   - R1: The original outer environment type expected by the ReaderReaderIOResult
//   - R2: The new input outer environment type
//
// Parameters:
//   - f: An IOResult Kleisli arrow that transforms R2 to R1 with IO effects that can fail
//
// Returns:
//   - A function that takes a ReaderReaderIOResult[R1, A] and returns a ReaderReaderIOResult[R2, A]
//
//go:inline
func LocalIOResultK[A, R1, R2 any](f ioresult.Kleisli[R2, R1]) func(ReaderReaderIOResult[R1, A]) ReaderReaderIOResult[R2, A] {
	return RRIOE.LocalIOEitherK[context.Context, A](f)
}

//go:inline
func LocalResultK[A, R1, R2 any](f result.Kleisli[R2, R1]) func(ReaderReaderIOResult[R1, A]) ReaderReaderIOResult[R2, A] {
	return RRIOE.LocalEitherK[context.Context, A](f)
}

// LocalReaderIOEitherK transforms the outer environment of a ReaderReaderIOResult using a ReaderIOResult-based Kleisli arrow.
// It allows you to modify the outer environment through a computation that depends on the inner context
// and can perform IO effects that may fail.
//
// This is useful when the outer environment transformation requires access to the inner context (e.g., context.Context)
// and may perform IO operations that can fail, such as database queries, API calls, or file operations.
//
// The transformation happens in three stages:
//  1. The ReaderIOResult effect f is executed with the R2 outer environment and inner context
//  2. If successful (Ok), the R1 value is passed as the outer environment to the ReaderReaderIOResult[R1, A]
//  3. If failed (Err), the error is propagated without executing the ReaderReaderIOResult
//
// Type Parameters:
//   - A: The success type produced by the ReaderReaderIOResult
//   - R1: The original outer environment type expected by the ReaderReaderIOResult
//   - R2: The new input outer environment type
//
// Parameters:
//   - f: A ReaderIOResult Kleisli arrow that transforms R2 to R1 with context-aware IO effects that can fail
//
// Returns:
//   - A function that takes a ReaderReaderIOResult[R1, A] and returns a ReaderReaderIOResult[R2, A]
//
//go:inline
func LocalReaderIOEitherK[A, R1, R2 any](f readerioresult.Kleisli[R2, R1]) func(ReaderReaderIOResult[R1, A]) ReaderReaderIOResult[R2, A] {
	return RRIOE.LocalReaderIOEitherK[A](f)
}

// LocalReaderIOResultK transforms the outer environment of a ReaderReaderIOResult using a ReaderIOResult-based Kleisli arrow.
// This is a type-safe alias for LocalReaderIOEitherK specialized for Result types (which use error as the error type).
//
// It allows you to modify the outer environment through a computation that depends on the inner context
// and can perform IO effects that may fail.
//
// The transformation happens in three stages:
//  1. The ReaderIOResult effect f is executed with the R2 outer environment and inner context
//  2. If successful (Ok), the R1 value is passed as the outer environment to the ReaderReaderIOResult[R1, A]
//  3. If failed (Err), the error is propagated without executing the ReaderReaderIOResult
//
// Type Parameters:
//   - A: The success type produced by the ReaderReaderIOResult
//   - R1: The original outer environment type expected by the ReaderReaderIOResult
//   - R2: The new input outer environment type
//
// Parameters:
//   - f: A ReaderIOResult Kleisli arrow that transforms R2 to R1 with context-aware IO effects that can fail
//
// Returns:
//   - A function that takes a ReaderReaderIOResult[R1, A] and returns a ReaderReaderIOResult[R2, A]
//
//go:inline
func LocalReaderIOResultK[A, R1, R2 any](f readerioresult.Kleisli[R2, R1]) func(ReaderReaderIOResult[R1, A]) ReaderReaderIOResult[R2, A] {
	return RRIOE.LocalReaderIOEitherK[A](f)
}

//go:inline
func LocalReaderReaderIOEitherK[A, R1, R2 any](f Kleisli[R2, R2, R1]) func(ReaderReaderIOResult[R1, A]) ReaderReaderIOResult[R2, A] {
	return RRIOE.LocalReaderReaderIOEitherK[A](f)
}
