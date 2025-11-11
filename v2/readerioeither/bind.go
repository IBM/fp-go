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
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	L "github.com/IBM/fp-go/v2/optics/lens"
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
	return G.Let[ReaderIOEither[R, E, S1], ReaderIOEither[R, E, S2], IOE.IOEither[E, S1], IOE.IOEither[E, S2], R, E, S1, S2, T](setter, f)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
//
//go:inline
func LetTo[R, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) Operator[R, E, S1, S2] {
	return G.LetTo[ReaderIOEither[R, E, S1], ReaderIOEither[R, E, S2], IOE.IOEither[E, S1], IOE.IOEither[E, S2], R, E, S1, S2, T](setter, b)
}

// BindTo initializes a new state [S1] from a value [T]
//
//go:inline
func BindTo[R, E, S1, T any](
	setter func(T) S1,
) Operator[R, E, T, S1] {
	return G.BindTo[ReaderIOEither[R, E, S1], ReaderIOEither[R, E, T], IOE.IOEither[E, S1], IOE.IOEither[E, T], R, E, S1, T](setter)
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
	return G.ApS[ReaderIOEither[R, E, func(T) S2], ReaderIOEither[R, E, S1], ReaderIOEither[R, E, S2], ReaderIOEither[R, E, T], IOE.IOEither[E, func(T) S2], IOE.IOEither[E, S1], IOE.IOEither[E, S2], IOE.IOEither[E, T], R, E, S1, S2, T](setter, fa)
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
	return Bind[R, E, S, S, T](lens.Set, F.Flow2(lens.Get, f))
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
	return Let[R, E, S, S, T](lens.Set, F.Flow2(lens.Get, f))
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
	return LetTo[R, E, S, S, T](lens.Set, b)
}
