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

package record

import (
	L "github.com/IBM/fp-go/v2/optics/lens"
	G "github.com/IBM/fp-go/v2/optics/lens/record/generic"
	O "github.com/IBM/fp-go/v2/option"
)

// AtRecord returns a lens that focusses on a value in a record
func AtRecord[V any, K comparable](key K) L.Lens[map[K]V, O.Option[V]] {
	return G.AtRecord[map[K]V](key)
}

// AtKey returns a `Lens` focused on a required key of a `ReadonlyRecord`
func AtKey[S any, V any, K comparable](key K) func(sa L.Lens[S, map[K]V]) L.Lens[S, O.Option[V]] {
	return G.AtKey[map[K]V, S](key)
}
