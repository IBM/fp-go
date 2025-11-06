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
	"bytes"
	"context"
	"io"

	E "github.com/IBM/fp-go/v2/either"
)

type (
	readerWithContext struct {
		ctx      context.Context
		delegate io.Reader
	}
)

func (rdr *readerWithContext) Read(p []byte) (int, error) {
	// check for cancellarion
	if err := rdr.ctx.Err(); err != nil {
		return 0, err
	}
	// simply dispatch
	return rdr.delegate.Read(p)
}

// MakeReader creates a context aware reader
func MakeReader(ctx context.Context, rdr io.Reader) io.Reader {
	return &readerWithContext{ctx, rdr}
}

// ReadAll reads the content of a reader and allows it to be canceled
func ReadAll(ctx context.Context, rdr io.Reader) E.Either[error, []byte] {
	var buffer bytes.Buffer
	_, err := io.Copy(&buffer, MakeReader(ctx, rdr))
	return E.TryCatchError(buffer.Bytes(), err)
}
