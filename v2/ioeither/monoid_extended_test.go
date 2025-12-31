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

	E "github.com/IBM/fp-go/v2/either"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

func TestApplicativeMonoidSeq(t *testing.T) {
	m := ApplicativeMonoidSeq[error](S.Monoid)

	t.Run("concatenates two Right values sequentially", func(t *testing.T) {
		result := m.Concat(Of[error]("hello"), Of[error](" world"))()
		assert.Equal(t, E.Of[error]("hello world"), result)
	})

	t.Run("empty with Right value", func(t *testing.T) {
		result := m.Concat(Of[error]("test"), m.Empty())()
		assert.Equal(t, E.Of[error]("test"), result)
	})

	t.Run("Right value with empty", func(t *testing.T) {
		result := m.Concat(m.Empty(), Of[error]("test"))()
		assert.Equal(t, E.Of[error]("test"), result)
	})

	t.Run("Left value short-circuits", func(t *testing.T) {
		err := errors.New("error")
		result := m.Concat(Left[string](err), Of[error]("test"))()
		assert.True(t, E.IsLeft(result))
	})

	t.Run("Right with Left returns Left", func(t *testing.T) {
		err := errors.New("error")
		result := m.Concat(Of[error]("test"), Left[string](err))()
		assert.True(t, E.IsLeft(result))
	})
}

func TestApplicativeMonoidPar(t *testing.T) {
	m := ApplicativeMonoidPar[error](S.Monoid)

	t.Run("concatenates two Right values in parallel", func(t *testing.T) {
		result := m.Concat(Of[error]("hello"), Of[error](" world"))()
		assert.Equal(t, E.Of[error]("hello world"), result)
	})

	t.Run("empty with Right value", func(t *testing.T) {
		result := m.Concat(Of[error]("test"), m.Empty())()
		assert.Equal(t, E.Of[error]("test"), result)
	})

	t.Run("Right value with empty", func(t *testing.T) {
		result := m.Concat(m.Empty(), Of[error]("test"))()
		assert.Equal(t, E.Of[error]("test"), result)
	})

	t.Run("Left value returns Left", func(t *testing.T) {
		err := errors.New("error")
		result := m.Concat(Left[string](err), Of[error]("test"))()
		assert.True(t, E.IsLeft(result))
	})

	t.Run("Right with Left returns Left", func(t *testing.T) {
		err := errors.New("error")
		result := m.Concat(Of[error]("test"), Left[string](err))()
		assert.True(t, E.IsLeft(result))
	})

	t.Run("multiple concatenations", func(t *testing.T) {
		result := m.Concat(
			m.Concat(Of[error]("a"), Of[error]("b")),
			m.Concat(Of[error]("c"), Of[error]("d")),
		)()
		assert.Equal(t, E.Of[error]("abcd"), result)
	})
}
