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
)

// ---------------------------------------------------------------------------
// Variants
//
// Examples for the uncurried (Monad*), value-preserving (ChainFirst*, Tap*) and
// Either-named forms of the idiomatic operators. They reuse the fixtures from
// idiomatic_example_test.go. The side-effect fixtures below print, so each
// example's output shows both that the effect ran and which value survived.
// ---------------------------------------------------------------------------

// exVarValidate parses s and reports success - the shape of a validation step.
func exVarValidate(s string) (int, error) {
	n, err := strconv.Atoi(s)
	if err == nil {
		fmt.Println("valid:", n)
	}
	return n, err
}

// exVarAuthorize reads the tenant from the context and reports the access.
func exVarAuthorize(resource string) func(context.Context) (string, error) {
	return func(ctx context.Context) (string, error) {
		tenant, ok := ctx.Value(exTenantKey{}).(string)
		if !ok {
			return "", errors.New("no tenant in context")
		}
		fmt.Println("authorized:", tenant+"/"+resource)
		return tenant, nil
	}
}

// exVarFetchLogged reads the context, performs IO and reports what it fetched.
func exVarFetchLogged(path string) func(context.Context) func() (string, error) {
	return func(ctx context.Context) func() (string, error) {
		return func() (string, error) {
			tenant, ok := ctx.Value(exTenantKey{}).(string)
			if !ok {
				return "", errors.New("no tenant in context")
			}
			fmt.Println("fetched:", tenant+"/"+path)
			return tenant + "/" + path, nil
		}
	}
}

// ---------------------------------------------------------------------------
// ReaderIOResult level
// ---------------------------------------------------------------------------

// ExampleFromReaderIOResultI demonstrates the conversion that FromIdiomatic is
// an alias of, including the failure path.
func ExampleFromReaderIOResultI() {
	fmt.Println(exRun(exCtx("acme"), FromReaderIOResultI(exFetch("users"))))
	fmt.Println(exRun(context.Background(), FromReaderIOResultI(exFetch("users"))))

	// Output:
	// https://acme.example.com/users <nil>
	//  no tenant in context
}

// ExampleMonadChainI demonstrates the uncurried form of ChainI.
func ExampleMonadChainI() {
	fmt.Println(exRun(exCtx("acme"), MonadChainI(Of("users"), exFetch)))

	// Output:
	// https://acme.example.com/users <nil>
}

// ExampleMonadChainFirstI demonstrates running a context-aware, effectful
// function while keeping the original value.
func ExampleMonadChainFirstI() {
	fmt.Println(exRun(exCtx("acme"), MonadChainFirstI(Of("users"), exVarFetchLogged)))

	// Output:
	// fetched: acme/users
	// users <nil>
}

// ExampleMonadTapI demonstrates the Tap alias of MonadChainFirstI.
func ExampleMonadTapI() {
	fmt.Println(exRun(exCtx("acme"), MonadTapI(Of("users"), exVarFetchLogged)))

	// Output:
	// fetched: acme/users
	// users <nil>
}

// ExampleTapI demonstrates that a failing side effect fails the whole
// computation, even though its successful result would have been discarded.
func ExampleTapI() {
	fmt.Println(exRun(exCtx("acme"), F.Pipe1(Of("users"), TapI(exVarFetchLogged))))
	fmt.Println(exRun(context.Background(), F.Pipe1(Of("users"), TapI(exVarFetchLogged))))

	// Output:
	// fetched: acme/users
	// users <nil>
	//  no tenant in context
}

// ---------------------------------------------------------------------------
// Result level
// ---------------------------------------------------------------------------

// ExampleMonadChainEitherIK demonstrates the Either-named alias of MonadChainResultIK.
func ExampleMonadChainEitherIK() {
	fmt.Println(exRun(context.Background(), MonadChainEitherIK(Of("42"), exParse)))

	// Output:
	// 42 <nil>
}

// ExampleChainEitherIK demonstrates the Either-named alias of ChainResultIK.
func ExampleChainEitherIK() {
	fmt.Println(exRun(context.Background(), F.Pipe1(Of("42"), ChainEitherIK(exParse))))

	// Output:
	// 42 <nil>
}

// ExampleMonadChainFirstEitherIK demonstrates a validation step in uncurried form.
func ExampleMonadChainFirstEitherIK() {
	fmt.Println(exRun(context.Background(), MonadChainFirstEitherIK(Of("42"), exVarValidate)))

	// Output:
	// valid: 42
	// 42 <nil>
}

// ExampleMonadChainFirstResultIK demonstrates the Result-named alias of
// MonadChainFirstEitherIK.
func ExampleMonadChainFirstResultIK() {
	fmt.Println(exRun(context.Background(), MonadChainFirstResultIK(Of("42"), exVarValidate)))

	// Output:
	// valid: 42
	// 42 <nil>
}

// ExampleMonadTapEitherIK demonstrates the Tap alias of MonadChainFirstEitherIK.
func ExampleMonadTapEitherIK() {
	fmt.Println(exRun(context.Background(), MonadTapEitherIK(Of("42"), exVarValidate)))

	// Output:
	// valid: 42
	// 42 <nil>
}

// ExampleMonadTapResultIK demonstrates the Tap alias of MonadChainFirstResultIK.
func ExampleMonadTapResultIK() {
	fmt.Println(exRun(context.Background(), MonadTapResultIK(Of("42"), exVarValidate)))

	// Output:
	// valid: 42
	// 42 <nil>
}

// ExampleChainFirstEitherIK demonstrates a validation step that keeps the
// original value.
func ExampleChainFirstEitherIK() {
	fmt.Println(exRun(context.Background(), F.Pipe1(Of("42"), ChainFirstEitherIK(exVarValidate))))

	// Output:
	// valid: 42
	// 42 <nil>
}

// ExampleChainFirstResultIK demonstrates that a failed validation fails the
// computation.
func ExampleChainFirstResultIK() {
	fmt.Println(exRun(context.Background(), F.Pipe1(Of("42"), ChainFirstResultIK(exVarValidate))))
	fmt.Println(exRun(context.Background(), F.Pipe1(Of("nope"), ChainFirstResultIK(exVarValidate))))

	// Output:
	// valid: 42
	// 42 <nil>
	//  strconv.Atoi: parsing "nope": invalid syntax
}

// ExampleTapEitherIK demonstrates the Tap alias of ChainFirstEitherIK.
func ExampleTapEitherIK() {
	fmt.Println(exRun(context.Background(), F.Pipe1(Of("42"), TapEitherIK(exVarValidate))))

	// Output:
	// valid: 42
	// 42 <nil>
}

// ExampleTapResultIK demonstrates the Tap alias of ChainFirstResultIK.
func ExampleTapResultIK() {
	fmt.Println(exRun(context.Background(), F.Pipe1(Of("42"), TapResultIK(exVarValidate))))

	// Output:
	// valid: 42
	// 42 <nil>
}

// ---------------------------------------------------------------------------
// IOResult level
// ---------------------------------------------------------------------------

// ExampleMonadChainIOEitherIK demonstrates the Either-named, uncurried form of
// ChainIOResultIK.
func ExampleMonadChainIOEitherIK() {
	fmt.Println(exRun(context.Background(), MonadChainIOEitherIK(Of(7), exLoad)))

	// Output:
	// record-7 <nil>
}

// ExampleMonadChainIOResultIK demonstrates the uncurried form of ChainIOResultIK.
func ExampleMonadChainIOResultIK() {
	fmt.Println(exRun(context.Background(), MonadChainIOResultIK(Of(7), exLoad)))

	// Output:
	// record-7 <nil>
}

// ExampleChainIOEitherIK demonstrates the Either-named alias of ChainIOResultIK,
// including its failure path.
func ExampleChainIOEitherIK() {
	fmt.Println(exRun(context.Background(), F.Pipe1(Of(7), ChainIOEitherIK(exLoad))))
	fmt.Println(exRun(context.Background(), F.Pipe1(Of(-1), ChainIOEitherIK(exLoad))))

	// Output:
	// record-7 <nil>
	//  cannot load -1
}

// ExampleChainFirstIOEitherIK demonstrates running a deferred side effect while
// keeping the original value.
func ExampleChainFirstIOEitherIK() {
	fmt.Println(exRun(context.Background(), F.Pipe1(Of(21), ChainFirstIOEitherIK(exAudit))))

	// Output:
	// audited: 21
	// 21 <nil>
}

// ExampleChainFirstIOResultIK demonstrates the Result-named alias of
// ChainFirstIOEitherIK.
func ExampleChainFirstIOResultIK() {
	fmt.Println(exRun(context.Background(), F.Pipe1(Of(21), ChainFirstIOResultIK(exAudit))))

	// Output:
	// audited: 21
	// 21 <nil>
}

// ExampleTapIOEitherIK demonstrates the Tap alias of ChainFirstIOEitherIK.
func ExampleTapIOEitherIK() {
	fmt.Println(exRun(context.Background(), F.Pipe1(Of(21), TapIOEitherIK(exAudit))))

	// Output:
	// audited: 21
	// 21 <nil>
}

// ---------------------------------------------------------------------------
// ReaderResult level
// ---------------------------------------------------------------------------

// ExampleMonadChainReaderEitherIK demonstrates the Either-named, uncurried form
// of ChainReaderResultIK.
func ExampleMonadChainReaderEitherIK() {
	fmt.Println(exRun(exCtx("acme"), MonadChainReaderEitherIK(Of("orders"), exTenant)))

	// Output:
	// acme/orders <nil>
}

// ExampleMonadChainReaderResultIK demonstrates the uncurried form of
// ChainReaderResultIK.
func ExampleMonadChainReaderResultIK() {
	fmt.Println(exRun(exCtx("acme"), MonadChainReaderResultIK(Of("orders"), exTenant)))

	// Output:
	// acme/orders <nil>
}

// ExampleChainReaderEitherIK demonstrates the Either-named alias of
// ChainReaderResultIK.
func ExampleChainReaderEitherIK() {
	fmt.Println(exRun(exCtx("acme"), F.Pipe1(Of("orders"), ChainReaderEitherIK(exTenant))))

	// Output:
	// acme/orders <nil>
}

// ExampleMonadChainFirstReaderEitherIK demonstrates a context-dependent check in
// uncurried form that keeps the original value.
func ExampleMonadChainFirstReaderEitherIK() {
	fmt.Println(exRun(exCtx("acme"), MonadChainFirstReaderEitherIK(Of("orders"), exVarAuthorize)))

	// Output:
	// authorized: acme/orders
	// orders <nil>
}

// ExampleMonadChainFirstReaderResultIK demonstrates the Result-named alias of
// MonadChainFirstReaderEitherIK.
func ExampleMonadChainFirstReaderResultIK() {
	fmt.Println(exRun(exCtx("acme"), MonadChainFirstReaderResultIK(Of("orders"), exVarAuthorize)))

	// Output:
	// authorized: acme/orders
	// orders <nil>
}

// ExampleMonadTapReaderEitherIK demonstrates the Tap alias of
// MonadChainFirstReaderEitherIK.
func ExampleMonadTapReaderEitherIK() {
	fmt.Println(exRun(exCtx("acme"), MonadTapReaderEitherIK(Of("orders"), exVarAuthorize)))

	// Output:
	// authorized: acme/orders
	// orders <nil>
}

// ExampleMonadTapReaderResultIK demonstrates the Tap alias of
// MonadChainFirstReaderResultIK.
func ExampleMonadTapReaderResultIK() {
	fmt.Println(exRun(exCtx("acme"), MonadTapReaderResultIK(Of("orders"), exVarAuthorize)))

	// Output:
	// authorized: acme/orders
	// orders <nil>
}

// ExampleChainFirstReaderEitherIK demonstrates a context-dependent check that
// keeps the original value.
func ExampleChainFirstReaderEitherIK() {
	fmt.Println(exRun(exCtx("acme"), F.Pipe1(Of("orders"), ChainFirstReaderEitherIK(exVarAuthorize))))

	// Output:
	// authorized: acme/orders
	// orders <nil>
}

// ExampleChainFirstReaderResultIK demonstrates that a failed context-dependent
// check fails the computation.
func ExampleChainFirstReaderResultIK() {
	fmt.Println(exRun(exCtx("acme"), F.Pipe1(Of("orders"), ChainFirstReaderResultIK(exVarAuthorize))))
	fmt.Println(exRun(context.Background(), F.Pipe1(Of("orders"), ChainFirstReaderResultIK(exVarAuthorize))))

	// Output:
	// authorized: acme/orders
	// orders <nil>
	//  no tenant in context
}

// ExampleTapReaderEitherIK demonstrates the Tap alias of ChainFirstReaderEitherIK.
func ExampleTapReaderEitherIK() {
	fmt.Println(exRun(exCtx("acme"), F.Pipe1(Of("orders"), TapReaderEitherIK(exVarAuthorize))))

	// Output:
	// authorized: acme/orders
	// orders <nil>
}
