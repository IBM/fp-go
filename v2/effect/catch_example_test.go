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
	"errors"
	"fmt"
	"io/fs"

	F "github.com/IBM/fp-go/v2/function"
)

// exCatchDeps are the dependencies of the error-type examples.
type exCatchDeps struct {
	DefaultPort int
}

// exMissingKeyError is the error of a failed lookup in the examples.
type exMissingKeyError struct{ key string }

func (e *exMissingKeyError) Error() string { return fmt.Sprintf("missing key %q", e.key) }

// exLookupPort fails with a wrapped exMissingKeyError.
func exLookupPort(key string) Effect[exCatchDeps, int] {
	return Fail[exCatchDeps, int](fmt.Errorf("config: %w", &exMissingKeyError{key: key}))
}

// exDefaultPort reads the default port from the dependencies.
func exDefaultPort(deps exCatchDeps) int { return deps.DefaultPort }

// ExampleCatchAs recovers from a missing key with the default port from the dependencies.
func ExampleCatchAs() {
	port := F.Pipe1(
		exLookupPort("port"),
		CatchAs(func(*exMissingKeyError) Effect[exCatchDeps, int] {
			return Asks(exDefaultPort)
		}),
	)

	fmt.Println(port(exCatchDeps{DefaultPort: 8080})(context.Background())())
	// Output:
	// Right[int](8080)
}

// ExampleCatchIs recovers from a sentinel error.
func ExampleCatchIs() {
	readConfig := F.Pipe1(
		Fail[exCatchDeps, string](fmt.Errorf("open config.yaml: %w", fs.ErrNotExist)),
		CatchIs(fs.ErrNotExist, func(error) Effect[exCatchDeps, string] {
			return Of[exCatchDeps]("defaults")
		}),
	)

	fmt.Println(readConfig(exCatchDeps{})(context.Background())())
	// Output:
	// Right[string](defaults)
}

// ExampleMapErrorAs replaces a missing key error by a more general error.
func ExampleMapErrorAs() {
	errInvalidConfig := errors.New("invalid configuration")

	port := F.Pipe1(
		exLookupPort("port"),
		MapErrorAs[*exMissingKeyError, exCatchDeps, int](func(e *exMissingKeyError) error {
			return fmt.Errorf("%w: %s", errInvalidConfig, e)
		}),
	)

	fmt.Println(port(exCatchDeps{})(context.Background())())
	// Output:
	// Left[*fmt.wrapError](invalid configuration: missing key "port")
}

// ExampleTapErrorAs reports a missing key and keeps the error.
func ExampleTapErrorAs() {
	port := F.Pipe1(
		exLookupPort("port"),
		TapErrorAs[*exMissingKeyError, exCatchDeps, int](func(e *exMissingKeyError) Effect[exCatchDeps, Void] {
			return FromIO[exCatchDeps](func() Void {
				fmt.Println("please set", e.key)
				return F.VOID
			})
		}),
	)

	fmt.Println(port(exCatchDeps{})(context.Background())())
	// Output:
	// please set port
	// Left[*fmt.wrapError](config: missing key "port")
}
