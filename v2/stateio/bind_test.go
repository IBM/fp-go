// Copyright (c) 2024 - 2025 IBM Corp.
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

package stateio

import (
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/stretchr/testify/assert"
)

type BindTestState struct {
	Name  string
	Age   int
	Email string
}

func TestDo(t *testing.T) {
	initial := BindTestState{}
	computation := Do[BindTestState](BindTestState{Name: "initial"})

	result := computation(initial)()

	assert.Equal(t, BindTestState{Name: "initial"}, pair.Tail(result))
	assert.Equal(t, initial, pair.Head(result))
}

func TestBind(t *testing.T) {
	initial := BindTestState{}

	getName := func(s BindTestState) StateIO[BindTestState, string] {
		return Of[BindTestState]("John")
	}

	computation := F.Pipe2(
		Do[BindTestState](BindTestState{}),
		Bind(
			func(name string) func(BindTestState) BindTestState {
				return func(s BindTestState) BindTestState {
					s.Name = name
					return s
				}
			},
			getName,
		),
		Map[BindTestState](func(s BindTestState) string { return s.Name }),
	)

	result := computation(initial)()

	assert.Equal(t, "John", pair.Tail(result))
}

func TestBindMultiple(t *testing.T) {
	initial := BindTestState{}

	getName := func(s BindTestState) StateIO[BindTestState, string] {
		return Of[BindTestState]("Jane")
	}

	getAge := func(s BindTestState) StateIO[BindTestState, int] {
		return Of[BindTestState](30)
	}

	computation := F.Pipe3(
		Do[BindTestState](BindTestState{}),
		Bind(
			func(name string) func(BindTestState) BindTestState {
				return func(s BindTestState) BindTestState {
					s.Name = name
					return s
				}
			},
			getName,
		),
		Bind(
			func(age int) func(BindTestState) BindTestState {
				return func(s BindTestState) BindTestState {
					s.Age = age
					return s
				}
			},
			getAge,
		),
		Map[BindTestState](func(s BindTestState) BindTestState { return s }),
	)

	result := computation(initial)()
	finalState := pair.Tail(result)

	assert.Equal(t, "Jane", finalState.Name)
	assert.Equal(t, 30, finalState.Age)
}

func TestLet(t *testing.T) {
	initial := BindTestState{Age: 25}

	computation := F.Pipe2(
		Do[BindTestState](BindTestState{Age: 25}),
		Let[BindTestState](
			func(isAdult bool) func(BindTestState) BindTestState {
				return func(s BindTestState) BindTestState {
					if isAdult {
						s.Email = "adult@example.com"
					} else {
						s.Email = "minor@example.com"
					}
					return s
				}
			},
			func(s BindTestState) bool { return s.Age >= 18 },
		),
		Map[BindTestState](func(s BindTestState) string { return s.Email }),
	)

	result := computation(initial)()

	assert.Equal(t, "adult@example.com", pair.Tail(result))
}

func TestLetTo(t *testing.T) {
	initial := BindTestState{}

	computation := F.Pipe2(
		Do[BindTestState](BindTestState{}),
		LetTo[BindTestState](
			func(email string) func(BindTestState) BindTestState {
				return func(s BindTestState) BindTestState {
					s.Email = email
					return s
				}
			},
			"constant@example.com",
		),
		Map[BindTestState](func(s BindTestState) string { return s.Email }),
	)

	result := computation(initial)()

	assert.Equal(t, "constant@example.com", pair.Tail(result))
}

func TestBindTo(t *testing.T) {
	initial := BindTestState{}

	computation := F.Pipe2(
		Of[BindTestState]("Alice"),
		BindTo[BindTestState](func(name string) BindTestState {
			return BindTestState{Name: name}
		}),
		Map[BindTestState](func(s BindTestState) string { return s.Name }),
	)

	result := computation(initial)()

	assert.Equal(t, "Alice", pair.Tail(result))
}

func TestApS(t *testing.T) {
	initial := BindTestState{}

	computation := F.Pipe3(
		Do[BindTestState](BindTestState{}),
		ApS(
			func(name string) func(BindTestState) BindTestState {
				return func(s BindTestState) BindTestState {
					s.Name = name
					return s
				}
			},
			Of[BindTestState]("Bob"),
		),
		ApS(
			func(age int) func(BindTestState) BindTestState {
				return func(s BindTestState) BindTestState {
					s.Age = age
					return s
				}
			},
			Of[BindTestState](40),
		),
		Map[BindTestState](func(s BindTestState) BindTestState { return s }),
	)

	result := computation(initial)()
	finalState := pair.Tail(result)

	assert.Equal(t, "Bob", finalState.Name)
	assert.Equal(t, 40, finalState.Age)
}

func TestComplexDoNotation(t *testing.T) {
	initial := BindTestState{}

	fetchName := func(s BindTestState) StateIO[BindTestState, string] {
		return Of[BindTestState]("Charlie")
	}

	fetchAge := func(s BindTestState) StateIO[BindTestState, int] {
		return Of[BindTestState](35)
	}

	computation := F.Pipe4(
		Do[BindTestState](BindTestState{}),
		Bind(
			func(name string) func(BindTestState) BindTestState {
				return func(s BindTestState) BindTestState {
					s.Name = name
					return s
				}
			},
			fetchName,
		),
		Bind(
			func(age int) func(BindTestState) BindTestState {
				return func(s BindTestState) BindTestState {
					s.Age = age
					return s
				}
			},
			fetchAge,
		),
		Let[BindTestState](
			func(email string) func(BindTestState) BindTestState {
				return func(s BindTestState) BindTestState {
					s.Email = email
					return s
				}
			},
			func(s BindTestState) string {
				return fmt.Sprintf("%s@example.com", s.Name)
			},
		),
		Map[BindTestState](func(s BindTestState) BindTestState { return s }),
	)

	result := computation(initial)()
	finalState := pair.Tail(result)

	assert.Equal(t, "Charlie", finalState.Name)
	assert.Equal(t, 35, finalState.Age)
	assert.Equal(t, "Charlie@example.com", finalState.Email)
}

// Lens-based tests
type NestedState struct {
	User BindTestState
	ID   int
}

var userLens = lens.MakeLensCurried(
	func(s NestedState) BindTestState { return s.User },
	func(user BindTestState) func(NestedState) NestedState {
		return func(s NestedState) NestedState {
			s.User = user
			return s
		}
	},
)

var nameLens = lens.MakeLensCurried(
	func(s BindTestState) string { return s.Name },
	func(name string) func(BindTestState) BindTestState {
		return func(s BindTestState) BindTestState {
			s.Name = name
			return s
		}
	},
)

func TestApSL(t *testing.T) {
	initial := NestedState{User: BindTestState{}, ID: 1}

	computation := F.Pipe2(
		Do[NestedState](NestedState{User: BindTestState{}, ID: 1}),
		ApSL(
			userLens,
			Of[NestedState](BindTestState{Name: "David", Age: 28}),
		),
		Map[NestedState](func(s NestedState) string { return s.User.Name }),
	)

	result := computation(initial)()

	assert.Equal(t, "David", pair.Tail(result))
}

func TestBindL(t *testing.T) {
	initial := NestedState{User: BindTestState{Name: "Eve"}, ID: 2}

	updateUser := func(user BindTestState) StateIO[NestedState, BindTestState] {
		return Of[NestedState](BindTestState{
			Name:  user.Name + " Updated",
			Age:   user.Age + 1,
			Email: user.Email,
		})
	}

	computation := F.Pipe2(
		Do[NestedState](NestedState{User: BindTestState{Name: "Eve", Age: 20}, ID: 2}),
		BindL(userLens, updateUser),
		Map[NestedState](func(s NestedState) string { return s.User.Name }),
	)

	result := computation(initial)()

	assert.Equal(t, "Eve Updated", pair.Tail(result))
}

func TestLetL(t *testing.T) {
	initial := NestedState{User: BindTestState{Name: "Frank"}, ID: 3}

	uppercase := func(name string) string {
		return fmt.Sprintf("%s (UPPERCASE)", name)
	}

	composedLens := F.Pipe1(userLens, lens.Compose[NestedState](nameLens))

	computation := F.Pipe2(
		Do[NestedState](NestedState{User: BindTestState{Name: "Frank"}, ID: 3}),
		LetL[NestedState](
			composedLens,
			uppercase,
		),
		Map[NestedState](func(s NestedState) string { return s.User.Name }),
	)

	result := computation(initial)()

	assert.Equal(t, "Frank (UPPERCASE)", pair.Tail(result))
}

func TestLetToL(t *testing.T) {
	initial := NestedState{User: BindTestState{}, ID: 4}

	composedLens := F.Pipe1(userLens, lens.Compose[NestedState](nameLens))

	computation := F.Pipe2(
		Do[NestedState](NestedState{User: BindTestState{}, ID: 4}),
		LetToL[NestedState](
			composedLens,
			"Grace",
		),
		Map[NestedState](func(s NestedState) string { return s.User.Name }),
	)

	result := computation(initial)()

	assert.Equal(t, "Grace", pair.Tail(result))
}

func TestDoNotationWithStatefulOperations(t *testing.T) {
	type Counter struct {
		Value int
	}

	initial := Counter{Value: 0}

	increment := func(c Counter) StateIO[Counter, int] {
		return func(s Counter) IO[Pair[Counter, int]] {
			return func() Pair[Counter, int] {
				newState := Counter{Value: s.Value + 1}
				return pair.MakePair(newState, newState.Value)
			}
		}
	}

	computation := F.Pipe3(
		Do[Counter](Counter{Value: 0}),
		Bind(
			func(v int) func(Counter) Counter {
				return func(c Counter) Counter {
					return Counter{Value: v}
				}
			},
			increment,
		),
		Bind(
			func(v int) func(Counter) Counter {
				return func(c Counter) Counter {
					return Counter{Value: v}
				}
			},
			increment,
		),
		Map[Counter](func(c Counter) int { return c.Value }),
	)

	result := computation(initial)()

	// After two increments starting from 0, we should have 2
	assert.Equal(t, 2, pair.Tail(result))
	assert.Equal(t, 2, pair.Head(result).Value)
}
