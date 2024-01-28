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

package reader

import (
	"context"
	"testing"

	F "github.com/IBM/fp-go/function"
	"github.com/IBM/fp-go/internal/utils"
	"github.com/stretchr/testify/assert"
)

func getLastName(s utils.Initial) Reader[context.Context, string] {
	return Of[context.Context]("Doe")
}

func getGivenName(s utils.WithLastName) Reader[context.Context, string] {
	return Of[context.Context]("John")
}

func TestBind(t *testing.T) {

	res := F.Pipe3(
		Do[context.Context](utils.Empty),
		Bind(utils.SetLastName, getLastName),
		Bind(utils.SetGivenName, getGivenName),
		Map[context.Context](utils.GetFullName),
	)

	assert.Equal(t, res(context.Background()), "John Doe")
}

func TestApS(t *testing.T) {

	res := F.Pipe3(
		Do[context.Context](utils.Empty),
		ApS(utils.SetLastName, Of[context.Context]("Doe")),
		ApS(utils.SetGivenName, Of[context.Context]("John")),
		Map[context.Context](utils.GetFullName),
	)

	assert.Equal(t, res(context.Background()), "John Doe")
}
