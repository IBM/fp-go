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

package erasure

import (
	"fmt"

	A "github.com/IBM/fp-go/array"
	F "github.com/IBM/fp-go/function"
	I "github.com/IBM/fp-go/identity"
	IOE "github.com/IBM/fp-go/ioeither"
	O "github.com/IBM/fp-go/option"
	R "github.com/IBM/fp-go/record"
	T "github.com/IBM/fp-go/tuple"
)

func providerToEntry(p Provider) T.Tuple2[string, ProviderFactory] {
	return T.MakeTuple2(p.Provides().Id(), p.Factory())
}

func missingProviderError(name string) func() IOE.IOEither[error, any] {
	return func() IOE.IOEither[error, any] {
		return IOE.Left[any](fmt.Errorf("No provider for dependency [%s]", name))
	}
}

func MakeInjector(providers []Provider) InjectableFactory {

	type Result = IOE.IOEither[error, any]

	// provide a mapping for all providers
	factoryById := F.Pipe2(
		providers,
		A.Map(providerToEntry),
		R.FromEntries[string, ProviderFactory],
	)
	// the resolved map
	var resolved = R.Empty[string, Result]()
	// the callback
	var injFct InjectableFactory

	// lazy initialization, so we can cross reference it
	injFct = func(token Token) Result {

		hit := F.Pipe3(
			token,
			Token.Id,
			R.Lookup[Result, string],
			I.Ap[O.Option[Result]](resolved),
		)

		provFct := F.Pipe2(
			token,
			T.Replicate2[Token],
			T.Map2(F.Flow3(
				Token.Id,
				R.Lookup[ProviderFactory, string],
				I.Ap[O.Option[ProviderFactory]](factoryById),
			), F.Flow2(
				Token.String,
				missingProviderError,
			)),
		)

		x := F.Pipe4(
			token,
			Token.Id,
			R.Lookup[Result, string],
			I.Ap[O.Option[Result]](resolved),
			O.GetOrElse(F.Flow2()),
		)
	}

	return injFct
}
