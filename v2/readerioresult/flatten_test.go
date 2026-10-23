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
	"errors"
	"fmt"
	"strconv"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/lazy"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

// fromCtxValue reads the value of the test context.
func fromCtxValue(c testContext) int { return c.value }

func TestFlattenIO_Success(t *testing.T) {
	t.Run("runs the produced IO", func(t *testing.T) {
		calls := 0
		m := FlattenIO(Of[testContext](IO[int](func() int { calls++; return 42 })))
		assert.Equal(t, 0, calls)
		assert.Equal(t, result.Of(42), m(testContext{})())
		assert.Equal(t, 1, calls)
	})
}

func TestFlattenIO_Failure(t *testing.T) {
	t.Run("upstream failure skips the IO", func(t *testing.T) {
		err := errors.New("upstream")
		assert.Equal(t, result.Left[int](err), FlattenIO(Left[testContext, IO[int]](err))(testContext{})())
	})
}

func TestFlattenIOResult_Success(t *testing.T) {
	t.Run("runs the produced IOResult", func(t *testing.T) {
		assert.Equal(t, result.Of(42), FlattenIOResult(Of[testContext](ioresult.Of(42)))(testContext{})())
	})
}

func TestFlattenIOResult_Failure(t *testing.T) {
	err := errors.New("failed")

	t.Run("failed IOResult fails the computation", func(t *testing.T) {
		assert.Equal(t, result.Left[int](err), FlattenIOResult(Of[testContext](ioresult.Left[int](err)))(testContext{})())
	})

	t.Run("upstream failure", func(t *testing.T) {
		assert.Equal(t, result.Left[int](err), FlattenIOResult(Left[testContext, IOResult[int]](err))(testContext{})())
	})
}

func TestFromReaderResult_Success(t *testing.T) {
	t.Run("lifts a successful ReaderResult", func(t *testing.T) {
		m := FromReaderResult(F.Flow2(fromCtxValue, result.Of[int]))
		assert.Equal(t, result.Of(10), m(testContext{value: 10})())
	})
}

func TestFromReaderResult_Failure(t *testing.T) {
	t.Run("lifts a failed ReaderResult", func(t *testing.T) {
		err := errors.New("failed")
		m := FromReaderResult(F.Constant1[testContext](result.Left[int](err)))
		assert.Equal(t, result.Left[int](err), m(testContext{})())
	})
}

func TestFromLazy_Success(t *testing.T) {
	t.Run("evaluates lazily", func(t *testing.T) {
		calls := 0
		m := FromLazy[testContext](func() int { calls++; return 42 })
		assert.Equal(t, 0, calls)
		assert.Equal(t, result.Of(42), m(testContext{})())
		assert.Equal(t, 1, calls)
	})
}

// ExampleFlattenIO demonstrates running an IO action produced by a Map step.
func ExampleFlattenIO() {
	m := F.Pipe2(
		Of[testContext](21),
		Map[testContext](F.Flow2(N.Mul(2), io.Of[int])),
		FlattenIO[testContext, int],
	)
	fmt.Println(m(testContext{})())
	// Output: Right[int](42)
}

// ExampleFlattenIOResult demonstrates running an IOResult produced by a Map step.
func ExampleFlattenIOResult() {
	parse := F.Flow2(
		Map[testContext](ioresult.Eitherize1(strconv.Atoi)),
		FlattenIOResult[testContext, int],
	)
	fmt.Println(parse(Of[testContext]("42"))(testContext{})())
	// Output: Right[int](42)
}

// ExampleFromReaderResult demonstrates lifting a ReaderResult.
func ExampleFromReaderResult() {
	m := FromReaderResult(F.Flow2(fromCtxValue, result.Of[int]))
	fmt.Println(m(testContext{value: 7})())
	// Output: Right[int](7)
}

// ExampleFromLazy demonstrates lifting a lazy computation.
func ExampleFromLazy() {
	fmt.Println(FromLazy[testContext](lazy.Of(42))(testContext{})())
	// Output: Right[int](42)
}
