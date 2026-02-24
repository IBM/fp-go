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

import (
	"github.com/IBM/fp-go/v2/context/readerioresult"
	"github.com/IBM/fp-go/v2/function"
)

type (
	// ReaderIOResult represents a context-aware computation that performs side effects
	// and can fail with an error. This is the main type used throughout the file package
	// for all file operations.
	//
	// ReaderIOResult[A] is equivalent to:
	//   func(context.Context) func() Either[error, A]
	//
	// The computation:
	//   - Takes a context.Context for cancellation and timeouts
	//   - Performs side effects (IO operations)
	//   - Returns Either an error or a value of type A
	//
	// See Also:
	//   - readerioresult.ReaderIOResult: The underlying type definition
	ReaderIOResult[A any] = readerioresult.ReaderIOResult[A]

	// Void represents the absence of a meaningful value, similar to unit type in other languages.
	// It is used when a function performs side effects but doesn't return a meaningful result.
	//
	// Void is typically used as the success type in operations like Close that perform
	// an action but don't produce a useful value.
	//
	// Example:
	//   Close[*os.File](file) // Returns ReaderIOResult[Void]
	//
	// See Also:
	//   - function.Void: The underlying type definition
	Void = function.Void

	// Kleisli represents a Kleisli arrow for ReaderIOResult.
	// It is a function that takes a value of type A and returns a ReaderIOResult[B].
	//
	// Kleisli arrows are used for monadic composition, allowing you to chain operations
	// that produce ReaderIOResults. They are particularly useful with Chain and Bind operations.
	//
	// Kleisli[A, B] is equivalent to:
	//   func(A) ReaderIOResult[B]
	//
	// Example:
	//   // A Kleisli arrow that reads a file given its path
	//   var readFileK Kleisli[string, []byte] = ReadFile
	//
	// See Also:
	//   - readerioresult.Kleisli: The underlying type definition
	//   - Operator: For transforming ReaderIOResults
	Kleisli[A, B any] = readerioresult.Kleisli[A, B]

	// Operator represents a transformation from one ReaderIOResult to another.
	// This is useful for point-free style composition and building reusable transformations.
	//
	// Operator[A, B] is equivalent to:
	//   func(ReaderIOResult[A]) ReaderIOResult[B]
	//
	// Operators are used to transform computations without executing them, enabling
	// powerful composition patterns.
	//
	// Example:
	//   // An operator that maps over file contents
	//   var toUpper Operator[[]byte, string] = Map(func(data []byte) string {
	//       return strings.ToUpper(string(data))
	//   })
	//
	// See Also:
	//   - readerioresult.Operator: The underlying type definition
	//   - Kleisli: For functions that produce ReaderIOResults
	Operator[A, B any] = readerioresult.Operator[A, B]
)
