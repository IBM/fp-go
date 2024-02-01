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

package writer

import (
	"testing"

	A "github.com/IBM/fp-go/array"
	F "github.com/IBM/fp-go/function"
	"github.com/IBM/fp-go/internal/utils"
	M "github.com/IBM/fp-go/monoid"
	"github.com/stretchr/testify/assert"
)

var (
	monoid = A.Monoid[string]()
	sg     = M.ToSemigroup(monoid)
)

func getLastName(s utils.Initial) Writer[[]string, string] {
	return Of[string](monoid)("Doe")
}

func getGivenName(s utils.WithLastName) Writer[[]string, string] {
	return Of[string](monoid)("John")
}

func TestBind(t *testing.T) {

	res := F.Pipe3(
		Do[utils.Initial](monoid)(utils.Empty),
		Bind(utils.SetLastName, getLastName),
		Bind(utils.SetGivenName, getGivenName),
		Map[[]string](utils.GetFullName),
	)

	assert.Equal(t, res(), Of[string](monoid)("John Doe")())
}

func TestApS(t *testing.T) {

	res := F.Pipe3(
		Do[utils.Initial](monoid)(utils.Empty),
		ApS(utils.SetLastName, Of[string](monoid)("Doe")),
		ApS(utils.SetGivenName, Of[string](monoid)("John")),
		Map[[]string](utils.GetFullName),
	)

	assert.Equal(t, res(), Of[string](monoid)("John Doe")())
}
