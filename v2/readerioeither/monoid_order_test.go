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

package readerioeither

import (
	"errors"
	"testing"

	"github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	M "github.com/IBM/fp-go/v2/monoid"
	"github.com/stretchr/testify/assert"
)

// The tests in this file use the non-commutative string monoid and effect logs to pin
// the argument order of the curried Map, Ap and Alt operations used by the monoids.

var (
	orderMonoid   = M.MakeMonoid(func(a, b string) string { return a + b }, "")
	errOrderA     = errors.New("a")
	errOrderB     = errors.New("b")
	errOrderEmpty = errors.New("empty")
)

// orderTagged returns a computation that appends tag to its environment
func orderTagged(tag string) ReaderIOEither[string, error, string] {
	return func(env string) func() either.Either[error, string] {
		return func() either.Either[error, string] {
			return either.Right[error](env + tag)
		}
	}
}

// orderLogged returns a computation that appends tag to its environment and records tag in log
func orderLogged(log *[]string, tag string) ReaderIOEither[string, error, string] {
	return func(env string) func() either.Either[error, string] {
		return func() either.Either[error, string] {
			*log = append(*log, tag)
			return either.Right[error](env + tag)
		}
	}
}

func TestApplicativeMonoidOrder(t *testing.T) {
	for name, m := range map[string]Monoid[string, error, string]{
		"default":    ApplicativeMonoid[string, error](orderMonoid),
		"sequential": ApplicativeMonoidSeq[string, error](orderMonoid),
		"parallel":   ApplicativeMonoidPar[string, error](orderMonoid),
	} {
		t.Run(name, func(t *testing.T) {
			a, b, c := orderTagged(":a"), orderTagged(":b"), orderTagged(":c")

			assert.Equal(t, either.Right[error]("x:ax:b"), m.Concat(a, b)("x")())
			assert.Equal(t, either.Right[error](""), m.Empty()("x")())
			assert.Equal(t, either.Right[error]("x:ax:bx:c"), m.Concat(m.Concat(a, b), c)("x")())
			assert.Equal(t, either.Right[error]("x:ax:bx:c"), m.Concat(a, m.Concat(b, c))("x")())
			assert.Equal(t, either.Left[string](errOrderA), m.Concat(Left[string, string](errOrderA), b)("x")())
			assert.Equal(t, either.Left[string](errOrderB), m.Concat(a, Left[string, string](errOrderB))("x")())
		})
	}
}

func TestApplicativeMonoidSeqEffectOrder(t *testing.T) {
	var log []string
	m := ApplicativeMonoidSeq[string, error](orderMonoid)

	assert.Equal(t, either.Right[error]("x:ax:b"), m.Concat(orderLogged(&log, ":a"), orderLogged(&log, ":b"))("x")())
	assert.Equal(t, []string{":a", ":b"}, log, "runs the effect of the first argument first")

	res := m.Concat(Left[string, string](errOrderA), Left[string, string](errOrderB))("x")()
	assert.Equal(t, either.Left[string](errOrderA), res)
}

func TestAlternativeMonoidOrder(t *testing.T) {
	m := AlternativeMonoid[string, error](orderMonoid)

	t.Run("combines two successes in argument order", func(t *testing.T) {
		assert.Equal(t, either.Right[error]("x:ax:b"), m.Concat(orderTagged(":a"), orderTagged(":b"))("x")())
	})

	t.Run("falls back to the success", func(t *testing.T) {
		assert.Equal(t, either.Right[error]("x:b"), m.Concat(Left[string, string](errOrderA), orderTagged(":b"))("x")())
		assert.Equal(t, either.Right[error]("x:a"), m.Concat(orderTagged(":a"), Left[string, string](errOrderB))("x")())
	})

	t.Run("reports the error of the second when both fail", func(t *testing.T) {
		res := m.Concat(Left[string, string](errOrderA), Left[string, string](errOrderB))("x")()
		assert.Equal(t, either.Left[string](errOrderB), res)
	})

	t.Run("empty yields the empty value", func(t *testing.T) {
		assert.Equal(t, either.Right[error](""), m.Empty()("x")())
	})
}

func TestAltMonoidOrder(t *testing.T) {
	m := AltMonoid[string, error, string](F.Constant(Left[string, string](errOrderEmpty)))

	t.Run("does not run the second effect when the first succeeds", func(t *testing.T) {
		var log []string
		assert.Equal(t, either.Right[error]("x:a"), m.Concat(orderLogged(&log, ":a"), orderLogged(&log, ":b"))("x")())
		assert.Equal(t, []string{":a"}, log)
	})

	t.Run("falls back to the second argument", func(t *testing.T) {
		assert.Equal(t, either.Right[error]("x:b"), m.Concat(Left[string, string](errOrderA), orderTagged(":b"))("x")())
	})

	t.Run("empty is the zero value and an identity", func(t *testing.T) {
		assert.Equal(t, either.Left[string](errOrderEmpty), m.Empty()("x")())
		assert.Equal(t, either.Right[error]("x:a"), m.Concat(m.Empty(), orderTagged(":a"))("x")())
		assert.Equal(t, either.Right[error]("x:a"), m.Concat(orderTagged(":a"), m.Empty())("x")())
	})
}
