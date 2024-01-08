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
	L "github.com/IBM/fp-go/io/generic"
	O "github.com/IBM/fp-go/option"
	T "github.com/IBM/fp-go/tuple"
)

// FromLazy returns an iterator on top of a lazy function
func FromLazy[GU ~func() O.Option[T.Tuple2[GU, U]], LZ ~func() U, U any](l LZ) GU {
	return F.Pipe1(
		l,
		L.Map[LZ, GU](F.Flow2(
			F.Bind1st(T.MakeTuple2[GU, U], Empty[GU]()),
			O.Of[T.Tuple2[GU, U]],
		)),
	)
}
