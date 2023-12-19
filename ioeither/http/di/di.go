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
	"net/http"

	DI "github.com/IBM/fp-go/di"
	IOE "github.com/IBM/fp-go/ioeither"
	IOEH "github.com/IBM/fp-go/ioeither/http"
)

var (
	// InjHttpClient is the [DI.InjectionToken] for the [http.DefaultClient]
	InjHttpClient = DI.MakeTokenWithDefault0("HTTP_CLIENT", IOE.Of[error](http.DefaultClient))

	// InjClient is the [DI.InjectionToken] for the default [IOEH.Client]
	InjClient = DI.MakeTokenWithDefault1("CLIENT", InjHttpClient.IOEither(), IOE.Map[error](IOEH.MakeClient))
)
