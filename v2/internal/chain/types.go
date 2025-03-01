// Copyright (c) 2024 IBM Corp.
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

package chain

import (
	"github.com/IBM/fp-go/v2/internal/apply"
	"github.com/IBM/fp-go/v2/internal/functor"
)

type Chainable[A, B, HKTA, HKTB, HKTFAB any] interface {
	apply.Apply[A, B, HKTA, HKTB, HKTFAB]
	Chain(func(A) HKTB) func(HKTA) HKTB
}

// ToFunctor converts from [Chainable] to [functor.Functor]
func ToFunctor[A, B, HKTA, HKTB, HKTFAB any](ap Chainable[A, B, HKTA, HKTB, HKTFAB]) functor.Functor[A, B, HKTA, HKTB] {
	return ap
}

// ToApply converts from [Chainable] to [functor.Functor]
func ToApply[A, B, HKTA, HKTB, HKTFAB any](ap Chainable[A, B, HKTA, HKTB, HKTFAB]) apply.Apply[A, B, HKTA, HKTB, HKTFAB] {
	return ap
}
