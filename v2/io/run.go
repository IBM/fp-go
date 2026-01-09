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

package io

// Run executes an IO computation and returns its result.
//
// This function is the primary way to execute IO computations. It takes an IO[A]
// (a lazy computation) and immediately evaluates it, returning the computed value.
//
// Run is the bridge between the pure functional world (where computations are
// described but not executed) and the imperative world (where side effects occur).
// It should typically be called at the edges of your application, such as in main()
// or in test code.
//
// Parameters:
//   - fa: The IO computation to execute
//
// Returns:
//   - The result of executing the IO computation
//
// Example:
//
//	// Create a lazy computation
//	greeting := io.Of("Hello, World!")
//
//	// Execute it and get the result
//	result := io.Run(greeting) // result == "Hello, World!"
//
// Example with side effects:
//
//	// Create a computation that prints and returns a value
//	computation := func() string {
//	    fmt.Println("Computing...")
//	    return "Done"
//	}
//
//	// Nothing is printed yet
//	io := io.MakeIO(computation)
//
//	// Now the computation runs and "Computing..." is printed
//	result := io.Run(io) // result == "Done"
//
// Example with composition:
//
//	result := io.Run(
//	    pipe.Pipe2(
//	        io.Of(5),
//	        io.Map(func(x int) int { return x * 2 }),
//	        io.Map(func(x int) int { return x + 1 }),
//	    ),
//	) // result == 11
//
// Note: Run should be used sparingly in application code. Prefer composing
// IO computations and only calling Run at the application boundaries.
func Run[A any](fa IO[A]) A {
	return fa()
}
