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

package readerioresult

import (
	F "github.com/IBM/fp-go/v2/function"
	RIOR "github.com/IBM/fp-go/v2/readerioresult"
)

// Bracket makes sure that a resource is cleaned up in the event of an error. The release action is called regardless of
// whether the body action returns and error or not.
//
// The release action runs uninterruptibly (see Uninterruptible): it receives a context that keeps the values of the
// caller's context but is never cancelled. This guarantees that the resource is released even if the context was
// cancelled or timed out while the resource was in use.
func Bracket[
	A, B, ANY any](

	acquire ReaderIOResult[A],
	use Kleisli[A, B],
	release func(A, Either[B]) ReaderIOResult[ANY],
) ReaderIOResult[B] {
	return RIOR.Bracket(acquire, F.Flow2(use, WithContext), func(a A, b Either[B]) ReaderIOResult[ANY] {
		return Uninterruptible(release(a, b))
	})
}
