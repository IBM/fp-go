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

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	RIE "github.com/IBM/fp-go/v2/readerioeither/generic"
)

// WithResource constructs a function that creates a resource, then operates on it and then releases the resource
func WithResource[
	GRA ~func(context.Context) GIOA,
	GRR ~func(context.Context) GIOR,
	GRANY ~func(context.Context) GIOANY,
	GIOR ~func() E.Either[error, R],
	GIOA ~func() E.Either[error, A],
	GIOANY ~func() E.Either[error, ANY],
	R, A, ANY any](onCreate GRR, onRelease func(R) GRANY) func(func(R) GRA) GRA {
	// wraps the callback functions with a context check
	return F.Flow2(
		F.Bind2nd(F.Flow2[func(R) GRA, func(GRA) GRA, R, GRA, GRA], WithContext[GRA]),
		RIE.WithResource[GRA](WithContext(onCreate), onRelease),
	)
}
