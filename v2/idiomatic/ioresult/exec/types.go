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

package exec

import (
	"github.com/IBM/fp-go/v2/idiomatic/ioresult"
)

type (
	// IOResult represents an IO operation that may fail with an error.
	// It is defined as func() (T, error), following Go's idiomatic error handling pattern.
	//
	// This type is re-exported from the ioresult package for convenience when working
	// with command execution, allowing users to reference exec.IOResult instead of
	// importing the ioresult package separately.
	IOResult[T any] = ioresult.IOResult[T]

	// Kleisli represents a function from A to IOResult[B].
	// Named after Heinrich Kleisli, it represents a monadic arrow in the IOResult monad.
	//
	// Kleisli functions are useful for chaining operations where each step may perform
	// IO and may fail. They can be composed using IOResult's Chain function.
	//
	// Example:
	//
	//	type Kleisli[A, B any] func(A) IOResult[B]
	//
	//	parseConfig := func(path string) IOResult[Config] { ... }
	//	validateConfig := func(cfg Config) IOResult[Config] { ... }
	//
	//	// Compose Kleisli functions
	//	loadAndValidate := Chain(validateConfig)(parseConfig("/config.yml"))
	Kleisli[A, B any] = ioresult.Kleisli[A, B]

	// Operator represents a function that transforms one IOResult into another.
	// It maps IOResult[A] to IOResult[B], useful for defining reusable transformations.
	//
	// Example:
	//
	//	type Operator[A, B any] func(IOResult[A]) IOResult[B]
	//
	//	addLogging := func(io IOResult[string]) IOResult[string] {
	//	    return func() (string, error) {
	//	        result, err := io()
	//	        log.Printf("Result: %v, Error: %v", result, err)
	//	        return result, err
	//	    }
	//	}
	Operator[A, B any] = ioresult.Operator[A, B]
)
