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

package ioeither

import (
	"context"

	"github.com/IBM/fp-go/v2/either"
	IOE "github.com/IBM/fp-go/v2/ioeither"
)

// withContext wraps an existing IOEither and performs a context check for cancellation before delegating
func WithContext[A any](ctx context.Context, ma IOE.IOEither[error, A]) IOE.IOEither[error, A] {
	return func() either.Either[error, A] {
		if err := context.Cause(ctx); err != nil {
			return either.Left[A](err)
		}
		return ma()
	}
}
