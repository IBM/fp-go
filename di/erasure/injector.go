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
	"log"

	A "github.com/IBM/fp-go/array"
	"github.com/IBM/fp-go/errors"
	F "github.com/IBM/fp-go/function"
	I "github.com/IBM/fp-go/identity"
	IG "github.com/IBM/fp-go/identity/generic"
	IOE "github.com/IBM/fp-go/ioeither"
	L "github.com/IBM/fp-go/lazy"
	O "github.com/IBM/fp-go/option"
	RIOE "github.com/IBM/fp-go/readerioeither"
	R "github.com/IBM/fp-go/record"
	T "github.com/IBM/fp-go/tuple"

	"sync"
)

func providerToEntry(p Provider) T.Tuple2[string, ProviderFactory] {
	return T.MakeTuple2(p.Provides().Id(), p.Factory())
}

var missingProviderError = F.Flow3(
	errors.OnSome[string]("no provider for dependency [%s]"),
	RIOE.Left[InjectableFactory, any, error],
	F.Constant[ProviderFactory],
)

func logEntryExit(name string, token Dependency) func() {
	log.Printf("Entry: [%s] -> [%s]:[%s]", name, token.Id(), token.String())
	return func() {
		log.Printf("Exit:  [%s] -> [%s]:[%s]", name, token.Id(), token.String())
	}
}

func MakeInjector(providers []Provider) InjectableFactory {

	type Result = IOE.IOEither[error, any]
	type LazyResult = L.Lazy[Result]

	// resolved stores the values resolved so far, key is the string ID
	// of the token, value is a lazy result
	var resolved sync.Map

	// provide a mapping for all providers
	factoryById := F.Pipe2(
		providers,
		A.Map(providerToEntry),
		R.FromEntries[string, ProviderFactory],
	)

	// the actual factory, we need lazy initialization
	var injFct InjectableFactory

	// lazy initialization, so we can cross reference it
	injFct = func(token Dependency) Result {

		defer logEntryExit("inj", token)()

		key := token.Id()

		// according to https://github.com/golang/go/issues/44159 this
		// is the best way to use the sync map
		actual, loaded := resolved.Load(key)
		if !loaded {

			computeResult := func() Result {
				defer logEntryExit("computeResult", token)()
				return F.Pipe5(
					token,
					T.Replicate2[Dependency],
					T.Map2(F.Flow3(
						Dependency.Id,
						R.Lookup[ProviderFactory, string],
						I.Ap[O.Option[ProviderFactory]](factoryById),
					), F.Flow2(
						Dependency.String,
						missingProviderError,
					)),
					T.Tupled2(O.MonadGetOrElse[ProviderFactory]),
					IG.Ap[ProviderFactory](injFct),
					IOE.Memoize[error, any],
				)
			}

			actual, _ = resolved.LoadOrStore(key, F.Pipe1(
				computeResult,
				L.Memoize[Result],
			))
		}

		return actual.(LazyResult)()
	}

	return injFct
}
