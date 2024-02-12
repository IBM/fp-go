// Copyright (c) 2024 IBM Corp.
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

package option

import (
	P "github.com/IBM/fp-go/pair"
	PG "github.com/IBM/fp-go/pair/generic"
)

// SequencePair converts a [Pair] of [Option[T]] into an [Option[Pair]].
func SequencePair[T1, T2 any](t P.Pair[Option[T1], Option[T2]]) Option[P.Pair[T1, T2]] {
	return PG.SequencePair(
		Map[T1, func(T2) P.Pair[T1, T2]],
		Ap[P.Pair[T1, T2], T2],
		t,
	)
}
