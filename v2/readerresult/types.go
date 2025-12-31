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

package readerresult

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/result"
)

type (
	// Endomorphism represents a function from a type to itself (A -> A).
	Endomorphism[A any] = endomorphism.Endomorphism[A]

	// Lazy represents a deferred computation that produces a value of type A.
	Lazy[A any] = lazy.Lazy[A]

	// Option represents an optional value that may or may not be present.
	Option[A any] = option.Option[A]

	// Either represents a value of one of two possible types (a disjoint union).
	Either[E, A any] = either.Either[E, A]

	// Result represents a computation that may fail with an error.
	Result[A any] = result.Result[A]

	// Reader represents a computation that depends on an environment R and produces a value A.
	Reader[R, A any] = reader.Reader[R, A]

	// ReaderResult represents a computation that depends on an environment R and may fail with an error.
	// It combines Reader (dependency injection) with Result (error handling).
	ReaderResult[R, A any] = Reader[R, Result[A]]

	// Monoid represents a monoid structure for ReaderResult values.
	Monoid[R, A any] = monoid.Monoid[ReaderResult[R, A]]

	// Kleisli represents a Kleisli arrow for the ReaderResult monad.
	// It's a function from A to ReaderResult[R, B], used for composing operations that
	// depend on an environment and may fail.
	Kleisli[R, A, B any] = Reader[A, ReaderResult[R, B]]

	// Operator represents a function that transforms one ReaderResult into another.
	// It takes a ReaderResult[R, A] and produces a ReaderResult[R, B].
	Operator[R, A, B any] = Kleisli[R, ReaderResult[R, A], B]

	// Predicate represents a function that tests a value of type A and returns a boolean.
	Predicate[A any] = predicate.Predicate[A]
)
