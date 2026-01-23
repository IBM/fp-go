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

package readereither

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
)

type (
	// Option represents an optional value that may or may not be present.
	Option[A any] = option.Option[A]

	// Either represents a value of one of two possible types (a disjoint union).
	Either[E, A any] = either.Either[E, A]

	// Reader represents a computation that depends on an environment R and produces a value A.
	Reader[R, A any] = reader.Reader[R, A]

	// ReaderEither represents a computation that depends on an environment R and can fail
	// with an error E or succeed with a value A.
	// It combines Reader (dependency injection) with Either (error handling).

	ReaderEither[R, E, A any] = Reader[R, Either[E, A]]
	// Kleisli represents a Kleisli arrow for the ReaderEither monad.
	// It's a function from A to ReaderEither[R, E, B], used for composing operations that
	// depend on an environment and may fail.
	Kleisli[R, E, A, B any] = Reader[A, ReaderEither[R, E, B]]

	// Operator represents a function that transforms one ReaderEither into another.
	// It takes a ReaderEither[R, E, A] and produces a ReaderEither[R, E, B].
	Operator[R, E, A, B any] = Kleisli[R, E, ReaderEither[R, E, A], B]

	Lazy[A any] = lazy.Lazy[A]
)
