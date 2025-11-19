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
	"io"

	"github.com/IBM/fp-go/v2/idiomatic/ioresult"
)

func onWriteAll[W io.Writer](data []byte) Kleisli[W, []byte] {
	return func(w W) IOResult[[]byte] {
		return func() ([]byte, error) {
			_, err := w.Write(data)
			return data, err
		}
	}
}

// WriteAll uses a generator function to create a stream, writes data to it and closes it
func WriteAll[W io.WriteCloser](data []byte) Operator[W, []byte] {
	onWrite := onWriteAll[W](data)
	return func(onCreate IOResult[W]) IOResult[[]byte] {
		return ioresult.WithResource[[]byte](
			onCreate,
			Close[W])(
			onWrite,
		)
	}
}

// Write uses a generator function to create a stream, writes data to it and closes it
func Write[R any, W io.WriteCloser](acquire IOResult[W]) Kleisli[Kleisli[W, R], R] {
	return ioresult.WithResource[R](
		acquire,
		Close[W])
}
