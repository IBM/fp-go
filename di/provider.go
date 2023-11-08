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
	A "github.com/IBM/fp-go/array"
	DIE "github.com/IBM/fp-go/di/erasure"
	E "github.com/IBM/fp-go/either"
	ER "github.com/IBM/fp-go/erasure"
	"github.com/IBM/fp-go/errors"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	T "github.com/IBM/fp-go/tuple"
)

func lookupAt[T any](idx int) func(params []any) E.Either[error, T] {
	return F.Flow3(
		A.Lookup[any](idx),
		E.FromOption[any](errors.OnNone("No parameter at position %d", idx)),
		E.Chain(ER.SafeUnerase[T]),
	)
}

func eraseProviderFactory0[R any](f func() IOE.IOEither[error, R]) func(params ...any) IOE.IOEither[error, any] {
	return func(params ...any) IOE.IOEither[error, any] {
		return F.Pipe1(
			f(),
			IOE.Map[error](ER.Erase[R]),
		)
	}
}

func eraseProviderFactory1[T1 any, R any](
	f func(T1) IOE.IOEither[error, R]) func(params ...any) IOE.IOEither[error, any] {
	ft := T.Tupled1(f)
	t1 := lookupAt[T1](0)
	return func(params ...any) IOE.IOEither[error, any] {
		return F.Pipe3(
			E.SequenceT1(t1(params)),
			IOE.FromEither[error, T.Tuple1[T1]],
			IOE.Chain(ft),
			IOE.Map[error](ER.Erase[R]),
		)
	}
}

func MakeProvider0[R any](
	token InjectionToken[R],
	fct func() IOE.IOEither[error, R],
) DIE.Provider {
	return DIE.MakeProvider(
		token,
		DIE.MakeProviderFactory(
			A.Empty[DIE.Token](),
			eraseProviderFactory0(fct),
		),
	)
}

func MakeProvider1[T1, R any](
	token InjectionToken[R],
	d1 InjectionToken[T1],
	fct func(T1) IOE.IOEither[error, R],
) DIE.Provider {

	return DIE.MakeProvider(
		token,
		DIE.MakeProviderFactory(
			A.From[DIE.Token](d1),
			eraseProviderFactory1(fct),
		),
	)
}
