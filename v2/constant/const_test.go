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

package constant

import (
	"strconv"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	N "github.com/IBM/fp-go/v2/number"
	S "github.com/IBM/fp-go/v2/string"

	"github.com/stretchr/testify/assert"
)

// TestMake tests the Make constructor
func TestMake(t *testing.T) {
	t.Run("creates Const with string value", func(t *testing.T) {
		c := Make[string, int]("hello")
		assert.Equal(t, "hello", Unwrap(c))
	})

	t.Run("creates Const with int value", func(t *testing.T) {
		c := Make[int, string](42)
		assert.Equal(t, 42, Unwrap(c))
	})

	t.Run("creates Const with struct value", func(t *testing.T) {
		type Config struct {
			Name string
			Port int
		}
		cfg := Config{Name: "server", Port: 8080}
		c := Make[Config, bool](cfg)
		assert.Equal(t, cfg, Unwrap(c))
	})
}

// TestUnwrap tests extracting values from Const
func TestUnwrap(t *testing.T) {
	t.Run("unwraps string value", func(t *testing.T) {
		c := Make[string, int]("world")
		value := Unwrap(c)
		assert.Equal(t, "world", value)
	})

	t.Run("unwraps empty string", func(t *testing.T) {
		c := Make[string, int]("")
		value := Unwrap(c)
		assert.Equal(t, "", value)
	})

	t.Run("unwraps zero value", func(t *testing.T) {
		c := Make[int, string](0)
		value := Unwrap(c)
		assert.Equal(t, 0, value)
	})
}

// TestOf tests the Of function
func TestOf(t *testing.T) {
	t.Run("creates Const with monoid empty value", func(t *testing.T) {
		of := Of[string, int](S.Monoid)
		c := of(42)
		assert.Equal(t, "", Unwrap(c))
	})

	t.Run("ignores input value", func(t *testing.T) {
		of := Of[string, int](S.Monoid)
		c1 := of(1)
		c2 := of(100)
		assert.Equal(t, Unwrap(c1), Unwrap(c2))
	})

	t.Run("works with int monoid", func(t *testing.T) {
		of := Of[int, string](N.MonoidSum[int]())
		c := of("ignored")
		assert.Equal(t, 0, Unwrap(c))
	})
}

// TestMap tests the Map function
func TestMap(t *testing.T) {
	t.Run("preserves wrapped value", func(t *testing.T) {
		fa := Make[string, int]("foo")
		result := F.Pipe1(fa, Map[string](utils.Double))
		assert.Equal(t, "foo", Unwrap(result))
	})

	t.Run("changes phantom type", func(t *testing.T) {
		fa := Make[string, int]("data")
		fb := Map[string, int, string](strconv.Itoa)(fa)
		// Value unchanged, but type changed from Const[string, int] to Const[string, string]
		assert.Equal(t, "data", Unwrap(fb))
	})

	t.Run("function is never called", func(t *testing.T) {
		called := false
		fa := Make[string, int]("test")
		fb := Map[string, int, string](func(i int) string {
			called = true
			return strconv.Itoa(i)
		})(fa)
		assert.False(t, called, "Map function should not be called")
		assert.Equal(t, "test", Unwrap(fb))
	})
}

// TestMonadMap tests the MonadMap function
func TestMonadMap(t *testing.T) {
	t.Run("preserves wrapped value", func(t *testing.T) {
		fa := Make[string, int]("original")
		fb := MonadMap(fa, func(i int) string { return strconv.Itoa(i) })
		assert.Equal(t, "original", Unwrap(fb))
	})

	t.Run("works with different types", func(t *testing.T) {
		fa := Make[int, string](42)
		fb := MonadMap(fa, func(s string) bool { return len(s) > 0 })
		assert.Equal(t, 42, Unwrap(fb))
	})
}

// TestAp tests the Ap function
func TestAp(t *testing.T) {
	t.Run("combines string values", func(t *testing.T) {
		fab := Make[string, int]("bar")
		fa := Make[string, func(int) int]("foo")
		result := Ap[string, int, int](S.Monoid)(fab)(fa)
		assert.Equal(t, "foobar", Unwrap(result))
	})

	t.Run("combines int values with sum", func(t *testing.T) {
		fab := Make[int, string](10)
		fa := Make[int, func(string) string](5)
		result := Ap[int, string, string](N.SemigroupSum[int]())(fab)(fa)
		assert.Equal(t, 15, Unwrap(result))
	})

	t.Run("combines int values with product", func(t *testing.T) {
		fab := Make[int, bool](3)
		fa := Make[int, func(bool) bool](4)
		result := Ap[int, bool, bool](N.SemigroupProduct[int]())(fab)(fa)
		assert.Equal(t, 12, Unwrap(result))
	})
}

// TestMonadAp tests the MonadAp function
func TestMonadAp(t *testing.T) {
	t.Run("combines values using semigroup", func(t *testing.T) {
		ap := MonadAp[string, int, int](S.Monoid)
		fab := Make[string, func(int) int]("hello")
		fa := Make[string, int]("world")
		result := ap(fab, fa)
		assert.Equal(t, "helloworld", Unwrap(result))
	})

	t.Run("works with empty strings", func(t *testing.T) {
		ap := MonadAp[string, int, int](S.Monoid)
		fab := Make[string, func(int) int]("")
		fa := Make[string, int]("test")
		result := ap(fab, fa)
		assert.Equal(t, "test", Unwrap(result))
	})
}

// TestMonoid tests the Monoid function
func TestMonoid(t *testing.T) {
	t.Run("always returns constant value", func(t *testing.T) {
		m := Monoid(42)
		assert.Equal(t, 42, m.Concat(1, 2))
		assert.Equal(t, 42, m.Concat(100, 200))
		assert.Equal(t, 42, m.Empty())
	})

	t.Run("works with strings", func(t *testing.T) {
		m := Monoid("constant")
		assert.Equal(t, "constant", m.Concat("a", "b"))
		assert.Equal(t, "constant", m.Empty())
	})

	t.Run("works with structs", func(t *testing.T) {
		type Point struct{ X, Y int }
		p := Point{X: 1, Y: 2}
		m := Monoid(p)
		assert.Equal(t, p, m.Concat(Point{X: 3, Y: 4}, Point{X: 5, Y: 6}))
		assert.Equal(t, p, m.Empty())
	})

	t.Run("satisfies monoid laws", func(t *testing.T) {
		m := Monoid(10)

		// Left identity: Concat(Empty(), x) = x (both return constant)
		assert.Equal(t, 10, m.Concat(m.Empty(), 5))

		// Right identity: Concat(x, Empty()) = x (both return constant)
		assert.Equal(t, 10, m.Concat(5, m.Empty()))

		// Associativity: Concat(Concat(x, y), z) = Concat(x, Concat(y, z))
		left := m.Concat(m.Concat(1, 2), 3)
		right := m.Concat(1, m.Concat(2, 3))
		assert.Equal(t, left, right)
		assert.Equal(t, 10, left)
	})
}

// TestConstFunctorLaws tests functor laws for Const
func TestConstFunctorLaws(t *testing.T) {
	t.Run("identity law", func(t *testing.T) {
		// map id = id
		fa := Make[string, int]("test")
		mapped := Map[string, int, int](F.Identity[int])(fa)
		assert.Equal(t, Unwrap(fa), Unwrap(mapped))
	})

	t.Run("composition law", func(t *testing.T) {
		// map (g . f) = map g . map f
		fa := Make[string, int]("data")
		f := func(i int) string { return strconv.Itoa(i) }
		g := func(s string) bool { return len(s) > 0 }

		// map (g . f)
		composed := Map[string, int, bool](func(i int) bool { return g(f(i)) })(fa)

		// map g . map f
		intermediate := F.Pipe1(fa, Map[string, int, string](f))
		chained := Map[string, string, bool](g)(intermediate)

		assert.Equal(t, Unwrap(composed), Unwrap(chained))
	})
}

// TestConstApplicativeLaws tests applicative laws for Const
func TestConstApplicativeLaws(t *testing.T) {
	t.Run("identity law", func(t *testing.T) {
		// For Const, ap combines the wrapped values using the semigroup
		// ap (of id) v combines empty (from of) with v's value
		v := Make[string, int]("value")
		ofId := Of[string, func(int) int](S.Monoid)(F.Identity[int])
		result := Ap[string, int, int](S.Monoid)(v)(ofId)
		// Result combines "" (from Of) with "value" using string monoid
		assert.Equal(t, "value", Unwrap(result))
	})

	t.Run("homomorphism law", func(t *testing.T) {
		// ap (of f) (of x) = of (f x)
		f := func(i int) string { return strconv.Itoa(i) }
		x := 42

		ofF := Of[string, func(int) string](S.Monoid)(f)
		ofX := Of[string, int](S.Monoid)(x)
		left := Ap[string, int, string](S.Monoid)(ofX)(ofF)

		right := Of[string, string](S.Monoid)(f(x))

		assert.Equal(t, Unwrap(left), Unwrap(right))
	})
}

// TestConstEdgeCases tests edge cases
func TestConstEdgeCases(t *testing.T) {
	t.Run("empty string values", func(t *testing.T) {
		c := Make[string, int]("")
		assert.Equal(t, "", Unwrap(c))

		mapped := Map[string, int, string](strconv.Itoa)(c)
		assert.Equal(t, "", Unwrap(mapped))
	})

	t.Run("zero values", func(t *testing.T) {
		c := Make[int, string](0)
		assert.Equal(t, 0, Unwrap(c))
	})

	t.Run("nil pointer", func(t *testing.T) {
		var ptr *int
		c := Make[*int, string](ptr)
		assert.Nil(t, Unwrap(c))
	})

	t.Run("multiple map operations", func(t *testing.T) {
		c := Make[string, int]("original")
		// Chain multiple map operations
		step1 := Map[string, int, string](strconv.Itoa)(c)
		step2 := Map[string, string, bool](func(s string) bool { return len(s) > 0 })(step1)
		result := Map[string, bool, int](func(b bool) int {
			if b {
				return 1
			}
			return 0
		})(step2)
		assert.Equal(t, "original", Unwrap(result))
	})
}

// BenchmarkMake benchmarks the Make constructor
func BenchmarkMake(b *testing.B) {
	b.ResetTimer()
	for b.Loop() {
		_ = Make[string, int]("test")
	}
}

// BenchmarkUnwrap benchmarks the Unwrap function
func BenchmarkUnwrap(b *testing.B) {
	c := Make[string, int]("test")
	b.ResetTimer()
	for b.Loop() {
		_ = Unwrap(c)
	}
}

// BenchmarkMap benchmarks the Map function
func BenchmarkMap(b *testing.B) {
	c := Make[string, int]("test")
	mapFn := Map[string, int, string](strconv.Itoa)
	b.ResetTimer()
	for b.Loop() {
		_ = mapFn(c)
	}
}

// BenchmarkAp benchmarks the Ap function
func BenchmarkAp(b *testing.B) {
	fab := Make[string, int]("hello")
	fa := Make[string, func(int) int]("world")
	apFn := Ap[string, int, int](S.Monoid)
	b.ResetTimer()
	for b.Loop() {
		_ = apFn(fab)(fa)
	}
}

// BenchmarkMonoid benchmarks the Monoid function
func BenchmarkMonoid(b *testing.B) {
	m := Monoid(42)
	b.ResetTimer()
	for b.Loop() {
		_ = m.Concat(1, 2)
	}
}
