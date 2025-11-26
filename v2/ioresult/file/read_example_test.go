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

package file_test

import (
	"fmt"
	"io"
	"os"

	FL "github.com/IBM/fp-go/v2/file"
	F "github.com/IBM/fp-go/v2/function"
	I "github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/ioresult"
	"github.com/IBM/fp-go/v2/ioresult/file"
)

// Example_read_basicUsage demonstrates basic usage of the Read function
// to read data from a file with automatic resource cleanup.
func Example_read_basicUsage() {
	// Create a temporary file for demonstration
	tmpFile, err := os.CreateTemp("", "example-*.txt")
	if err != nil {
		fmt.Printf("Error creating temp file: %v\n", err)
		return
	}
	defer os.Remove(tmpFile.Name())

	// Write some test data
	testData := "Hello, World! This is a test file."
	if _, err := tmpFile.WriteString(testData); err != nil {
		fmt.Printf("Error writing to temp file: %v\n", err)
		return
	}
	tmpFile.Close()

	// Define a reader function that reads the full file content
	readAll := F.Flow2(
		FL.ToReader[*os.File],
		ioresult.Eitherize1(io.ReadAll),
	)

	content := F.Pipe2(
		readAll,
		file.Read[[]byte](file.Open(tmpFile.Name())),
		ioresult.TapIOK(I.Printf[[]byte]("%s\n")),
	)

	content()

	// Output: Hello, World! This is a test file.
}
