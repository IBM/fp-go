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

package ord

import (
	O "github.com/IBM/fp-go/ord"
	S "github.com/IBM/fp-go/semigroup"
)

// Max gets a semigroup where `concat` will return the maximum, based on the provided order.
func Max[A any](o O.Ord[A]) S.Semigroup[A] {
	return S.MakeSemigroup(O.Max(o))
}

// Min gets a semigroup where `concat` will return the minimum, based on the provided order.
func Min[A any](o O.Ord[A]) S.Semigroup[A] {
	return S.MakeSemigroup(O.Min(o))
}
