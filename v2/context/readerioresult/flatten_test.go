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
	"errors"
	"fmt"
	"strconv"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioresult"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

func TestFlattenIO_Success(t *testing.T) {
	t.Run("runs the produced IO", func(t *testing.T) {
		calls := 0
		m := FlattenIO(Of(IO[int](func() int { calls++; return 42 })))
		assert.Equal(t, 0, calls)
		assert.Equal(t, result.Of(42), m(t.Context())())
		assert.Equal(t, 1, calls)
	})
}

func TestFlattenIO_Failure(t *testing.T) {
	t.Run("upstream failure skips the IO", func(t *testing.T) {
		err := errors.New("upstream")
		assert.Equal(t, result.Left[int](err), FlattenIO(Left[IO[int]](err))(t.Context())())
	})
}

func TestFlattenIOResult_Success(t *testing.T) {
	t.Run("runs the produced IOResult", func(t *testing.T) {
		assert.Equal(t, result.Of(42), FlattenIOResult(Of(ioresult.Of(42)))(t.Context())())
	})
}

func TestFlattenIOResult_Failure(t *testing.T) {
	err := errors.New("failed")

	t.Run("failed IOResult fails the computation", func(t *testing.T) {
		assert.Equal(t, result.Left[int](err), FlattenIOResult(Of(ioresult.Left[int](err)))(t.Context())())
	})

	t.Run("upstream failure", func(t *testing.T) {
		assert.Equal(t, result.Left[int](err), FlattenIOResult(Left[IOResult[int]](err))(t.Context())())
	})
}

// ExampleFlattenIO demonstrates running an IO action produced by a Map step.
func ExampleFlattenIO() {
	m := F.Pipe2(
		Of(21),
		Map(F.Flow2(N.Mul(2), io.Of[int])),
		FlattenIO[int],
	)
	fmt.Println(m(context.Background())())
	// Output: Right[int](42)
}

// ExampleFlattenIOResult demonstrates running an IOResult produced by a Map step.
func ExampleFlattenIOResult() {
	parse := F.Flow2(
		Map(ioresult.Eitherize1(strconv.Atoi)),
		FlattenIOResult[int],
	)
	fmt.Println(parse(Of("42"))(context.Background())())
	// Output: Right[int](42)
}
