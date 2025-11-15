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

package erasure

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	"github.com/stretchr/testify/assert"
)

func TestErase(t *testing.T) {
	t.Run("erases int value", func(t *testing.T) {
		value := 42
		erased := Erase(value)
		assert.NotNil(t, erased)
		// Verify it's a pointer to int
		ptr, ok := erased.(*int)
		assert.True(t, ok)
		assert.Equal(t, 42, *ptr)
	})

	t.Run("erases string value", func(t *testing.T) {
		value := "hello"
		erased := Erase(value)
		assert.NotNil(t, erased)
		ptr, ok := erased.(*string)
		assert.True(t, ok)
		assert.Equal(t, "hello", *ptr)
	})

	t.Run("erases struct value", func(t *testing.T) {
		type Person struct {
			Name string
			Age  int
		}
		value := Person{Name: "Alice", Age: 30}
		erased := Erase(value)
		assert.NotNil(t, erased)
		ptr, ok := erased.(*Person)
		assert.True(t, ok)
		assert.Equal(t, "Alice", ptr.Name)
		assert.Equal(t, 30, ptr.Age)
	})
}

func TestUnerase(t *testing.T) {
	t.Run("unerases int value", func(t *testing.T) {
		erased := Erase(42)
		value := Unerase[int](erased)
		assert.Equal(t, 42, value)
	})

	t.Run("unerases string value", func(t *testing.T) {
		erased := Erase("hello")
		value := Unerase[string](erased)
		assert.Equal(t, "hello", value)
	})

	t.Run("unerases complex type", func(t *testing.T) {
		type Data struct {
			Values []int
			Label  string
		}
		original := Data{Values: []int{1, 2, 3}, Label: "test"}
		erased := Erase(original)
		value := Unerase[Data](erased)
		assert.Equal(t, original, value)
	})
}

func TestSafeUnerase(t *testing.T) {
	t.Run("successfully unerases correct type", func(t *testing.T) {
		erased := Erase(42)
		result := SafeUnerase[int](erased)
		assert.True(t, E.IsRight(result))
		value := E.GetOrElse(func(error) int { return 0 })(result)
		assert.Equal(t, 42, value)
	})

	t.Run("returns error for wrong type", func(t *testing.T) {
		erased := Erase(42)
		result := SafeUnerase[string](erased)
		assert.True(t, E.IsLeft(result))
	})

	t.Run("returns error for non-erased value", func(t *testing.T) {
		notErased := "plain string"
		result := SafeUnerase[string](notErased)
		assert.True(t, E.IsLeft(result))
	})

	t.Run("successfully unerases string", func(t *testing.T) {
		erased := Erase("hello")
		result := SafeUnerase[string](erased)
		assert.True(t, E.IsRight(result))
		value := E.GetOrElse(func(error) string { return "" })(result)
		assert.Equal(t, "hello", value)
	})

	t.Run("successfully unerases complex type", func(t *testing.T) {
		type Config struct {
			Host string
			Port int
		}
		original := Config{Host: "localhost", Port: 8080}
		erased := Erase(original)
		result := SafeUnerase[Config](erased)
		assert.True(t, E.IsRight(result))
		value := E.GetOrElse(func(error) Config { return Config{} })(result)
		assert.Equal(t, original, value)
	})
}

func TestErase0(t *testing.T) {
	t.Run("erases nullary function returning int", func(t *testing.T) {
		typedFunc := func() int { return 42 }
		erasedFunc := Erase0(typedFunc)
		result := erasedFunc()
		assert.NotNil(t, result)
		value := Unerase[int](result)
		assert.Equal(t, 42, value)
	})

	t.Run("erases nullary function returning string", func(t *testing.T) {
		typedFunc := func() string { return "hello" }
		erasedFunc := Erase0(typedFunc)
		result := erasedFunc()
		assert.NotNil(t, result)
		value := Unerase[string](result)
		assert.Equal(t, "hello", value)
	})

	t.Run("erases nullary function returning struct", func(t *testing.T) {
		type Result struct {
			Success bool
			Message string
		}
		typedFunc := func() Result {
			return Result{Success: true, Message: "OK"}
		}
		erasedFunc := Erase0(typedFunc)
		result := erasedFunc()
		assert.NotNil(t, result)
		value := Unerase[Result](result)
		assert.True(t, value.Success)
		assert.Equal(t, "OK", value.Message)
	})
}

func TestErase1(t *testing.T) {
	t.Run("erases unary function int to string", func(t *testing.T) {
		typedFunc := strconv.Itoa
		erasedFunc := Erase1(typedFunc)
		result := erasedFunc(Erase(42))
		assert.NotNil(t, result)
		value := Unerase[string](result)
		assert.Equal(t, "42", value)
	})

	t.Run("erases unary function string to upper", func(t *testing.T) {
		typedFunc := strings.ToUpper
		erasedFunc := Erase1(typedFunc)
		result := erasedFunc(Erase("hello"))
		assert.NotNil(t, result)
		value := Unerase[string](result)
		assert.Equal(t, "HELLO", value)
	})

	t.Run("erases unary function with complex types", func(t *testing.T) {
		type Input struct {
			Value int
		}
		type Output struct {
			Result int
		}
		typedFunc := func(in Input) Output {
			return Output{Result: in.Value * 2}
		}
		erasedFunc := Erase1(typedFunc)
		result := erasedFunc(Erase(Input{Value: 21}))
		assert.NotNil(t, result)
		value := Unerase[Output](result)
		assert.Equal(t, 42, value.Result)
	})
}

func TestErase2(t *testing.T) {
	t.Run("erases binary function int addition", func(t *testing.T) {
		typedFunc := func(x, y int) int { return x + y }
		erasedFunc := Erase2(typedFunc)
		result := erasedFunc(Erase(10), Erase(32))
		assert.NotNil(t, result)
		value := Unerase[int](result)
		assert.Equal(t, 42, value)
	})

	t.Run("erases binary function string concatenation", func(t *testing.T) {
		typedFunc := func(x, y string) string { return x + y }
		erasedFunc := Erase2(typedFunc)
		result := erasedFunc(Erase("hello"), Erase(" world"))
		assert.NotNil(t, result)
		value := Unerase[string](result)
		assert.Equal(t, "hello world", value)
	})

	t.Run("erases binary function with different types", func(t *testing.T) {
		typedFunc := func(x int, y string) string {
			return fmt.Sprintf("%s: %d", y, x)
		}
		erasedFunc := Erase2(typedFunc)
		result := erasedFunc(Erase(42), Erase("answer"))
		assert.NotNil(t, result)
		value := Unerase[string](result)
		assert.Equal(t, "answer: 42", value)
	})

	t.Run("erases binary function with complex types", func(t *testing.T) {
		type Point struct {
			X, Y int
		}
		typedFunc := func(p1, p2 Point) Point {
			return Point{X: p1.X + p2.X, Y: p1.Y + p2.Y}
		}
		erasedFunc := Erase2(typedFunc)
		result := erasedFunc(Erase(Point{X: 1, Y: 2}), Erase(Point{X: 3, Y: 4}))
		assert.NotNil(t, result)
		value := Unerase[Point](result)
		assert.Equal(t, 4, value.X)
		assert.Equal(t, 6, value.Y)
	})
}

func TestEither(t *testing.T) {
	t.Run("integration test with Either", func(t *testing.T) {
		e1 := F.Pipe3(
			E.Of[error](Erase("Carsten")),
			E.Map[error](Erase1(strings.ToUpper)),
			E.GetOrElse(func(e error) any {
				return Erase("Error")
			}),
			Unerase[string],
		)

		assert.Equal(t, "CARSTEN", e1)
	})

	t.Run("integration test with Either and SafeUnerase", func(t *testing.T) {
		erased := Erase(42)
		result := F.Pipe1(
			SafeUnerase[int](erased),
			E.Map[error](N.Mul(2)),
		)

		assert.True(t, E.IsRight(result))
		value := E.GetOrElse(func(error) int { return 0 })(result)
		assert.Equal(t, 84, value)
	})

	t.Run("integration test with Either error case", func(t *testing.T) {
		erased := Erase(42)
		result := SafeUnerase[string](erased) // Wrong type

		assert.True(t, E.IsLeft(result))
	})
}

func TestRoundTrip(t *testing.T) {
	t.Run("round trip with int", func(t *testing.T) {
		original := 42
		erased := Erase(original)
		recovered := Unerase[int](erased)
		assert.Equal(t, original, recovered)
	})

	t.Run("round trip with string", func(t *testing.T) {
		original := "hello world"
		erased := Erase(original)
		recovered := Unerase[string](erased)
		assert.Equal(t, original, recovered)
	})

	t.Run("round trip with slice", func(t *testing.T) {
		original := []int{1, 2, 3, 4, 5}
		erased := Erase(original)
		recovered := Unerase[[]int](erased)
		assert.Equal(t, original, recovered)
	})

	t.Run("round trip with map", func(t *testing.T) {
		original := map[string]int{"a": 1, "b": 2}
		erased := Erase(original)
		recovered := Unerase[map[string]int](erased)
		assert.Equal(t, original, recovered)
	})
}
