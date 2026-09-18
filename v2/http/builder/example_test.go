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

package builder

import (
	"fmt"

	B "github.com/IBM/fp-go/v2/bytes"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	C "github.com/IBM/fp-go/v2/http/content"
	FD "github.com/IBM/fp-go/v2/http/form"
	H "github.com/IBM/fp-go/v2/http/headers"
	O "github.com/IBM/fp-go/v2/option"
)

var (
	// renderTarget renders the target URL of a [Builder], the error message on failure
	renderTarget = F.Flow2(
		(*Builder).GetTargetURL,
		E.GetOrElse(error.Error),
	)

	// renderBody renders the body of a [Builder] as a string
	renderBody = F.Flow3(
		(*Builder).GetBody,
		O.Map(E.Fold(error.Error, B.ToString)),
		O.GetOrElse(F.Constant("<no body>")),
	)
)

// ExampleDefault demonstrates the starting point of every request, an empty GET builder
func ExampleDefault() {
	fmt.Println(Default.GetMethod())
	fmt.Println(renderBody(Default))

	// Output:
	// GET
	// <no body>
}

// ExampleWithURL demonstrates how to assemble a request from composable operations.
// The builder is immutable, every operation returns a new instance.
func ExampleWithURL() {
	b := F.Pipe2(
		Default,
		WithURL("https://api.example.com/users"),
		WithHeader(H.Accept)(C.JSON),
	)

	fmt.Println(b.GetMethod())
	fmt.Println(renderTarget(b))
	fmt.Println(b.GetHeader(H.Accept))

	// Output:
	// GET
	// https://api.example.com/users
	// Some[string](application/json)
}

// ExampleBuilder_GetTargetURL demonstrates that the query parameters of the builder are
// merged into the query parameters already present in the URL
func ExampleBuilder_GetTargetURL() {
	b := F.Pipe3(
		Default,
		WithURL("https://api.example.com/search?lang=en"),
		WithQueryArg("q")("golang"),
		WithQueryArg("limit")("10"),
	)

	fmt.Println(renderTarget(b))

	// Output:
	// https://api.example.com/search?lang=en&limit=10&q=golang
}

// ExampleBuilder_GetTargetURL_invalid demonstrates that an unparseable URL surfaces as
// an error rather than as a panic
func ExampleBuilder_GetTargetURL_invalid() {
	fmt.Println(renderTarget(F.Pipe1(Default, WithURL("://not-a-url"))))

	// Output:
	// parse "://not-a-url": missing protocol scheme
}

// ExampleWithJSON demonstrates a JSON payload, which also sets the content type
func ExampleWithJSON() {
	b := F.Pipe3(
		Default,
		WithURL("https://api.example.com/users"),
		WithPost,
		WithJSON(map[string]string{"name": "John"}),
	)

	fmt.Println(b.GetMethod())
	fmt.Println(b.GetHeader(H.ContentType))
	fmt.Println(renderBody(b))

	// Output:
	// POST
	// Some[string](application/json)
	// {"name":"John"}
}

// ExampleWithFormData demonstrates a form encoded payload, which also sets the content type
func ExampleWithFormData() {
	data := F.Pipe2(
		FD.Default,
		FD.WithValue("name")("John"),
		FD.WithValue("city")("Berlin"),
	)

	b := F.Pipe2(
		Default,
		WithPost,
		WithFormData(data),
	)

	fmt.Println(b.GetHeader(H.ContentType))
	fmt.Println(renderBody(b))

	// Output:
	// Some[string](application/x-www-form-urlencoded)
	// city=Berlin&name=John
}

// ExampleWithBearer demonstrates the shortcut for a Bearer authorization header
func ExampleWithBearer() {
	fmt.Println(F.Pipe1(Default, WithBearer("s3cr3t")).GetHeader(H.Authorization))

	// Output:
	// Some[string](Bearer s3cr3t)
}

// ExampleHeader demonstrates the [Lens] for a single header. Setting the focus to none
// removes the header, the source builder is never modified.
func ExampleHeader() {
	accept := Header(H.Accept)

	withJSON := F.Pipe1(Default, accept.Set(O.Of(C.JSON)))
	without := F.Pipe1(withJSON, accept.Set(O.None[string]()))

	fmt.Println(accept)
	fmt.Println(accept.Get(withJSON))
	fmt.Println(accept.Get(without))
	fmt.Println(accept.Get(Default))

	// Output:
	// HttpHeader[accept]
	// Some[string](application/json)
	// None[string]
	// None[string]
}

// ExampleQueryArg demonstrates the [Lens] for a single query parameter
func ExampleQueryArg() {
	limit := QueryArg("limit")

	b := F.Pipe2(
		Default,
		WithURL("https://api.example.com/users"),
		limit.Set(O.Of("10")),
	)

	fmt.Println(limit.Get(b))
	fmt.Println(renderTarget(b))
	fmt.Println(renderTarget(F.Pipe1(b, limit.Set(O.None[string]()))))

	// Output:
	// Some[string](10)
	// https://api.example.com/users?limit=10
	// https://api.example.com/users
}

// ExampleMonoid demonstrates that builder operations are endomorphisms, so they can be
// combined into reusable bundles before being applied
func ExampleMonoid() {
	// a reusable bundle of operations
	jsonPost := Monoid.Concat(WithPost, WithContentType(C.JSON))

	b := F.Pipe2(
		Default,
		WithURL("https://api.example.com/users"),
		jsonPost,
	)

	fmt.Println(b.GetMethod())
	fmt.Println(b.GetHeader(H.ContentType))

	// Output:
	// POST
	// Some[string](application/json)
}

// ExampleMakeHash demonstrates the cache key of a builder. It only depends on the
// content of the request, not on the order in which the builder was assembled.
func ExampleMakeHash() {
	b1 := F.Pipe2(
		Default,
		WithURL("https://api.example.com/users"),
		WithContentType(C.JSON),
	)

	b2 := F.Pipe2(
		Default,
		WithContentType(C.JSON),
		WithURL("https://api.example.com/users"),
	)

	fmt.Println(MakeHash(b1) == MakeHash(b2))
	fmt.Println(MakeHash(b1) == MakeHash(Default))

	// Output:
	// true
	// false
}
