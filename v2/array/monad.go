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

package array

import (
	G "github.com/IBM/fp-go/v2/array/generic"
	"github.com/IBM/fp-go/v2/internal/monad"
)

// Monad returns the monadic operations for an array.
// This provides a structured way to access all monad operations (Map, Chain, Ap, Of)
// for arrays in a single interface.
//
// The Monad interface is useful when you need to pass monadic operations as parameters
// or when working with generic code that operates on any monad.
//
// Example:
//
//	m := array.Monad[int, string]()
//	result := m.Chain([]int{1, 2, 3}, func(x int) []string {
//	    return []string{fmt.Sprintf("%d", x), fmt.Sprintf("%d!", x)}
//	})
//	// Result: ["1", "1!", "2", "2!", "3", "3!"]
//
//go:inline
func Monad[A, B any]() monad.Monad[A, B, []A, []B, []func(A) B] {
	return G.Monad[A, B, []A, []B, []func(A) B]()
}
