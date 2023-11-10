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
	G "github.com/IBM/fp-go/iterator/stateless/generic"
)

// Any returns `true` if any element of the iterable is `true`. If the iterable is empty, return `false`
// Similar to the [https://docs.python.org/3/library/functions.html#any] function
func Any[U any](pred func(U) bool) func(ma Iterator[U]) bool {
	return G.Any[Iterator[U]](pred)
}
