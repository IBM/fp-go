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

package either

import (
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
)

type (
	// Option is a type alias for option.Option, provided for convenience
	// when working with Either and Option together.
	Option[A any] = option.Option[A]

	// Lens is an optic that focuses on a field of type T within a structure of type S.
	Lens[S, T any] = lens.Lens[S, T]

	// Endomorphism represents a function from a type to itself (T -> T).
	Endomorphism[T any] = endomorphism.Endomorphism[T]

	// Lazy represents a deferred computation that produces a value of type T.
	Lazy[T any] = lazy.Lazy[T]

	// Kleisli represents a Kleisli arrow for the Either monad.
	// It's a function from A to Either[E, B], used for composing operations that may fail.
	Kleisli[E, A, B any] = reader.Reader[A, Either[E, B]]

	// Operator represents a function that transforms one Either into another.
	// It takes an Either[E, A] and produces an Either[E, B].
	Operator[E, A, B any] = Kleisli[E, Either[E, A], B]

	// Monoid represents a monoid structure for Either values.
	Monoid[E, A any] = monoid.Monoid[Either[E, A]]

	// Predicate represents a function that tests a value of type A and returns a boolean.
	// It's commonly used for filtering and conditional operations.
	Predicate[A any] = predicate.Predicate[A]

	Pair[L, R any] = pair.Pair[L, R]
)
