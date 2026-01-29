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
	"context"

	"github.com/IBM/fp-go/v2/consumer"
	"github.com/IBM/fp-go/v2/context/ioresult"
	"github.com/IBM/fp-go/v2/context/readerresult"
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/endomorphism"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/ioref"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/optics/prism"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/predicate"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readereither"
	"github.com/IBM/fp-go/v2/readerio"
	RIOR "github.com/IBM/fp-go/v2/readerioresult"
	"github.com/IBM/fp-go/v2/readeroption"
	"github.com/IBM/fp-go/v2/result"
	"github.com/IBM/fp-go/v2/state"
	"github.com/IBM/fp-go/v2/tailrec"
)

type (
	// Option represents an optional value that may or may not be present.
	// It is used in operations that may not produce a value.
	Option[A any] = option.Option[A]

	// Either represents a computation that can result in either an error or a success value.
	// This is specialized to use [error] as the left (error) type, which is the standard
	// error type in Go.
	//
	// Either[A] is equivalent to Either[error, A] from the either package.
	Either[A any] = either.Either[error, A]

	Result[A any] = result.Result[A]

	// Lazy represents a deferred computation that produces a value of type A when executed.
	// The computation is not executed until explicitly invoked.
	Lazy[A any] = lazy.Lazy[A]

	// IO represents a side-effectful computation that produces a value of type A.
	// The computation is deferred and only executed when invoked.
	//
	// IO[A] is equivalent to func() A
	IO[A any] = io.IO[A]

	// IOEither represents a side-effectful computation that can fail with an error.
	// This combines IO (side effects) with Either (error handling).
	//
	// IOEither[A] is equivalent to func() Either[error, A]
	IOEither[A any] = ioeither.IOEither[error, A]

	IOResult[A any] = ioresult.IOResult[A]

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

	// ReaderIOResult is the main type of this package. It represents a computation that:
	//   - Depends on a [context.Context] (Reader aspect)
	//   - Performs side effects (IO aspect)
	//   - Can fail with an [error] (Either aspect)
	//   - Produces a value of type A on success
	//
	// This is a specialization of [readerioeither.ReaderIOResult] with:
	//   - Context type fixed to [context.Context]
	//   - Error type fixed to [error]
	//
	// The type is defined as:
	//   ReaderIOResult[A] = func(context.Context) func() Either[error, A]
	//
	// Example usage:
	//   func fetchUser(id string) ReaderIOResult[User] {
	//       return func(ctx context.Context) func() Either[error, User] {
	//           return func() Either[error, User] {
	//               user, err := userService.Get(ctx, id)
	//               if err != nil {
	//                   return either.Left[User](err)
	//               }
	//               return either.Right[error](user)
	//           }
	//       }
	//   }
	//
	// The computation is executed by providing a context and then invoking the result:
	//   ctx := t.Context()
	//   result := fetchUser("123")(ctx)()
	ReaderIOResult[A any] = RIOR.ReaderIOResult[context.Context, A]

	Kleisli[A, B any] = reader.Reader[A, ReaderIOResult[B]]

	// Operator represents a transformation from one ReaderIOResult to another.
	// This is useful for point-free style composition and building reusable transformations.
	//
	// Operator[A, B] is equivalent to Kleisli[ReaderIOResult[A], B]
	//
	// Example usage:
	//   // Define a reusable transformation
	//   var toUpper Operator[string, string] = Map(strings.ToUpper)
	//
	//   // Apply the transformation
	//   result := toUpper(computation)
	Operator[A, B any] = Kleisli[ReaderIOResult[A], B]

	ReaderResult[A any]       = readerresult.ReaderResult[A]
	ReaderEither[R, E, A any] = readereither.ReaderEither[R, E, A]
	ReaderOption[R, A any]    = readeroption.ReaderOption[R, A]

	Endomorphism[A any] = endomorphism.Endomorphism[A]

	Consumer[A any] = consumer.Consumer[A]

	Prism[S, T any] = prism.Prism[S, T]
	Lens[S, T any]  = lens.Lens[S, T]

	Trampoline[B, L any] = tailrec.Trampoline[B, L]

	Predicate[A any] = predicate.Predicate[A]

	Pair[A, B any] = pair.Pair[A, B]

	IORef[A any] = ioref.IORef[A]

	State[S, A any] = state.State[S, A]

	Void = function.Void
)
