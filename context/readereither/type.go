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

// Package readereither implements a specialization of the Reader monad assuming a golang context as the context of the monad and a standard golang error
package readereither

import (
	"context"

	RE "github.com/IBM/fp-go/readereither"
)

// ReaderEither is a specialization of the Reader monad for the typical golang scenario
type ReaderEither[A any] RE.ReaderEither[context.Context, error, A]
