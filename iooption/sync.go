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

package iooption

import (
	"context"

	G "github.com/IBM/fp-go/iooption/generic"
	L "github.com/IBM/fp-go/lazy"
)

// WithLock executes the provided IO operation in the scope of a lock
func WithLock[E, A any](lock L.Lazy[context.CancelFunc]) func(fa IOOption[A]) IOOption[A] {
	return G.WithLock[IOOption[A]](lock)
}
