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
	"os"

	RIOE "github.com/IBM/fp-go/v2/context/readerioresult"
	F "github.com/IBM/fp-go/v2/function"
	IO "github.com/IBM/fp-go/v2/io"
	IOF "github.com/IBM/fp-go/v2/io/file"
	IOEF "github.com/IBM/fp-go/v2/ioeither/file"
)

var (
	// onCreateTempFile creates a temp file with sensible defaults
	onCreateTempFile = CreateTemp("", "*")
	// destroy handler
	onReleaseTempFile = F.Flow4(
		IOF.Close[*os.File],
		IO.Map((*os.File).Name),
		RIOE.FromIO[string],
		RIOE.Chain(Remove),
	)
)

// CreateTemp created a temp file with proper parametrization
func CreateTemp(dir, pattern string) RIOE.ReaderIOResult[*os.File] {
	return F.Pipe2(
		IOEF.CreateTemp(dir, pattern),
		RIOE.FromIOEither[*os.File],
		RIOE.WithContext[*os.File],
	)
}

// WithTempFile creates a temporary file, then invokes a callback to create a resource based on the file, then close and remove the temp file
func WithTempFile[A any](f func(*os.File) RIOE.ReaderIOResult[A]) RIOE.ReaderIOResult[A] {
	return RIOE.WithResource[A](onCreateTempFile, onReleaseTempFile)(f)
}
