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
	L "github.com/IBM/fp-go/v2/optics/lens"

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
	f Kleisli[S1, T],
) Kleisli[Lazy[S1], S2] {
	return io.Bind(setter, f)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
func Let[S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) Kleisli[Lazy[S1], S2] {
	return io.Let(setter, f)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
func LetTo[S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) Kleisli[Lazy[S1], S2] {
	return io.LetTo(setter, b)
}

// BindTo initializes a new state [S1] from a value [T]
func BindTo[S1, T any](
	setter func(T) S1,
) Kleisli[Lazy[T], S1] {
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
) Kleisli[Lazy[S1], S2] {
	return io.ApS(setter, fa)
}

// ApSL is a variant of ApS that uses a lens to focus on a specific part of the context.
// This provides a more ergonomic API when working with nested structures, eliminating
// the need to manually write setter functions.
//
// The lens parameter provides both a getter and setter for a field of type T within
// the context S. This allows you to work with nested fields without manually managing
// the update logic.
//
// Example:
//
//	type Config struct {
//	    Host string
//	    Port int
//	}
//	type State struct {
//	    Config Config
//	    Data   string
//	}
//
//	configLens := L.Prop[State, Config]("Config")
//	getConfig := lazy.MakeLazy(func() Config { return Config{Host: "localhost", Port: 8080} })
//
//	result := F.Pipe2(
//	    lazy.Do(State{}),
//	    lazy.ApSL(configLens, getConfig),
//	)
func ApSL[S, T any](
	lens L.Lens[S, T],
	fa Lazy[T],
) Kleisli[Lazy[S], S] {
	return io.ApSL(lens, fa)
}

// BindL is a variant of Bind that uses a lens to focus on a specific part of the context.
// This provides a more ergonomic API when working with nested structures, eliminating
// the need to manually write setter functions.
//
// The lens parameter provides both a getter and setter for a field of type T within
// the context S. The function f receives the current value of the focused field and
// returns a new computation that produces an updated value.
//
// Example:
//
//	type Config struct {
//	    Host string
//	    Port int
//	}
//	type State struct {
//	    Config Config
//	    Data   string
//	}
//
//	configLens := L.Prop[State, Config]("Config")
//
//	result := F.Pipe2(
//	    lazy.Do(State{Config: Config{Host: "localhost"}}),
//	    lazy.BindL(configLens, func(cfg Config) lazy.Lazy[Config] {
//	        return lazy.MakeLazy(func() Config {
//	            cfg.Port = 8080
//	            return cfg
//	        })
//	    }),
//	)
func BindL[S, T any](
	lens L.Lens[S, T],
	f Kleisli[T, T],
) Kleisli[Lazy[S], S] {
	return io.BindL(lens, f)
}

// LetL is a variant of Let that uses a lens to focus on a specific part of the context.
// This provides a more ergonomic API when working with nested structures, eliminating
// the need to manually write setter functions.
//
// The lens parameter provides both a getter and setter for a field of type T within
// the context S. The function f receives the current value of the focused field and
// returns a new value (without wrapping in a monad).
//
// Example:
//
//	type Config struct {
//	    Host string
//	    Port int
//	}
//	type State struct {
//	    Config Config
//	    Data   string
//	}
//
//	configLens := L.Prop[State, Config]("Config")
//
//	result := F.Pipe2(
//	    lazy.Do(State{Config: Config{Host: "localhost"}}),
//	    lazy.LetL(configLens, func(cfg Config) Config {
//	        cfg.Port = 8080
//	        return cfg
//	    }),
//	)
func LetL[S, T any](
	lens L.Lens[S, T],
	f func(T) T,
) Kleisli[Lazy[S], S] {
	return io.LetL(lens, f)
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
//	type Config struct {
//	    Host string
//	    Port int
//	}
//	type State struct {
//	    Config Config
//	    Data   string
//	}
//
//	configLens := L.Prop[State, Config]("Config")
//	newConfig := Config{Host: "localhost", Port: 8080}
//
//	result := F.Pipe2(
//	    lazy.Do(State{}),
//	    lazy.LetToL(configLens, newConfig),
//	)
func LetToL[S, T any](
	lens L.Lens[S, T],
	b T,
) Kleisli[Lazy[S], S] {
	return io.LetToL(lens, b)
}
