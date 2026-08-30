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

package pair

import (
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/optics/common"
	"github.com/IBM/fp-go/v2/pair"
)

// Head returns a Lens that focuses on the head (left) element of a Pair[L, R].
//
// The head is the first component of the pair. Getting the lens extracts the head value;
// setting a new value replaces the head while leaving the tail unchanged.
//
// Type Parameters:
//   - L: The type of the head element
//   - R: The type of the tail element
//
// Returns:
//   - A Lens[Pair[L, R], L] focused on the head element
//
// See Also:
//   - Tail: Lens focused on the tail (right) element
func Head[L, R any]() Lens[Pair[L, R], L] {
	return common.MakeLensCurriedWithName(
		pair.Head[L, R],
		F.Flow2(
			F.Constant1[L, L],
			pair.MapHead[R, L, L],
		),
		"pair.Head",
	)
}

// Tail returns a Lens that focuses on the tail (right) element of a Pair[L, R].
//
// The tail is the second component of the pair. Getting the lens extracts the tail value;
// setting a new value replaces the tail while leaving the head unchanged.
//
// Type Parameters:
//   - L: The type of the head element
//   - R: The type of the tail element
//
// Returns:
//   - A Lens[Pair[L, R], R] focused on the tail element
//
// See Also:
//   - Head: Lens focused on the head (left) element
func Tail[L, R any]() Lens[Pair[L, R], R] {
	return common.MakeLensCurriedWithName(
		pair.Tail[L, R],
		F.Flow2(
			F.Constant1[R, R],
			pair.MapTail[L, R, R],
		),
		"pair.Tail",
	)
}
