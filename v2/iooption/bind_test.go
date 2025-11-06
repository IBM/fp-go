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
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

func getLastName(s utils.Initial) IOOption[string] {
	return Of("Doe")
}

func getGivenName(s utils.WithLastName) IOOption[string] {
	return Of("John")
}

func TestBind(t *testing.T) {

	res := F.Pipe3(
		Do(utils.Empty),
		Bind(utils.SetLastName, getLastName),
		Bind(utils.SetGivenName, getGivenName),
		Map(utils.GetFullName),
	)

	assert.Equal(t, res(), O.Of("John Doe"))
}

func TestApS(t *testing.T) {

	res := F.Pipe3(
		Do(utils.Empty),
		ApS(utils.SetLastName, Of("Doe")),
		ApS(utils.SetGivenName, Of("John")),
		Map(utils.GetFullName),
	)

	assert.Equal(t, res(), O.Of("John Doe"))
}
