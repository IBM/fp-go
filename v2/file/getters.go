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
	"path/filepath"
)

// Join appends a filename to a root path
func Join(name string) func(root string) string {
	return func(root string) string {
		return filepath.Join(root, name)
	}
}

// ToReader converts a [io.Reader]
func ToReader[R io.Reader](r R) io.Reader {
	return r
}

// ToWriter converts a [io.Writer]
func ToWriter[W io.Writer](w W) io.Writer {
	return w
}

// ToCloser converts a [io.Closer]
func ToCloser[C io.Closer](c C) io.Closer {
	return c
}
