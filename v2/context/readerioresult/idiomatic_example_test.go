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

package readerioresult

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/result"
)

// ---------------------------------------------------------------------------
// Shared fixtures
//
// Every function below has the shape ordinary Go code already has: it takes a
// context.Context where it needs one and returns (value, error). None of them
// knows about Result, IO or Reader - that is the point of the idiomatic bridge.
// ---------------------------------------------------------------------------

// exTenantKey is the context key the examples read their tenant from.
type exTenantKey struct{}

// exCtx builds a context carrying a tenant name.
func exCtx(tenant string) context.Context {
	return context.WithValue(context.Background(), exTenantKey{}, tenant)
}

// exRun executes a ReaderIOResult and returns the idiomatic (value, error) pair,
// so the examples can print it the way plain Go code would.
func exRun[A any](ctx context.Context, rio ReaderIOResult[A]) (A, error) {
	return result.Unwrap(rio(ctx)())
}

// exParse is an idiomatic Result Kleisli arrow: no context, no IO, may fail.
func exParse(s string) (int, error) {
	return strconv.Atoi(s)
}

// exTenant is an idiomatic ReaderResult Kleisli arrow: it reads the context and
// may fail, but performs no IO.
func exTenant(name string) func(context.Context) (string, error) {
	return func(ctx context.Context) (string, error) {
		tenant, ok := ctx.Value(exTenantKey{}).(string)
		if !ok {
			return "", errors.New("no tenant in context")
		}
		return tenant + "/" + name, nil
	}
}

// exLoad is an idiomatic IOResult Kleisli arrow: it performs a (here simulated)
// side effect and may fail, but does not read the context.
func exLoad(n int) func() (string, error) {
	return func() (string, error) {
		if n < 0 {
			return "", fmt.Errorf("cannot load %d", n)
		}
		return fmt.Sprintf("record-%d", n), nil
	}
}

// exFetch is an idiomatic ReaderIOResult Kleisli arrow - the full stack: it
// reads the context, performs a side effect and may fail.
func exFetch(path string) func(context.Context) func() (string, error) {
	return func(ctx context.Context) func() (string, error) {
		return func() (string, error) {
			tenant, ok := ctx.Value(exTenantKey{}).(string)
			if !ok {
				return "", errors.New("no tenant in context")
			}
			return "https://" + tenant + ".example.com/" + path, nil
		}
	}
}

// exAudit performs a side effect and reports nothing interesting - the typical
// shape of a function passed to Tap.
func exAudit(n int) func() (int, error) {
	return func() (int, error) {
		fmt.Println("audited:", n)
		return n, nil
	}
}

// ---------------------------------------------------------------------------
// Conversions
// ---------------------------------------------------------------------------

// ExampleFromIdiomatic demonstrates lifting a plain
// `func(context.Context) func() (A, error)` into a ReaderIOResult.
func ExampleFromIdiomatic() {
	rio := FromIdiomatic(exFetch("users"))

	fmt.Println(exRun(exCtx("acme"), rio))

	// Output:
	// https://acme.example.com/users <nil>
}

// ExampleFromResultI demonstrates lifting an already evaluated (value, error)
// pair, as produced by any ordinary Go call.
func ExampleFromResultI() {
	fmt.Println(exRun(context.Background(), FromResultI(strconv.Atoi("42"))))
	fmt.Println(exRun(context.Background(), FromResultI(strconv.Atoi("nope"))))

	// Output:
	// 42 <nil>
	// 0 strconv.Atoi: parsing "nope": invalid syntax
}

// ExampleFromReaderResultI demonstrates lifting a
// `func(context.Context) (A, error)`, the shape of most context-aware helpers.
func ExampleFromReaderResultI() {
	rio := FromReaderResultI(exTenant("orders"))

	fmt.Println(exRun(exCtx("acme"), rio))

	// Output:
	// acme/orders <nil>
}

// ExampleFromIOResultI demonstrates lifting a deferred `func() (A, error)`.
func ExampleFromIOResultI() {
	fmt.Println(exRun(context.Background(), FromIOResultI(exLoad(7))))

	// Output:
	// record-7 <nil>
}

// ---------------------------------------------------------------------------
// ChainI
// ---------------------------------------------------------------------------

// ExampleChainI demonstrates chaining an idiomatic context-aware, effectful
// function directly.
func ExampleChainI() {
	rio := F.Pipe1(
		Of("users"),
		ChainI(exFetch),
	)

	fmt.Println(exRun(exCtx("acme"), rio))

	// Output:
	// https://acme.example.com/users <nil>
}

// ExampleChainI_failure demonstrates that a missing context value surfaces as a
// failed computation.
func ExampleChainI_failure() {
	rio := F.Pipe1(
		Of("users"),
		ChainI(exFetch),
	)

	fmt.Println(exRun(context.Background(), rio))

	// Output:
	//  no tenant in context
}

// ExampleChainFirstI demonstrates running an idiomatic function for its effect
// while keeping the original value.
func ExampleChainFirstI() {
	rio := F.Pipe1(
		Of("users"),
		ChainFirstI(exFetch),
	)

	fmt.Println(exRun(exCtx("acme"), rio))

	// Output:
	// users <nil>
}

// ---------------------------------------------------------------------------
// ChainResultIK
// ---------------------------------------------------------------------------

// ExampleChainResultIK demonstrates chaining a plain `func(A) (B, error)`, the
// most common shape in existing Go code.
func ExampleChainResultIK() {
	rio := F.Pipe2(
		Of("42"),
		ChainResultIK(exParse),
		Map(N.Mul(2)),
	)

	fmt.Println(exRun(context.Background(), rio))

	// Output:
	// 84 <nil>
}

// ExampleChainResultIK_failure demonstrates that the returned error
// short-circuits the remainder of the pipeline.
func ExampleChainResultIK_failure() {
	rio := F.Pipe2(
		Of("not a number"),
		ChainResultIK(exParse),
		Map(N.Mul(2)), // never reached
	)

	fmt.Println(exRun(context.Background(), rio))

	// Output:
	// 0 strconv.Atoi: parsing "not a number": invalid syntax
}

// ExampleMonadChainResultIK demonstrates the uncurried form of ChainResultIK.
func ExampleMonadChainResultIK() {
	fmt.Println(exRun(context.Background(), MonadChainResultIK(Of("42"), exParse)))

	// Output:
	// 42 <nil>
}

// ---------------------------------------------------------------------------
// ChainReaderResultIK
// ---------------------------------------------------------------------------

// ExampleChainReaderResultIK demonstrates chaining a
// `func(A) func(context.Context) (B, error)`, i.e. a context-aware step that
// performs no IO.
func ExampleChainReaderResultIK() {
	rio := F.Pipe1(
		Of("orders"),
		ChainReaderResultIK(exTenant),
	)

	fmt.Println(exRun(exCtx("acme"), rio))

	// Output:
	// acme/orders <nil>
}

// ExampleChainReaderResultIK_failure demonstrates the error path.
func ExampleChainReaderResultIK_failure() {
	rio := F.Pipe1(
		Of("orders"),
		ChainReaderResultIK(exTenant),
	)

	fmt.Println(exRun(context.Background(), rio))

	// Output:
	//  no tenant in context
}

// ExampleTapReaderResultIK demonstrates a context-aware check that preserves the
// original value.
func ExampleTapReaderResultIK() {
	rio := F.Pipe1(
		Of("orders"),
		TapReaderResultIK(exTenant),
	)

	fmt.Println(exRun(exCtx("acme"), rio))

	// Output:
	// orders <nil>
}

// ---------------------------------------------------------------------------
// ChainIOResultIK
// ---------------------------------------------------------------------------

// ExampleChainIOResultIK demonstrates chaining a deferred
// `func(A) func() (B, error)`.
func ExampleChainIOResultIK() {
	rio := F.Pipe1(
		Of(7),
		ChainIOResultIK(exLoad),
	)

	fmt.Println(exRun(context.Background(), rio))

	// Output:
	// record-7 <nil>
}

// ExampleTapIOResultIK demonstrates an audit step: the side effect runs, its
// result is discarded and the original value flows on.
func ExampleTapIOResultIK() {
	rio := F.Pipe2(
		Of(21),
		TapIOResultIK(exAudit),
		Map(N.Mul(2)),
	)

	fmt.Println(exRun(context.Background(), rio))

	// Output:
	// audited: 21
	// 42 <nil>
}

// ---------------------------------------------------------------------------
// ChainOptionIK
// ---------------------------------------------------------------------------

// ExampleChainOptionIK demonstrates chaining a comma-ok lookup, turning a false
// flag into an error.
func ExampleChainOptionIK() {
	table := map[string]string{"host": "localhost"}
	lookup := func(k string) (string, bool) {
		v, ok := table[k]
		return v, ok
	}
	chain := ChainOptionIK[string, string](func() error {
		return errors.New("key not found")
	})

	fmt.Println(exRun(context.Background(), F.Pipe1(Of("host"), chain(lookup))))
	fmt.Println(exRun(context.Background(), F.Pipe1(Of("port"), chain(lookup))))

	// Output:
	// localhost <nil>
	//  key not found
}

// ---------------------------------------------------------------------------
// Putting it together
// ---------------------------------------------------------------------------

// Example_idiomaticPipeline mixes all four levels - Result, ReaderResult,
// IOResult and ReaderIOResult - in one pipeline built entirely from plain Go
// functions.
func Example_idiomaticPipeline() {
	rio := F.Pipe3(
		Of("7"),
		ChainResultIK(exParse),  // "7" -> 7
		ChainIOResultIK(exLoad), // 7 -> "record-7"
		ChainI(exFetch),         // -> "https://acme.example.com/record-7"
	)

	fmt.Println(exRun(exCtx("acme"), rio))

	// Output:
	// https://acme.example.com/record-7 <nil>
}
