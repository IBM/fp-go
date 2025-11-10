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

package readerioresult

import (
	L "github.com/IBM/fp-go/v2/optics/lens"
	RIOE "github.com/IBM/fp-go/v2/readerioeither"
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
func Do[R, S any](
	empty S,
) ReaderIOResult[R, S] {
	return RIOE.Do[R, error](empty)
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
//	        func(s State) readerioeither.ReaderIOResult[Env, error, User] {
//	            return readerioeither.Asks(func(env Env) ioeither.IOEither[error, User] {
//	                return env.UserRepo.FindUser()
//	            })
//	        },
//	    ),
//	    readerioeither.Bind(
//	        func(posts []Post) func(State) State {
//	            return func(s State) State { s.Posts = posts; return s }
//	        },
//	        func(s State) readerioeither.ReaderIOResult[Env, error, []Post] {
//	            // This can access s.User from the previous step
//	            return readerioeither.Asks(func(env Env) ioeither.IOEither[error, []Post] {
//	                return env.PostRepo.FindPostsByUser(s.User.ID)
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
	return RIOE.Bind(setter, f)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
//
//go:inline
func Let[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) Operator[R, S1, S2] {
	return RIOE.Let[R, error](setter, f)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
//
//go:inline
func LetTo[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) Operator[R, S1, S2] {
	return RIOE.LetTo[R, error](setter, b)
}

// BindTo initializes a new state [S1] from a value [T]
//
//go:inline
func BindTo[R, S1, T any](
	setter func(T) S1,
) Operator[R, T, S1] {
	return RIOE.BindTo[R, error](setter)
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
func ApS[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa ReaderIOResult[R, T],
) Operator[R, S1, S2] {
	return RIOE.ApS(setter, fa)
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
func ApSL[R, S, T any](
	lens L.Lens[S, T],
	fa ReaderIOResult[R, T],
) Operator[R, S, S] {
	return RIOE.ApSL(lens, fa)
}

// BindL is a variant of Bind that uses a lens to focus on a specific part of the context.
// This provides a more ergonomic API when working with nested structures, eliminating
// the need to manually write setter functions.
//
// The lens parameter provides both a getter and setter for a field of type T within
// the context S. The function f receives the current value of the focused field and
// returns a ReaderIOResult computation that produces an updated value.
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
//	    readerioeither.BindL(userLens, func(user User) readerioeither.ReaderIOResult[Env, error, User] {
//	        return readerioeither.Asks(func(env Env) ioeither.IOEither[error, User] {
//	            return env.UserRepo.FindUser()
//	        })
//	    }),
//	)
//
//go:inline
func BindL[R, S, T any](
	lens L.Lens[S, T],
	f Kleisli[R, T, T],
) Operator[R, S, S] {
	return RIOE.BindL(lens, f)
}

// LetL is a variant of Let that uses a lens to focus on a specific part of the context.
// This provides a more ergonomic API when working with nested structures, eliminating
// the need to manually write setter functions.
//
// The lens parameter provides both a getter and setter for a field of type T within
// the context S. The function f receives the current value of the focused field and
// returns a new value (without wrapping in a ReaderIOResult).
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
func LetL[R, S, T any](
	lens L.Lens[S, T],
	f func(T) T,
) Operator[R, S, S] {
	return RIOE.LetL[R, error](lens, f)
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
func LetToL[R, S, T any](
	lens L.Lens[S, T],
	b T,
) Operator[R, S, S] {
	return RIOE.LetToL[R, error](lens, b)
}
