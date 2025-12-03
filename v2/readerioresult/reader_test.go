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
	"log"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	R "github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {

	g := F.Pipe1(
		Of[context.Context](1),
		Map[context.Context](utils.Double),
	)

	assert.Equal(t, result.Of(2), g(context.Background())())
}

func TestAp(t *testing.T) {
	g := F.Pipe1(
		Right[context.Context](utils.Double),
		Ap[int](Right[context.Context](1)),
	)

	assert.Equal(t, result.Of(2), g(context.Background())())
}

func TestChainReaderK(t *testing.T) {

	g := F.Pipe1(
		Of[context.Context](1),
		ChainReaderK(func(v int) R.Reader[context.Context, string] {
			return R.Of[context.Context](fmt.Sprintf("%d", v))
		}),
	)

	assert.Equal(t, result.Of("1"), g(context.Background())())
}

func TestTapReaderIOK(t *testing.T) {

	rdr := Of[int]("TestTapReaderIOK")

	x := F.Pipe1(
		rdr,
		TapReaderIOK(func(a string) ReaderIO[int, any] {
			return func(ctx int) IO[any] {
				return func() any {
					log.Printf("Context: %d, Value: %s", ctx, a)
					return nil
				}
			}
		}),
	)

	x(10)()
}
