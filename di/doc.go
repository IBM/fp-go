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

// Package implements functions and data types supporting dependency injection patterns
//
// The fundamental building block is the concept of a [Dependency]. This describes the abstract concept of a function, service or value together with its type.
// Examples for dependencies can be as simple as configuration values such as the API URL for a service, the current username, a [Dependency] could be the map
// of the configuration environment, an http client or as complex as a service interface. Important is that a [Dependency] only defines the concept but
// not the implementation.
//
// The implementation of a [Dependency] is called a [Provider], the dependency of an `API URL` could e.g. be realized by a provider that consults the environment to read the information
// or a config file or simply hardcode it.
// In many cases the implementation of a [DIE.Provider] depends in turn on other [Dependency]s (but never directly on other [DIE.Provider]s), a provider for an `API URL` that reads
// the information from the environment would e.g. depend on a [Dependency] that represents this environment.
//
// It is the resposibility of the [DIE.InjectableFactory] to
//
// [Provider]: [github.com/IBM/fp-go/di/erasure.Provider]
package di

//go:generate go run .. di --count 10 --filename gen.go
