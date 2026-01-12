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

package state

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/lens"
	"github.com/IBM/fp-go/v2/pair"
	"github.com/stretchr/testify/assert"
)

// Test types for bind operations
type BindTestState struct {
	Counter    int
	Multiplier int
}

type Accumulator struct {
	Value   int
	Doubled int
	Status  string
}

// TestDo verifies that Do initializes a computation with an empty value
func TestDo(t *testing.T) {
	initial := BindTestState{Counter: 5, Multiplier: 2}
	emptyAcc := Accumulator{Value: 0, Doubled: 0, Status: ""}

	computation := Do[BindTestState](emptyAcc)
	result := computation(initial)

	assert.Equal(t, initial, pair.Head(result), "state should be unchanged")
	assert.Equal(t, emptyAcc, pair.Tail(result), "value should be empty accumulator")
}

// TestBind verifies that Bind sequences a computation and binds the result
func TestBind(t *testing.T) {
	initial := BindTestState{Counter: 10, Multiplier: 3}

	// Start with an accumulator
	startAcc := Accumulator{Value: 5, Doubled: 0, Status: ""}

	// Bind a computation that reads from state and updates the accumulator
	computation := Bind(
		func(doubled int) func(Accumulator) Accumulator {
			return func(acc Accumulator) Accumulator {
				acc.Doubled = doubled
				return acc
			}
		},
		func(acc Accumulator) State[BindTestState, int] {
			return Gets(func(s BindTestState) int {
				return acc.Value * s.Multiplier
			})
		},
	)

	result := computation(Of[BindTestState](startAcc))(initial)

	assert.Equal(t, initial, pair.Head(result), "state should be unchanged")
	assert.Equal(t, 5, pair.Tail(result).Value, "value should be preserved")
	assert.Equal(t, 15, pair.Tail(result).Doubled, "doubled should be 5 * 3 = 15")
}

// TestBindChaining verifies chaining multiple Bind operations
func TestBindChaining(t *testing.T) {
	initial := BindTestState{Counter: 10, Multiplier: 2}

	computation := F.Pipe3(
		Do[BindTestState](Accumulator{}),
		Bind(
			func(v int) func(Accumulator) Accumulator {
				return func(acc Accumulator) Accumulator {
					acc.Value = v
					return acc
				}
			},
			func(acc Accumulator) State[BindTestState, int] {
				return Gets(func(s BindTestState) int { return s.Counter })
			},
		),
		Bind(
			func(d int) func(Accumulator) Accumulator {
				return func(acc Accumulator) Accumulator {
					acc.Doubled = d
					return acc
				}
			},
			func(acc Accumulator) State[BindTestState, int] {
				return Gets(func(s BindTestState) int {
					return acc.Value * s.Multiplier
				})
			},
		),
		Bind(
			func(status string) func(Accumulator) Accumulator {
				return func(acc Accumulator) Accumulator {
					acc.Status = status
					return acc
				}
			},
			func(acc Accumulator) State[BindTestState, string] {
				return Of[BindTestState]("completed")
			},
		),
	)

	result := computation(initial)

	assert.Equal(t, 10, pair.Tail(result).Value, "value should be 10")
	assert.Equal(t, 20, pair.Tail(result).Doubled, "doubled should be 20")
	assert.Equal(t, "completed", pair.Tail(result).Status, "status should be completed")
}

// TestLet verifies that Let computes a pure value and binds it
func TestLet(t *testing.T) {
	initial := BindTestState{Counter: 5, Multiplier: 2}
	startAcc := Accumulator{Value: 10, Doubled: 20, Status: ""}

	// Compute sum from existing fields
	computation := Let[BindTestState](
		func(sum int) func(Accumulator) Accumulator {
			return func(acc Accumulator) Accumulator {
				acc.Status = "sum"
				acc.Value = sum
				return acc
			}
		},
		func(acc Accumulator) int {
			return acc.Value + acc.Doubled
		},
	)

	result := computation(Of[BindTestState](startAcc))(initial)

	assert.Equal(t, initial, pair.Head(result), "state should be unchanged")
	assert.Equal(t, 30, pair.Tail(result).Value, "value should be sum: 10 + 20 = 30")
	assert.Equal(t, "sum", pair.Tail(result).Status, "status should be set")
}

// TestLetTo verifies that LetTo binds a constant value
func TestLetTo(t *testing.T) {
	initial := BindTestState{Counter: 5, Multiplier: 2}
	startAcc := Accumulator{Value: 10, Doubled: 0, Status: ""}

	computation := LetTo[BindTestState](
		func(status string) func(Accumulator) Accumulator {
			return func(acc Accumulator) Accumulator {
				acc.Status = status
				return acc
			}
		},
		"initialized",
	)

	result := computation(Of[BindTestState](startAcc))(initial)

	assert.Equal(t, initial, pair.Head(result), "state should be unchanged")
	assert.Equal(t, "initialized", pair.Tail(result).Status, "status should be initialized")
	assert.Equal(t, 10, pair.Tail(result).Value, "value should be preserved")
}

// TestBindTo verifies that BindTo creates an initial accumulator
func TestBindTo(t *testing.T) {
	initial := BindTestState{Counter: 42, Multiplier: 2}

	computation := F.Pipe1(
		Of[BindTestState](100),
		BindTo[BindTestState](func(x int) Accumulator {
			return Accumulator{Value: x, Doubled: 0, Status: "created"}
		}),
	)

	result := computation(initial)

	assert.Equal(t, initial, pair.Head(result), "state should be unchanged")
	assert.Equal(t, 100, pair.Tail(result).Value, "value should be 100")
	assert.Equal(t, 0, pair.Tail(result).Doubled, "doubled should be 0")
	assert.Equal(t, "created", pair.Tail(result).Status, "status should be created")
}

// TestBindToWithPipeline verifies BindTo in a complete pipeline
func TestBindToWithPipeline(t *testing.T) {
	initial := BindTestState{Counter: 10, Multiplier: 3}

	computation := F.Pipe2(
		Of[BindTestState](5),
		BindTo[BindTestState](func(x int) Accumulator {
			return Accumulator{Value: x}
		}),
		Bind(
			func(d int) func(Accumulator) Accumulator {
				return func(acc Accumulator) Accumulator {
					acc.Doubled = d
					return acc
				}
			},
			func(acc Accumulator) State[BindTestState, int] {
				return Gets(func(s BindTestState) int {
					return acc.Value * s.Multiplier
				})
			},
		),
	)

	result := computation(initial)

	assert.Equal(t, 5, pair.Tail(result).Value, "value should be 5")
	assert.Equal(t, 15, pair.Tail(result).Doubled, "doubled should be 15")
}

// TestApS verifies applicative-style binding
func TestApS(t *testing.T) {
	initial := BindTestState{Counter: 7, Multiplier: 2}
	startAcc := Accumulator{Value: 10, Doubled: 0, Status: ""}

	// Independent computation that doesn't depend on accumulator
	getCounter := Gets(func(s BindTestState) int { return s.Counter })

	computation := ApS(
		func(c int) func(Accumulator) Accumulator {
			return func(acc Accumulator) Accumulator {
				acc.Doubled = c
				return acc
			}
		},
		getCounter,
	)

	result := computation(Of[BindTestState](startAcc))(initial)

	assert.Equal(t, initial, pair.Head(result), "state should be unchanged")
	assert.Equal(t, 10, pair.Tail(result).Value, "value should be preserved")
	assert.Equal(t, 7, pair.Tail(result).Doubled, "doubled should be counter value")
}

// TestApSL verifies lens-based applicative binding
func TestApSL(t *testing.T) {
	initial := BindTestState{Counter: 42, Multiplier: 2}

	// Create a lens for the Doubled field
	doubledLens := lens.MakeLens(
		func(acc Accumulator) int { return acc.Doubled },
		func(acc Accumulator, d int) Accumulator {
			acc.Doubled = d
			return acc
		},
	)

	getCounter := Gets(func(s BindTestState) int { return s.Counter })

	computation := F.Pipe1(
		Of[BindTestState](Accumulator{Value: 10}),
		ApSL(doubledLens, getCounter),
	)

	result := computation(initial)

	assert.Equal(t, 10, pair.Tail(result).Value, "value should be preserved")
	assert.Equal(t, 42, pair.Tail(result).Doubled, "doubled should be set to counter")
}

// TestBindL verifies lens-based monadic binding
func TestBindL(t *testing.T) {
	initial := BindTestState{Counter: 5, Multiplier: 3}

	// Create a lens for the Value field
	valueLens := lens.MakeLens(
		func(acc Accumulator) int { return acc.Value },
		func(acc Accumulator, v int) Accumulator {
			acc.Value = v
			return acc
		},
	)

	// Multiply the value by the state's multiplier
	computation := F.Pipe1(
		Of[BindTestState](Accumulator{Value: 10}),
		BindL(
			valueLens,
			func(v int) State[BindTestState, int] {
				return Gets(func(s BindTestState) int {
					return v * s.Multiplier
				})
			},
		),
	)

	result := computation(initial)

	assert.Equal(t, 30, pair.Tail(result).Value, "value should be 10 * 3 = 30")
}

// TestLetL verifies lens-based pure transformation
func TestLetL(t *testing.T) {
	initial := BindTestState{Counter: 5, Multiplier: 2}

	// Create a lens for the Value field
	valueLens := lens.MakeLens(
		func(acc Accumulator) int { return acc.Value },
		func(acc Accumulator, v int) Accumulator {
			acc.Value = v
			return acc
		},
	)

	// Double the value using a pure function
	computation := F.Pipe1(
		Of[BindTestState](Accumulator{Value: 21}),
		LetL[BindTestState](valueLens, func(v int) int { return v * 2 }),
	)

	result := computation(initial)

	assert.Equal(t, initial, pair.Head(result), "state should be unchanged")
	assert.Equal(t, 42, pair.Tail(result).Value, "value should be doubled to 42")
}

// TestLetToL verifies lens-based constant binding
func TestLetToL(t *testing.T) {
	initial := BindTestState{Counter: 5, Multiplier: 2}

	// Create a lens for the Status field
	statusLens := lens.MakeLens(
		func(acc Accumulator) string { return acc.Status },
		func(acc Accumulator, s string) Accumulator {
			acc.Status = s
			return acc
		},
	)

	computation := F.Pipe1(
		Of[BindTestState](Accumulator{Value: 10}),
		LetToL[BindTestState](statusLens, "completed"),
	)

	result := computation(initial)

	assert.Equal(t, 10, pair.Tail(result).Value, "value should be preserved")
	assert.Equal(t, "completed", pair.Tail(result).Status, "status should be set")
}

// TestComplexDoNotation verifies a complex do-notation pipeline
func TestComplexDoNotation(t *testing.T) {
	initial := BindTestState{Counter: 10, Multiplier: 2}

	// Create lenses
	valueLens := lens.MakeLens(
		func(acc Accumulator) int { return acc.Value },
		func(acc Accumulator, v int) Accumulator {
			acc.Value = v
			return acc
		},
	)

	doubledLens := lens.MakeLens(
		func(acc Accumulator) int { return acc.Doubled },
		func(acc Accumulator, d int) Accumulator {
			acc.Doubled = d
			return acc
		},
	)

	statusLens := lens.MakeLens(
		func(acc Accumulator) string { return acc.Status },
		func(acc Accumulator, s string) Accumulator {
			acc.Status = s
			return acc
		},
	)

	computation := F.Pipe5(
		Do[BindTestState](Accumulator{}),
		// Get counter from state and bind to Value
		Bind(
			valueLens.Set,
			func(acc Accumulator) State[BindTestState, int] {
				return Gets(func(s BindTestState) int { return s.Counter })
			},
		),
		// Compute doubled value using state
		BindL(
			doubledLens,
			func(d int) State[BindTestState, int] {
				return Gets(func(s BindTestState) int {
					return d * s.Multiplier
				})
			},
		),
		// Add a pure computation
		LetL[BindTestState](valueLens, func(v int) int { return v + 5 }),
		// Set a constant
		LetToL[BindTestState](statusLens, "processed"),
		// Extract final result
		Map[BindTestState](func(acc Accumulator) int {
			return acc.Value + acc.Doubled
		}),
	)

	result := computation(initial)

	// Value: 10 (counter) + 5 = 15
	// Doubled: 0 * 2 = 0 (doubled starts at 0)
	// Sum: 15 + 0 = 15
	assert.Equal(t, 15, pair.Tail(result), "final result should be 15")
}

// TestDoNotationWithModify verifies do-notation with state modifications
func TestDoNotationWithModify(t *testing.T) {
	initial := BindTestState{Counter: 0, Multiplier: 2}

	computation := F.Pipe3(
		Do[BindTestState](Accumulator{}),
		// Increment counter in state
		Bind(
			func(v int) func(Accumulator) Accumulator {
				return func(acc Accumulator) Accumulator {
					acc.Value = v
					return acc
				}
			},
			func(acc Accumulator) State[BindTestState, int] {
				return MonadChain(
					Modify(func(s BindTestState) BindTestState {
						s.Counter++
						return s
					}),
					func(_ Void) State[BindTestState, int] {
						return Gets(func(s BindTestState) int { return s.Counter })
					},
				)
			},
		),
		// Increment again
		Bind(
			func(d int) func(Accumulator) Accumulator {
				return func(acc Accumulator) Accumulator {
					acc.Doubled = d
					return acc
				}
			},
			func(acc Accumulator) State[BindTestState, int] {
				return MonadChain(
					Modify(func(s BindTestState) BindTestState {
						s.Counter++
						return s
					}),
					func(_ Void) State[BindTestState, int] {
						return Gets(func(s BindTestState) int { return s.Counter })
					},
				)
			},
		),
		Map[BindTestState](func(acc Accumulator) Accumulator { return acc }),
	)

	result := computation(initial)

	assert.Equal(t, 2, pair.Head(result).Counter, "counter should be incremented twice")
	assert.Equal(t, 1, pair.Tail(result).Value, "value should be 1")
	assert.Equal(t, 2, pair.Tail(result).Doubled, "doubled should be 2")
}

// TestLetWithComplexComputation verifies Let with complex pure computations
func TestLetWithComplexComputation(t *testing.T) {
	initial := BindTestState{Counter: 5, Multiplier: 2}

	computation := F.Pipe2(
		Do[BindTestState](Accumulator{Value: 10, Doubled: 20}),
		Let[BindTestState](
			func(sum int) func(Accumulator) Accumulator {
				return func(acc Accumulator) Accumulator {
					acc.Value = sum
					return acc
				}
			},
			func(acc Accumulator) int {
				// Complex computation from accumulator
				return (acc.Value + acc.Doubled) * 2
			},
		),
		Let[BindTestState](
			func(status string) func(Accumulator) Accumulator {
				return func(acc Accumulator) Accumulator {
					acc.Status = status
					return acc
				}
			},
			func(acc Accumulator) string {
				if acc.Value > 50 {
					return "high"
				}
				return "low"
			},
		),
	)

	result := computation(initial)

	// (10 + 20) * 2 = 60
	assert.Equal(t, 60, pair.Tail(result).Value, "value should be 60")
	assert.Equal(t, "high", pair.Tail(result).Status, "status should be high")
}

// TestMixedBindAndLet verifies mixing Bind and Let operations
func TestMixedBindAndLet(t *testing.T) {
	initial := BindTestState{Counter: 5, Multiplier: 3}

	computation := F.Pipe4(
		Do[BindTestState](Accumulator{}),
		// Bind: stateful operation
		Bind(
			func(v int) func(Accumulator) Accumulator {
				return func(acc Accumulator) Accumulator {
					acc.Value = v
					return acc
				}
			},
			func(acc Accumulator) State[BindTestState, int] {
				return Gets(func(s BindTestState) int { return s.Counter })
			},
		),
		// Let: pure computation
		Let[BindTestState](
			func(d int) func(Accumulator) Accumulator {
				return func(acc Accumulator) Accumulator {
					acc.Doubled = d
					return acc
				}
			},
			func(acc Accumulator) int {
				return acc.Value * 2
			},
		),
		// LetTo: constant
		LetTo[BindTestState](
			func(s string) func(Accumulator) Accumulator {
				return func(acc Accumulator) Accumulator {
					acc.Status = s
					return acc
				}
			},
			"done",
		),
		Map[BindTestState](func(acc Accumulator) Accumulator { return acc }),
	)

	result := computation(initial)

	assert.Equal(t, 5, pair.Tail(result).Value, "value should be 5")
	assert.Equal(t, 10, pair.Tail(result).Doubled, "doubled should be 10")
	assert.Equal(t, "done", pair.Tail(result).Status, "status should be done")
}
