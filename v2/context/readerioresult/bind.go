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
	"context"

	"github.com/IBM/fp-go/v2/internal/apply"
	L "github.com/IBM/fp-go/v2/optics/lens"
	RIOR "github.com/IBM/fp-go/v2/readerioresult"
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
//	result := readerioeither.Do(State{})
//
//go:inline
func Do[S any](
	empty S,
) ReaderIOResult[S] {
	return RIOR.Of[context.Context](empty)
}

// Bind attaches the result of a computation to a context [S1] to produce a context [S2].
// This enables sequential composition where each step can depend on the results of previous steps
// and access the context.Context from the environment.
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
//
//	result := F.Pipe2(
//	    readerioeither.Do(State{}),
//	    readerioeither.Bind(
//	        func(user User) func(State) State {
//	            return func(s State) State { s.User = user; return s }
//	        },
//	        func(s State) readerioeither.ReaderIOResult[User] {
//	            return func(ctx context.Context) ioeither.IOEither[error, User] {
//	                return ioeither.TryCatch(func() (User, error) {
//	                    return fetchUser(ctx)
//	                })
//	            }
//	        },
//	    ),
//	    readerioeither.Bind(
//	        func(cfg Config) func(State) State {
//	            return func(s State) State { s.Config = cfg; return s }
//	        },
//	        func(s State) readerioeither.ReaderIOResult[Config] {
//	            // This can access s.User from the previous step
//	            return func(ctx context.Context) ioeither.IOEither[error, Config] {
//	                return ioeither.TryCatch(func() (Config, error) {
//	                    return fetchConfigForUser(ctx, s.User.ID)
//	                })
//	            }
//	        },
//	    ),
//	)
//
//go:inline
func Bind[S1, S2, T any](
	setter func(T) func(S1) S2,
	f Kleisli[S1, T],
) Operator[S1, S2] {
	return RIOR.Bind[context.Context](setter, f)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
//
//go:inline
func Let[S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) Operator[S1, S2] {
	return RIOR.Let[context.Context](setter, f)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
//
//go:inline
func LetTo[S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) Operator[S1, S2] {
	return RIOR.LetTo[context.Context](setter, b)
}

// BindTo initializes a new state [S1] from a value [T]
//
//go:inline
func BindTo[S1, T any](
	setter func(T) S1,
) Operator[T, S1] {
	return RIOR.BindTo[context.Context](setter)
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
//
//	// These operations are independent and can be combined with ApS
//	getUser := func(ctx context.Context) ioeither.IOEither[error, User] {
//	    return ioeither.TryCatch(func() (User, error) {
//	        return fetchUser(ctx)
//	    })
//	}
//	getConfig := func(ctx context.Context) ioeither.IOEither[error, Config] {
//	    return ioeither.TryCatch(func() (Config, error) {
//	        return fetchConfig(ctx)
//	    })
//	}
//
//	result := F.Pipe2(
//	    readerioeither.Do(State{}),
//	    readerioeither.ApS(
//	        func(user User) func(State) State {
//	            return func(s State) State { s.User = user; return s }
//	        },
//	        getUser,
//	    ),
//	    readerioeither.ApS(
//	        func(cfg Config) func(State) State {
//	            return func(s State) State { s.Config = cfg; return s }
//	        },
//	        getConfig,
//	    ),
//	)
//
//go:inline
func ApS[S1, S2, T any](
	setter func(T) func(S1) S2,
	fa ReaderIOResult[T],
) Operator[S1, S2] {
	return apply.ApS(
		Ap,
		Map,
		setter,
		fa,
	)
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
//
//	userLens := lens.MakeLens(
//	    func(s State) User { return s.User },
//	    func(s State, u User) State { s.User = u; return s },
//	)
//
//	getUser := func(ctx context.Context) ioeither.IOEither[error, User] {
//	    return ioeither.TryCatch(func() (User, error) {
//	        return fetchUser(ctx)
//	    })
//	}
//	result := F.Pipe2(
//	    readerioeither.Of(State{}),
//	    readerioeither.ApSL(userLens, getUser),
//	)
//
//go:inline
func ApSL[S, T any](
	lens L.Lens[S, T],
	fa ReaderIOResult[T],
) Operator[S, S] {
	return ApS(lens.Set, fa)
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
//	    Config Config
//	}
//
//	userLens := lens.MakeLens(
//	    func(s State) User { return s.User },
//	    func(s State, u User) State { s.User = u; return s },
//	)
//
//	result := F.Pipe2(
//	    readerioeither.Do(State{}),
//	    readerioeither.BindL(userLens, func(user User) readerioeither.ReaderIOResult[User] {
//	        return func(ctx context.Context) ioeither.IOEither[error, User] {
//	            return ioeither.TryCatch(func() (User, error) {
//	                return fetchUser(ctx)
//	            })
//	        }
//	    }),
//	)
//
//go:inline
func BindL[S, T any](
	lens L.Lens[S, T],
	f Kleisli[T, T],
) Operator[S, S] {
	return RIOR.BindL[context.Context](lens, f)
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
//	    Config Config
//	}
//
//	userLens := lens.MakeLens(
//	    func(s State) User { return s.User },
//	    func(s State, u User) State { s.User = u; return s },
//	)
//
//	result := F.Pipe2(
//	    readerioeither.Do(State{User: User{Name: "Alice"}}),
//	    readerioeither.LetL(userLens, func(user User) User {
//	        user.Name = "Bob"
//	        return user
//	    }),
//	)
//
//go:inline
func LetL[S, T any](
	lens L.Lens[S, T],
	f func(T) T,
) Operator[S, S] {
	return RIOR.LetL[context.Context](lens, f)
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
//	userLens := lens.MakeLens(
//	    func(s State) User { return s.User },
//	    func(s State, u User) State { s.User = u; return s },
//	)
//
//	newUser := User{Name: "Bob", ID: 123}
//	result := F.Pipe2(
//	    readerioeither.Do(State{}),
//	    readerioeither.LetToL(userLens, newUser),
//	)
//
//go:inline
func LetToL[S, T any](
	lens L.Lens[S, T],
	b T,
) Operator[S, S] {
	return RIOR.LetToL[context.Context](lens, b)
}
