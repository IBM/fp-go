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

package iooption

import (
	"fmt"
	"os"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	I "github.com/IBM/fp-go/v2/io"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	assert.Equal(t, O.Of(2), F.Pipe1(
		Of(1),
		Map(utils.Double),
	)())

}

func TestChainOptionK(t *testing.T) {
	f := ChainOptionK(func(n int) Option[int] {
		if n > 0 {
			return O.Of(n)
		}
		return O.None[int]()

	})
	assert.Equal(t, O.Of(1), f(Of(1))())
	assert.Equal(t, O.None[int](), f(Of(-1))())
	assert.Equal(t, O.None[int](), f(None[int]())())
}

func TestFromOption(t *testing.T) {
	f := FromOption[int]
	assert.Equal(t, O.Of(1), f(O.Some(1))())
	assert.Equal(t, O.None[int](), f(O.None[int]())())
}

func TestChainIOK(t *testing.T) {
	f := ChainIOK(func(n int) I.IO[string] {
		return func() string {
			return fmt.Sprintf("%d", n)
		}
	})

	assert.Equal(t, O.Of("1"), f(Of(1))())
	assert.Equal(t, O.None[string](), f(None[int]())())
}

func TestEnv(t *testing.T) {
	env := Optionize1(os.LookupEnv)

	assert.True(t, O.IsSome(env("PATH")()))
	assert.False(t, O.IsSome(env("PATHxyz")()))
}
