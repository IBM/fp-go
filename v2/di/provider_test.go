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
	"fmt"
	"testing"
	"time"

	A "github.com/IBM/fp-go/v2/array"
	DIE "github.com/IBM/fp-go/v2/di/erasure"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/result"
	"github.com/stretchr/testify/assert"
)

var (
	INJ_KEY2 = MakeToken[string]("INJ_KEY2")
	INJ_KEY1 = MakeToken[string]("INJ_KEY1")
	INJ_KEY3 = MakeToken[string]("INJ_KEY3")
)

func TestSimpleProvider(t *testing.T) {

	var staticCount int

	staticValue := func(value string) IOResult[string] {
		return func() Result[string] {
			staticCount++
			return result.Of(fmt.Sprintf("Static based on [%s], at [%s]", value, time.Now()))
		}
	}

	var dynamicCount int

	dynamicValue := func(value string) IOResult[string] {
		return func() Result[string] {
			dynamicCount++
			return result.Of(fmt.Sprintf("Dynamic based on [%s] at [%s]", value, time.Now()))
		}
	}

	p1 := MakeProvider0(INJ_KEY1, staticValue("Carsten"))
	p2 := MakeProvider1(INJ_KEY2, INJ_KEY1.Identity(), dynamicValue)

	inj := DIE.MakeInjector(A.From(p1, p2))

	i1 := Resolve(INJ_KEY1)
	i2 := Resolve(INJ_KEY2)

	res := IOE.SequenceT4(
		i2(inj),
		i1(inj),
		i2(inj),
		i1(inj),
	)

	r := res()

	assert.True(t, E.IsRight(r))
	assert.Equal(t, 1, staticCount)
	assert.Equal(t, 1, dynamicCount)
}

func TestOptionalProvider(t *testing.T) {

	var staticCount int

	staticValue := func(value string) IOResult[string] {
		return func() Result[string] {
			staticCount++
			return result.Of(fmt.Sprintf("Static based on [%s], at [%s]", value, time.Now()))
		}
	}

	var dynamicCount int

	dynamicValue := func(value Option[string]) IOResult[string] {
		return func() Result[string] {
			dynamicCount++
			return result.Of(fmt.Sprintf("Dynamic based on [%s] at [%s]", value, time.Now()))
		}
	}

	p1 := MakeProvider0(INJ_KEY1, staticValue("Carsten"))
	p2 := MakeProvider1(INJ_KEY2, INJ_KEY1.Option(), dynamicValue)

	inj := DIE.MakeInjector(A.From(p1, p2))

	i1 := Resolve(INJ_KEY1)
	i2 := Resolve(INJ_KEY2)

	res := IOE.SequenceT4(
		i2(inj),
		i1(inj),
		i2(inj),
		i1(inj),
	)

	r := res()

	assert.True(t, E.IsRight(r))
	assert.Equal(t, 1, staticCount)
	assert.Equal(t, 1, dynamicCount)
}

func TestOptionalProviderMissingDependency(t *testing.T) {

	var dynamicCount int

	dynamicValue := func(value Option[string]) IOResult[string] {
		return func() Result[string] {
			dynamicCount++
			return result.Of(fmt.Sprintf("Dynamic based on [%s] at [%s]", value, time.Now()))
		}
	}

	p2 := MakeProvider1(INJ_KEY2, INJ_KEY1.Option(), dynamicValue)

	inj := DIE.MakeInjector(A.From(p2))

	i2 := Resolve(INJ_KEY2)

	res := IOE.SequenceT2(
		i2(inj),
		i2(inj),
	)

	r := res()

	assert.True(t, E.IsRight(r))
	assert.Equal(t, 1, dynamicCount)
}

func TestProviderMissingDependency(t *testing.T) {

	var dynamicCount int

	dynamicValue := func(value string) IOResult[string] {
		return func() Result[string] {
			dynamicCount++
			return result.Of(fmt.Sprintf("Dynamic based on [%s] at [%s]", value, time.Now()))
		}
	}

	p2 := MakeProvider1(INJ_KEY2, INJ_KEY1.Identity(), dynamicValue)

	inj := DIE.MakeInjector(A.From(p2))

	i2 := Resolve(INJ_KEY2)

	res := IOE.SequenceT2(
		i2(inj),
		i2(inj),
	)

	r := res()

	assert.True(t, E.IsLeft(r))
	assert.Equal(t, 0, dynamicCount)
}

func TestEagerAndLazyProvider(t *testing.T) {

	var staticCount int

	staticValue := func(value string) IOResult[string] {
		return func() Result[string] {
			staticCount++
			return result.Of(fmt.Sprintf("Static based on [%s], at [%s]", value, time.Now()))
		}
	}

	var dynamicCount int

	dynamicValue := func(value string) IOResult[string] {
		return func() Result[string] {
			dynamicCount++
			return result.Of(fmt.Sprintf("Dynamic based on [%s] at [%s]", value, time.Now()))
		}
	}

	var lazyEagerCount int

	lazyEager := func(laz IOResult[string], eager string) IOResult[string] {
		return F.Pipe1(
			laz,
			IOE.Chain(func(lazValue string) IOResult[string] {
				return func() Result[string] {
					lazyEagerCount++
					return result.Of(fmt.Sprintf("Dynamic based on [%s], [%s] at [%s]", lazValue, eager, time.Now()))
				}
			}),
		)
	}

	p1 := MakeProvider0(INJ_KEY1, staticValue("Carsten"))
	p2 := MakeProvider1(INJ_KEY2, INJ_KEY1.Identity(), dynamicValue)
	p3 := MakeProvider2(INJ_KEY3, INJ_KEY2.IOEither(), INJ_KEY1.Identity(), lazyEager)

	inj := DIE.MakeInjector(A.From(p1, p2, p3))

	i3 := Resolve(INJ_KEY3)

	r := i3(inj)()

	fmt.Println(r)

	assert.True(t, E.IsRight(r))
	assert.Equal(t, 1, staticCount)
	assert.Equal(t, 1, dynamicCount)
	assert.Equal(t, 1, lazyEagerCount)
}

func TestItemProvider(t *testing.T) {
	// define a multi token
	injMulti := MakeMultiToken[string]("configs")

	// provide some values
	v1 := ConstProvider(injMulti.Item(), "Value1")
	v2 := ConstProvider(injMulti.Item(), "Value2")
	// mix in non-multi values
	p1 := ConstProvider(INJ_KEY1, "Value3")
	p2 := ConstProvider(INJ_KEY2, "Value4")

	// populate the injector
	inj := DIE.MakeInjector(A.From(p1, v1, p2, v2))

	// access the multi value
	multi := Resolve(injMulti.Container())

	multiInj := multi(inj)

	value := multiInj()

	assert.Equal(t, result.Of(A.From("Value1", "Value2")), value)
}

func TestEmptyItemProvider(t *testing.T) {
	// define a multi token
	injMulti := MakeMultiToken[string]("configs")

	// mix in non-multi values
	p1 := ConstProvider(INJ_KEY1, "Value3")
	p2 := ConstProvider(INJ_KEY2, "Value4")

	// populate the injector
	inj := DIE.MakeInjector(A.From(p1, p2))

	// access the multi value
	multi := Resolve(injMulti.Container())

	multiInj := multi(inj)

	value := multiInj()

	assert.Equal(t, result.Of(A.Empty[string]()), value)
}

func TestDependencyOnMultiProvider(t *testing.T) {
	// define a multi token
	injMulti := MakeMultiToken[string]("configs")

	// provide some values
	v1 := ConstProvider(injMulti.Item(), "Value1")
	v2 := ConstProvider(injMulti.Item(), "Value2")
	// mix in non-multi values
	p1 := ConstProvider(INJ_KEY1, "Value3")
	p2 := ConstProvider(INJ_KEY2, "Value4")

	fromMulti := func(val string, multi []string) IOResult[string] {
		return ioresult.Of(fmt.Sprintf("Val: %s, Multi: %s", val, multi))
	}
	p3 := MakeProvider2(INJ_KEY3, INJ_KEY1.Identity(), injMulti.Container().Identity(), fromMulti)

	// populate the injector
	inj := DIE.MakeInjector(A.From(p1, p2, v1, v2, p3))

	r3 := Resolve(INJ_KEY3)

	v := r3(inj)()

	assert.Equal(t, result.Of("Val: Value3, Multi: [Value1 Value2]"), v)
}

func TestTokenWithDefaultProvider(t *testing.T) {
	// token without a default
	injToken1 := MakeToken[string]("Token1")
	// token with a default
	injToken2 := MakeTokenWithDefault0("Token2", ioresult.Of("Carsten"))
	// dependency
	injToken3 := MakeToken[string]("Token3")

	p3 := MakeProvider1(injToken3, injToken2.Identity(), func(data string) IOResult[string] {
		return ioresult.Of(fmt.Sprintf("Token: %s", data))
	})

	// populate the injector
	inj := DIE.MakeInjector(A.From(p3))

	// resolving injToken3 should work and use the default provider for injToken2
	r1 := Resolve(injToken1)
	r3 := Resolve(injToken3)

	// inj1 should not be available
	assert.True(t, E.IsLeft(r1(inj)()))
	// r3 should work
	assert.Equal(t, result.Of("Token: Carsten"), r3(inj)())
}

func TestTokenWithDefaultProviderAndOverride(t *testing.T) {
	// token with a default
	injToken2 := MakeTokenWithDefault0("Token2", ioresult.Of("Carsten"))
	// dependency
	injToken3 := MakeToken[string]("Token3")

	p2 := ConstProvider(injToken2, "Override")

	p3 := MakeProvider1(injToken3, injToken2.Identity(), func(data string) IOResult[string] {
		return ioresult.Of(fmt.Sprintf("Token: %s", data))
	})

	// populate the injector
	inj := DIE.MakeInjector(A.From(p2, p3))

	// resolving injToken3 should work and use the default provider for injToken2
	r3 := Resolve(injToken3)

	// r3 should work
	assert.Equal(t, result.Of("Token: Override"), r3(inj)())
}
