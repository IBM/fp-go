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

package option

import (
	T "github.com/IBM/fp-go/optics/traversal/generic"
	O "github.com/IBM/fp-go/option"
)

type (
	Traversal[S, A any] T.Traversal[S, A, O.Option[S], O.Option[A]]
)

func Compose[
	S, A, B any](ab Traversal[A, B]) func(Traversal[S, A]) Traversal[S, B] {
	return T.Compose[
		Traversal[A, B],
		Traversal[S, A],
		Traversal[S, B],
	](ab)
}
