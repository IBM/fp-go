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
	"os"

	"github.com/IBM/fp-go/v2/ioeither"
)

var (
	// Open opens a file for reading
	Open = ioeither.Eitherize1(os.Open)
	// Create opens a file for writing
	Create = ioeither.Eitherize1(os.Create)
	// ReadFile reads the context of a file
	ReadFile = ioeither.Eitherize1(os.ReadFile)
)

// WriteFile writes a data blob to a file
func WriteFile(dstName string, perm os.FileMode) func([]byte) ioeither.IOEither[error, []byte] {
	return func(data []byte) ioeither.IOEither[error, []byte] {
		return ioeither.TryCatchError(func() ([]byte, error) {
			return data, os.WriteFile(dstName, data, perm)
		})
	}
}

// Remove removes a file by name
func Remove(name string) ioeither.IOEither[error, string] {
	return ioeither.TryCatchError(func() (string, error) {
		return name, os.Remove(name)
	})
}

// Close closes an object
func Close[C io.Closer](c C) ioeither.IOEither[error, any] {
	return ioeither.TryCatchError(func() (any, error) {
		return c, c.Close()
	})
}
