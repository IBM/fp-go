// Copyright (c) 2023 IBM Corp.
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

	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	"github.com/IBM/fp-go/internal/utils"
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
