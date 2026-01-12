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
	"context"

	"github.com/IBM/fp-go/v2/consumer"
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
	"github.com/IBM/fp-go/v2/tailrec"
)

type (
	// Lazy represents a deferred computation that produces a value of type A when executed.
	// The computation is not executed until explicitly invoked.
	Lazy[A any] = lazy.Lazy[A]

	// IO represents a side-effectful computation that produces a value of type A.
	// The computation is deferred and only executed when invoked.
	//
	// IO[A] is equivalent to func() A
	IO[A any] = io.IO[A]

	// Reader represents a computation that depends on a context of type R.
	// This is used for dependency injection and accessing shared context.
	//
	// Reader[R, A] is equivalent to func(R) A
	Reader[R, A any] = reader.Reader[R, A]

	// ReaderIO represents a context-dependent computation that performs side effects.
	// This is specialized to use [context.Context] as the context type.
	//
	// ReaderIO[A] is equivalent to func(context.Context) func() A
	ReaderIO[A any] = readerio.ReaderIO[context.Context, A]

	// Kleisli represents a Kleisli arrow for the ReaderIO monad.
	// It is a function that takes a value of type A and returns a ReaderIO computation
	// that produces a value of type B.
	//
	// Kleisli arrows are used for composing monadic computations and are fundamental
	// to functional programming patterns involving effects and context.
	//
	// Kleisli[A, B] is equivalent to func(A) func(context.Context) func() B
	Kleisli[A, B any] = reader.Reader[A, ReaderIO[B]]

	// Operator represents a transformation from one ReaderIO computation to another.
	// It takes a ReaderIO[A] and returns a ReaderIO[B], allowing for the composition
	// of context-dependent, side-effectful computations.
	//
	// Operators are useful for building pipelines of ReaderIO computations where
	// each step can depend on the previous computation's result.
	//
	// Operator[A, B] is equivalent to func(ReaderIO[A]) func(context.Context) func() B
	Operator[A, B any] = Kleisli[ReaderIO[A], B]

	Consumer[A any] = consumer.Consumer[A]

	Either[E, A any] = either.Either[E, A]

	Trampoline[B, L any] = tailrec.Trampoline[B, L]

	Predicate[A any] = predicate.Predicate[A]

	Void = function.Void
)
