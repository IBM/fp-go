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
	I "github.com/IBM/fp-go/identity"
	IO "github.com/IBM/fp-go/io"
	IOG "github.com/IBM/fp-go/io/generic"
	IOE "github.com/IBM/fp-go/ioeither"
	IOO "github.com/IBM/fp-go/iooption"
	Int "github.com/IBM/fp-go/number/integer"
	O "github.com/IBM/fp-go/option"
	R "github.com/IBM/fp-go/record"
)

type InjectableFactory = func(Dependency) IOE.IOEither[error, any]
type ProviderFactory = func(InjectableFactory) IOE.IOEither[error, any]

type paramIndex = map[int]int
type paramValue = map[int]any
type handler = func(paramIndex) func([]IOE.IOEither[error, any]) IOE.IOEither[error, paramValue]

type Provider interface {
	fmt.Stringer
	// Provides returns the [Dependency] implemented by this provider
	Provides() Dependency
	// Factory returns s function that can create an instance of the dependency based on an [InjectableFactory]
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

func mapFromToken(idx int, token Dependency) map[int]paramIndex {
	return R.Singleton(token.Flag()&BehaviourMask, R.Singleton(idx, idx))
}

var (
	mergeTokenMaps = R.UnionMonoid[int](R.UnionLastSemigroup[int, int]())
	foldDeps       = A.FoldMapWithIndex[Dependency](mergeTokenMaps)(mapFromToken)
	mergeMaps      = R.UnionLastMonoid[int, any]()
	collectParams  = R.CollectOrd[any, any](Int.Ord)(F.SK[int, any])

	handlers = map[int]handler{
		Identity: func(mp paramIndex) func([]IOE.IOEither[error, any]) IOE.IOEither[error, paramValue] {
			return func(res []IOE.IOEither[error, any]) IOE.IOEither[error, paramValue] {
				return F.Pipe1(
					mp,
					IOE.TraverseRecord[int](getAt(res)),
				)
			}
		},
		Option: func(mp paramIndex) func([]IOE.IOEither[error, any]) IOE.IOEither[error, paramValue] {
			return func(res []IOE.IOEither[error, any]) IOE.IOEither[error, paramValue] {
				return F.Pipe3(
					mp,
					IOG.TraverseRecord[IO.IO[map[int]E.Either[error, any]], paramIndex](getAt(res)),
					IO.Map(R.Map[int](F.Flow2(
						E.ToOption[error, any],
						F.ToAny[O.Option[any]],
					))),
					IOE.FromIO[error, paramValue],
				)
			}
		},
		IOEither: func(mp paramIndex) func([]IOE.IOEither[error, any]) IOE.IOEither[error, paramValue] {
			return func(res []IOE.IOEither[error, any]) IOE.IOEither[error, paramValue] {
				return F.Pipe2(
					mp,
					R.Map[int](F.Flow2(
						getAt(res),
						F.ToAny[IOE.IOEither[error, any]],
					)),
					IOE.Of[error, paramValue],
				)
			}
		},
		IOOption: func(mp paramIndex) func([]IOE.IOEither[error, any]) IOE.IOEither[error, paramValue] {
			return func(res []IOE.IOEither[error, any]) IOE.IOEither[error, paramValue] {
				return F.Pipe2(
					mp,
					R.Map[int](F.Flow3(
						getAt(res),
						IOE.ToIOOption[error, any],
						F.ToAny[IOO.IOOption[any]],
					)),
					IOE.Of[error, paramValue],
				)
			}
		},
	}
)

type Mapping = map[int]paramIndex

func getAt[T any](ar []T) func(idx int) T {
	return func(idx int) T {
		return ar[idx]
	}
}

func handleMapping(mp Mapping) func(res []IOE.IOEither[error, any]) IOE.IOEither[error, []any] {
	preFct := F.Pipe1(
		mp,
		R.Collect(func(idx int, p paramIndex) func([]IOE.IOEither[error, any]) IOE.IOEither[error, paramValue] {
			return handlers[idx](p)
		}),
	)
	doFct := F.Flow2(
		I.Flap[IOE.IOEither[error, paramValue], []IOE.IOEither[error, any]],
		IOE.TraverseArray[error, func([]IOE.IOEither[error, any]) IOE.IOEither[error, paramValue], paramValue],
	)
	postFct := IOE.Map[error](F.Flow2(
		A.Fold(mergeMaps),
		collectParams,
	))

	return func(res []IOE.IOEither[error, any]) IOE.IOEither[error, []any] {
		return F.Pipe2(
			preFct,
			doFct(res),
			postFct,
		)
	}
}

// MakeProviderFactory constructs a [ProviderFactory] based on a set of [Dependency]s and
// a function that accepts the resolved dependencies to return a result
func MakeProviderFactory(
	deps []Dependency,
	fct func(param ...any) IOE.IOEither[error, any]) ProviderFactory {

	return F.Flow3(
		F.Curry2(A.MonadMap[Dependency, IOE.IOEither[error, any]])(deps),
		handleMapping(foldDeps(deps)),
		IOE.Chain(F.Unvariadic0(fct)),
	)
}
