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

package readereither

import (
	"testing"

	ET "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	"github.com/stretchr/testify/assert"
)

type MyContext string

const defaultContext MyContext = "default"

func TestMap(t *testing.T) {

	g := F.Pipe1(
		Of[MyContext, error](1),
		Map[MyContext, error](utils.Double),
	)

	assert.Equal(t, ET.Of[error](2), g(defaultContext))

}

func TestAp(t *testing.T) {
	g := F.Pipe1(
		Of[MyContext, error](utils.Double),
		Ap[int](Of[MyContext, error](1)),
	)
	assert.Equal(t, ET.Of[error](2), g(defaultContext))

}

func TestFlatten(t *testing.T) {

	g := F.Pipe1(
		Of[MyContext, string](Of[MyContext, string]("a")),
		Flatten[MyContext, string, string],
	)

	assert.Equal(t, ET.Of[string]("a"), g(defaultContext))
}
