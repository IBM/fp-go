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
	"context"

	E "github.com/IBM/fp-go/either"
	G "github.com/IBM/fp-go/internal/bracket"
	I "github.com/IBM/fp-go/readerio/generic"
)

// Bracket makes sure that a resource is cleaned up in the event of an error. The release action is called regardless of
// whether the body action returns and error or not.
func Bracket[
	GA ~func(context.Context) TA,
	GB ~func(context.Context) TB,
	GANY ~func(context.Context) TANY,

	TA ~func() E.Either[error, A],
	TB ~func() E.Either[error, B],
	TANY ~func() E.Either[error, ANY],

	A, B, ANY any](

	acquire GA,
	use func(A) GB,
	release func(A, E.Either[error, B]) GANY,
) GB {
	return G.Bracket[GA, GB, GANY, E.Either[error, B], A, B](
		I.Of[GB, TB, context.Context, E.Either[error, B]],
		MonadChain[GA, GB, TA, TB, A, B],
		I.MonadChain[GB, GB, TB, TB, context.Context, E.Either[error, B], E.Either[error, B]],
		MonadChain[GANY, GB, TANY, TB, ANY, B],

		acquire,
		use,
		release,
	)
}
