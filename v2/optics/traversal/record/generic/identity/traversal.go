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
	I "github.com/IBM/fp-go/v2/identity"
	G "github.com/IBM/fp-go/v2/optics/traversal/generic"
	RR "github.com/IBM/fp-go/v2/optics/traversal/record/generic"
)

// FromRecord returns a traversal from a record for the identity monad
func FromRecord[MA ~map[K]A, K comparable, A any]() G.Traversal[MA, A, MA, A] {
	return RR.FromRecord[MA](
		I.Of[MA],
		I.Map[MA, func(A) MA],
		I.Ap[MA, A],
	)
}
