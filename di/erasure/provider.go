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

package erasure

import (
	"fmt"

	A "github.com/IBM/fp-go/array"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IO "github.com/IBM/fp-go/io"
	IOG "github.com/IBM/fp-go/io/generic"
	IOE "github.com/IBM/fp-go/ioeither"
	O "github.com/IBM/fp-go/option"
	RIOE "github.com/IBM/fp-go/readerioeither"
	R "github.com/IBM/fp-go/record"
	T "github.com/IBM/fp-go/tuple"
)

type InjectableFactory = RIOE.ReaderIOEither[Token, error, any]
type ProviderFactory = RIOE.ReaderIOEither[InjectableFactory, error, any]

type Provider interface {
	fmt.Stringer
	Provides() Token
	Factory() ProviderFactory
}

type provider struct {
	provides Token
	factory  ProviderFactory
}

func (p *provider) Provides() Token {
	return p.provides
}

func (p *provider) Factory() ProviderFactory {
	return p.factory
}

func (p *provider) String() string {
	return fmt.Sprintf("Provider for [%s]", p.provides)
}

func MakeProvider(token Token, fct ProviderFactory) Provider {
	return &provider{token, fct}
}

func mapFromToken(idx int, token Token) map[TokenType]map[int]int {
	return R.Singleton(token.Type(), R.Singleton(idx, idx))
}

var mergeTokenMaps = R.UnionMonoid[TokenType](R.UnionLastSemigroup[int, int]())
var foldDeps = A.FoldMapWithIndex[Token](mergeTokenMaps)(mapFromToken)

var lookupMandatory = R.Lookup[map[int]int](Mandatory)
var lookupOption = R.Lookup[map[int]int](Option)

type Mapping = map[TokenType]map[int]int

func getAt[T any](ar []T) func(idx int) T {
	return func(idx int) T {
		return ar[idx]
	}
}

func handleMandatory(mp Mapping) func(res []IOE.IOEither[error, any]) IOE.IOEither[error, map[int]any] {

	onNone := F.Nullary2(R.Empty[int, any], IOE.Of[error, map[int]any])

	return func(res []IOE.IOEither[error, any]) IOE.IOEither[error, map[int]any] {
		return F.Pipe2(
			mp,
			lookupMandatory,
			O.Fold(
				onNone,
				IOE.TraverseRecord[int](getAt(res)),
			),
		)
	}
}

func handleOption(mp Mapping) func(res []IOE.IOEither[error, any]) IO.IO[map[int]O.Option[any]] {

	onNone := F.Nullary2(R.Empty[int, O.Option[any]], IO.Of[map[int]O.Option[any]])

	return func(res []IOE.IOEither[error, any]) IO.IO[map[int]O.Option[any]] {

		return F.Pipe2(
			mp,
			lookupOption,
			O.Fold(
				onNone,
				F.Flow2(
					IOG.TraverseRecord[IO.IO[map[int]E.Either[error, any]], map[int]int](getAt(res)),
					IO.Map(R.Map[int](E.ToOption[error, any])),
				),
			),
		)
	}
}

func mergeArguments(count int) func(
	mandatory IOE.IOEither[error, map[int]any],
	optonal IO.IO[map[int]O.Option[any]],
) IOE.IOEither[error, []any] {

	optMapToAny := R.Map[int](F.ToAny[O.Option[any]])
	mergeMaps := R.UnionLastMonoid[int, any]()

	return func(
		mandatory IOE.IOEither[error, map[int]any],
		optional IO.IO[map[int]O.Option[any]],
	) IOE.IOEither[error, []any] {

		return F.Pipe1(
			IOE.SequenceT2(mandatory, IOE.FromIO[error](optional)),
			IOE.Map[error](T.Tupled2(func(mnd map[int]any, opt map[int]O.Option[any]) []any {
				// merge all parameters
				merged := mergeMaps.Concat(mnd, optMapToAny(opt))

				return R.ReduceWithIndex(func(idx int, res []any, value any) []any {
					res[idx] = value
					return res
				}, make([]any, count))(merged)
			})),
		)
	}
}

func MakeProviderFactory(
	deps []Token,
	fct func(param ...any) IOE.IOEither[error, any]) ProviderFactory {

	mapping := foldDeps(deps)

	mandatory := handleMandatory(mapping)
	optional := handleOption(mapping)

	merge := mergeArguments(A.Size(deps))

	f := F.Unvariadic0(fct)

	return func(inj InjectableFactory) IOE.IOEither[error, any] {
		// resolve all dependencies
		resolved := A.MonadMap(deps, inj)
		// resolve dependencies
		return F.Pipe1(
			merge(mandatory(resolved), optional(resolved)),
			IOE.Chain(f),
		)
	}
}
