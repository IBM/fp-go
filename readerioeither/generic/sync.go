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

	ET "github.com/IBM/fp-go/either"
	G "github.com/IBM/fp-go/readerio/generic"
)

// WithLock executes the provided IO operation in the scope of a lock
func WithLock[GEA ~func(R) GIOA, GIOA ~func() ET.Either[E, A], R, E, A any](lock func() context.CancelFunc) func(fa GEA) GEA {
	return G.WithLock[GEA](lock)
}
