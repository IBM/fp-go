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

// Package di implements functions and data types supporting dependency injection patterns
//
// The fundamental building block is the concept of a [Dependency]. This describes the abstract concept of a function, service or value together with its type.
// Examples for dependencies can be as simple as configuration values such as the API URL for a service, the current username, a [Dependency] could be the map
// of the configuration environment, an http client or as complex as a service interface. Important is that a [Dependency] only defines the concept but
// not the implementation.
//
// The implementation of a [Dependency] is called a [Provider], the dependency of an `API URL` could e.g. be realized by a provider that consults the environment to read the information
// or a config file or simply hardcode it.
// In many cases the implementation of a [Provider] depends in turn on other [Dependency] (but never directly on other [Provider]s), a provider for an `API URL` that reads
// the information from the environment would e.g. depend on a [Dependency] that represents this environment.
//
// It is the resposibility of the [InjectableFactory] to create an instance of a [Dependency]. All instances are considered singletons. Create an [InjectableFactory] via the [MakeInjector] method. Use [Resolve] to create a strongly typed
// factory for a particular [InjectionToken].
//
// In most cases it is not required to use [InjectableFactory] directly, instead you would create a [Provider] via the [MakeProvider2] method (suffix indicates the number of dependencies). In this call
// you give a number of (strongly typed) [Dependency] identifiers and a (strongly typed) factory function for the implementation. The dependency injection framework makes
// sure to resolve the dependencies before calling the factory method.
//
// For convenience purposes it can be helpful to attach a default implementation of a [Dependency]. In this case use the [MakeTokenWithDefault2] method (suffix indicates the number of dependencies)
// to define the [InjectionToken].
//
// [Provider]: [github.com/IBM/fp-go/v2/di/erasure.Provider]
// [InjectableFactory]: [github.com/IBM/fp-go/v2/di/erasure.InjectableFactory]
// [MakeInjector]: [github.com/IBM/fp-go/v2/di/erasure.MakeInjector]
package di

//go:generate go run .. di --count 15 --filename gen.go
