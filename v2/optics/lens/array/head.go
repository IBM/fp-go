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

package array

import (
	L "github.com/IBM/fp-go/v2/optics/lens"
	G "github.com/IBM/fp-go/v2/optics/lens/array/generic"
	O "github.com/IBM/fp-go/v2/option"
)

// AtHead focusses on the head of an array. The setter works as follows
// - if the new value is none, the result will be an empty array
// - if the new value is some and the array is empty, it creates a new array with one element
// - if the new value is some and the array is not empty, it replaces the head
func AtHead[A any]() L.Lens[[]A, O.Option[A]] {
	return G.AtHead[[]A]()
}
