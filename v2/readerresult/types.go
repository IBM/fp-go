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
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/result"
)

type (
	Endomorphism[A any] = endomorphism.Endomorphism[A]
	Lazy[A any]         = lazy.Lazy[A]
	Option[A any]       = option.Option[A]
	Either[E, A any]    = either.Either[E, A]
	Result[A any]       = result.Result[A]
	Reader[R, A any]    = reader.Reader[R, A]

	ReaderResult[R, A any] = Reader[R, Result[A]]

	Kleisli[R, A, B any]  = Reader[A, ReaderResult[R, B]]
	Operator[R, A, B any] = Kleisli[R, ReaderResult[R, A], B]
)
