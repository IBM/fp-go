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

package reader

import (
	"context"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
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

func TestLet(t *testing.T) {
	type State struct {
		FirstName string
		LastName  string
		FullName  string
	}

	res := F.Pipe2(
		Do[context.Context](State{FirstName: "John", LastName: "Doe"}),
		Let[context.Context](
			func(full string) func(State) State {
				return func(s State) State { s.FullName = full; return s }
			},
			func(s State) string {
				return s.FirstName + " " + s.LastName
			},
		),
		Map[context.Context](func(s State) string { return s.FullName }),
	)

	assert.Equal(t, "John Doe", res(context.Background()))
}

func TestLetTo(t *testing.T) {
	type State struct {
		Name    string
		Version string
	}

	res := F.Pipe2(
		Do[context.Context](State{Name: "MyApp"}),
		LetTo[context.Context](
			func(v string) func(State) State {
				return func(s State) State { s.Version = v; return s }
			},
			"1.0.0",
		),
		Map[context.Context](func(s State) State { return s }),
	)

	result := res(context.Background())
	assert.Equal(t, "MyApp", result.Name)
	assert.Equal(t, "1.0.0", result.Version)
}

func TestBindTo(t *testing.T) {
	type State struct{ Name string }

	getName := Asks(func(c context.Context) string { return "TestName" })
	initState := BindTo[context.Context](func(name string) State {
		return State{Name: name}
	})
	result := initState(getName)

	state := result(context.Background())
	assert.Equal(t, "TestName", state.Name)
}
