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

package iterresult

import (
	"github.com/IBM/fp-go/v2/iterator/itereither"
)

// MonadApFirst combines two effectful actions, keeping only the result of the first.
//
// Marble diagram:
//
//	First:  ---R(1)---R(2)---|
//	Second: ---R(10)---R(20)---|
//	Output: ---R(1)---R(2)---|
//
// If either sequence contains a Left, the error is propagated:
//
//	First:  ---R(1)---L(e)---|
//	Second: ---R(10)---R(20)---|
//	Output: ---R(1)---L(e)---|
func MonadApFirst[A, B any](first SeqResult[A], second SeqResult[B]) SeqResult[A] {
	return itereither.MonadApFirst(first, second)
}

// ApFirst combines two effectful actions, keeping only the result of the first.
//
// Marble diagram:
//
//	First:  ---R(1)---R(2)---|
//	Second: ---R(10)---R(20)---|
//	Output: ---R(1)---R(2)---|
//
// If either sequence contains a Left, the error is propagated.
func ApFirst[A, B any](second SeqResult[B]) Operator[A, A] {
	return itereither.ApFirst[A](second)
}

// MonadApSecond combines two effectful actions, keeping only the result of the second.
//
// Marble diagram:
//
//	First:  ---R(1)---R(2)---|
//	Second: ---R(10)---R(20)---|
//	Output: ---R(10)---R(20)---|
//
// If either sequence contains a Left, the error is propagated:
//
//	First:  ---R(1)---L(e)---|
//	Second: ---R(10)---R(20)---|
//	Output: ---R(10)---L(e)---|
func MonadApSecond[A, B any](first SeqResult[A], second SeqResult[B]) SeqResult[B] {
	return itereither.MonadApSecond(first, second)
}

// ApSecond combines two effectful actions, keeping only the result of the second.
//
// Marble diagram:
//
//	First:  ---R(1)---R(2)---|
//	Second: ---R(10)---R(20)---|
//	Output: ---R(10)---R(20)---|
//
// If either sequence contains a Left, the error is propagated.
func ApSecond[A, B any](second SeqResult[B]) Operator[A, B] {
	return itereither.ApSecond[A](second)
}
