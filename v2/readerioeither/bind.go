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

package readerioeither

import (
	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	L "github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
	G "github.com/IBM/fp-go/v2/readerioeither/generic"
)

// Do creates an empty context of type [S] to be used with the [Bind] operation.
// This is the starting point for do-notation style composition.
//
// Example:
//
//	type State struct {
//	    User   User
//	    Posts  []Post
//	}
//	type Env struct {
//	    UserRepo UserRepository
//	    PostRepo PostRepository
//	}
//	result := readerioeither.Do[Env, error](State{})
//
//go:inline
func Do[R, E, S any](
	empty S,
) ReaderIOEither[R, E, S] {
	return G.Do[ReaderIOEither[R, E, S]](empty)
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
//	    Posts  []Post
//	}
//	type Env struct {
//	    UserRepo UserRepository
//	    PostRepo PostRepository
//	}
//
//	result := F.Pipe2(
//	    readerioeither.Do[Env, error](State{}),
//	    readerioeither.Bind(
//	        func(user User) func(State) State {
//	            return func(s State) State { s.User = user; return s }
//	        },
//	        func(s State) readerioeither.ReaderIOEither[Env, error, User] {
//	            return readerioeither.Asks(func(env Env) ioeither.IOEither[error, User] {
//	                return env.UserRepo.FindUser()
//	            })
//	        },
//	    ),
//	    readerioeither.Bind(
//	        func(posts []Post) func(State) State {
//	            return func(s State) State { s.Posts = posts; return s }
//	        },
//	        func(s State) readerioeither.ReaderIOEither[Env, error, []Post] {
//	            // This can access s.User from the previous step
//	            return readerioeither.Asks(func(env Env) ioeither.IOEither[error, []Post] {
//	                return env.PostRepo.FindPostsByUser(s.User.ID)
//	            })
//	        },
//	    ),
//	)
//
//go:inline
func Bind[R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) ReaderIOEither[R, E, T],
) Operator[R, E, S1, S2] {
	return G.Bind[ReaderIOEither[R, E, S1], ReaderIOEither[R, E, S2]](setter, f)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
//
//go:inline
func Let[R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) Operator[R, E, S1, S2] {
	return G.Let[ReaderIOEither[R, E, S1], ReaderIOEither[R, E, S2]](setter, f)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
//
//go:inline
func LetTo[R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) Operator[R, E, S1, S2] {
	return G.LetTo[ReaderIOEither[R, E, S1], ReaderIOEither[R, E, S2]](setter, b)
}

// BindTo initializes a new state [S1] from a value [T]
//
//go:inline
func BindTo[R, E, S1, T any](
	setter func(T) S1,
) Operator[R, E, T, S1] {
	return G.BindTo[ReaderIOEither[R, E, S1], ReaderIOEither[R, E, T]](setter)
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
//	    Posts  []Post
//	}
//	type Env struct {
//	    UserRepo UserRepository
//	    PostRepo PostRepository
//	}
//
//	// These operations are independent and can be combined with ApS
//	getUser := readerioeither.Asks(func(env Env) ioeither.IOEither[error, User] {
//	    return env.UserRepo.FindUser()
//	})
//	getPosts := readerioeither.Asks(func(env Env) ioeither.IOEither[error, []Post] {
//	    return env.PostRepo.FindPosts()
//	})
//
//	result := F.Pipe2(
//	    readerioeither.Do[Env, error](State{}),
//	    readerioeither.ApS(
//	        func(user User) func(State) State {
//	            return func(s State) State { s.User = user; return s }
//	        },
//	        getUser,
//	    ),
//	    readerioeither.ApS(
//	        func(posts []Post) func(State) State {
//	            return func(s State) State { s.Posts = posts; return s }
//	        },
//	        getPosts,
//	    ),
//	)
//
//go:inline
func ApS[R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa ReaderIOEither[R, E, T],
) Operator[R, E, S1, S2] {
	return G.ApS[ReaderIOEither[R, E, func(T) S2], ReaderIOEither[R, E, S1], ReaderIOEither[R, E, S2]](setter, fa)
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
//	    Posts  []Post
//	}
//	type Env struct {
//	    UserRepo UserRepository
//	    PostRepo PostRepository
//	}
//
//	userLens := lens.MakeLens(
//	    func(s State) User { return s.User },
//	    func(s State, u User) State { s.User = u; return s },
//	)
//
//	getUser := readerioeither.Asks(func(env Env) ioeither.IOEither[error, User] {
//	    return env.UserRepo.FindUser()
//	})
//	result := F.Pipe2(
//	    readerioeither.Of[Env, error](State{}),
//	    readerioeither.ApSL(userLens, getUser),
//	)
//
//go:inline
func ApSL[R, E, S, T any](
	lens L.Lens[S, T],
	fa ReaderIOEither[R, E, T],
) Operator[R, E, S, S] {
	return ApS(lens.Set, fa)
}

// BindL is a variant of Bind that uses a lens to focus on a specific part of the context.
// This provides a more ergonomic API when working with nested structures, eliminating
// the need to manually write setter functions.
//
// The lens parameter provides both a getter and setter for a field of type T within
// the context S. The function f receives the current value of the focused field and
// returns a ReaderIOEither computation that produces an updated value.
//
// Example:
//
//	type State struct {
//	    User   User
//	    Posts  []Post
//	}
//	type Env struct {
//	    UserRepo UserRepository
//	    PostRepo PostRepository
//	}
//
//	userLens := lens.MakeLens(
//	    func(s State) User { return s.User },
//	    func(s State, u User) State { s.User = u; return s },
//	)
//
//	result := F.Pipe2(
//	    readerioeither.Do[Env, error](State{}),
//	    readerioeither.BindL(userLens, func(user User) readerioeither.ReaderIOEither[Env, error, User] {
//	        return readerioeither.Asks(func(env Env) ioeither.IOEither[error, User] {
//	            return env.UserRepo.FindUser()
//	        })
//	    }),
//	)
//
//go:inline
func BindL[R, E, S, T any](
	lens L.Lens[S, T],
	f func(T) ReaderIOEither[R, E, T],
) Operator[R, E, S, S] {
	return Bind(lens.Set, F.Flow2(lens.Get, f))
}

// LetL is a variant of Let that uses a lens to focus on a specific part of the context.
// This provides a more ergonomic API when working with nested structures, eliminating
// the need to manually write setter functions.
//
// The lens parameter provides both a getter and setter for a field of type T within
// the context S. The function f receives the current value of the focused field and
// returns a new value (without wrapping in a ReaderIOEither).
//
// Example:
//
//	type State struct {
//	    User   User
//	    Posts  []Post
//	}
//
//	userLens := lens.MakeLens(
//	    func(s State) User { return s.User },
//	    func(s State, u User) State { s.User = u; return s },
//	)
//
//	result := F.Pipe2(
//	    readerioeither.Do[any, error](State{User: User{Name: "Alice"}}),
//	    readerioeither.LetL(userLens, func(user User) User {
//	        user.Name = "Bob"
//	        return user
//	    }),
//	)
//
//go:inline
func LetL[R, E, S, T any](
	lens L.Lens[S, T],
	f func(T) T,
) Operator[R, E, S, S] {
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
//	    Posts  []Post
//	}
//
//	userLens := lens.MakeLens(
//	    func(s State) User { return s.User },
//	    func(s State, u User) State { s.User = u; return s },
//	)
//
//	newUser := User{Name: "Bob", ID: 123}
//	result := F.Pipe2(
//	    readerioeither.Do[any, error](State{}),
//	    readerioeither.LetToL(userLens, newUser),
//	)
//
//go:inline
func LetToL[R, E, S, T any](
	lens L.Lens[S, T],
	b T,
) Operator[R, E, S, S] {
	return LetTo[R, E](lens.Set, b)
}

// BindIOEitherK is a variant of Bind that works with IOEither computations.
// It lifts an IOEither Kleisli arrow into the ReaderIOEither context, allowing you to
// compose IOEither operations within a do-notation chain.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - f: An IOEither Kleisli arrow (S1 -> IOEither[E, T])
//
// Returns:
//   - An Operator that can be used in a do-notation chain
func BindIOEitherK[R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f ioeither.Kleisli[E, S1, T],
) Operator[R, E, S1, S2] {
	return Bind(setter, F.Flow2(f, FromIOEither[R, E, T]))
}

// BindIOK is a variant of Bind that works with IO computations.
// It lifts an IO Kleisli arrow into the ReaderIOEither context.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - f: An IO Kleisli arrow (S1 -> IO[T])
func BindIOK[R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f io.Kleisli[S1, T],
) Operator[R, E, S1, S2] {
	return Bind(setter, F.Flow2(f, FromIO[R, E, T]))
}

// BindReaderK is a variant of Bind that works with Reader computations.
// It lifts a Reader Kleisli arrow into the ReaderIOEither context.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - f: A Reader Kleisli arrow (S1 -> Reader[R, T])
func BindReaderK[E, R, S1, S2, T any](
	setter func(T) func(S1) S2,
	f reader.Kleisli[R, S1, T],
) Operator[R, E, S1, S2] {
	return Bind(setter, F.Flow2(f, FromReader[E, R, T]))
}

// BindReaderIOK is a variant of Bind that works with ReaderIO computations.
// It lifts a ReaderIO Kleisli arrow into the ReaderIOEither context.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - f: A ReaderIO Kleisli arrow (S1 -> ReaderIO[R, T])
func BindReaderIOK[E, R, S1, S2, T any](
	setter func(T) func(S1) S2,
	f readerio.Kleisli[R, S1, T],
) Operator[R, E, S1, S2] {
	return Bind(setter, F.Flow2(f, FromReaderIO[E, R, T]))
}

// BindEitherK is a variant of Bind that works with Either computations.
// It lifts an Either Kleisli arrow into the ReaderIOEither context.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - f: An Either Kleisli arrow (S1 -> Either[E, T])
func BindEitherK[R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f either.Kleisli[E, S1, T],
) Operator[R, E, S1, S2] {
	return Bind(setter, F.Flow2(f, FromEither[R, E, T]))
}

// BindIOEitherKL is a lens-based variant of BindIOEitherK.
// It combines a lens with an IOEither Kleisli arrow, focusing on a specific field
// within the state structure.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - f: An IOEither Kleisli arrow (T -> IOEither[E, T])
func BindIOEitherKL[R, E, S, T any](
	lens L.Lens[S, T],
	f ioeither.Kleisli[E, T, T],
) Operator[R, E, S, S] {
	return BindL(lens, F.Flow2(f, FromIOEither[R, E, T]))
}

// BindIOKL is a lens-based variant of BindIOK.
// It combines a lens with an IO Kleisli arrow, focusing on a specific field
// within the state structure.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - f: An IO Kleisli arrow (T -> IO[T])
func BindIOKL[R, E, S, T any](
	lens L.Lens[S, T],
	f io.Kleisli[T, T],
) Operator[R, E, S, S] {
	return BindL(lens, F.Flow2(f, FromIO[R, E, T]))
}

// BindReaderKL is a lens-based variant of BindReaderK.
// It combines a lens with a Reader Kleisli arrow, focusing on a specific field
// within the state structure.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - f: A Reader Kleisli arrow (T -> Reader[R, T])
func BindReaderKL[E, R, S, T any](
	lens L.Lens[S, T],
	f reader.Kleisli[R, T, T],
) Operator[R, E, S, S] {
	return BindL(lens, F.Flow2(f, FromReader[E, R, T]))
}

// BindReaderIOKL is a lens-based variant of BindReaderIOK.
// It combines a lens with a ReaderIO Kleisli arrow, focusing on a specific field
// within the state structure.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - f: A ReaderIO Kleisli arrow (T -> ReaderIO[R, T])
func BindReaderIOKL[E, R, S, T any](
	lens L.Lens[S, T],
	f readerio.Kleisli[R, T, T],
) Operator[R, E, S, S] {
	return BindL(lens, F.Flow2(f, FromReaderIO[E, R, T]))
}

// ApIOEitherS is an applicative variant that works with IOEither values.
// Unlike BindIOEitherK, this uses applicative composition (ApS) instead of monadic
// composition (Bind), allowing independent computations to be combined.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - fa: An IOEither value
func ApIOEitherS[R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa IOEither[E, T],
) Operator[R, E, S1, S2] {
	return F.Bind2nd(F.Flow2[ReaderIOEither[R, E, S1], ioeither.Operator[E, S1, S2]], ioeither.ApS(setter, fa))
}

// ApIOS is an applicative variant that works with IO values.
// It lifts an IO value into the ReaderIOEither context using applicative composition.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - fa: An IO value
func ApIOS[R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa IO[T],
) Operator[R, E, S1, S2] {
	return ApS(setter, FromIO[R, E](fa))
}

// ApReaderS is an applicative variant that works with Reader values.
// It lifts a Reader value into the ReaderIOEither context using applicative composition.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - fa: A Reader value
func ApReaderS[R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Reader[R, T],
) Operator[R, E, S1, S2] {
	return ApS(setter, FromReader[E](fa))
}

// ApReaderIOS is an applicative variant that works with ReaderIO values.
// It lifts a ReaderIO value into the ReaderIOEither context using applicative composition.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - fa: A ReaderIO value
func ApReaderIOS[R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa ReaderIO[R, T],
) Operator[R, E, S1, S2] {
	return ApS(setter, FromReaderIO[E](fa))
}

// ApEitherS is an applicative variant that works with Either values.
// It lifts an Either value into the ReaderIOEither context using applicative composition.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - fa: An Either value
func ApEitherS[R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Either[E, T],
) Operator[R, E, S1, S2] {
	return ApS(setter, FromEither[R](fa))
}

// ApIOEitherSL is a lens-based variant of ApIOEitherS.
// It combines a lens with an IOEither value using applicative composition.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - fa: An IOEither value
func ApIOEitherSL[R, E, S, T any](
	lens L.Lens[S, T],
	fa IOEither[E, T],
) Operator[R, E, S, S] {
	return F.Bind2nd(F.Flow2[ReaderIOEither[R, E, S], ioeither.Operator[E, S, S]], ioeither.ApSL(lens, fa))
}

// ApIOSL is a lens-based variant of ApIOS.
// It combines a lens with an IO value using applicative composition.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - fa: An IO value
func ApIOSL[R, E, S, T any](
	lens L.Lens[S, T],
	fa IO[T],
) Operator[R, E, S, S] {
	return ApSL(lens, FromIO[R, E](fa))
}

// ApReaderSL is a lens-based variant of ApReaderS.
// It combines a lens with a Reader value using applicative composition.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - fa: A Reader value
func ApReaderSL[R, E, S, T any](
	lens L.Lens[S, T],
	fa Reader[R, T],
) Operator[R, E, S, S] {
	return ApSL(lens, FromReader[E](fa))
}

// ApReaderIOSL is a lens-based variant of ApReaderIOS.
// It combines a lens with a ReaderIO value using applicative composition.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - fa: A ReaderIO value
func ApReaderIOSL[R, E, S, T any](
	lens L.Lens[S, T],
	fa ReaderIO[R, T],
) Operator[R, E, S, S] {
	return ApSL(lens, FromReaderIO[E](fa))
}

// ApEitherSL is a lens-based variant of ApEitherS.
// It combines a lens with an Either value using applicative composition.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - fa: An Either value
func ApEitherSL[R, E, S, T any](
	lens L.Lens[S, T],
	fa Either[E, T],
) Operator[R, E, S, S] {
	return ApSL(lens, FromEither[R](fa))
}
