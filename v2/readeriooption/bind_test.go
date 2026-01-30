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

package readeriooption

import (
	"context"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

func getLastName(s utils.Initial) ReaderIOOption[context.Context, string] {
	return Of[context.Context]("Doe")
}

func getGivenName(s utils.WithLastName) ReaderIOOption[context.Context, string] {
	return Of[context.Context]("John")
}

func TestBind(t *testing.T) {

	res := F.Pipe3(
		Do[context.Context](utils.Empty),
		Bind(utils.SetLastName, getLastName),
		Bind(utils.SetGivenName, getGivenName),
		Map[context.Context](utils.GetFullName),
	)

	assert.Equal(t, O.Of("John Doe"), res(context.Background())())
}

func TestApS(t *testing.T) {

	res := F.Pipe3(
		Do[context.Context](utils.Empty),
		ApS(utils.SetLastName, Of[context.Context]("Doe")),
		ApS(utils.SetGivenName, Of[context.Context]("John")),
		Map[context.Context](utils.GetFullName),
	)

	assert.Equal(t, O.Of("John Doe"), res(context.Background())())
}

func TestLet(t *testing.T) {
	res := F.Pipe3(
		Do[context.Context](utils.Empty),
		Let[context.Context](utils.SetLastName, func(s utils.Initial) string {
			return "Doe"
		}),
		Let[context.Context](utils.SetGivenName, func(s utils.WithLastName) string {
			return "John"
		}),
		Map[context.Context](utils.GetFullName),
	)

	assert.Equal(t, O.Of("John Doe"), res(context.Background())())
}

func TestLetTo(t *testing.T) {
	res := F.Pipe3(
		Do[context.Context](utils.Empty),
		LetTo[context.Context](utils.SetLastName, "Doe"),
		LetTo[context.Context](utils.SetGivenName, "John"),
		Map[context.Context](utils.GetFullName),
	)

	assert.Equal(t, O.Of("John Doe"), res(context.Background())())
}

func TestBindTo(t *testing.T) {
	type State struct {
		Value int
	}

	res := F.Pipe1(
		Of[context.Context](42),
		BindTo[context.Context](func(v int) State {
			return State{Value: v}
		}),
	)

	assert.Equal(t, O.Of(State{Value: 42}), res(context.Background())())
}
