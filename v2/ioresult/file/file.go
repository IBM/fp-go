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
	// Stat returns [FileInfo] object
	Stat = file.Stat

	// UserCacheDir returns an [IOResult] that resolves to the default root directory
	// to use for user-specific cached data. Users should create their own application-specific
	// subdirectory within this one and use that.
	//
	// On Unix systems, it returns $XDG_CACHE_HOME as specified by
	// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html if
	// non-empty, else $HOME/.cache.
	// On Darwin, it returns $HOME/Library/Caches.
	// On Windows, it returns %LocalAppData%.
	// On Plan 9, it returns $home/lib/cache.
	//
	// If the location cannot be determined (for example, $HOME is not defined),
	// then it will return an error wrapped in [Err].
	UserCacheDir = file.UserCacheDir

	// UserConfigDir returns an [IOResult] that resolves to the default root directory
	// to use for user-specific configuration data. Users should create their own
	// application-specific subdirectory within this one and use that.
	//
	// On Unix systems, it returns $XDG_CONFIG_HOME as specified by
	// https://specifications.freedesktop.org/basedir-spec/basedir-spec-latest.html if
	// non-empty, else $HOME/.config.
	// On Darwin, it returns $HOME/Library/Application Support.
	// On Windows, it returns %AppData%.
	// On Plan 9, it returns $home/lib.
	//
	// If the location cannot be determined (for example, $HOME is not defined),
	// then it will return an error wrapped in [Err].
	UserConfigDir = file.UserConfigDir

	// UserHomeDir returns an [IOResult] that resolves to the current user's home directory.
	//
	// On Unix, including macOS, it returns the $HOME environment variable.
	// On Windows, it returns %USERPROFILE%.
	// On Plan 9, it returns the $home environment variable.
	//
	// If the location cannot be determined (for example, $HOME is not defined),
	// then it will return an error wrapped in [Err].
	UserHomeDir = file.UserHomeDir
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
func Close[C io.Closer](c C) IOResult[struct{}] {
	return file.Close(c)
}
