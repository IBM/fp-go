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

package reader

import (
	"github.com/IBM/fp-go/v2/internal/apply"
	"github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/functor"
)

// Do creates an empty context of type [S] to be used with the [Bind] operation.
// This is the starting point for the do-notation style of composing Reader computations.
//
// Example:
//
//	type State struct {
//	    Name string
//	    Age  int
//	}
//	type Config struct {
//	    DefaultName string
//	    DefaultAge  int
//	}
//
//	result := function.Pipe3(
//	    reader.Do[Config](State{}),
//	    reader.Bind(
//	        func(name string) func(State) State {
//	            return func(s State) State { s.Name = name; return s }
//	        },
//	        func(s State) reader.Reader[Config, string] {
//	            return reader.Asks(func(c Config) string { return c.DefaultName })
//	        },
//	    ),
//	    reader.Bind(
//	        func(age int) func(State) State {
//	            return func(s State) State { s.Age = age; return s }
//	        },
//	        func(s State) reader.Reader[Config, int] {
//	            return reader.Asks(func(c Config) int { return c.DefaultAge })
//	        },
//	    ),
//	)
func Do[R, S any](
	empty S,
) Reader[R, S] {
	return Of[R](empty)
}

// Bind attaches the result of a computation to a context [S1] to produce a context [S2].
// This enables building up complex computations in a pipeline where each step can depend
// on the results of previous steps and access the shared environment.
//
// The setter function takes the result of the computation and returns a function that
// updates the context from S1 to S2.
//
// Example:
//
//	type State struct { Value int }
//	type Config struct { Increment int }
//
//	addIncrement := reader.Bind(
//	    func(inc int) func(State) State {
//	        return func(s State) State { return State{Value: s.Value + inc} }
//	    },
//	    func(s State) reader.Reader[Config, int] {
//	        return reader.Asks(func(c Config) int { return c.Increment })
//	    },
//	)
func Bind[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) Reader[R, T],
) Operator[R, S1, S2] {
	return chain.Bind(
		Chain[R, S1, S2],
		Map[R, T, S2],
		setter,
		f,
	)
}

// Let attaches the result of a pure computation to a context [S1] to produce a context [S2].
// Unlike Bind, the computation function f does not return a Reader, just a plain value.
// This is useful for transformations that don't need to access the environment.
//
// Example:
//
//	type State struct {
//	    FirstName string
//	    LastName  string
//	    FullName  string
//	}
//
//	addFullName := reader.Let(
//	    func(full string) func(State) State {
//	        return func(s State) State { s.FullName = full; return s }
//	    },
//	    func(s State) string {
//	        return s.FirstName + " " + s.LastName
//	    },
//	)
func Let[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) Operator[R, S1, S2] {
	return functor.Let(
		Map[R, S1, S2],
		setter,
		f,
	)
}

// LetTo attaches a constant value to a context [S1] to produce a context [S2].
// This is useful for adding fixed values to the context without any computation.
//
// Example:
//
//	type State struct {
//	    Name    string
//	    Version string
//	}
//
//	addVersion := reader.LetTo(
//	    func(v string) func(State) State {
//	        return func(s State) State { s.Version = v; return s }
//	    },
//	    "1.0.0",
//	)
func LetTo[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) Operator[R, S1, S2] {
	return functor.LetTo(
		Map[R, S1, S2],
		setter,
		b,
	)
}

// BindTo initializes a new state [S1] from a value [T].
// This is typically used to start a binding chain by wrapping an initial Reader value
// into a state structure.
//
// Example:
//
//	type State struct { Name string }
//	type Config struct { DefaultName string }
//
//	getName := reader.Asks(func(c Config) string { return c.DefaultName })
//	initState := reader.BindTo(func(name string) State {
//	    return State{Name: name}
//	})
//	result := initState(getName)
func BindTo[R, S1, T any](
	setter func(T) S1,
) Operator[R, T, S1] {
	return chain.BindTo(
		Map[R, T, S1],
		setter,
	)
}

// ApS attaches a value to a context [S1] to produce a context [S2] by considering
// the context and the value concurrently (using Applicative rather than Monad).
//
// This is useful when you have independent computations that can be combined without
// one depending on the result of the other.
//
// Example:
//
//	type State struct {
//	    Host string
//	    Port int
//	}
//	type Config struct {
//	    Host string
//	    Port int
//	}
//
//	getPort := reader.Asks(func(c Config) int { return c.Port })
//	addPort := reader.ApS(
//	    func(port int) func(State) State {
//	        return func(s State) State { s.Port = port; return s }
//	    },
//	    getPort,
//	)
func ApS[R, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Reader[R, T],
) Operator[R, S1, S2] {
	return apply.ApS(
		Ap[S2, R, T],
		Map[R, S1, func(T) S2],
		setter,
		fa,
	)
}
