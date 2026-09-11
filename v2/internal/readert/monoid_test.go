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

package readert_test

import (
	"testing"

	"github.com/IBM/fp-go/v2/internal/readert"
	"github.com/IBM/fp-go/v2/reader"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// tagged returns a reader that appends tag to its environment
func tagged(tag string) reader.Reader[string, string] {
	return func(env string) string {
		return env + tag
	}
}

func TestApplySemigroup(t *testing.T) {
	sg := readert.ApplySemigroup(
		reader.Map[string, string, func(string) string],
		reader.Ap[string, string, string],
		S.Monoid,
	)

	t.Run("combines the results in argument order", func(t *testing.T) {
		assert.Equal(t, "x:ax:b", sg.Concat(tagged(":a"), tagged(":b"))("x"))
	})

	t.Run("is associative", func(t *testing.T) {
		a, b, c := tagged(":a"), tagged(":b"), tagged(":c")
		assert.Equal(t, sg.Concat(sg.Concat(a, b), c)("x"), sg.Concat(a, sg.Concat(b, c))("x"))
	})
}

func TestApplicativeMonoid(t *testing.T) {
	m := readert.ApplicativeMonoid(
		reader.Of[string, string],
		reader.Map[string, string, func(string) string],
		reader.Ap[string, string, string],
		S.Monoid,
	)

	t.Run("combines the results in argument order", func(t *testing.T) {
		assert.Equal(t, "x:ax:b", m.Concat(tagged(":a"), tagged(":b"))("x"))
	})

	t.Run("empty ignores the environment and yields the empty value", func(t *testing.T) {
		assert.Equal(t, "", m.Empty()("x"))
	})

	t.Run("empty is a left and right identity", func(t *testing.T) {
		assert.Equal(t, "x:a", m.Concat(m.Empty(), tagged(":a"))("x"))
		assert.Equal(t, "x:a", m.Concat(tagged(":a"), m.Empty())("x"))
	})

	t.Run("is associative", func(t *testing.T) {
		a, b, c := tagged(":a"), tagged(":b"), tagged(":c")
		assert.Equal(t, "x:ax:bx:c", m.Concat(m.Concat(a, b), c)("x"))
		assert.Equal(t, "x:ax:bx:c", m.Concat(a, m.Concat(b, c))("x"))
	})
}
