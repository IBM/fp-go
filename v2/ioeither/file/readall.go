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

	FL "github.com/IBM/fp-go/v2/file"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/ioeither"
)

var (
	// readAll is the adapted version of [io.ReadAll]
	readAll = ioeither.Eitherize1(io.ReadAll)
)

// ReadAll reads all data from a ReadCloser and ensures it is properly closed.
// It takes an IOEither that acquires the ReadCloser, reads all its content until EOF,
// and automatically closes the reader, even if an error occurs during reading.
//
// This is the recommended way to read entire files with proper resource management.
//
// Example:
//
//	readOp := ReadAll(Open("input.txt"))
//	result := readOp() // Either[error, []byte]
func ReadAll[R io.ReadCloser](acquire IOEither[error, R]) IOEither[error, []byte] {
	return F.Pipe1(
		F.Flow2(
			FL.ToReader[R],
			readAll,
		),
		ioeither.WithResource[[]byte](acquire, Close[R]),
	)
}
