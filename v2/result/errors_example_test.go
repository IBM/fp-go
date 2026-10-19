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

package result

import (
	"errors"
	"fmt"
	"io/fs"
)

// ExampleErrorAs extracts a *fs.PathError from a wrapped error.
func ExampleErrorAs() {
	err := fmt.Errorf("config: %w", &fs.PathError{Op: "open", Path: "config.yaml", Err: fs.ErrNotExist})

	fmt.Println(ErrorAs[*fs.PathError](err))
	fmt.Println(ErrorAs[*fs.PathError](errors.New("boom")))
	// Output:
	// Right[*fs.PathError](open config.yaml: file does not exist)
	// Left[*errors.errorString](boom)
}

// ExampleErrorIs selects the errors that wrap fs.ErrNotExist.
func ExampleErrorIs() {
	notExist := ErrorIs(fs.ErrNotExist)

	fmt.Println(notExist(fmt.Errorf("open config.yaml: %w", fs.ErrNotExist)))
	fmt.Println(notExist(errors.New("boom")))
	// Output:
	// Right[*fmt.wrapError](open config.yaml: file does not exist)
	// Left[*errors.errorString](boom)
}
