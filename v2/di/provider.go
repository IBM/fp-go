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

package di

import (
	A "github.com/IBM/fp-go/v2/array"
	DIE "github.com/IBM/fp-go/v2/di/erasure"
	E "github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/errors"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/ioresult"
)

func lookupAt[T any](idx int, token Dependency[T]) func(params []any) Result[T] {
	return F.Flow3(
		A.Lookup[any](idx),
		E.FromOption[any](errors.OnNone("No parameter at position %d", idx)),
		E.Chain(token.Unerase),
	)
}

func eraseTuple[A, R any](f func(A) IOResult[R]) func(Result[A]) IOResult[any] {
	return F.Flow3(
		IOE.FromEither[error, A],
		IOE.Chain(f),
		IOE.Map[error](F.ToAny[R]),
	)
}

func eraseProviderFactory0[R any](f IOResult[R]) func(params ...any) IOResult[any] {
	return func(_ ...any) IOResult[any] {
		return F.Pipe1(
			f,
			IOE.Map[error](F.ToAny[R]),
		)
	}
}

func MakeProviderFactory0[R any](
	fct IOResult[R],
) DIE.ProviderFactory {
	return DIE.MakeProviderFactory(
		A.Empty[DIE.Dependency](),
		eraseProviderFactory0(fct),
	)
}

// MakeTokenWithDefault0 creates a unique [InjectionToken] for a specific type with an attached default [DIE.Provider]
func MakeTokenWithDefault0[R any](name string, fct IOResult[R]) InjectionToken[R] {
	return MakeTokenWithDefault[R](name, MakeProviderFactory0(fct))
}

func MakeProvider0[R any](
	token InjectionToken[R],
	fct IOResult[R],
) DIE.Provider {
	return DIE.MakeProvider(
		token,
		MakeProviderFactory0(fct),
	)
}

// ConstProvider simple implementation for a provider with a constant value
func ConstProvider[R any](token InjectionToken[R], value R) DIE.Provider {
	return MakeProvider0(token, ioresult.Of(value))
}
