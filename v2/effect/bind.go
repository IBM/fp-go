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

package effect

import (
	"github.com/IBM/fp-go/v2/context/readerreaderioresult"
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
)

// Do creates an Effect with an initial state value.
// This is the starting point for do-notation style effect composition,
// allowing you to build up complex state transformations step by step.
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - S: The state type
//
// # Parameters
//
//   - empty: The initial state value
//
// # Returns
//
//   - Effect[C, S]: An effect that produces the initial state
//
// # Example
//
//	type State struct {
//		Name string
//		Age  int
//	}
//	eff := effect.Do[MyContext](State{})
//
//go:inline
func Do[C, S any](
	empty S,
) Effect[C, S] {
	return readerreaderioresult.Of[C](empty)
}

// Bind executes an effectful computation and binds its result to the state.
// This is the core operation for do-notation, allowing you to sequence effects
// while accumulating results in a state structure.
//
// # Type Parameters
//
//   - C: The context type required by the effects
//   - S1: The input state type
//   - S2: The output state type
//   - T: The type of value produced by the effect
//
// # Parameters
//
//   - setter: A function that takes the effect result and returns a state updater
//   - f: An effectful computation that depends on the current state
//
// # Returns
//
//   - Operator[C, S1, S2]: A function that transforms the state effect
//
// # Example
//
//	eff := effect.Bind(
//		func(age int) func(State) State {
//			return func(s State) State {
//				s.Age = age
//				return s
//			}
//		},
//		func(s State) Effect[MyContext, int] {
//			return effect.Of[MyContext](30)
//		},
//	)(effect.Do[MyContext](State{}))
//
//go:inline
func Bind[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	f Kleisli[C, S1, T],
) Operator[C, S1, S2] {
	return readerreaderioresult.Bind(setter, f)
}

// Let computes a pure value from the current state and binds it to the state.
// Unlike Bind, this doesn't perform any effects - it's for pure computations.
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - S1: The input state type
//   - S2: The output state type
//   - T: The type of computed value
//
// # Parameters
//
//   - setter: A function that takes the computed value and returns a state updater
//   - f: A pure function that computes a value from the current state
//
// # Returns
//
//   - Operator[C, S1, S2]: A function that transforms the state effect
//
// # Example
//
//	eff := effect.Let[MyContext](
//		func(nameLen int) func(State) State {
//			return func(s State) State {
//				s.NameLength = nameLen
//				return s
//			}
//		},
//		func(s State) int {
//			return len(s.Name)
//		},
//	)(stateEff)
//
//go:inline
func Let[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) Operator[C, S1, S2] {
	return readerreaderioresult.Let[C](setter, f)
}

// LetTo binds a constant value to the state.
// This is useful for setting fixed values in your state structure.
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - S1: The input state type
//   - S2: The output state type
//   - T: The type of the constant value
//
// # Parameters
//
//   - setter: A function that takes the constant and returns a state updater
//   - b: The constant value to bind
//
// # Returns
//
//   - Operator[C, S1, S2]: A function that transforms the state effect
//
// # Example
//
//	eff := effect.LetTo[MyContext](
//		func(age int) func(State) State {
//			return func(s State) State {
//				s.Age = age
//				return s
//			}
//		},
//		42,
//	)(stateEff)
//
//go:inline
func LetTo[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) Operator[C, S1, S2] {
	return readerreaderioresult.LetTo[C](setter, b)
}

// BindTo wraps a value in an initial state structure.
// This is typically used to start a bind chain by converting a simple value
// into a state structure.
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - S1: The state type to create
//   - T: The type of the input value
//
// # Parameters
//
//   - setter: A function that creates a state from the value
//
// # Returns
//
//   - Operator[C, T, S1]: A function that wraps the value in state
//
// # Example
//
//	eff := effect.BindTo[MyContext](func(name string) State {
//		return State{Name: name}
//	})(effect.Of[MyContext]("Alice"))
//
//go:inline
func BindTo[C, S1, T any](
	setter func(T) S1,
) Operator[C, T, S1] {
	return readerreaderioresult.BindTo[C](setter)
}

// ApS applies an effect and binds its result to the state using a setter function.
// This is similar to Bind but takes a pre-existing effect rather than a function
// that creates an effect from the state.
//
// # Type Parameters
//
//   - C: The context type required by the effects
//   - S1: The input state type
//   - S2: The output state type
//   - T: The type of value produced by the effect
//
// # Parameters
//
//   - setter: A function that takes the effect result and returns a state updater
//   - fa: The effect to apply
//
// # Returns
//
//   - Operator[C, S1, S2]: A function that transforms the state effect
//
// # Example
//
//	ageEffect := effect.Of[MyContext](30)
//	eff := effect.ApS(
//		func(age int) func(State) State {
//			return func(s State) State {
//				s.Age = age
//				return s
//			}
//		},
//		ageEffect,
//	)(stateEff)
//
//go:inline
func ApS[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Effect[C, T],
) Operator[C, S1, S2] {
	return readerreaderioresult.ApS(setter, fa)
}

// ApSL applies an effect and updates a field in the state using a lens.
// This provides a more ergonomic way to update nested state structures.
//
// # Type Parameters
//
//   - C: The context type required by the effects
//   - S: The state type
//   - T: The type of the field being updated
//
// # Parameters
//
//   - lens: A lens focusing on the field to update
//   - fa: The effect producing the new field value
//
// # Returns
//
//   - Operator[C, S, S]: A function that updates the state field
//
// # Example
//
//	ageLens := lens.MakeLens(
//		func(s State) int { return s.Age },
//		func(s State, age int) State { s.Age = age; return s },
//	)
//	ageEffect := effect.Of[MyContext](30)
//	eff := effect.ApSL(ageLens, ageEffect)(stateEff)
//
//go:inline
func ApSL[C, S, T any](
	lens Lens[S, T],
	fa Effect[C, T],
) Operator[C, S, S] {
	return readerreaderioresult.ApSL(lens, fa)
}

// BindL executes an effectful computation on a field and updates it using a lens.
// The effect function receives the current field value and produces a new value.
//
// # Type Parameters
//
//   - C: The context type required by the effects
//   - S: The state type
//   - T: The type of the field being updated
//
// # Parameters
//
//   - lens: A lens focusing on the field to update
//   - f: An effectful function that transforms the field value
//
// # Returns
//
//   - Operator[C, S, S]: A function that updates the state field
//
// # Example
//
//	ageLens := lens.MakeLens(
//		func(s State) int { return s.Age },
//		func(s State, age int) State { s.Age = age; return s },
//	)
//	eff := effect.BindL(
//		ageLens,
//		func(age int) Effect[MyContext, int] {
//			return effect.Of[MyContext](age + 1)
//		},
//	)(stateEff)
//
//go:inline
func BindL[C, S, T any](
	lens Lens[S, T],
	f func(T) Effect[C, T],
) Operator[C, S, S] {
	return readerreaderioresult.BindL(lens, f)
}

// LetL computes a new field value from the current value using a lens.
// This is a pure transformation of a field within the state.
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - S: The state type
//   - T: The type of the field being updated
//
// # Parameters
//
//   - lens: A lens focusing on the field to update
//   - f: A pure function that transforms the field value
//
// # Returns
//
//   - Operator[C, S, S]: A function that updates the state field
//
// # Example
//
//	ageLens := lens.MakeLens(
//		func(s State) int { return s.Age },
//		func(s State, age int) State { s.Age = age; return s },
//	)
//	eff := effect.LetL[MyContext](
//		ageLens,
//		func(age int) int { return age * 2 },
//	)(stateEff)
//
//go:inline
func LetL[C, S, T any](
	lens Lens[S, T],
	f func(T) T,
) Operator[C, S, S] {
	return readerreaderioresult.LetL[C](lens, f)
}

// LetToL sets a field to a constant value using a lens.
//
// # Type Parameters
//
//   - C: The context type required by the effect
//   - S: The state type
//   - T: The type of the field being updated
//
// # Parameters
//
//   - lens: A lens focusing on the field to update
//   - b: The constant value to set
//
// # Returns
//
//   - Operator[C, S, S]: A function that updates the state field
//
// # Example
//
//	ageLens := lens.MakeLens(
//		func(s State) int { return s.Age },
//		func(s State, age int) State { s.Age = age; return s },
//	)
//	eff := effect.LetToL[MyContext](ageLens, 42)(stateEff)
//
//go:inline
func LetToL[C, S, T any](
	lens Lens[S, T],
	b T,
) Operator[C, S, S] {
	return readerreaderioresult.LetToL[C](lens, b)
}

//go:inline
func BindIOEitherK[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	f ioeither.Kleisli[error, S1, T],
) Operator[C, S1, S2] {
	return readerreaderioresult.BindIOEitherK[C](setter, f)
}

//go:inline
func BindIOResultK[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	f ioresult.Kleisli[S1, T],
) Operator[C, S1, S2] {
	return readerreaderioresult.BindIOResultK[C](setter, f)
}

//go:inline
func BindIOK[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	f io.Kleisli[S1, T],
) Operator[C, S1, S2] {
	return readerreaderioresult.BindIOK[C](setter, f)
}

//go:inline
func BindReaderK[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	f reader.Kleisli[C, S1, T],
) Operator[C, S1, S2] {
	return readerreaderioresult.BindReaderK(setter, f)
}

//go:inline
func BindReaderIOK[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	f readerio.Kleisli[C, S1, T],
) Operator[C, S1, S2] {
	return readerreaderioresult.BindReaderIOK(setter, f)
}

//go:inline
func BindEitherK[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	f either.Kleisli[error, S1, T],
) Operator[C, S1, S2] {
	return readerreaderioresult.BindEitherK[C](setter, f)
}

//go:inline
func BindIOEitherKL[C, S, T any](
	lens Lens[S, T],
	f ioeither.Kleisli[error, T, T],
) Operator[C, S, S] {
	return readerreaderioresult.BindIOEitherKL[C](lens, f)
}

//go:inline
func BindIOKL[C, S, T any](
	lens Lens[S, T],
	f io.Kleisli[T, T],
) Operator[C, S, S] {
	return readerreaderioresult.BindIOKL[C](lens, f)
}

//go:inline
func BindReaderKL[C, S, T any](
	lens Lens[S, T],
	f reader.Kleisli[C, T, T],
) Operator[C, S, S] {
	return readerreaderioresult.BindReaderKL(lens, f)
}

//go:inline
func BindReaderIOKL[C, S, T any](
	lens Lens[S, T],
	f readerio.Kleisli[C, T, T],
) Operator[C, S, S] {
	return readerreaderioresult.BindReaderIOKL(lens, f)
}

//go:inline
func ApIOEitherS[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa IOEither[error, T],
) Operator[C, S1, S2] {
	return readerreaderioresult.ApIOEitherS[C](setter, fa)
}

//go:inline
func ApIOS[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa IO[T],
) Operator[C, S1, S2] {
	return readerreaderioresult.ApIOS[C](setter, fa)
}

//go:inline
func ApReaderS[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Reader[C, T],
) Operator[C, S1, S2] {
	return readerreaderioresult.ApReaderS(setter, fa)
}

//go:inline
func ApReaderIOS[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa ReaderIO[C, T],
) Operator[C, S1, S2] {
	return readerreaderioresult.ApReaderIOS(setter, fa)
}

//go:inline
func ApEitherS[C, S1, S2, T any](
	setter func(T) func(S1) S2,
	fa Either[error, T],
) Operator[C, S1, S2] {
	return readerreaderioresult.ApEitherS[C](setter, fa)
}

//go:inline
func ApIOEitherSL[C, S, T any](
	lens Lens[S, T],
	fa IOEither[error, T],
) Operator[C, S, S] {
	return readerreaderioresult.ApIOEitherSL[C](lens, fa)
}

//go:inline
func ApIOSL[C, S, T any](
	lens Lens[S, T],
	fa IO[T],
) Operator[C, S, S] {
	return readerreaderioresult.ApIOSL[C](lens, fa)
}

//go:inline
func ApReaderSL[C, S, T any](
	lens Lens[S, T],
	fa Reader[C, T],
) Operator[C, S, S] {
	return readerreaderioresult.ApReaderSL(lens, fa)
}

//go:inline
func ApReaderIOSL[C, S, T any](
	lens Lens[S, T],
	fa ReaderIO[C, T],
) Operator[C, S, S] {
	return readerreaderioresult.ApReaderIOSL(lens, fa)
}

//go:inline
func ApEitherSL[C, S, T any](
	lens Lens[S, T],
	fa Either[error, T],
) Operator[C, S, S] {
	return readerreaderioresult.ApEitherSL[C](lens, fa)
}
