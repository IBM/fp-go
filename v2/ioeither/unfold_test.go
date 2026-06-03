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
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/pair"
	"github.com/stretchr/testify/assert"
)

// countDownStep is a reusable step function that counts down from n to 1,
// stopping (Right(None)) when n reaches 0.
func countDownStep(n int) IOEither[error, O.Option[P.Pair[int, int]]] {
	return func() E.Either[error, O.Option[P.Pair[int, int]]] {
		if n == 0 {
			return E.Of[error](O.None[P.Pair[int, int]]())
		}
		return E.Of[error](O.Some(P.MakePair(n-1, n)))
	}
}

func collectUnfold[E, A any](seq Seq[Either[E, A]]) []Either[E, A] {
	var result []Either[E, A]
	for v := range seq {
		result = append(result, v)
	}
	return result
}

func TestUnfoldEmpty(t *testing.T) {
	// seed = 0, so f returns Right(None) immediately
	seq := Unfold(countDownStep, 0)
	assert.Equal(t, []Either[error, int](nil), collectUnfold(seq))
}

func TestUnfoldFiniteSequence(t *testing.T) {
	// seed = 3, expect [Right(3), Right(2), Right(1)]
	seq := Unfold(countDownStep, 3)
	assert.Equal(t, []Either[error, int]{
		E.Of[error](3),
		E.Of[error](2),
		E.Of[error](1),
	}, collectUnfold(seq))
}

func TestUnfoldErrorOnFirstStep(t *testing.T) {
	errBoom := errors.New("boom")
	f := func(_ int) IOEither[error, O.Option[P.Pair[int, int]]] {
		return func() E.Either[error, O.Option[P.Pair[int, int]]] {
			return E.Left[O.Option[P.Pair[int, int]]](errBoom)
		}
	}
	seq := Unfold(f, 5)
	assert.Equal(t, []Either[error, int]{
		E.Left[int](errBoom),
	}, collectUnfold(seq))
}

func TestUnfoldErrorMidSequence(t *testing.T) {
	// Emits Right(2), Right(1), then Left(error) when seed reaches 0
	errDone := errors.New("done")
	f := func(n int) IOEither[error, O.Option[P.Pair[int, int]]] {
		return func() E.Either[error, O.Option[P.Pair[int, int]]] {
			if n == 0 {
				return E.Left[O.Option[P.Pair[int, int]]](errDone)
			}
			return E.Of[error](O.Some(P.MakePair(n-1, n)))
		}
	}
	seq := Unfold(f, 2)
	assert.Equal(t, []Either[error, int]{
		E.Of[error](2),
		E.Of[error](1),
		E.Left[int](errDone),
	}, collectUnfold(seq))
}

func TestUnfoldEarlyTermination(t *testing.T) {
	// Consumer stops after 2 elements; f would produce 5 but only 2 are taken.
	seq := Unfold(countDownStep, 5)
	var result []Either[error, int]
	for v := range seq {
		result = append(result, v)
		if len(result) == 2 {
			break
		}
	}
	assert.Equal(t, []Either[error, int]{
		E.Of[error](5),
		E.Of[error](4),
	}, result)
}
