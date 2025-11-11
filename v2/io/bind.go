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

package io

import (
	F "github.com/IBM/fp-go/v2/function"
	INTA "github.com/IBM/fp-go/v2/internal/apply"
	INTC "github.com/IBM/fp-go/v2/internal/chain"
	INTF "github.com/IBM/fp-go/v2/internal/functor"
	L "github.com/IBM/fp-go/v2/optics/lens"
)

// Do creates an empty context of type S to be used with the Bind operation.
// This is the starting point for do-notation style composition.
//
// Example:
//
//	type State struct {
//	    user User
//	    posts []Post
//	}
//	result := pipe.Pipe2(
//	    io.Do(State{}),
//	    io.Bind("user", fetchUser),
//	    io.Bind("posts", func(s State) io.IO[[]Post] {
//	        return fetchPosts(s.user.Id)
//	    }),
//	)
func Do[S any](
	empty S,
) IO[S] {
	return Of(empty)
}

// Bind attaches the result of an IO computation to a context S1 to produce a context S2.
// This is used in do-notation style composition to build up state incrementally.
//
// The setter function takes the result T and returns a function that updates S1 to S2.
//
// Example:
//
//	io.Bind(func(user User) func(s State) State {
//	    return func(s State) State {
//	        s.user = user
//	        return s
//	    }
//	}, fetchUser)
func Bind[S1, S2, T any](
	setter func(T) func(S1) S2,
	f Kleisli[S1, T],
) Operator[S1, S2] {
	return INTC.Bind(
		Chain[S1, S2],
		Map[T, S2],
		setter,
		f,
	)
}

// Let attaches the result of a pure computation to a context S1 to produce a context S2.
// Similar to Bind, but for pure (non-IO) computations.
//
// Example:
//
//	io.Let(func(count int) func(s State) State {
//	    return func(s State) State {
//	        s.count = count
//	        return s
//	    }
//	}, func(s State) int { return len(s.items) })
func Let[S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) Operator[S1, S2] {
	return INTF.Let(
		Map[S1, S2],
		setter,
		f,
	)
}

// LetTo attaches a constant value to a context S1 to produce a context S2.
// Similar to Let, but with a constant value instead of a computation.
//
// Example:
//
//	io.LetTo(func(status string) func(s State) State {
//	    return func(s State) State {
//	        s.status = status
//	        return s
//	    }
//	}, "ready")
func LetTo[S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) Operator[S1, S2] {
	return INTF.LetTo(
		Map[S1, S2],
		setter,
		b,
	)
}

// BindTo initializes a new state S1 from a value T.
// This is typically used to start a do-notation chain from a single value.
//
// Example:
//
//	io.BindTo(func(user User) State {
//	    return State{user: user}
//	})
func BindTo[S1, T any](
	setter func(T) S1,
) Operator[T, S1] {
	return INTC.BindTo(
		Map[T, S1],
		setter,
	)
}

// ApS attaches a value to a context S1 to produce a context S2 by considering
// the context and the value concurrently (using applicative operations).
// This allows parallel execution of independent computations.
//
// Example:
//
//	io.ApS(func(posts []Post) func(s State) State {
//	    return func(s State) State {
//	        s.posts = posts
//	        return s
//	    }
//	}, fetchPosts())
func ApS[S1, S2, T any](
	setter func(T) func(S1) S2,
	fa IO[T],
) Operator[S1, S2] {
	return INTA.ApS(
		Ap[S2, T],
		Map[S1, func(T) S2],
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
//	    io.Of(Config{Host: "localhost"}),
//	    io.ApSL(portLens, io.Of(8080)),
//	)
func ApSL[S, T any](
	lens L.Lens[S, T],
	fa IO[T],
) Operator[S, S] {
	return ApS(lens.Set, fa)
}

// BindL attaches the result of a computation to a context using a lens-based setter.
// This is a convenience function that combines Bind with a lens, allowing you to use
// optics to update nested structures based on their current values.
//
// The lens parameter provides both the getter and setter for a field within the structure S.
// The computation function f receives the current value of the focused field and returns
// an IO that produces the new value.
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
//	// Increment the counter asynchronously
//	increment := func(v int) io.IO[int] {
//	    return io.Of(v + 1)
//	}
//
//	result := F.Pipe1(
//	    io.Of(Counter{Value: 42}),
//	    io.BindL(valueLens, increment),
//	) // IO[Counter{Value: 43}]
func BindL[S, T any](
	lens L.Lens[S, T],
	f Kleisli[T, T],
) Operator[S, S] {
	return Bind(lens.Set, F.Flow2(lens.Get, f))
}

// LetL attaches the result of a pure computation to a context using a lens-based setter.
// This is a convenience function that combines Let with a lens, allowing you to use
// optics to update nested structures with pure transformations.
//
// The lens parameter provides both the getter and setter for a field within the structure S.
// The transformation function f receives the current value of the focused field and returns
// the new value directly (not wrapped in IO).
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
//	// Double the counter value
//	double := func(v int) int { return v * 2 }
//
//	result := F.Pipe1(
//	    io.Of(Counter{Value: 21}),
//	    io.LetL(valueLens, double),
//	) // IO[Counter{Value: 42}]
func LetL[S, T any](
	lens L.Lens[S, T],
	f func(T) T,
) Operator[S, S] {
	return Let(lens.Set, F.Flow2(lens.Get, f))
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
//	    io.Of(Config{Debug: true, Timeout: 30}),
//	    io.LetToL(debugLens, false),
//	) // IO[Config{Debug: false, Timeout: 30}]
func LetToL[S, T any](
	lens L.Lens[S, T],
	b T,
) Operator[S, S] {
	return LetTo(lens.Set, b)
}
