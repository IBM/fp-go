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

package generic

import (
	"fmt"

	F "github.com/IBM/fp-go/v2/function"
	I "github.com/IBM/fp-go/v2/identity"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
	RR "github.com/IBM/fp-go/v2/record/generic"
)

// AtRecord returns a lens that focusses on a value in a record
func AtRecord[M ~map[K]V, V any, K comparable](key K) L.Lens[M, O.Option[V]] {
	addKey := F.Bind1of2(RR.UpsertAt[M, K, V])(key)
	delKey := F.Bind1of1(RR.DeleteAt[M, K, V])(key)
	fold := O.Fold(
		delKey,
		addKey,
	)
	return L.MakeLensWithName(
		RR.Lookup[M](key),
		func(m M, v O.Option[V]) M {
			return F.Pipe2(
				v,
				fold,
				I.Ap[M](m),
			)
		},
		fmt.Sprintf("At[%v]", key),
	)
}

// AtKey returns a `Lens` focused on a required key of a `ReadonlyRecord`
func AtKey[M ~map[K]V, S any, V any, K comparable](key K) func(sa L.Lens[S, M]) L.Lens[S, O.Option[V]] {
	return L.Compose[S](AtRecord[M](key))
}
