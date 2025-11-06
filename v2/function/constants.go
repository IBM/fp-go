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

package function

// ConstTrue is a nullary function that always returns true.
//
// This is a pre-defined constant function created using Constant(true).
// It's useful as a default predicate or when you need a function that
// always evaluates to true.
//
// Example:
//
//	result := ConstTrue()  // true
//
//	// Use as a default predicate
//	filter := func(pred func() bool) string {
//	    if pred() {
//	        return "accepted"
//	    }
//	    return "rejected"
//	}
//	result := filter(ConstTrue)  // "accepted"
var ConstTrue = Constant(true)

// ConstFalse is a nullary function that always returns false.
//
// This is a pre-defined constant function created using Constant(false).
// It's useful as a default predicate or when you need a function that
// always evaluates to false.
//
// Example:
//
//	result := ConstFalse()  // false
//
//	// Use as a default predicate
//	filter := func(pred func() bool) string {
//	    if pred() {
//	        return "accepted"
//	    }
//	    return "rejected"
//	}
//	result := filter(ConstFalse)  // "rejected"
var ConstFalse = Constant(false)

// ConstNil returns a nil pointer of the specified type.
//
// This function creates a nil pointer for any type A. It's useful when you need
// to provide a nil value in a generic context or when initializing optional fields.
//
// Type Parameters:
//   - A: The type for which to create a nil pointer
//
// Returns:
//   - A nil pointer of type *A
//
// Example:
//
//	nilInt := ConstNil[int]()      // (*int)(nil)
//	nilString := ConstNil[string]() // (*string)(nil)
//
//	// Useful for optional fields
//	type Config struct {
//	    Timeout *int
//	    MaxRetries *int
//	}
//	config := Config{
//	    Timeout: ConstNil[int](),
//	    MaxRetries: Ref(3),
//	}
func ConstNil[A any]() *A {
	return (*A)(nil)
}
