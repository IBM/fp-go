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

package generic

import (
	ET "github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/function"
	A "github.com/IBM/fp-go/v2/internal/apply"
	C "github.com/IBM/fp-go/v2/internal/chain"
	FE "github.com/IBM/fp-go/v2/internal/fromeither"
	FR "github.com/IBM/fp-go/v2/internal/fromreader"
	F "github.com/IBM/fp-go/v2/internal/functor"
)

// Do creates an empty context of type [S] to be used with the [Bind] operation.
// This is the starting point for do-notation style composition.
//
// Example:
//
//	type State struct {
//	    Config Config
//	    User   User
//	}
//	type Env struct {
//	    ConfigService ConfigService
//	    UserService   UserService
//	}
//	result := generic.Do[ReaderEither[Env, error, State], Env, error, State](State{})
//
//go:inline
func Do[GS ~func(R) ET.Either[E, S], R, E, S any](
	empty S,
) GS {
	return Of[GS](empty)
}

// Bind attaches the result of a computation to a context [S1] to produce a context [S2].
// This enables sequential composition where each step can depend on the results of previous steps
// and access the shared environment.
//
// The setter function takes the result of the computation and returns a function that
// updates the context from S1 to S2.
//
// Example:
//
//	type State struct {
//	    Config Config
//	    User   User
//	}
//	type Env struct {
//	    ConfigService ConfigService
//	    UserService   UserService
//	}
//
//	result := F.Pipe2(
//	    generic.Do[ReaderEither[Env, error, State], Env, error, State](State{}),
//	    generic.Bind[ReaderEither[Env, error, State], ReaderEither[Env, error, State], ReaderEither[Env, error, Config], Env, error, State, State, Config](
//	        func(cfg Config) func(State) State {
//	            return func(s State) State { s.Config = cfg; return s }
//	        },
//	        func(s State) ReaderEither[Env, error, Config] {
//	            return func(env Env) either.Either[error, Config] {
//	                return env.ConfigService.Load()
//	            }
//	        },
//	    ),
//	    generic.Bind[ReaderEither[Env, error, State], ReaderEither[Env, error, State], ReaderEither[Env, error, User], Env, error, State, State, User](
//	        func(user User) func(State) State {
//	            return func(s State) State { s.User = user; return s }
//	        },
//	        func(s State) ReaderEither[Env, error, User] {
//	            // This can access s.Config from the previous step
//	            return func(env Env) either.Either[error, User] {
//	                return env.UserService.GetUserForConfig(s.Config)
//	            }
//	        },
//	    ),
//	)
//
//go:inline
func Bind[
	GS1 ~func(R) ET.Either[E, S1],
	GS2 ~func(R) ET.Either[E, S2],
	GT ~func(R) ET.Either[E, T], R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) GT,
) func(GS1) GS2 {
	return C.Bind(
		Chain[GS1, GS2, E, R, S1, S2],
		Map[GT, GS2, E, R, T, S2],
		setter,
		f,
	)
}

//go:inline
func BindReaderK[
	GS1 ~func(R) ET.Either[E, S1],
	GS2 ~func(R) ET.Either[E, S2],
	GRT ~func(R) T,
	R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) GRT,
) func(GS1) GS2 {
	return FR.BindReaderK(
		Chain[GS1, GS2, E, R, S1, S2],
		Map[func(R) ET.Either[E, T], GS2, E, R, T, S2],
		FromReader[GRT, func(R) ET.Either[E, T]],
		setter,
		f,
	)
}

//go:inline
func BindEitherK[
	GS1 ~func(R) ET.Either[E, S1],
	GS2 ~func(R) ET.Either[E, S2],
	R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) ET.Either[E, T],
) func(GS1) GS2 {
	return FE.BindEitherK(
		Chain[GS1, GS2, E, R, S1, S2],
		Map[func(R) ET.Either[E, T], GS2, E, R, T, S2],
		FromEither[func(R) ET.Either[E, T]],
		setter,
		f,
	)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
func Let[GS1 ~func(R) ET.Either[E, S1], GS2 ~func(R) ET.Either[E, S2], R, E, S1, S2, T any](
	key func(T) func(S1) S2,
	f func(S1) T,
) func(GS1) GS2 {
	return F.Let(
		Map[GS1, GS2, E, R, S1, S2],
		key,
		f,
	)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
func LetTo[GS1 ~func(R) ET.Either[E, S1], GS2 ~func(R) ET.Either[E, S2], R, E, S1, S2, B any](
	key func(B) func(S1) S2,
	b B,
) func(GS1) GS2 {
	return F.LetTo(
		Map[GS1, GS2, E, R, S1, S2],
		key,
		b,
	)
}

// BindTo initializes a new state [S1] from a value [T]
func BindTo[GS1 ~func(R) ET.Either[E, S1], GT ~func(R) ET.Either[E, T], R, E, S1, T any](
	setter func(T) S1,
) func(GT) GS1 {
	return C.BindTo(
		Map[GT, GS1, E, R, T, S1],
		setter,
	)
}

//go:inline
func BindToReader[
	GS1 ~func(R) ET.Either[E, S1],
	GT ~func(R) T,
	R, E, S1, T any](
	setter func(T) S1,
) func(GT) GS1 {
	return function.Flow2(
		FromReader[GT, func(R) ET.Either[E, T]],
		BindTo[GS1, func(R) ET.Either[E, T]](setter),
	)
}

//go:inline
func BindToEither[
	GS1 ~func(R) ET.Either[E, S1],
	R, E, S1, T any](
	setter func(T) S1,
) func(ET.Either[E, T]) GS1 {
	return function.Flow2(
		FromEither[func(R) ET.Either[E, T]],
		BindTo[GS1, func(R) ET.Either[E, T]](setter),
	)
}

// ApS attaches a value to a context [S1] to produce a context [S2] by considering
// the context and the value concurrently (using Applicative rather than Monad).
// This allows independent computations to be combined without one depending on the result of the other.
//
// Unlike Bind, which sequences operations, ApS can be used when operations are independent
// and can conceptually run in parallel.
//
// Example:
//
//	type State struct {
//	    Config Config
//	    User   User
//	}
//	type Env struct {
//	    ConfigService ConfigService
//	    UserService   UserService
//	}
//
//	// These operations are independent and can be combined with ApS
//	getConfig := func(env Env) either.Either[error, Config] {
//	    return env.ConfigService.Load()
//	}
//	getUser := func(env Env) either.Either[error, User] {
//	    return env.UserService.GetCurrent()
//	}
//
//	result := F.Pipe2(
//	    generic.Do[ReaderEither[Env, error, State], Env, error, State](State{}),
//	    generic.ApS[...](
//	        func(cfg Config) func(State) State {
//	            return func(s State) State { s.Config = cfg; return s }
//	        },
//	        getConfig,
//	    ),
//	    generic.ApS[...](
//	        func(user User) func(State) State {
//	            return func(s State) State { s.User = user; return s }
//	        },
//	        getUser,
//	    ),
//	)
func ApS[GS1 ~func(R) ET.Either[E, S1], GS2 ~func(R) ET.Either[E, S2], GT ~func(R) ET.Either[E, T], R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa GT,
) func(GS1) GS2 {
	return A.ApS(
		Ap[GT, GS2, func(R) ET.Either[E, func(T) S2], E, R, T, S2],
		Map[GS1, func(R) ET.Either[E, func(T) S2], E, R, S1, func(T) S2],
		setter,
		fa,
	)
}

//go:inline
func ApReaderS[
	GS1 ~func(R) ET.Either[E, S1],
	GS2 ~func(R) ET.Either[E, S2],
	GT ~func(R) T,
	R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa GT,
) func(GS1) GS2 {
	return ApS[GS1, GS2](
		setter,
		FromReader[GT, func(R) ET.Either[E, T]](fa),
	)
}

//go:inline
func ApEitherS[
	GS1 ~func(R) ET.Either[E, S1],
	GS2 ~func(R) ET.Either[E, S2],
	R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa ET.Either[E, T],
) func(GS1) GS2 {
	return ApS[GS1, GS2](
		setter,
		FromEither[func(R) ET.Either[E, T]](fa),
	)
}
