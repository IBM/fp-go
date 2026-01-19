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

package readerreaderioeither

import (
	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/apply"
	"github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
)

// Do creates an empty context of type [S] to be used with the [Bind] operation.
// This is the starting point for do-notation style composition with two reader contexts.
//
// Example:
//
//	type State struct {
//	    User   User
//	    Posts  []Post
//	}
//	type OuterEnv struct {
//	    Database string
//	}
//	type InnerEnv struct {
//	    UserRepo UserRepository
//	    PostRepo PostRepository
//	}
//	result := readerreaderioeither.Do[OuterEnv, InnerEnv, error](State{})
//
//go:inline
func Do[R, C, E, S any](
	empty S,
) ReaderReaderIOEither[R, C, E, S] {
	return Of[R, C, E](empty)
}

// Bind attaches the result of a computation to a context [S1] to produce a context [S2].
// This enables sequential composition where each step can depend on the results of previous steps
// and access both the outer (R) and inner (C) reader environments.
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
//	type OuterEnv struct {
//	    Database string
//	}
//	type InnerEnv struct {
//	    UserRepo UserRepository
//	    PostRepo PostRepository
//	}
//
//	result := F.Pipe2(
//	    readerreaderioeither.Do[OuterEnv, InnerEnv, error](State{}),
//	    readerreaderioeither.Bind(
//	        func(user User) func(State) State {
//	            return func(s State) State { s.User = user; return s }
//	        },
//	        func(s State) readerreaderioeither.ReaderReaderIOEither[OuterEnv, InnerEnv, error, User] {
//	            return func(outer OuterEnv) readerioeither.ReaderIOEither[InnerEnv, error, User] {
//	                return readerioeither.Asks(func(inner InnerEnv) ioeither.IOEither[error, User] {
//	                    return inner.UserRepo.FindUser(outer.Database)
//	                })
//	            }
//	        },
//	    ),
//	    readerreaderioeither.Bind(
//	        func(posts []Post) func(State) State {
//	            return func(s State) State { s.Posts = posts; return s }
//	        },
//	        func(s State) readerreaderioeither.ReaderReaderIOEither[OuterEnv, InnerEnv, error, []Post] {
//	            return func(outer OuterEnv) readerioeither.ReaderIOEither[InnerEnv, error, []Post] {
//	                return readerioeither.Asks(func(inner InnerEnv) ioeither.IOEither[error, []Post] {
//	                    return inner.PostRepo.FindPostsByUser(outer.Database, s.User.ID)
//	                })
//	            }
//	        },
//	    ),
//	)
//
//go:inline
func Bind[R, C, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) ReaderReaderIOEither[R, C, E, T],
) Operator[R, C, E, S1, S2] {
	return chain.Bind(
		Chain[R, C, E, S1, S2],
		Map[R, C, E, T, S2],
		setter,
		f,
	)
}

// Let attaches a pure computation result to a context [S1] to produce a context [S2].
// Unlike [Bind], the computation function f is pure (doesn't perform effects).
//
// Example:
//
//	readerreaderioeither.Let(
//	    func(fullName string) func(State) State {
//	        return func(s State) State { s.FullName = fullName; return s }
//	    },
//	    func(s State) string {
//	        return s.FirstName + " " + s.LastName
//	    },
//	)
//
//go:inline
func Let[R, C, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) Operator[R, C, E, S1, S2] {
	return functor.Let(
		Map[R, C, E, S1, S2],
		setter,
		f,
	)
}

// LetTo attaches a constant value to a context [S1] to produce a context [S2].
//
// Example:
//
//	readerreaderioeither.LetTo(
//	    func(status string) func(State) State {
//	        return func(s State) State { s.Status = status; return s }
//	    },
//	    "active",
//	)
//
//go:inline
func LetTo[R, C, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) Operator[R, C, E, S1, S2] {
	return functor.LetTo(
		Map[R, C, E, S1, S2],
		setter,
		b,
	)
}

// BindTo wraps a value of type T into a context S1 using the provided setter function.
// This is typically used as the first operation after [Do] to initialize the context.
//
// Example:
//
//	F.Pipe1(
//	    readerreaderioeither.Of[OuterEnv, InnerEnv, error](42),
//	    readerreaderioeither.BindTo(func(n int) State { return State{Count: n} }),
//	)
//
//go:inline
func BindTo[R, C, E, S1, T any](
	setter func(T) S1,
) Operator[R, C, E, T, S1] {
	return chain.BindTo(
		Map[R, C, E, T, S1],
		setter,
	)
}

// ApS applies a computation in parallel (applicative style) and attaches its result to the context.
// Unlike [Bind], this doesn't allow the computation to depend on the current context state.
//
// Example:
//
//	readerreaderioeither.ApS(
//	    func(count int) func(State) State {
//	        return func(s State) State { s.Count = count; return s }
//	    },
//	    getCount, // ReaderReaderIOEither[OuterEnv, InnerEnv, error, int]
//	)
//
//go:inline
func ApS[R, C, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa ReaderReaderIOEither[R, C, E, T],
) Operator[R, C, E, S1, S2] {
	return apply.ApS(
		Ap[S2, R, C, E, T],
		Map[R, C, E, S1, func(T) S2],
		setter,
		fa,
	)
}

// ApSL is a lens-based version of [ApS] that uses a lens to focus on a specific field in the context.
//
//go:inline
func ApSL[R, C, E, S, T any](
	lens Lens[S, T],
	fa ReaderReaderIOEither[R, C, E, T],
) Operator[R, C, E, S, S] {
	return ApS(lens.Set, fa)
}

// BindL is a lens-based version of [Bind] that uses a lens to focus on a specific field in the context.
//
//go:inline
func BindL[R, C, E, S, T any](
	lens Lens[S, T],
	f func(T) ReaderReaderIOEither[R, C, E, T],
) Operator[R, C, E, S, S] {
	return Bind(lens.Set, F.Flow2(lens.Get, f))
}

// LetL is a lens-based version of [Let] that uses a lens to focus on a specific field in the context.
//
//go:inline
func LetL[R, C, E, S, T any](
	lens Lens[S, T],
	f func(T) T,
) Operator[R, C, E, S, S] {
	return Let[R, C, E](lens.Set, F.Flow2(lens.Get, f))
}

// LetToL is a lens-based version of [LetTo] that uses a lens to focus on a specific field in the context.
//
//go:inline
func LetToL[R, C, E, S, T any](
	lens Lens[S, T],
	b T,
) Operator[R, C, E, S, S] {
	return LetTo[R, C, E](lens.Set, b)
}

// BindIOEitherK binds a computation that returns an IOEither to the context.
// The Kleisli function is automatically lifted into ReaderReaderIOEither.
//
//go:inline
func BindIOEitherK[R, C, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f ioeither.Kleisli[E, S1, T],
) Operator[R, C, E, S1, S2] {
	return Bind(setter, F.Flow2(f, FromIOEither[R, C, E, T]))
}

// BindIOK binds a computation that returns an IO to the context.
// The Kleisli function is automatically lifted into ReaderReaderIOEither.
//
//go:inline
func BindIOK[R, C, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f io.Kleisli[S1, T],
) Operator[R, C, E, S1, S2] {
	return Bind(setter, F.Flow2(f, FromIO[R, C, E, T]))
}

// BindReaderK binds a computation that returns a Reader to the context.
// The Kleisli function is automatically lifted into ReaderReaderIOEither.
//
//go:inline
func BindReaderK[C, E, R, S1, S2, T any](
	setter func(T) func(S1) S2,
	f reader.Kleisli[R, S1, T],
) Operator[R, C, E, S1, S2] {
	return Bind(setter, F.Flow2(f, FromReader[C, E, R, T]))
}

// BindReaderIOK binds a computation that returns a ReaderIO to the context.
// The Kleisli function is automatically lifted into ReaderReaderIOEither.
//
//go:inline
func BindReaderIOK[C, E, R, S1, S2, T any](
	setter func(T) func(S1) S2,
	f readerio.Kleisli[R, S1, T],
) Operator[R, C, E, S1, S2] {
	return Bind(setter, F.Flow2(f, FromReaderIO[C, E, R, T]))
}

// BindEitherK binds a computation that returns an Either to the context.
// The Kleisli function is automatically lifted into ReaderReaderIOEither.
//
//go:inline
func BindEitherK[R, C, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f either.Kleisli[E, S1, T],
) Operator[R, C, E, S1, S2] {
	return Bind(setter, F.Flow2(f, FromEither[R, C, E, T]))
}

// BindIOEitherKL is a lens-based version of [BindIOEitherK].
//
//go:inline
func BindIOEitherKL[R, C, E, S, T any](
	lens Lens[S, T],
	f ioeither.Kleisli[E, T, T],
) Operator[R, C, E, S, S] {
	return BindL(lens, F.Flow2(f, FromIOEither[R, C, E, T]))
}

// BindIOKL is a lens-based version of [BindIOK].
//
//go:inline
func BindIOKL[R, C, E, S, T any](
	lens Lens[S, T],
	f io.Kleisli[T, T],
) Operator[R, C, E, S, S] {
	return BindL(lens, F.Flow2(f, FromIO[R, C, E, T]))
}

// BindReaderKL is a lens-based version of [BindReaderK].
//
//go:inline
func BindReaderKL[C, E, R, S, T any](
	lens Lens[S, T],
	f reader.Kleisli[R, T, T],
) Operator[R, C, E, S, S] {
	return BindL(lens, F.Flow2(f, FromReader[C, E, R, T]))
}

// BindReaderIOKL is a lens-based version of [BindReaderIOK].
//
//go:inline
func BindReaderIOKL[C, E, R, S, T any](
	lens Lens[S, T],
	f readerio.Kleisli[R, T, T],
) Operator[R, C, E, S, S] {
	return BindL(lens, F.Flow2(f, FromReaderIO[C, E, R, T]))
}

// ApIOEitherS applies an IOEither computation and attaches its result to the context.
//
//go:inline
func ApIOEitherS[R, C, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa IOEither[E, T],
) Operator[R, C, E, S1, S2] {
	return ApS(setter, FromIOEither[R, C](fa))
}

// ApIOS applies an IO computation and attaches its result to the context.
//
//go:inline
func ApIOS[R, C, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa IO[T],
) Operator[R, C, E, S1, S2] {
	return ApS(setter, FromIO[R, C, E](fa))
}

// ApReaderS applies a Reader computation and attaches its result to the context.
//
//go:inline
func ApReaderS[C, E, R, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Reader[R, T],
) Operator[R, C, E, S1, S2] {
	return ApS(setter, FromReader[C, E](fa))
}

// ApReaderIOS applies a ReaderIO computation and attaches its result to the context.
//
//go:inline
func ApReaderIOS[C, E, R, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa ReaderIO[R, T],
) Operator[R, C, E, S1, S2] {
	return ApS(setter, FromReaderIO[C, E](fa))
}

// ApEitherS applies an Either computation and attaches its result to the context.
//
//go:inline
func ApEitherS[R, C, E, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Either[E, T],
) Operator[R, C, E, S1, S2] {
	return ApS(setter, FromEither[R, C](fa))
}

// ApIOEitherSL is a lens-based version of [ApIOEitherS].
//
//go:inline
func ApIOEitherSL[R, C, E, S, T any](
	lens Lens[S, T],
	fa IOEither[E, T],
) Operator[R, C, E, S, S] {
	return ApIOEitherS[R, C](lens.Set, fa)
}

// ApIOSL is a lens-based version of [ApIOS].
//
//go:inline
func ApIOSL[R, C, E, S, T any](
	lens Lens[S, T],
	fa IO[T],
) Operator[R, C, E, S, S] {
	return ApSL(lens, FromIO[R, C, E](fa))
}

// ApReaderSL is a lens-based version of [ApReaderS].
//
//go:inline
func ApReaderSL[C, E, R, S, T any](
	lens Lens[S, T],
	fa Reader[R, T],
) Operator[R, C, E, S, S] {
	return ApReaderS[C, E](lens.Set, fa)
}

// ApReaderIOSL is a lens-based version of [ApReaderIOS].
//
//go:inline
func ApReaderIOSL[C, E, R, S, T any](
	lens Lens[S, T],
	fa ReaderIO[R, T],
) Operator[R, C, E, S, S] {
	return ApReaderIOS[C, E](lens.Set, fa)
}

// ApEitherSL is a lens-based version of [ApEitherS].
//
//go:inline
func ApEitherSL[R, C, E, S, T any](
	lens Lens[S, T],
	fa Either[E, T],
) Operator[R, C, E, S, S] {
	return ApEitherS[R, C](lens.Set, fa)
}
