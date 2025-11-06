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

package examples

type HKT[T any] struct {
}

// Pointed
func Of[A any](A) HKT[A] { return HKT[A]{} }

// Functor
func Map[A, B any](func(A) B) func(HKT[A]) HKT[B] { return func(HKT[A]) HKT[B] { return HKT[B]{} } }
func MapTo[A, B any](A) func(HKT[A]) HKT[B]       { return func(HKT[A]) HKT[B] { return HKT[B]{} } }

// Chain
func Chain[A, B any](func(A) HKT[B]) func(HKT[A]) HKT[B] {
	return func(HKT[A]) HKT[B] { return HKT[B]{} }
}
func ChainTo[A, B any](HKT[B]) func(HKT[A]) HKT[B] {
	return func(HKT[A]) HKT[B] { return HKT[B]{} }
}
func ChainFirst[A, B any](func(A) HKT[B]) func(HKT[A]) HKT[A] {
	return func(HKT[A]) HKT[A] { return HKT[A]{} }
}

// Apply
func Ap[A, B any](HKT[A]) func(HKT[func(A) B]) HKT[B] {
	return func(HKT[func(A) B]) HKT[B] { return HKT[B]{} }
}

// Derived
func Flatten[A, B any](HKT[HKT[A]]) HKT[A] {
	return HKT[A]{}
}
func Reduce[A, B any](func(B, A) B, B) func(HKT[A]) HKT[B] {
	return func(HKT[A]) HKT[B] { return HKT[B]{} }
}
