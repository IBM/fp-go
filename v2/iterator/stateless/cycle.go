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
)

// DropWhile creates an [Iterator] that drops elements from the [Iterator] as long as the predicate is true; afterwards, returns every element.
// Note, the [Iterator] does not produce any output until the predicate first becomes false
func Cycle[U any](ma Iterator[U]) Iterator[U] {
	return G.Cycle[Iterator[U]](ma)
}
