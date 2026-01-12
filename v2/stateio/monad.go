// Copyright (c) 2024 - 2025 IBM Corp.
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

package stateio

import (
	"github.com/IBM/fp-go/v2/internal/applicative"
	"github.com/IBM/fp-go/v2/internal/functor"
	"github.com/IBM/fp-go/v2/internal/monad"
	"github.com/IBM/fp-go/v2/internal/pointed"
)

// StateIOPointed implements the Pointed typeclass for StateIO.
// It provides the 'Of' operation to lift pure values into the StateIO context.
type StateIOPointed[
	S, A any,
] struct{}

// StateIOFunctor implements the Functor typeclass for StateIO.
// It provides the 'Map' operation to transform values within the StateIO context.
type StateIOFunctor[
	S, A, B any,
] struct{}

// StateIOApplicative implements the Applicative typeclass for StateIO.
// It provides 'Of', 'Map', and 'Ap' operations for applicative composition.
type StateIOApplicative[
	S, A, B any,
] struct{}

// StateIOMonad implements the Monad typeclass for StateIO.
// It provides 'Of', 'Map', 'Chain', and 'Ap' operations for monadic composition.
type StateIOMonad[
	S, A, B any,
] struct{}

// Of lifts a pure value into the StateIO context.
func (o *StateIOPointed[S, A]) Of(a A) StateIO[S, A] {
	return Of[S](a)
}

// Of lifts a pure value into the StateIO context.
func (o *StateIOMonad[S, A, B]) Of(a A) StateIO[S, A] {
	return Of[S](a)
}

// Of lifts a pure value into the StateIO context.
func (o *StateIOApplicative[S, A, B]) Of(a A) StateIO[S, A] {
	return Of[S](a)
}

// Map transforms the value within a StateIO using the provided function.
func (o *StateIOMonad[S, A, B]) Map(f func(A) B) Operator[S, A, B] {
	return Map[S](f)
}

// Map transforms the value within a StateIO using the provided function.
func (o *StateIOApplicative[S, A, B]) Map(f func(A) B) Operator[S, A, B] {
	return Map[S](f)
}

// Map transforms the value within a StateIO using the provided function.
func (o *StateIOFunctor[S, A, B]) Map(f func(A) B) Operator[S, A, B] {
	return Map[S](f)
}

// Chain sequences two StateIO computations, threading state through both.
func (o *StateIOMonad[S, A, B]) Chain(f Kleisli[S, A, B]) Operator[S, A, B] {
	return Chain(f)
}

// Ap applies a function wrapped in StateIO to a value wrapped in StateIO.
func (o *StateIOMonad[S, A, B]) Ap(fa StateIO[S, A]) Operator[S, func(A) B, B] {
	return Ap[B](fa)
}

// Ap applies a function wrapped in StateIO to a value wrapped in StateIO.
func (o *StateIOApplicative[S, A, B]) Ap(fa StateIO[S, A]) Operator[S, func(A) B, B] {
	return Ap[B](fa)
}

// Pointed returns a Pointed instance for StateIO.
// The Pointed typeclass provides the 'Of' operation to lift pure values
// into the StateIO context.
//
// Example:
//
//	p := Pointed[AppState, int]()
//	result := p.Of(42)
func Pointed[
	S, A any,
]() pointed.Pointed[A, StateIO[S, A]] {
	return &StateIOPointed[S, A]{}
}

// Functor returns a Functor instance for StateIO.
// The Functor typeclass provides the 'Map' operation to transform values
// within the StateIO context.
//
// Example:
//
//	f := Functor[AppState, int, string]()
//	result := f.Map(strconv.Itoa)(Of[AppState](42))
func Functor[
	S, A, B any,
]() functor.Functor[A, B, StateIO[S, A], StateIO[S, B]] {
	return &StateIOFunctor[S, A, B]{}
}

// Applicative returns an Applicative instance for StateIO.
// The Applicative typeclass provides 'Of', 'Map', and 'Ap' operations
// for applicative-style composition of StateIO computations.
//
// Example:
//
//	app := Applicative[AppState, int, string]()
//	fab := Of[AppState](func(x int) string { return strconv.Itoa(x) })
//	fa := Of[AppState](42)
//	result := app.Ap(fa)(fab)
func Applicative[
	S, A, B any,
]() applicative.Applicative[A, B, StateIO[S, A], StateIO[S, B], StateIO[S, func(A) B]] {
	return &StateIOApplicative[S, A, B]{}
}

// Monad returns a Monad instance for StateIO.
// The Monad typeclass provides 'Of', 'Map', 'Chain', and 'Ap' operations
// for monadic composition of StateIO computations.
//
// Example:
//
//	m := Monad[AppState, int, string]()
//	result := m.Chain(func(x int) StateIO[AppState, string] {
//	    return Of[AppState](strconv.Itoa(x))
//	})(Of[AppState](42))
func Monad[
	S, A, B any,
]() monad.Monad[A, B, StateIO[S, A], StateIO[S, B], StateIO[S, func(A) B]] {
	return &StateIOMonad[S, A, B]{}
}
