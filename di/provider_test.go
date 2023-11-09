// Copyright (c) 2023 IBM Corp.
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

	A "github.com/IBM/fp-go/array"
	DIE "github.com/IBM/fp-go/di/erasure"
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	IOE "github.com/IBM/fp-go/ioeither"
	O "github.com/IBM/fp-go/option"
	"github.com/stretchr/testify/assert"
)

var (
	INJ_KEY2 = MakeToken[string]("INJ_KEY2")
	INJ_KEY1 = MakeToken[string]("INJ_KEY1")
	INJ_KEY3 = MakeToken[string]("INJ_KEY3")
)

func TestSimpleProvider(t *testing.T) {

	var staticCount int

	staticValue := func(value string) func() IOE.IOEither[error, string] {
		return func() IOE.IOEither[error, string] {
			return func() E.Either[error, string] {
				staticCount++
				return E.Of[error](fmt.Sprintf("Static based on [%s], at [%s]", value, time.Now()))
			}
		}
	}

	var dynamicCount int

	dynamicValue := func(value string) IOE.IOEither[error, string] {
		return func() E.Either[error, string] {
			dynamicCount++
			return E.Of[error](fmt.Sprintf("Dynamic based on [%s] at [%s]", value, time.Now()))
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

	staticValue := func(value string) func() IOE.IOEither[error, string] {
		return func() IOE.IOEither[error, string] {
			return func() E.Either[error, string] {
				staticCount++
				return E.Of[error](fmt.Sprintf("Static based on [%s], at [%s]", value, time.Now()))
			}
		}
	}

	var dynamicCount int

	dynamicValue := func(value O.Option[string]) IOE.IOEither[error, string] {
		return func() E.Either[error, string] {
			dynamicCount++
			return E.Of[error](fmt.Sprintf("Dynamic based on [%s] at [%s]", value, time.Now()))
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

	dynamicValue := func(value O.Option[string]) IOE.IOEither[error, string] {
		return func() E.Either[error, string] {
			dynamicCount++
			return E.Of[error](fmt.Sprintf("Dynamic based on [%s] at [%s]", value, time.Now()))
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

	dynamicValue := func(value string) IOE.IOEither[error, string] {
		return func() E.Either[error, string] {
			dynamicCount++
			return E.Of[error](fmt.Sprintf("Dynamic based on [%s] at [%s]", value, time.Now()))
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

	staticValue := func(value string) func() IOE.IOEither[error, string] {
		return func() IOE.IOEither[error, string] {
			return func() E.Either[error, string] {
				staticCount++
				return E.Of[error](fmt.Sprintf("Static based on [%s], at [%s]", value, time.Now()))
			}
		}
	}

	var dynamicCount int

	dynamicValue := func(value string) IOE.IOEither[error, string] {
		return func() E.Either[error, string] {
			dynamicCount++
			return E.Of[error](fmt.Sprintf("Dynamic based on [%s] at [%s]", value, time.Now()))
		}
	}

	var lazyEagerCount int

	lazyEager := func(laz IOE.IOEither[error, string], eager string) IOE.IOEither[error, string] {
		return F.Pipe1(
			laz,
			IOE.Chain(func(lazValue string) IOE.IOEither[error, string] {
				return func() E.Either[error, string] {
					lazyEagerCount++
					return E.Of[error](fmt.Sprintf("Dynamic based on [%s], [%s] at [%s]", lazValue, eager, time.Now()))
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
