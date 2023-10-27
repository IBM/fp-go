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

package fromeither

import (
	ET "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	C "github.com/IBM/fp-go/internal/chain"
	O "github.com/IBM/fp-go/option"
)

func FromOption[A, HKTEA, E any](fromEither func(ET.Either[E, A]) HKTEA, onNone func() E) func(ma O.Option[A]) HKTEA {
	return F.Flow2(ET.FromOption[A](onNone), fromEither)
}

func FromPredicate[E, A, HKTEA any](fromEither func(ET.Either[E, A]) HKTEA, pred func(A) bool, onFalse func(A) E) func(A) HKTEA {
	return F.Flow2(ET.FromPredicate(pred, onFalse), fromEither)
}

func MonadFromOption[E, A, HKTEA any](
	fromEither func(ET.Either[E, A]) HKTEA,
	onNone func() E,
	ma O.Option[A],
) HKTEA {
	return F.Pipe1(
		O.MonadFold(
			ma,
			F.Nullary2(onNone, ET.Left[A, E]),
			ET.Right[E, A],
		),
		fromEither,
	)
}

func FromOptionK[A, E, B, HKTEB any](
	fromEither func(ET.Either[E, B]) HKTEB,
	onNone func() E) func(f func(A) O.Option[B]) func(A) HKTEB {
	// helper
	return F.Bind2nd(F.Flow2[func(A) O.Option[B], func(O.Option[B]) HKTEB, A, O.Option[B], HKTEB], FromOption(fromEither, onNone))
}

func MonadChainEitherK[A, E, B, HKTEA, HKTEB any](
	mchain func(HKTEA, func(A) HKTEB) HKTEB,
	fromEither func(ET.Either[E, B]) HKTEB,
	ma HKTEA,
	f func(A) ET.Either[E, B]) HKTEB {
	return mchain(ma, F.Flow2(f, fromEither))
}

func ChainOptionK[A, E, B, HKTEA, HKTEB any](
	mchain func(HKTEA, func(A) HKTEB) HKTEB,
	fromEither func(ET.Either[E, B]) HKTEB,
	onNone func() E,
) func(f func(A) O.Option[B]) func(ma HKTEA) HKTEB {
	return F.Flow2(FromOptionK[A](fromEither, onNone), F.Bind1st(F.Bind2nd[HKTEA, func(A) HKTEB, HKTEB], mchain))
}

func MonadChainFirstEitherK[A, E, B, HKTEA, HKTEB any](
	mchain func(HKTEA, func(A) HKTEA) HKTEA,
	mmap func(HKTEB, func(B) A) HKTEA,
	fromEither func(ET.Either[E, B]) HKTEB,
	ma HKTEA,
	f func(A) ET.Either[E, B]) HKTEA {
	return C.MonadChainFirst(mchain, mmap, ma, F.Flow2(f, fromEither))
}

func ChainFirstEitherK[A, E, B, HKTEA, HKTEB any](
	mchain func(HKTEA, func(A) HKTEA) HKTEA,
	mmap func(HKTEB, func(B) A) HKTEA,
	fromEither func(ET.Either[E, B]) HKTEB,
	f func(A) ET.Either[E, B]) func(HKTEA) HKTEA {
	return C.ChainFirst(mchain, mmap, F.Flow2(f, fromEither))
}
