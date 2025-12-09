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
	"context"

	RR "github.com/IBM/fp-go/v2/idiomatic/readerresult"
	"github.com/IBM/fp-go/v2/idiomatic/result"
	C "github.com/IBM/fp-go/v2/internal/chain"
	L "github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/reader"
	RES "github.com/IBM/fp-go/v2/result"
)

// Do initializes a do-notation context with an empty state.
//
// This is the starting point for do-notation style composition, which allows
// imperative-style sequencing of ReaderResult computations while maintaining
// functional purity.
//
// Type Parameters:
//   - S: The state type
//
// Parameters:
//   - empty: The initial empty state
//
// Returns:
//   - A ReaderResult[S] containing the initial state
//
// Example:
//
//	type State struct {
//	    User  User
//	    Posts []Post
//	}
//
//	result := F.Pipe2(
//	    readerresult.Do(State{}),
//	    readerresult.Bind(
//	        func(u User) func(State) State {
//	            return func(s State) State { s.User = u; return s }
//	        },
//	        func(s State) readerresult.ReaderResult[User] {
//	            return getUser(42)
//	        },
//	    ),
//	    readerresult.Bind(
//	        func(posts []Post) func(State) State {
//	            return func(s State) State { s.Posts = posts; return s }
//	        },
//	        func(s State) readerresult.ReaderResult[[]Post] {
//	            return getPosts(s.User.ID)
//	        },
//	    ),
//	)
//
//go:inline
func Do[S any](
	empty S,
) ReaderResult[S] {
	return RR.Do[context.Context](empty)
}

// Bind sequences a ReaderResult computation and updates the state with its result.
//
// This is the core operation for do-notation, allowing you to chain computations
// where each step can depend on the accumulated state and update it with new values.
//
// Type Parameters:
//   - S1: The input state type
//   - S2: The output state type
//   - T: The type of value produced by the computation
//
// Parameters:
//   - setter: A function that takes the computation result and returns a state updater
//   - f: A Kleisli arrow that produces the next computation based on current state
//
// Returns:
//   - An Operator that transforms ReaderResult[S1] to ReaderResult[S2]
//
// Example:
//
//	readerresult.Bind(
//	    func(user User) func(State) State {
//	        return func(s State) State { s.User = user; return s }
//	    },
//	    func(s State) readerresult.ReaderResult[User] {
//	        return getUser(s.UserID)
//	    },
//	)
//
//go:inline
func Bind[S1, S2, T any](
	setter func(T) func(S1) S2,
	f Kleisli[S1, T],
) Operator[S1, S2] {
	return C.Bind(
		Chain[S1, S2],
		Map[T, S2],
		setter,
		f,
	)
}

// Let attaches the result of a pure computation to a state.
//
// Unlike Bind, Let works with pure functions (not ReaderResult computations).
// This is useful for deriving values from the current state without performing
// any effects.
//
// Type Parameters:
//   - S1: The input state type
//   - S2: The output state type
//   - T: The type of value computed
//
// Parameters:
//   - setter: A function that takes the computed value and returns a state updater
//   - f: A pure function that computes a value from the current state
//
// Returns:
//   - An Operator that transforms ReaderResult[S1] to ReaderResult[S2]
//
// Example:
//
//	readerresult.Let(
//	    func(fullName string) func(State) State {
//	        return func(s State) State { s.FullName = fullName; return s }
//	    },
//	    func(s State) string {
//	        return s.FirstName + " " + s.LastName
//	    },
//	)
//
//go:inline
func Let[S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) Operator[S1, S2] {
	return RR.Let[context.Context](setter, f)
}

// LetTo attaches a constant value to a state.
//
// This is a simplified version of Let for when you want to add a constant
// value to the state without computing it.
//
// Type Parameters:
//   - S1: The input state type
//   - S2: The output state type
//   - T: The type of the constant value
//
// Parameters:
//   - setter: A function that takes the constant and returns a state updater
//   - b: The constant value to attach
//
// Returns:
//   - An Operator that transforms ReaderResult[S1] to ReaderResult[S2]
//
// Example:
//
//	readerresult.LetTo(
//	    func(status string) func(State) State {
//	        return func(s State) State { s.Status = status; return s }
//	    },
//	    "active",
//	)
//
//go:inline
func LetTo[S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) Operator[S1, S2] {
	return RR.LetTo[context.Context](setter, b)
}

// BindTo initializes do-notation by binding a value to a state.
//
// This is typically used as the first operation after a computation to
// start building up a state structure.
//
// Type Parameters:
//   - S1: The state type to create
//   - T: The type of the initial value
//
// Parameters:
//   - setter: A function that creates the initial state from a value
//
// Returns:
//   - An Operator that transforms ReaderResult[T] to ReaderResult[S1]
//
// Example:
//
//	type State struct {
//	    User User
//	}
//
//	result := F.Pipe1(
//	    getUser(42),
//	    readerresult.BindTo(func(u User) State {
//	        return State{User: u}
//	    }),
//	)
//
//go:inline
func BindTo[S1, T any](
	setter func(T) S1,
) Operator[T, S1] {
	return RR.BindTo[context.Context](setter)
}

//go:inline
func ApS[S1, S2, T any](
	setter func(T) func(S1) S2,
	fa ReaderResult[T],
) Operator[S1, S2] {
	return RR.ApS[context.Context](setter, fa)
}

//go:inline
func ApSL[S, T any](
	lens L.Lens[S, T],
	fa ReaderResult[T],
) Operator[S, S] {
	return ApSL(lens, fa)
}

//go:inline
func BindL[S, T any](
	lens L.Lens[S, T],
	f Kleisli[T, T],
) Operator[S, S] {
	return RR.BindL(lens, f)
}

//go:inline
func LetL[S, T any](
	lens L.Lens[S, T],
	f Endomorphism[T],
) Operator[S, S] {
	return RR.LetL[context.Context](lens, f)
}

//go:inline
func LetToL[S, T any](
	lens L.Lens[S, T],
	b T,
) Operator[S, S] {
	return RR.LetToL[context.Context](lens, b)
}

//go:inline
func BindReaderK[S1, S2, T any](
	setter func(T) func(S1) S2,
	f reader.Kleisli[context.Context, S1, T],
) Operator[S1, S2] {
	return RR.BindReaderK(setter, f)
}

//go:inline
func BindEitherK[S1, S2, T any](
	setter func(T) func(S1) S2,
	f RES.Kleisli[S1, T],
) Operator[S1, S2] {
	return RR.BindEitherK[context.Context](setter, f)
}

//go:inline
func BindResultK[S1, S2, T any](
	setter func(T) func(S1) S2,
	f result.Kleisli[S1, T],
) Operator[S1, S2] {
	return RR.BindResultK[context.Context](setter, f)
}

//go:inline
func BindToReader[
	S1, T any](
	setter func(T) S1,
) func(Reader[context.Context, T]) ReaderResult[S1] {
	return RR.BindToReader[context.Context](setter)
}

//go:inline
func BindToEither[
	S1, T any](
	setter func(T) S1,
) func(Result[T]) ReaderResult[S1] {
	return RR.BindToEither[context.Context](setter)
}

//go:inline
func BindToResult[
	S1, T any](
	setter func(T) S1,
) func(T, error) ReaderResult[S1] {
	return RR.BindToResult[context.Context](setter)
}

//go:inline
func ApReaderS[
	S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Reader[context.Context, T],
) Operator[S1, S2] {
	return RR.ApReaderS(setter, fa)
}

//go:inline
func ApResultS[
	S1, S2, T any](
	setter func(T) func(S1) S2,
) func(T, error) Operator[S1, S2] {
	return RR.ApResultS[context.Context](setter)
}

//go:inline
func ApEitherS[
	S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Result[T],
) Operator[S1, S2] {
	return RR.ApEitherS[context.Context](setter, fa)
}
