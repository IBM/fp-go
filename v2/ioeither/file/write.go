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
	"io"

	"github.com/IBM/fp-go/v2/ioeither"
)

func onWriteAll[W io.Writer](data []byte) func(w W) ioeither.IOEither[error, []byte] {
	return func(w W) ioeither.IOEither[error, []byte] {
		return ioeither.TryCatchError(func() ([]byte, error) {
			_, err := w.Write(data)
			return data, err
		})
	}
}

// WriteAll uses a generator function to create a stream, writes data to it and closes it
func WriteAll[W io.WriteCloser](data []byte) func(acquire ioeither.IOEither[error, W]) ioeither.IOEither[error, []byte] {
	onWrite := onWriteAll[W](data)
	return func(onCreate ioeither.IOEither[error, W]) ioeither.IOEither[error, []byte] {
		return ioeither.WithResource[[]byte](
			onCreate,
			Close[W])(
			onWrite,
		)
	}
}

// Write uses a generator function to create a stream, writes data to it and closes it
func Write[R any, W io.WriteCloser](acquire ioeither.IOEither[error, W]) func(use func(W) ioeither.IOEither[error, R]) ioeither.IOEither[error, R] {
	return ioeither.WithResource[R](
		acquire,
		Close[W])
}
