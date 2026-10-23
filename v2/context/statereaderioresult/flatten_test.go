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

package statereaderioresult

import (
	"context"
	"errors"
	"fmt"
	"testing"

	RIORES "github.com/IBM/fp-go/v2/context/readerioresult"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	P "github.com/IBM/fp-go/v2/pair"
	RES "github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// incCounter increments the counter and returns the new value.
func incCounter(s testState) ReaderIOResult[Pair[testState, int]] {
	s.counter++
	return RIORES.Of(P.MakePair(s, s.counter))
}

func TestFlatten_Success(t *testing.T) {
	t.Run("runs the inner computation", func(t *testing.T) {
		res := Flatten(Of[testState](Of[testState](42)))(testState{counter: 1})(t.Context())()
		assert.Equal(t, RES.Of(P.MakePair(testState{counter: 1}, 42)), res)
	})

	t.Run("threads the state from the outer to the inner computation", func(t *testing.T) {
		inc := StateReaderIOResult[testState, int](incCounter)
		outer := MonadMap(inc, func(int) StateReaderIOResult[testState, int] { return inc })
		res := Flatten(outer)(testState{counter: 0})(t.Context())()
		assert.Equal(t, RES.Of(P.MakePair(testState{counter: 2}, 2)), res)
	})
}

func TestFlatten_Failure(t *testing.T) {
	err := errors.New("failed")

	t.Run("outer failure", func(t *testing.T) {
		res := Flatten(Left[testState, StateReaderIOResult[testState, int]](err))(testState{})(t.Context())()
		assert.Equal(t, RES.Left[Pair[testState, int]](err), res)
	})

	t.Run("inner failure", func(t *testing.T) {
		res := Flatten(Of[testState](Left[testState, int](err)))(testState{})(t.Context())()
		assert.Equal(t, RES.Left[Pair[testState, int]](err), res)
	})
}

// ExampleFlatten demonstrates removing one level of nesting from a StateReaderIOResult.
func ExampleFlatten() {
	nested := Of[testState](Of[testState](42))
	fmt.Println(RES.Unwrap(Flatten(nested)(testState{counter: 1})(context.Background())()))
	// Output: Pair[statereaderioresult.testState, int]({1}, 42) <nil>
}

type fromCtxKey struct{}

// ctxValue reads the test value from the context.
func ctxValue(ctx context.Context) string { return ctx.Value(fromCtxKey{}).(string) }

func TestFromReader_Success(t *testing.T) {
	t.Run("reads the context and keeps the state", func(t *testing.T) {
		ctx := context.WithValue(t.Context(), fromCtxKey{}, "v")
		res := FromReader[testState](ctxValue)(testState{counter: 1})(ctx)()
		assert.Equal(t, RES.Of(P.MakePair(testState{counter: 1}, "v")), res)
	})
}

func TestFromReaderIO_Success(t *testing.T) {
	t.Run("reads the context and keeps the state", func(t *testing.T) {
		ctx := context.WithValue(t.Context(), fromCtxKey{}, "v")
		res := FromReaderIO[testState](F.Flow2(ctxValue, io.Of[string]))(testState{counter: 1})(ctx)()
		assert.Equal(t, RES.Of(P.MakePair(testState{counter: 1}, "v")), res)
	})
}

func TestFromReaderResult_Success(t *testing.T) {
	t.Run("lifts a successful ReaderResult", func(t *testing.T) {
		ctx := context.WithValue(t.Context(), fromCtxKey{}, "v")
		res := FromReaderResult[testState](F.Flow2(ctxValue, RES.Of[string]))(testState{counter: 1})(ctx)()
		assert.Equal(t, RES.Of(P.MakePair(testState{counter: 1}, "v")), res)
	})
}

func TestFromReaderResult_Failure(t *testing.T) {
	t.Run("lifts a failed ReaderResult", func(t *testing.T) {
		err := errors.New("failed")
		res := FromReaderResult[testState](F.Constant1[context.Context](RES.Left[string](err)))(testState{})(t.Context())()
		assert.Equal(t, RES.Left[Pair[testState, string]](err), res)
	})
}

// ExampleFromReader demonstrates lifting a Reader on context.Context.
func ExampleFromReader() {
	ctx := context.WithValue(context.Background(), fromCtxKey{}, "v")
	fmt.Println(RES.Unwrap(FromReader[testState](ctxValue)(testState{counter: 1})(ctx)()))
	// Output: Pair[statereaderioresult.testState, string]({1}, v) <nil>
}

// ExampleFromReaderIO demonstrates lifting a ReaderIO on context.Context.
func ExampleFromReaderIO() {
	ctx := context.WithValue(context.Background(), fromCtxKey{}, "v")
	fmt.Println(RES.Unwrap(FromReaderIO[testState](F.Flow2(ctxValue, io.Of[string]))(testState{counter: 1})(ctx)()))
	// Output: Pair[statereaderioresult.testState, string]({1}, v) <nil>
}

// ExampleFromReaderResult demonstrates lifting a ReaderResult on context.Context.
func ExampleFromReaderResult() {
	ctx := context.WithValue(context.Background(), fromCtxKey{}, "v")
	fmt.Println(RES.Unwrap(FromReaderResult[testState](F.Flow2(ctxValue, RES.Of[string]))(testState{counter: 1})(ctx)()))
	// Output: Pair[statereaderioresult.testState, string]({1}, v) <nil>
}
