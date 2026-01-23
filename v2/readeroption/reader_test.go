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
	"github.com/IBM/fp-go/v2/lazy"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/stretchr/testify/assert"
)

// Test context type
type Config struct {
	Host    string
	Port    int
	Timeout int
}

var defaultConfig = Config{
	Host:    "localhost",
	Port:    8080,
	Timeout: 30,
}

// TestOf tests the Of function which wraps a value in a ReaderOption
func TestOf(t *testing.T) {
	ro := Of[Config](42)
	result := ro(defaultConfig)
	assert.Equal(t, O.Some(42), result)
}

// TestSome tests the Some function which is an alias for Of
func TestSome(t *testing.T) {
	ro := Some[Config](42)
	result := ro(defaultConfig)
	assert.Equal(t, O.Some(42), result)
}

// TestNone tests the None function which creates a ReaderOption representing no value
func TestNone(t *testing.T) {
	ro := None[Config, int]()
	result := ro(defaultConfig)
	assert.Equal(t, O.None[int](), result)
}

// TestFromOption tests lifting an Option into a ReaderOption
func TestFromOption(t *testing.T) {
	t.Run("Some value", func(t *testing.T) {
		opt := O.Some(42)
		ro := FromOption[Config](opt)
		result := ro(defaultConfig)
		assert.Equal(t, O.Some(42), result)
	})

	t.Run("None value", func(t *testing.T) {
		opt := O.None[int]()
		ro := FromOption[Config](opt)
		result := ro(defaultConfig)
		assert.Equal(t, O.None[int](), result)
	})
}

// TestFromReader tests lifting a Reader into a ReaderOption
func TestFromReader(t *testing.T) {
	r := reader.Of[Config](42)
	ro := FromReader(r)
	result := ro(defaultConfig)
	assert.Equal(t, O.Some(42), result)
}

// TestSomeReader tests lifting a Reader into a ReaderOption (alias for FromReader)
func TestSomeReader(t *testing.T) {
	r := reader.Of[Config](42)
	ro := SomeReader(r)
	result := ro(defaultConfig)
	assert.Equal(t, O.Some(42), result)
}

// TestMonadMap tests applying a function to the value inside a ReaderOption
func TestMonadMap(t *testing.T) {
	t.Run("Map over Some", func(t *testing.T) {
		ro := Of[Config](21)
		mapped := MonadMap(ro, utils.Double)
		result := mapped(defaultConfig)
		assert.Equal(t, O.Some(42), result)
	})

	t.Run("Map over None", func(t *testing.T) {
		ro := None[Config, int]()
		mapped := MonadMap(ro, utils.Double)
		result := mapped(defaultConfig)
		assert.Equal(t, O.None[int](), result)
	})
}

// TestMap tests the curried version of MonadMap
func TestMap(t *testing.T) {
	t.Run("Map over Some", func(t *testing.T) {
		result := F.Pipe1(
			Of[Config](21),
			Map[Config](utils.Double),
		)
		assert.Equal(t, O.Some(42), result(defaultConfig))
	})

	t.Run("Map over None", func(t *testing.T) {
		result := F.Pipe1(
			None[Config, int](),
			Map[Config](utils.Double),
		)
		assert.Equal(t, O.None[int](), result(defaultConfig))
	})
}

// TestMonadChain tests sequencing two ReaderOption computations
func TestMonadChain(t *testing.T) {
	t.Run("Chain with Some", func(t *testing.T) {
		ro := Of[Config](21)
		chained := MonadChain(ro, func(x int) ReaderOption[Config, int] {
			return Of[Config](x * 2)
		})
		result := chained(defaultConfig)
		assert.Equal(t, O.Some(42), result)
	})

	t.Run("Chain with None", func(t *testing.T) {
		ro := None[Config, int]()
		chained := MonadChain(ro, func(x int) ReaderOption[Config, int] {
			return Of[Config](x * 2)
		})
		result := chained(defaultConfig)
		assert.Equal(t, O.None[int](), result)
	})

	t.Run("Chain returning None", func(t *testing.T) {
		ro := Of[Config](21)
		chained := MonadChain(ro, func(x int) ReaderOption[Config, int] {
			return None[Config, int]()
		})
		result := chained(defaultConfig)
		assert.Equal(t, O.None[int](), result)
	})
}

// TestChain tests the curried version of MonadChain
func TestChain(t *testing.T) {
	t.Run("Chain with Some", func(t *testing.T) {
		result := F.Pipe1(
			Of[Config](21),
			Chain(func(x int) ReaderOption[Config, int] {
				return Of[Config](x * 2)
			}),
		)
		assert.Equal(t, O.Some(42), result(defaultConfig))
	})

	t.Run("Chain with None", func(t *testing.T) {
		result := F.Pipe1(
			None[Config, int](),
			Chain(func(x int) ReaderOption[Config, int] {
				return Of[Config](x * 2)
			}),
		)
		assert.Equal(t, O.None[int](), result(defaultConfig))
	})
}

// TestMonadAp tests applying a function wrapped in a ReaderOption
func TestMonadAp(t *testing.T) {
	t.Run("Ap with both Some", func(t *testing.T) {
		fab := Of[Config](utils.Double)
		fa := Of[Config](21)
		result := MonadAp(fab, fa)
		assert.Equal(t, O.Some(42), result(defaultConfig))
	})

	t.Run("Ap with None function", func(t *testing.T) {
		fab := None[Config, func(int) int]()
		fa := Of[Config](21)
		result := MonadAp(fab, fa)
		assert.Equal(t, O.None[int](), result(defaultConfig))
	})

	t.Run("Ap with None value", func(t *testing.T) {
		fab := Of[Config](utils.Double)
		fa := None[Config, int]()
		result := MonadAp(fab, fa)
		assert.Equal(t, O.None[int](), result(defaultConfig))
	})
}

// TestAp tests the curried version of MonadAp
func TestAp(t *testing.T) {
	t.Run("Ap with both Some", func(t *testing.T) {
		result := F.Pipe1(
			Of[Config](utils.Double),
			Ap[int](Of[Config](21)),
		)
		assert.Equal(t, O.Some(42), result(defaultConfig))
	})
}

// TestFromPredicate tests creating a Kleisli arrow that filters based on a predicate
func TestFromPredicate(t *testing.T) {
	isPositive := FromPredicate[Config](func(x int) bool { return x > 0 })

	t.Run("Predicate satisfied", func(t *testing.T) {
		result := F.Pipe1(
			Of[Config](42),
			Chain(isPositive),
		)
		assert.Equal(t, O.Some(42), result(defaultConfig))
	})

	t.Run("Predicate not satisfied", func(t *testing.T) {
		result := F.Pipe1(
			Of[Config](-5),
			Chain(isPositive),
		)
		assert.Equal(t, O.None[int](), result(defaultConfig))
	})
}

// TestFold tests extracting the value from a ReaderOption with handlers
func TestFold(t *testing.T) {
	t.Run("Fold with Some", func(t *testing.T) {
		ro := Of[Config](42)
		result := Fold(
			reader.Of[Config]("none"),
			func(x int) reader.Reader[Config, string] {
				return reader.Of[Config](fmt.Sprintf("%d", x))
			},
		)(ro)
		assert.Equal(t, "42", result(defaultConfig))
	})

	t.Run("Fold with None", func(t *testing.T) {
		ro := None[Config, int]()
		result := Fold(
			reader.Of[Config]("none"),
			func(x int) reader.Reader[Config, string] {
				return reader.Of[Config](fmt.Sprintf("%d", x))
			},
		)(ro)
		assert.Equal(t, "none", result(defaultConfig))
	})
}

// TestMonadFold tests the non-curried version of Fold
func TestMonadFold(t *testing.T) {
	t.Run("MonadFold with Some", func(t *testing.T) {
		ro := Of[Config](42)
		result := MonadFold(
			ro,
			reader.Of[Config]("none"),
			func(x int) reader.Reader[Config, string] {
				return reader.Of[Config](fmt.Sprintf("%d", x))
			},
		)
		assert.Equal(t, "42", result(defaultConfig))
	})

	t.Run("MonadFold with None", func(t *testing.T) {
		ro := None[Config, int]()
		result := MonadFold(
			ro,
			reader.Of[Config]("none"),
			func(x int) reader.Reader[Config, string] {
				return reader.Of[Config](fmt.Sprintf("%d", x))
			},
		)
		assert.Equal(t, "none", result(defaultConfig))
	})
}

// TestGetOrElse tests getting the value or a default
func TestGetOrElse(t *testing.T) {
	t.Run("GetOrElse with Some", func(t *testing.T) {
		ro := Of[Config](42)
		result := GetOrElse(reader.Of[Config](0))(ro)
		assert.Equal(t, 42, result(defaultConfig))
	})

	t.Run("GetOrElse with None", func(t *testing.T) {
		ro := None[Config, int]()
		result := GetOrElse(reader.Of[Config](99))(ro)
		assert.Equal(t, 99, result(defaultConfig))
	})
}

// TestAsk tests retrieving the current environment
func TestAsk(t *testing.T) {
	ro := Ask[Config]()
	result := ro(defaultConfig)
	assert.Equal(t, O.Some(defaultConfig), result)
}

// TestAsks tests applying a function to the environment
func TestAsks(t *testing.T) {
	getPort := Asks(func(cfg Config) int {
		return cfg.Port
	})
	result := getPort(defaultConfig)
	assert.Equal(t, O.Some(8080), result)
}

// TestMonadChainOptionK tests chaining with a function that returns an Option
func TestMonadChainOptionK(t *testing.T) {
	parsePositive := func(x int) O.Option[int] {
		if x > 0 {
			return O.Some(x)
		}
		return O.None[int]()
	}

	t.Run("ChainOptionK with valid value", func(t *testing.T) {
		ro := Of[Config](42)
		result := MonadChainOptionK(ro, parsePositive)
		assert.Equal(t, O.Some(42), result(defaultConfig))
	})

	t.Run("ChainOptionK with invalid value", func(t *testing.T) {
		ro := Of[Config](-5)
		result := MonadChainOptionK(ro, parsePositive)
		assert.Equal(t, O.None[int](), result(defaultConfig))
	})

	t.Run("ChainOptionK with None", func(t *testing.T) {
		ro := None[Config, int]()
		result := MonadChainOptionK(ro, parsePositive)
		assert.Equal(t, O.None[int](), result(defaultConfig))
	})
}

// TestChainOptionK tests the curried version of MonadChainOptionK
func TestChainOptionK(t *testing.T) {
	parsePositive := func(x int) O.Option[int] {
		if x > 0 {
			return O.Some(x)
		}
		return O.None[int]()
	}

	t.Run("ChainOptionK with valid value", func(t *testing.T) {
		result := F.Pipe1(
			Of[Config](42),
			ChainOptionK[Config](parsePositive),
		)
		assert.Equal(t, O.Some(42), result(defaultConfig))
	})

	t.Run("ChainOptionK with invalid value", func(t *testing.T) {
		result := F.Pipe1(
			Of[Config](-5),
			ChainOptionK[Config](parsePositive),
		)
		assert.Equal(t, O.None[int](), result(defaultConfig))
	})
}

// TestFlatten tests removing one level of nesting
func TestFlatten(t *testing.T) {
	t.Run("Flatten nested Some", func(t *testing.T) {
		nested := Of[Config](Of[Config](42))
		flattened := Flatten(nested)
		result := flattened(defaultConfig)
		assert.Equal(t, O.Some(42), result)
	})

	t.Run("Flatten outer None", func(t *testing.T) {
		nested := None[Config, ReaderOption[Config, int]]()
		flattened := Flatten(nested)
		result := flattened(defaultConfig)
		assert.Equal(t, O.None[int](), result)
	})

	t.Run("Flatten inner None", func(t *testing.T) {
		nested := Of[Config](None[Config, int]())
		flattened := Flatten(nested)
		result := flattened(defaultConfig)
		assert.Equal(t, O.None[int](), result)
	})
}

// TestLocal tests transforming the environment before passing it to a computation
func TestLocal(t *testing.T) {
	type GlobalConfig struct {
		DB Config
	}

	getPort := Asks(func(cfg Config) int {
		return cfg.Port
	})

	globalConfig := GlobalConfig{
		DB: defaultConfig,
	}

	result := Local[int](func(g GlobalConfig) Config {
		return g.DB
	})(getPort)

	assert.Equal(t, O.Some(8080), result(globalConfig))
}

// TestRead tests executing a ReaderOption with an environment
func TestRead(t *testing.T) {
	ro := Of[Config](42)
	result := Read[int](defaultConfig)(ro)
	assert.Equal(t, O.Some(42), result)
}

// TestReadOption tests executing a ReaderOption with an optional environment
func TestReadOption(t *testing.T) {
	ro := Of[Config](42)

	t.Run("ReadOption with Some environment", func(t *testing.T) {
		result := ReadOption[int](O.Some(defaultConfig))(ro)
		assert.Equal(t, O.Some(42), result)
	})

	t.Run("ReadOption with None environment", func(t *testing.T) {
		result := ReadOption[int](O.None[Config]())(ro)
		assert.Equal(t, O.None[int](), result)
	})
}

// TestMonadFlap tests applying a value to a function wrapped in a ReaderOption
func TestMonadFlap(t *testing.T) {
	t.Run("Flap with Some function", func(t *testing.T) {
		fab := Of[Config](utils.Double)
		result := MonadFlap(fab, 21)
		assert.Equal(t, O.Some(42), result(defaultConfig))
	})

	t.Run("Flap with None function", func(t *testing.T) {
		fab := None[Config, func(int) int]()
		result := MonadFlap(fab, 21)
		assert.Equal(t, O.None[int](), result(defaultConfig))
	})
}

// TestFlap tests the curried version of MonadFlap
func TestFlap(t *testing.T) {
	t.Run("Flap with Some function", func(t *testing.T) {
		result := F.Pipe1(
			Of[Config](utils.Double),
			Flap[Config, int](21),
		)
		assert.Equal(t, O.Some(42), result(defaultConfig))
	})
}

// TestMonadAlt tests providing an alternative ReaderOption
func TestMonadAlt(t *testing.T) {
	t.Run("Alt with first Some", func(t *testing.T) {
		primary := Of[Config](42)
		fallback := Of[Config](99)
		result := MonadAlt(primary, lazy.Of(fallback))
		assert.Equal(t, O.Some(42), result(defaultConfig))
	})

	t.Run("Alt with first None", func(t *testing.T) {
		primary := None[Config, int]()
		fallback := Of[Config](99)
		result := MonadAlt(primary, lazy.Of(fallback))
		assert.Equal(t, O.Some(99), result(defaultConfig))
	})

	t.Run("Alt with both None", func(t *testing.T) {
		primary := None[Config, int]()
		fallback := None[Config, int]()
		result := MonadAlt(primary, lazy.Of(fallback))
		assert.Equal(t, O.None[int](), result(defaultConfig))
	})
}

// TestAlt tests the curried version of MonadAlt
func TestAlt(t *testing.T) {
	t.Run("Alt with first Some", func(t *testing.T) {
		result := F.Pipe1(
			Of[Config](42),
			Alt(lazy.Of(Of[Config](99))),
		)
		assert.Equal(t, O.Some(42), result(defaultConfig))
	})

	t.Run("Alt with first None", func(t *testing.T) {
		result := F.Pipe1(
			None[Config, int](),
			Alt(lazy.Of(Of[Config](99))),
		)
		assert.Equal(t, O.Some(99), result(defaultConfig))
	})
}

// TestComplexChaining tests a complex chain of operations
func TestComplexChaining(t *testing.T) {
	// Simulate a complex workflow with environment access
	result := F.Pipe3(
		Ask[Config](),
		Map[Config](func(cfg Config) int { return cfg.Port }),
		Chain(func(port int) ReaderOption[Config, int] {
			if port > 0 {
				return Of[Config](port * 2)
			}
			return None[Config, int]()
		}),
		Map[Config](func(x int) string { return fmt.Sprintf("%d", x) }),
	)

	assert.Equal(t, O.Some("16160"), result(defaultConfig))
}

// TestEnvironmentDependentComputation tests computations that depend on environment
func TestEnvironmentDependentComputation(t *testing.T) {
	// A computation that uses the environment to make decisions
	validateTimeout := func(value int) ReaderOption[Config, int] {
		return func(cfg Config) O.Option[int] {
			if value <= cfg.Timeout {
				return O.Some(value)
			}
			return O.None[int]()
		}
	}

	t.Run("Value within timeout", func(t *testing.T) {
		result := F.Pipe1(
			Of[Config](20),
			Chain(validateTimeout),
		)
		assert.Equal(t, O.Some(20), result(defaultConfig))
	})

	t.Run("Value exceeds timeout", func(t *testing.T) {
		result := F.Pipe1(
			Of[Config](50),
			Chain(validateTimeout),
		)
		assert.Equal(t, O.None[int](), result(defaultConfig))
	})
}
