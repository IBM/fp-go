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

package lazy

import (
	"github.com/IBM/fp-go/v2/io"
)

// Do creates an empty context of type [S] to be used with the [Bind] operation.
// This is the starting point for do-notation style composition.
//
// Example:
//
//	type State struct {
//	    Config Config
//	    Data   Data
//	}
//	result := lazy.Do(State{})
func Do[S any](
	empty S,
) Lazy[S] {
	return io.Do(empty)
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
//	    Config Config
//	    Data   Data
//	}
//
//	result := F.Pipe2(
//	    lazy.Do(State{}),
//	    lazy.Bind(
//	        func(cfg Config) func(State) State {
//	            return func(s State) State { s.Config = cfg; return s }
//	        },
//	        func(s State) lazy.Lazy[Config] {
//	            return lazy.MakeLazy(func() Config { return loadConfig() })
//	        },
//	    ),
//	    lazy.Bind(
//	        func(data Data) func(State) State {
//	            return func(s State) State { s.Data = data; return s }
//	        },
//	        func(s State) lazy.Lazy[Data] {
//	            // This can access s.Config from the previous step
//	            return lazy.MakeLazy(func() Data { return loadData(s.Config) })
//	        },
//	    ),
//	)
func Bind[S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) Lazy[T],
) func(Lazy[S1]) Lazy[S2] {
	return io.Bind(setter, f)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
func Let[S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) func(Lazy[S1]) Lazy[S2] {
	return io.Let(setter, f)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
func LetTo[S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) func(Lazy[S1]) Lazy[S2] {
	return io.LetTo(setter, b)
}

// BindTo initializes a new state [S1] from a value [T]
func BindTo[S1, T any](
	setter func(T) S1,
) func(Lazy[T]) Lazy[S1] {
	return io.BindTo(setter)
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
//	    Config  Config
//	    Data    Data
//	}
//
//	// These operations are independent and can be combined with ApS
//	getConfig := lazy.MakeLazy(func() Config { return loadConfig() })
//	getData := lazy.MakeLazy(func() Data { return loadData() })
//
//	result := F.Pipe2(
//	    lazy.Do(State{}),
//	    lazy.ApS(
//	        func(cfg Config) func(State) State {
//	            return func(s State) State { s.Config = cfg; return s }
//	        },
//	        getConfig,
//	    ),
//	    lazy.ApS(
//	        func(data Data) func(State) State {
//	            return func(s State) State { s.Data = data; return s }
//	        },
//	        getData,
//	    ),
//	)
func ApS[S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Lazy[T],
) func(Lazy[S1]) Lazy[S2] {
	return io.ApS(setter, fa)
}
