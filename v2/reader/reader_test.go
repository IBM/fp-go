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

package reader

import (
	"fmt"
	"strconv"
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	N "github.com/IBM/fp-go/v2/number"
	T "github.com/IBM/fp-go/v2/tuple"
	"github.com/stretchr/testify/assert"

	"github.com/IBM/fp-go/v2/internal/utils"
)

type Config struct {
	Host       string
	Port       int
	Multiplier int
	Prefix     string
}

func TestAsk(t *testing.T) {
	config := Config{Host: "localhost", Port: 8080}
	r := Ask[Config]()
	result := r(config)
	assert.Equal(t, config, result)
}

func TestAsks(t *testing.T) {
	config := Config{Port: 8080}
	getPort := Asks(func(c Config) int { return c.Port })
	result := getPort(config)
	assert.Equal(t, 8080, result)
}

func TestAsksReader(t *testing.T) {
	config := Config{Host: "localhost"}
	r := AsksReader(func(c Config) Reader[Config, string] {
		if c.Host == "localhost" {
			return Of[Config]("local")
		}
		return Of[Config]("remote")
	})
	result := r(config)
	assert.Equal(t, "local", result)
}

func TestMap(t *testing.T) {
	assert.Equal(t, 2, F.Pipe1(Of[string](1), Map[string](utils.Double))(""))
}

func TestMonadMap(t *testing.T) {
	config := Config{Port: 8080}
	getPort := func(c Config) int { return c.Port }
	getPortStr := MonadMap(getPort, strconv.Itoa)
	result := getPortStr(config)
	assert.Equal(t, "8080", result)
}

func TestAp(t *testing.T) {
	assert.Equal(t, 2, F.Pipe1(Of[int](utils.Double), Ap[int](Of[int](1)))(0))
}

func TestMonadAp(t *testing.T) {
	config := Config{Port: 8080, Multiplier: 2}
	add := func(x int) func(int) int { return func(y int) int { return x + y } }
	getAdder := func(c Config) func(int) int { return add(c.Port) }
	getMultiplier := func(c Config) int { return c.Multiplier }
	result := MonadAp(getAdder, getMultiplier)(config)
	assert.Equal(t, 8082, result)
}

func TestOf(t *testing.T) {
	r := Of[Config]("constant")
	result := r(Config{Host: "any"})
	assert.Equal(t, "constant", result)
}

func TestChain(t *testing.T) {
	config := Config{Port: 8080}
	getPort := Asks(func(c Config) int { return c.Port })
	portToString := func(port int) Reader[Config, string] {
		return Of[Config](fmt.Sprintf("Port: %d", port))
	}
	r := Chain(portToString)(getPort)
	result := r(config)
	assert.Equal(t, "Port: 8080", result)
}

func TestMonadChain(t *testing.T) {
	config := Config{Port: 8080}
	getPort := func(c Config) int { return c.Port }
	portToString := func(port int) Reader[Config, string] {
		return func(c Config) string { return fmt.Sprintf("Port: %d", port) }
	}
	r := MonadChain(getPort, portToString)
	result := r(config)
	assert.Equal(t, "Port: 8080", result)
}

func TestFlatten(t *testing.T) {
	config := Config{Multiplier: 5}
	nested := func(c Config) Reader[Config, int] {
		return func(c2 Config) int { return c.Multiplier + c2.Multiplier }
	}
	flat := Flatten(nested)
	result := flat(config)
	assert.Equal(t, 10, result)
}

func TestCompose(t *testing.T) {
	type Env struct{ Config Config }
	env := Env{Config: Config{Port: 8080}}
	getConfig := func(e Env) Config { return e.Config }
	getPort := func(c Config) int { return c.Port }
	getPortFromEnv := Compose[int](getConfig)(getPort)
	result := getPortFromEnv(env)
	assert.Equal(t, 8080, result)
}

func TestFirst(t *testing.T) {
	double := N.Mul(2)
	r := First[int, int, string](double)
	result := r(T.MakeTuple2(5, "hello"))
	assert.Equal(t, T.MakeTuple2(10, "hello"), result)
}

func TestSecond(t *testing.T) {
	double := N.Mul(2)
	r := Second[string](double)
	result := r(T.MakeTuple2("hello", 5))
	assert.Equal(t, T.MakeTuple2("hello", 10), result)
}

func TestPromap(t *testing.T) {
	type Env struct{ Config Config }
	env := Env{Config: Config{Port: 8080}}
	getPort := func(c Config) int { return c.Port }
	extractConfig := func(e Env) Config { return e.Config }
	toString := strconv.Itoa
	r := Promap(extractConfig, toString)(getPort)
	result := r(env)
	assert.Equal(t, "8080", result)
}

func TestLocal(t *testing.T) {
	type DetailedConfig struct {
		Host string
		Port int
	}
	type SimpleConfig struct{ Host string }
	detailed := DetailedConfig{Host: "localhost", Port: 8080}
	getHost := func(c SimpleConfig) string { return c.Host }
	simplify := func(d DetailedConfig) SimpleConfig { return SimpleConfig{Host: d.Host} }
	r := Local[string](simplify)(getHost)
	result := r(detailed)
	assert.Equal(t, "localhost", result)
}

func TestRead(t *testing.T) {
	config := Config{Port: 8080}
	getPort := Asks(func(c Config) int { return c.Port })
	run := Read[int](config)
	port := run(getPort)
	assert.Equal(t, 8080, port)
}

func TestMonadFlap(t *testing.T) {
	config := Config{Multiplier: 3}
	getMultiplier := func(c Config) func(int) int {
		return func(x int) int { return x * c.Multiplier }
	}
	r := MonadFlap(getMultiplier, 5)
	result := r(config)
	assert.Equal(t, 15, result)
}

func TestFlap(t *testing.T) {
	config := Config{Multiplier: 3}
	getMultiplier := Asks(func(c Config) func(int) int {
		return func(x int) int { return x * c.Multiplier }
	})
	applyTo5 := Flap[Config, int, int](5)
	r := applyTo5(getMultiplier)
	result := r(config)
	assert.Equal(t, 15, result)
}
