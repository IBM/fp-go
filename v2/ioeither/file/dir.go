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

	"github.com/IBM/fp-go/v2/ioeither"
)

// MkdirAll creates a directory and all necessary parent directories with the specified permissions.
// If the directory already exists, MkdirAll does nothing and returns success.
// This is equivalent to the Unix command `mkdir -p`.
//
// The perm parameter specifies the Unix permission bits for the created directories.
// Common values include 0755 (rwxr-xr-x) for directories.
//
// Returns an IOEither that, when executed, creates the directory structure and returns
// the path on success or an error on failure.
//
// See [os.MkdirAll] for more details.
//
// Example:
//
//	mkdirOp := MkdirAll("/tmp/my/nested/dir", 0755)
//	result := mkdirOp() // Either[error, string]
func MkdirAll(path string, perm os.FileMode) IOEither[error, string] {
	return ioeither.TryCatchError(func() (string, error) {
		return path, os.MkdirAll(path, perm)
	})
}

// Mkdir creates a single directory with the specified permissions.
// Unlike MkdirAll, it returns an error if the parent directory does not exist
// or if the directory already exists.
//
// The perm parameter specifies the Unix permission bits for the created directory.
// Common values include 0755 (rwxr-xr-x) for directories.
//
// Returns an IOEither that, when executed, creates the directory and returns
// the path on success or an error on failure.
//
// See [os.Mkdir] for more details.
//
// Example:
//
//	mkdirOp := Mkdir("/tmp/mydir", 0755)
//	result := mkdirOp() // Either[error, string]
func Mkdir(path string, perm os.FileMode) IOEither[error, string] {
	return ioeither.TryCatchError(func() (string, error) {
		return path, os.Mkdir(path, perm)
	})
}
