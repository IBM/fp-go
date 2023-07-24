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

package ioeither

import (
	G "github.com/IBM/fp-go/ioeither/generic"
)

// MonadApFirst combines two effectful actions, keeping only the result of the first.
func MonadApFirst[A, E, B any](first IOEither[E, A], second IOEither[E, B]) IOEither[E, A] {
	return G.MonadApFirst[IOEither[E, A], IOEither[E, B], IOEither[E, func(B) A]](first, second)
}

// ApFirst combines two effectful actions, keeping only the result of the first.
func ApFirst[A, E, B any](second IOEither[E, B]) func(IOEither[E, A]) IOEither[E, A] {
	return G.ApFirst[IOEither[E, A], IOEither[E, B], IOEither[E, func(B) A]](second)
}

// MonadApSecond combines two effectful actions, keeping only the result of the second.
func MonadApSecond[A, E, B any](first IOEither[E, A], second IOEither[E, B]) IOEither[E, B] {
	return G.MonadApSecond[IOEither[E, A], IOEither[E, B], IOEither[E, func(B) B]](first, second)
}

// ApSecond combines two effectful actions, keeping only the result of the second.
func ApSecond[A, E, B any](second IOEither[E, B]) func(IOEither[E, A]) IOEither[E, B] {
	return G.ApSecond[IOEither[E, A], IOEither[E, B], IOEither[E, func(B) B]](second)
}
