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

package generic

import (
	F "github.com/IBM/fp-go/v2/function"
	O "github.com/IBM/fp-go/v2/option"
	P "github.com/IBM/fp-go/v2/pair"
)

// Compress returns an [Iterator] that filters elements from a data [Iterator] returning only those that have a corresponding element in selector [Iterator] that evaluates to `true`.
// Stops when either the data or selectors iterator has been exhausted.
func Compress[GU ~func() Option[Pair[GU, U]], GB ~func() Option[Pair[GB, bool]], CS ~func() Option[Pair[CS, Pair[U, bool]]], U any](sel GB) func(GU) GU {
	return F.Flow2(
		Zip[GU, GB, CS](sel),
		FilterMap[GU, CS](F.Flow2(
			O.FromPredicate(P.Tail[U, bool]),
			O.Map(P.Head[U, bool]),
		)),
	)
}
