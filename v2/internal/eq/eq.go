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

package eq

import (
	EQ "github.com/IBM/fp-go/v2/eq"
	F "github.com/IBM/fp-go/v2/function"
)

// Eq implements an equals predicate on the basis of `map` and `ap`
func Eq[HKTA, HKTABOOL, HKTBOOL, A any](
	fmap func(HKTA, func(A) func(A) bool) HKTABOOL,
	fap func(HKTABOOL, HKTA) HKTBOOL,

	e EQ.Eq[A],
) func(l, r HKTA) HKTBOOL {
	c := F.Curry2(e.Equals)
	return func(fl, fr HKTA) HKTBOOL {
		return fap(fmap(fl, c), fr)
	}
}
