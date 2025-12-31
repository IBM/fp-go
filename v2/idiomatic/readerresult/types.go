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
	// Endomorphism represents a function from type A to type A.
	Endomorphism[A any] = endomorphism.Endomorphism[A]

	// Lazy represents a deferred computation that produces a value of type A when evaluated.
	Lazy[A any] = lazy.Lazy[A]

	// Option represents an optional value that may or may not be present.
	Option[A any] = option.Option[A]

	// Either represents a value that can be one of two types: Left (E) or Right (A).
	Either[E, A any] = either.Either[E, A]

	// Result represents an Either with error as the left type, compatible with Go's (value, error) tuple.
	Result[A any] = result.Result[A]

	// Reader represents a computation that depends on a read-only environment of type R and produces a value of type A.
	Reader[R, A any] = reader.Reader[R, A]

	// ReaderResult represents a computation that depends on an environment R and may fail with an error.
	// It is equivalent to Reader[R, Result[A]] or func(R) (A, error).
	// This combines dependency injection with error handling in a functional style.
	ReaderResult[R, A any] = func(R) (A, error)

	// Monoid represents a monoid structure for ReaderResult values.
	Monoid[R, A any] = monoid.Monoid[ReaderResult[R, A]]

	// Kleisli represents a function from A to a ReaderResult of B.
	// It is used for chaining computations that depend on environment and may fail.
	Kleisli[R, A, B any] = Reader[A, ReaderResult[R, B]]

	// Operator represents a transformation from ReaderResult[R, A] to ReaderResult[R, B].
	// It is commonly used in function composition pipelines.
	Operator[R, A, B any] = Kleisli[R, ReaderResult[R, A], B]

	Predicate[A any] = predicate.Predicate[A]
)
