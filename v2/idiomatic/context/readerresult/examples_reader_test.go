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

	RES "github.com/IBM/fp-go/v2/result"
)

// ExampleFromEither demonstrates lifting a Result (Either) into a
// The resulting ReaderResult ignores the context and returns the Result value.
func ExampleFromEither() {
	res := RES.Of(42)
	rr := FromEither(res)
	value, err := rr(context.Background())
	fmt.Println(value, err)
	// Output:
	// 42 <nil>
}

// ExampleFromResult demonstrates creating a ReaderResult from a Go-style (value, error) tuple.
// This is useful for converting standard Go error handling into the ReaderResult monad.
func ExampleFromResult() {
	rr := FromResult(42, nil)
	value, err := rr(context.Background())
	fmt.Println(value, err)
	// Output:
	// 42 <nil>
}

// ExampleFromResult_error demonstrates creating a ReaderResult from an error case.
// The resulting ReaderResult will propagate the error when executed.
func ExampleFromResult_error() {
	rr := FromResult(0, errors.New("failed"))
	value, err := rr(context.Background())
	fmt.Println(value, err != nil)
	// Output:
	// 0 true
}

// ExampleLeft demonstrates creating a ReaderResult that always fails with an error.
// This is the error constructor for ReaderResult, analogous to Either's Left.
func ExampleLeft() {
	rr := Left[int](errors.New("failed"))
	value, err := rr(context.Background())
	fmt.Println(value, err != nil)
	// Output:
	// 0 true
}

// ExampleRight demonstrates creating a ReaderResult that always succeeds with a value.
// This is the success constructor for ReaderResult, analogous to Either's Right.
func ExampleRight() {
	rr := Right(42)
	value, err := rr(context.Background())
	fmt.Println(value, err)
	// Output:
	// 42 <nil>
}

// ExampleOf demonstrates the monadic return/pure operation for
// It creates a ReaderResult that always succeeds with the given value.
func ExampleOf() {
	rr := Of(42)
	value, err := rr(context.Background())
	fmt.Println(value, err)
	// Output:
	// 42 <nil>
}

// ExampleAsk demonstrates getting the context.Context environment.
// This returns a ReaderResult that provides access to the context itself.
func ExampleAsk() {
	rr := Ask()
	ctx := context.Background()
	value, err := rr(ctx)
	fmt.Println(value == ctx, err)
	// Output:
	// true <nil>
}

// ExampleAsks demonstrates extracting a value from the context using a function.
// This is useful for accessing configuration or other data stored in the context.
func ExampleAsks() {
	type Config struct {
		Port int
	}

	getPort := Asks(func(ctx context.Context) int {
		// In real code, extract config from context
		return 8080
	})

	value, err := getPort(context.Background())
	fmt.Println(value, err)
	// Output:
	// 8080 <nil>
}
