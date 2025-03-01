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

package option

import (
	F "github.com/IBM/fp-go/v2/function"
)

func toType[T any](a any) (T, bool) {
	b, ok := a.(T)
	return b, ok
}

func ToType[T any](src any) Option[T] {
	return F.Pipe1(
		src,
		Optionize1(toType[T]),
	)
}

func ToAny[T any](src T) Option[any] {
	return Of(any(src))
}
