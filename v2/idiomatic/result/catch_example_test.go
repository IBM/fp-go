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
func lookupPort(key string) (int, error) {
	return Left[int](fmt.Errorf("config: %w", &missingKeyError{key: key}))
}

// ExampleCatchAs recovers from a missing key with a default value.
func ExampleCatchAs() {
	port, err := CatchAs(func(*missingKeyError) (int, error) { return Of(8080) })(lookupPort("port"))

	fmt.Println(port, err)
	// Output:
	// 8080 <nil>
}

// ExampleCatchIs recovers from a sentinel error.
func ExampleCatchIs() {
	readConfig := CatchIs(fs.ErrNotExist, func(error) (string, error) { return Of("defaults") })

	fmt.Println(readConfig(Left[string](fmt.Errorf("open config.yaml: %w", fs.ErrNotExist))))
	// Output:
	// defaults <nil>
}

// ExampleMapErrorAs replaces a missing key error by a more general error.
func ExampleMapErrorAs() {
	errInvalidConfig := errors.New("invalid configuration")

	toInvalidConfig := MapErrorAs[*missingKeyError, int](func(e *missingKeyError) error {
		return fmt.Errorf("%w: %s", errInvalidConfig, e)
	})

	fmt.Println(toInvalidConfig(lookupPort("port")))
	// Output:
	// 0 invalid configuration: missing key "port"
}

// ExampleTapErrorAs reports a missing key and keeps the error.
func ExampleTapErrorAs() {
	reportMissing := TapErrorAs[*missingKeyError, int](func(e *missingKeyError) (F.Void, error) {
		fmt.Println("please set", e.key)
		return Of(F.VOID)
	})

	fmt.Println(reportMissing(lookupPort("port")))
	// Output:
	// please set port
	// 0 config: missing key "port"
}
