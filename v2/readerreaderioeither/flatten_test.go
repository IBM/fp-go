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

package readerreaderioeither

import (
	"fmt"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/stretchr/testify/assert"
)

func TestFlattenIO_Success(t *testing.T) {
	t.Run("runs the produced IO", func(t *testing.T) {
		calls := 0
		m := FlattenIO(Of[exCfg, exReq, string](IO[int](func() int { calls++; return 42 })))
		assert.Equal(t, 0, calls)
		assert.Equal(t, E.Right[string](42), exRun(m))
		assert.Equal(t, 1, calls)
	})
}

func TestFlattenIO_Failure(t *testing.T) {
	t.Run("upstream failure skips the IO", func(t *testing.T) {
		assert.Equal(t, E.Left[int]("boom"), exRun(FlattenIO(Left[exCfg, exReq, IO[int]]("boom"))))
	})
}

func TestFlattenIOEither_Success(t *testing.T) {
	t.Run("runs the produced IOEither", func(t *testing.T) {
		assert.Equal(t, E.Right[string](42), exRun(FlattenIOEither(Of[exCfg, exReq, string](IOE.Of[string](42)))))
	})
}

func TestFlattenIOEither_Failure(t *testing.T) {
	t.Run("failed IOEither fails the computation", func(t *testing.T) {
		assert.Equal(t, E.Left[int]("inner"), exRun(FlattenIOEither(Of[exCfg, exReq, string](IOE.Left[int]("inner")))))
	})

	t.Run("upstream failure", func(t *testing.T) {
		assert.Equal(t, E.Left[int]("outer"), exRun(FlattenIOEither(Left[exCfg, exReq, IOEither[string, int]]("outer"))))
	})
}

// ExampleFlattenIO demonstrates running an IO action produced by a Map step.
func ExampleFlattenIO() {
	fmt.Println(exRun(F.Pipe2(
		Of[exCfg, exReq, string](21),
		Map[exCfg, exReq, string](F.Flow2(N.Mul(2), io.Of[int])),
		FlattenIO[exCfg, exReq, string, int],
	)))
	// Output: Right[int](42)
}

// ExampleFlattenIOEither demonstrates running an IOEither produced by a Map step.
func ExampleFlattenIOEither() {
	fmt.Println(exRun(F.Pipe2(
		Of[exCfg, exReq, string](21),
		Map[exCfg, exReq, string](F.Flow2(N.Mul(2), IOE.Of[string, int])),
		FlattenIOEither[exCfg, exReq, string, int],
	)))
	// Output: Right[int](42)
}
