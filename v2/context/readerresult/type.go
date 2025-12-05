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

// package readerresult implements a specialization of the Reader monad assuming a golang context as the context of the monad and a standard golang error
package readerresult

import (
	"context"

	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readereither"
	"github.com/IBM/fp-go/v2/result"
)

type (
	Option[A any]    = option.Option[A]
	Either[A any]    = either.Either[error, A]
	Result[A any]    = result.Result[A]
	Reader[R, A any] = reader.Reader[R, A]
	// ReaderResult is a specialization of the Reader monad for the typical golang scenario
	ReaderResult[A any] = readereither.ReaderEither[context.Context, error, A]

	Kleisli[A, B any]  = reader.Reader[A, ReaderResult[B]]
	Operator[A, B any] = Kleisli[ReaderResult[A], B]
)
