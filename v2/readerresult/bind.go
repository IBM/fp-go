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
	F "github.com/IBM/fp-go/v2/function"
	RRI "github.com/IBM/fp-go/v2/idiomatic/readerresult"
	RI "github.com/IBM/fp-go/v2/idiomatic/result"
	L "github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/reader"
	G "github.com/IBM/fp-go/v2/readereither/generic"
	"github.com/IBM/fp-go/v2/result"
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
func Do[R, S any](
	empty S,
) ReaderResult[R, S] {
	return G.Do[ReaderResult[R, S]](empty)
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
//	        func(s State) readereither.ReaderResult[Env, error, User] {
//	            return readereither.Asks(func(env Env) either.Either[error, User] {
//	                return env.UserService.GetUser()
//	            })
//	        },
//	    ),
//	    readereither.Bind(
//	        func(cfg Config) func(State) State {
//	            return func(s State) State { s.Config = cfg; return s }
//	        },
//	        func(s State) readereither.ReaderResult[Env, error, Config] {
//	            // This can access s.User from the previous step
//	            return readereither.Asks(func(env Env) either.Either[error, Config] {
//	                return env.ConfigService.GetConfigForUser(s.User.ID)
//	            })
//	        },
//	    ),
//	)
//
//go:inline
func Bind[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	f Kleisli[R, S1, T],
) Operator[R, S1, S2] {
	return G.Bind[ReaderResult[R, S1], ReaderResult[R, S2]](setter, f)
}

// BindI attaches the result of an idiomatic computation to a context [S1] to produce a context [S2].
// This is the idiomatic version of Bind, where the computation returns (T, error) instead of Result[T].
// This enables sequential composition with Go's native error handling style where each step can depend
// on the results of previous steps and access the shared environment.
//
// Example:
//
//	type State struct {
//	    User   User
//	    Config Config
//	}
//	type Env struct {
//	    UserService UserService
//	}
//
//	// Idiomatic function returning (User, error)
//	getUser := func(s State) func(env Env) (User, error) {
//	    return func(env Env) (User, error) {
//	        return env.UserService.GetUser()
//	    }
//	}
//
//	result := F.Pipe1(
//	    readerresult.Do[Env](State{}),
//	    readerresult.BindI(
//	        func(user User) func(State) State {
//	            return func(s State) State { s.User = user; return s }
//	        },
//	        getUser,
//	    ),
//	)
//
//go:inline
func BindI[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	f RRI.Kleisli[R, S1, T],
) Operator[R, S1, S2] {
	return Bind(setter, fromReaderResultKleisliI(f))
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
//
//go:inline
func Let[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) Operator[R, S1, S2] {
	return G.Let[ReaderResult[R, S1], ReaderResult[R, S2]](setter, f)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
//
//go:inline
func LetTo[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) Operator[R, S1, S2] {
	return G.LetTo[ReaderResult[R, S1], ReaderResult[R, S2]](setter, b)
}

// BindTo initializes a new state [S1] from a value [T]
//
//go:inline
func BindTo[R, S1, T any](
	setter func(T) S1,
) Operator[R, T, S1] {
	return G.BindTo[ReaderResult[R, S1], ReaderResult[R, T]](setter)
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
//
//go:inline
func ApS[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa ReaderResult[R, T],
) Operator[R, S1, S2] {
	return G.ApS[ReaderResult[R, S1], ReaderResult[R, S2]](setter, fa)
}

// ApIS attaches a value from an idiomatic ReaderResult to a context [S1] to produce a context [S2].
// This is the idiomatic version of ApS, where the computation returns (T, error) instead of Result[T].
// Unlike BindI which sequences operations, ApIS uses applicative semantics, meaning the computation
// is independent of the current state and can conceptually run in parallel.
//
// Example:
//
//	type State struct {
//	    User   User
//	    Config Config
//	}
//	type Env struct {
//	    UserService UserService
//	}
//
//	// Idiomatic independent computation returning (User, error)
//	getUser := func(env Env) (User, error) {
//	    return env.UserService.GetUser()
//	}
//
//	result := F.Pipe1(
//	    readerresult.Do[Env](State{}),
//	    readerresult.ApIS(
//	        func(user User) func(State) State {
//	            return func(s State) State { s.User = user; return s }
//	        },
//	        getUser,
//	    ),
//	)
//
//go:inline
func ApIS[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa RRI.ReaderResult[R, T],
) Operator[R, S1, S2] {
	return ApS(setter, FromReaderResultI(fa))
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
//
//go:inline
func ApSL[R, S, T any](
	lens L.Lens[S, T],
	fa ReaderResult[R, T],
) Operator[R, S, S] {
	return ApS(lens.Set, fa)
}

// ApISL attaches a value from an idiomatic ReaderResult to a context using a lens-based setter.
// This is the idiomatic version of ApSL, where the computation returns (T, error) instead of Result[T].
// It combines ApIS with a lens, allowing you to use optics to update nested structures in a more composable way.
//
// Example:
//
//	type State struct {
//	    User   User
//	    Config Config
//	}
//	type Env struct {
//	    ConfigService ConfigService
//	}
//
//	configLens := lens.MakeLens(
//	    func(s State) Config { return s.Config },
//	    func(s State, c Config) State { s.Config = c; return s },
//	)
//
//	// Idiomatic computation returning (Config, error)
//	getConfig := func(env Env) (Config, error) {
//	    return env.ConfigService.GetConfig()
//	}
//
//	result := F.Pipe1(
//	    readerresult.Of[Env](State{}),
//	    readerresult.ApISL(configLens, getConfig),
//	)
//
//go:inline
func ApISL[R, S, T any](
	lens L.Lens[S, T],
	fa RRI.ReaderResult[R, T],
) Operator[R, S, S] {
	return ApS(lens.Set, FromReaderResultI(fa))
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
//	    readereither.BindL(userLens, func(user User) readereither.ReaderResult[Env, error, User] {
//	        return readereither.Asks(func(env Env) either.Either[error, User] {
//	            return env.UserService.GetUser()
//	        })
//	    }),
//	)
//
//go:inline
func BindL[R, S, T any](
	lens L.Lens[S, T],
	f Kleisli[R, T, T],
) Operator[R, S, S] {
	return Bind(lens.Set, F.Flow2(lens.Get, f))
}

// BindIL is a variant of BindI that uses a lens to focus on a specific part of the context.
// This is the idiomatic version of BindL, where the computation returns (T, error) instead of Result[T].
// It provides a more ergonomic API when working with nested structures, eliminating the need to manually
// write setter functions while supporting Go's native error handling pattern.
//
// Example:
//
//	type State struct {
//	    User   User
//	    Config Config
//	}
//	type Env struct {
//	    UserService UserService
//	}
//
//	userLens := lens.MakeLens(
//	    func(s State) User { return s.User },
//	    func(s State, u User) State { s.User = u; return s },
//	)
//
//	// Idiomatic function returning (User, error)
//	updateUser := func(user User) func(env Env) (User, error) {
//	    return func(env Env) (User, error) {
//	        return env.UserService.UpdateUser(user)
//	    }
//	}
//
//	result := F.Pipe1(
//	    readerresult.Do[Env](State{}),
//	    readerresult.BindIL(userLens, updateUser),
//	)
//
//go:inline
func BindIL[R, S, T any](
	lens L.Lens[S, T],
	f RRI.Kleisli[R, T, T],
) Operator[R, S, S] {
	return Bind(lens.Set, F.Flow3(lens.Get, f, FromReaderResultI[R, T]))
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
//
//go:inline
func LetL[R, S, T any](
	lens L.Lens[S, T],
	f Endomorphism[T],
) Operator[R, S, S] {
	return Let[R](lens.Set, F.Flow2(lens.Get, f))
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
//
//go:inline
func LetToL[R, S, T any](
	lens L.Lens[S, T],
	b T,
) Operator[R, S, S] {
	return LetTo[R](lens.Set, b)
}

// BindReaderK lifts a Reader Kleisli arrow into a ReaderResult context and binds it to the state.
// This allows you to integrate pure Reader computations (that don't have error handling)
// into a ReaderResult computation chain.
//
// The function f takes the current state S1 and returns a Reader[R, T] computation.
// The result T is then attached to the state using the setter to produce state S2.
//
// Example:
//
//	type Env struct {
//	    ConfigPath string
//	}
//	type State struct {
//	    Config string
//	}
//
//	// A pure Reader computation that reads from environment
//	getConfigPath := func(s State) reader.Reader[Env, string] {
//	    return func(env Env) string {
//	        return env.ConfigPath
//	    }
//	}
//
//	result := F.Pipe2(
//	    readerresult.Do[Env](State{}),
//	    readerresult.BindReaderK(
//	        func(path string) func(State) State {
//	            return func(s State) State { s.Config = path; return s }
//	        },
//	        getConfigPath,
//	    ),
//	)
//
//go:inline
func BindReaderK[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	f reader.Kleisli[R, S1, T],
) Operator[R, S1, S2] {
	return G.BindReaderK[ReaderResult[R, S1], ReaderResult[R, S2]](setter, f)
}

//go:inline
func BindEitherK[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	f result.Kleisli[S1, T],
) Operator[R, S1, S2] {
	return G.BindEitherK[ReaderResult[R, S1], ReaderResult[R, S2]](setter, f)
}

// BindEitherIK lifts an idiomatic Result Kleisli arrow into a ReaderResult context and binds it to the state.
// This is the idiomatic version of BindEitherK, where the function returns (T, error) instead of Result[T].
// It allows you to integrate idiomatic Result computations (that may fail but don't need environment access)
// into a ReaderResult computation chain.
//
// Example:
//
//	type State struct {
//	    Value       int
//	    ParsedValue int
//	}
//
//	// Idiomatic function returning (int, error)
//	parseValue := func(s State) (int, error) {
//	    if s.Value < 0 {
//	        return 0, errors.New("negative value")
//	    }
//	    return s.Value * 2, nil
//	}
//
//	result := F.Pipe1(
//	    readerresult.Do[context.Context](State{Value: 5}),
//	    readerresult.BindEitherIK[context.Context](
//	        func(parsed int) func(State) State {
//	            return func(s State) State { s.ParsedValue = parsed; return s }
//	        },
//	        parseValue,
//	    ),
//	)
//
//go:inline
func BindEitherIK[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	f RI.Kleisli[S1, T],
) Operator[R, S1, S2] {
	return BindEitherK[R](setter, fromResultKleisliI(f))
}

// BindResultK lifts a Result Kleisli arrow into a ReaderResult context and binds it to the state.
// This allows you to integrate Result computations (that may fail with an error but don't need
// environment access) into a ReaderResult computation chain.
//
// The function f takes the current state S1 and returns a Result[T] computation.
// If the Result is successful, the value T is attached to the state using the setter to produce state S2.
// If the Result is an error, the entire computation short-circuits with that error.
//
// Example:
//
//	type State struct {
//	    Value int
//	    ParsedValue int
//	}
//
//	// A Result computation that may fail
//	parseValue := func(s State) result.Result[int] {
//	    if s.Value < 0 {
//	        return result.Error[int](errors.New("negative value"))
//	    }
//	    return result.Of(s.Value * 2)
//	}
//
//	result := F.Pipe2(
//	    readerresult.Do[any](State{Value: 5}),
//	    readerresult.BindResultK(
//	        func(parsed int) func(State) State {
//	            return func(s State) State { s.ParsedValue = parsed; return s }
//	        },
//	        parseValue,
//	    ),
//	)
//
//go:inline
func BindResultK[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	f result.Kleisli[S1, T],
) Operator[R, S1, S2] {
	return G.BindEitherK[ReaderResult[R, S1], ReaderResult[R, S2]](setter, f)
}

// BindResultIK is an alias for BindEitherIK.
// It lifts an idiomatic Result Kleisli arrow into a ReaderResult context and binds it to the state.
// The function f returns (T, error) in Go's idiomatic style.
//
//go:inline
func BindResultIK[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	f RI.Kleisli[S1, T],
) Operator[R, S1, S2] {
	return BindResultK[R](setter, fromResultKleisliI(f))
}

// BindToReader initializes a new state S1 from a Reader[R, T] computation.
// This is used to start a ReaderResult computation chain from a pure Reader value.
//
// The setter function takes the result T from the Reader and initializes the state S1.
// This is useful when you want to begin a do-notation chain with a Reader computation
// that doesn't involve error handling.
//
// Example:
//
//	type Env struct {
//	    ConfigPath string
//	}
//	type State struct {
//	    Config string
//	}
//
//	// A Reader that just reads from the environment
//	getConfigPath := func(env Env) string {
//	    return env.ConfigPath
//	}
//
//	result := F.Pipe1(
//	    reader.Of[Env](getConfigPath),
//	    readerresult.BindToReader(func(path string) State {
//	        return State{Config: path}
//	    }),
//	)
//
//go:inline
func BindToReader[
	R, S1, T any](
	setter func(T) S1,
) func(Reader[R, T]) ReaderResult[R, S1] {
	return G.BindToReader[ReaderResult[R, S1], Reader[R, T]](setter)
}

//go:inline
func BindToEither[
	R, S1, T any](
	setter func(T) S1,
) func(Result[T]) ReaderResult[R, S1] {
	return G.BindToEither[ReaderResult[R, S1]](setter)
}

// BindToEitherI initializes a new state S1 from an idiomatic (value, error) pair.
// This is the idiomatic version of BindToEither, accepting Go's native error handling pattern.
// It's used to start a ReaderResult computation chain from an idiomatic Result that may contain an error.
//
// Example:
//
//	type State struct {
//	    Value int
//	}
//
//	// Idiomatic result from parsing
//	value, err := strconv.Atoi("42")
//
//	computation := readerresult.BindToEitherI[context.Context](
//	    func(value int) State {
//	        return State{Value: value}
//	    },
//	)(value, err)
//
//go:inline
func BindToEitherI[
	R, S1, T any](
	setter func(T) S1,
) func(T, error) ReaderResult[R, S1] {
	bte := BindToEither[R](setter)
	return func(t T, err error) ReaderResult[R, S1] {
		return bte(result.TryCatchError(t, err))
	}
}

// BindToResult initializes a new state S1 from a Result[T] value.
// This is used to start a ReaderResult computation chain from a Result that may contain an error.
//
// The setter function takes the successful result T and initializes the state S1.
// If the Result is an error, the entire computation will carry that error forward.
// This is useful when you want to begin a do-notation chain with a Result computation
// that doesn't need environment access.
//
// Example:
//
//	type State struct {
//	    Value int
//	}
//
//	// A Result that might contain an error
//	parseResult := result.TryCatch(func() int {
//	    // some parsing logic that might fail
//	    return 42
//	})
//
//	computation := F.Pipe1(
//	    parseResult,
//	    readerresult.BindToResult[any](func(value int) State {
//	        return State{Value: value}
//	    }),
//	)
//
//go:inline
func BindToResult[
	R, S1, T any](
	setter func(T) S1,
) func(Result[T]) ReaderResult[R, S1] {
	return G.BindToEither[ReaderResult[R, S1]](setter)
}

// BindToResultI is an alias for BindToEitherI.
// It initializes a new state S1 from an idiomatic (value, error) pair.
//
//go:inline
func BindToResultI[
	R, S1, T any](
	setter func(T) S1,
) func(T, error) ReaderResult[R, S1] {
	return BindToEitherI[R](setter)
}

// ApReaderS attaches a value from a pure Reader computation to a context [S1] to produce a context [S2]
// using Applicative semantics (independent, non-sequential composition).
//
// Unlike BindReaderK which uses monadic bind (sequential), ApReaderS uses applicative apply,
// meaning the Reader computation fa is independent of the current state and can conceptually
// execute in parallel.
//
// This is useful when you want to combine a Reader computation with your ReaderResult state
// without creating a dependency between them.
//
// Example:
//
//	type Env struct {
//	    DefaultPort int
//	    DefaultHost string
//	}
//	type State struct {
//	    Port int
//	    Host string
//	}
//
//	getDefaultPort := func(env Env) int { return env.DefaultPort }
//	getDefaultHost := func(env Env) string { return env.DefaultHost }
//
//	result := F.Pipe2(
//	    readerresult.Do[Env](State{}),
//	    readerresult.ApReaderS(
//	        func(port int) func(State) State {
//	            return func(s State) State { s.Port = port; return s }
//	        },
//	        getDefaultPort,
//	    ),
//	    readerresult.ApReaderS(
//	        func(host string) func(State) State {
//	            return func(s State) State { s.Host = host; return s }
//	        },
//	        getDefaultHost,
//	    ),
//	)
//
//go:inline
func ApReaderS[
	R, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Reader[R, T],
) Operator[R, S1, S2] {
	return G.ApReaderS[ReaderResult[R, S1], ReaderResult[R, S2]](setter, fa)
}

//go:inline
func ApEitherS[
	R, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Result[T],
) Operator[R, S1, S2] {
	return G.ApEitherS[ReaderResult[R, S1], ReaderResult[R, S2]](setter, fa)
}

// ApEitherIS attaches a value from an idiomatic (value, error) pair to a context [S1] to produce a context [S2].
// This is the idiomatic version of ApEitherS, accepting Go's native error handling pattern.
// It uses Applicative semantics (independent, non-sequential composition).
//
// Example:
//
//	type State struct {
//	    Value1 int
//	    Value2 int
//	}
//
//	// Idiomatic parsing result
//	value, err := strconv.Atoi("42")
//
//	computation := F.Pipe1(
//	    readerresult.Do[context.Context](State{}),
//	    readerresult.ApEitherIS[context.Context](
//	        func(v int) func(State) State {
//	            return func(s State) State { s.Value1 = v; return s }
//	        },
//	    )(value, err),
//	)
//
//go:inline
func ApEitherIS[
	R, S1, S2, T any](
	setter func(T) func(S1) S2,
) func(T, error) Operator[R, S1, S2] {
	return func(t T, err error) Operator[R, S1, S2] {
		return ApEitherS[R](setter, result.TryCatchError(t, err))
	}
}

// ApResultS attaches a value from a Result to a context [S1] to produce a context [S2]
// using Applicative semantics (independent, non-sequential composition).
//
// Unlike BindResultK which uses monadic bind (sequential), ApResultS uses applicative apply,
// meaning the Result computation fa is independent of the current state and can conceptually
// execute in parallel.
//
// If the Result fa contains an error, the entire computation short-circuits with that error.
// This is useful when you want to combine a Result value with your ReaderResult state
// without creating a dependency between them.
//
// Example:
//
//	type State struct {
//	    Value1 int
//	    Value2 int
//	}
//
//	// Independent Result computations
//	parseValue1 := result.TryCatch(func() int { return 42 })
//	parseValue2 := result.TryCatch(func() int { return 100 })
//
//	computation := F.Pipe2(
//	    readerresult.Do[any](State{}),
//	    readerresult.ApResultS(
//	        func(v int) func(State) State {
//	            return func(s State) State { s.Value1 = v; return s }
//	        },
//	        parseValue1,
//	    ),
//	    readerresult.ApResultS(
//	        func(v int) func(State) State {
//	            return func(s State) State { s.Value2 = v; return s }
//	        },
//	        parseValue2,
//	    ),
//	)
//
//go:inline
func ApResultS[
	R, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Result[T],
) Operator[R, S1, S2] {
	return G.ApEitherS[ReaderResult[R, S1], ReaderResult[R, S2]](setter, fa)
}

// ApResultIS is an alias for ApEitherIS.
// It attaches a value from an idiomatic (value, error) pair to a context [S1] to produce a context [S2].
//
//go:inline
func ApResultIS[
	R, S1, S2, T any](
	setter func(T) func(S1) S2,
) func(T, error) Operator[R, S1, S2] {
	return ApEitherIS[R](setter)
}
