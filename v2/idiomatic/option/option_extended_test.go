// Copyright (c) 2025 IBM Corp.
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

package option

import (
	"fmt"
	"testing"

	N "github.com/IBM/fp-go/v2/number"
	"github.com/stretchr/testify/assert"
)

// Test FromNillable
func TestFromNillable(t *testing.T) {
	var nilPtr *int = nil
	AssertEq(None[*int]())(FromNillable(nilPtr))(t)

	val := 42
	ptr := &val
	result, resultok := FromNillable(ptr)
	assert.True(t, IsSome(result, resultok))
	assert.Equal(t, &val, result)
}

// Test MapTo
func TestMapTo(t *testing.T) {
	t.Run("positive case - replace value", func(t *testing.T) {
		replaceWith42 := MapTo[string](42)
		// Should replace value when Some
		AssertEq(Some(42))(replaceWith42(Some("hello")))(t)
		AssertEq(Some(42))(replaceWith42(Some("world")))(t)
	})

	t.Run("negative case - input is None", func(t *testing.T) {
		replaceWith42 := MapTo[string](42)
		// Should return None when input is None
		AssertEq(None[int]())(replaceWith42(None[string]()))(t)
	})
}

// Test GetOrElse
func TestGetOrElse(t *testing.T) {
	t.Run("positive case - extract value from Some", func(t *testing.T) {
		getOrZero := GetOrElse(func() int { return 0 })
		// Should extract value when Some
		assert.Equal(t, 42, getOrZero(Some(42)))
		assert.Equal(t, 100, getOrZero(Some(100)))
	})

	t.Run("negative case - use default for None", func(t *testing.T) {
		getOrZero := GetOrElse(func() int { return 0 })
		// Should return default when None
		assert.Equal(t, 0, getOrZero(None[int]()))
	})

	t.Run("positive case - custom default", func(t *testing.T) {
		getOrNegative := GetOrElse(func() int { return -1 })
		// Should use custom default
		assert.Equal(t, -1, getOrNegative(None[int]()))
		assert.Equal(t, 42, getOrNegative(Some(42)))
	})
}

// Test ChainTo
func TestChainTo(t *testing.T) {
	t.Run("positive case - replace with Some", func(t *testing.T) {
		replaceWith := ChainTo[int](Some("hello"))
		// Should replace any input with the fixed value
		AssertEq(Some("hello"))(replaceWith(Some(42)))(t)
		AssertEq(None[string]())(replaceWith(None[int]()))(t)
	})

	t.Run("negative case - replace with None", func(t *testing.T) {
		replaceWith := ChainTo[int](None[string]())
		// Should replace any input with None
		AssertEq(None[string]())(replaceWith(Some(42)))(t)
		AssertEq(None[string]())(replaceWith(None[int]()))(t)
	})
}

// Test ChainFirst
func TestChainFirst(t *testing.T) {
	t.Run("positive case - side effect succeeds", func(t *testing.T) {
		sideEffect := func(x int) (string, bool) {
			return Some(fmt.Sprintf("%d", x))
		}
		chainFirst := ChainFirst(sideEffect)

		// Should keep original value when side effect succeeds
		AssertEq(Some(5))(chainFirst(Some(5)))(t)
	})

	t.Run("negative case - side effect fails", func(t *testing.T) {
		sideEffect := func(x int) (string, bool) {
			if x < 0 {
				return None[string]()
			}
			return Some(fmt.Sprintf("%d", x))
		}
		chainFirst := ChainFirst(sideEffect)

		// Should return None when side effect fails
		AssertEq(None[int]())(chainFirst(Some(-5)))(t)
	})

	t.Run("negative case - input is None", func(t *testing.T) {
		sideEffect := func(x int) (string, bool) {
			return Some(fmt.Sprintf("%d", x))
		}
		chainFirst := ChainFirst(sideEffect)

		// Should return None when input is None
		AssertEq(None[int]())(chainFirst(None[int]()))(t)
	})
}

// Test Filter
func TestFilter(t *testing.T) {
	t.Run("positive case - predicate satisfied", func(t *testing.T) {
		isPositive := Filter(N.MoreThan(0))
		// Should keep value when predicate is satisfied
		AssertEq(Some(5))(isPositive(Some(5)))(t)
	})

	t.Run("negative case - predicate not satisfied", func(t *testing.T) {
		isPositive := Filter(N.MoreThan(0))
		// Should return None when predicate fails
		AssertEq(None[int]())(isPositive(Some(-1)))(t)
		AssertEq(None[int]())(isPositive(Some(0)))(t)
	})

	t.Run("negative case - input is None", func(t *testing.T) {
		isPositive := Filter(N.MoreThan(0))
		// Should return None when input is None
		AssertEq(None[int]())(isPositive(None[int]()))(t)
	})
}

// Test Flap
func TestFlap(t *testing.T) {
	t.Run("positive case - function is Some", func(t *testing.T) {
		applyFive := Flap[int](5)
		double := N.Mul(2)
		// Should apply value to function
		AssertEq(Some(10))(applyFive(Some(double)))(t)
	})

	t.Run("positive case - multiple operations", func(t *testing.T) {
		applyTen := Flap[int](10)
		triple := N.Mul(3)
		// Should work with different values
		AssertEq(Some(30))(applyTen(Some(triple)))(t)
	})

	t.Run("negative case - function is None", func(t *testing.T) {
		applyFive := Flap[int](5)
		// Should return None when function is None
		AssertEq(None[int]())(applyFive(None[func(int) int]()))(t)
	})
}

// Test String and Format
func TestStringFormat(t *testing.T) {
	str := ToString(Some(42))
	assert.Contains(t, str, "Some")
	assert.Contains(t, str, "42")

	str = ToString(None[int]())
	assert.Contains(t, str, "None")
}

// // Test Semigroup
// func TestSemigroup(t *testing.T) {
// 	intSemigroup := N.MonoidSum[int]()
// 	optSemigroup := Semigroup[int]()(intSemigroup)

// 	AssertEq(Some(5), optSemigroup.Concat(Some(2), Some(3)))
// 	AssertEq(Some(2), optSemigroup.Concat(Some(2), None[int]()))
// 	AssertEq(Some(3), optSemigroup.Concat(None[int](), Some(3)))
// 	AssertEq(None[int](), optSemigroup.Concat(None[int](), None[int]()))
// }

// // Test Monoid
// func TestMonoid(t *testing.T) {
// 	intSemigroup := N.MonoidSum[int]()
// 	optMonoid := Monoid[int]()(intSemigroup)

// 	AssertEq(Some(5), optMonoid.Concat(Some(2), Some(3)))
// 	AssertEq(None[int](), optMonoid.Empty())
// }

// // Test ApplySemigroup
// func TestApplySemigroup(t *testing.T) {
// 	intSemigroup := N.MonoidSum[int]()
// 	optSemigroup := ApplySemigroup(intSemigroup)

// 	AssertEq(Some(5), optSemigroup.Concat(Some(2), Some(3)))
// 	AssertEq(None[int](), optSemigroup.Concat(Some(2), None[int]()))
// }

// // Test ApplicativeMonoid
// func TestApplicativeMonoid(t *testing.T) {
// 	intMonoid := N.MonoidSum[int]()
// 	optMonoid := ApplicativeMonoid(intMonoid)

// 	AssertEq(Some(5), optMonoid.Concat(Some(2), Some(3)))
// 	AssertEq(Some(0), optMonoid.Empty())
// }

// // Test AlternativeMonoid
// func TestAlternativeMonoid(t *testing.T) {
// 	intMonoid := N.MonoidSum[int]()
// 	optMonoid := AlternativeMonoid(intMonoid)

// 	// AlternativeMonoid uses applicative semantics, so it combines values
// 	AssertEq(Some(5), optMonoid.Concat(Some(2), Some(3)))
// 	AssertEq(Some(3), optMonoid.Concat(None[int](), Some(3)))
// 	AssertEq(Some(0), optMonoid.Empty())
// }

// // Test AltMonoid
// func TestAltMonoid(t *testing.T) {
// 	optMonoid := AltMonoid[int]()

// 	AssertEq(Some(2), optMonoid.Concat(Some(2), Some(3)))
// 	AssertEq(Some(3), optMonoid.Concat(None[int](), Some(3)))
// 	AssertEq(None[int](), optMonoid.Empty())
// }

// Test Do, Let, LetTo, BindTo
func TestDoLetLetToBindTo(t *testing.T) {
	type State struct {
		x        int
		y        int
		computed int
		name     string
	}

	result, resultok := Pipe5(
		State{},
		Do,
		Let(func(c int) func(State) State {
			return func(s State) State { s.x = c; return s }
		}, func(s State) int { return 5 }),
		LetTo(func(n string) func(State) State {
			return func(s State) State { s.name = n; return s }
		}, "test"),
		Bind(func(y int) func(State) State {
			return func(s State) State { s.y = y; return s }
		}, func(s State) (int, bool) { return Some(10) }),
		Map(func(s State) State {
			s.computed = s.x + s.y
			return s
		}),
	)

	AssertEq(Some(State{x: 5, y: 10, computed: 15, name: "test"}))(result, resultok)(t)
}

// Test BindTo
func TestBindToFunction(t *testing.T) {
	type State struct {
		value int
	}

	result, resultok := Pipe2(
		42,
		Some,
		BindTo(func(x int) State { return State{value: x} }),
	)

	AssertEq(Some(State{value: 42}))(result, resultok)(t)
}

// // Test Functor
// func TestFunctor(t *testing.T) {
// 	f := Functor[int, string]()
// 	mapper := f.Map(strconv.Itoa)

// 	AssertEq(Some("42"), mapper(Some(42)))
// 	AssertEq(None[string](), mapper(None[int]()))
// }

// // Test Monad
// func TestMonad(t *testing.T) {
// 	m := Monad[int, string]()

// 	// Test Of
// 	AssertEq(Some(42), m.Of(42))

// 	// Test Map
// 	mapper := m.Map(strconv.Itoa)
// 	AssertEq(Some("42"), mapper(Some(42)))

// 	// Test Chain
// 	chainer := m.Chain(func(x int) (string, bool) {
// 		if x > 0 {
// 			return Some(fmt.Sprintf("%d", x))
// 		}
// 		return None[string]()
// 	})
// 	AssertEq(Some("42"), chainer(Some(42)))

// 	// Test Ap
// 	double := func(x int) string { return fmt.Sprintf("%d", x*2) }
// 	ap := m.Ap(Some(5))
// 	AssertEq(Some("10"), ap(Some(double)))
// }

// // Test Pointed
// func TestPointed(t *testing.T) {
// 	p := Pointed[int]()
// 	AssertEq(Some(42), p.Of(42))
// }

// Test ToAny
func TestToAny(t *testing.T) {
	result, resultok := ToAny(42)
	assert.True(t, IsSome(result, resultok))

	assert.Equal(t, 42, result)
}

// Test TraverseArray
func TestTraverseArray(t *testing.T) {
	validate := func(x int) (int, bool) {
		if x > 0 {
			return Some(x * 2)
		}
		return None[int]()
	}

	result, resultok := TraverseArray(validate)([]int{1, 2, 3})
	AssertEq(Some([]int{2, 4, 6}))(result, resultok)(t)

	result, resultok = TraverseArray(validate)([]int{1, -1, 3})
	AssertEq(None[[]int]())(result, resultok)(t)
}

// Test TraverseArrayWithIndex
func TestTraverseArrayWithIndex(t *testing.T) {
	f := func(i int, x int) (int, bool) {
		if x > i {
			return Some(x + i)
		}
		return None[int]()
	}

	result, resultok := TraverseArrayWithIndex(f)([]int{1, 2, 3})
	AssertEq(Some([]int{1, 3, 5}))(result, resultok)(t)
}

// Test TraverseRecord
func TestTraverseRecord(t *testing.T) {
	validate := func(x int) (string, bool) {
		if x > 0 {
			return Some(fmt.Sprintf("%d", x))
		}
		return None[string]()
	}

	input := map[string]int{"a": 1, "b": 2}
	result, resultok := TraverseRecord[string](validate)(input)

	AssertEq(Some(map[string]string{"a": "1", "b": "2"}))(result, resultok)(t)
}

// Test TraverseRecordWithIndex
func TestTraverseRecordWithIndex(t *testing.T) {
	f := func(k string, v int) (string, bool) {
		return Some(fmt.Sprintf("%s:%d", k, v))
	}

	input := map[string]int{"a": 1, "b": 2}
	result, resultok := TraverseRecordWithIndex(f)(input)

	assert.True(t, IsSome(result, resultok))
}

// Test TraverseTuple functions
func TestTraverseTuple2(t *testing.T) {
	f1 := func(x int) (int, bool) { return Some(x * 2) }
	f2 := func(s string) (string, bool) { return Some(s + "!") }

	traverse := TraverseTuple2(f1, f2)
	r1, r2, resultok := traverse(5, "hello")

	assert.True(t, resultok)
	assert.Equal(t, r1, 10)
	assert.Equal(t, r2, "hello!")
}

// Test FromStrictCompare
func TestFromStrictCompare(t *testing.T) {
	optOrd := FromStrictCompare[int]()

	assert.Equal(t, 0, optOrd(Some(5))(Some(5)))
	assert.Equal(t, -1, optOrd(Some(3))(Some(5)))
	assert.Equal(t, +1, optOrd(Some(5))(Some(3)))
	assert.Equal(t, -1, optOrd(None[int]())(Some(5)))
	assert.Equal(t, +1, optOrd(Some(5))(None[int]()))
}
