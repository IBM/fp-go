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
	"fmt"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/lazy"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/stretchr/testify/assert"
)

func TestFlattenIO_Success(t *testing.T) {
	t.Run("runs the produced IO", func(t *testing.T) {
		calls := 0
		m := FlattenIO(Of[testContext, string](IO[int](func() int { calls++; return 42 })))
		assert.Equal(t, 0, calls)
		assert.Equal(t, E.Right[string](42), m(testContext{})())
		assert.Equal(t, 1, calls)
	})
}

func TestFlattenIO_Failure(t *testing.T) {
	t.Run("upstream failure skips the IO", func(t *testing.T) {
		assert.Equal(t, E.Left[int]("boom"), FlattenIO(Left[testContext, IO[int]]("boom"))(testContext{})())
	})
}

func TestFlattenIOEither_Success(t *testing.T) {
	t.Run("runs the produced IOEither", func(t *testing.T) {
		m := FlattenIOEither(Of[testContext, string](IOE.Of[string](42)))
		assert.Equal(t, E.Right[string](42), m(testContext{})())
	})
}

func TestFlattenIOEither_Failure(t *testing.T) {
	t.Run("failed IOEither fails the computation", func(t *testing.T) {
		m := FlattenIOEither(Of[testContext, string](IOE.Left[int]("inner")))
		assert.Equal(t, E.Left[int]("inner"), m(testContext{})())
	})

	t.Run("upstream failure", func(t *testing.T) {
		m := FlattenIOEither(Left[testContext, IOEither[string, int]]("outer"))
		assert.Equal(t, E.Left[int]("outer"), m(testContext{})())
	})
}

func TestFromLazy_Success(t *testing.T) {
	t.Run("evaluates lazily", func(t *testing.T) {
		calls := 0
		m := FromLazy[testContext, string](func() int { calls++; return 42 })
		assert.Equal(t, 0, calls)
		assert.Equal(t, E.Right[string](42), m(testContext{})())
		assert.Equal(t, 1, calls)
	})
}

// ExampleFlattenIO demonstrates running an IO action produced by a Map step.
func ExampleFlattenIO() {
	m := F.Pipe2(
		Of[testContext, error](21),
		Map[testContext, error](F.Flow2(N.Mul(2), io.Of[int])),
		FlattenIO[testContext, error, int],
	)
	fmt.Println(m(testContext{})())
	// Output: Right[int](42)
}

// ExampleFlattenIOEither demonstrates running an IOEither produced by a Map step.
func ExampleFlattenIOEither() {
	m := F.Pipe2(
		Of[testContext, string](21),
		Map[testContext, string](F.Flow2(N.Mul(2), IOE.Of[string, int])),
		FlattenIOEither[testContext, string, int],
	)
	fmt.Println(m(testContext{})())
	// Output: Right[int](42)
}

// ExampleFromLazy demonstrates lifting a lazy computation.
func ExampleFromLazy() {
	fmt.Println(FromLazy[testContext, error](lazy.Of(42))(testContext{})())
	// Output: Right[int](42)
}
