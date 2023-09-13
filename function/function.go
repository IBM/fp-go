// Copyright (c) 2023 IBM Corp.
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

// Identity returns the value 'a'
func Identity[A any](a A) A {
	return a
}

// Constant creates a nullary function that returns the constant value 'a'
func Constant[A any](a A) func() A {
	return func() A {
		return a
	}
}

// Constant1 creates a unary function that returns the constant value 'a' and ignores its input
func Constant1[B, A any](a A) func(B) A {
	return func(_ B) A {
		return a
	}
}

// Constant2 creates a binary function that returns the constant value 'a' and ignores its inputs
func Constant2[B, C, A any](a A) func(B, C) A {
	return func(_ B, _ C) A {
		return a
	}
}

func IsNil[A any](a *A) bool {
	return a == nil
}

func IsNonNil[A any](a *A) bool {
	return a != nil
}

// Swap returns a new binary function that changes the order of input parameters
func Swap[T1, T2, R any](f func(T1, T2) R) func(T2, T1) R {
	return func(t2 T2, t1 T1) R {
		return f(t1, t2)
	}
}

// First returns the first out of two input values
func First[T1, T2 any](t1 T1, _ T2) T1 {
	return t1
}

// Second returns the second out of two input values
// Identical to [SK]
func Second[T1, T2 any](_ T1, t2 T2) T2 {
	return t2
}
