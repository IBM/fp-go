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

	"github.com/IBM/fp-go/v2/context/readerio"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/apply"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/reader"
	RIOR "github.com/IBM/fp-go/v2/readerioresult"
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
	return RIOR.Bind(setter, WithContextK(f))
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

//go:inline
func BindToP[S1, T any](
	setter Prism[S1, T],
) Operator[T, S1] {
	return BindTo(setter.ReverseGet)
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
	lens Lens[S, T],
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
	lens Lens[S, T],
	f Kleisli[T, T],
) Operator[S, S] {
	return RIOR.BindL(lens, WithContextK(f))
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
	lens Lens[S, T],
	f Endomorphism[T],
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
	lens Lens[S, T],
	b T,
) Operator[S, S] {
	return RIOR.LetToL[context.Context](lens, b)
}

// BindIOEitherK is a variant of Bind that works with IOEither computations.
// It lifts an IOEither Kleisli arrow into the ReaderIOResult context (with context.Context as environment).
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - f: An IOEither Kleisli arrow (S1 -> IOEither[error, T])
//
//go:inline
func BindIOEitherK[S1, S2, T any](
	setter func(T) func(S1) S2,
	f ioresult.Kleisli[S1, T],
) Operator[S1, S2] {
	return Bind(setter, F.Flow2(f, FromIOEither[T]))
}

// BindIOResultK is a variant of Bind that works with IOResult computations.
// This is an alias for BindIOEitherK for consistency with the Result naming convention.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - f: An IOResult Kleisli arrow (S1 -> IOResult[T])
//
//go:inline
func BindIOResultK[S1, S2, T any](
	setter func(T) func(S1) S2,
	f ioresult.Kleisli[S1, T],
) Operator[S1, S2] {
	return Bind(setter, F.Flow2(f, FromIOResult[T]))
}

// BindIOK is a variant of Bind that works with IO computations.
// It lifts an IO Kleisli arrow into the ReaderIOResult context.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - f: An IO Kleisli arrow (S1 -> IO[T])
//
//go:inline
func BindIOK[S1, S2, T any](
	setter func(T) func(S1) S2,
	f io.Kleisli[S1, T],
) Operator[S1, S2] {
	return Bind(setter, F.Flow2(f, FromIO[T]))
}

// BindReaderK is a variant of Bind that works with Reader computations.
// It lifts a Reader Kleisli arrow (with context.Context) into the ReaderIOResult context.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - f: A Reader Kleisli arrow (S1 -> Reader[context.Context, T])
//
//go:inline
func BindReaderK[S1, S2, T any](
	setter func(T) func(S1) S2,
	f reader.Kleisli[context.Context, S1, T],
) Operator[S1, S2] {
	return Bind(setter, F.Flow2(f, FromReader[T]))
}

// BindReaderIOK is a variant of Bind that works with ReaderIO computations.
// It lifts a ReaderIO Kleisli arrow (with context.Context) into the ReaderIOResult context.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - f: A ReaderIO Kleisli arrow (S1 -> ReaderIO[context.Context, T])
//
//go:inline
func BindReaderIOK[S1, S2, T any](
	setter func(T) func(S1) S2,
	f readerio.Kleisli[S1, T],
) Operator[S1, S2] {
	return Bind(setter, F.Flow2(f, FromReaderIO[T]))
}

// BindEitherK is a variant of Bind that works with Either (Result) computations.
// It lifts an Either Kleisli arrow into the ReaderIOResult context.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - f: An Either Kleisli arrow (S1 -> Either[error, T])
//
//go:inline
func BindEitherK[S1, S2, T any](
	setter func(T) func(S1) S2,
	f result.Kleisli[S1, T],
) Operator[S1, S2] {
	return Bind(setter, F.Flow2(f, FromEither[T]))
}

// BindResultK is a variant of Bind that works with Result computations.
// This is an alias for BindEitherK for consistency with the Result naming convention.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - f: A Result Kleisli arrow (S1 -> Result[T])
//
//go:inline
func BindResultK[S1, S2, T any](
	setter func(T) func(S1) S2,
	f result.Kleisli[S1, T],
) Operator[S1, S2] {
	return Bind(setter, F.Flow2(f, FromResult[T]))
}

// BindIOEitherKL is a lens-based variant of BindIOEitherK.
// It combines a lens with an IOEither Kleisli arrow, focusing on a specific field
// within the state structure.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - f: An IOEither Kleisli arrow (T -> IOEither[error, T])
//
//go:inline
func BindIOEitherKL[S, T any](
	lens Lens[S, T],
	f ioresult.Kleisli[T, T],
) Operator[S, S] {
	return BindL(lens, F.Flow2(f, FromIOEither[T]))
}

// BindIOResultKL is a lens-based variant of BindIOResultK.
// This is an alias for BindIOEitherKL for consistency with the Result naming convention.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - f: An IOResult Kleisli arrow (T -> IOResult[T])
//
//go:inline
func BindIOResultKL[S, T any](
	lens Lens[S, T],
	f ioresult.Kleisli[T, T],
) Operator[S, S] {
	return BindL(lens, F.Flow2(f, FromIOEither[T]))
}

// BindIOKL is a lens-based variant of BindIOK.
// It combines a lens with an IO Kleisli arrow, focusing on a specific field
// within the state structure.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - f: An IO Kleisli arrow (T -> IO[T])
//
//go:inline
func BindIOKL[S, T any](
	lens Lens[S, T],
	f io.Kleisli[T, T],
) Operator[S, S] {
	return BindL(lens, F.Flow2(f, FromIO[T]))
}

// BindReaderKL is a lens-based variant of BindReaderK.
// It combines a lens with a Reader Kleisli arrow (with context.Context), focusing on a specific field
// within the state structure.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - f: A Reader Kleisli arrow (T -> Reader[context.Context, T])
//
//go:inline
func BindReaderKL[S, T any](
	lens Lens[S, T],
	f reader.Kleisli[context.Context, T, T],
) Operator[S, S] {
	return BindL(lens, F.Flow2(f, FromReader[T]))
}

// BindReaderIOKL is a lens-based variant of BindReaderIOK.
// It combines a lens with a ReaderIO Kleisli arrow (with context.Context), focusing on a specific field
// within the state structure.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - f: A ReaderIO Kleisli arrow (T -> ReaderIO[context.Context, T])
//
//go:inline
func BindReaderIOKL[S, T any](
	lens Lens[S, T],
	f readerio.Kleisli[T, T],
) Operator[S, S] {
	return BindL(lens, F.Flow2(f, FromReaderIO[T]))
}

// ApIOEitherS is an applicative variant that works with IOEither values.
// Unlike BindIOEitherK, this uses applicative composition (ApS) instead of monadic
// composition (Bind), allowing independent computations to be combined.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - fa: An IOEither value
//
//go:inline
func ApIOEitherS[S1, S2, T any](
	setter func(T) func(S1) S2,
	fa IOResult[T],
) Operator[S1, S2] {
	return F.Bind2nd(F.Flow2[ReaderIOResult[S1], ioresult.Operator[S1, S2]], ioeither.ApS(setter, fa))
}

// ApIOResultS is an applicative variant that works with IOResult values.
// This is an alias for ApIOEitherS for consistency with the Result naming convention.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - fa: An IOResult value
//
//go:inline
func ApIOResultS[S1, S2, T any](
	setter func(T) func(S1) S2,
	fa IOResult[T],
) Operator[S1, S2] {
	return F.Bind2nd(F.Flow2[ReaderIOResult[S1], ioresult.Operator[S1, S2]], ioeither.ApS(setter, fa))
}

// ApIOS is an applicative variant that works with IO values.
// It lifts an IO value into the ReaderIOResult context using applicative composition.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - fa: An IO value
//
//go:inline
func ApIOS[S1, S2, T any](
	setter func(T) func(S1) S2,
	fa IO[T],
) Operator[S1, S2] {
	return ApS(setter, FromIO(fa))
}

// ApReaderS is an applicative variant that works with Reader values.
// It lifts a Reader value (with context.Context) into the ReaderIOResult context using applicative composition.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - fa: A Reader value
//
//go:inline
func ApReaderS[S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Reader[context.Context, T],
) Operator[S1, S2] {
	return ApS(setter, FromReader(fa))
}

// ApReaderIOS is an applicative variant that works with ReaderIO values.
// It lifts a ReaderIO value (with context.Context) into the ReaderIOResult context using applicative composition.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - fa: A ReaderIO value
//
//go:inline
func ApReaderIOS[S1, S2, T any](
	setter func(T) func(S1) S2,
	fa ReaderIO[T],
) Operator[S1, S2] {
	return ApS(setter, FromReaderIO(fa))
}

// ApEitherS is an applicative variant that works with Either (Result) values.
// It lifts an Either value into the ReaderIOResult context using applicative composition.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - fa: An Either value
//
//go:inline
func ApEitherS[S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Result[T],
) Operator[S1, S2] {
	return ApS(setter, FromEither(fa))
}

// ApResultS is an applicative variant that works with Result values.
// This is an alias for ApEitherS for consistency with the Result naming convention.
//
// Parameters:
//   - setter: Updates state from S1 to S2 using result T
//   - fa: A Result value
//
//go:inline
func ApResultS[S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Result[T],
) Operator[S1, S2] {
	return ApS(setter, FromResult(fa))
}

// ApIOEitherSL is a lens-based variant of ApIOEitherS.
// It combines a lens with an IOEither value using applicative composition.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - fa: An IOEither value
//
//go:inline
func ApIOEitherSL[S, T any](
	lens Lens[S, T],
	fa IOResult[T],
) Operator[S, S] {
	return F.Bind2nd(F.Flow2[ReaderIOResult[S], ioresult.Operator[S, S]], ioresult.ApSL(lens, fa))
}

// ApIOResultSL is a lens-based variant of ApIOResultS.
// This is an alias for ApIOEitherSL for consistency with the Result naming convention.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - fa: An IOResult value
//
//go:inline
func ApIOResultSL[S, T any](
	lens Lens[S, T],
	fa IOResult[T],
) Operator[S, S] {
	return F.Bind2nd(F.Flow2[ReaderIOResult[S], ioresult.Operator[S, S]], ioresult.ApSL(lens, fa))
}

// ApIOSL is a lens-based variant of ApIOS.
// It combines a lens with an IO value using applicative composition.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - fa: An IO value
//
//go:inline
func ApIOSL[S, T any](
	lens Lens[S, T],
	fa IO[T],
) Operator[S, S] {
	return ApSL(lens, FromIO(fa))
}

// ApReaderSL is a lens-based variant of ApReaderS.
// It combines a lens with a Reader value (with context.Context) using applicative composition.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - fa: A Reader value
//
//go:inline
func ApReaderSL[S, T any](
	lens Lens[S, T],
	fa Reader[context.Context, T],
) Operator[S, S] {
	return ApSL(lens, FromReader(fa))
}

// ApReaderIOSL is a lens-based variant of ApReaderIOS.
// It combines a lens with a ReaderIO value (with context.Context) using applicative composition.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - fa: A ReaderIO value
//
//go:inline
func ApReaderIOSL[S, T any](
	lens Lens[S, T],
	fa ReaderIO[T],
) Operator[S, S] {
	return ApSL(lens, FromReaderIO(fa))
}

// ApEitherSL is a lens-based variant of ApEitherS.
// It combines a lens with an Either value using applicative composition.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - fa: An Either value
//
//go:inline
func ApEitherSL[S, T any](
	lens Lens[S, T],
	fa Result[T],
) Operator[S, S] {
	return ApSL(lens, FromEither(fa))
}

// ApResultSL is a lens-based variant of ApResultS.
// This is an alias for ApEitherSL for consistency with the Result naming convention.
//
// Parameters:
//   - lens: A lens focusing on field T within state S
//   - fa: A Result value
//
//go:inline
func ApResultSL[S, T any](
	lens Lens[S, T],
	fa Result[T],
) Operator[S, S] {
	return ApSL(lens, FromResult(fa))
}
