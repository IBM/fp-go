// Copyright (c) 2024 - 2025 IBM Corp.
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

package io

import (
	"github.com/IBM/fp-go/v2/internal/functor"
)

type (
	ioFunctor[A, B any] struct{}

	// IOFunctor represents the functor type class for IO.
	// A functor allows mapping a function over a wrapped value without
	// unwrapping it, preserving the structure.
	IOFunctor[A, B any] = functor.Functor[A, B, IO[A], IO[B]]
)

func (o *ioFunctor[A, B]) Map(f func(A) B) Operator[A, B] {
	return Map(f)
}

// Functor returns an instance of the Functor type class for IO.
// This provides a structured way to access functor operations (Map)
// for IO computations.
//
// Example:
//
//	f := io.Functor[int, string]()
//	result := f.Map(strconv.Itoa)(io.Of(42))
func Functor[A, B any]() IOFunctor[A, B] {
	return &ioFunctor[A, B]{}
}
