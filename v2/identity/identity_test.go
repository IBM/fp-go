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

package identity

import (
	"fmt"
	"strconv"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	N "github.com/IBM/fp-go/v2/number"
	S "github.com/IBM/fp-go/v2/string"
	T "github.com/IBM/fp-go/v2/tuple"
	"github.com/stretchr/testify/assert"
)

func TestOf(t *testing.T) {
	t.Run("wraps int", func(t *testing.T) {
		result := Of(42)
		assert.Equal(t, 42, result)
	})

	t.Run("wraps string", func(t *testing.T) {
		result := Of("hello")
		assert.Equal(t, "hello", result)
	})

	t.Run("wraps struct", func(t *testing.T) {
		type Person struct{ Name string }
		p := Person{Name: "Alice"}
		result := Of(p)
		assert.Equal(t, p, result)
	})
}

func TestMap(t *testing.T) {
	t.Run("transforms int", func(t *testing.T) {
		result := F.Pipe1(1, Map(utils.Double))
		assert.Equal(t, 2, result)
	})

	t.Run("transforms string", func(t *testing.T) {
		result := F.Pipe1("hello", Map(S.Size))
		assert.Equal(t, 5, result)
	})

	t.Run("chains multiple maps", func(t *testing.T) {
		result := F.Pipe2(
			5,
			Map(N.Mul(2)),
			Map(N.Add(3)),
		)
		assert.Equal(t, 13, result)
	})
}

func TestMonadMap(t *testing.T) {
	t.Run("transforms value", func(t *testing.T) {
		result := MonadMap(10, N.Mul(3))
		assert.Equal(t, 30, result)
	})

	t.Run("changes type", func(t *testing.T) {
		result := MonadMap(42, S.Format[int]("Number: %d"))
		assert.Equal(t, "Number: 42", result)
	})
}

func TestMapTo(t *testing.T) {
	t.Run("replaces with constant int", func(t *testing.T) {
		result := F.Pipe1("ignored", MapTo[string](100))
		assert.Equal(t, 100, result)
	})

	t.Run("replaces with constant string", func(t *testing.T) {
		result := F.Pipe1(42, MapTo[int]("constant"))
		assert.Equal(t, "constant", result)
	})
}

func TestMonadMapTo(t *testing.T) {
	t.Run("replaces value", func(t *testing.T) {
		result := MonadMapTo("anything", 999)
		assert.Equal(t, 999, result)
	})
}

func TestChain(t *testing.T) {
	t.Run("chains computation", func(t *testing.T) {
		result := F.Pipe1(1, Chain(utils.Double))
		assert.Equal(t, 2, result)
	})

	t.Run("chains multiple operations", func(t *testing.T) {
		result := F.Pipe2(
			10,
			Chain(N.Mul(2)),
			Chain(N.Add(5)),
		)
		assert.Equal(t, 25, result)
	})

	t.Run("changes type", func(t *testing.T) {
		result := F.Pipe1(5, Chain(S.Format[int]("Value: %d")))
		assert.Equal(t, "Value: 5", result)
	})
}

func TestMonadChain(t *testing.T) {
	t.Run("chains computation", func(t *testing.T) {
		result := MonadChain(7, N.Mul(7))
		assert.Equal(t, 49, result)
	})
}

func TestChainFirst(t *testing.T) {
	t.Run("executes but keeps original", func(t *testing.T) {
		sideEffect := ""
		result := F.Pipe1(
			42,
			ChainFirst(func(n int) string {
				sideEffect = fmt.Sprintf("Processed: %d", n)
				return sideEffect
			}),
		)
		assert.Equal(t, 42, result)
		assert.Equal(t, "Processed: 42", sideEffect)
	})

	t.Run("chains with other operations", func(t *testing.T) {
		result := F.Pipe2(
			10,
			ChainFirst(func(n int) string { return "ignored" }),
			Map(N.Mul(2)),
		)
		assert.Equal(t, 20, result)
	})
}

func TestMonadChainFirst(t *testing.T) {
	t.Run("keeps original value", func(t *testing.T) {
		result := MonadChainFirst(100, strconv.Itoa)
		assert.Equal(t, 100, result)
	})
}

func TestAp(t *testing.T) {
	t.Run("applies function", func(t *testing.T) {
		result := F.Pipe1(utils.Double, Ap[int](1))
		assert.Equal(t, 2, result)
	})

	t.Run("applies curried function", func(t *testing.T) {
		add := N.Add[int]
		result := F.Pipe1(add(10), Ap[int](5))
		assert.Equal(t, 15, result)
	})

	t.Run("changes type", func(t *testing.T) {
		toString := S.Format[int]("Number: %d")
		result := F.Pipe1(toString, Ap[string](42))
		assert.Equal(t, "Number: 42", result)
	})
}

func TestMonadAp(t *testing.T) {
	t.Run("applies function to value", func(t *testing.T) {
		result := MonadAp(N.Mul(3), 7)
		assert.Equal(t, 21, result)
	})
}

func TestFlap(t *testing.T) {
	t.Run("flips application", func(t *testing.T) {
		double := N.Mul(2)
		result := F.Pipe1(double, Flap[int](5))
		assert.Equal(t, 10, result)
	})

	t.Run("with multiple functions", func(t *testing.T) {
		funcs := []func(int) int{
			N.Mul(2),
			N.Add(10),
			func(n int) int { return n * n },
		}

		results := make([]int, len(funcs))
		for i, f := range funcs {
			results[i] = Flap[int](5)(f)
		}

		assert.Equal(t, []int{10, 15, 25}, results)
	})
}

func TestMonadFlap(t *testing.T) {
	t.Run("applies value to function", func(t *testing.T) {
		result := MonadFlap(S.Format[int]("Value: %d"), 42)
		assert.Equal(t, "Value: 42", result)
	})
}

func TestDo(t *testing.T) {
	t.Run("initializes context", func(t *testing.T) {
		type State struct{ Value int }
		result := Do(State{Value: 10})
		assert.Equal(t, State{Value: 10}, result)
	})
}

func TestBind(t *testing.T) {
	t.Run("binds computation result", func(t *testing.T) {
		type State struct {
			X int
			Y int
		}

		result := F.Pipe2(
			Do(State{}),
			Bind(
				func(x int) func(State) State {
					return func(s State) State {
						s.X = x
						return s
					}
				},
				func(State) int { return 10 },
			),
			Bind(
				func(y int) func(State) State {
					return func(s State) State {
						s.Y = y
						return s
					}
				},
				func(State) int { return 20 },
			),
		)

		assert.Equal(t, State{X: 10, Y: 20}, result)
	})
}

func TestLet(t *testing.T) {
	t.Run("attaches computed value", func(t *testing.T) {
		type State struct {
			X   int
			Sum int
		}

		result := F.Pipe2(
			Do(State{X: 5}),
			Let(
				func(sum int) func(State) State {
					return func(s State) State {
						s.Sum = sum
						return s
					}
				},
				func(s State) int { return s.X * 2 },
			),
			Map(func(s State) State {
				s.Sum += 10
				return s
			}),
		)

		assert.Equal(t, State{X: 5, Sum: 20}, result)
	})
}

func TestLetTo(t *testing.T) {
	t.Run("attaches constant value", func(t *testing.T) {
		type State struct {
			Name string
		}

		result := F.Pipe1(
			Do(State{}),
			LetTo(
				func(name string) func(State) State {
					return func(s State) State {
						s.Name = name
						return s
					}
				},
				"Alice",
			),
		)

		assert.Equal(t, State{Name: "Alice"}, result)
	})
}

func TestBindTo(t *testing.T) {
	t.Run("initializes state from value", func(t *testing.T) {
		type State struct{ Value int }

		result := F.Pipe1(
			42,
			BindTo(func(v int) State {
				return State{Value: v}
			}),
		)

		assert.Equal(t, State{Value: 42}, result)
	})
}

func TestApS(t *testing.T) {
	t.Run("applies value in context", func(t *testing.T) {
		type State struct {
			X int
			Y int
		}

		result := F.Pipe1(
			Do(State{X: 10}),
			ApS(
				func(y int) func(State) State {
					return func(s State) State {
						s.Y = y
						return s
					}
				},
				20,
			),
		)

		assert.Equal(t, State{X: 10, Y: 20}, result)
	})
}

func TestSequenceT(t *testing.T) {
	t.Run("SequenceT2", func(t *testing.T) {
		result := SequenceT2(1, 2)
		assert.Equal(t, T.MakeTuple2(1, 2), result)
	})

	t.Run("SequenceT3", func(t *testing.T) {
		result := SequenceT3("a", "b", "c")
		assert.Equal(t, T.MakeTuple3("a", "b", "c"), result)
	})

	t.Run("SequenceT4", func(t *testing.T) {
		result := SequenceT4(1, 2, 3, 4)
		assert.Equal(t, T.MakeTuple4(1, 2, 3, 4), result)
	})
}

func TestSequenceTuple(t *testing.T) {
	t.Run("SequenceTuple2", func(t *testing.T) {
		tuple := T.MakeTuple2(10, 20)
		result := SequenceTuple2(tuple)
		assert.Equal(t, tuple, result)
	})

	t.Run("SequenceTuple3", func(t *testing.T) {
		tuple := T.MakeTuple3(1, 2, 3)
		result := SequenceTuple3(tuple)
		assert.Equal(t, tuple, result)
	})
}

func TestTraverseTuple(t *testing.T) {
	t.Run("TraverseTuple2", func(t *testing.T) {
		tuple := T.MakeTuple2(1, 2)
		result := TraverseTuple2(
			N.Mul(2),
			N.Mul(3),
		)(tuple)
		assert.Equal(t, T.MakeTuple2(2, 6), result)
	})

	t.Run("TraverseTuple3", func(t *testing.T) {
		tuple := T.MakeTuple3(1, 2, 3)
		result := TraverseTuple3(
			N.Add(10),
			func(n int) int { return n + 20 },
			func(n int) int { return n + 30 },
		)(tuple)
		assert.Equal(t, T.MakeTuple3(11, 22, 33), result)
	})

	t.Run("TraverseTuple2 with type change", func(t *testing.T) {
		tuple := T.MakeTuple2(5, 10)
		result := TraverseTuple2(
			func(n int) string { return fmt.Sprintf("A%d", n) },
			func(n int) string { return fmt.Sprintf("B%d", n) },
		)(tuple)
		assert.Equal(t, T.MakeTuple2("A5", "B10"), result)
	})
}

func TestMonad(t *testing.T) {
	t.Run("monad interface", func(t *testing.T) {
		m := Monad[int, string]()

		// Test Of
		value := m.Of(42)
		assert.Equal(t, 42, value)

		// Test Map
		mapped := m.Map(S.Format[int]("Number: %d"))(value)
		assert.Equal(t, "Number: 42", mapped)

		// Test Chain
		chained := m.Chain(S.Format[int]("Value: %d"))(value)
		assert.Equal(t, "Value: 42", chained)

		// Test Ap
		applied := m.Ap(10)(func(n int) string {
			return fmt.Sprintf("Result: %d", n)
		})
		assert.Equal(t, "Result: 10", applied)
	})
}

// Test monad laws
func TestMonadLaws(t *testing.T) {
	t.Run("left identity", func(t *testing.T) {
		// Of(a).Chain(f) === f(a)
		a := 42
		f := N.Mul(2)

		left := F.Pipe1(Of(a), Chain(f))
		right := f(a)

		assert.Equal(t, right, left)
	})

	t.Run("right identity", func(t *testing.T) {
		// m.Chain(Of) === m
		m := 42

		result := F.Pipe1(m, Chain(Of[int]))

		assert.Equal(t, m, result)
	})

	t.Run("associativity", func(t *testing.T) {
		// m.Chain(f).Chain(g) === m.Chain(x => f(x).Chain(g))
		m := 5
		f := N.Mul(2)
		g := N.Add(10)

		left := F.Pipe2(m, Chain(f), Chain(g))
		right := F.Pipe1(m, Chain(func(x int) int {
			return F.Pipe1(f(x), Chain(g))
		}))

		assert.Equal(t, right, left)
	})
}

// Test functor laws
func TestFunctorLaws(t *testing.T) {
	t.Run("identity", func(t *testing.T) {
		// Map(id) === id
		value := 42

		result := F.Pipe1(value, Map(F.Identity[int]))

		assert.Equal(t, value, result)
	})

	t.Run("composition", func(t *testing.T) {
		// Map(f).Map(g) === Map(g ∘ f)
		value := 5
		f := N.Mul(2)
		g := N.Add(10)

		left := F.Pipe2(value, Map(f), Map(g))
		right := F.Pipe1(value, Map(F.Flow2(f, g)))

		assert.Equal(t, right, left)
	})
}

func TestSequenceT1(t *testing.T) {
	t.Run("sequences single value", func(t *testing.T) {
		result := SequenceT1(42)
		assert.Equal(t, T.MakeTuple1(42), result)
	})
}

func TestSequenceTuple1(t *testing.T) {
	t.Run("sequences tuple1", func(t *testing.T) {
		tuple := T.MakeTuple1("hello")
		result := SequenceTuple1(tuple)
		assert.Equal(t, tuple, result)
	})
}

func TestTraverseTuple1(t *testing.T) {
	t.Run("traverses tuple1", func(t *testing.T) {
		tuple := T.MakeTuple1(5)
		result := TraverseTuple1(func(n int) int { return n * 10 })(tuple)
		assert.Equal(t, T.MakeTuple1(50), result)
	})
}

func TestSequenceTuple4(t *testing.T) {
	t.Run("sequences tuple4", func(t *testing.T) {
		tuple := T.MakeTuple4(1, 2, 3, 4)
		result := SequenceTuple4(tuple)
		assert.Equal(t, tuple, result)
	})
}

func TestTraverseTuple4(t *testing.T) {
	t.Run("traverses tuple4", func(t *testing.T) {
		tuple := T.MakeTuple4(1, 2, 3, 4)
		result := TraverseTuple4(
			N.Add(10),
			func(n int) int { return n + 20 },
			func(n int) int { return n + 30 },
			func(n int) int { return n + 40 },
		)(tuple)
		assert.Equal(t, T.MakeTuple4(11, 22, 33, 44), result)
	})
}

func TestSequenceT5(t *testing.T) {
	t.Run("sequences 5 values", func(t *testing.T) {
		result := SequenceT5(1, 2, 3, 4, 5)
		assert.Equal(t, T.MakeTuple5(1, 2, 3, 4, 5), result)
	})
}

func TestSequenceTuple5(t *testing.T) {
	t.Run("sequences tuple5", func(t *testing.T) {
		tuple := T.MakeTuple5(1, 2, 3, 4, 5)
		result := SequenceTuple5(tuple)
		assert.Equal(t, tuple, result)
	})
}

func TestTraverseTuple5(t *testing.T) {
	t.Run("traverses tuple5", func(t *testing.T) {
		tuple := T.MakeTuple5(1, 2, 3, 4, 5)
		result := TraverseTuple5(
			func(n int) int { return n * 1 },
			N.Mul(2),
			N.Mul(3),
			func(n int) int { return n * 4 },
			func(n int) int { return n * 5 },
		)(tuple)
		assert.Equal(t, T.MakeTuple5(1, 4, 9, 16, 25), result)
	})
}

func TestSequenceT6(t *testing.T) {
	t.Run("sequences 6 values", func(t *testing.T) {
		result := SequenceT6(1, 2, 3, 4, 5, 6)
		assert.Equal(t, T.MakeTuple6(1, 2, 3, 4, 5, 6), result)
	})
}

func TestSequenceTuple6(t *testing.T) {
	t.Run("sequences tuple6", func(t *testing.T) {
		tuple := T.MakeTuple6(1, 2, 3, 4, 5, 6)
		result := SequenceTuple6(tuple)
		assert.Equal(t, tuple, result)
	})
}

func TestTraverseTuple6(t *testing.T) {
	t.Run("traverses tuple6", func(t *testing.T) {
		tuple := T.MakeTuple6(1, 2, 3, 4, 5, 6)
		result := TraverseTuple6(
			N.Add(1),
			func(n int) int { return n + 2 },
			N.Add(3),
			func(n int) int { return n + 4 },
			N.Add(5),
			func(n int) int { return n + 6 },
		)(tuple)
		assert.Equal(t, T.MakeTuple6(2, 4, 6, 8, 10, 12), result)
	})
}

func TestSequenceT7(t *testing.T) {
	t.Run("sequences 7 values", func(t *testing.T) {
		result := SequenceT7(1, 2, 3, 4, 5, 6, 7)
		assert.Equal(t, T.MakeTuple7(1, 2, 3, 4, 5, 6, 7), result)
	})
}

func TestSequenceTuple7(t *testing.T) {
	t.Run("sequences tuple7", func(t *testing.T) {
		tuple := T.MakeTuple7(1, 2, 3, 4, 5, 6, 7)
		result := SequenceTuple7(tuple)
		assert.Equal(t, tuple, result)
	})
}

func TestTraverseTuple7(t *testing.T) {
	t.Run("traverses tuple7", func(t *testing.T) {
		tuple := T.MakeTuple7(1, 2, 3, 4, 5, 6, 7)
		result := TraverseTuple7(
			func(n int) int { return n * 10 },
			func(n int) int { return n * 10 },
			func(n int) int { return n * 10 },
			func(n int) int { return n * 10 },
			func(n int) int { return n * 10 },
			func(n int) int { return n * 10 },
			func(n int) int { return n * 10 },
		)(tuple)
		assert.Equal(t, T.MakeTuple7(10, 20, 30, 40, 50, 60, 70), result)
	})
}

func TestSequenceT8(t *testing.T) {
	t.Run("sequences 8 values", func(t *testing.T) {
		result := SequenceT8(1, 2, 3, 4, 5, 6, 7, 8)
		assert.Equal(t, T.MakeTuple8(1, 2, 3, 4, 5, 6, 7, 8), result)
	})
}

func TestSequenceTuple8(t *testing.T) {
	t.Run("sequences tuple8", func(t *testing.T) {
		tuple := T.MakeTuple8(1, 2, 3, 4, 5, 6, 7, 8)
		result := SequenceTuple8(tuple)
		assert.Equal(t, tuple, result)
	})
}

func TestTraverseTuple8(t *testing.T) {
	t.Run("traverses tuple8", func(t *testing.T) {
		tuple := T.MakeTuple8(1, 2, 3, 4, 5, 6, 7, 8)
		result := TraverseTuple8(
			func(n int) int { return n },
			func(n int) int { return n },
			func(n int) int { return n },
			func(n int) int { return n },
			func(n int) int { return n },
			func(n int) int { return n },
			func(n int) int { return n },
			func(n int) int { return n },
		)(tuple)
		assert.Equal(t, tuple, result)
	})
}

func TestSequenceT9(t *testing.T) {
	t.Run("sequences 9 values", func(t *testing.T) {
		result := SequenceT9(1, 2, 3, 4, 5, 6, 7, 8, 9)
		assert.Equal(t, T.MakeTuple9(1, 2, 3, 4, 5, 6, 7, 8, 9), result)
	})
}

func TestSequenceTuple9(t *testing.T) {
	t.Run("sequences tuple9", func(t *testing.T) {
		tuple := T.MakeTuple9(1, 2, 3, 4, 5, 6, 7, 8, 9)
		result := SequenceTuple9(tuple)
		assert.Equal(t, tuple, result)
	})
}

func TestTraverseTuple9(t *testing.T) {
	t.Run("traverses tuple9", func(t *testing.T) {
		tuple := T.MakeTuple9(1, 2, 3, 4, 5, 6, 7, 8, 9)
		result := TraverseTuple9(
			N.Add(1),
			N.Add(1),
			N.Add(1),
			N.Add(1),
			N.Add(1),
			N.Add(1),
			N.Add(1),
			N.Add(1),
			N.Add(1),
		)(tuple)
		assert.Equal(t, T.MakeTuple9(2, 3, 4, 5, 6, 7, 8, 9, 10), result)
	})
}

func TestSequenceT10(t *testing.T) {
	t.Run("sequences 10 values", func(t *testing.T) {
		result := SequenceT10(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		assert.Equal(t, T.MakeTuple10(1, 2, 3, 4, 5, 6, 7, 8, 9, 10), result)
	})
}

func TestSequenceTuple10(t *testing.T) {
	t.Run("sequences tuple10", func(t *testing.T) {
		tuple := T.MakeTuple10(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		result := SequenceTuple10(tuple)
		assert.Equal(t, tuple, result)
	})
}

func TestTraverseTuple10(t *testing.T) {
	t.Run("traverses tuple10", func(t *testing.T) {
		tuple := T.MakeTuple10(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
		result := TraverseTuple10(
			N.Mul(2),
			N.Mul(2),
			N.Mul(2),
			N.Mul(2),
			N.Mul(2),
			N.Mul(2),
			N.Mul(2),
			N.Mul(2),
			N.Mul(2),
			N.Mul(2),
		)(tuple)
		assert.Equal(t, T.MakeTuple10(2, 4, 6, 8, 10, 12, 14, 16, 18, 20), result)
	})
}

func TestExtract(t *testing.T) {
	t.Run("extracts int value", func(t *testing.T) {
		result := Extract(42)
		assert.Equal(t, 42, result)
	})

	t.Run("extracts string value", func(t *testing.T) {
		result := Extract("hello")
		assert.Equal(t, "hello", result)
	})

	t.Run("extracts struct value", func(t *testing.T) {
		type Person struct{ Name string }
		p := Person{Name: "Alice"}
		result := Extract(p)
		assert.Equal(t, p, result)
	})

	t.Run("extracts pointer value", func(t *testing.T) {
		value := 100
		ptr := &value
		result := Extract(ptr)
		assert.Equal(t, ptr, result)
		assert.Equal(t, 100, *result)
	})
}

func TestExtend(t *testing.T) {
	t.Run("extends with transformation", func(t *testing.T) {
		result := F.Pipe1(21, Extend(utils.Double))
		assert.Equal(t, 42, result)
	})

	t.Run("extends with type change", func(t *testing.T) {
		result := F.Pipe1(42, Extend(S.Format[int]("Number: %d")))
		assert.Equal(t, "Number: 42", result)
	})

	t.Run("chains multiple extends", func(t *testing.T) {
		result := F.Pipe2(
			5,
			Extend(N.Mul(2)),
			Extend(N.Add(10)),
		)
		assert.Equal(t, 20, result)
	})

	t.Run("extends with complex computation", func(t *testing.T) {
		result := F.Pipe1(
			10,
			Extend(func(n int) string {
				doubled := n * 2
				return fmt.Sprintf("Result: %d", doubled)
			}),
		)
		assert.Equal(t, "Result: 20", result)
	})
}

// Test Comonad laws
func TestComonadLaws(t *testing.T) {
	t.Run("left identity", func(t *testing.T) {
		// Extract(Extend(f)(w)) === f(w)
		w := 42
		f := N.Mul(2)

		left := Extract(F.Pipe1(w, Extend(f)))
		right := f(w)

		assert.Equal(t, right, left)
	})

	t.Run("right identity", func(t *testing.T) {
		// Extend(Extract)(w) === w
		w := 42

		result := F.Pipe1(w, Extend(Extract[int]))

		assert.Equal(t, w, result)
	})

	t.Run("associativity", func(t *testing.T) {
		// Extend(f)(Extend(g)(w)) === Extend(x => f(Extend(g)(x)))(w)
		w := 5
		f := N.Mul(2)
		g := N.Add(10)

		left := F.Pipe2(w, Extend(g), Extend(f))
		right := F.Pipe1(w, Extend(func(x int) int {
			return f(F.Pipe1(x, Extend(g)))
		}))

		assert.Equal(t, right, left)
	})
}
