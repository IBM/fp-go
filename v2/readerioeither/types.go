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
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
)

type (
	Either[E, A any]            = either.Either[E, A]
	Reader[R, A any]            = reader.Reader[R, A]
	ReaderIO[R, A any]          = readerio.ReaderIO[R, A]
	IOEither[E, A any]          = ioeither.IOEither[E, A]
	ReaderIOEither[R, E, A any] = Reader[R, IOEither[E, A]]

	Operator[R, E, A, B any] = reader.Reader[ReaderIOEither[R, E, A], ReaderIOEither[R, E, B]]
)
