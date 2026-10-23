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

package statereaderioeither

import (
	"errors"
	"fmt"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/io"
	P "github.com/IBM/fp-go/v2/pair"
	RIOE "github.com/IBM/fp-go/v2/readerioeither"
	"github.com/stretchr/testify/assert"
)

// incCounter increments the counter and returns the new value.
func incCounter(s testState) ReaderIOEither[testContext, error, Pair[testState, int]] {
	s.counter++
	return RIOE.Of[testContext, error](P.MakePair(s, s.counter))
}

func TestFlatten_Success(t *testing.T) {
	t.Run("runs the inner computation", func(t *testing.T) {
		res := Flatten(Of[testState, testContext, error](Of[testState, testContext, error](42)))(testState{counter: 1})(testContext{})()
		assert.Equal(t, E.Right[error](P.MakePair(testState{counter: 1}, 42)), res)
	})

	t.Run("threads the state from the outer to the inner computation", func(t *testing.T) {
		inc := StateReaderIOEither[testState, testContext, error, int](incCounter)
		outer := MonadMap(inc, func(int) StateReaderIOEither[testState, testContext, error, int] { return inc })
		res := Flatten(outer)(testState{counter: 0})(testContext{})()
		assert.Equal(t, E.Right[error](P.MakePair(testState{counter: 2}, 2)), res)
	})
}

func TestFlatten_Failure(t *testing.T) {
	err := errors.New("failed")

	t.Run("outer failure", func(t *testing.T) {
		res := Flatten(Left[testState, testContext, StateReaderIOEither[testState, testContext, error, int]](err))(testState{})(testContext{})()
		assert.Equal(t, E.Left[Pair[testState, int]](err), res)
	})

	t.Run("inner failure", func(t *testing.T) {
		res := Flatten(Of[testState, testContext, error](Left[testState, testContext, int](err)))(testState{})(testContext{})()
		assert.Equal(t, E.Left[Pair[testState, int]](err), res)
	})
}

// ExampleFlatten demonstrates removing one level of nesting from a StateReaderIOEither.
func ExampleFlatten() {
	nested := Of[testState, testContext, error](Of[testState, testContext, error](42))
	fmt.Println(E.Unwrap(Flatten(nested)(testState{counter: 1})(testContext{})()))
	// Output: Pair[statereaderioeither.testState, int]({1}, 42) <nil>
}

func TestFromReaderIO_Success(t *testing.T) {
	t.Run("reads the context and keeps the state", func(t *testing.T) {
		m := FromReaderIO[testState, error](func(c testContext) IO[int] { return io.Of(c.multiplier) })
		assert.Equal(t, E.Right[error](P.MakePair(testState{counter: 1}, 3)), m(testState{counter: 1})(testContext{multiplier: 3})())
	})
}

// ExampleFromReaderIO demonstrates lifting a ReaderIO into a StateReaderIOEither.
func ExampleFromReaderIO() {
	m := FromReaderIO[testState, error](func(c testContext) IO[int] { return io.Of(c.multiplier) })
	fmt.Println(E.Unwrap(m(testState{counter: 1})(testContext{multiplier: 3})()))
	// Output: Pair[statereaderioeither.testState, int]({1}, 3) <nil>
}
