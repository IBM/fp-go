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

	FL "github.com/IBM/fp-go/v2/file"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
)

var (
	// readAll is the adapted version of [io.ReadAll]
	readAll = IOE.Eitherize1(io.ReadAll)
)

// ReadAll uses a generator function to create a stream, reads it and closes it
func ReadAll[R io.ReadCloser](acquire IOE.IOEither[error, R]) IOE.IOEither[error, []byte] {
	return F.Pipe1(
		F.Flow2(
			FL.ToReader[R],
			readAll,
		),
		IOE.WithResource[[]byte](acquire, Close[R]),
	)
}
