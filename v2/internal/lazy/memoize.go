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

package lazy

import "sync"

// Memoize computes the value of the provided IO monad lazily but exactly once
func Memoize[GA ~func() A, A any](ma GA) GA {
	// synchronization primitives
	var once sync.Once
	var result A
	// callback
	gen := func() {
		result = ma()
	}
	// returns our memoized wrapper
	return func() A {
		once.Do(gen)
		return result
	}
}
