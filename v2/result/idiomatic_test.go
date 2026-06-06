// Copyright (c) 2025 IBM Corp.
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

package result

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// positiveToStr is a KleisliIdiomatic[int, string] for testing.
func positiveToStr(n int) (string, error) {
	if n > 0 {
		return strconv.Itoa(n), nil
	}
	return "", errors.New("non-positive")
}

// recoverNotFound is a KleisliIdiomatic[error, int] that recovers from "not found" errors.
func recoverNotFound(err error) (int, error) {
	if err.Error() == "not found" {
		return 0, nil
	}
	return 0, err
}

var errSentinel = errors.New("sentinel error")

func TestMonadChainIdiomatic(t *testing.T) {
	t.Run("Right input, function succeeds", func(t *testing.T) {
		assert.Equal(t, Right(42), MonadChainIdiomatic(Right("42"), strconv.Atoi))
	})
	t.Run("Right input, function fails", func(t *testing.T) {
		result := MonadChainIdiomatic(Right("abc"), strconv.Atoi)
		assert.True(t, IsLeft(result))
	})
	t.Run("Left input propagates", func(t *testing.T) {
		assert.Equal(t, Left[int](errSentinel), MonadChainIdiomatic(Left[string](errSentinel), strconv.Atoi))
	})
}

func TestChainIdiomatic(t *testing.T) {
	parse := ChainIdiomatic(strconv.Atoi)

	t.Run("Right input, function succeeds", func(t *testing.T) {
		assert.Equal(t, Right(7), parse(Right("7")))
	})
	t.Run("Right input, function fails", func(t *testing.T) {
		assert.True(t, IsLeft(parse(Right("bad"))))
	})
	t.Run("Left input propagates", func(t *testing.T) {
		assert.Equal(t, Left[int](errSentinel), parse(Left[string](errSentinel)))
	})
}

func TestMonadChainLeftIdiomatic(t *testing.T) {
	t.Run("Left input, function recovers", func(t *testing.T) {
		assert.Equal(t, Right(0), MonadChainLeftIdiomatic(Left[int](errors.New("not found")), recoverNotFound))
	})
	t.Run("Left input, function propagates", func(t *testing.T) {
		result := MonadChainLeftIdiomatic(Left[int](errSentinel), recoverNotFound)
		assert.Equal(t, Left[int](errSentinel), result)
	})
	t.Run("Right input passes through unchanged", func(t *testing.T) {
		assert.Equal(t, Right(42), MonadChainLeftIdiomatic(Right(42), recoverNotFound))
	})
}

func TestChainLeftIdiomatic(t *testing.T) {
	recover := ChainLeftIdiomatic(recoverNotFound)

	t.Run("Left input, function recovers", func(t *testing.T) {
		assert.Equal(t, Right(0), recover(Left[int](errors.New("not found"))))
	})
	t.Run("Left input, function propagates", func(t *testing.T) {
		assert.Equal(t, Left[int](errSentinel), recover(Left[int](errSentinel)))
	})
	t.Run("Right input passes through unchanged", func(t *testing.T) {
		assert.Equal(t, Right(42), recover(Right(42)))
	})
}

func TestMonadChainFirstIdiomatic(t *testing.T) {
	t.Run("Right input, function succeeds — original value kept", func(t *testing.T) {
		assert.Equal(t, Right(5), MonadChainFirstIdiomatic(Right(5), positiveToStr))
	})
	t.Run("Right input, function fails", func(t *testing.T) {
		result := MonadChainFirstIdiomatic(Right(-1), positiveToStr)
		assert.True(t, IsLeft(result))
	})
	t.Run("Left input propagates", func(t *testing.T) {
		assert.Equal(t, Left[int](errSentinel), MonadChainFirstIdiomatic(Left[int](errSentinel), positiveToStr))
	})
}

func TestChainFirstIdiomatic(t *testing.T) {
	keepIfPositive := ChainFirstIdiomatic(positiveToStr)

	t.Run("Right input, function succeeds — original value kept", func(t *testing.T) {
		assert.Equal(t, Right(5), keepIfPositive(Right(5)))
	})
	t.Run("Right input, function fails", func(t *testing.T) {
		assert.True(t, IsLeft(keepIfPositive(Right(-1))))
	})
	t.Run("Left input propagates", func(t *testing.T) {
		assert.Equal(t, Left[int](errSentinel), keepIfPositive(Left[int](errSentinel)))
	})
}
