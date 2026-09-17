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
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOR "github.com/IBM/fp-go/v2/ioresult"
	O "github.com/IBM/fp-go/v2/option"
	S "github.com/IBM/fp-go/v2/string"
)

// dependency is a minimal [Dependency]: an identity, a display name and a behaviour flag.
// The strongly typed [github.com/IBM/fp-go/v2/di] package generates these for you, the
// erasure layer only cares about the identity and the flag.
type dependency struct {
	name string
	flag int
}

func (d dependency) String() string                           { return d.name }
func (d dependency) Id() string                               { return d.name }
func (d dependency) Flag() int                                { return d.flag }
func (d dependency) ProviderFactory() Option[ProviderFactory] { return O.None[ProviderFactory]() }

var (
	// render turns an untyped [Result] into a string, the error message on failure
	render = E.Fold(error.Error, S.Format[any]("%v"))

	// renderOption renders an untyped [Option]
	renderOption = O.Fold(F.Constant("<none>"), S.Format[any]("%v"))
)

// ExampleMakeInjector demonstrates how an injector resolves a dependency by transitively
// resolving everything the corresponding provider needs
func ExampleMakeInjector() {
	// the dependencies of our little application
	greeting := dependency{"greeting", IDENTITY}
	subject := dependency{"subject", IDENTITY}
	message := dependency{"message", IDENTITY}

	injector := MakeInjector(A.From(
		MakeProvider(greeting, F.Constant1[InjectableFactory](IOR.Of[any]("Hello"))),
		MakeProvider(subject, F.Constant1[InjectableFactory](IOR.Of[any]("World"))),
		// the message is assembled from the two dependencies above
		MakeProvider(message, MakeProviderFactory(
			A.From[Dependency](greeting, subject),
			func(params ...any) IOResult[any] {
				return IOR.Of[any](fmt.Sprintf("%s, %s!", params[0], params[1]))
			},
		)),
	))

	// resolving returns an effect, running it produces the value
	fmt.Println(render(injector(message)()))

	// Output:
	// Hello, World!
}

// ExampleMakeInjector_missingProvider demonstrates that a dependency without a provider
// makes the resolution fail
func ExampleMakeInjector_missingProvider() {
	injector := MakeInjector(Empty)

	fmt.Println(render(injector(dependency{"database", IDENTITY})()))

	// Output:
	// no provider for dependency [database]
}

// ExampleMakeProviderFactory demonstrates the four behaviours a [Dependency] can request.
// All four refer to the same dependency, they only differ in how it gets injected.
func ExampleMakeProviderFactory() {
	// all four views share the same identity, the flag selects the behaviour
	required := dependency{"port", IDENTITY}
	optional := dependency{"port", OPTION}
	lazy := dependency{"port", IOEITHER}
	lazyOptional := dependency{"port", IOOPTION}

	factory := MakeProviderFactory(
		A.From[Dependency](required, optional, lazy, lazyOptional),
		func(params ...any) IOResult[any] {
			return IOR.Of[any](fmt.Sprintf(
				"required=%v, optional=%s, lazy=%s, lazyOptional=%s",
				params[0],
				renderOption(params[1].(Option[any])),
				render(params[2].(IOResult[any])()),
				renderOption(params[3].(IOOption[any])()),
			))
		},
	)

	injector := MakeInjector(A.From(
		MakeProvider(required, F.Constant1[InjectableFactory](IOR.Of[any](8080))),
	))

	fmt.Println(render(factory(injector)()))

	// Output:
	// required=8080, optional=8080, lazy=8080, lazyOptional=8080
}

// ExampleMakeProviderFactory_optional demonstrates that an optional dependency without a
// provider resolves to none instead of failing
func ExampleMakeProviderFactory_optional() {
	factory := MakeProviderFactory(
		A.From[Dependency](dependency{"port", OPTION}),
		func(params ...any) IOResult[any] {
			return IOR.Of[any]("optional=" + renderOption(params[0].(Option[any])))
		},
	)

	fmt.Println(render(factory(MakeInjector(Empty))()))

	// Output:
	// optional=<none>
}

// ExampleMakeInjector_multi demonstrates how the providers for an item dependency are
// aggregated into the array exposed by the matching container dependency
func ExampleMakeInjector_multi() {
	// item and container share the same identity
	item := dependency{"middleware", ITEM | IDENTITY}
	container := dependency{"middleware", MULTI | IDENTITY}

	injector := MakeInjector(A.From(
		MakeProvider(item, F.Constant1[InjectableFactory](IOR.Of[any]("logging"))),
		MakeProvider(item, F.Constant1[InjectableFactory](IOR.Of[any]("tracing"))),
	))

	fmt.Println(render(injector(container)()))

	// Output:
	// [logging tracing]
}
