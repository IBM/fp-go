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
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/pair"
)

func Cycle[GU ~func() O.Option[P.Pair[GU, U]], U any](ma GU) GU {
	// avoid cyclic references
	var m func(O.Option[P.Pair[GU, U]]) O.Option[P.Pair[GU, U]]

	recurse := func(mu GU) GU {
		return F.Nullary2(
			mu,
			m,
		)
	}

	m = O.Fold(func() O.Option[P.Pair[GU, U]] {
		return recurse(ma)()
	}, F.Flow2(
		P.BiMap(recurse, F.Identity[U]),
		O.Of[P.Pair[GU, U]],
	))

	return recurse(ma)
}
