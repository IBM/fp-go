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
	ET "github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/function"
	OI "github.com/IBM/fp-go/v2/idiomatic/option"
	RRI "github.com/IBM/fp-go/v2/idiomatic/readerresult"
	RI "github.com/IBM/fp-go/v2/idiomatic/result"
	"github.com/IBM/fp-go/v2/internal/eithert"
	"github.com/IBM/fp-go/v2/internal/fromeither"
	"github.com/IBM/fp-go/v2/internal/fromreader"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/readert"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/result"
)

func fromReaderResultKleisliI[R, A, B any](f RRI.Kleisli[R, A, B]) Kleisli[R, A, B] {
	return function.Flow2(f, FromReaderResultI[R, B])
}

func fromResultKleisliI[A, B any](f RI.Kleisli[A, B]) result.Kleisli[A, B] {
	return result.Eitherize1(f)
}

func fromOptionKleisliI[A, B any](f OI.Kleisli[A, B]) option.Kleisli[A, B] {
	return option.Optionize1(f)
}

// FromEither lifts a Result[A] into a ReaderResult[R, A] that ignores the environment.
// The resulting computation will always produce the same result regardless of the environment provided.
//
// Example:
//
//	res := result.Of(42)
//	rr := readerresult.FromEither[Config](res)
//	// rr(anyConfig) will always return result.Of(42)
//
//go:inline
func FromEither[R, A any](e Result[A]) ReaderResult[R, A] {
	return reader.Of[R](e)
}

// FromResult is an alias for FromEither.
// It lifts a Result[A] into a ReaderResult[R, A] that ignores the environment.
//
//go:inline
func FromResult[R, A any](e Result[A]) ReaderResult[R, A] {
	return reader.Of[R](e)
}

// FromResultI lifts an idiomatic Go (value, error) pair into a ReaderResult[R, A] that ignores the environment.
// This is the idiomatic version of FromResult, accepting Go's native error handling pattern.
// If err is non-nil, the resulting computation will always fail with that error.
// If err is nil, the resulting computation will always succeed with the value a.
//
// Example:
//
//	value, err := strconv.Atoi("42")
//	rr := readerresult.FromResultI[Config](value, err)
//	// rr(anyConfig) will return result.Of(42) if err is nil
//
//go:inline
func FromResultI[R, A any](a A, err error) ReaderResult[R, A] {
	return reader.Of[R](result.TryCatchError(a, err))
}

// FromReaderResultI converts an idiomatic ReaderResult (that returns (A, error)) into a functional ReaderResult (that returns Result[A]).
// This bridges the gap between Go's idiomatic error handling and functional programming style.
// The idiomatic RRI.ReaderResult[R, A] is a function R -> (A, error).
// The functional ReaderResult[R, A] is a function R -> Result[A].
//
// Example:
//
//	// Idiomatic ReaderResult
//	getUserID := func(cfg Config) (int, error) {
//	    if cfg.Valid {
//	        return 42, nil
//	    }
//	    return 0, errors.New("invalid config")
//	}
//	rr := readerresult.FromReaderResultI(getUserID)
//	// rr is now a functional ReaderResult[Config, int]
//
//go:inline
func FromReaderResultI[R, A any](rr RRI.ReaderResult[R, A]) ReaderResult[R, A] {
	return func(r R) Result[A] {
		return result.TryCatchError(rr(r))
	}
}

// RightReader lifts a Reader[R, A] into a ReaderResult[R, A] that always succeeds.
// The resulting computation reads a value from the environment and wraps it in a successful Result.
//
// Example:
//
//	getPort := func(cfg Config) int { return cfg.Port }
//	rr := readerresult.RightReader[Config](getPort)
//	// rr(cfg) returns result.Of(cfg.Port)
func RightReader[R, A any](r Reader[R, A]) ReaderResult[R, A] {
	return eithert.RightF(reader.MonadMap[R, A, Result[A]], r)
}

// LeftReader lifts a Reader[R, error] into a ReaderResult[R, A] that always fails.
// The resulting computation reads an error from the environment and wraps it in a failed Result.
//
// Example:
//
//	getError := func(cfg Config) error { return cfg.InitError }
//	rr := readerresult.LeftReader[User](getError)
//	// rr(cfg) returns result.Left[User](cfg.InitError)
func LeftReader[A, R any](l Reader[R, error]) ReaderResult[R, A] {
	return eithert.LeftF(reader.MonadMap[R, error, Result[A]], l)
}

// Left creates a ReaderResult that always fails with the given error, ignoring the environment.
//
// Example:
//
//	rr := readerresult.Left[Config, User](errors.New("not found"))
//	// rr(anyConfig) always returns result.Left[User](error)
//
//go:inline
func Left[R, A any](l error) ReaderResult[R, A] {
	return eithert.Left(reader.Of[R, Result[A]], l)
}

// Right creates a ReaderResult that always succeeds with the given value, ignoring the environment.
// This is the "pure" or "return" operation for the ReaderResult monad.
//
// Example:
//
//	rr := readerresult.Right[Config](42)
//	// rr(anyConfig) always returns result.Of(42)
//
//go:inline
func Right[R, A any](r A) ReaderResult[R, A] {
	return eithert.Right(reader.Of[R, Result[A]], r)
}

// FromReader is an alias for RightReader.
// It lifts a Reader[R, A] into a ReaderResult[R, A] that always succeeds.
//
//go:inline
func FromReader[R, A any](r Reader[R, A]) ReaderResult[R, A] {
	return RightReader(r)
}

// MonadMap transforms the success value of a ReaderResult using the given function.
// If the computation fails, the error is propagated unchanged.
//
// Example:
//
//	rr := readerresult.Of[Config](5)
//	doubled := readerresult.MonadMap(rr, N.Mul(2))
//	// doubled(cfg) returns result.Of(10)
//
//go:inline
func MonadMap[R, A, B any](fa ReaderResult[R, A], f func(A) B) ReaderResult[R, B] {
	return readert.MonadMap[ReaderResult[R, A], ReaderResult[R, B]](ET.MonadMap[error, A, B], fa, f)
}

// Map is the curried version of MonadMap.
// It returns an Operator that can be used in function composition pipelines.
//
// Example:
//
//	double := readerresult.Map[Config](N.Mul(2))
//	result := F.Pipe1(readerresult.Of[Config](5), double)
//
//go:inline
func Map[R, A, B any](f func(A) B) Operator[R, A, B] {
	return readert.Map[ReaderResult[R, A], ReaderResult[R, B]](ET.Map[error, A, B], f)
}

// MonadChain sequences two ReaderResult computations, where the second depends on the result of the first.
// This is also known as "flatMap" or "bind". If the first computation fails, the second is not executed.
//
// Example:
//
//	getUser := func(id int) readerresult.ReaderResult[DB, User] { ... }
//	getPosts := func(user User) readerresult.ReaderResult[DB, []Post] { ... }
//	userPosts := readerresult.MonadChain(getUser(42), getPosts)
//
//go:inline
func MonadChain[R, A, B any](ma ReaderResult[R, A], f Kleisli[R, A, B]) ReaderResult[R, B] {
	return readert.MonadChain(ET.MonadChain[error, A, B], ma, f)
}

//go:inline
func MonadChainReaderK[R, A, B any](ma ReaderResult[R, A], f reader.Kleisli[R, A, B]) ReaderResult[R, B] {
	return readert.MonadChain(ET.MonadChain[error, A, B], ma, function.Flow2(f, FromReader[R, B]))
}

// Chain is the curried version of MonadChain.
// It returns an Operator that can be used in function composition pipelines.
//
// Example:
//
//	getPosts := func(user User) readerresult.ReaderResult[DB, []Post] { ... }
//	result := F.Pipe1(getUser(42), readerresult.Chain[DB](getPosts))
//
//go:inline
func Chain[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, B] {
	return readert.Chain[ReaderResult[R, A]](ET.Chain[error, A, B], f)
}

//go:inline
func ChainReaderK[R, A, B any](f reader.Kleisli[R, A, B]) Operator[R, A, B] {
	return readert.Chain[ReaderResult[R, A]](ET.Chain[error, A, B], function.Flow2(f, FromReader[R, B]))
}

// MonadChainI sequences two ReaderResult computations, where the second is an idiomatic Kleisli arrow.
// This is the idiomatic version of MonadChain, allowing you to chain with functions that return (B, error).
// The idiomatic Kleisli arrow RRI.Kleisli[R, A, B] is a function A -> R -> (B, error).
// If the first computation fails, the second is not executed.
//
// Example:
//
//	getUserID := readerresult.Of[DB](42)
//	// Idiomatic function that returns (User, error)
//	fetchUser := func(id int) func(db DB) (User, error) {
//	    return func(db DB) (User, error) {
//	        return db.GetUser(id)  // returns (User, error)
//	    }
//	}
//	result := readerresult.MonadChainI(getUserID, fetchUser)
//
//go:inline
func MonadChainI[R, A, B any](ma ReaderResult[R, A], f RRI.Kleisli[R, A, B]) ReaderResult[R, B] {
	return MonadChain(ma, fromReaderResultKleisliI(f))
}

// ChainI is the curried version of MonadChainI.
// It allows chaining with idiomatic Kleisli arrows that return (B, error).
//
// Example:
//
//	// Idiomatic function that returns (User, error)
//	fetchUser := func(id int) func(db DB) (User, error) {
//	    return func(db DB) (User, error) {
//	        return db.GetUser(id)  // returns (User, error)
//	    }
//	}
//	result := F.Pipe1(getUserIDRR, readerresult.ChainI[DB](fetchUser))
//
//go:inline
func ChainI[R, A, B any](f RRI.Kleisli[R, A, B]) Operator[R, A, B] {
	return Chain(fromReaderResultKleisliI(f))
}

// Of creates a ReaderResult that always succeeds with the given value.
// This is an alias for Right and is the "pure" or "return" operation for the ReaderResult monad.
//
// Example:
//
//	rr := readerresult.Of[Config](42)
//	// rr(anyConfig) always returns result.Of(42)
//
//go:inline
func Of[R, A any](a A) ReaderResult[R, A] {
	return readert.MonadOf[ReaderResult[R, A]](ET.Of[error, A], a)
}

// MonadAp applies a function wrapped in a ReaderResult to a value wrapped in a ReaderResult.
// Both computations share the same environment. This is useful for combining independent
// computations that don't depend on each other's results.
//
// Example:
//
//	add := func(x int) func(int) int { return func(y int) int { return x + y } }
//	fabr := readerresult.Of[Config](add(5))
//	fa := readerresult.Of[Config](3)
//	result := readerresult.MonadAp(fabr, fa)  // Returns Of(8)
//
//go:inline
func MonadAp[B, R, A any](fab ReaderResult[R, func(A) B], fa ReaderResult[R, A]) ReaderResult[R, B] {
	return readert.MonadAp[ReaderResult[R, A], ReaderResult[R, B], ReaderResult[R, func(A) B], R, A](ET.MonadAp[B, error, A], fab, fa)
}

//go:inline
func MonadApReader[B, R, A any](fab ReaderResult[R, func(A) B], fa Reader[R, A]) ReaderResult[R, B] {
	return MonadAp(fab, FromReader(fa))
}

// Ap is the curried version of MonadAp.
// It returns an Operator that can be used in function composition pipelines.
//
//go:inline
func Ap[B, R, A any](fa ReaderResult[R, A]) Operator[R, func(A) B, B] {
	return readert.Ap[ReaderResult[R, A], ReaderResult[R, B], ReaderResult[R, func(A) B], R, A](ET.Ap[B, error, A], fa)
}

//go:inline
func ApReader[B, R, A any](fa Reader[R, A]) Operator[R, func(A) B, B] {
	return Ap[B](FromReader(fa))
}

// MonadApResult applies a function wrapped in a ReaderResult to a value wrapped in a plain Result.
// The Result value is independent of the environment, while the function may depend on it.
// This is useful when you have a pre-computed Result value that you want to apply a context-dependent function to.
//
// Example:
//
//	add := func(x int) func(int) int { return func(y int) int { return x + y } }
//	fabr := readerresult.Of[Config](add(5))
//	fa := result.Of(3)  // Pre-computed Result, independent of environment
//	result := readerresult.MonadApResult(fabr, fa)  // Returns Of(8)
//
//go:inline
func MonadApResult[B, R, A any](fab ReaderResult[R, func(A) B], fa result.Result[A]) ReaderResult[R, B] {
	return readert.MonadAp[ReaderResult[R, A], ReaderResult[R, B], ReaderResult[R, func(A) B], R, A](ET.MonadAp[B, error, A], fab, FromResult[R](fa))
}

// ApResult is the curried version of MonadApResult.
// It returns an Operator that applies a pre-computed Result value to a function in a ReaderResult context.
// This is useful in function composition pipelines when you have a static Result value.
//
// Example:
//
//	fa := result.Of(10)
//	result := F.Pipe1(
//	    readerresult.Of[Config](utils.Double),
//	    readerresult.ApResult[int, Config](fa),
//	)
//	// result(cfg) returns result.Of(20)
//
//go:inline
func ApResult[B, R, A any](fa Result[A]) Operator[R, func(A) B, B] {
	return readert.Ap[ReaderResult[R, A], ReaderResult[R, B], ReaderResult[R, func(A) B], R, A](ET.Ap[B, error, A], FromResult[R](fa))
}

// ApResultI is the curried idiomatic version of ApResult.
// It accepts a (value, error) pair directly and applies it to a function in a ReaderResult context.
// This bridges Go's idiomatic error handling with the functional ApResult operation.
//
// Example:
//
//	value, err := strconv.Atoi("10")  // Returns (10, nil)
//	result := F.Pipe1(
//	    readerresult.Of[Config](utils.Double),
//	    readerresult.ApResultI[int, Config](value, err),
//	)
//	// result(cfg) returns result.Of(20)
//
//go:inline
func ApResultI[B, R, A any](a A, err error) Operator[R, func(A) B, B] {
	return Ap[B](FromResultI[R](a, err))
}

// MonadApI applies a function wrapped in a ReaderResult to a value wrapped in an idiomatic ReaderResult.
// This is the idiomatic version of MonadAp, where the second parameter returns (A, error) instead of Result[A].
// Both computations share the same environment.
//
// Example:
//
//	add := func(x int) func(int) int { return func(y int) int { return x + y } }
//	fabr := readerresult.Of[Config](add(5))
//	// Idiomatic computation returning (int, error)
//	fa := func(cfg Config) (int, error) { return cfg.Port, nil }
//	result := readerresult.MonadApI(fabr, fa)
//
//go:inline
func MonadApI[B, R, A any](fab ReaderResult[R, func(A) B], fa RRI.ReaderResult[R, A]) ReaderResult[R, B] {
	return MonadAp(fab, FromReaderResultI(fa))
}

// ApI is the curried version of MonadApI.
// It allows applying to idiomatic ReaderResult values that return (A, error).
//
// Example:
//
//	// Idiomatic computation returning (int, error)
//	fa := func(cfg Config) (int, error) { return cfg.Port, nil }
//	result := F.Pipe1(fabr, readerresult.ApI[int, Config](fa))
//
//go:inline
func ApI[B, R, A any](fa RRI.ReaderResult[R, A]) Operator[R, func(A) B, B] {
	return Ap[B](FromReaderResultI(fa))
}

// FromPredicate creates a Kleisli arrow that tests a predicate and returns either the input value
// or an error. If the predicate returns true, the value is returned as a success. If false,
// the onFalse function is called to generate an error.
//
// Example:
//
//	isPositive := readerresult.FromPredicate[Config](
//	    N.MoreThan(0),
//	    func(x int) error { return fmt.Errorf("%d is not positive", x) },
//	)
//	result := isPositive(5)  // Returns ReaderResult that succeeds with 5
//	result := isPositive(-1) // Returns ReaderResult that fails with error
//
//go:inline
func FromPredicate[R, A any](pred func(A) bool, onFalse func(A) error) Kleisli[R, A, A] {
	return fromeither.FromPredicate(FromEither[R, A], pred, onFalse)
}

// Fold handles both success and failure cases by providing functions for each.
// The result is always a Reader[R, B] (without the error channel).
// This is useful for converting a ReaderResult into a plain Reader by handling the error case.
//
// Example:
//
//	handleError := func(err error) reader.Reader[Config, string] {
//	    return func(cfg Config) string { return "Error: " + err.Error() }
//	}
//	handleSuccess := func(user User) reader.Reader[Config, string] {
//	    return func(cfg Config) string { return user.Name }
//	}
//	result := readerresult.Fold[Config, User, string](handleError, handleSuccess)(getUserRR)
//
//go:inline
func Fold[R, A, B any](onLeft reader.Kleisli[R, error, B], onRight reader.Kleisli[R, A, B]) func(ReaderResult[R, A]) Reader[R, B] {
	return eithert.MatchE(reader.MonadChain[R, Result[A], B], onLeft, onRight)
}

// GetOrElse extracts the success value or computes a default value from the error.
// The result is a Reader[R, A] that always succeeds.
//
// Example:
//
//	defaultUser := func(err error) reader.Reader[Config, User] {
//	    return func(cfg Config) User { return User{Name: "Guest"} }
//	}
//	result := readerresult.GetOrElse[Config](defaultUser)(getUserRR)
//
//go:inline
func GetOrElse[R, A any](onLeft reader.Kleisli[R, error, A]) func(ReaderResult[R, A]) Reader[R, A] {
	return eithert.GetOrElse(reader.MonadChain[R, Result[A], A], reader.Of[R, A], onLeft)
}

// OrElse provides an alternative ReaderResult computation if the first one fails.
// This is useful for fallback logic or retry scenarios.
//
// Example:
//
//	getPrimaryUser := func(id int) readerresult.ReaderResult[DB, User] { ... }
//	getBackupUser := func(err error) readerresult.ReaderResult[DB, User] {
//	    return readerresult.Of[DB](User{Name: "Guest"})
//	}
//	result := F.Pipe1(getPrimaryUser(42), readerresult.OrElse[DB](getBackupUser))
//
//go:inline
func OrElse[R, A any](onLeft Kleisli[R, error, A]) Operator[R, A, A] {
	return eithert.OrElse(reader.MonadChain[R, Result[A], Result[A]], reader.Of[R, Result[A]], onLeft)
}

// OrElseI provides an alternative ReaderResult computation using an idiomatic Kleisli arrow if the first one fails.
// This is the idiomatic version of OrElse, where the fallback function returns (A, error) instead of Result[A].
// This is useful for fallback logic or retry scenarios with idiomatic Go functions.
//
// Example:
//
//	getPrimaryUser := func(id int) readerresult.ReaderResult[DB, User] { ... }
//	// Idiomatic fallback returning (User, error)
//	getBackupUser := func(err error) func(db DB) (User, error) {
//	    return func(db DB) (User, error) {
//	        return User{Name: "Guest"}, nil
//	    }
//	}
//	result := F.Pipe1(getPrimaryUser(42), readerresult.OrElseI[DB](getBackupUser))
//
//go:inline
func OrElseI[R, A any](onLeft RRI.Kleisli[R, error, A]) Operator[R, A, A] {
	return OrElse(fromReaderResultKleisliI(onLeft))
}

// OrLeft transforms the error value if the computation fails, leaving successful values unchanged.
// This is useful for error mapping or enriching error information.
//
// Example:
//
//	enrichError := func(err error) reader.Reader[Config, error] {
//	    return func(cfg Config) error {
//	        return fmt.Errorf("DB error on %s: %w", cfg.DBHost, err)
//	    }
//	}
//	result := F.Pipe1(getUserRR, readerresult.OrLeft[Config](enrichError))
//
//go:inline
func OrLeft[R, A any](onLeft reader.Kleisli[R, error, error]) Operator[R, A, A] {
	return eithert.OrLeft(
		reader.MonadChain[R, Result[A], Result[A]],
		reader.MonadMap[R, error, Result[A]],
		reader.Of[R, Result[A]],
		onLeft,
	)
}

// Ask retrieves the current environment as a successful ReaderResult.
// This is useful for accessing configuration or context within a ReaderResult computation.
//
// Example:
//
//	result := F.Pipe1(
//	    readerresult.Ask[Config](),
//	    readerresult.Map[Config](func(cfg Config) int { return cfg.Port }),
//	)
//	// result(cfg) returns result.Of(cfg.Port)
//
//go:inline
func Ask[R any]() ReaderResult[R, R] {
	return fromreader.Ask(FromReader[R, R])()
}

// Asks retrieves a value from the environment using the provided Reader function.
// This lifts a Reader computation into a ReaderResult that always succeeds.
//
// Example:
//
//	getPort := func(cfg Config) int { return cfg.Port }
//	result := readerresult.Asks[Config](getPort)
//	// result(cfg) returns result.Of(cfg.Port)
//
//go:inline
func Asks[R, A any](r Reader[R, A]) ReaderResult[R, A] {
	return fromreader.Asks(FromReader[R, A])(r)
}

// MonadChainEitherK chains a ReaderResult with a function that returns a plain Result.
// This is useful for integrating functions that don't need environment access.
//
// Example:
//
//	parseUser := func(data string) result.Result[User] { ... }
//	result := readerresult.MonadChainEitherK(getUserDataRR, parseUser)
//
//go:inline
func MonadChainEitherK[R, A, B any](ma ReaderResult[R, A], f result.Kleisli[A, B]) ReaderResult[R, B] {
	return fromeither.MonadChainEitherK(
		MonadChain[R, A, B],
		FromEither[R, B],
		ma,
		f,
	)
}

// MonadChainEitherIK chains a ReaderResult with an idiomatic function that returns (B, error).
// This is the idiomatic version of MonadChainEitherK, accepting functions in Go's native error handling pattern.
// The function f doesn't need environment access and returns a (value, error) pair.
//
// Example:
//
//	getUserDataRR := readerresult.Of[DB]("user_data")
//	// Idiomatic parser returning (User, error)
//	parseUser := func(data string) (User, error) {
//	    return json.Unmarshal([]byte(data), &User{})
//	}
//	result := readerresult.MonadChainEitherIK(getUserDataRR, parseUser)
//
//go:inline
func MonadChainEitherIK[R, A, B any](ma ReaderResult[R, A], f RI.Kleisli[A, B]) ReaderResult[R, B] {
	return MonadChainEitherK(ma, fromResultKleisliI(f))
}

// ChainEitherK is the curried version of MonadChainEitherK.
// It lifts a Result-returning function into a ReaderResult operator.
//
// Example:
//
//	parseUser := func(data string) result.Result[User] { ... }
//	result := F.Pipe1(getUserDataRR, readerresult.ChainEitherK[Config](parseUser))
//
//go:inline
func ChainEitherK[R, A, B any](f result.Kleisli[A, B]) Operator[R, A, B] {
	return fromeither.ChainEitherK(
		Chain[R, A, B],
		FromEither[R, B],
		f,
	)
}

// ChainEitherIK is the curried version of MonadChainEitherIK.
// It lifts an idiomatic function returning (B, error) into a ReaderResult operator.
//
// Example:
//
//	// Idiomatic parser returning (User, error)
//	parseUser := func(data string) (User, error) {
//	    return json.Unmarshal([]byte(data), &User{})
//	}
//	result := F.Pipe1(getUserDataRR, readerresult.ChainEitherIK[DB](parseUser))
//
//go:inline
func ChainEitherIK[R, A, B any](f RI.Kleisli[A, B]) Operator[R, A, B] {
	return ChainEitherK[R](fromResultKleisliI(f))
}

// ChainOptionK chains with a function that returns an Option, converting None to an error.
// This is useful for integrating functions that return optional values.
//
// Example:
//
//	findUser := func(id int) option.Option[User] { ... }
//	notFound := func() error { return errors.New("user not found") }
//	chain := readerresult.ChainOptionK[Config, int, User](notFound)
//	result := F.Pipe1(readerresult.Of[Config](42), chain(findUser))
//
//go:inline
func ChainOptionK[R, A, B any](onNone Lazy[error]) func(option.Kleisli[A, B]) Operator[R, A, B] {
	return fromeither.ChainOptionK(MonadChain[R, A, B], FromEither[R, B], onNone)
}

// ChainOptionIK chains with an idiomatic function that returns (Option[B], error), converting None to an error.
// This is the idiomatic version of ChainOptionK, accepting functions in Go's native error handling pattern.
// The onNone function is called when the Option is None to generate an error.
//
// Example:
//
//	// Idiomatic function returning (Option[User], error)
//	findUser := func(id int) (option.Option[User], error) {
//	    user, err := db.Query(id)
//	    if err != nil {
//	        return option.None[User](), err
//	    }
//	    if user == nil {
//	        return option.None[User](), nil
//	    }
//	    return option.Some(*user), nil
//	}
//	notFound := func() error { return errors.New("user not found") }
//	chain := readerresult.ChainOptionIK[Config, int, User](notFound)
//	result := F.Pipe1(readerresult.Of[Config](42), chain(findUser))
//
//go:inline
func ChainOptionIK[R, A, B any](onNone Lazy[error]) func(OI.Kleisli[A, B]) Operator[R, A, B] {
	return function.Flow2(fromOptionKleisliI[A, B], ChainOptionK[R, A, B](onNone))
}

// Flatten removes one level of nesting from a nested ReaderResult.
// This converts ReaderResult[R, ReaderResult[R, A]] into ReaderResult[R, A].
//
// Example:
//
//	nested := readerresult.Of[Config](readerresult.Of[Config](42))
//	flat := readerresult.Flatten(nested)
//	// flat(cfg) returns result.Of(42)
//
//go:inline
func Flatten[R, A any](mma ReaderResult[R, ReaderResult[R, A]]) ReaderResult[R, A] {
	return MonadChain(mma, function.Identity[ReaderResult[R, A]])
}

// FlattenI removes one level of nesting from a ReaderResult containing an idiomatic ReaderResult.
// This converts ReaderResult[R, RRI.ReaderResult[R, A]] into ReaderResult[R, A].
// The inner computation returns (A, error) in Go's idiomatic style.
//
// Example:
//
//	// Nested computation where inner returns (A, error)
//	nested := readerresult.Of[Config](func(cfg Config) (int, error) {
//	    return 42, nil
//	})
//	flat := readerresult.FlattenI(nested)
//	// flat(cfg) returns result.Of(42)
//
//go:inline
func FlattenI[R, A any](mma ReaderResult[R, RRI.ReaderResult[R, A]]) ReaderResult[R, A] {
	return MonadChain(mma, FromReaderResultI[R, A])
}

// MonadBiMap maps functions over both the error and success channels simultaneously.
// This transforms both the error type and the success type in a single operation.
//
// Example:
//
//	enrichErr := func(e error) error { return fmt.Errorf("enriched: %w", e) }
//	double := N.Mul(2)
//	result := readerresult.MonadBiMap(rr, enrichErr, double)
//
//go:inline
func MonadBiMap[R, A, B any](fa ReaderResult[R, A], f Endomorphism[error], g func(A) B) ReaderResult[R, B] {
	return eithert.MonadBiMap(reader.MonadMap[R, Result[A], Result[B]], fa, f, g)
}

// BiMap is the curried version of MonadBiMap.
// It maps a pair of functions over the error and success channels.
//
// Example:
//
//	enrichErr := func(e error) error { return fmt.Errorf("enriched: %w", e) }
//	double := N.Mul(2)
//	result := F.Pipe1(rr, readerresult.BiMap[Config](enrichErr, double))
//
//go:inline
func BiMap[R, A, B any](f Endomorphism[error], g func(A) B) Operator[R, A, B] {
	return eithert.BiMap(reader.Map[R, Result[A], Result[B]], f, g)
}

// Local changes the environment type during execution of a ReaderResult.
// This is similar to Contravariant's contramap and allows adapting computations
// to work with different environment types.
//
// Example:
//
//	// Convert DB environment to Config environment
//	toConfig := func(db DB) Config { return db.Config }
//	rr := readerresult.Of[Config](42)
//	adapted := readerresult.Local[int](toConfig)(rr)
//	// adapted now accepts DB instead of Config
//
//go:inline
func Local[A, R1, R2 any](f func(R2) R1) func(ReaderResult[R1, A]) ReaderResult[R2, A] {
	return reader.Local[Result[A]](f)
}

// Read applies an environment value to a ReaderResult to execute it and obtain the Result.
// This is the primary way to "run" a ReaderResult computation.
//
// Example:
//
//	rr := readerresult.Asks(func(cfg Config) int { return cfg.Port })
//	run := readerresult.Read[int](myConfig)
//	res := run(rr)  // Returns result.Result[int]
//
//go:inline
func Read[A, R any](r R) func(ReaderResult[R, A]) Result[A] {
	return reader.Read[Result[A]](r)
}

// MonadFlap applies a wrapped function to a concrete value (reverse of Ap).
// This is useful when you have a function in a context and a plain value.
//
// Example:
//
//	fabr := readerresult.Of[Config](N.Mul(2))
//	result := readerresult.MonadFlap(fabr, 5)  // Returns Of(10)
//
//go:inline
func MonadFlap[R, A, B any](fab ReaderResult[R, func(A) B], a A) ReaderResult[R, B] {
	return functor.MonadFlap(MonadMap[R, func(A) B, B], fab, a)
}

// Flap is the curried version of MonadFlap.
//
//go:inline
func Flap[R, B, A any](a A) Operator[R, func(A) B, B] {
	return functor.Flap(Map[R, func(A) B, B], a)
}

// MonadMapLeft transforms the error value without affecting successful results.
// This is useful for error enrichment, wrapping, or transformation.
//
// Example:
//
//	enrichErr := func(e error) error { return fmt.Errorf("DB error: %w", e) }
//	result := readerresult.MonadMapLeft(getUserRR, enrichErr)
//
//go:inline
func MonadMapLeft[R, A any](fa ReaderResult[R, A], f Endomorphism[error]) ReaderResult[R, A] {
	return eithert.MonadMapLeft(reader.MonadMap[R, Result[A], Result[A]], fa, f)
}

// MapLeft is the curried version of MonadMapLeft.
// It applies a mapping function to the error channel only.
//
// Example:
//
//	enrichErr := func(e error) error { return fmt.Errorf("DB error: %w", e) }
//	result := F.Pipe1(getUserRR, readerresult.MapLeft[Config](enrichErr))
//
//go:inline
func MapLeft[R, A any](f Endomorphism[error]) Operator[R, A, A] {
	return eithert.MapLeft(reader.Map[R, Result[A], Result[A]], f)
}

// MonadAlt tries the first computation, and if it fails, tries the second.
// This implements the Alternative pattern for error recovery.
//
//go:inline
func MonadAlt[R, A any](first ReaderResult[R, A], second Lazy[ReaderResult[R, A]]) ReaderResult[R, A] {
	return eithert.MonadAlt(
		reader.Of[R, Result[A]],
		reader.MonadChain[R, Result[A], Result[A]],

		first,
		second,
	)
}

// MonadAltI tries the first computation, and if it fails, tries the second idiomatic computation.
// This is the idiomatic version of MonadAlt, where the alternative computation returns (A, error).
// The second computation is lazy-evaluated and only executed if the first fails.
//
// Example:
//
//	primary := readerresult.Left[Config, int](errors.New("primary failed"))
//	// Idiomatic alternative returning (int, error)
//	alternative := func() func(cfg Config) (int, error) {
//	    return func(cfg Config) (int, error) {
//	        return 42, nil
//	    }
//	}
//	result := readerresult.MonadAltI(primary, alternative)
//
//go:inline
func MonadAltI[R, A any](first ReaderResult[R, A], second Lazy[RRI.ReaderResult[R, A]]) ReaderResult[R, A] {
	return MonadAlt(first, function.Pipe1(second, lazy.Map(FromReaderResultI[R, A])))
}

// Alt tries the first computation, and if it fails, tries the second.
// This implements the Alternative pattern for error recovery.
//
//go:inline
func Alt[R, A any](second Lazy[ReaderResult[R, A]]) Operator[R, A, A] {
	return eithert.Alt(
		reader.Of[R, Result[A]],
		reader.Chain[R, Result[A], Result[A]],

		second,
	)
}

// AltI is the curried version of MonadAltI.
// It tries the first computation, and if it fails, tries the idiomatic alternative that returns (A, error).
//
// Example:
//
//	// Idiomatic alternative returning (int, error)
//	alternative := func() func(cfg Config) (int, error) {
//	    return func(cfg Config) (int, error) {
//	        return 42, nil
//	    }
//	}
//	result := F.Pipe1(primary, readerresult.AltI[Config](alternative))
//
//go:inline
func AltI[R, A any](second Lazy[RRI.ReaderResult[R, A]]) Operator[R, A, A] {
	return Alt(function.Pipe1(second, lazy.Map(FromReaderResultI[R, A])))
}
