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
	"github.com/IBM/fp-go/v2/internal/pointed"
)

type (
	ioPointed[A any] struct{}
	// IOPointed represents the pointed functor type class for IO.
	// A pointed functor is a functor with the ability to lift a pure value
	// into the functor context (Of operation).
	IOPointed[A any] = pointed.Pointed[A, IO[A]]
)

func (o *ioPointed[A]) Of(a A) IO[A] {
	return Of(a)
}

// Pointed returns an instance of the Pointed type class for IO.
// This provides a structured way to access the Of operation for IO computations.
//
// Example:
//
//	p := io.Pointed[int]()
//	result := p.Of(42)
func Pointed[A any]() IOPointed[A] {
	return &ioPointed[A]{}
}
