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
	"github.com/IBM/fp-go/v2/either"
	A "github.com/IBM/fp-go/v2/internal/apply"
	C "github.com/IBM/fp-go/v2/internal/chain"
	F "github.com/IBM/fp-go/v2/internal/functor"
)

// Do creates an empty context of type [S] to be used with the [Bind] operation.
// This is the starting point for do-notation style composition.
//
// Example:
//
//	type State struct {
//	    User   User
//	    Config Config
//	}
//	type Env struct {
//	    UserRepo   UserRepository
//	    ConfigRepo ConfigRepository
//	}
//	result := generic.Do[ReaderIOEither[Env, error, State], IOEither[error, State], Env, error, State](State{})
func Do[GRS ~func(R) GS, GS ~func() either.Either[E, S], R, E, S any](
	empty S,
) GRS {
	return Of[GRS](empty)
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
//	    User   User
//	    Config Config
//	}
//	type Env struct {
//	    UserRepo   UserRepository
//	    ConfigRepo ConfigRepository
//	}
//
//	result := F.Pipe2(
//	    generic.Do[ReaderIOEither[Env, error, State], IOEither[error, State], Env, error, State](State{}),
//	    generic.Bind[ReaderIOEither[Env, error, State], ReaderIOEither[Env, error, State], ReaderIOEither[Env, error, User], IOEither[error, State], IOEither[error, State], IOEither[error, User], Env, error, State, State, User](
//	        func(user User) func(State) State {
//	            return func(s State) State { s.User = user; return s }
//	        },
//	        func(s State) ReaderIOEither[Env, error, User] {
//	            return func(env Env) ioeither.IOEither[error, User] {
//	                return env.UserRepo.FindUser()
//	            }
//	        },
//	    ),
//	    generic.Bind[ReaderIOEither[Env, error, State], ReaderIOEither[Env, error, State], ReaderIOEither[Env, error, Config], IOEither[error, State], IOEither[error, State], IOEither[error, Config], Env, error, State, State, Config](
//	        func(cfg Config) func(State) State {
//	            return func(s State) State { s.Config = cfg; return s }
//	        },
//	        func(s State) ReaderIOEither[Env, error, Config] {
//	            // This can access s.User from the previous step
//	            return func(env Env) ioeither.IOEither[error, Config] {
//	                return env.ConfigRepo.LoadConfigForUser(s.User.ID)
//	            }
//	        },
//	    ),
//	)
func Bind[GRS1 ~func(R) GS1, GRS2 ~func(R) GS2, GRT ~func(R) GT, GS1 ~func() either.Either[E, S1], GS2 ~func() either.Either[E, S2], GT ~func() either.Either[E, T], R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) GRT,
) func(GRS1) GRS2 {
	return C.Bind(
		Chain[GRS1, GRS2, GS1, GS2, R, E, S1, S2],
		Map[GRT, GRS2, GT, GS2, R, E, T, S2],
		setter,
		f,
	)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
func Let[GRS1 ~func(R) GS1, GRS2 ~func(R) GS2, GS1 ~func() either.Either[E, S1], GS2 ~func() either.Either[E, S2], R, E, S1, S2, T any](
	key func(T) func(S1) S2,
	f func(S1) T,
) func(GRS1) GRS2 {
	return F.Let(
		Map[GRS1, GRS2, GS1, GS2, R, E, S1, S2],
		key,
		f,
	)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
func LetTo[GRS1 ~func(R) GS1, GRS2 ~func(R) GS2, GS1 ~func() either.Either[E, S1], GS2 ~func() either.Either[E, S2], R, E, S1, S2, B any](
	key func(B) func(S1) S2,
	b B,
) func(GRS1) GRS2 {
	return F.LetTo(
		Map[GRS1, GRS2, GS1, GS2, R, E, S1, S2],
		key,
		b,
	)
}

// BindTo initializes a new state [S1] from a value [T]
func BindTo[GRS1 ~func(R) GS1, GRT ~func(R) GT, GS1 ~func() either.Either[E, S1], GT ~func() either.Either[E, T], R, E, S1, T any](
	setter func(T) S1,
) func(GRT) GRS1 {
	return C.BindTo(
		Map[GRT, GRS1, GT, GS1, R, E, T, S1],
		setter,
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
//	    User   User
//	    Config Config
//	}
//	type Env struct {
//	    UserRepo   UserRepository
//	    ConfigRepo ConfigRepository
//	}
//
//	// These operations are independent and can be combined with ApS
//	getUser := func(env Env) ioeither.IOEither[error, User] {
//	    return env.UserRepo.FindUser()
//	}
//	getConfig := func(env Env) ioeither.IOEither[error, Config] {
//	    return env.ConfigRepo.LoadConfig()
//	}
//
//	result := F.Pipe2(
//	    generic.Do[ReaderIOEither[Env, error, State], IOEither[error, State], Env, error, State](State{}),
//	    generic.ApS[...](
//	        func(user User) func(State) State {
//	            return func(s State) State { s.User = user; return s }
//	        },
//	        getUser,
//	    ),
//	    generic.ApS[...](
//	        func(cfg Config) func(State) State {
//	            return func(s State) State { s.Config = cfg; return s }
//	        },
//	        getConfig,
//	    ),
//	)
func ApS[GRTS1 ~func(R) GTS1, GRS1 ~func(R) GS1, GRS2 ~func(R) GS2, GRT ~func(R) GT, GTS1 ~func() either.Either[E, func(T) S2], GS1 ~func() either.Either[E, S1], GS2 ~func() either.Either[E, S2], GT ~func() either.Either[E, T], R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa GRT,
) func(GRS1) GRS2 {
	return A.ApS(
		Ap[GRT, GRS2, GRTS1, GT, GS2, GTS1, R, E, T, S2],
		Map[GRS1, GRTS1, GS1, GTS1, R, E, S1, func(T) S2],
		setter,
		fa,
	)
}
