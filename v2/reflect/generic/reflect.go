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
	R "reflect"
)

func Map[GA ~[]A, A any](f func(R.Value) A) func(R.Value) GA {
	return func(val R.Value) GA {
		l := val.Len()
		res := make(GA, l)
		for i := l - 1; i >= 0; i-- {
			res[i] = f(val.Index(i))
		}
		return res
	}
}
