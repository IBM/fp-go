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

package readeroption

import (
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/fromoption"
	"github.com/IBM/fp-go/v2/internal/fromreader"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/optiont"
	"github.com/IBM/fp-go/v2/internal/readert"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
)

// FromOption lifts an Option[A] into a ReaderOption[E, A].
// The resulting computation ignores the environment and returns the given option.
//
//go:inline
func FromOption[E, A any](e Option[A]) ReaderOption[E, A] {
	return reader.Of[E](e)
}

// Some wraps a value in a ReaderOption, representing a successful computation.
// This is equivalent to Of but more explicit about the Option semantics.
//
//go:inline
func Some[E, A any](r A) ReaderOption[E, A] {
	return optiont.Of(reader.Of[E, Option[A]], r)
}

// FromReader lifts a Reader[E, A] into a ReaderOption[E, A].
// The resulting computation always succeeds (returns Some).
//
//go:inline
func FromReader[E, A any](r Reader[E, A]) ReaderOption[E, A] {
	return SomeReader(r)
}

// SomeReader lifts a Reader[E, A] into a ReaderOption[E, A].
// The resulting computation always succeeds (returns Some).
//
//go:inline
func SomeReader[E, A any](r Reader[E, A]) ReaderOption[E, A] {
	return optiont.SomeF(reader.MonadMap[E, A, Option[A]], r)
}

// MonadMap applies a function to the value inside a ReaderOption.
// If the ReaderOption contains None, the function is not applied.
//
// Example:
//
//	ro := readeroption.Of[Config](42)
//	doubled := readeroption.MonadMap(ro, N.Mul(2))
//
//go:inline
func MonadMap[E, A, B any](fa ReaderOption[E, A], f func(A) B) ReaderOption[E, B] {
	return readert.MonadMap[ReaderOption[E, A], ReaderOption[E, B]](O.MonadMap[A, B], fa, f)
}

// Map returns a function that applies a transformation to the value inside a ReaderOption.
// This is the curried version of MonadMap, useful for composition with F.Pipe.
//
// Example:
//
//	doubled := F.Pipe1(
//	    readeroption.Of[Config](42),
//	    readeroption.Map[Config](N.Mul(2)),
//	)
//
//go:inline
func Map[E, A, B any](f func(A) B) Operator[E, A, B] {
	return readert.Map[ReaderOption[E, A], ReaderOption[E, B]](O.Map[A, B], f)
}

// MonadChain sequences two ReaderOption computations, where the second depends on the result of the first.
// If the first computation returns None, the second is not executed.
//
// Example:
//
//	findUser := func(id int) readeroption.ReaderOption[DB, User] { ... }
//	loadProfile := func(user User) readeroption.ReaderOption[DB, Profile] { ... }
//	result := readeroption.MonadChain(findUser(123), loadProfile)
//
//go:inline
func MonadChain[E, A, B any](ma ReaderOption[E, A], f Kleisli[E, A, B]) ReaderOption[E, B] {
	return readert.MonadChain(O.MonadChain[A, B], ma, f)
}

// Chain returns a function that sequences ReaderOption computations.
// This is the curried version of MonadChain, useful for composition with F.Pipe.
//
// Example:
//
//	result := F.Pipe1(
//	    findUser(123),
//	    readeroption.Chain(loadProfile),
//	)
//
//go:inline
func Chain[E, A, B any](f Kleisli[E, A, B]) Operator[E, A, B] {
	return readert.Chain[ReaderOption[E, A]](O.Chain[A, B], f)
}

// Of wraps a value in a ReaderOption, representing a successful computation.
// The resulting computation ignores the environment and returns Some(a).
//
// Example:
//
//	ro := readeroption.Of[Config](42)
//	result := ro(config) // Returns option.Some(42)
//
//go:inline
func Of[E, A any](a A) ReaderOption[E, A] {
	return readert.MonadOf[ReaderOption[E, A]](O.Of[A], a)
}

// None creates a ReaderOption representing a failed computation.
// The resulting computation ignores the environment and returns None.
//
// Example:
//
//	ro := readeroption.None[Config, int]()
//	result := ro(config) // Returns option.None[int]()
//
//go:inline
func None[E, A any]() ReaderOption[E, A] {
	return reader.Of[E](O.None[A]())
}

// MonadAp applies a function wrapped in a ReaderOption to a value wrapped in a ReaderOption.
// Both computations are executed with the same environment.
// If either computation returns None, the result is None.
//
//go:inline
func MonadAp[E, A, B any](fab ReaderOption[E, func(A) B], fa ReaderOption[E, A]) ReaderOption[E, B] {
	return readert.MonadAp[ReaderOption[E, A], ReaderOption[E, B], ReaderOption[E, func(A) B], E, A](O.MonadAp[B, A], fab, fa)
}

// Ap returns a function that applies a function wrapped in a ReaderOption to a value.
// This is the curried version of MonadAp.
//
//go:inline
func Ap[B, E, A any](fa ReaderOption[E, A]) Operator[E, func(A) B, B] {
	return readert.Ap[ReaderOption[E, A], ReaderOption[E, B], ReaderOption[E, func(A) B], E, A](O.Ap[B, A], fa)
}

// FromPredicate creates a Kleisli arrow that filters a value based on a predicate.
// If the predicate returns true, the value is wrapped in Some; otherwise, None is returned.
//
// Example:
//
//	isPositive := readeroption.FromPredicate[Config](N.MoreThan(0))
//	result := F.Pipe1(
//	    readeroption.Of[Config](42),
//	    readeroption.Chain(isPositive),
//	)
//
//go:inline
func FromPredicate[E, A any](pred Predicate[A]) Kleisli[E, A, A] {
	return fromoption.FromPredicate(FromOption[E, A], pred)
}

// Fold extracts the value from a ReaderOption by providing handlers for both cases.
// The onNone handler is called if the computation returns None.
// The onRight handler is called if the computation returns Some(a).
//
// Example:
//
//	result := readeroption.Fold(
//	    func() reader.Reader[Config, string] { return reader.Of[Config]("not found") },
//	    func(user User) reader.Reader[Config, string] { return reader.Of[Config](user.Name) },
//	)(findUser(123))
//
//go:inline
func Fold[E, A, B any](onNone Reader[E, B], onRight reader.Kleisli[E, A, B]) reader.Operator[E, Option[A], B] {
	return optiont.MatchE(reader.Chain[E, Option[A], B], function.Constant(onNone), onRight)
}

// MonadFold extracts the value from a ReaderOption by providing handlers for both cases.
// This is the non-curried version of Fold.
// The onNone handler is called if the computation returns None.
// The onRight handler is called if the computation returns Some(a).
//
// Example:
//
//	result := readeroption.MonadFold(
//	    findUser(123),
//	    reader.Of[Config]("not found"),
//	    func(user User) reader.Reader[Config, string] { return reader.Of[Config](user.Name) },
//	)
//
//go:inline
func MonadFold[E, A, B any](fa ReaderOption[E, A], onNone Reader[E, B], onRight reader.Kleisli[E, A, B]) Reader[E, B] {
	return optiont.MonadMatchE(fa, reader.MonadChain[E, Option[A], B], function.Constant(onNone), onRight)
}

// GetOrElse returns the value from a ReaderOption, or a default value if it's None.
//
// Example:
//
//	result := readeroption.GetOrElse(
//	    func() reader.Reader[Config, User] { return reader.Of[Config](defaultUser) },
//	)(findUser(123))
//
//go:inline
func GetOrElse[E, A any](onNone Reader[E, A]) reader.Operator[E, Option[A], A] {
	return optiont.GetOrElse(reader.Chain[E, Option[A], A], function.Constant(onNone), reader.Of[E, A])
}

// Ask retrieves the current environment as a ReaderOption.
// This always succeeds and returns Some(environment).
//
// Example:
//
//	getConfig := readeroption.Ask[Config]()
//	result := getConfig(myConfig) // Returns option.Some(myConfig)
//
//go:inline
func Ask[E any]() ReaderOption[E, E] {
	return fromreader.Ask(FromReader[E, E])()
}

// Asks creates a ReaderOption that applies a function to the environment.
// This always succeeds and returns Some(f(environment)).
//
// Example:
//
//	getTimeout := readeroption.Asks(func(cfg Config) int { return cfg.Timeout })
//	result := getTimeout(myConfig) // Returns option.Some(myConfig.Timeout)
//
//go:inline
func Asks[E, A any](r Reader[E, A]) ReaderOption[E, A] {
	return fromreader.Asks(FromReader[E, A])(r)
}

// MonadChainOptionK chains a ReaderOption with a function that returns an Option.
// This is useful for integrating functions that return Option directly.
//
// Example:
//
//	parseAge := func(s string) option.Option[int] { ... }
//	result := readeroption.MonadChainOptionK(
//	    readeroption.Of[Config]("25"),
//	    parseAge,
//	)
//
//go:inline
func MonadChainOptionK[E, A, B any](ma ReaderOption[E, A], f O.Kleisli[A, B]) ReaderOption[E, B] {
	return fromoption.MonadChainOptionK(
		MonadChain[E, A, B],
		FromOption[E, B],
		ma,
		f,
	)
}

// ChainOptionK returns a function that chains a ReaderOption with a function returning an Option.
// This is the curried version of MonadChainOptionK.
//
// Example:
//
//	parseAge := func(s string) option.Option[int] { ... }
//	result := F.Pipe1(
//	    readeroption.Of[Config]("25"),
//	    readeroption.ChainOptionK[Config](parseAge),
//	)
//
//go:inline
func ChainOptionK[E, A, B any](f O.Kleisli[A, B]) Operator[E, A, B] {
	return fromoption.ChainOptionK(
		Chain[E, A, B],
		FromOption[E, B],
		f,
	)
}

// Flatten removes one level of nesting from a ReaderOption.
// Converts ReaderOption[E, ReaderOption[E, A]] to ReaderOption[E, A].
//
// Example:
//
//	nested := readeroption.Of[Config](readeroption.Of[Config](42))
//	flattened := readeroption.Flatten(nested)
//
//go:inline
func Flatten[E, A any](mma ReaderOption[E, ReaderOption[E, A]]) ReaderOption[E, A] {
	return MonadChain(mma, function.Identity[ReaderOption[E, A]])
}

// Local changes the value of the local context during the execution of the action `ma` (similar to `Contravariant`'s
// `contramap`).
//
// This allows you to transform the environment before passing it to a computation.
//
// Example:
//
//	type GlobalConfig struct { DB DBConfig }
//	type DBConfig struct { Host string }
//
//	// A computation that needs DBConfig
//	query := func(cfg DBConfig) option.Option[User] { ... }
//
//	// Transform GlobalConfig to DBConfig
//	result := readeroption.Local(func(g GlobalConfig) DBConfig { return g.DB })(
//	    readeroption.Asks(query),
//	)
//
//go:inline
func Local[A, R1, R2 any](f func(R2) R1) func(ReaderOption[R1, A]) ReaderOption[R2, A] {
	return reader.Local[Option[A]](f)
}

// Read applies a context to a reader to obtain its value.
// This executes the ReaderOption computation with the given environment.
//
// Example:
//
//	ro := readeroption.Of[Config](42)
//	result := readeroption.Read[int](myConfig)(ro) // Returns option.Some(42)
//
//go:inline
func Read[A, E any](e E) func(ReaderOption[E, A]) Option[A] {
	return reader.Read[Option[A]](e)
}

// ReadOption executes a ReaderOption with an optional environment.
// If the environment is None, the result is None.
// If the environment is Some(e), the ReaderOption is executed with e.
//
// This is useful when the environment itself might not be available.
//
// Example:
//
//	ro := readeroption.Of[Config](42)
//	result1 := readeroption.ReadOption[int](option.Some(myConfig))(ro) // Returns option.Some(42)
//	result2 := readeroption.ReadOption[int](option.None[Config]())(ro) // Returns option.None[int]()
//
//go:inline
func ReadOption[A, E any](e Option[E]) func(ReaderOption[E, A]) Option[A] {
	return function.Flow2(
		O.Chain[E],
		Read[A](e),
	)
}

// MonadFlap applies a value to a function wrapped in a ReaderOption.
// This is the reverse of MonadAp.
//
//go:inline
func MonadFlap[E, A, B any](fab ReaderOption[E, func(A) B], a A) ReaderOption[E, B] {
	return functor.MonadFlap(MonadMap[E, func(A) B, B], fab, a)
}

// Flap returns a function that applies a value to a function wrapped in a ReaderOption.
// This is the curried version of MonadFlap.
//
//go:inline
func Flap[E, B, A any](a A) Operator[E, func(A) B, B] {
	return functor.Flap(Map[E, func(A) B, B], a)
}

// MonadAlt provides an alternative ReaderOption if the first one returns None.
// If fa returns Some(a), that value is returned; otherwise, the alternative computation is executed.
// This is useful for providing fallback behavior.
//
// Example:
//
//	primary := findUserInCache(123)
//	fallback := findUserInDB(123)
//	result := readeroption.MonadAlt(primary, fallback)
//
//go:inline
func MonadAlt[E, A any](first ReaderOption[E, A], second Lazy[ReaderOption[E, A]]) ReaderOption[E, A] {
	return optiont.MonadAlt(
		reader.Of[E, Option[A]],
		reader.MonadChain[E, Option[A], Option[A]],
		first,
		second,
	)
}

// Alt returns a function that provides an alternative ReaderOption if the first one returns None.
// This is the curried version of MonadAlt, useful for composition with F.Pipe.
//
// Example:
//
//	result := F.Pipe1(
//	    findUserInCache(123),
//	    readeroption.Alt(findUserInDB(123)),
//	)
//
//go:inline
func Alt[E, A any](second Lazy[ReaderOption[E, A]]) Operator[E, A, A] {
	return optiont.Alt(
		reader.Of[E, Option[A]],
		reader.Chain[E, Option[A], Option[A]],
		second,
	)
}
