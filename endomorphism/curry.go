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

package endomorphism

import (
	G "github.com/IBM/fp-go/endomorphism/generic"
)

// Curry2 curries a binary function
func Curry2[FCT ~func(T0, T1) T1, T0, T1 any](f FCT) func(T0) Endomorphism[T1] {
	return G.Curry2[Endomorphism[T1]](f)
}

// Curry3 curries a ternary function
func Curry3[FCT ~func(T0, T1, T2) T2, T0, T1, T2 any](f FCT) func(T0) func(T1) Endomorphism[T2] {
	return G.Curry3[Endomorphism[T2]](f)
}
