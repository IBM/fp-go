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

package iooption

import (
	"github.com/IBM/fp-go/v2/internal/apply"
	"github.com/IBM/fp-go/v2/internal/chain"
	"github.com/IBM/fp-go/v2/internal/functor"
)

// Bind creates an empty context of type [S] to be used with the [Bind] operation
func Do[S any](
	empty S,
) IOOption[S] {
	return Of(empty)
}

// Bind attaches the result of a computation to a context [S1] to produce a context [S2]
func Bind[S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) IOOption[T],
) func(IOOption[S1]) IOOption[S2] {
	return chain.Bind(
		Chain[S1, S2],
		Map[T, S2],
		setter,
		f,
	)
}

// Let attaches the result of a computation to a context [S1] to produce a context [S2]
func Let[S1, S2, T any](
	setter func(T) func(S1) S2,
	f func(S1) T,
) func(IOOption[S1]) IOOption[S2] {
	return functor.Let(
		Map[S1, S2],
		setter,
		f,
	)
}

// LetTo attaches the a value to a context [S1] to produce a context [S2]
func LetTo[S1, S2, T any](
	setter func(T) func(S1) S2,
	b T,
) func(IOOption[S1]) IOOption[S2] {
	return functor.LetTo(
		Map[S1, S2],
		setter,
		b,
	)
}

// BindTo initializes a new state [S1] from a value [T]
func BindTo[S1, T any](
	setter func(T) S1,
) func(IOOption[T]) IOOption[S1] {
	return chain.BindTo(
		Map[T, S1],
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
//	    Name  string
//	    Age   int
//	}
//
//	// These operations are independent and can be combined with ApS
//	getName := iooption.Some("Alice")
//	getAge := iooption.Some(30)
//
//	result := F.Pipe2(
//	    iooption.Do(State{}),
//	    iooption.ApS(
//	        func(name string) func(State) State {
//	            return func(s State) State { s.Name = name; return s }
//	        },
//	        getName,
//	    ),
//	    iooption.ApS(
//	        func(age int) func(State) State {
//	            return func(s State) State { s.Age = age; return s }
//	        },
//	        getAge,
//	    ),
//	)
func ApS[S1, S2, T any](
	setter func(T) func(S1) S2,
	fa IOOption[T],
) func(IOOption[S1]) IOOption[S2] {
	return apply.ApS(
		Ap[S2, T],
		Map[S1, func(T) S2],
		setter,
		fa,
	)
}
