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

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	R "github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/readerio"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {

	g := F.Pipe1(
		Of[context.Context, error](1),
		Map[context.Context, error](utils.Double),
	)

	assert.Equal(t, E.Of[error](2), g(context.Background())())
}

func TestOrLeft(t *testing.T) {
	f := OrLeft[int](func(s string) readerio.ReaderIO[context.Context, string] {
		return readerio.Of[context.Context](s + "!")
	})

	g1 := F.Pipe1(
		Right[context.Context, string](1),
		f,
	)

	g2 := F.Pipe1(
		Left[context.Context, int]("a"),
		f,
	)

	assert.Equal(t, E.Of[string](1), g1(context.Background())())
	assert.Equal(t, E.Left[int]("a!"), g2(context.Background())())
}

func TestAp(t *testing.T) {
	g := F.Pipe1(
		Right[context.Context, error](utils.Double),
		Ap[int](Right[context.Context, error](1)),
	)

	assert.Equal(t, E.Right[error](2), g(context.Background())())
}

func TestChainReaderK(t *testing.T) {

	g := F.Pipe1(
		Of[context.Context, error](1),
		ChainReaderK[error](func(v int) R.Reader[context.Context, string] {
			return R.Of[context.Context](fmt.Sprintf("%d", v))
		}),
	)

	assert.Equal(t, E.Right[error]("1"), g(context.Background())())
}
