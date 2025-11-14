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
	R "reflect"

	F "github.com/IBM/fp-go/v2/function"
	LG "github.com/IBM/fp-go/v2/io/generic"
	L "github.com/IBM/fp-go/v2/lazy"
	N "github.com/IBM/fp-go/v2/number"
	I "github.com/IBM/fp-go/v2/number/integer"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/pair"
)

func FromReflect[GR ~func() Option[Pair[GR, R.Value]]](val R.Value) GR {
	// recursive callback
	var recurse func(idx int) GR

	// limits the index
	fromPred := O.FromPredicate(I.Between(0, val.Len()))

	recurse = func(idx int) GR {
		return F.Pipe3(
			idx,
			L.Of[int],
			L.Map(fromPred),
			LG.Map[Lazy[Option[int]], GR](O.Map(
				F.Flow2(
					P.Of[int],
					P.BiMap(F.Flow2(N.Add(1), recurse), val.Index),
				),
			)),
		)
	}

	// start the recursion
	return recurse(0)
}
