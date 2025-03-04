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

package readerioeither

import (
	"context"

	CIOE "github.com/IBM/fp-go/v2/context/ioeither/generic"
	"github.com/IBM/fp-go/v2/ioeither"
)

// WithContext wraps an existing [ReaderIOEither] and performs a context check for cancellation before delegating
func WithContext[A any](ma ReaderIOEither[A]) ReaderIOEither[A] {
	return func(ctx context.Context) ioeither.IOEither[error, A] {
		if err := context.Cause(ctx); err != nil {
			return ioeither.Left[A](err)
		}
		return CIOE.WithContext(ctx, ma(ctx))
	}
}
