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

package erasure

import (
	"fmt"

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	I "github.com/IBM/fp-go/v2/identity"
	IO "github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	Int "github.com/IBM/fp-go/v2/number/integer"
	R "github.com/IBM/fp-go/v2/record"
	"github.com/IBM/fp-go/v2/result"
)

type (
	// InjectableFactory is a factory function that can create an untyped instance of a service based on its [Dependency] identifier
	InjectableFactory = ReaderIOResult[Dependency, any]
	ProviderFactory   = ReaderIOResult[InjectableFactory, any]

	paramIndex = map[int]int
	paramValue = map[int]any
	handler    = func(paramIndex) func([]IOResult[any]) IOResult[paramValue]
	mapping    = map[int]paramIndex

	Provider interface {
		fmt.Stringer
		// Provides returns the [Dependency] implemented by this provider
		Provides() Dependency
		// Factory returns s function that can create an instance of the dependency based on an [InjectableFactory]
		Factory() ProviderFactory
	}

	provider struct {
		provides Dependency
		factory  ProviderFactory
	}
)

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
	// Empty is the empty array of providers
	Empty = A.Empty[Provider]()

	mergeTokenMaps = R.UnionMonoid[int](R.UnionLastSemigroup[int, int]())
	foldDeps       = A.FoldMapWithIndex[Dependency](mergeTokenMaps)(mapFromToken)
	mergeMaps      = R.UnionLastMonoid[int, any]()
	collectParams  = R.CollectOrd[any, any](Int.Ord)(F.SK[int, any])

	mapDeps = F.Curry2(A.MonadMap[Dependency, IOResult[any]])

	handlers = map[int]handler{
		IDENTITY: func(mp paramIndex) func([]IOResult[any]) IOResult[paramValue] {
			return func(res []IOResult[any]) IOResult[paramValue] {
				return F.Pipe1(
					mp,
					IOE.TraverseRecord[int](getAt(res)),
				)
			}
		},
		OPTION: func(mp paramIndex) func([]IOResult[any]) IOResult[paramValue] {
			return func(res []IOResult[any]) IOResult[paramValue] {
				return F.Pipe3(
					mp,
					IO.TraverseRecord[int](getAt(res)),
					IO.Map(R.Map[int](F.Flow2(
						result.ToOption[any],
						F.ToAny[Option[any]],
					))),
					IOE.FromIO[error, paramValue],
				)
			}
		},
		IOEITHER: func(mp paramIndex) func([]IOResult[any]) IOResult[paramValue] {
			return func(res []IOResult[any]) IOResult[paramValue] {
				return F.Pipe2(
					mp,
					R.Map[int](F.Flow2(
						getAt(res),
						F.ToAny[IOResult[any]],
					)),
					IOE.Of[error, paramValue],
				)
			}
		},
		IOOPTION: func(mp paramIndex) func([]IOResult[any]) IOResult[paramValue] {
			return func(res []IOResult[any]) IOResult[paramValue] {
				return F.Pipe2(
					mp,
					R.Map[int](F.Flow3(
						getAt(res),
						IOE.ToIOOption[error, any],
						F.ToAny[IOOption[any]],
					)),
					IOE.Of[error, paramValue],
				)
			}
		},
	}
)

func getAt[T any](ar []T) func(idx int) T {
	return func(idx int) T {
		return ar[idx]
	}
}

func handleMapping(mp mapping) func(res []IOResult[any]) IOResult[[]any] {
	preFct := F.Pipe1(
		mp,
		R.Collect(func(idx int, p paramIndex) func([]IOResult[any]) IOResult[paramValue] {
			return handlers[idx](p)
		}),
	)
	doFct := F.Flow2(
		I.Flap[IOResult[paramValue], []IOResult[any]],
		IOE.TraverseArray[error, func([]IOResult[any]) IOResult[paramValue], paramValue],
	)
	postFct := IOE.Map[error](F.Flow2(
		A.Fold(mergeMaps),
		collectParams,
	))

	return func(res []IOResult[any]) IOResult[[]any] {
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
	fct func(param ...any) IOResult[any]) ProviderFactory {

	return F.Flow3(
		mapDeps(deps),
		handleMapping(foldDeps(deps)),
		IOE.Chain(F.Unvariadic0(fct)),
	)
}
