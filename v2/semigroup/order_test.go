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

package semigroup_test

import (
	"errors"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	IO "github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/semigroup"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

// The tests in this file use the non-commutative string semigroup and effect logs to pin
// the argument order of the curried map, ap and alt operations.

func TestApplySemigroupOrder(t *testing.T) {
	sg := semigroup.ApplySemigroup(O.Map[string, func(string) string], O.Ap[string, string], S.Monoid)

	t.Run("combines values in argument order", func(t *testing.T) {
		assert.Equal(t, O.Some("ab"), sg.Concat(O.Some("a"), O.Some("b")))
	})

	t.Run("is None when either argument is None", func(t *testing.T) {
		assert.Equal(t, O.None[string](), sg.Concat(O.None[string](), O.Some("b")))
		assert.Equal(t, O.None[string](), sg.Concat(O.Some("a"), O.None[string]()))
	})

	t.Run("is associative", func(t *testing.T) {
		a, b, c := O.Some("a"), O.Some("b"), O.Some("c")
		assert.Equal(t, sg.Concat(sg.Concat(a, b), c), sg.Concat(a, sg.Concat(b, c)))
		assert.Equal(t, O.Some("abc"), sg.Concat(a, sg.Concat(b, c)))
	})

	t.Run("runs the effect of the first argument first", func(t *testing.T) {
		var log []string
		emit := func(s string) IO.IO[string] {
			return func() string {
				log = append(log, s)
				return s
			}
		}
		seq := semigroup.ApplySemigroup(IO.Map[string, func(string) string], IO.ApSeq[string, string], S.Monoid)

		assert.Equal(t, "ab", seq.Concat(emit("a"), emit("b"))())
		assert.Equal(t, []string{"a", "b"}, log)
	})

	t.Run("reports the error of the first argument when both fail", func(t *testing.T) {
		errA, errB := errors.New("a"), errors.New("b")
		esg := semigroup.ApplySemigroup(E.Map[error, string, func(string) string], E.Ap[string, error, string], S.Monoid)

		assert.Equal(t, E.Left[string](errA), esg.Concat(E.Left[string](errA), E.Left[string](errB)))
	})
}

func TestAltSemigroupOrder(t *testing.T) {
	sg := semigroup.AltSemigroup(O.Alt[string])

	t.Run("prefers the first argument", func(t *testing.T) {
		assert.Equal(t, O.Some("a"), sg.Concat(O.Some("a"), O.Some("b")))
	})

	t.Run("falls back to the second argument", func(t *testing.T) {
		assert.Equal(t, O.Some("b"), sg.Concat(O.None[string](), O.Some("b")))
	})

	t.Run("is None when both arguments are None", func(t *testing.T) {
		assert.Equal(t, O.None[string](), sg.Concat(O.None[string](), O.None[string]()))
	})

	t.Run("runs the second effect only when the first fails", func(t *testing.T) {
		var log []string
		emit := func(s string) IOE.IOEither[error, string] {
			return func() E.Either[error, string] {
				log = append(log, s)
				return E.Right[error](s)
			}
		}
		fail := func(s string) IOE.IOEither[error, string] {
			return func() E.Either[error, string] {
				log = append(log, s)
				return E.Left[string](errors.New(s))
			}
		}
		esg := semigroup.AltSemigroup(IOE.Alt[error, string])

		assert.Equal(t, E.Right[error]("a"), esg.Concat(emit("a"), emit("b"))())
		assert.Equal(t, []string{"a"}, log)

		log = nil
		assert.Equal(t, E.Right[error]("b"), esg.Concat(fail("a"), emit("b"))())
		assert.Equal(t, []string{"a", "b"}, log)
	})
}
