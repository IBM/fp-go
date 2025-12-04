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

// Package mostlyadequate contains examples from the "Mostly Adequate Guide to Functional Programming"
// adapted to Go using fp-go. These examples demonstrate functional programming concepts in a practical way.
package mostlyadequate

import (
	"github.com/IBM/fp-go/v2/iooption"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/result"
)

// Type aliases for common monads used throughout the examples.
// These aliases simplify type signatures and make the code more readable.
// This pattern is recommended in fp-go v2 - define aliases once and use them throughout your codebase.
type (
	// Result represents a computation that may fail with an error.
	// It's an alias for result.Result[A] which is equivalent to either.Either[error, A].
	// Use this when you need error handling with a specific success type.
	//
	// Example:
	//   func divide(a, b int) Result[int] {
	//       if b == 0 {
	//           return result.Error[int](errors.New("division by zero"))
	//       }
	//       return result.Ok(a / b)
	//   }
	Result[A any] = result.Result[A]

	// IOOption represents a lazy computation that may not produce a value.
	// It combines IO (lazy evaluation) with Option (optional values).
	// Use this when you have side effects that might not return a value.
	//
	// Example:
	//   func readConfig() IOOption[Config] {
	//       return func() option.Option[Config] {
	//           // Read from file system (side effect)
	//           if fileExists {
	//               return option.Some(config)
	//           }
	//           return option.None[Config]()
	//       }
	//   }
	IOOption[A any] = iooption.IOOption[A]

	// Option represents an optional value - either Some(value) or None.
	// Use this instead of pointers or sentinel values to represent absence of a value.
	//
	// Example:
	//   func findUser(id int) Option[User] {
	//       if user, found := users[id]; found {
	//           return option.Some(user)
	//       }
	//       return option.None[User]()
	//   }
	Option[A any] = option.Option[A]

	// IOResult represents a lazy computation that may fail with an error.
	// It combines IO (lazy evaluation) with Result (error handling).
	// Use this for side effects that can fail, like file I/O or HTTP requests.
	//
	// Example:
	//   func readFile(path string) IOResult[[]byte] {
	//       return func() result.Result[[]byte] {
	//           data, err := os.ReadFile(path)
	//           if err != nil {
	//               return result.Error[[]byte](err)
	//           }
	//           return result.Ok(data)
	//       }
	//   }
	IOResult[A any] = ioresult.IOResult[A]
)
