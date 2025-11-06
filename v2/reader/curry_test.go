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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCurry0(t *testing.T) {
	config := Config{Port: 8080}
	getValue := func(c Config) int { return c.Port }
	r := Curry0(getValue)
	result := r(config)
	assert.Equal(t, 8080, result)
}

func TestCurry1(t *testing.T) {
	config := Config{Prefix: ">> "}
	addPrefix := func(c Config, s string) string { return c.Prefix + s }
	curried := Curry1(addPrefix)
	r := curried("hello")
	result := r(config)
	assert.Equal(t, ">> hello", result)
}

func TestCurry2(t *testing.T) {
	config := Config{Prefix: "-"}
	join := func(c Config, a, b string) string { return a + c.Prefix + b }
	curried := Curry2(join)
	r := curried("hello")("world")
	result := r(config)
	assert.Equal(t, "hello-world", result)
}

func TestCurry3(t *testing.T) {
	config := Config{Prefix: "-"}
	join := func(c Config, a, b, d string) string {
		return a + c.Prefix + b + c.Prefix + d
	}
	curried := Curry3(join)
	r := curried("a")("b")("c")
	result := r(config)
	assert.Equal(t, "a-b-c", result)
}

func TestCurry4(t *testing.T) {
	config := Config{Multiplier: 10}
	sum := func(c Config, a, b, d, e int) int {
		return (a + b + d + e) * c.Multiplier
	}
	curried := Curry4(sum)
	r := curried(1)(2)(3)(4)
	result := r(config)
	assert.Equal(t, 100, result)
}

func TestUncurry0(t *testing.T) {
	config := Config{Port: 8080}
	r := Of[Config](42)
	f := Uncurry0(r)
	result := f(config)
	assert.Equal(t, 42, result)
}

func TestUncurry1(t *testing.T) {
	config := Config{Prefix: ">> "}
	curried := func(s string) Reader[Config, string] {
		return Asks(func(c Config) string { return c.Prefix + s })
	}
	f := Uncurry1(curried)
	result := f(config, "hello")
	assert.Equal(t, ">> hello", result)
}

func TestUncurry2(t *testing.T) {
	config := Config{Prefix: "-"}
	curried := func(a string) func(string) Reader[Config, string] {
		return func(b string) Reader[Config, string] {
			return Asks(func(c Config) string { return a + c.Prefix + b })
		}
	}
	f := Uncurry2(curried)
	result := f(config, "hello", "world")
	assert.Equal(t, "hello-world", result)
}

func TestUncurry3(t *testing.T) {
	config := Config{Prefix: "-"}
	curried := func(a string) func(string) func(string) Reader[Config, string] {
		return func(b string) func(string) Reader[Config, string] {
			return func(c string) Reader[Config, string] {
				return Asks(func(cfg Config) string {
					return fmt.Sprintf("%s%s%s%s%s", a, cfg.Prefix, b, cfg.Prefix, c)
				})
			}
		}
	}
	f := Uncurry3(curried)
	result := f(config, "a", "b", "c")
	assert.Equal(t, "a-b-c", result)
}

func TestUncurry4(t *testing.T) {
	config := Config{Multiplier: 10}
	curried := func(a int) func(int) func(int) func(int) Reader[Config, int] {
		return func(b int) func(int) func(int) Reader[Config, int] {
			return func(c int) func(int) Reader[Config, int] {
				return func(d int) Reader[Config, int] {
					return Asks(func(cfg Config) int { return (a + b + c + d) * cfg.Multiplier })
				}
			}
		}
	}
	f := Uncurry4(curried)
	result := f(config, 1, 2, 3, 4)
	assert.Equal(t, 100, result)
}
