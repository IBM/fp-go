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
	F "github.com/IBM/fp-go/function"
	IO "github.com/IBM/fp-go/io/generic"
)

// WithLock executes the provided IO operation in the scope of a lock
func WithLock[
	GRA ~func(context.Context) GIOA,
	GIOA ~func() E.Either[error, A],
	GRCANCEL ~func(context.Context) GIOCANCEL,
	GIOCANCEL ~func() E.Either[error, context.CancelFunc],
	A any](lock GRCANCEL) func(fa GRA) GRA {

	type GRANY func(ctx context.Context) func() E.Either[error, any]
	type IOANY func() any

	return F.Flow2(
		F.Constant1[context.CancelFunc, GRA],
		WithResource[GRA, GRCANCEL, GRANY](lock, F.Flow2(
			IO.FromImpure[IOANY, context.CancelFunc],
			FromIO[GRANY, IOANY],
		)),
	)
}
