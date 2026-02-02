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

package effect

import (
	"errors"
	"testing"

	"github.com/IBM/fp-go/v2/monoid"
	"github.com/stretchr/testify/assert"
)

func TestApplicativeMonoid(t *testing.T) {
	t.Run("combines successful effects with string monoid", func(t *testing.T) {
		stringMonoid := monoid.MakeMonoid(
			func(a, b string) string { return a + b },
			"",
		)

		effectMonoid := ApplicativeMonoid[TestContext](stringMonoid)

		eff1 := Of[TestContext]("Hello")
		eff2 := Of[TestContext](" ")
		eff3 := Of[TestContext]("World")

		combined := effectMonoid.Concat(eff1, effectMonoid.Concat(eff2, eff3))
		result, err := runEffect(combined, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "Hello World", result)
	})

	t.Run("combines successful effects with int monoid", func(t *testing.T) {
		intMonoid := monoid.MakeMonoid(
			func(a, b int) int { return a + b },
			0,
		)

		effectMonoid := ApplicativeMonoid[TestContext](intMonoid)

		eff1 := Of[TestContext](10)
		eff2 := Of[TestContext](20)
		eff3 := Of[TestContext](30)

		combined := effectMonoid.Concat(eff1, effectMonoid.Concat(eff2, eff3))
		result, err := runEffect(combined, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, 60, result)
	})

	t.Run("returns empty value for empty monoid", func(t *testing.T) {
		stringMonoid := monoid.MakeMonoid(
			func(a, b string) string { return a + b },
			"empty",
		)

		effectMonoid := ApplicativeMonoid[TestContext](stringMonoid)

		result, err := runEffect(effectMonoid.Empty(), TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "empty", result)
	})

	t.Run("propagates first error", func(t *testing.T) {
		expectedErr := errors.New("first error")
		stringMonoid := monoid.MakeMonoid(
			func(a, b string) string { return a + b },
			"",
		)

		effectMonoid := ApplicativeMonoid[TestContext](stringMonoid)

		eff1 := Fail[TestContext, string](expectedErr)
		eff2 := Of[TestContext]("World")

		combined := effectMonoid.Concat(eff1, eff2)
		_, err := runEffect(combined, TestContext{Value: "test"})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("propagates second error", func(t *testing.T) {
		expectedErr := errors.New("second error")
		stringMonoid := monoid.MakeMonoid(
			func(a, b string) string { return a + b },
			"",
		)

		effectMonoid := ApplicativeMonoid[TestContext](stringMonoid)

		eff1 := Of[TestContext]("Hello")
		eff2 := Fail[TestContext, string](expectedErr)

		combined := effectMonoid.Concat(eff1, eff2)
		_, err := runEffect(combined, TestContext{Value: "test"})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("combines multiple effects", func(t *testing.T) {
		intMonoid := monoid.MakeMonoid(
			func(a, b int) int { return a * b },
			1,
		)

		effectMonoid := ApplicativeMonoid[TestContext](intMonoid)

		effects := []Effect[TestContext, int]{
			Of[TestContext](2),
			Of[TestContext](3),
			Of[TestContext](4),
			Of[TestContext](5),
		}

		combined := effectMonoid.Empty()
		for _, eff := range effects {
			combined = effectMonoid.Concat(combined, eff)
		}

		result, err := runEffect(combined, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, 120, result) // 1 * 2 * 3 * 4 * 5
	})

	t.Run("works with custom types", func(t *testing.T) {
		type Counter struct {
			Count int
		}

		counterMonoid := monoid.MakeMonoid(
			func(a, b Counter) Counter {
				return Counter{Count: a.Count + b.Count}
			},
			Counter{Count: 0},
		)

		effectMonoid := ApplicativeMonoid[TestContext](counterMonoid)

		eff1 := Of[TestContext](Counter{Count: 5})
		eff2 := Of[TestContext](Counter{Count: 10})
		eff3 := Of[TestContext](Counter{Count: 15})

		combined := effectMonoid.Concat(eff1, effectMonoid.Concat(eff2, eff3))
		result, err := runEffect(combined, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, 30, result.Count)
	})
}

func TestAlternativeMonoid(t *testing.T) {
	t.Run("combines successful effects with monoid", func(t *testing.T) {
		stringMonoid := monoid.MakeMonoid(
			func(a, b string) string { return a + b },
			"",
		)

		effectMonoid := AlternativeMonoid[TestContext](stringMonoid)

		eff1 := Of[TestContext]("First")
		eff2 := Of[TestContext]("Second")

		combined := effectMonoid.Concat(eff1, eff2)
		result, err := runEffect(combined, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "FirstSecond", result) // Alternative still combines when both succeed
	})

	t.Run("tries second effect if first fails", func(t *testing.T) {
		stringMonoid := monoid.MakeMonoid(
			func(a, b string) string { return a + b },
			"",
		)

		effectMonoid := AlternativeMonoid[TestContext](stringMonoid)

		eff1 := Fail[TestContext, string](errors.New("first failed"))
		eff2 := Of[TestContext]("Second")

		combined := effectMonoid.Concat(eff1, eff2)
		result, err := runEffect(combined, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "Second", result)
	})

	t.Run("returns error if all effects fail", func(t *testing.T) {
		expectedErr := errors.New("second error")
		stringMonoid := monoid.MakeMonoid(
			func(a, b string) string { return a + b },
			"",
		)

		effectMonoid := AlternativeMonoid[TestContext](stringMonoid)

		eff1 := Fail[TestContext, string](errors.New("first error"))
		eff2 := Fail[TestContext, string](expectedErr)

		combined := effectMonoid.Concat(eff1, eff2)
		_, err := runEffect(combined, TestContext{Value: "test"})

		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)
	})

	t.Run("returns empty value for empty monoid", func(t *testing.T) {
		stringMonoid := monoid.MakeMonoid(
			func(a, b string) string { return a + b },
			"default",
		)

		effectMonoid := AlternativeMonoid[TestContext](stringMonoid)

		result, err := runEffect(effectMonoid.Empty(), TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "default", result)
	})

	t.Run("chains multiple alternatives", func(t *testing.T) {
		intMonoid := monoid.MakeMonoid(
			func(a, b int) int { return a + b },
			0,
		)

		effectMonoid := AlternativeMonoid[TestContext](intMonoid)

		eff1 := Fail[TestContext, int](errors.New("error 1"))
		eff2 := Fail[TestContext, int](errors.New("error 2"))
		eff3 := Of[TestContext](42)
		eff4 := Of[TestContext](100)

		combined := effectMonoid.Concat(
			effectMonoid.Concat(eff1, eff2),
			effectMonoid.Concat(eff3, eff4),
		)

		result, err := runEffect(combined, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, 142, result) // Combines successful values: 42 + 100
	})

	t.Run("works with custom types", func(t *testing.T) {
		type Result struct {
			Value string
			Code  int
		}

		resultMonoid := monoid.MakeMonoid(
			func(a, b Result) Result {
				return Result{Value: a.Value + b.Value, Code: a.Code + b.Code}
			},
			Result{Value: "", Code: 0},
		)

		effectMonoid := AlternativeMonoid[TestContext](resultMonoid)

		eff1 := Fail[TestContext, Result](errors.New("failed"))
		eff2 := Of[TestContext](Result{Value: "success", Code: 200})
		eff3 := Of[TestContext](Result{Value: "backup", Code: 201})

		combined := effectMonoid.Concat(effectMonoid.Concat(eff1, eff2), eff3)
		result, err := runEffect(combined, TestContext{Value: "test"})

		assert.NoError(t, err)
		assert.Equal(t, "successbackup", result.Value) // Combines both successful values
		assert.Equal(t, 401, result.Code)              // 200 + 201
	})
}

func TestMonoidComparison(t *testing.T) {
	t.Run("ApplicativeMonoid vs AlternativeMonoid with all success", func(t *testing.T) {
		stringMonoid := monoid.MakeMonoid(
			func(a, b string) string { return a + "," + b },
			"",
		)

		applicativeMonoid := ApplicativeMonoid[TestContext](stringMonoid)
		alternativeMonoid := AlternativeMonoid[TestContext](stringMonoid)

		eff1 := Of[TestContext]("A")
		eff2 := Of[TestContext]("B")

		// Applicative combines values
		applicativeResult, err1 := runEffect(
			applicativeMonoid.Concat(eff1, eff2),
			TestContext{Value: "test"},
		)

		// Alternative takes first
		alternativeResult, err2 := runEffect(
			alternativeMonoid.Concat(eff1, eff2),
			TestContext{Value: "test"},
		)

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, "A,B", applicativeResult) // Combined with comma separator
		assert.Equal(t, "A,B", alternativeResult) // Also combined (Alternative uses Alt semigroup)
	})

	t.Run("ApplicativeMonoid vs AlternativeMonoid with failures", func(t *testing.T) {
		intMonoid := monoid.MakeMonoid(
			func(a, b int) int { return a + b },
			0,
		)

		applicativeMonoid := ApplicativeMonoid[TestContext](intMonoid)
		alternativeMonoid := AlternativeMonoid[TestContext](intMonoid)

		eff1 := Fail[TestContext, int](errors.New("error 1"))
		eff2 := Of[TestContext](42)

		// Applicative fails on first error
		_, err1 := runEffect(
			applicativeMonoid.Concat(eff1, eff2),
			TestContext{Value: "test"},
		)

		// Alternative tries second on first failure
		result2, err2 := runEffect(
			alternativeMonoid.Concat(eff1, eff2),
			TestContext{Value: "test"},
		)

		assert.Error(t, err1)
		assert.NoError(t, err2)
		assert.Equal(t, 42, result2)
	})
}
