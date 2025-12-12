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
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

func getLastName(s utils.Initial) ReaderIOResult[string] {
	return Of("Doe")
}

func getGivenName(s utils.WithLastName) ReaderIOResult[string] {
	return Of("John")
}

func TestBind(t *testing.T) {

	res := F.Pipe3(
		Do(utils.Empty),
		Bind(utils.SetLastName, getLastName),
		Bind(utils.SetGivenName, getGivenName),
		Map(utils.GetFullName),
	)

	assert.Equal(t, res(t.Context())(), E.Of[error]("John Doe"))
}

func TestApS(t *testing.T) {

	res := F.Pipe3(
		Do(utils.Empty),
		ApS(utils.SetLastName, Of("Doe")),
		ApS(utils.SetGivenName, Of("John")),
		Map(utils.GetFullName),
	)

	assert.Equal(t, res(t.Context())(), E.Of[error]("John Doe"))
}

func TestApS_WithError(t *testing.T) {
	// Test that ApS propagates errors correctly
	testErr := assert.AnError

	res := F.Pipe3(
		Do(utils.Empty),
		ApS(utils.SetLastName, Left[string](testErr)),
		ApS(utils.SetGivenName, Of("John")),
		Map(utils.GetFullName),
	)

	result := res(t.Context())()
	assert.True(t, E.IsLeft(result))
	assert.Equal(t, testErr, E.ToError(result))
}

func TestApS_WithSecondError(t *testing.T) {
	// Test that ApS propagates errors from the second operation
	testErr := assert.AnError

	res := F.Pipe3(
		Do(utils.Empty),
		ApS(utils.SetLastName, Of("Doe")),
		ApS(utils.SetGivenName, Left[string](testErr)),
		Map(utils.GetFullName),
	)

	result := res(t.Context())()
	assert.True(t, E.IsLeft(result))
	assert.Equal(t, testErr, E.ToError(result))
}

func TestApS_MultipleFields(t *testing.T) {
	// Test ApS with more than two fields
	type Person struct {
		FirstName  string
		MiddleName string
		LastName   string
		Age        int
	}

	setFirstName := func(s string) func(Person) Person {
		return func(p Person) Person {
			p.FirstName = s
			return p
		}
	}

	setMiddleName := func(s string) func(Person) Person {
		return func(p Person) Person {
			p.MiddleName = s
			return p
		}
	}

	setLastName := func(s string) func(Person) Person {
		return func(p Person) Person {
			p.LastName = s
			return p
		}
	}

	setAge := func(a int) func(Person) Person {
		return func(p Person) Person {
			p.Age = a
			return p
		}
	}

	res := F.Pipe5(
		Do(Person{}),
		ApS(setFirstName, Of("John")),
		ApS(setMiddleName, Of("Q")),
		ApS(setLastName, Of("Doe")),
		ApS(setAge, Of(42)),
		Map(func(p Person) Person { return p }),
	)

	result := res(t.Context())()
	assert.True(t, E.IsRight(result))
	person := E.ToOption(result)
	assert.True(t, O.IsSome(person))
	p, _ := O.Unwrap(person)
	assert.Equal(t, "John", p.FirstName)
	assert.Equal(t, "Q", p.MiddleName)
	assert.Equal(t, "Doe", p.LastName)
	assert.Equal(t, 42, p.Age)
}

func TestApS_WithDifferentTypes(t *testing.T) {
	// Test ApS with different value types
	type State struct {
		Name  string
		Count int
		Flag  bool
	}

	setName := func(s string) func(State) State {
		return func(st State) State {
			st.Name = s
			return st
		}
	}

	setCount := func(c int) func(State) State {
		return func(st State) State {
			st.Count = c
			return st
		}
	}

	setFlag := func(f bool) func(State) State {
		return func(st State) State {
			st.Flag = f
			return st
		}
	}

	res := F.Pipe4(
		Do(State{}),
		ApS(setName, Of("test")),
		ApS(setCount, Of(100)),
		ApS(setFlag, Of(true)),
		Map(func(s State) State { return s }),
	)

	result := res(t.Context())()
	assert.True(t, E.IsRight(result))
	stateOpt := E.ToOption(result)
	assert.True(t, O.IsSome(stateOpt))
	state, _ := O.Unwrap(stateOpt)
	assert.Equal(t, "test", state.Name)
	assert.Equal(t, 100, state.Count)
	assert.True(t, state.Flag)
}

func TestApS_EmptyState(t *testing.T) {
	// Test ApS starting with an empty state
	type Empty struct{}

	res := Do(Empty{})

	result := res(t.Context())()
	assert.True(t, E.IsRight(result))
	emptyOpt := E.ToOption(result)
	assert.Equal(t, O.Of(Empty{}), emptyOpt)
}

func TestApS_ChainedWithBind(t *testing.T) {
	// Test mixing ApS with Bind operations
	type State struct {
		Independent string
		Dependent   string
	}

	setIndependent := func(s string) func(State) State {
		return func(st State) State {
			st.Independent = s
			return st
		}
	}

	setDependent := func(s string) func(State) State {
		return func(st State) State {
			st.Dependent = s
			return st
		}
	}

	getDependentValue := func(s State) ReaderIOResult[string] {
		// This depends on the Independent field
		return Of(s.Independent + "-dependent")
	}

	res := F.Pipe3(
		Do(State{}),
		ApS(setIndependent, Of("value")),
		Bind(setDependent, getDependentValue),
		Map(func(s State) State { return s }),
	)

	result := res(t.Context())()
	assert.True(t, E.IsRight(result))
	stateOpt := E.ToOption(result)
	assert.True(t, O.IsSome(stateOpt))
	state, _ := O.Unwrap(stateOpt)
	assert.Equal(t, "value", state.Independent)
	assert.Equal(t, "value-dependent", state.Dependent)
}

func TestApS_WithContextCancellation(t *testing.T) {
	// Test that ApS respects context cancellation
	type State struct {
		Value string
	}

	setValue := func(s string) func(State) State {
		return func(st State) State {
			st.Value = s
			return st
		}
	}

	// Create a computation that would succeed
	computation := ApS(setValue, Of("test"))(Do(State{}))

	// Create a cancelled context
	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	result := computation(ctx)()
	assert.True(t, E.IsLeft(result))
}
