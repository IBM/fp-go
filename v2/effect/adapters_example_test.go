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

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	R "github.com/IBM/fp-go/v2/result"
	S "github.com/IBM/fp-go/v2/string"
)

// ---------------------------------------------------------------------------
// Eitherize / Eitherize1
// ---------------------------------------------------------------------------

// exDBConfig is the dependency used by the adapter examples.
type exDBConfig struct {
	DSN string
}

// exLoadSchema is an ordinary Go function in (dependency, ctx) -> (value, error)
// form — the shape almost every Go API already has.
func exLoadSchema(cfg exDBConfig, _ context.Context) (string, error) {
	if cfg.DSN == "" {
		return "", fmt.Errorf("no DSN configured")
	}
	return "schema@" + cfg.DSN, nil
}

// exLoadTable takes one extra argument beyond the dependency and the context.
func exLoadTable(cfg exDBConfig, _ context.Context, name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("empty table name")
	}
	return cfg.DSN + "/" + name, nil
}

// ExampleEitherize demonstrates lifting an idiomatic Go (value, error) function
// into an Effect, so it composes with the rest of the package.
func ExampleEitherize() {
	eff := Eitherize(exLoadSchema)

	fmt.Println(RunSync(Provide[string](exDBConfig{DSN: "postgres://localhost"})(eff))(context.Background()))
	// Output:
	// schema@postgres://localhost <nil>
}

// ExampleEitherize_failure demonstrates that a returned error becomes a failed
// effect rather than a value the caller has to remember to check.
func ExampleEitherize_failure() {
	eff := Eitherize(exLoadSchema)

	fmt.Println(RunSync(Provide[string](exDBConfig{})(eff))(context.Background()))
	// Output:
	//  no DSN configured
}

// ExampleEitherize_composition demonstrates that the lifted function is a
// first-class Effect and can be piped point-free like any other.
func ExampleEitherize_composition() {
	eff := F.Pipe1(
		Eitherize(exLoadSchema),
		Map[exDBConfig](S.ToUpperCase),
	)

	fmt.Println(RunSync(Provide[string](exDBConfig{DSN: "postgres://localhost"})(eff))(context.Background()))
	// Output:
	// SCHEMA@POSTGRES://LOCALHOST <nil>
}

// ExampleEitherize1 demonstrates lifting a one-argument Go function into a
// Kleisli arrow, ready for Chain.
func ExampleEitherize1() {
	loadTable := Eitherize1(exLoadTable)

	eff := F.Pipe1(
		Of[exDBConfig]("users"),
		Chain(loadTable),
	)

	fmt.Println(RunSync(Provide[string](exDBConfig{DSN: "postgres://localhost"})(eff))(context.Background()))
	// Output:
	// postgres://localhost/users <nil>
}

// ExampleEitherize1_failure demonstrates that the Kleisli arrow's error short
// circuits the rest of the pipeline.
func ExampleEitherize1_failure() {
	loadTable := Eitherize1(exLoadTable)

	eff := F.Pipe2(
		Of[exDBConfig](""),
		Chain(loadTable),
		Map[exDBConfig](S.ToUpperCase), // never reached
	)

	fmt.Println(RunSync(Provide[string](exDBConfig{DSN: "postgres://localhost"})(eff))(context.Background()))
	// Output:
	//  empty table name
}

// ---------------------------------------------------------------------------
// Memoize / ContramapMemoize
// ---------------------------------------------------------------------------

// exMemoConfig must be comparable to be usable as a cache key.
type exMemoConfig struct {
	Tenant string
}

// exMemoTenant reads the cache key out of the environment.
func exMemoTenant(c exMemoConfig) string { return c.Tenant }

// ExampleMemoize demonstrates caching an effect by its environment value: the
// underlying work runs once per distinct environment, no matter how often the
// effect is executed.
func ExampleMemoize() {
	calls := 0

	// FromThunk defers the work until the effect is actually run.
	expensive := Memoize(FromThunk[exMemoConfig](func(context.Context) func() Result[int] {
		return func() Result[int] {
			calls++
			return R.Of(42)
		}
	}))

	thunk := Provide[int](exMemoConfig{Tenant: "acme"})(expensive)

	fmt.Println(RunSync(thunk)(context.Background()))
	fmt.Println(RunSync(thunk)(context.Background()))
	fmt.Println("underlying calls:", calls)
	// Output:
	// 42 <nil>
	// 42 <nil>
	// underlying calls: 1
}

// ExampleMemoize_perEnvironment demonstrates that the cache is keyed by the
// environment: a different environment triggers the work again.
func ExampleMemoize_perEnvironment() {
	calls := 0

	expensive := Memoize(FromThunk[exMemoConfig](func(context.Context) func() Result[int] {
		return func() Result[int] {
			calls++
			return R.Of(calls)
		}
	}))

	acme, _ := RunSync(Provide[int](exMemoConfig{Tenant: "acme"})(expensive))(context.Background())
	globex, _ := RunSync(Provide[int](exMemoConfig{Tenant: "globex"})(expensive))(context.Background())
	acmeAgain, _ := RunSync(Provide[int](exMemoConfig{Tenant: "acme"})(expensive))(context.Background())

	fmt.Println(acme, globex, acmeAgain)
	fmt.Println("underlying calls:", calls)
	// Output:
	// 1 2 1
	// underlying calls: 2
}

// ExampleContramapMemoize demonstrates caching by a key *derived* from the
// environment, which is what you need when the environment itself is not
// comparable or carries fields irrelevant to the result.
func ExampleContramapMemoize() {
	calls := 0

	eff := F.Pipe1(
		FromThunk[exMemoConfig](func(context.Context) func() Result[int] {
			return func() Result[int] {
				calls++
				return R.Of(calls)
			}
		}),
		// Cache on the Tenant field only.
		ContramapMemoize[int](exMemoTenant),
	)

	first, _ := RunSync(Provide[int](exMemoConfig{Tenant: "acme"})(eff))(context.Background())
	second, _ := RunSync(Provide[int](exMemoConfig{Tenant: "acme"})(eff))(context.Background())
	other, _ := RunSync(Provide[int](exMemoConfig{Tenant: "globex"})(eff))(context.Background())

	fmt.Println(first, second, other)
	fmt.Println("underlying calls:", calls)
	// Output:
	// 1 1 2
	// underlying calls: 2
}

// ---------------------------------------------------------------------------
// ProMap / Promap
// ---------------------------------------------------------------------------

// exOuter is the wider environment the caller has.
type exOuter struct {
	Inner exInner
}

// exInner is the narrower environment the effect actually needs.
type exInner struct {
	Multiplier int
}

// exOuterInner projects the wide environment onto the narrow one.
func exOuterInner(o exOuter) exInner { return o.Inner }

// exInnerMultiplier reads the field the effect depends on.
func exInnerMultiplier(i exInner) int { return i.Multiplier }

// ExampleProMap demonstrates adapting an effect on both ends at once: the
// dependency is mapped contravariantly (exOuter → exInner) while the result is
// mapped covariantly (int → string).
func ExampleProMap() {
	// An effect that only knows about exInner.
	inner := F.Pipe1(
		Asks(exInnerMultiplier),
		Map[exInner](N.Mul(3)),
	)

	// Adapt it to run in exOuter and report a string.
	adapted := ProMap(exOuterInner, S.Format[int]("value=%d"))(inner)

	fmt.Println(RunSync(Provide[string](exOuter{Inner: exInner{Multiplier: 7}})(adapted))(context.Background()))
	// Output:
	// value=21 <nil>
}

// ExampleProMap_errorsPropagate demonstrates that the result mapping is only
// applied on success.
func ExampleProMap_errorsPropagate() {
	inner := Fail[exInner, int](fmt.Errorf("inner failed"))

	adapted := ProMap(exOuterInner, S.Format[int]("value=%d"))(inner)

	fmt.Println(RunSync(Provide[string](exOuter{})(adapted))(context.Background()))
	// Output:
	//  inner failed
}

// ExamplePromap demonstrates the lower-case spelling, which is an alias of
// ProMap kept for callers who prefer the Haskell-style name.
func ExamplePromap() {
	inner := Asks(exInnerMultiplier)

	adapted := Promap(exOuterInner, S.Format[int]("m=%d"))(inner)

	fmt.Println(RunSync(Provide[string](exOuter{Inner: exInner{Multiplier: 4}})(adapted))(context.Background()))
	// Output:
	// m=4 <nil>
}
