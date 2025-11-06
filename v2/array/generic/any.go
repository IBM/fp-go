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
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
)

// AnyWithIndex tests if any of the elements in the array matches the predicate
func AnyWithIndex[AS ~[]A, PRED ~func(int, A) bool, A any](pred PRED) func(AS) bool {
	return F.Flow2(
		FindFirstWithIndex[AS](pred),
		O.IsSome[A],
	)
}

// Any tests if any of the elements in the array matches the predicate
func Any[AS ~[]A, PRED ~func(A) bool, A any](pred PRED) func(AS) bool {
	return AnyWithIndex[AS](F.Ignore1of2[int](pred))
}
