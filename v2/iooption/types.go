// Copyright (c) 2025 IBM Corp.
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

package iooption

import (
	"github.com/IBM/fp-go/v2/io"
	I "github.com/IBM/fp-go/v2/io"
	"github.com/IBM/fp-go/v2/lazy"
	"github.com/IBM/fp-go/v2/option"
)

type (
	Option[A any] = option.Option[A]
	IO[A any]     = io.IO[A]
	Lazy[A any]   = lazy.Lazy[A]

	// IOOption represents a synchronous computation that may fail
	// refer to [https://andywhite.xyz/posts/2021-01-27-rte-foundations/#ioeitherlte-agt] for more details
	IOOption[A any] = I.IO[Option[A]]
)
