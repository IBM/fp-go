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

package iooption

import (
	G "github.com/IBM/fp-go/iooption/generic"
	O "github.com/IBM/fp-go/option"
)

// Bracket makes sure that a resource is cleaned up in the event of an error. The release action is called regardless of
// whether the body action returns and error or not.
func Bracket[E, A, B, ANY any](
	acquire IOOption[A],
	use func(A) IOOption[B],
	release func(A, O.Option[B]) IOOption[ANY],
) IOOption[B] {
	return G.Bracket(acquire, use, release)
}
