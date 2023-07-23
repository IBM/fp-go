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

package number

import (
	M "github.com/IBM/fp-go/magma"
)

func MagmaSub[A int | int8 | int16 | int32 | int64 | float32 | float64 | complex64 | complex128]() M.Magma[A] {
	return M.MakeMagma(func(first A, second A) A {
		return first - second
	})
}
