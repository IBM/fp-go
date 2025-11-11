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
	C "github.com/IBM/fp-go/v2/constant"
	M "github.com/IBM/fp-go/v2/monoid"
	G "github.com/IBM/fp-go/v2/optics/traversal/generic"
	RR "github.com/IBM/fp-go/v2/optics/traversal/record/generic"
)

// FromRecord returns a traversal from an array for the const monad
func FromRecord[MA ~map[K]A, E, K comparable, A any](m M.Monoid[E]) G.Traversal[MA, A, C.Const[E, MA], C.Const[E, A]] {
	return RR.FromRecord[MA](
		C.Of[E, MA](m),
		C.Map[E, MA, func(A) MA],
		C.Ap[E, A, MA](m),
	)
}
