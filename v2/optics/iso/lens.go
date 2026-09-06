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

import "github.com/IBM/fp-go/v2/internal/common"

// AsLens converts an Iso[S, A] into a Lens[S, A].
//
// An Iso is a strictly stronger optic than a Lens: every Iso is a valid Lens
// (Get and ReverseGet map directly onto Get and Set), but not every Lens is an
// Iso.  This function performs that downcast.
//
// Type Parameters:
//   - S: The source type
//   - A: The target type
//
// Parameters:
//   - sa: The isomorphism to convert
//
// Returns:
//   - A Lens[S, A] whose Get is sa.Get and whose Set is derived from sa.ReverseGet
//
// See Also:
//   - ComposePrism: compose an Iso with a Prism
//   - Compose: compose two Isos
func AsLens[S, A any](sa Iso[S, A]) common.Lens[S, A] {
	return common.IsoAsLens(sa)
}
