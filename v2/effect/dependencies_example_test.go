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

package effect

import (
	"context"
	"fmt"

	ctxreader "github.com/IBM/fp-go/v2/context/reader"
	ctxthunk "github.com/IBM/fp-go/v2/context/readerioresult"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioresult"
	N "github.com/IBM/fp-go/v2/number"
	R "github.com/IBM/fp-go/v2/result"
	S "github.com/IBM/fp-go/v2/string"
)

// ---------------------------------------------------------------------------
// Shared fixtures
// ---------------------------------------------------------------------------

// exDeps is the narrow dependency an effect actually needs.
type exDeps struct {
	Multiplier int
}

// exApp is the wide dependency the application carries around.
type exApp struct {
	Deps exDeps
	Name string
}

// exAppDeps projects the wide dependency onto the narrow one.
func exAppDeps(a exApp) exDeps { return a.Deps }

// exDepsMultiplier reads the field the effect depends on.
func exDepsMultiplier(d exDeps) int { return d.Multiplier }

// exTripled is the effect under adaptation: it needs exDeps and triples its
// multiplier.  Built point-free from Asks.
func exTripled() Effect[exDeps, int] {
	return F.Pipe1(Asks(exDepsMultiplier), Map[exDeps](N.Mul(3)))
}

// ---------------------------------------------------------------------------
// Ask
// ---------------------------------------------------------------------------

// ExampleAsk demonstrates retrieving the whole dependency as a value.
func ExampleAsk() {
	eff := Ask[exDeps]()

	fmt.Println(RunSync(Provide[exDeps](exDeps{Multiplier: 7})(eff))(context.Background()))
	// Output:
	// {7} <nil>
}

// ExampleAsk_projection demonstrates that Asks(getter) is the point-free
// shorthand for Ask followed by Map.
func ExampleAsk_projection() {
	viaAsk := F.Pipe1(Ask[exDeps](), Map[exDeps](exDepsMultiplier))
	viaAsks := Asks(exDepsMultiplier)

	cfg := exDeps{Multiplier: 7}
	a, _ := RunSync(Provide[int](cfg)(viaAsk))(context.Background())
	b, _ := RunSync(Provide[int](cfg)(viaAsks))(context.Background())

	fmt.Println(a, b, a == b)
	// Output:
	// 7 7 true
}

// ---------------------------------------------------------------------------
// Local / ContraMap
// ---------------------------------------------------------------------------

// ExampleLocal demonstrates running an effect that needs a narrow dependency
// inside a wider one, by supplying a pure projection.
func ExampleLocal() {
	adapted := Local[int](exAppDeps)(exTripled())

	fmt.Println(RunSync(Provide[int](exApp{Deps: exDeps{Multiplier: 7}, Name: "app"})(adapted))(context.Background()))
	// Output:
	// 21 <nil>
}

// ExampleLocal_composition demonstrates chaining projections to descend through
// several layers of configuration.
func ExampleLocal_composition() {
	type Root struct{ App exApp }

	// Two projections compose point-free into one.
	rootToDeps := F.Flow2(func(r Root) exApp { return r.App }, exAppDeps)

	adapted := Local[int](rootToDeps)(exTripled())

	root := Root{App: exApp{Deps: exDeps{Multiplier: 5}}}
	fmt.Println(RunSync(Provide[int](root)(adapted))(context.Background()))
	// Output:
	// 15 <nil>
}

// ExampleContraMap demonstrates that ContraMap is an alias of Local, named for
// the contravariant direction the dependency travels in.
func ExampleContraMap() {
	adapted := ContraMap[int](exAppDeps)(exTripled())

	fmt.Println(RunSync(Provide[int](exApp{Deps: exDeps{Multiplier: 4}})(adapted))(context.Background()))
	// Output:
	// 12 <nil>
}

// ---------------------------------------------------------------------------
// Local*K — effectful dependency construction
// ---------------------------------------------------------------------------

// ExampleLocalReaderK demonstrates building the narrow dependency from a
// context.Context, which is what Reader means in this package.
func ExampleLocalReaderK() {
	// exApp -> Reader[exDeps]: may consult the runtime context.Context.
	loadDeps := func(a exApp) ctxreader.Reader[exDeps] {
		return func(context.Context) exDeps { return a.Deps }
	}

	adapted := LocalReaderK[int](loadDeps)(exTripled())

	fmt.Println(RunSync(Provide[int](exApp{Deps: exDeps{Multiplier: 6}})(adapted))(context.Background()))
	// Output:
	// 18 <nil>
}

// ExampleLocalIOK demonstrates building the narrow dependency with a side
// effect that cannot fail.
func ExampleLocalIOK() {
	// exApp -> IO[exDeps]
	loadDeps := func(a exApp) io.IO[exDeps] {
		return io.Of(a.Deps)
	}

	adapted := LocalIOK[int](loadDeps)(exTripled())

	fmt.Println(RunSync(Provide[int](exApp{Deps: exDeps{Multiplier: 8}})(adapted))(context.Background()))
	// Output:
	// 24 <nil>
}

// ExampleLocalResultK demonstrates building the narrow dependency with a pure
// computation that may fail — validation, for instance.
func ExampleLocalResultK() {
	// exApp -> Result[exDeps]: reject a zero multiplier.
	validate := func(a exApp) Result[exDeps] {
		if a.Deps.Multiplier == 0 {
			return R.Left[exDeps](fmt.Errorf("multiplier must not be zero"))
		}
		return R.Of(a.Deps)
	}

	adapted := LocalResultK[int](validate)(exTripled())

	fmt.Println(RunSync(Provide[int](exApp{Deps: exDeps{Multiplier: 9}})(adapted))(context.Background()))
	fmt.Println(RunSync(Provide[int](exApp{})(adapted))(context.Background()))
	// Output:
	// 27 <nil>
	// 0 multiplier must not be zero
}

// ExampleLocalIOResultK demonstrates building the narrow dependency with a side
// effect that may fail.
func ExampleLocalIOResultK() {
	// exApp -> IOResult[exDeps]
	loadDeps := func(a exApp) ioresult.IOResult[exDeps] {
		return ioresult.Of(a.Deps)
	}

	adapted := LocalIOResultK[int](loadDeps)(exTripled())

	fmt.Println(RunSync(Provide[int](exApp{Deps: exDeps{Multiplier: 2}})(adapted))(context.Background()))
	// Output:
	// 6 <nil>
}

// ExampleLocalThunkK demonstrates building the narrow dependency with a Thunk,
// which is an IO that may fail *and* sees the runtime context.Context.
func ExampleLocalThunkK() {
	loadDeps := func(a exApp) ctxthunk.ReaderIOResult[exDeps] {
		return ctxthunk.Of(a.Deps)
	}

	adapted := LocalThunkK[int](loadDeps)(exTripled())

	fmt.Println(RunSync(Provide[int](exApp{Deps: exDeps{Multiplier: 10}})(adapted))(context.Background()))
	// Output:
	// 30 <nil>
}

// ExampleLocalEffectK demonstrates the most general form: the narrow dependency
// is produced by a full Effect, so building it can itself read the wide
// dependency, perform IO and fail.
func ExampleLocalEffectK() {
	// exApp -> Effect[exApp, exDeps]: derive the multiplier from the app name.
	deriveDeps := func(a exApp) Effect[exApp, exDeps] {
		return Of[exApp](exDeps{Multiplier: len(a.Name)})
	}

	adapted := LocalEffectK[int](deriveDeps)(exTripled())

	fmt.Println(RunSync(Provide[int](exApp{Name: "abcd"})(adapted))(context.Background()))
	// Output:
	// 12 <nil>
}

// ExampleLocalEffectK_failure demonstrates that a failure while building the
// dependency fails the whole effect.
func ExampleLocalEffectK_failure() {
	deriveDeps := func(exApp) Effect[exApp, exDeps] {
		return Fail[exApp, exDeps](fmt.Errorf("cannot resolve dependencies"))
	}

	adapted := LocalEffectK[int](deriveDeps)(exTripled())

	fmt.Println(RunSync(Provide[int](exApp{})(adapted))(context.Background()))
	// Output:
	// 0 cannot resolve dependencies
}

// ---------------------------------------------------------------------------
// TraverseArray
// ---------------------------------------------------------------------------

// ExampleTraverseArray demonstrates running one effectful step per element and
// collecting the results, failing fast if any element fails.
func ExampleTraverseArray() {
	// Kleisli arrow: int -> Effect[exDeps, string], built point-free.
	label := F.Flow2(
		Of[exDeps, int],
		Map[exDeps](S.Format[int]("n=%d")),
	)

	eff := F.Pipe1(
		Succeed[exDeps]([]int{1, 2, 3}),
		Chain(TraverseArray(label)),
	)

	fmt.Println(RunSync(Provide[[]string](exDeps{})(eff))(context.Background()))
	// Output:
	// [n=1 n=2 n=3] <nil>
}

// ExampleTraverseArray_dependencyAware demonstrates that every step sees the
// same dependency.
func ExampleTraverseArray_dependencyAware() {
	// Scale each element by the multiplier read from the dependency.
	scale := func(n int) Effect[exDeps, int] {
		return F.Pipe1(Asks(exDepsMultiplier), Map[exDeps](N.Mul(n)))
	}

	eff := TraverseArray(scale)([]int{1, 2, 3})

	fmt.Println(RunSync(Provide[[]int](exDeps{Multiplier: 10})(eff))(context.Background()))
	// Output:
	// [10 20 30] <nil>
}

// ExampleTraverseArray_failFast demonstrates that the first failing element
// fails the whole traversal.
func ExampleTraverseArray_failFast() {
	positiveOnly := func(n int) Effect[exDeps, int] {
		if n < 0 {
			return Fail[exDeps, int](fmt.Errorf("negative element: %d", n))
		}
		return Of[exDeps](n)
	}

	eff := TraverseArray(positiveOnly)([]int{1, -2, 3})

	fmt.Println(RunSync(Provide[[]int](exDeps{})(eff))(context.Background()))
	// Output:
	// [] negative element: -2
}
