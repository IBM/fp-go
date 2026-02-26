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

package monoid

import "github.com/IBM/fp-go/v2/function"

// Void is an alias for function.Void, representing the unit type.
//
// The Void type (also known as Unit in functional programming) has exactly one value,
// making it useful for representing the absence of meaningful information. It's similar
// to void in other languages, but as a value rather than the absence of a value.
//
// This type alias is provided in the monoid package for convenience when working with
// VoidMonoid and other monoid operations that may use the unit type.
//
// Common use cases:
//   - As a return type for functions that perform side effects but don't return meaningful data
//   - As a placeholder type parameter when a type is required but no data needs to be passed
//   - In monoid operations where you need to track that operations occurred without caring about results
//
// See also:
//   - function.Void: The underlying type definition
//   - function.VOID: The single inhabitant of the Void type
//   - VoidMonoid: A monoid instance for the Void type
//
// Example:
//
//	// Function that performs an action but returns no meaningful data
//	func logMessage(msg string) Void {
//	    fmt.Println(msg)
//	    return function.VOID
//	}
//
//	// Using Void in monoid operations
//	m := VoidMonoid()
//	result := m.Concat(function.VOID, function.VOID)  // function.VOID
type (
	Void = function.Void
)
