// Copyright (c) 2025 IBM Corp.
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

package reader

type (
	// The purpose of the `Reader` monad is to avoid threading arguments through multiple functions in order to only get them where they are needed.
	// The first template argument `R` is the the context to read from, the second argument `A` is the return value of the monad
	Reader[R, A any] = func(R) A

	Mapper[R, A, B any] = func(Reader[R, A]) Reader[R, B]
)
