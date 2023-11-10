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

package generic

import (
	F "github.com/IBM/fp-go/function"
	IG "github.com/IBM/fp-go/identity/generic"
	O "github.com/IBM/fp-go/option"
)

func Switch[M1 ~map[K]V, M2 ~map[K]R, N ~map[K]FCT, FCT ~func(V) R, K comparable, V, R any](n N, d FCT) func(M1) M2 {
	return MapWithIndex[M1, M2](func(idx K, val V) R {
		return F.Pipe3(
			n,
			Lookup[N](idx),
			O.GetOrElse(F.Constant(d)),
			IG.Flap[FCT, R](val),
		)
	})
}
