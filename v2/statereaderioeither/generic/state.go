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

package generic

import (
	ET "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	ST "github.com/IBM/fp-go/v2/internal/statet"
	P "github.com/IBM/fp-go/v2/pair"
	G "github.com/IBM/fp-go/v2/readerioeither/generic"
)

func Left[
	SRIOEA ~func(S) RIOEA,
	RIOEA ~func(R) IOEA,
	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	S, R, E, A any,
](e E) SRIOEA {
	return F.Constant1[S](G.Left[RIOEA](e))
}

func Right[
	SRIOEA ~func(S) RIOEA,
	RIOEA ~func(R) IOEA,
	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	S, R, E, A any,
](a A) SRIOEA {
	return ST.Of[SRIOEA](
		G.Of[RIOEA],
		a,
	)
}

func Of[
	SRIOEA ~func(S) RIOEA,
	RIOEA ~func(R) IOEA,
	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	S, R, E, A any,
](a A) SRIOEA {
	return Right[SRIOEA](a)
}

func MonadMap[
	SRIOEA ~func(S) RIOEA,
	SRIOEB ~func(S) RIOEB,
	RIOEA ~func(R) IOEA,
	RIOEB ~func(R) IOEB,
	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	IOEB ~func() ET.Either[E, P.Pair[B, S]],
	S, R, E, A, B any,
](fa SRIOEA, f func(A) B) SRIOEB {
	return ST.MonadMap[SRIOEA, SRIOEB](
		G.MonadMap[RIOEA, RIOEB],
		fa,
		f,
	)
}

func Map[
	SRIOEA ~func(S) RIOEA,
	SRIOEB ~func(S) RIOEB,
	RIOEA ~func(R) IOEA,
	RIOEB ~func(R) IOEB,
	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	IOEB ~func() ET.Either[E, P.Pair[B, S]],
	S, R, E, A, B any,
](f func(A) B) func(SRIOEA) SRIOEB {
	return ST.Map[SRIOEA, SRIOEB](
		G.Map[RIOEA, RIOEB],
		f,
	)
}

func MonadChain[
	SRIOEA ~func(S) RIOEA,
	SRIOEB ~func(S) RIOEB,
	RIOEA ~func(R) IOEA,
	RIOEB ~func(R) IOEB,
	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	IOEB ~func() ET.Either[E, P.Pair[B, S]],
	S, R, E, A, B any,
](fa SRIOEA, f func(A) SRIOEB) SRIOEB {
	return ST.MonadChain[SRIOEA, SRIOEB](
		G.MonadChain[RIOEA, RIOEB],
		fa,
		f,
	)
}

func Chain[
	SRIOEA ~func(S) RIOEA,
	SRIOEB ~func(S) RIOEB,
	RIOEA ~func(R) IOEA,
	RIOEB ~func(R) IOEB,
	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	IOEB ~func() ET.Either[E, P.Pair[B, S]],
	S, R, E, A, B any,
](f func(A) SRIOEB) func(SRIOEA) SRIOEB {
	return ST.Chain[SRIOEA, SRIOEB](
		G.Chain[RIOEA, RIOEB],
		f,
	)
}

func MonadAp[
	SRIOEA ~func(S) RIOEA,
	SRIOEB ~func(S) RIOEB,
	SRIOEAB ~func(S) RIOEAB,
	RIOEA ~func(R) IOEA,
	RIOEB ~func(R) IOEB,
	RIOEAB ~func(R) IOEAB,
	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	IOEB ~func() ET.Either[E, P.Pair[B, S]],
	IOEAB ~func() ET.Either[E, P.Pair[func(A) B, S]],
	S, R, E, A, B any,
](fab SRIOEAB, fa SRIOEA) SRIOEB {
	return ST.MonadAp[SRIOEA, SRIOEB, SRIOEAB](
		G.MonadMap[RIOEA, RIOEB],
		G.MonadChain[RIOEAB, RIOEB],
		fab,
		fa,
	)
}

func Ap[
	SRIOEA ~func(S) RIOEA,
	SRIOEB ~func(S) RIOEB,
	SRIOEAB ~func(S) RIOEAB,
	RIOEA ~func(R) IOEA,
	RIOEB ~func(R) IOEB,
	RIOEAB ~func(R) IOEAB,
	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	IOEB ~func() ET.Either[E, P.Pair[B, S]],
	IOEAB ~func() ET.Either[E, P.Pair[func(A) B, S]],
	S, R, E, A, B any,
](fa SRIOEA) func(SRIOEAB) SRIOEB {
	return ST.Ap[SRIOEA, SRIOEB, SRIOEAB](
		G.Map[RIOEA, RIOEB],
		G.Chain[RIOEAB, RIOEB],
		fa,
	)
}

// Conversions

func FromReaderIOEither[
	SRIOEA ~func(S) RIOEA,
	RIOEA_IN ~func(R) IOEA_IN,

	RIOEA ~func(R) IOEA,

	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	IOEA_IN ~func() ET.Either[E, A],

	S, R, E, A any,
](fa RIOEA_IN) SRIOEA {
	return ST.FromF[SRIOEA](
		G.MonadMap[RIOEA_IN, RIOEA],
		fa,
	)
}

func FromReaderEither[
	SRIOEA ~func(S) RIOEA,
	RIOEA_IN ~func(R) IOEA_IN,
	RIOEA ~func(R) IOEA,

	REA_IN ~func(R) ET.Either[E, A],

	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	IOEA_IN ~func() ET.Either[E, A],

	S, R, E, A any,
](fa REA_IN) SRIOEA {
	return FromReaderIOEither[SRIOEA](G.FromReaderEither[REA_IN, RIOEA_IN](fa))
}

func FromIOEither[
	SRIOEA ~func(S) RIOEA,
	RIOEA_IN ~func(R) IOEA_IN,
	IOEA_IN ~func() ET.Either[E, A],
	RIOEA ~func(R) IOEA,

	IOEA ~func() ET.Either[E, P.Pair[A, S]],

	S, R, E, A any,
](fa IOEA_IN) SRIOEA {
	return FromReaderIOEither[SRIOEA](G.FromIOEither[RIOEA_IN](fa))
}

func FromIO[
	SRIOEA ~func(S) RIOEA,
	RIOEA_IN ~func(R) IOEA_IN,

	IO_IN ~func() A,

	RIOEA ~func(R) IOEA,

	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	IOEA_IN ~func() ET.Either[E, A],

	S, R, E, A any,
](fa IO_IN) SRIOEA {
	return FromReaderIOEither[SRIOEA](G.FromIO[RIOEA_IN](fa))
}

func FromReader[
	SRIOEA ~func(S) RIOEA,
	RIOEA_IN ~func(R) IOEA_IN,

	R_IN ~func(R) A,

	RIOEA ~func(R) IOEA,

	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	IOEA_IN ~func() ET.Either[E, A],

	S, R, E, A any,
](fa R_IN) SRIOEA {
	return FromReaderIOEither[SRIOEA](G.FromReader[R_IN, RIOEA_IN](fa))
}

func FromEither[
	SRIOEA ~func(S) RIOEA,

	RIOEA ~func(R) IOEA,

	IOEA ~func() ET.Either[E, P.Pair[A, S]],

	S, R, E, A any,
](ma ET.Either[E, A]) SRIOEA {
	return ET.MonadFold(ma, Left[SRIOEA], Right[SRIOEA])
}

func FromState[
	SRIOEA ~func(S) RIOEA,
	STATE ~func(S) P.Pair[A, S],
	RIOEA ~func(R) IOEA,

	IOEA ~func() ET.Either[E, P.Pair[A, S]],

	S, R, E, A any,
](fa STATE) SRIOEA {
	return ST.FromState[SRIOEA](G.Of[RIOEA], fa)
}

// Combinators

func Local[
	SR1IOEA ~func(S) R1IOEA,
	SR2IOEA ~func(S) R2IOEA,
	R1IOEA ~func(R1) IOEA,
	R2IOEA ~func(R2) IOEA,
	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	S, R1, R2, E, A any,
](f func(R2) R1) func(SR1IOEA) SR2IOEA {
	return func(ma SR1IOEA) SR2IOEA {
		return F.Flow2(ma, G.Local[R1IOEA, R2IOEA](f))
	}
}

func Asks[
	SRIOEA ~func(S) RIOEA,
	RIOEA ~func(R) IOEA,
	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	S, R, E, A any,
](f func(R) SRIOEA) SRIOEA {
	return func(s S) RIOEA {
		return func(r R) IOEA {
			return f(r)(s)(r)
		}
	}
}

func FromIOEitherK[
	SRIOEB ~func(S) RIOEB,
	RIOEB_IN ~func(R) IOEB_IN,
	IOEB_IN ~func() ET.Either[E, B],
	RIOEB ~func(R) IOEB,
	IOEB ~func() ET.Either[E, P.Pair[B, S]],
	S, R, E, A, B any,
](f func(A) IOEB_IN) func(A) SRIOEB {
	return F.Flow2(
		f,
		FromIOEither[SRIOEB, RIOEB_IN],
	)
}

func FromEitherK[
	SRIOEB ~func(S) RIOEB,
	RIOEB ~func(R) IOEB,
	IOEB ~func() ET.Either[E, P.Pair[B, S]],
	S, R, E, A, B any,
](f func(A) ET.Either[E, B]) func(A) SRIOEB {
	return F.Flow2(
		f,
		FromEither[SRIOEB],
	)
}

func FromIOK[
	SRIOEB ~func(S) RIOEB,
	RIOEB_IN ~func(R) IOEB_IN,

	IOB_IN ~func() B,

	RIOEB ~func(R) IOEB,

	IOEB ~func() ET.Either[E, P.Pair[B, S]],
	IOEB_IN ~func() ET.Either[E, B],

	S, R, E, A, B any,
](f func(A) IOB_IN) func(A) SRIOEB {
	return F.Flow2(
		f,
		FromIO[SRIOEB, RIOEB_IN, IOB_IN],
	)
}

func FromReaderIOEitherK[
	SRIOEB ~func(S) RIOEB,
	RIOEB_IN ~func(R) IOEB_IN,
	IOEB_IN ~func() ET.Either[E, B],
	RIOEB ~func(R) IOEB,
	IOEB ~func() ET.Either[E, P.Pair[B, S]],
	S, R, E, A, B any,
](f func(A) RIOEB_IN) func(A) SRIOEB {
	return F.Flow2(
		f,
		FromReaderIOEither[SRIOEB, RIOEB_IN],
	)
}

func MonadChainReaderIOEitherK[
	SRIOEA ~func(S) RIOEA,
	SRIOEB ~func(S) RIOEB,
	RIOEB_IN ~func(R) IOEB_IN,
	RIOEA ~func(R) IOEA,
	RIOEB ~func(R) IOEB,
	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	IOEB ~func() ET.Either[E, P.Pair[B, S]],
	IOEB_IN ~func() ET.Either[E, B],
	S, R, E, A, B any,
](ma SRIOEA, f func(A) RIOEB_IN) SRIOEB {
	return MonadChain(ma, FromReaderIOEitherK[SRIOEB, RIOEB_IN](f))
}

func ChainReaderIOEitherK[
	SRIOEA ~func(S) RIOEA,
	SRIOEB ~func(S) RIOEB,
	RIOEB_IN ~func(R) IOEB_IN,
	RIOEA ~func(R) IOEA,
	RIOEB ~func(R) IOEB,
	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	IOEB ~func() ET.Either[E, P.Pair[B, S]],
	IOEB_IN ~func() ET.Either[E, B],
	S, R, E, A, B any,
](f func(A) RIOEB_IN) func(SRIOEA) SRIOEB {
	return Chain[SRIOEA](FromReaderIOEitherK[SRIOEB, RIOEB_IN](f))
}

func MonadChainIOEitherK[
	SRIOEA ~func(S) RIOEA,
	SRIOEB ~func(S) RIOEB,
	RIOEB_IN ~func(R) IOEB_IN,
	IOEB_IN ~func() ET.Either[E, B],
	RIOEA ~func(R) IOEA,
	RIOEB ~func(R) IOEB,
	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	IOEB ~func() ET.Either[E, P.Pair[B, S]],
	S, R, E, A, B any,
](ma SRIOEA, f func(A) IOEB_IN) SRIOEB {
	return MonadChain(ma, FromIOEitherK[SRIOEB, RIOEB_IN](f))
}

func ChainIOEitherK[
	SRIOEA ~func(S) RIOEA,
	SRIOEB ~func(S) RIOEB,
	RIOEB_IN ~func(R) IOEB_IN,
	IOEB_IN ~func() ET.Either[E, B],
	RIOEA ~func(R) IOEA,
	RIOEB ~func(R) IOEB,
	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	IOEB ~func() ET.Either[E, P.Pair[B, S]],
	S, R, E, A, B any,
](f func(A) IOEB_IN) func(SRIOEA) SRIOEB {
	return Chain[SRIOEA](FromIOEitherK[SRIOEB, RIOEB_IN](f))
}

func MonadChainEitherK[
	SRIOEA ~func(S) RIOEA,
	SRIOEB ~func(S) RIOEB,
	RIOEA ~func(R) IOEA,
	RIOEB ~func(R) IOEB,
	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	IOEB ~func() ET.Either[E, P.Pair[B, S]],
	S, R, E, A, B any,
](ma SRIOEA, f func(A) ET.Either[E, B]) SRIOEB {
	return MonadChain(ma, FromEitherK[SRIOEB](f))
}

func ChainEitherK[
	SRIOEA ~func(S) RIOEA,
	SRIOEB ~func(S) RIOEB,
	RIOEA ~func(R) IOEA,
	RIOEB ~func(R) IOEB,
	IOEA ~func() ET.Either[E, P.Pair[A, S]],
	IOEB ~func() ET.Either[E, P.Pair[B, S]],
	S, R, E, A, B any,
](f func(A) ET.Either[E, B]) func(SRIOEA) SRIOEB {
	return Chain[SRIOEA](FromEitherK[SRIOEB](f))
}
