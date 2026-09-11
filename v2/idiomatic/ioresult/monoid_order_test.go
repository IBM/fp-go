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

package ioresult

import (
	"errors"
	"testing"

	M "github.com/IBM/fp-go/v2/monoid"
	"github.com/stretchr/testify/assert"
)

// The tests in this file use the non-commutative string monoid and effect logs to pin
// the argument order of the curried Map, Ap and Alt operations used by the monoids.

var (
	orderMonoid = M.MakeMonoid(func(a, b string) string { return a + b }, "")
	errOrderA   = errors.New("a")
	errOrderB   = errors.New("b")
)

func TestApplicativeMonoidOrder(t *testing.T) {
	for name, m := range map[string]Monoid[string]{
		"default":    ApplicativeMonoid(orderMonoid),
		"sequential": ApplicativeMonoidSeq(orderMonoid),
		"parallel":   ApplicativeMonoidPar(orderMonoid),
	} {
		t.Run(name, func(t *testing.T) {
			a, b, c := Of("a"), Of("b"), Of("c")

			v, err := m.Concat(a, b)()
			assert.NoError(t, err)
			assert.Equal(t, "ab", v)

			v, err = m.Empty()()
			assert.NoError(t, err)
			assert.Equal(t, "", v)

			v, err = m.Concat(m.Concat(a, b), c)()
			assert.NoError(t, err)
			assert.Equal(t, "abc", v)

			v, err = m.Concat(a, m.Concat(b, c))()
			assert.NoError(t, err)
			assert.Equal(t, "abc", v)

			_, err = m.Concat(Left[string](errOrderA), b)()
			assert.ErrorIs(t, err, errOrderA)

			_, err = m.Concat(a, Left[string](errOrderB))()
			assert.ErrorIs(t, err, errOrderB)
		})
	}
}

func TestApplicativeMonoidSeqEffectOrder(t *testing.T) {
	var log []string
	emit := func(s string) IOResult[string] {
		return func() (string, error) {
			log = append(log, s)
			return s, nil
		}
	}
	m := ApplicativeMonoidSeq(orderMonoid)

	v, err := m.Concat(emit("a"), emit("b"))()
	assert.NoError(t, err)
	assert.Equal(t, "ab", v)
	assert.Equal(t, []string{"a", "b"}, log, "runs the effect of the first argument first")

	_, err = m.Concat(Left[string](errOrderA), Left[string](errOrderB))()
	assert.ErrorIs(t, err, errOrderA)
}

func TestAltSemigroupOrder(t *testing.T) {
	var log []string
	emit := func(s string) IOResult[string] {
		return func() (string, error) {
			log = append(log, s)
			return s, nil
		}
	}
	sg := AltSemigroup[string]()

	t.Run("does not run the second effect when the first succeeds", func(t *testing.T) {
		log = nil
		v, err := sg.Concat(emit("a"), emit("b"))()
		assert.NoError(t, err)
		assert.Equal(t, "a", v)
		assert.Equal(t, []string{"a"}, log)
	})

	t.Run("falls back to the second argument", func(t *testing.T) {
		v, err := sg.Concat(Left[string](errOrderA), Of("b"))()
		assert.NoError(t, err)
		assert.Equal(t, "b", v)
	})

	t.Run("reports the error of the second when both fail", func(t *testing.T) {
		_, err := sg.Concat(Left[string](errOrderA), Left[string](errOrderB))()
		assert.ErrorIs(t, err, errOrderB)
	})
}
