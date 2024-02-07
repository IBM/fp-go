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

var (
	// missingProviderError returns a [ProviderFactory] that fails due to a missing dependency
	missingProviderError = F.Flow4(
		Dependency.String,
		errors.OnSome[string]("no provider for dependency [%s]"),
		IOE.Left[any, error],
		F.Constant1[InjectableFactory, IOE.IOEither[error, any]],
	)

	// missingProviderErrorOrDefault returns the default [ProviderFactory] or an error
	missingProviderErrorOrDefault = F.Flow3(
		T.Replicate2[Dependency],
		T.Map2(Dependency.ProviderFactory, F.Flow2(missingProviderError, F.Constant[ProviderFactory])),
		T.Tupled2(O.MonadGetOrElse[ProviderFactory]),
	)

	emptyMulti any = A.Empty[any]()

	// emptyMultiDependency returns a [ProviderFactory] for an empty, multi dependency
	emptyMultiDependency = F.Constant1[Dependency](F.Constant1[InjectableFactory](IOE.Of[error](emptyMulti)))

	// handleMissingProvider covers the case of a missing provider. It either
	// returns an error or an empty multi value provider
	handleMissingProvider = F.Flow2(
		F.Ternary(isMultiDependency, emptyMultiDependency, missingProviderErrorOrDefault),
		F.Constant[ProviderFactory],
	)

	// mergeItemProviders is a monoid for item provider factories
	mergeItemProviders = R.UnionMonoid[string](A.Semigroup[ProviderFactory]())

	// mergeProviders is a monoid for provider factories
	mergeProviders = R.UnionLastMonoid[string, ProviderFactory]()

	// collectItemProviders create a provider map for item providers
	collectItemProviders = F.Flow2(
		A.FoldMap[Provider](mergeItemProviders)(itemProviderToMap),
		R.Map[string](itemProviderFactory),
	)

	// collectProviders collects non-item providers
	collectProviders = F.Flow2(
		A.Map(providerToEntry),
		R.FromEntries[string, ProviderFactory],
	)

	// assembleProviders constructs the provider map for item and non-item providers
	assembleProviders = F.Flow3(
		A.Partition(isItemProvider),
		T.Map2(collectProviders, collectItemProviders),
		T.Tupled2(mergeProviders.Concat),
	)
)

// isMultiDependency tests if a dependency is a container dependency
func isMultiDependency(dep Dependency) bool {
	return dep.Flag()&Multi == Multi
}

// isItemProvider tests if a provivder provides a single item
func isItemProvider(provider Provider) bool {
	return provider.Provides().Flag()&Item == Item
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

// MakeInjector creates an [InjectableFactory] based on a set of [Provider]s
//
// The resulting [InjectableFactory] can then be used to retrieve service instances given their [Dependency]. The implementation
// makes sure to transitively resolve the required dependencies.
func MakeInjector(providers []Provider) InjectableFactory {

	type Result = IOE.IOEither[error, any]
	type LazyResult = L.Lazy[Result]

	// resolved stores the values resolved so far, key is the string ID
	// of the token, value is a lazy result
	var resolved sync.Map

	// provide a mapping for all providers
	factoryByID := assembleProviders(providers)

	// the actual factory, we need lazy initialization
	var injFct InjectableFactory

	// lazy initialization, so we can cross reference it
	injFct = func(token Dependency) Result {

		key := token.Id()

		// according to https://github.com/golang/go/issues/44159 this
		// is the best way to use the sync map
		actual, loaded := resolved.Load(key)
		if !loaded {

			computeResult := L.MakeLazy(func() Result {
				return F.Pipe5(
					token,
					T.Replicate2[Dependency],
					T.Map2(F.Flow3(
						Dependency.Id,
						R.Lookup[ProviderFactory, string],
						I.Ap[O.Option[ProviderFactory]](factoryByID),
					), handleMissingProvider),
					T.Tupled2(O.MonadGetOrElse[ProviderFactory]),
					IG.Ap[ProviderFactory](injFct),
					IOE.Memoize[error, any],
				)
			})

			actual, _ = resolved.LoadOrStore(key, F.Pipe1(
				computeResult,
				L.Memoize[Result],
			))
		}

		return actual.(LazyResult)()
	}

	return injFct
}
