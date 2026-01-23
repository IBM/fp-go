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

package readereither

import (
	ET "github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/eithert"
	"github.com/IBM/fp-go/v2/internal/fromeither"
	"github.com/IBM/fp-go/v2/internal/fromreader"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/readert"
	"github.com/IBM/fp-go/v2/reader"
)

func FromEither[E, L, A any](e Either[L, A]) ReaderEither[E, L, A] {
	return reader.Of[E](e)
}

func RightReader[L, E, A any](r Reader[E, A]) ReaderEither[E, L, A] {
	return eithert.RightF(reader.MonadMap[E, A, Either[L, A]], r)
}

func LeftReader[A, E, L any](l Reader[E, L]) ReaderEither[E, L, A] {
	return eithert.LeftF(reader.MonadMap[E, L, Either[L, A]], l)
}

func Left[E, A, L any](l L) ReaderEither[E, L, A] {
	return eithert.Left(reader.Of[E, Either[L, A]], l)
}

func Right[E, L, A any](r A) ReaderEither[E, L, A] {
	return eithert.Right(reader.Of[E, Either[L, A]], r)
}

func FromReader[L, E, A any](r Reader[E, A]) ReaderEither[E, L, A] {
	return RightReader[L](r)
}

func MonadMap[E, L, A, B any](fa ReaderEither[E, L, A], f func(A) B) ReaderEither[E, L, B] {
	return readert.MonadMap[ReaderEither[E, L, A], ReaderEither[E, L, B]](ET.MonadMap[L, A, B], fa, f)
}

func Map[E, L, A, B any](f func(A) B) func(ReaderEither[E, L, A]) ReaderEither[E, L, B] {
	return readert.Map[ReaderEither[E, L, A], ReaderEither[E, L, B]](ET.Map[L, A, B], f)
}

func MonadChain[E, L, A, B any](ma ReaderEither[E, L, A], f func(A) ReaderEither[E, L, B]) ReaderEither[E, L, B] {
	return readert.MonadChain(ET.MonadChain[L, A, B], ma, f)
}

func Chain[E, L, A, B any](f func(A) ReaderEither[E, L, B]) func(ReaderEither[E, L, A]) ReaderEither[E, L, B] {
	return readert.Chain[ReaderEither[E, L, A]](ET.Chain[L, A, B], f)
}

func MonadChainReaderK[L, E, A, B any](ma ReaderEither[E, L, A], f reader.Kleisli[E, A, B]) ReaderEither[E, L, B] {
	return MonadChain(ma, function.Flow2(f, FromReader[L, E, B]))
}

func ChainReaderK[L, E, A, B any](f reader.Kleisli[E, A, B]) func(ReaderEither[E, L, A]) ReaderEither[E, L, B] {
	return Chain(function.Flow2(f, FromReader[L, E, B]))
}

func Of[E, L, A any](a A) ReaderEither[E, L, A] {
	return readert.MonadOf[ReaderEither[E, L, A]](ET.Of[L, A], a)
}

func MonadAp[B, E, L, A any](fab ReaderEither[E, L, func(A) B], fa ReaderEither[E, L, A]) ReaderEither[E, L, B] {
	return readert.MonadAp[ReaderEither[E, L, A], ReaderEither[E, L, B], ReaderEither[E, L, func(A) B], E, A](ET.MonadAp[B, L, A], fab, fa)
}

func Ap[B, E, L, A any](fa ReaderEither[E, L, A]) func(ReaderEither[E, L, func(A) B]) ReaderEither[E, L, B] {
	return readert.Ap[ReaderEither[E, L, A], ReaderEither[E, L, B], ReaderEither[E, L, func(A) B], E, A](ET.Ap[B, L, A], fa)
}

func FromPredicate[E, L, A any](pred func(A) bool, onFalse func(A) L) func(A) ReaderEither[E, L, A] {
	return fromeither.FromPredicate(FromEither[E, L, A], pred, onFalse)
}

func Fold[E, L, A, B any](onLeft func(L) Reader[E, B], onRight func(A) Reader[E, B]) func(ReaderEither[E, L, A]) Reader[E, B] {
	return eithert.MatchE(reader.MonadChain[E, Either[L, A], B], onLeft, onRight)
}

func GetOrElse[E, L, A any](onLeft func(L) Reader[E, A]) func(ReaderEither[E, L, A]) Reader[E, A] {
	return eithert.GetOrElse(reader.MonadChain[E, Either[L, A], A], reader.Of[E, A], onLeft)
}

func OrLeft[A, L1, E, L2 any](onLeft func(L1) Reader[E, L2]) func(ReaderEither[E, L1, A]) ReaderEither[E, L2, A] {
	return eithert.OrLeft(
		reader.MonadChain[E, Either[L1, A], Either[L2, A]],
		reader.MonadMap[E, L2, Either[L2, A]],
		reader.Of[E, Either[L2, A]],
		onLeft,
	)
}

func Ask[E, L any]() ReaderEither[E, L, E] {
	return fromreader.Ask(FromReader[L, E, E])()
}

func Asks[L, E, A any](r Reader[E, A]) ReaderEither[E, L, A] {
	return fromreader.Asks(FromReader[L, E, A])(r)
}

func MonadChainEitherK[E, L, A, B any](ma ReaderEither[E, L, A], f func(A) Either[L, B]) ReaderEither[E, L, B] {
	return fromeither.MonadChainEitherK(
		MonadChain[E, L, A, B],
		FromEither[E, L, B],
		ma,
		f,
	)
}

func ChainEitherK[E, L, A, B any](f func(A) Either[L, B]) func(ma ReaderEither[E, L, A]) ReaderEither[E, L, B] {
	return fromeither.ChainEitherK(
		Chain[E, L, A, B],
		FromEither[E, L, B],
		f,
	)
}

func ChainOptionK[E, A, B, L any](onNone func() L) func(func(A) Option[B]) func(ReaderEither[E, L, A]) ReaderEither[E, L, B] {
	return fromeither.ChainOptionK(MonadChain[E, L, A, B], FromEither[E, L, B], onNone)
}

func Flatten[E, L, A any](mma ReaderEither[E, L, ReaderEither[E, L, A]]) ReaderEither[E, L, A] {
	return MonadChain(mma, function.Identity[ReaderEither[E, L, A]])
}

func MonadBiMap[E, E1, E2, A, B any](fa ReaderEither[E, E1, A], f func(E1) E2, g func(A) B) ReaderEither[E, E2, B] {
	return eithert.MonadBiMap(reader.MonadMap[E, Either[E1, A], Either[E2, B]], fa, f, g)
}

// BiMap maps a pair of functions over the two type arguments of the bifunctor.
func BiMap[E, E1, E2, A, B any](f func(E1) E2, g func(A) B) func(ReaderEither[E, E1, A]) ReaderEither[E, E2, B] {
	return eithert.BiMap(reader.Map[E, Either[E1, A], Either[E2, B]], f, g)
}

// Local changes the value of the local context during the execution of the action `ma` (similar to `Contravariant`'s
// `contramap`).
func Local[E, A, R1, R2 any](f func(R2) R1) func(ReaderEither[R1, E, A]) ReaderEither[R2, E, A] {
	return reader.Local[Either[E, A]](f)
}

// Read applies a context to a reader to obtain its value
func Read[E1, A, E any](e E) func(ReaderEither[E, E1, A]) Either[E1, A] {
	return reader.Read[Either[E1, A]](e)
}

// ReadEither applies a context wrapped in an Either to a ReaderEither to obtain its result.
// This function is useful when the context itself may be absent or invalid (represented as Left),
// allowing you to conditionally execute a ReaderEither computation based on the availability
// of the required context.
//
// If the context Either is Left, it short-circuits and returns Left without executing the ReaderEither.
// If the context Either is Right, it extracts the context value and applies it to the ReaderEither,
// returning the resulting Either.
//
// This is particularly useful in scenarios where:
//   - Configuration or dependencies may be missing or invalid
//   - You want to chain context validation with computation execution
//   - You need to propagate context errors through your computation pipeline
//
// Type Parameters:
//   - E1: The error type (Left value) of both the input Either and the ReaderEither result
//   - A: The success type (Right value) of the ReaderEither result
//   - E: The context/environment type that the ReaderEither depends on
//
// Parameters:
//   - e: An Either[E1, E] representing the context that may or may not be available
//
// Returns:
//   - A function that takes a ReaderEither[E, E1, A] and returns Either[E1, A]
//
// Example:
//
//	type Config struct{ apiKey string }
//	type ConfigError struct{ msg string }
//
//	// A computation that needs config
//	fetchData := func(cfg Config) either.Either[ConfigError, string] {
//	    if cfg.apiKey == "" {
//	        return either.Left[string](ConfigError{"missing API key"})
//	    }
//	    return either.Right[ConfigError]("data from API")
//	}
//
//	// Context may be invalid
//	validConfig := either.Right[ConfigError](Config{apiKey: "secret"})
//	invalidConfig := either.Left[Config](ConfigError{"config not found"})
//
//	computation := readereither.FromReader[ConfigError](fetchData)
//
//	// With valid config - executes computation
//	result1 := readereither.ReadEither(validConfig)(computation)
//	// result1 = Right("data from API")
//
//	// With invalid config - short-circuits without executing
//	result2 := readereither.ReadEither(invalidConfig)(computation)
//	// result2 = Left(ConfigError{"config not found"})
//
//go:inline
func ReadEither[E1, A, E any](e Either[E1, E]) func(ReaderEither[E, E1, A]) Either[E1, A] {
	return function.Flow2(
		ET.Chain[E1, E],
		Read[E1, A](e),
	)
}

func MonadFlap[L, E, A, B any](fab ReaderEither[L, E, func(A) B], a A) ReaderEither[L, E, B] {
	return functor.MonadFlap(MonadMap[L, E, func(A) B, B], fab, a)
}

func Flap[L, E, B, A any](a A) func(ReaderEither[L, E, func(A) B]) ReaderEither[L, E, B] {
	return functor.Flap(Map[L, E, func(A) B, B], a)
}

func MonadMapLeft[C, E1, E2, A any](fa ReaderEither[C, E1, A], f func(E1) E2) ReaderEither[C, E2, A] {
	return eithert.MonadMapLeft(reader.MonadMap[C, Either[E1, A], Either[E2, A]], fa, f)
}

// MapLeft applies a mapping function to the error channel
func MapLeft[C, E1, E2, A any](f func(E1) E2) func(ReaderEither[C, E1, A]) ReaderEither[C, E2, A] {
	return eithert.MapLeft(reader.Map[C, Either[E1, A], Either[E2, A]], f)
}

// OrElse recovers from a Left (error) by providing an alternative computation with access to the reader context.
// If the ReaderEither is Right, it returns the value unchanged.
// If the ReaderEither is Left, it applies the provided function to the error value,
// which returns a new ReaderEither that replaces the original.
//
// This is useful for error recovery, fallback logic, or chaining alternative computations
// that need access to configuration or dependencies. The error type can be widened from E1 to E2.
//
// Example:
//
//	type Config struct{ fallbackValue int }
//
//	// Recover using config-dependent fallback
//	recover := readereither.OrElse(func(err error) readereither.ReaderEither[Config, error, int] {
//	    if err.Error() == "not found" {
//	        return readereither.Asks[error](func(cfg Config) either.Either[error, int] {
//	            return either.Right[error](cfg.fallbackValue)
//	        })
//	    }
//	    return readereither.Left[Config, int](err)
//	})
//	result := recover(readereither.Left[Config, int](errors.New("not found")))(Config{fallbackValue: 42}) // Right(42)
//
//go:inline
func OrElse[R, E1, E2, A any](onLeft Kleisli[R, E2, E1, A]) Kleisli[R, E2, ReaderEither[R, E1, A], A] {
	return Fold(onLeft, Of[R, E2, A])
}

// MonadChainLeft chains a computation on the left (error) side of a ReaderEither.
// If the input is a Left value, it applies the function f to transform the error and potentially
// change the error type from EA to EB. If the input is a Right value, it passes through unchanged.
//
// This is useful for error recovery or error transformation scenarios where you want to handle
// errors by performing another computation that may also fail, with access to configuration context.
//
// Note: This is functionally identical to the uncurried form of [OrElseW]. Use [ChainLeft] when
// emphasizing the monadic chaining perspective, and [OrElseW] for error recovery semantics.
//
// Parameters:
//   - fa: The input ReaderEither that may contain an error of type EA
//   - f: A Kleisli function that takes an error of type EA and returns a ReaderEither with error type EB
//
// Returns:
//   - A ReaderEither with the potentially transformed error type EB
//
// Example:
//
//	type Config struct{ fallbackValue int }
//	type ValidationError struct{ field string }
//	type SystemError struct{ code int }
//
//	// Recover from validation errors using config
//	result := MonadChainLeft(
//	    Left[Config, int](ValidationError{"username"}),
//	    func(ve ValidationError) readereither.ReaderEither[Config, SystemError, int] {
//	        if ve.field == "username" {
//	            return Asks[SystemError](func(cfg Config) either.Either[SystemError, int] {
//	                return either.Right[SystemError](cfg.fallbackValue)
//	            })
//	        }
//	        return Left[Config, int](SystemError{400})
//	    },
//	)
//
//go:inline
func MonadChainLeft[R, EA, EB, A any](fa ReaderEither[R, EA, A], f Kleisli[R, EB, EA, A]) ReaderEither[R, EB, A] {
	return func(r R) Either[EB, A] {
		return ET.Fold(
			func(ea EA) Either[EB, A] { return f(ea)(r) },
			ET.Right[EB, A],
		)(fa(r))
	}
}

// ChainLeft is the curried version of [MonadChainLeft].
// It returns a function that chains a computation on the left (error) side of a ReaderEither.
//
// This is particularly useful in functional composition pipelines where you want to handle
// errors by performing another computation that may also fail, with access to configuration context.
//
// Note: This is functionally identical to [OrElseW]. They are different names for the same operation.
// Use [ChainLeft] when emphasizing the monadic chaining perspective on the error channel,
// and [OrElseW] when emphasizing error recovery/fallback semantics.
//
// Parameters:
//   - f: A Kleisli function that takes an error of type EA and returns a ReaderEither with error type EB
//
// Returns:
//   - A function that transforms a ReaderEither with error type EA to one with error type EB
//
// Example:
//
//	type Config struct{ retryLimit int }
//
//	// Create a reusable error handler with config access
//	recoverFromError := ChainLeft(func(err string) readereither.ReaderEither[Config, int, string] {
//	    if strings.Contains(err, "retryable") {
//	        return Asks[int](func(cfg Config) either.Either[int, string] {
//	            if cfg.retryLimit > 0 {
//	                return either.Right[int]("recovered")
//	            }
//	            return either.Left[string](500)
//	        })
//	    }
//	    return Left[Config, string](404)
//	})
//
//	result := F.Pipe1(
//	    Left[Config, string]("retryable error"),
//	    recoverFromError,
//	)(Config{retryLimit: 3})
//
//go:inline
func ChainLeft[R, EA, EB, A any](f Kleisli[R, EB, EA, A]) func(ReaderEither[R, EA, A]) ReaderEither[R, EB, A] {
	return func(fa ReaderEither[R, EA, A]) ReaderEither[R, EB, A] {
		return MonadChainLeft(fa, f)
	}
}

// MonadChainFirstLeft chains a computation on the left (error) side but always returns the original error.
// If the input is a Left value, it applies the function f to the error and executes the resulting computation,
// but always returns the original Left error regardless of what f returns (Left or Right).
// If the input is a Right value, it passes through unchanged without calling f.
//
// This is useful for side effects on errors (like logging or metrics) where you want to perform an action
// when an error occurs but always propagate the original error, ensuring the error path is preserved.
//
// Parameters:
//   - ma: The input ReaderEither that may contain an error of type EA
//   - f: A function that takes an error of type EA and returns a ReaderEither (typically for side effects)
//
// Returns:
//   - A ReaderEither with the original error preserved if input was Left, or the original Right value
//
// Example:
//
//	type Config struct{ loggingEnabled bool }
//
//	// Log errors but preserve the original error
//	result := MonadChainFirstLeft(
//	    Left[Config, int]("database error"),
//	    func(err string) readereither.ReaderEither[Config, string, int] {
//	        return Asks[string](func(cfg Config) either.Either[string, int] {
//	            if cfg.loggingEnabled {
//	                log.Printf("Error: %s", err)
//	            }
//	            return either.Right[string](0)
//	        })
//	    },
//	)
//	// result will always be Left("database error")
//
//go:inline
func MonadChainFirstLeft[A, R, EA, EB, B any](ma ReaderEither[R, EA, A], f Kleisli[R, EB, EA, B]) ReaderEither[R, EA, A] {
	return eithert.MonadChainFirstLeft(
		reader.MonadChain[R, Either[EA, A], Either[EA, A]],
		reader.MonadMap[R, Either[EB, B], Either[EA, A]],
		reader.Of[R, Either[EA, A]],
		ma,
		f,
	)
}

//go:inline
func MonadTapLeft[A, R, EA, EB, B any](ma ReaderEither[R, EA, A], f Kleisli[R, EB, EA, B]) ReaderEither[R, EA, A] {
	return MonadChainFirstLeft(ma, f)
}

// ChainFirstLeft is the curried version of [MonadChainFirstLeft].
// It returns a function that chains a computation on the left (error) side while always preserving the original error.
//
// This is particularly useful for adding error handling side effects (like logging, metrics, or notifications)
// in a functional pipeline. The original error is always returned regardless of what f returns (Left or Right),
// ensuring the error path is preserved.
//
// Parameters:
//   - f: A function that takes an error of type EA and returns a ReaderEither (typically for side effects)
//
// Returns:
//   - A function that performs the side effect but always returns the original error if input was Left
//
// Example:
//
//	type Config struct{ metricsEnabled bool }
//
//	// Create a reusable error logger
//	logError := ChainFirstLeft(func(err string) readereither.ReaderEither[Config, any, int] {
//	    return Asks[any](func(cfg Config) either.Either[any, int] {
//	        if cfg.metricsEnabled {
//	            metrics.RecordError(err)
//	        }
//	        return either.Right[any](0)
//	    })
//	})
//
//	result := F.Pipe1(
//	    Left[Config, int]("validation failed"),
//	    logError, // records the error in metrics
//	)
//	// result is always Left("validation failed")
//
//go:inline
func ChainFirstLeft[A, R, EA, EB, B any](f Kleisli[R, EB, EA, B]) Operator[R, EA, A, A] {
	return eithert.ChainFirstLeft(
		reader.Chain[R, Either[EA, A], Either[EA, A]],
		reader.Map[R, Either[EB, B], Either[EA, A]],
		reader.Of[R, Either[EA, A]],
		f,
	)
}

//go:inline
func TapLeft[A, R, EA, EB, B any](f Kleisli[R, EB, EA, B]) Operator[R, EA, A, A] {
	return ChainFirstLeft[A](f)
}

// MonadFold applies one of two functions depending on the Either value.
// If Left, applies onLeft function. If Right, applies onRight function.
// Both functions return a Reader[E, B].
//
//go:inline
func MonadFold[E, L, A, B any](ma ReaderEither[E, L, A], onLeft func(L) Reader[E, B], onRight func(A) Reader[E, B]) Reader[E, B] {
	return Fold(onLeft, onRight)(ma)
}

//go:inline
func MonadAlt[R, E, A any](first ReaderEither[R, E, A], second Lazy[ReaderEither[R, E, A]]) ReaderEither[R, E, A] {
	return eithert.MonadAlt(
		reader.Of[R, Either[E, A]],
		reader.MonadChain[R, Either[E, A], Either[E, A]],

		first,
		second,
	)
}

//go:inline
func Alt[R, E, A any](second Lazy[ReaderEither[R, E, A]]) Operator[R, E, A, A] {
	return eithert.Alt(
		reader.Of[R, Either[E, A]],
		reader.Chain[R, Either[E, A], Either[E, A]],

		second,
	)
}
