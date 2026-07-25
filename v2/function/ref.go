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
// Useful for obtaining a pointer to a literal or a local variable, particularly
// when working with APIs that require pointer types.
//
// Type Parameters:
//   - A: the type of the value
//
// Parameters:
//   - a: the value to create a pointer to
//
// Returns:
//   - a non-nil pointer to a copy of a
func Ref[A any](a A) *A {
	return &a
}

// Deref dereferences a pointer and returns the value it points to.
//
// Panics if a is nil. Use DerefSafe when the pointer may be nil.
//
// Type Parameters:
//   - A: the type of the value
//
// Parameters:
//   - a: the non-nil pointer to dereference
//
// Returns:
//   - the value pointed to by a
func Deref[A any](a *A) A {
	return *a
}

// DerefSafe dereferences a pointer and returns the value it points to,
// or the result of the zero thunk if the pointer is nil.
//
// The zero thunk is only called when the pointer is nil, so it is safe to use
// with expensive or side-effectful fallback computations. Use Constant to wrap
// a plain value as a thunk: DerefSafe(Constant(0)).
//
// Type Parameters:
//   - A: the type of the value
//
// Parameters:
//   - zero: a thunk evaluated to produce the fallback value when the pointer is nil
//
// Returns:
//   - a function from *A to A that returns *a when a is non-nil, or zero() otherwise
func DerefSafe[A any](zero func() A) func(*A) A {
	return func(a *A) A {
		if a == nil {
			return zero()
		}
		return *a
	}
}
