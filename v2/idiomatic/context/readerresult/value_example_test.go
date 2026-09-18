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

package readerresult

import (
	"context"
	"errors"
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
)

type exampleKey string

const exampleUserKey exampleKey = "user"

func ExampleAskValue() {
	getUser := AskValue[string](exampleUserKey)

	ctx := context.WithValue(context.Background(), exampleUserKey, "Alice")
	fmt.Println(getUser(ctx))
	fmt.Println(getUser(context.Background()))

	// Output:
	// Some[string](Alice) <nil>
	// None[string] <nil>
}

// A value that must be present can be turned into an error with ChainOptionK.
func ExampleAskValue_required() {
	errNoUser := errors.New("no user in context")

	requireUser := F.Pipe1(
		AskValue[string](exampleUserKey),
		ChainOptionK[O.Option[string], string](F.Constant(errNoUser))(O.Unwrap[string]),
	)

	ctx := context.WithValue(context.Background(), exampleUserKey, "Alice")
	fmt.Println(requireUser(ctx))
	fmt.Println(requireUser(context.Background()))

	// Output:
	// Alice <nil>
	//  no user in context
}

func ExampleWithValue() {
	greet := F.Pipe1(
		AskValue[string](exampleUserKey),
		Map(O.GetOrElse(F.Constant("anonymous"))),
	)

	fmt.Println(greet(context.Background()))
	fmt.Println(F.Pipe1(greet, WithValue[string](exampleUserKey, "Alice"))(context.Background()))

	// Output:
	// anonymous <nil>
	// Alice <nil>
}
