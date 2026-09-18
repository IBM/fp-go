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
	"errors"
	"sync"
	"sync/atomic"
	"testing"

	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOR "github.com/IBM/fp-go/v2/ioresult"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// testDependency is a minimal [Dependency] implementation used to exercise the
// erasure layer without going through the strongly typed [di] package
type testDependency struct {
	name string
	id   string
	flag int
	def  Option[ProviderFactory]
}

func (d *testDependency) String() string                           { return d.name }
func (d *testDependency) Id() string                               { return d.id }
func (d *testDependency) Flag() int                                { return d.flag }
func (d *testDependency) ProviderFactory() Option[ProviderFactory] { return d.def }

// makeDep creates a [Dependency] without an attached default [ProviderFactory]
func makeDep(name, id string, flag int) *testDependency {
	return &testDependency{name, id, flag, O.None[ProviderFactory]()}
}

// constFactory creates a [ProviderFactory] that ignores the injector and returns a constant
func constFactory(value any) ProviderFactory {
	return F.Constant1[InjectableFactory](IOR.Of(value))
}

// countingFactory creates a [ProviderFactory] whose effect increments a counter when run
func countingFactory(count *int32, value any) ProviderFactory {
	return F.Constant1[InjectableFactory](IOResult[any](func() result.Result[any] {
		atomic.AddInt32(count, 1)
		return result.Of(value)
	}))
}

func TestGetAt(t *testing.T) {
	at := getAt(A.From("a", "b", "c"))

	assert.Equal(t, "a", at(0))
	assert.Equal(t, "c", at(2))
}

func TestIsMultiDependency(t *testing.T) {
	assert.True(t, isMultiDependency(makeDep("multi", "id", MULTI|IDENTITY)))
	assert.False(t, isMultiDependency(makeDep("single", "id", IDENTITY)))
	assert.False(t, isMultiDependency(makeDep("item", "id", ITEM|IDENTITY)))
}

func TestIsItemProvider(t *testing.T) {
	assert.True(t, isItemProvider(MakeProvider(makeDep("item", "id", ITEM|IDENTITY), constFactory("a"))))
	assert.False(t, isItemProvider(MakeProvider(makeDep("single", "id", IDENTITY), constFactory("a"))))
}

func TestMapFromToken(t *testing.T) {
	assert.Equal(t, mapping{IDENTITY: paramIndex{2: 2}}, mapFromToken(2, makeDep("dep", "id", IDENTITY)))
	assert.Equal(t, mapping{IOOPTION: paramIndex{0: 0}}, mapFromToken(0, makeDep("dep", "id", MULTI|IOOPTION)))
}

func TestProvider(t *testing.T) {
	dep := makeDep("Service", "id", IDENTITY)
	fct := constFactory("value")

	p := MakeProvider(dep, fct)

	assert.Equal(t, dep, p.Provides())
	assert.Equal(t, result.Of[any]("value"), p.Factory()(MakeInjector(Empty))())
	assert.Equal(t, "Provider for [Service]", p.String())
}

func TestHandlerForFlag(t *testing.T) {
	for _, flag := range A.From(IDENTITY, OPTION, IOEITHER, IOOPTION) {
		assert.NotNil(t, handlerForFlag(flag), "no handler for flag %d", flag)
	}
}

// TestHandleMappingOrdersParameters verifies that the parameters are collected in the
// order of their positions, independently of the behaviour they were grouped under
func TestHandleMappingOrdersParameters(t *testing.T) {
	res := A.From(IOR.Of[any]("zero"), IOR.Of[any]("one"), IOR.Of[any]("two"))

	// position 0 and 2 are required, position 1 is optional
	mp := mapping{
		IDENTITY: paramIndex{0: 0, 2: 2},
		OPTION:   paramIndex{1: 1},
	}

	assert.Equal(t, result.Of(A.From[any]("zero", O.Of[any]("one"), "two")), handleMapping(mp)(res)())
}

// TestHandleMappingRequiredFailure verifies that a failing required parameter fails the mapping
func TestHandleMappingRequiredFailure(t *testing.T) {
	res := A.From(IOR.Left[any](errors.New("boom")))

	assert.True(t, E.IsLeft(handleMapping(mapping{IDENTITY: paramIndex{0: 0}})(res)()))
}

// TestHandleMappingOptionalFailure verifies that a failing optional parameter becomes none
func TestHandleMappingOptionalFailure(t *testing.T) {
	res := A.From(IOR.Left[any](errors.New("boom")))

	assert.Equal(t, result.Of(A.From[any](O.None[any]())), handleMapping(mapping{OPTION: paramIndex{0: 0}})(res)())
}

func TestHandleMappingEmpty(t *testing.T) {
	assert.Equal(t, result.Of(A.Empty[any]()), handleMapping(mapping{})(A.Empty[IOResult[any]]())())
}

// TestLookupFactory covers the three ways a [ProviderFactory] can be obtained for a [Dependency]
func TestLookupFactory(t *testing.T) {
	inj := MakeInjector(Empty)
	lookup := lookupFactory(map[string]ProviderFactory{"registered": constFactory("registered value")})

	t.Run("registered provider", func(t *testing.T) {
		assert.Equal(t, result.Of[any]("registered value"), lookup(makeDep("Registered", "registered", IDENTITY))(inj)())
	})

	t.Run("default on the dependency", func(t *testing.T) {
		dep := &testDependency{"WithDefault", "withDefault", IDENTITY, O.Of(constFactory("default value"))}
		assert.Equal(t, result.Of[any]("default value"), lookup(dep)(inj)())
	})

	t.Run("missing provider fails", func(t *testing.T) {
		res := lookup(makeDep("Missing", "missing", IDENTITY))(inj)()

		assert.True(t, E.IsLeft(res))
		assert.ErrorContains(t, E.Fold(F.Identity[error], F.Constant1[any, error](nil))(res), "no provider for dependency [Missing]")
	})

	t.Run("missing multi provider is empty", func(t *testing.T) {
		assert.Equal(t, result.Of(emptyMulti), lookup(makeDep("MissingMulti", "missingMulti", MULTI|IDENTITY))(inj)())
	})
}

func TestMakeInjectorTransitive(t *testing.T) {
	depA := makeDep("A", "a", IDENTITY)
	depB := makeDep("B", "b", IDENTITY)

	inj := MakeInjector(A.From(
		MakeProvider(depA, constFactory("hello")),
		MakeProvider(depB, MakeProviderFactory(
			A.From[Dependency](depA),
			func(params ...any) IOResult[any] {
				return IOR.Of[any](params[0].(string) + " world")
			},
		)),
	))

	assert.Equal(t, result.Of[any]("hello world"), inj(depB)())
}

// TestMakeInjectorMemoizes verifies that the effect of a provider runs at most once,
// no matter how often the dependency is resolved
func TestMakeInjectorMemoizes(t *testing.T) {
	var count int32

	dep := makeDep("Counted", "counted", IDENTITY)
	inj := MakeInjector(A.From(MakeProvider(dep, countingFactory(&count, "value"))))

	first := inj(dep)
	second := inj(dep)

	// resolution alone must not run the effect
	assert.Equal(t, int32(0), atomic.LoadInt32(&count))

	assert.Equal(t, result.Of[any]("value"), first())
	assert.Equal(t, result.Of[any]("value"), second())
	assert.Equal(t, result.Of[any]("value"), inj(dep)())

	assert.Equal(t, int32(1), atomic.LoadInt32(&count))
}

// TestMakeInjectorConcurrentResolution verifies that the memoization of the injector is
// safe to use from multiple goroutines
func TestMakeInjectorConcurrentResolution(t *testing.T) {
	var count int32

	dep := makeDep("Concurrent", "concurrent", IDENTITY)
	inj := MakeInjector(A.From(MakeProvider(dep, countingFactory(&count, "value"))))

	var wg sync.WaitGroup
	for range 64 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			assert.Equal(t, result.Of[any]("value"), inj(dep)())
		}()
	}
	wg.Wait()

	assert.Equal(t, int32(1), atomic.LoadInt32(&count))
}

func TestMakeInjectorMissingProvider(t *testing.T) {
	inj := MakeInjector(Empty)

	assert.True(t, E.IsLeft(inj(makeDep("Missing", "missing", IDENTITY))()))
}

func TestMakeInjectorMissingMultiProvider(t *testing.T) {
	inj := MakeInjector(Empty)

	assert.Equal(t, result.Of(emptyMulti), inj(makeDep("Empty", "empty", MULTI|IDENTITY))())
}

func TestMakeInjectorDefaultProviderFactory(t *testing.T) {
	dep := &testDependency{"WithDefault", "withDefault", IDENTITY, O.Of(constFactory("default value"))}

	assert.Equal(t, result.Of[any]("default value"), MakeInjector(Empty)(dep)())
}

// TestMakeInjectorDefaultProviderFactoryOverride verifies that a registered provider wins
// over the default attached to the dependency
func TestMakeInjectorDefaultProviderFactoryOverride(t *testing.T) {
	dep := &testDependency{"WithDefault", "withDefault", IDENTITY, O.Of(constFactory("default value"))}

	inj := MakeInjector(A.From(MakeProvider(dep, constFactory("explicit value"))))

	assert.Equal(t, result.Of[any]("explicit value"), inj(dep)())
}

// TestMakeInjectorItemProviders verifies that all providers for an item token are
// aggregated into the array exposed by the corresponding container token
func TestMakeInjectorItemProviders(t *testing.T) {
	item := makeDep("Item", "shared", ITEM|IDENTITY)
	container := makeDep("Container", "shared", MULTI|IDENTITY)

	inj := MakeInjector(A.From(
		MakeProvider(item, constFactory("a")),
		MakeProvider(item, constFactory("b")),
	))

	assert.Equal(t, result.Of[any](A.From[any]("a", "b")), inj(container)())
}

// TestMakeInjectorItemProviderFailure verifies that a failing item fails the container
func TestMakeInjectorItemProviderFailure(t *testing.T) {
	item := makeDep("Item", "shared", ITEM|IDENTITY)
	container := makeDep("Container", "shared", MULTI|IDENTITY)

	inj := MakeInjector(A.From(
		MakeProvider(item, constFactory("a")),
		MakeProvider(item, F.Constant1[InjectableFactory](IOR.Left[any](errors.New("boom")))),
	))

	assert.True(t, E.IsLeft(inj(container)()))
}

// TestMakeProviderFactoryBehaviours verifies that each behaviour flag injects the
// parameter in its expected shape and at its expected position
func TestMakeProviderFactoryBehaviours(t *testing.T) {
	identity := makeDep("V", "v", IDENTITY)
	optional := makeDep("Option[V]", "v", OPTION)
	lazy := makeDep("IOEither[V]", "v", IOEITHER)
	lazyOptional := makeDep("IOOption[V]", "v", IOOPTION)

	var params []any

	fct := MakeProviderFactory(
		A.From[Dependency](identity, optional, lazy, lazyOptional),
		func(p ...any) IOResult[any] {
			params = p
			return IOR.Of[any](len(p))
		},
	)

	inj := MakeInjector(A.From(MakeProvider(identity, constFactory("value"))))

	assert.Equal(t, result.Of[any](4), fct(inj)())
	assert.Len(t, params, 4)
	assert.Equal(t, "value", params[0])
	assert.Equal(t, O.Of[any]("value"), params[1])
	assert.Equal(t, result.Of[any]("value"), params[2].(IOResult[any])())
	assert.Equal(t, O.Of[any]("value"), params[3].(IOOption[any])())
}

// TestMakeProviderFactoryMissingBehaviours verifies the behaviour of each flag when the
// dependency cannot be resolved
func TestMakeProviderFactoryMissingBehaviours(t *testing.T) {
	inj := MakeInjector(Empty)

	t.Run("required fails", func(t *testing.T) {
		fct := MakeProviderFactory(
			A.From[Dependency](makeDep("Missing", "missing", IDENTITY)),
			func(p ...any) IOResult[any] { return IOR.Of[any](p[0]) },
		)

		assert.True(t, E.IsLeft(fct(inj)()))
	})

	t.Run("optional becomes none", func(t *testing.T) {
		var params []any
		fct := MakeProviderFactory(
			A.From[Dependency](makeDep("Option[Missing]", "missing", OPTION)),
			func(p ...any) IOResult[any] {
				params = p
				return IOR.Of[any]("ok")
			},
		)

		assert.Equal(t, result.Of[any]("ok"), fct(inj)())
		assert.Equal(t, O.None[any](), params[0])
	})

	t.Run("lazy optional becomes none", func(t *testing.T) {
		var params []any
		fct := MakeProviderFactory(
			A.From[Dependency](makeDep("IOOption[Missing]", "missing", IOOPTION)),
			func(p ...any) IOResult[any] {
				params = p
				return IOR.Of[any]("ok")
			},
		)

		assert.Equal(t, result.Of[any]("ok"), fct(inj)())
		assert.Equal(t, O.None[any](), params[0].(IOOption[any])())
	})

	t.Run("lazy required defers the failure", func(t *testing.T) {
		var params []any
		fct := MakeProviderFactory(
			A.From[Dependency](makeDep("IOEither[Missing]", "missing", IOEITHER)),
			func(p ...any) IOResult[any] {
				params = p
				return IOR.Of[any]("ok")
			},
		)

		// the factory itself succeeds, the failure only surfaces when the parameter is run
		assert.Equal(t, result.Of[any]("ok"), fct(inj)())
		assert.True(t, E.IsLeft(params[0].(IOResult[any])()))
	})
}

// TestMakeProviderFactoryWithoutDependencies verifies the degenerate case of a factory
// that does not have any dependency
func TestMakeProviderFactoryWithoutDependencies(t *testing.T) {
	fct := MakeProviderFactory(
		A.Empty[Dependency](),
		func(p ...any) IOResult[any] { return IOR.Of[any](len(p)) },
	)

	assert.Equal(t, result.Of[any](0), fct(MakeInjector(Empty))())
}
