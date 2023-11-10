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
	F "github.com/IBM/fp-go/function"
	O "github.com/IBM/fp-go/option"
	T "github.com/IBM/fp-go/tuple"
)

// Any returns `true` if any element of the iterable is `true`. If the iterable is empty, return `false`
func Any[GU ~func() O.Option[T.Tuple2[GU, U]], FCT ~func(U) bool, U any](pred FCT) func(ma GU) bool {
	return F.Flow3(
		Filter[GU](pred),
		First[GU],
		O.Fold(F.ConstFalse, F.Constant1[U](true)),
	)
}
