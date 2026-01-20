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

package readerioresult

import (
	"context"
	"fmt"
	"testing"
	"time"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	"github.com/stretchr/testify/assert"
)

func TestInnerContextCancelSemantics(t *testing.T) {
	// start with a simple context
	outer := t.Context()

	parent, parentCancel := context.WithCancel(outer)
	defer parentCancel()

	inner, innerCancel := context.WithCancel(parent)
	defer innerCancel()

	assert.NoError(t, parent.Err())
	assert.NoError(t, inner.Err())

	innerCancel()

	assert.NoError(t, parent.Err())
	assert.Error(t, inner.Err())

}

func TestOuterContextCancelSemantics(t *testing.T) {
	// start with a simple context
	outer := t.Context()

	parent, outerCancel := context.WithCancel(outer)
	defer outerCancel()

	inner, innerCancel := context.WithCancel(parent)
	defer innerCancel()

	assert.NoError(t, parent.Err())
	assert.NoError(t, inner.Err())

	outerCancel()

	assert.Error(t, parent.Err())
	assert.Error(t, inner.Err())

}

func TestOuterAndInnerContextCancelSemantics(t *testing.T) {
	// start with a simple context
	outer := t.Context()

	parent, outerCancel := context.WithCancel(outer)
	defer outerCancel()

	inner, innerCancel := context.WithCancel(parent)
	defer innerCancel()

	assert.NoError(t, parent.Err())
	assert.NoError(t, inner.Err())

	outerCancel()
	innerCancel()

	assert.Error(t, parent.Err())
	assert.Error(t, inner.Err())

	outerCancel()
	innerCancel()

	assert.Error(t, parent.Err())
	assert.Error(t, inner.Err())
}

func TestCancelCauseSemantics(t *testing.T) {
	// start with a simple context
	outer := t.Context()

	parent, outerCancel := context.WithCancelCause(outer)
	defer outerCancel(nil)

	inner := context.WithValue(parent, "key", "value")

	assert.NoError(t, parent.Err())
	assert.NoError(t, inner.Err())

	err := fmt.Errorf("test error")

	outerCancel(err)

	assert.Error(t, parent.Err())
	assert.Error(t, inner.Err())

	assert.Equal(t, err, context.Cause(parent))
	assert.Equal(t, err, context.Cause(inner))
}

func TestTimer(t *testing.T) {
	delta := 3 * time.Second
	timer := Timer(delta)
	ctx := t.Context()

	t0 := time.Now()
	res := timer(ctx)()
	t1 := time.Now()

	assert.WithinDuration(t, t0.Add(delta), t1, time.Second)
	assert.True(t, E.IsRight(res))
}

func TestCanceledApply(t *testing.T) {
	// our error
	err := fmt.Errorf("TestCanceledApply")
	// the actual apply value errors out after some time
	errValue := F.Pipe1(
		Left[string](err),
		Delay[string](time.Second),
	)
	// function never resolves
	fct := Never[func(string) string]()
	// apply the values, we expect an error after 1s

	applied := F.Pipe1(
		fct,
		Ap[string](errValue),
	)

	res := applied(t.Context())()
	assert.Equal(t, E.Left[string](err), res)
}

func TestRegularApply(t *testing.T) {
	value := Of("Carsten")
	fct := Of(utils.Upper)

	applied := F.Pipe1(
		fct,
		Ap[string](value),
	)

	res := applied(t.Context())()
	assert.Equal(t, E.Of[error]("CARSTEN"), res)
}

func TestWithResourceNoErrors(t *testing.T) {
	var countAcquire, countBody, countRelease int

	acquire := FromLazy(func() int {
		countAcquire++
		return countAcquire
	})

	release := func(int) ReaderIOResult[int] {
		return FromLazy(func() int {
			countRelease++
			return countRelease
		})
	}

	body := func(int) ReaderIOResult[int] {
		return FromLazy(func() int {
			countBody++
			return countBody
		})
	}

	resRIOE := WithResource[int](acquire, release)(body)

	res := resRIOE(t.Context())()

	assert.Equal(t, 1, countAcquire)
	assert.Equal(t, 1, countBody)
	assert.Equal(t, 1, countRelease)
	assert.Equal(t, E.Of[error](1), res)
}

func TestWithResourceErrorInBody(t *testing.T) {
	var countAcquire, countBody, countRelease int

	acquire := FromLazy(func() int {
		countAcquire++
		return countAcquire
	})

	release := func(int) ReaderIOResult[int] {
		return FromLazy(func() int {
			countRelease++
			return countRelease
		})
	}

	err := fmt.Errorf("error in body")
	body := func(int) ReaderIOResult[int] {
		return Left[int](err)
	}

	resRIOE := WithResource[int](acquire, release)(body)

	res := resRIOE(t.Context())()

	assert.Equal(t, 1, countAcquire)
	assert.Equal(t, 0, countBody)
	assert.Equal(t, 1, countRelease)
	assert.Equal(t, E.Left[int](err), res)
}

func TestWithResourceErrorInAcquire(t *testing.T) {
	var countAcquire, countBody, countRelease int

	err := fmt.Errorf("error in acquire")
	acquire := Left[int](err)

	release := func(int) ReaderIOResult[int] {
		return FromLazy(func() int {
			countRelease++
			return countRelease
		})
	}

	body := func(int) ReaderIOResult[int] {
		return FromLazy(func() int {
			countBody++
			return countBody
		})
	}

	resRIOE := WithResource[int](acquire, release)(body)

	res := resRIOE(t.Context())()

	assert.Equal(t, 0, countAcquire)
	assert.Equal(t, 0, countBody)
	assert.Equal(t, 0, countRelease)
	assert.Equal(t, E.Left[int](err), res)
}

func TestWithResourceErrorInRelease(t *testing.T) {
	var countAcquire, countBody, countRelease int

	acquire := FromLazy(func() int {
		countAcquire++
		return countAcquire
	})

	err := fmt.Errorf("error in release")
	release := func(int) ReaderIOResult[int] {
		return Left[int](err)
	}

	body := func(int) ReaderIOResult[int] {
		return FromLazy(func() int {
			countBody++
			return countBody
		})
	}

	resRIOE := WithResource[int](acquire, release)(body)

	res := resRIOE(t.Context())()

	assert.Equal(t, 1, countAcquire)
	assert.Equal(t, 1, countBody)
	assert.Equal(t, 0, countRelease)
	assert.Equal(t, E.Left[int](err), res)
}

func TestMonadChainFirstLeft(t *testing.T) {
	ctx := t.Context()

	// Test with Left value - function returns Left, always preserves original error
	t.Run("Left value with function returning Left preserves original error", func(t *testing.T) {
		sideEffectCalled := false
		originalErr := fmt.Errorf("original error")
		result := MonadChainFirstLeft(
			Left[int](originalErr),
			func(e error) ReaderIOResult[int] {
				sideEffectCalled = true
				return Left[int](fmt.Errorf("new error")) // This error is ignored
			},
		)
		actualResult := result(ctx)()
		assert.True(t, sideEffectCalled)
		assert.Equal(t, E.Left[int](originalErr), actualResult)
	})

	// Test with Left value - function returns Right, still returns original Left
	t.Run("Left value with function returning Right still returns original Left", func(t *testing.T) {
		var capturedError error
		originalErr := fmt.Errorf("validation failed")
		result := MonadChainFirstLeft(
			Left[int](originalErr),
			func(e error) ReaderIOResult[int] {
				capturedError = e
				return Right(999) // This Right value is ignored
			},
		)
		actualResult := result(ctx)()
		assert.Equal(t, originalErr, capturedError)
		assert.Equal(t, E.Left[int](originalErr), actualResult)
	})

	// Test with Right value - should pass through without calling function
	t.Run("Right value passes through", func(t *testing.T) {
		sideEffectCalled := false
		result := MonadChainFirstLeft(
			Right(42),
			func(e error) ReaderIOResult[int] {
				sideEffectCalled = true
				return Left[int](fmt.Errorf("should not be called"))
			},
		)
		assert.False(t, sideEffectCalled)
		assert.Equal(t, E.Right[error](42), result(ctx)())
	})

	// Test that side effects are executed but original error is always preserved
	t.Run("Side effects executed but original error preserved", func(t *testing.T) {
		effectCount := 0
		originalErr := fmt.Errorf("original error")
		result := MonadChainFirstLeft(
			Left[int](originalErr),
			func(e error) ReaderIOResult[int] {
				effectCount++
				// Try to return Right, but original Left should still be returned
				return Right(999)
			},
		)
		actualResult := result(ctx)()
		assert.Equal(t, 1, effectCount)
		assert.Equal(t, E.Left[int](originalErr), actualResult)
	})
}

func TestChainFirstLeft(t *testing.T) {
	ctx := t.Context()

	// Test with Left value - function returns Left, always preserves original error
	t.Run("Left value with function returning Left preserves error", func(t *testing.T) {
		var captured error
		originalErr := fmt.Errorf("test error")
		chainFn := ChainFirstLeft[int](func(e error) ReaderIOResult[int] {
			captured = e
			return Left[int](fmt.Errorf("ignored error"))
		})
		result := F.Pipe1(
			Left[int](originalErr),
			chainFn,
		)
		actualResult := result(ctx)()
		assert.Equal(t, originalErr, captured)
		assert.Equal(t, E.Left[int](originalErr), actualResult)
	})

	// Test with Left value - function returns Right, still returns original Left
	t.Run("Left value with function returning Right still returns original Left", func(t *testing.T) {
		var captured error
		originalErr := fmt.Errorf("test error")
		chainFn := ChainFirstLeft[int](func(e error) ReaderIOResult[int] {
			captured = e
			return Right(42) // This Right is ignored
		})
		result := F.Pipe1(
			Left[int](originalErr),
			chainFn,
		)
		actualResult := result(ctx)()
		assert.Equal(t, originalErr, captured)
		assert.Equal(t, E.Left[int](originalErr), actualResult)
	})

	// Test with Right value - should pass through without calling function
	t.Run("Right value passes through", func(t *testing.T) {
		called := false
		chainFn := ChainFirstLeft[int](func(e error) ReaderIOResult[int] {
			called = true
			return Right(0)
		})
		result := F.Pipe1(
			Right(100),
			chainFn,
		)
		assert.False(t, called)
		assert.Equal(t, E.Right[error](100), result(ctx)())
	})

	// Test that original error is always preserved regardless of what f returns
	t.Run("Original error always preserved", func(t *testing.T) {
		originalErr := fmt.Errorf("original")
		chainFn := ChainFirstLeft[int](func(e error) ReaderIOResult[int] {
			// Try to return Right, but original Left should still be returned
			return Right(999)
		})

		result := F.Pipe1(
			Left[int](originalErr),
			chainFn,
		)
		assert.Equal(t, E.Left[int](originalErr), result(ctx)())
	})

	// Test logging with Left preservation
	t.Run("Logging with Left preservation", func(t *testing.T) {
		errorLog := []string{}
		originalErr := fmt.Errorf("step1")
		logError := ChainFirstLeft[string](func(e error) ReaderIOResult[string] {
			errorLog = append(errorLog, "Logged: "+e.Error())
			return Left[string](fmt.Errorf("log entry")) // This is ignored
		})

		result := F.Pipe2(
			Left[string](originalErr),
			logError,
			ChainLeft(func(e error) ReaderIOResult[string] {
				return Left[string](fmt.Errorf("step2"))
			}),
		)

		actualResult := result(ctx)()
		assert.Equal(t, []string{"Logged: step1"}, errorLog)
		assert.Equal(t, E.Left[string](fmt.Errorf("step2")), actualResult)
	})
}
