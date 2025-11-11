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

package stateless

import (
	IO "github.com/IBM/fp-go/v2/io"
	G "github.com/IBM/fp-go/v2/iterator/stateless/generic"
	L "github.com/IBM/fp-go/v2/lazy"
)

// FromLazy returns an [Iterator] on top of a lazy function
func FromLazy[U any](l L.Lazy[U]) Iterator[U] {
	return G.FromLazy[Iterator[U]](l)
}

// FromIO returns an [Iterator] on top of an IO function
func FromIO[U any](io IO.IO[U]) Iterator[U] {
	return G.FromLazy[Iterator[U]](io)
}
