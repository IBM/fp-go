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
	"testing"

	"github.com/IBM/fp-go/v2/eq"
	N "github.com/IBM/fp-go/v2/number"
	L "github.com/IBM/fp-go/v2/optics/lens"
	"github.com/stretchr/testify/assert"
)

// Test Alt function
func TestAlt(t *testing.T) {
	t.Run("Some value - returns original", func(t *testing.T) {
		withDefault := Alt(func() (int, bool) { return 100, true })
		AssertEq(Some(42))(withDefault(Some(42)))(t)
	})

	t.Run("None value - returns alternative Some", func(t *testing.T) {
		withDefault := Alt(func() (int, bool) { return 100, true })
		AssertEq(Some(100))(withDefault(None[int]()))(t)
	})

	t.Run("None value - alternative is also None", func(t *testing.T) {
		withDefault := Alt(func() (int, bool) { return None[int]() })
		AssertEq(None[int]())(withDefault(None[int]()))(t)
	})
}

// Test Reduce function
func TestReduce(t *testing.T) {
	t.Run("Some value - applies reducer", func(t *testing.T) {
		sum := Reduce(func(acc, val int) int { return acc + val }, 10)
		result := sum(Some(5))
		assert.Equal(t, 15, result)
	})

	t.Run("None value - returns initial", func(t *testing.T) {
		sum := Reduce(func(acc, val int) int { return acc + val }, 10)
		result := sum(None[int]())
		assert.Equal(t, 10, result)
	})

	t.Run("string concatenation", func(t *testing.T) {
		concat := Reduce(func(acc, val string) string { return acc + val }, "prefix:")
		result := concat(Some("test"))
		assert.Equal(t, "prefix:test", result)
	})
}

// Test FromZero function
func TestFromZero(t *testing.T) {
	t.Run("zero value - returns Some", func(t *testing.T) {
		AssertEq(Some(0))(FromZero[int]()(0))(t)
	})

	t.Run("non-zero value - returns None", func(t *testing.T) {
		AssertEq(None[int]())(FromZero[int]()(5))(t)
	})

	t.Run("empty string - returns Some", func(t *testing.T) {
		AssertEq(Some(""))(FromZero[string]()(""))(t)
	})

	t.Run("non-empty string - returns None", func(t *testing.T) {
		AssertEq(None[string]())(FromZero[string]()("hello"))(t)
	})
}

// Test FromNonZero function
func TestFromNonZero(t *testing.T) {
	t.Run("non-zero value - returns Some", func(t *testing.T) {
		AssertEq(Some(5))(FromNonZero[int]()(5))(t)
	})

	t.Run("zero value - returns None", func(t *testing.T) {
		AssertEq(None[int]())(FromNonZero[int]()(0))(t)
	})

	t.Run("non-empty string - returns Some", func(t *testing.T) {
		AssertEq(Some("hello"))(FromNonZero[string]()("hello"))(t)
	})

	t.Run("empty string - returns None", func(t *testing.T) {
		AssertEq(None[string]())(FromNonZero[string]()(""))(t)
	})
}

// Test FromEq function
func TestFromEq(t *testing.T) {
	t.Run("matching value - returns Some", func(t *testing.T) {
		equals42 := FromEq(eq.FromStrictEquals[int]())(42)
		AssertEq(Some(42))(equals42(42))(t)
	})

	t.Run("non-matching value - returns None", func(t *testing.T) {
		equals42 := FromEq(eq.FromStrictEquals[int]())(42)
		AssertEq(None[int]())(equals42(10))(t)
	})

	t.Run("string equality", func(t *testing.T) {
		equalsHello := FromEq(eq.FromStrictEquals[string]())("hello")
		assert.True(t, IsSome(equalsHello("hello")))
		assert.True(t, IsNone(equalsHello("world")))
	})
}

// Test Pipe and Flow functions
func TestPipe1(t *testing.T) {
	double := func(x int) (int, bool) { return x * 2, true }
	AssertEq(Some(10))(Pipe1(5, double))(t)
}

func TestFlow1(t *testing.T) {
	double := func(x int, ok bool) (int, bool) { return x * 2, ok }
	flow := Flow1(double)
	AssertEq(Some(10))(flow(Some(5)))(t)
}

func TestFlow2(t *testing.T) {
	double := func(x int, ok bool) (int, bool) { return x * 2, ok }
	add10 := func(x int, ok bool) (int, bool) {
		if ok {
			return x + 10, true
		}
		return 0, false
	}
	flow := Flow2(double, add10)
	AssertEq(Some(20))(flow(Some(5)))(t)
}

func TestPipe3(t *testing.T) {
	double := func(x int) (int, bool) { return x * 2, true }
	add10 := func(x int, ok bool) (int, bool) {
		if ok {
			return x + 10, true
		}
		return 0, false
	}
	mul3 := func(x int, ok bool) (int, bool) {
		if ok {
			return x * 3, true
		}
		return 0, false
	}
	AssertEq(Some(60))(Pipe3(5, double, add10, mul3))(t) // (5 * 2 + 10) * 3 = 60
}

func TestPipe4(t *testing.T) {
	double := func(x int) (int, bool) { return x * 2, true }
	add10 := func(x int, ok bool) (int, bool) {
		if ok {
			return x + 10, true
		}
		return 0, false
	}
	mul3 := func(x int, ok bool) (int, bool) {
		if ok {
			return x * 3, true
		}
		return 0, false
	}
	sub5 := func(x int, ok bool) (int, bool) {
		if ok {
			return x - 5, true
		}
		return 0, false
	}
	AssertEq(Some(55))(Pipe4(5, double, add10, mul3, sub5))(t) // ((5 * 2 + 10) * 3) - 5 = 55
}

func TestFlow4(t *testing.T) {
	f1 := func(x int, ok bool) (int, bool) { return x + 1, ok }
	f2 := func(x int, ok bool) (int, bool) { return x * 2, ok }
	f3 := func(x int, ok bool) (int, bool) { return x - 5, ok }
	f4 := func(x int, ok bool) (int, bool) { return x * 10, ok }
	flow := Flow4(f1, f2, f3, f4)
	AssertEq(Some(70))(flow(Some(5)))(t) // ((5 + 1) * 2 - 5) * 10 = 70
}

func TestFlow5(t *testing.T) {
	f1 := func(x int, ok bool) (int, bool) { return x + 1, ok }
	f2 := func(x int, ok bool) (int, bool) { return x * 2, ok }
	f3 := func(x int, ok bool) (int, bool) { return x - 5, ok }
	f4 := func(x int, ok bool) (int, bool) { return x * 10, ok }
	f5 := func(x int, ok bool) (int, bool) { return x + 100, ok }
	flow := Flow5(f1, f2, f3, f4, f5)
	AssertEq(Some(170))(flow(Some(5)))(t) // (((5 + 1) * 2 - 5) * 10) + 100 = 170
}

// Test Functor and Pointed
func TestMakeFunctor(t *testing.T) {
	t.Run("Map with functor", func(t *testing.T) {
		f := MakeFunctor[int, int]()
		double := f.Map(N.Mul(2))
		AssertEq(Some(42))(double(Some(21)))(t)
	})

	t.Run("Map with None", func(t *testing.T) {
		f := MakeFunctor[int, int]()
		double := f.Map(N.Mul(2))
		AssertEq(None[int]())(double(None[int]()))(t)
	})
}

func TestMakePointed(t *testing.T) {
	t.Run("Of with value", func(t *testing.T) {
		p := MakePointed[int]()
		AssertEq(Some(42))(p.Of(42))(t)
	})

	t.Run("Of with string", func(t *testing.T) {
		p := MakePointed[string]()
		AssertEq(Some("hello"))(p.Of("hello"))(t)
	})
}

// Test lens-based operations
type TestStruct struct {
	Value int
	Name  string
}

func TestApSL(t *testing.T) {
	valueLens := L.MakeLens(
		func(s TestStruct) int { return s.Value },
		func(s TestStruct, v int) TestStruct { s.Value = v; return s },
	)

	t.Run("Some struct, Some value", func(t *testing.T) {
		applyValue := ApSL(valueLens)
		v, ok := applyValue(Some(42))(Some(TestStruct{Value: 0, Name: "test"}))
		assert.True(t, ok)
		assert.Equal(t, 42, v.Value)
		assert.Equal(t, "test", v.Name)
	})

	t.Run("Some struct, None value", func(t *testing.T) {
		applyValue := ApSL(valueLens)
		AssertEq(None[TestStruct]())(applyValue(None[int]())(Some(TestStruct{Value: 10, Name: "test"})))(t)
	})

	t.Run("None struct, Some value", func(t *testing.T) {
		applyValue := ApSL(valueLens)
		AssertEq(None[TestStruct]())(applyValue(Some(42))(None[TestStruct]()))(t)
	})
}

func TestBindL(t *testing.T) {
	valueLens := L.MakeLens(
		func(s TestStruct) int { return s.Value },
		func(s TestStruct, v int) TestStruct { s.Value = v; return s },
	)

	t.Run("increment value with validation", func(t *testing.T) {
		increment := func(v int) (int, bool) {
			if v < 100 {
				return v + 1, true
			}
			return 0, false
		}
		bindIncrement := BindL(valueLens, increment)
		v, ok := bindIncrement(Some(TestStruct{Value: 42, Name: "test"}))
		assert.True(t, ok)
		assert.Equal(t, 43, v.Value)
		assert.Equal(t, "test", v.Name)
	})

	t.Run("validation fails", func(t *testing.T) {
		increment := func(v int) (int, bool) {
			if v < 100 {
				return v + 1, true
			}
			return 0, false
		}
		bindIncrement := BindL(valueLens, increment)
		AssertEq(None[TestStruct]())(bindIncrement(Some(TestStruct{Value: 100, Name: "test"})))(t)
	})

	t.Run("None input", func(t *testing.T) {
		increment := func(v int) (int, bool) { return v + 1, true }
		bindIncrement := BindL(valueLens, increment)
		AssertEq(None[TestStruct]())(bindIncrement(None[TestStruct]()))(t)
	})
}

func TestLetL(t *testing.T) {
	valueLens := L.MakeLens(
		func(s TestStruct) int { return s.Value },
		func(s TestStruct, v int) TestStruct { s.Value = v; return s },
	)

	t.Run("double value", func(t *testing.T) {
		double := func(v int) int { return v * 2 }
		letDouble := LetL(valueLens, double)
		v, ok := letDouble(Some(TestStruct{Value: 21, Name: "test"}))
		assert.True(t, ok)
		assert.Equal(t, 42, v.Value)
		assert.Equal(t, "test", v.Name)
	})

	t.Run("None input", func(t *testing.T) {
		double := func(v int) int { return v * 2 }
		letDouble := LetL(valueLens, double)
		AssertEq(None[TestStruct]())(letDouble(None[TestStruct]()))(t)
	})
}

func TestLetToL(t *testing.T) {
	valueLens := L.MakeLens(
		func(s TestStruct) int { return s.Value },
		func(s TestStruct, v int) TestStruct { s.Value = v; return s },
	)

	t.Run("set constant value", func(t *testing.T) {
		setValue := LetToL(valueLens, 100)
		v, ok := setValue(Some(TestStruct{Value: 42, Name: "test"}))
		assert.True(t, ok)
		assert.Equal(t, 100, v.Value)
		assert.Equal(t, "test", v.Name)
	})

	t.Run("None input", func(t *testing.T) {
		setValue := LetToL(valueLens, 100)
		AssertEq(None[TestStruct]())(setValue(None[TestStruct]()))(t)
	})
}

// Test tuple traversals
func TestTraverseTuple5(t *testing.T) {
	double := func(x int) (int, bool) { return x * 2, true }
	v1, v2, v3, v4, v5, ok := TraverseTuple5(double, double, double, double, double)(1, 2, 3, 4, 5)
	assert.True(t, ok)
	assert.Equal(t, 2, v1)
	assert.Equal(t, 4, v2)
	assert.Equal(t, 6, v3)
	assert.Equal(t, 8, v4)
	assert.Equal(t, 10, v5)
}

func TestTraverseTuple6(t *testing.T) {
	double := func(x int) (int, bool) { return x * 2, true }
	v1, v2, v3, v4, v5, v6, ok := TraverseTuple6(double, double, double, double, double, double)(1, 2, 3, 4, 5, 6)
	assert.True(t, ok)
	assert.Equal(t, 2, v1)
	assert.Equal(t, 4, v2)
	assert.Equal(t, 6, v3)
	assert.Equal(t, 8, v4)
	assert.Equal(t, 10, v5)
	assert.Equal(t, 12, v6)
}

func TestTraverseTuple7(t *testing.T) {
	double := func(x int) (int, bool) { return x * 2, true }
	v1, v2, v3, v4, v5, v6, v7, ok := TraverseTuple7(double, double, double, double, double, double, double)(1, 2, 3, 4, 5, 6, 7)
	assert.True(t, ok)
	assert.Equal(t, 2, v1)
	assert.Equal(t, 4, v2)
	assert.Equal(t, 6, v3)
	assert.Equal(t, 8, v4)
	assert.Equal(t, 10, v5)
	assert.Equal(t, 12, v6)
	assert.Equal(t, 14, v7)
}

func TestTraverseTuple8(t *testing.T) {
	double := func(x int) (int, bool) { return x * 2, true }
	v1, v2, v3, v4, v5, v6, v7, v8, ok := TraverseTuple8(double, double, double, double, double, double, double, double)(1, 2, 3, 4, 5, 6, 7, 8)
	assert.True(t, ok)
	assert.Equal(t, 2, v1)
	assert.Equal(t, 4, v2)
	assert.Equal(t, 6, v3)
	assert.Equal(t, 8, v4)
	assert.Equal(t, 10, v5)
	assert.Equal(t, 12, v6)
	assert.Equal(t, 14, v7)
	assert.Equal(t, 16, v8)
}

func TestTraverseTuple9(t *testing.T) {
	double := func(x int) (int, bool) { return x * 2, true }
	v1, v2, v3, v4, v5, v6, v7, v8, v9, ok := TraverseTuple9(double, double, double, double, double, double, double, double, double)(1, 2, 3, 4, 5, 6, 7, 8, 9)
	assert.True(t, ok)
	assert.Equal(t, 2, v1)
	assert.Equal(t, 4, v2)
	assert.Equal(t, 6, v3)
	assert.Equal(t, 8, v4)
	assert.Equal(t, 10, v5)
	assert.Equal(t, 12, v6)
	assert.Equal(t, 14, v7)
	assert.Equal(t, 16, v8)
	assert.Equal(t, 18, v9)
}

func TestTraverseTuple10(t *testing.T) {
	double := func(x int) (int, bool) { return x * 2, true }
	v1, v2, v3, v4, v5, v6, v7, v8, v9, v10, ok := TraverseTuple10(double, double, double, double, double, double, double, double, double, double)(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
	assert.True(t, ok)
	assert.Equal(t, 2, v1)
	assert.Equal(t, 4, v2)
	assert.Equal(t, 6, v3)
	assert.Equal(t, 8, v4)
	assert.Equal(t, 10, v5)
	assert.Equal(t, 12, v6)
	assert.Equal(t, 14, v7)
	assert.Equal(t, 16, v8)
	assert.Equal(t, 18, v9)
	assert.Equal(t, 20, v10)
}

// Test tuple traversals with failure cases
func TestTraverseTuple5_Failure(t *testing.T) {
	validate := func(x int) (int, bool) {
		if x > 0 {
			return x, true
		}
		return 0, false
	}
	_, _, _, _, _, ok := TraverseTuple5(validate, validate, validate, validate, validate)(1, -2, 3, 4, 5)
	assert.False(t, ok)
}
