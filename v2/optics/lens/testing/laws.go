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

package testing

import (
	"testing"

	E "github.com/IBM/fp-go/v2/eq"
	L "github.com/IBM/fp-go/v2/optics/lens"
	"github.com/stretchr/testify/assert"
)

// LensGet tests the law:
// get(set(a)(s)) = a
func LensGet[S, A any](
	t *testing.T,
	eqa E.Eq[A],
) func(l L.Lens[S, A]) func(s S, a A) bool {

	return func(l L.Lens[S, A]) func(s S, a A) bool {

		return func(s S, a A) bool {
			return assert.True(t, eqa.Equals(l.Get(l.Set(a)(s)), a), "Lens get(set(a)(s)) = a")
		}
	}
}

// LensSet tests the laws:
// set(get(s))(s) = s
// set(a)(set(a)(s)) = set(a)(s)
func LensSet[S, A any](
	t *testing.T,
	eqs E.Eq[S],
) func(l L.Lens[S, A]) func(s S, a A) bool {

	return func(l L.Lens[S, A]) func(s S, a A) bool {

		return func(s S, a A) bool {
			return assert.True(t, eqs.Equals(l.Set(l.Get(s))(s), s), "Lens set(get(s))(s) = s") && assert.True(t, eqs.Equals(l.Set(a)(l.Set(a)(s)), l.Set(a)(s)), "Lens set(a)(set(a)(s)) = set(a)(s)")
		}
	}
}

// AssertLaws tests the lens laws
//
// get(set(a)(s)) = a
// set(get(s))(s) = s
// set(a)(set(a)(s)) = set(a)(s)
func AssertLaws[S, A any](
	t *testing.T,
	eqa E.Eq[A],
	eqs E.Eq[S],
) func(l L.Lens[S, A]) func(s S, a A) bool {

	lenGet := LensGet[S](t, eqa)
	lenSet := LensSet[S, A](t, eqs)

	return func(l L.Lens[S, A]) func(s S, a A) bool {

		get := lenGet(l)
		set := lenSet(l)

		return func(s S, a A) bool {
			return get(s, a) && set(s, a)
		}
	}
}
