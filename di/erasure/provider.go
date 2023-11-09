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
	IOO "github.com/IBM/fp-go/iooption"
	Int "github.com/IBM/fp-go/number/integer"
	O "github.com/IBM/fp-go/option"
	RIOE "github.com/IBM/fp-go/readerioeither"
	R "github.com/IBM/fp-go/record"
)

type InjectableFactory = RIOE.ReaderIOEither[Dependency, error, any]
type ProviderFactory = RIOE.ReaderIOEither[InjectableFactory, error, any]
type paramIndex = map[int]int

type Provider interface {
	fmt.Stringer
	Provides() Dependency
	Factory() ProviderFactory
}

type provider struct {
	provides Dependency
	factory  ProviderFactory
}

func (p *provider) Provides() Dependency {
	return p.provides
}

func (p *provider) Factory() ProviderFactory {
	return p.factory
}

func (p *provider) String() string {
	return fmt.Sprintf("Provider for [%s]", p.provides)
}

func MakeProvider(token Dependency, fct ProviderFactory) Provider {
	return &provider{token, fct}
}

func mapFromToken(idx int, token Dependency) map[TokenType]paramIndex {
	return R.Singleton(token.Type(), R.Singleton(idx, idx))
}

var mergeTokenMaps = R.UnionMonoid[TokenType](R.UnionLastSemigroup[int, int]())
var foldDeps = A.FoldMapWithIndex[Dependency](mergeTokenMaps)(mapFromToken)

var lookupIdentity = R.Lookup[paramIndex](Identity)
var lookupOption = R.Lookup[paramIndex](Option)
var lookupIOEither = R.Lookup[paramIndex](IOEither)
var lookupIOOption = R.Lookup[paramIndex](IOOption)

type Mapping = map[TokenType]paramIndex

func getAt[T any](ar []T) func(idx int) T {
	return func(idx int) T {
		return ar[idx]
	}
}

type identityResult = IOE.IOEither[error, map[int]any]

func handleIdentity(mp Mapping) func(res []IOE.IOEither[error, any]) identityResult {

	onNone := F.Nullary2(R.Empty[int, any], IOE.Of[error, map[int]any])

	return func(res []IOE.IOEither[error, any]) identityResult {
		return F.Pipe2(
			mp,
			lookupIdentity,
			O.Fold(
				onNone,
				IOE.TraverseRecord[int](getAt(res)),
			),
		)
	}
}

type optionResult = IO.IO[map[int]O.Option[any]]

func handleOption(mp Mapping) func(res []IOE.IOEither[error, any]) optionResult {

	onNone := F.Nullary2(R.Empty[int, O.Option[any]], IO.Of[map[int]O.Option[any]])

	return func(res []IOE.IOEither[error, any]) optionResult {

		return F.Pipe2(
			mp,
			lookupOption,
			O.Fold(
				onNone,
				F.Flow2(
					IOG.TraverseRecord[IO.IO[map[int]E.Either[error, any]], paramIndex](getAt(res)),
					IO.Map(R.Map[int](E.ToOption[error, any])),
				),
			),
		)
	}
}

type ioeitherResult = IO.IO[map[int]IOE.IOEither[error, any]]

func handleIOEither(mp Mapping) func(res []IOE.IOEither[error, any]) ioeitherResult {

	onNone := F.Nullary2(R.Empty[int, IOE.IOEither[error, any]], IO.Of[map[int]IOE.IOEither[error, any]])

	return func(res []IOE.IOEither[error, any]) ioeitherResult {

		return F.Pipe2(
			mp,
			lookupIOEither,
			O.Fold(
				onNone,
				F.Flow2(
					R.Map[int](getAt(res)),
					IO.Of[map[int]IOE.IOEither[error, any]],
				),
			),
		)
	}
}

type iooptionResult = IO.IO[map[int]IOO.IOOption[any]]

func handleIOOption(mp Mapping) func(res []IOE.IOEither[error, any]) iooptionResult {

	onNone := F.Nullary2(R.Empty[int, IOO.IOOption[any]], IO.Of[map[int]IOO.IOOption[any]])

	return func(res []IOE.IOEither[error, any]) iooptionResult {

		return F.Pipe2(
			mp,
			lookupIOOption,
			O.Fold(
				onNone,
				F.Flow2(
					R.Map[int](F.Flow2(
						getAt(res),
						IOO.FromIOEither[error, any],
					)),
					IO.Of[map[int]IOO.IOOption[any]],
				),
			),
		)
	}
}

var optionMapToAny = R.Map[int](F.ToAny[O.Option[any]])
var ioeitherMapToAny = R.Map[int](F.ToAny[IOE.IOEither[error, any]])
var iooptionMapToAny = R.Map[int](F.ToAny[IOO.IOOption[any]])
var mergeMaps = R.UnionLastMonoid[int, any]()
var collectParams = R.CollectOrd[any, any](Int.Ord)(F.SK[int, any])

func mergeArguments(
	identity identityResult,
	option optionResult,
	ioeither ioeitherResult,
	iooption iooptionResult,
) IOE.IOEither[error, []any] {

	return F.Pipe2(
		A.From(
			identity,
			F.Pipe2(
				option,
				IO.Map(optionMapToAny),
				IOE.FromIO[error, map[int]any],
			),
			F.Pipe2(
				ioeither,
				IO.Map(ioeitherMapToAny),
				IOE.FromIO[error, map[int]any],
			),
			F.Pipe2(
				iooption,
				IO.Map(iooptionMapToAny),
				IOE.FromIO[error, map[int]any],
			),
		),
		IOE.SequenceArray[error, map[int]any],
		IOE.Map[error](F.Flow2(
			A.Fold(mergeMaps),
			collectParams,
		)),
	)
}

func MakeProviderFactory(
	deps []Dependency,
	fct func(param ...any) IOE.IOEither[error, any]) ProviderFactory {

	mapping := foldDeps(deps)

	identity := handleIdentity(mapping)
	optional := handleOption(mapping)
	ioeither := handleIOEither(mapping)
	iooption := handleIOOption(mapping)

	f := F.Unvariadic0(fct)

	return func(inj InjectableFactory) IOE.IOEither[error, any] {
		// resolve all dependencies
		resolved := A.MonadMap(deps, inj)
		// resolve dependencies
		return F.Pipe1(
			mergeArguments(
				identity(resolved),
				optional(resolved),
				ioeither(resolved),
				iooption(resolved),
			),
			IOE.Chain(f),
		)
	}
}
