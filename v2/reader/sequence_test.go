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
	"testing"

	T "github.com/IBM/fp-go/v2/tuple"
	"github.com/stretchr/testify/assert"
)

func TestSequenceT1(t *testing.T) {
	config := Config{Port: 8080}
	r := Asks(func(c Config) int { return c.Port })
	result := SequenceT1(r)
	tuple := result(config)
	assert.Equal(t, T.MakeTuple1(8080), tuple)
}

func TestSequenceT2(t *testing.T) {
	config := Config{Host: "localhost", Port: 8080}
	getHost := Asks(func(c Config) string { return c.Host })
	getPort := Asks(func(c Config) int { return c.Port })
	result := SequenceT2(getHost, getPort)
	tuple := result(config)
	assert.Equal(t, T.MakeTuple2("localhost", 8080), tuple)
}

func TestSequenceT3(t *testing.T) {
	config := Config{Host: "localhost", Port: 8080, Multiplier: 2}
	getHost := Asks(func(c Config) string { return c.Host })
	getPort := Asks(func(c Config) int { return c.Port })
	getMultiplier := Asks(func(c Config) int { return c.Multiplier })
	result := SequenceT3(getHost, getPort, getMultiplier)
	tuple := result(config)
	assert.Equal(t, T.MakeTuple3("localhost", 8080, 2), tuple)
}

func TestSequenceT4(t *testing.T) {
	config := Config{
		Host:       "localhost",
		Port:       8080,
		Multiplier: 2,
		Prefix:     ">>",
	}
	getHost := Asks(func(c Config) string { return c.Host })
	getPort := Asks(func(c Config) int { return c.Port })
	getMultiplier := Asks(func(c Config) int { return c.Multiplier })
	getPrefix := Asks(func(c Config) string { return c.Prefix })
	result := SequenceT4(getHost, getPort, getMultiplier, getPrefix)
	tuple := result(config)
	expected := T.MakeTuple4("localhost", 8080, 2, ">>")
	assert.Equal(t, expected, tuple)
}
