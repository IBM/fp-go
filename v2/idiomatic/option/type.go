// Copyright (c) 2025 IBM Corp.
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

package option

func toType[T any](a any) (T, bool) {
	b, ok := a.(T)
	return b, ok
}

// ToType attempts to convert a value of type any to a specific type T using type assertion.
// Returns Some(value) if the type assertion succeeds, None if it fails.
//
// Example:
//
//	var x any = 42
//	result := ToType[int](x) // Some(42)
//
//	var y any = "hello"
//	result := ToType[int](y) // None (wrong type)
//
//go:inline
func ToType[T any](src any) (T, bool) {
	return toType[T](src)
}

// ToAny converts a value of any type to Option[any].
// This always succeeds and returns Some containing the value as any.
//
// Example:
//
//	result := ToAny(42) // Some(any(42))
//	result := ToAny("hello") // Some(any("hello"))
//
//go:inline
func ToAny[T any](src T) (any, bool) {
	return Of(any(src))
}
