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
	"log"

	A "github.com/IBM/fp-go/array"
	"github.com/IBM/fp-go/errors"
	F "github.com/IBM/fp-go/function"
	I "github.com/IBM/fp-go/identity"
	IG "github.com/IBM/fp-go/identity/generic"
	IOE "github.com/IBM/fp-go/ioeither"
	L "github.com/IBM/fp-go/lazy"
	O "github.com/IBM/fp-go/option"
	R "github.com/IBM/fp-go/record"
	T "github.com/IBM/fp-go/tuple"

	"sync"
)

func providerToEntry(p Provider) T.Tuple2[string, ProviderFactory] {
	return T.MakeTuple2(p.Provides().Id(), p.Factory())
}

func itemProviderToMap(p Provider) map[string][]ProviderFactory {
	return R.Singleton(p.Provides().Id(), A.Of(p.Factory()))
}

var missingProviderError = F.Flow4(
	Dependency.String,
	errors.OnSome[string]("no provider for dependency [%s]"),
	IOE.Left[any, error],
	F.Constant1[InjectableFactory, IOE.IOEither[error, any]],
)

var emptyMulti any = A.Empty[any]()

var emptyMultiDependency = F.Constant1[Dependency](F.Constant1[InjectableFactory](IOE.Of[error](emptyMulti)))

func logEntryExit(name string, token Dependency) func() {
	log.Printf("Entry: [%s] -> [%s]:[%s]", name, token.Id(), token.String())
	return func() {
		log.Printf("Exit:  [%s] -> [%s]:[%s]", name, token.Id(), token.String())
	}
}

// isMultiDependency tests if a dependency is a container dependency
func isMultiDependency(dep Dependency) bool {
	return dep.Type() == Multi
}

var handleMissingProvider = F.Flow2(
	F.Ternary(isMultiDependency, emptyMultiDependency, missingProviderError),
	F.Constant[ProviderFactory],
)

// isItemProvider tests if a provivder provides a single item
func isItemProvider(provider Provider) bool {
	return provider.Provides().Type() == Item
}

// itemProviderFactory combines multiple factories into one, returning an array
func itemProviderFactory(fcts []ProviderFactory) ProviderFactory {
	return func(inj InjectableFactory) IOE.IOEither[error, any] {
		return F.Pipe2(
			fcts,
			IOE.TraverseArray(I.Flap[IOE.IOEither[error, any]](inj)),
			IOE.Map[error](F.ToAny[[]any]),
		)
	}
}

var mergeItemProviders = R.UnionMonoid[string](A.Semigroup[ProviderFactory]())

// collectItemProviders create a provider map for item providers
var collectItemProviders = F.Flow2(
	A.FoldMap[Provider](mergeItemProviders)(itemProviderToMap),
	R.Map[string](itemProviderFactory),
)

// collectProviders collects non-item providers
var collectProviders = F.Flow2(
	A.Map(providerToEntry),
	R.FromEntries[string, ProviderFactory],
)

var mergeProviders = R.UnionLastMonoid[string, ProviderFactory]()

// assembleProviders constructs the provider map for item and non-item providers
var assembleProviders = F.Flow3(
	A.Partition(isItemProvider),
	T.Map2(collectProviders, collectItemProviders),
	T.Tupled2(mergeProviders.Concat),
)

func MakeInjector(providers []Provider) InjectableFactory {

	type Result = IOE.IOEither[error, any]
	type LazyResult = L.Lazy[Result]

	// resolved stores the values resolved so far, key is the string ID
	// of the token, value is a lazy result
	var resolved sync.Map

	// provide a mapping for all providers
	factoryById := assembleProviders(providers)

	// the actual factory, we need lazy initialization
	var injFct InjectableFactory

	// lazy initialization, so we can cross reference it
	injFct = func(token Dependency) Result {

		defer logEntryExit("inj", token)()

		key := token.Id()

		// according to https://github.com/golang/go/issues/44159 this
		// is the best way to use the sync map
		actual, loaded := resolved.Load(key)
		if !loaded {

			computeResult := func() Result {
				defer logEntryExit("computeResult", token)()
				return F.Pipe5(
					token,
					T.Replicate2[Dependency],
					T.Map2(F.Flow3(
						Dependency.Id,
						R.Lookup[ProviderFactory, string],
						I.Ap[O.Option[ProviderFactory]](factoryById),
					), handleMissingProvider),
					T.Tupled2(O.MonadGetOrElse[ProviderFactory]),
					IG.Ap[ProviderFactory](injFct),
					IOE.Memoize[error, any],
				)
			}

			actual, _ = resolved.LoadOrStore(key, F.Pipe1(
				computeResult,
				L.Memoize[Result],
			))
		}

		return actual.(LazyResult)()
	}

	return injFct
}
