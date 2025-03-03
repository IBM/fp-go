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

package ioeither

import (
	"github.com/IBM/fp-go/v2/semigroup"
)

type (
	Semigroup[E, A any] = semigroup.Semigroup[IOEither[E, A]]
)

// AltSemigroup is a [Semigroup] that tries the first item and then the second one using an alternative
func AltSemigroup[E, A any]() semigroup.Semigroup[IOEither[E, A]] {
	return semigroup.AltSemigroup(
		MonadAlt[E, A],
	)
}
