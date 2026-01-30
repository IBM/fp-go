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

package readeriooption

import (
	"context"
	"fmt"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	RIO "github.com/IBM/fp-go/v2/readerio"
	"github.com/stretchr/testify/assert"
)

func TestOf(t *testing.T) {
	ro := Of[context.Context](42)
	result := ro(context.Background())()
	assert.Equal(t, O.Of(42), result)
}

func TestSome(t *testing.T) {
	ro := Some[context.Context](42)
	result := ro(context.Background())()
	assert.Equal(t, O.Of(42), result)
}

func TestNone(t *testing.T) {
	ro := None[context.Context, int]()
	result := ro(context.Background())()
	assert.Equal(t, O.None[int](), result)
}

func TestFromOption_Some(t *testing.T) {
	opt := O.Of(42)
	ro := FromOption[context.Context](opt)
	result := ro(context.Background())()
	assert.Equal(t, O.Of(42), result)
}

func TestFromOption_None(t *testing.T) {
	opt := O.None[int]()
	ro := FromOption[context.Context](opt)
	result := ro(context.Background())()
	assert.Equal(t, O.None[int](), result)
}

func TestFromReader(t *testing.T) {
	type Config struct {
		Value int
	}

	r := func(cfg Config) int {
		return cfg.Value * 2
	}

	ro := FromReader[Config](r)
	cfg := Config{Value: 21}
	result := ro(cfg)()

	assert.Equal(t, O.Of(42), result)
}

func TestSomeReader(t *testing.T) {
	type Config struct {
		Value int
	}

	r := func(cfg Config) int {
		return cfg.Value * 2
	}

	ro := SomeReader[Config](r)
	cfg := Config{Value: 21}
	result := ro(cfg)()

	assert.Equal(t, O.Of(42), result)
}

func TestMonadMap_Some(t *testing.T) {
	ro := Of[context.Context](21)
	result := MonadMap(ro, func(x int) int { return x * 2 })
	assert.Equal(t, O.Of(42), result(context.Background())())
}

func TestMonadMap_None(t *testing.T) {
	ro := None[context.Context, int]()
	result := MonadMap(ro, func(x int) int { return x * 2 })
	assert.Equal(t, O.None[int](), result(context.Background())())
}

func TestMap_Some(t *testing.T) {
	result := F.Pipe1(
		Of[context.Context](21),
		Map[context.Context](func(x int) int { return x * 2 }),
	)
	assert.Equal(t, O.Of(42), result(context.Background())())
}

func TestMap_None(t *testing.T) {
	result := F.Pipe1(
		None[context.Context, int](),
		Map[context.Context](func(x int) int { return x * 2 }),
	)
	assert.Equal(t, O.None[int](), result(context.Background())())
}

func TestMonadChain_BothSome(t *testing.T) {
	ro1 := Of[context.Context](21)
	ro2 := func(x int) ReaderIOOption[context.Context, int] {
		return Of[context.Context](x * 2)
	}

	result := MonadChain(ro1, ro2)
	assert.Equal(t, O.Of(42), result(context.Background())())
}

func TestMonadChain_FirstNone(t *testing.T) {
	ro1 := None[context.Context, int]()
	ro2 := func(x int) ReaderIOOption[context.Context, int] {
		return Of[context.Context](x * 2)
	}

	result := MonadChain(ro1, ro2)
	assert.Equal(t, O.None[int](), result(context.Background())())
}

func TestMonadChain_SecondNone(t *testing.T) {
	ro1 := Of[context.Context](21)
	ro2 := func(x int) ReaderIOOption[context.Context, int] {
		return None[context.Context, int]()
	}

	result := MonadChain(ro1, ro2)
	assert.Equal(t, O.None[int](), result(context.Background())())
}

func TestChain(t *testing.T) {
	result := F.Pipe1(
		Of[context.Context](21),
		Chain(func(x int) ReaderIOOption[context.Context, int] {
			return Of[context.Context](x * 2)
		}),
	)
	assert.Equal(t, O.Of(42), result(context.Background())())
}

func TestMonadAp_BothSome(t *testing.T) {
	fab := Of[context.Context](func(x int) int { return x * 2 })
	fa := Of[context.Context](21)

	result := MonadAp(fab, fa)
	assert.Equal(t, O.Of(42), result(context.Background())())
}

func TestMonadAp_FunctionNone(t *testing.T) {
	fab := None[context.Context, func(int) int]()
	fa := Of[context.Context](21)

	result := MonadAp(fab, fa)
	assert.Equal(t, O.None[int](), result(context.Background())())
}

func TestMonadAp_ValueNone(t *testing.T) {
	fab := Of[context.Context](func(x int) int { return x * 2 })
	fa := None[context.Context, int]()

	result := MonadAp(fab, fa)
	assert.Equal(t, O.None[int](), result(context.Background())())
}

func TestAp(t *testing.T) {
	fa := Of[context.Context](21)
	result := F.Pipe1(
		Of[context.Context](func(x int) int { return x * 2 }),
		Ap[int](fa),
	)
	assert.Equal(t, O.Of(42), result(context.Background())())
}

func TestFromPredicate_True(t *testing.T) {
	isPositive := FromPredicate[context.Context](func(x int) bool { return x > 0 })

	result := F.Pipe1(
		Of[context.Context](42),
		Chain(isPositive),
	)

	assert.Equal(t, O.Of(42), result(context.Background())())
}

func TestFromPredicate_False(t *testing.T) {
	isPositive := FromPredicate[context.Context](func(x int) bool { return x > 0 })

	result := F.Pipe1(
		Of[context.Context](-42),
		Chain(isPositive),
	)

	assert.Equal(t, O.None[int](), result(context.Background())())
}

func TestAsk(t *testing.T) {
	type Config struct {
		Value int
	}

	ro := Ask[Config]()
	cfg := Config{Value: 42}
	result := ro(cfg)()

	assert.Equal(t, O.Of(cfg), result)
}

func TestAsks(t *testing.T) {
	type Config struct {
		Value int
	}

	ro := Asks(func(cfg Config) int {
		return cfg.Value * 2
	})

	cfg := Config{Value: 21}
	result := ro(cfg)()

	assert.Equal(t, O.Of(42), result)
}

func TestMonadChainOptionK_Some(t *testing.T) {
	parsePositive := func(x int) O.Option[int] {
		if x > 0 {
			return O.Of(x)
		}
		return O.None[int]()
	}

	result := MonadChainOptionK(
		Of[context.Context](42),
		parsePositive,
	)

	assert.Equal(t, O.Of(42), result(context.Background())())
}

func TestMonadChainOptionK_None(t *testing.T) {
	parsePositive := func(x int) O.Option[int] {
		if x > 0 {
			return O.Of(x)
		}
		return O.None[int]()
	}

	result := MonadChainOptionK(
		Of[context.Context](-42),
		parsePositive,
	)

	assert.Equal(t, O.None[int](), result(context.Background())())
}

func TestChainOptionK(t *testing.T) {
	parsePositive := func(x int) O.Option[int] {
		if x > 0 {
			return O.Of(x)
		}
		return O.None[int]()
	}

	result := F.Pipe1(
		Of[context.Context](42),
		ChainOptionK[context.Context](parsePositive),
	)

	assert.Equal(t, O.Of(42), result(context.Background())())
}

func TestFlatten(t *testing.T) {
	nested := Of[context.Context](Of[context.Context](42))
	result := Flatten(nested)
	assert.Equal(t, O.Of(42), result(context.Background())())
}

func TestLocal(t *testing.T) {
	type GlobalConfig struct {
		Factor int
	}

	type LocalConfig struct {
		Multiplier int
	}

	// Computation that needs LocalConfig
	computation := func(cfg LocalConfig) IOOption[int] {
		return func() O.Option[int] {
			return O.Of(10 * cfg.Multiplier)
		}
	}

	// Adapt to work with GlobalConfig
	adapted := Local[int](func(g GlobalConfig) LocalConfig {
		return LocalConfig{Multiplier: g.Factor}
	})(computation)

	globalCfg := GlobalConfig{Factor: 5}
	result := adapted(globalCfg)()

	assert.Equal(t, O.Of(50), result)
}

func TestRead(t *testing.T) {
	type Config struct {
		Value int
	}

	ro := func(cfg Config) IOOption[int] {
		return func() O.Option[int] {
			return O.Of(cfg.Value * 2)
		}
	}

	cfg := Config{Value: 21}
	result := Read[int](cfg)(ro)()

	assert.Equal(t, O.Of(42), result)
}

func TestMonadFlap(t *testing.T) {
	fab := Of[context.Context](func(x int) int { return x * 2 })
	result := MonadFlap(fab, 21)
	assert.Equal(t, O.Of(42), result(context.Background())())
}

func TestFlap(t *testing.T) {
	result := F.Pipe1(
		Of[context.Context](func(x int) int { return x * 2 }),
		Flap[context.Context, int](21),
	)
	assert.Equal(t, O.Of(42), result(context.Background())())
}

func TestMonadAlt_FirstSome(t *testing.T) {
	first := Of[context.Context](42)
	second := func() ReaderIOOption[context.Context, int] {
		return Of[context.Context](100)
	}

	result := MonadAlt(first, second)
	assert.Equal(t, O.Of(42), result(context.Background())())
}

func TestMonadAlt_FirstNone(t *testing.T) {
	first := None[context.Context, int]()
	second := func() ReaderIOOption[context.Context, int] {
		return Of[context.Context](100)
	}

	result := MonadAlt(first, second)
	assert.Equal(t, O.Of(100), result(context.Background())())
}

func TestMonadAlt_BothNone(t *testing.T) {
	first := None[context.Context, int]()
	second := func() ReaderIOOption[context.Context, int] {
		return None[context.Context, int]()
	}

	result := MonadAlt(first, second)
	assert.Equal(t, O.None[int](), result(context.Background())())
}

func TestAlt(t *testing.T) {
	result := F.Pipe1(
		None[context.Context, int](),
		Alt(func() ReaderIOOption[context.Context, int] {
			return Of[context.Context](42)
		}),
	)
	assert.Equal(t, O.Of(42), result(context.Background())())
}

func TestGetOrElse_Some(t *testing.T) {
	ro := Of[context.Context](42)
	result := MonadFold(ro, RIO.Of[context.Context](100), func(x int) RIO.ReaderIO[context.Context, int] {
		return RIO.Of[context.Context](x)
	})(context.Background())()
	assert.Equal(t, 42, result)
}

func TestGetOrElse_None(t *testing.T) {
	ro := None[context.Context, int]()
	result := MonadFold(ro, RIO.Of[context.Context](100), func(x int) RIO.ReaderIO[context.Context, int] {
		return RIO.Of[context.Context](x)
	})(context.Background())()
	assert.Equal(t, 100, result)
}

func TestMonadFold_Some(t *testing.T) {
	ro := Of[context.Context](42)
	result := MonadFold(
		ro,
		RIO.Of[context.Context]("none"),
		func(x int) RIO.ReaderIO[context.Context, string] {
			return RIO.Of[context.Context]("value: " + fmt.Sprintf("%d", x))
		},
	)(context.Background())()
	assert.Equal(t, "value: 42", result)
}

func TestMonadFold_None(t *testing.T) {
	ro := None[context.Context, int]()
	result := MonadFold(
		ro,
		RIO.Of[context.Context]("none"),
		func(x int) RIO.ReaderIO[context.Context, string] {
			return RIO.Of[context.Context]("value: " + fmt.Sprintf("%d", x))
		},
	)(context.Background())()
	assert.Equal(t, "none", result)
}

func TestComplexChain(t *testing.T) {
	// Test a complex chain of operations
	type Config struct {
		Factor int
	}

	result := F.Pipe3(
		Of[Config](10),
		Map[Config](func(x int) int { return x * 2 }), // 20
		Chain(func(x int) ReaderIOOption[Config, int] {
			return Asks(func(cfg Config) int {
				return x * cfg.Factor
			})
		}),
		Chain(func(x int) ReaderIOOption[Config, int] {
			if x > 50 {
				return Of[Config](x)
			}
			return None[Config, int]()
		}),
	)

	cfg := Config{Factor: 5}
	assert.Equal(t, O.Of(100), result(cfg)())

	cfg2 := Config{Factor: 2}
	assert.Equal(t, O.None[int](), result(cfg2)())
}
