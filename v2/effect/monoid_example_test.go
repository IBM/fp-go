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

	A "github.com/IBM/fp-go/v2/array"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/lazy"
	N "github.com/IBM/fp-go/v2/number"
	S "github.com/IBM/fp-go/v2/string"
)

// exMonoidConfig is the dependency used by the monoid examples.
type exMonoidConfig struct{}

// ExampleApplicativeMonoid demonstrates combining effects by running both and
// merging their successful values with the underlying monoid.
func ExampleApplicativeMonoid() {
	m := ApplicativeMonoid[exMonoidConfig](S.Monoid)

	combined := m.Concat(Of[exMonoidConfig]("Hello, "), Of[exMonoidConfig]("World"))

	fmt.Println(RunSync(Provide[string](exMonoidConfig{})(combined))(context.Background()))
	// Output:
	// Hello, World <nil>
}

// ExampleApplicativeMonoid_empty demonstrates the identity element: it succeeds
// with the underlying monoid's empty value.
func ExampleApplicativeMonoid_empty() {
	sum := ApplicativeMonoid[exMonoidConfig](N.MonoidSum[int]())
	concat := ApplicativeMonoid[exMonoidConfig](S.Monoid)

	zero, _ := RunSync(Provide[int](exMonoidConfig{})(sum.Empty()))(context.Background())
	blank, _ := RunSync(Provide[string](exMonoidConfig{})(concat.Empty()))(context.Background())

	fmt.Println(zero)
	fmt.Printf("%q\n", blank)
	// Output:
	// 0
	// ""
}

// ExampleApplicativeMonoid_failure demonstrates applicative semantics: if either
// side fails, the combination fails.
func ExampleApplicativeMonoid_failure() {
	m := ApplicativeMonoid[exMonoidConfig](S.Monoid)

	combined := m.Concat(
		Fail[exMonoidConfig, string](fmt.Errorf("left failed")),
		Of[exMonoidConfig]("World"),
	)

	fmt.Println(RunSync(Provide[string](exMonoidConfig{})(combined))(context.Background()))
	// Output:
	//  left failed
}

// ExampleApplicativeMonoid_fold demonstrates folding a whole slice of effects
// into one, which is what a Monoid instance is really for.
func ExampleApplicativeMonoid_fold() {
	m := ApplicativeMonoid[exMonoidConfig](N.MonoidSum[int]())

	total := F.Pipe1(
		A.Map(Of[exMonoidConfig, int])([]int{1, 2, 3, 4}),
		A.Fold(m),
	)

	fmt.Println(RunSync(Provide[int](exMonoidConfig{})(total))(context.Background()))
	// Output:
	// 10 <nil>
}

// ExampleAlternativeMonoid demonstrates alternative semantics: a failing side is
// recovered from, so the combination succeeds as long as one side does.
func ExampleAlternativeMonoid() {
	m := AlternativeMonoid[exMonoidConfig](S.Monoid)

	// Both succeed → values are combined.
	both := m.Concat(Of[exMonoidConfig]("Hello, "), Of[exMonoidConfig]("World"))
	fmt.Println(RunSync(Provide[string](exMonoidConfig{})(both))(context.Background()))

	// The left side fails → the right side still contributes.
	recovered := m.Concat(
		Fail[exMonoidConfig, string](fmt.Errorf("left failed")),
		Of[exMonoidConfig]("World"),
	)
	fmt.Println(RunSync(Provide[string](exMonoidConfig{})(recovered))(context.Background()))
	// Output:
	// Hello, World <nil>
	// World <nil>
}

// ExampleAltMonoid demonstrates first-success-wins semantics: values are never
// combined, the first effect that succeeds is the result.
func ExampleAltMonoid() {
	// The zero element is the effect used when there is nothing to choose from.
	m := AltMonoid(lazy.Of(Fail[exMonoidConfig, string](fmt.Errorf("no candidate succeeded"))))

	// The first success wins, the second is not used.
	first := m.Concat(Of[exMonoidConfig]("primary"), Of[exMonoidConfig]("fallback"))
	fmt.Println(RunSync(Provide[string](exMonoidConfig{})(first))(context.Background()))

	// A failing head falls through to the next candidate.
	fallback := m.Concat(
		Fail[exMonoidConfig, string](fmt.Errorf("primary unavailable")),
		Of[exMonoidConfig]("fallback"),
	)
	fmt.Println(RunSync(Provide[string](exMonoidConfig{})(fallback))(context.Background()))

	// With no candidate at all the zero element surfaces.
	fmt.Println(RunSync(Provide[string](exMonoidConfig{})(m.Empty()))(context.Background()))
	// Output:
	// primary <nil>
	// fallback <nil>
	//  no candidate succeeded
}
