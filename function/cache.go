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

import (
	G "github.com/IBM/fp-go/function/generic"
)

// Memoize converts a unary function into a unary function that caches the value depending on the parameter
func Memoize[K comparable, T any](f func(K) T) func(K) T {
	return G.Memoize(f)
}

// ContramapMemoize converts a unary function into a unary function that caches the value depending on the parameter
func ContramapMemoize[A any, K comparable, T any](kf func(A) K) func(func(A) T) func(A) T {
	return G.ContramapMemoize[func(A) T](kf)
}
