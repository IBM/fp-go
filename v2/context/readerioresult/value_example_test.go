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
	"time"

	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
)

type exampleKey string

const (
	exampleUserKey      exampleKey = "user"
	exampleRequestIDKey exampleKey = "requestID"
)

func ExampleAskValue() {
	getUser := AskValue[string](exampleUserKey)

	ctx := context.WithValue(context.Background(), exampleUserKey, "Alice")
	fmt.Println(getUser(ctx)())
	fmt.Println(getUser(context.Background())())

	// Output:
	// Right[common.Option[string]](Some[string](Alice))
	// Right[common.Option[string]](None[string])
}

// A value that must be present can be turned into an error with FromOption.
func ExampleAskValue_required() {
	errNoUser := errors.New("no user in context")

	requireUser := F.Pipe1(
		AskValue[string](exampleUserKey),
		Chain(FromOption[string](F.Constant(errNoUser))),
	)

	ctx := context.WithValue(context.Background(), exampleUserKey, "Alice")
	fmt.Println(requireUser(ctx)())
	fmt.Println(requireUser(context.Background())())

	// Output:
	// Right[string](Alice)
	// Left[*errors.errorString](no user in context)
}

func ExampleWithValue() {
	greet := F.Pipe1(
		AskValue[string](exampleUserKey),
		Map(F.Flow2(
			O.GetOrElse(F.Constant("anonymous")),
			func(name string) string { return "Hello, " + name },
		)),
	)

	fmt.Println(greet(context.Background())())
	fmt.Println(F.Pipe1(greet, WithValue[string](exampleUserKey, "Alice"))(context.Background())())

	// Output:
	// Right[string](Hello, anonymous)
	// Right[string](Hello, Alice)
}

// WithValue, WithTimeout and friends compose as ordinary operators, so request
// scoping can be expressed declaratively around a handler.
func ExampleWithValue_requestScope() {
	handler := F.Pipe1(
		Ask(),
		Map(func(ctx context.Context) string {
			_, hasDeadline := ctx.Deadline()
			return fmt.Sprintf("user=%v request=%v deadline=%v",
				ctx.Value(exampleUserKey), ctx.Value(exampleRequestIDKey), hasDeadline)
		}),
	)

	scoped := F.Pipe3(
		handler,
		WithTimeout[string](time.Second),
		WithValue[string](exampleUserKey, "Alice"),
		WithValue[string](exampleRequestIDKey, 42),
	)

	fmt.Println(scoped(context.Background())())

	// Output:
	// Right[string](user=Alice request=42 deadline=true)
}
