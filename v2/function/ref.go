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

// Ref creates a pointer to a value.
//
// This function takes a value and returns a pointer to it. It's useful when you need
// to convert a value to a pointer, particularly when working with APIs that require
// pointers or when you need to create optional values.
//
// Type Parameters:
//   - A: The type of the value
//
// Parameters:
//   - a: The value to create a pointer to
//
// Returns:
//   - A pointer to the value
//
// Example:
//
//	value := 42
//	ptr := Ref(value)
//	fmt.Println(*ptr)  // 42
//
//	// Useful for creating pointers to literals
//	strPtr := Ref("hello")
//	fmt.Println(*strPtr)  // "hello"
//
//	// Creating optional values
//	type Config struct {
//	    Timeout *int
//	}
//	config := Config{Timeout: Ref(30)}
func Ref[A any](a A) *A {
	return &a
}

// Deref dereferences a pointer to get its value.
//
// This function takes a pointer and returns the value it points to. It will panic
// if the pointer is nil, so it should only be used when you're certain the pointer
// is not nil. For safe dereferencing, check with IsNonNil first.
//
// Type Parameters:
//   - A: The type of the value
//
// Parameters:
//   - a: The pointer to dereference
//
// Returns:
//   - The value pointed to by the pointer
//
// Example:
//
//	value := 42
//	ptr := &value
//	result := Deref(ptr)  // 42
//
//	// Safe usage with nil check
//	var ptr *int
//	if IsNonNil(ptr) {
//	    result := Deref(ptr)
//	    fmt.Println(result)
//	}
//
//	// Chaining with Ref
//	original := "hello"
//	copy := Deref(Ref(original))  // "hello"
func Deref[A any](a *A) A {
	return *a
}
