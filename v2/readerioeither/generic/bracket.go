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
	"github.com/IBM/fp-go/v2/either"
	G "github.com/IBM/fp-go/v2/internal/bracket"
	I "github.com/IBM/fp-go/v2/readerio/generic"
)

// Bracket makes sure that a resource is cleaned up in the event of an error. The release action is called regardless of
// whether the body action returns and error or not.
func Bracket[
	GA ~func(R) TA,
	GB ~func(R) TB,
	GANY ~func(R) TANY,

	TA ~func() either.Either[E, A],
	TB ~func() either.Either[E, B],
	TANY ~func() either.Either[E, ANY],

	R, E, A, B, ANY any](

	acquire GA,
	use func(A) GB,
	release func(A, either.Either[E, B]) GANY,
) GB {
	return G.MonadBracket[GA, GB, GANY, either.Either[E, B], A, B](
		I.Of[GB, TB, R, either.Either[E, B]],
		MonadChain[GA, GB, TA, TB, R, E, A, B],
		I.MonadChain[GB, GB, TB, TB, R, either.Either[E, B], either.Either[E, B]],
		MonadChain[GANY, GB, TANY, TB, R, E, ANY, B],

		acquire,
		use,
		release,
	)
}
