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

package bytes

import (
	"bytes"

	A "github.com/IBM/fp-go/v2/array"
	O "github.com/IBM/fp-go/v2/ord"
)

var (
	// monoid for byte arrays
	Monoid = A.Monoid[byte]()

	// ConcatAll concatenates all bytes
	ConcatAll = A.ArrayConcatAll[byte]

	// Ord implements the default ordering on bytes
	Ord = O.MakeOrd(bytes.Compare, bytes.Equal)
)
