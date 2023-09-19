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

	IOE "github.com/IBM/fp-go/ioeither"
)

func onReadAll[R io.Reader](r R) IOE.IOEither[error, []byte] {
	return IOE.TryCatchError(func() ([]byte, error) {
		return io.ReadAll(r)
	})
}

// ReadAll uses a generator function to create a stream, reads it and closes it
func ReadAll[R io.ReadCloser](acquire IOE.IOEither[error, R]) IOE.IOEither[error, []byte] {
	return IOE.WithResource[[]byte](
		acquire,
		Close[R])(
		onReadAll[R],
	)
}
