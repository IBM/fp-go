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

package iterresult

import (
	"github.com/IBM/fp-go/v2/iterator/itereither"
)

// Bracket makes sure that a resource is cleaned up in the event of an error.
// The release action is called regardless of whether the body action returns an error or not.
//
// Marble diagram (successful case):
//
//	Acquire: ---R(resource)---|
//	Use:     ---R(result1)---R(result2)---|
//	Release: (called with resource and Right(result2))
//	Output:  ---R(result1)---R(result2)---|
//
// Marble diagram (error case):
//
//	Acquire: ---R(resource)---|
//	Use:     ---R(result1)---L(error)---|
//	Release: (called with resource and Left(error))
//	Output:  ---R(result1)---L(error)---|
//
// The release function is always called to clean up the resource,
// even when an error occurs during use.
func Bracket[A, B, ANY any](
	acquire SeqResult[A],
	use Kleisli[A, B],
	release func(A, Result[B]) SeqResult[ANY],
) SeqResult[B] {
	return itereither.Bracket(acquire, use, release)
}
