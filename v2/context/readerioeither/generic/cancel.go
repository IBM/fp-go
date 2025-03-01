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

	CIOE "github.com/IBM/fp-go/v2/context/ioeither/generic"
	E "github.com/IBM/fp-go/v2/either"
	IOE "github.com/IBM/fp-go/v2/ioeither/generic"
)

// withContext wraps an existing ReaderIOEither and performs a context check for cancellation before delegating
func WithContext[GRA ~func(context.Context) GIOA, GIOA ~func() E.Either[error, A], A any](ma GRA) GRA {
	return func(ctx context.Context) GIOA {
		if err := context.Cause(ctx); err != nil {
			return IOE.Left[GIOA](err)
		}
		return CIOE.WithContext(ctx, ma(ctx))
	}
}
