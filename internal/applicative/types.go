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

package applicative

import (
	"github.com/IBM/fp-go/internal/apply"
	"github.com/IBM/fp-go/internal/functor"
	"github.com/IBM/fp-go/internal/pointed"
)

type Applicative[A, B, HKTA, HKTB, HKTFAB any] interface {
	apply.Apply[A, B, HKTA, HKTB, HKTFAB]
	pointed.Pointed[A, HKTA]
}

// ToFunctor converts from [Applicative] to [functor.Functor]
func ToFunctor[A, B, HKTA, HKTB, HKTFAB any](ap Applicative[A, B, HKTA, HKTB, HKTFAB]) functor.Functor[A, B, HKTA, HKTB] {
	return ap
}

// ToApply converts from [Applicative] to [apply.Apply]
func ToApply[A, B, HKTA, HKTB, HKTFAB any](ap Applicative[A, B, HKTA, HKTB, HKTFAB]) apply.Apply[A, B, HKTA, HKTB, HKTFAB] {
	return ap
}

// ToPointed converts from [Applicative] to [pointed.Pointed]
func ToPointed[A, B, HKTA, HKTB, HKTFAB any](ap Applicative[A, B, HKTA, HKTB, HKTFAB]) pointed.Pointed[A, HKTA] {
	return ap
}
