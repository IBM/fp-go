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

package readeroption

import (
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/utils"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/stretchr/testify/assert"
)

type MyContext string

const defaultContext MyContext = "default"

func TestMap(t *testing.T) {

	g := F.Pipe1(
		Of[MyContext](1),
		Map[MyContext](utils.Double),
	)

	assert.Equal(t, O.Of(2), g(defaultContext))

}

func TestAp(t *testing.T) {
	g := F.Pipe1(
		Of[MyContext](utils.Double),
		Ap[int](Of[MyContext](1)),
	)
	assert.Equal(t, O.Of(2), g(defaultContext))

}

func TestFlatten(t *testing.T) {

	g := F.Pipe1(
		Of[MyContext](Of[MyContext]("a")),
		Flatten[MyContext, string],
	)

	assert.Equal(t, O.Of("a"), g(defaultContext))
}

func TestFromOption(t *testing.T) {
	// Test with Some
	opt1 := O.Of(42)
	ro1 := FromOption[MyContext](opt1)
	assert.Equal(t, O.Of(42), ro1(defaultContext))

	// Test with None
	opt2 := O.None[int]()
	ro2 := FromOption[MyContext](opt2)
	assert.Equal(t, O.None[int](), ro2(defaultContext))
}

func TestSome(t *testing.T) {
	ro := Some[MyContext](42)
	assert.Equal(t, O.Of(42), ro(defaultContext))
}

func TestFromReader(t *testing.T) {
	reader := func(ctx MyContext) int {
		return 42
	}
	ro := FromReader(reader)
	assert.Equal(t, O.Of(42), ro(defaultContext))
}

func TestOf(t *testing.T) {
	ro := Of[MyContext](42)
	assert.Equal(t, O.Of(42), ro(defaultContext))
}

func TestNone(t *testing.T) {
	ro := None[MyContext, int]()
	assert.Equal(t, O.None[int](), ro(defaultContext))
}

func TestChain(t *testing.T) {
	double := func(x int) ReaderOption[MyContext, int] {
		return Of[MyContext](x * 2)
	}

	g := F.Pipe1(
		Of[MyContext](21),
		Chain(double),
	)

	assert.Equal(t, O.Of(42), g(defaultContext))

	// Test with None
	g2 := F.Pipe1(
		None[MyContext, int](),
		Chain(double),
	)
	assert.Equal(t, O.None[int](), g2(defaultContext))
}

func TestFromPredicate(t *testing.T) {
	isPositive := FromPredicate[MyContext](func(x int) bool {
		return x > 0
	})

	// Test with positive number
	g1 := F.Pipe1(
		Of[MyContext](42),
		Chain(isPositive),
	)
	assert.Equal(t, O.Of(42), g1(defaultContext))

	// Test with negative number
	g2 := F.Pipe1(
		Of[MyContext](-5),
		Chain(isPositive),
	)
	assert.Equal(t, O.None[int](), g2(defaultContext))
}

func TestFold(t *testing.T) {
	onNone := func() Reader[MyContext, string] {
		return func(ctx MyContext) string {
			return "none"
		}
	}
	onSome := func(x int) Reader[MyContext, string] {
		return func(ctx MyContext) string {
			return fmt.Sprintf("%d", x)
		}
	}

	// Test with Some
	g1 := Fold(onNone, onSome)(Of[MyContext](42))
	assert.Equal(t, "42", g1(defaultContext))

	// Test with None
	g2 := Fold(onNone, onSome)(None[MyContext, int]())
	assert.Equal(t, "none", g2(defaultContext))
}

func TestGetOrElse(t *testing.T) {
	defaultValue := func() Reader[MyContext, int] {
		return func(ctx MyContext) int {
			return 0
		}
	}

	// Test with Some
	g1 := GetOrElse(defaultValue)(Of[MyContext](42))
	assert.Equal(t, 42, g1(defaultContext))

	// Test with None
	g2 := GetOrElse(defaultValue)(None[MyContext, int]())
	assert.Equal(t, 0, g2(defaultContext))
}

func TestAsk(t *testing.T) {
	ro := Ask[MyContext, any]()
	result := ro(defaultContext)
	assert.Equal(t, O.Of(defaultContext), result)
}

func TestAsks(t *testing.T) {
	reader := func(ctx MyContext) string {
		return string(ctx)
	}
	ro := Asks(reader)
	result := ro(defaultContext)
	assert.Equal(t, O.Of("default"), result)
}

func TestChainOptionK(t *testing.T) {
	parsePositive := func(x int) O.Option[int] {
		if x > 0 {
			return O.Of(x)
		}
		return O.None[int]()
	}

	// Test with positive number
	g1 := F.Pipe1(
		Of[MyContext](42),
		ChainOptionK[MyContext](parsePositive),
	)
	assert.Equal(t, O.Of(42), g1(defaultContext))

	// Test with negative number
	g2 := F.Pipe1(
		Of[MyContext](-5),
		ChainOptionK[MyContext](parsePositive),
	)
	assert.Equal(t, O.None[int](), g2(defaultContext))
}

func TestLocal(t *testing.T) {
	type GlobalContext struct {
		Value string
	}

	// A computation that needs a string context
	ro := Asks(func(s string) string {
		return "Hello, " + s
	})

	// Transform GlobalContext to string
	transformed := Local[string](func(g GlobalContext) string {
		return g.Value
	})(ro)

	result := transformed(GlobalContext{Value: "World"})
	assert.Equal(t, O.Of("Hello, World"), result)
}

func TestRead(t *testing.T) {
	ro := Of[MyContext](42)
	result := Read[int](defaultContext)(ro)
	assert.Equal(t, O.Of(42), result)
}

func TestFlap(t *testing.T) {
	addFunc := func(x int) int {
		return x + 10
	}

	g := F.Pipe1(
		Of[MyContext](addFunc),
		Flap[MyContext, int](32),
	)

	assert.Equal(t, O.Of(42), g(defaultContext))
}
