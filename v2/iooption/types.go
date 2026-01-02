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

package iooption

import (
	"github.com/IBM/fp-go/v2/consumer"
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/optics/prism"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/tailrec"
)

type (
	// Either represents a value of one of two possible types (a disjoint union).
	Either[E, A any] = either.Either[E, A]

	// Option represents an optional value that may or may not be present.
	Option[A any] = option.Option[A]

	// IO represents a synchronous computation that cannot fail.
	IO[A any] = io.IO[A]

	// Lazy represents a deferred computation that produces a value of type A.
	Lazy[A any] = lazy.Lazy[A]

	// IOOption represents a synchronous computation that may not produce a value.
	// It combines IO (side effects) with Option (optional values).
	// Refer to [https://andywhite.xyz/posts/2021-01-27-rte-foundations/#ioeitherlte-agt] for more details.
	IOOption[A any] = io.IO[Option[A]]

	// Kleisli represents a Kleisli arrow for the IOOption monad.
	// It's a function from A to IOOption[B], used for composing operations that may not produce a value.
	Kleisli[A, B any] = reader.Reader[A, IOOption[B]]

	// Operator represents a function that transforms one IOOption into another.
	// It takes an IOOption[A] and produces an IOOption[B].
	Operator[A, B any] = Kleisli[IOOption[A], B]

	// Consumer represents a function that consumes a value of type A.
	// It's typically used for side effects like logging or updating state.
	Consumer[A any] = consumer.Consumer[A]

	// Lens is an optic that focuses on a field of type T within a structure of type S.
	Lens[S, T any] = lens.Lens[S, T]

	// Prism is an optic that focuses on a case of a sum type.
	Prism[S, T any] = prism.Prism[S, T]

	// Trampoline represents a tail-recursive computation that can be evaluated safely
	// without stack overflow. It's used for implementing stack-safe recursive algorithms.
	Trampoline[B, L any] = tailrec.Trampoline[B, L]

	// Predicate represents a function that tests a value of type A and returns a boolean.
	// It's commonly used for filtering and conditional operations.
	Predicate[A any] = predicate.Predicate[A]
)
