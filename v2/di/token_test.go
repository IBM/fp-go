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

package di

import (
	"errors"
	"testing"

	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	IOO "github.com/IBM/fp-go/v2/iooption"
	"github.com/IBM/fp-go/v2/ioresult"
	O "github.com/IBM/fp-go/v2/option"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

func TestMakeToken(t *testing.T) {
	token := MakeToken[string]("TestToken")

	assert.NotNil(t, token)
	assert.Equal(t, "TestToken", token.String())
	assert.NotEmpty(t, token.Id())
}

func TestTokenIdentity(t *testing.T) {
	token := MakeToken[int]("IntToken")
	identity := token.Identity()

	assert.NotNil(t, identity)
	assert.Equal(t, token.Id(), identity.Id())
	assert.Equal(t, token.String(), identity.String())
}

func TestTokenOption(t *testing.T) {
	token := MakeToken[int]("IntToken")
	option := token.Option()

	assert.NotNil(t, option)
	assert.Contains(t, option.String(), "Option")
	assert.Equal(t, token.Id(), option.Id())
}

func TestTokenIOEither(t *testing.T) {
	token := MakeToken[int]("IntToken")
	ioeither := token.IOEither()

	assert.NotNil(t, ioeither)
	assert.Contains(t, ioeither.String(), "IOEither")
	assert.Equal(t, token.Id(), ioeither.Id())
}

func TestTokenIOOption(t *testing.T) {
	token := MakeToken[int]("IntToken")
	iooption := token.IOOption()

	assert.NotNil(t, iooption)
	assert.Contains(t, iooption.String(), "IOOption")
	assert.Equal(t, token.Id(), iooption.Id())
}

func TestTokenUnerase(t *testing.T) {
	token := MakeToken[int]("IntToken")

	// Test successful unerase
	res := token.Unerase(42)
	assert.True(t, E.IsRight(res))
	assert.Equal(t, result.Of(42), res)

	// Test failed unerase (wrong type)
	result2 := token.Unerase("not an int")
	assert.True(t, E.IsLeft(result2))
}

func TestTokenFlag(t *testing.T) {
	token := MakeToken[int]("IntToken")

	// Flags should be set for different dependency types
	optionFlag := token.Option().Flag()
	ioeitherFlag := token.IOEither().Flag()
	iooptionFlag := token.IOOption().Flag()

	// Flags should be different for different types
	assert.NotEqual(t, optionFlag, ioeitherFlag)
	assert.NotEqual(t, optionFlag, iooptionFlag)
	assert.NotEqual(t, ioeitherFlag, iooptionFlag)
}

func TestTokenProviderFactory(t *testing.T) {
	// Token without default
	token1 := MakeToken[int]("Token1")
	assert.True(t, O.IsNone(token1.ProviderFactory()))

	// Token with default
	token2 := MakeTokenWithDefault0("Token2", ioresult.Of(42))
	assert.True(t, O.IsSome(token2.ProviderFactory()))
}

func TestMakeMultiToken(t *testing.T) {
	multiToken := MakeMultiToken[string]("MultiToken")

	assert.NotNil(t, multiToken)
	assert.NotNil(t, multiToken.Container())
	assert.NotNil(t, multiToken.Item())

	// Container and Item should share the same ID
	assert.Equal(t, multiToken.Container().Id(), multiToken.Item().Id())

	// But have different names
	assert.Contains(t, multiToken.Container().String(), "Container")
	assert.Contains(t, multiToken.Item().String(), "Item")
}

func TestMultiTokenFlags(t *testing.T) {
	multiToken := MakeMultiToken[string]("MultiToken")

	// Container should have Multi flag
	containerFlag := multiToken.Container().Flag()
	assert.NotZero(t, containerFlag)

	// Item should have Item flag
	itemFlag := multiToken.Item().Flag()
	assert.NotZero(t, itemFlag)
}

func TestTokenUniqueness(t *testing.T) {
	token1 := MakeToken[int]("Token")
	token2 := MakeToken[int]("Token")

	// Even with same name, IDs should be different
	assert.NotEqual(t, token1.Id(), token2.Id())
}

func TestOptionTokenUnerase(t *testing.T) {
	token := MakeToken[int]("IntToken")
	optionToken := token.Option()

	// Test successful unerase with Some
	res := optionToken.Unerase(O.Of[any](42))
	assert.True(t, E.IsRight(res))

	// Test successful unerase with None
	noneResult := optionToken.Unerase(O.None[any]())
	assert.True(t, E.IsRight(noneResult))
	assert.Equal(t, result.Of(O.None[int]()), noneResult)

	// Test failed unerase (wrong type)
	badResult := optionToken.Unerase(42) // Not an Option
	assert.True(t, E.IsLeft(badResult))
}

func TestIOEitherTokenUnerase(t *testing.T) {
	token := MakeToken[int]("IntToken")
	ioeitherToken := token.IOEither()

	// Test successful unerase
	ioValue := ioresult.Of(any(42))
	result := ioeitherToken.Unerase(ioValue)
	assert.True(t, E.IsRight(result))

	// Execute the IOEither to verify it works
	if E.IsRight(result) {
		ioe := E.ToOption(result)
		if O.IsSome(ioe) {
			executed := O.GetOrElse(F.Constant(IOE.Left[int](errors.New("fail"))))(ioe)()
			assert.True(t, E.IsRight(executed))
		}
	}

	// Test failed unerase (wrong type)
	badResult := ioeitherToken.Unerase(42)
	assert.True(t, E.IsLeft(badResult))
}

func TestIOOptionTokenUnerase(t *testing.T) {
	token := MakeToken[int]("IntToken")
	iooptionToken := token.IOOption()

	// Test successful unerase
	ioValue := IOO.Of(any(42))
	result := iooptionToken.Unerase(ioValue)
	assert.True(t, E.IsRight(result))

	// Test failed unerase (wrong type)
	badResult := iooptionToken.Unerase(42)
	assert.True(t, E.IsLeft(badResult))
}

func TestMultiTokenContainerUnerase(t *testing.T) {
	multiToken := MakeMultiToken[int]("MultiToken")
	container := multiToken.Container()

	// Test successful unerase with array
	arrayValue := []any{1, 2, 3}
	result := container.Unerase(arrayValue)
	assert.True(t, E.IsRight(result))

	if E.IsRight(result) {
		arr := E.ToOption(result)
		if O.IsSome(arr) {
			values := O.GetOrElse(F.Constant([]int{}))(arr)
			assert.Equal(t, []int{1, 2, 3}, values)
		}
	}

	// Test failed unerase (wrong type in array)
	badArray := []any{1, "not an int", 3}
	badResult := container.Unerase(badArray)
	assert.True(t, E.IsLeft(badResult))
}

func TestMakeTokenWithDefault(t *testing.T) {
	factory := MakeProviderFactory0(ioresult.Of(42))
	token := MakeTokenWithDefault[int]("TokenWithDefault", factory)

	assert.NotNil(t, token)
	assert.True(t, O.IsSome(token.ProviderFactory()))
}

func TestTokenStringRepresentation(t *testing.T) {
	token := MakeToken[int]("MyToken")

	assert.Equal(t, "MyToken", token.String())
	assert.Contains(t, token.Option().String(), "Option[MyToken]")
	assert.Contains(t, token.IOEither().String(), "IOEither[MyToken]")
	assert.Contains(t, token.IOOption().String(), "IOOption[MyToken]")
}

func TestMultiTokenStringRepresentation(t *testing.T) {
	multiToken := MakeMultiToken[int]("MyMulti")

	assert.Contains(t, multiToken.Container().String(), "Container[MyMulti]")
	assert.Contains(t, multiToken.Item().String(), "Item[MyMulti]")
}

// Benchmark tests
func BenchmarkMakeToken(b *testing.B) {
	for b.Loop() {
		MakeToken[int]("BenchToken")
	}
}

func BenchmarkTokenUnerase(b *testing.B) {
	token := MakeToken[int]("BenchToken")
	value := any(42)

	b.ResetTimer()
	for b.Loop() {
		token.Unerase(value)
	}
}

func BenchmarkMakeMultiToken(b *testing.B) {
	for b.Loop() {
		MakeMultiToken[int]("BenchMulti")
	}
}
