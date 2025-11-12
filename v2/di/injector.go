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
	"github.com/IBM/fp-go/v2/identity"
	IOR "github.com/IBM/fp-go/v2/ioresult"
	RIOR "github.com/IBM/fp-go/v2/readerioresult"
)

// Resolve performs a type safe resolution of a dependency
func Resolve[T any](token InjectionToken[T]) RIOR.ReaderIOResult[DIE.InjectableFactory, T] {
	return F.Flow2(
		identity.Ap[IOResult[any]](asDependency(token)),
		IOR.ChainResultK(token.Unerase),
	)
}
