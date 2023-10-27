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

	F "github.com/IBM/fp-go/function"
	G "github.com/IBM/fp-go/io/generic"
)

// WithLock executes the provided IO operation in the scope of a lock
func WithLock[GEA ~func(R) GIOA, GIOA ~func() A, R, A any](lock func() context.CancelFunc) func(fa GEA) GEA {
	l := G.WithLock[GIOA](lock)
	return func(fa GEA) GEA {
		return F.Flow2(
			fa,
			l,
		)
	}
}
