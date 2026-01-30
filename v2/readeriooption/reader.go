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

package readeriooption

import (
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/fromoption"
	"github.com/IBM/fp-go/v2/internal/fromreader"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/optiont"
	"github.com/IBM/fp-go/v2/iooption"
	"github.com/IBM/fp-go/v2/lazy"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
)

// FromOption lifts an Option[A] into a ReaderIOOption[R, A].
// The resulting computation ignores the environment and returns the given option.
//
//go:inline
func FromOption[R, A any](t Option[A]) ReaderIOOption[R, A] {
	return readerio.Of[R](t)
}

// Some wraps a value in a ReaderIOOption, representing a successful computation.
// This is equivalent to Of but more explicit about the Option semantics.
//
//go:inline
func Some[R, A any](r A) ReaderIOOption[R, A] {
	return optiont.Of(readerio.Of[R, Option[A]], r)
}

// FromReader lifts a Reader[R, A] into a ReaderIOOption[R, A].
// The resulting computation always succeeds (returns Some).
//
//go:inline
func FromReader[R, A any](r Reader[R, A]) ReaderIOOption[R, A] {
	return SomeReader(r)
}

// SomeReader lifts a Reader[R, A] into a ReaderIOOption[R, A].
// The resulting computation always succeeds (returns Some).
//
//go:inline
func SomeReader[R, A any](r Reader[R, A]) ReaderIOOption[R, A] {
	return function.Flow2(r, iooption.Some[A])
}

// MonadMap applies a function to the value inside a ReaderIOOption.
// If the ReaderIOOption contains None, the function is not applied.
//
// Example:
//
//	ro := readeroption.Of[Config](42)
//	doubled := readeroption.MonadMap(ro, N.Mul(2))
//
//go:inline
func MonadMap[R, A, B any](fa ReaderIOOption[R, A], f func(A) B) ReaderIOOption[R, B] {
	return optiont.MonadMap(readerio.MonadMap[R, Option[A], Option[B]], fa, f)
}

// Map returns a function that applies a transformation to the value inside a ReaderIOOption.
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
func Map[R, A, B any](f func(A) B) Operator[R, A, B] {
	return optiont.Map(readerio.Map[R, Option[A], Option[B]], f)
}

// MonadChain sequences two ReaderIOOption computations, where the second depends on the result of the first.
// If the first computation returns None, the second is not executed.
//
// Example:
//
//	findUser := func(id int) readeroption.ReaderIOOption[DB, User] { ... }
//	loadProfile := func(user User) readeroption.ReaderIOOption[DB, Profile] { ... }
//	result := readeroption.MonadChain(findUser(123), loadProfile)
//
//go:inline
func MonadChain[R, A, B any](ma ReaderIOOption[R, A], f Kleisli[R, A, B]) ReaderIOOption[R, B] {
	return optiont.MonadChain(
		readerio.MonadChain[R, Option[A], Option[B]],
		readerio.Of[R, Option[B]],
		ma,
		f,
	)
}

// Chain returns a function that sequences ReaderIOOption computations.
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
func Chain[R, A, B any](f Kleisli[R, A, B]) Operator[R, A, B] {
	return optiont.Chain(
		readerio.Chain[R, Option[A], Option[B]],
		readerio.Of[R, Option[B]],
		f,
	)
}

// Of wraps a value in a ReaderIOOption, representing a successful computation.
// The resulting computation ignores the environment and returns Some(a).
//
// Example:
//
//	ro := readeroption.Of[Config](42)
//	result := ro(config) // Returns option.Some(42)
//
//go:inline
func Of[R, A any](a A) ReaderIOOption[R, A] {
	return Some[R](a)
}

// None creates a ReaderIOOption representing a failed computation.
// The resulting computation ignores the environment and returns None.
//
// Example:
//
//	ro := readeroption.None[Config, int]()
//	result := ro(config) // Returns option.None[int]()
//
//go:inline
func None[R, A any]() ReaderIOOption[R, A] {
	return readerio.Of[R](O.None[A]())
}

// MonadAp applies a function wrapped in a ReaderIOOption to a value wrapped in a ReaderIOOption.
// Both computations are executed with the same environment.
// If either computation returns None, the result is None.
//
//go:inline
func MonadAp[R, A, B any](fab ReaderIOOption[R, func(A) B], fa ReaderIOOption[R, A]) ReaderIOOption[R, B] {
	return optiont.MonadAp(
		readerio.MonadAp[Option[B], R, Option[A]],
		readerio.MonadMap[R, Option[func(A) B], func(Option[A]) Option[B]],
		fab,
		fa,
	)
}

// Ap returns a function that applies a function wrapped in a ReaderIOOption to a value.
// This is the curried version of MonadAp.
//
//go:inline
func Ap[B, R, A any](fa ReaderIOOption[R, A]) Operator[R, func(A) B, B] {
	return optiont.Ap(
		readerio.Ap[Option[B], R, Option[A]],
		readerio.Map[R, Option[func(A) B], func(Option[A]) Option[B]],
		fa,
	)
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
func FromPredicate[R, A any](pred Predicate[A]) Kleisli[R, A, A] {
	return fromoption.FromPredicate(FromOption[R, A], pred)
}

// Fold extracts the value from a ReaderIOOption by providing handlers for both cases.
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
func Fold[R, A, B any](onNone Reader[R, B], onRight reader.Kleisli[R, A, B]) reader.Operator[R, Option[A], B] {
	return optiont.MatchE(reader.Chain[R, Option[A], B], lazy.Of(onNone), onRight)
}

// MonadFold extracts the value from a ReaderIOOption by providing handlers for both cases.
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
func MonadFold[R, A, B any](fa ReaderIOOption[R, A], onNone ReaderIO[R, B], onRight readerio.Kleisli[R, A, B]) ReaderIO[R, B] {
	return optiont.MonadMatchE(fa, readerio.MonadChain[R, Option[A], B], lazy.Of(onNone), onRight)
}

// GetOrElse returns the value from a ReaderIOOption, or a default value if it's None.
//
// Example:
//
//	result := readeroption.GetOrElse(
//	    func() reader.Reader[Config, User] { return reader.Of[Config](defaultUser) },
//	)(findUser(123))
//
//go:inline
func GetOrElse[R, A any](onNone Reader[R, A]) reader.Operator[R, Option[A], A] {
	return optiont.GetOrElse(reader.Chain[R, Option[A], A], lazy.Of(onNone), reader.Of[R, A])
}

// Ask retrieves the current environment as a ReaderIOOption.
// This always succeeds and returns Some(environment).
//
// Example:
//
//	getConfig := readeroption.Ask[Config]()
//	result := getConfig(myConfig) // Returns option.Some(myConfig)
//
//go:inline
func Ask[R any]() ReaderIOOption[R, R] {
	return fromreader.Ask(FromReader[R, R])()
}

// Asks creates a ReaderIOOption that applies a function to the environment.
// This always succeeds and returns Some(f(environment)).
//
// Example:
//
//	getTimeout := readeroption.Asks(func(cfg Config) int { return cfg.Timeout })
//	result := getTimeout(myConfig) // Returns option.Some(myConfig.Timeout)
//
//go:inline
func Asks[R, A any](r Reader[R, A]) ReaderIOOption[R, A] {
	return fromreader.Asks(FromReader[R, A])(r)
}

// MonadChainOptionK chains a ReaderIOOption with a function that returns an Option.
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
func MonadChainOptionK[R, A, B any](ma ReaderIOOption[R, A], f O.Kleisli[A, B]) ReaderIOOption[R, B] {
	return fromoption.MonadChainOptionK(
		MonadChain[R, A, B],
		FromOption[R, B],
		ma,
		f,
	)
}

// ChainOptionK returns a function that chains a ReaderIOOption with a function returning an Option.
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
func ChainOptionK[R, A, B any](f O.Kleisli[A, B]) Operator[R, A, B] {
	return fromoption.ChainOptionK(
		Chain[R, A, B],
		FromOption[R, B],
		f,
	)
}

// Flatten removes one level of nesting from a ReaderIOOption.
// Converts ReaderIOOption[R, ReaderIOOption[R, A]] to ReaderIOOption[R, A].
//
// Example:
//
//	nested := readeroption.Of[Config](readeroption.Of[Config](42))
//	flattened := readeroption.Flatten(nested)
//
//go:inline
func Flatten[R, A any](mma ReaderIOOption[R, ReaderIOOption[R, A]]) ReaderIOOption[R, A] {
	return MonadChain(mma, function.Identity[ReaderIOOption[R, A]])
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
func Local[A, R1, R2 any](f func(R2) R1) func(ReaderIOOption[R1, A]) ReaderIOOption[R2, A] {
	return reader.Local[IOOption[A]](f)
}

// Read applies a context to a reader to obtain its value.
// This executes the ReaderIOOption computation with the given environment.
//
// Example:
//
//	ro := readeroption.Of[Config](42)
//	result := readeroption.Read[int](myConfig)(ro) // Returns option.Some(42)
//
//go:inline
func Read[A, R any](e R) func(ReaderIOOption[R, A]) IOOption[A] {
	return reader.Read[IOOption[A]](e)
}

// ReadOption executes a ReaderIOOption with an optional environment.
// If the environment is None, the result is None.
// If the environment is Some(e), the ReaderIOOption is executed with e.
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
// TOGGLE
// func ReadOption[A, R any](e Option[R]) func(ReaderIOOption[R, A]) IOOption[A] {
// 	return function.Flow2(
// 		optiont.Chain,
// 		Read[A](e),
// 	)
// }

// MonadFlap applies a value to a function wrapped in a ReaderIOOption.
// This is the reverse of MonadAp.
//
//go:inline
func MonadFlap[R, A, B any](fab ReaderIOOption[R, func(A) B], a A) ReaderIOOption[R, B] {
	return functor.MonadFlap(MonadMap[R, func(A) B, B], fab, a)
}

// Flap returns a function that applies a value to a function wrapped in a ReaderIOOption.
// This is the curried version of MonadFlap.
//
//go:inline
func Flap[R, B, A any](a A) Operator[R, func(A) B, B] {
	return functor.Flap(Map[R, func(A) B, B], a)
}

// MonadAlt provides an alternative ReaderIOOption if the first one returns None.
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
func MonadAlt[R, A any](first ReaderIOOption[R, A], second Lazy[ReaderIOOption[R, A]]) ReaderIOOption[R, A] {
	return optiont.MonadAlt(
		readerio.Of[R, Option[A]],
		readerio.MonadChain[R, Option[A], Option[A]],
		first,
		second,
	)
}

// Alt returns a function that provides an alternative ReaderIOOption if the first one returns None.
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
func Alt[R, A any](second Lazy[ReaderIOOption[R, A]]) Operator[R, A, A] {
	return optiont.Alt(
		readerio.Of[R, Option[A]],
		readerio.Chain[R, Option[A], Option[A]],
		second,
	)
}
