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
	"strconv"
	"sync/atomic"

	DIE "github.com/IBM/fp-go/v2/di/erasure"
	IO "github.com/IBM/fp-go/v2/io"
	O "github.com/IBM/fp-go/v2/option"
)

// Dependency describes the relationship to a service, that has a type and
// a behaviour such as required, option or lazy
type Dependency[T any] interface {
	DIE.Dependency
	// Unerase converts a value with erased type signature into a strongly typed value
	Unerase(val any) Result[T]
}

// InjectionToken uniquely identifies a dependency by giving it an Id, Type and name
type InjectionToken[T any] interface {
	Dependency[T]
	// Identity idenifies this dependency as a mandatory, required dependency, it will be resolved eagerly and injected as `T`.
	// If the dependency cannot be resolved, the resolution process fails
	Identity() Dependency[T]
	// Option identifies this dependency as optional, it will be resolved eagerly and injected as [Option[T]].
	// If the dependency cannot be resolved, the resolution process continues and the dependency is represented as [O.None[T]]
	Option() Dependency[Option[T]]
	// IOEither identifies this dependency as mandatory but it will be resolved lazily as a [IOResult[T]]. This
	// value is memoized to make sure the dependency is a singleton.
	// If the dependency cannot be resolved, the resolution process fails
	IOEither() Dependency[IOResult[T]]
	// IOOption identifies this dependency as optional but it will be resolved lazily as a [IOOption[T]]. This
	// value is memoized to make sure the dependency is a singleton.
	// If the dependency cannot be resolved, the resolution process continues and the dependency is represented as the none value.
	IOOption() Dependency[IOOption[T]]
}

// MultiInjectionToken uniquely identifies a dependency by giving it an Id, Type and name that can have multiple implementations.
// Implementations are provided via the [MultiInjectionToken.Item] injection token.
type MultiInjectionToken[T any] interface {
	// Container returns the injection token used to request an array of all provided items
	Container() InjectionToken[[]T]
	// Item returns the injection token used to provide an item
	Item() InjectionToken[T]
}

// makeID creates a generator of unique string IDs
func makeID() IO.IO[string] {
	var count atomic.Int64
	return func() string {
		return strconv.FormatInt(count.Add(1), 16)
	}
}

// genID is the common generator of unique string IDs
var genID = makeID()

type tokenBase struct {
	name            string
	id              string
	flag            int
	providerFactory Option[DIE.ProviderFactory]
}

type token[T any] struct {
	base   *tokenBase
	toType func(val any) Result[T]
}

func (t *token[T]) Id() string {
	return t.base.id
}

func (t *token[T]) Flag() int {
	return t.base.flag
}

func (t *token[T]) String() string {
	return t.base.name
}

func (t *token[T]) Unerase(val any) Result[T] {
	return t.toType(val)
}

func (t *token[T]) ProviderFactory() Option[DIE.ProviderFactory] {
	return t.base.providerFactory
}
func makeTokenBase(name, id string, typ int, providerFactory Option[DIE.ProviderFactory]) *tokenBase {
	return &tokenBase{name, id, typ, providerFactory}
}

func makeToken[T any](name, id string, typ int, unerase func(val any) Result[T], providerFactory Option[DIE.ProviderFactory]) Dependency[T] {
	return &token[T]{makeTokenBase(name, id, typ, providerFactory), unerase}
}

type injectionToken[T any] struct {
	token[T]
	option   Dependency[Option[T]]
	ioeither Dependency[IOResult[T]]
	iooption Dependency[IOOption[T]]
}

type multiInjectionToken[T any] struct {
	container *injectionToken[[]T]
	item      *injectionToken[T]
}

func (i *injectionToken[T]) Identity() Dependency[T] {
	return i
}

func (i *injectionToken[T]) Option() Dependency[Option[T]] {
	return i.option
}

func (i *injectionToken[T]) IOEither() Dependency[IOResult[T]] {
	return i.ioeither
}

func (i *injectionToken[T]) IOOption() Dependency[IOOption[T]] {
	return i.iooption
}

func (i *injectionToken[T]) ProviderFactory() Option[DIE.ProviderFactory] {
	return i.base.providerFactory
}

func (m *multiInjectionToken[T]) Container() InjectionToken[[]T] {
	return m.container
}

func (m *multiInjectionToken[T]) Item() InjectionToken[T] {
	return m.item
}

// makeToken create a unique [InjectionToken] for a specific type
func makeInjectionToken[T any](name string, providerFactory Option[DIE.ProviderFactory]) InjectionToken[T] {
	id := genID()
	toIdentity := toType[T]()
	return &injectionToken[T]{
		token[T]{makeTokenBase(name, id, DIE.IDENTITY, providerFactory), toIdentity},
		makeToken(fmt.Sprintf("Option[%s]", name), id, DIE.OPTION, toOptionType(toIdentity), providerFactory),
		makeToken(fmt.Sprintf("IOEither[%s]", name), id, DIE.IOEITHER, toIOEitherType(toIdentity), providerFactory),
		makeToken(fmt.Sprintf("IOOption[%s]", name), id, DIE.IOOPTION, toIOOptionType(toIdentity), providerFactory),
	}
}

// MakeToken create a unique [InjectionToken] for a specific type
func MakeToken[T any](name string) InjectionToken[T] {
	return makeInjectionToken[T](name, O.None[DIE.ProviderFactory]())
}

// MakeToken create a unique [InjectionToken] for a specific type
func MakeTokenWithDefault[T any](name string, providerFactory DIE.ProviderFactory) InjectionToken[T] {
	return makeInjectionToken[T](name, O.Of(providerFactory))
}

// MakeMultiToken creates a [MultiInjectionToken]
func MakeMultiToken[T any](name string) MultiInjectionToken[T] {
	id := genID()
	toItem := toType[T]()
	toContainer := toArrayType(toItem)
	containerName := fmt.Sprintf("Container[%s]", name)
	itemName := fmt.Sprintf("Item[%s]", name)
	// empty factory
	providerFactory := O.None[DIE.ProviderFactory]()
	// container
	container := &injectionToken[[]T]{
		token[[]T]{makeTokenBase(containerName, id, DIE.MULTI|DIE.IDENTITY, providerFactory), toContainer},
		makeToken(fmt.Sprintf("Option[%s]", containerName), id, DIE.MULTI|DIE.OPTION, toOptionType(toContainer), providerFactory),
		makeToken(fmt.Sprintf("IOEither[%s]", containerName), id, DIE.OPTION|DIE.IOEITHER, toIOEitherType(toContainer), providerFactory),
		makeToken(fmt.Sprintf("IOOption[%s]", containerName), id, DIE.OPTION|DIE.IOOPTION, toIOOptionType(toContainer), providerFactory),
	}
	// item
	item := &injectionToken[T]{
		token[T]{makeTokenBase(itemName, id, DIE.ITEM|DIE.IDENTITY, providerFactory), toItem},
		makeToken(fmt.Sprintf("Option[%s]", itemName), id, DIE.ITEM|DIE.OPTION, toOptionType(toItem), providerFactory),
		makeToken(fmt.Sprintf("IOEither[%s]", itemName), id, DIE.ITEM|DIE.IOEITHER, toIOEitherType(toItem), providerFactory),
		makeToken(fmt.Sprintf("IOOption[%s]", itemName), id, DIE.ITEM|DIE.IOOPTION, toIOOptionType(toItem), providerFactory),
	}
	// returns the token
	return &multiInjectionToken[T]{container, item}
}
