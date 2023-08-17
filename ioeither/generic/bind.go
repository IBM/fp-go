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
	ET "github.com/IBM/fp-go/either"
	G "github.com/IBM/fp-go/internal/bindt"
)

// Bind applies a function to an input state and merges the result into that state
func Bind[
	GS1 ~func() ET.Either[E, S1],
	GS2 ~func() ET.Either[E, S2],
	GA ~func() ET.Either[E, A],
	FCT ~func(S1) GA,
	E any,
	SET ~func(A) func(S1) S2,
	A, S1, S2 any](s SET) func(FCT) func(GS1) GS2 {
	return func(f FCT) func(GS1) GS2 {
		return G.Bind(
			Chain[GS1, GS2, E, S1, S2],
			Map[GA, GS2, E, A, S2],
			s,
			f,
		)
	}
}
