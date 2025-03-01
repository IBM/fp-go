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

package stateless

import (
	G "github.com/IBM/fp-go/v2/iterator/stateless/generic"
	P "github.com/IBM/fp-go/v2/pair"
)

// ZipWith applies a function to pairs of elements at the same index in two iterators, collecting the results in a new iterator. If one
// input iterator is short, excess elements of the longer iterator are discarded.
func ZipWith[FCT ~func(A, B) C, A, B, C any](fa Iterator[A], fb Iterator[B], f FCT) Iterator[C] {
	return G.ZipWith[Iterator[A], Iterator[B], Iterator[C]](fa, fb, f)
}

// Zip takes two iterators and returns an iterators of corresponding pairs. If one input iterators is short, excess elements of the
// longer iterator are discarded
func Zip[A, B any](fb Iterator[B]) func(Iterator[A]) Iterator[P.Pair[A, B]] {
	return G.Zip[Iterator[A], Iterator[B], Iterator[P.Pair[A, B]]](fb)
}
