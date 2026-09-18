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
	"time"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
)

type exampleKey string

const exampleUserKey exampleKey = "user"

// slowOperation waits for the given duration unless its context is done first
func slowOperation(d time.Duration) ReaderResult[string] {
	return func(ctx context.Context) Result[string] {
		select {
		case <-time.After(d):
			return E.Of[error]("done")
		case <-ctx.Done():
			return E.Left[string](ctx.Err())
		}
	}
}

func ExampleAskValue() {
	getUser := AskValue[string](exampleUserKey)

	ctx := context.WithValue(context.Background(), exampleUserKey, "Alice")
	fmt.Println(getUser(ctx))
	fmt.Println(getUser(context.Background()))

	// Output:
	// Right[common.Option[string]](Some[string](Alice))
	// Right[common.Option[string]](None[string])
}

// A value that must be present can be turned into an error with ChainOptionK.
func ExampleAskValue_required() {
	errNoUser := errors.New("no user in context")

	requireUser := F.Pipe1(
		AskValue[string](exampleUserKey),
		ChainOptionK[O.Option[string], string](F.Constant(errNoUser))(F.Identity[O.Option[string]]),
	)

	ctx := context.WithValue(context.Background(), exampleUserKey, "Alice")
	fmt.Println(requireUser(ctx))
	fmt.Println(requireUser(context.Background()))

	// Output:
	// Right[string](Alice)
	// Left[*errors.errorString](no user in context)
}

func ExampleWithValue() {
	greet := F.Pipe1(
		AskValue[string](exampleUserKey),
		Map(O.GetOrElse(F.Constant("anonymous"))),
	)

	fmt.Println(greet(context.Background()))
	fmt.Println(F.Pipe1(greet, WithValue[string](exampleUserKey, "Alice"))(context.Background()))

	// Output:
	// Right[string](anonymous)
	// Right[string](Alice)
}

func ExampleWithTimeout() {
	fast := F.Pipe1(slowOperation(time.Millisecond), WithTimeout[string](time.Second))
	slow := F.Pipe1(slowOperation(time.Second), WithTimeout[string](10*time.Millisecond))

	fmt.Println(fast(context.Background()))
	fmt.Println(slow(context.Background()))

	// Output:
	// Right[string](done)
	// Left[context.deadlineExceededError](context deadline exceeded)
}

func ExampleWithDeadline() {
	expired := F.Pipe1(
		slowOperation(time.Second),
		WithDeadline[string](time.Now().Add(-time.Second)),
	)

	fmt.Println(expired(context.Background()))

	// Output:
	// Left[context.deadlineExceededError](context deadline exceeded)
}
