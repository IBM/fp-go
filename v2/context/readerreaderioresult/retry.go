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
package readerreaderioresult

import (
	"context"

	RIOE "github.com/IBM/fp-go/v2/context/readerioresult"
	F "github.com/IBM/fp-go/v2/function"
	"github.com/IBM/fp-go/v2/reader"
	"github.com/IBM/fp-go/v2/retry"
)

//go:inline
func Retrying[R, A any](
	policy retry.RetryPolicy,
	action Kleisli[R, retry.RetryStatus, A],
	check Predicate[Result[A]],
) ReaderReaderIOResult[R, A] {
	// get an implementation for the types
	return func(r R) ReaderIOResult[context.Context, A] {
		return RIOE.Retrying(policy, F.Pipe1(action, reader.Map[retry.RetryStatus](reader.Read[ReaderIOResult[context.Context, A]](r))), check)
	}
}
