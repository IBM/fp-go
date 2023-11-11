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

package function

import (
	G "github.com/IBM/fp-go/function/generic"
)

// Switch applies a handler to different cases. The handers are stored in a map. A key function
// extracts the case from a value.
func Switch[K comparable, T, R any](kf func(T) K, n map[K]func(T) R, d func(T) R) func(T) R {
	return G.Switch(kf, n, d)
}
