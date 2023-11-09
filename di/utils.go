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
package di

import (
	E "github.com/IBM/fp-go/either"
	"github.com/IBM/fp-go/errors"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	IOO "github.com/IBM/fp-go/iooption"
	O "github.com/IBM/fp-go/option"
)

// toType converts and any to a T
func toType[T any]() func(t any) E.Either[error, T] {
	return E.ToType[T](errors.OnSome[any]("Value of type [%T] cannot be converted."))
}

// toOptionType converts an any to an Option[any] and then to an Option[T]
func toOptionType[T any]() func(t any) E.Either[error, O.Option[T]] {
	return F.Flow2(
		toType[O.Option[any]](),
		E.Chain(O.Fold(
			F.Nullary2(O.None[T], E.Of[error, O.Option[T]]),
			F.Flow2(
				toType[T](),
				E.Map[error](O.Of[T]),
			),
		)),
	)
}

// toIOEitherType converts an any to an IOEither[error, any] and then to an IOEither[error, T]
func toIOEitherType[T any]() func(t any) E.Either[error, IOE.IOEither[error, T]] {
	return F.Flow2(
		toType[IOE.IOEither[error, any]](),
		E.Map[error](IOE.ChainEitherK(toType[T]())),
	)
}

// toIOOptionType converts an any to an IOOption[any] and then to an IOOption[T]
func toIOOptionType[T any]() func(t any) E.Either[error, IOO.IOOption[T]] {
	return F.Flow2(
		toType[IOO.IOOption[any]](),
		E.Map[error](IOO.ChainOptionK(F.Flow2(
			toType[T](),
			E.ToOption[error, T],
		))),
	)
}
