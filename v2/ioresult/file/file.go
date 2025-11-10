// Copyright (c) 2023 - 2025 IBM Corp.
// All rights reserved.
//
// Licensed under the Apache LicensVersion 2.0 (the "License");
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

	"github.com/IBM/fp-go/v2/ioeither/file"
)

var (
	// Open opens a file for reading
	Open = file.Open
	// Create opens a file for writing
	Create = file.Create
	// ReadFile reads the context of a file
	ReadFile = file.ReadFile
)

// WriteFile writes a data blob to a file
//
//go:inline
func WriteFile(dstName string, perm os.FileMode) Kleisli[[]byte, []byte] {
	return file.WriteFile(dstName, perm)
}

// Remove removes a file by name
//
//go:inline
func Remove(name string) IOResult[string] {
	return file.Remove(name)
}

// Close closes an object
//
//go:inline
func Close[C io.Closer](c C) IOResult[any] {
	return file.Close(c)
}
