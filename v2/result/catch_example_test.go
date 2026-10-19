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

package result

import (
	"errors"
	"fmt"
	"io/fs"

	F "github.com/IBM/fp-go/v2/function"
)

// missingKeyError is the error of a failed lookup in the examples.
type missingKeyError struct{ key string }

func (e *missingKeyError) Error() string { return fmt.Sprintf("missing key %q", e.key) }

// lookupPort fails with a wrapped missingKeyError.
func lookupPort(key string) Result[int] {
	return Left[int](fmt.Errorf("config: %w", &missingKeyError{key: key}))
}

// ExampleCatchAs recovers from a missing key with a default value.
func ExampleCatchAs() {
	port := F.Pipe1(
		lookupPort("port"),
		CatchAs(func(*missingKeyError) Result[int] { return Of(8080) }),
	)

	fmt.Println(port)
	// Output:
	// Right[int](8080)
}

// ExampleCatchIs recovers from a sentinel error.
func ExampleCatchIs() {
	readConfig := F.Pipe1(
		Left[string](fmt.Errorf("open config.yaml: %w", fs.ErrNotExist)),
		CatchIs(fs.ErrNotExist, func(error) Result[string] { return Of("defaults") }),
	)

	fmt.Println(readConfig)
	// Output:
	// Right[string](defaults)
}

// ExampleMapErrorAs replaces a missing key error by a more general error.
func ExampleMapErrorAs() {
	errInvalidConfig := errors.New("invalid configuration")

	port := F.Pipe1(
		lookupPort("port"),
		MapErrorAs[*missingKeyError, int](func(e *missingKeyError) error {
			return fmt.Errorf("%w: %s", errInvalidConfig, e)
		}),
	)

	fmt.Println(port)
	// Output:
	// Left[*fmt.wrapError](invalid configuration: missing key "port")
}

// ExampleTapErrorAs reports a missing key and keeps the error.
func ExampleTapErrorAs() {
	port := F.Pipe1(
		lookupPort("port"),
		TapErrorAs[*missingKeyError, int](func(e *missingKeyError) Result[F.Void] {
			fmt.Println("please set", e.key)
			return Of(F.VOID)
		}),
	)

	fmt.Println(port)
	// Output:
	// please set port
	// Left[*fmt.wrapError](config: missing key "port")
}
