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

package exec

import (
	P "github.com/IBM/fp-go/v2/pair"
)

type (
	// CommandOutput represents the output of executing a command. The first field in the [Tuple2] is
	// stdout, the second one is stderr. Use [StdOut] and [StdErr] to access these fields
	CommandOutput = P.Pair[[]byte, []byte]
)

var (
	// StdOut returns the field of a [CommandOutput] representing `stdout`
	StdOut = P.Head[[]byte, []byte]
	// StdErr returns the field of a [CommandOutput] representing `stderr`
	StdErr = P.Tail[[]byte, []byte]
)
