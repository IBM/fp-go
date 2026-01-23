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

package cli

import (
	C "github.com/urfave/cli/v3"
)

const (
	keyFilename = "filename"
	keyCount    = "count"
)

var (
	flagFilename = &C.StringFlag{
		Name:  keyFilename,
		Value: "gen.go",
		Usage: "Name of the generated file",
	}

	flagCount = &C.IntFlag{
		Name:  keyCount,
		Value: 20,
		Usage: "Number of variations to create",
	}
)
