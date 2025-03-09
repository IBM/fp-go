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
	E "github.com/IBM/fp-go/v2/either"
	"github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/option"
)

type (
	Either[A any] = E.Either[error, A]
	Option[A any] = option.Option[A]
)

func ToTypeE[A any](src any) Either[A] {
	return function.Pipe2(
		src,
		Marshal[any],
		E.Chain(Unmarshal[A]),
	)
}

func ToTypeO[A any](src any) Option[A] {
	return function.Pipe1(
		ToTypeE[A](src),
		E.ToOption[error, A],
	)
}
