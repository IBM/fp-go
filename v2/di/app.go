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
	DIE "github.com/IBM/fp-go/v2/di/erasure"
	F "github.com/IBM/fp-go/v2/function"
	IO "github.com/IBM/fp-go/v2/io"
	IOR "github.com/IBM/fp-go/v2/ioresult"
)

var (
	// InjMain is the [InjectionToken] for the main application
	InjMain = MakeToken[any]("APP")

	// Main is the resolver for the main application
	Main = Resolve(InjMain)
)

// RunMain runs the main application from a set of [DIE.Provider]s
var RunMain = F.Flow3(
	DIE.MakeInjector,
	Main,
	IOR.Fold(IO.Of[error], F.Constant1[any](IO.Of[error](nil))),
)
