// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache LicensVersion 2.0 (the "License");
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

package ioresult

import (
	"github.com/IBM/fp-go/v2/ioeither"
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
//
//go:inline
func Do[S any](
	empty S,
) IOResult[S] {
	return ioeither.Do[error](empty)
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
//	        func(s State) ioeither.IOResult[error, User] {
//	            return ioeither.TryCatch(func() (User, error) {
//	                return fetchUser()
//	            })
//	        },
//	    ),
//	    ioeither.Bind(
//	        func(posts []Post) func(State) State {
//	            return func(s State) State { s.Posts = posts; return s }
//	        },
//	        func(s State) ioeither.IOResult[error, []Post] {
//	            // This can access s.User from the previous step
//	            return ioeither.TryCatch(func() ([]Post, error) {
//	                return fetchPostsForUser(s.User.ID)
//	            })
//	        },
//	    ),
//	)
//
//go:inline
func Bind[S1, S2, T any](
	setter func(T) func(S1) S2,
	f Kleisli[S1, T],
) Operator[S1, S2] {
	return ioeither.Bind(setter, f)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
//
//go:inline
func Let[S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) Operator[S1, S2] {
	return ioeither.Let[error](setter, f)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
//
//go:inline
func LetTo[S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) Operator[S1, S2] {
	return ioeither.LetTo[error](setter, b)
}

// BindTo initializes a new state [S1] from a value [T]
//
//go:inline
func BindTo[S1, T any](
	setter func(T) S1,
) Operator[T, S1] {
	return ioeither.BindTo[error](setter)
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
//
//go:inline
func ApS[S1, S2, T any](
	setter func(T) func(S1) S2,
	fa IOResult[T],
) Operator[S1, S2] {
	return ioeither.ApS(setter, fa)
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
//
//go:inline
func ApSL[S, T any](
	lens L.Lens[S, T],
	fa IOResult[T],
) Operator[S, S] {
	return ioeither.ApSL(lens, fa)
}

// BindL attaches the result of a computation to a context using a lens-based setter.
// This is a convenience function that combines Bind with a lens, allowing you to use
// optics to update nested structures based on their current values.
//
// The lens parameter provides both the getter and setter for a field within the structure S.
// The computation function f receives the current value of the focused field and returns
// an IOResult that produces the new value.
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
//	increment := func(v int) ioeither.IOResult[error, int] {
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
//
//go:inline
func BindL[S, T any](
	lens L.Lens[S, T],
	f Kleisli[T, T],
) Operator[S, S] {
	return ioeither.BindL(lens, f)
}

// LetL attaches the result of a pure computation to a context using a lens-based setter.
// This is a convenience function that combines Let with a lens, allowing you to use
// optics to update nested structures with pure transformations.
//
// The lens parameter provides both the getter and setter for a field within the structure S.
// The transformation function f receives the current value of the focused field and returns
// the new value directly (not wrapped in IOResult).
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
//
//go:inline
func LetL[S, T any](
	lens L.Lens[S, T],
	f Endomorphism[T],
) Operator[S, S] {
	return ioeither.LetL[error](lens, f)
}

// LetToL attaches a constant value to a context using a lens-based setter.
// This is a convenience function that combines LetTo with a lens, allowing you to use
// optics to set nested fields to specific values.
//
// The lens parameter provides the setter for a field within the structure S.
// Unlike LetL which transforms the current valuLetToL simply replaces it with
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
//	    ioeither.Of[error](Config{Debug: truTimeout: 30}),
//	    ioeither.LetToL(debugLens, false),
//	)
//
//go:inline
func LetToL[S, T any](
	lens L.Lens[S, T],
	b T,
) Operator[S, S] {
	return ioeither.LetToL[error](lens, b)
}
