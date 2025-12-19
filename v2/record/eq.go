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
	E "github.com/IBM/fp-go/v2/eq"
	G "github.com/IBM/fp-go/v2/record/generic"
)

func Eq[K comparable, V any](e E.Eq[V]) E.Eq[Record[K, V]] {
	return G.Eq[Record[K, V]](e)
}

// FromStrictEquals constructs an [EQ.Eq] from the canonical comparison function
func FromStrictEquals[K, V comparable]() E.Eq[Record[K, V]] {
	return G.FromStrictEquals[Record[K, V]]()
}
