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

package generic

import (
	I "github.com/IBM/fp-go/identity/generic"
)

func MonadAp[GA ~func(A) A, A any](fab GA, fa A) A {
	return I.MonadAp[GA, A, A](fab, fa)
}

func Ap[GA ~func(A) A, A any](fa A) func(GA) A {
	return I.Ap[GA, A, A](fa)
}

func MonadFlap[GA ~func(A) A, A any](fab GA, a A) A {
	return I.MonadFlap[GA, A, A](fab, a)
}

func Flap[GA ~func(A) A, A any](a A) func(GA) A {
	return I.Flap[GA, A, A](a)
}

func MonadMap[GA ~func(A) A, A any](fa A, f GA) A {
	return I.MonadMap[GA, A, A](fa, f)
}

func Map[GA ~func(A) A, A any](f GA) GA {
	return I.Map[GA, A, A](f)
}

func MonadChain[GA ~func(A) A, A any](ma A, f GA) A {
	return I.MonadChain[GA, A, A](ma, f)
}

func Chain[GA ~func(A) A, A any](f GA) GA {
	return I.Chain[GA, A](f)
}

func MonadChainFirst[GA ~func(A) A, A any](fa A, f GA) A {
	return I.MonadChainFirst[GA, A, A](fa, f)
}

func ChainFirst[GA ~func(A) A, A any](f GA) GA {
	return I.ChainFirst[GA, A, A](f)
}
