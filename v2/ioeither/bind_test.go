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

package ioeither

import (
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	"github.com/stretchr/testify/assert"
)

func getLastName(s utils.Initial) IOEither[error, string] {
	return Of[error]("Doe")
}

func getGivenName(s utils.WithLastName) IOEither[error, string] {
	return Of[error]("John")
}

func TestBind(t *testing.T) {

	res := F.Pipe3(
		Do[error](utils.Empty),
		Bind(utils.SetLastName, getLastName),
		Bind(utils.SetGivenName, getGivenName),
		Map[error](utils.GetFullName),
	)

	assert.Equal(t, res(), E.Of[error]("John Doe"))
}

func TestApS(t *testing.T) {

	res := F.Pipe3(
		Do[error](utils.Empty),
		ApS(utils.SetLastName, Of[error]("Doe")),
		ApS(utils.SetGivenName, Of[error]("John")),
		Map[error](utils.GetFullName),
	)

	assert.Equal(t, res(), E.Of[error]("John Doe"))
}
