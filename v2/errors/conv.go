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

package errors

import (
	"errors"

	O "github.com/IBM/fp-go/v2/option"
)

// As tries to extract the error of desired type from the given error
func As[A error]() func(error) O.Option[A] {
	return O.FromValidation(func(err error) (A, bool) {
		var a A
		ok := errors.As(err, &a)
		return a, ok
	})
}
