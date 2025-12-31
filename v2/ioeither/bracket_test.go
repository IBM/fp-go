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
	"github.com/stretchr/testify/assert"
)

func TestBracket(t *testing.T) {
	t.Run("successful acquisition, use, and release", func(t *testing.T) {
		acquired := false
		used := false
		released := false

		acquire := func() IOEither[error, string] {
			return func() E.Either[error, string] {
				acquired = true
				return E.Right[error]("resource")
			}
		}()

		use := func(r string) IOEither[error, int] {
			return func() E.Either[error, int] {
				used = true
				assert.Equal(t, "resource", r)
				return E.Right[error](42)
			}
		}

		release := func(r string, result E.Either[error, int]) IOEither[error, any] {
			return func() E.Either[error, any] {
				released = true
				assert.Equal(t, "resource", r)
				assert.True(t, E.IsRight(result))
				return E.Right[error, any](nil)
			}
		}

		result := Bracket(acquire, use, release)()

		assert.True(t, acquired)
		assert.True(t, used)
		assert.True(t, released)
		assert.Equal(t, E.Right[error](42), result)
	})

	t.Run("acquisition fails", func(t *testing.T) {
		used := false
		released := false

		acquire := Left[string](errors.New("acquisition failed"))

		use := func(r string) IOEither[error, int] {
			used = true
			return Of[error](42)
		}

		release := func(r string, result E.Either[error, int]) IOEither[error, any] {
			released = true
			return Of[error, any](nil)
		}

		result := Bracket(acquire, use, release)()

		assert.False(t, used)
		assert.False(t, released)
		assert.True(t, E.IsLeft(result))
	})

	t.Run("use fails but release is called", func(t *testing.T) {
		acquired := false
		released := false
		var releaseResult E.Either[error, int]

		acquire := func() IOEither[error, string] {
			return func() E.Either[error, string] {
				acquired = true
				return E.Right[error]("resource")
			}
		}()

		use := func(r string) IOEither[error, int] {
			return Left[int](errors.New("use failed"))
		}

		release := func(r string, result E.Either[error, int]) IOEither[error, any] {
			return func() E.Either[error, any] {
				released = true
				releaseResult = result
				assert.Equal(t, "resource", r)
				return E.Right[error, any](nil)
			}
		}

		result := Bracket(acquire, use, release)()

		assert.True(t, acquired)
		assert.True(t, released)
		assert.True(t, E.IsLeft(result))
		assert.True(t, E.IsLeft(releaseResult))
	})

	t.Run("release is called even when use succeeds", func(t *testing.T) {
		releaseCallCount := 0

		acquire := Of[error]("resource")

		use := func(r string) IOEither[error, int] {
			return Of[error](100)
		}

		release := func(r string, result E.Either[error, int]) IOEither[error, any] {
			return func() E.Either[error, any] {
				releaseCallCount++
				return E.Right[error, any](nil)
			}
		}

		result := Bracket(acquire, use, release)()

		assert.Equal(t, 1, releaseCallCount)
		assert.Equal(t, E.Right[error](100), result)
	})

	t.Run("release error overrides successful result", func(t *testing.T) {
		acquire := Of[error]("resource")

		use := func(r string) IOEither[error, int] {
			return Of[error](42)
		}

		release := func(r string, result E.Either[error, int]) IOEither[error, any] {
			return Left[any](errors.New("release failed"))
		}

		result := Bracket(acquire, use, release)()

		// According to bracket semantics, release errors are propagated
		assert.True(t, E.IsLeft(result))
	})
}
