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
	"github.com/IBM/fp-go/v2/internal/eithert"
	"github.com/IBM/fp-go/v2/internal/fromeither"
	"github.com/IBM/fp-go/v2/internal/fromreader"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/readert"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/result"
)

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
//	doubled := readerresult.MonadMap(rr, func(x int) int { return x * 2 })
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
//	double := readerresult.Map[Config](func(x int) int { return x * 2 })
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

// Ap is the curried version of MonadAp.
// It returns an Operator that can be used in function composition pipelines.
//
//go:inline
func Ap[B, R, A any](fa ReaderResult[R, A]) Operator[R, func(A) B, B] {
	return readert.Ap[ReaderResult[R, A], ReaderResult[R, B], ReaderResult[R, func(A) B], R, A](ET.Ap[B, error, A], fa)
}

// FromPredicate creates a Kleisli arrow that tests a predicate and returns either the input value
// or an error. If the predicate returns true, the value is returned as a success. If false,
// the onFalse function is called to generate an error.
//
// Example:
//
//	isPositive := readerresult.FromPredicate[Config](
//	    func(x int) bool { return x > 0 },
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

// MonadBiMap maps functions over both the error and success channels simultaneously.
// This transforms both the error type and the success type in a single operation.
//
// Example:
//
//	enrichErr := func(e error) error { return fmt.Errorf("enriched: %w", e) }
//	double := func(x int) int { return x * 2 }
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
//	double := func(x int) int { return x * 2 }
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
func Local[A, R2, R1 any](f func(R2) R1) func(ReaderResult[R1, A]) ReaderResult[R2, A] {
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
//	fabr := readerresult.Of[Config](func(x int) int { return x * 2 })
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
