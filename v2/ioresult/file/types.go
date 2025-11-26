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

package file

import "github.com/IBM/fp-go/v2/ioresult"

type (
	// IOResult represents a synchronous computation that may fail, returning a Result type.
	// It is an alias for ioresult.IOResult[T] which is equivalent to IO[Result[T]].
	//
	// IOResult[T] is a function that when executed returns Result[T], which is Either[error, T].
	// This provides a functional approach to handling IO operations that may fail, with
	// automatic resource management and composable error handling.
	//
	// Example:
	//   var readFile IOResult[[]byte] = func() Result[[]byte] {
	//       data, err := os.ReadFile("config.json")
	//       return result.TryCatchError(data, err)
	//   }
	//
	//   // Execute the IO operation
	//   result := readFile()
	//   data, err := E.UnwrapError(result)
	IOResult[T any] = ioresult.IOResult[T]

	// Kleisli represents a function that takes a value of type A and returns an IOResult[B].
	// It is an alias for ioresult.Kleisli[A, B] which is equivalent to Reader[A, IOResult[B]].
	//
	// Kleisli functions are the building blocks of monadic composition in the IOResult context.
	// They allow for chaining operations that may fail while maintaining functional purity.
	//
	// Example:
	//   // A Kleisli function that reads from a file handle
	//   var readAll Kleisli[*os.File, []byte] = func(f *os.File) IOResult[[]byte] {
	//       return TryCatchError(func() ([]byte, error) {
	//           return io.ReadAll(f)
	//       })
	//   }
	//
	//   // Can be composed with other Kleisli functions
	//   var processFile = F.Pipe1(
	//       file.Open("data.txt"),
	//       file.Read[[]byte, *os.File],
	//   )(readAll)
	Kleisli[A, B any] = ioresult.Kleisli[A, B]

	// Operator represents a function that transforms one IOResult into another.
	// It is an alias for ioresult.Operator[A, B] which is equivalent to Kleisli[IOResult[A], B].
	//
	// Operators are used for transforming and composing IOResult values, providing a way
	// to build complex data processing pipelines while maintaining error handling semantics.
	//
	// Example:
	//   // An operator that converts bytes to string
	//   var bytesToString Operator[[]byte, string] = Map(func(data []byte) string {
	//       return string(data)
	//   })
	//
	//   // An operator that validates JSON
	//   var validateJSON Operator[string, map[string]interface{}] = ChainEitherK(
	//       func(s string) Result[map[string]interface{}] {
	//           var result map[string]interface{}
	//           err := json.Unmarshal([]byte(s), &result)
	//           return result.TryCatchError(result, err)
	//       },
	//   )
	//
	//   // Compose operators in a pipeline
	//   var processJSON = F.Pipe2(
	//       readFileOperation,
	//       bytesToString,
	//       validateJSON,
	//   )
	Operator[A, B any] = ioresult.Operator[A, B]
)
