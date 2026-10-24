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

package prism

import (
	"github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/tailrec"
)

// FromLand creates a prism that focuses on the final result of a landed
// Trampoline.
//
// The prism's GetOption returns Some(result) if the Trampoline has landed and
// None if it bounces. Its ReverseGet always succeeds and wraps a value into a
// landed Trampoline.
//
// Type Parameters:
//   - B: The intermediate state type (bounce type)
//   - L: The final result type (land type)
//
// Returns:
//   - Prism[tailrec.Trampoline[B, L], L]: A prism that extracts landed results
//
// See Also:
//   - FromBounce: the prism for the bounce case
//   - tailrec.GetLand: the comma-ok getter used by GetOption
//   - tailrec.Land: the constructor used by ReverseGet
func FromLand[B, L any]() Prism[tailrec.Trampoline[B, L], L] {
	return MakePrismWithName(
		option.FromValidation(tailrec.GetLand[B, L]),
		tailrec.Land[B, L],
		"PrismFromLand",
	)
}

// FromBounce creates a prism that focuses on the intermediate state of a
// bouncing Trampoline.
//
// The prism's GetOption returns Some(state) if the Trampoline bounces and
// None if it has landed. Its ReverseGet always succeeds and wraps a value into
// a bouncing Trampoline.
//
// Type Parameters:
//   - B: The intermediate state type (bounce type)
//   - L: The final result type (land type)
//
// Returns:
//   - Prism[tailrec.Trampoline[B, L], B]: A prism that extracts bounce states
//
// See Also:
//   - FromLand: the prism for the land case
//   - tailrec.GetBounce: the comma-ok getter used by GetOption
//   - tailrec.Bounce: the constructor used by ReverseGet
func FromBounce[B, L any]() Prism[tailrec.Trampoline[B, L], B] {
	return MakePrismWithName(
		option.FromValidation(tailrec.GetBounce[B, L]),
		tailrec.Bounce[L, B],
		"PrismFromBounce",
	)
}
