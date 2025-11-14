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

// Scan takes an [Iterator] and returns a new [Iterator] of the same length, where the values
// of the new [Iterator] are the result of the application of `f` to the value of the
// source iterator with the previously accumulated value
func Scan[FCT ~func(V, U) V, U, V any](f FCT, initial V) Operator[U, V] {
	return G.Scan[Iterator[V], Iterator[U]](f, initial)
}
