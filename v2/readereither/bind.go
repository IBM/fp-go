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
	F "github.com/IBM/fp-go/v2/function"
	L "github.com/IBM/fp-go/v2/optics/lens"
	G "github.com/IBM/fp-go/v2/readereither/generic"
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
//	    UserService   UserService
//	    ConfigService ConfigService
//	}
//	result := readereither.Do[Env, error](State{})
//
//go:inline
func Do[R, E, S any](
	empty S,
) ReaderEither[R, E, S] {
	return G.Do[ReaderEither[R, E, S]](empty)
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
//	    UserService   UserService
//	    ConfigService ConfigService
//	}
//
//	result := F.Pipe2(
//	    readereither.Do[Env, error](State{}),
//	    readereither.Bind(
//	        func(user User) func(State) State {
//	            return func(s State) State { s.User = user; return s }
//	        },
//	        func(s State) readereither.ReaderEither[Env, error, User] {
//	            return readereither.Asks(func(env Env) either.Either[error, User] {
//	                return env.UserService.GetUser()
//	            })
//	        },
//	    ),
//	    readereither.Bind(
//	        func(cfg Config) func(State) State {
//	            return func(s State) State { s.Config = cfg; return s }
//	        },
//	        func(s State) readereither.ReaderEither[Env, error, Config] {
//	            // This can access s.User from the previous step
//	            return readereither.Asks(func(env Env) either.Either[error, Config] {
//	                return env.ConfigService.GetConfigForUser(s.User.ID)
//	            })
//	        },
//	    ),
//	)
//
//go:inline
func Bind[R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) ReaderEither[R, E, T],
) func(ReaderEither[R, E, S1]) ReaderEither[R, E, S2] {
	return G.Bind[ReaderEither[R, E, S1], ReaderEither[R, E, S2]](setter, f)
}

//go:inline
func BindReaderK[R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) Reader[R, T],
) func(ReaderEither[R, E, S1]) ReaderEither[R, E, S2] {
	return G.BindReaderK[ReaderEither[R, E, S1], ReaderEither[R, E, S2]](setter, f)
}

//go:inline
func BindEitherK[R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) Either[E, T],
) func(ReaderEither[R, E, S1]) ReaderEither[R, E, S2] {
	return G.BindEitherK[ReaderEither[R, E, S1], ReaderEither[R, E, S2]](setter, f)
}

//go:inline
func BindToReader[
	R, E, S1, T any](
	setter func(T) S1,
) func(Reader[R, T]) ReaderEither[R, E, S1] {
	return G.BindToReader[ReaderEither[R, E, S1], Reader[R, T]](setter)
}

//go:inline
func BindToEither[
	R, E, S1, T any](
	setter func(T) S1,
) func(ET.Either[E, T]) ReaderEither[R, E, S1] {
	return G.BindToEither[ReaderEither[R, E, S1]](setter)
}

//go:inline
func ApReaderS[
	R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Reader[R, T],
) Operator[R, E, S1, S2] {
	return G.ApReaderS[ReaderEither[R, E, S1], ReaderEither[R, E, S2]](setter, fa)
}

//go:inline
func ApEitherS[
	R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa ET.Either[E, T],
) Operator[R, E, S1, S2] {
	return G.ApEitherS[ReaderEither[R, E, S1], ReaderEither[R, E, S2]](setter, fa)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
//
//go:inline
func Let[R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) func(ReaderEither[R, E, S1]) ReaderEither[R, E, S2] {
	return G.Let[ReaderEither[R, E, S1], ReaderEither[R, E, S2]](setter, f)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
func LetTo[R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) func(ReaderEither[R, E, S1]) ReaderEither[R, E, S2] {
	return G.LetTo[ReaderEither[R, E, S1], ReaderEither[R, E, S2]](setter, b)
}

// BindTo initializes a new state [S1] from a value [T]
func BindTo[R, E, S1, T any](
	setter func(T) S1,
) func(ReaderEither[R, E, T]) ReaderEither[R, E, S1] {
	return G.BindTo[ReaderEither[R, E, S1], ReaderEither[R, E, T]](setter)
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
//	    UserService   UserService
//	    ConfigService ConfigService
//	}
//
//	// These operations are independent and can be combined with ApS
//	getUser := readereither.Asks(func(env Env) either.Either[error, User] {
//	    return env.UserService.GetUser()
//	})
//	getConfig := readereither.Asks(func(env Env) either.Either[error, Config] {
//	    return env.ConfigService.GetConfig()
//	})
//
//	result := F.Pipe2(
//	    readereither.Do[Env, error](State{}),
//	    readereither.ApS(
//	        func(user User) func(State) State {
//	            return func(s State) State { s.User = user; return s }
//	        },
//	        getUser,
//	    ),
//	    readereither.ApS(
//	        func(cfg Config) func(State) State {
//	            return func(s State) State { s.Config = cfg; return s }
//	        },
//	        getConfig,
//	    ),
//	)
func ApS[R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa ReaderEither[R, E, T],
) func(ReaderEither[R, E, S1]) ReaderEither[R, E, S2] {
	return G.ApS[ReaderEither[R, E, S1], ReaderEither[R, E, S2]](setter, fa)
}

// ApSL attaches a value to a context using a lens-based setter.
// This is a convenience function that combines ApS with a lens, allowing you to use
// optics to update nested structures in a more composable way.
//
// The lens parameter provides both the getter and setter for a field within the structure S.
// This eliminates the need to manually write setter functions.
//
// Example:
//
//	type State struct {
//	    User   User
//	    Config Config
//	}
//	type Env struct {
//	    UserService   UserService
//	    ConfigService ConfigService
//	}
//
//	configLens := lens.MakeLens(
//	    func(s State) Config { return s.Config },
//	    func(s State, c Config) State { s.Config = c; return s },
//	)
//
//	getConfig := readereither.Asks(func(env Env) either.Either[error, Config] {
//	    return env.ConfigService.GetConfig()
//	})
//	result := F.Pipe2(
//	    readereither.Of[Env, error](State{}),
//	    readereither.ApSL(configLens, getConfig),
//	)
func ApSL[R, E, S, T any](
	lens L.Lens[S, T],
	fa ReaderEither[R, E, T],
) func(ReaderEither[R, E, S]) ReaderEither[R, E, S] {
	return ApS(lens.Set, fa)
}

// BindL is a variant of Bind that uses a lens to focus on a specific part of the context.
// This provides a more ergonomic API when working with nested structures, eliminating
// the need to manually write setter functions.
//
// The lens parameter provides both a getter and setter for a field of type T within
// the context S. The function f receives the current value of the focused field and
// returns a ReaderEither computation that produces an updated value.
//
// Example:
//
//	type State struct {
//	    User   User
//	    Config Config
//	}
//	type Env struct {
//	    UserService   UserService
//	    ConfigService ConfigService
//	}
//
//	userLens := lens.MakeLens(
//	    func(s State) User { return s.User },
//	    func(s State, u User) State { s.User = u; return s },
//	)
//
//	result := F.Pipe2(
//	    readereither.Do[Env, error](State{}),
//	    readereither.BindL(userLens, func(user User) readereither.ReaderEither[Env, error, User] {
//	        return readereither.Asks(func(env Env) either.Either[error, User] {
//	            return env.UserService.GetUser()
//	        })
//	    }),
//	)
func BindL[R, E, S, T any](
	lens L.Lens[S, T],
	f func(T) ReaderEither[R, E, T],
) func(ReaderEither[R, E, S]) ReaderEither[R, E, S] {
	return Bind(lens.Set, F.Flow2(lens.Get, f))
}

// LetL is a variant of Let that uses a lens to focus on a specific part of the context.
// This provides a more ergonomic API when working with nested structures, eliminating
// the need to manually write setter functions.
//
// The lens parameter provides both a getter and setter for a field of type T within
// the context S. The function f receives the current value of the focused field and
// returns a new value (without wrapping in a ReaderEither).
//
// Example:
//
//	type State struct {
//	    User   User
//	    Config Config
//	}
//
//	configLens := lens.MakeLens(
//	    func(s State) Config { return s.Config },
//	    func(s State, c Config) State { s.Config = c; return s },
//	)
//
//	result := F.Pipe2(
//	    readereither.Do[any, error](State{Config: Config{Host: "localhost"}}),
//	    readereither.LetL(configLens, func(cfg Config) Config {
//	        cfg.Port = 8080
//	        return cfg
//	    }),
//	)
func LetL[R, E, S, T any](
	lens L.Lens[S, T],
	f func(T) T,
) func(ReaderEither[R, E, S]) ReaderEither[R, E, S] {
	return Let[R, E](lens.Set, F.Flow2(lens.Get, f))
}

// LetToL is a variant of LetTo that uses a lens to focus on a specific part of the context.
// This provides a more ergonomic API when working with nested structures, eliminating
// the need to manually write setter functions.
//
// The lens parameter provides both a getter and setter for a field of type T within
// the context S. The value b is set directly to the focused field.
//
// Example:
//
//	type State struct {
//	    User   User
//	    Config Config
//	}
//
//	configLens := lens.MakeLens(
//	    func(s State) Config { return s.Config },
//	    func(s State, c Config) State { s.Config = c; return s },
//	)
//
//	newConfig := Config{Host: "localhost", Port: 8080}
//	result := F.Pipe2(
//	    readereither.Do[any, error](State{}),
//	    readereither.LetToL(configLens, newConfig),
//	)
func LetToL[R, E, S, T any](
	lens L.Lens[S, T],
	b T,
) func(ReaderEither[R, E, S]) ReaderEither[R, E, S] {
	return LetTo[R, E](lens.Set, b)
}
