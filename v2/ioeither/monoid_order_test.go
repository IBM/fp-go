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

package ioeither

import (
	"errors"
	"testing"

	"github.com/IBM/fp-go/v2/either"
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
	for name, m := range map[string]Monoid[error, string]{
		"default":    ApplicativeMonoid[error](orderMonoid),
		"sequential": ApplicativeMonoidSeq[error](orderMonoid),
		"parallel":   ApplicativeMonoidPar[error](orderMonoid),
	} {
		t.Run(name, func(t *testing.T) {
			a, b, c := Of[error]("a"), Of[error]("b"), Of[error]("c")

			assert.Equal(t, either.Right[error]("ab"), m.Concat(a, b)())
			assert.Equal(t, either.Right[error](""), m.Empty()())
			assert.Equal(t, either.Right[error]("abc"), m.Concat(m.Concat(a, b), c)())
			assert.Equal(t, either.Right[error]("abc"), m.Concat(a, m.Concat(b, c))())
			assert.Equal(t, either.Left[string](errOrderA), m.Concat(Left[string](errOrderA), b)())
			assert.Equal(t, either.Left[string](errOrderB), m.Concat(a, Left[string](errOrderB))())
		})
	}
}

func TestApplicativeMonoidSeqEffectOrder(t *testing.T) {
	var log []string
	emit := func(s string) IOEither[error, string] {
		return func() either.Either[error, string] {
			log = append(log, s)
			return either.Right[error](s)
		}
	}
	m := ApplicativeMonoidSeq[error](orderMonoid)

	assert.Equal(t, either.Right[error]("ab"), m.Concat(emit("a"), emit("b"))())
	assert.Equal(t, []string{"a", "b"}, log, "runs the effect of the first argument first")

	assert.Equal(t, either.Left[string](errOrderA), m.Concat(Left[string](errOrderA), Left[string](errOrderB))())
}

func TestAltSemigroupOrder(t *testing.T) {
	var log []string
	emit := func(s string) IOEither[error, string] {
		return func() either.Either[error, string] {
			log = append(log, s)
			return either.Right[error](s)
		}
	}
	sg := AltSemigroup[error, string]()

	t.Run("does not run the second effect when the first succeeds", func(t *testing.T) {
		log = nil
		assert.Equal(t, either.Right[error]("a"), sg.Concat(emit("a"), emit("b"))())
		assert.Equal(t, []string{"a"}, log)
	})

	t.Run("falls back to the second argument", func(t *testing.T) {
		log = nil
		assert.Equal(t, either.Right[error]("b"), sg.Concat(Left[string](errOrderA), emit("b"))())
		assert.Equal(t, []string{"b"}, log)
	})

	t.Run("reports the error of the second when both fail", func(t *testing.T) {
		assert.Equal(t, either.Left[string](errOrderB), sg.Concat(Left[string](errOrderA), Left[string](errOrderB))())
	})
}
