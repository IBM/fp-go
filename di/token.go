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
	"strconv"
	"sync/atomic"

	DIE "github.com/IBM/fp-go/di/erasure"
	E "github.com/IBM/fp-go/either"
	IO "github.com/IBM/fp-go/io"
	IOE "github.com/IBM/fp-go/ioeither"
	IOO "github.com/IBM/fp-go/iooption"
	O "github.com/IBM/fp-go/option"
)

// Dependency describes the relationship to a service, that has a type and
// a behaviour such as required, option or lazy
type Dependency[T any] interface {
	DIE.Dependency
	// Unerase converts a value with erased type signature into a strongly typed value
	Unerase(val any) E.Either[error, T]
}

// InjectionToken uniquely identifies a dependency by giving it an Id, Type and name
type InjectionToken[T any] interface {
	Dependency[T]
	// Identity idenifies this dependency as a mandatory, required dependency, it will be resolved eagerly and injected as `T`.
	// If the dependency cannot be resolved, the resolution process fails
	Identity() Dependency[T]
	// Option identifies this dependency as optional, it will be resolved eagerly and injected as `O.Option[T]`.
	// If the dependency cannot be resolved, the resolution process continues and the dependency is represented as `O.None[T]`
	Option() Dependency[O.Option[T]]
	// IOEither identifies this dependency as mandatory but it will be resolved lazily as a `IOE.IOEither[error, T]`. This
	// value is memoized to make sure the dependency is a singleton.
	// If the dependency cannot be resolved, the resolution process fails
	IOEither() Dependency[IOE.IOEither[error, T]]
	// IOOption identifies this dependency as optional but it will be resolved lazily as a `IOO.IOOption[T]`. This
	// value is memoized to make sure the dependency is a singleton.
	// If the dependency cannot be resolved, the resolution process continues and the dependency is represented as the none value.
	IOOption() Dependency[IOO.IOOption[T]]
}

// MultiInjectionToken uniquely identifies a dependency by giving it an Id, Type and name.
type MultiInjectionToken[T any] interface {
	// Container returns the injection token used to request an array of all provided items
	Container() InjectionToken[[]T]
	// Item returns the injection token used to provide an item
	Item() InjectionToken[T]
}

// makeID creates a generator of unique string IDs
func makeId() IO.IO[string] {
	var count atomic.Int64
	return IO.MakeIO(func() string {
		return strconv.FormatInt(count.Add(1), 16)
	})
}

// genId is the common generator of unique string IDs
var genId = makeId()

type token[T any] struct {
	name            string
	id              string
	flag            int
	toType          func(val any) E.Either[error, T]
	providerFactory O.Option[DIE.ProviderFactory]
}

func (t *token[T]) Id() string {
	return t.id
}

func (t *token[T]) Flag() int {
	return t.flag
}

func (t *token[T]) String() string {
	return t.name
}

func (t *token[T]) Unerase(val any) E.Either[error, T] {
	return t.toType(val)
}

func (t *token[T]) ProviderFactory() O.Option[DIE.ProviderFactory] {
	return t.providerFactory
}

func makeToken[T any](name string, id string, typ int, unerase func(val any) E.Either[error, T], providerFactory O.Option[DIE.ProviderFactory]) Dependency[T] {
	return &token[T]{name, id, typ, unerase, providerFactory}
}

type injectionToken[T any] struct {
	token[T]
	option   Dependency[O.Option[T]]
	ioeither Dependency[IOE.IOEither[error, T]]
	iooption Dependency[IOO.IOOption[T]]
}

type multiInjectionToken[T any] struct {
	container *injectionToken[[]T]
	item      *injectionToken[T]
}

func (i *injectionToken[T]) Identity() Dependency[T] {
	return i
}

func (i *injectionToken[T]) Option() Dependency[O.Option[T]] {
	return i.option
}

func (i *injectionToken[T]) IOEither() Dependency[IOE.IOEither[error, T]] {
	return i.ioeither
}

func (i *injectionToken[T]) IOOption() Dependency[IOO.IOOption[T]] {
	return i.iooption
}

func (i *injectionToken[T]) ProviderFactory() O.Option[DIE.ProviderFactory] {
	return i.providerFactory
}

func (m *multiInjectionToken[T]) Container() InjectionToken[[]T] {
	return m.container
}

func (m *multiInjectionToken[T]) Item() InjectionToken[T] {
	return m.item
}

// makeToken create a unique `InjectionToken` for a specific type
func makeInjectionToken[T any](name string, providerFactory O.Option[DIE.ProviderFactory]) InjectionToken[T] {
	id := genId()
	toIdentity := toType[T]()
	return &injectionToken[T]{
		token[T]{name, id, DIE.Identity, toIdentity, providerFactory},
		makeToken[O.Option[T]](fmt.Sprintf("Option[%s]", name), id, DIE.Option, toOptionType(toIdentity), providerFactory),
		makeToken[IOE.IOEither[error, T]](fmt.Sprintf("IOEither[%s]", name), id, DIE.IOEither, toIOEitherType(toIdentity), providerFactory),
		makeToken[IOO.IOOption[T]](fmt.Sprintf("IOOption[%s]", name), id, DIE.IOOption, toIOOptionType(toIdentity), providerFactory),
	}
}

// MakeToken create a unique `InjectionToken` for a specific type
func MakeToken[T any](name string) InjectionToken[T] {
	return makeInjectionToken[T](name, O.None[DIE.ProviderFactory]())
}

// MakeToken create a unique `InjectionToken` for a specific type
func MakeTokenWithDefault[T any](name string, providerFactory DIE.ProviderFactory) InjectionToken[T] {
	return makeInjectionToken[T](name, O.Of(providerFactory))
}

// MakeMultiToken creates a [MultiInjectionToken]
func MakeMultiToken[T any](name string) MultiInjectionToken[T] {
	id := genId()
	toItem := toType[T]()
	toContainer := toArrayType(toItem)
	containerName := fmt.Sprintf("Container[%s]", name)
	itemName := fmt.Sprintf("Item[%s]", name)
	// empty factory
	providerFactory := O.None[DIE.ProviderFactory]()
	// container
	container := &injectionToken[[]T]{
		token[[]T]{containerName, id, DIE.Multi | DIE.Identity, toContainer, providerFactory},
		makeToken[O.Option[[]T]](fmt.Sprintf("Option[%s]", containerName), id, DIE.Multi|DIE.Option, toOptionType(toContainer), providerFactory),
		makeToken[IOE.IOEither[error, []T]](fmt.Sprintf("IOEither[%s]", containerName), id, DIE.Multi|DIE.IOEither, toIOEitherType(toContainer), providerFactory),
		makeToken[IOO.IOOption[[]T]](fmt.Sprintf("IOOption[%s]", containerName), id, DIE.Multi|DIE.IOOption, toIOOptionType(toContainer), providerFactory),
	}
	// item
	item := &injectionToken[T]{
		token[T]{itemName, id, DIE.Item | DIE.Identity, toItem, providerFactory},
		makeToken[O.Option[T]](fmt.Sprintf("Option[%s]", itemName), id, DIE.Item|DIE.Option, toOptionType(toItem), providerFactory),
		makeToken[IOE.IOEither[error, T]](fmt.Sprintf("IOEither[%s]", itemName), id, DIE.Item|DIE.IOEither, toIOEitherType(toItem), providerFactory),
		makeToken[IOO.IOOption[T]](fmt.Sprintf("IOOption[%s]", itemName), id, DIE.Item|DIE.IOOption, toIOOptionType(toItem), providerFactory),
	}
	// returns the token
	return &multiInjectionToken[T]{container, item}
}
