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

package readerio

import (
	G "github.com/IBM/fp-go/v2/readerio/generic"
)

// MonadApFirst combines two effectful actions, keeping only the result of the first.
func MonadApFirst[A, R, B any](first ReaderIO[R, A], second ReaderIO[R, B]) ReaderIO[R, A] {
	return G.MonadApFirst[ReaderIO[R, A], ReaderIO[R, B], ReaderIO[R, func(B) A]](first, second)
}

// ApFirst combines two effectful actions, keeping only the result of the first.
func ApFirst[A, R, B any](second ReaderIO[R, B]) func(ReaderIO[R, A]) ReaderIO[R, A] {
	return G.ApFirst[ReaderIO[R, A], ReaderIO[R, B], ReaderIO[R, func(B) A]](second)
}

// MonadApSecond combines two effectful actions, keeping only the result of the second.
func MonadApSecond[A, R, B any](first ReaderIO[R, A], second ReaderIO[R, B]) ReaderIO[R, B] {
	return G.MonadApSecond[ReaderIO[R, A], ReaderIO[R, B], ReaderIO[R, func(B) B]](first, second)
}

// ApSecond combines two effectful actions, keeping only the result of the second.
func ApSecond[A, R, B any](second ReaderIO[R, B]) func(ReaderIO[R, A]) ReaderIO[R, B] {
	return G.ApSecond[ReaderIO[R, A], ReaderIO[R, B], ReaderIO[R, func(B) B]](second)
}
