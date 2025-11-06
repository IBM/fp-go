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

package number

import (
	M "github.com/IBM/fp-go/v2/monoid"
)

// MonoidSum is the [Monoid] that adds elements with a zero empty element
func MonoidSum[A Number]() M.Monoid[A] {
	s := SemigroupSum[A]()
	return M.MakeMonoid(
		s.Concat,
		0,
	)
}

// MonoidProduct is the [Monoid] that multiplies elements with a one empty element
func MonoidProduct[A Number]() M.Monoid[A] {
	s := SemigroupProduct[A]()
	return M.MakeMonoid(
		s.Concat,
		1,
	)
}
