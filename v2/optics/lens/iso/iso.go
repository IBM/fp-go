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

package iso

import (
	F "github.com/IBM/fp-go/v2/function"
	I "github.com/IBM/fp-go/v2/optics/iso"
	IL "github.com/IBM/fp-go/v2/optics/iso/lens"
	L "github.com/IBM/fp-go/v2/optics/lens"
	O "github.com/IBM/fp-go/v2/option"
)

// FromNillable converts a nillable value to an option and back
func FromNillable[T any]() Iso[*T, Option[T]] {
	return I.MakeIso(F.Flow2(
		O.FromPredicate(F.IsNonNil[T]),
		O.Map(F.Deref[T]),
	),
		O.Fold(F.ConstNil[T], F.Ref[T]),
	)
}

// Compose converts a Lens to a property of `A` into a lens to a property of type `B`
// the transformation is done via an ISO
//
//go:inline
func Compose[S, A, B any](ab Iso[A, B]) Operator[S, A, B] {
	return F.Pipe2(
		ab,
		IL.IsoAsLens[A, B],
		L.Compose[S, A, B],
	)
}
