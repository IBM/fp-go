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

package file

import (
	"context"
	"io"
	"os"

	RIOE "github.com/IBM/fp-go/v2/context/readerioresult"
	ET "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/internal/file"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	IOEF "github.com/IBM/fp-go/v2/ioeither/file"
)

var (
	// Open opens a file for reading within the given context
	Open = F.Flow3(
		IOEF.Open,
		RIOE.FromIOEither[*os.File],
		RIOE.WithContext[*os.File],
	)

	// Remove removes a file by name
	Remove = F.Flow2(
		IOEF.Remove,
		RIOE.FromIOEither[string],
	)
)

// Close closes an object
func Close[C io.Closer](c C) RIOE.ReaderIOResult[struct{}] {
	return F.Pipe2(
		c,
		IOEF.Close[C],
		RIOE.FromIOEither[struct{}],
	)
}

// ReadFile reads a file in the scope of a context
func ReadFile(path string) RIOE.ReaderIOResult[[]byte] {
	return RIOE.WithResource[[]byte](Open(path), Close[*os.File])(func(r *os.File) RIOE.ReaderIOResult[[]byte] {
		return func(ctx context.Context) IOE.IOEither[error, []byte] {
			return func() ET.Either[error, []byte] {
				return file.ReadAll(ctx, r)
			}
		}
	})
}
