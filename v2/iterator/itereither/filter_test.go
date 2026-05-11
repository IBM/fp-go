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

package itereither

import (
	"fmt"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/iterator/iter"
	N "github.com/IBM/fp-go/v2/number"
	S "github.com/IBM/fp-go/v2/string"
	"github.com/stretchr/testify/assert"
)

func TestFilterOrElse_PredicateTrue(t *testing.T) {
	isPositive := N.MoreThan(0)
	onFalse := S.Format[int]("%d is not positive")

	filter := FilterOrElse(isPositive, onFalse)
	seq := iter.From(E.Right[string](42))
	result := collectEithers(filter(seq))

	assert.Equal(t, []Either[string, int]{E.Right[string](42)}, result)
}

func TestFilterOrElse_PredicateFalse(t *testing.T) {
	isPositive := N.MoreThan(0)
	onFalse := S.Format[int]("%d is not positive")

	filter := FilterOrElse(isPositive, onFalse)
	seq := iter.From(E.Right[string](-5))
	result := collectEithers(filter(seq))

	assert.Equal(t, []Either[string, int]{E.Left[int]("-5 is not positive")}, result)
}

func TestFilterOrElse_LeftPassesThrough(t *testing.T) {
	isPositive := N.MoreThan(0)
	onFalse := S.Format[int]("%d is not positive")

	filter := FilterOrElse(isPositive, onFalse)
	seq := iter.From(E.Left[int]("original error"))
	result := collectEithers(filter(seq))

	assert.Equal(t, []Either[string, int]{E.Left[int]("original error")}, result)
}

func TestFilterOrElse_ZeroValue(t *testing.T) {
	isNonZero := func(n int) bool { return n != 0 }
	onZero := func(n int) string { return "value is zero" }

	filter := FilterOrElse(isNonZero, onZero)
	seq := iter.From(E.Right[string](0))
	result := collectEithers(filter(seq))

	assert.Equal(t, []Either[string, int]{E.Left[int]("value is zero")}, result)
}

func TestFilterOrElse_StringValidation(t *testing.T) {
	isNonEmpty := S.IsNonEmpty
	onEmpty := func(s string) error { return fmt.Errorf("string is empty") }

	filter := FilterOrElse(isNonEmpty, onEmpty)

	t.Run("non-empty string passes", func(t *testing.T) {
		seq := iter.From(E.Right[error]("hello"))
		result := collectEithers(filter(seq))
		assert.Equal(t, []Either[error, string]{E.Right[error]("hello")}, result)
	})

	t.Run("empty string fails", func(t *testing.T) {
		seq := iter.From(E.Right[error](""))
		result := collectEithers(filter(seq))
		assert.Len(t, result, 1)
		assert.True(t, E.IsLeft(result[0]))
	})
}

func TestFilterOrElse_MultipleValues(t *testing.T) {
	isEven := func(n int) bool { return n%2 == 0 }
	onOdd := S.Format[int]("%d is odd")

	filter := FilterOrElse(isEven, onOdd)
	seq := iter.From(
		E.Right[string](2),
		E.Right[string](3),
		E.Right[string](4),
		E.Left[int]("error"),
		E.Right[string](5),
	)
	result := collectEithers(filter(seq))

	expected := []Either[string, int]{
		E.Right[string](2),
		E.Left[int]("3 is odd"),
		E.Right[string](4),
		E.Left[int]("error"),
		E.Left[int]("5 is odd"),
	}
	assert.Equal(t, expected, result)
}

func TestFilterOrElse_InPipeline(t *testing.T) {
	isPositive := N.MoreThan(0)
	onNegative := S.Format[int]("%d is not positive")

	result := F.Pipe2(
		iter.From(1, -2, 3, -4, 5),
		FromSeq[string],
		FilterOrElse(isPositive, onNegative),
	)

	collected := collectEithers(result)
	expected := []Either[string, int]{
		E.Right[string](1),
		E.Left[int]("-2 is not positive"),
		E.Right[string](3),
		E.Left[int]("-4 is not positive"),
		E.Right[string](5),
	}
	assert.Equal(t, expected, collected)
}

func TestFilterOrElse_WithChain(t *testing.T) {
	isPositive := N.MoreThan(0)
	onNegative := S.Format[int]("%d is not positive")

	result := F.Pipe3(
		iter.From(1, 2, 3),
		FromSeq[string],
		Map[string](func(n int) int { return n - 2 }),
		FilterOrElse(isPositive, onNegative),
	)

	collected := collectEithers(result)
	expected := []Either[string, int]{
		E.Left[int]("-1 is not positive"),
		E.Left[int]("0 is not positive"),
		E.Right[string](1),
	}
	assert.Equal(t, expected, collected)
}

func TestFilterOrElse_CustomPredicate(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}

	isAdult := func(u User) bool { return u.Age >= 18 }
	onMinor := func(u User) string {
		return fmt.Sprintf("%s is only %d years old", u.Name, u.Age)
	}

	filter := FilterOrElse(isAdult, onMinor)
	seq := iter.From(
		E.Right[string](User{"Alice", 25}),
		E.Right[string](User{"Bob", 16}),
		E.Right[string](User{"Charlie", 30}),
	)

	result := collectEithers(filter(seq))
	assert.Len(t, result, 3)
	assert.True(t, E.IsRight(result[0]))
	assert.True(t, E.IsLeft(result[1]))
	assert.True(t, E.IsRight(result[2]))
}

// Made with Bob
