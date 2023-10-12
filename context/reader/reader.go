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

package reader

import (
	"context"

	R "github.com/IBM/fp-go/reader/generic"
)

func MonadMap[A, B any](fa Reader[A], f func(A) B) Reader[B] {
	return R.MonadMap[Reader[A], Reader[B]](fa, f)
}

func Map[A, B any](f func(A) B) func(Reader[A]) Reader[B] {
	return R.Map[Reader[A], Reader[B]](f)
}

func MonadChain[A, B any](ma Reader[A], f func(A) Reader[B]) Reader[B] {
	return R.MonadChain(ma, f)
}

func Chain[A, B any](f func(A) Reader[B]) func(Reader[A]) Reader[B] {
	return R.Chain[Reader[A]](f)
}

func Of[A any](a A) Reader[A] {
	return R.Of[Reader[A]](a)
}

func MonadAp[A, B any](fab Reader[func(A) B], fa Reader[A]) Reader[B] {
	return R.MonadAp[Reader[A], Reader[B]](fab, fa)
}

func Ap[A, B any](fa Reader[A]) func(Reader[func(A) B]) Reader[B] {
	return R.Ap[Reader[A], Reader[B], Reader[func(A) B]](fa)
}

func Ask() Reader[context.Context] {
	return R.Ask[Reader[context.Context]]()
}

func Asks[A any](r Reader[A]) Reader[A] {
	return R.Asks(r)
}

func MonadFlap[B, A any](fab Reader[func(A) B], a A) Reader[B] {
	return R.MonadFlap[Reader[func(A) B], Reader[B]](fab, a)
}

func Flap[B, A any](a A) func(Reader[func(A) B]) Reader[B] {
	return R.Flap[Reader[func(A) B], Reader[B]](a)
}
