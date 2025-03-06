// Copyright (c) 2023 IBM Corp.
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

package readerioeither

import (
	"context"

	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
	RE "github.com/IBM/fp-go/v2/readerioeither"
)

type (
	// ReaderIOEither is a specialization of the [RE.ReaderIOEither] monad for the typical golang scenario in which the
	// left value is an [error] and the context is a [context.Context]
	Either[A any]         = either.Either[error, A]
	Lazy[A any]           = lazy.Lazy[A]
	IO[A any]             = io.IO[A]
	IOEither[A any]       = ioeither.IOEither[error, A]
	Reader[R, A any]      = reader.Reader[R, A]
	ReaderIO[R, A any]    = readerio.ReaderIO[R, A]
	ReaderIOEither[A any] = RE.ReaderIOEither[context.Context, error, A]

	Operator[A, B any] = Reader[ReaderIOEither[A], ReaderIOEither[B]]
)
