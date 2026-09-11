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

package monoid_test

import (
	"errors"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	IO "github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/lazy"
	M "github.com/IBM/fp-go/v2/monoid"
	O "github.com/IBM/fp-go/v2/option"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// The tests in this file use the non-commutative string monoid and effect logs to pin
// the argument order of the curried map, ap and alt operations.

func TestApplicativeMonoidOrder(t *testing.T) {
	m := M.ApplicativeMonoid(O.Of[string], O.Map[string, func(string) string], O.Ap[string, string], S.Monoid)

	t.Run("combines values in argument order", func(t *testing.T) {
		assert.Equal(t, O.Some("ab"), m.Concat(O.Some("a"), O.Some("b")))
	})

	t.Run("empty lifts the empty value", func(t *testing.T) {
		assert.Equal(t, O.Some(""), m.Empty())
	})

	t.Run("empty is a left and right identity", func(t *testing.T) {
		assert.Equal(t, O.Some("a"), m.Concat(m.Empty(), O.Some("a")))
		assert.Equal(t, O.Some("a"), m.Concat(O.Some("a"), m.Empty()))
	})

	t.Run("is associative", func(t *testing.T) {
		a, b, c := O.Some("a"), O.Some("b"), O.Some("c")
		assert.Equal(t, m.Concat(m.Concat(a, b), c), m.Concat(a, m.Concat(b, c)))
		assert.Equal(t, O.Some("abc"), m.Concat(a, m.Concat(b, c)))
	})

	t.Run("runs the effect of the first argument first", func(t *testing.T) {
		var log []string
		emit := func(s string) IO.IO[string] {
			return func() string {
				log = append(log, s)
				return s
			}
		}
		seq := M.ApplicativeMonoid(IO.Of[string], IO.Map[string, func(string) string], IO.ApSeq[string, string], S.Monoid)

		assert.Equal(t, "ab", seq.Concat(emit("a"), emit("b"))())
		assert.Equal(t, []string{"a", "b"}, log)
	})
}

func TestAlternativeMonoidOrder(t *testing.T) {
	m := M.AlternativeMonoid(
		E.Of[error, string],
		E.Map[error, string, func(string) string],
		E.Ap[string, error, string],
		E.Alt[error, string],
		S.Monoid,
	)
	errA, errB := errors.New("a"), errors.New("b")
	a, b := E.Right[error]("a"), E.Right[error]("b")
	failA, failB := E.Left[string](errA), E.Left[string](errB)

	t.Run("combines two successes in argument order", func(t *testing.T) {
		assert.Equal(t, E.Right[error]("ab"), m.Concat(a, b))
	})

	t.Run("falls back to the second when the first fails", func(t *testing.T) {
		assert.Equal(t, b, m.Concat(failA, b))
	})

	t.Run("falls back to the first when the second fails", func(t *testing.T) {
		assert.Equal(t, a, m.Concat(a, failB))
	})

	t.Run("reports the error of the second when both fail", func(t *testing.T) {
		assert.Equal(t, failB, m.Concat(failA, failB))
	})

	t.Run("empty lifts the empty value", func(t *testing.T) {
		assert.Equal(t, E.Right[error](""), m.Empty())
	})

	t.Run("does not evaluate the fallback when both succeed", func(t *testing.T) {
		var log []string
		emit := func(s string) IOE.IOEither[error, string] {
			return func() E.Either[error, string] {
				log = append(log, s)
				return E.Right[error](s)
			}
		}
		seq := M.AlternativeMonoid(
			IOE.Of[error, string],
			IOE.Map[error, string, func(string) string],
			IOE.ApSeq[string, error, string],
			IOE.Alt[error, string],
			S.Monoid,
		)

		assert.Equal(t, E.Right[error]("ab"), seq.Concat(emit("a"), emit("b"))())
		assert.Equal(t, []string{"a", "b"}, log)
	})
}

func TestAltMonoidOrder(t *testing.T) {
	m := M.AltMonoid(O.None[string], O.Alt[string])

	t.Run("prefers the first argument", func(t *testing.T) {
		assert.Equal(t, O.Some("a"), m.Concat(O.Some("a"), O.Some("b")))
	})

	t.Run("falls back to the second argument", func(t *testing.T) {
		assert.Equal(t, O.Some("b"), m.Concat(O.None[string](), O.Some("b")))
	})

	t.Run("empty is the zero value and an identity", func(t *testing.T) {
		assert.Equal(t, O.None[string](), m.Empty())
		assert.Equal(t, O.Some("a"), m.Concat(m.Empty(), O.Some("a")))
		assert.Equal(t, O.Some("a"), m.Concat(O.Some("a"), m.Empty()))
	})

	t.Run("runs the second effect only when the first fails", func(t *testing.T) {
		var log []string
		emit := func(s string) IOE.IOEither[error, string] {
			return func() E.Either[error, string] {
				log = append(log, s)
				return E.Right[error](s)
			}
		}
		alt := M.AltMonoid(lazy.Of(IOE.Left[string](errors.New("empty"))), IOE.Alt[error, string])

		assert.Equal(t, E.Right[error]("a"), alt.Concat(emit("a"), emit("b"))())
		assert.Equal(t, []string{"a"}, log)

		log = nil
		assert.Equal(t, E.Right[error]("b"), alt.Concat(alt.Empty(), emit("b"))())
		assert.Equal(t, []string{"b"}, log)
	})
}
