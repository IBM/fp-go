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

package reader

import (
	"testing"

	M "github.com/IBM/fp-go/v2/monoid"
	"github.com/stretchr/testify/assert"
)

func TestApplicativeMonoidOrder(t *testing.T) {
	m := ApplicativeMonoid[string](M.MakeMonoid(func(a, b string) string { return a + b }, ""))

	tagged := func(tag string) Reader[string, string] {
		return func(env string) string {
			return env + tag
		}
	}

	t.Run("combines the results in argument order", func(t *testing.T) {
		assert.Equal(t, "x:ax:b", m.Concat(tagged(":a"), tagged(":b"))("x"))
	})

	t.Run("empty yields the empty value", func(t *testing.T) {
		assert.Equal(t, "", m.Empty()("x"))
	})

	t.Run("is associative", func(t *testing.T) {
		a, b, c := tagged(":a"), tagged(":b"), tagged(":c")
		assert.Equal(t, "x:ax:bx:c", m.Concat(m.Concat(a, b), c)("x"))
		assert.Equal(t, "x:ax:bx:c", m.Concat(a, m.Concat(b, c))("x"))
	})
}
