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

package lenses

import (
	"context"
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/lazy"
	O "github.com/IBM/fp-go/v2/option"
)

// ExampleAtContext_get demonstrates reading a typed value from a context.
// Get returns Some when the key is present with the right type, None otherwise.
func ExampleAtContext_get() {
	type reqKey string

	lens := AtContext[string, reqKey]("user")

	// Key is present and has the expected type.
	ctx := context.WithValue(context.Background(), reqKey("user"), "alice")
	fmt.Println(lens.Get(ctx))

	// Key is absent.
	fmt.Println(lens.Get(context.Background()))

	// Output:
	// Some[string](alice)
	// None[string]
}

// ExampleAtContext_set demonstrates writing a typed value into a context.
// Set(Some(v)) creates a child context carrying the value; Set(None) is a no-op.
func ExampleAtContext_set() {
	type reqKey string

	lens := AtContext[string, reqKey]("role")
	ctx := context.Background()

	// Set(Some("admin")) derives a child context that carries the value.
	withRole := lens.Set(O.Some("admin"))(ctx)
	fmt.Println(lens.Get(withRole))

	// Set(None) leaves the context unchanged.
	unchanged := lens.Set(O.None[string]())(withRole)
	fmt.Println(lens.Get(unchanged))

	// Output:
	// Some[string](admin)
	// Some[string](admin)
}

// ExampleAtContext_typeMismatch demonstrates that Get returns None when the stored
// value has an incompatible type, making retrieval type-safe.
func ExampleAtContext_typeMismatch() {
	type reqKey string

	// Store a string under the key, but read it back as int.
	ctx := context.WithValue(context.Background(), reqKey("count"), "not-a-number")

	fmt.Println(AtContext[int, reqKey]("count").Get(ctx))

	// Output:
	// None[int]
}

// ExampleAtContext_getOrElse demonstrates extracting a value with a fallback
// using option.GetOrElse together with the lens.
func ExampleAtContext_getOrElse() {
	type reqKey string

	lens := AtContext[int, reqKey]("timeout")

	withTimeout := context.WithValue(context.Background(), reqKey("timeout"), 30)
	withoutTimeout := context.Background()

	extract := F.Flow2(lens.Get, O.GetOrElse(lazy.Of(5)))

	fmt.Println(extract(withTimeout))
	fmt.Println(extract(withoutTimeout))

	// Output:
	// 30
	// 5
}

// ExampleAtContext_multipleKeys demonstrates that lenses for different keys
// are independent and can be composed sequentially.
func ExampleAtContext_multipleKeys() {
	type reqKey string

	userLens  := AtContext[string, reqKey]("user")
	tokenLens := AtContext[string, reqKey]("token")

	ctx := F.Pipe1(
		context.Background(),
		F.Flow2(userLens.Set(O.Some("alice")), tokenLens.Set(O.Some("tok-abc"))),
	)

	fmt.Println(userLens.Get(ctx))
	fmt.Println(tokenLens.Get(ctx))

	// Output:
	// Some[string](alice)
	// Some[string](tok-abc)
}
