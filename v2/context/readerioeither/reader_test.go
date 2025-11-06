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

package readerioeither

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
	outer := context.Background()

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
	outer := context.Background()

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
	outer := context.Background()

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
	outer := context.Background()

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
	ctx := context.Background()

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
		Ap[string, string](errValue),
	)

	res := applied(context.Background())()
	assert.Equal(t, E.Left[string](err), res)
}

func TestRegularApply(t *testing.T) {
	value := Of("Carsten")
	fct := Of(utils.Upper)

	applied := F.Pipe1(
		fct,
		Ap[string, string](value),
	)

	res := applied(context.Background())()
	assert.Equal(t, E.Of[error]("CARSTEN"), res)
}

func TestWithResourceNoErrors(t *testing.T) {
	var countAcquire, countBody, countRelease int

	acquire := FromLazy(func() int {
		countAcquire++
		return countAcquire
	})

	release := func(int) ReaderIOEither[int] {
		return FromLazy(func() int {
			countRelease++
			return countRelease
		})
	}

	body := func(int) ReaderIOEither[int] {
		return FromLazy(func() int {
			countBody++
			return countBody
		})
	}

	resRIOE := WithResource[int](acquire, release)(body)

	res := resRIOE(context.Background())()

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

	release := func(int) ReaderIOEither[int] {
		return FromLazy(func() int {
			countRelease++
			return countRelease
		})
	}

	err := fmt.Errorf("error in body")
	body := func(int) ReaderIOEither[int] {
		return Left[int](err)
	}

	resRIOE := WithResource[int](acquire, release)(body)

	res := resRIOE(context.Background())()

	assert.Equal(t, 1, countAcquire)
	assert.Equal(t, 0, countBody)
	assert.Equal(t, 1, countRelease)
	assert.Equal(t, E.Left[int](err), res)
}

func TestWithResourceErrorInAcquire(t *testing.T) {
	var countAcquire, countBody, countRelease int

	err := fmt.Errorf("error in acquire")
	acquire := Left[int](err)

	release := func(int) ReaderIOEither[int] {
		return FromLazy(func() int {
			countRelease++
			return countRelease
		})
	}

	body := func(int) ReaderIOEither[int] {
		return FromLazy(func() int {
			countBody++
			return countBody
		})
	}

	resRIOE := WithResource[int](acquire, release)(body)

	res := resRIOE(context.Background())()

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
	release := func(int) ReaderIOEither[int] {
		return Left[int](err)
	}

	body := func(int) ReaderIOEither[int] {
		return FromLazy(func() int {
			countBody++
			return countBody
		})
	}

	resRIOE := WithResource[int](acquire, release)(body)

	res := resRIOE(context.Background())()

	assert.Equal(t, 1, countAcquire)
	assert.Equal(t, 1, countBody)
	assert.Equal(t, 0, countRelease)
	assert.Equal(t, E.Left[int](err), res)
}
