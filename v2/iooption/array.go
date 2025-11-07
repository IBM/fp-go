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

package iooption

import (
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/option"
)

// TraverseArray transforms an array
func TraverseArray[A, B any](f Kleisli[A, B]) Kleisli[[]A, []B] {
	return function.Flow2(
		io.TraverseArray(f),
		io.Map(option.SequenceArray[B]),
	)
}

// TraverseArrayWithIndex transforms an array
func TraverseArrayWithIndex[A, B any](f func(int, A) IOOption[B]) Kleisli[[]A, []B] {
	return function.Flow2(
		io.TraverseArrayWithIndex(f),
		io.Map(option.SequenceArray[B]),
	)
}

// SequenceArray converts a homogeneous sequence of either into an either of sequence
func SequenceArray[A any](ma []IOOption[A]) IOOption[[]A] {
	return TraverseArray(function.Identity[IOOption[A]])(ma)
}
