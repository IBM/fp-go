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

package readerio

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
//	    Host string
//	    Port int
//	}
//	type Config struct {
//	    DefaultHost string
//	    DefaultPort int
//	}
//	result := readerio.Do[Config](State{})
func Do[R, S any](
	empty S,
) ReaderIO[R, S] {
	return Of[R](empty)
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
//	    Host string
//	    Port int
//	}
//	type Config struct {
//	    DefaultHost string
//	    DefaultPort int
//	}
//
//	result := F.Pipe2(
//	    readerio.Do[Config](State{}),
//	    readerio.Bind(
//	        func(host string) func(State) State {
//	            return func(s State) State { s.Host = host; return s }
//	        },
//	        func(s State) readerio.ReaderIO[Config, string] {
//	            return readerio.Asks(func(c Config) io.IO[string] {
//	                return io.Of(c.DefaultHost)
//	            })
//	        },
//	    ),
//	    readerio.Bind(
//	        func(port int) func(State) State {
//	            return func(s State) State { s.Port = port; return s }
//	        },
//	        func(s State) readerio.ReaderIO[Config, int] {
//	            // This can access s.Host from the previous step
//	            return readerio.Asks(func(c Config) io.IO[int] {
//	                return io.Of(c.DefaultPort)
//	            })
//	        },
//	    ),
//	)
func Bind[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) ReaderIO[R, T],
) func(ReaderIO[R, S1]) ReaderIO[R, S2] {
	return chain.Bind(
		Chain[R, S1, S2],
		Map[R, T, S2],
		setter,
		f,
	)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
func Let[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) func(ReaderIO[R, S1]) ReaderIO[R, S2] {
	return functor.Let(
		Map[R, S1, S2],
		setter,
		f,
	)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
func LetTo[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) func(ReaderIO[R, S1]) ReaderIO[R, S2] {
	return functor.LetTo(
		Map[R, S1, S2],
		setter,
		b,
	)
}

// BindTo initializes a new state [S1] from a value [T]
func BindTo[R, S1, T any](
	setter func(T) S1,
) func(ReaderIO[R, T]) ReaderIO[R, S1] {
	return chain.BindTo(
		Map[R, T, S1],
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
//	    Host string
//	    Port int
//	}
//	type Config struct {
//	    DefaultHost string
//	    DefaultPort int
//	}
//
//	// These operations are independent and can be combined with ApS
//	getHost := readerio.Asks(func(c Config) io.IO[string] {
//	    return io.Of(c.DefaultHost)
//	})
//	getPort := readerio.Asks(func(c Config) io.IO[int] {
//	    return io.Of(c.DefaultPort)
//	})
//
//	result := F.Pipe2(
//	    readerio.Do[Config](State{}),
//	    readerio.ApS(
//	        func(host string) func(State) State {
//	            return func(s State) State { s.Host = host; return s }
//	        },
//	        getHost,
//	    ),
//	    readerio.ApS(
//	        func(port int) func(State) State {
//	            return func(s State) State { s.Port = port; return s }
//	        },
//	        getPort,
//	    ),
//	)
func ApS[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa ReaderIO[R, T],
) func(ReaderIO[R, S1]) ReaderIO[R, S2] {
	return apply.ApS(
		Ap[S2, R, T],
		Map[R, S1, func(T) S2],
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
//	    Host string
//	    Port int
//	}
//	type Config struct {
//	    DefaultHost string
//	    DefaultPort int
//	}
//
//	portLens := lens.MakeLens(
//	    func(s State) int { return s.Port },
//	    func(s State, p int) State { s.Port = p; return s },
//	)
//
//	getPort := readerio.Asks(func(c Config) io.IO[int] {
//	    return io.Of(c.DefaultPort)
//	})
//	result := F.Pipe2(
//	    readerio.Of[Config](State{Host: "localhost"}),
//	    readerio.ApSL(portLens, getPort),
//	)
func ApSL[R, S, T any](
	lens L.Lens[S, T],
	fa ReaderIO[R, T],
) func(ReaderIO[R, S]) ReaderIO[R, S] {
	return ApS(lens.Set, fa)
}

// BindL is a variant of Bind that uses a lens to focus on a specific part of the context.
// This provides a more ergonomic API when working with nested structures, eliminating
// the need to manually write setter functions.
//
// The lens parameter provides both a getter and setter for a field of type T within
// the context S. The function f receives the current value of the focused field and
// returns a ReaderIO computation that produces an updated value.
//
// Example:
//
//	type State struct {
//	    Host string
//	    Port int
//	}
//	type Config struct {
//	    DefaultHost string
//	    DefaultPort int
//	}
//
//	portLens := lens.MakeLens(
//	    func(s State) int { return s.Port },
//	    func(s State, p int) State { s.Port = p; return s },
//	)
//
//	result := F.Pipe2(
//	    readerio.Do[Config](State{Host: "localhost"}),
//	    readerio.BindL(portLens, func(port int) readerio.ReaderIO[Config, int] {
//	        return readerio.Asks(func(c Config) io.IO[int] {
//	            return io.Of(c.DefaultPort)
//	        })
//	    }),
//	)
func BindL[R, S, T any](
	lens L.Lens[S, T],
	f func(T) ReaderIO[R, T],
) func(ReaderIO[R, S]) ReaderIO[R, S] {
	return Bind(lens.Set, F.Flow2(lens.Get, f))
}

// LetL is a variant of Let that uses a lens to focus on a specific part of the context.
// This provides a more ergonomic API when working with nested structures, eliminating
// the need to manually write setter functions.
//
// The lens parameter provides both a getter and setter for a field of type T within
// the context S. The function f receives the current value of the focused field and
// returns a new value (without wrapping in a ReaderIO).
//
// Example:
//
//	type State struct {
//	    Host string
//	    Port int
//	}
//
//	portLens := lens.MakeLens(
//	    func(s State) int { return s.Port },
//	    func(s State, p int) State { s.Port = p; return s },
//	)
//
//	result := F.Pipe2(
//	    readerio.Do[any](State{Host: "localhost", Port: 8080}),
//	    readerio.LetL(portLens, func(port int) int {
//	        return port + 1
//	    }),
//	)
func LetL[R, S, T any](
	lens L.Lens[S, T],
	f func(T) T,
) func(ReaderIO[R, S]) ReaderIO[R, S] {
	return Let[R](lens.Set, F.Flow2(lens.Get, f))
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
//	    Host string
//	    Port int
//	}
//
//	portLens := lens.MakeLens(
//	    func(s State) int { return s.Port },
//	    func(s State, p int) State { s.Port = p; return s },
//	)
//
//	result := F.Pipe2(
//	    readerio.Do[any](State{Host: "localhost"}),
//	    readerio.LetToL(portLens, 8080),
//	)
func LetToL[R, S, T any](
	lens L.Lens[S, T],
	b T,
) func(ReaderIO[R, S]) ReaderIO[R, S] {
	return LetTo[R](lens.Set, b)
}
