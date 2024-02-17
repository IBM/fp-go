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
	ET "github.com/IBM/fp-go/either"
	ST "github.com/IBM/fp-go/internal/statet"
	P "github.com/IBM/fp-go/pair"
	G "github.com/IBM/fp-go/readerioeither/generic"
)

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
