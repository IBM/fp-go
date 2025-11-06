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

// ToAny converts a value of any type to the any (interface{}) type.
//
// This function performs an explicit type conversion to the any type, which can be
// useful when you need to store values of different types in a homogeneous collection
// or when interfacing with APIs that require any/interface{}.
//
// Type Parameters:
//   - A: The type of the input value
//
// Parameters:
//   - a: The value to convert
//
// Returns:
//   - The value as type any
//
// Example:
//
//	value := 42
//	anyValue := ToAny(value)  // any(42)
//
//	str := "hello"
//	anyStr := ToAny(str)  // any("hello")
//
//	// Useful for creating heterogeneous collections
//	values := []any{
//	    ToAny(42),
//	    ToAny("hello"),
//	    ToAny(true),
//	}
func ToAny[A any](a A) any {
	return any(a)
}
