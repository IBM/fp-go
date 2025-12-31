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
	"context"

	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/monoid"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/optics/prism"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/result"
	"github.com/IBM/fp-go/v2/tailrec"
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

	// ReaderResult represents a computation that depends on a context.Context and produces either a value of type A or an error.
	// It combines the Reader pattern with Result (error handling), making it suitable for context-aware operations that may fail.
	ReaderResult[A any] = func(context.Context) (A, error)

	// Monoid represents a monoid structure for ReaderResult values.
	Monoid[A any] = monoid.Monoid[ReaderResult[A]]

	// Kleisli represents a Kleisli arrow from A to ReaderResult[B].
	// It's a function that takes a value of type A and returns a computation that produces B or an error in a context.
	Kleisli[A, B any] = Reader[A, ReaderResult[B]]

	// Operator represents a Kleisli arrow that operates on ReaderResult values.
	// It transforms a ReaderResult[A] into a ReaderResult[B], useful for composing context-aware operations.
	Operator[A, B any] = Kleisli[ReaderResult[A], B]

	// Lens represents an optic that focuses on a field of type A within a structure of type S.
	Lens[S, A any] = lens.Lens[S, A]

	// Prism represents an optic that focuses on a case of type A within a sum type S.
	Prism[S, A any] = prism.Prism[S, A]

	// Trampoline represents a tail-recursive computation that can be evaluated iteratively.
	// It's used to implement stack-safe recursion.
	Trampoline[A, B any] = tailrec.Trampoline[A, B]

	Predicate[A any] = predicate.Predicate[A]
)
