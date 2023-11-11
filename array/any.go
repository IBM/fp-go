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

package array

import (
	G "github.com/IBM/fp-go/array/generic"
)

// AnyWithIndex tests if any of the elements in the array matches the predicate
func AnyWithIndex[A any](pred func(int, A) bool) func([]A) bool {
	return G.AnyWithIndex[[]A](pred)
}

// Any tests if any of the elements in the array matches the predicate
func Any[A any](pred func(A) bool) func([]A) bool {
	return G.Any[[]A](pred)
}
