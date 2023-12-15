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
	"github.com/IBM/fp-go/errors"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	T "github.com/IBM/fp-go/tuple"
)

func lookupAt[T any](idx int, token Dependency[T]) func(params []any) E.Either[error, T] {
	return F.Flow3(
		A.Lookup[any](idx),
		E.FromOption[any](errors.OnNone("No parameter at position %d", idx)),
		E.Chain(token.Unerase),
	)
}

func eraseProviderFactory0[R any](f func() IOE.IOEither[error, R]) func(params ...any) IOE.IOEither[error, any] {
	return func(params ...any) IOE.IOEither[error, any] {
		return F.Pipe1(
			f(),
			IOE.Map[error](F.ToAny[R]),
		)
	}
}

func eraseProviderFactory1[T1 any, R any](
	d1 Dependency[T1],
	f func(T1) IOE.IOEither[error, R]) func(params ...any) IOE.IOEither[error, any] {
	ft := T.Tupled1(f)
	t1 := lookupAt[T1](0, d1)
	return func(params ...any) IOE.IOEither[error, any] {
		return F.Pipe3(
			E.SequenceT1(t1(params)),
			IOE.FromEither[error, T.Tuple1[T1]],
			IOE.Chain(ft),
			IOE.Map[error](F.ToAny[R]),
		)
	}
}

func eraseProviderFactory2[T1, T2 any, R any](
	d1 Dependency[T1],
	d2 Dependency[T2],
	f func(T1, T2) IOE.IOEither[error, R]) func(params ...any) IOE.IOEither[error, any] {
	ft := T.Tupled2(f)
	t1 := lookupAt[T1](0, d1)
	t2 := lookupAt[T2](1, d2)
	return func(params ...any) IOE.IOEither[error, any] {
		return F.Pipe3(
			E.SequenceT2(t1(params), t2(params)),
			IOE.FromEither[error, T.Tuple2[T1, T2]],
			IOE.Chain(ft),
			IOE.Map[error](F.ToAny[R]),
		)
	}
}

func eraseProviderFactory3[T1, T2, T3 any, R any](
	d1 Dependency[T1],
	d2 Dependency[T2],
	d3 Dependency[T3],
	f func(T1, T2, T3) IOE.IOEither[error, R]) func(params ...any) IOE.IOEither[error, any] {
	ft := T.Tupled3(f)
	t1 := lookupAt[T1](0, d1)
	t2 := lookupAt[T2](1, d2)
	t3 := lookupAt[T3](2, d3)
	return func(params ...any) IOE.IOEither[error, any] {
		return F.Pipe3(
			E.SequenceT3(t1(params), t2(params), t3(params)),
			IOE.FromEither[error, T.Tuple3[T1, T2, T3]],
			IOE.Chain(ft),
			IOE.Map[error](F.ToAny[R]),
		)
	}
}

func eraseProviderFactory4[T1, T2, T3, T4 any, R any](
	d1 Dependency[T1],
	d2 Dependency[T2],
	d3 Dependency[T3],
	d4 Dependency[T4],
	f func(T1, T2, T3, T4) IOE.IOEither[error, R]) func(params ...any) IOE.IOEither[error, any] {
	ft := T.Tupled4(f)
	t1 := lookupAt[T1](0, d1)
	t2 := lookupAt[T2](1, d2)
	t3 := lookupAt[T3](2, d3)
	t4 := lookupAt[T4](3, d4)
	return func(params ...any) IOE.IOEither[error, any] {
		return F.Pipe3(
			E.SequenceT4(t1(params), t2(params), t3(params), t4(params)),
			IOE.FromEither[error, T.Tuple4[T1, T2, T3, T4]],
			IOE.Chain(ft),
			IOE.Map[error](F.ToAny[R]),
		)
	}
}

func MakeProviderFactory0[R any](
	fct func() IOE.IOEither[error, R],
) DIE.ProviderFactory {
	return DIE.MakeProviderFactory(
		A.Empty[DIE.Dependency](),
		eraseProviderFactory0(fct),
	)
}

// MakeTokenWithDefault0 create a unique `InjectionToken` for a specific type with an attached default provider
func MakeTokenWithDefault0[R any](name string, fct func() IOE.IOEither[error, R]) InjectionToken[R] {
	return MakeTokenWithDefault[R](name, MakeProviderFactory0(fct))
}

func MakeProvider0[R any](
	token InjectionToken[R],
	fct func() IOE.IOEither[error, R],
) DIE.Provider {
	return DIE.MakeProvider(
		token,
		MakeProviderFactory0(fct),
	)
}

func MakeProviderFactory1[T1, R any](
	d1 Dependency[T1],
	fct func(T1) IOE.IOEither[error, R],
) DIE.ProviderFactory {

	return DIE.MakeProviderFactory(
		A.From[DIE.Dependency](d1),
		eraseProviderFactory1(d1, fct),
	)
}

// MakeTokenWithDefault1 create a unique `InjectionToken` for a specific type with an attached default provider
func MakeTokenWithDefault1[T1, R any](name string,
	d1 Dependency[T1],
	fct func(T1) IOE.IOEither[error, R]) InjectionToken[R] {
	return MakeTokenWithDefault[R](name, MakeProviderFactory1(d1, fct))
}

func MakeProvider1[T1, R any](
	token InjectionToken[R],
	d1 Dependency[T1],
	fct func(T1) IOE.IOEither[error, R],
) DIE.Provider {

	return DIE.MakeProvider(
		token,
		MakeProviderFactory1(d1, fct),
	)
}

func MakeProviderFactory2[T1, T2, R any](
	d1 Dependency[T1],
	d2 Dependency[T2],
	fct func(T1, T2) IOE.IOEither[error, R],
) DIE.ProviderFactory {

	return DIE.MakeProviderFactory(
		A.From[DIE.Dependency](d1, d2),
		eraseProviderFactory2(d1, d2, fct),
	)
}

// MakeTokenWithDefault2 create a unique `InjectionToken` for a specific type with an attached default provider
func MakeTokenWithDefault2[T1, T2, R any](name string,
	d1 Dependency[T1],
	d2 Dependency[T2],
	fct func(T1, T2) IOE.IOEither[error, R]) InjectionToken[R] {
	return MakeTokenWithDefault[R](name, MakeProviderFactory2(d1, d2, fct))
}

func MakeProvider2[T1, T2, R any](
	token InjectionToken[R],
	d1 Dependency[T1],
	d2 Dependency[T2],
	fct func(T1, T2) IOE.IOEither[error, R],
) DIE.Provider {

	return DIE.MakeProvider(
		token,
		MakeProviderFactory2(d1, d2, fct),
	)
}

func MakeProviderFactory3[T1, T2, T3, R any](
	d1 Dependency[T1],
	d2 Dependency[T2],
	d3 Dependency[T3],
	fct func(T1, T2, T3) IOE.IOEither[error, R],
) DIE.ProviderFactory {

	return DIE.MakeProviderFactory(
		A.From[DIE.Dependency](d1, d2, d3),
		eraseProviderFactory3(d1, d2, d3, fct),
	)
}

// MakeTokenWithDefault3 create a unique `InjectionToken` for a specific type with an attached default provider
func MakeTokenWithDefault3[T1, T2, T3, R any](name string,
	d1 Dependency[T1],
	d2 Dependency[T2],
	d3 Dependency[T3],
	fct func(T1, T2, T3) IOE.IOEither[error, R]) InjectionToken[R] {
	return MakeTokenWithDefault[R](name, MakeProviderFactory3(d1, d2, d3, fct))
}

func MakeProvider3[T1, T2, T3, R any](
	token InjectionToken[R],
	d1 Dependency[T1],
	d2 Dependency[T2],
	d3 Dependency[T3],
	fct func(T1, T2, T3) IOE.IOEither[error, R],
) DIE.Provider {

	return DIE.MakeProvider(
		token,
		MakeProviderFactory3(d1, d2, d3, fct),
	)
}

func MakeProviderFactory4[T1, T2, T3, T4, R any](
	d1 Dependency[T1],
	d2 Dependency[T2],
	d3 Dependency[T3],
	d4 Dependency[T4],
	fct func(T1, T2, T3, T4) IOE.IOEither[error, R],
) DIE.ProviderFactory {

	return DIE.MakeProviderFactory(
		A.From[DIE.Dependency](d1, d2, d3, d4),
		eraseProviderFactory4(d1, d2, d3, d4, fct),
	)
}

// MakeTokenWithDefault4 create a unique `InjectionToken` for a specific type with an attached default provider
func MakeTokenWithDefault4[T1, T2, T3, T4, R any](name string,
	d1 Dependency[T1],
	d2 Dependency[T2],
	d3 Dependency[T3],
	d4 Dependency[T4],
	fct func(T1, T2, T3, T4) IOE.IOEither[error, R]) InjectionToken[R] {
	return MakeTokenWithDefault[R](name, MakeProviderFactory4(d1, d2, d3, d4, fct))
}

func MakeProvider4[T1, T2, T3, T4, R any](
	token InjectionToken[R],
	d1 Dependency[T1],
	d2 Dependency[T2],
	d3 Dependency[T3],
	d4 Dependency[T4],
	fct func(T1, T2, T3, T4) IOE.IOEither[error, R],
) DIE.Provider {

	return DIE.MakeProvider(
		token,
		MakeProviderFactory4(d1, d2, d3, d4, fct),
	)
}

// ConstProvider simple implementation for a provider with a constant value
func ConstProvider[R any](token InjectionToken[R], value R) DIE.Provider {
	return MakeProvider0[R](token, F.Constant(IOE.Of[error](value)))
}
