//   Copyright (c) 2023 IBM Corp.
//   All rights reserved.
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package readerioresult

import (
	"github.com/IBM/fp-go/v2/semigroup"
)

type (
	Semigroup[A any] = semigroup.Semigroup[ReaderIOResult[A]]
)

// AltSemigroup is a [Semigroup] that tries the first item and then the second one using an alternative.
// This creates a semigroup where combining two ReaderIOResult values means trying the first one,
// and if it fails, trying the second one. This is useful for implementing fallback behavior.
//
// Returns a Semigroup for ReaderIOResult[A] with Alt-based combination.
//
// Example:
//
//	sg := AltSemigroup[int]()
//	result := sg.Concat(
//	    Left[int](errors.New("first failed")),
//	    Right[int](42),
//	) // Returns Right(42)
func AltSemigroup[A any]() Semigroup[A] {
	return semigroup.AltSemigroup(
		MonadAlt[A],
	)
}
