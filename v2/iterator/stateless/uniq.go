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

package stateless

import (
	G "github.com/IBM/fp-go/v2/iterator/stateless/generic"
)

// StrictUniq converts an [Iterator] of arbitrary items into an [Iterator] or unique items
// where uniqueness is determined by the built-in uniqueness constraint
func StrictUniq[A comparable](as Iterator[A]) Iterator[A] {
	return G.StrictUniq(as)
}

// Uniq converts an [Iterator] of arbitrary items into an [Iterator] or unique items
// where uniqueness is determined based on a key extractor function
func Uniq[A any, K comparable](f func(A) K) Operator[A, A] {
	return G.Uniq[Iterator[A]](f)
}
