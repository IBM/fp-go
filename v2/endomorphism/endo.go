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

package endomorphism

import (
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/identity"
)

func MonadAp[A any](fab Endomorphism[A], fa A) A {
	return identity.MonadAp(fab, fa)
}

func Ap[A any](fa A) func(Endomorphism[A]) A {
	return identity.Ap[A](fa)
}

func Compose[A any](f1, f2 Endomorphism[A]) Endomorphism[A] {
	return function.Flow2(f1, f2)
}

func MonadChain[A any](ma Endomorphism[A], f Endomorphism[A]) Endomorphism[A] {
	return Compose(ma, f)
}

func Chain[A any](f Endomorphism[A]) Endomorphism[Endomorphism[A]] {
	return function.Bind2nd(MonadChain, f)
}
