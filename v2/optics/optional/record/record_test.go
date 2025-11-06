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

package record

import (
	"testing"

	F "github.com/IBM/fp-go/v2/function"
	OP "github.com/IBM/fp-go/v2/optics/optional"
	O "github.com/IBM/fp-go/v2/option"
	ON "github.com/IBM/fp-go/v2/option/number"
	RR "github.com/IBM/fp-go/v2/record"

	"github.com/stretchr/testify/assert"
)

type (
	GenericMap = map[string]any
)

func TestOptionalRecord(t *testing.T) {
	// sample record
	r := RR.Singleton("key", "value")

	// extract values
	optKey := AtKey[string, string]("key")
	optKey1 := AtKey[string, string]("key1")

	// check if we can get the key
	assert.Equal(t, O.Of("value"), optKey.GetOption(r))
	assert.Equal(t, O.None[string](), optKey1.GetOption(r))

	// check if we can set a value
	r1 := optKey1.Set("value1")(r)

	// check if we can get the key
	assert.Equal(t, O.Of("value"), optKey.GetOption(r))
	assert.Equal(t, O.None[string](), optKey1.GetOption(r))
	// check if we can get the key
	assert.Equal(t, O.Of("value"), optKey.GetOption(r1))
	assert.Equal(t, O.Of("value1"), optKey1.GetOption(r1))
}

func TestOptionalWithType(t *testing.T) {
	// sample record
	r := RR.Singleton("key", "1")
	// convert between string and int
	// writes a key
	optStringKey := AtKey[string, string]("key")
	optIntKey := F.Pipe1(
		optStringKey,
		OP.IChain[map[string]string](ON.Atoi, ON.Itoa),
	)
	// test the scenarions
	assert.Equal(t, O.Of("1"), optStringKey.GetOption(r))
	assert.Equal(t, O.Of(1), optIntKey.GetOption(r))
	// modify
	r1 := optIntKey.Set(2)(r)
	assert.Equal(t, O.Of("2"), optStringKey.GetOption(r1))
	assert.Equal(t, O.Of(2), optIntKey.GetOption(r1))
}

// func TestNestedRecord(t *testing.T) {
// 	// some sample data
// 	x := GenericMap{
// 		"a": GenericMap{
// 			"b": "1",
// 		},
// 	}
// 	// accessor for first level
// 	optA := F.Pipe1(
// 		AtKey[string, any]("a"),
// 		OP.IChainAny[GenericMap, GenericMap](),
// 	)
// 	optB := F.Pipe2(
// 		AtKey[string, any]("b"),
// 		OP.IChainAny[GenericMap, string](),
// 		OP.IChain[GenericMap](ON.Atoi, ON.Itoa),
// 	)
// 	// go directly to b
// 	optAB := F.Pipe1(
// 		optA,
// 		OP.Compose[GenericMap](optB),
// 	)
// 	// access the value of b
// 	assert.Equal(t, O.Of(1), optAB.GetOption(x))
// }
