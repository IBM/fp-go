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

package json

import (
	E "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	O "github.com/IBM/fp-go/option"
)

func ToTypeE[A any](src any) E.Either[error, A] {
	return F.Pipe2(
		src,
		Marshal[any],
		E.Chain(Unmarshal[A]),
	)
}

func ToTypeO[A any](src any) O.Option[A] {
	return F.Pipe1(
		ToTypeE[A](src),
		E.ToOption[error, A](),
	)
}
