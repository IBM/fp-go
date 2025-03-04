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

	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
)

// WithLock executes the provided IO operation in the scope of a lock
func WithLock[A any](lock ReaderIOEither[context.CancelFunc]) func(fa ReaderIOEither[A]) ReaderIOEither[A] {
	return function.Flow2(
		function.Constant1[context.CancelFunc, ReaderIOEither[A]],
		WithResource[A](lock, function.Flow2(
			io.FromImpure[context.CancelFunc],
			FromIO[any],
		)),
	)
}
