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

package readerresult

import (
	"context"
	"errors"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/tuple"
	"github.com/stretchr/testify/assert"
)

// TestSequenceT1 tests the SequenceT1 function
func TestSequenceT1(t *testing.T) {
	t.Run("wraps single success value in tuple", func(t *testing.T) {
		rr := Of(42)
		tupled := SequenceT1(rr)
		result := tupled(t.Context())

		assert.True(t, E.IsRight(result))
		val, _ := E.Unwrap(result)
		assert.Equal(t, 42, val.F1)
	})

	t.Run("preserves error", func(t *testing.T) {
		testErr := errors.New("test error")
		rr := Left[int](testErr)
		tupled := SequenceT1(rr)
		result := tupled(t.Context())

		assert.True(t, E.IsLeft(result))
		_, err := E.UnwrapError(result)
		assert.Equal(t, testErr, err)
	})

	t.Run("respects context cancellation", func(t *testing.T) {
		rr := func(ctx context.Context) E.Either[error, int] {
			if ctx.Err() != nil {
				return E.Left[int](ctx.Err())
			}
			return E.Of[error](42)
		}

		tupled := SequenceT1(rr)

		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		result := tupled(ctx)
		assert.True(t, E.IsLeft(result))
	})
}

// TestSequenceT2 tests the SequenceT2 function
func TestSequenceT2(t *testing.T) {
	t.Run("combines two success values into tuple", func(t *testing.T) {
		getName := Of("Alice")
		getAge := Of(30)

		combined := SequenceT2(getName, getAge)
		result := combined(t.Context())

		assert.True(t, E.IsRight(result))
		val, _ := E.Unwrap(result)
		assert.Equal(t, "Alice", val.F1)
		assert.Equal(t, 30, val.F2)
	})

	t.Run("fails if first ReaderResult fails", func(t *testing.T) {
		testErr := errors.New("name not found")
		getName := Left[string](testErr)
		getAge := Of(30)

		combined := SequenceT2(getName, getAge)
		result := combined(t.Context())

		assert.True(t, E.IsLeft(result))
		_, err := E.UnwrapError(result)
		assert.Equal(t, testErr, err)
	})

	t.Run("fails if second ReaderResult fails", func(t *testing.T) {
		testErr := errors.New("age not found")
		getName := Of("Alice")
		getAge := Left[int](testErr)

		combined := SequenceT2(getName, getAge)
		result := combined(t.Context())

		assert.True(t, E.IsLeft(result))
		_, err := E.UnwrapError(result)
		assert.Equal(t, testErr, err)
	})

	t.Run("executes both ReaderResults with same context", func(t *testing.T) {
		type ctxKey string
		ctx := context.WithValue(t.Context(), ctxKey("key"), "shared")

		getName := func(ctx context.Context) E.Either[error, string] {
			val := ctx.Value(ctxKey("key"))
			if val != nil {
				return E.Of[error](val.(string))
			}
			return E.Left[string](errors.New("key not found"))
		}

		getAge := func(ctx context.Context) E.Either[error, int] {
			val := ctx.Value(ctxKey("key"))
			if val != nil {
				return E.Of[error](len(val.(string)))
			}
			return E.Left[int](errors.New("key not found"))
		}

		combined := SequenceT2(getName, getAge)
		result := combined(ctx)

		assert.True(t, E.IsRight(result))
		val, _ := E.Unwrap(result)
		assert.Equal(t, "shared", val.F1)
		assert.Equal(t, 6, val.F2) // len("shared")
	})

	t.Run("executes all ReaderResults even if one fails", func(t *testing.T) {
		firstExecuted := false
		secondExecuted := false

		first := func(ctx context.Context) E.Either[error, int] {
			firstExecuted = true
			return E.Left[int](errors.New("first failed"))
		}

		second := func(ctx context.Context) E.Either[error, string] {
			secondExecuted = true
			return E.Of[error]("second")
		}

		combined := SequenceT2(first, second)
		result := combined(t.Context())

		assert.True(t, firstExecuted, "first should be executed")
		assert.True(t, secondExecuted, "second should be executed (applicative semantics)")
		assert.True(t, E.IsLeft(result))
	})
}

// TestSequenceT3 tests the SequenceT3 function
func TestSequenceT3(t *testing.T) {
	t.Run("combines three success values into tuple", func(t *testing.T) {
		getUserID := Of(123)
		getUserName := Of("Alice")
		getUserEmail := Of("alice@example.com")

		combined := SequenceT3(getUserID, getUserName, getUserEmail)
		result := combined(t.Context())

		assert.True(t, E.IsRight(result))
		val, _ := E.Unwrap(result)
		assert.Equal(t, 123, val.F1)
		assert.Equal(t, "Alice", val.F2)
		assert.Equal(t, "alice@example.com", val.F3)
	})

	t.Run("fails if any ReaderResult fails", func(t *testing.T) {
		testErr := errors.New("email not found")
		getUserID := Of(123)
		getUserName := Of("Alice")
		getUserEmail := Left[string](testErr)

		combined := SequenceT3(getUserID, getUserName, getUserEmail)
		result := combined(t.Context())

		assert.True(t, E.IsLeft(result))
		_, err := E.UnwrapError(result)
		assert.Equal(t, testErr, err)
	})

	t.Run("executes all ReaderResults even if one fails", func(t *testing.T) {
		firstExecuted := false
		secondExecuted := false
		thirdExecuted := false

		first := func(ctx context.Context) E.Either[error, int] {
			firstExecuted = true
			return E.Of[error](1)
		}

		second := func(ctx context.Context) E.Either[error, int] {
			secondExecuted = true
			return E.Left[int](errors.New("second failed"))
		}

		third := func(ctx context.Context) E.Either[error, int] {
			thirdExecuted = true
			return E.Of[error](3)
		}

		combined := SequenceT3(first, second, third)
		result := combined(t.Context())

		assert.True(t, firstExecuted, "first should be executed")
		assert.True(t, secondExecuted, "second should be executed")
		assert.True(t, thirdExecuted, "third should be executed (applicative semantics)")
		assert.True(t, E.IsLeft(result))
	})

	t.Run("respects context cancellation", func(t *testing.T) {
		getUserID := func(ctx context.Context) E.Either[error, int] {
			if ctx.Err() != nil {
				return E.Left[int](ctx.Err())
			}
			return E.Of[error](123)
		}

		getUserName := Of("Alice")
		getUserEmail := Of("alice@example.com")

		combined := SequenceT3(getUserID, getUserName, getUserEmail)

		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		result := combined(ctx)
		assert.True(t, E.IsLeft(result))
	})
}

// TestSequenceT4 tests the SequenceT4 function
func TestSequenceT4(t *testing.T) {
	t.Run("combines four success values into tuple", func(t *testing.T) {
		getID := Of(123)
		getName := Of("Alice")
		getEmail := Of("alice@example.com")
		getAge := Of(30)

		combined := SequenceT4(getID, getName, getEmail, getAge)
		result := combined(t.Context())

		assert.True(t, E.IsRight(result))
		val, _ := E.Unwrap(result)
		assert.Equal(t, 123, val.F1)
		assert.Equal(t, "Alice", val.F2)
		assert.Equal(t, "alice@example.com", val.F3)
		assert.Equal(t, 30, val.F4)
	})

	t.Run("fails if any ReaderResult fails", func(t *testing.T) {
		testErr := errors.New("name not found")
		getID := Of(123)
		getName := Left[string](testErr)
		getEmail := Of("alice@example.com")
		getAge := Of(30)

		combined := SequenceT4(getID, getName, getEmail, getAge)
		result := combined(t.Context())

		assert.True(t, E.IsLeft(result))
		_, err := E.UnwrapError(result)
		assert.Equal(t, testErr, err)
	})

	t.Run("executes all ReaderResults even if one fails", func(t *testing.T) {
		firstExecuted := false
		secondExecuted := false
		thirdExecuted := false
		fourthExecuted := false

		first := func(ctx context.Context) E.Either[error, int] {
			firstExecuted = true
			return E.Of[error](1)
		}

		second := func(ctx context.Context) E.Either[error, int] {
			secondExecuted = true
			return E.Left[int](errors.New("second failed"))
		}

		third := func(ctx context.Context) E.Either[error, int] {
			thirdExecuted = true
			return E.Of[error](3)
		}

		fourth := func(ctx context.Context) E.Either[error, int] {
			fourthExecuted = true
			return E.Of[error](4)
		}

		combined := SequenceT4(first, second, third, fourth)
		result := combined(t.Context())

		assert.True(t, firstExecuted, "first should be executed")
		assert.True(t, secondExecuted, "second should be executed")
		assert.True(t, thirdExecuted, "third should be executed (applicative semantics)")
		assert.True(t, fourthExecuted, "fourth should be executed (applicative semantics)")
		assert.True(t, E.IsLeft(result))
	})

	t.Run("can be used to build complex structures", func(t *testing.T) {
		type UserProfile struct {
			ID    int
			Name  string
			Email string
			Age   int
		}

		fetchUserData := SequenceT4(
			Of(123),
			Of("Alice"),
			Of("alice@example.com"),
			Of(30),
		)

		buildProfile := Map(func(t tuple.Tuple4[int, string, string, int]) UserProfile {
			return UserProfile{
				ID:    t.F1,
				Name:  t.F2,
				Email: t.F3,
				Age:   t.F4,
			}
		})

		userProfile := func(ctx context.Context) E.Either[error, UserProfile] {
			tupleResult := fetchUserData(ctx)
			if E.IsLeft(tupleResult) {
				_, err := E.UnwrapError(tupleResult)
				return E.Left[UserProfile](err)
			}
			tupleVal, _ := E.Unwrap(tupleResult)
			return buildProfile(Of(tupleVal))(ctx)
		}

		result := userProfile(t.Context())

		assert.True(t, E.IsRight(result))
		profile, _ := E.Unwrap(result)
		assert.Equal(t, 123, profile.ID)
		assert.Equal(t, "Alice", profile.Name)
		assert.Equal(t, "alice@example.com", profile.Email)
		assert.Equal(t, 30, profile.Age)
	})

	t.Run("executes all with same context", func(t *testing.T) {
		type ctxKey string
		ctx := context.WithValue(t.Context(), ctxKey("multiplier"), 2)

		getBase := func(ctx context.Context) E.Either[error, int] {
			return E.Of[error](10)
		}

		multiply := func(ctx context.Context) E.Either[error, int] {
			mult := ctx.Value(ctxKey("multiplier")).(int)
			return E.Of[error](mult)
		}

		getResult := func(ctx context.Context) E.Either[error, int] {
			mult := ctx.Value(ctxKey("multiplier")).(int)
			return E.Of[error](10 * mult)
		}

		getDescription := func(ctx context.Context) E.Either[error, string] {
			return E.Of[error]("calculated")
		}

		combined := SequenceT4(getBase, multiply, getResult, getDescription)
		result := combined(ctx)

		assert.True(t, E.IsRight(result))
		val, _ := E.Unwrap(result)
		assert.Equal(t, 10, val.F1)
		assert.Equal(t, 2, val.F2)
		assert.Equal(t, 20, val.F3)
		assert.Equal(t, "calculated", val.F4)
	})
}

// TestSequenceIntegration tests integration scenarios
func TestSequenceIntegration(t *testing.T) {
	t.Run("SequenceT2 with Map to transform tuple", func(t *testing.T) {
		getName := Of("Alice")
		getAge := Of(30)

		combined := SequenceT2(getName, getAge)
		formatted := Map(func(t tuple.Tuple2[string, int]) string {
			return t.F1 + " is " + string(rune(t.F2+48)) + " years old"
		})

		pipeline := func(ctx context.Context) E.Either[error, string] {
			tupleResult := combined(ctx)
			if E.IsLeft(tupleResult) {
				_, err := E.UnwrapError(tupleResult)
				return E.Left[string](err)
			}
			tupleVal, _ := E.Unwrap(tupleResult)
			return formatted(Of(tupleVal))(ctx)
		}

		result := pipeline(t.Context())
		assert.True(t, E.IsRight(result))
	})

	t.Run("SequenceT3 with Chain for dependent operations", func(t *testing.T) {
		getX := Of(10)
		getY := Of(20)
		getZ := Of(30)

		combined := SequenceT3(getX, getY, getZ)

		sumTuple := func(t tuple.Tuple3[int, int, int]) ReaderResult[int] {
			return Of(t.F1 + t.F2 + t.F3)
		}

		pipeline := func(ctx context.Context) E.Either[error, int] {
			tupleResult := combined(ctx)
			if E.IsLeft(tupleResult) {
				_, err := E.UnwrapError(tupleResult)
				return E.Left[int](err)
			}
			tupleVal, _ := E.Unwrap(tupleResult)
			return sumTuple(tupleVal)(ctx)
		}

		result := pipeline(t.Context())
		assert.True(t, E.IsRight(result))
		val, _ := E.Unwrap(result)
		assert.Equal(t, 60, val) // 10 + 20 + 30
	})

	t.Run("nested sequences", func(t *testing.T) {
		// Create two pairs
		pair1 := SequenceT2(Of(1), Of(2))
		pair2 := SequenceT2(Of(3), Of(4))

		// Combine the pairs
		combined := SequenceT2(pair1, pair2)

		result := combined(t.Context())
		assert.True(t, E.IsRight(result))

		val, _ := E.Unwrap(result)
		assert.Equal(t, 1, val.F1.F1)
		assert.Equal(t, 2, val.F1.F2)
		assert.Equal(t, 3, val.F2.F1)
		assert.Equal(t, 4, val.F2.F2)
	})
}
