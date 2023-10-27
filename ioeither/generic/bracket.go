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

package generic

import (
	ET "github.com/IBM/fp-go/either"
	G "github.com/IBM/fp-go/internal/bracket"
	I "github.com/IBM/fp-go/io/generic"
)

// Bracket makes sure that a resource is cleaned up in the event of an error. The release action is called regardless of
// whether the body action returns and error or not.
func Bracket[
	GA ~func() ET.Either[E, A],
	GB ~func() ET.Either[E, B],
	GANY ~func() ET.Either[E, ANY],
	E, A, B, ANY any](

	acquire GA,
	use func(A) GB,
	release func(A, ET.Either[E, B]) GANY,
) GB {
	return G.Bracket[GA, GB, GANY, ET.Either[E, B], A, B](
		I.Of[GB, ET.Either[E, B]],
		MonadChain[GA, GB, E, A, B],
		I.MonadChain[GB, GB, ET.Either[E, B], ET.Either[E, B]],
		MonadChain[GANY, GB, E, ANY, B],

		acquire,
		use,
		release,
	)
}
