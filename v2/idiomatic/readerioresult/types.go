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

package readerioresult

import (
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/idiomatic/ioresult"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
	"github.com/IBM/fp-go/v2/result"
)

type (
	// Endomorphism represents a function from type A to type A.
	Endomorphism[A any] = endomorphism.Endomorphism[A]

	// Lazy represents a deferred computation that produces a value of type A when evaluated.
	Lazy[A any] = lazy.Lazy[A]

	// Option represents an optional value that may or may not be present.
	Option[A any] = option.Option[A]

	// Result represents an Either with error as the left type, compatible with Go's (value, error) tuple.
	Result[A any] = result.Result[A]

	// Reader represents a computation that depends on a read-only environment of type R and produces a value of type A.
	Reader[R, A any] = reader.Reader[R, A]

	// IO represents a computation that performs side effects and returns a value of type A.
	IO[A any] = io.IO[A]

	// IOResult represents a computation that performs IO and may fail with an error.
	IOResult[A any] = ioresult.IOResult[A]

	// ReaderIOResult represents a computation that depends on an environment R,
	// performs IO operations, and may fail with an error.
	// It is equivalent to Reader[R, IOResult[A]] or func(R) func() (A, error).
	ReaderIOResult[R, A any] = Reader[R, IOResult[A]]

	// ReaderIO represents a computation that depends on an environment R and performs side effects.
	ReaderIO[R, A any] = readerio.ReaderIO[R, A]

	// Monoid represents a monoid structure for ReaderIOResult values.
	Monoid[R, A any] = monoid.Monoid[ReaderIOResult[R, A]]

	// Kleisli represents a function from A to a ReaderIOResult of B.
	// It is used for chaining computations that depend on environment, perform IO, and may fail.
	Kleisli[R, A, B any] = Reader[A, ReaderIOResult[R, B]]

	// Operator represents a transformation from ReaderIOResult[R, A] to ReaderIOResult[R, B].
	// It is commonly used in function composition pipelines.
	Operator[R, A, B any] = Kleisli[R, ReaderIOResult[R, A], B]

	Predicate[A any] = predicate.Predicate[A]
)
