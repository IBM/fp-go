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
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/apply"
	"github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/functor"
	L "github.com/IBM/fp-go/v2/optics/lens"
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
	f Kleisli[E, S1, T],
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

// ApSL attaches a value to a context using a lens-based setter.
// This is a convenience function that combines ApS with a lens, allowing you to use
// optics to update nested structures in a more composable way.
//
// The lens parameter provides both the getter and setter for a field within the structure S.
// This eliminates the need to manually write setter functions.
//
// Example:
//
//	type Config struct {
//	    Host string
//	    Port int
//	}
//
//	portLens := lens.MakeLens(
//	    func(c Config) int { return c.Port },
//	    func(c Config, p int) Config { c.Port = p; return c },
//	)
//
//	result := F.Pipe2(
//	    ioeither.Of[error](Config{Host: "localhost"}),
//	    ioeither.ApSL(portLens, ioeither.Of[error](8080)),
//	)
func ApSL[E, S, T any](
	lens L.Lens[S, T],
	fa IOEither[E, T],
) Operator[E, S, S] {
	return ApS(lens.Set, fa)
}

// BindL attaches the result of a computation to a context using a lens-based setter.
// This is a convenience function that combines Bind with a lens, allowing you to use
// optics to update nested structures based on their current values.
//
// The lens parameter provides both the getter and setter for a field within the structure S.
// The computation function f receives the current value of the focused field and returns
// an IOEither that produces the new value.
//
// Example:
//
//	type Counter struct {
//	    Value int
//	}
//
//	valueLens := lens.MakeLens(
//	    func(c Counter) int { return c.Value },
//	    func(c Counter, v int) Counter { c.Value = v; return c },
//	)
//
//	increment := func(v int) ioeither.IOEither[error, int] {
//	    return ioeither.TryCatch(func() (int, error) {
//	        if v >= 100 {
//	            return 0, errors.New("overflow")
//	        }
//	        return v + 1, nil
//	    })
//	}
//
//	result := F.Pipe1(
//	    ioeither.Of[error](Counter{Value: 42}),
//	    ioeither.BindL(valueLens, increment),
//	)
func BindL[E, S, T any](
	lens L.Lens[S, T],
	f Kleisli[E, T, T],
) Operator[E, S, S] {
	return Bind(lens.Set, F.Flow2(lens.Get, f))
}

// LetL attaches the result of a pure computation to a context using a lens-based setter.
// This is a convenience function that combines Let with a lens, allowing you to use
// optics to update nested structures with pure transformations.
//
// The lens parameter provides both the getter and setter for a field within the structure S.
// The transformation function f receives the current value of the focused field and returns
// the new value directly (not wrapped in IOEither).
//
// Example:
//
//	type Counter struct {
//	    Value int
//	}
//
//	valueLens := lens.MakeLens(
//	    func(c Counter) int { return c.Value },
//	    func(c Counter, v int) Counter { c.Value = v; return c },
//	)
//
//	double := func(v int) int { return v * 2 }
//
//	result := F.Pipe1(
//	    ioeither.Of[error](Counter{Value: 21}),
//	    ioeither.LetL(valueLens, double),
//	)
func LetL[E, S, T any](
	lens L.Lens[S, T],
	f func(T) T,
) Operator[E, S, S] {
	return Let[E](lens.Set, F.Flow2(lens.Get, f))
}

// LetToL attaches a constant value to a context using a lens-based setter.
// This is a convenience function that combines LetTo with a lens, allowing you to use
// optics to set nested fields to specific values.
//
// The lens parameter provides the setter for a field within the structure S.
// Unlike LetL which transforms the current value, LetToL simply replaces it with
// the provided constant value b.
//
// Example:
//
//	type Config struct {
//	    Debug   bool
//	    Timeout int
//	}
//
//	debugLens := lens.MakeLens(
//	    func(c Config) bool { return c.Debug },
//	    func(c Config, d bool) Config { c.Debug = d; return c },
//	)
//
//	result := F.Pipe1(
//	    ioeither.Of[error](Config{Debug: true, Timeout: 30}),
//	    ioeither.LetToL(debugLens, false),
//	)
func LetToL[E, S, T any](
	lens L.Lens[S, T],
	b T,
) Operator[E, S, S] {
	return LetTo[E](lens.Set, b)
}
