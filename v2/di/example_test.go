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

package di_test

import (
	"fmt"
	"strings"

	A "github.com/IBM/fp-go/v2/array"
	"github.com/IBM/fp-go/v2/di"
	DIE "github.com/IBM/fp-go/v2/di/erasure"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/ioresult"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/result"
	S "github.com/IBM/fp-go/v2/string"
)

var (
	// renderString renders a [di.Result] of a string, the error message on failure
	renderString = E.Fold(error.Error, F.Identity[string])

	// renderInt renders a [di.Result] of an int, the error message on failure
	renderInt = E.Fold(error.Error, S.Format[int]("%d"))
)

// ExampleResolve demonstrates the basic wiring: tokens identify the dependencies,
// providers implement them and the injector resolves them transitively
func ExampleResolve() {
	// the tokens carry the type of the dependency
	injHost := di.MakeToken[string]("host")
	injPort := di.MakeToken[int]("port")
	injURL := di.MakeToken[string]("url")

	injector := DIE.MakeInjector(A.From(
		di.ConstProvider(injHost, "example.com"),
		di.ConstProvider(injPort, 8443),
		// the URL is assembled from the host and the port
		di.MakeProvider2(
			injURL,
			injHost.Identity(),
			injPort.Identity(),
			func(host string, port int) di.IOResult[string] {
				return ioresult.Of(fmt.Sprintf("https://%s:%d", host, port))
			},
		),
	))

	// resolving returns an effect, running it produces the value
	fmt.Println(renderString(di.Resolve(injURL)(injector)()))

	// Output:
	// https://example.com:8443
}

// ExampleResolve_missingProvider demonstrates that resolving a dependency without a
// provider fails instead of panicking
func ExampleResolve_missingProvider() {
	injDatabase := di.MakeToken[string]("database")

	fmt.Println(renderString(di.Resolve(injDatabase)(DIE.MakeInjector(DIE.Empty))()))

	// Output:
	// no provider for dependency [database]
}

// ExampleResolve_singleton demonstrates that every dependency is created at most once,
// no matter how many other dependencies refer to it
func ExampleResolve_singleton() {
	injCount := di.MakeToken[int]("counter")
	injLeft := di.MakeToken[string]("left")
	injRight := di.MakeToken[string]("right")
	injBoth := di.MakeToken[string]("both")

	var created int

	injector := DIE.MakeInjector(A.From(
		di.MakeProvider0(injCount, di.IOResult[int](func() di.Result[int] {
			created++
			return result.Of(created)
		})),
		di.MakeProvider1(injLeft, injCount.Identity(), func(c int) di.IOResult[string] {
			return ioresult.Of(fmt.Sprintf("left sees %d", c))
		}),
		di.MakeProvider1(injRight, injCount.Identity(), func(c int) di.IOResult[string] {
			return ioresult.Of(fmt.Sprintf("right sees %d", c))
		}),
		di.MakeProvider2(injBoth, injLeft.Identity(), injRight.Identity(),
			func(left, right string) di.IOResult[string] {
				return ioresult.Of(left + ", " + right)
			}),
	))

	fmt.Println(renderString(di.Resolve(injBoth)(injector)()))
	fmt.Println("created:", created)

	// Output:
	// left sees 1, right sees 1
	// created: 1
}

// ExampleInjectionToken_Option demonstrates an optional dependency, it resolves to none
// instead of failing when no provider is registered
func ExampleInjectionToken_Option() {
	injTimeout := di.MakeToken[int]("timeout")
	injClient := di.MakeToken[string]("client")

	client := di.MakeProvider1(
		injClient,
		injTimeout.Option(),
		func(timeout di.Option[int]) di.IOResult[string] {
			return ioresult.Of("timeout=" + O.Fold(F.Constant("default"), S.Format[int]("%d"))(timeout))
		},
	)

	// without a provider for the timeout
	fmt.Println(renderString(di.Resolve(injClient)(DIE.MakeInjector(A.From(client)))()))

	// with a provider for the timeout
	withTimeout := DIE.MakeInjector(A.From(client, di.ConstProvider(injTimeout, 30)))
	fmt.Println(renderString(di.Resolve(injClient)(withTimeout)()))

	// Output:
	// timeout=default
	// timeout=30
}

// ExampleInjectionToken_IOEither demonstrates a lazy, required dependency. The dependency
// is injected as an effect, so it is only created when the consumer actually runs it.
func ExampleInjectionToken_IOEither() {
	injExpensive := di.MakeToken[string]("expensive")
	injLazy := di.MakeToken[string]("lazy")

	injector := DIE.MakeInjector(A.From(
		di.MakeProvider0(injExpensive, di.IOResult[string](func() di.Result[string] {
			fmt.Println("creating the expensive dependency")
			return result.Of("expensive value")
		})),
		di.MakeProvider1(injLazy, injExpensive.IOEither(),
			func(expensive di.IOResult[string]) di.IOResult[string] {
				fmt.Println("the expensive dependency has not been created yet")
				return expensive
			}),
	))

	fmt.Println(renderString(di.Resolve(injLazy)(injector)()))

	// Output:
	// the expensive dependency has not been created yet
	// creating the expensive dependency
	// expensive value
}

// ExampleMakeTokenWithDefault0 demonstrates a token that carries its own default
// implementation, which applies unless an explicit provider overrides it
func ExampleMakeTokenWithDefault0() {
	injRetries := di.MakeTokenWithDefault0("retries", ioresult.Of(3))

	// no explicit provider, the default applies
	fmt.Println(renderInt(di.Resolve(injRetries)(DIE.MakeInjector(DIE.Empty))()))

	// an explicit provider wins over the default
	overridden := DIE.MakeInjector(A.From(di.ConstProvider(injRetries, 10)))
	fmt.Println(renderInt(di.Resolve(injRetries)(overridden)()))

	// Output:
	// 3
	// 10
}

// ExampleMakeMultiToken demonstrates a dependency with multiple implementations. Each
// implementation provides the item token, consumers depend on the container token.
func ExampleMakeMultiToken() {
	middleware := di.MakeMultiToken[string]("middleware")
	injChain := di.MakeToken[string]("chain")

	injector := DIE.MakeInjector(A.From(
		di.ConstProvider(middleware.Item(), "logging"),
		di.ConstProvider(middleware.Item(), "tracing"),
		di.MakeProvider1(injChain, middleware.Container().Identity(),
			func(items []string) di.IOResult[string] {
				return ioresult.Of(strings.Join(items, " -> "))
			}),
	))

	fmt.Println(renderString(di.Resolve(injChain)(injector)()))

	// Output:
	// logging -> tracing
}

// ExampleMakeMultiToken_empty demonstrates that a container token without any item
// provider resolves to the empty array rather than failing
func ExampleMakeMultiToken_empty() {
	middleware := di.MakeMultiToken[string]("middleware")
	injChain := di.MakeToken[string]("chain")

	injector := DIE.MakeInjector(A.From(
		di.MakeProvider1(injChain, middleware.Container().Identity(),
			func(items []string) di.IOResult[string] {
				return ioresult.Of(fmt.Sprintf("%d middlewares", len(items)))
			}),
	))

	fmt.Println(renderString(di.Resolve(injChain)(injector)()))

	// Output:
	// 0 middlewares
}

// ExampleRunMain demonstrates how to run an application assembled from providers. The
// entry point provides [di.InjMain] and the result is the error of the application.
func ExampleRunMain() {
	injGreeting := di.MakeToken[string]("greeting")

	err := di.RunMain(A.From(
		di.ConstProvider(injGreeting, "Hello, World!"),
		di.MakeProvider1(di.InjMain, injGreeting.Identity(),
			func(greeting string) di.IOResult[any] {
				return func() di.Result[any] {
					fmt.Println(greeting)
					// the main token is typed as [any], it must not resolve to nil
					return result.Of[any](greeting)
				}
			}),
	))()

	fmt.Println("error:", err)

	// Output:
	// Hello, World!
	// error: <nil>
}
