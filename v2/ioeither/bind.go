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

package ioeither

import (
	"github.com/IBM/fp-go/v2/internal/apply"
	"github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/functor"
)

// Do creates an empty context of type [S] to be used with the [Bind] operation.
// This is the starting point for do-notation style composition.
//
// Example:
//
//	type State struct {
//	    User  User
//	    Posts []Post
//	}
//	result := ioeither.Do[error](State{})
func Do[E, S any](
	empty S,
) IOEither[E, S] {
	return Of[E](empty)
}

// Bind attaches the result of a computation to a context [S1] to produce a context [S2].
// This enables sequential composition where each step can depend on the results of previous steps.
//
// The setter function takes the result of the computation and returns a function that
// updates the context from S1 to S2.
//
// Example:
//
//	type State struct {
//	    User  User
//	    Posts []Post
//	}
//
//	result := F.Pipe2(
//	    ioeither.Do[error](State{}),
//	    ioeither.Bind(
//	        func(user User) func(State) State {
//	            return func(s State) State { s.User = user; return s }
//	        },
//	        func(s State) ioeither.IOEither[error, User] {
//	            return ioeither.TryCatch(func() (User, error) {
//	                return fetchUser()
//	            })
//	        },
//	    ),
//	    ioeither.Bind(
//	        func(posts []Post) func(State) State {
//	            return func(s State) State { s.Posts = posts; return s }
//	        },
//	        func(s State) ioeither.IOEither[error, []Post] {
//	            // This can access s.User from the previous step
//	            return ioeither.TryCatch(func() ([]Post, error) {
//	                return fetchPostsForUser(s.User.ID)
//	            })
//	        },
//	    ),
//	)
func Bind[E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) IOEither[E, T],
) Operator[E, S1, S2] {
	return chain.Bind(
		Chain[E, S1, S2],
		Map[E, T, S2],
		setter,
		f,
	)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
func Let[E, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) Operator[E, S1, S2] {
	return functor.Let(
		Map[E, S1, S2],
		setter,
		f,
	)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
func LetTo[E, S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) Operator[E, S1, S2] {
	return functor.LetTo(
		Map[E, S1, S2],
		setter,
		b,
	)
}

// BindTo initializes a new state [S1] from a value [T]
func BindTo[E, S1, T any](
	setter func(T) S1,
) Operator[E, T, S1] {
	return chain.BindTo(
		Map[E, T, S1],
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
//	    User  User
//	    Posts []Post
//	}
//
//	// These operations are independent and can be combined with ApS
//	getUser := ioeither.Right[error](User{ID: 1, Name: "Alice"})
//	getPosts := ioeither.Right[error]([]Post{{ID: 1, Title: "Hello"}})
//
//	result := F.Pipe2(
//	    ioeither.Do[error](State{}),
//	    ioeither.ApS(
//	        func(user User) func(State) State {
//	            return func(s State) State { s.User = user; return s }
//	        },
//	        getUser,
//	    ),
//	    ioeither.ApS(
//	        func(posts []Post) func(State) State {
//	            return func(s State) State { s.Posts = posts; return s }
//	        },
//	        getPosts,
//	    ),
//	)
func ApS[E, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa IOEither[E, T],
) Operator[E, S1, S2] {
	return apply.ApS(
		Ap[S2, E, T],
		Map[E, S1, func(T) S2],
		setter,
		fa,
	)
}
