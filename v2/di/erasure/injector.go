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
	A "github.com/IBM/fp-go/v2/array"
	"github.com/IBM/fp-go/v2/errors"
	F "github.com/IBM/fp-go/v2/function"
	I "github.com/IBM/fp-go/v2/identity"
	IOR "github.com/IBM/fp-go/v2/ioresult"
	L "github.com/IBM/fp-go/v2/lazy"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/pair"
	P "github.com/IBM/fp-go/v2/predicate"
	R "github.com/IBM/fp-go/v2/record"
	S "github.com/IBM/fp-go/v2/string"
	T "github.com/IBM/fp-go/v2/tuple"

	"slices"
	"sync"
	"sync/atomic"
)

func providerToEntry(p Provider) Entry[string, ProviderFactory] {
	return pair.MakePair(p.Provides().Id(), p.Factory())
}

func itemProviderToMap(p Provider) map[string][]ProviderFactory {
	return R.Singleton(p.Provides().Id(), A.Of(p.Factory()))
}

var (
	// missingProviderError returns a [ProviderFactory] that fails due to a missing dependency
	missingProviderError = F.Flow4(
		Dependency.String,
		errors.OnSome[string]("no provider for dependency [%s]"),
		IOR.Left[any],
		F.Constant1[InjectableFactory, IOResult[any]],
	)

	// missingProviderErrorOrDefault returns the default [ProviderFactory] or an error
	missingProviderErrorOrDefault = F.Flow3(
		T.Replicate2[Dependency],
		T.Map2(Dependency.ProviderFactory, F.Flow2(missingProviderError, F.Constant[ProviderFactory])),
		T.Tupled2(O.MonadGetOrElse[ProviderFactory]),
	)

	// circularDependencyError names the chain of dependencies that closes a circle, in resolution order
	circularDependencyError = F.Flow3(
		A.Map(Dependency.String),
		A.Intercalate(S.Monoid)(" -> "),
		errors.OnSome[string]("circular dependency detected: %s"),
	)

	emptyMulti any = A.Empty[any]()

	// emptyMultiDependency returns a [ProviderFactory] for an empty, multi dependency
	emptyMultiDependency = F.Constant1[Dependency](F.Constant1[InjectableFactory](IOR.Of(emptyMulti)))

	// handleMissingProvider covers the case of a missing provider. It either
	// returns an error or an empty multi value provider
	handleMissingProvider = F.Flow2(
		P.Fold(missingProviderErrorOrDefault, emptyMultiDependency)(isMultiDependency),
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
		pair.BiMap(collectProviders, collectItemProviders),
		pair.Paired(mergeProviders.Concat),
	)
)

// isMultiDependency tests if a dependency is a container dependency
func isMultiDependency(dep Dependency) bool {
	return dep.Flag()&MULTI == MULTI
}

// isItemProvider tests if a provivder provides a single item
func isItemProvider(provider Provider) bool {
	return provider.Provides().Flag()&ITEM == ITEM
}

// lookupFactory resolves the [ProviderFactory] for a [Dependency] against the map of
// registered providers, falling back to the default attached to the dependency or to
// an error provider via [handleMissingProvider]
func lookupFactory(factoryByID map[string]ProviderFactory) func(Dependency) ProviderFactory {
	return F.Flow3(
		T.Replicate2[Dependency],
		T.Map2(F.Flow3(
			Dependency.Id,
			R.Lookup[ProviderFactory, string],
			I.Ap[Option[ProviderFactory]](factoryByID),
		), handleMissingProvider),
		T.Tupled2(O.MonadGetOrElse[ProviderFactory]),
	)
}

// itemProviderFactory combines multiple factories into one, returning an array
func itemProviderFactory(fcts []ProviderFactory) ProviderFactory {
	return func(inj InjectableFactory) IOResult[any] {
		return F.Pipe2(
			fcts,
			IOR.TraverseArray(I.Flap[IOResult[any]](inj)),
			IOR.Map(F.ToAny[[]any]),
		)
	}
}

// MakeInjector creates an [InjectableFactory] based on a set of [Provider]s
//
// The resulting [InjectableFactory] can then be used to retrieve service instances given their [Dependency]. The implementation
// makes sure to transitively resolve the required dependencies.
//
// A dependency that refers back to itself while its provider factory runs, directly or through other dependencies,
// resolves to an error naming the chain that closes the circle instead of blocking forever. The detection follows
// one resolution at a time, so two goroutines resolving different tokens of the same cyclic graph can still block
// on each other, and so can a lookup that a provider's effect makes after its factory has returned.
func MakeInjector(providers []Provider) InjectableFactory {

	type Result = IOResult[any]
	type LazyResult = L.Lazy[Result]

	// resolved stores the values resolved so far, key is the string ID
	// of the token, value is a lazy result
	var resolved sync.Map

	// resolve the [ProviderFactory] for a [Dependency], this map is constant
	// for the lifetime of the injector
	factoryFor := lookupFactory(assembleProviders(providers))

	// the injector handed to a provider factory carries the chain of dependencies it is
	// resolving, until the factory returns and later lookups start a fresh chain
	var root InjectableFactory
	var resolveWith func([]Dependency) InjectableFactory

	resolveWith = func(path []Dependency) InjectableFactory {
		return func(token Dependency) Result {

			key := token.Id()

			if i := slices.IndexFunc(path, func(dep Dependency) bool { return dep.Id() == key }); i >= 0 {
				return IOR.Left[any](circularDependencyError(append(slices.Clip(path[i:]), token)))
			}

			// according to https://github.com/golang/go/issues/44159 this
			// is the best way to use the sync map
			actual, loaded := resolved.Load(key)
			if !loaded {
				chain := resolveWith(append(slices.Clip(path), token))
				actual, _ = resolved.LoadOrStore(key, L.Memoize(func() Result {
					var built atomic.Bool
					defer built.Store(true)
					return IOR.Memoize(factoryFor(token)(func(dep Dependency) Result {
						if built.Load() {
							return root(dep)
						}
						return chain(dep)
					}))
				}))
			}

			return actual.(LazyResult)()
		}
	}

	root = resolveWith(nil)
	return root
}
