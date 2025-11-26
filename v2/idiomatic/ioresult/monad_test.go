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

package ioresult

import (
	"errors"
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/stretchr/testify/assert"
)

// Helper function to compare IOResult values
func ioResultEqual[T comparable](a, b IOResult[T]) bool {
	valA, errA := a()
	valB, errB := b()

	if errA != nil && errB != nil {
		return errA.Error() == errB.Error()
	}
	if errA != nil || errB != nil {
		return false
	}
	return valA == valB
}

// TestPointedOf tests that Pointed().Of creates a successful IOResult
func TestPointedOf(t *testing.T) {
	t.Run("Creates successful IOResult with integer", func(t *testing.T) {
		pointed := Pointed[int]()
		result := pointed.Of(42)

		val, err := result()
		assert.NoError(t, err)
		assert.Equal(t, 42, val)
	})

	t.Run("Creates successful IOResult with string", func(t *testing.T) {
		pointed := Pointed[string]()
		result := pointed.Of("hello")

		val, err := result()
		assert.NoError(t, err)
		assert.Equal(t, "hello", val)
	})

	t.Run("Creates successful IOResult with struct", func(t *testing.T) {
		type User struct {
			Name string
			Age  int
		}
		pointed := Pointed[User]()
		user := User{Name: "Alice", Age: 30}
		result := pointed.Of(user)

		val, err := result()
		assert.NoError(t, err)
		assert.Equal(t, user, val)
	})
}

// TestFunctorMap tests that Functor().Map correctly transforms values
func TestFunctorMap(t *testing.T) {
	t.Run("Maps over successful value", func(t *testing.T) {
		functor := Functor[int, int]()
		io := Of(5)
		mapped := functor.Map(N.Mul(2))(io)

		val, err := mapped()
		assert.NoError(t, err)
		assert.Equal(t, 10, val)
	})

	t.Run("Maps over error preserves error", func(t *testing.T) {
		functor := Functor[int, int]()
		io := Left[int](errors.New("test error"))
		mapped := functor.Map(N.Mul(2))(io)

		_, err := mapped()
		assert.Error(t, err)
		assert.Equal(t, "test error", err.Error())
	})

	t.Run("Maps with type transformation", func(t *testing.T) {
		functor := Functor[int, string]()
		io := Of(42)
		mapped := functor.Map(func(x int) string { return fmt.Sprintf("value: %d", x) })(io)

		val, err := mapped()
		assert.NoError(t, err)
		assert.Equal(t, "value: 42", val)
	})
}

// TestMonadChain tests that Monad().Chain correctly chains computations
func TestMonadChain(t *testing.T) {
	t.Run("Chains successful computations", func(t *testing.T) {
		monad := Monad[int, int]()
		io := monad.Of(5)
		chained := monad.Chain(func(x int) IOResult[int] {
			return Of(x * 2)
		})(io)

		val, err := chained()
		assert.NoError(t, err)
		assert.Equal(t, 10, val)
	})

	t.Run("Chains with error in first computation", func(t *testing.T) {
		monad := Monad[int, int]()
		io := Left[int](errors.New("initial error"))
		chained := monad.Chain(func(x int) IOResult[int] {
			return Of(x * 2)
		})(io)

		_, err := chained()
		assert.Error(t, err)
		assert.Equal(t, "initial error", err.Error())
	})

	t.Run("Chains with error in second computation", func(t *testing.T) {
		monad := Monad[int, int]()
		io := monad.Of(5)
		chained := monad.Chain(func(x int) IOResult[int] {
			return Left[int](errors.New("chain error"))
		})(io)

		_, err := chained()
		assert.Error(t, err)
		assert.Equal(t, "chain error", err.Error())
	})

	t.Run("Chains with type transformation", func(t *testing.T) {
		monad := Monad[int, string]()
		io := Of(42)
		chained := monad.Chain(func(x int) IOResult[string] {
			return Of(fmt.Sprintf("value: %d", x))
		})(io)

		val, err := chained()
		assert.NoError(t, err)
		assert.Equal(t, "value: 42", val)
	})
}

// TestMonadAp tests the applicative functionality
func TestMonadAp(t *testing.T) {
	t.Run("Applies function to value", func(t *testing.T) {
		monad := Monad[int, int]()
		fn := Of(N.Mul(2))
		val := monad.Of(5)
		result := monad.Ap(val)(fn)

		res, err := result()
		assert.NoError(t, err)
		assert.Equal(t, 10, res)
	})

	t.Run("Error in function", func(t *testing.T) {
		monad := Monad[int, int]()
		fn := Left[func(int) int](errors.New("function error"))
		val := monad.Of(5)
		result := monad.Ap(val)(fn)

		_, err := result()
		assert.Error(t, err)
		assert.Equal(t, "function error", err.Error())
	})

	t.Run("Error in value", func(t *testing.T) {
		monad := Monad[int, int]()
		fn := Of(N.Mul(2))
		val := Left[int](errors.New("value error"))
		result := monad.Ap(val)(fn)

		_, err := result()
		assert.Error(t, err)
		assert.Equal(t, "value error", err.Error())
	})
}

// Monad Law Tests

// TestMonadLeftIdentity verifies: Chain(Of(a), f) == f(a)
// The left identity law states that wrapping a value with Of and then chaining
// with a function f should be the same as just applying f to the value.
func TestMonadLeftIdentity(t *testing.T) {
	t.Run("Left identity with successful function", func(t *testing.T) {
		monad := Monad[int, string]()
		a := 42
		f := func(x int) IOResult[string] {
			return Of(fmt.Sprintf("value: %d", x))
		}

		// Chain(Of(a), f)
		left := monad.Chain(f)(monad.Of(a))

		// f(a)
		right := f(a)

		// Both should produce the same result
		leftVal, leftErr := left()
		rightVal, rightErr := right()

		assert.Equal(t, rightErr, leftErr)
		assert.Equal(t, rightVal, leftVal)
	})

	t.Run("Left identity with error-returning function", func(t *testing.T) {
		monad := Monad[int, string]()
		a := -1
		f := func(x int) IOResult[string] {
			if x < 0 {
				return Left[string](errors.New("negative value"))
			}
			return Of(fmt.Sprintf("value: %d", x))
		}

		left := monad.Chain(f)(monad.Of(a))
		right := f(a)

		leftVal, leftErr := left()
		rightVal, rightErr := right()

		assert.Equal(t, rightErr, leftErr)
		assert.Equal(t, rightVal, leftVal)
	})

	t.Run("Left identity with multiple values", func(t *testing.T) {
		testCases := []int{0, 1, 42, 100, -5}
		monad := Monad[int, int]()
		f := func(x int) IOResult[int] {
			return Of(x * 2)
		}

		for _, a := range testCases {
			left := monad.Chain(f)(monad.Of(a))
			right := f(a)

			leftVal, leftErr := left()
			rightVal, rightErr := right()

			assert.Equal(t, rightErr, leftErr, "Errors should match for value %d", a)
			assert.Equal(t, rightVal, leftVal, "Values should match for value %d", a)
		}
	})
}

// TestMonadRightIdentity verifies: Chain(m, Of) == m
// The right identity law states that chaining an IOResult with Of should
// return the original IOResult unchanged.
func TestMonadRightIdentity(t *testing.T) {
	t.Run("Right identity with successful value", func(t *testing.T) {
		monad := Monad[int, int]()
		m := Of(42)

		// Chain(m, Of)
		chained := monad.Chain(func(x int) IOResult[int] {
			return monad.Of(x)
		})(m)

		// Should be equivalent to m
		mVal, mErr := m()
		chainedVal, chainedErr := chained()

		assert.Equal(t, mErr, chainedErr)
		assert.Equal(t, mVal, chainedVal)
	})

	t.Run("Right identity with error", func(t *testing.T) {
		monad := Monad[int, int]()
		m := Left[int](errors.New("test error"))

		chained := monad.Chain(func(x int) IOResult[int] {
			return monad.Of(x)
		})(m)

		mVal, mErr := m()
		chainedVal, chainedErr := chained()

		assert.Equal(t, mErr, chainedErr)
		assert.Equal(t, mVal, chainedVal)
	})

	t.Run("Right identity with different types", func(t *testing.T) {
		monadStr := Monad[string, string]()
		m := Of("hello")

		chained := monadStr.Chain(func(x string) IOResult[string] {
			return monadStr.Of(x)
		})(m)

		mVal, mErr := m()
		chainedVal, chainedErr := chained()

		assert.Equal(t, mErr, chainedErr)
		assert.Equal(t, mVal, chainedVal)
	})
}

// TestMonadAssociativity verifies: Chain(Chain(m, f), g) == Chain(m, x => Chain(f(x), g))
// The associativity law states that the order of nesting chains doesn't matter.
func TestMonadAssociativity(t *testing.T) {
	t.Run("Associativity with successful computations", func(t *testing.T) {
		monadIntInt := Monad[int, int]()
		monadIntStr := Monad[int, string]()

		m := Of(5)
		f := func(x int) IOResult[int] {
			return Of(x * 2)
		}
		g := func(y int) IOResult[string] {
			return Of(fmt.Sprintf("result: %d", y))
		}

		// Chain(Chain(m, f), g)
		left := monadIntStr.Chain(g)(monadIntInt.Chain(f)(m))

		// Chain(m, x => Chain(f(x), g))
		right := monadIntStr.Chain(func(x int) IOResult[string] {
			return monadIntStr.Chain(g)(f(x))
		})(m)

		leftVal, leftErr := left()
		rightVal, rightErr := right()

		assert.Equal(t, rightErr, leftErr)
		assert.Equal(t, rightVal, leftVal)
	})

	t.Run("Associativity with error in first function", func(t *testing.T) {
		monadIntInt := Monad[int, int]()
		monadIntStr := Monad[int, string]()

		m := Of(5)
		f := func(x int) IOResult[int] {
			return Left[int](errors.New("error in f"))
		}
		g := func(y int) IOResult[string] {
			return Of(fmt.Sprintf("result: %d", y))
		}

		left := monadIntStr.Chain(g)(monadIntInt.Chain(f)(m))
		right := monadIntStr.Chain(func(x int) IOResult[string] {
			return monadIntStr.Chain(g)(f(x))
		})(m)

		leftVal, leftErr := left()
		rightVal, rightErr := right()

		assert.Equal(t, rightErr, leftErr)
		assert.Equal(t, rightVal, leftVal)
	})

	t.Run("Associativity with error in second function", func(t *testing.T) {
		monadIntInt := Monad[int, int]()
		monadIntStr := Monad[int, string]()

		m := Of(5)
		f := func(x int) IOResult[int] {
			return Of(x * 2)
		}
		g := func(y int) IOResult[string] {
			return Left[string](errors.New("error in g"))
		}

		left := monadIntStr.Chain(g)(monadIntInt.Chain(f)(m))
		right := monadIntStr.Chain(func(x int) IOResult[string] {
			return monadIntStr.Chain(g)(f(x))
		})(m)

		leftVal, leftErr := left()
		rightVal, rightErr := right()

		assert.Equal(t, rightErr, leftErr)
		assert.Equal(t, rightVal, leftVal)
	})

	t.Run("Associativity with complex chain", func(t *testing.T) {
		monad1 := Monad[int, int]()
		monad2 := Monad[int, int]()

		m := Of(2)
		f := func(x int) IOResult[int] { return Of(x + 3) }
		g := func(y int) IOResult[int] { return Of(y * 4) }

		// (2 + 3) * 4 = 20
		left := monad2.Chain(g)(monad1.Chain(f)(m))
		right := monad1.Chain(func(x int) IOResult[int] {
			return monad2.Chain(g)(f(x))
		})(m)

		leftVal, leftErr := left()
		rightVal, rightErr := right()

		assert.NoError(t, leftErr)
		assert.NoError(t, rightErr)
		assert.Equal(t, 20, leftVal)
		assert.Equal(t, 20, rightVal)
	})
}

// TestFunctorComposition verifies: Map(f . g) == Map(f) . Map(g)
// The functor composition law states that mapping a composition of functions
// should be the same as composing the maps of those functions.
func TestFunctorComposition(t *testing.T) {
	t.Run("Functor composition law", func(t *testing.T) {
		functor1 := Functor[int, int]()
		functor2 := Functor[int, string]()

		m := Of(5)
		f := N.Mul(2)
		g := func(x int) string { return fmt.Sprintf("value: %d", x) }

		// Map(g . f)
		composed := functor2.Map(F.Flow2(f, g))(m)

		// Map(g) . Map(f)
		separate := functor2.Map(g)(functor1.Map(f)(m))

		composedVal, composedErr := composed()
		separateVal, separateErr := separate()

		assert.Equal(t, composedErr, separateErr)
		assert.Equal(t, composedVal, separateVal)
	})

	t.Run("Functor composition with error", func(t *testing.T) {
		functor1 := Functor[int, int]()
		functor2 := Functor[int, string]()

		m := Left[int](errors.New("test error"))
		f := N.Mul(2)
		g := func(x int) string { return fmt.Sprintf("value: %d", x) }

		composed := functor2.Map(F.Flow2(f, g))(m)
		separate := functor2.Map(g)(functor1.Map(f)(m))

		composedVal, composedErr := composed()
		separateVal, separateErr := separate()

		assert.Equal(t, composedErr, separateErr)
		assert.Equal(t, composedVal, separateVal)
	})
}

// TestFunctorIdentity verifies: Map(id) == id
// The functor identity law states that mapping the identity function
// should return the original IOResult unchanged.
func TestFunctorIdentity(t *testing.T) {
	t.Run("Functor identity with successful value", func(t *testing.T) {
		functor := Functor[int, int]()
		m := Of(42)

		// Map(id)
		mapped := functor.Map(F.Identity[int])(m)

		mVal, mErr := m()
		mappedVal, mappedErr := mapped()

		assert.Equal(t, mErr, mappedErr)
		assert.Equal(t, mVal, mappedVal)
	})

	t.Run("Functor identity with error", func(t *testing.T) {
		functor := Functor[int, int]()
		m := Left[int](errors.New("test error"))

		mapped := functor.Map(F.Identity[int])(m)

		mVal, mErr := m()
		mappedVal, mappedErr := mapped()

		assert.Equal(t, mErr, mappedErr)
		assert.Equal(t, mVal, mappedVal)
	})
}

// TestMonadParVsSeq tests that MonadPar and MonadSeq produce the same results
func TestMonadParVsSeq(t *testing.T) {
	t.Run("Par and Seq produce same results for Map", func(t *testing.T) {
		monadPar := MonadPar[int, int]()
		monadSeq := MonadSeq[int, int]()

		io := Of(5)
		f := N.Mul(2)

		par := monadPar.Map(f)(io)
		seq := monadSeq.Map(f)(io)

		parVal, parErr := par()
		seqVal, seqErr := seq()

		assert.Equal(t, parErr, seqErr)
		assert.Equal(t, parVal, seqVal)
	})

	t.Run("Par and Seq produce same results for Chain", func(t *testing.T) {
		monadPar := MonadPar[int, string]()
		monadSeq := MonadSeq[int, string]()

		io := Of(42)
		f := func(x int) IOResult[string] {
			return Of(fmt.Sprintf("value: %d", x))
		}

		par := monadPar.Chain(f)(io)
		seq := monadSeq.Chain(f)(io)

		parVal, parErr := par()
		seqVal, seqErr := seq()

		assert.Equal(t, parErr, seqErr)
		assert.Equal(t, parVal, seqVal)
	})

	t.Run("Default Monad uses parallel execution", func(t *testing.T) {
		monadDefault := Monad[int, int]()
		monadPar := MonadPar[int, int]()

		io := Of(5)
		f := N.Mul(2)

		def := monadDefault.Map(f)(io)
		par := monadPar.Map(f)(io)

		defVal, defErr := def()
		parVal, parErr := par()

		assert.Equal(t, parErr, defErr)
		assert.Equal(t, parVal, defVal)
	})
}

// TestMonadIntegration tests complete workflows using the monad interface
func TestMonadIntegration(t *testing.T) {
	t.Run("Complex pipeline using monad operations", func(t *testing.T) {
		monad1 := Monad[int, int]()
		monad2 := Monad[int, string]()

		// Build a pipeline: multiply by 2, add 3, then format
		result := F.Pipe2(
			monad1.Of(5),
			monad1.Map(N.Mul(2)),
			monad1.Chain(func(x int) IOResult[int] {
				return Of(x + 3)
			}),
		)

		// Continue with type change
		formatted := monad2.Map(func(x int) string {
			return fmt.Sprintf("Final: %d", x)
		})(result)

		val, err := formatted()
		assert.NoError(t, err)
		assert.Equal(t, "Final: 13", val) // (5 * 2) + 3 = 13
	})

	t.Run("Error handling in complex pipeline", func(t *testing.T) {
		monad1 := Monad[int, int]()
		monad2 := Monad[int, string]()

		result := F.Pipe2(
			monad1.Of(5),
			monad1.Map(N.Mul(2)),
			monad1.Chain(func(x int) IOResult[int] {
				if x > 5 {
					return Left[int](errors.New("value too large"))
				}
				return Of(x + 3)
			}),
		)

		formatted := monad2.Map(func(x int) string {
			return fmt.Sprintf("Final: %d", x)
		})(result)

		_, err := formatted()
		assert.Error(t, err)
		assert.Equal(t, "value too large", err.Error())
	})
}
