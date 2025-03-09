// Copyright (c) 2024 IBM Corp.
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

package statereaderioeither

import (
	"github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readereither"
	"github.com/IBM/fp-go/v2/readerioeither"
	"github.com/IBM/fp-go/v2/state"
)

type (
	State[S, A any]                     = state.State[S, A]
	Pair[L, R any]                      = pair.Pair[L, R]
	Reader[R, A any]                    = reader.Reader[R, A]
	Either[E, A any]                    = either.Either[E, A]
	IO[A any]                           = io.IO[A]
	IOEither[E, A any]                  = ioeither.IOEither[E, A]
	ReaderIOEither[R, E, A any]         = readerioeither.ReaderIOEither[R, E, A]
	ReaderEither[R, E, A any]           = readereither.ReaderEither[R, E, A]
	StateReaderIOEither[S, R, E, A any] = Reader[S, ReaderIOEither[R, E, Pair[S, A]]]
	Operator[S, R, E, A, B any]         = Reader[StateReaderIOEither[S, R, E, A], StateReaderIOEither[S, R, E, B]]
)
