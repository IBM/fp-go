// Copyright (c) 2025 IBM Corp.
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
	"github.com/IBM/fp-go/v2/consumer"
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/reader"
)

type (
	// IO represents a lazy computation that performs side effects and produces a value of type A.
	// It's an alias for io.IO[A] and encapsulates effectful operations.
	IO[A any] = io.IO[A]

	Either[E, A any] = either.Either[E, A]

	// Reader represents a computation that depends on an environment of type R and produces a value of type A.
	// It's an alias for reader.Reader[R, A] and is used for dependency injection patterns.
	Reader[R, A any] = reader.Reader[R, A]

	// ReaderIO combines Reader and IO monads. It represents a computation that:
	// 1. Depends on an environment of type R (Reader aspect)
	// 2. Performs side effects and produces a value of type A (IO aspect)
	// This is useful for operations that need both dependency injection and effect management.
	ReaderIO[R, A any] = Reader[R, IO[A]]

	// Kleisli represents a Kleisli arrow for the ReaderIO monad.
	// It's a function from A to ReaderIO[R, B], which allows composition of
	// monadic functions. This is the fundamental building block for chaining
	// operations in the ReaderIO context.
	Kleisli[R, A, B any] = Reader[A, ReaderIO[R, B]]

	// Operator is a specialized Kleisli arrow that operates on ReaderIO values.
	// It transforms a ReaderIO[R, A] into a ReaderIO[R, B], making it useful
	// for building pipelines of ReaderIO operations. This is commonly used for
	// middleware-style transformations and operation composition.
	Operator[R, A, B any] = Kleisli[R, ReaderIO[R, A], B]

	Consumer[A any] = consumer.Consumer[A]
)
