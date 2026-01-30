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
	"github.com/stretchr/testify/assert"
)

func TestTraverseArray_AllSuccess(t *testing.T) {
	// Test traversing an array where all operations succeed
	double := func(x int) ReaderIOOption[context.Context, int] {
		return Of[context.Context](x * 2)
	}

	input := []int{1, 2, 3, 4, 5}
	result := TraverseArray[context.Context](double)(input)

	expected := O.Of([]int{2, 4, 6, 8, 10})
	assert.Equal(t, expected, result(context.Background())())
}

func TestTraverseArray_OneFailure(t *testing.T) {
	// Test traversing an array where one operation fails
	failOnThree := func(x int) ReaderIOOption[context.Context, int] {
		if x == 3 {
			return None[context.Context, int]()
		}
		return Of[context.Context](x * 2)
	}

	input := []int{1, 2, 3, 4, 5}
	result := TraverseArray[context.Context](failOnThree)(input)

	expected := O.None[[]int]()
	assert.Equal(t, expected, result(context.Background())())
}

func TestTraverseArray_EmptyArray(t *testing.T) {
	// Test traversing an empty array
	double := func(x int) ReaderIOOption[context.Context, int] {
		return Of[context.Context](x * 2)
	}

	input := []int{}
	result := TraverseArray[context.Context](double)(input)

	expected := O.Of([]int{})
	assert.Equal(t, expected, result(context.Background())())
}

func TestTraverseArray_WithEnvironment(t *testing.T) {
	// Test that the environment is properly passed through
	type Config struct {
		Multiplier int
	}

	multiply := func(x int) ReaderIOOption[Config, int] {
		return func(cfg Config) IOOption[int] {
			return func() Option[int] {
				return O.Of(x * cfg.Multiplier)
			}
		}
	}

	input := []int{1, 2, 3}
	result := TraverseArray[Config](multiply)(input)

	cfg := Config{Multiplier: 10}
	expected := O.Of([]int{10, 20, 30})
	assert.Equal(t, expected, result(cfg)())
}

func TestTraverseArray_ChainedOperation(t *testing.T) {
	// Test traversing as part of a chain
	type Config struct {
		Factor int
	}

	multiplyByFactor := func(x int) ReaderIOOption[Config, int] {
		return func(cfg Config) IOOption[int] {
			return func() Option[int] {
				return O.Of(x * cfg.Factor)
			}
		}
	}

	result := F.Pipe1(
		Of[Config]([]int{1, 2, 3, 4}),
		Chain(TraverseArray[Config](multiplyByFactor)),
	)

	cfg := Config{Factor: 5}
	expected := O.Of([]int{5, 10, 15, 20})
	assert.Equal(t, expected, result(cfg)())
}

func TestTraverseArrayWithIndex_AllSuccess(t *testing.T) {
	// Test traversing with index where all operations succeed
	addIndex := func(idx int, x string) ReaderIOOption[context.Context, string] {
		return Of[context.Context](fmt.Sprintf("%d:%s", idx, x))
	}

	input := []string{"a", "b", "c"}
	result := TraverseArrayWithIndex[context.Context](addIndex)(input)

	expected := O.Of([]string{"0:a", "1:b", "2:c"})
	assert.Equal(t, expected, result(context.Background())())
}

func TestTraverseArrayWithIndex_OneFailure(t *testing.T) {
	// Test traversing with index where one operation fails
	failOnIndex := func(idx int, x string) ReaderIOOption[context.Context, string] {
		if idx == 1 {
			return None[context.Context, string]()
		}
		return Of[context.Context](fmt.Sprintf("%d:%s", idx, x))
	}

	input := []string{"a", "b", "c"}
	result := TraverseArrayWithIndex[context.Context](failOnIndex)(input)

	expected := O.None[[]string]()
	assert.Equal(t, expected, result(context.Background())())
}

func TestTraverseArrayWithIndex_EmptyArray(t *testing.T) {
	// Test traversing an empty array with index
	addIndex := func(idx int, x string) ReaderIOOption[context.Context, string] {
		return Of[context.Context](fmt.Sprintf("%d:%s", idx, x))
	}

	input := []string{}
	result := TraverseArrayWithIndex[context.Context](addIndex)(input)

	expected := O.Of([]string{})
	assert.Equal(t, expected, result(context.Background())())
}

func TestTraverseArrayWithIndex_WithEnvironment(t *testing.T) {
	// Test that environment is properly passed with index
	type Config struct {
		Prefix string
	}

	formatWithIndex := func(idx int, x string) ReaderIOOption[Config, string] {
		return func(cfg Config) IOOption[string] {
			return func() Option[string] {
				return O.Of(fmt.Sprintf("%s%d:%s", cfg.Prefix, idx, x))
			}
		}
	}

	input := []string{"a", "b", "c"}
	result := TraverseArrayWithIndex[Config](formatWithIndex)(input)

	cfg := Config{Prefix: "item-"}
	expected := O.Of([]string{"item-0:a", "item-1:b", "item-2:c"})
	assert.Equal(t, expected, result(cfg)())
}

func TestTraverseArrayWithIndex_IndexUsedInLogic(t *testing.T) {
	// Test using index in computation logic
	multiplyByIndex := func(idx int, x int) ReaderIOOption[context.Context, int] {
		return Of[context.Context](x * idx)
	}

	input := []int{10, 20, 30, 40}
	result := TraverseArrayWithIndex[context.Context](multiplyByIndex)(input)

	// 10*0=0, 20*1=20, 30*2=60, 40*3=120
	expected := O.Of([]int{0, 20, 60, 120})
	assert.Equal(t, expected, result(context.Background())())
}

func TestTraverseArray_ComplexType(t *testing.T) {
	// Test traversing with complex types
	type User struct {
		ID   int
		Name string
	}

	type UserProfile struct {
		UserID      int
		DisplayName string
	}

	loadProfile := func(user User) ReaderIOOption[context.Context, UserProfile] {
		return Of[context.Context](UserProfile{
			UserID:      user.ID,
			DisplayName: "Profile: " + user.Name,
		})
	}

	users := []User{
		{ID: 1, Name: "Alice"},
		{ID: 2, Name: "Bob"},
		{ID: 3, Name: "Charlie"},
	}

	result := TraverseArray[context.Context](loadProfile)(users)

	expected := O.Of([]UserProfile{
		{UserID: 1, DisplayName: "Profile: Alice"},
		{UserID: 2, DisplayName: "Profile: Bob"},
		{UserID: 3, DisplayName: "Profile: Charlie"},
	})
	assert.Equal(t, expected, result(context.Background())())
}

func TestTraverseArray_ConditionalFailure(t *testing.T) {
	// Test conditional failure based on environment
	type Config struct {
		MaxValue int
	}

	validateAndDouble := func(x int) ReaderIOOption[Config, int] {
		return func(cfg Config) IOOption[int] {
			return func() Option[int] {
				if x > cfg.MaxValue {
					return O.None[int]()
				}
				return O.Of(x * 2)
			}
		}
	}

	input := []int{1, 2, 3, 4, 5}

	// With MaxValue=3, should fail on 4 and 5
	cfg1 := Config{MaxValue: 3}
	result1 := TraverseArray[Config](validateAndDouble)(input)
	assert.Equal(t, O.None[[]int](), result1(cfg1)())

	// With MaxValue=10, all should succeed
	cfg2 := Config{MaxValue: 10}
	result2 := TraverseArray[Config](validateAndDouble)(input)
	expected := O.Of([]int{2, 4, 6, 8, 10})
	assert.Equal(t, expected, result2(cfg2)())
}
