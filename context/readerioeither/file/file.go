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

package file

import (
	"context"
	"io"
	"os"

	RIOE "github.com/IBM/fp-go/context/readerioeither"
	ET "github.com/IBM/fp-go/either"
	F "github.com/IBM/fp-go/function"
	"github.com/IBM/fp-go/internal/file"
	IOE "github.com/IBM/fp-go/ioeither"
)

var (
	openIOE = IOE.Eitherize1(os.Open)
	// Open opens a file for reading within the given context
	Open = F.Flow3(
		openIOE,
		RIOE.FromIOEither[*os.File],
		RIOE.WithContext[*os.File],
	)
)

// Close closes an object
func Close[C io.Closer](c C) RIOE.ReaderIOEither[any] {
	return RIOE.FromIOEither(func() ET.Either[error, any] {
		return ET.TryCatchError(func() (any, error) {
			return c, c.Close()
		})
	})
}

// ReadFile reads a file in the scope of a context
func ReadFile(path string) RIOE.ReaderIOEither[[]byte] {
	return RIOE.WithResource[*os.File, []byte](Open(path), Close[*os.File])(func(r *os.File) RIOE.ReaderIOEither[[]byte] {
		return func(ctx context.Context) IOE.IOEither[error, []byte] {
			return IOE.MakeIO(func() ET.Either[error, []byte] {
				return file.ReadAll(ctx, r)
			})
		}
	})
}
